package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
)

type RecipeRepository interface {
	FindAll(pagination m.PaginationRequest, requesterID uuid.UUID, scope m.RecipeListScope) ([]m.Recipe, int64, error)
	FindSingle(recipe m.Recipe) (m.Recipe, error)
	Create(recipe m.Recipe, outbox rmq.OutboxEvent) (m.Recipe, error)
	Update(recipe m.Recipe, outbox rmq.OutboxEvent) (m.Recipe, error)
	Delete(recipe m.Recipe, outbox rmq.OutboxEvent) error
}

type RecipeService struct {
	recipe RecipeRepository
	logger m.LoggerInterface
}

// NewRecipeService creates a new RecipeService instance. Publishing recipe
// events is no longer done here directly - Create/Update/Delete write an
// outbox row in the same DB transaction as the change (see
// RecipeRepository.Create/Update/Delete), and a shared/rabbitmq.Relay
// publishes from the outbox in the background (wired in config.go).
func NewRecipeService(recipeRepo RecipeRepository, logger m.LoggerInterface) *RecipeService {
	return &RecipeService{
		recipe: recipeRepo,
		logger: logger,
	}
}

// FindAll contains the business logic to get all recipes visible to
// requesterID, narrowed by scope (KIT-30: ?scope=mine|shared|all).
func (s RecipeService) FindAll(pagination m.PaginationRequest, requesterID uuid.UUID, scope m.RecipeListScope) (m.PaginatedResponse[m.RecipeDTO], error) {
	if scope == "" {
		scope = m.ScopeAll
	} else if !scope.Valid() {
		return m.PaginatedResponse[m.RecipeDTO]{}, errors.New("invalid scope")
	}

	recipes, total, err := s.recipe.FindAll(pagination, requesterID, scope)
	if err != nil {
		return m.PaginatedResponse[m.RecipeDTO]{}, errors.New("internal server error")
	}

	return m.PaginatedResponse[m.RecipeDTO]{
		Data:       m.Recipe{}.ConvertAllToDTO(recipes),
		Pagination: m.NewPaginationMetadata(pagination.Page, pagination.Limit, total),
	}, nil
}

// FindSingle contains the business logic to get a specific recipe.
// A private recipe not owned by requesterID is reported as not found,
// same as if it didn't exist, so its existence isn't leaked to other
// users. "shared" visibility is treated as owner-only until the
// recipe_shares model (KIT-32) exists to resolve explicit grants.
func (s RecipeService) FindSingle(recipeDTO m.RecipeDTO, requesterID uuid.UUID) (m.RecipeDTO, error) {
	recipe, err := s.recipe.FindSingle(recipeDTO.ConvertFromDTO())
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.RecipeDTO{}, err
		default:
			return m.RecipeDTO{}, errors.New("internal server error")
		}
	}

	if !canRead(recipe, requesterID) {
		return m.RecipeDTO{}, errors.New("not found")
	}

	return recipe.ConvertToDTO(), nil
}

// canRead reports whether requesterID may read recipe, given its owner and
// visibility. Shares (KIT-32) aren't modeled yet, so "shared" visibility
// grants nothing beyond the owner for now.
func canRead(recipe m.Recipe, requesterID uuid.UUID) bool {
	if recipe.OwnerID == requesterID {
		return true
	}
	return recipe.Visibility == m.VisibilityPublic
}

// Create handles the business logic for the creation of a recipe and passes the recipe object to the recipe repo for processing.
// ownerID comes from the authenticated Principal, never from the client-supplied DTO.
func (s RecipeService) Create(recipeDTO m.RecipeDTO, ownerID uuid.UUID) (m.RecipeDTO, error) {
	s.logger.Infof("%+v", recipeDTO)

	if recipeDTO.ID != uuid.Nil {
		return m.RecipeDTO{}, errors.New("existing id on new element is not allowed")
	}

	if recipeDTO.Name == "" {
		return m.RecipeDTO{}, errors.New("name is empty")
	}

	if recipeDTO.Description == "" {
		return m.RecipeDTO{}, errors.New("description is empty")
	}

	if recipeDTO.ServingCount == 0 {
		return m.RecipeDTO{}, errors.New("serving count 0 is not allowed")
	}

	if recipeDTO.Visibility == "" {
		recipeDTO.Visibility = m.VisibilityPrivate
	} else if !recipeDTO.Visibility.Valid() {
		return m.RecipeDTO{}, errors.New("invalid visibility")
	}

	recipe := recipeDTO.ConvertFromDTO()
	// Generated here (not left to the DB default) so the outbox event
	// below can carry the real recipe ID before the row is even inserted.
	recipe.ID = uuid.New()
	recipe.OwnerID = ownerID

	outbox, err := buildOutboxEntry(rmq.RecipeCreatedEvent, recipe, time.Now(), false)
	if err != nil {
		return m.RecipeDTO{}, errors.New("internal server error")
	}

	createdRecipe, err := s.recipe.Create(recipe, outbox)
	if err != nil {
		return m.RecipeDTO{}, err
	}

	return createdRecipe.ConvertToDTO(), nil
}

