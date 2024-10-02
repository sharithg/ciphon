package storage

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	Client *minio.Client
}

func SetupMinio() *Minio {
	endpoint := os.Getenv("MINIIO_HOST")
	accessKeyID := os.Getenv("MINIIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIIO_SECRET_ACCESS_KEY")
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	minio := &Minio{Client: minioClient}

	minio.CreateBucket("node-pem-files")

	return &Minio{Client: minioClient}
}

func (m *Minio) Upload(ctx context.Context, bucketName string, objectName string, filePath string, contentType string) (*minio.UploadInfo, error) {
	info, err := m.Client.FPutObject(ctx, "node-pem-files", objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (m *Minio) CreateBucket(bucketName string) {
	ctx := context.Background()

	err := m.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		exists, errBucketExists := m.Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
}
