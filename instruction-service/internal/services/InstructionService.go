package services

import (
	"errors"
	m "instruction-service/internal/models"

	"github.com/google/uuid"
)

type InstructionRepository interface {
	FindInstruction(recipeID uuid.UUID) (m.Instruction, error)
	CreateInstruction(instruction m.Instruction) (m.Instruction, error)
	UpdateInstruction(instruction m.Instruction) (m.Instruction, error)
	DeleteInstruction(instruction m.Instruction) error
}

type InstructionService struct {
	repo InstructionRepository
}

// NewRecipeService creates a new RecipeService instance
func NewRecipeService(instructionRepo InstructionRepository) *InstructionService {
	return &InstructionService{
		repo: instructionRepo,
	}
}

func (s InstructionService) FindInstruction(recipeID uuid.UUID) (m.Instruction, error) {
	// TODO create logic
	instruction, err := s.repo.FindInstruction(recipeID)
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.Instruction{}, err
		default:
			return m.Instruction{}, errors.New("internal server error")
		}
	}

	return instruction, nil
}

func (s InstructionService) CreateInstruction(instruction m.Instruction, recipeID uuid.UUID) (m.Instruction, error) {
	// TODO create logic
	instruction, err := s.repo.CreateInstruction(instruction)
	if err != nil {
		return m.Instruction{}, err
	}

	return instruction, nil
}

func (s InstructionService) UpdateInstruction(instruction m.Instruction, recipeID uuid.UUID) (m.Instruction, error) {
	var err error
	if _, err = s.repo.FindInstruction(recipeID); err != nil {
		return m.Instruction{}, errors.New("unable to find existing instruction for the given recipe id")
	}

	if instruction.RecipeID != recipeID {
		instruction.RecipeID = recipeID
	}

	updated, err := s.repo.UpdateInstruction(instruction)
	if err != nil {
		return m.Instruction{}, err
	}

	return updated, nil
}

func (s InstructionService) DeleteInstruction(recipeID uuid.UUID) error {
	var err error

	instruction, err := s.repo.FindInstruction(recipeID)
	if err != nil {
		return errors.New("unable to find existing instruction for the given recipe id")
	}

	err = s.repo.DeleteInstruction(instruction)
	if err != nil {
		return err
	}

	return nil
}
