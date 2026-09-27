package auth

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepository struct {
	user                User
	createdPasswordHash string
	createdTokenHash    string
	rotatedCurrentHash  string
	resetTokenHash      string
	resetPasswordHash   string
	sessionActive       bool
}

func (r *fakeRepository) CreateUser(_ context.Context, email, hash, name string) (User, error) {
	r.createdPasswordHash = hash
	r.user = User{ID: "user-1", Email: email, PasswordHash: hash, DisplayName: name, Status: "ACTIVE", Roles: []string{"customer"}}
	return r.user, nil
}
func (r *fakeRepository) UserByEmail(context.Context, string) (User, error) {
	if r.user.ID == "" {
		return User{}, ErrInvalidCredentials
	}
	return r.user, nil
}
func (r *fakeRepository) UserByIdentifier(context.Context, string) (User, error) {
	if r.user.ID == "" {
		return User{}, ErrInvalidCredentials
	}
	return r.user, nil
}
func (r *fakeRepository) UserByID(context.Context, string) (User, error) { return r.user, nil }
func (r *fakeRepository) CreateSession(_ context.Context, userID, hash, family string, expires time.Time) (Session, error) {
	r.createdTokenHash = hash
	r.sessionActive = true
	return Session{ID: "session-1", UserID: userID, FamilyID: family, ExpiresAt: expires}, nil
}
func (r *fakeRepository) RotateSession(_ context.Context, current, next string, expires time.Time) (Session, error) {
	r.rotatedCurrentHash = current
	r.createdTokenHash = next
	return Session{ID: "session-2", UserID: r.user.ID, FamilyID: "family-1", ExpiresAt: expires}, nil
}
func (*fakeRepository) RevokeSession(context.Context, string) error { return nil }
func (r *fakeRepository) SessionActive(context.Context, string, time.Time) (bool, error) {
	return r.sessionActive, nil
}
func (r *fakeRepository) StorePasswordReset(_ context.Context, _ string, hash string, _ time.Time) error {
	r.resetTokenHash = hash
	return nil
}
func (r *fakeRepository) UpdatePassword(_ context.Context, _ string, passwordHash string, _ time.Time) error {
	r.resetPasswordHash = passwordHash
	return nil
}
func (r *fakeRepository) ResetPassword(_ context.Context, tokenHash, passwordHash string, _ time.Time) error {
	r.resetTokenHash = tokenHash
	r.resetPasswordHash = passwordHash
	return nil
}

type fakeNotifier struct{ email, token string }

func (n *fakeNotifier) SendPasswordReset(_ context.Context, _ string, email, token string) error {
	n.email = email
	n.token = token
	return nil
}

func TestRegisterHashesPasswordAndReturnsUsableTokens(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, nil, "test-signing-key")
	user, tokens, err := service.Register(context.Background(), " USER@Example.COM ", "long-password-value", "Amina")
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != "user@example.com" || repo.createdPasswordHash == "long-password-value" {
		t.Fatalf("registration did not normalize and hash: %#v", user)
	}
	if tokens.RefreshToken == "" || repo.createdTokenHash == tokens.RefreshToken {
		t.Fatal("refresh token must be returned and stored only as a hash")
	}
	principal, err := service.Authenticate(context.Background(), tokens.AccessToken)
	if err != nil || principal.UserID != user.ID {
		t.Fatalf("authenticate returned %#v, %v", principal, err)
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, nil, "key")
	_, _, err := service.Register(context.Background(), "user@example.com", "long-password-value", "Amina")
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = service.Login(context.Background(), "user@example.com", "wrong-password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestLoginAcceptsPhoneIdentifier(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, nil, "key")
	user, _, err := service.Register(context.Background(), "user@example.com", "long-password-value", "Amina")
	if err != nil {
		t.Fatal(err)
	}
	repo.user.Phone = "0555123456"
	if _, _, err = service.Login(context.Background(), "0555123456", "long-password-value"); err != nil {
		t.Fatalf("phone identifier should authenticate user %q: %v", user.ID, err)
	}
}

func TestAuthenticateRejectsRevokedBackingSession(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, nil, "test-signing-key")
	_, tokens, err := service.Register(context.Background(), "user@example.com", "long-password-value", "Amina")
	if err != nil {
		t.Fatal(err)
	}
	repo.sessionActive = false
	if _, err := service.Authenticate(context.Background(), tokens.AccessToken); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected revoked session to be rejected, got %v", err)
	}
}

func TestRefreshRotatesHashedToken(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, nil, "key")
	_, tokens, err := service.Register(context.Background(), "user@example.com", "long-password-value", "Amina")
	if err != nil {
		t.Fatal(err)
	}
	next, err := service.Refresh(context.Background(), tokens.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}
	if next.RefreshToken == tokens.RefreshToken || repo.rotatedCurrentHash != hashToken(tokens.RefreshToken) {
		t.Fatal("refresh token was not rotated by hash")
	}
}

func TestForgotPasswordDoesNotRevealUnknownEmail(t *testing.T) {
	service := NewService(&fakeRepository{}, nil, "key")
	if err := service.ForgotPassword(context.Background(), "missing@example.com"); err != nil {
		t.Fatalf("unknown email must receive neutral success: %v", err)
	}
}

func TestForgotAndResetUseHashedTokens(t *testing.T) {
	repo := &fakeRepository{user: User{ID: "user-1", Email: "user@example.com", Status: "ACTIVE"}}
	notifier := &fakeNotifier{}
	service := NewService(repo, notifier, "key")
	if err := service.ForgotPassword(context.Background(), repo.user.Email); err != nil {
		t.Fatal(err)
	}
	if notifier.token == "" || repo.resetTokenHash == notifier.token {
		t.Fatal("reset token must be delivered raw and persisted hashed")
	}
	if err := service.ResetPassword(context.Background(), notifier.token, "replacement-password"); err != nil {
		t.Fatal(err)
	}
	if repo.resetPasswordHash == "replacement-password" {
		t.Fatal("replacement password must be hashed")
	}
}

func TestChangePasswordVerifiesCurrentPasswordAndHashesReplacement(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, nil, "key")
	user, _, err := service.Register(context.Background(), "admin@example.com", "current-password", "Admin")
	if err != nil {
		t.Fatal(err)
	}
	if err = service.ChangePassword(context.Background(), Principal{UserID: user.ID}, "current-password", "replacement-password"); err != nil {
		t.Fatal(err)
	}
	if repo.resetPasswordHash == "replacement-password" || repo.resetPasswordHash == "" {
		t.Fatal("replacement password must be hashed")
	}
	if err = service.ChangePassword(context.Background(), Principal{UserID: user.ID}, "wrong-password", "another-password"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("wrong current password error=%v", err)
	}
}
