package CacheService

import (
	"context"

	"github.com/ihulsbus/cookbook/shared/cache"
	m "github.com/ihulsbus/cookbook/shared/models"
	rc "github.com/ihulsbus/cookbook/shared/recipeclient"
	"github.com/olric-data/olric"
)

type CacheService struct {
	cache  *olric.DMap
	client *rc.RecipeAPIClient
	logger m.LoggerInterface
}

func NewCacheService(ctx context.Context, client *rc.RecipeAPIClient, logger m.LoggerInterface) (*CacheService, error) {
	cacheModule, err := cache.NewCacheModule(ctx, logger)
	if err != nil {
		return nil, err
	}

	cacheClient := cacheModule.NewEmbeddedClient()
	dmap, err := cacheClient.NewDMap("recipes")
	if err != nil {
		return nil, err
	}

	service := &CacheService{
		cache:  &dmap,
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

	recipes, err := s.client.GetAllRecipes() // or whatever method
	if err != nil {
		return err
	}

	for _, recipe := range recipes {
		if err := s.cache.Put(recipe.ID, recipe); err != nil {
			s.logger.Errorf("failed to cache recipe %s: %v", recipe.ID, err)
		}
	}

	s.logger.Infof("cached %d recipes", len(recipes))
	return nil
}

// Called by RabbitMQ handlers
func (s *CacheService) AddRecipe(recipe m.RecipeDTO) error {
	return s.cache.Put(recipe.ID, recipe)
}

func (s *CacheService) RemoveRecipe(id string) error {
	return s.cache.Delete(id)
}
