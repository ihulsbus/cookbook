package services

import (
	"errors"
	"testing"

	m "ingredient-service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	findAllUnit m.Unit = m.Unit{
		ID:        uuid.New(),
		FullName:  "unit",
		ShortName: "u",
	}
	unit m.Unit = m.Unit{
		ID:        uuid.New(),
		FullName:  "unit",
		ShortName: "u",
	}
)

type UnitRepositoryMock struct{}

func (UnitRepositoryMock) FindAll() ([]m.Unit, error) {
	switch findAllUnit.FullName {
	case "findall": // OK
		var units []m.Unit
		units = append(units, unit)
		return units, nil
	case "notfound":
		return nil, errors.New("not found")
	default: // ERR
		return nil, errors.New("error")
	}
}

func (UnitRepositoryMock) FindSingle(unitInput m.Unit) (m.Unit, error) {
	switch unitInput.FullName {
	case "find":
		return unit, nil
	case "update":
		return unit, nil
	case "updateerror":
		return unit, nil
	case "delete":
		return unit, nil
	case "deleteerror":
		return unit, nil
	case "":
		return unit, nil
	case "notfound":
		return m.Unit{}, errors.New("not found")
	default:
		return m.Unit{}, errors.New("error")
	}
}

func (UnitRepositoryMock) Create(unitInput m.Unit) (m.Unit, error) {
	switch unitInput.FullName {
	case "create":
		return unit, nil
	default:
		return m.Unit{}, errors.New("error")
	}
}

func (UnitRepositoryMock) Update(unitInput m.Unit) (m.Unit, error) {
	switch unitInput.FullName {
	case "update":
		return unit, nil
	case "unit":
		return unit, nil
	default:
		return m.Unit{}, errors.New("error")
	}
}

func (UnitRepositoryMock) Delete(unitInput m.Unit) error {
	switch unitInput.FullName {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

func TestUnitFindAll_OK(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	findAllUnit.FullName = "findall"

	result, err := s.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "unit", result[0].FullName)
}

func TestUnitFindAll_Err(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	findAllUnit.FullName = "error"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "internal server error")
}

func TestUnitFindAll_NotFound(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	findAllUnit.FullName = "notfound"

	result, err := s.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
}

func TestUnitFindSingle_OK(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		ID:       unit.ID,
		FullName: "find",
	}
	result, err := s.FindSingle(unitDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.UnitDTO{}, result)
	assert.Equal(t, "unit", result.FullName)
	assert.Equal(t, unit.ID, result.ID)
}

func TestUnitFindSingle_FindErr(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		ID:       unit.ID,
		FullName: "error",
	}
	result, err := s.FindSingle(unitDTO)

	assert.Error(t, err)
	assert.IsType(t, m.UnitDTO{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestUnitFindSingle_NotFoundErr(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		ID:       unit.ID,
		FullName: "notfound",
	}
	result, err := s.FindSingle(unitDTO)

	assert.Error(t, err)
	assert.IsType(t, m.UnitDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestUnitCreate_OK(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		FullName: "create",
	}
	result, err := s.Create(unitDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.UnitDTO{}, result)
	assert.Equal(t, "unit", result.FullName)
}

func TestUnitCreate_IDErr(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		ID:       unit.ID,
		FullName: "create",
	}
	result, err := s.Create(unitDTO)

	assert.Error(t, err)
	assert.IsType(t, m.UnitDTO{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestUnitCreate_ExistsErr(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		FullName: "find",
	}
	result, err := s.Create(unitDTO)

	assert.Error(t, err)
	assert.IsType(t, m.UnitDTO{}, result)
	assert.EqualError(t, err, "unit already exists")
}

func TestUnitCreate_Err(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		FullName: "error",
	}
	result, err := s.Create(unitDTO)

	assert.Error(t, err)
	assert.IsType(t, m.UnitDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestUnitCreate_NoName(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		FullName: "",
	}
	result, err := s.Create(unitDTO)

	assert.Error(t, err)
	assert.IsType(t, m.UnitDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestUnitUpdate_Ok(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		ID:       unit.ID,
		FullName: "update",
	}
	result, err := s.Update(unitDTO)

	assert.NoError(t, err)
	assert.IsType(t, m.UnitDTO{}, result)
	assert.Equal(t, result.FullName, "unit")
}

func TestUnitUpdate_NotFoundErr(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		ID:       unit.ID,
		FullName: "notfound",
	}
	result, err := s.Update(unitDTO)

	assert.Error(t, err)
	assert.IsType(t, m.UnitDTO{}, result)
	assert.EqualError(t, err, "unit does not exist. nothing to update")
}

func TestUnitUpdate_Err(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		ID:       unit.ID,
		FullName: "updateerror",
	}
	result, err := s.Update(unitDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, m.UnitDTO{}, result)
}

func TestUnitDelete_Ok(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		ID:       unit.ID,
		FullName: "delete",
	}
	err := s.Delete(unitDTO)

	assert.NoError(t, err)
}

func TestUnitDelete_NotFoundErr(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		ID:       unit.ID,
		FullName: "notfound",
	}
	err := s.Delete(unitDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "unit does not exist. nothing to delete")
}

func TestUnitDelete_Err(t *testing.T) {
	s := NewUnitService(&UnitRepositoryMock{})

	unitDTO := m.UnitDTO{
		ID:       unit.ID,
		FullName: "deleteerror",
	}
	err := s.Delete(unitDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
