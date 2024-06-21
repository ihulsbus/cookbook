package services

import (
	"errors"
	"testing"

	m "ingredient-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	findAllIngredient m.Ingredient = m.Ingredient{
		ID:   uuid.New(),
		Name: "ingredient",
	}
	ingredient m.Ingredient = m.Ingredient{
		ID:   uuid.New(),
		Name: "ingredient",
	}

	findAllUnit m.Unit = m.Unit{
		ID:        uuid.New(),
		FullName:  "unit",
		ShortName: "u",
	}
	unit m.Unit = m.Unit{
		ID:        uuid.New(),
		FullName:  "unit",
		ShortName: "u",
	}
)

type IngredientRepositoryMock struct{}

func (IngredientRepositoryMock) FindAll() ([]m.Ingredient, error) {
	switch findAllIngredient.Name {
	case "findall": // OK
		var ingredients []m.Ingredient
		ingredients = append(ingredients, ingredient)
		return ingredients, nil
	case "notfound":
		return nil, errors.New("not found")
	default: // ERR
		return nil, errors.New("error")
	}
}

func (IngredientRepositoryMock) FindSingle(ingredientInput m.Ingredient) (m.Ingredient, error) {
	switch ingredientInput.Name {
	case "find":
		return ingredient, nil
	case "update":
		return ingredient, nil
	case "updateerror":
		return ingredient, nil
	case "delete":
		return ingredient, nil
	case "deleteerror":
		return ingredient, nil
	case "":
		return ingredient, nil
	case "notfound":
		return m.Ingredient{}, errors.New("not found")
	default:
		return m.Ingredient{}, errors.New("error")
	}
}

func (IngredientRepositoryMock) Create(ingredientInput m.Ingredient) (m.Ingredient, error) {
	switch ingredientInput.Name {
	case "create":
		return ingredient, nil
	default:
		return m.Ingredient{}, errors.New("error")
	}
}

func (IngredientRepositoryMock) Update(ingredientInput m.Ingredient) (m.Ingredient, error) {
	switch ingredientInput.Name {
	case "update":
		return ingredient, nil
	case "ingredient":
		return ingredient, nil
	default:
		return m.Ingredient{}, errors.New("error")
	}
}

func (IngredientRepositoryMock) Delete(ingredientInput m.Ingredient) error {
	switch ingredientInput.Name {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

func TestIngredientFindAll_OK(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	findAllIngredient.Name = "findall"

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "ingredient", result[0].Name)
}

func TestIngredientFindAll_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	findAllIngredient.Name = "error"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "internal server error")
}

func TestRecipeFindAll_NotFound(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	findAllIngredient.Name = "notfound"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
}

func TestIngredientFindSingle_OK(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		ID:   ingredient.ID,
		Name: "find",
	}
	result, err := s.FindSingle(ingredientDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.IngredientDTO{}, result)
	assert.Equal(t, "ingredient", result.Name)
	assert.Equal(t, ingredient.ID, result.ID)
}

func TestIngredientFindSingle_FindErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		ID:   ingredient.ID,
		Name: "error",
	}
	result, err := s.FindSingle(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, m.IngredientDTO{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestIngredientFindSingle_NotFoundErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		ID:   ingredient.ID,
		Name: "notfound",
	}
	result, err := s.FindSingle(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, m.IngredientDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestIngredientCreate_OK(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		Name: "create",
	}
	result, err := s.Create(ingredientDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.IngredientDTO{}, result)
	assert.Equal(t, "ingredient", result.Name)
}

func TestIngredientCreate_IDErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		ID:   ingredient.ID,
		Name: "create",
	}
	result, err := s.Create(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, m.IngredientDTO{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestIngredientCreate_ExistsErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		Name: "find",
	}
	result, err := s.Create(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, m.IngredientDTO{}, result)
	assert.EqualError(t, err, "ingredient already exists")
}

func TestIngredientCreate_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		Name: "error",
	}
	result, err := s.Create(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, m.IngredientDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestIngredientCreate_NoName(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		Name: "",
	}
	result, err := s.Create(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, m.IngredientDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestIngredientUpdate_Ok(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		ID:   ingredient.ID,
		Name: "update",
	}
	result, err := s.Update(ingredientDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.IngredientDTO{}, result)
	assert.Equal(t, result.Name, "ingredient")
}

func TestIngredientUpdate_NotFoundErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		ID:   ingredient.ID,
		Name: "notfound",
	}
	result, err := s.Update(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, m.IngredientDTO{}, result)
	assert.EqualError(t, err, "ingredient does not exist. nothing to update")
}

func TestIngredientUpdate_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		ID:   ingredient.ID,
		Name: "updateerror",
	}
	result, err := s.Update(ingredientDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, m.IngredientDTO{}, result)
}

func TestIngredientDelete_Ok(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		ID:   ingredient.ID,
		Name: "delete",
	}
	err := s.Delete(ingredientDTO)

	assert.NoError(t, err)
}

func TestIngredientDelete_NotFoundErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		ID:   ingredient.ID,
		Name: "notfound",
	}
	err := s.Delete(ingredientDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "ingredient does not exist. nothing to delete")
}

func TestIngredientDelete_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := m.IngredientDTO{
		ID:   ingredient.ID,
		Name: "deleteerror",
	}
	err := s.Delete(ingredientDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
