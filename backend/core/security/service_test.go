package authorization

import (
	"context"
	"errors"
	"testing"

	"github.com/aisha-platform/aisha/backend/features/auth"
)

func TestAuthorizeUsesCasbinPolicies(t *testing.T) {
	service, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if err := service.AddPolicy("customer", "/api/v1/orders", "read"); err != nil {
		t.Fatal(err)
	}
	principal := auth.Principal{UserID: "user-1", Roles: []string{"customer"}}
	if err := service.Authorize(context.Background(), principal, "/api/v1/orders", "read"); err != nil {
		t.Fatalf("expected allow, got %v", err)
	}
	if err := service.Authorize(context.Background(), principal, "/api/v1/admin/orders", "read"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected deny, got %v", err)
	}
}

func TestAuthorizeChecksEveryAssignedRole(t *testing.T) {
	service, err := New()
	if err != nil {
		t.Fatal(err)
	}
	principal := auth.Principal{Roles: []string{"customer", "artisan"}}
	if err := service.Authorize(context.Background(), principal, "/api/v1/artisan/profile", "write"); err != nil {
		t.Fatal(err)
	}
}

func TestDefaultPoliciesAllowAuthenticatedRolesToReadOwnPrincipal(t *testing.T) {
	service, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for _, role := range []string{"customer", "artisan", "moderator", "warehouse_agent", "administrator"} {
		principal := auth.Principal{UserID: "user-1", Roles: []string{role}}
		if err := service.Authorize(context.Background(), principal, "/api/v1/me", "read"); err != nil {
			t.Fatalf("expected %s to read /me: %v", role, err)
		}
	}
	visitor := auth.Principal{Roles: []string{"visitor"}}
	if err := service.Authorize(context.Background(), visitor, "/api/v1/me", "read"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected visitor principal to be denied, got %v", err)
	}
}
