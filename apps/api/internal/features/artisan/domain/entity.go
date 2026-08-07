package domain

import (
	"context"
	"errors"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

var (
	ErrNotFound          = errors.New("artisan application not found")
	ErrValidation        = errors.New("artisan application validation failed")
	ErrInvalidTransition = errors.New("invalid artisan application transition")
)

type Translation struct {
	Locale    string `json:"locale"`
	Biography string `json:"biography"`
}
type ApplicationInput struct {
	PublicDisplayName, InternalName, WorkshopName, Wilaya, Location, ContactEmail, ContactPhone, ContactVisibility string
	CategoryIDs                                                                                                    []string
	Translations                                                                                                   []Translation
}
type Application struct {
	ID, UserID, PublicDisplayName, InternalName, WorkshopName, Wilaya, Location, ContactEmail, ContactPhone, ContactVisibility, Status, ReviewReason string
	CategoryIDs                                                                                                                                      []string
	Translations                                                                                                                                     []Translation
}
type Document struct {
	ID               string `json:"id"`
	DocumentType     string `json:"documentType"`
	ObjectKey        string `json:"-"`
	OriginalFilename string `json:"originalFilename"`
	MediaType        string `json:"mediaType"`
	SizeBytes        int64  `json:"sizeBytes"`
	URL              string `json:"url,omitempty"`
}
type Media struct {
	ID, MediaKind, ObjectKey, OriginalFilename, MediaType string
	SizeBytes                                             int64
	SortOrder                                             int
	Visibility, URL                                       string
}
type DocumentUploadInput struct {
	DocumentType, ObjectKey, OriginalFilename, MediaType, Checksum string
	SizeBytes                                                      int64
}
type MediaUploadInput struct {
	MediaKind, ObjectKey, OriginalFilename, MediaType, Checksum string
	SizeBytes                                                   int64
	SortOrder                                                   int
}
type Repository interface {
	Submit(context.Context, string, ApplicationInput) (Application, error)
	Mine(context.Context, string) (Application, error)
	UpdateApproved(context.Context, string, ApplicationInput) (Application, error)
	List(context.Context, string, int, int) ([]Application, int, error)
	Decide(context.Context, string, string, string, string) (Application, error)
	Documents(context.Context, string) ([]Document, error)
}
type Authorizer interface {
	Authorize(context.Context, auth.Principal, string, string) error
}

type MediaRepository interface {
	Profile(context.Context, string) (string, string, error)
	Documents(context.Context, string) ([]Document, error)
	AddDocument(context.Context, string, DocumentUploadInput) (Document, error)
	OwnDocuments(context.Context, string) ([]Document, error)
	AddMedia(context.Context, string, MediaUploadInput) (Media, error)
	OwnMedia(context.Context, string) ([]Media, error)
}
