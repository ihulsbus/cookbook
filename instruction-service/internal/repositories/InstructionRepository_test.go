package repositories

import (
	"database/sql/driver"
	"errors"
	m "instruction-service/internal/models"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	instruction m.Instruction = m.Instruction{
		Description: "instruction",
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

func TestFindInstruction_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "instructions" WHERE "instructions"."recipe_id" = 1]`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := r.FindInstruction(uuid.UUID{})

	assert.NoError(t, err)
	assert.IsType(t, m.Instruction{}, result)

}

func TestFindInstruction_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectQuery(`[SELECT * FROM "instructions" WHERE "instructions"."recipe_id" = 1]`).
		WillReturnError(errors.New("error"))

	_, err := r.FindInstruction(uuid.UUID{})

	assert.Error(t, err)
}

func TestCreateInstruction_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`[INSERT INTO "instructions" ("created_at","updated_at","deleted_at","recipe_id","description") VALUES ($1,$2,$3,$4,$6) RETURNING "id"]`).
		WithArgs(
			AnyTime{},
			AnyTime{},
			nil,
			1,
			"instruction",
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	result, err := r.CreateInstruction(instruction)

	assert.NoError(t, err)
	assert.IsType(t, m.Instruction{}, result)
}

func TestCreateInstruction_Err(t *testing.T) {
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

	result, err := r.CreateInstruction(instruction)

	assert.Error(t, err)
	assert.IsType(t, m.Instruction{}, result)
}

func TestUpdateInstruction_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "instructions" SET "updated_at"=$1,"recipe_id"=$2,"description"=$3 WHERE recipe_id = $4 AND "instructions"."deleted_at" IS NULL]`).
		WithArgs(
			AnyTime{},
			1,
			"instruction",
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := r.UpdateInstruction(instruction)

	assert.NoError(t, err)
	assert.IsType(t, m.Instruction{}, result)
}

func TestUpdateInstruction_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "instructions" SET "updated_at"=$1,"recipe_id"=$2,"description"=$3 WHERE recipe_id = $4 AND "instructions"."deleted_at" IS NULL]`).
		WithArgs(
			AnyTime{},
			1,
			"instruction",
			1,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.UpdateInstruction(instruction)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, m.Instruction{}, result)
}

func TestDeleteInstruction_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "instructions" SET "deleted_at"=$1 WHERE "instructions"."recipe_id" = $2 AND "instructions"."deleted_at" IS NULL]`).
		WithArgs(
			AnyTime{},
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.DeleteInstruction(instruction)

	assert.NoError(t, err)
}

func TestDeleteInstruction_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewRecipeRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`[UPDATE "instructions" SET "deleted_at"=$1 WHERE "instructions"."recipe_id" = $2 AND "instructions"."deleted_at" IS NULL]`).
		WithArgs(
			AnyTime{},
			1,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectCommit()

	err := r.DeleteInstruction(instruction)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
