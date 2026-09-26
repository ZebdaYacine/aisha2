package application

import (
	"context"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/cart/domain"
)

type cartAuthorizer struct{ err error }

func (a cartAuthorizer) Authorize(context.Context, auth.Principal, string, string) error {
	return a.err
}

type cartRepository struct{ merged []domain.Input }

func (r *cartRepository) List(context.Context, string) ([]domain.Item, error) { return nil, nil }
func (r *cartRepository) Add(context.Context, string, domain.Input) ([]domain.Item, error) {
	return nil, nil
}
func (r *cartRepository) Set(context.Context, string, string, int) ([]domain.Item, error) {
	return nil, nil
}
func (r *cartRepository) Remove(context.Context, string, string) ([]domain.Item, error) {
	return nil, nil
}
func (r *cartRepository) Merge(_ context.Context, _ string, in []domain.Input) ([]domain.Item, error) {
	r.merged = in
	return nil, nil
}

func TestMergeValidatesProductAndQuantity(t *testing.T) {
	r := &cartRepository{}
	s := NewService(r, cartAuthorizer{})
	_, err := s.Merge(context.Background(), auth.Principal{UserID: "u"}, []domain.Input{{ProductID: "bad", Quantity: 1}})
	if err != domain.ErrValidation {
		t.Fatalf("err=%v", err)
	}
	_, err = s.Merge(context.Background(), auth.Principal{UserID: "u"}, []domain.Input{{ProductID: "00000000-0000-0000-0000-000000000001", Quantity: 101}})
	if err != domain.ErrValidation {
		t.Fatalf("err=%v", err)
	}
	if len(r.merged) != 0 {
		t.Fatal("repository should not be called")
	}
}

func TestMergePassesValidatedItems(t *testing.T) {
	r := &cartRepository{}
	s := NewService(r, cartAuthorizer{})
	_, err := s.Merge(context.Background(), auth.Principal{UserID: "u"}, []domain.Input{{ProductID: "00000000-0000-0000-0000-000000000001", Quantity: 2}})
	if err != nil {
		t.Fatal(err)
	}
	if len(r.merged) != 1 || r.merged[0].Quantity != 2 {
		t.Fatalf("merged=%v", r.merged)
	}
}
