package services

import (
	"errors"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
)

type RecipeRepository interface {
	FindAll() ([]m.Recipe, error)
	FindSingle(recipe m.Recipe) (m.Recipe, error)
	Create(recipe m.Recipe) (m.Recipe, error)
	Update(recipe m.Recipe) (m.Recipe, error)
	Delete(recipe m.Recipe) error
}

type RabbitMQRepository interface {
	RecipeCreatedEvent(recipe m.Recipe) error
	RecipeUpdatedEvent(recipe m.Recipe) error
	RecipeDeletedEvent(recipeID uuid.UUID) error
}

type RecipeService struct {
	recipe   RecipeRepository
	rabbitmq RabbitMQRepository
	logger   m.LoggerInterface
}

// NewRecipeService creates a new RecipeService instance
func NewRecipeService(recipeRepo RecipeRepository, rabbitMQRepo RabbitMQRepository, logger m.LoggerInterface) *RecipeService {
	return &RecipeService{
		recipe:   recipeRepo,
		rabbitmq: rabbitMQRepo,
		logger:   logger,
	}
}

// Find contains the business logic to get all recipes
func (s RecipeService) FindAll() ([]m.RecipeDTO, error) {
	var recipes []m.Recipe

	recipes, err := s.recipe.FindAll()
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

	recipe, err := s.recipe.FindSingle(recipeDTO.ConvertFromDTO())
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
	s.logger.Infof("%+v", recipeDTO)

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

	createdRecipe, err := s.recipe.Create(recipeDTO.ConvertFromDTO())
	if err != nil {
		return m.RecipeDTO{}, err
	}

	if err := s.rabbitmq.RecipeCreatedEvent(createdRecipe); err != nil {
		s.logger.Errorf("failed publishing recipe created event to servicebus")
	}

	return createdRecipe.ConvertToDTO(), nil
}

func (s RecipeService) Update(recipeDTO m.RecipeDTO) (m.RecipeDTO, error) {
	var updatedRecipe m.Recipe
	var originalRecipe m.Recipe

	originalRecipe, err := s.recipe.FindSingle(recipeDTO.ConvertFromDTO())
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

	updatedRecipe, err = s.recipe.Update(recipeDTO.ConvertFromDTO())
	if err != nil {
		return m.RecipeDTO{}, err
	}

	if err := s.rabbitmq.RecipeUpdatedEvent(updatedRecipe); err != nil {
		s.logger.Errorf("failed publishing recipe update event to servicebus")
	}

	return updatedRecipe.ConvertToDTO(), nil
}

func (s RecipeService) Delete(recipeDTO m.RecipeDTO) error {
	// TODO create safety logic
	if err := s.recipe.Delete(recipeDTO.ConvertFromDTO()); err != nil {
		return err
	}

	if err := s.rabbitmq.RecipeDeletedEvent(recipeDTO.ID); err != nil {
		s.logger.Errorf("failed publishing recipe deleted event to servicebus")
	}

	return nil
}
