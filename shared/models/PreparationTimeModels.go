package models

import (
	"time"

	"github.com/google/uuid"
)

type PreparationTime struct {
	ID       uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primary_key" json:"id" example:"23582396-12a3-425b-a597-8a22052823da"`
	Duration time.Duration `gorm:"not null" json:"duration" example:"1h30m0s"`
}
