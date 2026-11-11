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
	case "find":
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

// Dispatches on name, defaulting to "nothing found" so any other test
// name (e.g. "create") passes through unless it opts into one of these triggers.
func (IngredientRepositoryMock) FindByNormalizedName(name string) (models.Ingredient, error) {
	switch name {
	case "duplicate":
		return ingredient, nil
	case "findbynormalizederror":
		return models.Ingredient{}, errors.New("db error")
	default:
		return models.Ingredient{}, nil
	}
}

func (IngredientRepositoryMock) FindSimilar(name string, limit int) ([]models.Ingredient, error) {
	switch name {
	case "similar":
		return []models.Ingredient{ingredient}, nil
	case "findsimilarerror":
		return nil, errors.New("db error")
	default:
		return nil, nil
	}
}

func (IngredientRepositoryMock) Create(ingredientInput models.Ingredient) (models.Ingredient, error) {
	switch ingredientInput.Name {
	case "create", "similar":
		return ingredient, nil
	case "dupecreate":
		return models.Ingredient{}, errors.New("duplicate")
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
	result, candidates, err := s.Create(ingredientDTO, uuid.New(), false)

	assert.NoError(t, err)
	assert.Nil(t, candidates)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.Equal(t, "ingredient", result.Name)
}

func TestIngredientCreate_IDErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		ID:   ingredient.ID,
		Name: "create",
	}
	result, _, err := s.Create(ingredientDTO, uuid.New(), false)

	assert.Error(t, err)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestIngredientCreate_NoName(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		Name: "",
	}
	result, _, err := s.Create(ingredientDTO, uuid.New(), false)

	assert.Error(t, err)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

// A normalized-name match always blocks, regardless of force, with the
// existing ingredient as the sole candidate.
func TestIngredientCreate_ExactDuplicate(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		Name: "duplicate",
	}
	result, candidates, err := s.Create(ingredientDTO, uuid.New(), true)

	assert.EqualError(t, err, "duplicate")
	assert.Equal(t, models.IngredientDTO{}, result)
	assert.Len(t, candidates, 1)
	assert.Equal(t, ingredient.ID, candidates[0].ID)
}

func TestIngredientCreate_FindByNormalizedNameErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		Name: "findbynormalizederror",
	}
	_, candidates, err := s.Create(ingredientDTO, uuid.New(), false)

	assert.EqualError(t, err, "internal server error")
	assert.Nil(t, candidates)
}

// A near-duplicate (pg_trgm similarity > 0.6) blocks unless force=true.
func TestIngredientCreate_SimilarBlocksWithoutForce(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		Name: "similar",
	}
	_, candidates, err := s.Create(ingredientDTO, uuid.New(), false)

	assert.EqualError(t, err, "similar ingredients exist")
	assert.Len(t, candidates, 1)
}

func TestIngredientCreate_SimilarForceBypasses(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		Name: "similar",
	}
	result, candidates, err := s.Create(ingredientDTO, uuid.New(), true)

	assert.NoError(t, err)
	assert.Nil(t, candidates)
	assert.Equal(t, "ingredient", result.Name)
}

func TestIngredientCreate_FindSimilarErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		Name: "findsimilarerror",
	}
	_, _, err := s.Create(ingredientDTO, uuid.New(), false)

	assert.EqualError(t, err, "internal server error")
}

// Neither pre-check caught it, but the DB's unique index did (a
// concurrent create) - the repository's "duplicate" sentinel passes
// through rather than being masked as "internal server error".
func TestIngredientCreate_RaceDuplicate(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		Name: "dupecreate",
	}
	_, _, err := s.Create(ingredientDTO, uuid.New(), true)

	assert.EqualError(t, err, "duplicate")
}

func TestIngredientCreate_RepoErr(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	ingredientDTO := models.IngredientDTO{
		Name: "error",
	}
	result, _, err := s.Create(ingredientDTO, uuid.New(), true)

	assert.Error(t, err)
	assert.IsType(t, models.IngredientDTO{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestIngredientSuggest_OK(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	result, err := s.Suggest("similar")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestIngredientSuggest_EmptyQuery(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	_, err := s.Suggest("")

	assert.EqualError(t, err, "query is required")
}

func TestIngredientSuggest_NoMatches(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	result, err := s.Suggest("nomatch")

	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestIngredientSuggest_Err(t *testing.T) {
	s := NewIngredientService(&IngredientRepositoryMock{})

	_, err := s.Suggest("findsimilarerror")

	assert.EqualError(t, err, "internal server error")
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
