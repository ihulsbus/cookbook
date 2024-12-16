package rabbitmq

import (
	"github.com/wagslane/go-rabbitmq"
)

type Consumer struct {
	consumer *rabbitmq.Consumer
}

// Setup a new RabbitMQ consumer that handles incoming messages.
func NewConsumer(connection *rabbitmq.Conn, queueName string, routingKeys []string, exchangeName string, handler func(d rabbitmq.Delivery) rabbitmq.Action) (*Consumer, error) {
	options := []func(*rabbitmq.ConsumerOptions){
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsQueueQuorum,
		rabbitmq.WithConsumerOptionsExchangeName(exchangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeDurable,
	}

	// add routing keys dynamically
	for _, key := range routingKeys {
		options = append(options, rabbitmq.WithConsumerOptionsRoutingKey(key))
	}

	consumer, err := rabbitmq.NewConsumer(
		connection,
		queueName,
		options...,
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
