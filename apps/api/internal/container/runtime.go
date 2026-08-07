package container

import (
	"context"

	"github.com/aisha-platform/aisha/apps/api/internal/config"
	"github.com/aisha-platform/aisha/apps/api/internal/features/admin"
	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/catalogue"
	"github.com/aisha-platform/aisha/apps/api/internal/features/product"
	"github.com/aisha-platform/aisha/apps/api/internal/features/user"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/cache"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/database"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/health"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
	httpapi "github.com/aisha-platform/aisha/apps/api/internal/server"
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewRuntime)

type Runtime struct {
	App            *fiber.App
	infrastructure Infrastructure
	cleanup        func()
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
	infrastructure := Infrastructure{Database: pool, Redis: redisClient, Storage: minioClient}
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
	artisanRepository := artisan.NewPostgresRepository(pool)
	artisanService := artisan.NewService(artisanRepository, authorizationService)
	artisanMediaService := artisan.NewMediaService(artisanRepository, authorizationService, storage.NewObjectStore(minioClient), cfg.MinIOArtisanBucket, cfg.ArtisanDocumentMaxBytes, cfg.ArtisanMediaMaxBytes)
	productService := product.NewService(product.NewPostgresRepository(pool), authorizationService, storage.NewObjectStore(minioClient), cfg.MinIOPrivateBucket, cfg.ProductMediaMaxBytes)
	adminService := admin.NewService(admin.NewPostgresRepository(pool), authorizationService)
	features := Features{Auth: authService, Customer: customerService, Artisan: artisanService, Catalogue: catalogueService, Product: productService, Admin: adminService}
	return &Runtime{
		App:            httpapi.NewWithProduct(cfg, healthService, features.Auth, authorizationService, rateLimiter, features.Artisan, features.Customer, features.Product, artisanMediaService, features.Admin, features.Catalogue),
		infrastructure: infrastructure,
		cleanup: func() {
			closeInfrastructure(infrastructure)
		},
	}, nil
}

func (r *Runtime) Close() {
	if r.cleanup != nil {
		r.cleanup()
	}
}
