package auth

import "github.com/aisha-platform/aisha/apps/api/internal/features/auth/domain"

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
	ErrValidation         = domain.ErrValidation
)
