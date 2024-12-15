package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/wagslane/go-rabbitmq"
)

// PublisherInterface defines the behavior of a publisher
type PublisherInterface interface {
	Publish(body []byte, routingKeys []string, options ...func(*rabbitmq.PublishOptions)) error
	Close()
}

type Publisher struct {
	publisher PublisherInterface
}

func NewPublisher(conn *rabbitmq.Conn, exchangeName string) (*Publisher, error) {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(exchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeDurable,
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{publisher: publisher}, nil
}

func (p *Publisher) publish(event Event) error {
	message, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload for event %s: %w", event.Name, err)
	}

	return p.publisher.Publish(
		message,
		[]string{event.RoutingKey},
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)
}

func (p *Publisher) Close() {
	if p.publisher != nil {
		p.publisher.Close()
	}
}

func (p *Publisher) PublishRecipeCreated(payload RecipePayload) error {
	var event = RecipeCreatedEvent
	event.Payload = payload
	return p.publish(event)
}

func (p *Publisher) PublishRecipeUpdated(payload RecipePayload) error {
	var event = RecipeUpdatedEvent
	event.Payload = payload
	return p.publish(event)
}

func (p *Publisher) PublishRecipeDeleted(payload RecipePayload) error {
	var event = RecipeDeletedEvent
	event.Payload = payload
	return p.publish(event)
}

func (p *Publisher) PublishImageUpdated(payload ImagePayload) error {
	var event = ImageUpdatedEvent
	event.Payload = payload
	return p.publish(event)
}
