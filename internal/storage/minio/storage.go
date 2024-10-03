package minio

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
)

type Storage struct {
	Client *minio.Client
}

func NewStorage(client *minio.Client) *Storage {
	return &Storage{Client: client}
}

func (s *Storage) Upload(ctx context.Context, bucketName string, objectName string, filePath string, contentType string) (*minio.UploadInfo, error) {
	info, err := s.Client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (s *Storage) GetObjectBytes(ctx context.Context, bucketName string, objectName string, filePath string, contentType string) ([]byte, error) {
	tempFile, err := os.CreateTemp("", "key-*.pem")
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	err = s.Client.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}

func (s *Storage) CreateBucket(ctx context.Context, bucketName string) error {

	err := s.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		return err
	}

	log.Printf("Successfully created %s\n", bucketName)
	return nil
}

func (s *Storage) SetupBuckets() error {
	ctx := context.Background()

	bucketName := "node-pem-files"

	exists, err := s.Client.BucketExists(ctx, bucketName)

	if err != nil {
		return err
	}

	if !exists {
		if err := s.CreateBucket(ctx, bucketName); err != nil {
			return err
		}
	}

	return nil
}
