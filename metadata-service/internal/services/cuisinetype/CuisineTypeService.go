package services

import (
	"errors"

	m "metadata-service/internal/models"

	"github.com/google/uuid"
)

type CuisineTypeRepository interface {
	FindAll() ([]m.CuisineType, error)
	FindSingle(cuisineType m.CuisineType) (m.CuisineType, error)
	Create(cuisineType m.CuisineType) (m.CuisineType, error)
	Update(cuisineType m.CuisineType) (m.CuisineType, error)
	Delete(cuisineType m.CuisineType) error
}
type CuisineTypeService struct {
	repo CuisineTypeRepository
}

// NewCuisineTypeService creates a new CuisineTypeService instance
func NewCuisineTypeService(cuisineTypeRepo CuisineTypeRepository) *CuisineTypeService {
	return &CuisineTypeService{
		repo: cuisineTypeRepo,
	}
}

func (s CuisineTypeService) FindAll() ([]m.CuisineTypeDTO, error) {

	cuisineTypes, err := s.repo.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	result := m.CuisineType{}.ConvertAllToDTO(cuisineTypes)
	return result, nil
}

func (s CuisineTypeService) FindSingle(cuisineTypeDTO m.CuisineTypeDTO) (m.CuisineTypeDTO, error) {

	cuisineType, err := s.repo.FindSingle(cuisineTypeDTO.ConvertFromDTO())
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.CuisineTypeDTO{}, err
		default:
			return m.CuisineTypeDTO{}, errors.New("internal server error")
		}
	}

	return cuisineType.ConvertToDTO(), nil
}

func (s CuisineTypeService) Create(cuisineTypeDTO m.CuisineTypeDTO) (m.CuisineTypeDTO, error) {

	if cuisineTypeDTO.ID != uuid.Nil {
		return m.CuisineTypeDTO{}, errors.New("existing id on new element is not allowed")
	}

	if cuisineTypeDTO.Name == "" {
		return m.CuisineTypeDTO{}, errors.New("name is empty")
	}

	created, err := s.repo.Create(cuisineTypeDTO.ConvertFromDTO())
	if err != nil {
		return m.CuisineTypeDTO{}, err
	}

	return created.ConvertToDTO(), nil
}

func (s CuisineTypeService) Update(cuisineTypeDTO m.CuisineTypeDTO) (m.CuisineTypeDTO, error) {

	if cuisineTypeDTO.Name == "" {
		return m.CuisineTypeDTO{}, errors.New("name is empty")
	}

	updatedCuisineType, err := s.repo.Update(cuisineTypeDTO.ConvertFromDTO())
	if err != nil {
		return m.CuisineTypeDTO{}, err
	}

	return updatedCuisineType.ConvertToDTO(), nil
}

func (s CuisineTypeService) Delete(cuisineTypeDTO m.CuisineTypeDTO) error {

	err := s.repo.Delete(cuisineTypeDTO.ConvertFromDTO())
	if err != nil {
		return err
	}

	return nil
}
