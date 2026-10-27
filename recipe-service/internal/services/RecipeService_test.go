package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
)

var (
	findAllRecipe m.Recipe = m.Recipe{
		ID:           uuid.New(),
		Name:         "recipe",
		Description:  "description",
		ServingCount: 1,
	}
	recipe m.Recipe = m.Recipe{
		ID:           uuid.New(),
		Name:         "recipe",
		Description:  "description",
		ServingCount: 1,
	}
)

type RecipeRepositoryMock struct{}
type RabbitMQRepositoryMock struct{}
type LoggerMock struct{}

func (RecipeRepositoryMock) FindAll(pagination m.PaginationRequest) ([]m.Recipe, int64, error) {
	switch findAllRecipe.Name {
	case "findall":
		return []m.Recipe{recipe}, 1, nil
	case "notfound":
		return nil, 0, errors.New("not found")
	default:
		return nil, 0, errors.New("error")
	}
}

func (RabbitMQRepositoryMock) RecipeCreatedEvent(r m.Recipe) error   { return nil }
func (RabbitMQRepositoryMock) RecipeUpdatedEvent(r m.Recipe) error   { return nil }
func (RabbitMQRepositoryMock) RecipeDeletedEvent(id uuid.UUID) error { return nil }

func (LoggerMock) Debugf(format string, args ...interface{}) {}
func (LoggerMock) Infof(format string, args ...interface{})  {}
func (LoggerMock) Warnf(format string, args ...interface{})  {}
func (LoggerMock) Errorf(format string, args ...interface{}) {}

func (RecipeRepositoryMock) FindSingle(recipeInput m.Recipe) (m.Recipe, error) {
	switch recipeInput.Name {
	case "find":
		return recipe, nil
	case "update":
		return recipe, nil
	case "updateerror":
		return recipe, nil
	case "":
		return recipe, nil
	case "notfound":
		return m.Recipe{}, errors.New("not found")
	default:
		return m.Recipe{}, errors.New("error")
	}
}

func (RecipeRepositoryMock) Create(recipeInput m.Recipe) (m.Recipe, error) {
	switch recipeInput.Name {
	case "create":
		return recipe, nil
	default:
		return m.Recipe{}, errors.New("error")
	}
}

func (RecipeRepositoryMock) Update(recipeInput m.Recipe) (m.Recipe, error) {
	switch recipeInput.Name {
	case "update":
		return recipe, nil
	case "recipe":
		return recipe, nil
	default:
		return m.Recipe{}, errors.New("error")
	}
}

func (RecipeRepositoryMock) Delete(recipeInput m.Recipe) error {
	switch recipeInput.Name {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

func TestRecipeFindAll_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	findAllRecipe.Name = "findall"

	result, err := s.FindAll(m.NormalizePagination(1, 25))

	assert.NoError(t, err)
	assert.IsType(t, m.PaginatedResponse[m.RecipeDTO]{}, result)
	assert.Len(t, result.Data, 1)
}

func TestRecipeFindAll_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	findAllRecipe.Name = "error"

	result, err := s.FindAll(m.NormalizePagination(1, 25))

	assert.Error(t, err)
	assert.IsType(t, m.PaginatedResponse[m.RecipeDTO]{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestRecipeFindAll_NotFound(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	findAllRecipe.Name = "notfound"

	result, err := s.FindAll(m.NormalizePagination(1, 25))

	assert.Error(t, err)
	assert.IsType(t, m.PaginatedResponse[m.RecipeDTO]{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestRecipeFindSingle_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "find",
	}
	result, err := s.FindSingle(recipeDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.Equal(t, "recipe", result.Name)
	assert.Equal(t, recipe.ID, result.ID)
}

func TestRecipeFindSingle_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "error",
	}
	result, err := s.FindSingle(recipeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestRecipeFindSingle_NotFound(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "notfound",
	}
	result, err := s.FindSingle(recipeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestRecipeCreate_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		Name:         "create",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Create(recipeDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.Equal(t, "recipe", result.Name)
	assert.Equal(t, recipe.ID, result.ID)
}

func TestRecipeCreate_IDErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "create",
	}
	result, err := s.Create(recipeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestRecipeCreate_NoName(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		Name: "",
	}
	result, err := s.Create(recipeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestRecipeCreate_NoDescription(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		Name:        recipe.Name,
		Description: "",
	}
	result, err := s.Create(recipeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "description is empty")
}

func TestRecipeCreate_NoServingCount(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		Name:         recipe.Name,
		Description:  recipe.Description,
		ServingCount: 0,
	}
	result, err := s.Create(recipeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "serving count 0 is not allowed")
}

func TestRecipeCreate_CreateErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		Name:         "error",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Create(recipeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestRecipeUpdate_Ok(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "update",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Update(recipeDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.Equal(t, recipe.ID, result.ID)
	assert.Equal(t, recipe.Name, result.Name)
	assert.Equal(t, recipe.Description, result.Description)
	assert.Equal(t, recipe.ServingCount, result.ServingCount)
}

func TestRecipeUpdate_NoNameErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Update(recipeDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
}

func TestRecipeUpdate_NoDescription(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "update",
		Description:  "",
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Update(recipeDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
}

func TestRecipeUpdate_NoServingCount(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "update",
		Description:  recipe.Description,
		ServingCount: 0,
	}
	result, err := s.Update(recipeDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
}

func TestRecipeUpdate_FindErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "notfound",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Update(recipeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestRecipeUpdate_UpdateErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "updateerror",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Update(recipeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestRecipeDelete_Ok(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "delete",
	}
	err := s.Delete(recipeDTO)

	assert.NoError(t, err)
}

func TestRecipeDelete_DeleteErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &RabbitMQRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "error",
	}
	err := s.Delete(recipeDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
