package models

import "github.com/google/uuid"

// Amount struct to hold recipe ingredient data
type Amount struct {
	RecipeID     uuid.UUID `gorm:"primaryKey"`
	IngredientID uuid.UUID `gorm:"primaryKey"`
	Quantity     int       `json:"Quantity"`
	UnitID       uuid.UUID `json:"UnitID"`
	Unit         Unit      `gorm:"references:ID"`
}

func (r Amount) ConvertToDTO() AmountDTO {
	return AmountDTO{
		IngredientID: r.IngredientID,
		Quantity:     r.Quantity,
		Unit:         r.Unit.ConvertToDTO(),
	}
}

func (r Amount) ConvertAllToDTO(recipeIngredients []Amount) []AmountDTO {
	var data []AmountDTO

	for _, ri := range recipeIngredients {
		data = append(data, ri.ConvertToDTO())
	}

	return data
}

type AmountDTO struct {
	IngredientID uuid.UUID `json:"IngredientID" example:"23582396-12a3-425b-a597-8a22052823da"`
	Quantity     int       `json:"Quantity" example:"40"`
	Unit         UnitDTO   `json:"unit"`
}

func (r AmountDTO) ConvertFromDTO() Amount {
	return Amount{
		IngredientID: r.IngredientID,
		Quantity:     r.Quantity,
		UnitID:       r.Unit.ID,
		Unit:         r.Unit.ConvertFromDTO(),
	}
}

func (r AmountDTO) ConvertAllFromDTO(recipeIngredients []AmountDTO) []Amount {
	var data []Amount

	for _, ri := range recipeIngredients {
		data = append(data, ri.ConvertFromDTO())
	}

	return data
}
