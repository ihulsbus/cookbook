package repositories

import (
	"errors"

	m "ingredient-service/internal/models"

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

	if len(ingredients) <= 0 {
		return nil, errors.New("not found")
	}

	return ingredients, nil
}

func (r IngredientRepository) FindUnits() ([]m.Unit, error) {
	var units []m.Unit

	if err := r.db.Find(&units).Error; err != nil {
		return nil, err
	}

	if len(units) <= 0 {
		return nil, errors.New("not found")

	}

	return units, nil
}

func (r IngredientRepository) FindSingle(ingredientID uint) (m.Ingredient, error) {
	var ingredient m.Ingredient

	result := r.db.Where(whereID, ingredientID).First(&ingredient)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.Ingredient{}, errors.New("not found")
		} else {
			return m.Ingredient{}, result.Error
		}
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

		if err := tx.Where(whereID, ingredient.ID).Updates(&ingredient).Error; err != nil {
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
