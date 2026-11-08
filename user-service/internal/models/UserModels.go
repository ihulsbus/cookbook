// Package models holds user-service's own local tables - a cache of
// display names and household metadata (KIT-74). Unlike shared/models,
// these aren't exposed to other services as Go types: other services
// read them through user-service's own HTTP API (KIT-33), not by
// importing this package.
package models

import (
	"time"

	"github.com/google/uuid"
)

// User is a local cache of a display name for an Authentik user ID
// (OIDC sub). Populated and kept current against the Authentik API by
// KIT-33 - this table exists so the rest of the system never calls out
// to Authentik directly on every request, matching the project's
// no-service-calls-another-service-at-request-time preference wherever
// a cache will do (see cookbook/README.md's "No service calls another
// service at startup" section for the startup-time version of the same
// principle).
type User struct {
	// ID is the Authentik user UUID (the OIDC `sub` claim, per KIT-33's
	// "Subject mode: Based on the User's UUID" decision) - never
	// generated locally.
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	DisplayName string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// Household groups users for later use (pantry, shared grocery lists,
// etc.) - only metadata lives here if Authentik's own group attributes
// turn out not to be enough (see KIT-74's description). Membership
// changes are expected to publish household.member.added/removed
// (shared/rabbitmq) once something actually manages membership.
type Household struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// HouseholdMember is one user's membership in a household. Composite
// primary key (HouseholdID, UserID): a user belongs to a household at
// most once.
type HouseholdMember struct {
	HouseholdID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	JoinedAt    time.Time `gorm:"autoCreateTime"`
}
