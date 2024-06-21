package repositories

import (
	"database/sql/driver"
	"errors"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	m "ingredient-service/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ingredient m.Ingredient = m.Ingredient{
		ID:   uuid.New(),
		Name: "ingredient",
	}
	unit m.Unit = m.Unit{
		ID:        uuid.New(),
		FullName:  "unit",
		ShortName: "u",
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

func TestIngredientFindAll_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ingredients" WHERE "ingredients"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(
				ingredient.ID,
				ingredient.Name,
			))

	result, err := r.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestIngredientFindAll_NotFoundErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ingredients" WHERE "ingredients"."deleted_at" IS NULL`)).
		WillReturnRows(&sqlmock.Rows{})

	result, err := r.FindAll()

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.Len(t, result, 0)
}

func TestIngredientFindAll_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ingredients" WHERE "ingredients"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("error"))

	result, err := r.FindAll()

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.Len(t, result, 0)
}

func TestIngredientFindUnits_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "units" WHERE "units"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "FullName", "ShortName"}).
			AddRow(
				unit.ID,
				unit.FullName,
				unit.ShortName,
			))

	result, err := r.FindUnits()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestIngredientFindUnits_NotFoundErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "units" WHERE "units"."deleted_at" IS NULL`)).
		WillReturnRows(&sqlmock.Rows{})

	result, err := r.FindUnits()

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.Len(t, result, 0)
}

func TestIngredientFindUnits_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "units" WHERE "units"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("error"))

	result, err := r.FindUnits()

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.Len(t, result, 0)
}

func TestIngredientFindSingle_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ingredients" WHERE "ingredients"."deleted_at" IS NULL AND "ingredients"."id" = $1 ORDER BY "ingredients"."id" LIMIT $2`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(
				ingredient.ID,
				ingredient.Name,
			))

	result, err := r.FindSingle(ingredient)

	assert.NoError(t, err)
	assert.Equal(t, ingredient, result)
}

func TestIngredientFindSingle_NotFoundErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ingredients" WHERE "ingredients"."deleted_at" IS NULL AND "ingredients"."id" = $1 ORDER BY "ingredients"."id" LIMIT $2`)).
		WillReturnRows(&sqlmock.Rows{})

	result, err := r.FindSingle(ingredient)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.IsType(t, m.Ingredient{}, result)
}

func TestIngredientFindSingle_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ingredients" WHERE "ingredients"."deleted_at" IS NULL AND "ingredients"."id" = $1 ORDER BY "ingredients"."id" LIMIT $2`)).
		WillReturnError(errors.New("error"))

	result, err := r.FindSingle(ingredient)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.Equal(t, m.Ingredient{}, result)
}

func TestIngredientCreate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ingredients" ("name","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			ingredient.Name,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(ingredient.ID))
	mock.ExpectCommit()
	result, err := r.Create(ingredient)

	assert.NoError(t, err)
	assert.IsType(t, m.Ingredient{}, result)
}

func TestIngredientCreate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ingredients" ("name","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(
			ingredient.Name,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			ingredient.ID,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	_, err := r.Create(ingredient)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestIngredientUpdate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ingredients" SET "name"=$1,"updated_at"=$2 WHERE id = $3 AND "ingredients"."deleted_at" IS NULL AND "id" = $4`)).
		WithArgs(
			ingredient.Name,
			sqlmock.AnyArg(),
			ingredient.ID,
			ingredient.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	result, err := r.Update(ingredient)

	assert.NoError(t, err)
	assert.IsType(t, m.Ingredient{}, result)
}

func TestIngredientUpdate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ingredients" SET "name"=$1,"updated_at"=$2 WHERE id = $3 AND "ingredients"."deleted_at" IS NULL AND "id" = $4`)).
		WithArgs(
			ingredient.Name,
			sqlmock.AnyArg(),
			ingredient.ID,
			ingredient.ID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	_, err := r.Update(ingredient)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestIngredientDelete_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ingredients" SET "deleted_at"=$1 WHERE "ingredients"."id" = $2 AND "ingredients"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			ingredient.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.Delete(ingredient)

	assert.NoError(t, err)
}

func TestIngredientDelete_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ingredients" SET "deleted_at"=$1 WHERE "ingredients"."id" = $2 AND "ingredients"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			ingredient.ID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	err := r.Delete(ingredient)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
