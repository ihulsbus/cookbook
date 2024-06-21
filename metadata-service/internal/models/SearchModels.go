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
