package services

import (
	"errors"
	"testing"

	m "metadata-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	findAllDifficultyLevel m.DifficultyLevel = m.DifficultyLevel{
		ID:    uuid.New(),
		Level: 1,
	}
	difficultyLevel m.DifficultyLevel = m.DifficultyLevel{
		ID:    uuid.New(),
		Level: 1,
	}
)

type DifficultyLevelRepositoryMock struct{}

func (*DifficultyLevelRepositoryMock) FindAll() ([]m.DifficultyLevel, error) {
	switch findAllDifficultyLevel.Level {
	case 1:
		var difficultyLevels []m.DifficultyLevel
		difficultyLevels = append(difficultyLevels, findAllDifficultyLevel)
		return difficultyLevels, nil
	case 2:
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (*DifficultyLevelRepositoryMock) FindSingle(difficultyLevel m.DifficultyLevel) (m.DifficultyLevel, error) {
	switch difficultyLevel.Level {
	case 1:
		return difficultyLevel, nil
	case 2:
		return m.DifficultyLevel{}, errors.New("not found")
	default:
		return m.DifficultyLevel{}, errors.New("error")
	}
}

func (*DifficultyLevelRepositoryMock) Create(difficultyLevel m.DifficultyLevel) (m.DifficultyLevel, error) {
	switch difficultyLevel.Level {
	case 1:
		return difficultyLevel, nil
	default:
		return m.DifficultyLevel{}, errors.New("error")
	}
}

func (*DifficultyLevelRepositoryMock) Update(difficultyLevel m.DifficultyLevel) (m.DifficultyLevel, error) {
	switch difficultyLevel.Level {
	case 1:
		return difficultyLevel, nil
	default:
		return m.DifficultyLevel{}, errors.New("error")
	}
}

func (*DifficultyLevelRepositoryMock) Delete(difficultyLevel m.DifficultyLevel) error {
	switch difficultyLevel.Level {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

// ======================================================================

func TestDifficultyLevelFindAll_OK(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})
	findAllDifficultyLevel.Level = 1

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.IsType(t, []m.DifficultyLevelDTO{}, result)
	assert.Len(t, result, 1)
}

func TestDifficultyLevelFindAll_Err(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})
	findAllDifficultyLevel.Level = 3

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "internal server error")
}

func TestDifficultyLevelFindAll_NotFound(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})
	findAllDifficultyLevel.Level = 2

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
}

func TestDifficultyLevelFindSingle_OK(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})

	difficultyLevelDTO := m.DifficultyLevelDTO{
		ID:    difficultyLevel.ID,
		Level: 1,
	}
	result, err := s.FindSingle(difficultyLevelDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.DifficultyLevelDTO{}, result)
	assert.Equal(t, 1, result.Level)
	assert.Equal(t, difficultyLevel.ID, result.ID)
}

func TestDifficultyLevelFindSingle_Err(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})

	difficultyLevelDTO := m.DifficultyLevelDTO{
		ID:    difficultyLevel.ID,
		Level: 3,
	}
	result, err := s.FindSingle(difficultyLevelDTO)

	assert.Error(t, err)
	assert.IsType(t, m.DifficultyLevelDTO{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestDifficultyLevelFindSingle_NotFound(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})

	difficultyLevelDTO := m.DifficultyLevelDTO{
		ID:    difficultyLevel.ID,
		Level: 2,
	}
	result, err := s.FindSingle(difficultyLevelDTO)

	assert.Error(t, err)
	assert.IsType(t, m.DifficultyLevelDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestDifficultyLevelCreate_OK(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})

	difficultyLevelDTO := m.DifficultyLevelDTO{
		Level: 1,
	}
	result, err := s.Create(difficultyLevelDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.DifficultyLevelDTO{}, result)
}

func TestDifficultyLevelCreate_IDErr(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})

	difficultyLevelDTO := m.DifficultyLevelDTO{
		ID:    difficultyLevel.ID,
		Level: 1,
	}
	result, err := s.Create(difficultyLevelDTO)

	assert.Error(t, err)
	assert.IsType(t, m.DifficultyLevelDTO{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestDifficultyLevelCreate_NoLevelErr(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})

	difficultyLevelDTO := m.DifficultyLevelDTO{
		Level: 0,
	}
	result, err := s.Create(difficultyLevelDTO)

	assert.Error(t, err)
	assert.IsType(t, m.DifficultyLevelDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestDifficultyLevelCreate_CreateErr(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})

	difficultyLevelDTO := m.DifficultyLevelDTO{
		Level: 3,
	}
	result, err := s.Create(difficultyLevelDTO)

	assert.Error(t, err)
	assert.IsType(t, m.DifficultyLevelDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestDifficultyLevelUpdate_OK(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})

	difficultyLevelDTO := m.DifficultyLevelDTO{
		ID:    difficultyLevel.ID,
		Level: 1,
	}
	result, err := s.Update(difficultyLevelDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.DifficultyLevelDTO{}, result)
}

func TestDifficultyLevelUpdate_NoLevelErr(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})

	difficultyLevelDTO := m.DifficultyLevelDTO{
		ID:    difficultyLevel.ID,
		Level: 0,
	}
	result, err := s.Update(difficultyLevelDTO)

	assert.Error(t, err)
	assert.IsType(t, m.DifficultyLevelDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestDifficultyLevelUpdate_UpdateErr(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})

	difficultyLevelDTO := m.DifficultyLevelDTO{
		ID:    difficultyLevel.ID,
		Level: 3,
	}
	result, err := s.Update(difficultyLevelDTO)

	assert.Error(t, err)
	assert.IsType(t, m.DifficultyLevelDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestDifficultyLevelDelete_OK(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})

	difficultyLevelDTO := m.DifficultyLevelDTO{
		ID:    difficultyLevel.ID,
		Level: 1,
	}
	err := s.Delete(difficultyLevelDTO)

	assert.NoError(t, err)
}

func TestDifficultyLevelDelete_DeleteErr(t *testing.T) {
	s := NewDifficultyLevelService(&DifficultyLevelRepositoryMock{})

	difficultyLevelDTO := m.DifficultyLevelDTO{
		ID:    difficultyLevel.ID,
		Level: 3,
	}
	err := s.Delete(difficultyLevelDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
