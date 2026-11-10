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

// Household and HouseholdMember remain unused as of KIT-33/KIT-34:
// Authentik's own group attributes turned out to be enough (see
// KIT-74's "if Authentik group attributes aren't enough" - they were).
// A household is an Authentik group under a parent "households" group,
// with attributes {cookbook_household: true, owner: <user uuid>} and
// its `users` field as the membership list - see
// user-service/internal/services/UserService.go and KIT-34's
// description. Kept here rather than deleted in case a local
// denormalized read-model is ever needed (e.g. for the future pantry
// work KIT-74 anticipates).
type Household struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// HouseholdMember is one user's membership in a household. Composite
// primary key (HouseholdID, UserID): a user belongs to a household at
// most once. See Household's comment - unused as of KIT-33/KIT-34.
type HouseholdMember struct {
	HouseholdID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	JoinedAt    time.Time `gorm:"autoCreateTime"`
}

// UserDTO is a user as returned by GET /api/v2/user/users?query= (search)
// and GET /api/v2/user/users/:id (resolve) - KIT-33. Username is only
// populated by the search endpoint; the resolve endpoint only ever
// needs (and only ever has cached) a display name.
type UserDTO struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username,omitempty"`
	DisplayName string    `json:"display_name"`
}

// HouseholdDTO is a household as returned by the household endpoints
// (KIT-34). ID is the household's Authentik group UUID.
type HouseholdDTO struct {
	ID      string               `json:"id"`
	Name    string               `json:"name"`
	OwnerID uuid.UUID            `json:"owner_id"`
	Members []HouseholdMemberDTO `json:"members"`
}

// HouseholdMemberDTO is one member of a HouseholdDTO.
type HouseholdMemberDTO struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
}
