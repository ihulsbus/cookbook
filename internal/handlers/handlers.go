package handlers

import (
	log "github.com/sirupsen/logrus"

	s "github.com/ihulsbus/cookbook/internal/services"
)

type Handlers struct {
	recipeService     *s.RecipeService
	ingredientService *s.IngredientService
	imageService      *s.ImageService
	logger            *log.Logger
}

func NewHandlers(recipes *s.RecipeService, ingredients *s.IngredientService, imageService *s.ImageService, logger *log.Logger) *Handlers {
	return &Handlers{
		recipeService:     recipes,
		imageService:      imageService,
		ingredientService: ingredients,
		logger:            logger,
	}
}
