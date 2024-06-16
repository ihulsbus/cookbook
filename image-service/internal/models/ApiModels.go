package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	EntityType string         `gorm:"type:varchar(50);not null"` // e.g., "recipe" or "ingredient"
	EntityID   uuid.UUID      `gorm:"type:uuid;not null"`
	URL        string         `gorm:"type:text;not null"`
	Size       int64          `gorm:"not null"`                  // size in bytes
	Type       string         `gorm:"type:varchar(50);not null"` // e.g., "image/jpeg"
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (image *Image) BeforeCreate(tx *gorm.DB) (err error) {
	image.ID = uuid.New()
	return
}
