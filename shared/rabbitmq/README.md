# shared/rabbitmq

Event publishing and consumption for Cookbook's services, plus a
transactional outbox and relay for reliable delivery.

## Event-carried state transfer (ECST)

Recipe events are **ECST**: each one carries the recipe's full current
state, not a delta. A consumer that only cares about part of that state
(e.g. `shared/authz`, which only needs ownership/visibility/shares) can
ignore the rest; a consumer that missed an earlier event, or receives one
out of order, can still end up correct by applying whichever event it
sees with the highest `Version` - there's nothing to "catch up" on, because
every event is already the whole picture.

### Routing keys

| Routing key              | Payload               | Published when |
|---------------------------|-----------------------|-----------------|
| `recipe.created`          | `RecipePayload`       | A recipe is created |
| `recipe.updated`          | `RecipePayload`       | A recipe's content, owner, or visibility changes |
| `recipe.deleted`          | `RecipePayload`       | A recipe is deleted (`Deleted: true`) |
| `recipe.access.changed`   | `RecipeAccessState`   | Reserved for a future endpoint that changes sharing/visibility without a full recipe update (e.g. a share-management API) |
| `image.updated`           | `ImagePayload`        | An image is created or replaced |
| `image.deleted`           | `ImagePayload`        | An image is deleted |

### Payload shape

`RecipePayload` carries both the recipe's content (name, description,
serving count) and its full access state:

```go
type RecipePayload struct {
    ID           uuid.UUID
    Name         string
    Description  string
    ServingCount int

    OwnerID    string       // OIDC sub of the owner
    Visibility string       // "private" | "shared" | "public"
    Shares     []ShareEntry
    Version    int64        // monotonically increasing per recipe
    Deleted    bool
}

type ShareEntry struct {
    TargetType string // "user" | "group"
    TargetID   string
    Permission string // "read" | "write"
}
```

`RecipeAccessState` is the same access-state subset, for
`recipe.access.changed` and for consumers (like `shared/authz`) that only
need it - `RecipePayload.ToAccessState()` extracts it from a content event.

### Idempotent, out-of-order consumption

Every consumer that projects this state should apply the rule: **drop the
event when its `Version` is not strictly greater than what's already
stored.** This makes redelivery (at-least-once, from the outbox relay) and
out-of-order delivery both safe without any other coordination.
`shared/authz`'s `Projector.Handle` implements exactly this.

### Transport metadata

`event_id` (a UUID) and `schema_version` (currently `1`, see
`CurrentSchemaVersion`) are set as AMQP message headers by `Publish`, not
in the JSON body - they're for tracing/correlation and future schema
evolution, kept separate from the business payload.

## Reliable publishing: the outbox

A direct `Publisher.Publish` call can fail (RabbitMQ down, network
partition) and lose the event - previously the caller just logged it.
`OutboxEvent` plus `Relay` fix that:

1. The caller writes an `OutboxEvent` row (via `NewOutboxEntry`) in the
   **same DB transaction** as the business change it describes, so the
   event is never recorded unless the change it describes actually
   committed.
2. A `Relay` polls for `OutboxEvent` rows where `published_at IS NULL`,
   publishes each, and marks it published on success.

This is at-least-once delivery: a crash between a successful publish and
marking the row published will re-publish that event on the next poll.
That's fine because every consumer is expected to be idempotent via
`Version`, per the rule above.

```go
publisher, _ := rabbitmq.NewPublisher(conn, "cookbook")
relay := rabbitmq.NewRelay(db, publisher, logger, 2*time.Second, 50)
go relay.Start(ctx)

// ... inside the same DB transaction as creating a Recipe:
outbox, _ := rabbitmq.NewOutboxEntry(rabbitmq.Event{
    RoutingKey: rabbitmq.RecipeCreatedRoutingKey,
    Payload:    payload, // a RecipePayload
})
tx.Create(&recipe)
tx.Create(&outbox)
```

Call `rabbitmq.AutoMigrateOutbox(db)` once at startup, alongside the
service's own model migrations.
