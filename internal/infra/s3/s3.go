package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/estella-studio/atr-backend/internal/infra/env"
)

type S3Itf interface {
	Upload(ctx context.Context, objectKey string, object []byte) error
}

type S3 struct {
	Client          *s3.Client
	bucketName      string
	accountID       string
	accessKeyID     string
	accessKeySecret string
}

func NewS3(env *env.Env) S3Itf {
	S3 := S3{
		bucketName:      env.S3BucketName,
		accountID:       env.S3AccountID,
		accessKeyID:     env.S3AccessKeyID,
		accessKeySecret: env.S3AccessKeySecret,
	}

	client := New(&S3)

	S3.Client = client

	return &S3
}

func New(s *S3) *s3.Client {
	config, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s.accessKeyID, s.accessKeySecret, "")),
		config.WithRegion("auto"))
	if err != nil {
		log.Panic(err)
	}

	client := s3.NewFromConfig(config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", s.accountID))
	})

	return client
}

func (s *S3) Upload(ctx context.Context, objectKey string, object []byte) error {
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(object),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			log.Printf("S3: %v\n", "EntityTooLarge")
		} else {
			log.Printf("S3: %v\n", "can't upload file")
		}
	}

	return err
}
