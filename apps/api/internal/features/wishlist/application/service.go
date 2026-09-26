package application

import (
	"context"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/wishlist/domain"
	"github.com/google/uuid"
)

const resource = "/api/v1/wishlist"

type Service struct {
	repository domain.Repository
	authorizer domain.Authorizer
}
type Item = domain.Item

var ErrValidation = domain.ErrValidation

func NewService(r domain.Repository, a domain.Authorizer) *Service {
	return &Service{repository: r, authorizer: a}
}
func (s *Service) List(ctx context.Context, p auth.Principal) ([]domain.Item, error) {
	if err := s.authorizer.Authorize(ctx, p, resource, "read"); err != nil {
		return nil, err
	}
	return s.repository.List(ctx, p.UserID)
}
func (s *Service) Add(ctx context.Context, p auth.Principal, id string) ([]domain.Item, error) {
	if err := s.authorizer.Authorize(ctx, p, resource, "write"); err != nil {
		return nil, err
	}
	if _, err := uuid.Parse(id); err != nil {
		return nil, domain.ErrValidation
	}
	return s.repository.Add(ctx, p.UserID, id)
}
func (s *Service) Remove(ctx context.Context, p auth.Principal, id string) ([]domain.Item, error) {
	if err := s.authorizer.Authorize(ctx, p, resource, "write"); err != nil {
		return nil, err
	}
	if _, err := uuid.Parse(id); err != nil {
		return nil, domain.ErrValidation
	}
	return s.repository.Remove(ctx, p.UserID, id)
}
