package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
)

type RecipeRepository interface {
	FindAll(pagination m.PaginationRequest, requesterID uuid.UUID, requesterGroups []string, scope m.RecipeListScope) ([]m.Recipe, int64, error)
	FindSingle(recipe m.Recipe) (m.Recipe, error)
	FindShares(recipeID uuid.UUID) ([]m.RecipeShare, error)
	Create(recipe m.Recipe, outbox rmq.OutboxEvent) (m.Recipe, error)
	Update(recipe m.Recipe, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) (m.Recipe, error)
	Delete(recipe m.Recipe, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) error
	UpsertShare(share m.RecipeShare, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) error
	DeleteShare(recipeID uuid.UUID, targetType m.ShareTargetType, targetID string, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) error
	UpdateVisibility(recipe m.Recipe, buildOutbox func(shares []m.RecipeShare) (rmq.OutboxEvent, error)) (m.Recipe, error)
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
// requesterGroups is only consulted for ScopeShared (KIT-35), to also
// surface recipes shared with a group requesterID belongs to.
func (s RecipeService) FindAll(pagination m.PaginationRequest, requesterID uuid.UUID, requesterGroups []string, scope m.RecipeListScope) (m.PaginatedResponse[m.RecipeDTO], error) {
	if scope == "" {
		scope = m.ScopeAll
	} else if !scope.Valid() {
		return m.PaginatedResponse[m.RecipeDTO]{}, errors.New("invalid scope")
	}

	recipes, total, err := s.recipe.FindAll(pagination, requesterID, requesterGroups, scope)
	if err != nil {
		return m.PaginatedResponse[m.RecipeDTO]{}, errors.New("internal server error")
	}

	return m.PaginatedResponse[m.RecipeDTO]{
		Data:       m.Recipe{}.ConvertAllToDTO(recipes),
		Pagination: m.NewPaginationMetadata(pagination.Page, pagination.Limit, total),
	}, nil
}

// FindSingle contains the business logic to get a specific recipe.
// A private recipe not owned by requesterID, and not shared with them or
// one of requesterGroups (KIT-32), is reported as not found, same as if
// it didn't exist, so its existence isn't leaked to other users.
func (s RecipeService) FindSingle(recipeDTO m.RecipeDTO, requesterID uuid.UUID, requesterGroups []string) (m.RecipeDTO, error) {
	recipe, err := s.recipe.FindSingle(recipeDTO.ConvertFromDTO())
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.RecipeDTO{}, err
		default:
			return m.RecipeDTO{}, errors.New("internal server error")
		}
	}

	if recipe.OwnerID != requesterID && recipe.Visibility != m.VisibilityPublic {
		shares, err := s.recipe.FindShares(recipe.ID)
		if err != nil {
			return m.RecipeDTO{}, errors.New("internal server error")
		}
		if !canRead(recipe, shares, requesterID, requesterGroups) {
			return m.RecipeDTO{}, errors.New("not found")
		}
	}

	return recipe.ConvertToDTO(), nil
}

// canRead reports whether requesterID (member of requesterGroups) may
// read recipe, given its owner, visibility, and explicit shares
// (KIT-32). Mirrors shared/authz's evaluate() so the two stay consistent:
// recipe-service enforces this on its own routes, shared/authz on every
// other service's.
func canRead(recipe m.Recipe, shares []m.RecipeShare, requesterID uuid.UUID, requesterGroups []string) bool {
	if recipe.OwnerID == requesterID {
		return true
	}
	if recipe.Visibility == m.VisibilityPublic {
		return true
	}
	return hasShareAccess(shares, requesterID, requesterGroups, m.SharePermissionRead)
}

