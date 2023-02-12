package repositories

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	m "github.com/ihulsbus/cookbook/internal/models"
	"github.com/stretchr/testify/assert"
)

var (
	recipe m.Recipe = m.Recipe{
		RecipeName:      "recipe",
		Description:     "description",
		DifficultyLevel: 1,
		CookingTime:     1,
		ServingCount:    1,
	}
	instruction m.Instruction = m.Instruction{
		RecipeID:    1,
		StepNumber:  1,
		Description: "instruction",
	}
)

func TestRecipeFindAll_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "recipe_category" WHERE "recipe_category"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(`[SELECT * FROM "recipe_ingredient" WHERE "recipe_ingredient"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(`[SELECT * FROM "recipe_tag" WHERE "recipe_tag"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(`[SELECT * FROM "recipes" WHERE "recipes"."deleted_at" IS NULL]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := r.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestRecipeFindAll_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "recipe_category" WHERE "recipe_category"."recipe_id" = 1]`).WillReturnError(errors.New("error"))
	result, err := r.FindAll()

	assert.Error(t, err)
	assert.Len(t, result, 0)
}

func TestRecipeFindSingle_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "recipe_category" WHERE "recipe_category"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(`[SELECT * FROM "recipe_ingredient" WHERE "recipe_ingredient"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(`[SELECT * FROM "recipe_tag" WHERE "recipe_tag"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(`[SELECT * FROM "recipes" WHERE "recipes"."deleted_at" IS NULL AND "recipes"."id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := r.FindSingle(1)

	assert.NoError(t, err)
	assert.IsType(t, result, m.Recipe{})
}

func TestRecipeFindSingle_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "recipe_category" WHERE "recipe_category"."recipe_id" = 1]`).WillReturnError(errors.New("error"))
	result, err := r.FindSingle(1)

	assert.Error(t, err)
	assert.IsType(t, result, m.Recipe{})
}

func TestRecipeFindInstruction_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "instructions" WHERE "instructions"."recipe_id" = 1]`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := r.FindInstruction(1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)

}

func TestRecipeFindInstruction_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "instructions" WHERE "instructions"."recipe_id" = 1]`).
		WillReturnError(errors.New("error"))

	_, err := r.FindInstruction(1)

	assert.Error(t, err)
}

func TestRecipeCreateInstruction_OK(t *testing.T) {
	var createInstructions []m.Instruction
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`[INSERT INTO "instructions" ("created_at","updated_at","deleted_at","recipe_id","step_number","description") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"]`).
		WithArgs(
			AnyTime{},
			AnyTime{},
			nil,
			1,
			1,
			"instruction",
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	createInstructions = append(createInstructions, instruction)

	result, err := r.CreateInstruction(createInstructions)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestRecipeCreateInstruction_Err(t *testing.T) {
	var createInstructions []m.Instruction
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`[INSERT INTO "instructions" ("created_at","updated_at","deleted_at","recipe_id","step_number","description") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"]`).
		WithArgs(
			AnyTime{},
			AnyTime{},
			nil,
			1,
			1,
			"instruction",
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	createInstructions = append(createInstructions, instruction)

	result, err := r.CreateInstruction(createInstructions)

	assert.Error(t, err)
	assert.Len(t, result, 0)
}

func TestRecipeCreate_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`[INSERT INTO "recipes" ("created_at","updated_at","deleted_at","recipe_name","description","difficulty_level","cooking_time","serving_count","image_name") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"]`).
		WithArgs(
			AnyTime{},
			AnyTime{},
			nil,
			"recipe",
			"description",
			1,
			1,
			1,
			"",
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	result, err := r.Create(recipe)

	assert.NoError(t, err)
	assert.IsType(t, result, m.Recipe{})
}

func TestRecipeCreate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`[INSERT INTO "recipes" ("created_at","updated_at","deleted_at","recipe_name","description","difficulty_level","cooking_time","serving_count","image_name") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"]`).
		WithArgs(
			AnyTime{},
			AnyTime{},
			nil,
			"recipe",
			"description",
			1,
			1,
			1,
			"",
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()
	result, err := r.Create(recipe)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, m.Recipe{})
}

func TestRecipeUpdate_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	updateRecipe := recipe
	updateRecipe.ID = 1

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "recipes" SET "updated_at"=$1,"recipe_name"=$2,"description"=$3,"difficulty_level"=$4,"cooking_time"=$5,"serving_count"=$6 WHERE "id" = $7]`).
		WithArgs(
			AnyTime{},
			"recipe",
			"description",
			1,
			1,
			1,
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	result, err := r.Update(updateRecipe)

	assert.NoError(t, err)
	assert.IsType(t, result, m.Recipe{})
}

func TestRecipeUpdate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	updateRecipe := recipe
	updateRecipe.ID = 1

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "recipes" SET "updated_at"=$1,"recipe_name"=$2,"description"=$3,"difficulty_level"=$4,"cooking_time"=$5,"serving_count"=$6 WHERE "id" = $7]`).
		WithArgs(
			AnyTime{},
			"recipe",
			"description",
			1,
			1,
			1,
			1,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectCommit()

	result, err := r.Update(updateRecipe)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, m.Recipe{})
}

func TestRecipeDelete_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	updateRecipe := recipe
	updateRecipe.ID = 1

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "recipes" SET "deleted_at"=$1 WHERE "recipes"."id" = $2 AND "recipes"."deleted_at" IS NULL]`).
		WithArgs(
			AnyTime{},
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.Delete(updateRecipe)

	assert.NoError(t, err)
}

func TestRecipeDelete_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	updateRecipe := recipe
	updateRecipe.ID = 1

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "recipes" SET "deleted_at"=$1 WHERE "recipes"."id" = $2 AND "recipes"."deleted_at" IS NULL]`).
		WithArgs(
			AnyTime{},
			1,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectCommit()

	err := r.Delete(updateRecipe)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
