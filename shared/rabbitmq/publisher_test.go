package rabbitmq

import (
	"encoding/json"
	"errors"
	"testing"

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

	testMessage := map[string]string{"key": "value"}
	body, _ := json.Marshal(testMessage)
	routingKey := "test-routing"

	mockPublisher.On("Publish", body, []string{routingKey}, mock.Anything).Return(nil)

	err := publisher.Publish(routingKey, testMessage)

	assert.NoError(t, err)
	mockPublisher.AssertCalled(t, "Publish", body, []string{routingKey}, mock.Anything)
}

func TestPublishMessage_MarshalError(t *testing.T) {
	publisher := &Publisher{}

	invalidMessage := make(chan int) // Channels cannot be marshaled to JSON

	err := producer.Publish("test-routing", invalidMessage)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "json: unsupported type")
}

func TestPublishMessage_PublishError(t *testing.T) {
	mockPublisher := new(MockPublisher)
	publisher := &Publisher{publisher: mockPublisher}

	testMessage := map[string]string{"key": "value"}
	body, _ := json.Marshal(testMessage)
	routingKey := "test-routing"

	mockPublisher.On("Publish", body, []string{routingKey}, mock.Anything).Return(errors.New("publish failed"))

	err := publisher.Publish(routingKey, testMessage)

	assert.Error(t, err)
	assert.Equal(t, "publish failed", err.Error())
	mockPublisher.AssertCalled(t, "Publish", body, []string{routingKey}, mock.Anything)
}

func TestClose(t *testing.T) {
	mockPublisher := new(MockPublisher)
	publisher := &Publisher{publisher: mockPublisher}

	mockPublisher.On("Close").Return()

	publisher.Close()

	mockPublisher.AssertCalled(t, "Close")
}
