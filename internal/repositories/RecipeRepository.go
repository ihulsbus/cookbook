package repositories

import (
	"errors"

	m "github.com/ihulsbus/cookbook/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	whereRecipeID = "recipe_id = ?"
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

func (r RecipeRepository) FindInstruction(recipeID uint) (m.Instruction, error) {
	var instruction m.Instruction
	instruction.RecipeID = recipeID

	result := r.db.Where(whereRecipeID, recipeID).First(&instruction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.Instruction{}, errors.New("not found")
		} else {
			return m.Instruction{}, result.Error
		}
	}

	return instruction, nil
}

func (r RecipeRepository) CreateInstruction(instruction m.Instruction) (m.Instruction, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Create(&instruction).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return instruction, err
	}
	return instruction, nil
}

func (r RecipeRepository) UpdateInstruction(instruction m.Instruction) (m.Instruction, error) {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Where(whereRecipeID, &instruction.RecipeID).Updates(&instruction).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return instruction, err
	}

	return instruction, nil
}

func (r RecipeRepository) DeleteInstruction(instruction m.Instruction) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Where(whereRecipeID, &instruction.RecipeID).Delete(&instruction).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// FindRecipeIngredients finds all ingredients associated to a recipe and returns them in a slice
func (r RecipeRepository) FindIngredientLink(recipeID uint) ([]m.RecipeIngredient, error) {
	var recipeIngredients []m.RecipeIngredient

	if err := r.db.Where(whereRecipeID, &recipeID).Joins("Unit").Find(&recipeIngredients).Error; err != nil {
		return nil, err
	}

	if len(recipeIngredients) <= 0 {
		return nil, errors.New("not found")
	}

	return recipeIngredients, nil
}

func (r RecipeRepository) CreateIngredientLink(link []m.RecipeIngredient) ([]m.RecipeIngredient, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		for i := range link {
			if err = tx.Create(&link[i]).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return link, nil
}

func (r RecipeRepository) UpdateIngredientLink(link []m.RecipeIngredient) ([]m.RecipeIngredient, error) {
	if err := r.db.Transaction(func(tx *gorm.DB) error {

		for i := range link {

			result := tx.Where(whereRecipeID, &link[i].RecipeID).Updates(&link[i])
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				if err := tx.Create(&link[i]).Error; err != nil {
					return err
				}
			}

		}
		return nil
	}); err != nil {
		return link, err
	}

	return link, nil
}

func (r RecipeRepository) DeleteIngredientLink(link []m.RecipeIngredient) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {

		for i := range link {

			if err := tx.Where(whereRecipeID, &link[i].RecipeID).Delete(&link[i]).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
