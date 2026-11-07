# Cookbook

[![CodeQL](https://github.com/ihulsbus/cookbook/actions/workflows/codeql-analysis.yml/badge.svg?branch=master)](https://github.com/ihulsbus/cookbook/actions/workflows/codeql-analysis.yml)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=ihulsbus_cookbook&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=ihulsbus_cookbook)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ihulsbus_cookbook&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ihulsbus_cookbook)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=ihulsbus_cookbook&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=ihulsbus_cookbook)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ihulsbus_cookbook&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ihulsbus_cookbook)

A modern microservices-based backend API for the Gourmedy recipe management application. Built with Go and designed for scalability, maintainability, and ease of deployment.

## Overview

Cookbook is a comprehensive recipe management system that provides APIs for creating, storing, and managing recipes, ingredients, cooking instructions, and metadata. The application follows a microservices architecture with event-driven communication between services.

### Key Features

- **Recipe Management**: Create, read, update, and delete recipes with full metadata
- **Ingredient Tracking**: Manage ingredients, units, and amounts
- **Instruction Management**: Step-by-step cooking instructions
- **Image Handling**: Upload and manage recipe images with S3-compatible storage
- **Full-Text Search**: Search across recipes, ingredients, and instructions
- **Metadata Management**: Categories, cuisine types, difficulty levels, tags, and preparation times
- **Pagination Support**: All list endpoints support pagination with customizable page size
- **Event-Driven Architecture**: Services communicate via RabbitMQ for loose coupling
- **Authentication**: OAuth2/OIDC integration with Authentik

## Architecture

The project consists of the following microservices:

### Core Services

- **recipe-service**: Manages recipe entities and lifecycle
- **ingredient-service**: Handles ingredients, units, and amounts
- **instruction-service**: Manages cooking instructions for recipes
- **metadata-service**: Manages recipe metadata (categories, tags, cuisine types, difficulty levels, preparation times)
- **image-service**: Handles image uploads and storage with S3-compatible backends
- **search-service**: Provides full-text search capabilities across all entities
- **notification-service**: Handles notifications and events (in development)

### Shared Components

The `shared/` directory contains common libraries used across all services:

- **models**: Shared data models and DTOs with conversion methods
- **authz**: Recipe access control - a local `recipe_access` projection kept
  current by consuming recipe-service's ECST events, plus the
  `Authorizer`/`RequireRecipeAccess` middleware each service uses to check it
- **healthchecks**: Health check utilities
- **http**: Common HTTP utilities and middleware
- **httpclient**: Enhanced HTTP client with retry and circuit breaker patterns
- **oidc**: Provider-agnostic OIDC authentication (validates tokens against Authentik)
- **rabbitmq**: RabbitMQ publisher/consumer utilities, plus a transactional
  outbox and relay for reliable event publishing
- **imageclient**: Client library for image service integration

### No service calls another service at startup

A service's startup must only depend on its own database and its own
RabbitMQ connection (with retry/backoff) - never on another service being
up and reachable. Cross-service data comes from one of two places:

- **Events**: a service keeps a local projection (e.g. `shared/authz`'s
  `recipe_access` table) current by consuming another service's events in
  the background, starting from whatever state it finds at boot.
- **Lazy, circuit-broken calls at request time**: if a projection is
  missing a row (e.g. the read landed right after the row's owner was
  created, before its event was consumed), fetch it on demand - never at
  startup, and always through a client with retry/circuit-breaker
  (`shared/httpclient`), never a raw, unretried call.

A readiness probe may report degraded while a dependency is unavailable;
liveness must not fail because of it. This replaced an earlier pattern
(KIT-79) where ingredient/instruction/metadata/image-service each called
recipe-service synchronously during startup to warm a cache, which meant
all four went into a restart loop whenever recipe-service was down -
backwards for services that don't otherwise depend on it being up.

### Runbook: rebuilding a recipe-access projection (KIT-31)

Each consumer of `shared/authz` (ingredient/instruction/metadata/image-
service) keeps its own `recipe_access` projection current by consuming
recipe-service's `recipe.*`/`recipe.access.changed` events as they
happen. That's enough for normal operation, but it only covers events
published *after* a consumer started listening - it does nothing for a
consumer that's missing history, for example:

- a **new consumer** (a service added after recipe-service already had
  recipes), whose projection table starts empty;
- a consumer **recovering from data loss** (its `recipe_access` table was
  dropped, truncated, or restored from a backup older than recipe-
  service's);
- a consumer that's known to have **missed events** (e.g. it was down
  past RabbitMQ's retention, or an outbox bug dropped rows before this
  fix).

`recipe-service --republish-access` fixes this by replaying the current
access state of every recipe it has - including soft-deleted ones, so a
consumer that missed the original delete still learns about it - as a
`recipe.access.changed` event. It's a one-shot CLI flag, not a long-
running mode: it writes one outbox row per recipe and exits.

**When to run it**: after bringing up a new consumer, or restoring a
consumer's `recipe_access` table, or after fixing an outage that's known
to have caused missed events. Not part of routine deployment - normal
event consumption is sufficient otherwise.

It's the same binary as the live service, so it needs the same
environment/config as a normal recipe-service start (database, RabbitMQ,
and OIDC all get initialised before the flag is even checked) - run it
with recipe-service's own config for the target environment, not a
stripped-down one.

**How to run it**:

```bash
# From recipe-service's own directory, pointed at the target environment's
# config (same config.yaml / env vars the real service would use):
go run ./cmd/recipe-service --republish-access
```

or, against a running container image:

```bash
docker run --rm <recipe-service-image> --republish-access
```

It connects to recipe-service's own database directly (the same
`DatabaseClient` the live service uses) and only inserts outbox rows -
it does **not** publish to RabbitMQ itself and does not need recipe-
service's RabbitMQ consumer or HTTP servers running in the same process.
The already-running recipe-service instance's outbox Relay drains those
rows in the background; if no instance of recipe-service is currently
running, start (or leave running) one afterwards so the queued rows
actually get published.

It prints a one-line summary when done:

```
republish-access: <total> recipes found, <queued> queued, <failed> failed
```

If `failed` is nonzero, check the logs for which recipe IDs failed (e.g.
a marshalling error) and re-run - re-running is safe and idempotent on
the consumer side (see below), though note each run inserts fresh outbox
rows rather than deduplicating against a prior run's.

**How to verify it worked**: compare row counts between recipe-service's
own recipe count and the consumer's projection table, including soft-
deleted rows on both sides:

```sql
-- On recipe-service's database:
SELECT count(*) FROM recipes;

-- On the consumer's database (e.g. ingredient-service):
SELECT count(*) FROM recipe_accesses;
```

The two counts should match once the consumer has drained the queued
events (check `rabbitmqctl list_queues` or the RabbitMQ management UI for
the consumer's queue depth returning to zero). A replay is safe to run
against a consumer that's already up to date: each event's `Version` is
the recipe's own historical `UpdatedAt`/`DeletedAt`, not "now", so a
consumer already at that version correctly drops the replayed row as
stale instead of re-applying it - row counts (and the data itself) are
unaffected either way.

## Technology Stack

- **Language**: Go 1.26+
- **Database**: PostgreSQL 15
- **Message Queue**: RabbitMQ 3
- **Object Storage**: S3
- **Authentication**: Authentik / OAuth2 OIDC
- **Reverse Proxy**: Traefik
- **API Documentation**: OpenAPI 3.1

## Getting Started

### Prerequisites

- Go 1.26 or later
- Docker and Docker Compose
- PostgreSQL 15 (or use Docker Compose)
- MinIO or S3-compatible storage (or use Docker Compose)
- RabbitMQ (or use Docker Compose)

### Quick Start with Docker Compose

The easiest way to run the entire stack locally is using Docker Compose:

```bash
# Clone the repository
git clone https://github.com/ihulsbus/cookbook.git
cd cookbook

# Start all infrastructure services
docker-compose up -d

# The following services will be available:
# - PostgreSQL: localhost:5432
# - MinIO: localhost:9000 (API), localhost:9001 (Console)
# - RabbitMQ: localhost:5672 (AMQP), localhost:15672 (Management UI)
# - Traefik Dashboard: localhost:8080
#
# Authentication is handled by Authentik, run separately (not part of this
# compose file) - see authentik/README.md for the OIDC provider setup.
```

### Running Services Locally

Each service can be run independently for development:

```bash
# Navigate to a service directory
cd recipe-service

# Install dependencies
go mod download

# Run the service
go run cmd/main.go
```

### Environment Variables

Each service requires the following environment variables. See individual service directories for specific configuration needs.

#### Database Configuration

| Variable | Type | Example | Description |
|----------|------|---------|-------------|
| `cbb_database_host` | string | `localhost` | PostgreSQL host |
| `cbb_database_port` | int | `5432` | PostgreSQL port |
| `cbb_database_database` | string | `recipes` | Database name |
| `cbb_database_username` | string | `admin` | Database user |
| `cbb_database_password` | string | `password` | Database password |
| `cbb_database_sslmode` | string | `disable` | SSL mode (disable/require) |
| `cbb_database_timezone` | string | `Europe/Amsterdam` | Database timezone |

#### S3/Object Storage Configuration (image-service)

| Variable | Type | Example | Description |
|----------|------|---------|-------------|
| `cbb_s3_endpoint` | string | `http://localhost:9000` | S3 endpoint URL |
| `cbb_s3_key` | string | `admin` | S3 access key |
| `cbb_s3_secret` | string | `adminpassword` | S3 secret key |
| `cbb_s3_bucket` | string | `cbhbe` | S3 bucket name |

#### RabbitMQ Configuration

| Variable | Type | Example | Description |
|----------|------|---------|-------------|
| `cbb_rabbitmq_host` | string | `localhost` | RabbitMQ host |
| `cbb_rabbitmq_port` | int | `5672` | RabbitMQ port |
| `cbb_rabbitmq_username` | string | `admin` | RabbitMQ user |
| `cbb_rabbitmq_password` | string | `password` | RabbitMQ password |

#### Authentication Configuration

| Variable | Type | Example | Description |
|----------|------|---------|-------------|
| `cbb_auth0_domain` | string | `https://example.eu.auth0.com/` | OIDC provider domain |
| `cbb_auth0_clientid` | string | `client_id` | OAuth2 client ID |
| `cbb_auth0_audience` | string | `api_audience` | OAuth2 audience |

#### General Configuration

| Variable | Type | Example | Description |
|----------|------|---------|-------------|
| `cbb_debug` | bool | `true` | Enable debug logging |

## API Documentation

The API is documented using OpenAPI 3.1 specification. The complete API documentation is available in `Cookbook.yaml`.

### Key Endpoints

#### Recipes
- `GET /api/v2/recipe` - List all recipes (paginated)
- `GET /api/v2/recipe/{id}` - Get a specific recipe
- `POST /api/v2/recipe` - Create a new recipe
- `PUT /api/v2/recipe/{id}` - Update a recipe
- `DELETE /api/v2/recipe/{id}` - Delete a recipe

#### Ingredients
- `GET /api/v2/ingredient` - List all ingredients (paginated)
- `GET /api/v2/ingredient/{id}` - Get a specific ingredient
- `POST /api/v2/ingredient` - Create a new ingredient
- `PUT /api/v2/ingredient/{id}` - Update an ingredient
- `DELETE /api/v2/ingredient/{id}` - Delete an ingredient

#### Search
- `GET /api/v2/search/recipe?q={query}` - Search recipes
- `GET /api/v2/search/ingredient?q={query}` - Search ingredients
- `GET /api/v2/search/instruction?q={query}` - Search instructions

All list endpoints support pagination via `?page={page}&limit={limit}` query parameters.

## Development

### Project Structure

```
cookbook/
├── recipe-service/          # Recipe management service
├── ingredient-service/      # Ingredient management service
├── instruction-service/     # Instruction management service
├── metadata-service/        # Metadata management service
├── image-service/          # Image storage service
├── search-service/         # Search service
├── notification-service/   # Notification service (in development)
├── shared/                 # Shared libraries and utilities
│   ├── models/            # Common data models
│   ├── cache/             # Caching utilities
│   ├── healthchecks/      # Health check utilities
│   ├── http/              # HTTP utilities
│   ├── httpclient/        # Enhanced HTTP client
│   ├── oidc/              # Authentication (Authentik)
│   ├── rabbitmq/          # Message queue utilities
│   ├── recipeclient/      # Inter-service client
│   └── imageclient/       # Image service client
├── docker-compose.yml      # Docker Compose configuration
├── Cookbook.yaml          # OpenAPI specification
├── go.work                # Go workspace configuration
└── init-db.sh             # Database initialization script
```

### Running Tests

Each service contains comprehensive unit and integration tests:

```bash
# Run all tests for a specific service
cd recipe-service
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests for all services (from root)
go test ./...
```

### Database Setup

The project uses PostgreSQL with separate databases for each service. The `init-db.sh` script automatically creates all required databases:

- `recipes` - Recipe service database
- `ingredients` - Ingredient service database
- `instructions` - Instruction service database
- `metadata` - Metadata service database
- `images` - Image service database
- `search` - Search service database

### Code Quality

The project includes:
- Static analysis with SonarCloud
- Security scanning with CodeQL
- Linting with golangci-lint
- Comprehensive test coverage
- API specification validation

## Deployment

### Docker

Each service includes a `Dockerfile` for containerized deployment:

```bash
# Build a service image
cd recipe-service
docker build -t cookbook/recipe-service:latest .

# Run the service
docker run -d \
  --name recipe-service \
  -e cbb_database_host=postgres \
  -e cbb_database_port=5432 \
  -e cbb_database_database=recipes \
  -e cbb_database_username=admin \
  -e cbb_database_password=password \
  cookbook/recipe-service:latest
```

### Kubernetes

Kubernetes manifests are available in the `k8s/` directory for each service. See `DOCKER_SETUP.md` and `WORKSPACE_SETUP.md` for detailed deployment instructions.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Coding Standards

- Follow Go best practices and idioms
- Write comprehensive tests for all new features
- Update documentation for API changes
- Ensure all tests pass before submitting PR
- Follow the existing code structure and patterns

## License

This project is licensed under the AGPL-3.0 License - see the [LICENSE](LICENSE) file for details.

