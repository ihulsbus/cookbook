package services

import (
	"context"
	"fmt"

	m "search-service/internal/models"
)

type RecipeRepository interface {
	Query(ctx context.Context, req m.SearchRequest) ([]m.RecipeResult, error)
}

type IngredientRepository interface {
	Query(ctx context.Context, req m.SearchRequest) ([]m.IngredientResult, error)
}

type MetadataRepository interface {
	Query(ctx context.Context, req m.SearchRequest) ([]m.MetadataSearchResult, error)
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

func (s *SearchService) Search(ctx context.Context, req m.SearchRequest) (m.SearchResult, error) {
	results := m.SearchResult{}

	// Query Recipe Service
	recipeResults, err := s.recipe.Query(ctx, req)
	if err != nil {
		return results, fmt.Errorf("failed to query recipe service: %w", err)
	}
	results.Recipes = recipeResults

	// Query Ingredient Service
	ingredientResults, err := s.ingredient.Query(ctx, req)
	if err != nil {
		return results, fmt.Errorf("failed to query ingredient service: %w", err)
	}
	results.Ingredients = ingredientResults

	// Query Metadata Service
	metadataResults, err := s.metadata.Query(ctx, req)
	if err != nil {
		return results, fmt.Errorf("failed to query metadata service: %w", err)
	}
	results.Metadata = metadataResults

	return results, nil
}
