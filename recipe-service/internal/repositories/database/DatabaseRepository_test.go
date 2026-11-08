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
	"github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
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
		OwnerID:      uuid.New(),
		Visibility:   models.VisibilityPrivate,
	}

	testOutbox rmq.OutboxEvent = rmq.OutboxEvent{
		EventID:       uuid.New(),
		RoutingKey:    rmq.RecipeCreatedRoutingKey,
		Payload:       []byte(`{}`),
		SchemaVersion: 1,
	}

	// noopBuildOutbox is the buildOutbox callback for tests that don't
	// care about the shares a mutation's transaction reads - it just
	// returns testOutbox regardless.
	noopBuildOutbox = func(shares []models.RecipeShare) (rmq.OutboxEvent, error) {
		return testOutbox, nil
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

func TestRecipeFindAll_OK(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)
	pagination := models.NormalizePagination(1, 25)
	requesterID := recipe.OwnerID

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "recipes" WHERE (owner_id = $1 OR visibility = $2) AND "recipes"."deleted_at" IS NULL`)).
		WithArgs(requesterID, models.VisibilityPublic).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" WHERE (owner_id = $1 OR visibility = $2) AND "recipes"."deleted_at" IS NULL LIMIT $3`)).
		WithArgs(requesterID, models.VisibilityPublic, 25).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "serving_count", "owner_id", "visibility"}).
			AddRow(
				recipe.ID,
				recipe.Name,
				recipe.Description,
				recipe.ServingCount,
				recipe.OwnerID,
				recipe.Visibility,
			))

	data, total, err := r.FindAll(pagination, requesterID, nil, models.ScopeAll)

	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, recipe.ID, data[0].ID)
	assert.Equal(t, recipe.Name, data[0].Name)
	assert.Equal(t, recipe.Description, data[0].Description)
}

func TestRecipeFindAll_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)
	pagination := models.NormalizePagination(1, 25)
	requesterID := recipe.OwnerID

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "recipes" WHERE (owner_id = $1 OR visibility = $2) AND "recipes"."deleted_at" IS NULL`)).
		WithArgs(requesterID, models.VisibilityPublic).
		WillReturnError(errors.New("error"))

	data, total, err := r.FindAll(pagination, requesterID, nil, models.ScopeAll)

	assert.Error(t, err)
	assert.Len(t, data, 0)
	assert.Equal(t, int64(0), total)
}

func TestRecipeFindAll_Empty(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)
	pagination := models.NormalizePagination(1, 25)
	requesterID := recipe.OwnerID

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "recipes" WHERE (owner_id = $1 OR visibility = $2) AND "recipes"."deleted_at" IS NULL`)).
		WithArgs(requesterID, models.VisibilityPublic).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" WHERE (owner_id = $1 OR visibility = $2) AND "recipes"."deleted_at" IS NULL LIMIT $3`)).
		WithArgs(requesterID, models.VisibilityPublic, 25).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "serving_count", "owner_id", "visibility"}))

	data, total, err := r.FindAll(pagination, requesterID, nil, models.ScopeAll)

	assert.NoError(t, err)
	assert.Len(t, data, 0)
	assert.Equal(t, int64(0), total)
}

func TestRecipeFindAll_ScopeMine(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)
	pagination := models.NormalizePagination(1, 25)
	requesterID := recipe.OwnerID

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "recipes" WHERE owner_id = $1 AND "recipes"."deleted_at" IS NULL`)).
		WithArgs(requesterID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" WHERE owner_id = $1 AND "recipes"."deleted_at" IS NULL LIMIT $2`)).
		WithArgs(requesterID, 25).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "serving_count", "owner_id", "visibility"}).
			AddRow(
				recipe.ID,
				recipe.Name,
				recipe.Description,
				recipe.ServingCount,
				recipe.OwnerID,
				recipe.Visibility,
			))

	data, total, err := r.FindAll(pagination, requesterID, nil, models.ScopeMine)

	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, int64(1), total)
}

// TestRecipeFindAll_ScopeShared_NoGroups: requesterGroups empty - only the
// direct-user-share subquery branch should run (KIT-35).
func TestRecipeFindAll_ScopeShared_NoGroups(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)
	pagination := models.NormalizePagination(1, 25)
	requesterID := recipe.OwnerID

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "recipes" WHERE (owner_id != $1 AND id IN (SELECT recipe_id FROM recipe_shares WHERE target_type = 'user' AND target_id = $2)) AND "recipes"."deleted_at" IS NULL`)).
		WithArgs(requesterID, requesterID.String()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" WHERE (owner_id != $1 AND id IN (SELECT recipe_id FROM recipe_shares WHERE target_type = 'user' AND target_id = $2)) AND "recipes"."deleted_at" IS NULL LIMIT $3`)).
		WithArgs(requesterID, requesterID.String(), 25).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "serving_count", "owner_id", "visibility"}).
			AddRow(
				recipe.ID,
				recipe.Name,
				recipe.Description,
				recipe.ServingCount,
				recipe.OwnerID,
				recipe.Visibility,
			))

	data, total, err := r.FindAll(pagination, requesterID, nil, models.ScopeShared)

	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, int64(1), total)
}

