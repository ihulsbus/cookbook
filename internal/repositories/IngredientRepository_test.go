package repositories

import (
	"database/sql/driver"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	m "github.com/ihulsbus/cookbook/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ingredient = m.Ingredient{}
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

	mock.ExpectQuery(`[SELECT * FROM ingredients]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := r.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestIngredientFindAll_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectQuery(`[SELECT * FROM ingredients]`).
		WillReturnError(errors.New("error"))

	result, err := r.FindAll()

	assert.EqualError(t, err, "error")
	assert.Len(t, result, 0)
}

func TestIngredientFindUnits_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectQuery(`[SELECT * FROM units]`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := r.FindUnits()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestIngredientFindUnits_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectQuery(`[SELECT * FROM units]`).
		WillReturnError(errors.New("error"))

	result, err := r.FindUnits()

	assert.EqualError(t, err, "error")
	assert.Len(t, result, 0)
}

func TestIngredientFindSingle_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)
	expectedIngredient := ingredient
	expectedIngredient.ID = 1

	mock.ExpectQuery(`[SELECT * FROM "ingredients" WHERE id = 1 AND "ingredients"."deleted_at" IS NULL]`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := r.FindSingle(1)

	assert.NoError(t, err)
	assert.Equal(t, result, expectedIngredient)
}

func TestIngredientFindSingle_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)
	// expectedIngredient := ingredient

	mock.ExpectQuery(`[SELECT * FROM "ingredients" WHERE id = 1 AND "ingredients"."deleted_at" IS NULL]`).
		WillReturnError(errors.New("error"))

	result, err := r.FindSingle(1)

	assert.EqualError(t, err, "error")
	assert.Equal(t, result, ingredient)
}

func TestIngredientCreate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)
	expectedIngredient := ingredient

	mock.ExpectBegin()
	mock.ExpectQuery(`[INSERT INTO "ingredients"("created_at","updated_at","deleted_at","ingredient_name") VALUES ($1,$2,$3,$4) RETURNING "id"]`).
		WithArgs(AnyTime{}, AnyTime{}, nil, "").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	_, err := r.Create(ingredient)

	expectedIngredient.ID = 1

	assert.NoError(t, err)
	// assert.Equal(t, result, expectedIngredient)
}

func TestIngredientCreate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`[INSERT INTO "ingredients"("created_at","updated_at","deleted_at","ingredient_name")]`).
		WithArgs(AnyTime{}, AnyTime{}, nil, "").WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	_, err := r.Create(ingredient)

	assert.EqualError(t, err, "error")
}

func TestIngredientUpdate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)
	expectedIngredient := ingredient

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "ingredients" SET "updated_at"=$1 WHERE ID = $2 AND "ingredients"."deleted_at" IS NULL]`).
		WithArgs(
			AnyTime{},
			0,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	_, err := r.Update(ingredient)

	expectedIngredient.ID = 1

	assert.NoError(t, err)
	// assert.Equal(t, result, expectedIngredient)
}

func TestIngredientUpdate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "ingredients" SET "updated_at"=$1 WHERE ID = $2 AND "ingredients"."deleted_at" IS NULL]`).
		WithArgs(
			AnyTime{},
			0,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	_, err := r.Update(ingredient)

	assert.EqualError(t, err, "error")
}

func TestIngredientDelete_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)
	deleteIngredient := ingredient
	deleteIngredient.ID = 1

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "ingredients" SET "deleted_at"=$1 WHERE "ingredients"."id" = $2 AND "ingredients"."deleted_at" IS NULL]`).
		WithArgs(
			AnyTime{},
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.Delete(deleteIngredient)

	assert.NoError(t, err)
}

func TestIngredientDelete_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewIngredientRepository(db)
	deleteIngredient := ingredient
	deleteIngredient.ID = 1

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "ingredients" SET "deleted_at"=$1 WHERE "ingredients"."id" = $2 AND "ingredients"."deleted_at" IS NULL]`).
		WithArgs(
			AnyTime{},
			1,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	err := r.Delete(deleteIngredient)

	assert.EqualError(t, err, "error")
}
