package application

import (
	"context"
	"net/mail"
	"strings"

	"github.com/aisha-platform/aisha/apps/api/internal/features/admin/domain"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"golang.org/x/crypto/bcrypt"
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

func (s *Service) CreateUser(ctx context.Context, p auth.Principal, input domain.UserInput) (domain.User, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/users", "write"); err != nil {
		return domain.User{}, err
	}
	if err := validateUserInput(input, true); err != nil {
		return domain.User{}, err
	}
	repository, ok := s.repository.(domain.UserManagementRepository)
	if !ok {
		return domain.User{}, domain.ErrValidation
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}
	input.Password = string(hash)
	return repository.CreateUser(ctx, p.UserID, input, p.UserID)
}

func (s *Service) UpdateUser(ctx context.Context, p auth.Principal, userID string, input domain.UserInput) (domain.User, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/users", "write"); err != nil {
		return domain.User{}, err
	}
	if strings.TrimSpace(userID) == "" || userID == p.UserID || (input.Password != "" && len(input.Password) < 12) {
		return domain.User{}, domain.ErrValidation
	}
	if err := validateUserInput(input, false); err != nil {
		return domain.User{}, err
	}
	repository, ok := s.repository.(domain.UserManagementRepository)
	if !ok {
		return domain.User{}, domain.ErrValidation
	}
	if input.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return domain.User{}, err
		}
		input.Password = string(hash)
	}
	return repository.UpdateUser(ctx, p.UserID, userID, input)
}

func validateUserInput(input domain.UserInput, requirePassword bool) error {
	if _, err := mail.ParseAddress(strings.TrimSpace(input.Email)); err != nil || !strings.Contains(strings.TrimSpace(input.Email), "@") {
		return domain.ErrValidation
	}
	if len(strings.TrimSpace(input.DisplayName)) < 2 || len(strings.TrimSpace(input.DisplayName)) > 120 || len(input.Phone) > 32 {
		return domain.ErrValidation
	}
	if requirePassword && len(input.Password) < 12 {
		return domain.ErrValidation
	}
	if len(input.Password) > 128 || len(input.Roles) == 0 {
		return domain.ErrValidation
	}
	seen := map[string]bool{}
	allowed := map[string]bool{"customer": true, "artisan": true, "moderator": true, "warehouse_agent": true, "administrator": true}
	for i, role := range input.Roles {
		input.Roles[i] = strings.TrimSpace(role)
		if input.Roles[i] == "" || seen[input.Roles[i]] || !allowed[input.Roles[i]] {
			return domain.ErrValidation
		}
		seen[input.Roles[i]] = true
	}
	return nil
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

func (s *Service) ListCategories(ctx context.Context, p auth.Principal, page, size int) ([]domain.Category, int, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/categories", "read"); err != nil {
		return nil, 0, err
	}
	repository, ok := s.repository.(domain.CategoryRepository)
	if !ok {
		return nil, 0, domain.ErrValidation
	}
	page, size = normalizePage(page, size)
	return repository.ListCategories(ctx, size, (page-1)*size)
}

func (s *Service) CreateCategory(ctx context.Context, p auth.Principal, input domain.CategoryInput) (domain.Category, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/categories", "write"); err != nil {
		return domain.Category{}, err
	}
	if err := validateCategory(input); err != nil {
		return domain.Category{}, err
	}
	repository, ok := s.repository.(domain.CategoryRepository)
	if !ok {
		return domain.Category{}, domain.ErrValidation
	}
	return repository.CreateCategory(ctx, p.UserID, input)
}

func (s *Service) UpdateCategory(ctx context.Context, p auth.Principal, id string, input domain.CategoryInput) (domain.Category, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/categories", "write"); err != nil {
		return domain.Category{}, err
	}
	if strings.TrimSpace(id) == "" {
		return domain.Category{}, domain.ErrValidation
	}
	if err := validateCategory(input); err != nil {
		return domain.Category{}, err
	}
	repository, ok := s.repository.(domain.CategoryRepository)
	if !ok {
		return domain.Category{}, domain.ErrValidation
	}
	return repository.UpdateCategory(ctx, p.UserID, id, input)
}

func (s *Service) DeleteCategory(ctx context.Context, p auth.Principal, id string) error {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/categories", "write"); err != nil {
		return err
	}
	if strings.TrimSpace(id) == "" {
		return domain.ErrValidation
	}
	repository, ok := s.repository.(domain.CategoryRepository)
	if !ok {
		return domain.ErrValidation
	}
	return repository.DeleteCategory(ctx, p.UserID, id)
}

func (s *Service) ListOrders(ctx context.Context, p auth.Principal, page, size int) ([]domain.Order, int, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/orders", "read"); err != nil {
		return nil, 0, err
	}
	repository, ok := s.repository.(domain.OrderRepository)
	if !ok {
		return nil, 0, domain.ErrValidation
	}
	page, size = normalizePage(page, size)
	return repository.ListOrders(ctx, size, (page-1)*size)
}

func validateCategory(input domain.CategoryInput) error {
	if len(strings.TrimSpace(input.Slug)) < 2 || len(input.Slug) > 80 || len(strings.TrimSpace(input.DisplayName)) < 2 || len(input.DisplayName) > 160 || input.BenefitRateBasisPoints < 0 || input.BenefitRateBasisPoints > 10000 {
		return domain.ErrValidation
	}
	for _, locale := range []string{"ar", "en", "fr", "es"} {
		if len(strings.TrimSpace(input.Translations[locale])) < 2 || len(input.Translations[locale]) > 160 {
			return domain.ErrValidation
		}
	}
	return nil
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
