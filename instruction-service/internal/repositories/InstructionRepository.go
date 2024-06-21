package repositories

import (
	"errors"

	m "instruction-service/internal/models"

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

func (r InstructionRepository) Find(instruction m.Instruction) (m.Instruction, error) {
	result := r.db.First(&instruction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.Instruction{}, errors.New("not found")
		} else {
			return m.Instruction{}, result.Error
		}
	}

	return instruction, nil
}

func (r InstructionRepository) Create(instruction m.Instruction) (m.Instruction, error) {

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

func (r InstructionRepository) Update(instruction m.Instruction) (m.Instruction, error) {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Updates(&instruction).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return instruction, err
	}

	return instruction, nil
}

func (r InstructionRepository) Delete(instruction m.Instruction) error {
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
