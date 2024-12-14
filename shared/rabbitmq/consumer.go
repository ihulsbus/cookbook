package rabbitmq

import (
	"github.com/wagslane/go-rabbitmq"
)

type Consumer struct {
	consumer *rabbitmq.Consumer
}

// Setup a new RabbitMQ consumer that handles incoming messages.
func NewConsumer(connection *rabbitmq.Conn, queueName, routingKey, exchangeName string, handler func(d rabbitmq.Delivery) rabbitmq.Action) (*Consumer, error) {
	consumer, err := rabbitmq.NewConsumer(
		connection,
		queueName,
		rabbitmq.WithConsumerOptionsRoutingKey(routingKey),
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsQueueQuorum,
		rabbitmq.WithConsumerOptionsExchangeName(exchangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, err
	}

	err = consumer.Run(handler)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: consumer}, nil
}

func (c *Consumer) Close() {
	if c.consumer != nil {
		c.consumer.Close()
	}
}
