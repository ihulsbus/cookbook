package CacheService

import (
	"context"

	"github.com/ihulsbus/cookbook/shared/cache"
	m "github.com/ihulsbus/cookbook/shared/models"
)

type RecipeClient interface {
	GetAllRecipes(pagination ...m.PaginationRequest) (m.PaginatedResponse[m.RecipeDTO], error)
}

type CacheService struct {
	cache  *cache.Cache
	client RecipeClient
	logger cache.Logger
}

func NewCacheService(ctx context.Context, client RecipeClient, logger cache.Logger, bindPort int) (*CacheService, error) {
	c, err := cache.New(ctx, logger, "recipes", bindPort)
	if err != nil {
		return nil, err
	}

	service := &CacheService{
		cache:  c,
		client: client,
		logger: logger,
	}

	if err := service.PopulateCache(); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *CacheService) PopulateCache() error {
	s.logger.Info("populating cache from database")

	result, err := s.client.GetAllRecipes()
	if err != nil {
		return err
	}

	for _, recipe := range result.Data {
		if err := s.cache.Put(recipe.ID.String(), recipe); err != nil {
			s.logger.Errorf("failed to cache recipe %s: %v", recipe.ID, err)
		}
	}

	s.logger.Infof("cached %d recipes", len(result.Data))
	return nil
}

func (s *CacheService) AddRecipe(recipe m.RecipeDTO) error {
	return s.cache.Put(recipe.ID.String(), recipe)
}

func (s *CacheService) RemoveRecipe(id string) error {
	return s.cache.Delete(id)
}
