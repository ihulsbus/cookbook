package services

import (
	"errors"

	m "recipe-service/internal/models"
)

type RecipeRepository interface {
	FindAll() ([]m.Recipe, error)
	FindSingle(recipeID uint) (m.Recipe, error)
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
func (s RecipeService) FindAll() ([]m.Recipe, error) {
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

	return recipes, nil
}

// Find contains the business logic to get a specific recipe
func (s RecipeService) FindSingle(recipeID uint) (m.Recipe, error) {
	var recipe m.Recipe

	recipe, err := s.repo.FindSingle(uint(recipeID))
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.Recipe{}, err
		default:
			return m.Recipe{}, errors.New("internal server error")
		}
	}

	return recipe, nil
}

// Create handles the business logic for the creation of a recipe and passes the recipe object to the recipe repo for processing
func (s RecipeService) Create(recipe m.Recipe) (m.Recipe, error) {

	recipe, err := s.repo.Create(recipe)
	if err != nil {
		return recipe, err
	}

	return recipe, nil
}

func (s RecipeService) Update(recipe m.Recipe, recipeID uint) (m.Recipe, error) {
	var updatedRecipe m.Recipe
	var originalRecipe m.Recipe

	originalRecipe, err := s.repo.FindSingle(recipeID)
	if err != nil {
		return updatedRecipe, err
	}

	if recipe.RecipeName == "" {
		recipe.RecipeName = originalRecipe.RecipeName
	}

	if recipe.Description == "" {
		recipe.Description = originalRecipe.Description
	}

	if recipe.CookingTime == 0 {
		recipe.CookingTime = originalRecipe.CookingTime
	}

	if recipe.ServingCount == 0 {
		recipe.ServingCount = originalRecipe.ServingCount
	}

	if recipe.AuthorID != originalRecipe.AuthorID {
		recipe.AuthorID = originalRecipe.AuthorID
	}

	if recipe.ImageName != originalRecipe.ImageName {
		recipe.ImageName = originalRecipe.ImageName
	}

	updatedRecipe, err = s.repo.Update(recipe)
	if err != nil {
		return updatedRecipe, err
	}

	return updatedRecipe, nil
}

func (s RecipeService) Delete(recipeID uint) error {
	var recipe m.Recipe

	recipe.ID = recipeID

	// TODO create safety logic
	if err := s.repo.Delete(recipe); err != nil {
		return err
	}

	return nil
}
