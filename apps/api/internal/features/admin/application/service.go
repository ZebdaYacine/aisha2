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
}

func NewService(repository domain.Repository, authorizer domain.Authorizer) *Service {
	return &Service{repository: repository, authorizer: authorizer}
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
