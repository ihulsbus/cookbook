package repositories

import (
	"errors"

	m "ingredient-service/internal/models"

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

func (r UnitRepository) FindAll() ([]m.Unit, error) {
	var units []m.Unit

	if err := r.db.Find(&units).Error; err != nil {
		return nil, err
	}

	if len(units) <= 0 {
		return nil, errors.New("not found")
	}

	return units, nil
}

func (r UnitRepository) FindSingle(unit m.Unit) (m.Unit, error) {

	result := r.db.First(&unit)
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
