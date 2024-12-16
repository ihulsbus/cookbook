package handlers

import (
	"errors"
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
	FindAll() ([]m.ImageDataDTO, error)
	Find(imageDTO m.ImageDataDTO) (m.ImageDataDTO, error)
	Create(imageDTO m.ImageFileDTO) (m.ImageDataDTO, error)
	Update(imageDTO m.ImageFileDTO) (m.ImageDataDTO, error)
	Delete(imageDTO m.ImageDataDTO) error
}

type HttpHandlers struct {
	imageService imageService
	logger       m.LoggerInterface
}

func NewHttpHandler(service imageService, logger m.LoggerInterface) *HttpHandlers {
	return &HttpHandlers{
		imageService: service,
		logger:       logger,
	}
}

func (h HttpHandlers) FindAll(ctx *gin.Context) {
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

func (h HttpHandlers) Find(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, imageDTO)
}

func (h HttpHandlers) Create(ctx *gin.Context) {
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

	var header *multipart.FileHeader
	imageFileDTO.File, header, err = ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer imageFileDTO.File.Close()

	err = h.verifyImage(imageFileDTO, header)
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

func (h HttpHandlers) Update(ctx *gin.Context) {
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

	var header *multipart.FileHeader
	imageFileDTO.File, header, err = ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer imageFileDTO.File.Close()

	err = h.verifyImage(imageFileDTO, header)
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

func (h HttpHandlers) Delete(ctx *gin.Context) {
	var imageDTO m.ImageDataDTO
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

func (h HttpHandlers) verifyImage(imageFileDTO m.ImageFileDTO, header *multipart.FileHeader) error {
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
		img, err = jpeg.Decode(imageFileDTO.File)
	case "image/png":
		img, err = png.Decode(imageFileDTO.File)
	default:
		return fmt.Errorf("unsupported image format: %s", imageFileDTO.Type)
	}
	if err != nil {
		return errors.New("invalid image")
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
