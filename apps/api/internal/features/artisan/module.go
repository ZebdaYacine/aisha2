package artisan

import application "github.com/aisha-platform/aisha/apps/api/internal/features/artisan/application"
import "github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"

type Service = application.Service
type MediaService = application.MediaService

func NewService(repository Repository, authorizer Authorizer) *Service {
	return application.NewService(repository, authorizer)
}

func NewMediaService(repository MediaRepository, authorizer Authorizer, store storage.ObjectStore, bucket string, documentMax, mediaMax int64) *MediaService {
	return application.NewMediaService(repository, authorizer, store, bucket, documentMax, mediaMax)
}
