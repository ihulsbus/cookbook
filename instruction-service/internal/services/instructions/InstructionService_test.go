package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
)

var (
	instruction = m.Instruction{
		ID:          uuid.New(),
		Sequence:    1,
		Description: "instruction",
		MediaID:     uuid.New(),
	}
	instructionArr = []m.Instruction{instruction}
)

type InstructionRepositoryMock struct{}

type RecipeClientMock struct{}
type ImageClientMock struct{}

func (InstructionRepositoryMock) Find(entityID uuid.UUID) (*[]m.Instruction, error) {
	switch instruction.Description {
	case "find":
		return &instructionArr, nil
	case "create":
		return nil, errors.New("not found")
	case "update":
		return &instructionArr, nil
	case "updateError":
		return &instructionArr, nil
	case "delete":
		return &instructionArr, nil
	case "deleteError":
		return &instructionArr, nil
	case "notfound":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (InstructionRepositoryMock) Create(instructionInput *[]m.Instruction) (*[]m.Instruction, error) {
	switch instruction.Description {
	case "create":
		return instructionInput, nil
	case "update":
		return instructionInput, nil
	default:
		return nil, errors.New("error")
	}
}

func (InstructionRepositoryMock) Delete(_ *[]m.Instruction) error {
	switch instruction.Description {
	case "delete":
		return nil
	case "update":
		return nil
	default:
		return errors.New("error")
	}
}

func (RecipeClientMock) RecipeExists(recipeID string) (bool, error) {
	return true, nil
}

// ========================================================================================================

func TestFindInstruction_OK(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{}, RecipeClientMock{})

	instruction.Description = "find"
	result, err := s.Find(instruction.EntityID)

	assert.NoError(t, err)
	assert.IsType(t, &[]m.InstructionDTO{}, result)
}

func TestFindInstruction_NotFoundErr(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{}, RecipeClientMock{})

	instruction.Description = "notfound"
	result, err := s.Find(instruction.EntityID)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.IsType(t, &[]m.InstructionDTO{}, result)
}

func TestFindInstruction_Err(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{}, RecipeClientMock{})

	instruction.Description = "error"
	result, err := s.Find(instruction.EntityID)

	assert.Error(t, err)
	assert.IsType(t, &[]m.InstructionDTO{}, result)
	assert.EqualError(t, err, "internal server error")

}

func TestCreateInstruction_OK(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{}, RecipeClientMock{})

	instruction.Description = "create"

	var instructionDTOArr []m.InstructionDTO
	instructionDTOArr = append(instructionDTOArr, instruction.ConvertToDTO())
	result, err := s.Create(instruction.EntityID, &instructionDTOArr)

	assert.NoError(t, err)
	assert.IsType(t, &[]m.InstructionDTO{}, result)
}

func TestCreateInstruction_Err(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{}, RecipeClientMock{})

	instruction.Description = "error"

	var instructionDTOArr []m.InstructionDTO
	instructionDTOArr = append(instructionDTOArr, instruction.ConvertToDTO())
	result, err := s.Create(instruction.EntityID, &instructionDTOArr)

	assert.Error(t, err)
	assert.IsType(t, &[]m.InstructionDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestUpdateInstruction_OK(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{}, RecipeClientMock{})

	instruction.Description = "update"

	var instructionDTOArr []m.InstructionDTO
	instructionDTOArr = append(instructionDTOArr, instruction.ConvertToDTO())
	result, err := s.Update(instruction.EntityID, &instructionDTOArr)

	assert.NoError(t, err)
	assert.IsType(t, &[]m.InstructionDTO{}, result)
}

func TestUpdateInstruction_FindErr(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{}, RecipeClientMock{})

	instruction.Description = "error"

	var instructionDTOArr []m.InstructionDTO
	instructionDTOArr = append(instructionDTOArr, instruction.ConvertToDTO())
	result, err := s.Update(instruction.EntityID, &instructionDTOArr)

	assert.Error(t, err)
	assert.IsType(t, &[]m.InstructionDTO{}, result)
	assert.EqualError(t, err, "unable to find existing recipe. cannot update something that does not exist")
}

func TestUpdateInstruction_UpdateErr(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{}, RecipeClientMock{})

	instruction.Description = "find"

	var instructionDTOArr []m.InstructionDTO
	instructionDTOArr = append(instructionDTOArr, instruction.ConvertToDTO())
	_, err := s.Update(instruction.EntityID, &instructionDTOArr)

	assert.Error(t, err)
	assert.EqualError(t, err, "an error occured deleting existing instructions error")
}

func TestDeleteInstruction_OK(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{}, RecipeClientMock{})

	instruction.Description = "delete"

	err := s.Delete(instruction.EntityID)

	assert.NoError(t, err)
}

func TestDeleteInstruction_FindErr(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{}, RecipeClientMock{})

	instruction.Description = "error"

	err := s.Delete(instruction.EntityID)

	assert.Error(t, err)
	assert.EqualError(t, err, "unable to find existing instruction. cannot delete something that does not exist")
}

func TestDeleteInstruction_DeleteErr(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{}, RecipeClientMock{})

	instruction.Description = "deleteError"

	err := s.Delete(instruction.EntityID)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
