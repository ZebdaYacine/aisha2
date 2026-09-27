package authorization

import (
	"context"
	"errors"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

func TestAuthorizeUsesCasbinPolicies(t *testing.T) {
	service, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if err := service.AddPolicy("customer", "/api/v1/orders/test", "read"); err != nil {
		t.Fatal(err)
	}
	principal := auth.Principal{UserID: "user-1", Roles: []string{"customer"}}
	if err := service.Authorize(context.Background(), principal, "/api/v1/orders/test", "read"); err != nil {
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

func TestAuthorizeArtisanRetainsCustomerAccountAccess(t *testing.T) {
	service, err := New()
	if err != nil {
		t.Fatal(err)
	}
	principal := auth.Principal{UserID: "artisan-user", Roles: []string{"artisan"}}
	for _, permission := range [][2]string{{"/api/v1/me", "write"}, {"/api/v1/addresses", "write"}, {"/api/v1/artisan-applications/me", "read"}, {"/api/v1/artisan-applications/me/documents", "read"}} {
		if err := service.Authorize(context.Background(), principal, permission[0], permission[1]); err != nil {
			t.Fatalf("artisan permission %s:%s denied: %v", permission[0], permission[1], err)
		}
	}
}

func TestBackOfficeRoleBoundaries(t *testing.T) {
	service, err := New()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name   string
		role   string
		allow  [][2]string
		denied [][2]string
	}{
		{
			name:   "moderator",
			role:   "moderator",
			allow:  [][2]string{{"/api/v1/admin/product-submissions", "write"}, {"/api/v1/admin/users", "read"}, {"/api/v1/admin/media", "read"}},
			denied: [][2]string{{"/api/v1/admin/users", "write"}, {"/api/v1/admin/artisan-applications", "write"}, {"/api/v1/warehouse/inventory", "write"}},
		},
		{
			name:   "warehouse agent",
			role:   "warehouse_agent",
			allow:  [][2]string{{"/api/v1/warehouse/receptions", "write"}, {"/api/v1/warehouse/inventory", "write"}, {"/api/v1/admin/users", "read"}, {"/api/v1/admin/artisan-applications", "read"}, {"/api/v1/admin/artisan-applications/1/media", "read"}},
			denied: [][2]string{{"/api/v1/admin/artisan-applications/1", "write"}, {"/api/v1/admin/media", "read"}, {"/api/v1/admin/product-submissions", "write"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			principal := auth.Principal{UserID: "user-1", Roles: []string{test.role}}
			for _, permission := range test.allow {
				if err := service.Authorize(context.Background(), principal, permission[0], permission[1]); err != nil {
					t.Fatalf("expected %s:%s to be allowed: %v", permission[0], permission[1], err)
				}
			}
			for _, permission := range test.denied {
				if err := service.Authorize(context.Background(), principal, permission[0], permission[1]); !errors.Is(err, ErrForbidden) {
					t.Fatalf("expected %s:%s to be denied, got %v", permission[0], permission[1], err)
				}
			}
		})
	}
}
