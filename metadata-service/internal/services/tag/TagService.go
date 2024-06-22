package services

import (
	"errors"

	m "metadata-service/internal/models"

	"github.com/google/uuid"
)

type TagRepository interface {
	FindAll() ([]m.Tag, error)
	FindSingle(tag m.Tag) (m.Tag, error)
	Create(tag m.Tag) (m.Tag, error)
	Update(tag m.Tag) (m.Tag, error)
	Delete(tag m.Tag) error
}
type TagService struct {
	repo TagRepository
}

// NewTagService creates a new TagService instance
func NewTagService(tagRepo TagRepository) *TagService {
	return &TagService{
		repo: tagRepo,
	}
}

func (s TagService) FindAll() ([]m.TagDTO, error) {

	tags, err := s.repo.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	result := m.Tag{}.ConvertAllToDTO(tags)
	return result, nil
}

func (s TagService) FindSingle(tagDTO m.TagDTO) (m.TagDTO, error) {

	tag, err := s.repo.FindSingle(tagDTO.ConvertFromDTO())
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.TagDTO{}, err
		default:
			return m.TagDTO{}, errors.New("internal server error")
		}
	}

	return tag.ConvertToDTO(), nil
}

func (s TagService) Create(tagDTO m.TagDTO) (m.TagDTO, error) {

	if tagDTO.ID != uuid.Nil {
		return m.TagDTO{}, errors.New("existing id on new element is not allowed")
	}

	if tagDTO.Name == "" {
		return m.TagDTO{}, errors.New("name is empty")
	}

	created, err := s.repo.Create(tagDTO.ConvertFromDTO())
	if err != nil {
		return m.TagDTO{}, err
	}

	return created.ConvertToDTO(), nil
}

func (s TagService) Update(tagDTO m.TagDTO) (m.TagDTO, error) {

	if tagDTO.Name == "" {
		return m.TagDTO{}, errors.New("name is empty")
	}

	updatedTag, err := s.repo.Update(tagDTO.ConvertFromDTO())
	if err != nil {
		return m.TagDTO{}, err
	}

	return updatedTag.ConvertToDTO(), nil
}

func (s TagService) Delete(tagDTO m.TagDTO) error {

	err := s.repo.Delete(tagDTO.ConvertFromDTO())
	if err != nil {
		return err
	}

	return nil
}
