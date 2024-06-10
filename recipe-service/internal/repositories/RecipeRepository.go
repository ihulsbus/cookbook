package repositories

import (
	"errors"

	m "recipe-service/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	whereRecipeID = "recipe_id = ?"
	whereID       = "id = ?"
)

type RecipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) *RecipeRepository {
	return &RecipeRepository{
		db: db,
	}
}

// FindAll retrieves all recipes from the database and returns them in a slice
func (r RecipeRepository) FindAll() ([]m.Recipe, error) {
	var recipes []m.Recipe

	if err := r.db.Preload(clause.Associations).Find(&recipes).Error; err != nil {
		return nil, err
	}
	if len(recipes) <= 0 {
		return nil, errors.New("not found")
	}

	return recipes, nil
}

// Find searches for a specific recipe in the database and returns it when found.
func (r RecipeRepository) FindSingle(recipeID uint) (m.Recipe, error) {
	var recipe m.Recipe
	recipe.ID = recipeID

	result := r.db.Preload(clause.Associations).Where(whereID, recipeID).First(&recipe)
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
func (r RecipeRepository) Create(recipe m.Recipe) (m.Recipe, error) {

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

func (r RecipeRepository) Update(recipe m.Recipe) (m.Recipe, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		// update the recipe while skipping associations as this is very very
		if err = tx.Omit("author_id", "image_name").Updates(&recipe).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return recipe, err
	}
	return recipe, nil
}

func (r RecipeRepository) Delete(recipe m.Recipe) error {

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
