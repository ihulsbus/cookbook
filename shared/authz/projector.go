package authz

import (
	"encoding/json"
	"errors"

	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	rabbitmqlib "github.com/wagslane/go-rabbitmq"
	"gorm.io/gorm"
)

// Projector keeps RecipeAccess rows in sync with recipe-service's ECST
// events.
type Projector struct {
	db *gorm.DB
}

// NewProjector creates a Projector backed by db. db must already have
// RecipeAccess migrated into it (see AutoMigrate).
func NewProjector(db *gorm.DB) *Projector {
	return &Projector{db: db}
}

// Handle upserts the projection for event.RecipeID, applying it only if
// event.Version is strictly greater than what's stored. Each ECST event
// carries the full access state (not a delta), so applying the
// highest-Version one seen is always correct regardless of delivery order
// or duplicate redelivery.
func (p *Projector) Handle(event rmq.RecipeAccessState) error {
	sharesJSON, err := json.Marshal(event.Shares)
	if err != nil {
		return err
	}

	incoming := RecipeAccess{
		RecipeID:   event.RecipeID,
		OwnerID:    event.OwnerID,
		Visibility: event.Visibility,
		Shares:     sharesJSON,
		Version:    event.Version,
		Deleted:    event.Deleted,
	}

	return p.db.Transaction(func(tx *gorm.DB) error {
		var existing RecipeAccess

		err := tx.First(&existing, "recipe_id = ?", event.RecipeID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tx.Create(&incoming).Error
			}
			return err
		}

		if event.Version <= existing.Version {
			// Stale or duplicate delivery; already applied.
			return nil
		}

		return tx.Model(&existing).
			Select("OwnerID", "Visibility", "Shares", "Version", "Deleted").
			Updates(incoming).Error
	})
}

// StartConsumer subscribes to recipe-service's lifecycle events and feeds
// each one to Handle.
func (p *Projector) StartConsumer(connection *rabbitmqlib.Conn, queueName, exchangeName string) (*rmq.Consumer, error) {
	routingKeys := []string{
		rmq.RecipeCreatedRoutingKey,
		rmq.RecipeUpdatedRoutingKey,
		rmq.RecipeDeletedRoutingKey,
		rmq.RecipeAccessChangedRoutingKey,
	}

	return rmq.NewConsumer(connection, queueName, routingKeys, exchangeName, p.handleDelivery)
}

func (p *Projector) handleDelivery(d rabbitmqlib.Delivery) rabbitmqlib.Action {
	var state rmq.RecipeAccessState

	switch d.RoutingKey {
	case rmq.RecipeCreatedRoutingKey, rmq.RecipeUpdatedRoutingKey, rmq.RecipeDeletedRoutingKey:
		var payload rmq.RecipePayload
		if err := json.Unmarshal(d.Body, &payload); err != nil {
			return rabbitmqlib.NackRequeue
		}
		state = payload.ToAccessState()
	case rmq.RecipeAccessChangedRoutingKey:
		if err := json.Unmarshal(d.Body, &state); err != nil {
			return rabbitmqlib.NackRequeue
		}
	default:
		return rabbitmqlib.NackDiscard
	}

	if err := p.Handle(state); err != nil {
		return rabbitmqlib.NackRequeue
	}

	return rabbitmqlib.Ack
}
