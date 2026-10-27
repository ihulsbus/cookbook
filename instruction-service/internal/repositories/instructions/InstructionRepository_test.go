package repositories

import (
	"errors"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	instruction = m.Instruction{
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

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "instructions" WHERE entity_id = $1 AND "instructions"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sequence", "description", "media_id"}).
			AddRow(
				instruction.ID,
				instruction.Sequence,
				instruction.Description,
				instruction.MediaID,
			))

	result, err := r.Find(instruction.EntityID)

	assert.NoError(t, err)
	derefResult := *result

	assert.Len(t, derefResult, 1)
	assert.IsType(t, &[]m.Instruction{}, result)
	assert.Equal(t, instruction.ID, derefResult[0].ID)
	assert.Equal(t, instruction.Sequence, derefResult[0].Sequence)
	assert.Equal(t, instruction.Description, derefResult[0].Description)
	assert.Equal(t, instruction.MediaID, derefResult[0].MediaID)
}

func TestFindInstruction_NotFoundErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewInstructionRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "instructions" WHERE entity_id = $1 AND "instructions"."deleted_at" IS NULL`)).
		WillReturnRows(&sqlmock.Rows{})

	_, err := r.Find(instruction.EntityID)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
}

func TestFindInstruction_FindErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewInstructionRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "instructions" WHERE entity_id = $1 AND "instructions"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("error"))

	_, err := r.Find(instruction.EntityID)

	assert.Error(t, err)
}

func TestCreateInstruction_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewInstructionRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "instructions" ("sequence","description","media_id","entity_id","entity_type","created_at","updated_at","deleted_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).
		WithArgs(
			instruction.Sequence,
			instruction.Description,
			instruction.MediaID,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(instruction.ID))
	mock.ExpectCommit()

	var createRequest []m.Instruction
	createRequest = append(createRequest, instruction)
	result, err := r.Create(&createRequest)

	assert.NoError(t, err)
	assert.IsType(t, &[]m.Instruction{}, result)
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

	var createRequest []m.Instruction
	createRequest = append(createRequest, instruction)
	result, err := r.Create(&createRequest)

	assert.Error(t, err)
	assert.IsType(t, &[]m.Instruction{}, result)
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

	var deleteRequest []m.Instruction
	deleteRequest = append(deleteRequest, instruction)
	err := r.Delete(&deleteRequest)

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

	var deleteRequest []m.Instruction
	deleteRequest = append(deleteRequest, instruction)
	err := r.Delete(&deleteRequest)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
