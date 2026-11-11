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

// RecipeAccess is this service's local, read-optimized projection of a
// recipe's ownership, visibility and shares, kept up to date by a
// Projector consuming recipe-service's ECST events. It is NOT the source
// of truth - recipe-service's own Recipe row is.
type RecipeAccess struct {
	RecipeID uuid.UUID `gorm:"type:uuid;primaryKey"`
	// String (not uuid.UUID) to match RecipeAccessState.OwnerID and
	// oidc.Principal.UserID.
	OwnerID string `gorm:"not null;index"`
	// Plain string mirroring shared/models.Visibility's values, to avoid
	// depending on shared/models for one enum.
	Visibility string `gorm:"type:varchar(10);not null;default:'private';index"`
	// JSON-encoded []rabbitmq.ShareEntry.
	Shares []byte `gorm:"type:jsonb"`
	// Projector.Handle only applies an update when its Version is
	// strictly greater, so this makes updates safe under out-of-order or
	// duplicate delivery.
	Version int64 `gorm:"not null"`
	Deleted bool  `gorm:"not null;default:false"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// AutoMigrate migrates the recipe_access table into db.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&RecipeAccess{})
}
