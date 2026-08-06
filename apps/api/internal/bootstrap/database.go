package bootstrap

import (
	"context"

	"github.com/aisha-platform/aisha/apps/api/internal/config"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OpenDatabase is the bootstrap-owned database entry point for migrations,
// health checks, and container wiring. It does not execute migrations.
func OpenDatabase(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	return database.Open(ctx, cfg)
}
