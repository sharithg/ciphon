package minio

import (
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func New() (*minio.Client, error) {
	endpoint := os.Getenv("MINIIO_HOST")
	accessKeyID := os.Getenv("MINIIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIIO_SECRET_ACCESS_KEY")
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
