package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/wagslane/go-rabbitmq"
)

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

// Publish marshals event.Payload and sends it to event.RoutingKey.
// event_id and schema_version travel as AMQP headers, not in the JSON
// body; a zero EventID gets one generated here so callers don't have to
// remember to set it.
func (p *Publisher) Publish(event Event) error {
	message, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload for event %s: %w", event.Name, err)
	}

	eventID := event.EventID
	if eventID == uuid.Nil {
		// Time-ordered (v7), not random v4; mirrors shared/models.NewID
		// without depending on that package from here.
		eventID = uuid.Must(uuid.NewV7())
	}

	return p.publisher.Publish(
		message,
		[]string{event.RoutingKey},
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsHeaders(rabbitmq.Table{
			"event_id":       eventID.String(),
			"schema_version": event.SchemaVersion,
		}),
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
	return p.Publish(event)
}

func (p *Publisher) PublishRecipeUpdated(payload RecipePayload) error {
	var event = RecipeUpdatedEvent
	event.Payload = payload
	return p.Publish(event)
}

func (p *Publisher) PublishRecipeDeleted(payload RecipePayload) error {
	var event = RecipeDeletedEvent
	event.Payload = payload
	return p.Publish(event)
}

func (p *Publisher) PublishImageUpdated(payload ImagePayload) error {
	var event = ImageUpdatedEvent
	event.Payload = payload
	return p.Publish(event)
}

func (p *Publisher) PublishImageDeleted(payload ImagePayload) error {
	var event = ImageDeletedEvent
	event.Payload = payload
	return p.Publish(event)
}

func (p *Publisher) PublishHouseholdMemberAdded(payload HouseholdMemberPayload) error {
	var event = HouseholdMemberAddedEvent
	event.Payload = payload
	return p.Publish(event)
}

func (p *Publisher) PublishHouseholdMemberRemoved(payload HouseholdMemberPayload) error {
	var event = HouseholdMemberRemovedEvent
	event.Payload = payload
	return p.Publish(event)
}
