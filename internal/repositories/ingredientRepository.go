package repositories

import (
	log "github.com/sirupsen/logrus"

	m "github.com/ihulsbus/cookbook/internal/models"
	"gorm.io/gorm"
)

type IngredientRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewIngredientRepository(db *gorm.DB, logger *log.Logger) *IngredientRepository {
	return &IngredientRepository{
		db:     db,
		logger: logger,
	}
}

func (r IngredientRepository) IngredientFindAll() ([]m.Ingredient, error) {
	var ingredients []m.Ingredient

	if err := r.db.Find(&ingredients).Error; err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (r IngredientRepository) IngredientFindSingle(ingredientID int) ([]m.Ingredient, error) {
	var ingredient []m.Ingredient

	if err := r.db.Where("id = ?", ingredientID).Find(&ingredient).Error; err != nil {
		return nil, err
	}

	return ingredient, nil
}

func (r IngredientRepository) CreateIngredient(ingredient m.Ingredient) (m.Ingredient, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&ingredient).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return ingredient, err
	}

	return ingredient, nil
}

func (r IngredientRepository) UpdateIngredient(ingredient m.Ingredient) (m.Ingredient, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&ingredient).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return ingredient, err
	}

	return ingredient, nil
}

func (r IngredientRepository) DeleteIngredient(ingredient m.Ingredient) error {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(&ingredient).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
