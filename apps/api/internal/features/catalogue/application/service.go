package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/catalogue/domain"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
)

type Service struct {
	repository domain.Repository
	store      storage.ObjectStore
	bucket     string
}

func NewService(repository domain.Repository) *Service { return &Service{repository: repository} }
func NewServiceWithMedia(repository domain.Repository, store storage.ObjectStore, bucket string) *Service {
	return &Service{repository: repository, store: store, bucket: bucket}
}

func (s *Service) Categories(ctx context.Context, request domain.PageRequest) (domain.Page[domain.Category], error) {
	return s.repository.Categories(ctx, request.Normalized())
}
func (s *Service) Products(ctx context.Context, request domain.PageRequest, category, query, workshop string) (domain.Page[domain.Product], error) {
	result, err := s.repository.Products(ctx, request.Normalized(), category, query, workshop)
	if err != nil {
		return result, err
	}
	for i := range result.Items {
		result.Items[i].Media, err = s.resolveMedia(ctx, result.Items[i].Media)
		if err != nil {
			return domain.Page[domain.Product]{}, err
		}
	}
	return result, nil
}
func (s *Service) Product(ctx context.Context, id, locale string) (domain.Product, error) {
	item, err := s.repository.Product(ctx, id, domain.PageRequest{Locale: locale}.Normalized().Locale)
	if err != nil {
		return item, err
	}
	item.Media, err = s.resolveMedia(ctx, item.Media)
	return item, err
}

func (s *Service) resolveMedia(ctx context.Context, media []string) ([]string, error) {
	if s.store == nil || s.bucket == "" {
		return media, nil
	}
	resolved := make([]string, len(media))
	for i, key := range media {
		resolved[i] = key
		if key == "" || strings.HasPrefix(key, "/") || strings.HasPrefix(key, "http://") || strings.HasPrefix(key, "https://") {
			continue
		}
		url, err := s.store.PresignedGet(ctx, s.bucket, key, 15*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("sign public product media: %w", err)
		}
		resolved[i] = url.String()
	}
	return resolved, nil
}
func (s *Service) Artisans(ctx context.Context, request domain.PageRequest) (domain.Page[domain.Artisan], error) {
	return s.repository.Artisans(ctx, request.Normalized())
}
func (s *Service) Artisan(ctx context.Context, id, locale string) (domain.Artisan, error) {
	return s.repository.Artisan(ctx, id, domain.PageRequest{Locale: locale}.Normalized().Locale)
}

func (s *Service) Workshops(ctx context.Context, request domain.PageRequest) (domain.Page[domain.Workshop], error) {
	return s.repository.Workshops(ctx, request.Normalized())
}

func (s *Service) Workshop(ctx context.Context, id, locale string) (domain.Workshop, error) {
	return s.repository.Workshop(ctx, id, domain.PageRequest{Locale: locale}.Normalized().Locale)
}
