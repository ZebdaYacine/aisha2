package httpapi

import (
	"context"
	"log/slog"

	"github.com/aisha-platform/aisha/backend/core/cache"
	"github.com/aisha-platform/aisha/backend/core/config"
	"github.com/aisha-platform/aisha/backend/core/database"
	"github.com/aisha-platform/aisha/backend/core/health"
	"github.com/aisha-platform/aisha/backend/core/security"
	"github.com/aisha-platform/aisha/backend/core/storage"
	"github.com/aisha-platform/aisha/backend/features/artisans"
	"github.com/aisha-platform/aisha/backend/features/auth"
	"github.com/aisha-platform/aisha/backend/features/catalogue"
	"github.com/aisha-platform/aisha/backend/features/users"
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
	authorizationService, err := authorization.New()
	if err != nil {
		_ = redisClient.Close()
		pool.Close()
		return nil, err
	}
	rateLimiter := cache.NewRateLimiter(redisClient)
	catalogueService := catalogue.NewService(catalogue.NewPostgresRepository(pool))
	customerService := customer.NewService(customer.NewPostgresRepository(pool), authorizationService)
	artisanService := artisan.NewService(artisan.NewPostgresRepository(pool), authorizationService)
	return &Runtime{
		App: New(cfg, healthService, authService, authorizationService, rateLimiter, artisanService, customerService, catalogueService),
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
