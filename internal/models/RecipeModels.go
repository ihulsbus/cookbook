package models

import (
	"gorm.io/gorm"
)

// Recipe struct to hold recipe data
type Recipe struct {
	gorm.Model
	RecipeName      string       `gorm:"not null" json:"RecipeName" example:"apple pie"`
	Description     string       `gorm:"not null" json:"Description" example:"pie with apples"`
	DifficultyLevel int          `gorm:"not null" json:"DifficultyLevel" example:"1"`
	CookingTime     int          `gorm:"default:0" json:"CookTime" example:"23"`
	ServingCount    int          `gorm:"default:0" json:"ServingCount" example:"4"`
	Ingredient      []Ingredient `gorm:"many2many:recipe_ingredient;" json:"Ingredients"`
	Category        []Category   `gorm:"many2many:recipe_category;" json:"Categories"`
	Tags            []Tag        `gorm:"many2many:recipe_tag;" json:"Tags"`
	ImageName       string       `json:"ImageName" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// Instruction struct to hold instruction data
type Instruction struct {
	gorm.Model
	RecipeID    uint   `json:"RecipeID" example:"1"`
	Recipe      Recipe `gorm:"references:ID" json:"-"`
	StepNumber  int    `json:"StepNumber" example:"1"`
	Description string `json:"Description" example:"lorem ipsum dolor sit amet"`
}

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

// Category struct to hold category data
type Category struct {
	gorm.Model
	CategoryName string `gorm:"not null;unique" json:"CategoryName" example:"desserts"`
}

// Tag struct to hold tag data
type Tag struct {
	gorm.Model
	TagName string `gorm:"not null;unique" json:"TagName" example:"pies"`
}

type Unit struct {
	ID        uint   `gorm:"primaryKey;not null;unique;index" json:"ID" example:"1"`
	FullName  string `gorm:"not null;unique" json:"FullName" example:"Fluid ounce"`
	ShortName string `gorm:"not null;unique" json:"ShortName" example:"fl oz"`
}
