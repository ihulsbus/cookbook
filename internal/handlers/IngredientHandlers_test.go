package handlers

import (
	"errors"

	m "github.com/ihulsbus/cookbook/internal/models"
)

type IngredientServiceMock struct {
}

func (s *IngredientServiceMock) FindAllIngredients() ([]m.Ingredient, error) {
	var ingredients []m.Ingredient

	ingredients = append(ingredients, m.Ingredient{}, m.Ingredient{})
	return ingredients, nil
}

func (s *IngredientServiceMock) FindSingleIngredient(ingredientID int) (m.Ingredient, error) {
	if ingredientID == 1 {
		return m.Ingredient{}, nil
	} else {
		return m.Ingredient{}, errors.New("error")
	}
}

func (s *IngredientServiceMock) CreateIngredient(ingredient m.Ingredient) (m.Ingredient, error) {
	if ingredient.IngredientName != "create" {
		return ingredient, errors.New("error")
	}
	return ingredient, nil
}

func (s *IngredientServiceMock) UpdateIngredient(ingredient m.Ingredient) (m.Ingredient, error) {
	if ingredient.IngredientName != "update" {
		return ingredient, errors.New("error")
	}
	return ingredient, nil
}

func (s *IngredientServiceMock) DeleteIngredient(ingredient m.Ingredient) error {
	if ingredient.ID != 1 {
		return errors.New("error")
	}
	return nil
}
