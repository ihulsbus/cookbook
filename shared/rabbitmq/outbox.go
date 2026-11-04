package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wagslane/go-rabbitmq"
	"gorm.io/gorm"
)

// LoggerInterface is the minimal logging surface outbox/relay code needs.
// Defined locally (rather than importing shared/models) so this module
// stays free of a dependency on a specific service's config stack.
type LoggerInterface interface {
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// OutboxEvent is a durable record of an event to publish, written in the
// same DB transaction as the business change it describes. A Relay polls
// for unpublished rows and publishes them, so a publish failure (e.g.
// RabbitMQ is down) no longer loses the event - it just stays unpublished
// and gets retried on the next poll, instead of being logged and dropped.
type OutboxEvent struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RoutingKey    string    `gorm:"not null"`
	Payload       []byte    `gorm:"type:jsonb;not null"`
	EventID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	SchemaVersion int       `gorm:"not null;default:1"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	PublishedAt   *time.Time
}

// NewOutboxEntry builds the outbox row for event, ready to be inserted in
// the same transaction as the change it describes. event.RoutingKey and
// event.Payload must already be set; EventID defaults to a fresh UUID and
// SchemaVersion to CurrentSchemaVersion when left zero.
func NewOutboxEntry(event Event) (OutboxEvent, error) {
	body, err := json.Marshal(event.Payload)
	if err != nil {
		return OutboxEvent{}, fmt.Errorf("failed to marshal payload for event %s: %w", event.Name, err)
	}

	eventID := event.EventID
	if eventID == uuid.Nil {
		eventID = uuid.New()
	}

	schemaVersion := event.SchemaVersion
	if schemaVersion == 0 {
		schemaVersion = CurrentSchemaVersion
	}

	return OutboxEvent{
		RoutingKey:    event.RoutingKey,
		Payload:       body,
		EventID:       eventID,
		SchemaVersion: schemaVersion,
	}, nil
}

// AutoMigrateOutbox migrates the outbox table into db. Call once at
// service startup alongside the service's own model migrations.
func AutoMigrateOutbox(db *gorm.DB) error {
	return db.AutoMigrate(&OutboxEvent{})
}

// Relay polls the outbox table for unpublished events and publishes them,
// marking each published on success. At-least-once delivery: a crash
// between a successful publish and marking it published re-publishes that
// event on the next poll, which is safe because every consumer in this
// system is expected to apply events idempotently via Version.
type Relay struct {
	db        *gorm.DB
	publisher *Publisher
	logger    LoggerInterface
	interval  time.Duration
	batchSize int
}

// NewRelay creates a Relay. interval is how often to poll for unpublished
// events; batchSize caps how many are published per poll.
func NewRelay(db *gorm.DB, publisher *Publisher, logger LoggerInterface, interval time.Duration, batchSize int) *Relay {
	if batchSize <= 0 {
		batchSize = 50
	}
	return &Relay{db: db, publisher: publisher, logger: logger, interval: interval, batchSize: batchSize}
}

// Start runs the poll loop until ctx is done. Intended to be run in its
// own goroutine by the caller.
func (r *Relay) Start(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.runOnce(ctx); err != nil {
				r.logger.Errorf("outbox relay: poll failed: %v", err)
			}
		}
	}
}

// runOnce publishes one batch of unpublished events and returns how many
// were published. Exported for tests; safe to call concurrently with Start
// only if the caller doesn't also run Start (it's not re-entrant-safe
// across processes, which is fine for this single-relay-per-service setup).
func (r *Relay) runOnce(ctx context.Context) error {
	var pending []OutboxEvent
	if err := r.db.WithContext(ctx).
		Where("published_at IS NULL").
		Order("created_at ASC").
		Limit(r.batchSize).
		Find(&pending).Error; err != nil {
		return fmt.Errorf("failed to load pending outbox events: %w", err)
	}

	for _, row := range pending {
		if err := r.publishRow(row); err != nil {
			r.logger.Warnf("outbox relay: publish failed for event %s (routing key %s): %v. Will retry.", row.EventID, row.RoutingKey, err)
			continue
		}

		now := time.Now()
		if err := r.db.WithContext(ctx).Model(&OutboxEvent{}).
			Where("id = ?", row.ID).
			Update("published_at", now).Error; err != nil {
			r.logger.Errorf("outbox relay: published event %s but failed to mark it published: %v", row.EventID, err)
		}
	}

	return nil
}

func (r *Relay) publishRow(row OutboxEvent) error {
	return r.publisher.publisher.Publish(
		row.Payload,
		[]string{row.RoutingKey},
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsHeaders(rabbitmq.Table{
			"event_id":       row.EventID.String(),
			"schema_version": row.SchemaVersion,
		}),
	)
}
