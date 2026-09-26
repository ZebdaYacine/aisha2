package application

import (
	"context"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/wishlist/domain"
	"testing"
)

type wishlistAuthorizer struct{ err error }

func (a wishlistAuthorizer) Authorize(context.Context, auth.Principal, string, string) error {
	return a.err
}

type wishlistRepository struct{ added string }

func (wishlistRepository) List(context.Context, string) ([]domain.Item, error) { return nil, nil }
func (r *wishlistRepository) Add(_ context.Context, _ string, id string) ([]domain.Item, error) {
	r.added = id
	return nil, nil
}
func (wishlistRepository) Remove(context.Context, string, string) ([]domain.Item, error) {
	return nil, nil
}
func TestAddRejectsInvalidProductID(t *testing.T) {
	_, err := NewService(&wishlistRepository{}, wishlistAuthorizer{}).Add(context.Background(), auth.Principal{UserID: "u"}, "bad")
	if err != domain.ErrValidation {
		t.Fatalf("err=%v", err)
	}
}
func TestAddUsesOwnerScopedRepository(t *testing.T) {
	r := &wishlistRepository{}
	id := "00000000-0000-0000-0000-000000000001"
	if _, err := NewService(r, wishlistAuthorizer{}).Add(context.Background(), auth.Principal{UserID: "owner"}, id); err != nil {
		t.Fatal(err)
	}
	if r.added != id {
		t.Fatalf("added=%s", r.added)
	}
}
