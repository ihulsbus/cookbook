package repositories

import (
	"errors"
	"fmt"
	tc "image-service/internal/test_common"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
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

type PresignerMock struct {
	url string
	err error
}

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

func (p PresignerMock) Presign(bucket, key string, expiresIn time.Duration) (string, error) {
	if p.err != nil {
		return "", p.err
	}
	return p.url, nil
}

// ========================================================================================================

func TestImageUpload_OK(t *testing.T) {

	r := NewS3Repository(&S3InterfaceMock{}, PresignerMock{}, &LoggerInterfaceMock{}, "bucket")
	filename = imgFile.ID.String()
	imgFile.File, _ = tc.CreateFile()

	err := r.UploadImage(imgFile)

	assert.NoError(t, err)
}

func TestImageUpload_PutErr(t *testing.T) {
	r := NewS3Repository(&S3InterfaceMock{}, PresignerMock{}, &LoggerInterfaceMock{}, "bucket")
	filename = "filename"
	imgFile.File, _ = tc.CreateFile()

	err := r.UploadImage(imgFile)

	assert.Error(t, err)
}

func TestImageDelete_OK(t *testing.T) {

	r := NewS3Repository(&S3InterfaceMock{}, PresignerMock{}, &LoggerInterfaceMock{}, "bucket")
	filename = imgFile.ID.String()
	imgFile.File, _ = tc.CreateFile()

	err := r.DeleteImage(imgData)

	assert.NoError(t, err)
}

func TestImageUpload_JpegOK(t *testing.T) {
	// image/jpeg (the real Content-Type a browser sends for a JPEG, as
	// opposed to the test fixtures' "image/jpg") must resolve to the
	// same ".jpg" extension - this is the bug fileExtension fixes.
	r := NewS3Repository(&S3InterfaceMock{}, PresignerMock{}, &LoggerInterfaceMock{}, "bucket")
	filename = imgFile.ID.String()

	jpegFile := imgFile
	jpegFile.Type = "image/jpeg"
	jpegFile.File, _ = tc.CreateFile()

	err := r.UploadImage(jpegFile)

	assert.NoError(t, err)
}

func TestPresignGetURL_OK(t *testing.T) {
	r := NewS3Repository(&S3InterfaceMock{}, PresignerMock{url: "https://example.com/signed"}, &LoggerInterfaceMock{}, "bucket")

	url, err := r.PresignGetURL(imgData, 15*time.Minute)

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/signed", url)
}

func TestPresignGetURL_Err(t *testing.T) {
	r := NewS3Repository(&S3InterfaceMock{}, PresignerMock{err: errors.New("presign error")}, &LoggerInterfaceMock{}, "bucket")

	_, err := r.PresignGetURL(imgData, 15*time.Minute)

	assert.Error(t, err)
	assert.EqualError(t, err, "presign error")
}
