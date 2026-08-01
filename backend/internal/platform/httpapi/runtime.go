package httpapi

import (
	"context"
	"log/slog"

	"github.com/aisha-platform/aisha/backend/internal/auth"
	"github.com/aisha-platform/aisha/backend/internal/platform/cache"
	"github.com/aisha-platform/aisha/backend/internal/platform/config"
	"github.com/aisha-platform/aisha/backend/internal/platform/database"
	"github.com/aisha-platform/aisha/backend/internal/platform/health"
	"github.com/aisha-platform/aisha/backend/internal/platform/storage"
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewRuntime)

type Runtime struct {
	App     *fiber.App
	cleanup func()
}

func NewRuntime(ctx context.Context, cfg config.Config) (*Runtime, error) {
	pool, err := database.Open(ctx, cfg)
	if err != nil {
		return nil, err
	}
	redisClient, err := cache.Open(cfg)
	if err != nil {
		pool.Close()
		return nil, err
	}
	minioClient, err := storage.Open(cfg)
	if err != nil {
		_ = redisClient.Close()
		pool.Close()
		return nil, err
	}
	healthService := health.New(
		health.PostgreSQLChecker{Pool: pool},
		health.RedisChecker{Client: redisClient},
		health.MinIOChecker{Client: minioClient},
	)
	authRepository := auth.NewPostgresRepository(pool)
	resetNotifier := auth.NewOutboxResetNotifier(pool)
	authService := auth.NewService(authRepository, resetNotifier, cfg.AuthSigningKey)
	return &Runtime{
		App: New(cfg, healthService, authService),
		cleanup: func() {
			if err := redisClient.Close(); err != nil {
				slog.Error("close redis client", "error", err)
			}
			pool.Close()
		},
	}, nil
}

func (r *Runtime) Close() {
	if r.cleanup != nil {
		r.cleanup()
	}
}
