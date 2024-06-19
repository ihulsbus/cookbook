package services

import (
	"errors"

	m "recipe-service/internal/models"

	"github.com/google/uuid"
)

type RecipeRepository interface {
	FindAll() ([]m.Recipe, error)
	FindSingle(recipe m.Recipe) (m.Recipe, error)
	Create(recipe m.Recipe) (m.Recipe, error)
	Update(recipe m.Recipe) (m.Recipe, error)
	Delete(recipe m.Recipe) error
}

type RecipeService struct {
	repo RecipeRepository
}

// NewRecipeService creates a new RecipeService instance
func NewRecipeService(recipeRepo RecipeRepository) *RecipeService {
	return &RecipeService{
		repo: recipeRepo,
	}
}

// Find contains the business logic to get all recipes
func (s RecipeService) FindAll() ([]m.RecipeDTO, error) {
	var recipes []m.Recipe

	recipes, err := s.repo.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	return m.Recipe{}.ConvertAllToDTO(recipes), nil
}

// Find contains the business logic to get a specific recipe
func (s RecipeService) FindSingle(recipeDTO m.RecipeDTO) (m.RecipeDTO, error) {
	var recipe m.Recipe

	recipe, err := s.repo.FindSingle(recipeDTO.ConvertFromDTO())
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.RecipeDTO{}, err
		default:
			return m.RecipeDTO{}, errors.New("internal server error")
		}
	}

	return recipe.ConvertToDTO(), nil
}

// Create handles the business logic for the creation of a recipe and passes the recipe object to the recipe repo for processing
func (s RecipeService) Create(recipeDTO m.RecipeDTO) (m.RecipeDTO, error) {

	if recipeDTO.ID != uuid.Nil {
		return m.RecipeDTO{}, errors.New("existing id on new element is not allowed")
	}

	if recipeDTO.Name == "" {
		return m.RecipeDTO{}, errors.New("name is empty")
	}

	if recipeDTO.Description == "" {
		return m.RecipeDTO{}, errors.New("description is empty")
	}

	if recipeDTO.ServingCount == 0 {
		return m.RecipeDTO{}, errors.New("serving count 0 is not allowed")
	}

	recipe, err := s.repo.Create(recipeDTO.ConvertFromDTO())
	if err != nil {
		return m.RecipeDTO{}, err
	}

	return recipe.ConvertToDTO(), nil
}

func (s RecipeService) Update(recipeDTO m.RecipeDTO) (m.RecipeDTO, error) {
	var updatedRecipe m.Recipe
	var originalRecipe m.Recipe

	originalRecipe, err := s.repo.FindSingle(recipeDTO.ConvertFromDTO())
	if err != nil {
		return m.RecipeDTO{}, err
	}

	if recipeDTO.Name == "" {
		recipeDTO.Name = originalRecipe.Name
	}

	if recipeDTO.Description == "" {
		recipeDTO.Description = originalRecipe.Description
	}

	if recipeDTO.ServingCount == 0 {
		recipeDTO.ServingCount = originalRecipe.ServingCount
	}

	updatedRecipe, err = s.repo.Update(recipeDTO.ConvertFromDTO())
	if err != nil {
		return m.RecipeDTO{}, err
	}

	return updatedRecipe.ConvertToDTO(), nil
}

func (s RecipeService) Delete(recipeDTO m.RecipeDTO) error {
	// TODO create safety logic
	if err := s.repo.Delete(recipeDTO.ConvertFromDTO()); err != nil {
		return err
	}

	return nil
}
