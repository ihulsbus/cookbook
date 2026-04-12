# Health Check Implementation

This service implements comprehensive health checks following Kubernetes best practices.

## Available Endpoints

### 1. Liveness Probe: `/health/live`
**Purpose**: Checks if the application is alive and running.

**When to use**: Kubernetes uses this to determine if the pod should be restarted.

**Response**:
```json
{
  "status": "UP",
  "timestamp": "2025-02-14T20:00:00Z",
  "description": "Service is alive"
}
```

**Status Codes**:
- `200 OK`: Service is alive

**Failure Action**: Kubernetes restarts the pod after 3 consecutive failures.

---

### 2. Readiness Probe: `/health/ready`
**Purpose**: Checks if the application is ready to accept traffic.

**When to use**: Kubernetes uses this to determine if the pod should receive traffic from the load balancer.

**Checks performed**:
- ✅ Database connectivity (2-second timeout)
- Add more checks as needed (RabbitMQ, Redis, etc.)

**Response when healthy**:
```json
{
  "status": "UP",
  "timestamp": "2025-02-14T20:00:00Z",
  "checks": {
    "database": "UP"
  },
  "description": "Service readiness check"
}
```

**Response when unhealthy**:
```json
{
  "status": "DOWN",
  "timestamp": "2025-02-14T20:00:00Z",
  "checks": {
    "database": "DOWN: connection refused"
  },
  "description": "Service readiness check"
}
```

**Status Codes**:
- `200 OK`: Service is ready
- `503 Service Unavailable`: Service is not ready

**Failure Action**: Kubernetes removes the pod from the service load balancer after 3 consecutive failures.

---

### 3. Startup Probe: `/health/startup`
**Purpose**: Checks if the application has finished starting up.

**When to use**: Kubernetes uses this for slow-starting containers before checking liveness/readiness.

**Checks performed**:
- ✅ Database connectivity
- Add startup-specific checks (migrations complete, cache warmed, etc.)

**Response**: Same format as readiness probe.

**Status Codes**:
- `200 OK`: Service has started
- `503 Service Unavailable`: Service is still starting

**Failure Action**: Kubernetes restarts the pod after 30 consecutive failures (150 seconds with 5-second intervals).

---

## Kubernetes Configuration

The health checks are configured in `k8s/values.yaml`:

```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: http
  initialDelaySeconds: 10    # Wait 10s before first check
  periodSeconds: 10           # Check every 10 seconds
  timeoutSeconds: 3           # Timeout after 3 seconds
  failureThreshold: 3         # Restart after 3 failures

readinessProbe:
  httpGet:
    path: /health/ready
    port: http
  initialDelaySeconds: 5      # Wait 5s before first check
  periodSeconds: 5            # Check every 5 seconds
  timeoutSeconds: 3           # Timeout after 3 seconds
  failureThreshold: 3         # Remove from LB after 3 failures

startupProbe:
  httpGet:
    path: /health/startup
    port: http
  initialDelaySeconds: 0      # Start checking immediately
  periodSeconds: 5            # Check every 5 seconds
  timeoutSeconds: 3           # Timeout after 3 seconds
  failureThreshold: 30        # Allow 150s for startup
```

---

## Testing Locally

### Test liveness:
```bash
curl http://localhost:8080/health/live
```

### Test readiness (with database):
```bash
curl http://localhost:8080/health/ready
```

### Test readiness failure (stop database):
```bash
# Stop your database
docker stop cookbook_db

# Should return 503
curl -i http://localhost:8080/health/ready
```

---

## Adding New Health Checks

To add a new dependency check (e.g., RabbitMQ, Redis):

1. **Add the check method** to `internal/handlers/health/HealthHandlers.go`:

```go
func (h *HealthHandlers) checkRabbitMQ() error {
    // Your RabbitMQ connectivity check
    return nil
}
```

2. **Update Readiness probe** to include the new check:

```go
func (h *HealthHandlers) Readiness(ctx *gin.Context) {
    checks := make(map[string]string)
    overallStatus := "UP"

    // Existing database check
    if err := h.checkDatabase(); err != nil {
        checks["database"] = "DOWN: " + err.Error()
        overallStatus = "DOWN"
    } else {
        checks["database"] = "UP"
    }

    // New RabbitMQ check
    if err := h.checkRabbitMQ(); err != nil {
        checks["rabbitmq"] = "DOWN: " + err.Error()
        overallStatus = "DOWN"
    } else {
        checks["rabbitmq"] = "UP"
    }

    // ... rest of the code
}
```

3. **Add tests** in `internal/handlers/health/HealthHandlers_test.go`

---

## Best Practices

### Liveness Probe
- ✅ Should be **lightweight**
- ✅ Should NOT check external dependencies
- ✅ Should only verify the app itself is responsive
- ❌ Don't check database, APIs, etc. (use readiness for that)

### Readiness Probe
- ✅ Should check **all critical dependencies**
- ✅ Can be more expensive (but keep under 3 seconds)
- ✅ Service can recover without restart
- ✅ Use for temporary issues (DB connection pool exhausted, etc.)

### Startup Probe
- ✅ Use for slow-starting services
- ✅ Prevents liveness from killing the pod too early
- ✅ Can have longer timeout
- ✅ Useful for services that need to warm up caches, run migrations, etc.

---

## Monitoring

Monitor health check metrics in your observability platform:
- Health check failure rate
- Time to become ready
- Frequency of restarts
- Time between restarts

---

## Replicating to Other Services

To add health checks to other services (ingredient-service, metadata-service, etc.):

1. Copy `internal/handlers/health/` directory
2. Update `internal/config/config.go` to initialize HealthHandler
3. Register routes in the service's main router
4. Update `k8s/values.yaml` with probe configuration
5. Run tests to verify

Example for ingredient-service:
```bash
cp -r recipe-service/internal/handlers/health ingredient-service/internal/handlers/
# Then update imports and config as shown above
```
