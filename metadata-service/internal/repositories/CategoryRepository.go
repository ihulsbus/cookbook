package repositories

import (
	"errors"

	m "metadata-service/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	whereCategoryID = "id = ?"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) FindAll() ([]m.Category, error) {
	var categories []m.Category

	if err := r.db.Preload(clause.Associations).Find(&categories).Error; err != nil {
		return nil, err
	}

	if len(categories) <= 0 {
		return nil, errors.New("not found")
	}

	return categories, nil
}

func (r *CategoryRepository) FindSingle(categoryID uint) (m.Category, error) {
	var category m.Category

	result := r.db.Preload(clause.Associations).Where(whereCategoryID, categoryID).First(&category)
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
