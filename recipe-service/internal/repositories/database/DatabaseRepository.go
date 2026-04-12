package repositories

import (
	"errors"

	m "github.com/ihulsbus/cookbook/shared/models"
	"gorm.io/gorm"
)

type DatabaseRepository struct {
	db *gorm.DB
}

func NewDatabaseRepository(db *gorm.DB) *DatabaseRepository {
	return &DatabaseRepository{
		db: db,
	}
}

// FindAll retrieves all recipes from the database and returns them in a slice
func (r DatabaseRepository) FindAll() ([]m.Recipe, error) {
	var recipes []m.Recipe

	if err := r.db.Find(&recipes).Error; err != nil {
		return nil, err
	}
	if len(recipes) <= 0 {
		return nil, errors.New("not found")
	}

	return recipes, nil
}

// Find searches for a specific recipe in the database and returns it when found.
func (r DatabaseRepository) FindSingle(recipe m.Recipe) (m.Recipe, error) {

	result := r.db.First(&recipe)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.Recipe{}, errors.New("not found")
		} else {
			return m.Recipe{}, result.Error
		}
	}

	return recipe, nil
}

// Create handles the creation of a recipe and stores the relevant information in the database
func (r DatabaseRepository) Create(recipe m.Recipe) (m.Recipe, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Create(&recipe).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return recipe, err
	}

	return recipe, nil
}

func (r DatabaseRepository) Update(recipe m.Recipe) (m.Recipe, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Updates(&recipe).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return recipe, err
	}
	return recipe, nil
}

func (r DatabaseRepository) Delete(recipe m.Recipe) error {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(&recipe).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
