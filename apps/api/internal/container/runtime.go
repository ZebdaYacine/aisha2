package container

import (
	"context"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/config"
	"github.com/aisha-platform/aisha/apps/api/internal/features/admin"
	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/cart"
	"github.com/aisha-platform/aisha/apps/api/internal/features/catalogue"
	"github.com/aisha-platform/aisha/apps/api/internal/features/inventory"
	"github.com/aisha-platform/aisha/apps/api/internal/features/moderation"
	"github.com/aisha-platform/aisha/apps/api/internal/features/order"
	"github.com/aisha-platform/aisha/apps/api/internal/features/product"
	"github.com/aisha-platform/aisha/apps/api/internal/features/user"
	"github.com/aisha-platform/aisha/apps/api/internal/features/warehouse"
	"github.com/aisha-platform/aisha/apps/api/internal/features/wishlist"
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
	minioPresignClient, err := storage.OpenAt(cfg, cfg.MinIOPublicEndpoint)
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
	objectStore := storage.NewObjectStoreWithPresigner(minioClient, minioPresignClient)
	artisanMediaService := artisan.NewMediaService(artisanRepository, authorizationService, objectStore, cfg.MinIOArtisanBucket, cfg.ArtisanDocumentMaxBytes, cfg.ArtisanMediaMaxBytes)
	productService := product.NewService(product.NewPostgresRepository(pool), authorizationService, objectStore, cfg.MinIOPrivateBucket, cfg.ProductMediaMaxBytes)
	adminMediaService := admin.NewMediaService(pool, authorizationService, objectStore, cfg.MinIOArtisanBucket, cfg.MinIOPrivateBucket)
	adminService := admin.NewService(admin.NewPostgresRepository(pool), authorizationService, adminMediaService)
	moderationService := moderation.NewServiceWithMedia(pool, authorizationService, objectStore, cfg.MinIOPrivateBucket)
	orderService := order.NewService(pool, authorizationService)
	warehouseService := warehouse.NewService(pool, authorizationService, objectStore, cfg.MinIOPrivateBucket, cfg.UploadMaxBytes)
	inventoryService := inventory.NewService(pool, authorizationService)
	workerCtx, cancelWorker := context.WithCancel(ctx)
	reservationWorker := inventory.NewReservationExpiryWorker(pool, time.Minute)
	go reservationWorker.Run(workerCtx)
	cartService := cart.NewService(pool, authorizationService)
	wishlistService := wishlist.NewService(pool, authorizationService)
	features := Features{Auth: authService, Customer: customerService, Artisan: artisanService, Catalogue: catalogueService, Product: productService, Admin: adminService, Moderation: moderationService, Order: orderService, Warehouse: warehouseService, Inventory: inventoryService, Cart: cartService, Wishlist: wishlistService}
	return &Runtime{
		App:            httpapi.NewWithWorkflows(cfg, healthService, features.Auth, authorizationService, rateLimiter, features.Artisan, features.Customer, features.Product, artisanMediaService, features.Admin, features.Moderation, features.Order, features.Warehouse, features.Inventory, features.Cart, features.Wishlist, features.Catalogue),
		infrastructure: infrastructure,
		cleanup: func() {
			cancelWorker()
			closeInfrastructure(infrastructure)
		},
	}, nil
}

func (r *Runtime) Close() {
	if r.cleanup != nil {
		r.cleanup()
	}
}
