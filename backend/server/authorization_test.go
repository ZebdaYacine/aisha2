package httpapi

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/aisha-platform/aisha/backend/core/security"
	"github.com/aisha-platform/aisha/backend/features/auth"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func TestCasbinMiddlewareDistinguishesUnauthenticatedAndForbidden(t *testing.T) {
	service, err := authorization.New()
	if err != nil {
		t.Fatal(err)
	}
	middleware := NewCasbinMiddleware(service)
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	app.Use(requestid.New())
	app.Get("/unauthenticated", middleware.Require("/api/v1/me", "read"), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})
	app.Get("/forbidden", func(c fiber.Ctx) error {
		c.Locals(principalLocal, auth.Principal{Roles: []string{"visitor"}})
		return c.Next()
	}, middleware.Require("/api/v1/me", "read"), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	assertErrorCode(t, app, "/unauthenticated", fiber.StatusUnauthorized, CodeAuthenticationRequired)
	assertErrorCode(t, app, "/forbidden", fiber.StatusForbidden, CodeForbidden)
}

func TestCasbinMiddlewareAllowsConfiguredRole(t *testing.T) {
	service, err := authorization.New()
	if err != nil {
		t.Fatal(err)
	}
	middleware := NewCasbinMiddleware(service)
	app := fiber.New()
	app.Get("/me", func(c fiber.Ctx) error {
		c.Locals(principalLocal, auth.Principal{Roles: []string{"customer"}})
		return c.Next()
	}, middleware.Require("/api/v1/me", "read"), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})
	response, err := app.Test(httptest.NewRequest("GET", "/me", nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected allowed request, got %d", response.StatusCode)
	}
}

func assertErrorCode(t *testing.T, app *fiber.App, path string, status int, code ErrorCode) {
	t.Helper()
	response, err := app.Test(httptest.NewRequest("GET", path, nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != status {
		t.Fatalf("expected %d, got %d", status, response.StatusCode)
	}
	var body ErrorResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Error.Code != code {
		t.Fatalf("expected %s, got %s", code, body.Error.Code)
	}
}
