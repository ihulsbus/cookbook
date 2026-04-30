# Cookbook API — Readiness Checklist

_Last reviewed: 2026-04-27 against `Cookbook.yaml` v2.0 and branch `feature/V2-microservices`_

Items are sorted highest to lowest priority. Each item includes the file(s) to change and a short description of what to do.

---

## Critical — Blocking functionality

### 1. Amount route path is broken
**File:** `ingredient-service/internal/ingredient-service/ingredient-service.go:49`

`v2.Group("amount")` is missing the leading `/`, which makes Gin resolve all `/amount` routes to `/api/v2amount` instead of `/api/v2/amount`. Change to `v2.Group("/amount")`.

---

### 2. Preparation time CRUD endpoints are entirely missing
**Files to create:**
- `metadata-service/internal/handlers/preparationtime/PreparationTimeHandlers.go`
- `metadata-service/internal/handlers/preparationtime/PreparationTimeHandlers_test.go`

**File to update:** `metadata-service/internal/metadata-service/MetadataService.go`

Five endpoints are defined in the spec but have no code at all:
- `GET /api/v2/metadata/preparationtime`
- `POST /api/v2/metadata/preparationtime`
- `GET /api/v2/metadata/preparationtime/{PreparationTimeID}`
- `PUT /api/v2/metadata/preparationtime/{PreparationTimeID}`
- `DELETE /api/v2/metadata/preparationtime/{PreparationTimeID}`

Implement handlers, wire a `PreparationTimeService` interface (following the same pattern as `DifficultyLevelHandlers`), and register the routes in `MetadataService.go`.

---

## High — Required before frontend development

### 3. All endpoints require `administrator` role — no read-only access
**Files:** every `*-service.go` router file across all services

Every single endpoint, including all GETs, is locked behind `c.KeycloakModule.Middleware("administrator")`. A frontend needs at minimum unauthenticated or lower-privilege read access to list and display data.

Split route groups: keep `administrator` for writes (POST/PUT/DELETE), and either remove auth or introduce a `viewer`/`user` role for read-only GETs.

---

### 4. Search endpoint uses `ShouldBindJSON` on a GET request
**File:** `search-service/internal/handlers/searchHandlers.go:31`
**File:** `search-service/internal/search-service/search-service.go:47`
**Spec:** `GET /search` with `query`, `limit`, `page` fields

The handler binds a JSON request body on a `GET` endpoint. Many HTTP clients, proxies, and browser `fetch` implementations strip or reject GET bodies. Change the handler to read `query`, `limit`, and `page` from query parameters (`ctx.Query(...)`), and update the spec to use `in: query` parameters instead of a `requestBody`.

---

## Medium — Spec non-compliance / undocumented behaviour

### 5. `GET /api/v2/metadata/search/all` is not in the spec
**File:** `metadata-service/internal/metadata-service/MetadataService.go:203–207`

The route `/metadata/search/all` → `SearchHandlers.GetAllMetadata` is registered and functional but absent from `Cookbook.yaml`. Either add the endpoint to the spec or remove it from the service. The `POST /metadata/search` route already covers filtered searches; decide if a separate "get all" search endpoint is needed and document it.

---

### 6. No pagination on any collection endpoint
**Spec note:** Every `GET` collection endpoint description says _"this endpoint does not support pagination. This will be added in the future."_

Currently there is zero pagination infrastructure. All list endpoints return every record. Before a frontend is usable at scale, add `page` / `limit` query parameters (or cursor-based pagination) to at minimum: `GET /recipe`, `GET /ingredient`, `GET /unit`, `GET /images`, `GET /metadata/tag`, `GET /metadata/category`, `GET /metadata/cuisinetype`, `GET /metadata/difficultylevel`, `GET /metadata/preparationtime`.

---

## Low — Code quality / minor issues

### 7. Image service: `GET /images/search` uses same `Find` method as `GET /images/{id}`
**File:** `image-service/internal/handlers/http/HttpHandler.go:83–106`

`SearchByRecipe` calls `h.imageService.Find(imageDTO)` — the same method used for fetching by ID. If the backing query logic is shared, ensure the repository distinguishes between lookup-by-ID and lookup-by-entity-type/id. Currently both paths are controlled by which fields are populated on the DTO, which is fragile.

---

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
