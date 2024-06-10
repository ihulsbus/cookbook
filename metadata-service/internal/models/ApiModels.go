package models

import (
	"gorm.io/gorm"
)

// Category struct to hold category data
type Category struct {
	gorm.Model
	CategoryName string `gorm:"not null;unique" json:"CategoryName" example:"desserts"`
}

// Tag struct to hold tag data
type Tag struct {
	gorm.Model
	TagName string `gorm:"not null;unique" json:"TagName" example:"pies"`
}
