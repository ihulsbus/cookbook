package services

import (
	"errors"
	"testing"

	m "github.com/ihulsbus/cookbook/internal/models"
	"github.com/stretchr/testify/assert"
)

var (
	recipe m.Recipe = m.Recipe{
		RecipeName: "recipe",
	}
	instruction m.Instruction = m.Instruction{
		RecipeID:    1,
		StepNumber:  1,
		Description: "instruction",
	}
)

type RecipeRepositoryMock struct{}

func (RecipeRepositoryMock) FindAll() ([]m.Recipe, error) {
	var recipes []m.Recipe
	recipes = append(recipes, recipe)

	return recipes, nil
}

func (RecipeRepositoryMock) FindSingle(recipeID uint) (m.Recipe, error) {
	switch recipeID {
	case 1:
		return recipe, nil

	default:
		return recipe, errors.New("error")
	}
}

func (RecipeRepositoryMock) FindInstruction(recipeID uint) ([]m.Instruction, error) {
	switch recipeID {
	case 1:
		var instructions []m.Instruction
		instructions = append(instructions, instruction)

		return instructions, nil
	default:
		return nil, errors.New("error")
	}
}

func (RecipeRepositoryMock) CreateInstruction(instruction []m.Instruction) ([]m.Instruction, error) {
	switch instruction[0].RecipeID {
	case 1:
		return instruction, nil
	default:
		return nil, errors.New("error")
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

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, result[0].RecipeName, "recipe")
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

func TestRecipeFindInstruction_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	result, err := s.FindInstruction(1)

	assert.NoError(t, err)
	assert.IsType(t, result, []m.Instruction{})
	assert.Len(t, result, 1)
	assert.Equal(t, result[0].StepNumber, 1)
}

func TestRecipeFindInstruction_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	result, err := s.FindInstruction(0)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "error")

}

func TestRecipeCreateInstruction_OK(t *testing.T) {
	var createInstruction []m.Instruction
	s := NewRecipeService(&RecipeRepositoryMock{})

	createInstruction = append(createInstruction, instruction)

	result, err := s.CreateInstruction(createInstruction)

	assert.NoError(t, err)
	assert.IsType(t, result, []m.Instruction{})
	assert.Equal(t, result[0].Description, "instruction")
}

func TestRecipeCreateInstruction_Err(t *testing.T) {
	var createInstruction []m.Instruction
	s := NewRecipeService(&RecipeRepositoryMock{})

	createInstruction = append(createInstruction, instruction)
	createInstruction[0].RecipeID = 0

	result, err := s.CreateInstruction(createInstruction)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, []m.Instruction{})
	assert.Len(t, result, 0)
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

	result, err := s.Update(updateRecipe)

	assert.NoError(t, err)
	assert.IsType(t, result, m.Recipe{})
	assert.Equal(t, result.RecipeName, "recipe")
}

func TestRecipeUpdate_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	updateRecipe := recipe
	updateRecipe.ID = 2

	result, err := s.Update(updateRecipe)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, m.Recipe{})
}

func TestRecipeDelete_Ok(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	deleteRecipe := recipe
	deleteRecipe.ID = 1

	err := s.Delete(deleteRecipe)

	assert.NoError(t, err)
}

func TestRecipeDelete_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	deleteRecipe := recipe
	deleteRecipe.ID = 2

	err := s.Delete(deleteRecipe)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