// TestRecipeFindAll_ScopeShared_WithGroups: requesterGroups non-empty -
// the group-share subquery branch is also included (KIT-35).
func TestRecipeFindAll_ScopeShared_WithGroups(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)
	pagination := models.NormalizePagination(1, 25)
	requesterID := recipe.OwnerID
	requesterGroups := []string{"cookbook-user"}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "recipes" WHERE (owner_id != $1 AND id IN (SELECT recipe_id FROM recipe_shares WHERE (target_type = 'user' AND target_id = $2) OR (target_type = 'group' AND target_id IN ($3)))) AND "recipes"."deleted_at" IS NULL`)).
		WithArgs(requesterID, requesterID.String(), "cookbook-user").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" WHERE (owner_id != $1 AND id IN (SELECT recipe_id FROM recipe_shares WHERE (target_type = 'user' AND target_id = $2) OR (target_type = 'group' AND target_id IN ($3)))) AND "recipes"."deleted_at" IS NULL LIMIT $4`)).
		WithArgs(requesterID, requesterID.String(), "cookbook-user", 25).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "serving_count", "owner_id", "visibility"}).
			AddRow(
				recipe.ID,
				recipe.Name,
				recipe.Description,
				recipe.ServingCount,
				recipe.OwnerID,
				recipe.Visibility,
			))

	data, total, err := r.FindAll(pagination, requesterID, requesterGroups, models.ScopeShared)

	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, int64(1), total)
}

func TestRecipeFindSingle_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipes" WHERE "recipes"."deleted_at" IS NULL AND "recipes"."id" = $1 ORDER BY "recipes"."id" LIMIT $2`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(recipe.ID))

	result, err := r.FindSingle(recipe)

	assert.NoError(t, err)
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeFindSingle_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_category" WHERE "recipe_category"."recipe_id" = 1`)).
		WillReturnError(errors.New("error"))
	result, err := r.FindSingle(recipe)

	assert.Error(t, err)
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeFindShares_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	share := models.RecipeShare{
		RecipeID:   recipe.ID,
		TargetType: models.ShareTargetUser,
		TargetID:   uuid.New().String(),
		Permission: models.SharePermissionRead,
		CreatedBy:  recipe.OwnerID,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_shares" WHERE recipe_id = $1 ORDER BY created_at`)).
		WithArgs(recipe.ID).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "target_type", "target_id", "permission", "created_by"}).
			AddRow(share.RecipeID, share.TargetType, share.TargetID, share.Permission, share.CreatedBy))

	result, err := r.FindShares(recipe.ID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, share.TargetID, result[0].TargetID)
}

func TestRecipeFindShares_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_shares" WHERE recipe_id = $1 ORDER BY created_at`)).
		WithArgs(recipe.ID).
		WillReturnError(errors.New("error"))

	result, err := r.FindShares(recipe.ID)

	assert.Error(t, err)
	assert.Len(t, result, 0)
}

func TestRecipeCreate_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "recipes"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(recipe.ID))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "outbox_events"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testOutbox.ID))
	mock.ExpectCommit()

	result, err := r.Create(recipe, testOutbox)

	assert.NoError(t, err)
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeCreate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "recipes"`)).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()
	result, err := r.Create(recipe, testOutbox)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeCreate_OutboxErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "recipes"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(recipe.ID))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "outbox_events"`)).
		WillReturnError(errors.New("outbox error"))
	mock.ExpectRollback()

	result, err := r.Create(recipe, testOutbox)

	assert.Error(t, err)
	assert.EqualError(t, err, "outbox error")
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeUpdate_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipes" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_shares" WHERE recipe_id = $1`)).
		WithArgs(recipe.ID).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "target_type", "target_id", "permission", "created_by"}))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "outbox_events"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testOutbox.ID))
	mock.ExpectCommit()

	result, err := r.Update(recipe, noopBuildOutbox)

	assert.NoError(t, err)
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeUpdate_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipes" SET`)).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.Update(recipe, noopBuildOutbox)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeUpdate_SharesErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipes" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_shares" WHERE recipe_id = $1`)).
		WithArgs(recipe.ID).
		WillReturnError(errors.New("shares error"))
	mock.ExpectRollback()

	result, err := r.Update(recipe, noopBuildOutbox)

	assert.Error(t, err)
	assert.EqualError(t, err, "shares error")
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeUpdate_BuildOutboxErr(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipes" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_shares" WHERE recipe_id = $1`)).
		WithArgs(recipe.ID).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "target_type", "target_id", "permission", "created_by"}))
	mock.ExpectRollback()

	result, err := r.Update(recipe, func(shares []models.RecipeShare) (rmq.OutboxEvent, error) {
		return rmq.OutboxEvent{}, errors.New("build outbox error")
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "build outbox error")
	assert.IsType(t, result, models.Recipe{})
}

