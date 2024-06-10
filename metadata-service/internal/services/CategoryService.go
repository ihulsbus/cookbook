package services

import (
	"errors"

	m "metadata-service/internal/models"
)

type CategoryRepository interface {
	FindAll() ([]m.Category, error)
	FindSingle(recipeID uint) (m.Category, error)
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

func (s CategoryService) FindAll() ([]m.Category, error) {

	categorys, err := s.repo.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	return categorys, nil
}

func (s CategoryService) FindSingle(categoryID uint) (m.Category, error) {

	category, err := s.repo.FindSingle(categoryID)
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.Category{}, err
		default:
			return m.Category{}, errors.New("internal server error")
		}
	}

	return category, nil
}

func (s CategoryService) Create(category m.Category) (m.Category, error) {

	if category.ID != 0 {
		return m.Category{}, errors.New("existing id on new element is not allowed")
	}

	if category.CategoryName == "" {
		return m.Category{}, errors.New("categoryname is empty")
	}

	created, err := s.repo.Create(category)
	if err != nil {
		return m.Category{}, err
	}

	return created, nil
}

func (s CategoryService) Update(category m.Category, categoryID uint) (m.Category, error) {

	if categoryID == 0 {
		return m.Category{}, errors.New("missing id of element to update")
	}

	if category.ID != categoryID {
		category.ID = categoryID
	}

	if category.CategoryName == "" {
		return m.Category{}, errors.New("categoryname is empty")
	}

	updatedCategory, err := s.repo.Update(category)
	if err != nil {
		return m.Category{}, err
	}

	return updatedCategory, nil
}

func (s CategoryService) Delete(categoryID uint) error {
	var category m.Category

	if categoryID == 0 {
		return errors.New("missing id of element to delete")
	}

	category.ID = categoryID

	err := s.repo.Delete(category)
	if err != nil {
		return err
	}

	return nil
}
