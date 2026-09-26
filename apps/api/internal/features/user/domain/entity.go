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

// AccountState is the small cross-capability read model used to compose the
// authenticated account summary. Artisan membership is intentionally read
// from the existing artisan profile until the membership model is completed
// by PRJ-EPIC-007.
type AccountState struct {
	UserStatus    string
	ArtisanStatus string
}

type AccountSummary struct {
	UserID          string
	UserStatus      string
	ArtisanStatus   string
	CustomerEnabled bool
	ArtisanEnabled  bool
	Capabilities    []string
}

type Repository interface {
	Profile(context.Context, string) (Profile, error)
	UpdateProfile(context.Context, string, string, string) (Profile, error)
	Addresses(context.Context, string) ([]Address, error)
	CreateAddress(context.Context, string, AddressInput) (Address, error)
	UpdateAddress(context.Context, string, string, AddressInput) (Address, error)
	DeleteAddress(context.Context, string, string) error
}

type AccountRepository interface {
	AccountState(context.Context, string) (AccountState, error)
}
type Authorizer interface {
	Authorize(context.Context, auth.Principal, string, string) error
}
