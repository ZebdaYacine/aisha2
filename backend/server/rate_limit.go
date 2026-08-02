package httpapi

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

type RateLimiter interface {
	Allow(context.Context, string, int64, time.Duration) (bool, time.Duration, error)
}

func authRateLimit(limiter RateLimiter, limit int64, window time.Duration, operation string) fiber.Handler {
	return func(c fiber.Ctx) error {
		if limiter == nil {
			return c.Next()
		}
		key := fmt.Sprintf("rate-limit:auth:%s:%s", operation, c.IP())
		allowed, retryAfter, err := limiter.Allow(c.Context(), key, limit, window)
		if err != nil {
			return WrapAPIError(err, CodeInternalError, "An unexpected error occurred.")
		}
		if !allowed {
			seconds := int64(retryAfter.Round(time.Second) / time.Second)
			if seconds < 1 {
				seconds = 1
			}
			c.Set(fiber.HeaderRetryAfter, strconv.FormatInt(seconds, 10))
			return NewAPIError(CodeRateLimited, "Too many requests. Please try again later.", nil)
		}
		return c.Next()
	}
}
