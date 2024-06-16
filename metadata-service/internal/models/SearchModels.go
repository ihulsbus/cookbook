package models

import "github.com/google/uuid"

// request
type MetadataSearchRequest struct {
	RecipeID          *uuid.UUID `json:"recipe_id,omitempty"`
	CategoryID        *uuid.UUID `json:"category_id,omitempty"`
	TagID             *uuid.UUID `json:"tag_id,omitempty"`
	DifficultyLevelID *uuid.UUID `json:"difficulty_level_id,omitempty"`
	CuisineTypeID     *uuid.UUID `json:"cuisine_type_id,omitempty"`
	MinPrepTime       *int       `json:"min_prep_time,omitempty"` // in minutes
	MaxPrepTime       *int       `json:"max_prep_time,omitempty"` // in minutes
}

type MetadataSearchRequestDTO struct {
	RecipeID        uuid.UUID `json:"recipe_id,omitempty"`
	CategoryID      uuid.UUID `json:"category_id,omitempty"`
	TagID           uuid.UUID `json:"tag_id,omitempty"`
	DifficultyLevel int       `json:"difficulty_level,omitempty"`
	CuisineType     string    `json:"cuisinetype,omitempty"`
	MinPrepTime     int       `json:"min_prep_time,omitempty"` // in minutes
	MaxPrepTime     int       `json:"max_prep_time,omitempty"` // in minutes
}

// result
type MetadataSearchResult struct {
	RecipeID          uuid.UUID   `json:"recipe_id,omitempty"`
	CategoryIDs       []uuid.UUID `json:"category_id,omitempty"`
	TagIDs            []uuid.UUID `json:"tag_id,omitempty"`
	DifficultyLevelID uuid.UUID   `json:"difficulty_level_id,omitempty"`
	CuisineTypeID     uuid.UUID   `json:"cuisine_type_id,omitempty"`
	PreparationTimeID uuid.UUID   `json:"preparation_time_id,omitempty"` // in minutes
}

type MetadataSearchResultDTO struct {
	RecipeID        uuid.UUID   `json:"recipe_id"`
	CategoryIDs     []uuid.UUID `json:"category_id"`
	TagIDs          []uuid.UUID `json:"tag_id"`
	DifficultyLevel int         `json:"difficulty_level"`
	PreparationTime int         `json:"preparation_time"` // in minutes
	CuisineType     string      `json:"cuisinetype"`
}
