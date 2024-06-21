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

func (ingredient *Ingredient) BeforeCreate(tx *gorm.DB) (err error) {
	ingredient.ID = uuid.New()
	return
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

type IngredientDTO struct {
	ID   uuid.UUID `json:"id" example:"23582396-12a3-425b-a597-8a22052823da"`
	Name string    `json:"name" example:"asparagus"`
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
