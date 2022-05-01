package models

import (
	"mime/multipart"

	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model
	Title             string              `gorm:"not null" json:"title"`
	Description       string              `gorm:"not null" json:"description"`
	Method            string              `json:"method"`
	PrepTime          int                 `gorm:"default:0" json:"preptime"`
	CookTime          int                 `gorm:"default:0" json:"cooktime"`
	Persons           int                 `gorm:"default:0" json:"persons"`
	Ingredients       []Ingredient        `gorm:"many2many:Recipe_Ingredients;foreignKey:ID" json:"ingredients"`
	IngredientAmounts []Recipe_Ingredient `json:"ingredientamounts"`
	Tags              []Tags              `gorm:"many2many:Recipe_Tags;foreignKey:ID" json:"tags"`
}

type Tags struct {
	ID   int    `gorm:"primaryKey;serial;unique;not null;autoIncrement" json:"id"`
	Name string `gorm:"not null;unique" json:"name"`
}

type RecipeFile struct {
	ID   int
	File *multipart.FileHeader
}

type Recipe_Ingredient struct {
	RecipeID     int    `json:"recipeid"`
	IngredientID int    `json:"ingredientid"`
	Amount       int    `gorm:"default:0" json:"amount"`
	Unit         string `gorm:"not null" json:"unit"`
}

type Recipe_Tag struct {
	RecipeID int `json:"recipeid"`
	TagsID   int `json:"ingredientid"`
}

type RecipeDTO struct {
	ID                uint                `json:"id"`
	Title             string              `json:"title"`
	Description       string              `json:"description"`
	Method            string              `json:"method"`
	PrepTime          int                 `json:"preptime"`
	CookTime          int                 `json:"cooktime"`
	Persons           int                 `json:"persons"`
	Ingredients       []Ingredient        `json:"ingredients"`
	IngredientAmounts []Recipe_Ingredient `json:"ingredientamounts"`
	Tags              []Tags              `json:"tags"`
}

// ConvertToDTO converts a single m.Recipe to a single m.RecipeDTO object
func (r Recipe) ConvertToDTO() RecipeDTO {
	return RecipeDTO{
		ID:                r.ID,
		Title:             r.Title,
		Description:       r.Description,
		Method:            r.Method,
		PrepTime:          r.PrepTime,
		CookTime:          r.CookTime,
		Persons:           r.Persons,
		Ingredients:       r.Ingredients,
		IngredientAmounts: r.IngredientAmounts,
		Tags:              r.Tags,
	}
}

// ConvertAllToDTO converts a slice of m.Recipe to a slice of m.RecipeDTO objects
func (r Recipe) ConvertAllToDTO(recipes []Recipe) []RecipeDTO {
	var data []RecipeDTO

	for _, recipe := range recipes {
		data = append(data, recipe.ConvertToDTO())
	}

	return data
}
