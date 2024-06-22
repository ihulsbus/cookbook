package repositories

import (
	"errors"
	m "metadata-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SearcRepository struct {
	db *gorm.DB
}

func NewSearchRepository(db *gorm.DB) *SearcRepository {
	return &SearcRepository{
		db: db,
	}
}

func (r *SearcRepository) SearchMetadata(request m.MetadataSearchRequest) ([]m.MetadataSearchResult, error) {
	var results []m.MetadataSearchResult

	// Start with base query. we do this on categories as all recipes need to have a category
	query := r.db.Table("recipe_categories").
		Select("recipe_categories.recipe_id").
		Joins("LEFT JOIN recipe_tags ON recipe_categories.recipe_id = recipe_tags.recipe_id").
		Joins("LEFT JOIN recipe_difficulty_levels ON recipe_categories.recipe_id = recipe_difficulty_levels.recipe_id").
		Joins("LEFT JOIN recipe_preparation_times ON recipe_categories.recipe_id = recipe_preparation_times.recipe_id").
		Joins("LEFT JOIN recipe_cuisine_types ON recipe_categories.recipe_id = recipe_cuisine_types.recipe_id").
		Group("recipe_categories.recipe_id")

		// Apply filters based on request
	if *request.CategoryID != uuid.Nil {
		query = query.Where("recipe_categories.category_id = ?", *request.CategoryID)
	}

	if *request.TagID != uuid.Nil {
		query = query.Where("recipe_tags.tag_id = ?", *request.TagID)
	}

	if *request.DifficultyLevelID != uuid.Nil {
		query = query.Where("recipe_difficulty_levels.difficulty_level_id = ?", *request.DifficultyLevelID)
	}

	if *request.CuisineTypeID != uuid.Nil {
		query = query.Where("recipe_cuisine_types.cuisine_type_id = ?", *request.CuisineTypeID)
	}

	// Handle preparation time range
	if request.MinPrepTime != nil || request.MaxPrepTime != nil {
		var prepTimeIDs []uuid.UUID

		prepTimeQuery := r.db.Table("preparation_times").Select("id")

		if request.MinPrepTime != nil {
			prepTimeQuery = prepTimeQuery.Where("duration >= ?", *request.MinPrepTime)
		}

		if request.MaxPrepTime != nil {
			prepTimeQuery = prepTimeQuery.Where("duration <= ?", *request.MaxPrepTime)
		}

		if err := prepTimeQuery.Pluck("id", &prepTimeIDs).Error; err != nil {
			return nil, err
		}

		query = query.Where("recipe_preparation_times.preparation_time_id IN ?", prepTimeIDs)
	}

	// Scan results into a slice of recipe IDs
	var recipeIDs []uuid.UUID
	if err := query.Pluck("recipe_categories.recipe_id", &recipeIDs).Error; err != nil {
		return nil, err
	}

	// Collect detailed metadata for each recipe ID
	for _, recipeID := range recipeIDs {
		var categoryID, tagID, difficultyLevelID, preparationTimeID, cuisineTypeID []uuid.UUID

		r.db.Table("recipe_categories").Where("recipe_id = ?", recipeID).Pluck("category_id", &categoryID)
		r.db.Table("recipe_tags").Where("recipe_id = ?", recipeID).Pluck("tag_id", &tagID)
		r.db.Table("recipe_difficulty_levels").Where("recipe_id = ?", recipeID).Pluck("difficulty_level_id", &difficultyLevelID)
		r.db.Table("recipe_preparation_times").Where("recipe_id = ?", recipeID).Pluck("preparation_time_id", &preparationTimeID)
		r.db.Table("recipe_cuisine_types").Where("recipe_id = ?", recipeID).Pluck("cuisine_type_id", &cuisineTypeID)

		if len(difficultyLevelID) < 1 {
			return nil, errors.New("no difficulty level associated with a recipe. this should be impossible")
		}

		if len(preparationTimeID) < 1 {
			return nil, errors.New("no preparation time associated with a recipe. this should be impossible")
		}

		if len(cuisineTypeID) < 1 {
			return nil, errors.New("no cuisine type associated with a recipe. this should be impossible")
		}

		results = append(results, m.MetadataSearchResult{
			RecipeID:          recipeID,
			CategoryIDs:       categoryID,
			TagIDs:            tagID,
			DifficultyLevelID: difficultyLevelID[0],
			PreparationTimeID: preparationTimeID[0],
			CuisineTypeID:     cuisineTypeID[0],
		})
	}

	return results, nil
}
