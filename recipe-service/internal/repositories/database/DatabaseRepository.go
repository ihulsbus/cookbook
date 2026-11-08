package repositories

import (
	"errors"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
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

// Every share/visibility mutation below takes a buildOutbox callback of
// type func([]m.RecipeShare) (rmq.OutboxEvent, error) - written out in
// full at each call site rather than as a named type, since a named
// func type is a distinct type for Go interface satisfaction purposes
// and RecipeRepository (services package) declares these methods with
// the literal func type. buildOutbox is called with the recipe's current
// shares, read inside the same transaction as the mutation itself -
// never as a separate query afterwards, which would race a concurrent
// share write and could publish a stale shares list.

// FindAll retrieves a paginated list of recipes visible to requesterID,
// narrowed by scope (KIT-30):
//   - ScopeMine: only recipes requesterID owns.
//   - ScopeShared: recipes shared with requesterID directly or via one
//     of requesterGroups, excluding ones requesterID owns (KIT-35),
//     resolved against the recipe_shares table (KIT-32).
//   - ScopeAll (default): recipes requesterID owns, plus every public
//     recipe - the set GET /recipe always returned before ?scope existed.
func (r DatabaseRepository) FindAll(pagination m.PaginationRequest, requesterID uuid.UUID, requesterGroups []string, scope m.RecipeListScope) ([]m.Recipe, int64, error) {
	var recipes []m.Recipe
	var total int64

	var where string
	var args []interface{}

	switch scope {
	case m.ScopeMine:
		where = "owner_id = ?"
		args = []interface{}{requesterID}
	case m.ScopeShared:
		if len(requesterGroups) > 0 {
			where = `owner_id != ? AND id IN (SELECT recipe_id FROM recipe_shares WHERE (target_type = 'user' AND target_id = ?) OR (target_type = 'group' AND target_id IN ?))`
			args = []interface{}{requesterID, requesterID.String(), requesterGroups}
		} else {
			where = `owner_id != ? AND id IN (SELECT recipe_id FROM recipe_shares WHERE target_type = 'user' AND target_id = ?)`
			args = []interface{}{requesterID, requesterID.String()}
		}
	default: // ScopeAll
		where = "owner_id = ? OR visibility = ?"
		args = []interface{}{requesterID, m.VisibilityPublic}
	}

	if err := r.db.Model(&m.Recipe{}).Where(where, args...).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where(where, args...).Limit(pagination.Limit).Offset(pagination.Offset()).Find(&recipes).Error; err != nil {
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

// FindShares returns every recipe_shares row for recipeID, oldest first.
func (r DatabaseRepository) FindShares(recipeID uuid.UUID) ([]m.RecipeShare, error) {
	var shares []m.RecipeShare
	err := r.db.Where("recipe_id = ?", recipeID).Order("created_at").Find(&shares).Error
	return shares, err
}

// Create handles the creation of a recipe and stores the relevant information in the database.
// outbox is inserted in the same transaction as the recipe row, so a
// recipe is never committed without its corresponding event becoming
// durable (the Relay publishes it from the outbox table afterwards). A
// brand new recipe can't have any recipe_shares rows yet, so unlike
// Update/Delete below, Create takes a precomputed outbox rather than a
// buildOutbox callback - there's nothing to read.
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

// Update applies changes to recipe and writes buildOutbox's event in the
// same transaction, having given it the recipe's current shares (read
// inside this transaction, not supplied by the caller) so the resulting
// event's access state is never stale relative to a concurrent share
// change.
func (r DatabaseRepository) Update(recipe m.Recipe, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) (m.Recipe, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(&recipe).Error; err != nil {
			return err
		}

		var shares []m.RecipeShare
		if err := tx.Where("recipe_id = ?", recipe.ID).Find(&shares).Error; err != nil {
			return err
		}

		outbox, err := buildOutbox(shares)
		if err != nil {
			return err
		}

		return tx.Create(&outbox).Error
	}); err != nil {
		return recipe, err
	}
	return recipe, nil
}

// Delete removes recipe and writes buildOutbox's event in the same
// transaction - see Update for why shares are read inside the
// transaction rather than passed in.
func (r DatabaseRepository) Delete(recipe m.Recipe, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&recipe).Error; err != nil {
			return err
		}

		var shares []m.RecipeShare
		if err := tx.Where("recipe_id = ?", recipe.ID).Find(&shares).Error; err != nil {
			return err
		}

		outbox, err := buildOutbox(shares)
		if err != nil {
			return err
		}

		return tx.Create(&outbox).Error
	})
}

// UpsertShare creates share, or - if a row already exists for the same
// (RecipeID, TargetType, TargetID) - updates its Permission, and writes
// buildOutbox's event in the same transaction.
func (r DatabaseRepository) UpsertShare(share m.RecipeShare, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "recipe_id"}, {Name: "target_type"}, {Name: "target_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"permission"}),
		}).Create(&share).Error; err != nil {
			return err
		}

		var shares []m.RecipeShare
		if err := tx.Where("recipe_id = ?", share.RecipeID).Find(&shares).Error; err != nil {
			return err
		}

		outbox, err := buildOutbox(shares)
		if err != nil {
			return err
		}

		return tx.Create(&outbox).Error
	})
}

// DeleteShare removes the recipe_shares row for (recipeID, targetType,
// targetID), if one exists, and writes buildOutbox's event in the same
// transaction. Returns an error reading "not found" if no such row
// existed.
func (r DatabaseRepository) DeleteShare(recipeID uuid.UUID, targetType m.ShareTargetType, targetID string, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("recipe_id = ? AND target_type = ? AND target_id = ?", recipeID, targetType, targetID).Delete(&m.RecipeShare{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("not found")
		}

		var shares []m.RecipeShare
		if err := tx.Where("recipe_id = ?", recipeID).Find(&shares).Error; err != nil {
			return err
		}

		outbox, err := buildOutbox(shares)
		if err != nil {
			return err
		}

		return tx.Create(&outbox).Error
	})
}

// UpdateVisibility changes only recipe's Visibility column (PUT
// /recipe/:id/visibility, KIT-32) and writes buildOutbox's event in the
// same transaction.
func (r DatabaseRepository) UpdateVisibility(recipe m.Recipe, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) (m.Recipe, error) {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&m.Recipe{}).Where("id = ?", recipe.ID).Update("visibility", recipe.Visibility).Error; err != nil {
			return err
		}

		var shares []m.RecipeShare
		if err := tx.Where("recipe_id = ?", recipe.ID).Find(&shares).Error; err != nil {
			return err
		}

		outbox, err := buildOutbox(shares)
		if err != nil {
			return err
		}

		return tx.Create(&outbox).Error
	}); err != nil {
		return recipe, err
	}
	return recipe, nil
}
