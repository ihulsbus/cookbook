package rabbitmq

import (
	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"github.com/wagslane/go-rabbitmq"
)

type RabbitMQRepository struct {
	publisher *rmq.Publisher
	logger    models.LoggerInterface
}

type RecipeDeletedEvent struct {
	RecipeID uuid.UUID `json:"recipe_id"`
}

func NewRabbitMQRepository(conn *rabbitmq.Conn, exchangeName string, logger models.LoggerInterface) (*RabbitMQRepository, error) {
	publisher, err := rmq.NewPublisher(conn, exchangeName)
	if err != nil {
		return nil, err
	}
	return &RabbitMQRepository{publisher: publisher, logger: logger}, nil
}

func (r RabbitMQRepository) ImageUpdatedEvent(image models.ImageData) error {
	var payload rmq.ImagePayload = rmq.ImagePayload{
		ID:         image.ID,
		EntityID:   image.EntityID,
		EntityType: image.EntityType,
		Size:       image.Size,
		Type:       image.Type,
	}

	return r.publisher.PublishImageUpdated(payload)
}
