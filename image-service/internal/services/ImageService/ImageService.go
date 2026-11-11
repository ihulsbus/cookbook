package ImageService

import (
	"errors"
	"time"

	m "github.com/ihulsbus/cookbook/shared/models"
)

const presignTTL = 15 * time.Minute

type S3Repository interface {
	UploadImage(img m.ImageFile) error
	DeleteImage(image m.ImageData) error
	PresignGetURL(image m.ImageData, expiresIn time.Duration) (string, error)
}

type DatabaseRepository interface {
	FindAll(pagination m.PaginationRequest) ([]m.ImageData, int64, error)
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

// Best-effort: a presign failure is logged and dto is returned without a URL
// rather than failing the whole request.
func (s ImageService) withPresignedURL(image m.ImageData, dto m.ImageDataDTO) m.ImageDataDTO {
	url, err := s.s3Repo.PresignGetURL(image, presignTTL)
	if err != nil {
		s.logger.Errorf("failed presigning URL for image %s: %v", image.ID, err)
		return dto
	}
	return dto.WithURL(url)
}

// Doesn't presign every row - the admin tooling that calls this doesn't need
// to load the image bytes.
func (s ImageService) FindAll(pagination m.PaginationRequest) (m.PaginatedResponse[m.ImageDataDTO], error) {
	images, total, err := s.databaseRepo.FindAll(pagination)
	if err != nil {
		return m.PaginatedResponse[m.ImageDataDTO]{}, errors.New("internal server error")
	}

	imageDTOs := m.ImageData{}.ConvertAllToDTO(images)

	return m.PaginatedResponse[m.ImageDataDTO]{
		Data:       imageDTOs,
		Pagination: m.NewPaginationMetadata(pagination.Page, pagination.Limit, total),
	}, nil
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

	return s.withPresignedURL(image, image.ConvertToDTO()), nil
}

func (s ImageService) Create(imageFileDTO m.ImageFileDTO) (m.ImageDataDTO, error) {
	var imageFile m.ImageFile = imageFileDTO.ConvertFromDTO()
	var err error

	// Time-ordered ID, not the DB's uuidv7() default.
	imageFile.ID = m.NewID()

	if err := s.s3Repo.UploadImage(imageFile); err != nil {
		return m.ImageDataDTO{}, err
	}

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

	return s.withPresignedURL(imageData, imageData.ConvertToDTO()), nil
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
		return s.withPresignedURL(image, image.ConvertToDTO()), nil
	}

	return s.withPresignedURL(image, image.ConvertToDTO()), nil
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
