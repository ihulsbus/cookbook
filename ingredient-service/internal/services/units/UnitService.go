package services

import (
	"errors"

	m "ingredient-service/internal/models"

	"github.com/google/uuid"
)

type UnitRepository interface {
	FindAll() ([]m.Unit, error)
	FindSingle(unit m.Unit) (m.Unit, error)
	Create(unit m.Unit) (m.Unit, error)
	Update(unit m.Unit) (m.Unit, error)
	Delete(unit m.Unit) error
}
type UnitService struct {
	repo UnitRepository
}

// NewUnitService creates a new UnitService instance
func NewUnitService(unitRepo UnitRepository) *UnitService {
	return &UnitService{
		repo: unitRepo,
	}
}

func (s UnitService) FindAll() ([]m.UnitDTO, error) {
	var units []m.Unit

	units, err := s.repo.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	return m.Unit{}.ConvertAllToDTO(units), nil
}

func (s UnitService) FindSingle(unitDTO m.UnitDTO) (m.UnitDTO, error) {

	unit, err := s.repo.FindSingle(unitDTO.ConvertFromDTO())
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.UnitDTO{}, err
		default:
			return m.UnitDTO{}, errors.New("internal server error")
		}
	}

	return unit.ConvertToDTO(), nil
}

func (s UnitService) Create(unitDTO m.UnitDTO) (m.UnitDTO, error) {
	var unit m.Unit

	if unitDTO.ID != uuid.Nil {
		return m.UnitDTO{}, errors.New("existing id on new element is not allowed")
	}

	if unitDTO.FullName == "" {
		return m.UnitDTO{}, errors.New("name is empty")
	}

	found, err := s.FindSingle(unitDTO)
	if err == nil || found.ID != uuid.Nil {
		return m.UnitDTO{}, errors.New("unit already exists")
	}

	unit, err = s.repo.Create(unitDTO.ConvertFromDTO())
	if err != nil {
		return m.UnitDTO{}, err
	}

	return unit.ConvertToDTO(), nil
}

func (s UnitService) Update(unitDTO m.UnitDTO) (m.UnitDTO, error) {
	var unit m.Unit

	_, err := s.FindSingle(unitDTO)
	if err != nil {
		return m.UnitDTO{}, errors.New("unit does not exist. nothing to update")
	}

	unit, err = s.repo.Update(unitDTO.ConvertFromDTO())
	if err != nil {
		return m.UnitDTO{}, err
	}

	return unit.ConvertToDTO(), nil
}

func (s UnitService) Delete(unitDTO m.UnitDTO) error {

	_, err := s.FindSingle(unitDTO)
	if err != nil {
		return errors.New("unit does not exist. nothing to delete")
	}

	// TODO: check if there are recipies using the unit. If so, an error should be returned and the unit should not be deleted.
	err = s.repo.Delete(unitDTO.ConvertFromDTO())
	if err != nil {
		return err
	}

	return nil
}
