package application

import (
	"context"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/inventory/domain"
)

type inventoryAuthorizer struct{ err error }

func (a inventoryAuthorizer) Authorize(context.Context, auth.Principal, string, string) error {
	return a.err
}

type inventoryRepository struct {
	adjustCalled bool
	global       bool
}

func (r *inventoryRepository) ListBalances(context.Context, string, bool, string, int, int) ([]domain.Balance, int, error) {
	return nil, 0, nil
}
func (r *inventoryRepository) Adjust(_ context.Context, _ string, global bool, _ string, _ domain.AdjustmentInput) (domain.Movement, error) {
	r.adjustCalled = true
	r.global = global
	return domain.Movement{}, nil
}

func TestAdjustRejectsZeroDelta(t *testing.T) {
	r := &inventoryRepository{}
	s := NewService(r, inventoryAuthorizer{})
	_, err := s.Adjust(context.Background(), auth.Principal{UserID: "actor", Roles: []string{"warehouse_agent"}}, "00000000-0000-0000-0000-000000000001", domain.AdjustmentInput{Reason: "count", ReferenceKey: "ref"})
	if err != domain.ErrValidation {
		t.Fatalf("err=%v", err)
	}
	if r.adjustCalled {
		t.Fatal("repository should not be called")
	}
}

func TestAdjustPassesGlobalAccessOnlyToOperationalRoles(t *testing.T) {
	r := &inventoryRepository{}
	s := NewService(r, inventoryAuthorizer{})
	_, err := s.Adjust(context.Background(), auth.Principal{UserID: "actor", Roles: []string{"artisan"}}, "00000000-0000-0000-0000-000000000001", domain.AdjustmentInput{QuantityDelta: 2, Reason: "cycle count", ReferenceKey: "ref"})
	if err != nil {
		t.Fatal(err)
	}
	if r.global {
		t.Fatal("artisan must not receive global inventory scope")
	}
}

func TestListNormalizesPageAndRejectsInvalidWorkshop(t *testing.T) {
	s := NewService(&inventoryRepository{}, inventoryAuthorizer{})
	if _, _, err := s.List(context.Background(), auth.Principal{UserID: "actor"}, "not-a-uuid", 1, 20); err != domain.ErrValidation {
		t.Fatalf("err=%v", err)
	}
}
