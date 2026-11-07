package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"github.com/stretchr/testify/assert"
)

var (
	ownerID uuid.UUID = uuid.New()

	findAllRecipe m.Recipe = m.Recipe{
		ID:           uuid.New(),
		Name:         "recipe",
		Description:  "description",
		ServingCount: 1,
		OwnerID:      ownerID,
		Visibility:   m.VisibilityPrivate,
	}
	recipe m.Recipe = m.Recipe{
		ID:           uuid.New(),
		Name:         "recipe",
		Description:  "description",
		ServingCount: 1,
		OwnerID:      ownerID,
		Visibility:   m.VisibilityPrivate,
	}
)

type RecipeRepositoryMock struct{}
type LoggerMock struct{}

func (RecipeRepositoryMock) FindAll(pagination m.PaginationRequest, requesterID uuid.UUID, scope m.RecipeListScope) ([]m.Recipe, int64, error) {
	switch findAllRecipe.Name {
	case "findall":
		return []m.Recipe{recipe}, 1, nil
	case "notfound":
		return nil, 0, errors.New("not found")
	default:
		return nil, 0, errors.New("error")
	}
}

func (LoggerMock) Debugf(format string, args ...interface{}) {}
func (LoggerMock) Infof(format string, args ...interface{})  {}
func (LoggerMock) Warnf(format string, args ...interface{})  {}
func (LoggerMock) Errorf(format string, args ...interface{}) {}

func (RecipeRepositoryMock) FindSingle(recipeInput m.Recipe) (m.Recipe, error) {
	switch recipeInput.Name {
	case "find":
		return recipe, nil
	case "update", "updateerror", "delete", "":
		return recipe, nil
	case "notfound":
		return m.Recipe{}, errors.New("not found")
	default:
		return m.Recipe{}, errors.New("error")
	}
}

func (RecipeRepositoryMock) Create(recipeInput m.Recipe, outbox rmq.OutboxEvent) (m.Recipe, error) {
	switch recipeInput.Name {
	case "create":
		return recipe, nil
	default:
		return m.Recipe{}, errors.New("error")
	}
}

func (RecipeRepositoryMock) Update(recipeInput m.Recipe, outbox rmq.OutboxEvent) (m.Recipe, error) {
	switch recipeInput.Name {
	case "update":
		return recipe, nil
	case "recipe":
		return recipe, nil
	default:
		return m.Recipe{}, errors.New("error")
	}
}

// Delete is called with whatever FindSingle returned (the service layer
// re-deletes the record it just looked up, not the caller's raw input),
// so unlike the other mocks this dispatches on a distinct trigger value
// ("repoerror") rather than on the same Name the caller used to reach
// this point - see TestRecipeDelete_RepoErr.
func (RecipeRepositoryMock) Delete(recipeInput m.Recipe, outbox rmq.OutboxEvent) error {
	switch recipeInput.Name {
	case "repoerror":
		return errors.New("repo error")
	default:
		return nil
	}
}

func TestRecipeFindAll_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	findAllRecipe.Name = "findall"

	result, err := s.FindAll(m.NormalizePagination(1, 25), ownerID, m.ScopeAll)

	assert.NoError(t, err)
	assert.IsType(t, m.PaginatedResponse[m.RecipeDTO]{}, result)
	assert.Len(t, result.Data, 1)
}

func TestRecipeFindAll_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	findAllRecipe.Name = "error"

	result, err := s.FindAll(m.NormalizePagination(1, 25), ownerID, m.ScopeAll)

	assert.Error(t, err)
	assert.IsType(t, m.PaginatedResponse[m.RecipeDTO]{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestRecipeFindAll_NotFound(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	findAllRecipe.Name = "notfound"

	result, err := s.FindAll(m.NormalizePagination(1, 25), ownerID, m.ScopeAll)

	assert.Error(t, err)
	assert.IsType(t, m.PaginatedResponse[m.RecipeDTO]{}, result)
	assert.EqualError(t, err, "internal server error")
}

// TestRecipeFindAll_EmptyScopeDefaultsToAll: an empty ?scope (the param
// omitted) must behave exactly like ScopeAll - GET /recipe's behavior
// before ?scope existed (KIT-30).
func TestRecipeFindAll_EmptyScopeDefaultsToAll(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	findAllRecipe.Name = "findall"

	result, err := s.FindAll(m.NormalizePagination(1, 25), ownerID, "")

	assert.NoError(t, err)
	assert.IsType(t, m.PaginatedResponse[m.RecipeDTO]{}, result)
	assert.Len(t, result.Data, 1)
}

func TestRecipeFindAll_InvalidScope(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	result, err := s.FindAll(m.NormalizePagination(1, 25), ownerID, "nonsense")

	assert.Error(t, err)
	assert.IsType(t, m.PaginatedResponse[m.RecipeDTO]{}, result)
	assert.EqualError(t, err, "invalid scope")
}

func TestRecipeFindSingle_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "find",
	}
	result, err := s.FindSingle(recipeDTO, ownerID)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.Equal(t, "recipe", result.Name)
	assert.Equal(t, recipe.ID, result.ID)
}

