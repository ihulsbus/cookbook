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

// Valid reports whether v is one of the defined Visibility values.
func (v Visibility) Valid() bool {
	switch v {
	case VisibilityPrivate, VisibilityShared, VisibilityPublic:
		return true
	default:
		return false
	}
}

// RecipeListScope narrows GET /recipe (KIT-30) to a subset of what's
// otherwise visible to the requester.
type RecipeListScope string

const (
	// ScopeMine returns only recipes the requester owns.
	ScopeMine RecipeListScope = "mine"
	// ScopeShared returns recipes shared with the requester (owner
	// excluded). Always empty until the recipe_shares model (KIT-32)
	// exists to resolve explicit grants.
	ScopeShared RecipeListScope = "shared"
	// ScopeAll (the default) returns recipes the requester owns, plus
	// every public recipe - the same set GET /recipe always returned
	// before ?scope existed.
	ScopeAll RecipeListScope = "all"
)

// Valid reports whether s is one of the defined RecipeListScope values.
func (s RecipeListScope) Valid() bool {
	switch s {
	case ScopeMine, ScopeShared, ScopeAll:
		return true
	default:
		return false
	}
}

// Recipe struct to hold recipe data
type Recipe struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Name         string         `gorm:"not null" json:"RecipeName" example:"apple pie"`
	Description  string         `gorm:"size:65535;not null" json:"Description" example:"pie with apples"`
	ServingCount int            `gorm:"default:0" json:"serving_count" example:"4"`
	// OwnerID is the OIDC sub of the recipe's owner. Set from the
	// authenticated Principal on create and never from client input.
	OwnerID uuid.UUID `gorm:"type:uuid;not null;index" json:"owner_id"`
	// Visibility controls read access for non-owners: private (owner only),
	// shared (owner plus explicit shares - see shared/authz), or public
	// (every authenticated user). Defaults to private.
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
