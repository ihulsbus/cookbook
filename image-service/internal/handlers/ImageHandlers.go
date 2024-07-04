package handlers

import (
	"fmt"
	"image"
	m "image-service/internal/models"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	MaxImageSize = 5 << 20 // 5 MB
)

type imageService interface {
	FindAll() ([]m.ImageDTO, error)
	Find(imageDTO m.ImageDTO) (m.ImageDTO, error)
	Create(imageDTO m.ImageDTO) (m.ImageDTO, error)
	Update(imageDTO m.ImageDTO) (m.ImageDTO, error)
	Delete(imageDTO m.ImageDTO) error
}

type ImageHandlers struct {
	imageService imageService
	logger       m.LoggerInterface
}

func NewImageHandlers(service imageService, logger m.LoggerInterface) *ImageHandlers {
	return &ImageHandlers{
		imageService: service,
		logger:       logger,
	}
}

func (h ImageHandlers) FindAll(ctx *gin.Context) {
	imageDTO, err := h.imageService.FindAll()
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

func (h ImageHandlers) Find(ctx *gin.Context) {
	var imageDTO m.ImageDTO
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

	ctx.JSON(http.StatusOK, imageDTO)
}

func (h ImageHandlers) Create(ctx *gin.Context) {
	var imageDTO m.ImageDTO
	var err error

	imageDTO.EntityID, err = uuid.Parse(ctx.Param("entityID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid entityID"})
		return
	}

	imageDTO.EntityType = ctx.Param("entityType")
	if imageDTO.EntityType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "entityType is required"})
		return
	}

	var header *multipart.FileHeader
	imageDTO.File, header, err = ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer imageDTO.File.Close()

	// Check file size
	imageDTO.Size = header.Size
	if imageDTO.Size > MaxImageSize {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "image file is too large"})
		return
	}

	// Decode the image regardless of the format
	var img image.Image
	imageDTO.Type = header.Header.Get("Content-Type")
	switch imageDTO.Type {
	case "image/jpeg":
		img, err = jpeg.Decode(imageDTO.File)
	case "image/png":
		img, err = png.Decode(imageDTO.File)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unsupported image format: %s", imageDTO.Type)})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid image"})
		return
	}

	// Check image dimensions
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width < 300 || height < 300 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "image dimensions are too small"})
		return
	}

	if width > 1000 || height > 1000 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "image dimensions are too big"})
		return
	}

	imageDTO, err = h.imageService.Create(imageDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, imageDTO)

}

func (h ImageHandlers) Update(ctx *gin.Context) {
	var imageDTO m.ImageDTO
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&imageDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// deliberaly set this to ensure the parameter ID is used instead of an accidental id in body
	// perhaps separate create/update DTO's are needed
	imageDTO.ID = id

	imageDTO, err = h.imageService.Update(imageDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, imageDTO)
}

func (h ImageHandlers) Delete(ctx *gin.Context) {
	var imageDTO m.ImageDTO
	var err error

	imageDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID"})
		return
	}

	err = h.imageService.Delete(imageDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
