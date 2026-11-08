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

	// notFoundID is a recipe ID FindSingle treats as missing regardless of
	// Name - the share/visibility service methods (ListShares,
	// UpsertShare, DeleteShare, UpdateVisibility) only have an ID to look
	// up with (m.Recipe{ID: recipeID}, Name always ""), so they can't be
	// dispatched by Name the way the older tests are.
	notFoundID = uuid.New()

	// repoErrID resolves normally via FindSingle (owned by ownerID), but
	// its mutating repo call is wired to fail - see
	// RecipeRepositoryMock.UpdateVisibility.
	repoErrID = uuid.New()

	// shareOwnerID/shareTargetRecipe/shareTargetShares exercise canRead/
	// canWrite granted by an explicit share (KIT-32), as opposed to
	// ownership or public visibility. shareTargetShares is mutated per
	// test to control what FindShares(shareTargetRecipe.ID) returns.
	shareOwnerID      uuid.UUID = uuid.New()
	shareTargetRecipe m.Recipe = m.Recipe{
		ID:           uuid.New(),
		Name:         "shareaccess",
		Description:  "description",
		ServingCount: 1,
		OwnerID:      shareOwnerID,
		Visibility:   m.VisibilityPrivate,
	}
	shareTargetShares []m.RecipeShare
	shareFindErr      error
)

type RecipeRepositoryMock struct{}
type LoggerMock struct{}

func (RecipeRepositoryMock) FindAll(pagination m.PaginationRequest, requesterID uuid.UUID, requesterGroups []string, scope m.RecipeListScope) ([]m.Recipe, int64, error) {
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
	if recipeInput.ID == notFoundID {
		return m.Recipe{}, errors.New("not found")
	}
	if recipeInput.ID == repoErrID {
		r := recipe
		r.ID = repoErrID
		return r, nil
	}
	// The share/visibility service methods (ListShares, UpsertShare,
	// DeleteShare, UpdateVisibility) only have an ID to look up with -
	// m.Recipe{ID: recipeID}, Name always "" - so shareTargetRecipe must
	// be resolvable by ID alone, not just via the Name == "shareaccess"
	// dispatch FindSingle/Update's tests use.
	if recipeInput.ID == shareTargetRecipe.ID {
		return shareTargetRecipe, nil
	}

	switch recipeInput.Name {
	case "find":
		return recipe, nil
	case "shareaccess":
		return shareTargetRecipe, nil
	case "update", "updateerror", "delete", "":
		return recipe, nil
	case "notfound":
		return m.Recipe{}, errors.New("not found")
	default:
		return m.Recipe{}, errors.New("error")
	}
}

func (RecipeRepositoryMock) FindShares(recipeID uuid.UUID) ([]m.RecipeShare, error) {
	if shareFindErr != nil {
		return nil, shareFindErr
	}
	if recipeID == shareTargetRecipe.ID {
		return shareTargetShares, nil
	}
	return nil, nil
}

func (RecipeRepositoryMock) Create(recipeInput m.Recipe, outbox rmq.OutboxEvent) (m.Recipe, error) {
	switch recipeInput.Name {
	case "create":
		return recipe, nil
	default:
		return m.Recipe{}, errors.New("error")
	}
}

func (RecipeRepositoryMock) Update(recipeInput m.Recipe, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) (m.Recipe, error) {
	switch recipeInput.Name {
	case "update":
		return recipe, nil
	case "recipe":
		return recipe, nil
	case "shareaccess":
		return shareTargetRecipe, nil
	default:
		return m.Recipe{}, errors.New("error")
	}
}

// Delete is called with whatever FindSingle returned (the service layer
// re-deletes the record it just looked up, not the caller's raw input),
// so unlike the other mocks this dispatches on a distinct trigger value
// ("repoerror") rather than on the same Name the caller used to reach
// this point - see TestRecipeDelete_RepoErr.
func (RecipeRepositoryMock) Delete(recipeInput m.Recipe, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) error {
	switch recipeInput.Name {
	case "repoerror":
		return errors.New("repo error")
	default:
		return nil
	}
}

