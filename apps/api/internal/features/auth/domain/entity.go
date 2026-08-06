package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserInactive       = errors.New("user is inactive")
	ErrValidation         = errors.New("validation error")
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	DisplayName  string
	Status       string
	Roles        []string
	CreatedAt    time.Time
}

type Session struct {
	ID        string
	UserID    string
	FamilyID  string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

type Repository interface {
	CreateUser(context.Context, string, string, string) (User, error)
	UserByEmail(context.Context, string) (User, error)
	UserByID(context.Context, string) (User, error)
	CreateSession(context.Context, string, string, string, time.Time) (Session, error)
	RotateSession(context.Context, string, string, time.Time) (Session, error)
	RevokeSession(context.Context, string) error
	SessionActive(context.Context, string, time.Time) (bool, error)
	StorePasswordReset(context.Context, string, string, time.Time) error
	ResetPassword(context.Context, string, string, time.Time) error
}

type ResetNotifier interface {
	SendPasswordReset(context.Context, string, string, string) error
}

type Principal struct {
	UserID string
	Roles  []string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