// Update applies changes to a recipe. Only the owner may update it -
// write access via a share (KIT-32) isn't modeled yet.
func (s RecipeService) Update(recipeDTO m.RecipeDTO, requesterID uuid.UUID) (m.RecipeDTO, error) {
	var updatedRecipe m.Recipe

	originalRecipe, err := s.recipe.FindSingle(recipeDTO.ConvertFromDTO())
	if err != nil {
		return m.RecipeDTO{}, err
	}

	if originalRecipe.OwnerID != requesterID {
		return m.RecipeDTO{}, errors.New("not found")
	}

	if recipeDTO.Name == "" {
		recipeDTO.Name = originalRecipe.Name
	}

	if recipeDTO.Description == "" {
		recipeDTO.Description = originalRecipe.Description
	}

	if recipeDTO.ServingCount == 0 {
		recipeDTO.ServingCount = originalRecipe.ServingCount
	}

	if recipeDTO.Visibility == "" {
		recipeDTO.Visibility = originalRecipe.Visibility
	} else if !recipeDTO.Visibility.Valid() {
		return m.RecipeDTO{}, errors.New("invalid visibility")
	}

	recipe := recipeDTO.ConvertFromDTO()
	recipe.OwnerID = originalRecipe.OwnerID

	outbox, err := buildOutboxEntry(rmq.RecipeUpdatedEvent, recipe, time.Now(), false)
	if err != nil {
		return m.RecipeDTO{}, errors.New("internal server error")
	}

	updatedRecipe, err = s.recipe.Update(recipe, outbox)
	if err != nil {
		return m.RecipeDTO{}, err
	}

	return updatedRecipe.ConvertToDTO(), nil
}

// Delete removes a recipe. Only the owner may delete it - write access via
// a share (KIT-32) isn't modeled yet.
func (s RecipeService) Delete(recipeDTO m.RecipeDTO, requesterID uuid.UUID) error {
	originalRecipe, err := s.recipe.FindSingle(recipeDTO.ConvertFromDTO())
	if err != nil {
		return err
	}

	if originalRecipe.OwnerID != requesterID {
		return errors.New("not found")
	}

	outbox, err := buildOutboxEntry(rmq.RecipeDeletedEvent, originalRecipe, time.Now(), true)
	if err != nil {
		return errors.New("internal server error")
	}

	if err := s.recipe.Delete(originalRecipe, outbox); err != nil {
		return err
	}

	return nil
}

// GetAccessState returns the authoritative RecipeAccessState for
// recipeID, straight from the Recipe row - no visibility filtering,
// because this is for the internal-only access endpoint (KIT-28) that
// other services' shared/authz instances call on their own
// projection-miss, not for a client request. Shares aren't populated yet
// (there's no recipe_shares table to read from until KIT-32).
func (s RecipeService) GetAccessState(recipeID uuid.UUID) (rmq.RecipeAccessState, error) {
	recipe, err := s.recipe.FindSingle(m.Recipe{ID: recipeID})
	if err != nil {
		return rmq.RecipeAccessState{}, err
	}

	return rmq.RecipeAccessState{
		RecipeID:   recipe.ID,
		OwnerID:    recipe.OwnerID.String(),
		Visibility: string(recipe.Visibility),
		Version:    recipe.UpdatedAt.UnixNano(),
		Deleted:    false,
	}, nil
}

// buildOutboxEntry fills eventTemplate's payload with recipe's full
// content and access state and returns the outbox row ready to insert
// alongside the recipe row in the same transaction.
//
// version is captured by the caller (time.Now(), not recipe.UpdatedAt):
// at this point the recipe row hasn't been written yet, so GORM's
// autoUpdateTime hook hasn't run and UpdatedAt is still zero. Capturing
// "now" here also guarantees deleted gets a version strictly after the
// last update's, which recipe.UpdatedAt could collide with for a soft
// delete (GORM's default soft-delete sets deleted_at only, not
// updated_at).
func buildOutboxEntry(eventTemplate rmq.Event, recipe m.Recipe, version time.Time, deleted bool) (rmq.OutboxEvent, error) {
	event := eventTemplate
	event.Payload = rmq.RecipePayload{
		ID:           recipe.ID,
		Name:         recipe.Name,
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
		OwnerID:      recipe.OwnerID.String(),
		Visibility:   string(recipe.Visibility),
		Version:      version.UnixNano(),
		Deleted:      deleted,
	}

	return rmq.NewOutboxEntry(event)
}
