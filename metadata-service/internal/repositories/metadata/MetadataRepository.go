package repositories

import (
	"errors"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type RecipeMetadataRepository struct {
	db *gorm.DB
}

// NewRecipeMetadataRepository constructs a new repository
func NewRecipeMetadataRepository(db *gorm.DB) *RecipeMetadataRepository {
	return &RecipeMetadataRepository{db: db}
}

// FindAll loads all metadata entries, including associations
func (r *RecipeMetadataRepository) FindAll() (*[]m.RecipeMetadata, error) {
	var metas []m.RecipeMetadata
	err := r.db.
		Preload(clause.Associations).
		Find(&metas).Error
	if err != nil {
		return nil, err
	}
	if len(metas) == 0 {
		return nil, errors.New("not found")
	}
	return &metas, nil
}

// FindSingle fetches metadata for a given RecipeID (expects meta.RecipeID set)
func (r *RecipeMetadataRepository) FindSingle(recipeID uuid.UUID) (*m.RecipeMetadata, error) {
	var meta m.RecipeMetadata
	err := r.db.
		Preload(clause.Associations).
		First(&meta, "recipe_id = ?", recipeID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &meta, nil
}

// Create inserts a new RecipeMetadata record with associations
func (r *RecipeMetadataRepository) Create(meta *m.RecipeMetadata) (*m.RecipeMetadata, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		// Create meta; many 2 many/join-tables will be handled automatically
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&meta).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return meta, nil
}

// Update saves changes to metadata and its associations
func (r *RecipeMetadataRepository) Update(meta *m.RecipeMetadata) (*m.RecipeMetadata, error) {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// FullSaveAssociations ensures tags/categories are replaced
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&meta).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return meta, nil
}

// Delete removes a metadata entry (by PK) and cascades through join tables
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
