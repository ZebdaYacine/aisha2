package storage

import (
	"fmt"

	"github.com/aisha-platform/aisha/apps/api/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Open(cfg config.Config) (*minio.Client, error) {
	return OpenAt(cfg, cfg.MinIOEndpoint)
}

func OpenAt(cfg config.Config, endpoint string) (*minio.Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
		// MinIO uses the default S3 region. Supplying it prevents the SDK
		// from making a bucket-location request to a browser-facing endpoint
		// that may not be reachable from inside the API container.
		Region: "us-east-1",
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}
	return client, nil
}
