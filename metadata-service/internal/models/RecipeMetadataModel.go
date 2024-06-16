package models

import "github.com/google/uuid"

// This is the combination model for all metadata associated to a recipe
type RecipeMetadataDTO struct {
	RecipeID        uuid.UUID `json:"recipe_id"`
	CategoryName    string    `json:"category_name"`
	CuisineTypeName string    `json:"cuisine_type_name"`
	Tags            []string  `json:"tags"`
	DifficultyLevel string    `json:"difficulty_level"`
	PreparationTime string    `json:"preparation_time"`
}
