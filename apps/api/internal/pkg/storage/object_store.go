package storage

import (
	"context"
	"io"
	"mime"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

type ObjectStore interface {
	Put(ctx context.Context, bucket, key, contentType string, body io.Reader, size int64) error
	Delete(ctx context.Context, bucket, key string) error
	PresignedGet(ctx context.Context, bucket, key string, expiry time.Duration) (*url.URL, error)
}

// DownloadPresigner is optional so existing object-store implementations can
// continue to provide preview URLs while MinIO adds a browser download hint.
type DownloadPresigner interface {
	PresignedDownload(ctx context.Context, bucket, key, filename string, expiry time.Duration) (*url.URL, error)
}

type MinIOStore struct {
	client        minioClient
	presignClient minioClient
}

type minioClient interface {
	PutObject(context.Context, string, string, io.Reader, int64, minio.PutObjectOptions) (minio.UploadInfo, error)
	RemoveObject(context.Context, string, string, minio.RemoveObjectOptions) error
	PresignedGetObject(context.Context, string, string, time.Duration, url.Values) (*url.URL, error)
}

func NewObjectStore(client *minio.Client) ObjectStore {
	return &MinIOStore{client: client, presignClient: client}
}

func NewObjectStoreWithPresigner(client, presignClient *minio.Client) ObjectStore {
	if presignClient == nil {
		presignClient = client
	}
	return &MinIOStore{client: client, presignClient: presignClient}
}

func (s *MinIOStore) Put(ctx context.Context, bucket, key, contentType string, body io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, bucket, key, body, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (s *MinIOStore) Delete(ctx context.Context, bucket, key string) error {
	return s.client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
}

func (s *MinIOStore) PresignedGet(ctx context.Context, bucket, key string, expiry time.Duration) (*url.URL, error) {
	return s.presignClient.PresignedGetObject(ctx, bucket, key, expiry, nil)
}

func (s *MinIOStore) PresignedDownload(ctx context.Context, bucket, key, filename string, expiry time.Duration) (*url.URL, error) {
	filename = strings.TrimSpace(strings.NewReplacer("\r", "", "\n", "").Replace(filename))
	if filename == "" {
		filename = "download"
	}
	disposition := mime.FormatMediaType("attachment", map[string]string{"filename": filename})
	return s.presignClient.PresignedGetObject(ctx, bucket, key, expiry, url.Values{
		"response-content-disposition": []string{disposition},
	})
}
