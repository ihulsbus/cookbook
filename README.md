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
- **Authentication**: OAuth2/OIDC integration with Keycloak

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
- **cache**: Redis-based caching layer
- **healthchecks**: Health check utilities
- **http**: Common HTTP utilities and middleware
- **httpclient**: Enhanced HTTP client with retry and circuit breaker patterns
- **keycloak**: Keycloak authentication integration
- **rabbitmq**: RabbitMQ publisher/consumer utilities
- **recipeclient**: Client library for inter-service communication
- **imageclient**: Client library for image service integration

## Technology Stack

- **Language**: Go 1.26+
- **Database**: PostgreSQL 15
- **Message Queue**: RabbitMQ 3
- **Object Storage**: S3
- **Authentication**: Keycloak / OAuth2 OIDC
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
# - Keycloak: localhost:8081
# - Traefik Dashboard: localhost:8080
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
│   ├── keycloak/          # Authentication
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
- `keycloak` - Keycloak authentication database

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

