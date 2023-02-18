package repositories

import (
	m "github.com/ihulsbus/cookbook/internal/models"
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
	var categorys []m.Category

	if err := r.db.Preload(clause.Associations).Find(&categorys).Error; err != nil {
		return nil, err
	}

	return categorys, nil
}

func (r *CategoryRepository) FindSingle(categoryID uint) (m.Category, error) {
	var category m.Category

	if err := r.db.Preload(clause.Associations).Where(whereCategoryID, categoryID).Find(&category).Error; err != nil {
		return m.Category{}, err
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
