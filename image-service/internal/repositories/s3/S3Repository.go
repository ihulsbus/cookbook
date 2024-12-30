package repositories

import (
	"bytes"
	"fmt"
	m "image-service/internal/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type LoggerInterface interface {
	Error(args ...interface{})
}

type S3Interface interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}

type S3Repository struct {
	BucketName string
	s3Client   S3Interface
	logger     LoggerInterface
}

func NewS3Repository(s3Client S3Interface, logger LoggerInterface, bucketName string) *S3Repository {
	return &S3Repository{
		BucketName: bucketName,
		s3Client:   s3Client,
		logger:     logger,
	}
}

func (r S3Repository) UploadImage(image m.ImageFile) error {
	var fileExtention string

	switch image.Type {
	case "image/jpg":
		fileExtention = "jpg"
	case "image/png":
		fileExtention = "png"
	}

	objectPath := fmt.Sprintf("img/%s.%s", image.ID.String(), fileExtention)
	fileReader := bytes.NewReader(image.File.Bytes())

	_, err := r.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(objectPath),
		Body:   fileReader,
		ACL:    aws.String("public-read"),
	})

	return err
}

func (r S3Repository) DeleteImage(image m.ImageData) error {
	var fileExtention string

	switch image.Type {
	case "image/jpg":
		fileExtention = "jpg"
	case "image/png":
		fileExtention = "png"
	}

	objectPath := fmt.Sprintf("img/%s.%s", image.ID.String(), fileExtention)

	_, err := r.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(objectPath),
	})

	return err
}
