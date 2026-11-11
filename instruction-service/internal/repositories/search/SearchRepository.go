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

// Prepared for future expansion of instruction search beyond a single filter.
func (r *SearchRepository) SearchInstruction(request m.InstructionSearchRequest) (m.InstructionSearchResult, error) {
	var result m.InstructionSearchResult
	result.RecipeID = request.RecipeID

	query := r.db.Table("instructions").
		Select("instructions.instruction_id")

	if request.RecipeID != uuid.Nil {
		query = query.Where("instructions.recipe_id = ?", request.RecipeID)
	}

	if err := query.Pluck("instructions.instruction_id", &result.InstructionIDs).Error; err != nil {
		return m.InstructionSearchResult{}, err
	}

	return result, nil
}
