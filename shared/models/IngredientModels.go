package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Ingredient struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuidv7();primaryKey"`
	Name string    `gorm:"unique; not null" json:"IngredientName"`
	// DensityGPerMl and GramsPerPiece let Convert (see Convert.go) cross
	// between mass, volume and count for this ingredient. Nil when
	// unknown, in which case Convert returns ErrNoConversion instead of
	// guessing.
	DensityGPerMl *decimal.Decimal `gorm:"type:numeric(10,4)" json:"DensityGPerMl,omitempty"`
	GramsPerPiece *decimal.Decimal `gorm:"type:numeric(10,4)" json:"GramsPerPiece,omitempty"`
	// CreatedBy is the OIDC sub of whoever created this catalog entry -
	// any authenticated user can, not just cookbook-admin.
	CreatedBy uuid.UUID      `gorm:"type:uuid;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// NormalizedName has no Go field on purpose: it's a Postgres STORED
	// GENERATED column (created via raw SQL, not GORM AutoMigrate), with
	// a unique index for exact-duplicate prevention and a pg_trgm GIN
	// index for similarity search. Queried by raw column name; nothing
	// here needs to read its value back into Go.
}

func (ingredient *Ingredient) BeforeCreate(tx *gorm.DB) (err error) {
	// Time-ordered, not the DB's uuidv7() default; guarded so
	// tests/imports that already supply an ID aren't overridden.
	if ingredient.ID == uuid.Nil {
		ingredient.ID = NewID()
	}
	return
}

func (i Ingredient) ConvertToDTO() IngredientDTO {
	return IngredientDTO{
		ID:            i.ID,
		Name:          i.Name,
		DensityGPerMl: i.DensityGPerMl,
		GramsPerPiece: i.GramsPerPiece,
		CreatedAt:     i.CreatedAt,
		UpdatedAt:     i.UpdatedAt,
	}
}

func (c Ingredient) ConvertAllToDTO(ingredients []Ingredient) []IngredientDTO {
	var data []IngredientDTO

	for _, ingredient := range ingredients {
		data = append(data, ingredient.ConvertToDTO())
	}

	return data
}

type IngredientDTO struct {
	ID            uuid.UUID        `json:"id" example:"23582396-12a3-425b-a597-8a22052823da"`
	Name          string           `json:"name" example:"asparagus"`
	DensityGPerMl *decimal.Decimal `json:"density_g_per_ml,omitempty" swaggertype:"string" example:"1.03"`
	GramsPerPiece *decimal.Decimal `json:"grams_per_piece,omitempty" swaggertype:"string" example:"50"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

func (i IngredientDTO) ConvertFromDTO() Ingredient {
	return Ingredient{
		ID:            i.ID,
		Name:          i.Name,
		DensityGPerMl: i.DensityGPerMl,
		GramsPerPiece: i.GramsPerPiece,
		CreatedAt:     i.CreatedAt,
		UpdatedAt:     i.UpdatedAt,
	}
}

func (c IngredientDTO) ConvertAllFromDTO(ingredientDTOs []IngredientDTO) []Ingredient {
	var data []Ingredient

	for _, ingredientDTO := range ingredientDTOs {
		data = append(data, ingredientDTO.ConvertFromDTO())
	}

	return data
}
