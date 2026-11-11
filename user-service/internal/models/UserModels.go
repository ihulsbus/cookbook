// Package models holds user-service's own local tables - a cache of
// display names and household metadata. Unlike shared/models, these
// aren't exposed to other services as Go types; other services read
// them through this service's HTTP API instead.
package models

import (
	"time"

	"github.com/google/uuid"
)

// User is a local cache of a display name for an Authentik user ID
// (OIDC sub), so the rest of the system doesn't call Authentik on
// every request.
type User struct {
	// ID is the Authentik user UUID (the OIDC `sub` claim) - never
	// generated locally.
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	DisplayName string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// Household and HouseholdMember are currently unused - Authentik's own
// group attributes turned out to be enough. Kept in case a local
// denormalized read-model is needed later.
type Household struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// HouseholdMember is one user's membership in a household. Composite
// primary key (HouseholdID, UserID) - a user belongs to a household at
// most once. Unused, see Household's comment.
type HouseholdMember struct {
	HouseholdID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	JoinedAt    time.Time `gorm:"autoCreateTime"`
}

// UserDTO is returned by the search and resolve user endpoints.
// Username is only populated by search; resolve only has a cached
// display name.
type UserDTO struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username,omitempty"`
	DisplayName string    `json:"display_name"`
}

// HouseholdDTO is a household as returned by the household endpoints.
// ID is the household's Authentik group UUID.
type HouseholdDTO struct {
	ID      string               `json:"id"`
	Name    string               `json:"name"`
	OwnerID uuid.UUID            `json:"owner_id"`
	Members []HouseholdMemberDTO `json:"members"`
}

type HouseholdMemberDTO struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
}
