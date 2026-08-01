package httpapi

import (
	"log/slog"
	"time"

	"github.com/aisha-platform/aisha/backend/internal/artisan"
	"github.com/aisha-platform/aisha/backend/internal/auth"
	"github.com/aisha-platform/aisha/backend/internal/authorization"
	"github.com/aisha-platform/aisha/backend/internal/catalogue"
	"github.com/aisha-platform/aisha/backend/internal/customer"
	"github.com/aisha-platform/aisha/backend/internal/platform/config"
	"github.com/aisha-platform/aisha/backend/internal/platform/health"
	"github.com/aisha-platform/aisha/backend/internal/platform/httpapi/dto"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func New(cfg config.Config, healthService *health.Service, authService *auth.Service, authorizationService *authorization.Service, rateLimiter RateLimiter, artisanService *artisan.Service, customerService *customer.Service, catalogueServices ...*catalogue.Service) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName + " API",
		BodyLimit:    2 * 1024 * 1024,
		ErrorHandler: errorHandler,
	})
	app.Use(requestid.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "X-Request-ID"},
	}))
	app.Use(func(c fiber.Ctx) error {
		started := time.Now()
		err := c.Next()
		slog.Info("http request",
			"request_id", c.GetRespHeader(fiber.HeaderXRequestID),
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", time.Since(started).Milliseconds(),
		)
		return err
	})
	app.Get("/health/live", func(c fiber.Ctx) error {
		return c.JSON(dto.HealthResponseFrom(healthService.Live()))
	})
	app.Get("/health/ready", func(c fiber.Ctx) error {
		result, err := healthService.Ready(c.Context())
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(dto.HealthResponseFrom(result))
		}
		return c.JSON(dto.HealthResponseFrom(result))
	})
	if authService != nil {
		handler := NewAuthHandler(authService, NewRequestValidator())
		api := app.Group("/api/v1")
		api.Post("/auth/register", authRateLimit(rateLimiter, cfg.AuthRateLimitMax, cfg.AuthRateLimitWindow, "register"), handler.Register)
		api.Post("/auth/login", authRateLimit(rateLimiter, cfg.AuthRateLimitMax, cfg.AuthRateLimitWindow, "login"), handler.Login)
		api.Post("/auth/refresh", handler.Refresh)
		api.Post("/auth/logout", handler.Logout)
		api.Post("/auth/forgot-password", authRateLimit(rateLimiter, cfg.AuthRateLimitMax, cfg.AuthRateLimitWindow, "forgot-password"), handler.ForgotPassword)
		api.Post("/auth/reset-password", authRateLimit(rateLimiter, cfg.AuthRateLimitMax, cfg.AuthRateLimitWindow, "reset-password"), handler.ResetPassword)
		if authorizationService != nil {
			casbin := NewCasbinMiddleware(authorizationService)
			api.Get("/me", handler.RequirePrincipal, casbin.Require("/api/v1/me", "read"), handler.Me)
			if customerService != nil {
				customerHandler := NewCustomerHandler(customerService, NewRequestValidator())
				api.Get("/me/profile", handler.RequirePrincipal, customerHandler.Profile)
				api.Patch("/me/profile", handler.RequirePrincipal, customerHandler.UpdateProfile)
				api.Get("/addresses", handler.RequirePrincipal, customerHandler.Addresses)
				api.Post("/addresses", handler.RequirePrincipal, customerHandler.CreateAddress)
				api.Patch("/addresses/:id", handler.RequirePrincipal, customerHandler.UpdateAddress)
				api.Delete("/addresses/:id", handler.RequirePrincipal, customerHandler.DeleteAddress)
			}
			if artisanService != nil {
				artisanHandler := NewArtisanHandler(artisanService, NewRequestValidator())
				api.Post("/artisan-applications", handler.RequirePrincipal, artisanHandler.Submit)
				api.Get("/artisan-applications/me", handler.RequirePrincipal, artisanHandler.Mine)
				api.Patch("/artisan/profile", handler.RequirePrincipal, artisanHandler.UpdateProfile)
				api.Get("/admin/artisan-applications", handler.RequirePrincipal, artisanHandler.List)
				api.Post("/admin/artisan-applications/:id/:decision", handler.RequirePrincipal, artisanHandler.Decide)
				api.Get("/admin/artisan-applications/:id/documents", handler.RequirePrincipal, artisanHandler.Documents)
			}
		}
	}
	if len(catalogueServices) > 0 && catalogueServices[0] != nil {
		handler := NewCatalogueHandler(catalogueServices[0])
		api := app.Group("/api/v1")
		api.Get("/categories", handler.Categories)
		api.Get("/products", handler.Products)
		api.Get("/products/:id", handler.Product)
		api.Get("/artisans", handler.Artisans)
		api.Get("/artisans/:id", handler.Artisan)
	}
	return app
}
