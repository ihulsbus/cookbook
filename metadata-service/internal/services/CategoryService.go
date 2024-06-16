package services

import (
	"errors"

	m "metadata-service/internal/models"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	FindAll() ([]m.Category, error)
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

func (s CategoryService) FindAll() ([]m.CategoryDTO, error) {

	categories, err := s.repo.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	categoryDTOs := m.Category{}.ConvertAllToDTO(categories)
	return categoryDTOs, nil
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
