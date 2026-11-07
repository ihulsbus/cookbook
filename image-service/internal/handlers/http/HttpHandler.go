package handlers

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	az "github.com/ihulsbus/cookbook/shared/authz"
	hh "github.com/ihulsbus/cookbook/shared/http"
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/ihulsbus/cookbook/shared/oidc"
)

const (
	MaxImageSize = 5 << 20 // 5 MB

	// entityTypeRecipe is the only EntityType recipe-access checks apply
	// to. Other entity types (e.g. "ingredient") are a shared catalog,
	// not owned by a user, so no recipe-level check is relevant - any
	// authenticated caller may act on them, consistent with how
	// ingredient-service's own catalog routes work.
	entityTypeRecipe = "recipe"
)

type imageService interface {
	FindAll(pagination m.PaginationRequest) (m.PaginatedResponse[m.ImageDataDTO], error)
	Find(imageDTO m.ImageDataDTO) (m.ImageDataDTO, error)
	Create(imageDTO m.ImageFileDTO) (m.ImageDataDTO, error)
	Update(imageDTO m.ImageFileDTO) (m.ImageDataDTO, error)
	Delete(imageDTO m.ImageDataDTO) error
}

// recipeAuthorizer is the subset of *authz.Authorizer this handler
// needs, kept as a local interface (like imageService above) so it can
// be tested without a real DB-backed Authorizer.
type recipeAuthorizer interface {
	CanRead(ctx context.Context, principal *oidc.Principal, recipeID uuid.UUID) (bool, error)
	CanWrite(ctx context.Context, principal *oidc.Principal, recipeID uuid.UUID) (bool, error)
}

type HttpHandler struct {
	imageService imageService
	authorizer   recipeAuthorizer
	logger       m.LoggerInterface
}

func NewHttpHandler(service imageService, authorizer recipeAuthorizer, logger m.LoggerInterface) *HttpHandler {
	return &HttpHandler{
		imageService: service,
		authorizer:   authorizer,
		logger:       logger,
	}
}

// principal extracts the authenticated caller from the Gin context (set
// by oidc.Module.RequireAuth, which every route in this service's router
// already runs).
func (h HttpHandler) principal(ctx *gin.Context) (*oidc.Principal, bool) {
	principal, ok := oidc.PrincipalFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated principal"})
		return nil, false
	}
	return principal, true
}

// checkRecipeAccess authorizes principal for perm on recipeID when
// entityType is "recipe"; any other entity type (e.g. "ingredient") has
// no recipe-level access model and is allowed through (RequireAuth
// already ran). On denial it writes the response itself (404, matching
// recipe-service's "don't reveal existence" behavior) and returns false.
func (h HttpHandler) checkRecipeAccess(ctx *gin.Context, principal *oidc.Principal, entityType string, recipeID uuid.UUID, perm az.Permission) bool {
	if entityType != entityTypeRecipe {
		return true
	}

	var allowed bool
	var err error
	switch perm {
	case az.PermissionRead:
		allowed, err = h.authorizer.CanRead(context.Background(), principal, recipeID)
	case az.PermissionWrite:
		allowed, err = h.authorizer.CanWrite(context.Background(), principal, recipeID)
	}

	if err != nil {
		if errors.Is(err, az.ErrProjectionMiss) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "recipe access is not yet known - try again shortly"})
			return false
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "authorization check failed"})
		return false
	}

	if !allowed {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no images found"})
		return false
	}

	return true
}

// FindAll lists every image regardless of owner - admin-only at the
// router level (KIT-29), since there's no cheap way to filter this list
// against the recipe_access projection row by row.
func (h HttpHandler) FindAll(ctx *gin.Context) {
	pagination, err := hh.ParsePagination(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imageDTO, err := h.imageService.FindAll(pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, imageDTO)
}

func (h HttpHandler) Find(ctx *gin.Context) {
	principal, ok := h.principal(ctx)
	if !ok {
		return
	}

	var imageDTO m.ImageDataDTO
	var err error

	imageDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID"})
		return
	}

	imageDTO, err = h.imageService.Find(imageDTO)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no images found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// :id is an IMAGE id, not a recipe id - the image must be resolved
	// first (above) to learn which recipe it belongs to before it can be
	// authorized (KIT-29).
	if !h.checkRecipeAccess(ctx, principal, imageDTO.EntityType, imageDTO.EntityID, az.PermissionRead) {
		return
	}

	ctx.JSON(http.StatusOK, imageDTO)
}

