package services

import (
	"mime/multipart"

	"github.com/google/uuid"
	r "github.com/ihulsbus/cookbook/internal/repositories"
	log "github.com/sirupsen/logrus"
)

type ImageService struct {
	repo   *r.S3Repository
	logger *log.Logger
}

func NewImageService(repo *r.S3Repository, logger *log.Logger) *ImageService {
	return &ImageService{
		repo:   repo,
		logger: logger,
	}
}

func (s ImageService) UploadImage(file multipart.File, recipeID int) bool {

	// generate the file name we use in storage
	filename, err := uuid.NewRandom()
	if err != nil {
		s.logger.Errorf("error generating uuid %s", err.Error())
		return false
	}

	if b := s.repo.UploadImage(file, filename.String(), recipeID); b {
		return true
	}

	return false
}
