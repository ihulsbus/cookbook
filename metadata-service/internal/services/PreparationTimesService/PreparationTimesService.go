package services

import (
	"errors"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
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

func (s PreparationTimeService) FindAll() ([]m.PreparationTime, error) {

	preparationTimes, err := s.repo.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	return preparationTimes, nil
}

func (s PreparationTimeService) FindSingle(preparationTime m.PreparationTime) (m.PreparationTime, error) {

	preparationTime, err := s.repo.FindSingle(preparationTime)
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.PreparationTime{}, err
		default:
			return m.PreparationTime{}, errors.New("internal server error")
		}
	}

	return preparationTime, nil
}

func (s PreparationTimeService) Create(preparationTime m.PreparationTime) (m.PreparationTime, error) {

	if preparationTime.ID != uuid.Nil {
		return m.PreparationTime{}, errors.New("existing id on new element is not allowed")
	}

	if preparationTime.Duration == 0 {
		return m.PreparationTime{}, errors.New("name is empty")
	}

	created, err := s.repo.Create(preparationTime)
	if err != nil {
		return m.PreparationTime{}, err
	}

	return created, nil
}

func (s PreparationTimeService) Update(preparationTime m.PreparationTime) (m.PreparationTime, error) {

	if preparationTime.Duration == 0 {
		return m.PreparationTime{}, errors.New("name is empty")
	}

	updatedPreparationTime, err := s.repo.Update(preparationTime)
	if err != nil {
		return m.PreparationTime{}, err
	}

	return updatedPreparationTime, nil
}

func (s PreparationTimeService) Delete(preparationTime m.PreparationTime) error {

	err := s.repo.Delete(preparationTime)
	if err != nil {
		return err
	}

	return nil
}
