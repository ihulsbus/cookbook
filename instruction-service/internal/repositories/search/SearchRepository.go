package repositories

import (
	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"gorm.io/gorm"
)

type SearchRepository struct {
	db *gorm.DB
}

func NewSearchRepository(db *gorm.DB) *SearchRepository {
	return &SearchRepository{
		db: db,
	}
}

// this is currently prepared for possible future expansion on instruction usage
func (r *SearchRepository) SearchInstruction(request m.InstructionSearchRequest) (m.InstructionSearchResult, error) {
	var result m.InstructionSearchResult
	result.RecipeID = request.RecipeID

	// Start with base query. we do this on categories as all recipes need to have a category
	query := r.db.Table("instructions").
		Select("instructions.instruction_id")

	// Apply filters based on request
	if request.RecipeID != uuid.Nil {
		query = query.Where("instructions.recipe_id = ?", request.RecipeID)
	}

	// Scan results into a slice of recipe IDs
	if err := query.Pluck("instructions.instruction_id", &result.InstructionIDs).Error; err != nil {
		return m.InstructionSearchResult{}, err
	}

	return result, nil
}
