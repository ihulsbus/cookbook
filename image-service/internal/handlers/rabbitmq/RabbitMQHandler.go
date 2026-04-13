package handlers

import (
	"context"
	"encoding/json"

	"github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"github.com/olric-data/olric"
	"github.com/wagslane/go-rabbitmq"
)

type imageService interface {
	FindAll() ([]models.ImageDataDTO, error)
	Find(imageDTO models.ImageDataDTO) (models.ImageDataDTO, error)
	Delete(imageDTO models.ImageDataDTO) error
}

type RabbitMQHandler struct {
	service  imageService
	cache    olric.DMap
	logger   models.LoggerInterface
	ctx      *context.Context
	consumer *rmq.Consumer // created in this package
}

func NewRabbitMQHandler(imageService imageService, cache *olric.DMap, ctx *context.Context, logger models.LoggerInterface) (*RabbitMQHandler, error) {
	return &RabbitMQHandler{service: imageService, cache: *cache, ctx: ctx, logger: logger}, nil
}

func (c *RabbitMQHandler) StartConsuming(connection *rabbitmq.Conn, queueName, exchangeName string) error {
	var err error

	var routingKeys = []string{"image.findall", "image.find", "recipe.created", "recipe.deleted"}

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
	case "image.findall":
		_, err = c.service.FindAll()
	case "image.find":
		_, err = c.service.Find(models.ImageDataDTO{})
	case "recipe.created":
		var event models.RecipeDTO
		if err = json.Unmarshal(d.Body, &event); err != nil {
			c.logger.Errorf("Failed to unmarshal message: %v", err)
			return rabbitmq.NackRequeue
		}

		if err := c.cache.Put(*c.ctx, event.ID.String(), event); err != nil {
			c.logger.Errorf("Failed to put event into cache: %v", err)
			return rabbitmq.NackRequeue
		}

		c.logger.Debugf("added recipe %s into the cache", event.ID.String())
	case "recipe.deleted":
		var event models.RecipeDTO
		if err = json.Unmarshal(d.Body, &event); err != nil {
			c.logger.Errorf("Failed to unmarshal message: %v", err)
			return rabbitmq.NackRequeue
		}

		count, err := c.cache.Delete(*c.ctx, event.ID.String())
		if err != nil {
			c.logger.Errorf("Failed to put event into cache: %v", err)
			return rabbitmq.NackRequeue
		}

		c.logger.Debugf("deleted %d instance(s) from the cache", count)
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
