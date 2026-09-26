package application

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

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
	if administrator(p) {
		return domain.Application{}, domain.ErrValidation
	}
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan-applications", "write"); err != nil {
		return domain.Application{}, err
	}
	if !valid(input) {
		return domain.Application{}, domain.ErrValidation
	}
	// Keep the legacy endpoint compatible with the staged workflow. When the
	// repository supports drafts, it must still pass through the same final
	// server-side document/media checks as the new endpoint.
	if _, ok := s.repository.(domain.DraftRepository); ok {
		if _, err := s.SaveDraft(ctx, p, input); err != nil {
			return domain.Application{}, err
		}
		return s.FinalizeSubmission(ctx, p)
	}
	return s.repository.Submit(ctx, p.UserID, normalize(input))
}

func (s *Service) SaveDraft(ctx context.Context, p auth.Principal, input domain.ApplicationInput) (domain.Application, error) {
	if administrator(p) {
		return domain.Application{}, domain.ErrValidation
	}
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan-applications", "write"); err != nil {
		return domain.Application{}, err
	}
	if !valid(input) {
		return domain.Application{}, domain.ErrValidation
	}
	r, ok := s.repository.(domain.DraftRepository)
	if !ok {
		return domain.Application{}, domain.ErrInvalidTransition
	}
	return r.SaveDraft(ctx, p.UserID, normalize(input))
}

