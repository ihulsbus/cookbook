package rabbitmq

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/wagslane/go-rabbitmq"
)

type RabbitMQPublisher struct {
	publisher *rabbitmq.Publisher
}

type RecipeDeletedEvent struct {
	RecipeID uuid.UUID `json:"recipe_id"`
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

// TODO: implement Context
func (r RabbitMQPublisher) publishMessageToExchange(message interface{}, routingKeys []string) error {

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = r.publisher.Publish(data, routingKeys)
	if err != nil {
		return err
	}

	return nil
}

func (r RabbitMQPublisher) PublishRecipeDeleted(recipeID uuid.UUID) error {
	var event RecipeDeletedEvent = RecipeDeletedEvent{RecipeID: recipeID}
	var routingKeys []string = []string{"recipe.deleted"}

	err := r.publishMessageToExchange(event, routingKeys)
	if err != nil {
		return err
	}

	return nil
}
