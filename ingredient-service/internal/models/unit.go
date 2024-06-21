package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	ID        uuid.UUID `gorm:"primaryKey;not null;unique;index" json:"ID" example:"1"`
	FullName  string    `gorm:"not null;unique" json:"FullName" example:"Fluid ounce"`
	ShortName string    `gorm:"not null;unique" json:"ShortName" example:"fl oz"`
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
