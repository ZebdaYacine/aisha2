package application

import (
	"context"

	"github.com/aisha-platform/aisha/apps/api/internal/features/catalogue/domain"
)

type Service struct{ repository domain.Repository }

func NewService(repository domain.Repository) *Service { return &Service{repository: repository} }
func (s *Service) Categories(ctx context.Context, request domain.PageRequest) (domain.Page[domain.Category], error) {
	return s.repository.Categories(ctx, request.Normalized())
}
func (s *Service) Products(ctx context.Context, request domain.PageRequest, category string) (domain.Page[domain.Product], error) {
	return s.repository.Products(ctx, request.Normalized(), category)
}
func (s *Service) Product(ctx context.Context, id, locale string) (domain.Product, error) {
	return s.repository.Product(ctx, id, domain.PageRequest{Locale: locale}.Normalized().Locale)
}
func (s *Service) Artisans(ctx context.Context, request domain.PageRequest) (domain.Page[domain.Artisan], error) {
	return s.repository.Artisans(ctx, request.Normalized())
}
func (s *Service) Artisan(ctx context.Context, id, locale string) (domain.Artisan, error) {
	return s.repository.Artisan(ctx, id, domain.PageRequest{Locale: locale}.Normalized().Locale)
}
