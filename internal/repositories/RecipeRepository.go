package repositories

import (
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

	return recipes, nil
}

// Find searches for a specific recipe in the database and returns it when found.
func (r RecipeRepository) FindSingle(recipeID uint) (m.Recipe, error) {
	var recipe m.Recipe
	recipe.ID = recipeID

	if err := r.db.Preload(clause.Associations).Find(&recipe).Error; err != nil {
		return recipe, err
	}

	return recipe, nil
}

// FindRecipeIngredients finds all ingredients associated to a recipe and returns them in a slice
// func findRecipeIngredients(r *RecipeRepository, recipeID int) ([]m.Recipe_Ingredient, error) {
// 	var recipeIngredients []m.Recipe_Ingredient

// 	return recipeIngredients, nil
// }

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

	// if err := r.db.Transaction(func(tx *gorm.DB) error {
	// 	var err error

	// 	recipe, err = updateIngredientRelations(tx, recipe)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return nil
	// }); err != nil {
	// 	return recipe, err
	// }

	return recipe, nil
}

func (r RecipeRepository) Update(recipe m.Recipe) (m.Recipe, error) {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		// update the recipe while skipping associations as this is very very
		if err = tx.Updates(&recipe).Error; err != nil {
			return err
		}

		// recipe, err = updateIngredientRelations(tx, recipe)
		// if err != nil {
		// 	return err
		// }

		// recipe, err = updateTagRelations(tx, recipe)
		// if err != nil {
		// 	return err
		// }

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

	if err := r.db.Find(&instruction).Where(whereRecipeID, recipeID).Error; err != nil {
		return instruction, err
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

// func updateIngredientRelations(tx *gorm.DB, recipe m.Recipe) (m.Recipe, error) {
// 	var ingredientIDs []int

// 	for i := range recipe.IngredientAmounts {
// 		ingredientIDs = append(ingredientIDs, recipe.IngredientAmounts[i].IngredientID)
// 	}

// 	// GORM update relations is very broken... we'll have to do it ourselves
// 	var existing []m.Recipe_Ingredient

// 	if err := tx.Raw("SELECT * FROM recipe_ingredients WHERE recipe_id = ?", recipe.RecipeID).Scan(&existing).Error; err != nil {
// 		return recipe, err
// 	}

// 	for _, ingredientAmount := range recipe.IngredientAmounts {
// 		args := map[string]interface{}{"recipeid": recipe.RecipeID, "ingredientid": ingredientAmount.IngredientID, "amount": ingredientAmount.Amount, "unit": ingredientAmount.Unit}

// 		idx := slices.IndexFunc(existing, func(r m.Recipe_Ingredient) bool { return r.IngredientID == ingredientAmount.IngredientID })

// 		if idx == -1 {
// 			if err := tx.Raw("INSERT INTO recipe_ingredients (recipe_id, ingredient_id, amount, unit) VALUES(@recipeid, @ingredientid, @amount, @unit)", args).Scan(&recipe.IngredientAmounts).Error; err != nil {
// 				return recipe, err
// 			}
// 		} else {
// 			if err := tx.Raw("UPDATE recipe_ingredients SET amount = @amount, unit = @unit WHERE recipe_id = @recipeid AND ingredient_id = @ingredientid", args).Scan(&recipe.IngredientAmounts).Error; err != nil {
// 				return recipe, err
// 			}
// 		}
// 	}

// 	if err := tx.Unscoped().Exec("DELETE FROM recipe_ingredients WHERE recipe_id = ? AND ingredient_id NOT IN ?", recipe.RecipeID, ingredientIDs).Error; err != nil {
// 		return recipe, err
// 	}

// 	return recipe, nil
// }

// func updateTagRelations(tx *gorm.DB, recipe m.Recipe) (m.Recipe, error) {
// 	var tagIDs []int

// 	for i := range recipe.Tags {
// 		tagIDs = append(tagIDs, recipe.Tags[i].TagID)
// 	}

// 	// GORM update relations is very broken... we'll have to do it ourselves
// 	var existing []m.RecipeTag

// 	if err := tx.Raw("SELECT * FROM recipe_tags WHERE recipe_id = ?", recipe.RecipeID).Scan(&existing).Error; err != nil {
// 		return recipe, err
// 	}

// 	for _, tag := range recipe.Tags {
// 		args := map[string]interface{}{"recipeid": recipe.RecipeID, "tagid": tag.TagID}

// 		idx := slices.IndexFunc(existing, func(r m.RecipeTag) bool { return r.TagID == tag.TagID })

// 		if idx == -1 {
// 			if err := tx.Raw("INSERT INTO recipe_tags (recipe_id, tags_id) VALUES(@recipeid, @tagid)", args).Scan(&recipe.IngredientAmounts).Error; err != nil {
// 				return recipe, err
// 			}
// 		} else {
// 			return recipe, nil
// 		}
// 	}

// 	if err := tx.Unscoped().Exec("DELETE FROM recipe_tags WHERE recipe_id = ? AND tags_id NOT IN ?", recipe.RecipeID, tagIDs).Error; err != nil {
// 		return recipe, err
// 	}

// 	return recipe, nil
// }
