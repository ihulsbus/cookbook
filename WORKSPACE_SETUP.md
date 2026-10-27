# Go Workspace Setup

This project uses Go workspaces to manage shared modules across microservices, providing a seamless local development experience while maintaining compatibility with Docker and CI/CD.

## Overview

The Cookbook project is organized as a Go workspace monorepo with 7 microservices and 9 shared modules. This approach allows you to work on all services simultaneously without needing to publish or version shared code during development.

## Structure

```
cookbook/
├── go.work                         # Workspace configuration
├── go.work.sum                     # Workspace checksums
│
├── recipe-service/
│   ├── go.mod                      # Module with replace directives
│   ├── go.sum
│   └── cmd/main.go
│
├── ingredient-service/
│   ├── go.mod                      # Module with replace directives
│   ├── go.sum
│   └── cmd/main.go
│
├── instruction-service/
│   ├── go.mod
│   ├── go.sum
│   └── cmd/main.go
│
├── metadata-service/
│   ├── go.mod
│   ├── go.sum
│   └── cmd/main.go
│
├── image-service/
│   ├── go.mod
│   ├── go.sum
│   └── cmd/main.go
│
├── search-service/
│   ├── go.mod
│   ├── go.sum
│   └── cmd/main.go
│
├── notification-service/           # In development
│   ├── go.mod
│   ├── go.sum
│   └── cmd/main.go
│
└── shared/                         # Shared modules
    ├── cache/
    │   ├── go.mod
    │   └── cache.go
    ├── healthchecks/
    │   ├── go.mod
    │   └── healthchecks.go
    ├── http/
    │   ├── go.mod
    │   └── http.go
    ├── httpclient/
    │   ├── go.mod
    │   └── httpclient.go
    ├── imageclient/
    │   ├── go.mod
    │   └── imageclient.go
    ├── keycloak/
    │   ├── go.mod
    │   └── keycloak.go
    ├── models/
    │   ├── go.mod
    │   └── *.go                    # All shared data models
    ├── rabbitmq/
    │   ├── go.mod
    │   └── rabbitmq.go
    └── recipeclient/
        ├── go.mod
        └── recipeclient.go
```

## How It Works

### 1. Workspace Resolution (Local Development)

The `go.work` file at the root tells Go to use local versions of all modules:

```go
go 1.26.0

use (
    ./image-service
    ./ingredient-service
    ./instruction-service
    ./metadata-service
    ./notification-service
    ./recipe-service
    ./search-service
    ./shared/cache
    ./shared/healthchecks
    ./shared/http
    ./shared/httpclient
    ./shared/imageclient
    ./shared/keycloak
    ./shared/models
    ./shared/rabbitmq
    ./shared/recipeclient
)
```

When you build or run services, Go automatically uses the local shared module code without needing to fetch from GitHub or publish versions.

### 2. Replace Directives (Fallback & Docker)

Each service's `go.mod` has replace directives as a fallback mechanism:

```go
module recipe-service

go 1.26.0

require (
    github.com/ihulsbus/cookbook/shared/models v0.0.0
    github.com/ihulsbus/cookbook/shared/keycloak v0.0.0
    github.com/ihulsbus/cookbook/shared/rabbitmq v0.0.0
    // ... other dependencies
)

replace (
    github.com/ihulsbus/cookbook/shared/models => ../shared/models
    github.com/ihulsbus/cookbook/shared/keycloak => ../shared/keycloak
    github.com/ihulsbus/cookbook/shared/rabbitmq => ../shared/rabbitmq
    github.com/ihulsbus/cookbook/shared/recipeclient => ../shared/recipeclient
)
```

This ensures:
- ✅ **Docker builds work** - Workspace isn't active in Docker, replace directives handle paths
- ✅ **CI/CD works** - Can build services without workspace file
- ✅ **Backward compatibility** - Works in non-workspace environments
- ✅ **Relative paths** - No need for absolute paths or symlinks

## Common Commands

### Building Services

```bash
# Build from workspace root
cd /path/to/cookbook
go build ./recipe-service/cmd/main.go

# Build from within a service directory
cd recipe-service
go build ./cmd/main.go

# Build all services
for dir in *-service; do
    echo "Building $dir..."
    go build ./$dir/cmd/main.go
done

# Build with output binary
go build -o bin/recipe-service ./recipe-service/cmd/main.go
```

### Running Services

```bash
# Run from workspace root
go run ./recipe-service/cmd/main.go

# Run from within service directory
cd recipe-service
go run ./cmd/main.go

# Run with environment variables
export cbb_debug=true
export cbb_database_host=localhost
go run ./recipe-service/cmd/main.go
```

### Testing

