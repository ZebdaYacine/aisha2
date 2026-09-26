package application

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan/domain"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
	"github.com/google/uuid"
)

type MediaService struct {
	repository  domain.MediaRepository
	authorizer  domain.Authorizer
	store       storage.ObjectStore
	bucket      string
	documentMax int64
	mediaMax    int64
}

func NewMediaService(repository domain.MediaRepository, authorizer domain.Authorizer, store storage.ObjectStore, bucket string, documentMax, mediaMax int64) *MediaService {
	return &MediaService{repository: repository, authorizer: authorizer, store: store, bucket: bucket, documentMax: documentMax, mediaMax: mediaMax}
}

func (s *MediaService) UploadDocument(ctx context.Context, p auth.Principal, documentType, filename, declaredType string, data []byte) (domain.Document, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan-applications/me/documents", "write"); err != nil {
		return domain.Document{}, err
	}
	if s.store == nil {
		return domain.Document{}, errors.New("artisan media storage is unavailable")
	}
	if strings.TrimSpace(documentType) == "" {
		return domain.Document{}, domain.ErrValidation
	}
	file, err := storage.ValidateFile(data, declaredType, s.documentMax, map[string]bool{"DOCUMENT": true, "IMAGE": true})
	if err != nil {
		return domain.Document{}, err
	}
	profileID, status, err := s.repository.Profile(ctx, p.UserID)
	if err != nil {
		return domain.Document{}, err
	}
	if status == "SUSPENDED" {
		return domain.Document{}, domain.ErrInvalidTransition
	}
	key := fmt.Sprintf("artisans/%s/documents/%s", profileID, uuid.NewString())
	if err = s.store.Put(ctx, s.bucket, key, file.MediaType, bytes.NewReader(file.Bytes), int64(len(file.Bytes))); err != nil {
		return domain.Document{}, fmt.Errorf("store artisan document: %w", err)
	}
	item, err := s.repository.AddDocument(ctx, p.UserID, domain.DocumentUploadInput{DocumentType: strings.TrimSpace(documentType), ObjectKey: key, OriginalFilename: filename, MediaType: file.MediaType, SizeBytes: int64(len(file.Bytes)), Checksum: file.Checksum})
	if err != nil {
		_ = s.store.Delete(ctx, s.bucket, key)
		return domain.Document{}, err
	}
	if verifier, ok := s.repository.(interface {
		MarkVerificationPending(context.Context, string) error
	}); ok {
		if err = verifier.MarkVerificationPending(ctx, p.UserID); err != nil {
			return domain.Document{}, err
		}
	}
	return s.documentURL(ctx, item)
}

func (s *MediaService) OwnDocuments(ctx context.Context, p auth.Principal) ([]domain.Document, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan-applications/me/documents", "read"); err != nil {
		return nil, err
	}
	items, err := s.repository.OwnDocuments(ctx, p.UserID)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i], err = s.documentURL(ctx, items[i])
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (s *MediaService) UploadProfileMedia(ctx context.Context, p auth.Principal, filename, declaredType string, data []byte) (domain.Media, error) {
	return s.UploadProfileMediaWithOptions(ctx, p, "", filename, declaredType, data)
}

func (s *MediaService) UploadProfileMediaWithOptions(ctx context.Context, p auth.Principal, mediaKind, filename, declaredType string, data []byte) (domain.Media, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan/profile/media", "write"); err != nil {
		return domain.Media{}, err
	}
	if s.store == nil {
		return domain.Media{}, errors.New("artisan media storage is unavailable")
	}
	file, err := storage.ValidateFile(data, declaredType, s.mediaMax, map[string]bool{"IMAGE": true, "VIDEO": true})
	if err != nil {
		return domain.Media{}, err
	}
	if strings.TrimSpace(mediaKind) != "" && strings.ToUpper(strings.TrimSpace(mediaKind)) != file.Kind {
		return domain.Media{}, domain.ErrValidation
	}
	profileID, status, err := s.repository.Profile(ctx, p.UserID)
	if err != nil {
		return domain.Media{}, err
	}
	if status == "SUSPENDED" {
		return domain.Media{}, domain.ErrInvalidTransition
	}
	key := fmt.Sprintf("artisans/%s/media/%s", profileID, uuid.NewString())
	if err = s.store.Put(ctx, s.bucket, key, file.MediaType, bytes.NewReader(file.Bytes), int64(len(file.Bytes))); err != nil {
		return domain.Media{}, fmt.Errorf("store artisan profile media: %w", err)
	}
	item, err := s.repository.AddMedia(ctx, p.UserID, domain.MediaUploadInput{MediaKind: file.Kind, ObjectKey: key, OriginalFilename: filename, MediaType: file.MediaType, SizeBytes: int64(len(file.Bytes)), Checksum: file.Checksum})
	if err != nil {
		_ = s.store.Delete(ctx, s.bucket, key)
		return domain.Media{}, err
	}
	return s.mediaURL(ctx, item)
}

