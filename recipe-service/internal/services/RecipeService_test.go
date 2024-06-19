package services

import (
	"errors"
	"testing"

	m "recipe-service/internal/models"

	"github.com/google/uuid"
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

func (RecipeRepositoryMock) FindAll() ([]m.Recipe, error) {
	switch findAllRecipe.Name {
	case "findall":
		var recipes []m.Recipe
		recipes = append(recipes, recipe)
		return recipes, nil
	case "notfound":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

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
	s := NewRecipeService(&RecipeRepositoryMock{})

	findAllRecipe.Name = "findall"

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.IsType(t, []m.RecipeDTO{}, result)
	assert.Len(t, result, 1)
}

func TestRecipeFindAll_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	findAllRecipe.Name = "error"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "internal server error")
}

func TestRecipeFindAll_NotFound(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	findAllRecipe.Name = "notfound"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
}

func TestRecipeFindSingle_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

	recipeDTO := m.RecipeDTO{
		Name: "",
	}
	result, err := s.Create(recipeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestRecipeCreate_NoDescription(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

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
	s := NewRecipeService(&RecipeRepositoryMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "delete",
	}
	err := s.Delete(recipeDTO)

	assert.NoError(t, err)
}

func TestRecipeDelete_DeleteErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "error",
	}
	err := s.Delete(recipeDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
