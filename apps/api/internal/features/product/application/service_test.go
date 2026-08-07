package application

import (
	"context"
	"errors"
	"io"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/product/domain"
)

type productRepositoryStub struct {
	item          domain.Product
	createCalled  bool
	submitCalled  bool
	media         domain.Media
	addMediaInput domain.Media
}

func (r *productRepositoryStub) Create(context.Context, string, domain.Input) (domain.Product, error) {
	r.createCalled = true
	return r.item, nil
}
func (r *productRepositoryStub) ListOwned(context.Context, string, string, int, int) ([]domain.Product, int, error) {
	return []domain.Product{r.item}, 1, nil
}
func (r *productRepositoryStub) GetOwned(context.Context, string, string) (domain.Product, error) {
	return r.item, nil
}
func (r *productRepositoryStub) Update(context.Context, string, string, domain.Input) (domain.Product, error) {
	return r.item, nil
}
func (r *productRepositoryStub) Submit(context.Context, string, string) (domain.Product, error) {
	r.submitCalled = true
	r.item.Status = "PENDING_REVIEW"
	return r.item, nil
}
func (r *productRepositoryStub) AddMedia(_ context.Context, _ string, _ string, media domain.Media) (domain.Media, error) {
	r.addMediaInput = media
	if r.media.ID == "" {
		r.media = media
	}
	return r.media, nil
}
func (r *productRepositoryStub) DeleteMedia(context.Context, string, string, string) (domain.Media, error) {
	return r.media, nil
}

type productAuthorizerStub struct{ err error }

func (a productAuthorizerStub) Authorize(context.Context, auth.Principal, string, string) error {
	return a.err
}

type objectStoreStub struct {
	bucket string
	key    string
	body   []byte
}

func (s *objectStoreStub) Put(_ context.Context, bucket, key, _ string, body io.Reader, _ int64) error {
	s.bucket = bucket
	s.key = key
	var err error
	s.body, err = io.ReadAll(body)
	return err
}
func (s *objectStoreStub) Delete(context.Context, string, string) error { return nil }
func (s *objectStoreStub) PresignedGet(context.Context, string, string, time.Duration) (*url.URL, error) {
	return url.Parse("https://media.example.test/private-object")
}

func completeProduct() domain.Product {
	return domain.Product{
		ID:         "product-1",
		Status:     "DRAFT",
		PriceMinor: 12500,
		Translations: []domain.Translation{
			{Locale: "ar", Name: "منتج", Description: "وصف"},
			{Locale: "fr", Name: "Produit", Description: "Description"},
			{Locale: "en", Name: "Product", Description: "Description"},
			{Locale: "es", Name: "Producto", Description: "Descripción"},
		},
		Media: []domain.Media{{ID: "media-1", ObjectKey: "products/product-1/image", MediaKind: "IMAGE"}},
	}
}

func TestCreateStopsBeforeRepositoryWhenForbidden(t *testing.T) {
	repository := &productRepositoryStub{}
	denied := errors.New("denied")
	service := NewService(repository, productAuthorizerStub{err: denied}, nil, "private", 1024)

	_, err := service.Create(context.Background(), auth.Principal{UserID: "user-1"}, domain.Input{CategoryID: "category-1", ProductType: "ARTISAN_SPECIFIC", PriceMinor: 1, Currency: "DZD"})
	if !errors.Is(err, denied) || repository.createCalled {
		t.Fatalf("err=%v createCalled=%v", err, repository.createCalled)
	}
}

func TestSubmitRequiresEveryLocaleAndMedia(t *testing.T) {
	repository := &productRepositoryStub{item: completeProduct()}
	repository.item.Media = nil
	service := NewService(repository, productAuthorizerStub{}, nil, "private", 1024)

	_, err := service.Submit(context.Background(), auth.Principal{UserID: "user-1"}, repository.item.ID)
	if !errors.Is(err, domain.ErrValidation) || repository.submitCalled {
		t.Fatalf("err=%v submitCalled=%v", err, repository.submitCalled)
	}

	repository.item.Media = completeProduct().Media
	item, err := service.Submit(context.Background(), auth.Principal{UserID: "user-1"}, repository.item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !repository.submitCalled || item.Status != "PENDING_REVIEW" {
		t.Fatalf("submitCalled=%v status=%q", repository.submitCalled, item.Status)
	}
}

func TestUploadMediaValidatesAndStoresPrivateObject(t *testing.T) {
	repository := &productRepositoryStub{media: domain.Media{ID: "media-1", ProductID: "product-1", ObjectKey: "products/product-1/image", MediaKind: "IMAGE"}}
	store := &objectStoreStub{}
	service := NewService(repository, productAuthorizerStub{}, store, "aisha-product-private", 1024)

	data := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 1, 2, 3}
	item, err := service.UploadMedia(context.Background(), auth.Principal{UserID: "user-1"}, "product-1", "image.png", "image/png", data, "Front view")
	if err != nil {
		t.Fatal(err)
	}
	if store.bucket != "aisha-product-private" || !strings.HasPrefix(store.key, "products/product-1/") {
		t.Fatalf("stored object = %q/%q", store.bucket, store.key)
	}
	if string(store.body) != string(data) || repository.addMediaInput.Visibility != "PRIVATE" {
		t.Fatalf("media storage metadata/body not private: %#v", repository.addMediaInput)
	}
	if item.URL == "" {
		t.Fatal("expected a presigned URL")
	}
}
