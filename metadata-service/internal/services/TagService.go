package services

import (
	"errors"

	m "metadata-service/internal/models"
)

type TagRepository interface {
	FindAll() ([]m.Tag, error)
	FindSingle(recipeID uint) (m.Tag, error)
	Create(recipe m.Tag) (m.Tag, error)
	Update(recipe m.Tag) (m.Tag, error)
	Delete(recipe m.Tag) error
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

func (s TagService) FindAll() ([]m.Tag, error) {

	tags, err := s.repo.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	return tags, nil
}

func (s TagService) FindSingle(tagID uint) (m.Tag, error) {

	tag, err := s.repo.FindSingle(tagID)
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.Tag{}, err
		default:
			return m.Tag{}, errors.New("internal server error")
		}
	}

	return tag, nil
}

func (s TagService) Create(tag m.Tag) (m.Tag, error) {

	if tag.ID != 0 {
		return m.Tag{}, errors.New("existing id on new element is not allowed")
	}

	if tag.TagName == "" {
		return m.Tag{}, errors.New("tagname is empty")
	}

	created, err := s.repo.Create(tag)
	if err != nil {
		return m.Tag{}, err
	}

	return created, nil
}

func (s TagService) Update(tag m.Tag, tagID uint) (m.Tag, error) {

	if tagID == 0 {
		return m.Tag{}, errors.New("missing id of element to update")
	}

	if tag.ID != tagID {
		tag.ID = tagID
	}

	if tag.TagName == "" {
		return m.Tag{}, errors.New("tagname is empty")
	}

	updatedTag, err := s.repo.Update(tag)
	if err != nil {
		return m.Tag{}, err
	}

	return updatedTag, nil
}

func (s TagService) Delete(tagID uint) error {
	var tag m.Tag

	if tagID == 0 {
		return errors.New("missing id of element to delete")
	}

	tag.ID = tagID

	err := s.repo.Delete(tag)
	if err != nil {
		return err
	}

	return nil
}