func (s *MediaService) ReplaceProfileMedia(ctx context.Context, p auth.Principal, id, mediaKind, filename, declaredType string, data []byte) (domain.Media, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan/profile/media", "write"); err != nil {
		return domain.Media{}, err
	}
	if s.store == nil {
		return domain.Media{}, errors.New("artisan media storage is unavailable")
	}
	file, err := storage.ValidateFile(data, declaredType, s.mediaMax, map[string]bool{"IMAGE": true, "VIDEO": true})
	if err != nil {
		return domain.Media{}, err
	}
	if strings.TrimSpace(mediaKind) != "" && strings.ToUpper(strings.TrimSpace(mediaKind)) != file.Kind {
		return domain.Media{}, domain.ErrValidation
	}
	profileID, status, err := s.repository.Profile(ctx, p.UserID)
	if err != nil {
		return domain.Media{}, err
	}
	if status == "SUSPENDED" {
		return domain.Media{}, domain.ErrInvalidTransition
	}
	key := fmt.Sprintf("artisans/%s/media/%s", profileID, uuid.NewString())
	if err = s.store.Put(ctx, s.bucket, key, file.MediaType, bytes.NewReader(file.Bytes), int64(len(file.Bytes))); err != nil {
		return domain.Media{}, fmt.Errorf("store artisan profile media: %w", err)
	}
	item, oldKey, err := s.repository.ReplaceMedia(ctx, p.UserID, id, domain.MediaUploadInput{MediaKind: file.Kind, ObjectKey: key, OriginalFilename: filename, MediaType: file.MediaType, SizeBytes: int64(len(file.Bytes)), Checksum: file.Checksum})
	if err != nil {
		_ = s.store.Delete(ctx, s.bucket, key)
		return domain.Media{}, err
	}
	if oldKey != "" && oldKey != key {
		_ = s.store.Delete(ctx, s.bucket, oldKey)
	}
	return s.mediaURL(ctx, item)
}

func (s *MediaService) DeleteProfileMedia(ctx context.Context, p auth.Principal, id string) error {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan/profile/media", "write"); err != nil {
		return err
	}
	if _, status, err := s.repository.Profile(ctx, p.UserID); err != nil {
		return err
	} else if status == "SUSPENDED" {
		return domain.ErrInvalidTransition
	}
	objectKey, err := s.repository.DeleteMedia(ctx, p.UserID, id)
	if err != nil {
		return err
	}
	if s.store != nil && objectKey != "" {
		if err = s.store.Delete(ctx, s.bucket, objectKey); err != nil {
			return fmt.Errorf("delete artisan profile media object: %w", err)
		}
	}
	return nil
}

func (s *MediaService) OwnMedia(ctx context.Context, p auth.Principal) ([]domain.Media, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/artisan/profile/media", "read"); err != nil {
		return nil, err
	}
	items, err := s.repository.OwnMedia(ctx, p.UserID)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i], err = s.mediaURL(ctx, items[i])
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (s *MediaService) AdminDocuments(ctx context.Context, p auth.Principal, id string) ([]domain.Document, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/artisan-applications/*/documents", "read"); err != nil {
		return nil, err
	}
	items, err := s.repository.Documents(ctx, id)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i], err = s.documentURL(ctx, items[i])
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (s *MediaService) AdminMedia(ctx context.Context, p auth.Principal, id string) ([]domain.Media, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/artisan-applications/*/media", "read"); err != nil {
		return nil, err
	}
	repository, ok := s.repository.(interface {
		Media(context.Context, string) ([]domain.Media, error)
	})
	if !ok {
		return nil, domain.ErrInvalidTransition
	}
	items, err := repository.Media(ctx, id)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i], err = s.mediaURL(ctx, items[i])
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (s *MediaService) documentURL(ctx context.Context, item domain.Document) (domain.Document, error) {
	if strings.HasPrefix(item.ObjectKey, "/images/") {
		item.URL = item.ObjectKey
		return item, nil
	}
	if s.store == nil {
		return item, nil
	}
	value, err := s.store.PresignedGet(ctx, s.bucket, item.ObjectKey, 15*time.Minute)
	if err != nil {
		return item, err
	}
	item.URL = value.String()
	return item, nil
}

func (s *MediaService) mediaURL(ctx context.Context, item domain.Media) (domain.Media, error) {
	if item.Visibility == "PUBLIC" && strings.HasPrefix(item.ObjectKey, "/images/") {
		item.URL = item.ObjectKey
		return item, nil
	}
	if s.store == nil {
		return item, nil
	}
	value, err := s.store.PresignedGet(ctx, s.bucket, item.ObjectKey, 15*time.Minute)
	if err != nil {
		return item, err
	}
	item.URL = value.String()
	return item, nil
}
