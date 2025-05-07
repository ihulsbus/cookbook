package services

import (
	"errors"
	m "instruction-service/internal/models"

	"github.com/google/uuid"
)

type InstructionRepository interface {
	Find(instruction m.Instruction) (m.Instruction, error)
	Create(instruction m.Instruction) (m.Instruction, error)
	Update(instruction m.Instruction) (m.Instruction, error)
	Delete(instruction m.Instruction) error
}

type RecipeRepository interface {
	RecipeExists(recipeID string) (bool, error)
}

type ImageRepository interface {
	ImageExists(imageID string) (bool, error)
}

type InstructionService struct {
	repo   InstructionRepository
	recipe RecipeRepository
	image  ImageRepository
}

// NewInstructionService creates a new RecipeService instance
func NewInstructionService(instructionRepo InstructionRepository, recipe RecipeRepository, image ImageRepository) *InstructionService {
	return &InstructionService{
		repo:   instructionRepo,
		recipe: recipe,
		image:  image,
	}
}

func (s InstructionService) Find(recipeID uuid.UUID) (m.InstructionDTO, error) {
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

func (s InstructionService) Create(recipeID uuid.UUID, instructionDTO m.InstructionDTO) (m.InstructionDTO, error) {

	ok, err := s.recipe.RecipeExists(instructionDTO.EntityID.String())
	if err != nil {
		return m.InstructionDTO{}, err
	}

	if !ok {
		return m.InstructionDTO{}, errors.New("provided recipe does not exist")
	}

	if instructionDTO.MediaID != uuid.Nil {
		ok, err = s.image.ImageExists(instructionDTO.MediaID.String())
		if err != nil {
			return m.InstructionDTO{}, err
		}

		if !ok {
			return m.InstructionDTO{}, errors.New("provided recipe does not exist")
		}
	}

	instruction, err := s.repo.Create(instructionDTO.ConvertFromDTO())
	if err != nil {
		return m.InstructionDTO{}, err
	}

	return instruction.ConvertToDTO(), nil
}

func (s InstructionService) Update(recipeID uuid.UUID, instructionDTO m.InstructionDTO) (m.InstructionDTO, error) {
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

func (s InstructionService) Delete(recipeID uuid.UUID) error {
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
