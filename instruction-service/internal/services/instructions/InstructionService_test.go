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
		ID:          uuid.New(),
		Sequence:    1,
		Description: "instruction",
		MediaID:     uuid.New(),
	}
)

type InstructionRepositoryMock struct{}

func (InstructionRepositoryMock) Find(instructionInput m.Instruction) (m.Instruction, error) {
	switch instructionInput.Description {
	case "find":
		return instruction, nil
	case "create":
		return instruction, nil
	case "update":
		return instruction, nil
	case "updateerror":
		return instruction, nil
	case "delete":
		return instruction, nil
	case "deleteerror":
		return instruction, nil
	case "notfound":
		return m.Instruction{}, errors.New("not found")
	default:
		return m.Instruction{}, errors.New("error")
	}
}

func (InstructionRepositoryMock) Create(instructionInput m.Instruction) (m.Instruction, error) {
	switch instructionInput.Description {
	case "create":
		return instruction, nil
	default:
		return instruction, errors.New("error")
	}
}

func (InstructionRepositoryMock) Update(instructionInput m.Instruction) (m.Instruction, error) {
	switch instructionInput.Description {
	case "update":
		return instruction, nil
	default:
		return instruction, errors.New("error")
	}
}

func (InstructionRepositoryMock) Delete(instructionInput m.Instruction) error {
	switch instructionInput.Description {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

// ========================================================================================================

func TestFindInstruction_OK(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{})

	instructionDTO := m.InstructionDTO{
		ID:          instruction.ID,
		Description: "find",
	}
	result, err := s.Find(instructionDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.InstructionDTO{}, result)
	assert.Equal(t, "instruction", result.Description)
}

func TestFindInstruction_NotFoundErr(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{})

	instructionDTO := m.InstructionDTO{
		ID:          instruction.ID,
		Description: "notfound",
	}
	result, err := s.Find(instructionDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.IsType(t, m.InstructionDTO{}, result)
}

func TestFindInstruction_Err(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{})

	instructionDTO := m.InstructionDTO{
		ID:          instruction.ID,
		Description: "error",
	}
	result, err := s.Find(instructionDTO)

	assert.Error(t, err)
	assert.IsType(t, m.InstructionDTO{}, result)
	assert.EqualError(t, err, "internal server error")

}

func TestCreateInstruction_OK(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{})

	instructionDTO := m.InstructionDTO{
		Sequence:    instruction.Sequence,
		Description: "create",
		MediaID:     instruction.MediaID,
	}
	result, err := s.Create(instructionDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.InstructionDTO{}, result)
	assert.Equal(t, result.Description, "instruction")
}

func TestCreateInstruction_Err(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{})

	instructionDTO := m.InstructionDTO{
		Sequence:    instruction.Sequence,
		Description: "error",
		MediaID:     instruction.MediaID,
	}
	result, err := s.Create(instructionDTO)

	assert.Error(t, err)
	assert.IsType(t, m.InstructionDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestUpdateInstruction_OK(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{})

	instructionDTO := m.InstructionDTO{
		ID:          instruction.ID,
		Description: "update",
	}
	result, err := s.Update(instructionDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.InstructionDTO{}, result)
	assert.Equal(t, instruction.ID, result.ID)
	assert.Equal(t, instruction.Description, result.Description)
}

func TestUpdateInstruction_FindErr(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{})

	instructionDTO := m.InstructionDTO{
		ID:          instruction.ID,
		Description: "error",
	}
	result, err := s.Update(instructionDTO)

	assert.Error(t, err)
	assert.Equal(t, m.InstructionDTO{}, result)
	assert.EqualError(t, err, "unable to find existing instruction. cannot update something that does not exist")
}

func TestUpdateInstruction_UpdateErr(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{})

	instructionDTO := m.InstructionDTO{
		ID:          instruction.ID,
		Description: "find",
	}
	result, err := s.Update(instructionDTO)

	assert.Error(t, err)
	assert.Equal(t, m.InstructionDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestDeleteInstruction_OK(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{})

	instructionDTO := m.InstructionDTO{
		ID:          instruction.ID,
		Description: "delete",
	}
	err := s.Delete(instructionDTO)

	assert.NoError(t, err)
}

func TestDeleteInstruction_FindErr(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{})

	instructionDTO := m.InstructionDTO{
		ID:          instruction.ID,
		Description: "error",
	}
	err := s.Delete(instructionDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "unable to find existing instruction. cannot delete something that does not exist")
}

func TestDeleteInstruction_DeleteErr(t *testing.T) {
	s := NewInstructionService(&InstructionRepositoryMock{})

	instructionDTO := m.InstructionDTO{
		ID:          instruction.ID,
		Description: "deleteerror",
	}
	err := s.Delete(instructionDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
