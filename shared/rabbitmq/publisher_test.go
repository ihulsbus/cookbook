package rabbitmq

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wagslane/go-rabbitmq"
)

// MockPublisher is a mock implementation of the rabbitmq.Publisher
type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(body []byte, routingKeys []string, options ...func(*rabbitmq.PublishOptions)) error {
	args := m.Called(body, routingKeys, options)
	return args.Error(0)
}

func (m *MockPublisher) Close() {
	m.Called()
}

func TestPublishMessage_Success(t *testing.T) {
	mockPublisher := new(MockPublisher)
	publisher := &Publisher{publisher: mockPublisher}

	testMessage := Event{
		Name:       "",
		RoutingKey: "test-routing",
		Payload:    "test",
	}
	message, _ := json.Marshal(testMessage.Payload)

	mockPublisher.On("Publish", message, []string{testMessage.RoutingKey}, mock.Anything).Return(nil)

	err := publisher.Publish(testMessage)

	assert.NoError(t, err)
	mockPublisher.AssertCalled(t, "Publish", message, []string{testMessage.RoutingKey}, mock.Anything)
}

func TestPublishMessage_MarshalError(t *testing.T) {
	publisher := &Publisher{}

	testMessage := Event{
		Name:       "",
		RoutingKey: "test-routing",
		Payload:    make(chan int), // Channels cannot be marshaled to JSON
	}

	err := publisher.Publish(testMessage)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json: unsupported type")
}

func TestPublishMessage_PublishError(t *testing.T) {
	mockPublisher := new(MockPublisher)
	publisher := &Publisher{publisher: mockPublisher}

	testMessage := Event{
		Name:       "",
		RoutingKey: "test-routing",
		Payload:    "test",
	}

	message, _ := json.Marshal(testMessage.Payload)

	mockPublisher.On("Publish", message, []string{testMessage.RoutingKey}, mock.Anything).Return(errors.New("publish failed"))

	err := publisher.Publish(testMessage)

	assert.Error(t, err)
	assert.Equal(t, "publish failed", err.Error())
	mockPublisher.AssertCalled(t, "Publish", message, []string{testMessage.RoutingKey}, mock.Anything)
}

func TestPublishRecipeCreatedMessage_Success(t *testing.T) {
	mockPublisher := new(MockPublisher)
	publisher := &Publisher{publisher: mockPublisher}

	var payload RecipePayload = RecipePayload{
		ID:           uuid.New(),
		Name:         "test",
		Description:  "description",
		ServingCount: 1,
	}
	message, _ := json.Marshal(payload)

	mockPublisher.On("Publish", message, []string{RecipeCreatedRoutingKey}, mock.Anything).Return(nil)

	err := publisher.PublishRecipeCreated(payload)

	assert.NoError(t, err)
	mockPublisher.AssertCalled(t, "Publish", message, []string{RecipeCreatedRoutingKey}, mock.Anything)
}

func TestPublishRecipeUpdatedMessage_Success(t *testing.T) {
	mockPublisher := new(MockPublisher)
	publisher := &Publisher{publisher: mockPublisher}

	var payload RecipePayload = RecipePayload{
		ID:           uuid.New(),
		Name:         "test",
		Description:  "description",
		ServingCount: 1,
	}
	message, _ := json.Marshal(payload)

	mockPublisher.On("Publish", message, []string{RecipeUpdatedRoutingKey}, mock.Anything).Return(nil)

	err := publisher.PublishRecipeUpdated(payload)

	assert.NoError(t, err)
	mockPublisher.AssertCalled(t, "Publish", message, []string{RecipeUpdatedRoutingKey}, mock.Anything)
}

func TestPublishRecipeDeletedMessage_Success(t *testing.T) {
	mockPublisher := new(MockPublisher)
	publisher := &Publisher{publisher: mockPublisher}

	var payload RecipePayload = RecipePayload{
		ID:           uuid.New(),
		Name:         "test",
		Description:  "description",
		ServingCount: 1,
	}
	message, _ := json.Marshal(payload)

	mockPublisher.On("Publish", message, []string{RecipeDeletedRoutingKey}, mock.Anything).Return(nil)

	err := publisher.PublishRecipeDeleted(payload)

	assert.NoError(t, err)
	mockPublisher.AssertCalled(t, "Publish", message, []string{RecipeDeletedRoutingKey}, mock.Anything)
}

func TestPublishImageUpdatedMessage_Success(t *testing.T) {
	mockPublisher := new(MockPublisher)
	publisher := &Publisher{publisher: mockPublisher}

	var payload ImagePayload = ImagePayload{
		ID:         uuid.New(),
		EntityType: "recipe",
		EntityID:   uuid.New(),
		Size:       0,
		Type:       "png",
	}
	message, _ := json.Marshal(payload)

	mockPublisher.On("Publish", message, []string{ImageUpdatedRoutingKey}, mock.Anything).Return(nil)

	err := publisher.PublishImageUpdated(payload)

	assert.NoError(t, err)
	mockPublisher.AssertCalled(t, "Publish", message, []string{ImageUpdatedRoutingKey}, mock.Anything)
}

func TestClose(t *testing.T) {
	mockPublisher := new(MockPublisher)
	publisher := &Publisher{publisher: mockPublisher}

	mockPublisher.On("Close").Return()

	publisher.Close()

	mockPublisher.AssertCalled(t, "Close")
}
