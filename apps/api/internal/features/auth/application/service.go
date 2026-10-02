package application

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const accessTTL = 15 * time.Minute
const refreshTTL = 30 * 24 * time.Hour
const resetTTL = time.Hour

type User = domain.User
type Session = domain.Session
type Repository = domain.Repository
type ResetNotifier = domain.ResetNotifier
type Principal = domain.Principal
type Tokens = domain.Tokens

var (
	ErrInvalidCredentials = domain.ErrInvalidCredentials
	ErrEmailExists        = domain.ErrEmailExists
	ErrInvalidToken       = domain.ErrInvalidToken
	ErrUserInactive       = domain.ErrUserInactive
	ErrEmailUnverified    = domain.ErrEmailUnverified
	ErrValidation         = domain.ErrValidation
)

type Service struct {
	repo       Repository
	notifier   ResetNotifier
	secret     []byte
	now        func() time.Time
	webBaseURL string
}

func NewService(repo Repository, notifier ResetNotifier, secret string, webBaseURLs ...string) *Service {
	webBaseURL := "http://localhost:3033"
	if len(webBaseURLs) > 0 && strings.TrimSpace(webBaseURLs[0]) != "" {
		webBaseURL = strings.TrimRight(strings.TrimSpace(webBaseURLs[0]), "/")
	}
	return &Service{repo: repo, notifier: notifier, secret: []byte(secret), now: func() time.Time { return time.Now().UTC() }, webBaseURL: webBaseURL}
}

func (s *Service) Register(ctx context.Context, email, password, name string) (User, Tokens, error) {
	email = normalizeEmail(email)
	if !validEmail(email) || len(password) < 12 || len(strings.TrimSpace(name)) < 2 {
		return User{}, Tokens{}, fmt.Errorf("register validation: %w", ErrValidation)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, Tokens{}, fmt.Errorf("hash password: %w", err)
	}
	verificationToken, verificationHash, err := randomToken()
	if err != nil {
		return User{}, Tokens{}, fmt.Errorf("create email verification token: %w", err)
	}
	activationURL := s.webBaseURL + "/en/activate?token=" + url.QueryEscape(verificationToken)
	user, err := s.repo.CreateUser(ctx, email, string(hash), strings.TrimSpace(name), verificationHash, activationURL)
	if err != nil {
		return User{}, Tokens{}, fmt.Errorf("create user: %w", err)
	}
	tokens, err := s.newSession(ctx, user.ID, uuid.NewString())
	return user, tokens, err
}

