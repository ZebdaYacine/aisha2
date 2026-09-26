package application

import (
	"context"
	"strings"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/cart/domain"
	"github.com/google/uuid"
)

const resource = "/api/v1/cart"

type Service struct {
	repository domain.Repository
	authorizer domain.Authorizer
}

type Item = domain.Item
type Input = domain.Input
type Authorizer = domain.Authorizer

var (
	ErrValidation  = domain.ErrValidation
	ErrNotFound    = domain.ErrNotFound
	ErrUnavailable = domain.ErrUnavailable
)

func NewService(r domain.Repository, a domain.Authorizer) *Service {
	return &Service{repository: r, authorizer: a}
}

func (s *Service) List(ctx context.Context, p auth.Principal) ([]domain.Item, error) {
	if err := s.authorizer.Authorize(ctx, p, resource, "read"); err != nil {
		return nil, err
	}
	return s.repository.List(ctx, p.UserID)
}

func (s *Service) Add(ctx context.Context, p auth.Principal, in domain.Input) ([]domain.Item, error) {
	if err := s.authorizer.Authorize(ctx, p, resource, "write"); err != nil {
		return nil, err
	}
	if err := validate(in); err != nil {
		return nil, err
	}
	return s.repository.Add(ctx, p.UserID, in)
}

func (s *Service) Set(ctx context.Context, p auth.Principal, productID string, quantity int) ([]domain.Item, error) {
	if err := s.authorizer.Authorize(ctx, p, resource, "write"); err != nil {
		return nil, err
	}
	if _, err := uuid.Parse(productID); err != nil || quantity < 1 || quantity > 100 {
		return nil, domain.ErrValidation
	}
	return s.repository.Set(ctx, p.UserID, productID, quantity)
}

func (s *Service) Remove(ctx context.Context, p auth.Principal, productID string) ([]domain.Item, error) {
	if err := s.authorizer.Authorize(ctx, p, resource, "write"); err != nil {
		return nil, err
	}
	if _, err := uuid.Parse(productID); err != nil {
		return nil, domain.ErrValidation
	}
	return s.repository.Remove(ctx, p.UserID, productID)
}

func (s *Service) Merge(ctx context.Context, p auth.Principal, inputs []domain.Input) ([]domain.Item, error) {
	if err := s.authorizer.Authorize(ctx, p, resource, "write"); err != nil {
		return nil, err
	}
	if len(inputs) > 100 {
		return nil, domain.ErrValidation
	}
	for _, in := range inputs {
		if err := validate(in); err != nil {
			return nil, err
		}
	}
	return s.repository.Merge(ctx, p.UserID, inputs)
}

func validate(in domain.Input) error {
	if _, err := uuid.Parse(strings.TrimSpace(in.ProductID)); err != nil || in.Quantity < 1 || in.Quantity > 100 {
		return domain.ErrValidation
	}
	return nil
}
