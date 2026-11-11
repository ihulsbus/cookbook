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

// RepublishAccess replays every recipe's current access state (including
// soft-deleted ones) as a recipe.access.changed outbox row, for
// bootstrapping or recovering a consumer's projection - not for normal
// operation. Rows go through the same outbox/Relay path as every other
// recipe event, so a crash partway through is safe to re-run.
//
// Each row's Version is the recipe's own UpdatedAt (or DeletedAt if
// soft-deleted), not "now", so a replay is a true no-op for a consumer
// already at that version rather than forcing its projection to bump.
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
