package domain

import (
	"context"
	"errors"
	"time"

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
	MembershipStatus                                                                                                                                 string `json:"membershipStatus"`
}
type Document struct {
	ID               string    `json:"id"`
	DocumentType     string    `json:"documentType"`
	ObjectKey        string    `json:"-"`
	OriginalFilename string    `json:"originalFilename"`
	MediaType        string    `json:"mediaType"`
	SizeBytes        int64     `json:"sizeBytes"`
	CreatedAt        time.Time `json:"createdAt"`
	URL              string    `json:"url,omitempty"`
}
type Media struct {
	ID               string    `json:"id"`
	MediaKind        string    `json:"mediaKind"`
	ObjectKey        string    `json:"-"`
	OriginalFilename string    `json:"originalFilename"`
	MediaType        string    `json:"mediaType"`
	SizeBytes        int64     `json:"sizeBytes"`
	SortOrder        int       `json:"sortOrder"`
	Visibility       string    `json:"visibility"`
	CreatedAt        time.Time `json:"createdAt"`
	URL              string    `json:"url,omitempty"`
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
type Workshop struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	Wilaya       string `json:"wilaya,omitempty"`
	Location     string `json:"location,omitempty"`
	Status       string `json:"status"`
	IsDefault    bool   `json:"isDefault"`
	IsPublic     bool   `json:"isPublic"`
	ProductCount int    `json:"productCount"`
}
type WorkshopInput struct {
	Name, Description, Wilaya, Location string
	IsPublic                            bool
}
type Verification struct {
	ID           string     `json:"id"`
	MembershipID string     `json:"membershipId"`
	Status       string     `json:"status"`
	Reason       string     `json:"reason"`
	DecidedBy    string     `json:"decidedBy,omitempty"`
	DecidedAt    *time.Time `json:"decidedAt,omitempty"`
}
type VerificationDecision struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type Repository interface {
	Submit(context.Context, string, ApplicationInput) (Application, error)
	Mine(context.Context, string) (Application, error)
	UpdateApproved(context.Context, string, ApplicationInput) (Application, error)
	List(context.Context, string, int, int) ([]Application, int, error)
	Decide(context.Context, string, string, string, string) (Application, error)
	Documents(context.Context, string) ([]Document, error)
}
type DraftRepository interface {
	SaveDraft(context.Context, string, ApplicationInput) (Application, error)
	FinalizeSubmission(context.Context, string) (Application, error)
}
type WorkflowRepository interface {
	ActivateMembership(context.Context, string, WorkshopInput, string, string) (Application, error)
	ListWorkshops(context.Context, string) ([]Workshop, error)
	CreateWorkshop(context.Context, string, WorkshopInput, string, string) (Workshop, error)
	UpdateWorkshop(context.Context, string, string, WorkshopInput) (Workshop, error)
	SetWorkshopStatus(context.Context, string, string, string) (Workshop, error)
	SetWorkshopStatusAdmin(context.Context, string, string, string, string) (Workshop, error)
	DeleteWorkshop(context.Context, string, string) error
	SetMembershipStatus(context.Context, string, string, string, string) (Application, error)
	MineVerification(context.Context, string) (Verification, error)
	ListVerifications(context.Context, string, int, int) ([]Verification, int, error)
	DecideVerification(context.Context, string, VerificationDecision) (Verification, error)
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
	ReplaceMedia(context.Context, string, string, MediaUploadInput) (Media, string, error)
	DeleteMedia(context.Context, string, string) (string, error)
}
