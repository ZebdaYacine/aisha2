package application

import (
	"context"
	"strings"

	"github.com/aisha-platform/aisha/apps/api/internal/features/admin/domain"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

type Service struct {
	repository domain.Repository
	authorizer domain.Authorizer
	media      *MediaService
}

func NewService(repository domain.Repository, authorizer domain.Authorizer, media ...*MediaService) *Service {
	var mediaService *MediaService
	if len(media) > 0 {
		mediaService = media[0]
	}
	return &Service{repository: repository, authorizer: authorizer, media: mediaService}
}

func (s *Service) ListUserMedia(ctx context.Context, p auth.Principal, page, size int) ([]domain.UserMedia, int, error) {
	if s.media == nil {
		return nil, 0, domain.ErrValidation
	}
	return s.media.ListUserMedia(ctx, p, page, size)
}

func (s *Service) ListProductMedia(ctx context.Context, p auth.Principal, page, size int) ([]domain.ProductMedia, int, error) {
	if s.media == nil {
		return nil, 0, domain.ErrValidation
	}
	return s.media.ListProductMedia(ctx, p, page, size)
}

func (s *Service) DeleteUserMedia(ctx context.Context, p auth.Principal, id string) error {
	if s.media == nil {
		return domain.ErrValidation
	}
	return s.media.DeleteUserMedia(ctx, p, id)
}

func (s *Service) DeleteProductMedia(ctx context.Context, p auth.Principal, id string) error {
	if s.media == nil {
		return domain.ErrValidation
	}
	return s.media.DeleteProductMedia(ctx, p, id)
}

func (s *Service) ListUsers(ctx context.Context, p auth.Principal, role, status string, page, size int) ([]domain.User, int, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/users", "read"); err != nil {
		return nil, 0, err
	}
	page, size = normalizePage(page, size)
	return s.repository.ListUsers(ctx, role, status, size, (page-1)*size)
}

func (s *Service) SetRoles(ctx context.Context, p auth.Principal, userID string, roles []string) (domain.User, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/users", "write"); err != nil {
		return domain.User{}, err
	}
	if len(roles) == 0 {
		return domain.User{}, domain.ErrValidation
	}
	seen := map[string]bool{}
	for i, role := range roles {
		roles[i] = strings.TrimSpace(role)
		if roles[i] == "" || seen[roles[i]] {
			return domain.User{}, domain.ErrValidation
		}
		seen[roles[i]] = true
	}
	return s.repository.SetRoles(ctx, p.UserID, userID, roles)
}

func (s *Service) SetUserStatus(ctx context.Context, p auth.Principal, userID, status, reason string) (domain.User, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/users", "write"); err != nil {
		return domain.User{}, err
	}
	status = strings.ToUpper(strings.TrimSpace(status))
	reason = strings.TrimSpace(reason)
	if userID == "" || (status != "ACTIVE" && status != "SUSPENDED" && status != "DISABLED") || (status != "ACTIVE" && len(reason) < 2) || userID == p.UserID {
		return domain.User{}, domain.ErrValidation
	}
	return s.repository.SetUserStatus(ctx, p.UserID, userID, status, reason)
}

func (s *Service) AuditEvents(ctx context.Context, p auth.Principal, filter domain.AuditFilter, page, size int) ([]domain.AuditEvent, int, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/audit-events", "read"); err != nil {
		return nil, 0, err
	}
	page, size = normalizePage(page, size)
	return s.repository.AuditEvents(ctx, filter, size, (page-1)*size)
}

func normalizePage(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return page, size
}
