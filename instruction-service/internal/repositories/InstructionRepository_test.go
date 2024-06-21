package repositories

import (
	"errors"
	m "instruction-service/internal/models"
	"log"
	"os"
	"regexp"
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
		ID:          uuid.New(),
		Sequence:    1,
		Description: "instruction",
		MediaID:     uuid.New(),
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

// ========================================================================================================

func TestFindInstruction_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewInstructionRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "instructions" WHERE "instructions"."deleted_at" IS NULL AND "instructions"."id" = $1 ORDER BY "instructions"."id" LIMIT $2`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sequence", "description", "media_id"}).
			AddRow(
				instruction.ID,
				instruction.Sequence,
				instruction.Description,
				instruction.MediaID,
			))

	result, err := r.Find(instruction)

	assert.NoError(t, err)
	assert.IsType(t, m.Instruction{}, result)
	assert.Equal(t, instruction.ID, result.ID)
	assert.Equal(t, instruction.Sequence, result.Sequence)
	assert.Equal(t, instruction.Description, result.Description)
	assert.Equal(t, instruction.MediaID, result.MediaID)
}

func TestFindInstruction_NotFoundErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewInstructionRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "instructions" WHERE "instructions"."deleted_at" IS NULL AND "instructions"."id" = $1 ORDER BY "instructions"."id" LIMIT $2`)).
		WillReturnRows(&sqlmock.Rows{})

	_, err := r.Find(instruction)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
}

func TestFindInstruction_FindErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewInstructionRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "instructions" WHERE "instructions"."deleted_at" IS NULL AND "instructions"."id" = $1 ORDER BY "instructions"."id" LIMIT $2`)).
		WillReturnError(errors.New("error"))

	_, err := r.Find(instruction)

	assert.Error(t, err)
}

func TestCreateInstruction_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewInstructionRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "instructions" ("sequence","description","media_id","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(
			instruction.Sequence,
			instruction.Description,
			instruction.MediaID,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(instruction.ID))
	mock.ExpectCommit()

	result, err := r.Create(instruction)

	assert.NoError(t, err)
	assert.IsType(t, m.Instruction{}, result)
}

func TestCreateInstruction_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewInstructionRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "instructions" ("created_at","updated_at","deleted_at","recipe_id","step_number","description") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			1,
			1,
			"instruction",
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Create(instruction)

	assert.Error(t, err)
	assert.IsType(t, m.Instruction{}, result)
}

func TestUpdateInstruction_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewInstructionRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "instructions" SET "sequence"=$1,"description"=$2,"media_id"=$3,"updated_at"=$4 WHERE "instructions"."deleted_at" IS NULL AND "id" = $5`)).
		WithArgs(
			instruction.Sequence,
			instruction.Description,
			instruction.MediaID,
			sqlmock.AnyArg(),
			instruction.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := r.Update(instruction)

	assert.NoError(t, err)
	assert.IsType(t, m.Instruction{}, result)
}

func TestUpdateInstruction_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewInstructionRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "instructions" SET "sequence"=$1,"description"=$2,"media_id"=$3,"updated_at"=$4 WHERE "instructions"."deleted_at" IS NULL AND "id" = $5`)).
		WithArgs(
			instruction.Sequence,
			instruction.Description,
			instruction.MediaID,
			sqlmock.AnyArg(),
			instruction.ID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Update(instruction)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, m.Instruction{}, result)
}

func TestDeleteInstruction_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewInstructionRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "instructions" SET "deleted_at"=$1 WHERE "instructions"."id" = $2 AND "instructions"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			instruction.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := r.Delete(instruction)

	assert.NoError(t, err)
}

func TestDeleteInstruction_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewInstructionRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "instructions" SET "deleted_at"=$1 WHERE "instructions"."id" = $2 AND "instructions"."deleted_at" IS NULL`)).
		WithArgs(
			sqlmock.AnyArg(),
			instruction.ID,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectCommit()

	err := r.Delete(instruction)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
