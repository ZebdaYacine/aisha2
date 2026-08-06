package domain

import (
	"context"
	"errors"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

var (
	ErrNotFound   = errors.New("customer resource not found")
	ErrValidation = errors.New("customer validation failed")
)

type Profile struct{ ID, Email, DisplayName, Phone string }
type Address struct {
	ID, FullName, Phone, Line1, Line2, City, PostalCode, Country string
	Default                                                      bool
}
type AddressInput struct {
	FullName, Phone, Line1, Line2, City, PostalCode, Country string
	Default                                                  bool
}
type Repository interface {
	Profile(context.Context, string) (Profile, error)
	UpdateProfile(context.Context, string, string, string) (Profile, error)
	Addresses(context.Context, string) ([]Address, error)
	CreateAddress(context.Context, string, AddressInput) (Address, error)
	UpdateAddress(context.Context, string, string, AddressInput) (Address, error)
	DeleteAddress(context.Context, string, string) error
}
type Authorizer interface {
	Authorize(context.Context, auth.Principal, string, string) error
}
