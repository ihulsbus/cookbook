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

// Presigner generates a short-lived, signed GET URL for a private object -
// the mechanism behind KIT-72: images are no longer uploaded public-read,
// so the only way to read one back is a URL the service itself vouches
// for, handed out after the caller's access has already been checked
// (see ImageService.presignedURL).
type Presigner interface {
	Presign(bucket, key string, expiresIn time.Duration) (string, error)
}

// S3Presigner implements Presigner over a real *s3.S3 client, using the
// v1 SDK's GetObjectRequest+Presign pattern. Kept separate from
// S3Interface (rather than adding GetObjectRequest to it) so Presigner
// stays trivial to fake in tests - no need to construct a working
// *request.Request by hand.
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

// fileExtension maps a stored Content-Type to the extension object keys
// use. Accepts both "image/jpg" and the standard "image/jpeg" - the
// handler's verifyImage sets Content-Type from the upload's own header,
// which for a real JPEG is "image/jpeg", so only recognizing "image/jpg"
// silently produced an extensionless object key for every JPEG upload.
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

// UploadImage stores image privately (no ACL - KIT-72: images attached to
// a private recipe must not be readable by anyone with the URL). Reading
// one back goes through PresignGetURL, only after the caller's access has
// been checked.
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

// PresignGetURL returns a short-lived signed URL for reading image's
// object directly from storage. Call only after the caller's access to
// the owning recipe has already been checked (KIT-72) - this grants
// access to the bytes, unconditionally, for expiresIn.
func (r S3Repository) PresignGetURL(image m.ImageData, expiresIn time.Duration) (string, error) {
	objectPath := objectKey(image.ID, image.Type)
	return r.presigner.Presign(r.BucketName, objectPath, expiresIn)
}
