package models

import "github.com/google/uuid"

type InstructionSearchRequest struct {
	RecipeID uuid.UUID
}

type InstructionSearchRequestDTO struct {
	RecipeID uuid.UUID `json:"recipe_id"`
}

type InstructionSearchResult struct {
	RecipeID       uuid.UUID
	InstructionIDs []uuid.UUID
}

type InstructionSearchResultDTO struct {
	RecipeID       uuid.UUID   `json:"recipe_id"`
	InstructionIDs []uuid.UUID `json:"instruction_ids"`
}

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
	RecipeID          uuid.UUID `json:"recipe_id,omitempty"`
	CategoryID        uuid.UUID `json:"category_id,omitempty"`
	TagID             uuid.UUID `json:"tag_id,omitempty"`
	DifficultyLevelID uuid.UUID `json:"difficulty_level,omitempty"`
	CuisineTypeID     uuid.UUID `json:"cuisine_type,omitempty"`
	MinPrepTime       *int      `json:"min_prep_time,omitempty"` // in minutes
	MaxPrepTime       *int      `json:"max_prep_time,omitempty"` // in minutes
}

func (s MetadataSearchRequestDTO) ConvertFromDTO() MetadataSearchRequest {
	return MetadataSearchRequest{
		RecipeID:          &s.RecipeID,
		CategoryID:        &s.CategoryID,
		TagID:             &s.TagID,
		DifficultyLevelID: &s.DifficultyLevelID,
		CuisineTypeID:     &s.CuisineTypeID,
		MinPrepTime:       s.MinPrepTime,
		MaxPrepTime:       s.MaxPrepTime,
	}
}

// result
type MetadataSearchResult struct {
	RecipeID          uuid.UUID   `json:"recipe_id,omitempty"`
	CategoryIDs       []uuid.UUID `json:"category_ids,omitempty"`
	TagIDs            []uuid.UUID `json:"tag_ids,omitempty"`
	DifficultyLevelID uuid.UUID   `json:"difficulty_level_id,omitempty"`
	CuisineTypeID     uuid.UUID   `json:"cuisine_type_id,omitempty"`
	PreparationTimeID uuid.UUID   `json:"preparation_time_id,omitempty"` // in minutes
}

type MetadataSearchResultDTO struct {
	RecipeID          uuid.UUID   `json:"recipe_id"`
	CategoryIDs       []uuid.UUID `json:"category_ids"`
	TagIDs            []uuid.UUID `json:"tag_ids"`
	DifficultyLevelID uuid.UUID   `json:"difficulty_level"`
	PreparationTimeID uuid.UUID   `json:"preparation_time"` // in minutes
	CuisineTypeID     uuid.UUID   `json:"cuisine_type"`
}

func (s MetadataSearchResult) ConvertToDTO() MetadataSearchResultDTO {
	return MetadataSearchResultDTO{
		RecipeID:          s.RecipeID,
		CategoryIDs:       s.CategoryIDs,
		TagIDs:            s.TagIDs,
		DifficultyLevelID: s.DifficultyLevelID,
		PreparationTimeID: s.PreparationTimeID,
		CuisineTypeID:     s.CuisineTypeID,
	}
}

func (s MetadataSearchResult) ConvertAllToDTO(searchResults []MetadataSearchResult) []MetadataSearchResultDTO {
	var data []MetadataSearchResultDTO

	for _, searchResult := range searchResults {
		data = append(data, searchResult.ConvertToDTO())
	}

	return data
}
