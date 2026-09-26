package application

import (
	"context"
	"strings"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/inventory/domain"
	"github.com/google/uuid"
)

const resource = "/api/v1/warehouse/inventory"

type Service struct {
	repository domain.Repository
	authorizer domain.Authorizer
}

type Balance = domain.Balance
type Movement = domain.Movement
type AdjustmentInput = domain.AdjustmentInput
type Authorizer = domain.Authorizer

var (
	ErrValidation        = domain.ErrValidation
	ErrNotFound          = domain.ErrNotFound
	ErrDuplicate         = domain.ErrDuplicate
	ErrInsufficientStock = domain.ErrInsufficientStock
)

func NewService(repository domain.Repository, authorizer domain.Authorizer) *Service {
	return &Service{repository: repository, authorizer: authorizer}
}

func (s *Service) List(ctx context.Context, principal auth.Principal, workshopID string, page, size int) ([]domain.Balance, int, error) {
	if err := s.authorizer.Authorize(ctx, principal, resource, "read"); err != nil {
		return nil, 0, err
	}
	if workshopID != "" {
		if _, err := uuid.Parse(workshopID); err != nil {
			return nil, 0, domain.ErrValidation
		}
	}
	page, size = normalizePage(page, size)
	return s.repository.ListBalances(ctx, principal.UserID, hasGlobalAccess(principal), workshopID, size, (page-1)*size)
}

func (s *Service) Adjust(ctx context.Context, principal auth.Principal, productID string, input domain.AdjustmentInput) (domain.Movement, error) {
	if err := s.authorizer.Authorize(ctx, principal, resource, "write"); err != nil {
		return domain.Movement{}, err
	}
	if _, err := uuid.Parse(productID); err != nil {
		return domain.Movement{}, domain.ErrValidation
	}
	input.Reason = strings.TrimSpace(input.Reason)
	input.ReferenceKey = strings.TrimSpace(input.ReferenceKey)
	if input.QuantityDelta == 0 || input.Reason == "" || len(input.Reason) > 2000 || input.ReferenceKey == "" || len(input.ReferenceKey) > 160 {
		return domain.Movement{}, domain.ErrValidation
	}
	return s.repository.Adjust(ctx, principal.UserID, hasGlobalAccess(principal), productID, input)
}

func hasGlobalAccess(principal auth.Principal) bool {
	for _, role := range principal.Roles {
		if role == "administrator" || role == "warehouse_agent" {
			return true
		}
	}
	return false
}

func normalizePage(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return page, size
}
