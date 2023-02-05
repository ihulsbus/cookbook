package repositories

import (
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	m "github.com/ihulsbus/cookbook/internal/models"

	"gorm.io/gorm"
)

type LoggerInterface interface {
	Error(args ...interface{})
}

type S3Interface interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

type S3Repository struct {
	db       *gorm.DB
	s3Config m.S3Config
	s3Client S3Interface
	logger   LoggerInterface
}

func NewS3Repository(db *gorm.DB, s3Config m.S3Config, s3Client S3Interface, logger LoggerInterface) *S3Repository {
	return &S3Repository{
		db:       db,
		s3Config: s3Config,
		s3Client: s3Client,
		logger:   logger,
	}
}

func (r S3Repository) UploadImage(file multipart.File, filename string, recipeID int) bool {

	objectPath := fmt.Sprintf("img/%s.jpg", filename)

	_, err := r.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(r.s3Config.BucketName),
		Key:    aws.String(objectPath),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&m.Recipe{}).Where("ID = ?", recipeID).Update("image_name", filename).Error; err != nil {
			r.logger.Error(err)
			return err
		}

		return nil
	}); err != nil {
		return false
	}

	return true
}
