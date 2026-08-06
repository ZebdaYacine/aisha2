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
	ID, DocumentType, OriginalFilename, MediaType string
	SizeBytes                                     int64
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
