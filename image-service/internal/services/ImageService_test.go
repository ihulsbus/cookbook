package services

import (
	"errors"
	"image"
	m "image-service/internal/models"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	imageDTO m.ImageDTO
)

type S3RepositoryMock struct{}
type LoggerInterfaceMock struct{}

func (S3RepositoryMock) UploadImage(img m.Image) error {
	switch imageDTO.EntityType {
	case "create":
		return nil
	default:
		return errors.New("error")
	}
}

func (LoggerInterfaceMock) Errorf(format string, args ...interface{}) {}

func TestUploadImage_OK(t *testing.T) {
	imageDTO.File = createFile(t)
	s := NewImageService(&S3RepositoryMock{}, &LoggerInterfaceMock{})

	result, err := s.Create(imageDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.ImageDTO{}, result)
}

func TestUploadImage_Err(t *testing.T) {
	imageDTO.File = createFile(t)
	s := NewImageService(&S3RepositoryMock{}, &LoggerInterfaceMock{})

	result, err := s.Create(imageDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.ImageDTO{}, result)
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