```bash
# Test all services and shared modules
go test ./...

# Test specific service
go test ./recipe-service/...

# Test shared module
go test ./shared/models/...

# Test with verbose output
go test -v ./recipe-service/...

# Test with coverage
go test -cover ./...

# Test specific package
go test ./recipe-service/internal/handlers/...

# Run tests in parallel
go test -parallel 4 ./...
```

### Managing Dependencies

```bash
# Add a new dependency to a service
cd recipe-service
go get github.com/some/package@v1.2.3

# Update dependencies
cd recipe-service
go get -u ./...

# Tidy dependencies (remove unused, add missing)
cd recipe-service
go mod tidy

# Download dependencies
cd recipe-service
go mod download

# Verify dependencies
cd recipe-service
go mod verify
```

### Working with Shared Modules

```bash
# No special commands needed!
# Just edit the code in shared/ and rebuild services

# Example: Update the models package
vim shared/models/RecipeModels.go

# All services automatically see the changes
go build ./recipe-service/cmd/main.go
go test ./recipe-service/...
```

### Workspace Management

```bash
# View workspace info
go work use

# Sync workspace (rarely needed, happens automatically)
go work sync

# Edit workspace
go work use -r ./new-service

# Remove from workspace
go work use -r ./old-service
```

## Making Changes to Shared Modules

### 1. Edit the shared module code

```bash
# Edit any file in the shared modules
vim shared/models/PaginationModels.go
vim shared/httpclient/httpclient.go
vim shared/keycloak/keycloak.go
```

### 2. All services see changes immediately

No need to:
- Publish to GitHub
- Create git tags
- Update version numbers
- Run go get in services
- Restart anything

Just build and test:
```bash
go build ./recipe-service/cmd/main.go
go test ./recipe-service/...
```

### 3. Verify changes across services

```bash
# Test the shared module
go test ./shared/models/...

# Test all services that depend on it
go test ./recipe-service/...
go test ./ingredient-service/...
go test ./metadata-service/...

# Or test everything
go test ./...
```

### 4. Commit changes

```bash
# Commit both shared code and service code together
git add shared/models/PaginationModels.go
git add recipe-service/internal/handlers/RecipeHandlers.go
git commit -m "Add pagination support to recipe handlers"
```

## Docker Builds

When building Docker images, the workspace is **NOT active**. The `replace` directives in each service's `go.mod` handle module resolution.

### Example Dockerfile

```dockerfile
# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Copy entire monorepo (needed for replace directives)
COPY . .

# Build specific service
WORKDIR /build/recipe-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o recipe-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/recipe-service/recipe-service .

EXPOSE 8080

CMD ["./recipe-service"]
```

### Building Docker Images

```bash
# Build from repository root (important!)
cd /path/to/cookbook

# Build recipe service
docker build -t cookbook/recipe-service:latest -f recipe-service/Dockerfile .

# Build all services
for service in recipe ingredient instruction metadata image search; do
    docker build -t cookbook/${service}-service:latest -f ${service}-service/Dockerfile .
done
```

**Important**: Always build from the repository root, not from individual service directories, so the COPY command can access shared modules.

## Troubleshooting

### Issue: `unknown revision shared/models/v0.0.0`

**Cause**: Go is trying to fetch from GitHub instead of using local modules.

**Solution**:
1. Ensure you're in the workspace root
2. Check that `go.work` exists and lists the module
3. Verify replace directives in service's `go.mod`

```bash
cd /path/to/cookbook
cat go.work  # Should list all services and shared modules
cd recipe-service
cat go.mod   # Should have replace directives
```

### Issue: Changes to shared modules not reflected

**Cause**: Stale build cache or wrong directory.

**Solution**: Clean cache and rebuild
```bash
go clean -cache -testcache -modcache
go build ./recipe-service/cmd/main.go
```

### Issue: `go.work` not found

**Cause**: Running commands from wrong directory.

**Solution**: Navigate to workspace root
```bash
cd /home/ianhulsbus/Documents/workspace/personal/cookbook
# Or set GOWORK environment variable
export GOWORK=/path/to/cookbook/go.work
```

### Issue: Docker build fails with "module not found"

**Cause**: Building from wrong directory or incorrect COPY command.

**Solution**: Build from repository root
```bash
cd /path/to/cookbook  # Repository root
docker build -t myservice -f service/Dockerfile .
# Note the . at the end - context is repository root
```

### Issue: Tests fail after updating shared module

**Cause**: Breaking changes in shared module API.

**Solution**: Update all affected services
```bash
# Find usages
grep -r "OldFunction" *-service/

# Update and test each service
for dir in *-service; do
    echo "Testing $dir..."
    go test ./$dir/...
done
```

### Issue: go.mod and go.work.sum conflicts

**Cause**: Outdated checksums.

**Solution**: Sync and tidy
```bash
go work sync
cd recipe-service && go mod tidy
cd ../ingredient-service && go mod tidy
# Repeat for all services
```

