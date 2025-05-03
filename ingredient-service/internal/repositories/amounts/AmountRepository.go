package repositories

import (
	"errors"
	"github.com/google/uuid"

	"gorm.io/gorm"
	m "ingredient-service/internal/models"
)

type AmountRepository struct {
	db *gorm.DB
}

func NewAmountRepository(db *gorm.DB) *AmountRepository {
	return &AmountRepository{
		db: db,
	}
}

func (r AmountRepository) Find(recipeID uuid.UUID) (*[]m.Amount, error) {
	var amount []m.Amount

	result := r.db.Find(&amount, "recipe_id = ?", recipeID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		} else {
			return nil, result.Error
		}
	}

	return &amount, nil
}

func (r AmountRepository) Create(amounts *[]m.Amount) (*[]m.Amount, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(amounts).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return amounts, nil
}

func (r AmountRepository) Update(amounts *[]m.Amount) (*[]m.Amount, error) {
	input := *amounts

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Updates(input).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return amounts, nil
}

func (r AmountRepository) Delete(amounts *[]m.Amount) error {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(amounts).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
