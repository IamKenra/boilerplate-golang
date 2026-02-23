package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint      string
	AccessKey     string
	SecretKey     string
	Bucket        string
	Region        string
	UseSSL        bool
	PublicBaseURL string
}

type FileStorage interface {
	Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string, metadata map[string]string) error
	Delete(ctx context.Context, objectName string) error
	PresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
	ObjectURL(objectName string) string
	BucketName() string
}

type minioStorage struct {
	client        *minio.Client
	bucket        string
	publicBaseURL string
}

func NewMinioStorage(cfg Config) (FileStorage, error) {
	if cfg.Endpoint == "" {
		return nil, errors.New("minio endpoint is required")
	}
	if cfg.AccessKey == "" || cfg.SecretKey == "" {
		return nil, errors.New("minio credentials are required")
	}
	if cfg.Bucket == "" {
		return nil, errors.New("minio bucket is required")
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	storage := &minioStorage{
		client:        client,
		bucket:        cfg.Bucket,
		publicBaseURL: cfg.PublicBaseURL,
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check minio bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{Region: cfg.Region}); err != nil {
			return nil, fmt.Errorf("failed to create minio bucket: %w", err)
		}
	}

	return storage, nil
}

func (s *minioStorage) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string, metadata map[string]string) error {
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}
	if len(metadata) > 0 {
		opts.UserMetadata = metadata
	}
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, opts)
	return err
}

func (s *minioStorage) Delete(ctx context.Context, objectName string) error {
	return s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
}

func (s *minioStorage) PresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *minioStorage) ObjectURL(objectName string) string {
	if s.publicBaseURL == "" {
		return ""
	}
	return strings.TrimRight(s.publicBaseURL, "/") + "/" + strings.TrimLeft(objectName, "/")
}

func (s *minioStorage) BucketName() string {
	return s.bucket
}
