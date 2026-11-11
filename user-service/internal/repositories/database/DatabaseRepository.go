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

// GetDisplayName returns the cached display name for id, or
// "user not found" if nothing is cached yet; the service layer falls
// back to Authentik and re-caches via UpsertDisplayName.
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

func (r DatabaseRepository) UpsertDisplayName(id uuid.UUID, displayName string) error {
	user := um.User{ID: id, DisplayName: displayName}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"display_name", "updated_at"}),
	}).Create(&user).Error
}
