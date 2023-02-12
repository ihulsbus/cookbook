package services

import (
	m "github.com/ihulsbus/cookbook/internal/models"
)

type IngredientRepository interface {
	FindAll() ([]m.Ingredient, error)
	FindUnits() ([]m.Unit, error)
	FindSingle(ingredientID int) (m.Ingredient, error)
	Create(ingredient m.Ingredient) (m.Ingredient, error)
	Update(ingredient m.Ingredient) (m.Ingredient, error)
	Delete(ingredient m.Ingredient) error
}
type IngredientService struct {
	repo IngredientRepository
}

// NewRecipeService creates a new RecipeService instance
func NewIngredientService(ingredientRepo IngredientRepository) *IngredientService {
	return &IngredientService{
		repo: ingredientRepo,
	}
}

func (s IngredientService) FindAll() ([]m.Ingredient, error) {
	var ingredients []m.Ingredient

	ingredients, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (s IngredientService) FindUnits() ([]m.Unit, error) {
	var units []m.Unit

	units, err := s.repo.FindUnits()
	if err != nil {
		return nil, err
	}

	return units, nil
}

func (s IngredientService) FindSingle(ingredientID int) (m.Ingredient, error) {
	var ingredient m.Ingredient

	ingredient, err := s.repo.FindSingle(ingredientID)
	if err != nil {
		return ingredient, err
	}

	return ingredient, nil
}

func (s IngredientService) Create(ingredient m.Ingredient) (m.Ingredient, error) {
	var response m.Ingredient

	response, err := s.repo.Create(ingredient)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s IngredientService) Update(ingredient m.Ingredient) (m.Ingredient, error) {
	var response m.Ingredient

	response, err := s.repo.Update(ingredient)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s IngredientService) Delete(ingredient m.Ingredient) error {

	// TODO: check if there are recipies using the ingredient. If so, an error should be returned and the ingredient should not be deleted.
	err := s.repo.Delete(ingredient)
	if err != nil {
		return err
	}

	return nil
}
