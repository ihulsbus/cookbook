package repositories

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	m "image-service/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	filename string
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

func newMockDatabase(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {

	var mockDB *gorm.DB

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	sqlMockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sql mock init failed: %v", err.Error())
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlMockDB,
		PreferSimpleProtocol: true,
	})

	mockDB, err = gorm.Open(dialector, &gorm.Config{
		NowFunc: timeFunc,
		Logger:  newLogger,
	})
	if err != nil {
		t.Fatalf("gorm mock init failed: %v", err.Error())
	}

	return mockDB, mock
}

func timeFunc() time.Time {
	time, _ := time.Parse("2006-01-02 15:04", "2023-02-04 18:00")
	return time
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// ========================================================================================================

func TestImageUpload_OK(t *testing.T) {
	var err error
	db, mock := newMockDatabase(t)
	r := NewS3Repository(db, m.S3Config{}, &S3InterfaceMock{}, &LoggerInterfaceMock{})
	filename = "filename"
	file := createFile(t)

	if err != nil {
		t.Error(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "recipes" SET "image_name"=$1,"updated_at"=$2 WHERE ID = $3 AND "recipes"."deleted_at" IS NULL]`).
		WithArgs(
			"filename",
			AnyTime{},
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result := r.UploadImage(file, filename, 1)

	assert.True(t, result)
}

func TestImageUpload_PutErr(t *testing.T) {
	var err error
	db, mock := newMockDatabase(t)
	r := NewS3Repository(db, m.S3Config{}, &S3InterfaceMock{}, &LoggerInterfaceMock{})
	filename = "filename"
	file := createFile(t)

	if err != nil {
		t.Error(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "recipes" SET "image_name"=$1,"updated_at"=$2 WHERE ID = $3 AND "recipes"."deleted_at" IS NULL]`).
		WithArgs(
			"filename",
			AnyTime{},
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result := r.UploadImage(file, "falsefilename", 1)

	assert.False(t, result)
}

func TestImageUpload_DBErr(t *testing.T) {
	var err error
	db, mock := newMockDatabase(t)
	r := NewS3Repository(db, m.S3Config{}, &S3InterfaceMock{}, &LoggerInterfaceMock{})
	filename = "filename"
	file := createFile(t)

	if err != nil {
		t.Error(err)
	}

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "recipes" SET "image_name"=$1,"updated_at"=$2 WHERE ID = $3 AND "recipes"."deleted_at" IS NULL]`).
		WithArgs(
			"filename",
			AnyTime{},
			1,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result := r.UploadImage(file, filename, 1)

	assert.False(t, result)
}

// ====== Helpers ======
func createImage() *image.RGBA {
	width := 200
	height := 100

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{100, 200, 200, 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}

	return img
}

func createFile(t *testing.T) multipart.File {
	// Set up a pipe to avoid buffering
	pr, pw := io.Pipe()
	// This writer is going to transform
	// what we pass to it to multipart form data
	// and write it to our io.Pipe
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		// We create the form data field 'fileupload'
		// which returns another writer to write the actual file
		part, err := writer.CreateFormFile("file", "someimg.png")
		if err != nil {
			t.Error(err)
		}

		// https://yourbasic.org/golang/create-image/
		img := createImage()

		// Encode() takes an io.Writer.
		// We pass the multipart field
		// 'fileupload' that we defined
		// earlier which, in turn, writes
		// to our io.Pipe
		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	req := httptest.NewRequest("POST", "http://example.com/v1/recipe/1/upload", pr)
	file, _, _ := req.FormFile("file")

	return file

}
