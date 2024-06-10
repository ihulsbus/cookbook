package models

import "gorm.io/gorm"

// Recipe struct to hold recipe data
type Recipe struct {
	gorm.Model
	RecipeName      string `gorm:"not null" json:"RecipeName" example:"apple pie"`
	Description     string `gorm:"size:65535;not null" json:"Description" example:"pie with apples"`
	DifficultyLevel int    `gorm:"not null" json:"DifficultyLevel" example:"1"`
	CookingTime     int    `gorm:"default:0" json:"CookTime" example:"23"`
	ServingCount    int    `gorm:"default:0" json:"ServingCount" example:"4"`
	ImageName       string `json:"ImageName" example:"123e4567-e89b-12d3-a456-426614174000"`
	AuthorID        string `json:"Author" example:"123e4567-e89b-12d3-a456-426614174000"`
}
