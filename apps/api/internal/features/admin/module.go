package admin

import (
	"github.com/aisha-platform/aisha/apps/api/internal/features/admin/application"
	"github.com/aisha-platform/aisha/apps/api/internal/features/admin/data/repositories"
	"github.com/aisha-platform/aisha/apps/api/internal/features/admin/domain"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service = application.Service
type MediaService = application.MediaService
type UserMedia = domain.UserMedia
type ProductMedia = domain.ProductMedia

func NewService(repository Repository, authorizer Authorizer, media ...*MediaService) *Service {
	return application.NewService(repository, authorizer, media...)
}

func NewMediaService(pool *pgxpool.Pool, authorizer Authorizer, store storage.ObjectStore, artisanBucket, productBucket string) *MediaService {
	return application.NewMediaService(repositories.NewPostgresRepository(pool), authorizer, store, artisanBucket, productBucket)
}

func NewPostgresRepository(pool *pgxpool.Pool) *repositories.PostgresRepository {
	return repositories.NewPostgresRepository(pool)
}
