package handlers

import (
	"net/http"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	hh "github.com/ihulsbus/cookbook/shared/http"
	"github.com/ihulsbus/cookbook/shared/models"
	"github.com/ihulsbus/cookbook/shared/oidc"
)

type IngredientService interface {
	FindAll(pagination models.PaginationRequest) (models.PaginatedResponse[models.IngredientDTO], error)
	FindSingle(ingredientDTO models.IngredientDTO) (models.IngredientDTO, error)
	Create(ingredientDTO models.IngredientDTO, createdBy uuid.UUID, force bool) (models.IngredientDTO, []models.IngredientDTO, error)
	Update(ingredientDTO models.IngredientDTO) (models.IngredientDTO, error)
	Delete(ingredientDTO models.IngredientDTO) error
	Suggest(q string) ([]models.IngredientDTO, error)
}

type IngredientHandlers struct {
	ingredientService IngredientService
	logger            m.LoggerInterface
}

func NewIngredientHandlers(ingredients IngredientService, logger m.LoggerInterface) *IngredientHandlers {
	return &IngredientHandlers{
		ingredientService: ingredients,
		logger:            logger,
	}
}

// principalID extracts the authenticated caller's user ID, needed for
// the ingredient's CreatedBy audit field. Every route reaching here has
// already run RequireAuth.
func (h IngredientHandlers) principalID(ctx *gin.Context) (uuid.UUID, bool) {
	p, ok := oidc.PrincipalFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated principal"})
		return uuid.Nil, false
	}

	id, err := uuid.Parse(p.UserID)
	if err != nil {
		h.logger.Errorf("principal sub %q is not a UUID: %s", p.UserID, err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid principal"})
		return uuid.Nil, false
	}

	return id, true
}

func (h IngredientHandlers) GetAll(ctx *gin.Context) {
	pagination, err := hh.ParsePagination(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingredientDTO, err := h.ingredientService.FindAll(pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ingredientDTO)
}

func (h IngredientHandlers) GetSingle(ctx *gin.Context) {
	var ingredientDTO models.IngredientDTO
	var err error

	ingredientDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ingredient ID"})
		return
	}

	ingredientDTO, err = h.ingredientService.FindSingle(ingredientDTO)
	if err != nil {
		switch err.Error() {
		case "not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no ingredient found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, ingredientDTO)
}

// Create adds an ingredient to the shared catalog. Any authenticated
// user may do this - duplicate prevention, not an admin gate, keeps the
// catalog clean. ?force=true skips the near-duplicate block; an exact
// normalized duplicate always blocks regardless.
func (h IngredientHandlers) Create(ctx *gin.Context) {
	createdBy, ok := h.principalID(ctx)
	if !ok {
		return
	}

	var ingredientDTO models.IngredientDTO
	if err := ctx.ShouldBindJSON(&ingredientDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unexpected JSON input"})
		return
	}

	force := ctx.Query("force") == "true"

	created, candidates, err := h.ingredientService.Create(ingredientDTO, createdBy, force)
	if err != nil {
		switch err.Error() {
		case "name is empty", "existing id on new element is not allowed":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "duplicate":
			ctx.JSON(http.StatusConflict, gin.H{"error": "ingredient already exists", "candidates": candidates})
		case "similar ingredients exist":
			ctx.JSON(http.StatusConflict, gin.H{"error": "similar ingredient exists - did you mean one of these? pass force=true to create anyway", "candidates": candidates})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

// Suggest returns up to 5 ingredients similar to q, for a type-ahead
// "pick an existing ingredient" UI.
func (h IngredientHandlers) Suggest(ctx *gin.Context) {
	candidates, err := h.ingredientService.Suggest(ctx.Query("q"))
	if err != nil {
		switch err.Error() {
		case "query is required":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, candidates)
}

func (h IngredientHandlers) Update(ctx *gin.Context) {
	var ingredientDTO models.IngredientDTO
	var err error

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ingredient ID"})
		return
	}

	if err = ctx.ShouldBindJSON(&ingredientDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Deliberately overwritten so the URL param ID wins over any id in the body.
	ingredientDTO.ID = id

	ingredientDTO, err = h.ingredientService.Update(ingredientDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ingredientDTO)
}

func (h IngredientHandlers) Delete(ctx *gin.Context) {
	var ingredientDTO models.IngredientDTO
	var err error

	ingredientDTO.ID, err = uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ingredient ID"})
		return
	}

	err = h.ingredientService.Delete(ingredientDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
