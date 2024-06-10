package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Instruction struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	RecipeID    uuid.UUID      `gorm:"type:uuid;not null"`
	Sequence    int            `gorm:"not null"`
	Description string         `gorm:"type:text;not null"`
	MediaURL    string         `gorm:"type:text"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (instruction *Instruction) BeforeCreate(tx *gorm.DB) (err error) {
	instruction.ID = uuid.New()
	return
}

// YEET
type User struct{}
