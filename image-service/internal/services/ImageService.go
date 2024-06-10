package services

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type S3Repository interface {
	UploadImage(file multipart.File, filename string, recipeID int) bool
}

type LoggerInterface interface {
	Errorf(format string, args ...interface{})
}

type ImageService struct {
	repo   S3Repository
	logger LoggerInterface
}

func NewImageService(repo S3Repository, logger LoggerInterface) *ImageService {
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
