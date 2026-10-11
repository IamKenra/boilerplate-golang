package config

import (
	"log"
	"os"
	"strings"

	"svc-boilerplate-fiber-golang/utils/storage"
)

func SetupMinioStorage() storage.FileStorage {
	fileStorage, err := storage.NewMinioStorage(storage.Config{
		Endpoint:      os.Getenv("MINIO_ENDPOINT"),
		AccessKey:     os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey:     os.Getenv("MINIO_SECRET_KEY"),
		Bucket:        os.Getenv("MINIO_BUCKET"),
		Region:        os.Getenv("MINIO_REGION"),
		UseSSL:        strings.EqualFold(os.Getenv("MINIO_USE_SSL"), "true"),
		PublicBaseURL: os.Getenv("MINIO_PUBLIC_URL"),
	})
	if err != nil {
		log.Printf("[WARN] failed to connect to MinIO, asset upload disabled: %v", err)
		return nil
	}
	return fileStorage
}
