package application

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/admin/domain"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"io"
)

type mediaRepositoryStub struct{}

func (mediaRepositoryStub) ListUsers(context.Context, string, string, int, int) ([]domain.User, int, error) {
	return nil, 0, nil
}
func (mediaRepositoryStub) SetRoles(context.Context, string, string, []string) (domain.User, error) {
	return domain.User{}, nil
}
func (mediaRepositoryStub) SetUserStatus(context.Context, string, string, string, string) (domain.User, error) {
	return domain.User{}, nil
}
func (mediaRepositoryStub) AuditEvents(context.Context, domain.AuditFilter, int, int) ([]domain.AuditEvent, int, error) {
	return nil, 0, nil
}
func (mediaRepositoryStub) ListUserMedia(context.Context, int, int) ([]domain.UserMedia, int, error) {
	return []domain.UserMedia{{ID: "user-media", ObjectKey: "artisans/user/media/image"}}, 1, nil
}
func (mediaRepositoryStub) ListProductMedia(context.Context, int, int) ([]domain.ProductMedia, int, error) {
	return []domain.ProductMedia{{ID: "product-media", ObjectKey: "products/product/image"}}, 1, nil
}
func (mediaRepositoryStub) DeleteUserMedia(context.Context, string, string) (domain.UserMedia, error) {
	return domain.UserMedia{ID: "user-media", ObjectKey: "artisans/user/media/image"}, nil
}
func (mediaRepositoryStub) DeleteProductMedia(context.Context, string, string) (domain.ProductMedia, error) {
	return domain.ProductMedia{ID: "product-media", ObjectKey: "products/product/image"}, nil
}

type mediaStoreStub struct{}

func (mediaStoreStub) Put(context.Context, string, string, string, io.Reader, int64) error {
	return nil
}
func (mediaStoreStub) Delete(context.Context, string, string) error { return nil }
func (mediaStoreStub) PresignedGet(_ context.Context, bucket, key string, _ time.Duration) (*url.URL, error) {
	return url.Parse("https://media.test/" + bucket + "/" + key)
}

func TestAdminMediaServiceSignsUserAndProductObjects(t *testing.T) {
	service := NewMediaService(mediaRepositoryStub{}, adminAuthorizerStub{}, mediaStoreStub{}, "artisan-private", "product-private")
	userItems, total, err := service.ListUserMedia(context.Background(), auth.Principal{UserID: "admin"}, 1, 20)
	if err != nil || total != 1 || userItems[0].URL == "" {
		t.Fatalf("user media=%#v total=%d err=%v", userItems, total, err)
	}
	productItems, total, err := service.ListProductMedia(context.Background(), auth.Principal{UserID: "admin"}, 1, 20)
	if err != nil || total != 1 || productItems[0].URL == "" {
		t.Fatalf("product media=%#v total=%d err=%v", productItems, total, err)
	}
}

func TestAdminMediaServiceDeletesUserAndProductObjects(t *testing.T) {
	service := NewMediaService(mediaRepositoryStub{}, adminAuthorizerStub{}, mediaStoreStub{}, "artisan-private", "product-private")
	principal := auth.Principal{UserID: "admin"}
	if err := service.DeleteUserMedia(context.Background(), principal, "user-media"); err != nil {
		t.Fatalf("delete user media: %v", err)
	}
	if err := service.DeleteProductMedia(context.Background(), principal, "product-media"); err != nil {
		t.Fatalf("delete product media: %v", err)
	}
}
