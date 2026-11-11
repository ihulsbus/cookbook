package models

import (
	"bytes"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ImageData) TableName() string {
	return "images"
}

// ImageData is an image's metadata; it does not contain the file itself.
type ImageData struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuidv7();primary_key"`
	EntityType string         `gorm:"type:varchar(50);not null"` // e.g., "recipe" or "ingredient"
	EntityID   uuid.UUID      `gorm:"type:uuid;not null"`
	Size       int64          `gorm:"not null"`                  // bytes
	Type       string         `gorm:"type:varchar(50);not null"` // e.g., "image/jpeg"
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (i ImageData) ConvertToDTO() ImageDataDTO {
	return ImageDataDTO{
		ID:         i.ID,
		EntityType: i.EntityType,
		EntityID:   i.EntityID,
		Size:       i.Size,
		Type:       i.Type,
	}
}

func (i ImageData) ConvertAllToDTO(images []ImageData) []ImageDataDTO {
	var data []ImageDataDTO

	for _, image := range images {
		data = append(data, image.ConvertToDTO())
	}

	return data
}

type ImageDataDTO struct {
	ID         uuid.UUID `json:"id"`
	EntityType string    `json:"entity_type"`
	EntityID   uuid.UUID `json:"entity_id"`
	Size       int64     `json:"size"`
	Type       string    `json:"type"`
	// URL is a short-lived signed GET URL, set by the service layer only
	// after access has been checked. Empty if not generated (e.g. a list
	// response that doesn't presign every row).
	URL string `json:"url,omitempty"`
}

// WithURL returns a copy of i carrying url; never populated here.
func (i ImageDataDTO) WithURL(url string) ImageDataDTO {
	i.URL = url
	return i
}

func (i ImageDataDTO) ConvertFromDTO() ImageData {
	return ImageData{
		ID:         i.ID,
		EntityType: i.EntityType,
		EntityID:   i.EntityID,
		Size:       i.Size,
		Type:       i.Type,
	}
}

type ImageFile struct {
	ID         uuid.UUID // No json tag as it will never be filled through an unmarshal
	EntityType string    `json:"entity_type"`
	EntityID   uuid.UUID `json:"entity_id"`
	Size       int64     `json:"size"`
	Type       string    `json:"type"`
	File       bytes.Buffer
}

func (i ImageFile) ConvertToDTO() ImageFileDTO {
	return ImageFileDTO(i)
}

type ImageFileDTO struct {
	ID         uuid.UUID // No json tag as it will never be filled through an unmarshal
	EntityType string    `json:"entity_type"`
	EntityID   uuid.UUID `json:"entity_id"`
	Size       int64     `json:"size"`
	Type       string    `json:"type"`
	File       bytes.Buffer
}

func (i ImageFileDTO) ConvertFromDTO() ImageFile {
	return ImageFile(i)
}