// canWrite reports whether requesterID (member of requesterGroups) may
// write recipe's content, given its owner and explicit shares (KIT-32).
// Unlike canRead, visibility never grants write - only ownership or a
// write share does.
func canWrite(recipe m.Recipe, shares []m.RecipeShare, requesterID uuid.UUID, requesterGroups []string) bool {
	if recipe.OwnerID == requesterID {
		return true
	}
	return hasShareAccess(shares, requesterID, requesterGroups, m.SharePermissionWrite)
}

// hasShareAccess reports whether any share in shares targets requesterID
// directly or one of requesterGroups, with at least perm. A write share
// also satisfies a read check, same as shared/authz's evaluate().
func hasShareAccess(shares []m.RecipeShare, requesterID uuid.UUID, requesterGroups []string, perm m.SharePermission) bool {
	for _, share := range shares {
		var targetsRequester bool
		switch share.TargetType {
		case m.ShareTargetUser:
			targetsRequester = share.TargetID == requesterID.String()
		case m.ShareTargetGroup:
			targetsRequester = containsString(requesterGroups, share.TargetID)
		}

		if !targetsRequester {
			continue
		}
		if perm == m.SharePermissionRead || share.Permission == m.SharePermissionWrite {
			return true
		}
	}
	return false
}

func containsString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
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

	// A brand new recipe can't have any shares yet - nothing to look up.
	outbox, err := buildOutboxEntry(rmq.RecipeCreatedEvent, recipe, time.Now(), false, nil)
	if err != nil {
		return m.RecipeDTO{}, errors.New("internal server error")
	}

	createdRecipe, err := s.recipe.Create(recipe, outbox)
	if err != nil {
		return m.RecipeDTO{}, err
	}

	return createdRecipe.ConvertToDTO(), nil
}