func (h HttpHandler) SearchByRecipe(ctx *gin.Context) {
	principal, ok := h.principal(ctx)
	if !ok {
		return
	}

	var imageDTO m.ImageDataDTO
	var err error

	imageDTO.EntityType = ctx.Query("entityType")
	imageDTO.EntityID, err = uuid.Parse(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID"})
		return
	}

	// entityType/entityID are already known from the query here, so the
	// access check can run before touching the database at all.
	if !h.checkRecipeAccess(ctx, principal, imageDTO.EntityType, imageDTO.EntityID, az.PermissionRead) {
		return
	}

	imageDTO, err = h.imageService.Find(imageDTO)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no images found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, imageDTO)
}

func (h HttpHandler) Create(ctx *gin.Context) {
	principal, ok := h.principal(ctx)
	if !ok {
		return
	}

	var file multipart.File
	var imageFileDTO m.ImageFileDTO
	var err error

	imageFileDTO.EntityID, err = uuid.Parse(ctx.Param("entityID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid entityID"})
		return
	}

	imageFileDTO.EntityType = ctx.Param("entityType")
	if imageFileDTO.EntityType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "entityType is required"})
		return
	}

	// entityType/entityID come straight from the path, so write access
	// can be checked before accepting the upload at all.
	if !h.checkRecipeAccess(ctx, principal, imageFileDTO.EntityType, imageFileDTO.EntityID, az.PermissionWrite) {
		return
	}

	var header *multipart.FileHeader
	file, header, err = ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	image := io.TeeReader(file, &imageFileDTO.File)

	err = h.verifyImage(&imageFileDTO, image, header)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imageDataDTO, err := h.imageService.Create(imageFileDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, imageDataDTO)

}

func (h HttpHandler) Update(ctx *gin.Context) {
	principal, ok := h.principal(ctx)
	if !ok {
		return
	}

	var file multipart.File
	var imageFileDTO m.ImageFileDTO
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID"})
		return
	}

	// deliberaly set this to ensure the parameter ID is used instead of an accidental id in body
	// perhaps separate create/update DTO's are needed
	imageFileDTO.ID = id

	// :id is an IMAGE id - resolve it first to learn which recipe it
	// belongs to before authorizing the write (KIT-29).
	existing, err := h.imageService.Find(m.ImageDataDTO{ID: id})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no images found"})
		return
	}
	if !h.checkRecipeAccess(ctx, principal, existing.EntityType, existing.EntityID, az.PermissionWrite) {
		return
	}
	imageFileDTO.EntityID = existing.EntityID
	imageFileDTO.EntityType = existing.EntityType

	var header *multipart.FileHeader
	file, header, err = ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	image := io.TeeReader(file, &imageFileDTO.File)

	err = h.verifyImage(&imageFileDTO, image, header)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imageDataDTO, err := h.imageService.Update(imageFileDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, imageDataDTO)
}

func (h HttpHandler) Delete(ctx *gin.Context) {
	principal, ok := h.principal(ctx)
	if !ok {
		return
	}

	var imageDTO m.ImageDataDTO
	var err error

	imageDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID"})
		return
	}

	// :id is an IMAGE id - resolve it first to learn which recipe it
	// belongs to before authorizing the delete (KIT-29).
	existing, err := h.imageService.Find(imageDTO)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no images found"})
		return
	}
	if !h.checkRecipeAccess(ctx, principal, existing.EntityType, existing.EntityID, az.PermissionWrite) {
		return
	}

	err = h.imageService.Delete(imageDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h HttpHandler) verifyImage(imageFileDTO *m.ImageFileDTO, file io.Reader, header *multipart.FileHeader) error {
	var err error

	// Check file size
	imageFileDTO.Size = header.Size
	if imageFileDTO.Size > MaxImageSize {

		return errors.New("image file is too large")
	}

	// Decode the image regardless of the format
	var img image.Image
	imageFileDTO.Type = header.Header.Get("Content-Type")
	switch imageFileDTO.Type {
	case "image/jpeg":
		img, err = jpeg.Decode(file)
	case "image/png":
		img, err = png.Decode(file)
	default:
		return fmt.Errorf("unsupported image format: %s", imageFileDTO.Type)
	}
	if err != nil {
		return fmt.Errorf("invalid image: %v", err)
	}

	// Check image dimensions
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width < 300 || height < 300 {
		return errors.New("image dimensions are too small. Dimensions need to be between 300x300 and 1000x1000")
	}

	if width > 1000 || height > 1000 {
		return errors.New("image dimensions are too big. Dimensions need to be between 300x300 and 1000x1000")
	}

	return nil
}
