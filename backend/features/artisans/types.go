package artisan

import (
	"context"
	"errors"

	"github.com/aisha-platform/aisha/backend/features/auth"
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
type Service struct {
	repository Repository
	authorizer Authorizer
}

func NewService(repository Repository, authorizer Authorizer) *Service {
	return &Service{repository: repository, authorizer: authorizer}
}

func (s *Service) Submit(ctx context.Context, p auth.Principal, input ApplicationInput) (Application, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan-applications", "write"); err != nil {
		return Application{}, err
	}
	if !valid(input) {
		return Application{}, ErrValidation
	}
	return s.repository.Submit(ctx, p.UserID, normalize(input))
}
func (s *Service) Mine(ctx context.Context, p auth.Principal) (Application, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan-applications/me", "read"); err != nil {
		return Application{}, err
	}
	return s.repository.Mine(ctx, p.UserID)
}
func (s *Service) UpdateProfile(ctx context.Context, p auth.Principal, input ApplicationInput) (Application, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan/profile", "write"); err != nil {
		return Application{}, err
	}
	if !valid(input) {
		return Application{}, ErrValidation
	}
	return s.repository.UpdateApproved(ctx, p.UserID, normalize(input))
}
func (s *Service) List(ctx context.Context, p auth.Principal, status string, page, size int) ([]Application, int, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/artisan-applications", "read"); err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return s.repository.List(ctx, status, size, (page-1)*size)
}
func (s *Service) Decide(ctx context.Context, p auth.Principal, id, decision, reason string) (Application, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/artisan-applications/*", "write"); err != nil {
		return Application{}, err
	}
	if decision != "APPROVED" && decision != "CHANGES_REQUESTED" && decision != "REJECTED" {
		return Application{}, ErrValidation
	}
	if decision != "APPROVED" && len(trim(reason)) < 2 {
		return Application{}, ErrValidation
	}
	return s.repository.Decide(ctx, p.UserID, id, decision, trim(reason))
}
func (s *Service) Documents(ctx context.Context, p auth.Principal, id string) ([]Document, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/artisan-applications/*/documents", "read"); err != nil {
		return nil, err
	}
	return s.repository.Documents(ctx, id)
}
func valid(i ApplicationInput) bool {
	if len(trim(i.PublicDisplayName)) < 2 || trim(i.Wilaya) == "" || len(i.CategoryIDs) == 0 || (i.ContactVisibility != "PRIVATE" && i.ContactVisibility != "PUBLIC") {
		return false
	}
	seen := map[string]bool{}
	for _, t := range i.Translations {
		if seen[t.Locale] || (t.Locale != "ar" && t.Locale != "fr" && t.Locale != "en" && t.Locale != "es") {
			return false
		}
		seen[t.Locale] = true
	}
	return true
}
func normalize(i ApplicationInput) ApplicationInput {
	i.PublicDisplayName = trim(i.PublicDisplayName)
	i.InternalName = trim(i.InternalName)
	i.WorkshopName = trim(i.WorkshopName)
	i.Wilaya = trim(i.Wilaya)
	i.Location = trim(i.Location)
	i.ContactEmail = trim(i.ContactEmail)
	i.ContactPhone = trim(i.ContactPhone)
	return i
}
