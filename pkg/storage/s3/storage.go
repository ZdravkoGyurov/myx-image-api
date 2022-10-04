package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/ZdravkoGyurov/myx-image-api/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	s3session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	bucketName = aws.String("myx-images")
)

type Storage struct {
	svc      *s3.S3
	uploader *s3manager.Uploader
	cfg      config.S3
}

func New(cfg config.S3) *Storage {
	return &Storage{
		cfg: cfg,
	}
}

func (s *Storage) Connect() error {
	session, err := s3session.NewSessionWithOptions(s3session.Options{
		Config: aws.Config{
			S3ForcePathStyle: aws.Bool(true),
			Region:           aws.String(s.cfg.Region),
			Endpoint:         aws.String(s.cfg.Endpoint),
		},
	})
	if err != nil {
		return err
	}

	s.uploader = s3manager.NewUploader(session)
	s.svc = s3.New(session)
	return nil
}

func (s *Storage) StoreFile(ctx context.Context, fileName string, file io.Reader) (string, error) {
	result, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: bucketName,
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to store file in s3: %w", err)
	}

	return result.Location, nil
}

func (s *Storage) DeleteFile(ctx context.Context, fileName string) error {
	_, err := s.svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: bucketName,
		Key:    aws.String(fileName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from s3: %w", err)
	}

	err = s.svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: bucketName,
		Key:    aws.String(fileName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from s3: %w", err)
	}

	return nil
}
