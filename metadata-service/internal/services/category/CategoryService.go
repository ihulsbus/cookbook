package services

import (
	"errors"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
)

type CategoryRepository interface {
	FindAll(pagination m.PaginationRequest) ([]m.Category, int64, error)
	FindSingle(recipe m.Category) (m.Category, error)
	Create(recipe m.Category) (m.Category, error)
	Update(recipe m.Category) (m.Category, error)
	Delete(recipe m.Category) error
}
type CategoryService struct {
	repo CategoryRepository
}

// NewCategoryService creates a new CategoryService instance
func NewCategoryService(categoryRepo CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: categoryRepo,
	}
}

func (s CategoryService) FindAll(pagination m.PaginationRequest) (m.PaginatedResponse[m.CategoryDTO], error) {

	categories, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return m.PaginatedResponse[m.CategoryDTO]{}, errors.New("internal server error")
	}

	return m.PaginatedResponse[m.CategoryDTO]{
		Data:       m.Category{}.ConvertAllToDTO(categories),
		Pagination: m.NewPaginationMetadata(pagination.Page, pagination.Limit, total),
	}, nil
}

func (s CategoryService) FindSingle(categoryDTO m.CategoryDTO) (m.CategoryDTO, error) {

	category, err := s.repo.FindSingle(categoryDTO.ConvertFromDTO())
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.CategoryDTO{}, err
		default:
			return m.CategoryDTO{}, errors.New("internal server error")
		}
	}

	return category.ConvertToDTO(), nil
}

func (s CategoryService) Create(categoryDTO m.CategoryDTO) (m.CategoryDTO, error) {

	if categoryDTO.ID != uuid.Nil {
		return m.CategoryDTO{}, errors.New("existing id on new element is not allowed")
	}

	if categoryDTO.Name == "" {
		return m.CategoryDTO{}, errors.New("name is empty")
	}

	category, err := s.repo.Create(categoryDTO.ConvertFromDTO())
	if err != nil {
		return m.CategoryDTO{}, err
	}

	return category.ConvertToDTO(), nil
}

func (s CategoryService) Update(categoryDTO m.CategoryDTO) (m.CategoryDTO, error) {

	if categoryDTO.Name == "" {
		return m.CategoryDTO{}, errors.New("name is empty")
	}

	category, err := s.repo.Update(categoryDTO.ConvertFromDTO())
	if err != nil {
		return m.CategoryDTO{}, err
	}

	return category.ConvertToDTO(), nil
}

func (s CategoryService) Delete(categoryDTO m.CategoryDTO) error {

	err := s.repo.Delete(categoryDTO.ConvertFromDTO())
	if err != nil {
		return err
	}

	return nil
}
