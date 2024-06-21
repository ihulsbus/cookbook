package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Recipe struct to hold recipe data
type Recipe struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Name         string         `gorm:"not null" json:"RecipeName" example:"apple pie"`
	Description  string         `gorm:"size:65535;not null" json:"Description" example:"pie with apples"`
	ServingCount int            `gorm:"default:0" json:"ServingCount" example:"4"`
}

func (r Recipe) ConvertToDTO() RecipeDTO {
	return RecipeDTO{
		ID:           r.ID,
		Name:         r.Name,
		Description:  r.Description,
		ServingCount: r.ServingCount,
	}
}

func (c Recipe) ConvertAllToDTO(recipes []Recipe) []RecipeDTO {
	var data []RecipeDTO

	for _, recipe := range recipes {
		data = append(data, recipe.ConvertToDTO())
	}

	return data
}

type RecipeDTO struct {
	ID           uuid.UUID
	Name         string `gorm:"not null" json:"name" example:"apple pie"`
	Description  string `gorm:"size:65535;not null" json:"description" example:"pie with apples"`
	ServingCount int    `gorm:"default:0" json:"servingcount" example:"4"`
}

func (r RecipeDTO) ConvertFromDTO() Recipe {
	return Recipe{
		ID:           r.ID,
		Name:         r.Name,
		Description:  r.Description,
		ServingCount: r.ServingCount,
	}
}

func (t RecipeDTO) ConvertAllFromDTO(recipeDTOs []RecipeDTO) []Recipe {
	var data []Recipe

	for _, recipeDTO := range recipeDTOs {
		data = append(data, recipeDTO.ConvertFromDTO())
	}

	return data
}
