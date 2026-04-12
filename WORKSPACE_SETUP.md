# Go Workspace Setup

This project uses Go workspaces to manage shared modules across microservices.

## Structure

```
cookbook/
├── go.work                     # Workspace configuration
├── recipe-service/
│   └── go.mod                  # Uses replace directives
├── ingredient-service/
│   └── go.mod                  # Uses replace directives
├── metadata-service/
│   └── go.mod                  # Uses replace directives
├── ...
└── shared/
    ├── httpclient/
    │   └── go.mod
    ├── recipeclient/
    │   └── go.mod
    ├── imageclient/
    │   └── go.mod
    ├── keycloak/
    │   └── go.mod
    └── rabbitmq/
        └── go.mod
```

## How It Works

### 1. Workspace Resolution (Local Development)

The `go.work` file tells Go to use local versions of shared modules:

```go
use (
    ./recipe-service
    ./ingredient-service
    ./shared/httpclient
    ./shared/keycloak
    // ...
)
```

When you build or run services, Go automatically uses the local shared module code.

### 2. Replace Directives (Fallback & Docker)

Each service's `go.mod` has replace directives as a fallback:

```go
replace (
    github.com/ihulsbus/cookbook/shared/keycloak => ../shared/keycloak
    github.com/ihulsbus/cookbook/shared/rabbitmq => ../shared/rabbitmq
)
```

This ensures:
- **Docker builds work** (workspace isn't active in Docker)
- **CI/CD works** without needing the workspace
- **Backward compatibility** with non-workspace environments

## Common Commands

### ✅ DO Use These:

```bash
# Build a service (from any directory in the workspace)
go build ./recipe-service/cmd/recipe-service

# Run tests
go test ./...

# Run a specific service's tests
go test ./recipe-service/...

# Build from within a service directory
cd recipe-service
go build ./cmd/recipe-service

# Add a dependency to a service
cd recipe-service
go get github.com/some/package@v1.2.3

# Tidy up dependencies
cd recipe-service
go mod tidy
```

### ❌ DON'T Use These:

```bash
# Don't use go mod download - the workspace handles dependencies
go mod download  # ❌ Will try to fetch from GitHub

# Don't manually edit go.sum files
# They're managed automatically
```

## Making Changes to Shared Modules

### 1. Edit the shared module code

```bash
# Edit shared/httpclient/httpclient.go
vim shared/httpclient/httpclient.go
```

### 2. All services see the change immediately

```bash
# No need to publish, tag, or update versions
# Just build and test
go build ./recipe-service/cmd/recipe-service
```

### 3. Run tests to verify

```bash
# Test the shared module
go test ./shared/httpclient/...

# Test services that use it
go test ./recipe-service/...
go test ./metadata-service/...
```

## Docker Builds

When building Docker images, the workspace is NOT active. The `replace` directives in each service's `go.mod` handle this.

**Example Dockerfile pattern:**

```dockerfile
FROM golang:1.23 AS builder

WORKDIR /app

# Copy the entire monorepo (needed for replace directives)
COPY . .

# Build the specific service
WORKDIR /app/recipe-service
RUN go build -o /app/bin/recipe-service ./cmd/recipe-service

FROM alpine:latest
COPY --from=builder /app/bin/recipe-service /app/recipe-service
CMD ["/app/recipe-service"]
```

## Troubleshooting

### Issue: `unknown revision shared/keycloak/v0.0.0`

**Cause**: Go is trying to fetch from GitHub instead of using local modules.

**Solution**: Make sure you're in the workspace root directory and the `replace` directives are in place.

```bash
cd /path/to/cookbook
go build ./recipe-service/cmd/recipe-service
```

### Issue: Changes to shared module not reflected

**Cause**: Stale build cache.

**Solution**: Clean and rebuild:

```bash
go clean -cache
go build ./recipe-service/cmd/recipe-service
```

### Issue: `go.work` not found

**Cause**: You're in the wrong directory.

**Solution**: Navigate to the workspace root:

```bash
cd /home/ianhulsbus/Documents/workspace/personal/cookbook
```

## Benefits of This Setup

✅ **No version conflicts** - All services use the same local shared code
✅ **Instant feedback** - Change shared code, see results immediately
✅ **No publishing required** - Don't need to tag/publish for local development
✅ **Works with Docker** - Replace directives handle builds without workspace
✅ **Works with CI/CD** - Can build services independently
✅ **Easy testing** - Test shared modules and services together
✅ **Simpler workflow** - No need for `go get`, manual version bumps, etc.

## Migration Path to Separate Repos (Future)

When you're ready to split shared modules into separate repositories:

1. Create separate GitHub repos for each shared module
2. Tag initial versions (e.g., `v0.1.0`)
3. Update service `go.mod` files to reference the GitHub URLs
4. Remove `replace` directives
5. Remove shared modules from workspace
6. Keep services in monorepo or split them too

The workspace approach makes this migration optional and low-pressure.
