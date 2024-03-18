package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	endpoint := os.Getenv("BUCKET_ENDPOINT")
	accessKeyID := os.Getenv("BUCKET_ACCESS_KEY")
	secretAccessKey := os.Getenv("BUCKET_SECRET_ACCESS_KEY")
	useSSL := true

	var err error
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Printf("Failed to initialize MinIO client: %s\n", err)
		os.Exit(1)
	}
}

func UploadImageToS3(bucketName string, file io.Reader, fileName string) (string, error) {

	key := filepath.Join("images", uuid.New().String()+filepath.Ext(fileName))
	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	if contentType == "" {
		contentType = "image/jpeg"
	}
	_, err := minioClient.PutObject(context.Background(), bucketName, key, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Printf("Failed to upload to S3: %v", err)
		return "", err
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "inline")

	presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucketName, key, 24*time.Hour, reqParams)
	if err != nil {
		log.Printf("Failed to generate presigned URL: %v", err)
		return "", err
	}

	log.Printf("Successfully uploaded %s and generated presigned URL\n", key)
	return presignedURL.String(), nil
}
