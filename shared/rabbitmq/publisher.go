package rabbitmq

import (
	"encoding/json"
	"log"

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

func NewProducer(conn *rabbitmq.Conn, exchangeName string) (*Publisher, error) {
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

func (p *Publisher) PublishMessage(routingKey string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = p.publisher.Publish(
		body,
		[]string{routingKey},
		rabbitmq.WithPublishOptionsPersistentDelivery,
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Printf("Message published to routing key: %s", routingKey)
	return nil
}

func (p *Publisher) Close() {
	if p.publisher != nil {
		p.publisher.Close()
	}
}
