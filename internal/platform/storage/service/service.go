package service

import (
	"context"
	"mime/multipart"

	"enterprise_resource_planning/internal/platform/db"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Service struct {
	client *db.R2Client
}

func (s *Service) Upload(ctx context.Context, objKey string, file *multipart.FileHeader) error {
	src, err := file.Open()

	if err != nil {
		return err
	}

	defer src.Close()

	if _, err := s.client.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.client.Bucket),
		Key:         aws.String(objKey),
		Body:        src,
		ContentType: aws.String(file.Header.Get("Content-Type")),
	}); err != nil {
		return err
	}

	return nil
}

func New(client *db.R2Client) *Service {
	return &Service{client: client}
}
