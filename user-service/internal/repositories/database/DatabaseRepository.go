package repositories

import (
	"errors"

	"github.com/google/uuid"
	um "user-service/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DatabaseRepository struct {
	db *gorm.DB
}

func NewDatabaseRepository(db *gorm.DB) *DatabaseRepository {
	return &DatabaseRepository{
		db: db,
	}
}

// GetDisplayName returns the locally cached display name for id, or
// "user not found" if nothing has been cached for it yet - the service
// layer then falls back to Authentik (see
// services.UserService.GetDisplayName) and caches the result via
// UpsertDisplayName.
func (r DatabaseRepository) GetDisplayName(id uuid.UUID) (string, error) {
	var user um.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	return user.DisplayName, nil
}

// UpsertDisplayName inserts id's display name into the cache, or
// refreshes it if already present.
func (r DatabaseRepository) UpsertDisplayName(id uuid.UUID, displayName string) error {
	user := um.User{ID: id, DisplayName: displayName}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"display_name", "updated_at"}),
	}).Create(&user).Error
}
