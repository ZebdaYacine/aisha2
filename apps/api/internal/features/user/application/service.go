package application

import (
	"context"
	"strings"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/user/domain"
)

type Service struct {
	repository domain.Repository
	authorizer domain.Authorizer
}

func NewService(repository domain.Repository, authorizer domain.Authorizer) *Service {
	return &Service{repository: repository, authorizer: authorizer}
}
func (s *Service) Profile(ctx context.Context, p auth.Principal) (domain.Profile, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/me", "read"); err != nil {
		return domain.Profile{}, err
	}
	return s.repository.Profile(ctx, p.UserID)
}
func (s *Service) UpdateProfile(ctx context.Context, p auth.Principal, name, phone string) (domain.Profile, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/me", "write"); err != nil {
		return domain.Profile{}, err
	}
	name, phone = strings.TrimSpace(name), strings.TrimSpace(phone)
	if len(name) < 2 {
		return domain.Profile{}, domain.ErrValidation
	}
	return s.repository.UpdateProfile(ctx, p.UserID, name, phone)
}
func (s *Service) Addresses(ctx context.Context, p auth.Principal) ([]domain.Address, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/addresses", "read"); err != nil {
		return nil, err
	}
	return s.repository.Addresses(ctx, p.UserID)
}
func (s *Service) CreateAddress(ctx context.Context, p auth.Principal, input domain.AddressInput) (domain.Address, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/addresses", "write"); err != nil {
		return domain.Address{}, err
	}
	if !validAddress(input) {
		return domain.Address{}, domain.ErrValidation
	}
	return s.repository.CreateAddress(ctx, p.UserID, input)
}
func (s *Service) UpdateAddress(ctx context.Context, p auth.Principal, id string, input domain.AddressInput) (domain.Address, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/addresses", "write"); err != nil {
		return domain.Address{}, err
	}
	if !validAddress(input) {
		return domain.Address{}, domain.ErrValidation
	}
	return s.repository.UpdateAddress(ctx, p.UserID, id, input)
}
func (s *Service) DeleteAddress(ctx context.Context, p auth.Principal, id string) error {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/addresses", "write"); err != nil {
		return err
	}
	return s.repository.DeleteAddress(ctx, p.UserID, id)
}
func validAddress(i domain.AddressInput) bool {
	return strings.TrimSpace(i.FullName) != "" && strings.TrimSpace(i.Line1) != "" && strings.TrimSpace(i.City) != "" && strings.TrimSpace(i.PostalCode) != "" && strings.TrimSpace(i.Country) != ""
}
