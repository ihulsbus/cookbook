package repositories

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	m "github.com/ihulsbus/cookbook/shared/models"
)

type LoggerInterface interface {
	Error(args ...interface{})
}

type S3Interface interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}

// Presigner generates a short-lived, signed GET URL for a private object,
// since images are no longer uploaded public-read.
type Presigner interface {
	Presign(bucket, key string, expiresIn time.Duration) (string, error)
}

// S3Presigner is kept separate from S3Interface so Presigner stays trivial
// to fake in tests.
type S3Presigner struct {
	Client *s3.S3
}

func (p S3Presigner) Presign(bucket, key string, expiresIn time.Duration) (string, error) {
	req, _ := p.Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(expiresIn)
}

type S3Repository struct {
	BucketName string
	s3Client   S3Interface
	presigner  Presigner
	logger     LoggerInterface
}

func NewS3Repository(s3Client S3Interface, presigner Presigner, logger LoggerInterface, bucketName string) *S3Repository {
	return &S3Repository{
		BucketName: bucketName,
		s3Client:   s3Client,
		presigner:  presigner,
		logger:     logger,
	}
}

// Accepts both "image/jpg" and the standard "image/jpeg"; real JPEG uploads
// send "image/jpeg", so omitting it produced extensionless object keys.
func fileExtension(contentType string) string {
	switch contentType {
	case "image/jpg", "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	default:
		return ""
	}
}

func objectKey(id fmt.Stringer, contentType string) string {
	return fmt.Sprintf("img/%s.%s", id.String(), fileExtension(contentType))
}

// Stores the image with no ACL; reading it back goes through PresignGetURL
// only after the caller's access has been checked.
func (r S3Repository) UploadImage(image m.ImageFile) error {
	objectPath := objectKey(image.ID, image.Type)
	fileReader := bytes.NewReader(image.File.Bytes())

	_, err := r.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(objectPath),
		Body:   fileReader,
	})

	return err
}

func (r S3Repository) DeleteImage(image m.ImageData) error {
	objectPath := objectKey(image.ID, image.Type)

	_, err := r.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(objectPath),
	})

	return err
}

// Call only after the caller's access has already been checked - the
// returned URL grants access to the bytes unconditionally for expiresIn.
func (r S3Repository) PresignGetURL(image m.ImageData, expiresIn time.Duration) (string, error) {
	objectPath := objectKey(image.ID, image.Type)
	return r.presigner.Presign(r.BucketName, objectPath, expiresIn)
}
