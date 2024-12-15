package rabbitmq

import (
	"github.com/google/uuid"
)

// Event represents the metadata for a RabbitMQ event.
type Event struct {
	Name       string
	RoutingKey string
	Payload    interface{}
}

// Predefined routing keys
const (
	// Recipes
	RecipeCreatedRoutingKey = "recipe.created"
	RecipeUpdatedRoutingKey = "recipe.updated"
	RecipeDeletedRoutingKey = "recipe.deleted"

	// Images
	ImageUpdatedRoutingKey = "image.updated"
	ImageDeletedRoutingKey = "image.deleted"
)

// Predefined events
var (
	// Recipes
	RecipeCreatedEvent = Event{
		Name:       "RecipeCreated",
		RoutingKey: RecipeCreatedRoutingKey,
		Payload:    RecipePayload{},
	}
	RecipeUpdatedEvent = Event{
		Name:       "RecipeCreated",
		RoutingKey: RecipeUpdatedRoutingKey,
		Payload:    RecipePayload{},
	}
	RecipeDeletedEvent = Event{
		Name:       "RecipeCreated",
		RoutingKey: RecipeDeletedRoutingKey,
		Payload:    RecipePayload{},
	}

	// Images
	ImageUpdatedEvent = Event{
		Name:       "ImageUpdated",
		RoutingKey: ImageUpdatedRoutingKey,
		Payload:    ImagePayload{},
	}
)

// Recipe payload definitions
type RecipePayload struct {
	ID           uuid.UUID `json:"id" example:"505BBD89-5E8A-447F-9B4A-462CF56B124B"`
	Name         string    `json:"name" example:"apple pie"`
	Description  string    `json:"description" example:"pie with apples"`
	ServingCount int       `json:"servingcount" example:"4"`
}

// Image payload definitions
type ImagePayload struct {
	ID         uuid.UUID `json:"id" example:"505BBD89-5E8A-447F-9B4A-462CF56B124B"`
	EntityType string    `json:"entity_type" example:"recipe"`
	EntityID   uuid.UUID `json:"entity_id" example:"505BBD89-5E8A-447F-9B4A-462CF56B124B"`
	Size       int64     `json:"size" example:"12345"`
	Type       string    `json:"type" example:"jpeg"`
}
