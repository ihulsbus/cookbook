package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	DefaultCategories = []Category{
		{Name: "Fish"},
		{Name: "Meat"},
		{Name: "Vegetarian"},
		{Name: "Vegan"},
		{Name: "Soup"},
		{Name: "Stew"},
		{Name: "Curry"},
		{Name: "Pasta"},
		{Name: "Rice & Risotto"},
		{Name: "Salad"},
		{Name: "Bread"},
		{Name: "Fruit"},
		{Name: "Dessert"},
		{Name: "Baked treats"},
		{Name: "Sauce"},
		{Name: "Bouillon"},
		{Name: "Dough"},
	}
)

type Category struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuidv7();primaryKey"`
	Name      string         `gorm:"type:varchar(100);unique;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
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

type CategoryDTO struct {
	ID   uuid.UUID `json:"id,omitempty" binding:"uuid"` // ID can be omitted for create operations
	Name string    `json:"name" binding:"required,min=1,max=255"`
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
