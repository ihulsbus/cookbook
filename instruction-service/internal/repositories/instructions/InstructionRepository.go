package repositories

import (
	"errors"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"

	"gorm.io/gorm"
)

type InstructionRepository struct {
	db *gorm.DB
}

func NewInstructionRepository(db *gorm.DB) *InstructionRepository {
	return &InstructionRepository{
		db: db,
	}
}

func (r InstructionRepository) Find(entityID uuid.UUID) (*[]m.Instruction, error) {
	var instruction []m.Instruction

	result := r.db.Find(&instruction, "entity_id = ?", entityID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		} else {
			return nil, result.Error
		}
	}

	if len(instruction) == 0 {
		// no rows => treat as “not found”
		return nil, errors.New("not found")
	}

	return &instruction, nil
}

func (r InstructionRepository) Create(instruction *[]m.Instruction) (*[]m.Instruction, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(instruction).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return instruction, nil
}

func (r InstructionRepository) Delete(instruction *[]m.Instruction) error {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(&instruction).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
