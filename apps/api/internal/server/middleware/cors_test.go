package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestCORSAllowsConfiguredLocalAndVPSOrigins(t *testing.T) {
	app := fiber.New()
	app.Use(CORS([]string{"http://localhost:3033", "http://167.86.79.16"}))
	app.Get("/health", func(c fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	for _, origin := range []string{"http://localhost:3033", "http://167.86.79.16"} {
		request := httptest.NewRequest("GET", "/health", nil)
		request.Header.Set("Origin", origin)
		response, err := app.Test(request)
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != fiber.StatusNoContent {
			t.Fatalf("origin %s returned status %d", origin, response.StatusCode)
		}
		if got := response.Header.Get("Access-Control-Allow-Origin"); got != origin {
			t.Fatalf("origin %s returned allow-origin %q", origin, got)
		}
		if got := response.Header.Get("Access-Control-Allow-Credentials"); got != "true" {
			t.Fatalf("origin %s returned allow-credentials %q", origin, got)
		}
	}
}
