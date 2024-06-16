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
