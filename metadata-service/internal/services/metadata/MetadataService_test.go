package services

import (
	"errors"
	"testing"

	m "metadata-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	findAllRecipeMetadata m.RecipeMetadata = m.RecipeMetadata{
		RecipeID:      uuid.New(),
		CuisineTypeID: "RecipeMetadata",
	}
	RecipeMetadata m.RecipeMetadata = m.RecipeMetadata{
		ID:   uuid.New(),
		Name: "RecipeMetadata",
	}
)

type RecipeMetadataRepositoryMock struct{}

func (*RecipeMetadataRepositoryMock) FindAll() (*[]m.RecipeMetadata, error) {
	switch findAllRecipeMetadata.Name {
	case "findall":
		var categories []m.RecipeMetadata
		categories = append(categories, findAllRecipeMetadata)
		return categories, nil
	case "not found":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (*RecipeMetadataRepositoryMock) FindSingle(recipeID uuid.UUID) (*m.RecipeMetadata, error) {
	switch RecipeMetadata.Name {
	case "find":
		return RecipeMetadata, nil
	case "not found":
		return m.RecipeMetadata{}, errors.New("not found")
	default:
		return m.RecipeMetadata{}, errors.New("error")
	}
}

func (*RecipeMetadataRepositoryMock) Create(meta *m.RecipeMetadata) (*m.RecipeMetadata, error) {
	RecipeMetadataC := RecipeMetadata
	switch RecipeMetadata.Name {
	case "create":
		RecipeMetadataC.Name = "create"
		return RecipeMetadataC, nil
	default:
		return m.RecipeMetadata{}, errors.New("error")
	}
}

func (*RecipeMetadataRepositoryMock) Update(meta *m.RecipeMetadata) (*m.RecipeMetadata, error) {
	switch RecipeMetadata.Name {
	case "update":
		return RecipeMetadata, nil
	default:
		return m.RecipeMetadata{}, errors.New("error")
	}
}

func (*RecipeMetadataRepositoryMock) Delete(meta *m.RecipeMetadata) erro {
	switch RecipeMetadata.Name {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

// ======================================================================

func TestRecipeMetadataFindAll_OK(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})
	findAllRecipeMetadata.Name = "findall"

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.IsType(t, []m.RecipeMetadataDTO{}, result)
	assert.Len(t, result, 1)
}

func TestRecipeMetadataFindAll_err(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})
	findAllRecipeMetadata.Name = "fail"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "internal server error")
}

func TestRecipeMetadataFindAll_NotFound(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})
	findAllRecipeMetadata.Name = "not found"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
}

func TestRecipeMetadataFindSingle_OK(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})

	RecipeMetadataDTO := m.RecipeMetadataDTO{
		ID:   RecipeMetadata.ID,
		Name: "find",
	}
	result, err := s.FindSingle(RecipeMetadataDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeMetadataDTO{}, result)
	assert.Equal(t, "find", result.Name)
	assert.Equal(t, result.ID, RecipeMetadata.ID)
}

func TestRecipeMetadataFindSingle_Err(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})

	RecipeMetadataDTO := m.RecipeMetadataDTO{
		ID:   RecipeMetadata.ID,
		Name: RecipeMetadata.Name,
	}
	result, err := s.FindSingle(RecipeMetadataDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeMetadataDTO{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestRecipeMetadataFindSingle_NotFound(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})

	RecipeMetadataDTO := m.RecipeMetadataDTO{
		ID:   RecipeMetadata.ID,
		Name: "not found",
	}
	result, err := s.FindSingle(RecipeMetadataDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeMetadataDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestRecipeMetadataCreate_OK(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})

	RecipeMetadataDTO := m.RecipeMetadataDTO{
		Name: "create",
	}
	result, err := s.Create(RecipeMetadataDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeMetadataDTO{}, result)
}

func TestRecipeMetadataCreate_IDErr(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})

	RecipeMetadataDTO := m.RecipeMetadataDTO{
		ID:   RecipeMetadata.ID,
		Name: RecipeMetadata.Name,
	}
	result, err := s.Create(RecipeMetadataDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeMetadataDTO{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestRecipeMetadataCreate_NoNameErr(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})

	RecipeMetadataDTO := m.RecipeMetadataDTO{
		Name: "",
	}
	result, err := s.Create(RecipeMetadataDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeMetadataDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestRecipeMetadataCreate_CreateErr(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})

	RecipeMetadataDTO := m.RecipeMetadataDTO{
		Name: RecipeMetadata.Name,
	}
	result, err := s.Create(RecipeMetadataDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeMetadataDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestRecipeMetadataUpdate_OK(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})

	RecipeMetadataDTO := m.RecipeMetadataDTO{
		ID:   RecipeMetadata.ID,
		Name: "update",
	}
	result, err := s.Update(RecipeMetadataDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeMetadataDTO{}, result)
}

func TestRecipeMetadataUpdate_NoNameErr(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})

	RecipeMetadataDTO := m.RecipeMetadataDTO{
		ID: RecipeMetadata.ID,
	}
	result, err := s.Update(RecipeMetadataDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeMetadataDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestRecipeMetadataUpdate_UpdateErr(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})

	RecipeMetadataDTO := m.RecipeMetadataDTO{
		ID:   RecipeMetadata.ID,
		Name: RecipeMetadata.Name,
	}
	result, err := s.Update(RecipeMetadataDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeMetadataDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestRecipeMetadataDelete_OK(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})

	RecipeMetadataDTO := m.RecipeMetadataDTO{
		ID:   RecipeMetadata.ID,
		Name: "delete",
	}
	err := s.Delete(RecipeMetadataDTO)

	assert.NoError(t, err)
}

func TestRecipeMetadataDelete_DeleteErr(t *testing.T) {
	s := NewRecipeMetadataService(&RecipeMetadataRepositoryMock{})

	RecipeMetadataDTO := m.RecipeMetadataDTO{
		ID:   RecipeMetadata.ID,
		Name: RecipeMetadata.Name,
	}
	err := s.Delete(RecipeMetadataDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
