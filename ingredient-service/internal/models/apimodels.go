package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Ingredient struct to hold ingredient data
type Ingredient struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string         `gorm:"unique; not null" json:"IngredientName"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type IngredientDTO struct {
	ID   uuid.UUID `json:"id" example:"23582396-12a3-425b-a597-8a22052823da"`
	Name string    `json:"ingredientName" example:"asparagus"`
}

func (i Ingredient) ConvertToDTO() IngredientDTO {
	return IngredientDTO{
		ID:   i.ID,
		Name: i.Name,
	}
}

func (c Ingredient) ConvertAllToDTO(ingredients []Ingredient) []IngredientDTO {
	var data []IngredientDTO

	for _, ingredient := range ingredients {
		data = append(data, ingredient.ConvertToDTO())
	}

	return data
}

func (i IngredientDTO) ConvertFromDTO() Ingredient {
	return Ingredient{
		ID:   i.ID,
		Name: i.Name,
	}
}

func (c IngredientDTO) ConvertAllFromDTO(ingredientDTOs []IngredientDTO) []Ingredient {
	var data []Ingredient

	for _, ingredientDTO := range ingredientDTOs {
		data = append(data, ingredientDTO.ConvertFromDTO())
	}

	return data
}

// RecipeIngredient struct to hold recipe ingredient data
type RecipeIngredient struct {
	RecipeID     uuid.UUID `gorm:"primaryKey"`
	IngredientID uuid.UUID `gorm:"primaryKey"`
	Quantity     int       `json:"Quantity"`
	UnitID       uuid.UUID `json:"UnitID"`
	Unit         Unit      `gorm:"references:ID"`
}

type RecipeIngredientDTO struct {
	RecipeID     uuid.UUID `json:"RecipeID" example:"23582396-12a3-425b-a597-8a22052823da"`
	IngredientID uuid.UUID `json:"IngredientID" example:"23582396-12a3-425b-a597-8a22052823da"`
	Quantity     int       `json:"Quantity" example:"40"`
	Unit         UnitDTO   `json:"unit"`
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

type Unit struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FullName  string         `gorm:"not null;unique" json:"FullName" example:"Fluid ounce"`
	ShortName string         `gorm:"not null;unique" json:"ShortName" example:"fl oz"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UnitDTO struct {
	ID        uuid.UUID `gorm:"primaryKey;not null;unique;index" json:"ID" example:"1"`
	FullName  string    `gorm:"not null;unique" json:"FullName" example:"Fluid ounce"`
	ShortName string    `gorm:"not null;unique" json:"ShortName" example:"fl oz"`
}

func (u Unit) ConvertToDTO() UnitDTO {
	return UnitDTO{
		ID:        u.ID,
		FullName:  u.FullName,
		ShortName: u.ShortName,
	}
}

func (c Unit) ConvertAllToDTO(units []Unit) []UnitDTO {
	var data []UnitDTO

	for _, unit := range units {
		data = append(data, unit.ConvertToDTO())
	}

	return data
}

func (u UnitDTO) ConvertFromDTO() Unit {
	return Unit{
		ID:        u.ID,
		FullName:  u.FullName,
		ShortName: u.ShortName,
	}
}

func (c UnitDTO) ConvertAllFromDTO(unitDTOs []UnitDTO) []Unit {
	var data []Unit

	for _, unitDTO := range unitDTOs {
		data = append(data, unitDTO.ConvertFromDTO())
	}

	return data
}
