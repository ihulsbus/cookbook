package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Instruction struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Sequence    int            `gorm:"not null"`
	Description string         `gorm:"type:text;not null"`
	MediaID     uuid.UUID      `gorm:"type:uuid; not null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (instruction *Instruction) BeforeCreate(tx *gorm.DB) (err error) {
	instruction.ID = uuid.New()
	return
}

type InstructionDTO struct {
	ID          uuid.UUID `json:"id" example:"23582396-12a3-425b-a597-8a22052823da"`
	Sequence    int       `json:"sequence" example:"1"`
	Description string    `json:"description" example:"description"`
	MediaID     uuid.UUID `json:"media_url" example:"23582396-12a3-425b-a597-8a22052823da"`
}

func (i Instruction) ConvertToDTO() InstructionDTO {
	return InstructionDTO{
		ID:          i.ID,
		Sequence:    i.Sequence,
		Description: i.Description,
		MediaID:     i.MediaID,
	}
}

func (i Instruction) ConvertAllToDTO(instructions []Instruction) []InstructionDTO {
	var data []InstructionDTO

	for _, i := range instructions {
		data = append(data, i.ConvertToDTO())
	}

	return data
}

func (i InstructionDTO) ConvertFromDTO() Instruction {
	return Instruction{
		ID:          i.ID,
		Sequence:    i.Sequence,
		Description: i.Description,
		MediaID:     i.MediaID,
	}
}

// Association model
type RecipeInstruction struct {
	RecipeID      uuid.UUID      `gorm:"type:uuid;primaryKey"`
	InstructionID uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
