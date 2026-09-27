package catalogue

import application "github.com/aisha-platform/aisha/apps/api/internal/features/catalogue/application"
import "github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"

type Service = application.Service

func NewService(repository Repository) *Service { return application.NewService(repository) }
func NewServiceWithMedia(repository Repository, store storage.ObjectStore, bucket string) *Service {
	return application.NewServiceWithMedia(repository, store, bucket)
}
