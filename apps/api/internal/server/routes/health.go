package routes

import (
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/health"
	"github.com/aisha-platform/aisha/apps/api/internal/server/response"
	"github.com/gofiber/fiber/v3"
)

func RegisterHealth(app fiber.Router, service *health.Service) {
	app.Get("/health/live", func(c fiber.Ctx) error {
		return c.JSON(response.HealthResponseFrom(service.Live()))
	})
	app.Get("/health/ready", func(c fiber.Ctx) error {
		result, err := service.Ready(c.Context())
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(response.HealthResponseFrom(result))
		}
		return c.JSON(response.HealthResponseFrom(result))
	})
}
