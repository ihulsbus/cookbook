package services

import (
	"errors"
	"testing"

	m "github.com/ihulsbus/cookbook/internal/models"
	"github.com/stretchr/testify/assert"
)

var (
	findAllCategory m.Category = m.Category{
		CategoryName: "category",
	}
	category m.Category = m.Category{
		CategoryName: "category",
	}
)

type CategoryRepositoryMock struct{}

func (*CategoryRepositoryMock) FindAll() ([]m.Category, error) {
	switch findAllCategory.ID {
	case 1:
		var categorys []m.Category
		categorys = append(categorys, findAllCategory)
		return categorys, nil
	default:
		return nil, errors.New("error")
	}
}

func (*CategoryRepositoryMock) FindSingle(categoryID uint) (m.Category, error) {
	switch categoryID {
	case 1:
		return category, nil
	default:
		return m.Category{}, errors.New("error")
	}
}

func (*CategoryRepositoryMock) Create(createCategory m.Category) (m.Category, error) {
	switch createCategory.CategoryName {
	case "create":
		createCategory.ID = 1
		return createCategory, nil
	default:
		return m.Category{}, errors.New("error")
	}
}

func (*CategoryRepositoryMock) Update(category m.Category) (m.Category, error) {
	switch category.ID {
	case 1:
		return category, nil
	default:
		return m.Category{}, errors.New("error")
	}
}

func (*CategoryRepositoryMock) Delete(category m.Category) error {
	switch category.ID {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

// ======================================================================

func TestCategoryFindAll_OK(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})
	findAllCategory.ID = 1

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.IsType(t, []m.Category{}, result)
	assert.Len(t, result, 1)
}

func TestCategoryFindAll_err(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})
	findAllCategory.ID = 2

	result, err := s.FindAll()

	findAllCategory.ID = 1

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "error")
}

func TestCategoryFindSingle_OK(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	category.ID = 1

	result, err := s.FindSingle(uint(1))

	assert.NoError(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.Equal(t, "category", result.CategoryName)
	assert.Equal(t, uint(1), result.ID)
}

func TestCategoryFindSingle_Err(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	result, err := s.FindSingle(uint(2))

	assert.Error(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.EqualError(t, err, "error")
}

func TestCategoryCreate_OK(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	createCategory := category
	createCategory.ID = 0
	createCategory.CategoryName = "create"

	result, err := s.Create(createCategory)

	assert.NoError(t, err)
	assert.IsType(t, m.Category{}, result)
}

func TestCategoryCreate_IDErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	createCategory := category
	createCategory.ID = 1
	createCategory.CategoryName = "create"

	result, err := s.Create(createCategory)

	assert.Error(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestCategoryCreate_NoNameErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	createCategory := category
	createCategory.ID = 0
	createCategory.CategoryName = ""

	result, err := s.Create(createCategory)

	assert.Error(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.EqualError(t, err, "categoryname is empty")
}

func TestCategoryCreate_CreateErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	createCategory := category
	createCategory.ID = 0
	createCategory.CategoryName = "fails"

	result, err := s.Create(createCategory)

	assert.Error(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.EqualError(t, err, "error")
}

func TestCategoryUpdate_OK(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	updateCategory := category
	updateCategory.ID = 1
	updateCategory.CategoryName = "update"

	result, err := s.Update(updateCategory, uint(1))

	assert.NoError(t, err)
	assert.IsType(t, m.Category{}, result)
}

func TestCategoryUpdate_IDErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	updateCategory := category
	updateCategory.ID = 0
	updateCategory.CategoryName = "update"

	result, err := s.Update(updateCategory, uint(0))

	assert.Error(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.EqualError(t, err, "missing id of element to update")
}

func TestCategoryUpdate_IDNotEqual(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	updateCategory := category
	updateCategory.ID = 2
	updateCategory.CategoryName = "update"

	result, err := s.Update(updateCategory, uint(1))

	assert.NoError(t, err)
	assert.IsType(t, m.Category{}, result)
}

func TestCategoryUpdate_NoNameErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	updateCategory := category
	updateCategory.ID = 1
	updateCategory.CategoryName = ""

	result, err := s.Update(updateCategory, uint(1))

	assert.Error(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.EqualError(t, err, "categoryname is empty")
}

func TestCategoryUpdate_UpdateErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	updateCategory := category
	updateCategory.ID = 2
	updateCategory.CategoryName = "fails"

	result, err := s.Update(updateCategory, uint(2))

	assert.Error(t, err)
	assert.IsType(t, m.Category{}, result)
	assert.EqualError(t, err, "error")
}

func TestCategoryDelete_OK(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	err := s.Delete(uint(1))

	assert.NoError(t, err)
}

func TestCategoryDelete_IDErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	err := s.Delete(uint(0))

	assert.Error(t, err)
	assert.EqualError(t, err, "missing id of element to delete")
}

func TestCategoryDelete_DeleteErr(t *testing.T) {
	s := NewCategoryService(&CategoryRepositoryMock{})

	err := s.Delete(uint(2))

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
