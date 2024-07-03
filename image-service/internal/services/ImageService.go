package services

import (
	"errors"
	m "image-service/internal/models"

	"github.com/google/uuid"
)

type S3Repository interface {
	UploadImage(img m.Image) error
	DeleteImage(image m.Image) error
}

type ImageRepository interface {
	FindAll() ([]m.Image, error)
	Find(image m.Image) (m.Image, error)
	Create(image m.Image) (m.Image, error)
	Update(image m.Image) (m.Image, error)
	Delete(image m.Image) error
}

type LoggerInterface interface {
	Errorf(format string, args ...interface{})
}

type ImageService struct {
	imageRepo ImageRepository
	s3Repo    S3Repository
	logger    LoggerInterface
}

func NewImageService(imageRepo ImageRepository, s3Repo S3Repository, logger LoggerInterface) *ImageService {
	return &ImageService{
		imageRepo: imageRepo,
		s3Repo:    s3Repo,
		logger:    logger,
	}
}

func (s ImageService) FindAll() ([]m.ImageDTO, error) {
	var images []m.Image

	images, err := s.imageRepo.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	return m.Image{}.ConvertAllToDTO(images), nil
}

func (s ImageService) Find(imageDTO m.ImageDTO) (m.ImageDTO, error) {
	var image m.Image

	image, err := s.imageRepo.Find(imageDTO.ConvertFromDTO())
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.ImageDTO{}, err
		default:
			return m.ImageDTO{}, errors.New("internal server error")
		}
	}

	return image.ConvertToDTO(), nil
}

func (s ImageService) Create(imageDTO m.ImageDTO) (m.ImageDTO, error) {
	var image m.Image = imageDTO.ConvertFromDTO()
	var err error

	// generate the image ID we will use to identify the file in storage
	image.ID = uuid.New()

	if err := s.s3Repo.UploadImage(image); err != nil {
		return m.ImageDTO{}, err
	}

	image, err = s.imageRepo.Create(image)
	if err != nil {
		return m.ImageDTO{}, err
	}

	return image.ConvertToDTO(), nil
}

func (s ImageService) Update(imageDTO m.ImageDTO) (m.ImageDTO, error) {
	var image m.Image = imageDTO.ConvertFromDTO()
	var err error

	if _, err = s.imageRepo.Find(image); err != nil {
		return m.ImageDTO{}, errors.New("unable to find existing image. cannot update something that does not exist")
	}

	if err = s.s3Repo.UploadImage(image); err != nil {
		return m.ImageDTO{}, err
	}

	image, err = s.imageRepo.Update(image)
	if err != nil {
		return m.ImageDTO{}, err
	}

	return image.ConvertToDTO(), nil
}

func (s ImageService) Delete(imageDTO m.ImageDTO) error {
	var err error

	_, err = s.imageRepo.Find(imageDTO.ConvertFromDTO())
	if err != nil {
		return errors.New("unable to find existing image. cannot delete something that does not exist")
	}

	if err = s.s3Repo.DeleteImage(imageDTO.ConvertFromDTO()); err != nil {
		return err
	}

	err = s.imageRepo.Delete(imageDTO.ConvertFromDTO())
	if err != nil {
		return err
	}

	return nil
}
