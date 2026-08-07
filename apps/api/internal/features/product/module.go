package product

import (
	"github.com/aisha-platform/aisha/apps/api/internal/features/product/application"
	"github.com/aisha-platform/aisha/apps/api/internal/features/product/data/repositories"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service = application.Service

func NewService(repository Repository, authorizer Authorizer, store storage.ObjectStore, bucket string, mediaLimit int64) *Service {
	return application.NewService(repository, authorizer, store, bucket, mediaLimit)
}

func NewPostgresRepository(pool *pgxpool.Pool) *repositories.PostgresRepository {
	return repositories.NewPostgresRepository(pool)
}
