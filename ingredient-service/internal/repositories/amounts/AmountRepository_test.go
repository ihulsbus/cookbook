package repositories

import (
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
	amount m.Amount = m.Amount{
		RecipeID:     uuid.New(),
		IngredientID: uuid.New(),
		Quantity:     1,
		UnitID:       uuid.New(),
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
	timestamp, _ := time.Parse("2006-01-02 15:04", "2023-02-04 18:00")
	return timestamp
}

func TestAmountFind_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewAmountRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "amounts" WHERE recipe_id = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "ingredient_id", "quantity", "unit_id"}).
			AddRow(
				amount.RecipeID,
				amount.IngredientID,
				amount.Quantity,
				amount.UnitID,
			))

	result, err := r.Find(amount.RecipeID)

	var expectedResponse []m.Amount
	expectedResponse = append(expectedResponse, amount)

	assert.NoError(t, err)
	assert.Equal(t, &expectedResponse, result)
}

func TestAmountFind_NotFoundErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewAmountRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "amounts" WHERE recipe_id = $1`)).
		WillReturnError(errors.New("not found"))

	result, err := r.Find(amount.RecipeID)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	assert.IsType(t, &[]m.Amount{}, result)
}

func TestAmountFind_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewAmountRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "amounts" WHERE recipe_id = $1`)).
		WillReturnError(errors.New("error"))

	_, err := r.Find(amount.RecipeID)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestAmountCreate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewAmountRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "amounts" ("recipe_id","ingredient_id","quantity","unit_id") VALUES ($1,$2,$3,$4)`)).
		WithArgs(
			amount.RecipeID,
			amount.IngredientID,
			amount.Quantity,
			amount.UnitID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	var createRequest []m.Amount
	createRequest = append(createRequest, amount)
	result, err := r.Create(&createRequest)

	assert.NoError(t, err)
	assert.IsType(t, &[]m.Amount{}, result)
}

func TestAmountCreate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewAmountRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "amounts" ("recipe_id","ingredient_id","quantity","unit_id") VALUES ($1,$2,$3,$4)`)).
		WithArgs(
			amount.RecipeID,
			amount.IngredientID,
			amount.Quantity,
			amount.UnitID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	var createRequest []m.Amount
	createRequest = append(createRequest, amount)
	_, err := r.Create(&createRequest)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestAmountUpdate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewAmountRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "amounts" SET "quantity"=$1,"unit_id"=$2 WHERE "recipe_id" = $3 AND "ingredient_id" = $4`)).
		WithArgs(
			amount.Quantity,
			amount.UnitID,
			amount.RecipeID,
			amount.IngredientID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	var updateRequest []m.Amount
	updateRequest = append(updateRequest, amount)
	result, err := r.Update(&updateRequest)

	assert.NoError(t, err)
	assert.IsType(t, &[]m.Amount{}, result)
}

func TestAmountUpdate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewAmountRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "amounts" ("recipe_id","ingredient_id","quantity","unit_id") VALUES ($1,$2,$3,$4)`)).
		WithArgs(
			amount.RecipeID,
			amount.IngredientID,
			amount.Quantity,
			amount.UnitID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	var updateRequest []m.Amount
	updateRequest = append(updateRequest, amount)
	_, err := r.Update(&updateRequest)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestAmountDelete_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewAmountRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "amounts" WHERE ("amounts"."recipe_id","amounts"."ingredient_id") IN (($1,$2))`)).
		WithArgs(
			amount.RecipeID,
			amount.IngredientID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	var deleteRequest []m.Amount
	deleteRequest = append(deleteRequest, amount)
	err := r.Delete(&deleteRequest)

	assert.NoError(t, err)
}

func TestAmountDelete_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewAmountRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "amounts" WHERE ("amounts"."recipe_id","amounts"."ingredient_id") IN (($1,$2))`)).
		WithArgs(
			amount.RecipeID,
			amount.IngredientID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	var deleteRequest []m.Amount
	deleteRequest = append(deleteRequest, amount)
	err := r.Delete(&deleteRequest)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
