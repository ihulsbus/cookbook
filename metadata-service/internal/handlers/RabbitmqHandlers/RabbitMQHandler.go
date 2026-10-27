package handlers

import (
	"context"
	"encoding/json"

	m "github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"github.com/wagslane/go-rabbitmq"
)

type CacheService interface {
	AddRecipe(recipe m.RecipeDTO) error
	RemoveRecipe(id string) error
}

type RabbitMQHandler struct {
	cache    CacheService
	logger   m.LoggerInterface
	ctx      *context.Context
	consumer *rmq.Consumer // created in this package
}

func NewRabbitMQHandler(cache CacheService, ctx *context.Context, logger m.LoggerInterface) (*RabbitMQHandler, error) {
	return &RabbitMQHandler{cache: cache, ctx: ctx, logger: logger}, nil
}

func (c *RabbitMQHandler) StartConsuming(connection *rabbitmq.Conn, queueName, exchangeName string) error {
	var err error

	var routingKeys = []string{"recipe.created", "recipe.deleted"}

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

	// TODO: Implement response feature
	switch routingKey {
	case "recipe.created":
		var event m.RecipeDTO
		if err = json.Unmarshal(d.Body, &event); err != nil {
			c.logger.Errorf("Failed to unmarshal message: %v", err)
			return rabbitmq.NackRequeue
		}

		if err := c.cache.AddRecipe(event); err != nil {
			c.logger.Errorf("Failed to put event into cache: %v", err)
			return rabbitmq.NackRequeue
		}

		c.logger.Debugf("added recipe %s into the cache", event.ID.String())
	case "recipe.deleted":
		var event m.RecipeDTO
		if err = json.Unmarshal(d.Body, &event); err != nil {
			c.logger.Errorf("Failed to unmarshal message: %v", err)
			return rabbitmq.NackRequeue
		}

		err := c.cache.RemoveRecipe(event.ID.String())
		if err != nil {
			c.logger.Errorf("Failed to put event into cache: %v", err)
			return rabbitmq.NackRequeue
		}
	default:
		c.logger.Warnf("Discarding message. Unknown routing key received: %s", routingKey)
		return rabbitmq.NackDiscard
	}

	if err != nil {
		c.logger.Errorf("Failed to process event: %v", err)
		return rabbitmq.NackRequeue
	}

	return rabbitmq.Ack
}
