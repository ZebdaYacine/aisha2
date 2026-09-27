package application

import (
	"context"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/order/domain"
	"testing"
)

type orderAuth struct{ err error }

func (a orderAuth) Authorize(context.Context, auth.Principal, string, string) error { return a.err }

type orderRepo struct{}

func (orderRepo) Checkout(context.Context, string, string, []domain.CartItem, string, string) (domain.Order, error) {
	return domain.Order{ID: "o"}, nil
}
func (orderRepo) ListMine(context.Context, string, int, int) ([]domain.Order, int, error) {
	return nil, 0, nil
}
func (orderRepo) GetMine(context.Context, string, string) (domain.Order, error) {
	return domain.Order{}, nil
}
func (orderRepo) ListSeller(context.Context, string, int, int) ([]domain.SellerItem, int, error) {
	return nil, 0, nil
}
func (orderRepo) Cancel(context.Context, string, string) (domain.Order, error) {
	return domain.Order{}, nil
}
func (orderRepo) RecordReturn(context.Context, string, string, string) (domain.Return, error) {
	return domain.Return{}, nil
}
func TestCheckoutRejectsDuplicateLinesAndMissingIdempotencyKey(t *testing.T) {
	s := NewService(orderRepo{}, orderAuth{})
	p := auth.Principal{UserID: "u"}
	items := []domain.CartItem{{ProductID: "p", Quantity: 1}, {ProductID: "p", Quantity: 1}}
	if _, err := s.Checkout(context.Background(), p, "a", items, "key"); err != domain.ErrValidation {
		t.Fatalf("duplicate err=%v", err)
	}
	if _, err := s.Checkout(context.Background(), p, "a", []domain.CartItem{{ProductID: "p", Quantity: 1}}, ""); err != domain.ErrValidation {
		t.Fatalf("key err=%v", err)
	}
}

func TestCheckoutValidatesAddressAndProductIdentifiersBeforeRepository(t *testing.T) {
	s := NewService(orderRepo{}, orderAuth{})
	p := auth.Principal{UserID: "u"}
	items := []domain.CartItem{{ProductID: "00000000-0000-0000-0000-000000000001", Quantity: 1}}
	if _, err := s.Checkout(context.Background(), p, "not-a-uuid", items, "checkout-1"); err != domain.ErrValidation {
		t.Fatalf("address err=%v", err)
	}
	if _, err := s.Checkout(context.Background(), p, "00000000-0000-0000-0000-000000000002", []domain.CartItem{{ProductID: "not-a-uuid", Quantity: 1}}, "checkout-1"); err != domain.ErrValidation {
		t.Fatalf("product err=%v", err)
	}
}

func TestOrderIdentifiersAreValidatedBeforeRepository(t *testing.T) {
	s := NewService(orderRepo{}, orderAuth{})
	p := auth.Principal{UserID: "u"}
	if _, err := s.GetMine(context.Background(), p, "not-a-uuid"); err != domain.ErrValidation {
		t.Fatalf("get err=%v", err)
	}
	if _, err := s.Cancel(context.Background(), p, "not-a-uuid"); err != domain.ErrValidation {
		t.Fatalf("cancel err=%v", err)
	}
	if _, err := s.RecordReturn(context.Background(), p, "not-a-uuid", "damaged"); err != domain.ErrValidation {
		t.Fatalf("return err=%v", err)
	}
}
