package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"github.com/wagslane/go-rabbitmq"
)

type imageService interface {
	FindAll() ([]models.ImageDataDTO, error)
	Find(imageDTO models.ImageDataDTO) (models.ImageDataDTO, error)
	Delete(imageDTO models.ImageDataDTO) error
}

type RabbitMQHandler struct {
	service  imageService
	logger   models.LoggerInterface
	consumer *rmq.Consumer
}

func NewRabbitMQHandler(imageService imageService, logger models.LoggerInterface) (*RabbitMQHandler, error) {
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

func (c *RabbitMQHandler) StopConsuming() {
	if c.consumer != nil {
		c.consumer.Close()
	}
}

func (c *RabbitMQHandler) rabbitMQConsumerHandler(d rabbitmq.Delivery) rabbitmq.Action {
	var err error

	routingKey := d.RoutingKey
	c.logger.Infof("Received message with routing key: %s", routingKey)

	var event models.ImageDataDTO
	if err = json.Unmarshal(d.Body, &event); err != nil {
		c.logger.Errorf("Failed to unmarshal message: %v", err)
		return rabbitmq.NackRequeue
	}

	// TODO: Implement response feature
	switch routingKey {
	case "image.findall":
		_, err = c.service.FindAll()
	case "image.find":
		_, err = c.service.Find(models.ImageDataDTO{})
	case "recipe.deleted":
		fmt.Println(string(d.Body))
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
