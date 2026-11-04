# Cookbook API — Readiness Checklist

_Last reviewed: 2026-04-27 against `Cookbook.yaml` v2.0 and branch `feature/V2-microservices`_

Items are sorted highest to lowest priority. Each item includes the file(s) to change and a short description of what to do.

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
- ✅ All endpoints required `administrator` role with no read-only access — split into `cookbook-user`/`cookbook-admin` (KIT-21/KIT-22)