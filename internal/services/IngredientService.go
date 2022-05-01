package services

import (
	m "github.com/ihulsbus/cookbook/internal/models"
	r "github.com/ihulsbus/cookbook/internal/repositories"
)

type FindAllIngredients func(s *IngredientService) ([]m.Ingredient, error)
type FindSingleIngredient func(s *IngredientService, ingredientID int) ([]m.Ingredient, error)
type FindAllSections func(s *IngredientService) error
type FindSingleSection func(s *IngredientService, sectionID int) error
type FindRecipeIngredients func(s *IngredientService, recipeID int) ([]m.Recipe_Ingredient, error)
type CreateIngredient func(s *IngredientService, ingredient m.Ingredient) (m.Ingredient, error)
type UpdateIngredient func(s *IngredientService, ingredient m.Ingredient) (m.Ingredient, error)
type DeleteIngredient func(s *IngredientService, ingredient m.Ingredient) error

type IngredientService struct {
	repo *r.IngredientRepository

	FindAllIngredients   FindAllIngredients
	FindSingleIngredient FindSingleIngredient
	CreateIngredient     CreateIngredient
	UpdateIngredient     UpdateIngredient
	DeleteIngredient     DeleteIngredient

	FindAllSections   FindAllSections
	FindSingleSection FindSingleSection

	FindRecipeIngredients FindRecipeIngredients
}

// NewRecipeService creates a new RecipeService instance
func NewIngredientService(ingredientRepo *r.IngredientRepository) *IngredientService {
	return &IngredientService{
		repo: ingredientRepo,

		FindAllIngredients:   findAllIngredients,
		FindSingleIngredient: findSingleIngredient,
		CreateIngredient:     createIngredient,
		UpdateIngredient:     updateIngredient,
		DeleteIngredient:     deleteIngredient,
	}
}

func findAllIngredients(s *IngredientService) ([]m.Ingredient, error) {
	var ingredients []m.Ingredient

	ingredients, err := s.repo.IngredientFindAll(s.repo)
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

func findSingleIngredient(s *IngredientService, ingredientID int) ([]m.Ingredient, error) {
	var ingredients []m.Ingredient

	ingredients, err := s.repo.IngredientFindSingle(s.repo, ingredientID)
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

func createIngredient(s *IngredientService, ingredient m.Ingredient) (m.Ingredient, error) {
	var response m.Ingredient

	response, err := s.repo.CreateIngredient(s.repo, ingredient)
	if err != nil {
		return response, err
	}

	return response, nil
}

func updateIngredient(s *IngredientService, ingredient m.Ingredient) (m.Ingredient, error) {
	var response m.Ingredient

	response, err := s.repo.UpdateIngredient(s.repo, ingredient)
	if err != nil {
		return response, err
	}

	return response, nil
}

func deleteIngredient(s *IngredientService, ingredient m.Ingredient) error {

	// TODO: check if there are recipies using the ingredient. If so, an error should be returned and the ingredient should not be deleted.
	err := s.repo.DeleteIngredient(s.repo, ingredient)
	if err != nil {
		return err
	}

	return nil
}
