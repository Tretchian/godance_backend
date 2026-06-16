package storage

import (
	"context"
	"errors"
	"os"
	"time"

	"godance/internal/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// MinIO реализует domain.VideoStorage поверх S3-совместимого MinIO.
//
// Используются два клиента: internal (MINIO_ENDPOINT) — для серверных операций
// (проверка/создание бакета, HEAD объекта) и presign (MINIO_PUBLIC_ENDPOINT) —
// для подписи ссылок, по которым ходит клиент. Это важно в docker: подпись
// валидна только для того хоста, на который реально пойдёт браузер.
type MinIO struct {
	client  *s3.Client
	presign *s3.PresignClient
	bucket  string
}

var _ domain.VideoStorage = (*MinIO)(nil)

func NewMinIO() (*MinIO, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	publicEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT")
	if publicEndpoint == "" {
		publicEndpoint = endpoint
	}
	region := os.Getenv("MINIO_REGION")
	if region == "" {
		region = "us-east-1"
	}

	creds := credentials.NewStaticCredentialsProvider(
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		"",
	)

	internalClient := s3.New(s3.Options{
		Region:       region,
		Credentials:  creds,
		BaseEndpoint: aws.String(endpoint),
		UsePathStyle: true,
	})
	publicClient := s3.New(s3.Options{
		Region:       region,
		Credentials:  creds,
		BaseEndpoint: aws.String(publicEndpoint),
		UsePathStyle: true,
	})

	m := &MinIO{
		client:  internalClient,
		presign: s3.NewPresignClient(publicClient),
		bucket:  os.Getenv("MINIO_BUCKET"),
	}
	if err := m.ensureBucket(context.Background()); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MinIO) ensureBucket(ctx context.Context) error {
	if _, err := m.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(m.bucket),
	}); err == nil {
		return nil
	}
	_, err := m.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(m.bucket),
	})
	var owned *s3types.BucketAlreadyOwnedByYou
	var exists *s3types.BucketAlreadyExists
	if errors.As(err, &owned) || errors.As(err, &exists) {
		return nil
	}
	return err
}

func (m *MinIO) PresignPut(ctx context.Context, key, contentType string, ttl time.Duration) (string, error) {
	req, err := m.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(m.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (m *MinIO) PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error) {
	req, err := m.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (m *MinIO) Stat(ctx context.Context, key string) (int64, error) {
	out, err := m.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if isNotFound(err) {
			return 0, domain.ErrVideoObjectMissing
		}
		return 0, err
	}
	if out.ContentLength != nil {
		return *out.ContentLength, nil
	}
	return 0, nil
}

func isNotFound(err error) bool {
	var notFound *s3types.NotFound
	if errors.As(err, &notFound) {
		return true
	}
	var respErr *smithyhttp.ResponseError
	return errors.As(err, &respErr) && respErr.HTTPStatusCode() == 404
}
