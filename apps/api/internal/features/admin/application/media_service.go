package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/admin/domain"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
)

const mediaResource = "/api/v1/admin/media"

type MediaService struct {
	repository    domain.MediaRepository
	authorizer    domain.Authorizer
	store         storage.ObjectStore
	artisanBucket string
	productBucket string
}

func NewMediaService(repository domain.MediaRepository, authorizer domain.Authorizer, store storage.ObjectStore, artisanBucket, productBucket string) *MediaService {
	return &MediaService{repository: repository, authorizer: authorizer, store: store, artisanBucket: artisanBucket, productBucket: productBucket}
}

func (s *MediaService) ListUserMedia(ctx context.Context, principal auth.Principal, page, size int) ([]domain.UserMedia, int, error) {
	if err := s.authorizer.Authorize(ctx, principal, mediaResource, "read"); err != nil {
		return nil, 0, err
	}
	page, size = normalizePage(page, size)
	items, total, err := s.repository.ListUserMedia(ctx, size, (page-1)*size)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].URL, err = s.mediaURL(ctx, s.artisanBucket, items[i].ObjectKey)
		if err != nil {
			return nil, 0, err
		}
	}
	return items, total, nil
}

func (s *MediaService) ListProductMedia(ctx context.Context, principal auth.Principal, page, size int) ([]domain.ProductMedia, int, error) {
	if err := s.authorizer.Authorize(ctx, principal, mediaResource, "read"); err != nil {
		return nil, 0, err
	}
	page, size = normalizePage(page, size)
	items, total, err := s.repository.ListProductMedia(ctx, size, (page-1)*size)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].URL, err = s.mediaURL(ctx, s.productBucket, items[i].ObjectKey)
		if err != nil {
			return nil, 0, err
		}
	}
	return items, total, nil
}

func (s *MediaService) DeleteUserMedia(ctx context.Context, principal auth.Principal, id string) error {
	if err := s.authorizer.Authorize(ctx, principal, mediaResource, "write"); err != nil {
		return err
	}
	if strings.TrimSpace(id) == "" {
		return domain.ErrValidation
	}
	item, err := s.repository.DeleteUserMedia(ctx, principal.UserID, id)
	if err != nil {
		return err
	}
	return s.deleteObject(ctx, s.artisanBucket, item.ObjectKey)
}

func (s *MediaService) DeleteProductMedia(ctx context.Context, principal auth.Principal, id string) error {
	if err := s.authorizer.Authorize(ctx, principal, mediaResource, "write"); err != nil {
		return err
	}
	if strings.TrimSpace(id) == "" {
		return domain.ErrValidation
	}
	item, err := s.repository.DeleteProductMedia(ctx, principal.UserID, id)
	if err != nil {
		return err
	}
	return s.deleteObject(ctx, s.productBucket, item.ObjectKey)
}

func (s *MediaService) deleteObject(ctx context.Context, bucket, key string) error {
	if s.store == nil || key == "" || strings.HasPrefix(key, "/images/") {
		return nil
	}
	if err := s.store.Delete(ctx, bucket, key); err != nil {
		return fmt.Errorf("delete media object: %w", err)
	}
	return nil
}

func (s *MediaService) mediaURL(ctx context.Context, bucket, key string) (string, error) {
	if strings.HasPrefix(key, "/images/") {
		return key, nil
	}
	if s.store == nil || key == "" {
		return "", nil
	}
	value, err := s.store.PresignedGet(ctx, bucket, key, 15*time.Minute)
	if err != nil {
		return "", fmt.Errorf("create admin media URL: %w", err)
	}
	return value.String(), nil
}
