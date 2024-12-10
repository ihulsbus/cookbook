package rabbitmq

import (
	"github.com/wagslane/go-rabbitmq"
)

type RabbitMQPublisher struct {
	publisher *rabbitmq.Publisher
}

func NewRabbitMQPublisher(conn *rabbitmq.Conn, exchangeName string) (*RabbitMQPublisher, error) {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(exchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQPublisher{publisher: publisher}, nil
}
