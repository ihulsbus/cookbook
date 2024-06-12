package services

import (
	"errors"
	"testing"

	m "recipe-service/internal/models"

	"github.com/stretchr/testify/assert"
)

var (
	findAllRecipe          = recipe
	recipe        m.Recipe = m.Recipe{
		RecipeName: "recipe",
	}
)

type RecipeRepositoryMock struct{}

func (RecipeRepositoryMock) FindAll() ([]m.Recipe, error) {
	switch findAllRecipe.ID {
	case 1:
		var recipes []m.Recipe
		recipes = append(recipes, recipe)

		return recipes, nil
	default:
		return nil, errors.New("error")
	}
}

func (RecipeRepositoryMock) FindSingle(recipeID uint) (m.Recipe, error) {
	switch recipeID {
	case 1:
		return recipe, nil
	case 3:
		return recipe, nil
	default:
		return recipe, errors.New("error")
	}
}

func (RecipeRepositoryMock) Create(recipe m.Recipe) (m.Recipe, error) {
	switch recipe.ID {
	case 1:
		return recipe, nil
	default:
		return recipe, errors.New("error")
	}
}

func (RecipeRepositoryMock) Update(recipe m.Recipe) (m.Recipe, error) {
	switch recipe.ID {
	case 1:
		return recipe, nil
	case 2:
		return recipe, nil
	default:
		return recipe, errors.New("error")
	}
}

func (RecipeRepositoryMock) Delete(recipe m.Recipe) error {
	switch recipe.ID {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

func TestRecipeFindAll_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	findAllRecipe.ID = 1

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, result[0].RecipeName, "recipe")
}

func TestRecipeFindAll_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	findAllRecipe.ID = 2

	result, err := s.FindAll()

	findAllRecipe.ID = 1

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "error")
}

func TestRecipeFindSingle_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	result, err := s.FindSingle(1)

	assert.NoError(t, err)
	assert.IsType(t, result, m.Recipe{})
	assert.Equal(t, result.RecipeName, "recipe")
}

func TestRecipeFindSingle_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	result, err := s.FindSingle(2)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, m.Recipe{})
}

func TestRecipeCreate_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	createRecipe := recipe
	createRecipe.ID = 1

	result, err := s.Create(createRecipe)

	assert.NoError(t, err)
	assert.IsType(t, result, m.Recipe{})
	assert.Equal(t, result.RecipeName, "recipe")
}

func TestRecipeCreate_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	createRecipe := recipe
	createRecipe.ID = 2

	result, err := s.Create(createRecipe)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, m.Recipe{})
}

func TestRecipeUpdate_Ok(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	updateRecipe := recipe
	updateRecipe.ID = 1
	updateRecipe.RecipeName = ""

	result, err := s.Update(updateRecipe, uint(1))

	assert.NoError(t, err)
	assert.IsType(t, result, m.Recipe{})
	assert.Equal(t, result.RecipeName, "recipe")
}

func TestRecipeUpdate_FindErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	updateRecipe := recipe
	updateRecipe.ID = 2

	result, err := s.Update(updateRecipe, uint(2))

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, m.Recipe{})
}

func TestRecipeUpdate_UpdateErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	updateRecipe := recipe
	updateRecipe.ID = 3

	result, err := s.Update(updateRecipe, uint(3))

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, m.Recipe{})
}

func TestRecipeDelete_Ok(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	deleteRecipe := recipe
	deleteRecipe.ID = 1

	err := s.Delete(uint(1))

	assert.NoError(t, err)
}

func TestRecipeDelete_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	deleteRecipe := recipe
	deleteRecipe.ID = 2

	err := s.Delete(uint(2))

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}