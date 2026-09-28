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
func (r *adminRepositoryStub) SetUserStatus(context.Context, string, string, string, string) (domain.User, error) {
	return domain.User{Status: "SUSPENDED"}, nil
}
func (r *adminRepositoryStub) AuditEvents(context.Context, domain.AuditFilter, int, int) ([]domain.AuditEvent, int, error) {
	return nil, 0, nil
}

type categoryRepositoryStub struct {
	adminRepositoryStub
	input domain.CategoryInput
}

func (r *categoryRepositoryStub) ListCategories(context.Context, int, int) ([]domain.Category, int, error) {
	return nil, 0, nil
}

func (r *categoryRepositoryStub) CreateCategory(_ context.Context, _ string, input domain.CategoryInput) (domain.Category, error) {
	r.input = input
	return domain.Category{Slug: input.Slug, Translations: input.Translations}, nil
}

func (r *categoryRepositoryStub) UpdateCategory(_ context.Context, _, _ string, input domain.CategoryInput) (domain.Category, error) {
	r.input = input
	return domain.Category{Slug: input.Slug, Translations: input.Translations}, nil
}

func (r *categoryRepositoryStub) DeleteCategory(context.Context, string, string) error { return nil }

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

func TestSetUserStatusRequiresReasonAndProtectsSelf(t *testing.T) {
	repository := &adminRepositoryStub{}
	service := NewService(repository, adminAuthorizerStub{})
	if _, err := service.SetUserStatus(context.Background(), auth.Principal{UserID: "admin-1"}, "user-1", "SUSPENDED", ""); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("missing reason error=%v", err)
	}
	if _, err := service.SetUserStatus(context.Background(), auth.Principal{UserID: "admin-1"}, "admin-1", "SUSPENDED", "policy"); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("self status error=%v", err)
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

func TestCreateCategoryRequiresAllSupportedTranslations(t *testing.T) {
	repository := &categoryRepositoryStub{}
	service := NewService(repository, adminAuthorizerStub{})

	_, err := service.CreateCategory(context.Background(), auth.Principal{UserID: "admin-1"}, domain.CategoryInput{
		Slug: "jewellery", DisplayName: "Jewellery", Translations: map[string]string{"en": "Jewellery"},
	})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("missing translations error=%v", err)
	}
}

func TestCreateCategoryPassesBenefitRateAndTranslationsToRepository(t *testing.T) {
	repository := &categoryRepositoryStub{}
	service := NewService(repository, adminAuthorizerStub{})
	input := domain.CategoryInput{
		Slug: "jewellery", DisplayName: "Jewellery", BenefitRateBasisPoints: 1250,
		Translations: map[string]string{"ar": "مجوهرات", "en": "Jewellery", "fr": "Bijoux", "es": "Joyería"},
	}

	if _, err := service.CreateCategory(context.Background(), auth.Principal{UserID: "admin-1"}, input); err != nil {
		t.Fatal(err)
	}
	if repository.input.BenefitRateBasisPoints != 1250 || repository.input.Translations["ar"] != "مجوهرات" {
		t.Fatalf("input=%+v", repository.input)
	}
}
