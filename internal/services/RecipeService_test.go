package services

import (
	"errors"
	"testing"

	m "github.com/ihulsbus/cookbook/internal/models"
	"github.com/stretchr/testify/assert"
)

var (
	findAllRecipe          = recipe
	recipe        m.Recipe = m.Recipe{
		RecipeName: "recipe",
	}
	instruction m.Instruction = m.Instruction{
		RecipeID:    1,
		Description: "instruction",
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

func (RecipeRepositoryMock) FindInstruction(recipeID uint) (m.Instruction, error) {
	switch recipeID {
	case 1:
		return instruction, nil
	case 3:
		ni := instruction
		ni.RecipeID = 3
		return ni, nil
	default:
		return instruction, errors.New("error")
	}
}

func (RecipeRepositoryMock) CreateInstruction(instruction m.Instruction) (m.Instruction, error) {
	switch instruction.RecipeID {
	case 1:
		return instruction, nil
	default:
		return instruction, errors.New("error")
	}
}

func (RecipeRepositoryMock) UpdateInstruction(instruction m.Instruction) (m.Instruction, error) {
	switch instruction.RecipeID {
	case 1:
		return instruction, nil
	default:
		return instruction, errors.New("error")
	}
}

func (RecipeRepositoryMock) DeleteInstruction(instruction m.Instruction) error {
	switch instruction.RecipeID {
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

func TestFindInstruction_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	result, err := s.FindInstruction(1)

	assert.NoError(t, err)
	assert.IsType(t, m.Instruction{}, result)
}

func TestFindInstruction_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	result, err := s.FindInstruction(0)

	assert.Error(t, err)
	assert.IsType(t, m.Instruction{}, result)
	assert.EqualError(t, err, "error")

}

func TestCreateInstruction_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	result, err := s.CreateInstruction(instruction)

	assert.NoError(t, err)
	assert.IsType(t, m.Instruction{}, result)
	assert.Equal(t, result.Description, "instruction")
}

func TestCreateInstruction_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	createInstruction := instruction
	createInstruction.RecipeID = 0

	result, err := s.CreateInstruction(createInstruction)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, m.Instruction{}, result)
}

func TestUpdateInstruction_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	result, err := s.UpdateInstruction(instruction, uint(1))

	assert.NoError(t, err)
	assert.IsType(t, m.Instruction{}, result)
	assert.Equal(t, instruction.RecipeID, result.RecipeID)
	assert.Equal(t, instruction.Description, result.Description)
}

func TestUpdateInstruction_FindErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	result, err := s.UpdateInstruction(instruction, uint(2))

	assert.Error(t, err)
	assert.EqualError(t, err, "unable to find existing instruction for the given recipe id")
	assert.Equal(t, m.Instruction{}, result)
}

func TestUpdateInstruction_UpdateErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	result, err := s.UpdateInstruction(instruction, uint(3))

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.Equal(t, m.Instruction{}, result)
}

func TestDeleteInstruction_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	err := s.DeleteInstruction(uint(1))

	assert.NoError(t, err)
}

func TestDeleteInstruction_FindErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	err := s.DeleteInstruction(uint(2))

	assert.Error(t, err)
	assert.EqualError(t, err, "unable to find existing instruction for the given recipe id")
}

func TestDeleteInstruction_DeleteErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{})

	err := s.DeleteInstruction(uint(3))

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
