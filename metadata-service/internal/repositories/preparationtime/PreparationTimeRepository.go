package repositories

import (
	"errors"

	m "github.com/ihulsbus/cookbook/shared/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PreparationTimeRepository struct {
	db *gorm.DB
}

func NewPreparationTimeRepository(db *gorm.DB) *PreparationTimeRepository {
	return &PreparationTimeRepository{
		db: db,
	}
}

func (r *PreparationTimeRepository) FindAll(pagination m.PaginationRequest) ([]m.PreparationTime, int64, error) {
	var preparationTimes []m.PreparationTime
	var total int64

	if err := r.db.Model(&m.PreparationTime{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload(clause.Associations).Limit(pagination.Limit).Offset(pagination.Offset()).Find(&preparationTimes).Error; err != nil {
		return nil, 0, err
	}

	return preparationTimes, total, nil
}

func (r *PreparationTimeRepository) FindSingle(preparationTime m.PreparationTime) (m.PreparationTime, error) {

	result := r.db.First(&preparationTime)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.PreparationTime{}, errors.New("not found")
		} else {
			return m.PreparationTime{}, result.Error
		}
	}

	return preparationTime, nil
}

func (r *PreparationTimeRepository) Create(preparationTime m.PreparationTime) (m.PreparationTime, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Create(&preparationTime).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return m.PreparationTime{}, err
	}

	return preparationTime, nil
}

func (r *PreparationTimeRepository) Update(preparationTime m.PreparationTime) (m.PreparationTime, error) {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var err error

		if err = tx.Updates(&preparationTime).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return m.PreparationTime{}, err
	}

	return preparationTime, nil
}

func (r *PreparationTimeRepository) Delete(preparationTime m.PreparationTime) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&preparationTime).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