func (RecipeRepositoryMock) UpsertShare(share m.RecipeShare, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) error {
	switch share.TargetID {
	case "repoerror":
		return errors.New("repo error")
	default:
		return nil
	}
}

func (RecipeRepositoryMock) DeleteShare(recipeID uuid.UUID, targetType m.ShareTargetType, targetID string, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) error {
	switch targetID {
	case "notfound":
		return errors.New("not found")
	case "repoerror":
		return errors.New("repo error")
	default:
		return nil
	}
}

func (RecipeRepositoryMock) UpdateVisibility(recipeInput m.Recipe, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) (m.Recipe, error) {
	if recipeInput.ID == repoErrID {
		return m.Recipe{}, errors.New("repo error")
	}
	return recipeInput, nil
}

func TestRecipeFindAll_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	findAllRecipe.Name = "findall"

	result, err := s.FindAll(m.NormalizePagination(1, 25), ownerID, nil, m.ScopeAll)

	assert.NoError(t, err)
	assert.IsType(t, m.PaginatedResponse[m.RecipeDTO]{}, result)
	assert.Len(t, result.Data, 1)
}

func TestRecipeFindAll_Err(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	findAllRecipe.Name = "error"

	result, err := s.FindAll(m.NormalizePagination(1, 25), ownerID, nil, m.ScopeAll)

	assert.Error(t, err)
	assert.IsType(t, m.PaginatedResponse[m.RecipeDTO]{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestRecipeFindAll_NotFound(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	findAllRecipe.Name = "notfound"

	result, err := s.FindAll(m.NormalizePagination(1, 25), ownerID, nil, m.ScopeAll)

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

	result, err := s.FindAll(m.NormalizePagination(1, 25), ownerID, nil, "")

	assert.NoError(t, err)
	assert.IsType(t, m.PaginatedResponse[m.RecipeDTO]{}, result)
	assert.Len(t, result.Data, 1)
}

func TestRecipeFindAll_InvalidScope(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	result, err := s.FindAll(m.NormalizePagination(1, 25), ownerID, nil, "nonsense")

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
	result, err := s.FindSingle(recipeDTO, ownerID, nil)

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
	result, err := s.FindSingle(recipeDTO, ownerID, nil)

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
	result, err := s.FindSingle(recipeDTO, ownerID, nil)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "not found")
}

// TestRecipeFindSingle_PrivateNotOwner is the KIT-25 acceptance case: a
// private recipe not owned by the requester, and not shared with them
// (KIT-32), is reported as not found, not as a permission error, so its
// existence isn't leaked.
func TestRecipeFindSingle_PrivateNotOwner(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:   recipe.ID,
		Name: "find",
	}
	result, err := s.FindSingle(recipeDTO, uuid.New(), nil) // not the owner

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
	result, err := s.FindSingle(recipeDTO, uuid.New(), nil) // not the owner, but public

	assert.NoError(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.Equal(t, recipe.ID, result.ID)
}

// TestRecipeFindSingle_UserShare_Allowed: a private recipe not owned by
// the requester, but explicitly shared with them by user ID (KIT-32),
// grants read.
func TestRecipeFindSingle_UserShare_Allowed(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	requesterID := uuid.New()
	shareTargetShares = []m.RecipeShare{
		{RecipeID: shareTargetRecipe.ID, TargetType: m.ShareTargetUser, TargetID: requesterID.String(), Permission: m.SharePermissionRead},
	}
	defer func() { shareTargetShares = nil }()

	recipeDTO := m.RecipeDTO{ID: shareTargetRecipe.ID, Name: "shareaccess"}
	result, err := s.FindSingle(recipeDTO, requesterID, nil)

	assert.NoError(t, err)
	assert.Equal(t, shareTargetRecipe.ID, result.ID)
}

// TestRecipeFindSingle_GroupShare_Allowed: same as above, via a group the
// requester belongs to rather than their own user ID.
func TestRecipeFindSingle_GroupShare_Allowed(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareTargetShares = []m.RecipeShare{
		{RecipeID: shareTargetRecipe.ID, TargetType: m.ShareTargetGroup, TargetID: "cookbook-user", Permission: m.SharePermissionRead},
	}
	defer func() { shareTargetShares = nil }()

	recipeDTO := m.RecipeDTO{ID: shareTargetRecipe.ID, Name: "shareaccess"}
	result, err := s.FindSingle(recipeDTO, uuid.New(), []string{"cookbook-user"})

	assert.NoError(t, err)
	assert.Equal(t, shareTargetRecipe.ID, result.ID)
}

// TestRecipeFindSingle_UnrelatedShare_NotFound: shares exist on the
// recipe, but none target this requester or their groups - still not
// found, same as no shares at all.
func TestRecipeFindSingle_UnrelatedShare_NotFound(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareTargetShares = []m.RecipeShare{
		{RecipeID: shareTargetRecipe.ID, TargetType: m.ShareTargetUser, TargetID: uuid.New().String(), Permission: m.SharePermissionRead},
	}
	defer func() { shareTargetShares = nil }()

	recipeDTO := m.RecipeDTO{ID: shareTargetRecipe.ID, Name: "shareaccess"}
	result, err := s.FindSingle(recipeDTO, uuid.New(), nil)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestRecipeFindSingle_SharesLookupErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareFindErr = errors.New("db error")
	defer func() { shareFindErr = nil }()

	recipeDTO := m.RecipeDTO{ID: shareTargetRecipe.ID, Name: "shareaccess"}
	result, err := s.FindSingle(recipeDTO, uuid.New(), nil)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "internal server error")
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
	result, err := s.Update(recipeDTO, ownerID, nil)

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
	result, err := s.Update(recipeDTO, uuid.New(), nil) // not the owner, no share either

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "not found")
}

