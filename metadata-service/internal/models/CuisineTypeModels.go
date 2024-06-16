package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Database model
type CuisineType struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string         `gorm:"type:varchar(100);unique;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (cuisineType *CuisineType) BeforeCreate(tx *gorm.DB) (err error) {
	cuisineType.ID = uuid.New()
	return
}

// Association model
type RecipeCuisineType struct {
	RecipeID      uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CuisineTypeID uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// DTO model
type CuisineTypeDTO struct {
	ID   uuid.UUID `json:"id,omitempty" binding:"uuid"` // ID can be omitted for create operations
	Name string    `json:"name" binding:"required,min=1,max=255"`
}

func (c CuisineType) ConvertToDTO() CuisineTypeDTO {
	return CuisineTypeDTO{
		ID:   c.ID,
		Name: c.Name,
	}
}

func (c CuisineType) ConvertAllToDTO(cuisineTypes []CuisineType) []CuisineTypeDTO {
	var data []CuisineTypeDTO

	for _, cuisineType := range cuisineTypes {
		data = append(data, cuisineType.ConvertToDTO())
	}

	return data
}

func (c CuisineTypeDTO) ConvertFromDTO() CuisineType {
	return CuisineType{
		ID:   c.ID,
		Name: c.Name,
	}
}

func (c CuisineTypeDTO) ConvertAllFromDTO(cuisineTypes []CuisineTypeDTO) []CuisineType {
	var data []CuisineType

	for _, cuisineType := range cuisineTypes {
		data = append(data, cuisineType.ConvertFromDTO())
	}

	return data
}
