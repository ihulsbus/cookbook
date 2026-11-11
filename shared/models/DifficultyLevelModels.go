package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	DefaultDifficultyLevels = []DifficultyLevel{
		{Level: 1},
		{Level: 2},
		{Level: 3},
		{Level: 4},
		{Level: 5},
	}
)

type DifficultyLevel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuidv7();primaryKey"`
	Level     int            `gorm:"type:integer;unique;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
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

type DifficultyLevelDTO struct {
	ID    uuid.UUID `json:"id,omitempty" binding:"uuid"`
	Level int       `json:"level" binding:"required,numeric,min=1,max=5"`
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
