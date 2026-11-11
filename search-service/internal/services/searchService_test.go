package services

import (
	"context"
	"errors"
	"testing"

	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
)

type recipeRepoMock struct{ gotToken string }
type ingredientRepoMock struct{ gotToken string }
type metadataRepoMock struct{ gotToken string }

var (
	failRecipe     bool
	failIngredient bool
	failMetadata   bool
)

func (r *recipeRepoMock) Query(ctx context.Context, req m.SearchRequest, token string) ([]m.RecipeResult, error) {
	r.gotToken = token
	if failRecipe {
		return nil, errors.New("recipe down")
	}
	return []m.RecipeResult{{Name: "apple pie"}}, nil
}

func (r *ingredientRepoMock) Query(ctx context.Context, req m.SearchRequest, token string) ([]m.IngredientResult, error) {
	r.gotToken = token
	if failIngredient {
		return nil, errors.New("ingredient down")
	}
	return []m.IngredientResult{{Name: "apple"}}, nil
}

func (r *metadataRepoMock) Query(ctx context.Context, req m.SearchRequest, token string) ([]m.MetadataSearchResult, error) {
	r.gotToken = token
	if failMetadata {
		return nil, errors.New("metadata down")
	}
	return []m.MetadataSearchResult{{}}, nil
}

func TestSearchService_Search_OK(t *testing.T) {
	failRecipe, failIngredient, failMetadata = false, false, false

	recipeRepo := &recipeRepoMock{}
	ingredientRepo := &ingredientRepoMock{}
	metadataRepo := &metadataRepoMock{}
	s := NewSearchService(recipeRepo, ingredientRepo, metadataRepo)

	result, err := s.Search(context.Background(), m.SearchRequest{Query: "apple"}, "Bearer test-token")

	assert.NoError(t, err)
	assert.Len(t, result.Recipes, 1)
	assert.Len(t, result.Ingredients, 1)
	assert.Len(t, result.Metadata, 1)
	assert.Equal(t, "Bearer test-token", recipeRepo.gotToken)
	assert.Equal(t, "Bearer test-token", ingredientRepo.gotToken)
	assert.Equal(t, "Bearer test-token", metadataRepo.gotToken)
}

func TestSearchService_Search_RecipeErr(t *testing.T) {
	failRecipe, failIngredient, failMetadata = true, false, false
	defer func() { failRecipe = false }()

	s := NewSearchService(&recipeRepoMock{}, &ingredientRepoMock{}, &metadataRepoMock{})

	_, err := s.Search(context.Background(), m.SearchRequest{Query: "apple"}, "Bearer test-token")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "recipe down")
}

func TestSearchService_Search_IngredientErr(t *testing.T) {
	failRecipe, failIngredient, failMetadata = false, true, false
	defer func() { failIngredient = false }()

	s := NewSearchService(&recipeRepoMock{}, &ingredientRepoMock{}, &metadataRepoMock{})

	_, err := s.Search(context.Background(), m.SearchRequest{Query: "apple"}, "Bearer test-token")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ingredient down")
}

func TestSearchService_Search_MetadataErr(t *testing.T) {
	failRecipe, failIngredient, failMetadata = false, false, true
	defer func() { failMetadata = false }()

	s := NewSearchService(&recipeRepoMock{}, &ingredientRepoMock{}, &metadataRepoMock{})

	_, err := s.Search(context.Background(), m.SearchRequest{Query: "apple"}, "Bearer test-token")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "metadata down")
}
