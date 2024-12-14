package handlers

import (
	"encoding/json"
	m "image-service/internal/models"

	"github.com/wagslane/go-rabbitmq"
)

type imageService interface {
	FindAll() ([]m.ImageDataDTO, error)
	Find(imageDTO m.ImageDataDTO) (m.ImageDataDTO, error)
	Delete(imageDTO m.ImageDataDTO) error
}

type RabbitMQConsumer struct {
	service imageService
	logger  m.LoggerInterface
}

func NewRabbitMQConsumer(imageService imageService, logger m.LoggerInterface) (*RabbitMQConsumer, error) {
	return &RabbitMQConsumer{service: imageService, logger: logger}, nil
}

func (c *RabbitMQConsumer) StartConsuming(connection *rabbitmq.Conn, queueName, exchangeName string) error {
	consumer, err := rabbitmq.NewConsumer(
		connection,
		queueName,
		rabbitmq.WithConsumerOptionsRoutingKey("image.created"),
		rabbitmq.WithConsumerOptionsExchangeName(exchangeName),
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsQueueQuorum,
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}

	err = consumer.Run(c.rabbitMQConsumerHandler)
	if err != nil {
		return err
	}

	defer consumer.Close()

	return nil
}

func (c *RabbitMQConsumer) rabbitMQConsumerHandler(d rabbitmq.Delivery) rabbitmq.Action {
	var err error

	routingKey := d.RoutingKey
	c.logger.Infof("Received message with routing key: %s", routingKey)

	var event m.ImageDataDTO
	if err = json.Unmarshal(d.Body, &event); err != nil {
		c.logger.Errorf("Failed to unmarshal message: %v", err)
		return rabbitmq.NackRequeue
	}

	switch routingKey {
	case "image.find":
		_, err = c.service.Find(m.ImageDataDTO{})
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
