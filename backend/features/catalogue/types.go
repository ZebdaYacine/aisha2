package catalogue

import (
	"context"
	"time"
)

var supportedLocales = map[string]bool{"ar": true, "en": true, "es": true, "fr": true}

type PageRequest struct {
	Locale   string
	Page     int
	PageSize int
}

func (p PageRequest) Normalized() PageRequest {
	if !supportedLocales[p.Locale] {
		p.Locale = "en"
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p
}

type Category struct{ ID, Slug, Name string }
type Product struct {
	ID, ArtisanID, ArtisanName, CategoryID, CategorySlug, Name, Description, Story string
	Materials, ProductionMethod, Region, Currency, Status                          string
	PriceMinor                                                                     int64
	Media                                                                          []string
	PublishedAt                                                                    *time.Time
}
type Artisan struct {
	ID, Name, Workshop, Wilaya, Location, Biography string
	Media                                           []string
	ProductCount                                    int
}
type Page[T any] struct {
	Items                 []T
	Page, PageSize, Total int
}

type Repository interface {
	Categories(context.Context, PageRequest) (Page[Category], error)
	Products(context.Context, PageRequest, string) (Page[Product], error)
	Product(context.Context, string, string) (Product, error)
	Artisans(context.Context, PageRequest) (Page[Artisan], error)
	Artisan(context.Context, string, string) (Artisan, error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }
func (s *Service) Categories(ctx context.Context, request PageRequest) (Page[Category], error) {
	return s.repository.Categories(ctx, request.Normalized())
}
func (s *Service) Products(ctx context.Context, request PageRequest, category string) (Page[Product], error) {
	return s.repository.Products(ctx, request.Normalized(), category)
}
func (s *Service) Product(ctx context.Context, id, locale string) (Product, error) {
	return s.repository.Product(ctx, id, PageRequest{Locale: locale}.Normalized().Locale)
}
func (s *Service) Artisans(ctx context.Context, request PageRequest) (Page[Artisan], error) {
	return s.repository.Artisans(ctx, request.Normalized())
}
func (s *Service) Artisan(ctx context.Context, id, locale string) (Artisan, error) {
	return s.repository.Artisan(ctx, id, PageRequest{Locale: locale}.Normalized().Locale)
}
