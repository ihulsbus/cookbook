package handlers

import (
	"encoding/json"
	"fmt"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"github.com/wagslane/go-rabbitmq"
)

type recipeService interface {
	FindAll(models.PaginationRequest) (models.PaginatedResponse[models.RecipeDTO], error)
	FindSingle(recipeDTO models.RecipeDTO) (models.RecipeDTO, error)
	Create(recipeDTO models.RecipeDTO) (models.RecipeDTO, error)
	Update(recipeDTO models.RecipeDTO) (models.RecipeDTO, error)
	Delete(recipeDTO models.RecipeDTO) error
}

type RabbitMQHandler struct {
	service  recipeService
	logger   m.LoggerInterface
	consumer *rmq.Consumer
}

func NewRabbitMQHandler(recipeService recipeService, logger m.LoggerInterface) (*RabbitMQHandler, error) {
	return &RabbitMQHandler{service: recipeService, logger: logger}, nil
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
	case "image.created":
		fmt.Println(string(d.Body))
	case "image.updated":
		fmt.Println(string(d.Body))
	case "image.deleted":
		fmt.Println(string(d.Body))
	case "instruction.created":
		fmt.Println(string(d.Body))
	case "instruction.updated":
		fmt.Println(string(d.Body))
	case "instruction.deleted":
		fmt.Println(string(d.Body))
	case "metadata.created":
		fmt.Println(string(d.Body))
	case "metadata.updated":
		fmt.Println(string(d.Body))
	case "metadata.deleted":
		fmt.Println(string(d.Body))
	default:
		c.logger.Warnf("Discarding message. Unknown routing key received: %s", routingKey)
		return rabbitmq.NackDiscard
	}

	return rabbitmq.Ack
}
