package models

import "github.com/google/uuid"

type InstructionSearchRequest struct {
	RecipeID uuid.UUID
}

type InstructionSearchRequestDTO struct {
	RecipeID uuid.UUID `json:"recipe_id,omitempty"`
}

type InstructionSearchResult struct {
	RecipeID       uuid.UUID
	InstructionIDs []uuid.UUID
}

type InstructionSearchResultDTO struct {
	RecipeID       uuid.UUID   `json:"recipe_id,omitempty"`
	InstructionIDs []uuid.UUID `json:"instruction_ids,omitempty"`
}
