package repositories

import (
	m "github.com/ihulsbus/cookbook/internal/models"
	"gorm.io/gorm"
)

const (
	whereID = "id = ?"
)

type IngredientRepository struct {
	db *gorm.DB
}

func NewIngredientRepository(db *gorm.DB) *IngredientRepository {
	return &IngredientRepository{
		db: db,
	}
}

func (r IngredientRepository) FindAll() ([]m.Ingredient, error) {
	var ingredients []m.Ingredient

	if err := r.db.Find(&ingredients).Error; err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (r IngredientRepository) FindUnits() ([]m.Unit, error) {
	var units []m.Unit

	if err := r.db.Find(&units).Error; err != nil {
		return nil, err
	}

	return units, nil
}

func (r IngredientRepository) FindSingle(ingredientID int) (m.Ingredient, error) {
	var ingredient m.Ingredient

	if err := r.db.Where("id = ?", ingredientID).Find(&ingredient).Error; err != nil {
		return ingredient, err
	}

	return ingredient, nil
}

func (r IngredientRepository) Create(ingredient m.Ingredient) (m.Ingredient, error) {

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

func (r IngredientRepository) Update(ingredient m.Ingredient) (m.Ingredient, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("ID = ?", ingredient.ID).Updates(&ingredient).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return ingredient, err
	}

	return ingredient, nil
}

func (r IngredientRepository) Delete(ingredient m.Ingredient) error {

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
