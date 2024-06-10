package repositories

import (
	"database/sql/driver"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	m "metadata-service/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	category m.Category = m.Category{
		CategoryName: "category",
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

// ========================================================================================================

func TestCategoryFindAll_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "categorys" WHERE "categorys"."deleted_at" IS NULL]`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := r.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestCategoryFindAll_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "categorys" WHERE "categorys"."deleted_at" IS NULL]`).
		WillReturnError(errors.New("error"))

	result, err := r.FindAll()

	assert.Error(t, err)
	assert.Len(t, result, 0)
	assert.EqualError(t, err, "error")
}

func TestCategoryFindSingle_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "categorys" WHERE category_id = $1 AND "categorys"."deleted_at" IS NULL]`).
		WithArgs(
			1,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := r.FindSingle(uint(1))

	assert.NoError(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.Equal(t, uint(1), result.ID)
}

func TestCategoryFindSingle_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "categorys" WHERE category_id = $1 AND "categorys"."deleted_at" IS NULL]`).
		WithArgs(
			1,
		).
		WillReturnError(errors.New("error"))

	_, err := r.FindSingle(uint(1))

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestCategoryCreate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`[INSERT INTO "categorys" ("created_at","updated_at","deleted_at","category_name") VALUES ($1,$2,$3,$4) RETURNING "id"]`).
		WithArgs(
			AnyTime{},
			AnyTime{},
			nil,
			"category",
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	result, err := r.Create(category)

	assert.NoError(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.Equal(t, uint(1), result.ID)

}

func TestCategoryCreate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`[INSERT INTO "categorys" ("created_at","updated_at","deleted_at","category_name") VALUES ($1,$2,$3,$4) RETURNING "id"]`).
		WithArgs(
			AnyTime{},
			AnyTime{},
			nil,
			"category",
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Create(category)

	assert.Error(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.EqualError(t, err, "error")

}

func TestCategoryUpdate_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "categorys" SET "updated_at"=$1,"category_name"=$2 WHERE "id" = $3]`).
		WithArgs(
			AnyTime{},
			"category",
			1,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updateCategory := category
	updateCategory.ID = 1

	result, err := r.Update(updateCategory)

	assert.NoError(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.Equal(t, uint(1), result.ID)

}

func TestCategoryUpdate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "categorys" SET "updated_at"=$1,"category_name"=$2 WHERE "id" = $3]`).
		WithArgs(
			AnyTime{},
			"category",
			1,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	updateCategory := category
	updateCategory.ID = 1

	result, err := r.Update(updateCategory)

	assert.Error(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.EqualError(t, err, "error")

}

func TestCategoryDelete_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "categorys" SET "deleted_at"=$1 WHERE category_id = $2 AND "categorys"."deleted_at" IS NULL]`).
		WithArgs(
			AnyTime{},
			1,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	deleteCategory := category
	deleteCategory.ID = 1

	err := r.Delete(deleteCategory)

	assert.NoError(t, err)
}

func TestCategoryDelete_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewCategoryRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "categorys" SET "updated_at"=$1,"category_name"=$2 WHERE "id" = $3]`).
		WithArgs(
			AnyTime{},
			1,
		).WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	deleteCategory := category
	deleteCategory.ID = 1

	err := r.Delete(deleteCategory)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")

}
