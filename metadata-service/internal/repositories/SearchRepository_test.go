package repositories

import (
	"regexp"
	"testing"

	m "metadata-service/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	id         uuid.UUID = uuid.New()
	difficulty int       = 1
	cuisine    string    = ""
	minprep    int       = 1
	maxprep    int       = 2

	searchRequest m.MetadataSearchRequest = m.MetadataSearchRequest{
		RecipeID:          &id,
		CategoryID:        &id,
		TagID:             &id,
		DifficultyLevelID: &id,
		CuisineTypeID:     &id,
		MinPrepTime:       &minprep,
		MaxPrepTime:       &maxprep,
	}
)

func TestSearch_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewSearchRepository(db)

	// Define expected query
	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT recipe_categories.recipe_id FROM "recipe_categories" 
				LEFT JOIN recipe_tags ON recipe_categories.recipe_id = recipe_tags.recipe_id 
				LEFT JOIN recipe_difficulty_levels ON recipe_categories.recipe_id = recipe_difficulty_levels.recipe_id 
				LEFT JOIN recipe_preparation_times ON recipe_categories.recipe_id = recipe_preparation_times.recipe_id 
				LEFT JOIN recipe_cuisine_types ON recipe_categories.recipe_id = recipe_cuisine_types.recipe_id 
				WHERE recipe_categories.category_id = $1 
				AND recipe_tags.tag_id = $2 
				AND recipe_difficulty_levels.difficulty_level_id = $3 
				AND recipe_preparation_times.preparation_time >= $4 
				AND recipe_preparation_times.preparation_time <= $5 
				AND recipe_cuisine_types.cuisine_type_id = $6 
				GROUP BY "recipe_categories"."recipe_id"
    `)).
		WithArgs(
			searchRequest.CategoryID,
			searchRequest.TagID,
			searchRequest.DifficultyLevelID,
			searchRequest.MinPrepTime,
			searchRequest.MaxPrepTime,
			searchRequest.CuisineTypeID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id"}).AddRow(searchRequest.RecipeID))

	result, err := r.SearchMetadata(searchRequest)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}
