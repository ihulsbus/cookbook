package replay

import (
	"context"
	"encoding/json"
	stdlog "log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type testLogger struct{}

func (testLogger) Errorf(format string, args ...interface{}) {}
func (testLogger) Infof(format string, args ...interface{})  {}

func newMockDatabase(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	var mockDB *gorm.DB

	newLogger := logger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
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

	mockDB, err = gorm.Open(dialector, &gorm.Config{Logger: newLogger})
	if err != nil {
		t.Fatalf("gorm mock init failed: %v", err.Error())
	}

	return mockDB, mock
}

func TestRepublishAccess_OK(t *testing.T) {
	db, mock := newMockDatabase(t)

	activeRecipe := models.Recipe{
		ID:         uuid.New(),
		Name:       "apple pie",
		OwnerID:    uuid.New(),
		Visibility: models.VisibilityPrivate,
		UpdatedAt:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	deletedRecipe := models.Recipe{
		ID:         uuid.New(),
		Name:       "gone pie",
		OwnerID:    uuid.New(),
		Visibility: models.VisibilityPrivate,
		UpdatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:  gorm.DeletedAt{Time: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), Valid: true},
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" ORDER BY "recipes"."id" LIMIT $1`)).
		WithArgs(200).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "owner_id", "visibility", "updated_at", "deleted_at"}).
			AddRow(activeRecipe.ID, activeRecipe.Name, activeRecipe.OwnerID, activeRecipe.Visibility, activeRecipe.UpdatedAt, nil).
			AddRow(deletedRecipe.ID, deletedRecipe.Name, deletedRecipe.OwnerID, deletedRecipe.Visibility, deletedRecipe.UpdatedAt, deletedRecipe.DeletedAt.Time))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "outbox_events"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(uuid.New()).
			AddRow(uuid.New()))
	mock.ExpectCommit()

	err := RepublishAccess(context.Background(), db, testLogger{})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepublishAccess_LoadErr(t *testing.T) {
	db, mock := newMockDatabase(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" ORDER BY "recipes"."id" LIMIT $1`)).
		WithArgs(200).
		WillReturnError(assert.AnError)

	err := RepublishAccess(context.Background(), db, testLogger{})

	assert.Error(t, err)
}

func TestRepublishAccess_PublishErr(t *testing.T) {
	db, mock := newMockDatabase(t)

	activeRecipe := models.Recipe{
		ID:         uuid.New(),
		Name:       "apple pie",
		OwnerID:    uuid.New(),
		Visibility: models.VisibilityPrivate,
		UpdatedAt:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" ORDER BY "recipes"."id" LIMIT $1`)).
		WithArgs(200).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "owner_id", "visibility", "updated_at", "deleted_at"}).
			AddRow(activeRecipe.ID, activeRecipe.Name, activeRecipe.OwnerID, activeRecipe.Visibility, activeRecipe.UpdatedAt, nil))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "outbox_events"`)).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := RepublishAccess(context.Background(), db, testLogger{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "1 of 1 events failed to queue")
}

func TestAccessOutboxEntry_UsesUpdatedAtVersion(t *testing.T) {
	recipe := models.Recipe{
		ID:         uuid.New(),
		OwnerID:    uuid.New(),
		Visibility: models.VisibilityPublic,
		UpdatedAt:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	entry, err := accessOutboxEntry(recipe)

	assert.NoError(t, err)
	assert.Equal(t, rmq.RecipeAccessChangedRoutingKey, entry.RoutingKey)

	var state rmq.RecipeAccessState
	assert.NoError(t, json.Unmarshal(entry.Payload, &state))
	assert.Equal(t, recipe.ID, state.RecipeID)
	assert.Equal(t, recipe.UpdatedAt.UnixNano(), state.Version)
	assert.False(t, state.Deleted)
}

func TestAccessOutboxEntry_DeletedUsesDeletedAtVersion(t *testing.T) {
	recipe := models.Recipe{
		ID:         uuid.New(),
		OwnerID:    uuid.New(),
		Visibility: models.VisibilityPrivate,
		UpdatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:  gorm.DeletedAt{Time: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), Valid: true},
	}

	entry, err := accessOutboxEntry(recipe)

	assert.NoError(t, err)

	var state rmq.RecipeAccessState
	assert.NoError(t, json.Unmarshal(entry.Payload, &state))
	assert.Equal(t, recipe.DeletedAt.Time.UnixNano(), state.Version)
	assert.True(t, state.Deleted)
}
