package services

import (
	"errors"
	"testing"

	m "metadata-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	findAllPreparationTime m.PreparationTime = m.PreparationTime{
		ID:       uuid.New(),
		Duration: 1,
	}
	preparationTime m.PreparationTime = m.PreparationTime{
		ID:       uuid.New(),
		Duration: 1,
	}
)

type PreparationTimeRepositoryMock struct{}

func (*PreparationTimeRepositoryMock) FindAll() ([]m.PreparationTime, error) {
	switch findAllPreparationTime.Duration {
	case 1:
		var preparationTimes []m.PreparationTime
		preparationTimes = append(preparationTimes, findAllPreparationTime)
		return preparationTimes, nil
	case 2:
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (*PreparationTimeRepositoryMock) FindSingle(preparationTime m.PreparationTime) (m.PreparationTime, error) {
	switch preparationTime.Duration {
	case 1:
		return preparationTime, nil
	case 2:
		return m.PreparationTime{}, errors.New("not found")
	default:
		return m.PreparationTime{}, errors.New("error")
	}
}

func (*PreparationTimeRepositoryMock) Create(preparationTime m.PreparationTime) (m.PreparationTime, error) {
	switch preparationTime.Duration {
	case 1:
		return preparationTime, nil
	default:
		return m.PreparationTime{}, errors.New("error")
	}
}

func (*PreparationTimeRepositoryMock) Update(preparationTime m.PreparationTime) (m.PreparationTime, error) {
	switch preparationTime.Duration {
	case 1:
		return preparationTime, nil
	default:
		return m.PreparationTime{}, errors.New("error")
	}
}

func (*PreparationTimeRepositoryMock) Delete(preparationTime m.PreparationTime) error {
	switch preparationTime.Duration {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

// ======================================================================

func TestPreparationTimeFindAll_OK(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})
	findAllPreparationTime.Duration = 1

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.IsType(t, []m.PreparationTimeDTO{}, result)
	assert.Len(t, result, 1)
}

func TestPreparationTimeFindAll_Err(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})
	findAllPreparationTime.Duration = 3

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "internal server error")
}

func TestPreparationTimeFindAll_NotFound(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})
	findAllPreparationTime.Duration = 2

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
}

func TestPreparationTimeFindSingle_OK(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTimeDTO := m.PreparationTimeDTO{
		ID:       preparationTime.ID,
		Duration: 1,
	}
	result, err := s.FindSingle(preparationTimeDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.PreparationTimeDTO{}, result)
	assert.Equal(t, 1, result.Duration)
	assert.Equal(t, preparationTime.ID, result.ID)
}

func TestPreparationTimeFindSingle_Err(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTimeDTO := m.PreparationTimeDTO{
		ID:       preparationTime.ID,
		Duration: 3,
	}
	result, err := s.FindSingle(preparationTimeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTimeDTO{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestPreparationTimeFindSingle_NotFound(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTimeDTO := m.PreparationTimeDTO{
		ID:       preparationTime.ID,
		Duration: 2,
	}
	result, err := s.FindSingle(preparationTimeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTimeDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestPreparationTimeCreate_OK(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTimeDTO := m.PreparationTimeDTO{
		Duration: 1,
	}
	result, err := s.Create(preparationTimeDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.PreparationTimeDTO{}, result)
}

func TestPreparationTimeCreate_IDErr(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTimeDTO := m.PreparationTimeDTO{
		ID:       preparationTime.ID,
		Duration: 1,
	}
	result, err := s.Create(preparationTimeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTimeDTO{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestPreparationTimeCreate_NoDurationErr(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTimeDTO := m.PreparationTimeDTO{
		Duration: 0,
	}
	result, err := s.Create(preparationTimeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTimeDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestPreparationTimeCreate_CreateErr(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTimeDTO := m.PreparationTimeDTO{
		Duration: 3,
	}
	result, err := s.Create(preparationTimeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTimeDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestPreparationTimeUpdate_OK(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTimeDTO := m.PreparationTimeDTO{
		ID:       preparationTime.ID,
		Duration: 1,
	}
	result, err := s.Update(preparationTimeDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.PreparationTimeDTO{}, result)
}

func TestPreparationTimeUpdate_NoDurationErr(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTimeDTO := m.PreparationTimeDTO{
		ID:       preparationTime.ID,
		Duration: 0,
	}
	result, err := s.Update(preparationTimeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTimeDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestPreparationTimeUpdate_UpdateErr(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTimeDTO := m.PreparationTimeDTO{
		ID:       preparationTime.ID,
		Duration: 3,
	}
	result, err := s.Update(preparationTimeDTO)

	assert.Error(t, err)
	assert.IsType(t, m.PreparationTimeDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestPreparationTimeDelete_OK(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTimeDTO := m.PreparationTimeDTO{
		ID:       preparationTime.ID,
		Duration: 1,
	}
	err := s.Delete(preparationTimeDTO)

	assert.NoError(t, err)
}

func TestPreparationTimeDelete_DeleteErr(t *testing.T) {
	s := NewPreparationTimeService(&PreparationTimeRepositoryMock{})

	preparationTimeDTO := m.PreparationTimeDTO{
		ID:       preparationTime.ID,
		Duration: 3,
	}
	err := s.Delete(preparationTimeDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
