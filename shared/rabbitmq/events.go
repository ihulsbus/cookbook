package rabbitmq

import (
	"github.com/google/uuid"
)

// Event represents the metadata for a RabbitMQ event.
//
// EventID and SchemaVersion ride as AMQP message headers (set by Publish,
// see publisher.go), not in the JSON body, so they stay available for
// tracing/dedup without changing the payload's own shape.
type Event struct {
	Name          string
	RoutingKey    string
	Payload       interface{}
	EventID       uuid.UUID
	SchemaVersion int
}

// Predefined routing keys
const (
	// Recipes
	RecipeCreatedRoutingKey       = "recipe.created"
	RecipeUpdatedRoutingKey       = "recipe.updated"
	RecipeDeletedRoutingKey       = "recipe.deleted"
	RecipeAccessChangedRoutingKey = "recipe.access.changed"

	// Images
	ImageUpdatedRoutingKey = "image.updated"
	ImageDeletedRoutingKey = "image.deleted"
)

// CurrentSchemaVersion is the schema version stamped on every event
// published by this module version. Bump it when a payload's field shapes
// change in a way consumers need to branch on.
const CurrentSchemaVersion = 1

// Predefined events
var (
	// Recipes
	RecipeCreatedEvent = Event{
		Name:          "RecipeCreated",
		RoutingKey:    RecipeCreatedRoutingKey,
		Payload:       RecipePayload{},
		SchemaVersion: CurrentSchemaVersion,
	}
	RecipeUpdatedEvent = Event{
		Name:          "RecipeUpdated",
		RoutingKey:    RecipeUpdatedRoutingKey,
		Payload:       RecipePayload{},
		SchemaVersion: CurrentSchemaVersion,
	}
	RecipeDeletedEvent = Event{
		Name:          "RecipeDeleted",
		RoutingKey:    RecipeDeletedRoutingKey,
		Payload:       RecipePayload{},
		SchemaVersion: CurrentSchemaVersion,
	}
	RecipeAccessChangedEvent = Event{
		Name:          "RecipeAccessChanged",
		RoutingKey:    RecipeAccessChangedRoutingKey,
		Payload:       RecipeAccessState{},
		SchemaVersion: CurrentSchemaVersion,
	}

	// Images
	ImageUpdatedEvent = Event{
		Name:          "ImageUpdated",
		RoutingKey:    ImageUpdatedRoutingKey,
		Payload:       ImagePayload{},
		SchemaVersion: CurrentSchemaVersion,
	}

	ImageDeletedEvent = Event{
		Name:          "ImageDeleted",
		RoutingKey:    ImageDeletedRoutingKey,
		Payload:       ImagePayload{},
		SchemaVersion: CurrentSchemaVersion,
	}
)

// ShareEntry grants a permission on a recipe to a user or a group.
type ShareEntry struct {
	TargetType string `json:"target_type" example:"user"` // "user" | "group"
	TargetID   string `json:"target_id"`
	Permission string `json:"permission" example:"read"` // "read" | "write"
}

// RecipePayload is the body for recipe.created/updated/deleted. It carries
// the recipe's full content AND its full access state (not a delta), so
// consumers - both the existing content caches (image-service,
// ingredient-service) and shared/authz's projection - can apply it
// idempotently and out of order: drop the event when Version <= the
// consumer's own stored Version.
type RecipePayload struct {
	ID           uuid.UUID `json:"id" example:"505BBD89-5E8A-447F-9B4A-462CF56B124B"`
	Name         string    `json:"name" example:"apple pie"`
	Description  string    `json:"description" example:"pie with apples"`
	ServingCount int       `json:"servingcount" example:"4"`

	OwnerID    string       `json:"owner_id"`
	Visibility string       `json:"visibility" example:"private"`
	Shares     []ShareEntry `json:"shares"`
	Version    int64        `json:"version"`
	Deleted    bool         `json:"deleted"`
}

// ToAccessState extracts the access-state subset of a RecipePayload, for
// consumers (shared/authz) that only care about ownership/visibility/shares
// and not the recipe's content fields.
func (p RecipePayload) ToAccessState() RecipeAccessState {
	return RecipeAccessState{
		RecipeID:   p.ID,
		OwnerID:    p.OwnerID,
		Visibility: p.Visibility,
		Shares:     p.Shares,
		Version:    p.Version,
		Deleted:    p.Deleted,
	}
}

// RecipeAccessState is the body for recipe.access.changed - published when
// a recipe's ownership, visibility or shares change without a full content
// update (e.g. a future share-management endpoint). It carries the full
// access state, not a delta, for the same idempotent/out-of-order reasons
// as RecipePayload.
type RecipeAccessState struct {
	RecipeID   uuid.UUID    `json:"recipe_id"`
	OwnerID    string       `json:"owner_id"`
	Visibility string       `json:"visibility" example:"private"`
	Shares     []ShareEntry `json:"shares"`
	Version    int64        `json:"version"`
	Deleted    bool         `json:"deleted"`
}

// Image payload definitions
type ImagePayload struct {
	ID         uuid.UUID `json:"id" example:"505BBD89-5E8A-447F-9B4A-462CF56B124B"`
	EntityType string    `json:"entity_type" example:"recipe"`
	EntityID   uuid.UUID `json:"entity_id" example:"505BBD89-5E8A-447F-9B4A-462CF56B124B"`
	Size       int64     `json:"size" example:"12345"`
	Type       string    `json:"type" example:"jpeg"`
}
