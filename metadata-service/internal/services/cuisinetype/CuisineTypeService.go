package services

import (
	"errors"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
)

type CuisineTypeRepository interface {
	FindAll(pagination m.PaginationRequest) ([]m.CuisineType, int64, error)
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

func (s CuisineTypeService) FindAll(pagination m.PaginationRequest) (m.PaginatedResponse[m.CuisineTypeDTO], error) {

	cuisineTypes, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return m.PaginatedResponse[m.CuisineTypeDTO]{}, errors.New("internal server error")
	}

	return m.PaginatedResponse[m.CuisineTypeDTO]{
		Data:       m.CuisineType{}.ConvertAllToDTO(cuisineTypes),
		Pagination: m.NewPaginationMetadata(pagination.Page, pagination.Limit, total),
	}, nil
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
