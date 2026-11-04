package repositories

import (
	"errors"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"gorm.io/gorm"
)

type DatabaseRepository struct {
	db *gorm.DB
}

func NewDatabaseRepository(db *gorm.DB) *DatabaseRepository {
	return &DatabaseRepository{
		db: db,
	}
}

// FindAll retrieves a paginated list of recipes visible to requesterID:
// recipes they own, plus every public recipe. ("shared" visibility isn't
// included here yet - that needs the recipe_shares model from KIT-32.)
func (r DatabaseRepository) FindAll(pagination m.PaginationRequest, requesterID uuid.UUID) ([]m.Recipe, int64, error) {
	var recipes []m.Recipe
	var total int64

	if err := r.db.Model(&m.Recipe{}).Where("owner_id = ? OR visibility = ?", requesterID, m.VisibilityPublic).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("owner_id = ? OR visibility = ?", requesterID, m.VisibilityPublic).Limit(pagination.Limit).Offset(pagination.Offset()).Find(&recipes).Error; err != nil {
		return nil, 0, err
	}

	return recipes, total, nil
}

// Find searches for a specific recipe in the database and returns it when found.
func (r DatabaseRepository) FindSingle(recipe m.Recipe) (m.Recipe, error) {

	result := r.db.First(&recipe)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.Recipe{}, errors.New("not found")
		} else {
			return m.Recipe{}, result.Error
		}
	}

	return recipe, nil
}

// Create handles the creation of a recipe and stores the relevant information in the database.
// outbox is inserted in the same transaction as the recipe row, so a
// recipe is never committed without its corresponding event becoming
// durable (the Relay publishes it from the outbox table afterwards).
func (r DatabaseRepository) Create(recipe m.Recipe, outbox rmq.OutboxEvent) (m.Recipe, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Create(&recipe).Error; err != nil {
			return err
		}

		if err = tx.Create(&outbox).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return recipe, err
	}

	return recipe, nil
}

func (r DatabaseRepository) Update(recipe m.Recipe, outbox rmq.OutboxEvent) (m.Recipe, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Updates(&recipe).Error; err != nil {
			return err
		}

		if err = tx.Create(&outbox).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return recipe, err
	}
	return recipe, nil
}

func (r DatabaseRepository) Delete(recipe m.Recipe, outbox rmq.OutboxEvent) error {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(&recipe).Error; err != nil {
			return err
		}

		if err := tx.Create(&outbox).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
