package repositories

import (
	"errors"

	m "metadata-service/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DifficultyLevelRepository struct {
	db *gorm.DB
}

func NewDifficultyLevelRepository(db *gorm.DB) *DifficultyLevelRepository {
	return &DifficultyLevelRepository{
		db: db,
	}
}

func (r *DifficultyLevelRepository) FindAll() ([]m.DifficultyLevel, error) {
	var difficultyLevels []m.DifficultyLevel

	if err := r.db.Preload(clause.Associations).Find(&difficultyLevels).Error; err != nil {
		return nil, err
	}

	if len(difficultyLevels) <= 0 {
		return nil, errors.New("not found")
	}

	return difficultyLevels, nil
}

func (r *DifficultyLevelRepository) FindSingle(difficultyLevel m.DifficultyLevel) (m.DifficultyLevel, error) {

	result := r.db.First(&difficultyLevel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.DifficultyLevel{}, errors.New("not found")
		} else {
			return m.DifficultyLevel{}, result.Error
		}
	}

	return difficultyLevel, nil
}

func (r *DifficultyLevelRepository) Create(difficultyLevel m.DifficultyLevel) (m.DifficultyLevel, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Create(&difficultyLevel).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return m.DifficultyLevel{}, err
	}

	return difficultyLevel, nil
}

func (r *DifficultyLevelRepository) Update(difficultyLevel m.DifficultyLevel) (m.DifficultyLevel, error) {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Updates(&difficultyLevel).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return m.DifficultyLevel{}, err
	}

	return difficultyLevel, nil
}

func (r *DifficultyLevelRepository) Delete(difficultyLevel m.DifficultyLevel) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(&difficultyLevel).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
