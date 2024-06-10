package repositories

import (
	"database/sql/driver"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"recipe-service/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	recipe models.Recipe = models.Recipe{
		RecipeName:      "recipe",
		Description:     "description",
		DifficultyLevel: 1,
		CookingTime:     1,
		ServingCount:    1,
	}
)

func newMockDatabase(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {

	var mockDB *gorm.DB

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	sqlMockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sql mock init failed: %v", err.Error())
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlMockDB,
		PreferSimpleProtocol: true,
	})

	mockDB, err = gorm.Open(dialector, &gorm.Config{
		NowFunc: timeFunc,
		Logger:  newLogger,
	})
	if err != nil {
		t.Fatalf("gorm mock init failed: %v", err.Error())
	}

	return mockDB, mock
}

func timeFunc() time.Time {
	time, _ := time.Parse("2006-01-02 15:04", "2023-02-04 18:00")
	return time
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestRecipeFindAll_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "recipe_category" WHERE "recipe_category"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(`[SELECT * FROM "recipe_ingredient" WHERE "recipe_ingredient"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(`[SELECT * FROM "recipe_tag" WHERE "recipe_tag"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(`[SELECT * FROM "recipes" WHERE "recipes"."deleted_at" IS NULL]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := r.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestRecipeFindAll_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "recipe_category" WHERE "recipe_category"."recipe_id" = 1]`).WillReturnError(errors.New("error"))
	result, err := r.FindAll()

	assert.Error(t, err)
	assert.Len(t, result, 0)
}

func TestRecipeFindAll_NotFoundErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "recipe_category" WHERE "recipe_category"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
	mock.ExpectQuery(`[SELECT * FROM "recipe_ingredient" WHERE "recipe_ingredient"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
	mock.ExpectQuery(`[SELECT * FROM "recipe_tag" WHERE "recipe_tag"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
	mock.ExpectQuery(`[SELECT * FROM "recipes" WHERE "recipes"."deleted_at" IS NULL]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))
	result, err := r.FindAll()

	assert.Error(t, err)
	assert.Len(t, result, 0)
}

func TestRecipeFindSingle_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "recipe_category" WHERE "recipe_category"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"recipe_id"}).AddRow(1))
	mock.ExpectQuery(`[SELECT * FROM "recipe_ingredient" WHERE "recipe_ingredient"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"recipe_id"}).AddRow(1))
	mock.ExpectQuery(`[SELECT * FROM "recipe_tag" WHERE "recipe_tag"."recipe_id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"recipe_id"}).AddRow(1))
	mock.ExpectQuery(`[SELECT * FROM "recipes" WHERE "recipes"."deleted_at" IS NULL AND "recipes"."id" = 1]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := r.FindSingle(1)

	assert.NoError(t, err)
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeFindSingle_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "recipe_category" WHERE "recipe_category"."recipe_id" = 1]`).WillReturnError(errors.New("error"))
	result, err := r.FindSingle(1)

	assert.Error(t, err)
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeCreate_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`[INSERT INTO "recipes" ("created_at","updated_at","deleted_at","recipe_name","description","difficulty_level","cooking_time","serving_count","image_name","author_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"]`).
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
			"",
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	result, err := r.Create(recipe)

	assert.NoError(t, err)
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeCreate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`[INSERT INTO "recipes" ("created_at","updated_at","deleted_at","recipe_name","description","difficulty_level","cooking_time","serving_count","image_name","author_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"]`).
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
			"",
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()
	result, err := r.Create(recipe)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, models.Recipe{})
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
	assert.IsType(t, result, models.Recipe{})
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
	assert.IsType(t, result, models.Recipe{})
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
