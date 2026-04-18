package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wagslane/go-rabbitmq"
)

// MockLogger is a mock for the LoggerInterface.
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(format, args)
}

// MockCacheService is a mock for the CacheService interface.
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) AddRecipe(recipe models.RecipeDTO) error {
	args := m.Called(recipe)
	return args.Error(0)
}

func (m *MockCacheService) RemoveRecipe(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// ==================================================================================================

func TestRabbitMQConsumerHandler_RecipeCreated(t *testing.T) {
	// Arrange
	mockCache := new(MockCacheService)
	mockLogger := new(MockLogger)
	mockCtx := context.Background()
	consumer, _ := NewRabbitMQHandler(mockCache, &mockCtx, mockLogger)

	event := models.RecipeDTO{
		ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
	}
	body, _ := json.Marshal(event)

	delivery := rabbitmq.Delivery{}
	delivery.RoutingKey = "recipe.created"
	delivery.Body = body

	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockCache.On("AddRecipe", event).Return(nil)

	// Act
	result := consumer.rabbitMQConsumerHandler(delivery)

	// Assert
	assert.Equal(t, rabbitmq.Ack, result)
	mockCache.AssertCalled(t, "AddRecipe", event)
}

func TestRabbitMQConsumerHandler_RecipeCreated_UnmarshalError(t *testing.T) {
	// Arrange
	mockCache := new(MockCacheService)
	mockLogger := new(MockLogger)
	mockCtx := context.Background()
	consumer, _ := NewRabbitMQHandler(mockCache, &mockCtx, mockLogger)

	delivery := rabbitmq.Delivery{}
	delivery.RoutingKey = "recipe.created"
	delivery.Body = []byte("invalid-json")

	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	// Act
	result := consumer.rabbitMQConsumerHandler(delivery)

	// Assert
	assert.Equal(t, rabbitmq.NackRequeue, result)
	mockLogger.AssertCalled(t, "Errorf", mock.Anything, mock.Anything)
}

func TestRabbitMQConsumerHandler_RecipeCreated_CacheError(t *testing.T) {
	// Arrange
	mockCache := new(MockCacheService)
	mockLogger := new(MockLogger)
	mockCtx := context.Background()
	consumer, _ := NewRabbitMQHandler(mockCache, &mockCtx, mockLogger)

	event := models.RecipeDTO{
		ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
	}
	body, _ := json.Marshal(event)

	delivery := rabbitmq.Delivery{}
	delivery.RoutingKey = "recipe.created"
	delivery.Body = body

	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()
	mockCache.On("AddRecipe", event).Return(errors.New("cache error"))

	// Act
	result := consumer.rabbitMQConsumerHandler(delivery)

	// Assert
	assert.Equal(t, rabbitmq.NackRequeue, result)
	mockCache.AssertCalled(t, "AddRecipe", event)
	mockLogger.AssertCalled(t, "Errorf", mock.Anything, mock.Anything)
}

func TestRabbitMQConsumerHandler_RecipeDeleted(t *testing.T) {
	// Arrange
	mockCache := new(MockCacheService)
	mockLogger := new(MockLogger)
	mockCtx := context.Background()
	consumer, _ := NewRabbitMQHandler(mockCache, &mockCtx, mockLogger)

	id := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	event := models.RecipeDTO{ID: id}
	body, _ := json.Marshal(event)

	delivery := rabbitmq.Delivery{}
	delivery.RoutingKey = "recipe.deleted"
	delivery.Body = body

	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
	mockCache.On("RemoveRecipe", id.String()).Return(nil)

	// Act
	result := consumer.rabbitMQConsumerHandler(delivery)

	// Assert
	assert.Equal(t, rabbitmq.Ack, result)
	mockCache.AssertCalled(t, "RemoveRecipe", id.String())
}

func TestRabbitMQConsumerHandler_RecipeDeleted_UnmarshalError(t *testing.T) {
	// Arrange
	mockCache := new(MockCacheService)
	mockLogger := new(MockLogger)
	mockCtx := context.Background()
	consumer, _ := NewRabbitMQHandler(mockCache, &mockCtx, mockLogger)

	delivery := rabbitmq.Delivery{}
	delivery.RoutingKey = "recipe.deleted"
	delivery.Body = []byte("invalid-json")

	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	// Act
	result := consumer.rabbitMQConsumerHandler(delivery)

	// Assert
	assert.Equal(t, rabbitmq.NackRequeue, result)
	mockLogger.AssertCalled(t, "Errorf", mock.Anything, mock.Anything)
}

func TestRabbitMQConsumerHandler_RecipeDeleted_CacheError(t *testing.T) {
	// Arrange
	mockCache := new(MockCacheService)
	mockLogger := new(MockLogger)
	mockCtx := context.Background()
	consumer, _ := NewRabbitMQHandler(mockCache, &mockCtx, mockLogger)

	id := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	event := models.RecipeDTO{ID: id}
	body, _ := json.Marshal(event)

	delivery := rabbitmq.Delivery{}
	delivery.RoutingKey = "recipe.deleted"
	delivery.Body = body

	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()
	mockCache.On("RemoveRecipe", id.String()).Return(errors.New("cache error"))

	// Act
	result := consumer.rabbitMQConsumerHandler(delivery)

	// Assert
	assert.Equal(t, rabbitmq.NackRequeue, result)
	mockCache.AssertCalled(t, "RemoveRecipe", id.String())
	mockLogger.AssertCalled(t, "Errorf", mock.Anything, mock.Anything)
}

func TestRabbitMQConsumerHandler_UnknownRoutingKey(t *testing.T) {
	// Arrange
	mockCache := new(MockCacheService)
	mockLogger := new(MockLogger)
	mockCtx := context.Background()
	consumer, _ := NewRabbitMQHandler(mockCache, &mockCtx, mockLogger)

	delivery := rabbitmq.Delivery{}
	delivery.RoutingKey = "unknown.key"
	delivery.Body = []byte(`{}`)

	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Return()

	// Act
	result := consumer.rabbitMQConsumerHandler(delivery)

	// Assert
	assert.Equal(t, rabbitmq.NackDiscard, result)
	mockLogger.AssertCalled(t, "Warnf", mock.Anything, mock.Anything)
}
