package repositories

import (
	"errors"

	m "metadata-service/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{
		db: db,
	}
}

func (r *TagRepository) FindAll() ([]m.Tag, error) {
	var tags []m.Tag

	if err := r.db.Preload(clause.Associations).Find(&tags).Error; err != nil {
		return nil, err
	}

	if len(tags) <= 0 {
		return nil, errors.New("not found")
	}

	return tags, nil
}

func (r *TagRepository) FindSingle(tag m.Tag) (m.Tag, error) {

	result := r.db.First(&tag)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.Tag{}, errors.New("not found")
		} else {
			return m.Tag{}, result.Error
		}
	}

	return tag, nil
}

func (r *TagRepository) Create(tag m.Tag) (m.Tag, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Create(&tag).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return m.Tag{}, err
	}

	return tag, nil
}

func (r *TagRepository) Update(tag m.Tag) (m.Tag, error) {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Updates(&tag).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return m.Tag{}, err
	}

	return tag, nil
}

func (r *TagRepository) Delete(tag m.Tag) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(&tag).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
