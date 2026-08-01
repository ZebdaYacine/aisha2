package customer

import (
	"context"
	"errors"
	"strings"

	"github.com/aisha-platform/aisha/backend/internal/auth"
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
type Service struct {
	repository Repository
	authorizer Authorizer
}

func NewService(repository Repository, authorizer Authorizer) *Service {
	return &Service{repository: repository, authorizer: authorizer}
}
func (s *Service) Profile(ctx context.Context, p auth.Principal) (Profile, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/me", "read"); err != nil {
		return Profile{}, err
	}
	return s.repository.Profile(ctx, p.UserID)
}
func (s *Service) UpdateProfile(ctx context.Context, p auth.Principal, name, phone string) (Profile, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/me", "write"); err != nil {
		return Profile{}, err
	}
	name = strings.TrimSpace(name)
	phone = strings.TrimSpace(phone)
	if len(name) < 2 {
		return Profile{}, ErrValidation
	}
	return s.repository.UpdateProfile(ctx, p.UserID, name, phone)
}
func (s *Service) Addresses(ctx context.Context, p auth.Principal) ([]Address, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/addresses", "read"); err != nil {
		return nil, err
	}
	return s.repository.Addresses(ctx, p.UserID)
}
func (s *Service) CreateAddress(ctx context.Context, p auth.Principal, input AddressInput) (Address, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/addresses", "write"); err != nil {
		return Address{}, err
	}
	if !validAddress(input) {
		return Address{}, ErrValidation
	}
	return s.repository.CreateAddress(ctx, p.UserID, input)
}
func (s *Service) UpdateAddress(ctx context.Context, p auth.Principal, id string, input AddressInput) (Address, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/addresses", "write"); err != nil {
		return Address{}, err
	}
	if !validAddress(input) {
		return Address{}, ErrValidation
	}
	return s.repository.UpdateAddress(ctx, p.UserID, id, input)
}
func (s *Service) DeleteAddress(ctx context.Context, p auth.Principal, id string) error {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/addresses", "write"); err != nil {
		return err
	}
	return s.repository.DeleteAddress(ctx, p.UserID, id)
}
func validAddress(i AddressInput) bool {
	return strings.TrimSpace(i.FullName) != "" && strings.TrimSpace(i.Line1) != "" && strings.TrimSpace(i.City) != "" && strings.TrimSpace(i.PostalCode) != "" && strings.TrimSpace(i.Country) != ""
}
