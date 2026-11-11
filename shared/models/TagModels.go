package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuidv7();primaryKey"`
	Name      string         `gorm:"type:varchar(100);unique;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (t Tag) ConvertToDTO() TagDTO {
	return TagDTO{
		ID:   t.ID,
		Name: t.Name,
	}
}

func (t Tag) ConvertAllToDTO(tags []Tag) []TagDTO {
	var data []TagDTO

	for _, tag := range tags {
		data = append(data, tag.ConvertToDTO())
	}

	return data
}

type TagDTO struct {
	ID   uuid.UUID `json:"id,omitempty" binding:"uuid"` // ID can be omitted for create operations
	Name string    `json:"name" binding:"required,min=1,max=255"`
}

func (t TagDTO) ConvertFromDTO() Tag {
	return Tag{
		ID:   t.ID,
		Name: t.Name,
	}
}

func (t TagDTO) ConvertAllFromDTO(tags []TagDTO) []Tag {
	var data []Tag

	for _, tag := range tags {
		data = append(data, tag.ConvertFromDTO())
	}

	return data
}
