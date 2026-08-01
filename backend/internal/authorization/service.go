package authorization

import (
	"context"
	"fmt"

	"github.com/aisha-platform/aisha/backend/internal/auth"
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
	{"customer", "/api/v1/artisan-applications", "write"},
	{"customer", "/api/v1/artisan-applications/me", "read"},
	{"artisan", "/api/v1/artisan/profile", "write"},
	{"administrator", "/api/v1/admin/artisan-applications", "read"},
	{"administrator", "/api/v1/admin/artisan-applications/*", "write"},
	{"administrator", "/api/v1/admin/artisan-applications/*/documents", "read"},
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
