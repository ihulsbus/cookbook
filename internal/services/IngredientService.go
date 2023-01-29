package services

import (
	m "github.com/ihulsbus/cookbook/internal/models"
	r "github.com/ihulsbus/cookbook/internal/repositories"
	log "github.com/sirupsen/logrus"
)

type IngredientService struct {
	repo   *r.IngredientRepository
	logger *log.Logger
}

// NewRecipeService creates a new RecipeService instance
func NewIngredientService(ingredientRepo *r.IngredientRepository, logger *log.Logger) *IngredientService {
	return &IngredientService{
		repo:   ingredientRepo,
		logger: logger,
	}
}

func (s IngredientService) FindAllIngredients() ([]m.Ingredient, error) {
	var ingredients []m.Ingredient

	ingredients, err := s.repo.IngredientFindAll()
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (s IngredientService) FindSingleIngredient(ingredientID int) (m.Ingredient, error) {
	var ingredient m.Ingredient

	ingredient, err := s.repo.IngredientFindSingle(ingredientID)
	if err != nil {
		return ingredient, err
	}

	return ingredient, nil
}

func (s IngredientService) CreateIngredient(ingredient m.Ingredient) (m.Ingredient, error) {
	var response m.Ingredient

	response, err := s.repo.CreateIngredient(ingredient)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s IngredientService) UpdateIngredient(ingredient m.Ingredient) (m.Ingredient, error) {
	var response m.Ingredient

	response, err := s.repo.UpdateIngredient(ingredient)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s IngredientService) DeleteIngredient(ingredient m.Ingredient) error {

	// TODO: check if there are recipies using the ingredient. If so, an error should be returned and the ingredient should not be deleted.
	err := s.repo.DeleteIngredient(ingredient)
	if err != nil {
		return err
	}

	return nil
}
