package application

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/order/domain"
	"github.com/google/uuid"
	"strings"
)

type Service struct {
	repository domain.Repository
	authorizer domain.Authorizer
}
type CartItem = domain.CartItem
type Order = domain.Order
type SellerItem = domain.SellerItem
type Return = domain.Return
type Authorizer = domain.Authorizer

var (
	ErrValidation          = domain.ErrValidation
	ErrNotFound            = domain.ErrNotFound
	ErrOutOfStock          = domain.ErrOutOfStock
	ErrPriceChanged        = domain.ErrPriceChanged
	ErrInvalidTransition   = domain.ErrInvalidTransition
	ErrIdempotencyConflict = domain.ErrIdempotencyConflict
)

func NewService(r domain.Repository, a domain.Authorizer) *Service {
	return &Service{repository: r, authorizer: a}
}
func (s *Service) Checkout(ctx context.Context, p auth.Principal, address string, items []domain.CartItem, key string) (domain.Order, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/checkout", "write"); err != nil {
		return domain.Order{}, err
	}
	key = strings.TrimSpace(key)
	address = strings.TrimSpace(address)
	if key == "" || len(key) > 160 || uuid.Validate(address) != nil || len(items) == 0 || len(items) > 50 {
		return domain.Order{}, domain.ErrValidation
	}
	seen := map[string]bool{}
	for _, i := range items {
		i.ProductID = strings.TrimSpace(i.ProductID)
		if uuid.Validate(i.ProductID) != nil || i.Quantity < 1 || i.Quantity > 100 || seen[i.ProductID] {
			return domain.Order{}, domain.ErrValidation
		}
		seen[i.ProductID] = true
		if i.ExpectedCurrency != "" && len(i.ExpectedCurrency) != 3 {
			return domain.Order{}, domain.ErrValidation
		}
	}
	return s.repository.Checkout(ctx, p.UserID, address, items, key, hash(items, address, key))
}
func (s *Service) ListMine(ctx context.Context, p auth.Principal, page, size int) ([]domain.Order, int, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/orders", "read"); err != nil {
		return nil, 0, err
	}
	page, size = normalizePage(page, size)
	return s.repository.ListMine(ctx, p.UserID, size, (page-1)*size)
}
func (s *Service) GetMine(ctx context.Context, p auth.Principal, id string) (domain.Order, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/orders", "read"); err != nil {
		return domain.Order{}, err
	}
	if uuid.Validate(id) != nil {
		return domain.Order{}, domain.ErrValidation
	}
	return s.repository.GetMine(ctx, p.UserID, id)
}
func (s *Service) ListSeller(ctx context.Context, p auth.Principal, page, size int) ([]domain.SellerItem, int, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/orders/seller", "read"); err != nil {
		return nil, 0, err
	}
	page, size = normalizePage(page, size)
	return s.repository.ListSeller(ctx, p.UserID, size, (page-1)*size)
}
func (s *Service) Cancel(ctx context.Context, p auth.Principal, id string) (domain.Order, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/orders", "write"); err != nil {
		return domain.Order{}, err
	}
	if uuid.Validate(id) != nil {
		return domain.Order{}, domain.ErrValidation
	}
	return s.repository.Cancel(ctx, p.UserID, id)
}
func (s *Service) RecordReturn(ctx context.Context, p auth.Principal, id, reason string) (domain.Return, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/orders/*/returns", "write"); err != nil {
		return domain.Return{}, err
	}
	if strings.TrimSpace(reason) == "" {
		return domain.Return{}, domain.ErrValidation
	}
	if uuid.Validate(id) != nil {
		return domain.Return{}, domain.ErrValidation
	}
	return s.repository.RecordReturn(ctx, p.UserID, id, strings.TrimSpace(reason))
}
func normalizePage(p, s int) (int, int) {
	if p < 1 {
		p = 1
	}
	if s < 1 {
		s = 20
	}
	if s > 100 {
		s = 100
	}
	return p, s
}
func hash(v any, address, key string) string {
	b, _ := json.Marshal(struct {
		V            any
		Address, Key string
	}{v, address, key})
	return fmt.Sprintf("%x", sha256.Sum256(b))
}
