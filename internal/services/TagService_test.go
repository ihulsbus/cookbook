package services

import (
	"errors"
	"testing"

	m "github.com/ihulsbus/cookbook/internal/models"
	"github.com/stretchr/testify/assert"
)

var (
	findAllTag m.Tag = m.Tag{
		TagName: "tag",
	}
	tag m.Tag = m.Tag{
		TagName: "tag",
	}
)

type TagRepositoryMock struct{}

func (*TagRepositoryMock) FindAll() ([]m.Tag, error) {
	switch findAllTag.ID {
	case 1:
		var tags []m.Tag
		tags = append(tags, findAllTag)
		return tags, nil
	default:
		return nil, errors.New("error")
	}
}

func (*TagRepositoryMock) FindSingle(tagID uint) (m.Tag, error) {
	switch tagID {
	case 1:
		return tag, nil
	default:
		return m.Tag{}, errors.New("error")
	}
}

func (*TagRepositoryMock) Create(createTag m.Tag) (m.Tag, error) {
	switch createTag.TagName {
	case "create":
		createTag.ID = 1
		return createTag, nil
	default:
		return m.Tag{}, errors.New("error")
	}
}

func (*TagRepositoryMock) Update(tag m.Tag) (m.Tag, error) {
	switch tag.ID {
	case 1:
		return tag, nil
	default:
		return m.Tag{}, errors.New("error")
	}
}

func (*TagRepositoryMock) Delete(tag m.Tag) error {
	switch tag.ID {
	case 1:
		return nil
	default:
		return errors.New("error")
	}
}

// ======================================================================

func TestTagFindAll_OK(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})
	findAllTag.ID = 1

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.IsType(t, []m.Tag{}, result)
	assert.Len(t, result, 1)
}

func TestTagFindAll_err(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})
	findAllTag.ID = 2

	result, err := s.FindAll()

	findAllTag.ID = 1

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "error")
}

func TestTagFindSingle_OK(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	tag.ID = 1

	result, err := s.FindSingle(uint(1))

	assert.NoError(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.Equal(t, "tag", result.TagName)
	assert.Equal(t, uint(1), result.ID)
}

func TestTagFindSingle_Err(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	result, err := s.FindSingle(uint(2))

	assert.Error(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.EqualError(t, err, "error")
}

func TestTagCreate_OK(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	createTag := tag
	createTag.ID = 0
	createTag.TagName = "create"

	result, err := s.Create(createTag)

	assert.NoError(t, err)
	assert.IsType(t, m.Tag{}, result)
}

func TestTagCreate_IDErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	createTag := tag
	createTag.ID = 1
	createTag.TagName = "create"

	result, err := s.Create(createTag)

	assert.Error(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestTagCreate_NoNameErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	createTag := tag
	createTag.ID = 0
	createTag.TagName = ""

	result, err := s.Create(createTag)

	assert.Error(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.EqualError(t, err, "tagname is empty")
}

func TestTagCreate_CreateErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	createTag := tag
	createTag.ID = 0
	createTag.TagName = "fails"

	result, err := s.Create(createTag)

	assert.Error(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.EqualError(t, err, "error")
}

func TestTagUpdate_OK(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	updateTag := tag
	updateTag.ID = 1
	updateTag.TagName = "update"

	result, err := s.Update(updateTag, uint(1))

	assert.NoError(t, err)
	assert.IsType(t, m.Tag{}, result)
}

func TestTagUpdate_IDErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	updateTag := tag
	updateTag.ID = 0
	updateTag.TagName = "update"

	result, err := s.Update(updateTag, uint(0))

	assert.Error(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.EqualError(t, err, "missing id of element to update")
}

func TestTagUpdate_IDNotEqual(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	updateTag := tag
	updateTag.ID = 2
	updateTag.TagName = "update"

	result, err := s.Update(updateTag, uint(1))

	assert.NoError(t, err)
	assert.IsType(t, m.Tag{}, result)
}

func TestTagUpdate_NoNameErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	updateTag := tag
	updateTag.ID = 1
	updateTag.TagName = ""

	result, err := s.Update(updateTag, uint(1))

	assert.Error(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.EqualError(t, err, "tagname is empty")
}

func TestTagUpdate_UpdateErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	updateTag := tag
	updateTag.ID = 2
	updateTag.TagName = "fails"

	result, err := s.Update(updateTag, uint(2))

	assert.Error(t, err)
	assert.IsType(t, m.Tag{}, result)
	assert.EqualError(t, err, "error")
}

func TestTagDelete_OK(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	deleteTag := tag
	deleteTag.ID = 1
	deleteTag.TagName = "delete"

	err := s.Delete(deleteTag, uint(1))

	assert.NoError(t, err)
}

func TestTagDelete_IDErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	deleteTag := tag
	deleteTag.ID = 0
	deleteTag.TagName = "delete"

	err := s.Delete(deleteTag, uint(0))

	assert.Error(t, err)
	assert.EqualError(t, err, "missing id of element to delete")
}

func TestTagDelete_IDNotEqual(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	deleteTag := tag
	deleteTag.ID = 2
	deleteTag.TagName = "delete"

	err := s.Delete(deleteTag, uint(1))

	assert.NoError(t, err)
}

func TestTagDelete_DeleteErr(t *testing.T) {
	s := NewTagService(&TagRepositoryMock{})

	deleteTag := tag
	deleteTag.ID = 2
	deleteTag.TagName = "delete"

	err := s.Delete(deleteTag, uint(2))

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
