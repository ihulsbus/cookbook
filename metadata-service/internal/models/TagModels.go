package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Database model
type Tag struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string         `gorm:"type:varchar(100);unique;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (tag *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	tag.ID = uuid.New()
	return
}

// Association model
type RecipeTag struct {
	RecipeID  uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TagID     uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// DTO model
type TagDTO struct {
	ID   uuid.UUID `json:"id,omitempty" binding:"uuid"` // ID can be omitted for create operations
	Name string    `json:"name" binding:"required,min=1,max=255"`
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
