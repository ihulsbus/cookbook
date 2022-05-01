package repositories

import (
	m "github.com/ihulsbus/cookbook/internal/models"
	"gorm.io/gorm"
)

type IngredientFindAll func(r *IngredientRepository) ([]m.Ingredient, error)
type IngredientFindSingle func(r *IngredientRepository, ingredientID int) ([]m.Ingredient, error)
type SectionFindAll func(r *IngredientRepository) error
type SectionFindSingle func(r *IngredientRepository, sectionID int) error
type FindRecipeIngredients func(r *IngredientRepository, recipeID int) ([]m.Recipe_Ingredient, error)
type CreateIngredient func(r *IngredientRepository, ingredient m.Ingredient) (m.Ingredient, error)
type UpdateIngredient func(r *IngredientRepository, ingredient m.Ingredient) (m.Ingredient, error)
type DeleteIngredient func(r *IngredientRepository, ingredient m.Ingredient) error

type IngredientRepository struct {
	db *gorm.DB

	// Ingredients
	IngredientFindAll    IngredientFindAll
	IngredientFindSingle IngredientFindSingle
	CreateIngredient     CreateIngredient
	UpdateIngredient     UpdateIngredient
	DeleteIngredient     DeleteIngredient

	// Sections
	SectionFindAll    SectionFindAll
	SectionFindSingle SectionFindSingle

	// Recipes
	FindRecipeIngredients FindRecipeIngredients
}

func NewIngredientRepository(db *gorm.DB) *IngredientRepository {
	return &IngredientRepository{
		db: db,

		// Ingredients
		IngredientFindAll:    ingredientFindAll,
		IngredientFindSingle: ingredientFindSingle,
		CreateIngredient:     createIngredient,
		UpdateIngredient:     updateIngredient,
		DeleteIngredient:     deleteIngredient,
	}
}

func ingredientFindAll(r *IngredientRepository) ([]m.Ingredient, error) {
	var ingredients []m.Ingredient

	if err := r.db.Find(&ingredients).Error; err != nil {
		return nil, err
	}

	return ingredients, nil
}

func ingredientFindSingle(r *IngredientRepository, ingredientID int) ([]m.Ingredient, error) {
	var ingredient []m.Ingredient

	if err := r.db.Where("id = ?", ingredientID).Find(&ingredient).Error; err != nil {
		return nil, err
	}

	return ingredient, nil
}

func createIngredient(r *IngredientRepository, ingredient m.Ingredient) (m.Ingredient, error) {

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

func updateIngredient(r *IngredientRepository, ingredient m.Ingredient) (m.Ingredient, error) {

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

func deleteIngredient(r *IngredientRepository, ingredient m.Ingredient) error {

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
