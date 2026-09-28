package notification

import (
	"log/slog"
	"time"

	application "github.com/aisha-platform/aisha/apps/api/internal/features/notification/application"
	"github.com/aisha-platform/aisha/apps/api/internal/features/notification/data/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service = application.Service
type Hub = application.Hub
type Worker = application.Worker

func NewService(pool *pgxpool.Pool) *Service {
	return application.NewService(repositories.NewPostgresRepository(pool))
}

func NewWorker(pool *pgxpool.Pool, hub *Hub, interval time.Duration, logger *slog.Logger) *Worker {
	return application.NewWorker(repositories.NewPostgresRepository(pool), hub, interval, logger)
}

func NewHub() *Hub { return application.NewHub() }
