package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Database model
type DifficultyLevel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Level     int            `gorm:"type:integer;unique;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (difficultyLevel *DifficultyLevel) BeforeCreate(tx *gorm.DB) (err error) {
	difficultyLevel.ID = uuid.New()
	return
}

// Association model
type RecipeDifficultyLevel struct {
	RecipeID          uuid.UUID      `gorm:"type:uuid;primaryKey;unique"`
	DifficultyLevelID uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

// DTO model
type DifficultyLevelDTO struct {
	ID    uuid.UUID `json:"id,omitempty" binding:"uuid"` // ID can be omitted for create operations
	Level int       `json:"name" binding:"required,numeric,min=1,max=5"`
}

func (d DifficultyLevel) ConvertToDTO() DifficultyLevelDTO {
	return DifficultyLevelDTO{
		ID:    d.ID,
		Level: d.Level,
	}
}

func (d DifficultyLevel) ConvertAllToDTO(difficultyLevels []DifficultyLevel) []DifficultyLevelDTO {
	var data []DifficultyLevelDTO

	for _, level := range difficultyLevels {
		data = append(data, level.ConvertToDTO())
	}

	return data
}

func (d DifficultyLevelDTO) ConvertFromDTO() DifficultyLevel {
	return DifficultyLevel{
		ID:    d.ID,
		Level: d.Level,
	}
}

func (d DifficultyLevelDTO) ConvertAllFromDTO(difficultyLevels []DifficultyLevelDTO) []DifficultyLevel {
	var data []DifficultyLevel

	for _, level := range difficultyLevels {
		data = append(data, level.ConvertFromDTO())
	}

	return data
}
