package services

import (
	"errors"
	"testing"

	m "instruction-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	instruction m.Instruction = m.Instruction{
		Description: "instruction",
	}
)

type InstructionRepositoryMock struct{}

func (InstructionRepositoryMock) FindInstruction(recipeID uuid.UUID) (m.Instruction, error) {
	switch recipeID {
	case uuid.UUID{}:
		return instruction, nil
	case uuid.UUID{}:
		ni := instruction
		ni.RecipeID = uuid.UUID{}
		return ni, nil
	default:
		return instruction, errors.New("error")
	}
}

func (InstructionRepositoryMock) CreateInstruction(instruction m.Instruction) (m.Instruction, error) {
	switch instruction.RecipeID {
	case uuid.UUID{}:
		return instruction, nil
	default:
		return instruction, errors.New("error")
	}
}

func (InstructionRepositoryMock) UpdateInstruction(instruction m.Instruction) (m.Instruction, error) {
	switch instruction.RecipeID {
	case uuid.UUID{}:
		return instruction, nil
	default:
		return instruction, errors.New("error")
	}
}

func (InstructionRepositoryMock) DeleteInstruction(instruction m.Instruction) error {
	switch instruction.RecipeID {
	case uuid.UUID{}:
		return nil
	default:
		return errors.New("error")
	}
}

// ========================================================================================================

func TestFindInstruction_OK(t *testing.T) {
	s := NewRecipeService(&InstructionRepositoryMock{})

	result, err := s.FindInstruction(uuid.UUID{})

	assert.NoError(t, err)
	assert.IsType(t, m.Instruction{}, result)
}

func TestFindInstruction_Err(t *testing.T) {
	s := NewRecipeService(&InstructionRepositoryMock{})

	result, err := s.FindInstruction(uuid.UUID{})

	assert.Error(t, err)
	assert.IsType(t, m.Instruction{}, result)
	assert.EqualError(t, err, "error")

}

func TestCreateInstruction_OK(t *testing.T) {
	s := NewRecipeService(&InstructionRepositoryMock{})

	result, err := s.CreateInstruction(instruction, uuid.UUID{})

	assert.NoError(t, err)
	assert.IsType(t, m.Instruction{}, result)
	assert.Equal(t, result.Description, "instruction")
}

func TestCreateInstruction_Err(t *testing.T) {
	s := NewRecipeService(&InstructionRepositoryMock{})

	createInstruction := instruction
	createInstruction.RecipeID = uuid.UUID{}

	result, err := s.CreateInstruction(createInstruction, uuid.UUID{})

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, m.Instruction{}, result)
}

func TestUpdateInstruction_OK(t *testing.T) {
	s := NewRecipeService(&InstructionRepositoryMock{})

	result, err := s.UpdateInstruction(instruction, uuid.UUID{})

	assert.NoError(t, err)
	assert.IsType(t, m.Instruction{}, result)
	assert.Equal(t, instruction.RecipeID, result.RecipeID)
	assert.Equal(t, instruction.Description, result.Description)
}

func TestUpdateInstruction_FindErr(t *testing.T) {
	s := NewRecipeService(&InstructionRepositoryMock{})

	result, err := s.UpdateInstruction(instruction, uuid.UUID{})

	assert.Error(t, err)
	assert.EqualError(t, err, "unable to find existing instruction for the given recipe id")
	assert.Equal(t, m.Instruction{}, result)
}

func TestUpdateInstruction_UpdateErr(t *testing.T) {
	s := NewRecipeService(&InstructionRepositoryMock{})

	result, err := s.UpdateInstruction(instruction, uuid.UUID{})

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.Equal(t, m.Instruction{}, result)
}

func TestDeleteInstruction_OK(t *testing.T) {
	s := NewRecipeService(&InstructionRepositoryMock{})

	err := s.DeleteInstruction(uuid.UUID{})

	assert.NoError(t, err)
}

func TestDeleteInstruction_FindErr(t *testing.T) {
	s := NewRecipeService(&InstructionRepositoryMock{})

	err := s.DeleteInstruction(uuid.UUID{})

	assert.Error(t, err)
	assert.EqualError(t, err, "unable to find existing instruction for the given recipe id")
}

func TestDeleteInstruction_DeleteErr(t *testing.T) {
	s := NewRecipeService(&InstructionRepositoryMock{})

	err := s.DeleteInstruction(uuid.UUID{})

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
