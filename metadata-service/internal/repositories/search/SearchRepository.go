package repositories

import (
	"errors"
	m "metadata-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SearchRepository struct {
	db *gorm.DB
}

func NewSearchRepository(db *gorm.DB) *SearchRepository {
	return &SearchRepository{
		db: db,
	}
}

func (r *SearchRepository) GetAllRecipeMetadata() (*[]m.MetadataSearchResult, error) {
	var recipeIDs []uuid.UUID

	// Step 1: fetch all recipe IDs from base table
	if err := r.db.
		Table("recipe_difficulty_levels").
		Distinct().
		Pluck("recipe_id", &recipeIDs).Error; err != nil {
		return nil, err
	}

	if len(recipeIDs) == 0 {
		return &[]m.MetadataSearchResult{}, nil
	}

	// Step 2: bulk fetch all related data
	type row struct {
		RecipeID uuid.UUID
		ValueID  uuid.UUID
	}

	// Helper map builders
	fetchMany := func(table, column string) (map[uuid.UUID][]uuid.UUID, error) {
		var rows []row
		result := make(map[uuid.UUID][]uuid.UUID)

		err := r.db.Table(table).
			Select("recipe_id, "+column+" as value_id").
			Where("recipe_id IN ?", recipeIDs).
			Scan(&rows).Error
		if err != nil {
			return nil, err
		}

		for _, r := range rows {
			result[r.RecipeID] = append(result[r.RecipeID], r.ValueID)
		}

		return result, nil
	}

	fetchOne := func(table, column string) (map[uuid.UUID]uuid.UUID, error) {
		multi, err := fetchMany(table, column)
		if err != nil {
			return nil, err
		}
		single := make(map[uuid.UUID]uuid.UUID)
		for k, v := range multi {
			if len(v) > 0 {
				single[k] = v[0]
			}
		}
		return single, nil
	}

	// Fetch related metadata
	categories, err := fetchMany("recipe_categories", "category_id")
	if err != nil {
		return nil, err
	}
	tags, err := fetchMany("recipe_tags", "tag_id")
	if err != nil {
		return nil, err
	}
	difficulties, err := fetchOne("recipe_difficulty_levels", "difficulty_level_id")
	if err != nil {
		return nil, err
	}
	prepTimes, err := fetchOne("recipe_preparation_times", "preparation_time_id")
	if err != nil {
		return nil, err
	}
	cuisines, err := fetchOne("recipe_cuisine_types", "cuisine_type_id")
	if err != nil {
		return nil, err
	}

	// Step 3: Assemble results
	results := make([]m.MetadataSearchResult, 0, len(recipeIDs))
	for _, id := range recipeIDs {
		results = append(results, m.MetadataSearchResult{
			RecipeID:          id,
			CategoryIDs:       categories[id],
			TagIDs:            tags[id],
			DifficultyLevelID: difficulties[id],
			PreparationTimeID: prepTimes[id],
			CuisineTypeID:     cuisines[id],
		})
	}

	return &results, nil
}

func (r *SearchRepository) SearchMetadata(request m.MetadataSearchRequest) ([]m.MetadataSearchResult, error) {
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
