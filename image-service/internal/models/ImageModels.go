package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	EntityType string         `gorm:"type:varchar(50);not null"` // e.g., "recipe" or "ingredient"
	EntityID   uuid.UUID      `gorm:"type:uuid;not null"`
	Size       int64          `gorm:"not null"`                  // size in bytes
	Type       string         `gorm:"type:varchar(50);not null"` // e.g., "image/jpeg"
	File       multipart.File `gorm:"-:all"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (i Image) ConvertToDTO() ImageDTO {
	return ImageDTO{
		ID:         i.ID,
		EntityType: i.EntityType,
		EntityID:   i.EntityID,
		Size:       i.Size,
		Type:       i.Type,
		File:       i.File,
	}
}

func (i Image) ConvertAllToDTO(images []Image) []ImageDTO {
	var data []ImageDTO

	for _, image := range images {
		data = append(data, image.ConvertToDTO())
	}

	return data
}

type ImageDTO struct {
	ID         uuid.UUID `json:"id"`
	EntityType string    `json:"entity_type"`
	EntityID   uuid.UUID `json:"entity_id"`
	Size       int64     `json:"size"`
	Type       string    `json:"type"`
	File       multipart.File
}

func (i ImageDTO) ConvertFromDTO() Image {
	return Image{
		ID:         i.ID,
		EntityType: i.EntityType,
		EntityID:   i.EntityID,
		Size:       i.Size,
		Type:       i.Type,
		File:       i.File,
	}
}
