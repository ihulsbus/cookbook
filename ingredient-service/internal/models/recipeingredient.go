package models

import "github.com/google/uuid"

// RecipeIngredient struct to hold recipe ingredient data
type RecipeIngredient struct {
	RecipeID     uuid.UUID `gorm:"primaryKey"`
	IngredientID uuid.UUID `gorm:"primaryKey"`
	Quantity     int       `json:"Quantity"`
	UnitID       uuid.UUID `json:"UnitID"`
	Unit         Unit      `gorm:"references:ID"`
}

func (r RecipeIngredient) ConvertToDTO() RecipeIngredientDTO {
	return RecipeIngredientDTO{
		RecipeID:     r.RecipeID,
		IngredientID: r.IngredientID,
		Quantity:     r.Quantity,
		Unit:         r.Unit.ConvertToDTO(),
	}
}

func (r RecipeIngredient) ConvertAllToDTO(recipeIngredients []RecipeIngredient) []RecipeIngredientDTO {
	var data []RecipeIngredientDTO

	for _, ri := range recipeIngredients {
		data = append(data, ri.ConvertToDTO())
	}

	return data
}

type RecipeIngredientDTO struct {
	RecipeID     uuid.UUID `json:"RecipeID" example:"23582396-12a3-425b-a597-8a22052823da"`
	IngredientID uuid.UUID `json:"IngredientID" example:"23582396-12a3-425b-a597-8a22052823da"`
	Quantity     int       `json:"Quantity" example:"40"`
	Unit         UnitDTO   `json:"unit"`
}

func (r RecipeIngredientDTO) ConvertFromDTO() RecipeIngredient {
	return RecipeIngredient{
		RecipeID:     r.RecipeID,
		IngredientID: r.IngredientID,
		Quantity:     r.Quantity,
		UnitID:       r.Unit.ID,
		Unit:         r.Unit.ConvertFromDTO(),
	}
}

func (r RecipeIngredientDTO) ConvertAllFromDTO(recipeIngredients []RecipeIngredientDTO) []RecipeIngredient {
	var data []RecipeIngredient

	for _, ri := range recipeIngredients {
		data = append(data, ri.ConvertFromDTO())
	}

	return data
}
