package application

import (
	"context"
	"errors"
	"io"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan/domain"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

type mediaRepositoryStub struct {
	profileID  string
	status     string
	document   domain.Document
	media      domain.Media
	mediaInput domain.MediaUploadInput
}

func (r mediaRepositoryStub) Profile(context.Context, string) (string, string, error) {
	return r.profileID, r.status, nil
}
func (r mediaRepositoryStub) Documents(context.Context, string) ([]domain.Document, error) {
	return []domain.Document{r.document}, nil
}
func (r *mediaRepositoryStub) AddDocument(context.Context, string, domain.DocumentUploadInput) (domain.Document, error) {
	return r.document, nil
}
func (r mediaRepositoryStub) OwnDocuments(context.Context, string) ([]domain.Document, error) {
	return []domain.Document{r.document}, nil
}
func (r *mediaRepositoryStub) AddMedia(_ context.Context, _ string, input domain.MediaUploadInput) (domain.Media, error) {
	r.mediaInput = input
	if r.media.ID == "" {
		r.media = domain.Media{ID: "media-1", ObjectKey: input.ObjectKey, MediaKind: input.MediaKind, OriginalFilename: input.OriginalFilename, MediaType: input.MediaType, SizeBytes: input.SizeBytes, Visibility: "PRIVATE"}
	}
	return r.media, nil
}
func (r mediaRepositoryStub) OwnMedia(context.Context, string) ([]domain.Media, error) {
	return nil, nil
}
func (r mediaRepositoryStub) ReplaceMedia(context.Context, string, string, domain.MediaUploadInput) (domain.Media, string, error) {
	return r.media, "old-key", nil
}
func (r mediaRepositoryStub) DeleteMedia(context.Context, string, string) (string, error) {
	return "old-key", nil
}

type mediaAuthorizerStub struct{ err error }

func (a mediaAuthorizerStub) Authorize(context.Context, auth.Principal, string, string) error {
	return a.err
}

type mediaStoreStub struct{ key string }

func (s *mediaStoreStub) Put(_ context.Context, _ string, key, _ string, body io.Reader, _ int64) error {
	s.key = key
	_, err := io.ReadAll(body)
	return err
}
func (s *mediaStoreStub) Delete(context.Context, string, string) error { return nil }
func (s *mediaStoreStub) PresignedGet(context.Context, string, string, time.Duration) (*url.URL, error) {
	return url.Parse("https://media.example.test/private-document")
}

func TestUploadDocumentUsesPrivateProfileKeyAndPresignedURL(t *testing.T) {
	store := &mediaStoreStub{}
	repository := &mediaRepositoryStub{profileID: "profile-1", status: "SUBMITTED", document: domain.Document{ID: "document-1", ObjectKey: "artisans/profile-1/documents/document"}}
	service := NewMediaService(repository, mediaAuthorizerStub{}, store, "aisha-artisan-private", 1024, 1024)

	item, err := service.UploadDocument(context.Background(), auth.Principal{UserID: "user-1"}, "BUSINESS_REGISTRATION", "registration.pdf", "application/pdf", []byte("%PDF-1.7 document"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(store.key, "artisans/profile-1/documents/") || item.URL == "" {
		t.Fatalf("key=%q url=%q", store.key, item.URL)
	}
}

func TestUploadProfileMediaUsesPrivateBucketAndPresignedURL(t *testing.T) {
	store := &mediaStoreStub{}
	repository := &mediaRepositoryStub{profileID: "profile-1", status: "APPROVED"}
	service := NewMediaService(repository, mediaAuthorizerStub{}, store, "aisha-artisan-private", 1024, 1024)

	item, err := service.UploadProfileMedia(context.Background(), auth.Principal{UserID: "user-1"}, "profile.png", "image/png", []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 1})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(store.key, "artisans/profile-1/media/") || item.URL == "" {
		t.Fatalf("key=%q url=%q", store.key, item.URL)
	}
	if item.Visibility != "PRIVATE" {
		t.Fatalf("profile media visibility=%q", item.Visibility)
	}
	if repository.mediaInput.ObjectKey != store.key || repository.mediaInput.MediaType != "image/png" {
		t.Fatalf("media input=%#v store key=%q", repository.mediaInput, store.key)
	}
}

func TestReplaceProfileMediaStoresNewObjectAndReturnsSignedURL(t *testing.T) {
	store := &mediaStoreStub{}
	repository := &mediaRepositoryStub{profileID: "profile-1", status: "APPROVED", media: domain.Media{ID: "media-1", ObjectKey: "artisans/profile-1/media/old"}}
	service := NewMediaService(repository, mediaAuthorizerStub{}, store, "aisha-artisan-private", 1024, 1024)

	item, err := service.ReplaceProfileMedia(context.Background(), auth.Principal{UserID: "user-1"}, "media-1", "IMAGE", "updated.png", "image/png", []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 1})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(store.key, "artisans/profile-1/media/") || item.URL == "" {
		t.Fatalf("key=%q url=%q", store.key, item.URL)
	}
}

func TestDeleteProfileMediaRemovesOwnedObject(t *testing.T) {
	service := NewMediaService(&mediaRepositoryStub{profileID: "profile-1", status: "APPROVED"}, mediaAuthorizerStub{}, &mediaStoreStub{}, "aisha-artisan-private", 1024, 1024)
	if err := service.DeleteProfileMedia(context.Background(), auth.Principal{UserID: "user-1"}, "media-1"); err != nil {
		t.Fatal(err)
	}
}

func TestUploadDocumentStopsWhenForbidden(t *testing.T) {
	denied := errors.New("denied")
	service := NewMediaService(&mediaRepositoryStub{profileID: "profile-1"}, mediaAuthorizerStub{err: denied}, &mediaStoreStub{}, "private", 1024, 1024)

	_, err := service.UploadDocument(context.Background(), auth.Principal{UserID: "user-1"}, "ID", "id.pdf", "application/pdf", []byte("%PDF-1.7"))
	if !errors.Is(err, denied) {
		t.Fatalf("err=%v", err)
	}
}
