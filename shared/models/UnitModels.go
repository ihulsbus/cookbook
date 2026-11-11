package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Dimension is the kind of quantity a Unit measures (KIT-37). Convert
// only crosses dimensions via an Ingredient's density/grams-per-piece
// (KIT-38) - within one Dimension, each Unit's ToBase factor is enough.
type Dimension string

const (
	DimensionMass   Dimension = "mass"
	DimensionVolume Dimension = "volume"
	DimensionCount  Dimension = "count"
)

// Valid reports whether d is one of the defined Dimension values.
func (d Dimension) Valid() bool {
	switch d {
	case DimensionMass, DimensionVolume, DimensionCount:
		return true
	default:
		return false
	}
}

var (
	// DefaultUnits is seeded on startup (ingredient-service's
	// initUnits, OnConflict DoNothing on the unique FullName/ShortName
	// columns). ToBase is each unit's conversion factor to its
	// Dimension's base unit - g for mass, ml for volume, piece for
	// count (KIT-37) - so g, ml, cloves and pcs all have ToBase 1.
	DefaultUnits = []Unit{
		// Volume (base: ml)
		{FullName: "Teaspoon", ShortName: "tsp", Dimension: DimensionVolume, ToBase: decimal.RequireFromString("4.92892")},
		{FullName: "Tablespoon", ShortName: "tbsp", Dimension: DimensionVolume, ToBase: decimal.RequireFromString("14.7868")},
		{FullName: "Fluid Ounce", ShortName: "fl oz", Dimension: DimensionVolume, ToBase: decimal.RequireFromString("29.5735")},
		{FullName: "Cup", ShortName: "c", Dimension: DimensionVolume, ToBase: decimal.RequireFromString("236.588")},
		{FullName: "Pint", ShortName: "pt", Dimension: DimensionVolume, ToBase: decimal.RequireFromString("473.176")},
		{FullName: "Quart", ShortName: "qt", Dimension: DimensionVolume, ToBase: decimal.RequireFromString("946.353")},
		{FullName: "Gallon", ShortName: "gal", Dimension: DimensionVolume, ToBase: decimal.RequireFromString("3785.41")},
		{FullName: "Milliliter", ShortName: "ml", Dimension: DimensionVolume, ToBase: decimal.NewFromInt(1)},
		{FullName: "Deciliter", ShortName: "dl", Dimension: DimensionVolume, ToBase: decimal.NewFromInt(100)},
		{FullName: "Liter", ShortName: "l", Dimension: DimensionVolume, ToBase: decimal.NewFromInt(1000)},
		// pinch has no standard definition - this is the commonly cited
		// approximation (~1/16 tsp), good enough since recipes calling
		// for a pinch were never precise to begin with.
		{FullName: "Pinch", ShortName: "pn", Dimension: DimensionVolume, ToBase: decimal.RequireFromString("0.31")},
		// Mass (base: g)
		{FullName: "Ounce", ShortName: "oz", Dimension: DimensionMass, ToBase: decimal.RequireFromString("28.3495")},
		{FullName: "Pound", ShortName: "lb", Dimension: DimensionMass, ToBase: decimal.RequireFromString("453.592")},
		{FullName: "Milligram", ShortName: "mg", Dimension: DimensionMass, ToBase: decimal.RequireFromString("0.001")},
		{FullName: "Gram", ShortName: "g", Dimension: DimensionMass, ToBase: decimal.NewFromInt(1)},
		{FullName: "Kilogram", ShortName: "kg", Dimension: DimensionMass, ToBase: decimal.NewFromInt(1000)},
		// Count (base: piece)
		{FullName: "Cloves", ShortName: "cloves", Dimension: DimensionCount, ToBase: decimal.NewFromInt(1)},
		{FullName: "Pieces", ShortName: "pcs", Dimension: DimensionCount, ToBase: decimal.NewFromInt(1)},
	}
)

type Unit struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FullName  string    `gorm:"not null;unique" json:"FullName" example:"Fluid ounce"`
	ShortName string    `gorm:"not null;unique" json:"ShortName" example:"fl oz"`
	// Dimension and ToBase (KIT-37) are what Convert (see Convert.go)
	// uses to translate a quantity from one Unit to another - directly
	// when both share a Dimension, via an Ingredient's density/
	// grams-per-piece (KIT-38) otherwise.
	Dimension Dimension       `gorm:"type:varchar(10);not null" json:"Dimension" example:"volume"`
	ToBase    decimal.Decimal `gorm:"type:numeric(16,6);not null" json:"ToBase" swaggertype:"string" example:"29.5735"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt  `gorm:"index"`
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
		Dimension: u.Dimension,
		ToBase:    u.ToBase,
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
	ID        uuid.UUID       `gorm:"primaryKey;not null;unique;index" json:"id" example:"1"`
	FullName  string          `gorm:"not null;unique" json:"full_name" example:"Fluid ounce"`
	ShortName string          `gorm:"not null;unique" json:"short_name" example:"fl oz"`
	Dimension Dimension       `json:"dimension" example:"volume"`
	ToBase    decimal.Decimal `json:"to_base" swaggertype:"string" example:"29.5735"`
}

func (u UnitDTO) ConvertFromDTO() Unit {
	return Unit{
		ID:        u.ID,
		FullName:  u.FullName,
		ShortName: u.ShortName,
		Dimension: u.Dimension,
		ToBase:    u.ToBase,
	}
}

func (c UnitDTO) ConvertAllFromDTO(unitDTOs []UnitDTO) []Unit {
	var data []Unit

	for _, unitDTO := range unitDTOs {
		data = append(data, unitDTO.ConvertFromDTO())
	}

	return data
}
