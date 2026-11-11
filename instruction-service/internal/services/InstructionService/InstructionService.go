package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
)

type InstructionRepository interface {
	Find(entityID uuid.UUID) (*[]m.Instruction, error)
	Create(instruction *[]m.Instruction) (*[]m.Instruction, error)
	Delete(instruction *[]m.Instruction) error
}

type InstructionService struct {
	repo InstructionRepository
}

// Recipe existence/access is not checked here - the router's
// authz.RequireRecipeAccess middleware already confirmed both
// before this is ever called.
func NewInstructionService(instructionRepo InstructionRepository) *InstructionService {
	return &InstructionService{
		repo: instructionRepo,
	}
}

func (s InstructionService) Find(recipeID uuid.UUID) (*[]m.InstructionDTO, error) {
	instruction, err := s.repo.Find(recipeID)
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	instructionDTO := m.Instruction{}.ConvertAllToDTO(*instruction)
	return &instructionDTO, nil
}

func (s InstructionService) Create(entityID uuid.UUID, instructionDTO *[]m.InstructionDTO) (*[]m.InstructionDTO, error) {
	_, err := s.Find(entityID)
	if err == nil {
		return nil, errors.New("recipe has existing instructions. Use update or delete first")
	}

	instructions := m.InstructionDTO{}.ConvertAllFromDTO(*instructionDTO)

	// override any client-supplied entityID/type with the provided values
	for i := range instructions {
		instructions[i].EntityID = entityID
		instructions[i].EntityType = "recipe"
	}

	instructionResponse, err := s.repo.Create(&instructions)
	if err != nil {
		return nil, err
	}

	instructionResponseDTO := m.Instruction{}.ConvertAllToDTO(*instructionResponse)
	return &instructionResponseDTO, nil
}

func (s InstructionService) Update(entityID uuid.UUID, instructionDTO *[]m.InstructionDTO) (*[]m.InstructionDTO, error) {
	var err error

	if _, err = s.repo.Find(entityID); err != nil {
		return nil, errors.New("unable to find existing recipe. cannot update something that does not exist")
	}

	existingInstructions, err := s.repo.Find(entityID)
	if err != nil {
		return nil, err
	}

	if err = s.Delete(entityID); err != nil {
		return nil, fmt.Errorf("an error occured deleting existing instructions %s", err.Error())
	}

	instructions := m.InstructionDTO{}.ConvertAllFromDTO(*instructionDTO)

	// override any client-supplied entityID/type with the provided values
	for i := range instructions {
		instructions[i].EntityID = entityID
		instructions[i].EntityType = "recipe"
	}

	instructionResponse, err := s.repo.Create(&instructions)
	if err != nil {
		err = fmt.Errorf("an error occured creating new instructions: %s", err.Error())

		_, createErr := s.repo.Create(existingInstructions)
		if createErr != nil {
			err = fmt.Errorf("%s\nadditionally, an second error occured restoring previous instructions: %s",
				err.Error(),
				createErr.Error())
			return nil, err
		}
	}

	instructionResponseDTO := m.Instruction{}.ConvertAllToDTO(*instructionResponse)
	return &instructionResponseDTO, nil
}

func (s InstructionService) Delete(recipeID uuid.UUID) error {
	var err error

	instructions, err := s.repo.Find(recipeID)
	if err != nil {
		return errors.New("unable to find existing instruction. cannot delete something that does not exist")
	}

	err = s.repo.Delete(instructions)
	if err != nil {
		return err
	}

	return nil
}
