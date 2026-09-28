package authorization

import (
	"context"
	"fmt"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

type Service struct{ enforcer *casbin.Enforcer }

var defaultPolicies = [][]string{
	{"customer", "/api/v1/me", "read"},
	{"artisan", "/api/v1/me", "read"},
	{"moderator", "/api/v1/me", "read"},
	{"warehouse_agent", "/api/v1/me", "read"},
	{"administrator", "/api/v1/me", "read"},
	{"customer", "/api/v1/me", "write"},
	{"customer", "/api/v1/addresses", "read"},
	{"customer", "/api/v1/addresses", "write"},
	// Artisan membership is additive: artisan principals retain the same
	// customer profile and address permissions as every registered user.
	{"artisan", "/api/v1/me", "write"},
	{"artisan", "/api/v1/addresses", "read"},
	{"artisan", "/api/v1/addresses", "write"},
	{"moderator", "/api/v1/me", "write"},
	{"moderator", "/api/v1/addresses", "read"},
	{"moderator", "/api/v1/addresses", "write"},
	{"warehouse_agent", "/api/v1/me", "write"},
	{"warehouse_agent", "/api/v1/addresses", "read"},
	{"warehouse_agent", "/api/v1/addresses", "write"},
	{"administrator", "/api/v1/me", "write"},
	{"administrator", "/api/v1/addresses", "read"},
	{"administrator", "/api/v1/addresses", "write"},
	{"customer", "/api/v1/artisan-applications", "write"},
	{"customer", "/api/v1/artisan-applications/draft", "write"},
	{"customer", "/api/v1/artisan-applications/me", "read"},
	{"customer", "/api/v1/artisan-applications/me/submit", "write"},
	{"customer", "/api/v1/artisan-applications/me/documents", "read"},
	{"customer", "/api/v1/artisan-applications/me/documents", "write"},
	// An approved artisan must be able to read the application that owns the
	// existing workshop instead of being sent back to the application form.
	{"artisan", "/api/v1/artisan-applications/me", "read"},
	{"artisan", "/api/v1/artisan-applications/me/documents", "read"},
	{"customer", "/api/v1/artisan/profile/media", "read"},
	{"customer", "/api/v1/artisan/profile/media", "write"},
	{"artisan", "/api/v1/artisan/profile", "write"},
	{"artisan", "/api/v1/artisan/profile/media", "read"},
	{"artisan", "/api/v1/artisan/profile/media", "write"},
	{"artisan", "/api/v1/artisan/products", "read"},
	{"artisan", "/api/v1/artisan/workshops", "read"},
	{"artisan", "/api/v1/artisan/workshops", "write"},
	{"customer", "/api/v1/artisan/membership", "write"},
	{"artisan", "/api/v1/artisan/membership", "write"},
	{"artisan", "/api/v1/artisan/verification", "read"},
	{"administrator", "/api/v1/admin/artisan-memberships", "write"},
	{"administrator", "/api/v1/admin/artisan-verifications", "read"},
	{"administrator", "/api/v1/admin/artisan-verifications", "write"},
	{"artisan", "/api/v1/artisan/products", "write"},
	{"administrator", "/api/v1/admin/artisan-applications", "read"},
	{"administrator", "/api/v1/admin/artisan-applications/*", "write"},
	{"administrator", "/api/v1/admin/artisan-applications/*/documents", "read"},
	{"administrator", "/api/v1/admin/artisan-applications/*/media", "read"},
	// Moderators can review products and inspect user media, but cannot manage users.
	{"moderator", "/api/v1/admin/users", "read"},
	{"moderator", "/api/v1/admin/media", "read"},
	// Warehouse agents need read-only identity/application context while their
	// write access remains limited to warehouse reception and inventory.
	{"warehouse_agent", "/api/v1/admin/users", "read"},
	{"warehouse_agent", "/api/v1/admin/artisan-applications", "read"},
	{"warehouse_agent", "/api/v1/admin/artisan-applications/*/documents", "read"},
	{"warehouse_agent", "/api/v1/admin/artisan-applications/*/media", "read"},
	{"administrator", "/api/v1/admin/users", "read"},
	{"administrator", "/api/v1/admin/users", "write"},
	{"administrator", "/api/v1/admin/workshops", "write"},
	{"administrator", "/api/v1/admin/audit-events", "read"},
	{"administrator", "/api/v1/admin/media", "read"},
	{"administrator", "/api/v1/admin/media", "write"},
	{"administrator", "/api/v1/admin/categories", "read"},
	{"administrator", "/api/v1/admin/categories", "write"},
	{"administrator", "/api/v1/admin/orders", "read"},
	{"moderator", "/api/v1/admin/product-submissions", "read"},
	{"moderator", "/api/v1/admin/product-submissions", "write"},
	{"administrator", "/api/v1/admin/product-submissions", "read"},
	{"administrator", "/api/v1/admin/product-submissions", "write"},
	{"customer", "/api/v1/checkout", "write"},
	{"artisan", "/api/v1/checkout", "write"},
	{"customer", "/api/v1/cart", "read"},
	{"customer", "/api/v1/cart", "write"},
	{"customer", "/api/v1/wishlist", "read"},
	{"customer", "/api/v1/wishlist", "write"},
	{"artisan", "/api/v1/cart", "read"},
	{"artisan", "/api/v1/cart", "write"},
	{"artisan", "/api/v1/wishlist", "read"},
	{"artisan", "/api/v1/wishlist", "write"},
	{"artisan", "/api/v1/orders", "read"},
	{"artisan", "/api/v1/orders", "write"},
	{"customer", "/api/v1/orders", "read"},
	{"customer", "/api/v1/orders", "write"},
	{"artisan", "/api/v1/orders/seller", "read"},
	{"customer", "/api/v1/notifications", "read"},
	{"customer", "/api/v1/notifications", "write"},
	{"artisan", "/api/v1/notifications", "read"},
	{"artisan", "/api/v1/notifications", "write"},
	{"moderator", "/api/v1/notifications", "read"},
	{"moderator", "/api/v1/notifications", "write"},
	{"warehouse_agent", "/api/v1/notifications", "read"},
	{"warehouse_agent", "/api/v1/notifications", "write"},
	{"administrator", "/api/v1/notifications", "read"},
	{"administrator", "/api/v1/notifications", "write"},
	{"administrator", "/api/v1/admin/orders/*/returns", "write"},
	{"warehouse_agent", "/api/v1/admin/orders/*/returns", "write"},
	{"warehouse_agent", "/api/v1/warehouse/receptions", "read"},
	{"warehouse_agent", "/api/v1/warehouse/receptions", "write"},
	{"administrator", "/api/v1/warehouse/receptions", "read"},
	{"administrator", "/api/v1/warehouse/receptions", "write"},
	{"warehouse_agent", "/api/v1/warehouse/inventory", "read"},
	{"warehouse_agent", "/api/v1/warehouse/inventory", "write"},
	{"administrator", "/api/v1/warehouse/inventory", "read"},
	{"administrator", "/api/v1/warehouse/inventory", "write"},
	{"artisan", "/api/v1/warehouse/inventory", "read"},
}

func New() (*Service, error) {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "r.sub == p.sub && keyMatch(r.obj, p.obj) && r.act == p.act")
	e, err := casbin.NewEnforcer(m)
	if err != nil {
		return nil, fmt.Errorf("create casbin enforcer: %w", err)
	}
	if _, err := e.AddPolicies(defaultPolicies); err != nil {
		return nil, fmt.Errorf("load default casbin policies: %w", err)
	}
	return &Service{enforcer: e}, nil
}
func (s *Service) AddPolicy(role, resource, action string) error {
	ok, err := s.enforcer.AddPolicy(role, resource, action)
	if err != nil {
		return fmt.Errorf("add casbin policy: %w", err)
	}
	if !ok {
		return fmt.Errorf("casbin policy already exists")
	}
	return nil
}
func (s *Service) Authorize(_ context.Context, principal auth.Principal, resource, action string) error {
	for _, role := range principal.Roles {
		allowed, err := s.enforcer.Enforce(role, resource, action)
		if err != nil {
			return fmt.Errorf("enforce casbin policy: %w", err)
		}
		if allowed {
			return nil
		}
	}
	return ErrForbidden
}
