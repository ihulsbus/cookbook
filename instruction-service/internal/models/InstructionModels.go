package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Instruction struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Sequence    int            `gorm:"not null"`
	Description string         `gorm:"type:text; not null"`
	MediaID     uuid.UUID      `gorm:"type:uuid; not null"`
	EntityID    uuid.UUID      `gorm:"type:uuid; not null"`
	EntityType  string         `gorm:"type:text; not null"`
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
	MediaID     uuid.UUID `json:"media_id" example:"23582396-12a3-425b-a597-8a22052823da"`
	EntityID    uuid.UUID `json:"entity_id"`
	EntityType  string    `json:"entity_type"`
}

func (i Instruction) ConvertToDTO() InstructionDTO {
	return InstructionDTO{
		ID:          i.ID,
		Sequence:    i.Sequence,
		Description: i.Description,
		MediaID:     i.MediaID,
		EntityID:    i.EntityID,
		EntityType:  i.EntityType,
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

func (i InstructionDTO) ConvertAllFromDTO(instructions []InstructionDTO) []Instruction {
	var data []Instruction

	for _, i := range instructions {
		data = append(data, i.ConvertFromDTO())
	}

	return data
}