### Issue: Service doesn't see latest shared module changes

**Cause**: Module cache or IDE cache.

**Solution**:
```bash
# Clear Go cache
go clean -modcache

# Restart IDE/LSP (if using VS Code, GoLand, etc.)
# For VS Code: Cmd/Ctrl + Shift + P -> "Go: Restart Language Server"

# Rebuild
go build ./recipe-service/cmd/main.go
```

## Benefits of This Setup

✅ **No version conflicts** - All services always use the same shared code version
✅ **Instant feedback** - Change shared code, see results immediately across all services
✅ **No publishing required** - Don't need to tag/publish/wait for local development
✅ **Works with Docker** - Replace directives handle builds without workspace
✅ **Works with CI/CD** - Can build services independently in pipelines
✅ **Easy testing** - Test shared modules and services together
✅ **Atomic commits** - Commit shared code and service changes together
✅ **Simpler workflow** - No manual version bumps, no go get, no confusion
✅ **Better IDE support** - Go To Definition works across module boundaries
✅ **Faster iteration** - No publish/pull cycle for shared code changes

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Test All Services

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.26'

      - name: Test all services
        run: go test ./...

      - name: Build all services
        run: |
          go build ./recipe-service/cmd/main.go
          go build ./ingredient-service/cmd/main.go
          go build ./instruction-service/cmd/main.go
          go build ./metadata-service/cmd/main.go
          go build ./image-service/cmd/main.go
          go build ./search-service/cmd/main.go
```

The workspace file works automatically in CI/CD environments.

## Migration Path to Separate Repos (Future)

When you're ready to split shared modules into separate repositories:

### Option 1: Keep Monorepo
Continue using the workspace - it scales well even for large projects.

### Option 2: Extract Shared Modules

1. Create separate GitHub repos for each shared module
```bash
gh repo create ihulsbus/cookbook-models --public
gh repo create ihulsbus/cookbook-httpclient --public
```

2. Move code and tag versions
```bash
cd shared/models
git init
git remote add origin github.com/ihulsbus/cookbook-models
git add .
git commit -m "Initial commit"
git tag v0.1.0
git push origin main --tags
```

3. Update service `go.mod` files
```go
require (
    github.com/ihulsbus/cookbook-models v0.1.0
    github.com/ihulsbus/cookbook-httpclient v0.1.0
)

// Remove replace directives
```

4. Remove shared modules from workspace
```bash
# Edit go.work, remove shared module paths
```

5. Update services
```bash
cd recipe-service
go get github.com/ihulsbus/cookbook-models@v0.1.0
go mod tidy
```

The workspace approach makes this migration **optional** and **low-pressure** - only do it if you really need separate versioning.

## Best Practices

### Do:
✅ Make atomic commits that include both shared and service code changes
✅ Test shared modules thoroughly before committing
✅ Run `go test ./...` from root before pushing
✅ Keep shared modules focused and well-documented
✅ Use semantic commit messages that describe the full change
✅ Run `go mod tidy` in services after updating dependencies

### Don't:
❌ Don't manually edit `go.work.sum` or `go.sum` files
❌ Don't use `go get` with local shared modules (workspace handles it)
❌ Don't build Docker images from service directories (build from root)
❌ Don't commit with failing tests
❌ Don't make breaking changes to shared modules without updating all services
❌ Don't use absolute paths in replace directives

## IDE Configuration

### VS Code

Add to `.vscode/settings.json`:
```json
{
  "go.useLanguageServer": true,
  "go.toolsManagement.autoUpdate": true,
  "go.lintOnSave": "package",
  "go.buildOnSave": "off",
  "gopls": {
    "experimentalWorkspaceModule": true
  }
}
```

### GoLand / IntelliJ IDEA

1. Open the repository root folder
2. GoLand will automatically detect `go.work`
3. Enable Go Modules integration in Settings → Go → Go Modules
4. Mark shared folders as "Library Root" for better code completion

### Vim/Neovim with gopls

gopls automatically detects and uses `go.work` when present in the workspace root.

## Performance Tips

### Faster Builds
```bash
# Use build cache
go build -o bin/service ./service/cmd/main.go

# Parallel compilation (default, but can be explicit)
go build -p 8 ./service/cmd/main.go

# Skip tests when building
go build -tags=skiptest ./...
```

### Faster Tests
```bash
# Test in parallel
go test -parallel 8 ./...

# Cache test results
go test -count=1 ./...  # Disable cache
go test ./...           # Use cache

# Short tests only
go test -short ./...
```

## Summary

The Go workspace setup provides the best of both worlds:
- **Local development**: Fast, seamless, no version management needed
- **Production deployment**: Works with Docker, CI/CD, and standard Go tooling
- **Flexibility**: Easy to migrate to separate repos if needed in the future

Just edit code and run `go build` or `go test` - the workspace handles the rest!
