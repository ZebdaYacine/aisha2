package health

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type PostgreSQLChecker struct{ Pool *pgxpool.Pool }

func (PostgreSQLChecker) Name() string { return "postgres" }
func (c PostgreSQLChecker) Check(ctx context.Context) error {
	return c.Pool.Ping(ctx)
}

type RedisChecker struct{ Client *redis.Client }

func (RedisChecker) Name() string { return "redis" }
func (c RedisChecker) Check(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

type MinIOChecker struct{ Client *minio.Client }

func (MinIOChecker) Name() string { return "minio" }
func (c MinIOChecker) Check(ctx context.Context) error {
	_, err := c.Client.ListBuckets(ctx)
	return err
}