func (s *Service) FinalizeSubmission(ctx context.Context, p auth.Principal) (domain.Application, error) {
	if administrator(p) {
		return domain.Application{}, domain.ErrValidation
	}
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan-applications/me/submit", "write"); err != nil {
		return domain.Application{}, err
	}
	r, ok := s.repository.(domain.DraftRepository)
	if !ok {
		return domain.Application{}, domain.ErrInvalidTransition
	}
	return r.FinalizeSubmission(ctx, p.UserID)
}
func (s *Service) Mine(ctx context.Context, p auth.Principal) (domain.Application, error) {
	if administrator(p) {
		return domain.Application{}, domain.ErrValidation
	}
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan-applications/me", "read"); err != nil {
		return domain.Application{}, err
	}
	return s.repository.Mine(ctx, p.UserID)
}
func (s *Service) UpdateProfile(ctx context.Context, p auth.Principal, input domain.ApplicationInput) (domain.Application, error) {
	if administrator(p) {
		return domain.Application{}, domain.ErrValidation
	}
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

func (s *Service) workflow() (domain.WorkflowRepository, error) {
	r, ok := s.repository.(domain.WorkflowRepository)
	if !ok {
		return nil, domain.ErrInvalidTransition
	}
	return r, nil
}

func (s *Service) ActivateMembership(ctx context.Context, p auth.Principal, input domain.WorkshopInput, idempotencyKey string) (domain.Application, error) {
	if administrator(p) {
		return domain.Application{}, domain.ErrValidation
	}
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan/membership", "write"); err != nil {
		return domain.Application{}, err
	}
	if strings.TrimSpace(idempotencyKey) == "" || strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Wilaya) == "" {
		return domain.Application{}, domain.ErrValidation
	}
	r, err := s.workflow()
	if err != nil {
		return domain.Application{}, err
	}
	return r.ActivateMembership(ctx, p.UserID, input, idempotencyKey, fmt.Sprintf("%x", sha256.Sum256(mustJSON(input))))
}
func (s *Service) ListWorkshops(ctx context.Context, p auth.Principal) ([]domain.Workshop, error) {
	if administrator(p) {
		return nil, domain.ErrValidation
	}
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan/workshops", "read"); err != nil {
		return nil, err
	}
	r, err := s.workflow()
	if err != nil {
		return nil, err
	}
	return r.ListWorkshops(ctx, p.UserID)
}
func (s *Service) CreateWorkshop(ctx context.Context, p auth.Principal, input domain.WorkshopInput, key string) (domain.Workshop, error) {
	if administrator(p) {
		return domain.Workshop{}, domain.ErrValidation
	}
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan/workshops", "write"); err != nil {
		return domain.Workshop{}, err
	}
	if strings.TrimSpace(key) == "" || strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Wilaya) == "" {
		return domain.Workshop{}, domain.ErrValidation
	}
	r, err := s.workflow()
	if err != nil {
		return domain.Workshop{}, err
	}
	return r.CreateWorkshop(ctx, p.UserID, input, key, fmt.Sprintf("%x", sha256.Sum256(mustJSON(input))))
}
func (s *Service) UpdateWorkshop(ctx context.Context, p auth.Principal, id string, input domain.WorkshopInput) (domain.Workshop, error) {
	if administrator(p) {
		return domain.Workshop{}, domain.ErrValidation
	}
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan/workshops", "write"); err != nil {
		return domain.Workshop{}, err
	}
	if id == "" || strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Wilaya) == "" {
		return domain.Workshop{}, domain.ErrValidation
	}
	r, err := s.workflow()
	if err != nil {
		return domain.Workshop{}, err
	}
	return r.UpdateWorkshop(ctx, p.UserID, id, input)
}
func (s *Service) SetWorkshopStatus(ctx context.Context, p auth.Principal, id, status string) (domain.Workshop, error) {
	if administrator(p) {
		return domain.Workshop{}, domain.ErrValidation
	}
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan/workshops", "write"); err != nil {
		return domain.Workshop{}, err
	}
	status = strings.ToUpper(strings.TrimSpace(status))
	if status != "ACTIVE" && status != "INACTIVE" && status != "ARCHIVED" {
		return domain.Workshop{}, domain.ErrValidation
	}
	r, err := s.workflow()
	if err != nil {
		return domain.Workshop{}, err
	}
	return r.SetWorkshopStatus(ctx, p.UserID, id, status)
}
func (s *Service) SetWorkshopStatusAdmin(ctx context.Context, p auth.Principal, id, status, reason string) (domain.Workshop, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/workshops", "write"); err != nil {
		return domain.Workshop{}, err
	}
	status = strings.ToUpper(strings.TrimSpace(status))
	reason = strings.TrimSpace(reason)
	if id == "" || (status != "ACTIVE" && status != "INACTIVE" && status != "ARCHIVED") || (status != "ACTIVE" && len(reason) < 2) {
		return domain.Workshop{}, domain.ErrValidation
	}
	r, err := s.workflow()
	if err != nil {
		return domain.Workshop{}, err
	}
	return r.SetWorkshopStatusAdmin(ctx, p.UserID, id, status, reason)
}
func (s *Service) DeleteWorkshop(ctx context.Context, p auth.Principal, id string) error {
	if administrator(p) {
		return domain.ErrValidation
	}
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan/workshops", "write"); err != nil {
		return err
	}
	r, err := s.workflow()
	if err != nil {
		return err
	}
	return r.DeleteWorkshop(ctx, p.UserID, id)
}
func (s *Service) SetMembershipStatus(ctx context.Context, p auth.Principal, id, status, reason string) (domain.Application, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/artisan-memberships", "write"); err != nil {
		return domain.Application{}, err
	}
	status = strings.ToUpper(strings.TrimSpace(status))
	reason = strings.TrimSpace(reason)
	if status != "ACTIVE" && status != "SUSPENDED" && status != "CLOSED" || (status == "SUSPENDED" && reason == "") {
		return domain.Application{}, domain.ErrValidation
	}
	r, err := s.workflow()
	if err != nil {
		return domain.Application{}, err
	}
	return r.SetMembershipStatus(ctx, p.UserID, id, status, reason)
}
func (s *Service) MineVerification(ctx context.Context, p auth.Principal) (domain.Verification, error) {
	if administrator(p) {
		return domain.Verification{}, domain.ErrValidation
	}
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan/verification", "read"); err != nil {
		return domain.Verification{}, err
	}
	r, err := s.workflow()
	if err != nil {
		return domain.Verification{}, err
	}
	return r.MineVerification(ctx, p.UserID)
}
func (s *Service) ListVerifications(ctx context.Context, p auth.Principal, status string, page, size int) ([]domain.Verification, int, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/artisan-verifications", "read"); err != nil {
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
	r, err := s.workflow()
	if err != nil {
		return nil, 0, err
	}
	return r.ListVerifications(ctx, status, size, (page-1)*size)
}
func (s *Service) DecideVerification(ctx context.Context, p auth.Principal, input domain.VerificationDecision) (domain.Verification, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/artisan-verifications", "write"); err != nil {
		return domain.Verification{}, err
	}
	input.Status = strings.ToUpper(strings.TrimSpace(input.Status))
	input.Reason = strings.TrimSpace(input.Reason)
	if input.Status != "VERIFIED" && input.Status != "CHANGES_REQUESTED" && input.Status != "REJECTED" || input.ID == "" || input.Status != "VERIFIED" && input.Reason == "" {
		return domain.Verification{}, domain.ErrValidation
	}
	r, err := s.workflow()
	if err != nil {
		return domain.Verification{}, err
	}
	return r.DecideVerification(ctx, p.UserID, input)
}
func mustJSON(v any) []byte { b, _ := json.Marshal(v); return b }

func administrator(p auth.Principal) bool {
	for _, role := range p.Roles {
		if role == "administrator" {
			return true
		}
	}
	return false
}
