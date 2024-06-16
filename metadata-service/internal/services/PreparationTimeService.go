package services

import (
	"errors"

	m "metadata-service/internal/models"

	"github.com/google/uuid"
)

type PreparationTimeRepository interface {
	FindAll() ([]m.PreparationTime, error)
	FindSingle(preparationTime m.PreparationTime) (m.PreparationTime, error)
	Create(preparationTime m.PreparationTime) (m.PreparationTime, error)
	Update(preparationTime m.PreparationTime) (m.PreparationTime, error)
	Delete(preparationTime m.PreparationTime) error
}
type PreparationTimeService struct {
	repo PreparationTimeRepository
}

// NewPreparationTimeService creates a new PreparationTimeService instance
func NewPreparationTimeService(preparationTimeRepo PreparationTimeRepository) *PreparationTimeService {
	return &PreparationTimeService{
		repo: preparationTimeRepo,
	}
}

func (s PreparationTimeService) FindAll() ([]m.PreparationTimeDTO, error) {

	preparationTimes, err := s.repo.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	result := m.PreparationTime{}.ConvertAllToDTO(preparationTimes)
	return result, nil
}

func (s PreparationTimeService) FindSingle(preparationTimeDTO m.PreparationTimeDTO) (m.PreparationTimeDTO, error) {

	preparationTime, err := s.repo.FindSingle(preparationTimeDTO.ConvertFromDTO())
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.PreparationTimeDTO{}, err
		default:
			return m.PreparationTimeDTO{}, errors.New("internal server error")
		}
	}

	return preparationTime.ConvertToDTO(), nil
}

func (s PreparationTimeService) Create(preparationTimeDTO m.PreparationTimeDTO) (m.PreparationTimeDTO, error) {

	if preparationTimeDTO.ID != uuid.Nil {
		return m.PreparationTimeDTO{}, errors.New("existing id on new element is not allowed")
	}

	if preparationTimeDTO.Duration == 0 {
		return m.PreparationTimeDTO{}, errors.New("name is empty")
	}

	created, err := s.repo.Create(preparationTimeDTO.ConvertFromDTO())
	if err != nil {
		return m.PreparationTimeDTO{}, err
	}

	return created.ConvertToDTO(), nil
}

func (s PreparationTimeService) Update(preparationTimeDTO m.PreparationTimeDTO) (m.PreparationTimeDTO, error) {

	if preparationTimeDTO.Duration == 0 {
		return m.PreparationTimeDTO{}, errors.New("name is empty")
	}

	updatedPreparationTime, err := s.repo.Update(preparationTimeDTO.ConvertFromDTO())
	if err != nil {
		return m.PreparationTimeDTO{}, err
	}

	return updatedPreparationTime.ConvertToDTO(), nil
}

func (s PreparationTimeService) Delete(preparationTimeDTO m.PreparationTimeDTO) error {

	err := s.repo.Delete(preparationTimeDTO.ConvertFromDTO())
	if err != nil {
		return err
	}

	return nil
}