// TestRecipeUpdate_WriteShare_Allowed: a requester holding a write share
// (KIT-32) on a recipe they don't own may still update its content.
func TestRecipeUpdate_WriteShare_Allowed(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	requesterID := uuid.New()
	shareTargetShares = []m.RecipeShare{
		{RecipeID: shareTargetRecipe.ID, TargetType: m.ShareTargetUser, TargetID: requesterID.String(), Permission: m.SharePermissionWrite},
	}
	defer func() { shareTargetShares = nil }()

	recipeDTO := m.RecipeDTO{
		ID:           shareTargetRecipe.ID,
		Name:         "shareaccess",
		Description:  shareTargetRecipe.Description,
		ServingCount: shareTargetRecipe.ServingCount,
	}
	result, err := s.Update(recipeDTO, requesterID, nil)

	assert.NoError(t, err)
	assert.Equal(t, shareTargetRecipe.ID, result.ID)
}

// TestRecipeUpdate_ReadShareOnly_NotFound: a read-only share does not
// grant write, same masking as no share at all.
func TestRecipeUpdate_ReadShareOnly_NotFound(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	requesterID := uuid.New()
	shareTargetShares = []m.RecipeShare{
		{RecipeID: shareTargetRecipe.ID, TargetType: m.ShareTargetUser, TargetID: requesterID.String(), Permission: m.SharePermissionRead},
	}
	defer func() { shareTargetShares = nil }()

	recipeDTO := m.RecipeDTO{
		ID:           shareTargetRecipe.ID,
		Name:         "shareaccess",
		Description:  shareTargetRecipe.Description,
		ServingCount: shareTargetRecipe.ServingCount,
	}
	result, err := s.Update(recipeDTO, requesterID, nil)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "not found")
}

func TestRecipeUpdate_SharesLookupErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareFindErr = errors.New("db error")
	defer func() { shareFindErr = nil }()

	recipeDTO := m.RecipeDTO{
		ID:           shareTargetRecipe.ID,
		Name:         "shareaccess",
		Description:  shareTargetRecipe.Description,
		ServingCount: shareTargetRecipe.ServingCount,
	}
	result, err := s.Update(recipeDTO, uuid.New(), nil)

	assert.Error(t, err)
	assert.IsType(t, m.RecipeDTO{}, result)
	assert.EqualError(t, err, "internal server error")
}

