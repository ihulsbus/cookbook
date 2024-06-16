package services

import (
	"errors"
	"testing"

	m "metadata-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	findAllTag m.Tag = m.Tag{
		ID:   uuid.New(),
		Name: "tag",
	}
	tag m.Tag = m.Tag{
		ID:   uuid.New(),
		Name: "tag",
	}
)

type TagRepositoryMock struct{}

func (*TagRepositoryMock) FindAll() ([]m.Tag, error) {
	switch findAllTag.Name {
	case "findall":
		var tags []m.Tag
		tags = append(tags, findAllTag)
		return tags, nil
	case "not found":
		return nil, errors.New("not found")
	default:
		return nil, errors.New("error")
	}
}

func (*TagRepositoryMock) FindSingle(tag m.Tag) (m.Tag, error) {
	switch tag.Name {
	case "find":
		return tag, nil
	case "not found":
		return m.Tag{}, errors.New("not found")
	default:
		return m.Tag{}, errors.New("error")
	}
}

func (*TagRepositoryMock) Create(tag m.Tag) (m.Tag, error) {
	switch tag.Name {
	case "create":
		return tag, nil
	default:
		return m.Tag{}, errors.New("error")
	}
}

func (*TagRepositoryMock) Update(tag m.Tag) (m.Tag, error) {
	switch tag.Name {
	case "update":
		return tag, nil
	default:
		return m.Tag{}, errors.New("error")
	}
}

func (*TagRepositoryMock) Delete(tag m.Tag) error {
	switch tag.Name {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

// ======================================================================

func TestTagFindAll_OK(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})
	findAllTag.Name = "findall"

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.IsType(t, []m.TagDTO{}, result)
	assert.Len(t, result, 1)
}

func TestTagFindAll_Err(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})
	findAllTag.Name = "fail"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "internal server error")
}

func TestTagFindAll_NotFound(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})
	findAllTag.Name = "not found"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
}

func TestTagFindSingle_OK(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	tagDTO := m.TagDTO{
		ID:   tag.ID,
		Name: "find",
	}
	result, err := s.FindSingle(tagDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.TagDTO{}, result)
	assert.Equal(t, "find", result.Name)
	assert.Equal(t, tag.ID, result.ID)
}

func TestTagFindSingle_Err(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	tagDTO := m.TagDTO{
		ID:   tag.ID,
		Name: "error",
	}
	result, err := s.FindSingle(tagDTO)

	assert.Error(t, err)
	assert.IsType(t, m.TagDTO{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestTagFindSingle_NotFound(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	tagDTO := m.TagDTO{
		ID:   tag.ID,
		Name: "not found",
	}
	result, err := s.FindSingle(tagDTO)

	assert.Error(t, err)
	assert.IsType(t, m.TagDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestTagCreate_OK(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	tagDTO := m.TagDTO{
		Name: "create",
	}
	result, err := s.Create(tagDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.TagDTO{}, result)
}

func TestTagCreate_IDErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	tagDTO := m.TagDTO{
		ID:   tag.ID,
		Name: "create",
	}
	result, err := s.Create(tagDTO)

	assert.Error(t, err)
	assert.IsType(t, m.TagDTO{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestTagCreate_NoNameErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	tagDTO := m.TagDTO{
		Name: "",
	}
	result, err := s.Create(tagDTO)

	assert.Error(t, err)
	assert.IsType(t, m.TagDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestTagCreate_CreateErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	tagDTO := m.TagDTO{
		Name: "error",
	}
	result, err := s.Create(tagDTO)

	assert.Error(t, err)
	assert.IsType(t, m.TagDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestTagUpdate_OK(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	tagDTO := m.TagDTO{
		ID:   tag.ID,
		Name: "update",
	}
	result, err := s.Update(tagDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.TagDTO{}, result)
}

func TestTagUpdate_NoNameErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	tagDTO := m.TagDTO{
		ID:   tag.ID,
		Name: "",
	}
	result, err := s.Update(tagDTO)

	assert.Error(t, err)
	assert.IsType(t, m.TagDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestTagUpdate_UpdateErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	tagDTO := m.TagDTO{
		ID:   tag.ID,
		Name: "error",
	}
	result, err := s.Update(tagDTO)

	assert.Error(t, err)
	assert.IsType(t, m.TagDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestTagDelete_OK(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	tagDTO := m.TagDTO{
		ID:   tag.ID,
		Name: "delete",
	}
	err := s.Delete(tagDTO)

	assert.NoError(t, err)
}

func TestTagDelete_DeleteErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	tagDTO := m.TagDTO{
		ID:   tag.ID,
		Name: "error",
	}
	err := s.Delete(tagDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
