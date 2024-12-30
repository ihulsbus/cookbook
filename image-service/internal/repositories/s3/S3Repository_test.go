package repositories

import (
	"errors"
	"fmt"
	m "image-service/internal/models"
	tc "image-service/internal/test_common"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	filename string

	imgFile m.ImageFile = m.ImageFile{
		ID:   uuid.New(),
		Type: "image/jpg",
	}
	imgData m.ImageData = m.ImageData{
		ID:   imgFile.ID,
		Type: "image/jpg",
	}
)

type LoggerInterfaceMock struct{}

type S3InterfaceMock struct{}

func (LoggerInterfaceMock) Error(args ...interface{}) {}

func (S3InterfaceMock) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	name := fmt.Sprintf("img/%s.jpg", filename)
	switch *input.Key {
	case name:
		return nil, nil
	default:
		return nil, errors.New("error")
	}
}

func (S3InterfaceMock) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	name := fmt.Sprintf("img/%s.jpg", filename)
	switch *input.Key {
	case name:
		return nil, nil
	default:
		return nil, errors.New("error")
	}

}

// ========================================================================================================

func TestImageUpload_OK(t *testing.T) {

	r := NewS3Repository(&S3InterfaceMock{}, &LoggerInterfaceMock{}, "bucket")
	filename = imgFile.ID.String()
	imgFile.File = tc.CreateFile()

	err := r.UploadImage(imgFile)

	assert.NoError(t, err)
}

func TestImageUpload_PutErr(t *testing.T) {
	r := NewS3Repository(&S3InterfaceMock{}, &LoggerInterfaceMock{}, "bucket")
	filename = "filename"
	imgFile.File = tc.CreateFile()

	err := r.UploadImage(imgFile)

	assert.Error(t, err)
}

func TestImageDelete_OK(t *testing.T) {

	r := NewS3Repository(&S3InterfaceMock{}, &LoggerInterfaceMock{}, "bucket")
	filename = imgFile.ID.String()
	imgFile.File = tc.CreateFile()

	err := r.DeleteImage(imgData)

	assert.NoError(t, err)
}
