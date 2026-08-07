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
	profileID string
	status    string
	document  domain.Document
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
func (r mediaRepositoryStub) AddMedia(context.Context, string, domain.MediaUploadInput) (domain.Media, error) {
	return domain.Media{}, nil
}
func (r mediaRepositoryStub) OwnMedia(context.Context, string) ([]domain.Media, error) {
	return nil, nil
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

func TestUploadDocumentStopsWhenForbidden(t *testing.T) {
	denied := errors.New("denied")
	service := NewMediaService(&mediaRepositoryStub{profileID: "profile-1"}, mediaAuthorizerStub{err: denied}, &mediaStoreStub{}, "private", 1024, 1024)

	_, err := service.UploadDocument(context.Background(), auth.Principal{UserID: "user-1"}, "ID", "id.pdf", "application/pdf", []byte("%PDF-1.7"))
	if !errors.Is(err, denied) {
		t.Fatalf("err=%v", err)
	}
}
