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
type databaseRepositoryMock struct{}
type rabbitmqRepositoryMock struct{}
type LoggerInterfaceMock struct{}

func (s3RepositoryMock) UploadImage(imageInput m.ImageFile) error {
	switch imageInput.EntityType {
	case "uploadS3Fail":
		return errors.New("error")
	default:
		return nil
	}
}

func (s3RepositoryMock) DeleteImage(imageInput m.ImageData) error {
	switch imageInput.EntityType {
	case "deleteS3Fail":
		return errors.New("error")
	default:
		return nil
	}
}

func (i databaseRepositoryMock) FindAll() ([]m.ImageData, error) {
	switch imageDataDTO.EntityType {
	case "findallFail":
		return nil, errors.New("error")
	case "notfound":
		return []m.ImageData{}, errors.New("not found")
	default:
		return []m.ImageData{imageDataDTO.ConvertFromDTO()}, nil
	}
}

func (i databaseRepositoryMock) Find(imageInput m.ImageData) (m.ImageData, error) {
	switch imageInput.EntityType {
	case "findFail":
		return m.ImageData{}, errors.New("error")
	case "notfound":
		return m.ImageData{}, errors.New("not found")
	default:
		return imageDataDTO.ConvertFromDTO(), nil
	}
}

func (i databaseRepositoryMock) Create(imageInput m.ImageData) (m.ImageData, error) {
	switch imageInput.EntityType {
	case "createFail":
		return m.ImageData{}, errors.New("error")
	default:
		return imageDataDTO.ConvertFromDTO(), nil
	}
}

func (i databaseRepositoryMock) Update(imageInput m.ImageData) (m.ImageData, error) {
	switch imageInput.EntityType {
	case "updateFail":
		return m.ImageData{}, errors.New("error")
	default:
		return imageDataDTO.ConvertFromDTO(), nil
	}
}

func (i databaseRepositoryMock) Delete(imageInput m.ImageData) error {
	switch imageInput.EntityType {
	case "deleteFail":
		return errors.New("error")
	default:
		return nil
	}
}

func (r rabbitmqRepositoryMock) ImageUpdatedEvent(imageInput m.ImageData) error {
	switch imageInput.EntityType {
	case "rabbitmqFail":
		return errors.New("error")
	default:
		return nil
	}
}

func (LoggerInterfaceMock) Errorf(format string, args ...interface{}) {}

// ========================================================================================================

func TestFindAllImage_OK(t *testing.T) {
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "findall"
	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.IsType(t, []m.ImageDataDTO{}, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "findall", result[0].EntityType)
}

func TestFindAllImage_NotFoundErr(t *testing.T) {
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "notfound"
	result, err := s.FindAll()

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.IsType(t, []m.ImageDataDTO{}, result)
}

func TestFindAllImage_Err(t *testing.T) {
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "findallFail"
	result, err := s.FindAll()

	assert.Error(t, err)
	assert.IsType(t, []m.ImageDataDTO{}, result)
	assert.EqualError(t, err, "internal server error")

}

func TestFindImage_OK(t *testing.T) {
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "find"
	result, err := s.Find(imageDataDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.ImageDataDTO{}, result)
	assert.Equal(t, "find", result.EntityType)
}

func TestFindImage_NotFoundErr(t *testing.T) {
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "notfound"
	result, err := s.Find(imageDataDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.IsType(t, m.ImageDataDTO{}, result)
}

func TestFindImage_Err(t *testing.T) {
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "findFail"
	result, err := s.Find(imageDataDTO)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDataDTO{}, result)
	assert.EqualError(t, err, "internal server error")

}

func TestCreateImage_OK(t *testing.T) {
	imageFileDTO.File = tc.CreateFile()
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

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
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageFileDTO{
		EntityID:   imageDataDTO.EntityID,
		EntityType: "uploadS3Fail",
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
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageFileDTO{
		EntityID:   imageDataDTO.EntityID,
		EntityType: "createFail",
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
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

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
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageFileDTO{
		EntityID:   imageDataDTO.EntityID,
		EntityType: "findFail",
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
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageFileDTO{
		EntityID:   imageDataDTO.EntityID,
		EntityType: "uploadS3Fail",
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
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageFileDTO{
		EntityID:   imageDataDTO.EntityID,
		EntityType: "updateFail",
		Size:       imageDataDTO.Size,
		Type:       imageDataDTO.Type,
		File:       imageFileDTO.File,
	}
	result, err := s.Update(createImage)

	assert.Error(t, err)
	assert.IsType(t, m.ImageDataDTO{}, result)
}

func TestUpdateImage_RabbitmqErr(t *testing.T) {
	imageFileDTO.File = tc.CreateFile()
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	createImage := m.ImageFileDTO{
		EntityID:   imageDataDTO.EntityID,
		EntityType: "rabbitmqFail",
		Size:       imageDataDTO.Size,
		Type:       imageDataDTO.Type,
		File:       imageFileDTO.File,
	}
	result, err := s.Update(createImage)

	assert.NoError(t, err)
	assert.IsType(t, m.ImageDataDTO{}, result)
}

func TestDeleteImage_OK(t *testing.T) {
	// imageDataDTO.File = tc.CreateFile()
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "delete"
	err := s.Delete(imageDataDTO)

	assert.NoError(t, err)
}

func TestDeleteImage_FindErr(t *testing.T) {
	// imageDataDTO.File = tc.CreateFile()
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "findFail"
	err := s.Delete(imageDataDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "unable to find existing image. cannot delete something that does not exist")
}

func TestDeleteImage_DeleteS3Err(t *testing.T) {
	// imageDataDTO.File = tc.CreateFile()
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "deleteS3Fail"
	err := s.Delete(imageDataDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestDeleteImage_DeleteErr(t *testing.T) {
	// imageDataDTO.File = tc.CreateFile()
	s := NewImageService(&databaseRepositoryMock{}, &rabbitmqRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDataDTO.EntityType = "deleteFail"
	err := s.Delete(imageDataDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
