package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Visibility controls who besides the owner can see a recipe.
type Visibility string

const (
	VisibilityPrivate Visibility = "private"
	VisibilityShared  Visibility = "shared"
	VisibilityPublic  Visibility = "public"
)

func (v Visibility) Valid() bool {
	switch v {
	case VisibilityPrivate, VisibilityShared, VisibilityPublic:
		return true
	default:
		return false
	}
}

// RecipeListScope narrows GET /recipe to a subset of what's otherwise
// visible to the requester.
type RecipeListScope string

const (
	ScopeMine RecipeListScope = "mine"
	// ScopeShared returns recipes shared with the requester directly or
	// via a group (owner excluded), resolved against recipe_shares.
	ScopeShared RecipeListScope = "shared"
	// ScopeAll (the default) returns owned plus every public recipe -
	// the set GET /recipe always returned before ?scope existed.
	ScopeAll RecipeListScope = "all"
)

func (s RecipeListScope) Valid() bool {
	switch s {
	case ScopeMine, ScopeShared, ScopeAll:
		return true
	default:
		return false
	}
}

type Recipe struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuidv7();primaryKey"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Name         string         `gorm:"not null" json:"RecipeName" example:"apple pie"`
	Description  string         `gorm:"size:65535;not null" json:"Description" example:"pie with apples"`
	ServingCount int            `gorm:"default:0" json:"serving_count" example:"4"`
	// OwnerID is set from the authenticated Principal on create, never
	// from client input.
	OwnerID uuid.UUID `gorm:"type:uuid;not null;index" json:"owner_id"`
	// Visibility controls read access for non-owners: private (owner
	// only), shared (plus explicit shares - see shared/authz), or public
	// (any authenticated user).
	Visibility Visibility `gorm:"type:varchar(10);not null;default:'private';index" json:"visibility" example:"private"`
}

func (r Recipe) ConvertToDTO() RecipeDTO {
	return RecipeDTO{
		ID:           r.ID,
		Name:         r.Name,
		Description:  r.Description,
		ServingCount: r.ServingCount,
		OwnerID:      r.OwnerID,
		Visibility:   r.Visibility,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func (c Recipe) ConvertAllToDTO(recipes []Recipe) []RecipeDTO {
	var data []RecipeDTO

	for _, recipe := range recipes {
		data = append(data, recipe.ConvertToDTO())
	}

	return data
}

type RecipeDTO struct {
	ID           uuid.UUID  `json:"id" example:"23582396-12a3-425b-a597-8a22052823da"`
	Name         string     `gorm:"not null" json:"name" example:"apple pie"`
	Description  string     `gorm:"size:65535;not null" json:"description" example:"pie with apples"`
	ServingCount int        `gorm:"default:0" json:"serving_count" example:"4"`
	OwnerID      uuid.UUID  `json:"owner_id" example:"23582396-12a3-425b-a597-8a22052823da"`
	Visibility   Visibility `json:"visibility" example:"private"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (r RecipeDTO) ConvertFromDTO() Recipe {
	return Recipe{
		ID:           r.ID,
		Name:         r.Name,
		Description:  r.Description,
		ServingCount: r.ServingCount,
		OwnerID:      r.OwnerID,
		Visibility:   r.Visibility,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func (t RecipeDTO) ConvertAllFromDTO(recipeDTOs []RecipeDTO) []Recipe {
	var data []Recipe

	for _, recipeDTO := range recipeDTOs {
		data = append(data, recipeDTO.ConvertFromDTO())
	}

	return data
}

// MarshalBinary implements encoding.BinaryMarshaler
func (r RecipeDTO) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (r *RecipeDTO) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}

// SharePermission is the access level a recipe_shares row grants its
// target: read-only, or read+write.
type SharePermission string

const (
	SharePermissionRead  SharePermission = "read"
	SharePermissionWrite SharePermission = "write"
)

func (p SharePermission) Valid() bool {
	switch p {
	case SharePermissionRead, SharePermissionWrite:
		return true
	default:
		return false
	}
}

// ShareTargetType identifies what a recipe_shares row's TargetID refers
// to: a single user (OIDC sub) or every member of a group, identified
// however that group appears in the OIDC token's `groups` claim - a
// name for an ordinary Authentik group, or a UUID for a self-service
// household. See authentik/blueprints/cookbook-oidc.yaml's groups scope
// mapping.
type ShareTargetType string

const (
	ShareTargetUser  ShareTargetType = "user"
	ShareTargetGroup ShareTargetType = "group"
)

func (t ShareTargetType) Valid() bool {
	switch t {
	case ShareTargetUser, ShareTargetGroup:
		return true
	default:
		return false
	}
}

// RecipeShare is one explicit access grant on a recipe, to a user or a
// group - layered on top of Visibility, not a replacement for it.
//
// Composite primary key (RecipeID, TargetType, TargetID): at most one
// grant per target per recipe, so a repeat PUT for the same target
// upserts rather than duplicates.
type RecipeShare struct {
	RecipeID uuid.UUID `gorm:"type:uuid;primaryKey" json:"recipe_id"`
	// TargetType and TargetID together identify who this share is for -
	// see ShareTargetType. Both share one TargetID column, stored as
	// plain text regardless of which type it is.
	TargetType ShareTargetType `gorm:"type:varchar(10);primaryKey" json:"target_type" example:"user"`
	TargetID   string          `gorm:"primaryKey" json:"target_id"`
	Permission SharePermission `gorm:"type:varchar(10);not null" json:"permission" example:"read"`
	// CreatedBy is always the recipe's owner (only the owner may manage
	// shares), kept as its own column rather than inferred from the
	// recipe row so it survives a future ownership change.
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (s RecipeShare) ConvertToDTO() RecipeShareDTO {
	return RecipeShareDTO{
		TargetType: s.TargetType,
		TargetID:   s.TargetID,
		Permission: s.Permission,
		CreatedBy:  s.CreatedBy,
		CreatedAt:  s.CreatedAt,
	}
}

func (RecipeShare) ConvertAllToDTO(shares []RecipeShare) []RecipeShareDTO {
	data := make([]RecipeShareDTO, 0, len(shares))
	for _, share := range shares {
		data = append(data, share.ConvertToDTO())
	}
	return data
}

type RecipeShareDTO struct {
	TargetType ShareTargetType `json:"target_type" example:"user"`
	TargetID   string          `json:"target_id"`
	Permission SharePermission `json:"permission" example:"read"`
	CreatedBy  uuid.UUID       `json:"created_by"`
	CreatedAt  time.Time       `json:"created_at"`
}

func (d RecipeShareDTO) ConvertFromDTO(recipeID uuid.UUID, createdBy uuid.UUID) RecipeShare {
	return RecipeShare{
		RecipeID:   recipeID,
		TargetType: d.TargetType,
		TargetID:   d.TargetID,
		Permission: d.Permission,
		CreatedBy:  createdBy,
	}
}
