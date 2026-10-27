package repositories

import (
	"errors"

	m "github.com/ihulsbus/cookbook/shared/models"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) FindAll(pagination m.PaginationRequest) ([]m.Category, int64, error) {
	var categories []m.Category
	var total int64

	if err := r.db.Model(&m.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(pagination.Limit).Offset(pagination.Offset()).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *CategoryRepository) FindSingle(category m.Category) (m.Category, error) {

	result := r.db.First(&category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.Category{}, errors.New("not found")
		} else {
			return m.Category{}, result.Error
		}
	}

	return category, nil
}

func (r *CategoryRepository) Create(category m.Category) (m.Category, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Create(&category).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return m.Category{}, err
	}

	return category, nil
}

func (r *CategoryRepository) Update(category m.Category) (m.Category, error) {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Updates(&category).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return m.Category{}, err
	}

	return category, nil
}

func (r *CategoryRepository) Delete(category m.Category) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(&category).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
