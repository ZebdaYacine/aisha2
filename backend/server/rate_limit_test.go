package httpapi

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

type stubRateLimiter struct {
	allowed bool
	retry   time.Duration
	err     error
	key     string
}

func (l *stubRateLimiter) Allow(_ context.Context, key string, _ int64, _ time.Duration) (bool, time.Duration, error) {
	l.key = key
	return l.allowed, l.retry, l.err
}

func TestAuthRateLimitRejectsExcessRequests(t *testing.T) {
	limiter := &stubRateLimiter{retry: 12 * time.Second}
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	app.Use(requestid.New())
	app.Post("/login", authRateLimit(limiter, 5, time.Minute, "login"), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	response, err := app.Test(httptest.NewRequest("POST", "/login", nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusTooManyRequests || response.Header.Get(fiber.HeaderRetryAfter) != "12" {
		t.Fatalf("unexpected rate-limit response: status=%d retry=%q", response.StatusCode, response.Header.Get(fiber.HeaderRetryAfter))
	}
	var body ErrorResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Error.Code != CodeRateLimited || limiter.key == "" {
		t.Fatalf("unexpected rate-limit result: %#v key=%q", body, limiter.key)
	}
}

func TestAuthRateLimitAllowsRequest(t *testing.T) {
	limiter := &stubRateLimiter{allowed: true}
	app := fiber.New()
	app.Post("/login", authRateLimit(limiter, 5, time.Minute, "login"), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})
	response, err := app.Test(httptest.NewRequest("POST", "/login", nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected request to pass, got %d", response.StatusCode)
	}
}
