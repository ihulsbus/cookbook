package services

import (
	"errors"
	"testing"

	m "ingredient-service/internal/models"

	"github.com/stretchr/testify/assert"
)

var (
	findAllIngredient              = ingredient
	ingredient        m.Ingredient = m.Ingredient{IngredientName: "ingredient"}

	findAllUnit        = unit
	unit        m.Unit = m.Unit{ID: 1, FullName: "Fluid Ounce", ShortName: "fl oz"}
)

type IngredientRepositoryMock struct{}

func (IngredientRepositoryMock) FindAll() ([]m.Ingredient, error) {
	switch findAllIngredient.ID {
	case 1: // OK
		var ingredients []m.Ingredient
		ingredients = append(ingredients, ingredient)
		return ingredients, nil
	default: // ERR
		return nil, errors.New("error")
	}
}

func (IngredientRepositoryMock) FindUnits() ([]m.Unit, error) {
	switch findAllUnit.ID {
	case 1: // OK
		var units []m.Unit
		units = append(units, unit)
		return units, nil
	default: // ERR
		return nil, errors.New("error")
	}
}

func (IngredientRepositoryMock) FindSingle(ingredientID uint) (m.Ingredient, error) {
	switch ingredientID {
	case 0:
		ing := ingredient
		ing.ID = 0
		return ing, nil
	case 2: // NOT FOUND
		return m.Ingredient{}, errors.New("error")
	case 3:
		return m.Ingredient{}, errors.New("error")
	default: // OK
		ing := ingredient
		ing.ID = ingredientID
		return ing, nil
	}
}

func (IngredientRepositoryMock) Create(i m.Ingredient) (m.Ingredient, error) {
	switch i.ID {
	case 2: // OK
		ing := ingredient
		ing.ID = 2
		return ing, nil
	default: // ERR
		return ingredient, errors.New("error")
	}
}

func (IngredientRepositoryMock) Update(i m.Ingredient) (m.Ingredient, error) {
	switch i.ID {
	case 1: // OK
		ing := ingredient
		ing.ID = 1
		return ing, nil
	default: // ERR
		return ingredient, errors.New("error")
	}
}

func (IngredientRepositoryMock) Delete(i m.Ingredient) error {
	switch i.ID {
	case 1: // OK
		return nil
	default: // ERR
		return errors.New("error")
	}
}

func TestIngredientFindAll_OK(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})
	findAllIngredient.ID = 1

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "ingredient", result[0].IngredientName)
}

func TestIngredientFindAll_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	findAllIngredient.ID = 2

	result, err := s.FindAll()

	findAllIngredient.ID = 1

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "error")
}

func TestIngredientFindUnits_OK(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	findAllUnit.ID = 1

	result, err := s.FindUnits()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, uint(1), result[0].ID)
	assert.Equal(t, "Fluid Ounce", result[0].FullName)
	assert.Equal(t, "fl oz", result[0].ShortName)
}

func TestIngredientFindUnits_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	findAllUnit.ID = 2

	result, err := s.FindUnits()

	findAllUnit.ID = 1

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "error")
}

func TestIngredientFindSingle_OK(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	result, err := s.FindSingle(1)

	assert.NoError(t, err)
	assert.IsType(t, m.Ingredient{}, result)
	assert.Equal(t, "ingredient", result.IngredientName)
}

func TestIngredientFindSingle_FindErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	result, err := s.FindSingle(2)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, m.Ingredient{})
}

func TestIngredientFindSingle_NotFoundErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	result, err := s.FindSingle(0)

	assert.Error(t, err)
	assert.EqualError(t, err, "ingredient not found")
	assert.IsType(t, m.Ingredient{}, result)
}

func TestIngredientCreate_OK(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	createIngredient := ingredient
	createIngredient.ID = 2

	result, err := s.Create(createIngredient)

	assert.NoError(t, err)
	assert.IsType(t, m.Ingredient{}, result)
	assert.Equal(t, "ingredient", result.IngredientName)
}

func TestIngredientCreate_ExistsErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	createIngredient := ingredient
	createIngredient.ID = 1

	result, err := s.Create(createIngredient)

	assert.Error(t, err)
	assert.EqualError(t, err, "ingredient already exists")
	assert.IsType(t, m.Ingredient{}, result)
}

func TestIngredientCreate_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	createIngredient := ingredient
	createIngredient.ID = 3

	result, err := s.Create(createIngredient)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, m.Ingredient{}, result)
}

func TestIngredientUpdate_Ok(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	updateIngredient := ingredient
	updateIngredient.ID = 1

	result, err := s.Update(updateIngredient, 1)

	assert.NoError(t, err)
	assert.IsType(t, result, m.Ingredient{})
	assert.Equal(t, result.IngredientName, "ingredient")
}

func TestIngredientUpdate_NotFoundErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	updateIngredient := ingredient
	updateIngredient.ID = 2

	result, err := s.Update(updateIngredient, uint(2))

	assert.Error(t, err)
	assert.EqualError(t, err, "ingredient does not exist. nothing to update")
	assert.IsType(t, result, m.Ingredient{})
}

func TestIngredientUpdate_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	updateIngredient := ingredient
	updateIngredient.ID = 4

	result, err := s.Update(updateIngredient, uint(4))

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, m.Ingredient{})
}

func TestIngredientDelete_Ok(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	deleteIngredient := ingredient
	deleteIngredient.ID = 1

	err := s.Delete(uint(1))

	assert.NoError(t, err)
}

func TestIngredientDelete_NotFoundErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	deleteIngredient := ingredient
	deleteIngredient.ID = 2

	err := s.Delete(uint(2))

	assert.Error(t, err)
	assert.EqualError(t, err, "ingredient does not exist. nothing to delete")
}

func TestIngredientDelete_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	deleteIngredient := ingredient
	deleteIngredient.ID = 4

	err := s.Delete(uint(4))

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
