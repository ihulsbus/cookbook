package services

import (
	"errors"
	"testing"

	m "metadata-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	findAllCuisineType m.CuisineType = m.CuisineType{
		ID:   uuid.New(),
		Name: "cuisineType",
	}
	cuisineType m.CuisineType = m.CuisineType{
		ID:   uuid.New(),
		Name: "cuisineType",
	}
)

type CuisineTypeRepositoryMock struct{}

func (*CuisineTypeRepositoryMock) FindAll() ([]m.CuisineType, error) {
	switch findAllCuisineType.Name {
	case "findall":
		var cuisineTypes []m.CuisineType
		cuisineTypes = append(cuisineTypes, findAllCuisineType)
		return cuisineTypes, nil
	case "not found":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (*CuisineTypeRepositoryMock) FindSingle(cuisineType m.CuisineType) (m.CuisineType, error) {
	switch cuisineType.Name {
	case "find":
		return cuisineType, nil
	case "not found":
		return m.CuisineType{}, errors.New("not found")
	default:
		return m.CuisineType{}, errors.New("error")
	}
}

func (*CuisineTypeRepositoryMock) Create(cuisineType m.CuisineType) (m.CuisineType, error) {
	switch cuisineType.Name {
	case "create":
		return cuisineType, nil
	default:
		return m.CuisineType{}, errors.New("error")
	}
}

func (*CuisineTypeRepositoryMock) Update(cuisineType m.CuisineType) (m.CuisineType, error) {
	switch cuisineType.Name {
	case "update":
		return cuisineType, nil
	default:
		return m.CuisineType{}, errors.New("error")
	}
}

func (*CuisineTypeRepositoryMock) Delete(cuisineType m.CuisineType) error {
	switch cuisineType.Name {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

// ======================================================================

func TestCuisineTypeFindAll_OK(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})
	findAllCuisineType.Name = "findall"

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.IsType(t, []m.CuisineTypeDTO{}, result)
	assert.Len(t, result, 1)
}

func TestCuisineTypeFindAll_Err(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})
	findAllCuisineType.Name = "fail"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "internal server error")
}

func TestCuisineTypeFindAll_NotFound(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})
	findAllCuisineType.Name = "not found"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
}

func TestCuisineTypeFindSingle_OK(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})

	cuisineTypeDTO := m.CuisineTypeDTO{
		ID:   cuisineType.ID,
		Name: "find",
	}
	result, err := s.FindSingle(cuisineTypeDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.CuisineTypeDTO{}, result)
	assert.Equal(t, "find", result.Name)
	assert.Equal(t, cuisineType.ID, result.ID)
}

func TestCuisineTypeFindSingle_Err(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})

	cuisineTypeDTO := m.CuisineTypeDTO{
		ID:   cuisineType.ID,
		Name: "error",
	}
	result, err := s.FindSingle(cuisineTypeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CuisineTypeDTO{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestCuisineTypeFindSingle_NotFound(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})

	cuisineTypeDTO := m.CuisineTypeDTO{
		ID:   cuisineType.ID,
		Name: "not found",
	}
	result, err := s.FindSingle(cuisineTypeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CuisineTypeDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestCuisineTypeCreate_OK(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})

	cuisineTypeDTO := m.CuisineTypeDTO{
		Name: "create",
	}
	result, err := s.Create(cuisineTypeDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.CuisineTypeDTO{}, result)
}

func TestCuisineTypeCreate_IDErr(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})

	cuisineTypeDTO := m.CuisineTypeDTO{
		ID:   cuisineType.ID,
		Name: "create",
	}
	result, err := s.Create(cuisineTypeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CuisineTypeDTO{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestCuisineTypeCreate_NoNameErr(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})

	cuisineTypeDTO := m.CuisineTypeDTO{
		Name: "",
	}
	result, err := s.Create(cuisineTypeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CuisineTypeDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestCuisineTypeCreate_CreateErr(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})

	cuisineTypeDTO := m.CuisineTypeDTO{
		Name: "error",
	}
	result, err := s.Create(cuisineTypeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CuisineTypeDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestCuisineTypeUpdate_OK(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})

	cuisineTypeDTO := m.CuisineTypeDTO{
		ID:   cuisineType.ID,
		Name: "update",
	}
	result, err := s.Update(cuisineTypeDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.CuisineTypeDTO{}, result)
}

func TestCuisineTypeUpdate_NoNameErr(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})

	cuisineTypeDTO := m.CuisineTypeDTO{
		ID:   cuisineType.ID,
		Name: "",
	}
	result, err := s.Update(cuisineTypeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CuisineTypeDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestCuisineTypeUpdate_UpdateErr(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})

	cuisineTypeDTO := m.CuisineTypeDTO{
		ID:   cuisineType.ID,
		Name: "error",
	}
	result, err := s.Update(cuisineTypeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CuisineTypeDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestCuisineTypeDelete_OK(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})

	cuisineTypeDTO := m.CuisineTypeDTO{
		ID:   cuisineType.ID,
		Name: "delete",
	}
	err := s.Delete(cuisineTypeDTO)

	assert.NoError(t, err)
}

func TestCuisineTypeDelete_DeleteErr(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})

	cuisineTypeDTO := m.CuisineTypeDTO{
		ID:   cuisineType.ID,
		Name: "error",
	}
	err := s.Delete(cuisineTypeDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
