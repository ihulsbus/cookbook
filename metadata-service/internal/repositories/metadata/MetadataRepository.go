package repositories

import (
	"errors"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"

	"gorm.io/gorm"
)

type RecipeMetadataRepository struct {
	db *gorm.DB
}

func NewRecipeMetadataRepository(db *gorm.DB) *RecipeMetadataRepository {
	return &RecipeMetadataRepository{db: db}
}

func (r *RecipeMetadataRepository) FindAll(pagination m.PaginationRequest) ([]m.RecipeMetadata, int64, error) {
	var metas []m.RecipeMetadata
	var total int64

	if err := r.db.Model(&m.RecipeMetadata{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Preload("Categories").
		Preload("CuisineType").
		Preload("DifficultyLevel").
		Preload("Tags").
		Limit(pagination.Limit).Offset(pagination.Offset()).
		Find(&metas).Error
	if err != nil {
		return nil, 0, err
	}

	return metas, total, nil
}

func (r *RecipeMetadataRepository) FindSingle(recipeID uuid.UUID) (*m.RecipeMetadata, error) {
	var meta m.RecipeMetadata
	err := r.db.
		Preload("Categories").
		Preload("CuisineType").
		Preload("DifficultyLevel").
		Preload("Tags").
		First(&meta, "recipe_id = ?", recipeID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &meta, nil
}

func (r *RecipeMetadataRepository) Create(meta *m.RecipeMetadata) (*m.RecipeMetadata, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		// FullSaveAssociations handles the many2many join tables automatically.
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&meta).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return meta, nil
}

func (r *RecipeMetadataRepository) Update(meta *m.RecipeMetadata) (*m.RecipeMetadata, error) {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// FullSaveAssociations ensures tags/categories are replaced, not appended.
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&meta).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return meta, nil
}

// Delete also cascades through the join tables.
func (r *RecipeMetadataRepository) Delete(meta *m.RecipeMetadata) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&meta).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
