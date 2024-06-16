package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Database model
type Category struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string         `gorm:"type:varchar(100);unique;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (category *Category) BeforeCreate(tx *gorm.DB) (err error) {
	category.ID = uuid.New()
	return
}

// Association model
type RecipeCategory struct {
	RecipeID   uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CategoryID uuid.UUID      `gorm:"type:uuid;primaryKey"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// DTO model
type CategoryDTO struct {
	ID   uuid.UUID `json:"id,omitempty" binding:"uuid"` // ID can be omitted for create operations
	Name string    `json:"name" binding:"required,min=1,max=255"`
}

func (c Category) ConvertToDTO() CategoryDTO {
	return CategoryDTO{
		ID:   c.ID,
		Name: c.Name,
	}
}

func (c Category) ConvertAllToDTO(categories []Category) []CategoryDTO {
	var data []CategoryDTO

	for _, category := range categories {
		data = append(data, category.ConvertToDTO())
	}

	return data
}

func (c CategoryDTO) ConvertFromDTO() Category {
	return Category{
		ID:   c.ID,
		Name: c.Name,
	}
}

func (t CategoryDTO) ConvertAllFromDTO(categories []CategoryDTO) []Category {
	var data []Category

	for _, category := range categories {
		data = append(data, category.ConvertFromDTO())
	}

	return data
}