func TestRecipeUpdate_NoNameErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	recipeDTO := m.RecipeDTO{
		ID:           recipe.ID,
		Name:         "",
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
	}
	result, err := s.Update(recipeDTO, ownerID, nil)

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
	result, err := s.Update(recipeDTO, ownerID, nil)

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
	result, err := s.Update(recipeDTO, ownerID, nil)

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
	result, err := s.Update(recipeDTO, ownerID, nil)

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
	result, err := s.Update(recipeDTO, ownerID, nil)

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
	result, err := s.Update(recipeDTO, ownerID, nil)

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

// TestRecipeDelete_WriteShare_NotFound locks in a deliberate design
// choice: unlike Update, a write share (KIT-32) does not grant Delete -
// only the owner can delete a recipe outright.
func TestRecipeDelete_WriteShare_NotFound(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	requesterID := uuid.New()
	shareTargetShares = []m.RecipeShare{
		{RecipeID: shareTargetRecipe.ID, TargetType: m.ShareTargetUser, TargetID: requesterID.String(), Permission: m.SharePermissionWrite},
	}
	defer func() { shareTargetShares = nil }()

	recipeDTO := m.RecipeDTO{ID: shareTargetRecipe.ID, Name: "shareaccess"}
	err := s.Delete(recipeDTO, requesterID)

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

func TestListShares_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareTargetShares = []m.RecipeShare{
		{RecipeID: shareTargetRecipe.ID, TargetType: m.ShareTargetUser, TargetID: uuid.New().String(), Permission: m.SharePermissionRead},
	}
	defer func() { shareTargetShares = nil }()

	result, err := s.ListShares(shareTargetRecipe.ID, shareOwnerID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestListShares_NotOwner(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	_, err := s.ListShares(shareTargetRecipe.ID, uuid.New())

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
}

func TestListShares_RecipeNotFound(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	_, err := s.ListShares(notFoundID, ownerID)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
}

func TestListShares_SharesErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareFindErr = errors.New("db error")
	defer func() { shareFindErr = nil }()

	_, err := s.ListShares(shareTargetRecipe.ID, shareOwnerID)

	assert.Error(t, err)
	assert.EqualError(t, err, "internal server error")
}

func TestUpsertShare_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareDTO := m.RecipeShareDTO{
		TargetType: m.ShareTargetUser,
		TargetID:   uuid.New().String(),
		Permission: m.SharePermissionRead,
	}
	result, err := s.UpsertShare(shareTargetRecipe.ID, shareOwnerID, shareDTO)

	assert.NoError(t, err)
	assert.Equal(t, shareDTO.TargetID, result.TargetID)
	assert.Equal(t, shareOwnerID, result.CreatedBy)
}

func TestUpsertShare_NotOwner(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareDTO := m.RecipeShareDTO{
		TargetType: m.ShareTargetUser,
		TargetID:   uuid.New().String(),
		Permission: m.SharePermissionRead,
	}
	_, err := s.UpsertShare(shareTargetRecipe.ID, uuid.New(), shareDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
}

func TestUpsertShare_InvalidTargetType(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareDTO := m.RecipeShareDTO{
		TargetType: "nonsense",
		TargetID:   uuid.New().String(),
		Permission: m.SharePermissionRead,
	}
	_, err := s.UpsertShare(shareTargetRecipe.ID, shareOwnerID, shareDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid target type")
}

func TestUpsertShare_MissingTargetID(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareDTO := m.RecipeShareDTO{
		TargetType: m.ShareTargetUser,
		TargetID:   "",
		Permission: m.SharePermissionRead,
	}
	_, err := s.UpsertShare(shareTargetRecipe.ID, shareOwnerID, shareDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "target id is required")
}

func TestUpsertShare_UserTargetNotUUID(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareDTO := m.RecipeShareDTO{
		TargetType: m.ShareTargetUser,
		TargetID:   "not-a-uuid",
		Permission: m.SharePermissionRead,
	}
	_, err := s.UpsertShare(shareTargetRecipe.ID, shareOwnerID, shareDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "target id must be a valid user id")
}

// TestUpsertShare_GroupTargetNotUUID_OK: a group target need not be a
// UUID (it's a group name, e.g. "cookbook-user") - only "user" targets
// are validated as UUIDs.
func TestUpsertShare_GroupTargetNotUUID_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareDTO := m.RecipeShareDTO{
		TargetType: m.ShareTargetGroup,
		TargetID:   "cookbook-user",
		Permission: m.SharePermissionRead,
	}
	result, err := s.UpsertShare(shareTargetRecipe.ID, shareOwnerID, shareDTO)

	assert.NoError(t, err)
	assert.Equal(t, "cookbook-user", result.TargetID)
}

func TestUpsertShare_InvalidPermission(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareDTO := m.RecipeShareDTO{
		TargetType: m.ShareTargetUser,
		TargetID:   uuid.New().String(),
		Permission: "nonsense",
	}
	_, err := s.UpsertShare(shareTargetRecipe.ID, shareOwnerID, shareDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid permission")
}

func TestUpsertShare_RepoErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	// "repoerror" is a group target (not UUID-validated, unlike "user")
	// so it reaches the repository mock, which is wired to fail on this
	// exact TargetID.
	shareDTO := m.RecipeShareDTO{
		TargetType: m.ShareTargetGroup,
		TargetID:   "repoerror",
		Permission: m.SharePermissionRead,
	}
	_, err := s.UpsertShare(shareTargetRecipe.ID, shareOwnerID, shareDTO)

	assert.Error(t, err)
	assert.EqualError(t, err, "internal server error")
}

func TestDeleteShare_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	err := s.DeleteShare(shareTargetRecipe.ID, shareOwnerID, m.ShareTargetUser, uuid.New().String())

	assert.NoError(t, err)
}

