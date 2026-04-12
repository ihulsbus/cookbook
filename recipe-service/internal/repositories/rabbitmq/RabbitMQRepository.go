package rabbitmq

import (
	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"github.com/wagslane/go-rabbitmq"
)

type RabbitMQRepository struct {
	publisher *rmq.Publisher
	logger    m.LoggerInterface
}

type RecipeDeletedEvent struct {
	RecipeID uuid.UUID `json:"recipe_id"`
}

func NewRabbitMQRepository(conn *rabbitmq.Conn, exchangeName string, logger m.LoggerInterface) (*RabbitMQRepository, error) {
	publisher, err := rmq.NewPublisher(conn, exchangeName)
	if err != nil {
		return nil, err
	}
	return &RabbitMQRepository{publisher: publisher, logger: logger}, nil
}

func (r RabbitMQRepository) RecipeUpdatedEvent(recipe models.Recipe) error {
	var payload rmq.RecipePayload = rmq.RecipePayload{
		ID:           recipe.ID,
		Name:         recipe.Name,
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}

	return r.publisher.PublishRecipeUpdated(payload)
}
