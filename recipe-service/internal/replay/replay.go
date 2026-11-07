package replay

import (
	"context"
	"fmt"

	m "github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"gorm.io/gorm"
)

// republishLogger is the minimal logging surface RepublishAccess needs.
type republishLogger interface {
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
}

// RepublishAccess replays the current access state of every recipe in db -
// including soft-deleted ones, so a consumer that missed the original
// delete event still learns about it - as a recipe.access.changed outbox
// row (KIT-31). It's for bootstrapping a new consumer's projection, or
// recovering one that's known to have missed events, not for normal
// operation.
//
// Rows are written to the outbox table, the same as every other recipe
// event (see RecipeService.buildOutboxEntry) - never published directly -
// so delivery goes through the existing Relay and gets its
// retry-on-failure behaviour for free, and a crash partway through a
// republish just leaves the remaining recipes to pick up on a re-run
// (outbox inserts for already-processed recipes are already durable).
//
// Each row's Version is the recipe's own UpdatedAt (or DeletedAt for a
// soft-deleted recipe, which GORM's default soft delete leaves UpdatedAt
// unchanged for), i.e. the real historical point the state last changed,
// not "now". That makes a replay a true no-op for a consumer already at
// that Version - it drops the row as stale (Version <= stored) exactly
// like it would an out-of-order duplicate of the original event - rather
// than forcing every consumer's projection row to bump.
func RepublishAccess(ctx context.Context, db *gorm.DB, logger republishLogger) error {
	const batchSize = 200

	var total, queued, failed int

	var recipes []m.Recipe
	err := db.WithContext(ctx).Unscoped().FindInBatches(&recipes, batchSize, func(tx *gorm.DB, batchNum int) error {
		entries := make([]rmq.OutboxEvent, 0, len(recipes))
		for _, recipe := range recipes {
			total++

			entry, err := accessOutboxEntry(recipe)
			if err != nil {
				logger.Errorf("republish-access: failed to build event for recipe %s: %v", recipe.ID, err)
				failed++
				continue
			}
			entries = append(entries, entry)
		}

		if len(entries) == 0 {
			return nil
		}

		if err := db.WithContext(ctx).Create(&entries).Error; err != nil {
			logger.Errorf("republish-access: failed to write outbox rows for batch %d: %v", batchNum, err)
			failed += len(entries)
			return nil
		}

		queued += len(entries)
		return nil
	}).Error

	if err != nil {
		return fmt.Errorf("republish-access: failed to load recipes: %w", err)
	}

	logger.Infof("republish-access: done - %d recipes found, %d queued, %d failed", total, queued, failed)
	fmt.Printf("republish-access: %d recipes found, %d queued, %d failed\n", total, queued, failed)

	if failed > 0 {
		return fmt.Errorf("republish-access: %d of %d events failed to queue - see logs, re-run to retry", failed, total)
	}

	return nil
}

func accessOutboxEntry(recipe m.Recipe) (rmq.OutboxEvent, error) {
	deleted := !recipe.DeletedAt.Time.IsZero()

	version := recipe.UpdatedAt
	if deleted {
		version = recipe.DeletedAt.Time
	}

	event := rmq.RecipeAccessChangedEvent
	event.Payload = rmq.RecipeAccessState{
		RecipeID:   recipe.ID,
		OwnerID:    recipe.OwnerID.String(),
		Visibility: string(recipe.Visibility),
		Version:    version.UnixNano(),
		Deleted:    deleted,
	}

	return rmq.NewOutboxEntry(event)
}
