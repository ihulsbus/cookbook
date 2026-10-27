package services

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
)

var (
	findAllPreparationTime m.PreparationTime = m.PreparationTime{
		ID:       uuid.New(),
		Duration: 1 * time.Minute,
	}
	preparationTime m.PreparationTime = m.PreparationTime{
		ID:       uuid.New(),
		Duration: 1 * time.Minute,
	}
)

type PreparationTimeRepositoryMock struct{}

func (*PreparationTimeRepositoryMock) FindAll(pagination m.PaginationRequest) ([]m.PreparationTime, int64, error) {
	switch findAllPreparationTime.Duration {
	case 1 * time.Minute:
		var preparationTimes []m.PreparationTime
		preparationTimes = append(preparationTimes, findAllPreparationTime)
		return preparationTimes, 1, nil
	default:
		return nil, 0, errors.New("error")
	}
}

func (*PreparationTimeRepositoryMock) FindSingle(preparationTime m.PreparationTime) (m.PreparationTime, error) {
	switch preparationTime.Duration {
	case 1 * time.Minute:
		return preparationTime, nil
	case 2 * time.Minute:
		return m.PreparationTime{}, errors.New("not found")
	default:
		return m.PreparationTime{}, errors.New("error")
	}
}

func (*PreparationTimeRepositoryMock) Create(preparationTime m.PreparationTime) (m.PreparationTime, error) {
	switch preparationTime.Duration {
	case 1 * time.Minute:
		return preparationTime, nil
	default:
		return m.PreparationTime{}, errors.New("error")
	}
}

func (*PreparationTimeRepositoryMock) Update(preparationTime m.PreparationTime) (m.PreparationTime, error) {
	switch preparationTime.Duration {
	case 1 * time.Minute:
		return preparationTime, nil
	default:
		return m.PreparationTime{}, errors.New("error")
	}
}

func (*PreparationTimeRepositoryMock) Delete(preparationTime m.PreparationTime) error {
	switch preparationTime.Duration {
	case 1 * time.Minute:
		return nil
	default:
		return errors.New("error")
	}
}

// ======================================================================

func TestPreparationTimeFindAll_OK(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})
	findAllPreparationTime.Duration = 1 * time.Minute

	result, err := s.FindAll(m.NormalizePagination(1, 25))

	assert.NoError(t, err)
	assert.IsType(t, m.PaginatedResponse[m.PreparationTime]{}, result)
	assert.Len(t, result.Data, 1)
}

func TestPreparationTimeFindAll_Err(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})
	findAllPreparationTime.Duration = 3 * time.Minute

	result, err := s.FindAll(m.NormalizePagination(1, 25))

	assert.Error(t, err)
	assert.Equal(t, m.PaginatedResponse[m.PreparationTime]{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestPreparationTimeFindSingle_OK(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	result, err := s.FindSingle(preparationTime)

	assert.NoError(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.Equal(t, preparationTime.Duration, result.Duration)
	assert.Equal(t, preparationTime.ID, result.ID)
}

func TestPreparationTimeFindSingle_Err(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTime := m.PreparationTime{
		ID:       preparationTime.ID,
		Duration: 3 * time.Minute,
	}
	result, err := s.FindSingle(preparationTime)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestPreparationTimeFindSingle_NotFound(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTime := m.PreparationTime{
		ID:       preparationTime.ID,
		Duration: 2 * time.Minute,
	}
	result, err := s.FindSingle(preparationTime)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.EqualError(t, err, "not found")
}

func TestPreparationTimeCreate_OK(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTime := m.PreparationTime{
		Duration: 1 * time.Minute,
	}
	result, err := s.Create(preparationTime)

	assert.NoError(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
}

func TestPreparationTimeCreate_IDErr(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTime := m.PreparationTime{
		ID:       preparationTime.ID,
		Duration: 1 * time.Minute,
	}
	result, err := s.Create(preparationTime)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestPreparationTimeCreate_NoDurationErr(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTime := m.PreparationTime{
		Duration: 0 * time.Minute,
	}
	result, err := s.Create(preparationTime)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestPreparationTimeCreate_CreateErr(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTime := m.PreparationTime{
		Duration: 3 * time.Minute,
	}
	result, err := s.Create(preparationTime)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.EqualError(t, err, "error")
}

func TestPreparationTimeUpdate_OK(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTime := m.PreparationTime{
		ID:       preparationTime.ID,
		Duration: 1 * time.Minute,
	}
	result, err := s.Update(preparationTime)

	assert.NoError(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
}

func TestPreparationTimeUpdate_NoDurationErr(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTime := m.PreparationTime{
		ID:       preparationTime.ID,
		Duration: 0 * time.Minute,
	}
	result, err := s.Update(preparationTime)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestPreparationTimeUpdate_UpdateErr(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTime := m.PreparationTime{
		ID:       preparationTime.ID,
		Duration: 3 * time.Minute,
	}
	result, err := s.Update(preparationTime)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTime{}, result)
	assert.EqualError(t, err, "error")
}

func TestPreparationTimeDelete_OK(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTime := m.PreparationTime{
		ID:       preparationTime.ID,
		Duration: 1 * time.Minute,
	}
	err := s.Delete(preparationTime)

	assert.NoError(t, err)
}

func TestPreparationTimeDelete_DeleteErr(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTime := m.PreparationTime{
		ID:       preparationTime.ID,
		Duration: 3 * time.Minute,
	}
	err := s.Delete(preparationTime)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
