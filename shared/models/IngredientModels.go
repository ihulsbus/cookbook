package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Ingredient struct to hold ingredient data
type Ingredient struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string    `gorm:"unique; not null" json:"IngredientName"`
	// DensityGPerMl and GramsPerPiece (KIT-38) let Convert (see
	// Convert.go) cross between mass, volume and count for this
	// ingredient. Both nullable: left nil when unknown, in which case
	// Convert returns ErrNoConversion rather than guessing - the UI
	// then keeps showing the original unit.
	DensityGPerMl *decimal.Decimal `gorm:"type:numeric(10,4)" json:"DensityGPerMl,omitempty"`
	GramsPerPiece *decimal.Decimal `gorm:"type:numeric(10,4)" json:"GramsPerPiece,omitempty"`
	CreatedAt     time.Time        `gorm:"autoCreateTime"`
	UpdatedAt     time.Time        `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt   `gorm:"index"`
}

func (ingredient *Ingredient) BeforeCreate(tx *gorm.DB) (err error) {
	ingredient.ID = uuid.New()
	return
}

func (i Ingredient) ConvertToDTO() IngredientDTO {
	return IngredientDTO{
		ID:            i.ID,
		Name:          i.Name,
		DensityGPerMl: i.DensityGPerMl,
		GramsPerPiece: i.GramsPerPiece,
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
}

func (i IngredientDTO) ConvertFromDTO() Ingredient {
	return Ingredient{
		ID:            i.ID,
		Name:          i.Name,
		DensityGPerMl: i.DensityGPerMl,
		GramsPerPiece: i.GramsPerPiece,
	}
}

func (c IngredientDTO) ConvertAllFromDTO(ingredientDTOs []IngredientDTO) []Ingredient {
	var data []Ingredient

	for _, ingredientDTO := range ingredientDTOs {
		data = append(data, ingredientDTO.ConvertFromDTO())
	}

	return data
}
