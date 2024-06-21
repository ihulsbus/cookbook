package repositories

import (
	"database/sql/driver"
	"errors"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"recipe-service/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	recipe models.Recipe = models.Recipe{
		ID:           uuid.New(),
		Name:         "recipe",
		Description:  "description",
		ServingCount: 1,
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

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" WHERE "recipes"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "servingcount"}).
			AddRow(
				recipe.ID,
				recipe.Name,
				recipe.Description,
				recipe.ServingCount,
			))

	result, err := r.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, recipe.ID, result[0].ID)
	assert.Equal(t, recipe.Name, result[0].Name)
	assert.Equal(t, recipe.Description, result[0].Description)
}

func TestRecipeFindAll_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" WHERE "recipes"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("error"))
	result, err := r.FindAll()

	assert.Error(t, err)
	assert.Len(t, result, 0)
}

func TestRecipeFindAll_NotFoundErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" WHERE "recipes"."deleted_at" IS NULL`)).
		WillReturnRows(&sqlmock.Rows{})
	result, err := r.FindAll()

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.Len(t, result, 0)
}

func TestRecipeFindSingle_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" WHERE "recipes"."deleted_at" IS NULL AND "recipes"."id" = $1 ORDER BY "recipes"."id" LIMIT $2`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(recipe.ID))

	result, err := r.FindSingle(recipe)

	assert.NoError(t, err)
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeFindSingle_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_category" WHERE "recipe_category"."recipe_id" = 1`)).
		WillReturnError(errors.New("error"))
	result, err := r.FindSingle(recipe)

	assert.Error(t, err)
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeCreate_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "recipes" ("created_at","updated_at","deleted_at","name","description","serving_count","id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			recipe.Name,
			recipe.Description,
			recipe.ServingCount,
			recipe.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(recipe.ID))
	mock.ExpectCommit()

	result, err := r.Create(recipe)

	assert.NoError(t, err)
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeCreate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "recipes" ("created_at","updated_at","deleted_at","name","description","serving_count","id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			recipe.Name,
			recipe.Description,
			recipe.ServingCount,
			recipe.ID,
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

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipes" SET "updated_at"=$1,"name"=$2,"description"=$3,"serving_count"=$4 WHERE "recipes"."deleted_at" IS NULL AND "id" = $5`)).
		WithArgs(
			sqlmock.AnyArg(),
			recipe.Name,
			recipe.Description,
			recipe.ServingCount,
			recipe.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	result, err := r.Update(recipe)

	assert.NoError(t, err)
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeUpdate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipes" SET "updated_at"=$1,"name"=$2,"description"=$3,"serving_count"=$4 WHERE "recipes"."deleted_at" IS NULL AND "id" = $5`)).
		WithArgs(
			sqlmock.AnyArg(),
			recipe.Name,
			recipe.Description,
			recipe.ServingCount,
			recipe.ID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectCommit()

	result, err := r.Update(recipe)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeDelete_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipes" SET "deleted_at"=$1 WHERE "recipes"."id" = $2 AND "recipes"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			recipe.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.Delete(recipe)

	assert.NoError(t, err)
}

func TestRecipeDelete_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipes" SET "deleted_at"=$1 WHERE "recipes"."id" = $2 AND "recipes"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			recipe.ID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectCommit()

	err := r.Delete(recipe)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
