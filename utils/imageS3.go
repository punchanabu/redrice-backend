package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var s3Client *s3.Client

func init() {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           os.Getenv("BUCKET_ENDPOINT"),
			SigningRegion: os.Getenv("BUCKET_REGION"),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("BUCKET_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(os.Getenv("BUCKET_ACCESS_KEY"), os.Getenv("BUCKET_SECRET_ACCESS_KEY"), "")),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		fmt.Printf("Failed to load AWS config: %s\n", err)
		os.Exit(1)
	}

	s3Client = s3.NewFromConfig(cfg)
}

func UploadImageToS3(bucketName string, file io.Reader, fileName string) (string, error) {
	filename := uuid.New().String() + filepath.Ext(fileName)

	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	return "https://" + bucketName + ".s3.amazonaws.com/" + filename, nil
}
