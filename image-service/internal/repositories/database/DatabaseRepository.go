package repositories

import (
	"errors"
	m "image-service/internal/models"

	"gorm.io/gorm"
)

type DatabaseRepository struct {
	db *gorm.DB
}

func NewDatabaseRepository(db *gorm.DB) *DatabaseRepository {
	return &DatabaseRepository{
		db: db,
	}
}

func (r DatabaseRepository) FindAll() ([]m.ImageData, error) {
	var images []m.ImageData

	if err := r.db.Find(&images).Error; err != nil {
		return nil, err
	}

	if len(images) <= 0 {
		return nil, errors.New("not found")
	}

	return images, nil
}

func (r DatabaseRepository) Find(image m.ImageData) (m.ImageData, error) {

	result := r.db.First(&image)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.ImageData{}, errors.New("not found")
		} else {
			return m.ImageData{}, result.Error
		}
	}

	return image, nil
}

func (r DatabaseRepository) Create(image m.ImageData) (m.ImageData, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&image).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return m.ImageData{}, err
	}

	return image, nil
}

func (r DatabaseRepository) Update(image m.ImageData) (m.ImageData, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Updates(&image).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return m.ImageData{}, err
	}

	return image, nil

}

func (r DatabaseRepository) Delete(image m.ImageData) error {

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