func (s *Service) ActivateEmail(ctx context.Context, token string) (User, Tokens, error) {
	if strings.TrimSpace(token) == "" {
		return User{}, Tokens{}, ErrInvalidToken
	}
	userID, err := s.repo.ConsumeEmailVerification(ctx, hashToken(token), s.now())
	if err != nil {
		return User{}, Tokens{}, fmt.Errorf("consume email verification: %w", err)
	}
	user, err := s.repo.UserByID(ctx, userID)
	if err != nil {
		return User{}, Tokens{}, ErrInvalidToken
	}
	if user.Status != "ACTIVE" {
		return User{}, Tokens{}, ErrUserInactive
	}
	tokens, err := s.newSession(ctx, user.ID, uuid.NewString())
	if err != nil {
		return User{}, Tokens{}, err
	}
	return user, tokens, nil
}
func (s *Service) Login(ctx context.Context, identifier, password string) (User, Tokens, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return User{}, Tokens{}, ErrInvalidCredentials
	}
	user, err := s.repo.UserByIdentifier(ctx, identifier)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return User{}, Tokens{}, ErrInvalidCredentials
	}
	if user.Status != "ACTIVE" {
		return User{}, Tokens{}, ErrUserInactive
	}
	if user.Email != "" && user.EmailVerifiedAt == nil {
		return User{}, Tokens{}, ErrEmailUnverified
	}
	tokens, err := s.newSession(ctx, user.ID, uuid.NewString())
	return user, tokens, err
}
func (s *Service) newSession(ctx context.Context, userID, familyID string) (Tokens, error) {
	refresh, hash, err := randomToken()
	if err != nil {
		return Tokens{}, err
	}
	now := s.now()
	session, err := s.repo.CreateSession(ctx, userID, hash, familyID, now.Add(refreshTTL))
	if err != nil {
		return Tokens{}, fmt.Errorf("create session: %w", err)
	}
	expires := now.Add(accessTTL)
	access, err := signAccessToken(s.secret, userID, session.ID, expires)
	return Tokens{AccessToken: access, RefreshToken: refresh, ExpiresAt: expires}, err
}
func (s *Service) Refresh(ctx context.Context, refresh string) (Tokens, error) {
	next, nextHash, err := randomToken()
	if err != nil {
		return Tokens{}, err
	}
	now := s.now()
	session, err := s.repo.RotateSession(ctx, hashToken(refresh), nextHash, now.Add(refreshTTL))
	if err != nil {
		return Tokens{}, fmt.Errorf("rotate session: %w", err)
	}
	expires := now.Add(accessTTL)
	access, err := signAccessToken(s.secret, session.UserID, session.ID, expires)
	return Tokens{AccessToken: access, RefreshToken: next, ExpiresAt: expires}, err
}
func (s *Service) Logout(ctx context.Context, refresh string) error {
	if err := s.repo.RevokeSession(ctx, hashToken(refresh)); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	return nil
}
func (s *Service) Authenticate(ctx context.Context, access string) (Principal, error) {
	claims, err := verifyAccessToken(s.secret, access, s.now())
	if err != nil {
		return Principal{}, ErrInvalidToken
	}
	active, err := s.repo.SessionActive(ctx, claims.SessionID, s.now())
	if err != nil {
		return Principal{}, fmt.Errorf("check session: %w", err)
	}
	if !active {
		return Principal{}, ErrInvalidToken
	}
	user, err := s.repo.UserByID(ctx, claims.UserID)
	if err != nil {
		return Principal{}, ErrInvalidToken
	}
	if user.Status != "ACTIVE" {
		return Principal{}, ErrUserInactive
	}
	if user.Email != "" && user.EmailVerifiedAt == nil {
		return Principal{}, ErrEmailUnverified
	}
	return Principal{UserID: user.ID, Roles: user.Roles}, nil
}
func (s *Service) Me(ctx context.Context, principal Principal) (User, error) {
	return s.repo.UserByID(ctx, principal.UserID)
}
func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.UserByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return nil
	}
	token, hash, err := randomToken()
	if err != nil {
		return err
	}
	if err := s.repo.StorePasswordReset(ctx, user.ID, hash, s.now().Add(resetTTL)); err != nil {
		return fmt.Errorf("store reset token: %w", err)
	}
	if s.notifier != nil {
		if err := s.notifier.SendPasswordReset(ctx, user.ID, user.Email, token); err != nil {
			return fmt.Errorf("send reset token: %w", err)
		}
	}
	return nil
}
func (s *Service) ChangePassword(ctx context.Context, principal Principal, currentPassword, nextPassword string) error {
	if len(nextPassword) < 12 || len(nextPassword) > 128 || strings.TrimSpace(currentPassword) == "" {
		return ErrValidation
	}
	user, err := s.repo.UserByID(ctx, principal.UserID)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)) != nil {
		return ErrInvalidCredentials
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(nextPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := s.repo.UpdatePassword(ctx, principal.UserID, string(hash), s.now()); err != nil {
		return fmt.Errorf("change password: %w", err)
	}
	return nil
}
func (s *Service) ResetPassword(ctx context.Context, token, password string) error {
	if len(password) < 12 {
		return ErrInvalidCredentials
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := s.repo.ResetPassword(ctx, hashToken(token), string(hash), s.now()); err != nil {
		return fmt.Errorf("reset password: %w", err)
	}
	return nil
}
func normalizeEmail(value string) string { return strings.ToLower(strings.TrimSpace(value)) }
func validEmail(value string) bool {
	parts := strings.Split(value, "@")
	return len(parts) == 2 && parts[0] != "" && strings.Contains(parts[1], ".")
}
func Is(err, target error) bool { return errors.Is(err, target) }
