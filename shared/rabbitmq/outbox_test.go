package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/wagslane/go-rabbitmq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func newMockDatabase(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	sqlMockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sql mock init failed: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlMockDB,
		PreferSimpleProtocol: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{Logger: newLogger})
	if err != nil {
		t.Fatalf("gorm mock init failed: %v", err)
	}

	return db, mock
}

type relayLoggerMock struct{}

func (relayLoggerMock) Debugf(format string, args ...interface{}) {}
func (relayLoggerMock) Warnf(format string, args ...interface{})  {}
func (relayLoggerMock) Errorf(format string, args ...interface{}) {}

func TestNewOutboxEntry_Success(t *testing.T) {
	payload := RecipePayload{ID: uuid.New(), Name: "apple pie"}

	entry, err := NewOutboxEntry(Event{
		Name:       "RecipeCreated",
		RoutingKey: RecipeCreatedRoutingKey,
		Payload:    payload,
	})

	assert.NoError(t, err)
	assert.Equal(t, RecipeCreatedRoutingKey, entry.RoutingKey)
	assert.NotEqual(t, uuid.Nil, entry.EventID)
	assert.Equal(t, CurrentSchemaVersion, entry.SchemaVersion)

	var decoded RecipePayload
	assert.NoError(t, json.Unmarshal(entry.Payload, &decoded))
	assert.Equal(t, payload.ID, decoded.ID)
}

func TestNewOutboxEntry_KeepsExplicitEventIDAndSchemaVersion(t *testing.T) {
	eventID := uuid.New()

	entry, err := NewOutboxEntry(Event{
		RoutingKey:    RecipeUpdatedRoutingKey,
		Payload:       RecipePayload{},
		EventID:       eventID,
		SchemaVersion: 7,
	})

	assert.NoError(t, err)
	assert.Equal(t, eventID, entry.EventID)
	assert.Equal(t, 7, entry.SchemaVersion)
}

func TestNewOutboxEntry_MarshalError(t *testing.T) {
	_, err := NewOutboxEntry(Event{
		RoutingKey: RecipeCreatedRoutingKey,
		Payload:    make(chan int),
	})

	assert.Error(t, err)
}

// MockRelayPublisher implements PublisherInterface so Relay.runOnce can be
// tested without a real RabbitMQ connection.
type MockRelayPublisher struct {
	published []string // routing keys it was called with
	failOn    string   // routing key to fail on, if any
}

func (m *MockRelayPublisher) Publish(body []byte, routingKeys []string, options ...func(*rabbitmq.PublishOptions)) error {
	for _, k := range routingKeys {
		if k == m.failOn {
			return errors.New("publish failed")
		}
	}
	m.published = append(m.published, routingKeys...)
	return nil
}

func (m *MockRelayPublisher) Close() {}

func TestRelay_RunOnce_PublishesAndMarksPublished(t *testing.T) {
	db, mock := newMockDatabase(t)
	mockPub := &MockRelayPublisher{}
	publisher := &Publisher{publisher: mockPub}
	relay := NewRelay(db, publisher, relayLoggerMock{}, time.Second, 10)

	rowID := uuid.New()
	eventID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "outbox_events" WHERE published_at IS NULL ORDER BY created_at ASC LIMIT $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "routing_key", "payload", "event_id", "schema_version", "created_at", "published_at"}).
			AddRow(rowID, RecipeCreatedRoutingKey, []byte(`{}`), eventID, 1, time.Now(), nil))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "outbox_events" SET "published_at"=$1 WHERE id = $2`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := relay.runOnce(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, []string{RecipeCreatedRoutingKey}, mockPub.published)
}

func TestRelay_RunOnce_LeavesUnpublishedOnPublishFailure(t *testing.T) {
	db, mock := newMockDatabase(t)
	mockPub := &MockRelayPublisher{failOn: RecipeCreatedRoutingKey}
	publisher := &Publisher{publisher: mockPub}
	relay := NewRelay(db, publisher, relayLoggerMock{}, time.Second, 10)

	rowID := uuid.New()
	eventID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "outbox_events" WHERE published_at IS NULL ORDER BY created_at ASC LIMIT $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "routing_key", "payload", "event_id", "schema_version", "created_at", "published_at"}).
			AddRow(rowID, RecipeCreatedRoutingKey, []byte(`{}`), eventID, 1, time.Now(), nil))

	// No UPDATE expected: the publish failed, so the row must stay
	// unpublished for the next poll to retry.
	err := relay.runOnce(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, mockPub.published)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRelay_RunOnce_NoPendingEvents(t *testing.T) {
	db, mock := newMockDatabase(t)
	mockPub := &MockRelayPublisher{}
	publisher := &Publisher{publisher: mockPub}
	relay := NewRelay(db, publisher, relayLoggerMock{}, time.Second, 10)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "outbox_events" WHERE published_at IS NULL ORDER BY created_at ASC LIMIT $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "routing_key", "payload", "event_id", "schema_version", "created_at", "published_at"}))

	err := relay.runOnce(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, mockPub.published)
}
