package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
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

func (*CuisineTypeRepositoryMock) FindAll(pagination m.PaginationRequest) ([]m.CuisineType, int64, error) {
	switch findAllCuisineType.Name {
	case "findall":
		var cuisineTypes []m.CuisineType
		cuisineTypes = append(cuisineTypes, findAllCuisineType)
		return cuisineTypes, 1, nil
	default:
		return nil, 0, errors.New("error")
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

	result, err := s.FindAll(m.NormalizePagination(1, 25))

	assert.NoError(t, err)
	assert.IsType(t, m.PaginatedResponse[m.CuisineTypeDTO]{}, result)
	assert.Len(t, result.Data, 1)
}

func TestCuisineTypeFindAll_Err(t *testing.T) {
	s := NewCuisineTypeService(&CuisineTypeRepositoryMock{})
	findAllCuisineType.Name = "fail"

	result, err := s.FindAll(m.NormalizePagination(1, 25))

	assert.Error(t, err)
	assert.Equal(t, m.PaginatedResponse[m.CuisineTypeDTO]{}, result)
	assert.EqualError(t, err, "internal server error")
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
