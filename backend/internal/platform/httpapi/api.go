package httpapi

import (
	"log/slog"
	"time"

	"github.com/aisha-platform/aisha/backend/internal/platform/config"
	"github.com/aisha-platform/aisha/backend/internal/platform/health"
	"github.com/aisha-platform/aisha/backend/internal/platform/httpapi/dto"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func New(cfg config.Config, healthService *health.Service) *fiber.App {
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
	return app
}
