package handlers

import (
	log "github.com/sirupsen/logrus"

	s "github.com/ihulsbus/cookbook/internal/services"
)

type Handlers struct {
	recipeService     *s.RecipeService
	ingredientService *s.IngredientService
	logger            *log.Logger
}

func NewHandlers(recipes *s.RecipeService, ingredients *s.IngredientService, logger *log.Logger) *Handlers {
	return &Handlers{
		recipeService:     recipes,
		logger:            logger,
		ingredientService: ingredients,
	}
}
