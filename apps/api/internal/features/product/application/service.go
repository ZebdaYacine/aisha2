package application

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/product/domain"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
	"github.com/google/uuid"
)

type Service struct {
	repository domain.Repository
	authorizer domain.Authorizer
	store      storage.ObjectStore
	bucket     string
	mediaLimit int64
}

func NewService(repository domain.Repository, authorizer domain.Authorizer, store storage.ObjectStore, bucket string, mediaLimit int64) *Service {
	return &Service{repository: repository, authorizer: authorizer, store: store, bucket: bucket, mediaLimit: mediaLimit}
}

func (s *Service) Create(ctx context.Context, p auth.Principal, input domain.Input) (domain.Product, error) {
	if err := s.authorize(ctx, p, "/api/v1/artisan/products", "write"); err != nil {
		return domain.Product{}, err
	}
	if err := validateDraft(input); err != nil {
		return domain.Product{}, err
	}
	return s.repository.Create(ctx, p.UserID, normalize(input))
}

func (s *Service) List(ctx context.Context, p auth.Principal, status string, page, size int) ([]domain.Product, int, error) {
	if err := s.authorize(ctx, p, "/api/v1/artisan/products", "read"); err != nil {
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
	items, total, err := s.repository.ListOwned(ctx, p.UserID, status, size, (page-1)*size)
	if err != nil {
		return nil, 0, err
	}
	return s.withURLs(ctx, items), total, nil
}

func (s *Service) Get(ctx context.Context, p auth.Principal, id string) (domain.Product, error) {
	if err := s.authorize(ctx, p, "/api/v1/artisan/products", "read"); err != nil {
		return domain.Product{}, err
	}
	item, err := s.repository.GetOwned(ctx, p.UserID, id)
	if err != nil {
		return domain.Product{}, err
	}
	return s.withURL(ctx, item), nil
}

func (s *Service) Update(ctx context.Context, p auth.Principal, id string, input domain.Input) (domain.Product, error) {
	if err := s.authorize(ctx, p, "/api/v1/artisan/products", "write"); err != nil {
		return domain.Product{}, err
	}
	if err := validateDraft(input); err != nil {
		return domain.Product{}, err
	}
	item, err := s.repository.Update(ctx, p.UserID, id, normalize(input))
	if err != nil {
		return domain.Product{}, err
	}
	return s.withURL(ctx, item), nil
}

func (s *Service) Submit(ctx context.Context, p auth.Principal, id string) (domain.Product, error) {
	if err := s.authorize(ctx, p, "/api/v1/artisan/products", "write"); err != nil {
		return domain.Product{}, err
	}
	item, err := s.repository.GetOwned(ctx, p.UserID, id)
	if err != nil {
		return domain.Product{}, err
	}
	if err := validateSubmission(item); err != nil {
		return domain.Product{}, err
	}
	item, err = s.repository.Submit(ctx, p.UserID, id)
	if err != nil {
		return domain.Product{}, err
	}
	return s.withURL(ctx, item), nil
}

func (s *Service) UploadMedia(ctx context.Context, p auth.Principal, productID, filename, declaredType string, data []byte, altText string) (domain.Media, error) {
	if err := s.authorize(ctx, p, "/api/v1/artisan/products", "write"); err != nil {
		return domain.Media{}, err
	}
	if s.store == nil {
		return domain.Media{}, errors.New("product media storage is unavailable")
	}
	if s.mediaLimit <= 0 {
		return domain.Media{}, storage.ErrFileTooLarge
	}
	file, err := storage.ValidateFile(data, declaredType, s.mediaLimit, map[string]bool{"IMAGE": true, "VIDEO": true})
	if err != nil {
		return domain.Media{}, err
	}
	key := fmt.Sprintf("products/%s/%s", productID, uuid.NewString())
	if err = s.store.Put(ctx, s.bucket, key, file.MediaType, bytes.NewReader(file.Bytes), int64(len(file.Bytes))); err != nil {
		return domain.Media{}, fmt.Errorf("store product media: %w", err)
	}
	media, err := s.repository.AddMedia(ctx, p.UserID, productID, domain.Media{MediaKind: file.Kind, ObjectKey: key, Checksum: file.Checksum, OriginalFilename: filename, MediaType: file.MediaType, SizeBytes: int64(len(file.Bytes)), AltText: strings.TrimSpace(altText), Visibility: "PRIVATE"})
	if err != nil {
		_ = s.store.Delete(ctx, s.bucket, key)
		return domain.Media{}, err
	}
	return s.mediaURL(ctx, media)
}

func (s *Service) DeleteMedia(ctx context.Context, p auth.Principal, productID, mediaID string) error {
	if err := s.authorize(ctx, p, "/api/v1/artisan/products", "write"); err != nil {
		return err
	}
	media, err := s.repository.DeleteMedia(ctx, p.UserID, productID, mediaID)
	if err != nil {
		return err
	}
	if s.store != nil {
		if err = s.store.Delete(ctx, s.bucket, media.ObjectKey); err != nil {
			return fmt.Errorf("delete product media object: %w", err)
		}
	}
	return nil
}

func (s *Service) authorize(ctx context.Context, p auth.Principal, resource, action string) error {
	return s.authorizer.Authorize(ctx, p, resource, action)
}

func (s *Service) withURLs(ctx context.Context, items []domain.Product) []domain.Product {
	for i := range items {
		items[i] = s.withURL(ctx, items[i])
	}
	return items
}

func (s *Service) withURL(ctx context.Context, item domain.Product) domain.Product {
	for i := range item.Media {
		item.Media[i] = s.mediaURLOrEmpty(ctx, item.Media[i])
	}
	return item
}

func (s *Service) withURLMedia(ctx context.Context, media domain.Media) domain.Media {
	return s.mediaURLOrEmpty(ctx, media)
}

func (s *Service) mediaURL(ctx context.Context, media domain.Media) (domain.Media, error) {
	if s.store == nil {
		return media, nil
	}
	value, err := s.store.PresignedGet(ctx, s.bucket, media.ObjectKey, 15*time.Minute)
	if err != nil {
		return domain.Media{}, fmt.Errorf("create product media URL: %w", err)
	}
	media.URL = value.String()
	return media, nil
}

func (s *Service) mediaURLOrEmpty(ctx context.Context, media domain.Media) domain.Media {
	value, err := s.mediaURL(ctx, media)
	if err == nil {
		return value
	}
	return media
}

func validateDraft(input domain.Input) error {
	if strings.TrimSpace(input.CategoryID) == "" || (input.ProductType != "ARTISAN_SPECIFIC" && input.ProductType != "STANDARD_TRADITIONAL") || input.PriceMinor <= 0 || len(strings.TrimSpace(input.Currency)) != 3 {
		return domain.ErrValidation
	}
	for _, translation := range input.Translations {
		if !contains(domain.SupportedLocales, translation.Locale) || strings.TrimSpace(translation.Name) == "" && strings.TrimSpace(translation.Description) == "" {
			return domain.ErrValidation
		}
	}
	return nil
}

func validateSubmission(item domain.Product) error {
	if item.Status != "DRAFT" && item.Status != "CHANGES_REQUESTED" {
		return domain.ErrInvalidTransition
	}
	if item.PriceMinor <= 0 || len(item.Translations) != len(domain.SupportedLocales) || len(item.Media) == 0 {
		return domain.ErrValidation
	}
	seen := map[string]bool{}
	for _, translation := range item.Translations {
		if seen[translation.Locale] || strings.TrimSpace(translation.Name) == "" || strings.TrimSpace(translation.Description) == "" {
			return domain.ErrValidation
		}
		seen[translation.Locale] = true
	}
	for _, locale := range domain.SupportedLocales {
		if !seen[locale] {
			return domain.ErrValidation
		}
	}
	return nil
}

func normalize(input domain.Input) domain.Input {
	input.CategoryID = strings.TrimSpace(input.CategoryID)
	input.ProductType = strings.TrimSpace(input.ProductType)
	input.Currency = strings.ToUpper(strings.TrimSpace(input.Currency))
	input.Materials = strings.TrimSpace(input.Materials)
	input.ProductionMethod = strings.TrimSpace(input.ProductionMethod)
	input.IntendedUse = strings.TrimSpace(input.IntendedUse)
	input.Dimensions = strings.TrimSpace(input.Dimensions)
	input.CountryOfOrigin = strings.TrimSpace(input.CountryOfOrigin)
	input.RegionOfOrigin = strings.TrimSpace(input.RegionOfOrigin)
	for i := range input.Translations {
		input.Translations[i].Locale = strings.TrimSpace(input.Translations[i].Locale)
		input.Translations[i].Name = strings.TrimSpace(input.Translations[i].Name)
		input.Translations[i].Description = strings.TrimSpace(input.Translations[i].Description)
		input.Translations[i].Story = strings.TrimSpace(input.Translations[i].Story)
		input.Translations[i].CulturalContext = strings.TrimSpace(input.Translations[i].CulturalContext)
	}
	return input
}

func contains(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}
