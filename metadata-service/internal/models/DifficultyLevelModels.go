package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Database model
type DifficultyLevel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Level     string         `gorm:"type:varchar(50);unique;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Association model
type RecipeDifficultyLevel struct {
	RecipeID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	DifficultyLevelID uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

// DTO model
type DifficultyLevelDTO struct {
	ID   uuid.UUID `json:"id,omitempty" binding:"uuid"` // ID can be omitted for create operations
	Name string    `json:"name" binding:"required,oneof='Easy' 'Medium' 'Hard'"`
}
