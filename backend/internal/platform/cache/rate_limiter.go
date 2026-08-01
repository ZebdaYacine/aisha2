package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var incrementWithExpiry = redis.NewScript(`
local count = redis.call("INCR", KEYS[1])
if count == 1 then
  redis.call("PEXPIRE", KEYS[1], ARGV[1])
end
local ttl = redis.call("PTTL", KEYS[1])
return {count, ttl}
`)

type RateLimiter struct{ client *redis.Client }

func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{client: client}
}

func (l *RateLimiter) Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, time.Duration, error) {
	result, err := incrementWithExpiry.Run(ctx, l.client, []string{key}, window.Milliseconds()).Slice()
	if err != nil {
		return false, 0, fmt.Errorf("increment rate limit: %w", err)
	}
	if len(result) != 2 {
		return false, 0, fmt.Errorf("increment rate limit: unexpected redis response")
	}
	count, ok := result[0].(int64)
	if !ok {
		return false, 0, fmt.Errorf("increment rate limit: invalid count")
	}
	ttlMilliseconds, ok := result[1].(int64)
	if !ok {
		return false, 0, fmt.Errorf("increment rate limit: invalid ttl")
	}
	return count <= limit, time.Duration(ttlMilliseconds) * time.Millisecond, nil
}
