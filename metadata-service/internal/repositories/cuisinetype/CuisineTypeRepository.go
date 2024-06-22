package repositories

import (
	"errors"

	m "metadata-service/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CuisineTypeRepository struct {
	db *gorm.DB
}

func NewCuisineTypeRepository(db *gorm.DB) *CuisineTypeRepository {
	return &CuisineTypeRepository{
		db: db,
	}
}

func (r *CuisineTypeRepository) FindAll() ([]m.CuisineType, error) {
	var cuisineTypes []m.CuisineType

	if err := r.db.Preload(clause.Associations).Find(&cuisineTypes).Error; err != nil {
		return nil, err
	}

	if len(cuisineTypes) <= 0 {
		return nil, errors.New("not found")
	}

	return cuisineTypes, nil
}

func (r *CuisineTypeRepository) FindSingle(cuisineType m.CuisineType) (m.CuisineType, error) {

	result := r.db.First(&cuisineType)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.CuisineType{}, errors.New("not found")
		} else {
			return m.CuisineType{}, result.Error
		}
	}

	return cuisineType, nil
}

func (r *CuisineTypeRepository) Create(cuisineType m.CuisineType) (m.CuisineType, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Create(&cuisineType).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return m.CuisineType{}, err
	}

	return cuisineType, nil
}

func (r *CuisineTypeRepository) Update(cuisineType m.CuisineType) (m.CuisineType, error) {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Updates(&cuisineType).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return m.CuisineType{}, err
	}

	return cuisineType, nil
}

func (r *CuisineTypeRepository) Delete(cuisineType m.CuisineType) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(&cuisineType).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
