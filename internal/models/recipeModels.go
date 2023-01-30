package models

import (
	"gorm.io/gorm"
)

// Recipe struct to hold recipe data
type Recipe struct {
	gorm.Model
	RecipeName      string       `gorm:"not null" json:"RecipeName"`
	Description     string       `gorm:"not null" json:"Description"`
	DifficultyLevel int          `gorm:"not null" json:"DifficultyLevel"`
	CookingTime     int          `gorm:"default:0" json:"CookTime"`
	ServingCount    int          `gorm:"default:0" json:"ServingCount"`
	Ingredient      []Ingredient `gorm:"many2many:recipe_ingredient;" json:"Ingredients"`
	Category        []Category   `gorm:"many2many:recipe_category;" json:"Categories"`
	Tags            []Tag        `gorm:"many2many:recipe_tag;" json:"Tags"`
	ImageName       string
}

// Instruction struct to hold instruction data
type Instruction struct {
	gorm.Model
	RecipeID    int    `json:"RecipeID"`
	Recipe      Recipe `gorm:"references:ID"`
	StepNumber  int    `json:"StepNumber"`
	Description string `json:"Description"`
}

// Ingredient struct to hold ingredient data
type Ingredient struct {
	gorm.Model
	IngredientName string `gorm:"unique; not null" json:"IngredientName"`
}

// RecipeIngredient struct to hold recipe ingredient data
type RecipeIngredient struct {
	RecipeID     int    `gorm:"primaryKey" json:"RecipeID"`
	IngredientID int    `gorm:"primaryKey" json:"IngredientID"`
	Quantity     int    `json:"Quantity"`
	Unit         string `json:"Unit"`
}

// Category struct to hold category data
type Category struct {
	gorm.Model
	CategoryName string `gorm:"not null;unique" json:"CategoryName"`
}

// Tag struct to hold tag data
type Tag struct {
	gorm.Model
	TagName string `gorm:"not null;unique" json:"TagName"`
}