func TestRecipeFindSingle_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "error",
	}
	result, err := s.FindSingle(recipeDTO, ownerID)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestRecipeFindSingle_NotFound(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "notfound",
	}
	result, err := s.FindSingle(recipeDTO, ownerID)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "not found")
}

// TestRecipeFindSingle_PrivateNotOwner is the KIT-25 acceptance case: a
// private recipe not owned by the requester is reported as not found,
// not as a permission error, so its existence isn't leaked.
func TestRecipeFindSingle_PrivateNotOwner(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "find",
	}
	result, err := s.FindSingle(recipeDTO, uuid.New()) // not the owner

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestRecipeFindSingle_PublicNotOwner_Allowed(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	originalVisibility := recipe.Visibility
	recipe.Visibility = m.VisibilityPublic
	defer func() { recipe.Visibility = originalVisibility }()

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "find",
	}
	result, err := s.FindSingle(recipeDTO, uuid.New()) // not the owner, but public

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.Equal(t, recipe.ID, result.ID)
}

func TestRecipeCreate_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		Name:         "create",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Create(recipeDTO, ownerID)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.Equal(t, "recipe", result.Name)
	assert.Equal(t, recipe.ID, result.ID)
}

func TestRecipeCreate_IDErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "create",
	}
	result, err := s.Create(recipeDTO, ownerID)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "existing id on new element is not allowed")
}

func TestRecipeCreate_NoName(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		Name: "",
	}
	result, err := s.Create(recipeDTO, ownerID)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "name is empty")
}

func TestRecipeCreate_NoDescription(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		Name:        recipe.Name,
		Description: "",
	}
	result, err := s.Create(recipeDTO, ownerID)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "description is empty")
}

func TestRecipeCreate_NoServingCount(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		Name:         recipe.Name,
		Description:  recipe.Description,
		ServingCount: 0,
	}
	result, err := s.Create(recipeDTO, ownerID)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "serving count 0 is not allowed")
}

func TestRecipeCreate_InvalidVisibility(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		Name:         "create",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
		Visibility:   "nonsense",
	}
	result, err := s.Create(recipeDTO, ownerID)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "invalid visibility")
}

func TestRecipeCreate_CreateErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		Name:         "error",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Create(recipeDTO, ownerID)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestRecipeUpdate_Ok(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "update",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Update(recipeDTO, ownerID)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.Equal(t, recipe.ID, result.ID)
	assert.Equal(t, recipe.Name, result.Name)
	assert.Equal(t, recipe.Description, result.Description)
	assert.Equal(t, recipe.ServingCount, result.ServingCount)
}

func TestRecipeUpdate_NotOwner(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "update",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Update(recipeDTO, uuid.New()) // not the owner

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestRecipeUpdate_NoNameErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Update(recipeDTO, ownerID)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
}

func TestRecipeUpdate_NoDescription(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "update",
		Description:  "",
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Update(recipeDTO, ownerID)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
}

func TestRecipeUpdate_NoServingCount(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "update",
		Description:  recipe.Description,
		ServingCount: 0,
	}
	result, err := s.Update(recipeDTO, ownerID)

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
}

func TestRecipeUpdate_InvalidVisibility(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "update",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
		Visibility:   "nonsense",
	}
	result, err := s.Update(recipeDTO, ownerID)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "invalid visibility")
}

func TestRecipeUpdate_FindErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "notfound",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Update(recipeDTO, ownerID)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestRecipeUpdate_UpdateErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "updateerror",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Update(recipeDTO, ownerID)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "error")
}

func TestRecipeDelete_Ok(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "delete",
	}
	err := s.Delete(recipeDTO, ownerID)

	assert.NoError(t, err)
}

func TestRecipeDelete_NotOwner(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "delete",
	}
	err := s.Delete(recipeDTO, uuid.New()) // not the owner

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
}

func TestRecipeDelete_DeleteErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "error",
	}
	err := s.Delete(recipeDTO, ownerID)

	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}

// TestRecipeDelete_RepoErr exercises the repository-level failure path
// specifically (as opposed to TestRecipeDelete_DeleteErr, which fails at
// the FindSingle lookup before ever reaching Delete).
func TestRecipeDelete_RepoErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	original := recipe.Name
	recipe.Name = "repoerror"
	defer func() { recipe.Name = original }()

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "delete",
	}
	err := s.Delete(recipeDTO, ownerID)

	assert.Error(t, err)
	assert.EqualError(t, err, "repo error")
}

// TestGetAccessState_OK: GetAccessState builds its lookup with an empty
// Name (it only has an ID to go on), which lands on the mock's "" case -
// the same one FindSingle's other zero-Name callers hit - and returns
// the recipe fixture.
func TestGetAccessState_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	state, err := s.GetAccessState(recipe.ID)

	assert.NoError(t, err)
	assert.Equal(t, recipe.ID, state.RecipeID)
	assert.Equal(t, ownerID.String(), state.OwnerID)
	assert.Equal(t, string(m.VisibilityPrivate), state.Visibility)
	assert.False(t, state.Deleted)
}
