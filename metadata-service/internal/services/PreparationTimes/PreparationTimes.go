package services

import (
	"errors"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
)

type PreparationTimeRepository interface {
	FindAll(pagination m.PaginationRequest) ([]m.PreparationTime, int64, error)
	FindSingle(preparationTime m.PreparationTime) (m.PreparationTime, error)
	Create(preparationTime m.PreparationTime) (m.PreparationTime, error)
	Update(preparationTime m.PreparationTime) (m.PreparationTime, error)
	Delete(preparationTime m.PreparationTime) error
}
type PreparationTime struct {
	repo PreparationTimeRepository
}

// NewPreparationTimeService creates a new PreparationTime instance
func NewPreparationTimeService(preparationTimeRepo PreparationTimeRepository) *PreparationTime {
	return &PreparationTime{
		repo: preparationTimeRepo,
	}
}

func (s PreparationTime) FindAll(pagination m.PaginationRequest) (m.PaginatedResponse[m.PreparationTime], error) {

	preparationTimes, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return m.PaginatedResponse[m.PreparationTime]{}, errors.New("internal server error")
	}

	return m.PaginatedResponse[m.PreparationTime]{
		Data:       preparationTimes,
		Pagination: m.NewPaginationMetadata(pagination.Page, pagination.Limit, total),
	}, nil
}

func (s PreparationTime) FindSingle(preparationTime m.PreparationTime) (m.PreparationTime, error) {

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

func (s PreparationTime) Create(preparationTime m.PreparationTime) (m.PreparationTime, error) {

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

func (s PreparationTime) Update(preparationTime m.PreparationTime) (m.PreparationTime, error) {

	if preparationTime.Duration == 0 {
		return m.PreparationTime{}, errors.New("name is empty")
	}

	updatedPreparationTime, err := s.repo.Update(preparationTime)
	if err != nil {
		return m.PreparationTime{}, err
	}

	return updatedPreparationTime, nil
}

func (s PreparationTime) Delete(preparationTime m.PreparationTime) error {

	err := s.repo.Delete(preparationTime)
	if err != nil {
		return err
	}

	return nil
}
