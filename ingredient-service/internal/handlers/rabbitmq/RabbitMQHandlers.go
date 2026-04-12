package handlers

import (
	"encoding/json"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"github.com/wagslane/go-rabbitmq"
)

type RabbitMQHandler struct {
	logger   m.LoggerInterface
	consumer *rmq.Consumer
}

func NewRabbitMQHandler(logger m.LoggerInterface) (*RabbitMQHandler, error) {
	return &RabbitMQHandler{logger: logger}, nil
}

func (c *RabbitMQHandler) StartConsuming(connection *rabbitmq.Conn, queueName, exchangeName string) error {
	var err error

	var routingKeys []string = []string{"image.created", "image.updated", "image.deleted", "instruction.created", "instruction.updated", "instruction.deleted", "metadata.created", "metadata.updated", "metadata.deleted"}

	c.consumer, err = rmq.NewConsumer(connection, queueName, routingKeys, exchangeName, c.rabbitMQConsumerHandler)
	if err != nil {
		return err
	}

	return nil
}

func (c *RabbitMQHandler) StopConsuming() {
	if c.consumer != nil {
		c.consumer.Close()
	}
}

func (c *RabbitMQHandler) rabbitMQConsumerHandler(d rabbitmq.Delivery) rabbitmq.Action {
	var err error

	routingKey := d.RoutingKey
	c.logger.Infof("Received message with routing key: %s", routingKey)

	var event models.RecipeDTO
	if err = json.Unmarshal(d.Body, &event); err != nil {
		c.logger.Errorf("Failed to unmarshal message: %v", err)
		return rabbitmq.NackRequeue
	}

	// TODO: Implement response feature
	switch routingKey {
	case "recipe.created":
		c.logger.Infof("Recipe created: %s", string(d.Body))
	case "recipe.deleted":
		c.logger.Infof("Recipe deleted: %s", string(d.Body))
	default:
		c.logger.Warnf("Discarding message. Unknown routing key received: %s", routingKey)
		return rabbitmq.NackDiscard
	}

	return rabbitmq.Ack
}
