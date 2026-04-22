package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	MinioClient *minio.Client
	bucketName  string
	endpoint    string
	useSSL      bool
)

func InitMinio() {
	endpoint = os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL = os.Getenv("MINIO_USE_SSL") == "true"
	bucketName = os.Getenv("MINIO_BUCKET_NAME")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("failed to initialize minio: %v", err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("failed to check bucket existence: %v", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("failed to create bucket: %v", err)
		}
	}

	MinioClient = client
	
	// Set bucket policy to public read
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Action": ["s3:GetObject"],
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucketName)
	
	err = client.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		log.Printf("warning: failed to set bucket policy: %v", err)
	}

	log.Println("MinIO initialized successfully")
}

func Upload(ctx context.Context, reader io.Reader, size int64, objectName, contentType string) (string, error) {
	_, err := MinioClient.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}
	return objectName, nil
}

func UploadBase64(ctx context.Context, base64Str, objectName string) (string, error) {
	parts := strings.Split(base64Str, ";base64,")
	dataURI := base64Str
	if len(parts) == 2 {
		dataURI = parts[1]
	}

	contentType := "application/octet-stream"
	if len(parts) == 2 {
		mimePart := strings.Split(parts[0], ":")
		if len(mimePart) == 2 {
			contentType = mimePart[1]
		}
	}

	data, err := base64.StdEncoding.DecodeString(dataURI)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(data)
	return Upload(ctx, reader, int64(len(data)), objectName, contentType)
}

func GetURL(objectName string) string {
	if objectName == "" {
		return ""
	}
	protocol := "http"
	if useSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, endpoint, bucketName, objectName)
}
