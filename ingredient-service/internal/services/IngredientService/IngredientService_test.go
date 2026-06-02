package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ihulsbus/cookbook/shared/models"
	"github.com/stretchr/testify/assert"
)

var (
	findAllIngredient models.Ingredient = models.Ingredient{
		ID:   uuid.New(),
		Name: "ingredient",
	}
	ingredient models.Ingredient = models.Ingredient{
		ID:   uuid.New(),
		Name: "ingredient",
	}

	findAllUnit models.Unit = models.Unit{
		ID:        uuid.New(),
		FullName:  "unit",
		ShortName: "u",
	}
	unit models.Unit = models.Unit{
		ID:        uuid.New(),
		FullName:  "unit",
		ShortName: "u",
	}
)

type IngredientRepositoryMock struct{}

func (IngredientRepositoryMock) FindAll(pagination models.PaginationRequest) ([]models.Ingredient, int64, error) {
	switch findAllIngredient.Name {
	case "findall":
		return []models.Ingredient{ingredient}, 1, nil
	case "notfound":
		return nil, 0, errors.New("not found")
	default:
		return nil, 0, errors.New("error")
	}
}

func (IngredientRepositoryMock) FindByName(name string) (models.Ingredient, error) {
	switch name {
	case "ingredient":
		return ingredient, nil
	default:
		return models.Ingredient{}, nil
	}
}

func (IngredientRepositoryMock) FindSingle(ingredientInput models.Ingredient) (models.Ingredient, error) {
	switch ingredientInput.Name {
	case "find":
		return ingredient, nil
	case "update":
		return ingredient, nil
	case "updateerror":
		return ingredient, nil
	case "delete":
		return ingredient, nil
	case "deleteerror":
		return ingredient, nil
	case "":
		return ingredient, nil
	case "notfound":
		return models.Ingredient{}, errors.New("not found")
	default:
		return models.Ingredient{}, errors.New("error")
	}
}

func (IngredientRepositoryMock) Create(ingredientInput models.Ingredient) (models.Ingredient, error) {
	switch ingredientInput.Name {
	case "create":
		return ingredient, nil
	default:
		return models.Ingredient{}, errors.New("error")
	}
}

func (IngredientRepositoryMock) Update(ingredientInput models.Ingredient) (models.Ingredient, error) {
	switch ingredientInput.Name {
	case "update":
		return ingredient, nil
	case "ingredient":
		return ingredient, nil
	default:
		return models.Ingredient{}, errors.New("error")
	}
}

func (IngredientRepositoryMock) Delete(ingredientInput models.Ingredient) error {
	switch ingredientInput.Name {
	case "delete":
		return nil
	default:
		return errors.New("error")
	}
}

func TestIngredientFindAll_OK(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	findAllIngredient.Name = "findall"

	result, err := s.FindAll(models.NormalizePagination(1, 25))

	assert.NoError(t, err)
	assert.IsType(t, models.PaginatedResponse[models.IngredientDTO]{}, result)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, "ingredient", result.Data[0].Name)
}

func TestIngredientFindAll_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	findAllIngredient.Name = "error"

	result, err := s.FindAll(models.NormalizePagination(1, 25))

	assert.Error(t, err)
	assert.Equal(t, models.PaginatedResponse[models.IngredientDTO]{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestIngredientFindAll_NotFound(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	findAllIngredient.Name = "notfound"

	result, err := s.FindAll(models.NormalizePagination(1, 25))

	assert.Error(t, err)
	assert.Equal(t, models.PaginatedResponse[models.IngredientDTO]{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestIngredientFindSingle_OK(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		ID:   ingredient.ID,
		Name: "find",
	}
	result, err := s.FindSingle(ingredientDTO)

	assert.NoError(t, err)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.Equal(t, "ingredient", result.Name)
	assert.Equal(t, ingredient.ID, result.ID)
}

func TestIngredientFindSingle_FindErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		ID:   ingredient.ID,
		Name: "error",
	}
	result, err := s.FindSingle(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestIngredientFindSingle_NotFoundErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		ID:   ingredient.ID,
		Name: "notfound",
	}
	result, err := s.FindSingle(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestIngredientCreate_OK(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		Name: "create",
	}
	result, err := s.Create(ingredientDTO)

	assert.NoError(t, err)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.Equal(t, "ingredient", result.Name)
}

func TestIngredientCreate_IDErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		ID:   ingredient.ID,
		Name: "create",
	}
	result, err := s.Create(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestIngredientCreate_ExistsErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		Name: "find",
	}
	result, err := s.Create(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.EqualError(t, err, "ingredient already exists")
}

func TestIngredientCreate_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		Name: "error",
	}
	result, err := s.Create(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestIngredientCreate_NoName(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		Name: "",
	}
	result, err := s.Create(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestIngredientUpdate_Ok(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		ID:   ingredient.ID,
		Name: "update",
	}
	result, err := s.Update(ingredientDTO)

	assert.NoError(t, err)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.Equal(t, result.Name, "ingredient")
}

func TestIngredientUpdate_NotFoundErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		ID:   ingredient.ID,
		Name: "notfound",
	}
	result, err := s.Update(ingredientDTO)

	assert.Error(t, err)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.EqualError(t, err, "ingredient does not exist. nothing to update")
}

func TestIngredientUpdate_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		ID:   ingredient.ID,
		Name: "updateerror",
	}
	result, err := s.Update(ingredientDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
	assert.IsType(t, models.IngredientDTO{}, result)
}

func TestIngredientDelete_Ok(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		ID:   ingredient.ID,
		Name: "delete",
	}
	err := s.Delete(ingredientDTO)

	assert.NoError(t, err)
}

func TestIngredientDelete_NotFoundErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		ID:   ingredient.ID,
		Name: "notfound",
	}
	err := s.Delete(ingredientDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "ingredient does not exist. nothing to delete")
}

func TestIngredientDelete_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		ID:   ingredient.ID,
		Name: "deleteerror",
	}
	err := s.Delete(ingredientDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
