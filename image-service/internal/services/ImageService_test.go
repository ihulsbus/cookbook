package services

import (
	"errors"
	m "image-service/internal/models"
	tc "image-service/internal/test_common"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	imageDataDTO m.ImageDataDTO = m.ImageDataDTO{
		ID:         uuid.New(),
		EntityID:   uuid.New(),
		EntityType: "",
		Size:       0,
		Type:       "image/jpeg",
	}
	imageFileDTO m.ImageFileDTO = m.ImageFileDTO{
		ID:         imageDataDTO.ID,
		EntityID:   imageDataDTO.EntityID,
		EntityType: imageDataDTO.EntityType,
		Size:       imageDataDTO.Size,
		Type:       imageDataDTO.Type,
		File:       tc.CreateFile(),
	}
)

type s3RepositoryMock struct{}
type imageRepositoryMock struct{}
type LoggerInterfaceMock struct{}

func (s3RepositoryMock) UploadImage(imageInput m.ImageFile) error {
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

func (s3RepositoryMock) DeleteImage(imageInput m.ImageData) error {
	switch imageInput.EntityType {
	case "delete":
		return nil
	case "deleteErr":
		return nil
	default:
		return errors.New("error")
	}
}

func (i imageRepositoryMock) FindAll() ([]m.ImageData, error) {
	switch imageDataDTO.EntityType {
	case "findall":
		return []m.ImageData{imageDataDTO.ConvertFromDTO()}, nil
	case "notfound":
		return []m.ImageData{}, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (i imageRepositoryMock) Find(imageInput m.ImageData) (m.ImageData, error) {
	switch imageInput.EntityType {
	case "find":
		return imageDataDTO.ConvertFromDTO(), nil
	case "create":
		return imageDataDTO.ConvertFromDTO(), nil
	case "update":
		return imageDataDTO.ConvertFromDTO(), nil
	case "updateErr":
		return imageDataDTO.ConvertFromDTO(), nil
	case "updateS3Err":
		return imageDataDTO.ConvertFromDTO(), nil
	case "delete":
		return imageDataDTO.ConvertFromDTO(), nil
	case "deleteErr":
		return imageDataDTO.ConvertFromDTO(), nil
	case "deleteS3Err":
		return imageDataDTO.ConvertFromDTO(), nil
	case "notfound":
		return m.ImageData{}, errors.New("not found")
	default:
		return m.ImageData{}, errors.New("error")
	}
}

func (i imageRepositoryMock) Create(imageInput m.ImageData) (m.ImageData, error) {
	switch imageInput.EntityType {
	case "create":
		return imageDataDTO.ConvertFromDTO(), nil
	default:
		return m.ImageData{}, errors.New("error")
	}
}

func (i imageRepositoryMock) Update(imageInput m.ImageData) (m.ImageData, error) {
	switch imageInput.EntityType {
	case "update":
		return imageDataDTO.ConvertFromDTO(), nil
	default:
		return m.ImageData{}, errors.New("error")
	}
}

func (i imageRepositoryMock) Delete(imageInput m.ImageData) error {
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

	imageDataDTO.EntityType = "findall"
	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.IsType(t, []m.ImageDataDTO{}, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "findall", result[0].EntityType)
}

func TestFindAllImage_NotFoundErr(t *testing.T) {
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "notfound"
	result, err := s.FindAll()

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.IsType(t, []m.ImageDataDTO{}, result)
}

func TestFindAllImage_Err(t *testing.T) {
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "error"
	result, err := s.FindAll()

	assert.Error(t, err)
	assert.IsType(t, []m.ImageDataDTO{}, result)
	assert.EqualError(t, err, "internal server error")

}

func TestFindImage_OK(t *testing.T) {
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "find"
	result, err := s.Find(imageDataDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.ImageDataDTO{}, result)
	assert.Equal(t, "find", result.EntityType)
}

func TestFindImage_NotFoundErr(t *testing.T) {
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "notfound"
	result, err := s.Find(imageDataDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.IsType(t, m.ImageDataDTO{}, result)
}

func TestFindImage_Err(t *testing.T) {
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "error"
	result, err := s.Find(imageDataDTO)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDataDTO{}, result)
	assert.EqualError(t, err, "internal server error")

}

func TestCreateImage_OK(t *testing.T) {
	imageFileDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageFileDTO{
		EntityID:   imageDataDTO.EntityID,
		EntityType: "create",
		Size:       imageDataDTO.Size,
		Type:       imageDataDTO.Type,
		File:       imageFileDTO.File,
	}
	result, err := s.Create(createImage)

	assert.NoError(t, err)
	assert.IsType(t, m.ImageDataDTO{}, result)
}

func TestCreateImage_S3Err(t *testing.T) {
	imageFileDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageFileDTO{
		EntityID:   imageDataDTO.EntityID,
		EntityType: "error",
		Size:       imageDataDTO.Size,
		Type:       imageDataDTO.Type,
		File:       imageFileDTO.File,
	}
	result, err := s.Create(createImage)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDataDTO{}, result)
}

func TestCreateImage_Err(t *testing.T) {
	imageFileDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageFileDTO{
		EntityID:   imageDataDTO.EntityID,
		EntityType: "createErr",
		Size:       imageDataDTO.Size,
		Type:       imageDataDTO.Type,
		File:       imageFileDTO.File,
	}
	result, err := s.Create(createImage)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDataDTO{}, result)
}

func TestUpdateImage_OK(t *testing.T) {
	imageFileDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageFileDTO{
		EntityID:   imageDataDTO.EntityID,
		EntityType: "update",
		Size:       imageDataDTO.Size,
		Type:       imageDataDTO.Type,
		File:       imageFileDTO.File,
	}
	result, err := s.Update(createImage)

	assert.NoError(t, err)
	assert.IsType(t, m.ImageDataDTO{}, result)
}

func TestUpdateImage_FindErr(t *testing.T) {
	imageFileDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageFileDTO{
		EntityID:   imageDataDTO.EntityID,
		EntityType: "findErr",
		Size:       imageDataDTO.Size,
		Type:       imageDataDTO.Type,
		File:       imageFileDTO.File,
	}
	result, err := s.Update(createImage)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDataDTO{}, result)
}

func TestUpdateImage_S3Err(t *testing.T) {
	imageFileDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageFileDTO{
		EntityID:   imageDataDTO.EntityID,
		EntityType: "updateS3Err",
		Size:       imageDataDTO.Size,
		Type:       imageDataDTO.Type,
		File:       imageFileDTO.File,
	}
	result, err := s.Update(createImage)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDataDTO{}, result)
}

func TestUpdateImage_Err(t *testing.T) {
	imageFileDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageFileDTO{
		EntityID:   imageDataDTO.EntityID,
		EntityType: "updateErr",
		Size:       imageDataDTO.Size,
		Type:       imageDataDTO.Type,
		File:       imageFileDTO.File,
	}
	result, err := s.Update(createImage)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDataDTO{}, result)
}

func TestDeleteImage_OK(t *testing.T) {
	// imageDataDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "delete"
	err := s.Delete(imageDataDTO)

	assert.NoError(t, err)
}

func TestDeleteImage_FindErr(t *testing.T) {
	// imageDataDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "findError"
	err := s.Delete(imageDataDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "unable to find existing image. cannot delete something that does not exist")
}

func TestDeleteImage_DeleteS3Err(t *testing.T) {
	// imageDataDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "deleteS3Err"
	err := s.Delete(imageDataDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestDeleteImage_DeleteErr(t *testing.T) {
	// imageDataDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "deleteErr"
	err := s.Delete(imageDataDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
