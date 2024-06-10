package models

import (
	"gorm.io/gorm"
)

// Ingredient struct to hold ingredient data
type Ingredient struct {
	gorm.Model
	IngredientName string `gorm:"unique; not null" json:"IngredientName" example:"apple"`
}

// RecipeIngredient struct to hold recipe ingredient data
type RecipeIngredient struct {
	RecipeID     int  `gorm:"primaryKey" json:"RecipeID" example:"1"`
	IngredientID int  `gorm:"primaryKey" json:"IngredientID" example:"1"`
	Quantity     int  `json:"Quantity" example:"40"`
	UnitID       int  `json:"UnitID" example:"1"`
	Unit         Unit `gorm:"references:ID"`
}

type Unit struct {
	ID        uint   `gorm:"primaryKey;not null;unique;index" json:"ID" example:"1"`
	FullName  string `gorm:"not null;unique" json:"FullName" example:"Fluid ounce"`
	ShortName string `gorm:"not null;unique" json:"ShortName" example:"fl oz"`
}
