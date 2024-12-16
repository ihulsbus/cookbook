package handlers

import (
	"encoding/json"
	"fmt"
	m "image-service/internal/models"

	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"github.com/wagslane/go-rabbitmq"
)

type imageService interface {
	FindAll() ([]m.ImageDataDTO, error)
	Find(imageDTO m.ImageDataDTO) (m.ImageDataDTO, error)
	Delete(imageDTO m.ImageDataDTO) error
}

type RabbitMQHandler struct {
	service  imageService
	logger   m.LoggerInterface
	consumer *rmq.Consumer
}

func NewRabbitMQHandler(imageService imageService, logger m.LoggerInterface) (*RabbitMQHandler, error) {
	return &RabbitMQHandler{service: imageService, logger: logger}, nil
}

func (c *RabbitMQHandler) StartConsuming(connection *rabbitmq.Conn, queueName, exchangeName string) error {
	var err error

	var routingKeys []string = []string{"image.findall", "image.find", "recipe.deleted"}

	c.consumer, err = rmq.NewConsumer(connection, queueName, routingKeys, exchangeName, c.rabbitMQConsumerHandler)
	if err != nil {
		return err
	}

	return nil
}

func (c *RabbitMQHandler) rabbitMQConsumerHandler(d rabbitmq.Delivery) rabbitmq.Action {
	var err error

	routingKey := d.RoutingKey
	c.logger.Infof("Received message with routing key: %s", routingKey)

	var event m.ImageDataDTO
	if err = json.Unmarshal(d.Body, &event); err != nil {
		c.logger.Errorf("Failed to unmarshal message: %v", err)
		return rabbitmq.NackRequeue
	}

	// TODO: Implement response feature
	switch routingKey {
	case "image.findall":
		_, err = c.service.FindAll()
	case "image.find":
		_, err = c.service.Find(m.ImageDataDTO{})
	case "recipe.deleted":
		fmt.Println()
		return rabbitmq.Ack
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
