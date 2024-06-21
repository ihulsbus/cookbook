package models

import "github.com/google/uuid"

// This is the combination model for all metadata associated to a recipe
type RecipeMetadataDTO struct {
	RecipeID        uuid.UUID `json:"recipe_id"`
	CategoryName    string    `json:"category_name"`
	CuisineTypeName string    `json:"cuisine_type_name"`
	Tags            []string  `json:"tags"`
	DifficultyLevel int       `json:"difficulty_level"`
	PreparationTime int       `json:"preparation_time"`
}
