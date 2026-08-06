package cache

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestRedisRateLimiterEnforcesLimitAndExpiry(t *testing.T) {
	redisURL := os.Getenv("TEST_REDIS_URL")
	if redisURL == "" {
		t.Skip("TEST_REDIS_URL is required for the Redis rate-limit integration test")
	}
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		t.Fatal(err)
	}
	client := redis.NewClient(options)

	ctx := context.Background()
	key := "rate-limit:test:" + uuid.NewString()
	t.Cleanup(func() {
		_ = client.Del(context.Background(), key).Err()
		_ = client.Close()
	})
	limiter := NewRateLimiter(client)
	for request := 1; request <= 3; request++ {
		allowed, retryAfter, err := limiter.Allow(ctx, key, 2, time.Minute)
		if err != nil {
			t.Fatal(err)
		}
		if allowed != (request <= 2) {
			t.Fatalf("request %d allowed=%t", request, allowed)
		}
		if retryAfter <= 0 || retryAfter > time.Minute {
			t.Fatalf("unexpected retry duration: %s", retryAfter)
		}
	}
}
