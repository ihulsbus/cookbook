# Integration Example: Adding Circuit Breaker to Metadata Service

This guide shows how to integrate the enhanced HTTP client with circuit breaker and retry logic into the metadata-service.

## Step 1: Update Config Structure

Add retry and circuit breaker configuration to your config model:

```go
// shared/models/ConfigModels.go

type Config struct {
    Global   GlobalConfig   `mapstructure:"global"`
    Database DatabaseConfig `mapstructure:"database"`
    OAuth    OauthConfig    `mapstructure:"oauth"`
    // Add HTTP client configuration
    HTTPClient HTTPClientConfig `mapstructure:"httpclient"`
}

type HTTPClientConfig struct {
    Timeout        int                   `mapstructure:"timeout"`        // seconds
    Retry          RetryConfig           `mapstructure:"retry"`
    CircuitBreaker CircuitBreakerConfig  `mapstructure:"circuitbreaker"`
}

type RetryConfig struct {
    MaxRetries     int `mapstructure:"maxretries"`
    InitialBackoff int `mapstructure:"initialbackoff"` // milliseconds
    MaxBackoff     int `mapstructure:"maxbackoff"`     // milliseconds
    BackoffFactor  float64 `mapstructure:"backofffactor"`
}

type CircuitBreakerConfig struct {
    MaxRequests        uint32 `mapstructure:"maxrequests"`
    Interval           int    `mapstructure:"interval"`           // seconds
    Timeout            int    `mapstructure:"timeout"`            // seconds
    FailureThreshold   int    `mapstructure:"failurethreshold"`   // min requests before opening
    FailureRatePercent int    `mapstructure:"failureratepercent"` // percent (0-100)
}
```

## Step 2: Update config.yaml

Add HTTP client configuration with sensible defaults:

```yaml
# metadata-service/cmd/metadata-service/config.yaml

global:
  logLevel: INFO

httpclient:
  timeout: 5  # 5 seconds
  retry:
    maxretries: 3
    initialbackoff: 100  # 100ms
    maxbackoff: 5000     # 5s
    backofffactor: 2.0
  circuitbreaker:
    maxrequests: 3
    interval: 60           # 60 seconds
    timeout: 30            # 30 seconds
    failurethreshold: 10   # min 10 requests
    failureratepercent: 50 # 50% failure rate

oauth:
  clientid: metadata-service
  clientsecret: ${OAUTH_CLIENTSECRET}
  url: https://auth.example.com
  realm: production

database:
  # ... existing config
```

## Step 3: Initialize Enhanced HTTP Client

Update your config initialization to create the enhanced client:

```go
// internal/config/config.go or initFuncs.go

import (
    "time"
    "github.com/ihulsbus/cookbook/shared/httpclient"
    "github.com/sony/gobreaker"
)

var (
    Configuration m.Config
    Logger        *log.Logger
    DatabaseClient *gorm.DB

    // Add HTTP client
    HTTPClient     *httpclient.EnhancedHTTPClient
    RecipeClient   *recipeclient.RecipeAPIClient
)

func initHTTPClient() {
    cfg := Configuration.HTTPClient

    // Build retry config from settings
    retryConfig := &httpclient.RetryConfig{
        MaxRetries:     cfg.Retry.MaxRetries,
        InitialBackoff: time.Duration(cfg.Retry.InitialBackoff) * time.Millisecond,
        MaxBackoff:     time.Duration(cfg.Retry.MaxBackoff) * time.Millisecond,
        BackoffFactor:  cfg.Retry.BackoffFactor,
    }

    // Build circuit breaker config from settings
    cbConfig := &httpclient.CircuitBreakerConfig{
        MaxRequests: cfg.CircuitBreaker.MaxRequests,
        Interval:    time.Duration(cfg.CircuitBreaker.Interval) * time.Second,
        Timeout:     time.Duration(cfg.CircuitBreaker.Timeout) * time.Second,
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
            failurePercent := failureRatio * 100

            minRequests := uint32(cfg.CircuitBreaker.FailureThreshold)
            thresholdPercent := float64(cfg.CircuitBreaker.FailureRatePercent)

            return counts.Requests >= minRequests && failurePercent >= thresholdPercent
        },
    }

    // Create enhanced HTTP client
    HTTPClient = httpclient.NewEnhancedHTTPClient(
        time.Duration(cfg.Timeout)*time.Second,
        Configuration.OAuth.Url,
        Configuration.OAuth.Realm,
        Configuration.OAuth.ClientID,
        Configuration.OAuth.ClientSecret,
        Logger,
        retryConfig,
        cbConfig,
    )

    Logger.Info("Enhanced HTTP client initialized with circuit breaker and retry logic")
}

func initServiceClients() {
    var err error

    // Create recipe client with enhanced HTTP client
    RecipeClient, err = recipeclient.NewRecipeAPIClient(
        "http://recipe-service:8080/api/v2",
        HTTPClient,
    )
    if err != nil {
        Logger.Fatalf("Failed to initialize recipe client: %v", err)
    }

    Logger.Info("Service clients initialized")
}

func init() {
    initViper()
    initConfig()
    initLogging()
    initDatabase()
    initHTTPClient()     // Add this
    initServiceClients() // Add this
    // ... rest of init
}
```

## Step 4: Use Clients in Your Services

Now use the clients in your service layer:

