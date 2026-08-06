package application

import (
	"context"

	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan/domain"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

type Service struct {
	repository domain.Repository
	authorizer domain.Authorizer
}

func NewService(repository domain.Repository, authorizer domain.Authorizer) *Service {
	return &Service{repository: repository, authorizer: authorizer}
}

func (s *Service) Submit(ctx context.Context, p auth.Principal, input domain.ApplicationInput) (domain.Application, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan-applications", "write"); err != nil {
		return domain.Application{}, err
	}
	if !valid(input) {
		return domain.Application{}, domain.ErrValidation
	}
	return s.repository.Submit(ctx, p.UserID, normalize(input))
}
func (s *Service) Mine(ctx context.Context, p auth.Principal) (domain.Application, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan-applications/me", "read"); err != nil {
		return domain.Application{}, err
	}
	return s.repository.Mine(ctx, p.UserID)
}
func (s *Service) UpdateProfile(ctx context.Context, p auth.Principal, input domain.ApplicationInput) (domain.Application, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan/profile", "write"); err != nil {
		return domain.Application{}, err
	}
	if !valid(input) {
		return domain.Application{}, domain.ErrValidation
	}
	return s.repository.UpdateApproved(ctx, p.UserID, normalize(input))
}
func (s *Service) List(ctx context.Context, p auth.Principal, status string, page, size int) ([]domain.Application, int, error) {
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
func (s *Service) Decide(ctx context.Context, p auth.Principal, id, decision, reason string) (domain.Application, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/artisan-applications/*", "write"); err != nil {
		return domain.Application{}, err
	}
	if decision != "APPROVED" && decision != "CHANGES_REQUESTED" && decision != "REJECTED" {
		return domain.Application{}, domain.ErrValidation
	}
	if decision != "APPROVED" && len(trim(reason)) < 2 {
		return domain.Application{}, domain.ErrValidation
	}
	return s.repository.Decide(ctx, p.UserID, id, decision, trim(reason))
}
func (s *Service) Documents(ctx context.Context, p auth.Principal, id string) ([]domain.Document, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/artisan-applications/*/documents", "read"); err != nil {
		return nil, err
	}
	return s.repository.Documents(ctx, id)
}
func valid(i domain.ApplicationInput) bool {
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
func normalize(i domain.ApplicationInput) domain.ApplicationInput {
	i.PublicDisplayName = trim(i.PublicDisplayName)
	i.InternalName = trim(i.InternalName)
	i.WorkshopName = trim(i.WorkshopName)
	i.Wilaya = trim(i.Wilaya)
	i.Location = trim(i.Location)
	i.ContactEmail = trim(i.ContactEmail)
	i.ContactPhone = trim(i.ContactPhone)
	return i
}