func TestRecipeDelete_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipes" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_shares" WHERE recipe_id = $1`)).
		WithArgs(recipe.ID).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "target_type", "target_id", "permission", "created_by"}))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "outbox_events"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testOutbox.ID))
	mock.ExpectCommit()

	err := r.Delete(recipe, noopBuildOutbox)

	assert.NoError(t, err)
}

func TestRecipeDelete_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipes" SET`)).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	err := r.Delete(recipe, noopBuildOutbox)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestRecipeUpsertShare_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	share := models.RecipeShare{
		RecipeID:   recipe.ID,
		TargetType: models.ShareTargetUser,
		TargetID:   uuid.New().String(),
		Permission: models.SharePermissionRead,
		CreatedBy:  recipe.OwnerID,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "recipe_shares"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_shares" WHERE recipe_id = $1`)).
		WithArgs(recipe.ID).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "target_type", "target_id", "permission", "created_by"}).
			AddRow(share.RecipeID, share.TargetType, share.TargetID, share.Permission, share.CreatedBy))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "outbox_events"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testOutbox.ID))
	mock.ExpectCommit()

	err := r.UpsertShare(share, noopBuildOutbox)

	assert.NoError(t, err)
}

func TestRecipeUpsertShare_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	share := models.RecipeShare{
		RecipeID:   recipe.ID,
		TargetType: models.ShareTargetUser,
		TargetID:   uuid.New().String(),
		Permission: models.SharePermissionRead,
		CreatedBy:  recipe.OwnerID,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "recipe_shares"`)).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	err := r.UpsertShare(share, noopBuildOutbox)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestRecipeDeleteShare_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	targetID := uuid.New().String()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "recipe_shares" WHERE recipe_id = $1 AND target_type = $2 AND target_id = $3`)).
		WithArgs(recipe.ID, models.ShareTargetUser, targetID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_shares" WHERE recipe_id = $1`)).
		WithArgs(recipe.ID).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "target_type", "target_id", "permission", "created_by"}))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "outbox_events"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testOutbox.ID))
	mock.ExpectCommit()

	err := r.DeleteShare(recipe.ID, models.ShareTargetUser, targetID, noopBuildOutbox)

	assert.NoError(t, err)
}

func TestRecipeDeleteShare_NotFound(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	targetID := uuid.New().String()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "recipe_shares" WHERE recipe_id = $1 AND target_type = $2 AND target_id = $3`)).
		WithArgs(recipe.ID, models.ShareTargetUser, targetID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := r.DeleteShare(recipe.ID, models.ShareTargetUser, targetID, noopBuildOutbox)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
}

func TestRecipeDeleteShare_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	targetID := uuid.New().String()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "recipe_shares" WHERE recipe_id = $1 AND target_type = $2 AND target_id = $3`)).
		WithArgs(recipe.ID, models.ShareTargetUser, targetID).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	err := r.DeleteShare(recipe.ID, models.ShareTargetUser, targetID, noopBuildOutbox)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

func TestRecipeUpdateVisibility_Ok(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipes" SET "visibility"=$1,"updated_at"=$2 WHERE id = $3 AND "recipes"."deleted_at" IS NULL`)).
		WithArgs(models.VisibilityPublic, timeFunc(), recipe.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_shares" WHERE recipe_id = $1`)).
		WithArgs(recipe.ID).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "target_type", "target_id", "permission", "created_by"}))
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "outbox_events"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testOutbox.ID))
	mock.ExpectCommit()

	updated := recipe
	updated.Visibility = models.VisibilityPublic

	result, err := r.UpdateVisibility(updated, noopBuildOutbox)

	assert.NoError(t, err)
	assert.Equal(t, models.VisibilityPublic, result.Visibility)
}

func TestRecipeUpdateVisibility_Err(t *testing.T) {
	db, mock := newMockDatabase(t)
	r := NewDatabaseRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "recipes" SET "visibility"=$1,"updated_at"=$2 WHERE id = $3 AND "recipes"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	result, err := r.UpdateVisibility(recipe, noopBuildOutbox)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, models.Recipe{})
}
