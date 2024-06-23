package services

import (
	m "image-service/internal/models"

	"github.com/google/uuid"
)

type S3Repository interface {
	UploadImage(img m.Image) error
}

type ImageRepository interface {
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

func (s ImageService) Create(imageDTO m.ImageDTO) (m.ImageDTO, error) {
	var image m.Image = imageDTO.ConvertFromDTO()
	var err error

	// generate the image ID we will use to identify the file in storage
	image.ID, err = uuid.NewRandom()
	if err != nil {
		s.logger.Errorf("error generating uuid %s", err.Error())
		return m.ImageDTO{}, err
	}

	if err := s.repo.UploadImage(image); err != nil {
		return m.ImageDTO{}, err
	}

	return image.ConvertToDTO(), nil
}
