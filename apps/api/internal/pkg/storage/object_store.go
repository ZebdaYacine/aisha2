package storage

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type ObjectStore interface {
	Put(ctx context.Context, bucket, key, contentType string, body io.Reader, size int64) error
	Delete(ctx context.Context, bucket, key string) error
	PresignedGet(ctx context.Context, bucket, key string, expiry time.Duration) (*url.URL, error)
}

type MinIOStore struct{ client minioClient }

type minioClient interface {
	PutObject(context.Context, string, string, io.Reader, int64, minio.PutObjectOptions) (minio.UploadInfo, error)
	RemoveObject(context.Context, string, string, minio.RemoveObjectOptions) error
	PresignedGetObject(context.Context, string, string, time.Duration, url.Values) (*url.URL, error)
}

func NewObjectStore(client *minio.Client) ObjectStore {
	return &MinIOStore{client: client}
}

func (s *MinIOStore) Put(ctx context.Context, bucket, key, contentType string, body io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, bucket, key, body, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (s *MinIOStore) Delete(ctx context.Context, bucket, key string) error {
	return s.client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
}

func (s *MinIOStore) PresignedGet(ctx context.Context, bucket, key string, expiry time.Duration) (*url.URL, error) {
	return s.client.PresignedGetObject(ctx, bucket, key, expiry, nil)
}
