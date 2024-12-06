package rabbitmq

import (
	"github.com/wagslane/go-rabbitmq"
)

type RabbitMQRepository struct {
	conn      *rabbitmq.Conn
	Publisher *rabbitmq.Publisher
}

func NewRabbitMQRepository(conn *rabbitmq.Conn) *RabbitMQRepository {
	return &RabbitMQRepository{
		conn: conn,
	}
}

func (r RabbitMQRepository) NewPublisher(exchangeName string) error {
	var err error

	r.Publisher, err = rabbitmq.NewPublisher(
		r.conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(exchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r RabbitMQRepository) PublishMessage(queueName, routingKey, message string) error {
	err := r.Publisher.Publish(
		[]byte(message),
		[]string{"my_routing_key"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(queueName),
	)

	if err != nil {
		return err
	}

	return nil
}
