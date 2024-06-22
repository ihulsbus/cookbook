package services

import (
	"errors"
	"testing"

	m "metadata-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	findAllCategory m.Category = m.Category{
		ID:   uuid.New(),
		Name: "category",
	}
	category m.Category = m.Category{
		ID:   uuid.New(),
		Name: "category",
	}
)

type CategoryRepositoryMock struct{}

func (*CategoryRepositoryMock) FindAll() ([]m.Category, error) {
	switch findAllCategory.Name {
	case "findall":
		var categories []m.Category
		categories = append(categories, findAllCategory)
		return categories, nil
	case "not found":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (*CategoryRepositoryMock) FindSingle(category m.Category) (m.Category, error) {
	switch category.Name {
	case "find":
		return category, nil
	case "not found":
		return m.Category{}, errors.New("not found")
	default:
		return m.Category{}, errors.New("error")
	}
}

func (*CategoryRepositoryMock) Create(category m.Category) (m.Category, error) {
	categoryC := category
	switch category.Name {
	case "create":
		categoryC.Name = "create"
		return categoryC, nil
	default:
		return m.Category{}, errors.New("error")
	}
}

func (*CategoryRepositoryMock) Update(category m.Category) (m.Category, error) {
	switch category.Name {
	case "update":
		return category, nil
	default:
		return m.Category{}, errors.New("error")
	}
}

func (*CategoryRepositoryMock) Delete(category m.Category) error {
	switch category.Name {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

// ======================================================================

func TestCategoryFindAll_OK(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})
	findAllCategory.Name = "findall"

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.IsType(t, []m.CategoryDTO{}, result)
	assert.Len(t, result, 1)
}

func TestCategoryFindAll_err(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})
	findAllCategory.Name = "fail"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "internal server error")
}

func TestCategoryFindAll_NotFound(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})
	findAllCategory.Name = "not found"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
}

func TestCategoryFindSingle_OK(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	categoryDTO := m.CategoryDTO{
		ID:   category.ID,
		Name: "find",
	}
	result, err := s.FindSingle(categoryDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.CategoryDTO{}, result)
	assert.Equal(t, "find", result.Name)
	assert.Equal(t, result.ID, category.ID)
}

func TestCategoryFindSingle_Err(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	categoryDTO := m.CategoryDTO{
		ID:   category.ID,
		Name: category.Name,
	}
	result, err := s.FindSingle(categoryDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CategoryDTO{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestCategoryFindSingle_NotFound(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	categoryDTO := m.CategoryDTO{
		ID:   category.ID,
		Name: "not found",
	}
	result, err := s.FindSingle(categoryDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CategoryDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestCategoryCreate_OK(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	categoryDTO := m.CategoryDTO{
		Name: "create",
	}
	result, err := s.Create(categoryDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.CategoryDTO{}, result)
}

func TestCategoryCreate_IDErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	categoryDTO := m.CategoryDTO{
		ID:   category.ID,
		Name: category.Name,
	}
	result, err := s.Create(categoryDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CategoryDTO{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestCategoryCreate_NoNameErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	categoryDTO := m.CategoryDTO{
		Name: "",
	}
	result, err := s.Create(categoryDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CategoryDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestCategoryCreate_CreateErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	categoryDTO := m.CategoryDTO{
		Name: category.Name,
	}
	result, err := s.Create(categoryDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CategoryDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestCategoryUpdate_OK(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	categoryDTO := m.CategoryDTO{
		ID:   category.ID,
		Name: "update",
	}
	result, err := s.Update(categoryDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.CategoryDTO{}, result)
}

func TestCategoryUpdate_NoNameErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	categoryDTO := m.CategoryDTO{
		ID: category.ID,
	}
	result, err := s.Update(categoryDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CategoryDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestCategoryUpdate_UpdateErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	categoryDTO := m.CategoryDTO{
		ID:   category.ID,
		Name: category.Name,
	}
	result, err := s.Update(categoryDTO)

	assert.Error(t, err)
	assert.IsType(t, m.CategoryDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestCategoryDelete_OK(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	categoryDTO := m.CategoryDTO{
		ID:   category.ID,
		Name: "delete",
	}
	err := s.Delete(categoryDTO)

	assert.NoError(t, err)
}

func TestCategoryDelete_DeleteErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	categoryDTO := m.CategoryDTO{
		ID:   category.ID,
		Name: category.Name,
	}
	err := s.Delete(categoryDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
