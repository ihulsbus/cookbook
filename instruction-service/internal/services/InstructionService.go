package services

import (
	"errors"
	m "instruction-service/internal/models"
)

type InstructionRepository interface {
	Find(instruction m.Instruction) (m.Instruction, error)
	Create(instruction m.Instruction) (m.Instruction, error)
	Update(instruction m.Instruction) (m.Instruction, error)
	Delete(instruction m.Instruction) error
}

type InstructionService struct {
	repo InstructionRepository
}

// NewInstructionService creates a new RecipeService instance
func NewInstructionService(instructionRepo InstructionRepository) *InstructionService {
	return &InstructionService{
		repo: instructionRepo,
	}
}

func (s InstructionService) Find(instructionDTO m.InstructionDTO) (m.InstructionDTO, error) {
	// TODO create logic
	instruction, err := s.repo.Find(instructionDTO.ConvertFromDTO())
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.InstructionDTO{}, err
		default:
			return m.InstructionDTO{}, errors.New("internal server error")
		}
	}

	return instruction.ConvertToDTO(), nil
}

func (s InstructionService) Create(instructionDTO m.InstructionDTO) (m.InstructionDTO, error) {
	// TODO create logic
	instruction, err := s.repo.Create(instructionDTO.ConvertFromDTO())
	if err != nil {
		return m.InstructionDTO{}, err
	}

	return instruction.ConvertToDTO(), nil
}

func (s InstructionService) Update(instructionDTO m.InstructionDTO) (m.InstructionDTO, error) {
	var err error
	if _, err = s.repo.Find(instructionDTO.ConvertFromDTO()); err != nil {
		return m.InstructionDTO{}, errors.New("unable to find existing instruction. cannot update something that does not exist")
	}

	updated, err := s.repo.Update(instructionDTO.ConvertFromDTO())
	if err != nil {
		return m.InstructionDTO{}, err
	}

	return updated.ConvertToDTO(), nil
}

func (s InstructionService) Delete(instructionDTO m.InstructionDTO) error {
	var err error

	_, err = s.repo.Find(instructionDTO.ConvertFromDTO())
	if err != nil {
		return errors.New("unable to find existing instruction. cannot delete something that does not exist")
	}

	err = s.repo.Delete(instructionDTO.ConvertFromDTO())
	if err != nil {
		return err
	}

	return nil
}
