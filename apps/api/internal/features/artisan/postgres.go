package artisan

import (
	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan/data/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository = repositories.PostgresRepository

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return repositories.NewPostgresRepository(pool)
}
