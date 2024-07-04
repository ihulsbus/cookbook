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
	imageDTO m.ImageDTO = m.ImageDTO{
		ID:         uuid.New(),
		EntityType: "",
		EntityID:   uuid.New(),
		Size:       0,
		Type:       "image/jpeg",
		File:       tc.CreateFile(),
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
	imageDTO.File = tc.CreateFile()
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
	imageDTO.File = tc.CreateFile()
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
	imageDTO.File = tc.CreateFile()
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
	imageDTO.File = tc.CreateFile()
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
	imageDTO.File = tc.CreateFile()
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
	imageDTO.File = tc.CreateFile()
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
	imageDTO.File = tc.CreateFile()
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
	imageDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "delete"
	err := s.Delete(imageDTO)

	assert.NoError(t, err)
}

func TestDeleteImage_FindErr(t *testing.T) {
	imageDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "findError"
	err := s.Delete(imageDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "unable to find existing image. cannot delete something that does not exist")
}

func TestDeleteImage_DeleteS3Err(t *testing.T) {
	imageDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "deleteS3Err"
	err := s.Delete(imageDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestDeleteImage_DeleteErr(t *testing.T) {
	imageDTO.File = tc.CreateFile()
	s := NewImageService(&imageRepositoryMock{}, &s3RepositoryMock{}, &LoggerInterfaceMock{})

	imageDTO.EntityType = "deleteErr"
	err := s.Delete(imageDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
