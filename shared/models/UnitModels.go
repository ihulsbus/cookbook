package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	DefaultUnits = []Unit{
		// US units
		{FullName: "Teaspoon", ShortName: "tsp"},
		{FullName: "Tablespoon", ShortName: "tbsp"},
		{FullName: "Fluid Ounce", ShortName: "fl oz"},
		{FullName: "Ounce", ShortName: "oz"},
		{FullName: "Pound", ShortName: "lb"},
		{FullName: "Cup", ShortName: "c"},
		{FullName: "Pint", ShortName: "pt"},
		{FullName: "Quart", ShortName: "qt"},
		{FullName: "Gallon", ShortName: "gal"},
		// Metric units
		{FullName: "Milliliter", ShortName: "ml"},
		{FullName: "Deciliter", ShortName: "dl"},
		{FullName: "Liter", ShortName: "l"},
		{FullName: "Milligram", ShortName: "mg"},
		{FullName: "Gram", ShortName: "g"},
		{FullName: "Kilogram", ShortName: "kg"},
		// Generic units
		{FullName: "Pinch", ShortName: "pn"},
		{FullName: "Cloves", ShortName: "cloves"},
		{FullName: "Pieces", ShortName: "pcs"},
	}
)

type Unit struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FullName  string         `gorm:"not null;unique" json:"FullName" example:"Fluid ounce"`
	ShortName string         `gorm:"not null;unique" json:"ShortName" example:"fl oz"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (unit *Unit) BeforeCreate(tx *gorm.DB) (err error) {
	unit.ID = uuid.New()
	return
}

func (u Unit) ConvertToDTO() UnitDTO {
	return UnitDTO{
		ID:        u.ID,
		FullName:  u.FullName,
		ShortName: u.ShortName,
	}
}

func (c Unit) ConvertAllToDTO(units []Unit) []UnitDTO {
	var data []UnitDTO

	for _, unit := range units {
		data = append(data, unit.ConvertToDTO())
	}

	return data
}

type UnitDTO struct {
	ID        uuid.UUID `gorm:"primaryKey;not null;unique;index" json:"id" example:"1"`
	FullName  string    `gorm:"not null;unique" json:"full_name" example:"Fluid ounce"`
	ShortName string    `gorm:"not null;unique" json:"short_name" example:"fl oz"`
}

func (u UnitDTO) ConvertFromDTO() Unit {
	return Unit{
		ID:        u.ID,
		FullName:  u.FullName,
		ShortName: u.ShortName,
	}
}

func (c UnitDTO) ConvertAllFromDTO(unitDTOs []UnitDTO) []Unit {
	var data []Unit

	for _, unitDTO := range unitDTOs {
		data = append(data, unitDTO.ConvertFromDTO())
	}

	return data
}
