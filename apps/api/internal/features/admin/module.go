package admin

import (
	"github.com/aisha-platform/aisha/apps/api/internal/features/admin/application"
	"github.com/aisha-platform/aisha/apps/api/internal/features/admin/data/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service = application.Service

func NewService(repository Repository, authorizer Authorizer) *Service {
	return application.NewService(repository, authorizer)
}

func NewPostgresRepository(pool *pgxpool.Pool) *repositories.PostgresRepository {
	return repositories.NewPostgresRepository(pool)
}
