package storage

import (
	"fmt"

	"github.com/aisha-platform/aisha/backend/internal/platform/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Open(cfg config.Config) (*minio.Client, error) {
	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}
	return client, nil
}
