package handlers

type Handlers struct {
	recipeService     RecipeService
	ingredientService IngredientService
	imageService      ImageService
	logger            LoggerInterface
}

type LoggerInterface interface {
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

func NewHandlers(recipes RecipeService, ingredients IngredientService, imageService ImageService, logger LoggerInterface      ) *Handlers {
	return &Handlers{
		recipeService:     recipes,
		imageService:      imageService,
		ingredientService: ingredients,
		logger:            logger,
	}
}
