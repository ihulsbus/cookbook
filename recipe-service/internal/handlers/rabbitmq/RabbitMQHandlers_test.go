package handlers

import (
	"encoding/json"
	"testing"

	"github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wagslane/go-rabbitmq"
)

// MockRecipeService mocks the recipeService interface
type MockRecipeService struct {
	mock.Mock
}

func (m *MockRecipeService) FindAll(pagination models.PaginationRequest) (models.PaginatedResponse[models.RecipeDTO], error) {
	args := m.Called(pagination)
	return args.Get(0).(models.PaginatedResponse[models.RecipeDTO]), args.Error(1)
}

func (m *MockRecipeService) FindSingle(recipeDTO models.RecipeDTO) (models.RecipeDTO, error) {
	args := m.Called(recipeDTO)
	return args.Get(0).(models.RecipeDTO), args.Error(1)
}

func (m *MockRecipeService) Create(recipeDTO models.RecipeDTO) (models.RecipeDTO, error) {
	args := m.Called(recipeDTO)
	return args.Get(0).(models.RecipeDTO), args.Error(1)
}

func (m *MockRecipeService) Update(recipeDTO models.RecipeDTO) (models.RecipeDTO, error) {
	args := m.Called(recipeDTO)
	return args.Get(0).(models.RecipeDTO), args.Error(1)
}

func (m *MockRecipeService) Delete(recipeDTO models.RecipeDTO) error {
	args := m.Called(recipeDTO)
	return args.Error(0)
}

// MockLogger mocks the LoggerInterface
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

// MockConsumer mocks the rmq.Consumer
type MockConsumer struct {
	mock.Mock
}

func (m *MockConsumer) Run(handler func(d rabbitmq.Delivery) rabbitmq.Action) error {
	args := m.Called(handler)
	return args.Error(0)
}

func (m *MockConsumer) Close() {
	m.Called()
}

// TestRabbitMQConsumerHandler tests the consumer handler logic
func TestRabbitMQConsumerHandler(t *testing.T) {
	tests := []struct {
		name           string
		routingKey     string
		body           interface{}
		expectedAction rabbitmq.Action
		expectLog      string
	}{
		{"ImageCreated", "image.created", models.RecipeDTO{Name: "Test Recipe"}, rabbitmq.Ack, "Received message with routing key: image.created"},
		{"ImageUpdated", "image.updated", models.RecipeDTO{Name: "Updated Recipe"}, rabbitmq.Ack, "Received message with routing key: image.updated"},
		{"ImageDeleted", "image.deleted", models.RecipeDTO{Name: "Deleted Recipe"}, rabbitmq.Ack, "Received message with routing key: image.deleted"},
		{"InstructionCreated", "instruction.created", models.RecipeDTO{Name: "Instruction Created"}, rabbitmq.Ack, "Received message with routing key: instruction.created"},
		{"MetadataUpdated", "metadata.updated", models.RecipeDTO{Name: "Metadata Updated"}, rabbitmq.Ack, "Received message with routing key: metadata.updated"},
		{"InvalidJSON", "image.created", "invalid-json", rabbitmq.NackRequeue, "Failed to unmarshal message"},
		{"UnknownRoutingKey", "unknown.key", models.RecipeDTO{Name: "Unknown"}, rabbitmq.NackDiscard, "Discarding message. Unknown routing key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockRecipeService)
			mockLogger := new(MockLogger)

			handler, _ := NewRabbitMQHandler(mockService, mockLogger)

			var body []byte
			if strBody, ok := tt.body.(string); ok {
				body = []byte(strBody)
			} else {
				body, _ = json.Marshal(tt.body)
			}

			delivery := rabbitmq.Delivery{}
			delivery.RoutingKey = tt.routingKey
			delivery.Body = body

			// Expect log calls
			mockLogger.On("Infof", mock.Anything, mock.Anything).Maybe()
			mockLogger.On("Errorf", mock.Anything, mock.Anything).Maybe()
			mockLogger.On("Warnf", mock.Anything, mock.Anything).Maybe()

			// Act
			action := handler.rabbitMQConsumerHandler(delivery)

			// Assert
			assert.Equal(t, tt.expectedAction, action)

			if tt.expectLog != "" {
				mockLogger.AssertCalled(t, "Infof", "Received message with routing key: %s", []interface{}{tt.routingKey})
			}
			mockLogger.AssertExpectations(t)
		})
	}
}

// func TestRabbitMQConsumerHandler_ProcessError(t *testing.T) {
// 	// Arrange
// 	mockService := new(MockRecipeService)
// 	mockLogger := new(MockLogger)
// 	handler, _ := NewRabbitMQHandler(mockService, mockLogger)

// 	event := models.RecipeDTO{Name: "ErrorRecipe"}
// 	body, _ := json.Marshal(event)
// 	delivery := rabbitmq.Delivery{}
// 	delivery.RoutingKey = "image.created"
// 	delivery.Body = body

// 	mockLogger.On("Infof", mock.Anything, mock.Anything).Return()
// 	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()
// 	mockService.On("Create", event).Return(models.RecipeDTO{}, errors.New("processing error"))

// 	// Act
// 	result := handler.rabbitMQConsumerHandler(delivery)

// 	// Assert
// 	assert.Equal(t, rabbitmq.NackRequeue, result)
// 	mockService.AssertCalled(t, "Create", event)
// 	mockLogger.AssertCalled(t, "Errorf", mock.Anything, mock.Anything)
// }

// TestStartConsuming verifies that StartConsuming initializes the consumer
