// Package authz is a reusable module each child service embeds to
// authorize recipe access without owning the Recipe row itself. It keeps
// a local projection (RecipeAccess) built from the ECST events published
// by recipe-service (shared/rabbitmq), and exposes an Authorizer plus a
// RequireRecipeAccess Gin middleware that check it.
package authz

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RecipeAccess is this service's local projection of a recipe's
// ownership, visibility and shares, kept up to date by a Projector
// consuming recipe-service's ECST events. It is NOT the source of truth -
// recipe-service's own Recipe row is - this is a read-optimized copy so
// every service doesn't have to call recipe-service synchronously to
// authorize a request.
type RecipeAccess struct {
	RecipeID uuid.UUID `gorm:"type:uuid;primaryKey"`
	// OwnerID is the OIDC sub of the owner, stored as a string (not
	// uuid.UUID) to match shared/rabbitmq.RecipeAccessState.OwnerID and
	// oidc.Principal.UserID, both of which are plain strings.
	OwnerID string `gorm:"not null;index"`
	// Visibility mirrors shared/models.Visibility's values (private |
	// shared | public) but is kept as a plain string here too, so this
	// module doesn't need to depend on shared/models for one enum.
	Visibility string `gorm:"type:varchar(10);not null;default:'private';index"`
	// Shares is the JSON-encoded []rabbitmq.ShareEntry for this recipe.
	Shares []byte `gorm:"type:jsonb"`
	// Version is the event Version this row was last updated from.
	// Projector.Handle only applies an event when its Version is strictly
	// greater, making projection updates safe under out-of-order or
	// duplicate delivery.
	Version int64 `gorm:"not null"`
	Deleted bool  `gorm:"not null;default:false"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// AutoMigrate migrates the recipe_access table into db. Call once at the
// embedding service's startup, alongside its own model migrations.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&RecipeAccess{})
}
