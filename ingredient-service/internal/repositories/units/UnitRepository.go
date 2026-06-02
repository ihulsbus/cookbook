package repositories

import (
	"errors"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"gorm.io/gorm"
)

type UnitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) *UnitRepository {
	return &UnitRepository{
		db: db,
	}
}

func (r UnitRepository) FindAll(pagination m.PaginationRequest) ([]m.Unit, int64, error) {
	var units []m.Unit
	var total int64

	if err := r.db.Model(&m.Unit{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(pagination.Limit).Offset(pagination.Offset()).Find(&units).Error; err != nil {
		return nil, 0, err
	}

	return units, total, nil
}

func (r UnitRepository) FindSingle(unit m.Unit) (m.Unit, error) {
	var result *gorm.DB

	if unit.ID == uuid.Nil {
		result = r.db.First(&unit, "full_name = ?", unit.FullName)
	} else {
		result = r.db.First(&unit)
	}
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return m.Unit{}, errors.New("not found")
		} else {
			return m.Unit{}, result.Error
		}
	}

	return unit, nil
}

func (r UnitRepository) Create(unit m.Unit) (m.Unit, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&unit).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return unit, err
	}

	return unit, nil
}

func (r UnitRepository) Update(unit m.Unit) (m.Unit, error) {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Updates(&unit).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return unit, err
	}

	return unit, nil
}

func (r UnitRepository) Delete(unit m.Unit) error {

	if err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(&unit).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
