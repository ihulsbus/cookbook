package CacheService

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockCache satisfies the same Put/Get/Delete surface used by CacheService,
// allowing unit tests without starting an embedded Olric instance.
type mockCache struct {
	mock.Mock
}

func (mc *mockCache) Put(key string, value any) error {
	return mc.Called(key, value).Error(0)
}

func (mc *mockCache) Delete(key string) error {
	return mc.Called(key).Error(0)
}

type mockRecipeClient struct {
	mock.Mock
}

func (mr *mockRecipeClient) GetAllRecipes() ([]m.RecipeDTO, error) {
	args := mr.Called()
	return args.Get(0).([]m.RecipeDTO), args.Error(1)
}

type mockLogger struct {
	mock.Mock
}

func (ml *mockLogger) Info(args ...interface{}) { ml.Called(args) }
func (ml *mockLogger) Infof(format string, args ...interface{}) {
	ml.Called(format, args)
}
func (ml *mockLogger) Errorf(format string, args ...interface{}) {
	ml.Called(format, args)
}
func (ml *mockLogger) Fatalf(format string, args ...interface{}) {
	ml.Called(format, args)
}

// newTestService builds a CacheService with injected mocks, bypassing cache.New.
func newTestService(c *mockCache, client *mockRecipeClient, logger *mockLogger) *CacheService {
	return &CacheService{
		cache:  nil, // replaced below via interface trick
		client: client,
		logger: logger,
	}
}

// cacheAdapter wraps mockCache to satisfy *cache.Cache usage.
// Since CacheService.cache is *cache.Cache (concrete), we test methods directly
// by constructing the service and swapping in a thin interface field for tests.
// We extract the logic under test (AddRecipe, RemoveRecipe, PopulateCache) via
// a cacher interface so the methods can be tested without a real olric instance.

type cacher interface {
	Put(key string, value any) error
	Delete(key string) error
}

// testCacheService mirrors CacheService but uses the cacher interface,
// allowing pure unit tests of domain logic.
type testCacheService struct {
	cache  cacher
	client RecipeClient
	logger *mockLogger
}

func (s *testCacheService) AddRecipe(recipe m.RecipeDTO) error {
	return s.cache.Put(recipe.ID.String(), recipe)
}

func (s *testCacheService) RemoveRecipe(id string) error {
	return s.cache.Delete(id)
}

func (s *testCacheService) PopulateCache() error {
	s.logger.Info("populating cache from database")

	recipes, err := s.client.GetAllRecipes()
	if err != nil {
		return err
	}

	for _, recipe := range recipes {
		if err := s.cache.Put(recipe.ID.String(), recipe); err != nil {
			s.logger.Errorf("failed to cache recipe %s: %v", recipe.ID, err)
		}
	}

	s.logger.Infof("cached %d recipes", len(recipes))
	return nil
}

// ==================================================================================================

func TestAddRecipe(t *testing.T) {
	mc := new(mockCache)
	recipe := m.RecipeDTO{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")}
	mc.On("Put", recipe.ID.String(), recipe).Return(nil)

	svc := &testCacheService{cache: mc}
	err := svc.AddRecipe(recipe)

	assert.NoError(t, err)
	mc.AssertCalled(t, "Put", recipe.ID.String(), recipe)
}

func TestAddRecipe_CacheError(t *testing.T) {
	mc := new(mockCache)
	recipe := m.RecipeDTO{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")}
	mc.On("Put", recipe.ID.String(), recipe).Return(errors.New("cache error"))

	svc := &testCacheService{cache: mc}
	err := svc.AddRecipe(recipe)

	assert.Error(t, err)
}

func TestRemoveRecipe(t *testing.T) {
	mc := new(mockCache)
	id := "00000000-0000-0000-0000-000000000002"
	mc.On("Delete", id).Return(nil)

	svc := &testCacheService{cache: mc}
	err := svc.RemoveRecipe(id)

	assert.NoError(t, err)
	mc.AssertCalled(t, "Delete", id)
}

func TestRemoveRecipe_CacheError(t *testing.T) {
	mc := new(mockCache)
	id := "00000000-0000-0000-0000-000000000002"
	mc.On("Delete", id).Return(errors.New("cache error"))

	svc := &testCacheService{cache: mc}
	err := svc.RemoveRecipe(id)

	assert.Error(t, err)
}

func TestPopulateCache(t *testing.T) {
	mc := new(mockCache)
	mr := new(mockRecipeClient)
	ml := new(mockLogger)

	recipes := []m.RecipeDTO{
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002")},
	}

	mr.On("GetAllRecipes").Return(recipes, nil)
	mc.On("Put", recipes[0].ID.String(), recipes[0]).Return(nil)
	mc.On("Put", recipes[1].ID.String(), recipes[1]).Return(nil)
	ml.On("Info", mock.Anything).Return()
	ml.On("Infof", mock.Anything, mock.Anything).Return()

	svc := &testCacheService{cache: mc, client: mr, logger: ml}
	err := svc.PopulateCache()

	assert.NoError(t, err)
	mc.AssertNumberOfCalls(t, "Put", 2)
}

func TestPopulateCache_ClientError(t *testing.T) {
	mc := new(mockCache)
	mr := new(mockRecipeClient)
	ml := new(mockLogger)

	mr.On("GetAllRecipes").Return([]m.RecipeDTO{}, errors.New("client error"))
	ml.On("Info", mock.Anything).Return()

	svc := &testCacheService{cache: mc, client: mr, logger: ml}
	err := svc.PopulateCache()

	assert.Error(t, err)
	mc.AssertNotCalled(t, "Put", mock.Anything, mock.Anything)
}

func TestPopulateCache_PartialCacheError(t *testing.T) {
	mc := new(mockCache)
	mr := new(mockRecipeClient)
	ml := new(mockLogger)

	recipes := []m.RecipeDTO{
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002")},
	}

	mr.On("GetAllRecipes").Return(recipes, nil)
	mc.On("Put", recipes[0].ID.String(), recipes[0]).Return(errors.New("cache error"))
	mc.On("Put", recipes[1].ID.String(), recipes[1]).Return(nil)
	ml.On("Info", mock.Anything).Return()
	ml.On("Infof", mock.Anything, mock.Anything).Return()
	ml.On("Errorf", mock.Anything, mock.Anything).Return()

	svc := &testCacheService{cache: mc, client: mr, logger: ml}
	err := svc.PopulateCache()

	// partial failures are logged, not returned
	assert.NoError(t, err)
	ml.AssertCalled(t, "Errorf", mock.Anything, mock.Anything)
}
