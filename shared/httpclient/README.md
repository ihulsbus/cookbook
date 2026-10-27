# Enhanced HTTP Client with Circuit Breaker and Retry Logic

This package provides a resilient HTTP client with circuit breaker and retry logic for inter-service communication.

## Features

* ✅ **Exponential Backoff Retry** - Automatically retries failed requests with configurable backoff 
* ✅ **Circuit Breaker Pattern** - Prevents cascading failures by opening the circuit after consecutive failures 
* ✅ **Automatic Token Management** - Handles Keycloak authentication token refresh
* ✅ **Request Cloning** - Properly handles request body for retries
* ✅ **Comprehensive Logging** - Detailed logs for debugging and monitoring
* ✅ **Thread-Safe** - Safe for concurrent use
* ✅ **Configurable** - Sensible defaults with full customization

## Quick Start

### Basic Usage (with defaults)

```go
import "github.com/ihulsbus/cookbook/shared/httpclient"

// Create client with default retry and circuit breaker settings
client := httpclient.NewEnhancedHTTPClient(
    5*time.Second,                    // Request timeout
    "https://auth.example.com",       // Keycloak URL
    "my-realm",                       // Keycloak realm
    "my-client-id",                   // Client ID
    "my-client-secret",               // Client secret
    logger,                           // Logger interface
    nil,                              // Use default retry config
    nil,                              // Use default circuit breaker config
)

// Make requests - retries and circuit breaker work automatically
req, _ := http.NewRequest("GET", "http://recipe-service/api/v2/recipe/123", nil)
resp, err := client.Do(req)
```

### Custom Configuration

```go
// Custom retry configuration
retryConfig := &httpclient.RetryConfig{
    MaxRetries:     5,                    // Retry up to 5 times
    InitialBackoff: 200 * time.Millisecond, // Start with 200ms delay
    MaxBackoff:     10 * time.Second,     // Max delay of 10s
    BackoffFactor:  2.0,                  // Double delay each retry
}

// Custom circuit breaker configuration
cbConfig := &httpclient.CircuitBreakerConfig{
    MaxRequests: 5,              // Allow 5 requests in half-open state
    Interval:    120 * time.Second, // Clear stats every 2 minutes
    Timeout:     60 * time.Second,  // Stay open for 1 minute before trying
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        // Open circuit if 60% of last 20 requests failed
        failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
        return counts.Requests >= 20 && failureRatio >= 0.6
    },
}

client := httpclient.NewEnhancedHTTPClient(
    10*time.Second,
    "https://auth.example.com",
    "my-realm",
    "my-client-id",
    "my-client-secret",
    logger,
    retryConfig,
    cbConfig,
)
```

## How It Works

### Retry Logic

The client automatically retries requests that fail with:
- Network errors (connection refused, timeout, etc.)
- Server errors (5xx status codes)

**It does NOT retry:**
- Client errors (4xx status codes)
- Successful responses (2xx, 3xx)

**Retry Strategy:**
1. First attempt fails
2. Wait `InitialBackoff` (default: 100ms)
3. Second attempt fails
4. Wait `InitialBackoff * BackoffFactor` (default: 200ms)
5. Third attempt fails
6. Wait with exponential increase up to `MaxBackoff`
7. Continue up to `MaxRetries` attempts

### Circuit Breaker States

```
          ┌─────────┐
          │ CLOSED  │ ◄──── Normal operation
          └────┬────┘
               │
          Failures exceed
          threshold
               │
               ▼
          ┌─────────┐
          │  OPEN   │ ◄──── Reject all requests immediately
          └────┬────┘
               │
          Timeout expires
               │
               ▼
        ┌──────────────┐
        │  HALF-OPEN   │ ◄──── Test if service recovered
        └──────┬───────┘
               │
        Success │  Failure
               │
     ┌─────────┴──────────┐
     ▼                    ▼
 ┌─────────┐         ┌─────────┐
 │ CLOSED  │         │  OPEN   │
 └─────────┘         └─────────┘
```

**CLOSED (Normal):**
- All requests pass through
- Failures are counted
- Opens when failure threshold reached

**OPEN (Failing):**
- All requests fail immediately with `ErrOpenState`
- No requests reach the downstream service
- Stays open for `Timeout` duration

**HALF-OPEN (Testing):**
- Limited requests (`MaxRequests`) allowed through
- If requests succeed → Circuit closes
- If requests fail → Circuit opens again

### Default Configuration

```go
// Default Retry Configuration
RetryConfig{
    MaxRetries:     3,                    // 4 total attempts
    InitialBackoff: 100 * time.Millisecond,
    MaxBackoff:     5 * time.Second,
    BackoffFactor:  2.0,                  // Exponential backoff
}

// Default Circuit Breaker Configuration
CircuitBreakerConfig{
    MaxRequests: 3,                // 3 test requests in half-open
    Interval:    60 * time.Second, // Reset counters every minute
    Timeout:     30 * time.Second, // Stay open for 30 seconds
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        // Open if ≥10 requests and ≥50% failure rate
        failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
        return counts.Requests >= 10 && failureRatio >= 0.5
    },
}
```

## Integration with Service Clients

### RecipeClient Example

