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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	imageDTO m.ImageDTO = m.ImageDTO{
		ID:         uuid.New(),
		EntityType: "",
		EntityID:   uuid.New(),
		Size:       0,
		Type:       "image/jpeg",
		File:       createFile(),
	}
)

type s3RepositoryMock struct{}
type imageRepositoryMock struct{}
type LoggerInterfaceMock struct{}

func (s3RepositoryMock) UploadImage(imageInput m.Image) error {
	switch imageInput.EntityType {
	case "create":
		return nil
	case "update":
		return nil
	case "createErr":
		return nil
	case "updateErr":
		return nil
	default:
		return errors.New("error")
	}
}

func (s3RepositoryMock) DeleteImage(imageInput m.Image) error {
	switch imageInput.EntityType {
	case "delete":
		return nil
	case "deleteErr":
		return nil
	default:
		return errors.New("error")
	}
}

func (i imageRepositoryMock) FindAll() ([]m.Image, error) {
	switch imageDTO.EntityType {
	case "findall":
		return []m.Image{imageDTO.ConvertFromDTO()}, nil
	case "notfound":
		return []m.Image{}, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (i imageRepositoryMock) Find(imageInput m.Image) (m.Image, error) {
	switch imageInput.EntityType {
	case "find":
		return imageDTO.ConvertFromDTO(), nil
	case "create":
		return imageDTO.ConvertFromDTO(), nil
	case "update":
		return imageDTO.ConvertFromDTO(), nil
	case "updateErr":
		return imageDTO.ConvertFromDTO(), nil
	case "updateS3Err":
		return imageDTO.ConvertFromDTO(), nil
	case "delete":
		return imageDTO.ConvertFromDTO(), nil
	case "deleteErr":
		return imageDTO.ConvertFromDTO(), nil
	case "deleteS3Err":
		return imageDTO.ConvertFromDTO(), nil
	case "notfound":
		return m.Image{}, errors.New("not found")
	default:
		return m.Image{}, errors.New("error")
	}
}

func (i imageRepositoryMock) Create(imageInput m.Image) (m.Image, error) {
	switch imageInput.EntityType {
	case "create":
		return imageDTO.ConvertFromDTO(), nil
	default:
		return m.Image{}, errors.New("error")
	}
}

func (i imageRepositoryMock) Update(imageInput m.Image) (m.Image, error) {
	switch imageInput.EntityType {
	case "update":
		return imageDTO.ConvertFromDTO(), nil
	default:
		return m.Image{}, errors.New("error")
	}
}

func (i imageRepositoryMock) Delete(imageInput m.Image) error {
	switch imageInput.EntityType {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

func (LoggerInterfaceMock) Errorf(format string, args ...interface{}) {}

// ========================================================================================================

func TestFindAllImage_OK(t *testing.T) {
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "findall"
	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.IsType(t, []m.ImageDTO{}, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "findall", result[0].EntityType)
}

func TestFindAllImage_NotFoundErr(t *testing.T) {
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "notfound"
	result, err := s.FindAll()

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.IsType(t, []m.ImageDTO{}, result)
}

func TestFindAllImage_Err(t *testing.T) {
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "error"
	result, err := s.FindAll()

	assert.Error(t, err)
	assert.IsType(t, []m.ImageDTO{}, result)
	assert.EqualError(t, err, "internal server error")

}

func TestFindImage_OK(t *testing.T) {
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "find"
	result, err := s.Find(imageDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.ImageDTO{}, result)
	assert.Equal(t, "find", result.EntityType)
}

func TestFindImage_NotFoundErr(t *testing.T) {
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "notfound"
	result, err := s.Find(imageDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.IsType(t, m.ImageDTO{}, result)
}

func TestFindImage_Err(t *testing.T) {
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "error"
	result, err := s.Find(imageDTO)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDTO{}, result)
	assert.EqualError(t, err, "internal server error")

}

func TestCreateImage_OK(t *testing.T) {
	imageDTO.File = createFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageDTO{
		EntityID:   imageDTO.EntityID,
		EntityType: "create",
		Size:       imageDTO.Size,
		Type:       imageDTO.Type,
		File:       imageDTO.File,
	}
	result, err := s.Create(createImage)

	assert.NoError(t, err)
	assert.IsType(t, m.ImageDTO{}, result)
}

func TestCreateImage_S3Err(t *testing.T) {
	imageDTO.File = createFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageDTO{
		EntityID:   imageDTO.EntityID,
		EntityType: "error",
		Size:       imageDTO.Size,
		Type:       imageDTO.Type,
		File:       imageDTO.File,
	}
	result, err := s.Create(createImage)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDTO{}, result)
}

func TestCreateImage_Err(t *testing.T) {
	imageDTO.File = createFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageDTO{
		EntityID:   imageDTO.EntityID,
		EntityType: "createErr",
		Size:       imageDTO.Size,
		Type:       imageDTO.Type,
		File:       imageDTO.File,
	}
	result, err := s.Create(createImage)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDTO{}, result)
}

func TestUpdateImage_OK(t *testing.T) {
	imageDTO.File = createFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageDTO{
		EntityID:   imageDTO.EntityID,
		EntityType: "update",
		Size:       imageDTO.Size,
		Type:       imageDTO.Type,
		File:       imageDTO.File,
	}
	result, err := s.Update(createImage)

	assert.NoError(t, err)
	assert.IsType(t, m.ImageDTO{}, result)
}

func TestUpdateImage_FindErr(t *testing.T) {
	imageDTO.File = createFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageDTO{
		EntityID:   imageDTO.EntityID,
		EntityType: "findErr",
		Size:       imageDTO.Size,
		Type:       imageDTO.Type,
		File:       imageDTO.File,
	}
	result, err := s.Update(createImage)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDTO{}, result)
}

func TestUpdateImage_S3Err(t *testing.T) {
	imageDTO.File = createFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageDTO{
		EntityID:   imageDTO.EntityID,
		EntityType: "updateS3Err",
		Size:       imageDTO.Size,
		Type:       imageDTO.Type,
		File:       imageDTO.File,
	}
	result, err := s.Update(createImage)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDTO{}, result)
}

func TestUpdateImage_Err(t *testing.T) {
	imageDTO.File = createFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageDTO{
		EntityID:   imageDTO.EntityID,
		EntityType: "updateErr",
		Size:       imageDTO.Size,
		Type:       imageDTO.Type,
		File:       imageDTO.File,
	}
	result, err := s.Update(createImage)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDTO{}, result)
}

func TestDeleteImage_OK(t *testing.T) {
	imageDTO.File = createFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "delete"
	err := s.Delete(imageDTO)

	assert.NoError(t, err)
}

func TestDeleteImage_FindErr(t *testing.T) {
	imageDTO.File = createFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "findError"
	err := s.Delete(imageDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "unable to find existing image. cannot delete something that does not exist")
}

func TestDeleteImage_DeleteS3Err(t *testing.T) {
	imageDTO.File = createFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "deleteS3Err"
	err := s.Delete(imageDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestDeleteImage_DeleteErr(t *testing.T) {
	imageDTO.File = createFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "deleteErr"
	err := s.Delete(imageDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
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

func createFile() multipart.File {
	var part io.Writer
	var err error
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
		part, err = writer.CreateFormFile("file", "someimg.png")
		if err != nil {
			return
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
			return
		}
	}()
	if err != nil {
		return nil
	}

	req := httptest.NewRequest("POST", "http://example.com/v1/recipe/1/upload", pr)
	file, _, _ := req.FormFile("file")

	return file

}
