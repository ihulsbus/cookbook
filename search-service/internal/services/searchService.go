package services

import (
	"context"
	"fmt"

	m "github.com/ihulsbus/cookbook/shared/models"
)

type RecipeRepository interface {
	Query(ctx context.Context, req m.SearchRequest, token string) ([]m.RecipeResult, error)
}

type IngredientRepository interface {
	Query(ctx context.Context, req m.SearchRequest, token string) ([]m.IngredientResult, error)
}

type MetadataRepository interface {
	Query(ctx context.Context, req m.SearchRequest, token string) ([]m.MetadataSearchResult, error)
}

type SearchService struct {
	recipe     RecipeRepository
	ingredient IngredientRepository
	metadata   MetadataRepository
}

func NewSearchService(recipe RecipeRepository, ingredient IngredientRepository, metadata MetadataRepository) *SearchService {
	return &SearchService{
		recipe:     recipe,
		ingredient: ingredient,
		metadata:   metadata,
	}
}

// Search fans out to recipe-service, ingredient-service and
// metadata-service. token is the caller's own bearer token
// ("Bearer ..."), forwarded to each so every downstream service's own
// authorization/visibility filtering applies (KIT-30) - this service
// does no filtering of its own, by design, since each service is the
// authority on what it owns.
func (s *SearchService) Search(ctx context.Context, req m.SearchRequest, token string) (m.SearchResult, error) {
	results := m.SearchResult{}

	// Query Recipe Service
	recipeResults, err := s.recipe.Query(ctx, req, token)
	if err != nil {
		return results, fmt.Errorf("failed to query recipe service: %w", err)
	}
	results.Recipes = recipeResults

	// Query Ingredient Service
	ingredientResults, err := s.ingredient.Query(ctx, req, token)
	if err != nil {
		return results, fmt.Errorf("failed to query ingredient service: %w", err)
	}
	results.Ingredients = ingredientResults

	// Query Metadata Service
	metadataResults, err := s.metadata.Query(ctx, req, token)
	if err != nil {
		return results, fmt.Errorf("failed to query metadata service: %w", err)
	}
	results.Metadata = metadataResults

	return results, nil
}
