package repositories

import (
	"errors"
	m "image-service/internal/models"

	"gorm.io/gorm"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{
		db: db,
	}
}

func (r ImageRepository) FindAll() ([]m.Image, error) {
	var images []m.Image

	if err := r.db.Find(&images).Error; err != nil {
		return nil, err
	}

	if len(images) <= 0 {
		return nil, errors.New("not found")
	}

	return images, nil
}

func (r ImageRepository) Find(image m.Image) (m.Image, error) {

	result := r.db.First(&image)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.Image{}, errors.New("not found")
		} else {
			return m.Image{}, result.Error
		}
	}

	return image, nil
}

func (r ImageRepository) Create(image m.Image) (m.Image, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&image).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return m.Image{}, err
	}

	return image, nil
}

func (r ImageRepository) Update(image m.Image) (m.Image, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Updates(&image).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return m.Image{}, err
	}

	return image, nil

}

func (r ImageRepository) Delete(image m.Image) error {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(&image).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
