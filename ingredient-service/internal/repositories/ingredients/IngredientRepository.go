package repositories

import (
	"errors"

	m "github.com/ihulsbus/cookbook/shared/models"
	"gorm.io/gorm"
)

type IngredientRepository struct {
	db *gorm.DB
}

func NewIngredientRepository(db *gorm.DB) *IngredientRepository {
	return &IngredientRepository{
		db: db,
	}
}

func (r IngredientRepository) FindAll(pagination m.PaginationRequest) ([]m.Ingredient, int64, error) {
	var ingredients []m.Ingredient
	var total int64

	if err := r.db.Model(&m.Ingredient{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(pagination.Limit).Offset(pagination.Offset()).Find(&ingredients).Error; err != nil {
		return nil, 0, err
	}

	return ingredients, total, nil
}

func (r IngredientRepository) FindSingle(ingredient m.Ingredient) (m.Ingredient, error) {

	result := r.db.First(&ingredient)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.Ingredient{}, errors.New("not found")
		} else {
			return m.Ingredient{}, result.Error
		}
	}

	return ingredient, nil
}

func (r IngredientRepository) FindByName(name string) (m.Ingredient, error) {
	var ing m.Ingredient
	err := r.db.
		Where("name = ?", name).
		First(&ing).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// no record, signal “not found” to service
			return m.Ingredient{}, nil
		}
		// some other DB error
		return m.Ingredient{}, err
	}

	return ing, nil
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

		if err := tx.Updates(&ingredient).Error; err != nil {
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