```go
import (
    "github.com/ihulsbus/cookbook/shared/httpclient"
    "github.com/ihulsbus/cookbook/shared/recipeclient"
)

// Create enhanced HTTP client
httpClient := httpclient.NewEnhancedHTTPClient(
    5*time.Second,
    keycloakURL,
    keycloakRealm,
    clientID,
    clientSecret,
    logger,
    nil, // defaults
    nil, // defaults
)

// Create recipe client with enhanced HTTP client
recipeClient, err := recipeclient.NewRecipeAPIClient(
    "http://recipe-service:8080/api/v2",
    httpClient,
)

// Use it - retries and circuit breaker work automatically!
exists, err := recipeClient.RecipeExists("some-recipe-id")
```

### ImageClient Example

```go
import (
    "github.com/ihulsbus/cookbook/shared/httpclient"
    "github.com/ihulsbus/cookbook/shared/imageclient"
)

httpClient := httpclient.NewEnhancedHTTPClient(
    5*time.Second,
    keycloakURL,
    keycloakRealm,
    clientID,
    clientSecret,
    logger,
    nil,
    nil,
)

imageClient := imageclient.NewImageAPIClient(
    "http://image-service:8080",
    httpClient,
)

exists, err := imageClient.ImageExists("some-image-id")
```

## Monitoring

### Check Circuit Breaker State

```go
state := client.GetCircuitBreakerState()
switch state {
case gobreaker.StateClosed:
    log.Info("Circuit is closed - operating normally")
case gobreaker.StateOpen:
    log.Warn("Circuit is open - requests failing fast")
case gobreaker.StateHalfOpen:
    log.Info("Circuit is half-open - testing recovery")
}
```

### Get Circuit Breaker Statistics

```go
counts := client.GetCircuitBreakerCounts()
log.Infof("Circuit breaker stats: requests=%d, successes=%d, failures=%d, consecutive_failures=%d",
    counts.Requests,
    counts.TotalSuccesses,
    counts.TotalFailures,
    counts.ConsecutiveFailures,
)
```

### Metrics for Observability

Export these metrics to Prometheus/Grafana:

```go
// Circuit breaker state (0=closed, 1=open, 2=half-open)
circuitBreakerStateGauge.Set(float64(client.GetCircuitBreakerState()))

// Request counts
counts := client.GetCircuitBreakerCounts()
requestsTotal.Add(float64(counts.Requests))
successesTotal.Add(float64(counts.TotalSuccesses))
failuresTotal.Add(float64(counts.TotalFailures))
```

## Testing

Run tests:

```bash
cd shared/httpclient
go test -v
```

Test coverage:

```bash
go test -cover
```

## Migration Guide

### From Old HTTPClient to EnhancedHTTPClient

**Before:**
```go
client := httpclient.NewHTTPClient(
    5*time.Second,
    keycloakURL,
    keycloakRealm,
    clientID,
    clientSecret,
    logger,
)
```

**After:**
```go
client := httpclient.NewEnhancedHTTPClient(
    5*time.Second,
    keycloakURL,
    keycloakRealm,
    clientID,
    clientSecret,
    logger,
    nil, // Use defaults
    nil, // Use defaults
)
```

Both implement the same `Do(req *http.Request) (*http.Response, error)` interface, so existing code works without changes!

## When to Use Enhanced vs. Standard Client

### Use EnhancedHTTPClient for:
✅ Service-to-service communication
✅ External API calls
✅ Any HTTP call that might fail transiently
✅ Production environments

### Use standard HTTPClient (or net/http) for:
- Internal health checks (no retries needed)
- One-time setup requests
- Testing/development only

## Best Practices

1. **Configure per service** - Different services have different failure characteristics
2. **Monitor circuit breaker state** - Alert when circuits open frequently
3. **Log retry attempts** - Use Debug level to see retry behavior
4. **Set appropriate timeouts** - Request timeout should be less than sum of all retry attempts
5. **Test failure scenarios** - Simulate downstream failures to verify circuit breaker works
6. **Use metrics** - Export circuit breaker state and retry counts to monitoring

## Troubleshooting

### Issue: Too many retries slowing down requests

**Solution**: Reduce `MaxRetries` or increase `InitialBackoff`:

```go
retryConfig := &httpclient.RetryConfig{
    MaxRetries:     2,  // Reduce from default 3
    InitialBackoff: 50 * time.Millisecond,
    MaxBackoff:     2 * time.Second,
    BackoffFactor:  2.0,
}
```

### Issue: Circuit breaker opens too easily

**Solution**: Increase failure threshold in `ReadyToTrip`:

```go
cbConfig := &httpclient.CircuitBreakerConfig{
    // ...
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        // Require 70% failure rate instead of 50%
        failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
        return counts.Requests >= 20 && failureRatio >= 0.7
    },
}
```

### Issue: Circuit stays open too long

**Solution**: Reduce `Timeout`:

```go
cbConfig := &httpclient.CircuitBreakerConfig{
    Timeout: 10 * time.Second, // Reduce from default 30s
    // ...
}
```

## References

- [Circuit Breaker Pattern (Martin Fowler)](https://martinfowler.com/bliki/CircuitBreaker.html)
- [sony/gobreaker Documentation](https://github.com/sony/gobreaker)
- [Retry Pattern (Microsoft)](https://docs.microsoft.com/en-us/azure/architecture/patterns/retry)
