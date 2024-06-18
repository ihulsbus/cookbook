package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Database model
type PreparationTime struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Duration  int            `gorm:"type:integer;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (preparationTime *PreparationTime) BeforeCreate(tx *gorm.DB) (err error) {
	preparationTime.ID = uuid.New()
	return
}

// Association model
type RecipePreparationTime struct {
	RecipeID          uuid.UUID      `gorm:"type:uuid;primaryKey;unique"`
	PreparationTimeID uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

// DTO model
type PreparationTimeDTO struct {
	ID       uuid.UUID `json:"id,omitempty" binding:"uuid"` // ID can be omitted for create operations
	Duration int       `json:"name" binding:"required,min=1,max=100"`
}

func (p PreparationTime) ConvertToDTO() PreparationTimeDTO {
	return PreparationTimeDTO{
		ID:       p.ID,
		Duration: p.Duration,
	}
}

func (p PreparationTime) ConvertAllToDTO(cuisineTypes []PreparationTime) []PreparationTimeDTO {
	var data []PreparationTimeDTO

	for _, cuisineType := range cuisineTypes {
		data = append(data, cuisineType.ConvertToDTO())
	}

	return data
}

func (p PreparationTimeDTO) ConvertFromDTO() PreparationTime {
	return PreparationTime{
		ID:       p.ID,
		Duration: p.Duration,
	}
}

func (p PreparationTimeDTO) ConvertAllFromDTO(cuisineTypes []PreparationTimeDTO) []PreparationTime {
	var data []PreparationTime

	for _, cuisineType := range cuisineTypes {
		data = append(data, cuisineType.ConvertFromDTO())
	}

	return data
}
