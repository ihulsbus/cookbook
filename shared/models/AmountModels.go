package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Amount struct to hold recipe ingredient data
type Amount struct {
	RecipeID     uuid.UUID `gorm:"primaryKey"`
	IngredientID uuid.UUID `gorm:"primaryKey"`
	// Quantity (KIT-36) is an exact decimal, not a float: precision
	// matters down to e.g. 0.0025 g of an additive. Never convert it
	// through float64 anywhere in this path.
	Quantity decimal.Decimal `gorm:"type:numeric(12,4);not null" json:"Quantity" swaggertype:"string"`
	UnitID   uuid.UUID       `json:"UnitID"`
	Unit     Unit            `gorm:"references:ID"`
	// Note, Position and Optional (KIT-36) describe how this amount is
	// presented in the recipe's ingredient list - a free-text aside
	// ("finely chopped"), the display order, and whether it's skippable.
	Note     string `json:"Note"`
	Position int    `json:"Position"`
	Optional bool   `json:"Optional"`
}

func (r Amount) ConvertToDTO() AmountDTO {
	return AmountDTO{
		IngredientID: r.IngredientID,
		Quantity:     r.Quantity,
		UnitID:       r.UnitID,
		Note:         r.Note,
		Position:     r.Position,
		Optional:     r.Optional,
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
	IngredientID uuid.UUID       `json:"ingredientID" example:"23582396-12a3-425b-a597-8a22052823da"`
	Quantity     decimal.Decimal `json:"quantity" swaggertype:"string" example:"0.0025"`
	UnitID       uuid.UUID       `json:"unitID"`
	Note         string          `json:"note,omitempty" example:"finely chopped"`
	Position     int             `json:"position" example:"0"`
	Optional     bool            `json:"optional" example:"false"`
}

func (r AmountDTO) ConvertFromDTO() Amount {
	return Amount{
		IngredientID: r.IngredientID,
		Quantity:     r.Quantity,
		UnitID:       r.UnitID,
		Note:         r.Note,
		Position:     r.Position,
		Optional:     r.Optional,
	}
}

func (r AmountDTO) ConvertAllFromDTO(recipeIngredients []AmountDTO) []Amount {
	var data []Amount

	for _, ri := range recipeIngredients {
		data = append(data, ri.ConvertFromDTO())
	}

	return data
}
