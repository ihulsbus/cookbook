# Cookbook API — Readiness Checklist

_Last reviewed: 2026-04-27 against `Cookbook.yaml` v2.0 and branch `feature/V2-microservices`_

Items are sorted highest to lowest priority. Each item includes the file(s) to change and a short description of what to do.

---

## High — Required before frontend development

### 3. All endpoints require `administrator` role — no read-only access
**Files:** every `*-service.go` router file across all services

Every single endpoint, including all GETs, is locked behind `c.KeycloakModule.Middleware("administrator")`. A frontend needs at minimum unauthenticated or lower-privilege read access to list and display data.

Split route groups: keep `administrator` for writes (POST/PUT/DELETE), and either remove auth or introduce a `viewer`/`user` role for read-only GETs.

---

## Low — Code quality / minor issues

### 8. Search service has no RabbitMQ consumer
**File:** `search-service/internal/search-service/search-service.go`

Every other service subscribes to a RabbitMQ exchange so it can react to data-change events (e.g., cache invalidation). Search service has none. If the search service caches results or builds an index, add a consumer to invalidate/update on upstream changes. If it queries other services live on every request, document this as an intentional design decision.

---

## Completed (no action needed)

The following items from the previous review have been confirmed fixed:

- ✅ Search service `main.go` calls correct `SearchService()` function
- ✅ Search handler calls `h.service.Search()` and returns results
- ✅ `instruction-service`: `DELETE` returns `204 No Content` (not 200)
- ✅ `ImageService.go` package declared as `imageservice`
- ✅ `MetadataService.go` package declared as `metadataservice`
- ✅ Image service log message corrected to "image service"
- ✅ All `GetAll` handlers return `200` + empty array (not 404) when no records exist
- ✅ Ingredient service middleware applied to the correct route groups for POST/PUT/DELETE
- ✅ Amount route path is broken
- ✅ Search endpoint uses `ShouldBindJSON` on a GET request
- ✅ Preparation time CRUD endpoints are entirely missing
- ✅ `GET /api/v2/metadata/search/all` is not in the spec
- ✅ Image service: `GET /images/search` uses same `Find` method as `GET /images/{id}`
- ✅ No pagination on any collection endpoint