func TestDeleteShare_NotOwner(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	err := s.DeleteShare(shareTargetRecipe.ID, uuid.New(), m.ShareTargetUser, uuid.New().String())

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
}

func TestDeleteShare_RepoNotFound(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	err := s.DeleteShare(shareTargetRecipe.ID, shareOwnerID, m.ShareTargetUser, "notfound")

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
}

func TestDeleteShare_RepoErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	err := s.DeleteShare(shareTargetRecipe.ID, shareOwnerID, m.ShareTargetUser, "repoerror")

	assert.Error(t, err)
	assert.EqualError(t, err, "repo error")
}

func TestUpdateVisibility_OK(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	result, err := s.UpdateVisibility(recipe.ID, ownerID, m.VisibilityPublic)

	assert.NoError(t, err)
	assert.Equal(t, m.VisibilityPublic, result.Visibility)
}

func TestUpdateVisibility_NotOwner(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	_, err := s.UpdateVisibility(recipe.ID, uuid.New(), m.VisibilityPublic)

	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
}

func TestUpdateVisibility_InvalidVisibility(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	_, err := s.UpdateVisibility(recipe.ID, ownerID, "nonsense")

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid visibility")
}

func TestUpdateVisibility_RepoErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	_, err := s.UpdateVisibility(repoErrID, ownerID, m.VisibilityPublic)

	assert.Error(t, err)
	assert.EqualError(t, err, "internal server error")
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

// TestGetAccessState_IncludesShares: the access state's Shares (KIT-32)
// reflects the recipe's current recipe_shares rows, not an empty list.
func TestGetAccessState_IncludesShares(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	targetID := uuid.New().String()
	shareTargetShares = []m.RecipeShare{
		{RecipeID: shareTargetRecipe.ID, TargetType: m.ShareTargetUser, TargetID: targetID, Permission: m.SharePermissionWrite},
	}
	defer func() { shareTargetShares = nil }()

	state, err := s.GetAccessState(shareTargetRecipe.ID)

	assert.NoError(t, err)
	assert.Len(t, state.Shares, 1)
	assert.Equal(t, targetID, state.Shares[0].TargetID)
	assert.Equal(t, string(m.SharePermissionWrite), state.Shares[0].Permission)
}

func TestGetAccessState_SharesErr(t *testing.T) {
	s := NewRecipeService(&RecipeRepositoryMock{}, &LoggerMock{})

	shareFindErr = errors.New("db error")
	defer func() { shareFindErr = nil }()

	_, err := s.GetAccessState(shareTargetRecipe.ID)

	assert.Error(t, err)
	assert.EqualError(t, err, "internal server error")
}
