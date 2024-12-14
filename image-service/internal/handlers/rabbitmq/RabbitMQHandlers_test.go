package handlers

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wagslane/go-rabbitmq"

	"image-service/internal/models"
)

// MockImageService is a mock for the imageService interface.
type MockImageService struct {
	mock.Mock
}

func (m *MockImageService) FindAll() ([]models.ImageDataDTO, error) {
	args := m.Called()
	return args.Get(0).([]models.ImageDataDTO), args.Error(1)
}

func (m *MockImageService) Find(imageDTO models.ImageDataDTO) (models.ImageDataDTO, error) {
	args := m.Called(imageDTO)
	return args.Get(0).(models.ImageDataDTO), args.Error(1)
}

func (m *MockImageService) Create(imageDTO models.ImageDataDTO) (models.ImageDataDTO, error) {
	args := m.Called(imageDTO)
	return args.Get(0).(models.ImageDataDTO), args.Error(1)
}

func (m *MockImageService) Update(imageDTO models.ImageDataDTO) (models.ImageDataDTO, error) {
	args := m.Called(imageDTO)
	return args.Get(0).(models.ImageDataDTO), args.Error(1)
}

func (m *MockImageService) Delete(imageDTO models.ImageDataDTO) error {
	return m.Called(imageDTO).Error(0)
}

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

// ==================================================================================================

func TestRabbitMQConsumerHandler(t *testing.T) {
	// Arrange
	mockService := new(MockImageService)
	mockLogger := new(MockLogger)
	consumer, _ := NewRabbitMQHandler(mockService, mockLogger)

	event := models.ImageDataDTO{
		ID:         uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		EntityType: "",
		EntityID:   uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Size:       0,
		Type:       "",
	}
	body, _ := json.Marshal(event)

	delivery := rabbitmq.Delivery{}
	delivery.RoutingKey = "image.find"
	delivery.Body = body

	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()
	mockService.On("Find", event).Return(event, nil)

	// Act
	result := consumer.rabbitMQConsumerHandler(delivery)

	// Assert
	assert.Equal(t, rabbitmq.Ack, result)
	mockService.AssertCalled(t, "Find", event)
	mockLogger.AssertCalled(t, "Infof", mock.Anything, mock.Anything)
}

func TestRabbitMQConsumerHandler_UnmarshalError(t *testing.T) {
	// Arrange
	mockService := new(MockImageService)
	mockLogger := new(MockLogger)
	consumer, _ := NewRabbitMQHandler(mockService, mockLogger)

	delivery := rabbitmq.Delivery{}
	delivery.RoutingKey = "image.created"
	delivery.Body = []byte("invalid-json")

	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	// Act
	result := consumer.rabbitMQConsumerHandler(delivery)

	// Assert
	assert.Equal(t, rabbitmq.NackRequeue, result)
	mockLogger.AssertCalled(t, "Errorf", mock.Anything, mock.Anything)
}

func TestRabbitMQConsumerHandler_UnknownRoutingKey(t *testing.T) {
	// Arrange
	mockService := new(MockImageService)
	mockLogger := new(MockLogger)
	consumer, _ := NewRabbitMQHandler(mockService, mockLogger)

	delivery := rabbitmq.Delivery{}
	delivery.RoutingKey = "unknown.key"
	delivery.Body = []byte(`{"entity_type":"test-entity","entity_id":"5d8f29ab-c749-47c5-86ed-5ea4881383ea","size":12345,"type":"image/jpeg"}`)

	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Return()

	// Act
	result := consumer.rabbitMQConsumerHandler(delivery)

	// Assert
	assert.Equal(t, rabbitmq.NackDiscard, result)
	mockLogger.AssertCalled(t, "Warnf", mock.Anything, mock.Anything)
}

func TestRabbitMQConsumerHandler_ProcessError(t *testing.T) {
	// Arrange
	mockService := new(MockImageService)
	mockLogger := new(MockLogger)
	consumer, _ := NewRabbitMQHandler(mockService, mockLogger)

	event := models.ImageDataDTO{
		ID:         uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		EntityType: "",
		EntityID:   uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Size:       0,
		Type:       "",
	}
	body, _ := json.Marshal(event)

	delivery := rabbitmq.Delivery{}
	delivery.RoutingKey = "image.find"
	delivery.Body = body

	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()
	mockService.On("Find", event).Return(models.ImageDataDTO{}, errors.New("processing error"))
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	// Act
	result := consumer.rabbitMQConsumerHandler(delivery)

	// Assert
	assert.Equal(t, rabbitmq.NackRequeue, result)
	mockService.AssertCalled(t, "Find", event)
	mockLogger.AssertCalled(t, "Errorf", mock.Anything, mock.Anything)
}

// func TestRabbitMQConsumerHandler_CreateError(t *testing.T) {
// 	// Arrange
// 	mockService := new(MockImageService)
// 	mockLogger := new(MockLogger)
// 	consumer, _ := NewRabbitMQConsumer(mockService, mockLogger)

// 	event := models.ImageDataDTO{
// 		ID:         uuid.New(),
// 		EntityType: "test-entity",
// 		EntityID:   uuid.New(),
// 		Size:       12345,
// 		Type:       "image/jpeg",
// 	}
// 	body, _ := json.Marshal(event)

// 	delivery := rabbitmq.Delivery{}
// 	delivery.RoutingKey = "image.created"
// 	delivery.Body = body

// 	errorMessage := errors.New("creation failed")
// 	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
// 	mockService.On("Create", event).Return(models.ImageDataDTO{}, errorMessage)
// 	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

// 	// Act
// 	result := consumer.rabbitMQConsumerHandler(delivery)

// 	// Assert
// 	assert.Equal(t, rabbitmq.NackRequeue, result)
// 	mockService.AssertCalled(t, "Create", event)
// 	mockLogger.AssertCalled(t, "Errorf", mock.Anything, mock.Anything)
// }
