package repositories

import (
	"errors"

	m "instruction-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	whereRecipeID = "recipe_id = ?"
)

type InstructionRepository struct {
	db *gorm.DB
}

func NewInstructionRepository(db *gorm.DB) *InstructionRepository {
	return &InstructionRepository{
		db: db,
	}
}

func (r InstructionRepository) FindInstruction(recipeID uuid.UUID) (m.Instruction, error) {
	var instruction m.Instruction
	instruction.RecipeID = recipeID

	result := r.db.Where(whereRecipeID, recipeID).First(&instruction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.Instruction{}, errors.New("not found")
		} else {
			return m.Instruction{}, result.Error
		}
	}

	return instruction, nil
}

func (r InstructionRepository) CreateInstruction(instruction m.Instruction) (m.Instruction, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Create(&instruction).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return instruction, err
	}
	return instruction, nil
}

func (r InstructionRepository) UpdateInstruction(instruction m.Instruction) (m.Instruction, error) {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Where(whereRecipeID, &instruction.RecipeID).Updates(&instruction).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return instruction, err
	}

	return instruction, nil
}

func (r InstructionRepository) DeleteInstruction(instruction m.Instruction) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Where(whereRecipeID, &instruction.RecipeID).Delete(&instruction).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
