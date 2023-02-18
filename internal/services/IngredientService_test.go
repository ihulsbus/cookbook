package services

import (
	"errors"
	"testing"

	m "github.com/ihulsbus/cookbook/internal/models"
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
	case 1:
		var ingredients []m.Ingredient
		ingredients = append(ingredients, ingredient)

		return ingredients, nil
	default:
		return nil, errors.New("error")
	}
}

func (IngredientRepositoryMock) FindUnits() ([]m.Unit, error) {
	switch findAllUnit.ID {
	case 1:
		var units []m.Unit
		units = append(units, unit)

		return units, nil
	default:
		return nil, errors.New("error")
	}
}

func (IngredientRepositoryMock) FindSingle(ingredientID int) (m.Ingredient, error) {
	switch ingredientID {
	case 1:
		return ingredient, nil

	default:
		return ingredient, errors.New("error")
	}
}

func (IngredientRepositoryMock) Create(ingredient m.Ingredient) (m.Ingredient, error) {
	switch ingredient.ID {
	case 1:
		return ingredient, nil
	default:
		return ingredient, errors.New("error")
	}
}

func (IngredientRepositoryMock) Update(ingredient m.Ingredient) (m.Ingredient, error) {
	switch ingredient.ID {
	case 1:
		return ingredient, nil
	default:
		return ingredient, errors.New("error")
	}
}

func (IngredientRepositoryMock) Delete(ingredient m.Ingredient) error {
	switch ingredient.ID {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

func TestIngredientFindAll_OK(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})
	findAllIngredient.ID = 1

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, result[0].IngredientName, "ingredient")
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
	assert.Equal(t, result[0].ID, uint(1))
	assert.Equal(t, result[0].FullName, "Fluid Ounce")
	assert.Equal(t, result[0].ShortName, "fl oz")
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
	assert.IsType(t, result, m.Ingredient{})
	assert.Equal(t, result.IngredientName, "ingredient")
}

func TestIngredientFindSingle_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	result, err := s.FindSingle(2)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, m.Ingredient{})
}

func TestIngredientCreate_OK(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	createIngredient := ingredient
	createIngredient.ID = 1

	result, err := s.Create(createIngredient)

	assert.NoError(t, err)
	assert.IsType(t, result, m.Ingredient{})
	assert.Equal(t, result.IngredientName, "ingredient")
}

func TestIngredientCreate_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	createIngredient := ingredient
	createIngredient.ID = 2

	result, err := s.Create(createIngredient)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, m.Ingredient{})
}

func TestIngredientUpdate_Ok(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	updateIngredient := ingredient
	updateIngredient.ID = 1

	result, err := s.Update(updateIngredient)

	assert.NoError(t, err)
	assert.IsType(t, result, m.Ingredient{})
	assert.Equal(t, result.IngredientName, "ingredient")
}

func TestIngredientUpdate_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	updateIngredient := ingredient
	updateIngredient.ID = 2

	result, err := s.Update(updateIngredient)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, result, m.Ingredient{})
}

func TestIngredientDelete_Ok(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	deleteIngredient := ingredient
	deleteIngredient.ID = 1

	err := s.Delete(deleteIngredient)

	assert.NoError(t, err)
}

func TestIngredientDelete_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	deleteIngredient := ingredient
	deleteIngredient.ID = 2

	err := s.Delete(deleteIngredient)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
