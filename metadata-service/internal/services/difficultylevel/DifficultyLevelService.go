package services

import (
	"errors"

	m "metadata-service/internal/models"

	"github.com/google/uuid"
)

type DifficultyLevelRepository interface {
	FindAll() ([]m.DifficultyLevel, error)
	FindSingle(difficultyLevel m.DifficultyLevel) (m.DifficultyLevel, error)
	Create(difficultyLevel m.DifficultyLevel) (m.DifficultyLevel, error)
	Update(difficultyLevel m.DifficultyLevel) (m.DifficultyLevel, error)
	Delete(difficultyLevel m.DifficultyLevel) error
}
type DifficultyLevelService struct {
	repo DifficultyLevelRepository
}

// NewDifficultyLevelService creates a new DifficultyLevelService instance
func NewDifficultyLevelService(difficultyLevelRepo DifficultyLevelRepository) *DifficultyLevelService {
	return &DifficultyLevelService{
		repo: difficultyLevelRepo,
	}
}

func (s DifficultyLevelService) FindAll() ([]m.DifficultyLevelDTO, error) {

	difficultyLevels, err := s.repo.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	result := m.DifficultyLevel{}.ConvertAllToDTO(difficultyLevels)
	return result, nil
}

func (s DifficultyLevelService) FindSingle(difficultyLevelDTO m.DifficultyLevelDTO) (m.DifficultyLevelDTO, error) {

	difficultyLevel, err := s.repo.FindSingle(difficultyLevelDTO.ConvertFromDTO())
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.DifficultyLevelDTO{}, err
		default:
			return m.DifficultyLevelDTO{}, errors.New("internal server error")
		}
	}

	return difficultyLevel.ConvertToDTO(), nil
}

func (s DifficultyLevelService) Create(difficultyLevelDTO m.DifficultyLevelDTO) (m.DifficultyLevelDTO, error) {

	if difficultyLevelDTO.ID != uuid.Nil {
		return m.DifficultyLevelDTO{}, errors.New("existing id on new element is not allowed")
	}

	if difficultyLevelDTO.Level == 0 {
		return m.DifficultyLevelDTO{}, errors.New("name is empty")
	}

	created, err := s.repo.Create(difficultyLevelDTO.ConvertFromDTO())
	if err != nil {
		return m.DifficultyLevelDTO{}, err
	}

	return created.ConvertToDTO(), nil
}

func (s DifficultyLevelService) Update(difficultyLevelDTO m.DifficultyLevelDTO) (m.DifficultyLevelDTO, error) {

	if difficultyLevelDTO.Level == 0 {
		return m.DifficultyLevelDTO{}, errors.New("name is empty")
	}

	updatedDifficultyLevel, err := s.repo.Update(difficultyLevelDTO.ConvertFromDTO())
	if err != nil {
		return m.DifficultyLevelDTO{}, err
	}

	return updatedDifficultyLevel.ConvertToDTO(), nil
}

func (s DifficultyLevelService) Delete(difficultyLevelDTO m.DifficultyLevelDTO) error {

	err := s.repo.Delete(difficultyLevelDTO.ConvertFromDTO())
	if err != nil {
		return err
	}

	return nil
}