// Update applies changes to a recipe. The owner, or a requester holding a
// write share on it (KIT-32), may update its content.
func (s RecipeService) Update(recipeDTO m.RecipeDTO, requesterID uuid.UUID, requesterGroups []string) (m.RecipeDTO, error) {
	var updatedRecipe m.Recipe

	originalRecipe, err := s.recipe.FindSingle(recipeDTO.ConvertFromDTO())
	if err != nil {
		return m.RecipeDTO{}, err
	}

	if originalRecipe.OwnerID != requesterID {
		shares, err := s.recipe.FindShares(originalRecipe.ID)
		if err != nil {
			return m.RecipeDTO{}, errors.New("internal server error")
		}
		if !canWrite(originalRecipe, shares, requesterID, requesterGroups) {
			return m.RecipeDTO{}, errors.New("not found")
		}
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

	updatedRecipe, err = s.recipe.Update(recipe, func(shares []m.RecipeShare) (rmq.OutboxEvent, error) {
		return buildOutboxEntry(rmq.RecipeUpdatedEvent, recipe, time.Now(), false, shares)
	})
	if err != nil {
		return m.RecipeDTO{}, err
	}

	return updatedRecipe.ConvertToDTO(), nil
}

// Delete removes a recipe. Only the owner may delete it - a write share
// (KIT-32) grants editing the recipe's content, not deleting it outright.
func (s RecipeService) Delete(recipeDTO m.RecipeDTO, requesterID uuid.UUID) error {
	originalRecipe, err := s.recipe.FindSingle(recipeDTO.ConvertFromDTO())
	if err != nil {
		return err
	}

	if originalRecipe.OwnerID != requesterID {
		return errors.New("not found")
	}

	if err := s.recipe.Delete(originalRecipe, func(shares []m.RecipeShare) (rmq.OutboxEvent, error) {
		return buildOutboxEntry(rmq.RecipeDeletedEvent, originalRecipe, time.Now(), true, shares)
	}); err != nil {
		return err
	}

	return nil
}

// ListShares returns every explicit share on a recipe (KIT-32). Owner
// only - reported as "not found" for anyone else, same masking FindSingle
// uses, rather than a distinct permission error that would confirm the
// recipe's existence to a non-owner.
func (s RecipeService) ListShares(recipeID uuid.UUID, requesterID uuid.UUID) ([]m.RecipeShareDTO, error) {
	recipe, err := s.recipe.FindSingle(m.Recipe{ID: recipeID})
	if err != nil {
		return nil, err
	}

	if recipe.OwnerID != requesterID {
		return nil, errors.New("not found")
	}

	shares, err := s.recipe.FindShares(recipeID)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	return m.RecipeShare{}.ConvertAllToDTO(shares), nil
}

// UpsertShare creates or updates one explicit share on a recipe (KIT-32).
// Owner only; a write share does not grant re-sharing in v1 (ticket
// text), so a write-share holder gets the same "not found" as anyone
// else here, not a permission error.
func (s RecipeService) UpsertShare(recipeID uuid.UUID, requesterID uuid.UUID, shareDTO m.RecipeShareDTO) (m.RecipeShareDTO, error) {
	recipe, err := s.recipe.FindSingle(m.Recipe{ID: recipeID})
	if err != nil {
		return m.RecipeShareDTO{}, err
	}

	if recipe.OwnerID != requesterID {
		return m.RecipeShareDTO{}, errors.New("not found")
	}

	if !shareDTO.TargetType.Valid() {
		return m.RecipeShareDTO{}, errors.New("invalid target type")
	}

	if shareDTO.TargetID == "" {
		return m.RecipeShareDTO{}, errors.New("target id is required")
	}

	if shareDTO.TargetType == m.ShareTargetUser {
		if _, err := uuid.Parse(shareDTO.TargetID); err != nil {
			return m.RecipeShareDTO{}, errors.New("target id must be a valid user id")
		}
	}

	if !shareDTO.Permission.Valid() {
		return m.RecipeShareDTO{}, errors.New("invalid permission")
	}

	share := shareDTO.ConvertFromDTO(recipeID, requesterID)

	if err := s.recipe.UpsertShare(share, func(shares []m.RecipeShare) (rmq.OutboxEvent, error) {
		return buildAccessOutboxEntry(recipe, time.Now(), false, shares)
	}); err != nil {
		return m.RecipeShareDTO{}, errors.New("internal server error")
	}

	return share.ConvertToDTO(), nil
}

// DeleteShare removes one explicit share from a recipe (KIT-32). Owner
// only, same masking as UpsertShare/ListShares.
func (s RecipeService) DeleteShare(recipeID uuid.UUID, requesterID uuid.UUID, targetType m.ShareTargetType, targetID string) error {
	recipe, err := s.recipe.FindSingle(m.Recipe{ID: recipeID})
	if err != nil {
		return err
	}

	if recipe.OwnerID != requesterID {
		return errors.New("not found")
	}

	if err := s.recipe.DeleteShare(recipeID, targetType, targetID, func(shares []m.RecipeShare) (rmq.OutboxEvent, error) {
		return buildAccessOutboxEntry(recipe, time.Now(), false, shares)
	}); err != nil {
		return err
	}

	return nil
}

// UpdateVisibility changes a recipe's Visibility (PUT /recipe/:id/visibility,
// KIT-32). Owner only.
func (s RecipeService) UpdateVisibility(recipeID uuid.UUID, requesterID uuid.UUID, visibility m.Visibility) (m.RecipeDTO, error) {
	recipe, err := s.recipe.FindSingle(m.Recipe{ID: recipeID})
	if err != nil {
		return m.RecipeDTO{}, err
	}

	if recipe.OwnerID != requesterID {
		return m.RecipeDTO{}, errors.New("not found")
	}

	if !visibility.Valid() {
		return m.RecipeDTO{}, errors.New("invalid visibility")
	}

	recipe.Visibility = visibility

	updated, err := s.recipe.UpdateVisibility(recipe, func(shares []m.RecipeShare) (rmq.OutboxEvent, error) {
		return buildAccessOutboxEntry(recipe, time.Now(), false, shares)
	})
	if err != nil {
		return m.RecipeDTO{}, errors.New("internal server error")
	}

	return updated.ConvertToDTO(), nil
}

// GetAccessState returns the authoritative RecipeAccessState for
// recipeID, straight from the Recipe row and its current shares - no
// visibility filtering, because this is for the internal-only access
// endpoint (KIT-28) that other services' shared/authz instances call on
// their own projection-miss, not for a client request.
func (s RecipeService) GetAccessState(recipeID uuid.UUID) (rmq.RecipeAccessState, error) {
	recipe, err := s.recipe.FindSingle(m.Recipe{ID: recipeID})
	if err != nil {
		return rmq.RecipeAccessState{}, err
	}

	shares, err := s.recipe.FindShares(recipeID)
	if err != nil {
		return rmq.RecipeAccessState{}, errors.New("internal server error")
	}

	return rmq.RecipeAccessState{
		RecipeID:   recipe.ID,
		OwnerID:    recipe.OwnerID.String(),
		Visibility: string(recipe.Visibility),
		Shares:     toShareEntries(shares),
		Version:    recipe.UpdatedAt.UnixNano(),
		Deleted:    false,
	}, nil
}

// toShareEntries converts this service's RecipeShare rows to the
// rmq.ShareEntry shape ECST events carry.
func toShareEntries(shares []m.RecipeShare) []rmq.ShareEntry {
	entries := make([]rmq.ShareEntry, 0, len(shares))
	for _, share := range shares {
		entries = append(entries, rmq.ShareEntry{
			TargetType: string(share.TargetType),
			TargetID:   share.TargetID,
			Permission: string(share.Permission),
		})
	}
	return entries
}

// buildOutboxEntry fills eventTemplate's payload with recipe's full
// content and access state - including shares (KIT-32), read by the
// caller inside the same transaction the recipe row itself is written in
// (see DatabaseRepository.Update/Delete) - and returns the outbox row
// ready to insert alongside it.
//
// version is captured by the caller (time.Now(), not recipe.UpdatedAt):
// at this point the recipe row hasn't been written yet, so GORM's
// autoUpdateTime hook hasn't run and UpdatedAt is still zero. Capturing
// "now" here also guarantees deleted gets a version strictly after the
// last update's, which recipe.UpdatedAt could collide with for a soft
// delete (GORM's default soft-delete sets deleted_at only, not
// updated_at).
func buildOutboxEntry(eventTemplate rmq.Event, recipe m.Recipe, version time.Time, deleted bool, shares []m.RecipeShare) (rmq.OutboxEvent, error) {
	event := eventTemplate
	event.Payload = rmq.RecipePayload{
		ID:           recipe.ID,
		Name:         recipe.Name,
		Description:  recipe.Description,
		ServingCount: recipe.ServingCount,
		OwnerID:      recipe.OwnerID.String(),
		Visibility:   string(recipe.Visibility),
		Shares:       toShareEntries(shares),
		Version:      version.UnixNano(),
		Deleted:      deleted,
	}

	return rmq.NewOutboxEntry(event)
}

// buildAccessOutboxEntry is buildOutboxEntry for recipe.access.changed
// (KIT-32): used by the share/visibility endpoints, which change only a
// recipe's access state, never its content, so there's no reason to
// publish a full RecipePayload for them.
func buildAccessOutboxEntry(recipe m.Recipe, version time.Time, deleted bool, shares []m.RecipeShare) (rmq.OutboxEvent, error) {
	event := rmq.RecipeAccessChangedEvent
	event.Payload = rmq.RecipeAccessState{
		RecipeID:   recipe.ID,
		OwnerID:    recipe.OwnerID.String(),
		Visibility: string(recipe.Visibility),
		Shares:     toShareEntries(shares),
		Version:    version.UnixNano(),
		Deleted:    deleted,
	}

	return rmq.NewOutboxEntry(event)
}
