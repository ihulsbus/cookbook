package models

import "mime/multipart"

type Recipe struct {
	ID                int                 `gorm:"primaryKey;serial;unique;not null;autoIncrement" json:"id"`
	Title             string              `gorm:"not null" json:"title"`
	Description       string              `gorm:"not null" json:"description"`
	Method            string              `json:"method"`
	Ingredients       []Ingredient        `gorm:"many2many:Recipe_Ingredients;foreignKey:ID" json:"ingredients"`
	PrepTime          int                 `gorm:"default:0" json:"preptime"`
	CookTime          int                 `gorm:"default:0" json:"cooktime"`
	Amount_Persons    int                 `gorm:"default:0" json:"persons"`
	IngredientAmounts []Recipe_Ingredient `json:"ingredientamounts"`
}

type RecipeDTO struct {
	ID             int          `json:"id"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Ingredients    []Ingredient `json:"ingredients"`
	Method         string       `json:"method"`
	PrepTime       int          `json:"preptime"`
	CookTime       int          `json:"cooktime"`
	Amount_Persons int          `json:"persons"`
}

type RecipeFile struct {
	ID   int
	File *multipart.FileHeader
}

// ConvertToDTO converts a single m.Recipe to a single m.RecipeDTO object
func (r Recipe) ConvertToDTO() RecipeDTO {
	return RecipeDTO{
		ID:             r.ID,
		Title:          r.Title,
		Description:    r.Description,
		Ingredients:    r.Ingredients,
		Method:         r.Method,
		PrepTime:       r.PrepTime,
		CookTime:       r.CookTime,
		Amount_Persons: r.Amount_Persons,
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

// ConvertFromDTO converts a single m.RecipeDTO to a single m.Recipe object
func (r RecipeDTO) ConvertFromDTO() Recipe {
	return Recipe{
		ID:             r.ID,
		Title:          r.Title,
		Description:    r.Description,
		Method:         r.Method,
		PrepTime:       r.PrepTime,
		CookTime:       r.CookTime,
		Amount_Persons: r.Amount_Persons,
	}
}

// ConvertAllFromDTO converts a slice of m.RecipeDTO to a slice of m.Recipe objects
func (r RecipeDTO) ConvertAllFromDTO(recipes []RecipeDTO) []Recipe {
	var data []Recipe

	for _, recipe := range recipes {
		data = append(data, recipe.ConvertFromDTO())
	}

	return data
}
