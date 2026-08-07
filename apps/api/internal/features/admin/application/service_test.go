package application

import (
	"context"
	"errors"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/features/admin/domain"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

type adminRepositoryStub struct{ roles []string }

func (r *adminRepositoryStub) ListUsers(context.Context, string, string, int, int) ([]domain.User, int, error) {
	return nil, 0, nil
}
func (r *adminRepositoryStub) SetRoles(_ context.Context, _, _ string, roles []string) (domain.User, error) {
	r.roles = roles
	return domain.User{Roles: roles}, nil
}
func (r *adminRepositoryStub) AuditEvents(context.Context, domain.AuditFilter, int, int) ([]domain.AuditEvent, int, error) {
	return nil, 0, nil
}

type adminAuthorizerStub struct{ err error }

func (a adminAuthorizerStub) Authorize(context.Context, auth.Principal, string, string) error {
	return a.err
}

func TestSetRolesTrimsAndDeduplicatesBeforeRepository(t *testing.T) {
	repository := &adminRepositoryStub{}
	service := NewService(repository, adminAuthorizerStub{})

	item, err := service.SetRoles(context.Background(), auth.Principal{UserID: "admin-1"}, "user-1", []string{" admin", "artisan "})
	if err != nil {
		t.Fatal(err)
	}
	if len(item.Roles) != 2 || item.Roles[0] != "admin" || item.Roles[1] != "artisan" {
		t.Fatalf("roles=%v", item.Roles)
	}
}

func TestSetRolesRejectsDuplicateRole(t *testing.T) {
	repository := &adminRepositoryStub{}
	service := NewService(repository, adminAuthorizerStub{})

	_, err := service.SetRoles(context.Background(), auth.Principal{UserID: "admin-1"}, "user-1", []string{"admin", "admin"})
	if !errors.Is(err, domain.ErrValidation) || repository.roles != nil {
		t.Fatalf("err=%v roles=%v", err, repository.roles)
	}
}
