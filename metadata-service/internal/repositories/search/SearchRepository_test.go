package repositories

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"

	co "metadata-service/internal/common/test"
)

var (
	id      uuid.UUID = uuid.New()
	minprep int       = 1
	maxprep int       = 2

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
	db, mock := co.NewMockDatabase(t)
	r := NewSearchRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM "preparation_times" WHERE duration >= $1 AND duration <= $2`)).
		WithArgs(
			searchRequest.MinPrepTime,
			searchRequest.MaxPrepTime,
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

	requesterID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT recipe_categories.recipe_id FROM "recipe_categories"
				LEFT JOIN recipe_tags ON recipe_categories.recipe_id = recipe_tags.recipe_id
				LEFT JOIN recipe_difficulty_levels ON recipe_categories.recipe_id = recipe_difficulty_levels.recipe_id
				LEFT JOIN recipe_preparation_times ON recipe_categories.recipe_id = recipe_preparation_times.recipe_id
				LEFT JOIN recipe_cuisine_types ON recipe_categories.recipe_id = recipe_cuisine_types.recipe_id
				WHERE (recipe_categories.recipe_id IN (SELECT recipe_id FROM recipe_accesses WHERE (owner_id = $1 OR visibility = $2) AND NOT deleted))
				AND recipe_categories.category_id = $3
				AND recipe_tags.tag_id = $4
				AND recipe_difficulty_levels.difficulty_level_id = $5
				AND recipe_cuisine_types.cuisine_type_id = $6
				AND recipe_preparation_times.preparation_time_id IN ($7)
				GROUP BY "recipe_categories"."recipe_id"
    `)).
		WithArgs(
			requesterID.String(),
			"public",
			searchRequest.CategoryID,
			searchRequest.TagID,
			searchRequest.DifficultyLevelID,
			searchRequest.CuisineTypeID,
			id,
		).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id"}).AddRow(id))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "category_id" FROM "recipe_categories" WHERE recipe_id = $1`)).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "tag_id" FROM "recipe_tags" WHERE recipe_id = $1`)).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "difficulty_level_id" FROM "recipe_difficulty_levels" WHERE recipe_id = $1`)).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "preparation_time_id" FROM "recipe_preparation_times" WHERE recipe_id = $1`)).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "cuisine_type_id" FROM "recipe_cuisine_types" WHERE recipe_id = $1`)).WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

	result, err := r.SearchMetadata(searchRequest, requesterID)

	assert.NoError(t, err)
	assert.IsType(t, []m.MetadataSearchResult{}, result)
	assert.Len(t, result, 1)
	assert.Len(t, result[0].CategoryIDs, 1)
	assert.Len(t, result[0].TagIDs, 1)
	assert.Equal(t, result[0].RecipeID, id)
	assert.Equal(t, result[0].DifficultyLevelID, id)
	assert.Equal(t, result[0].CuisineTypeID, id)
	assert.Equal(t, result[0].PreparationTimeID, id)

}

// Without the recipe_accesses filter, this would return metadata for every
// recipe, including ones the caller cannot read.
func TestGetAllRecipeMetadata_FiltersToVisibleRecipes(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewSearchRepository(db)
	requesterID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT "recipe_id" FROM "recipe_difficulty_levels" WHERE recipe_id IN (SELECT recipe_id FROM recipe_accesses WHERE (owner_id = $1 OR visibility = $2) AND NOT deleted)`)).
		WithArgs(requesterID.String(), "public").
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id"}).AddRow(id))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT recipe_id, category_id as value_id FROM "recipe_categories" WHERE recipe_id IN ($1)`)).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "value_id"}).AddRow(id, id))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT recipe_id, tag_id as value_id FROM "recipe_tags" WHERE recipe_id IN ($1)`)).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "value_id"}).AddRow(id, id))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT recipe_id, difficulty_level_id as value_id FROM "recipe_difficulty_levels" WHERE recipe_id IN ($1)`)).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "value_id"}).AddRow(id, id))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT recipe_id, preparation_time_id as value_id FROM "recipe_preparation_times" WHERE recipe_id IN ($1)`)).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "value_id"}).AddRow(id, id))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT recipe_id, cuisine_type_id as value_id FROM "recipe_cuisine_types" WHERE recipe_id IN ($1)`)).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "value_id"}).AddRow(id, id))

	result, err := r.GetAllRecipeMetadata(requesterID)

	assert.NoError(t, err)
	assert.Len(t, *result, 1)
	assert.Equal(t, id, (*result)[0].RecipeID)
}

func TestGetAllRecipeMetadata_Empty(t *testing.T) {
	db, mock := co.NewMockDatabase(t)
	r := NewSearchRepository(db)
	requesterID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT "recipe_id" FROM "recipe_difficulty_levels" WHERE recipe_id IN (SELECT recipe_id FROM recipe_accesses WHERE (owner_id = $1 OR visibility = $2) AND NOT deleted)`)).
		WithArgs(requesterID.String(), "public").
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id"}))

	result, err := r.GetAllRecipeMetadata(requesterID)

	assert.NoError(t, err)
	assert.Len(t, *result, 0)
}
