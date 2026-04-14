package ImageService

import (
	"errors"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
)

type S3Repository interface {
	UploadImage(img m.ImageFile) error
	DeleteImage(image m.ImageData) error
}

type DatabaseRepository interface {
	FindAll() ([]m.ImageData, error)
	Find(image m.ImageData) (m.ImageData, error)
	Create(image m.ImageData) (m.ImageData, error)
	Update(image m.ImageData) (m.ImageData, error)
	Delete(image m.ImageData) error
}

type LoggerInterface interface {
	Errorf(format string, args ...interface{})
}

type RabbitMQRepository interface {
	ImageUpdatedEvent(image m.ImageData) error
}

type ImageService struct {
	databaseRepo DatabaseRepository
	rabbitmqRepo RabbitMQRepository
	s3Repo       S3Repository
	logger       LoggerInterface
}

func NewImageService(databaseRepo DatabaseRepository, rabbitmqRepo RabbitMQRepository, s3Repo S3Repository, logger LoggerInterface) *ImageService {
	return &ImageService{
		databaseRepo: databaseRepo,
		rabbitmqRepo: rabbitmqRepo,
		s3Repo:       s3Repo,
		logger:       logger,
	}
}

func (s ImageService) FindAll() ([]m.ImageDataDTO, error) {
	var images []m.ImageData

	images, err := s.databaseRepo.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	return m.ImageData{}.ConvertAllToDTO(images), nil
}

func (s ImageService) Find(imageDTO m.ImageDataDTO) (m.ImageDataDTO, error) {
	var image m.ImageData

	image, err := s.databaseRepo.Find(imageDTO.ConvertFromDTO())
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.ImageDataDTO{}, err
		default:
			return m.ImageDataDTO{}, errors.New("internal server error")
		}
	}

	return image.ConvertToDTO(), nil
}

func (s ImageService) Create(imageFileDTO m.ImageFileDTO) (m.ImageDataDTO, error) {
	var imageFile m.ImageFile = imageFileDTO.ConvertFromDTO()
	var err error

	// generate the image ID we will use to identify the file in storage
	imageFile.ID = uuid.New()

	if err := s.s3Repo.UploadImage(imageFile); err != nil {
		return m.ImageDataDTO{}, err
	}

	// Create the model for the database
	var imageData m.ImageData = m.ImageData{
		ID:         imageFile.ID,
		EntityID:   imageFile.EntityID,
		EntityType: imageFile.EntityType,
		Size:       imageFile.Size,
		Type:       imageFile.Type,
	}

	imageData, err = s.databaseRepo.Create(imageData)
	if err != nil {
		return m.ImageDataDTO{}, err
	}

	return imageData.ConvertToDTO(), nil
}

func (s ImageService) Update(imageFileDTO m.ImageFileDTO) (m.ImageDataDTO, error) {
	var imageFile m.ImageFile = imageFileDTO.ConvertFromDTO()
	var image m.ImageData
	var err error

	if image, err = s.databaseRepo.Find(m.ImageData{ID: imageFileDTO.ID, EntityID: imageFileDTO.EntityID, EntityType: imageFileDTO.EntityType}); err != nil {
		return m.ImageDataDTO{}, errors.New("unable to find existing image. cannot update something that does not exist")
	}

	if err = s.s3Repo.UploadImage(imageFile); err != nil {
		return m.ImageDataDTO{}, err
	}

	if image.Type != imageFile.Type {
		err = s.s3Repo.DeleteImage(image)
		if err != nil {
			s.logger.Errorf("error deleting old image: %v", err)
		}
	}

	image.Size = imageFile.Size
	image.Type = imageFile.Type

	image, err = s.databaseRepo.Update(image)
	if err != nil {
		return m.ImageDataDTO{}, err
	}

	err = s.rabbitmqRepo.ImageUpdatedEvent(image)
	if err != nil {
		s.logger.Errorf("failed publishing image update event to servicebus")
		return image.ConvertToDTO(), nil
	}

	return image.ConvertToDTO(), nil
}

func (s ImageService) Delete(imageDTO m.ImageDataDTO) error {
	var image m.ImageData
	var err error

	image, err = s.databaseRepo.Find(imageDTO.ConvertFromDTO())
	if err != nil {
		return errors.New("unable to find existing image. cannot delete something that does not exist")
	}

	if err = s.s3Repo.DeleteImage(image); err != nil {
		return err
	}

	err = s.databaseRepo.Delete(image)
	if err != nil {
		return err
	}

	return nil
}
