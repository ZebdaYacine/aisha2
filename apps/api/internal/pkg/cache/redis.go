package cache

import (
	"fmt"

	"github.com/aisha-platform/aisha/apps/api/internal/config"
	"github.com/redis/go-redis/v9"
)

func Open(cfg config.Config) (*redis.Client, error) {
	options, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis configuration: %w", err)
	}
	return redis.NewClient(options), nil
}
