package container

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

// Infrastructure groups process-owned adapters. Feature modules depend on
// their interfaces, while this package owns concrete lifecycle wiring.
type Infrastructure struct {
	Database *pgxpool.Pool
	Redis    *redis.Client
	Storage  *minio.Client
}