```go
// internal/services/metadata/MetadataService.go

type MetadataService struct {
    repo         MetadataRepository
    recipeClient *recipeclient.RecipeAPIClient
    logger       LoggerInterface
}

func NewMetadataService(
    repo MetadataRepository,
    recipeClient *recipeclient.RecipeAPIClient,
    logger LoggerInterface,
) *MetadataService {
    return &MetadataService{
        repo:         repo,
        recipeClient: recipeClient,
        logger:       logger,
    }
}

func (s *MetadataService) Create(recipeID uuid.UUID, meta *m.RecipeMetadataDTO) (*m.RecipeMetadataDTO, error) {
    // Verify recipe exists before creating metadata
    // This call automatically uses retry and circuit breaker!
    exists, err := s.recipeClient.RecipeExists(recipeID.String())
    if err != nil {
        s.logger.Errorf("Failed to verify recipe exists: %v", err)
        return nil, fmt.Errorf("unable to verify recipe: %w", err)
    }

    if !exists {
        return nil, fmt.Errorf("recipe %s not found", recipeID)
    }

    // Create metadata
    return s.repo.Create(recipeID, meta)
}
```

## Step 5: Update Config Initialization

Pass the recipe client to services:

```go
// internal/config/config.go

func init() {
    // ... previous init code

    initHTTPClient()
    initServiceClients()

    // Init services with recipe client
    MetadataService = services.NewMetadataService(
        MetadataRepository,
        RecipeClient,  // Pass the client here
        Logger,
    )

    // Init handlers
    MetadataHandlers = handlers.NewMetadataHandlers(MetadataService, Logger)
}
```

## Step 6: Add Health Check for Circuit Breaker

Add circuit breaker status to health checks:

```go
// internal/handlers/health/HealthHandlers.go

func (h *HealthHandlers) Readiness(ctx *gin.Context) {
    checks := make(map[string]string)
    overallStatus := "UP"

    // Check database
    if err := h.checkDatabase(); err != nil {
        checks["database"] = "DOWN: " + err.Error()
        overallStatus = "DOWN"
    } else {
        checks["database"] = "UP"
    }

    // Check circuit breaker state
    cbState := h.httpClient.GetCircuitBreakerState()
    checks["circuit_breaker"] = cbState.String()

    if cbState == gobreaker.StateOpen {
        checks["recipe_service"] = "DOWN: circuit breaker open"
        overallStatus = "DOWN"
        h.logger.Warnf("Recipe service circuit breaker is open")
    } else {
        checks["recipe_service"] = "UP"
    }

    // ... rest of readiness check
}
```

## Step 7: Add Monitoring Metrics (Optional)

Export circuit breaker metrics:

```go
// internal/metrics/circuitbreaker.go

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/sony/gobreaker"
)

var (
    circuitBreakerState = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "circuit_breaker_state",
            Help: "Circuit breaker state (0=closed, 1=open, 2=half-open)",
        },
        []string{"service"},
    )

    circuitBreakerRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "circuit_breaker_requests_total",
            Help: "Total requests through circuit breaker",
        },
        []string{"service", "result"},
    )
)

func init() {
    prometheus.MustRegister(circuitBreakerState)
    prometheus.MustRegister(circuitBreakerRequests)
}

func UpdateCircuitBreakerMetrics(serviceName string, cb *gobreaker.CircuitBreaker) {
    // Update state
    state := cb.State()
    var stateValue float64
    switch state {
    case gobreaker.StateClosed:
        stateValue = 0
    case gobreaker.StateOpen:
        stateValue = 1
    case gobreaker.StateHalfOpen:
        stateValue = 2
    }
    circuitBreakerState.WithLabelValues(serviceName).Set(stateValue)

    // Update counters
    counts := cb.Counts()
    circuitBreakerRequests.WithLabelValues(serviceName, "success").Add(float64(counts.TotalSuccesses))
    circuitBreakerRequests.WithLabelValues(serviceName, "failure").Add(float64(counts.TotalFailures))
}
```

## Step 8: Test the Integration

### Test Retry Logic

Stop the recipe-service and make a request:

```bash
# Terminal 1: Stop recipe service
docker stop recipe-service

# Terminal 2: Watch logs
kubectl logs -f metadata-service-xxx

# Terminal 3: Make request
curl -X POST http://metadata-service/api/v2/metadata/some-recipe-id

# You should see retry attempts in logs:
# "Request failed (attempt 1/4): connection refused. Retrying in 100ms..."
# "Request failed (attempt 2/4): connection refused. Retrying in 200ms..."
# etc.
```

### Test Circuit Breaker

Make multiple failing requests to open the circuit:

```bash
# Make 15 requests that will fail
for i in {1..15}; do
  curl http://metadata-service/api/v2/metadata/test-recipe-$i
  sleep 0.5
done

# Check health endpoint
curl http://metadata-service/health/ready

# Should show:
# {
#   "status": "DOWN",
#   "checks": {
#     "circuit_breaker": "open",
#     "recipe_service": "DOWN: circuit breaker open"
#   }
# }
```

## Complete Example: metadata-service

Here's the complete file structure after integration:

```
metadata-service/
├── cmd/
│   └── metadata-service/
│       ├── main.go
│       └── config.yaml  ← Updated with httpclient config
├── internal/
│   ├── config/
│   │   ├── config.go    ← Added HTTPClient, RecipeClient vars
│   │   └── initFuncs.go ← Added initHTTPClient(), initServiceClients()
│   ├── models/
│   │   └── config.go    ← Added HTTPClientConfig struct
│   ├── services/
│   │   └── metadata/
│   │       └── MetadataService.go ← Uses RecipeClient
│   ├── handlers/
│   │   └── health/
│   │       └── HealthHandlers.go  ← Added circuit breaker check
│   └── metrics/
│       └── circuitbreaker.go      ← Optional metrics
└── go.mod               ← Includes github.com/sony/gobreaker
```

## Rolling Out to Other Services

Once you've tested in metadata-service, roll out to:
1. instruction-service (uses recipeclient + imageclient)
2. Any future services that make inter-service calls

The pattern is identical for each service!
