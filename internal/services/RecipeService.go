package services

import (
	m "github.com/ihulsbus/cookbook/internal/models"
)

type RecipeRepository interface {
	FindAll() ([]m.Recipe, error)
	FindSingle(recipeID uint) (m.Recipe, error)
	FindInstruction(recipeID uint) ([]m.Instruction, error)
	CreateInstruction(instructions []m.Instruction) ([]m.Instruction, error)
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
		return nil, err
	}

	return recipes, nil
}

// Find contains the business logic to get a specific recipe
func (s RecipeService) FindSingle(recipeID int) (m.Recipe, error) {
	var recipe m.Recipe

	recipe, err := s.repo.FindSingle(uint(recipeID))
	if err != nil {
		return recipe, err
	}

	return recipe, nil
}

func (s RecipeService) FindInstruction(recipeID int) ([]m.Instruction, error) {
	// TODO create logic
	instructions, err := s.repo.FindInstruction(uint(recipeID))
	if err != nil {
		return nil, err
	}

	return instructions, nil
}

func (s RecipeService) CreateInstruction(instruction []m.Instruction) ([]m.Instruction, error) {
	// TODO create logic
	instruction, err := s.repo.CreateInstruction(instruction)
	if err != nil {
		return nil, err
	}

	return instruction, nil
}

// Create handles the business logic for the creation of a recipe and passes the recipe object to the recipe repo for processing
func (s RecipeService) Create(recipe m.Recipe) (m.Recipe, error) {

	recipe, err := s.repo.Create(recipe)
	if err != nil {
		return recipe, err
	}

	return recipe, nil
}

func (s RecipeService) Update(recipe m.Recipe) (m.Recipe, error) {
	var updatedRecipe m.Recipe
	var originalRecipe m.Recipe

	originalRecipe, err := s.repo.FindSingle(recipe.ID)
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

	updatedRecipe, err = s.repo.Update(recipe)
	if err != nil {
		return updatedRecipe, err
	}

	return updatedRecipe, nil
}

func (s RecipeService) Delete(recipe m.Recipe) error {
	// TODO create safety logic
	if err := s.repo.Delete(recipe); err != nil {
		return err
	}

	return nil
}
