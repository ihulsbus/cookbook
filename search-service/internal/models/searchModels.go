package models

import (
	"github.com/google/uuid"
)

type SearchRequest struct {
	Query string `json:"query" form:"query" binding:"required"`
	Limit int    `json:"limit" form:"limit" binding:"omitempty,min=1"`
	Page  int    `json:"page" form:"page" binding:"omitempty,min=1"`
}

type SearchResult struct {
	Recipes     []RecipeResult         `json:"recipes"`
	Ingredients []IngredientResult     `json:"ingredients"`
	Metadata    []MetadataSearchResult `json:"metadata"`
}

type RecipeResult struct {
	ID           uuid.UUID
	Name         string `gorm:"not null" json:"name" example:"apple pie"`
	Description  string `gorm:"size:65535;not null" json:"description" example:"pie with apples"`
	ServingCount int    `gorm:"default:0" json:"servingcount" example:"4"`
}

type IngredientResult struct {
	ID   uuid.UUID `json:"id" example:"23582396-12a3-425b-a597-8a22052823da"`
	Name string    `json:"name" example:"asparagus"`
}

type MetadataSearchRequest struct {
	RecipeID          uuid.UUID `json:"recipe_id,omitempty"`
	CategoryID        uuid.UUID `json:"category_id,omitempty"`
	TagID             uuid.UUID `json:"tag_id,omitempty"`
	DifficultyLevelID uuid.UUID `json:"difficulty_level,omitempty"`
	CuisineTypeID     uuid.UUID `json:"cuisine_type,omitempty"`
	MinPrepTime       *int      `json:"min_prep_time,omitempty"` // in minutes
	MaxPrepTime       *int      `json:"max_prep_time,omitempty"` // in minutes
}

type MetadataSearchResult struct {
	RecipeID          uuid.UUID   `json:"recipe_id"`
	CategoryIDs       []uuid.UUID `json:"category_ids"`
	TagIDs            []uuid.UUID `json:"tag_ids"`
	DifficultyLevelID uuid.UUID   `json:"difficulty_level"`
	PreparationTimeID uuid.UUID   `json:"preparation_time"` // in minutes
	CuisineTypeID     uuid.UUID   `json:"cuisine_type"`
}
