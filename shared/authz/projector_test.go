package authz

import (
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func newMockDatabase(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	sqlMockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sql mock init failed: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlMockDB,
		PreferSimpleProtocol: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{Logger: newLogger})
	if err != nil {
		t.Fatalf("gorm mock init failed: %v", err)
	}

	return db, mock
}

func TestProjectorHandle_CreatesWhenMissing(t *testing.T) {
	db, mock := newMockDatabase(t)
	p := NewProjector(db)

	recipeID := uuid.New()
	event := rmq.RecipeAccessState{
		RecipeID:   recipeID,
		OwnerID:    "owner-1",
		Visibility: "private",
		Version:    1000,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_accesses" WHERE recipe_id = $1`)).
		WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "recipe_accesses"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := p.Handle(event)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectorHandle_AppliesNewerVersion(t *testing.T) {
	db, mock := newMockDatabase(t)
	p := NewProjector(db)

	recipeID := uuid.New()
	event := rmq.RecipeAccessState{
		RecipeID:   recipeID,
		OwnerID:    "owner-1",
		Visibility: "public",
		Version:    2000,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_accesses" WHERE recipe_id = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "owner_id", "visibility", "version", "deleted"}).
			AddRow(recipeID, "owner-1", "private", 1000, false))
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipe_accesses" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := p.Handle(event)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestProjectorHandle_DropsOutOfOrderEvent checks that an event whose
// Version is not strictly greater than what's stored is a no-op (no
// UPDATE), so a duplicate or late redelivery never regresses the state.
func TestProjectorHandle_DropsOutOfOrderEvent(t *testing.T) {
	db, mock := newMockDatabase(t)
	p := NewProjector(db)

	recipeID := uuid.New()
	staleEvent := rmq.RecipeAccessState{
		RecipeID:   recipeID,
		OwnerID:    "owner-1",
		Visibility: "public",
		Version:    500, // older than what's stored (1000)
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_accesses" WHERE recipe_id = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "owner_id", "visibility", "version", "deleted"}).
			AddRow(recipeID, "owner-1", "private", 1000, false))
	mock.ExpectCommit()

	err := p.Handle(staleEvent)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProjectorHandle_DropsEqualVersion(t *testing.T) {
	db, mock := newMockDatabase(t)
	p := NewProjector(db)

	recipeID := uuid.New()
	duplicateEvent := rmq.RecipeAccessState{
		RecipeID:   recipeID,
		OwnerID:    "owner-1",
		Visibility: "private",
		Version:    1000, // exact redelivery of the already-applied event
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_accesses" WHERE recipe_id = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "owner_id", "visibility", "version", "deleted"}).
			AddRow(recipeID, "owner-1", "private", 1000, false))
	mock.ExpectCommit()

	err := p.Handle(duplicateEvent)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
