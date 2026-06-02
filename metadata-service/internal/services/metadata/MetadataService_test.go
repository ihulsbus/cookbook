package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
)

var (
	existingRecipe                  = "true"
	recipeMetadata m.RecipeMetadata = m.RecipeMetadata{
		RecipeID:          uuid.New(),
		CuisineTypeID:     uuid.New(),
		DifficultyLevelID: uuid.New(),
		PreparationTime:   99,
		ServingCount:      99,
	}
)

type RecipeMetadataRepositoryMock struct{}
type RecipeRepositoryMock struct{}

func (*RecipeMetadataRepositoryMock) FindAll(pagination m.PaginationRequest) ([]m.RecipeMetadata, int64, error) {
	switch recipeMetadata.ServingCount {
	case 1: // find all
		var findall []m.RecipeMetadata
		findall = append(findall, recipeMetadata)
		return findall, 1, nil
	default:
		return nil, 0, errors.New("error")
	}
}

func (*RecipeMetadataRepositoryMock) FindSingle(recipeID uuid.UUID) (*m.RecipeMetadata, error) {
	switch recipeMetadata.ServingCount {
	case 2: // find
		return &recipeMetadata, nil
	case 4: // update
		return &recipeMetadata, nil
	case 499: // update error
		return &recipeMetadata, nil
	case 5: // delete
		return &recipeMetadata, nil
	case 599: // delete error
		return &recipeMetadata, nil
	case 0: // not found
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (*RecipeMetadataRepositoryMock) Create(meta *m.RecipeMetadata) (*m.RecipeMetadata, error) {
	switch recipeMetadata.ServingCount {
	case 3: // create
		return &recipeMetadata, nil
	default:
		return nil, errors.New("error")
	}
}

func (*RecipeMetadataRepositoryMock) Update(meta *m.RecipeMetadata) (*m.RecipeMetadata, error) {
	switch recipeMetadata.ServingCount {
	case 4: // update
		return &recipeMetadata, nil
	default:
		return nil, errors.New("error")
	}
}

func (*RecipeMetadataRepositoryMock) Delete(meta *m.RecipeMetadata) error {
	switch recipeMetadata.ServingCount {
	case 5: // delete
		return nil
	default:
		return errors.New("error")
	}
}

func (*RecipeRepositoryMock) RecipeExists(recipeID string) (bool, error) {
	switch existingRecipe {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, errors.New("error")
	}
}

// ======================================================================

func TestRecipeMetadataFindAll_OK(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})
	recipeMetadata.ServingCount = 1

	result, err := s.FindAll(m.NormalizePagination(1, 25))

	assert.NoError(t, err)
	assert.IsType(t, m.PaginatedResponse[m.RecipeMetadataDTO]{}, result)
	assert.Len(t, result.Data, 1)
}

func TestRecipeMetadataFindAll_err(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})
	recipeMetadata.ServingCount = 99

	result, err := s.FindAll(m.NormalizePagination(1, 25))

	assert.Error(t, err)
	assert.Equal(t, m.PaginatedResponse[m.RecipeMetadataDTO]{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestRecipeMetadataFindSingle_OK(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 2
	result, err := s.Find(recipeMetadata.RecipeID)

	assert.NoError(t, err)
	assert.IsType(t, &m.RecipeMetadataDTO{}, result)
	assert.Equal(t, 2, result.ServingCount)
	assert.Equal(t, result.RecipeID, recipeMetadata.RecipeID)
}

func TestRecipeMetadataFindSingle_Err(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 99
	result, err := s.Find(recipeMetadata.RecipeID)

	assert.Error(t, err)
	assert.IsType(t, &m.RecipeMetadataDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestRecipeMetadataFindSingle_NotFound(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 0
	result, err := s.Find(recipeMetadata.RecipeID)

	assert.Error(t, err)
	assert.IsType(t, &m.RecipeMetadataDTO{}, result)
	assert.EqualError(t, err, "recipe not found")
}

func TestRecipeMetadataCreate_OK(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 3
	dto := recipeMetadata.ConvertToDTO()
	result, err := s.Create(recipeMetadata.RecipeID, &dto)

	assert.NoError(t, err)
	assert.IsType(t, &m.RecipeMetadataDTO{}, result)
}

func TestRecipeMetadataCreate_IDErr(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 3
	dto := recipeMetadata.ConvertToDTO()
	result, err := s.Create(uuid.Nil, &dto)

	assert.Error(t, err)
	assert.IsType(t, &m.RecipeMetadataDTO{}, result)
	assert.EqualError(t, err, "recipe id should not be nil")
}

func TestRecipeMetadataCreate_RecipeDoesNotExist(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 3
	existingRecipe = "false"
	dto := recipeMetadata.ConvertToDTO()
	result, err := s.Create(recipeMetadata.RecipeID, &dto)

	assert.Error(t, err)
	assert.IsType(t, &m.RecipeMetadataDTO{}, result)
	assert.Nil(t, result)
	assert.EqualError(t, err, "provided recipe does not exist. Cannot create metadata for a recipe that does not exist")
}

func TestRecipeMetadataCreate_RecipeExistsError(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 3
	existingRecipe = "error"
	dto := recipeMetadata.ConvertToDTO()
	result, err := s.Create(uuid.Nil, &dto)

	assert.Error(t, err)
	assert.IsType(t, &m.RecipeMetadataDTO{}, result)
	assert.Nil(t, result)
	assert.EqualError(t, err, "recipe id should not be nil")
}

func TestRecipeMetadataCreate_CreateErr(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 99
	dto := recipeMetadata.ConvertToDTO()
	result, err := s.Create(dto.RecipeID, &dto)

	assert.Error(t, err)
	assert.IsType(t, &m.RecipeMetadataDTO{}, result)
	assert.Nil(t, result)
	assert.EqualError(t, err, "error")
}

func TestRecipeMetadataUpdate_OK(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 4
	existingRecipe = "true"
	dto := recipeMetadata.ConvertToDTO()
	result, err := s.Update(dto.RecipeID, &dto)

	assert.NoError(t, err)
	assert.IsType(t, &m.RecipeMetadataDTO{}, result)
}

func TestRecipeMetadataUpdate_NoExistingRecipe(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 99
	dto := recipeMetadata.ConvertToDTO()
	result, err := s.Update(dto.RecipeID, &dto)

	assert.Error(t, err)
	assert.IsType(t, &m.RecipeMetadataDTO{}, result)
	assert.Nil(t, result)
	assert.EqualError(t, err, "provided recipe does not exist")
}

func TestRecipeMetadataUpdate_UpdateErr(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 499
	dto := recipeMetadata.ConvertToDTO()
	result, err := s.Update(dto.RecipeID, &dto)

	assert.Error(t, err)
	assert.IsType(t, &m.RecipeMetadataDTO{}, result)
	assert.Nil(t, result)
	assert.EqualError(t, err, "error")
}

func TestRecipeMetadataDelete_OK(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 5
	err := s.Delete(recipeMetadata.RecipeID)

	assert.NoError(t, err)
}

func TestRecipeMetadataDelete_FindError(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 299
	err := s.Delete(recipeMetadata.RecipeID)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestRecipeMetadataDelete_DeleteErr(t *testing.T) {
	s := NewMetadataService(&RecipeMetadataRepositoryMock{}, &RecipeRepositoryMock{})

	recipeMetadata.ServingCount = 99
	err := s.Delete(recipeMetadata.RecipeID)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
