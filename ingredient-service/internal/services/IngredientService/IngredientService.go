package services

import (
	"errors"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
)

type IngredientRepository interface {
	FindAll(pagination m.PaginationRequest) ([]m.Ingredient, int64, error)
	FindSingle(ingredient m.Ingredient) (m.Ingredient, error)
	FindByName(name string) (m.Ingredient, error)
	// See IngredientRepository's doc comments for how "normalized" is computed.
	FindByNormalizedName(name string) (m.Ingredient, error)
	FindSimilar(name string, limit int) ([]m.Ingredient, error)
	Create(ingredient m.Ingredient) (m.Ingredient, error)
	Update(ingredient m.Ingredient) (m.Ingredient, error)
	Delete(ingredient m.Ingredient) error
}

// similarityCandidateLimit is how many near-duplicates to surface, both
// from the create-time "did you mean ...?" block and from Suggest.
const similarityCandidateLimit = 5

type IngredientService struct {
	repo IngredientRepository
}

func NewIngredientService(ingredientRepo IngredientRepository) *IngredientService {
	return &IngredientService{
		repo: ingredientRepo,
	}
}

func (s IngredientService) FindAll(pagination m.PaginationRequest) (m.PaginatedResponse[m.IngredientDTO], error) {
	ingredients, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return m.PaginatedResponse[m.IngredientDTO]{}, errors.New("internal server error")
	}

	return m.PaginatedResponse[m.IngredientDTO]{
		Data:       m.Ingredient{}.ConvertAllToDTO(ingredients),
		Pagination: m.NewPaginationMetadata(pagination.Page, pagination.Limit, total),
	}, nil
}

func (s IngredientService) FindSingle(ingredientDTO m.IngredientDTO) (m.IngredientDTO, error) {

	ingredient, err := s.repo.FindSingle(ingredientDTO.ConvertFromDTO())
	if err != nil {
		switch err.Error() {
		case "not found":
			return m.IngredientDTO{}, err
		default:
			return m.IngredientDTO{}, errors.New("internal server error")
		}
	}

	return ingredient.ConvertToDTO(), nil
}

// Create adds a new ingredient to the shared catalog - any authenticated
// user may do this, not just cookbook-admin. createdBy is recorded for audit.
//
// Two things can block a create, returned via the candidates slice so
// the handler can shape a 409 with "did you mean ...?" data:
//   - "duplicate": the normalized_name already exists (always blocks).
//   - "similar ingredients exist": a near-duplicate exists (pg_trgm
//     similarity > 0.6) but isn't exact; blocks unless force is true.
func (s IngredientService) Create(ingredientDTO m.IngredientDTO, createdBy uuid.UUID, force bool) (m.IngredientDTO, []m.IngredientDTO, error) {
	if ingredientDTO.ID != uuid.Nil {
		return m.IngredientDTO{}, nil, errors.New("existing id on new element is not allowed")
	}

	if ingredientDTO.Name == "" {
		return m.IngredientDTO{}, nil, errors.New("name is empty")
	}

	exact, err := s.repo.FindByNormalizedName(ingredientDTO.Name)
	if err != nil {
		return m.IngredientDTO{}, nil, errors.New("internal server error")
	}
	if exact.ID != uuid.Nil {
		return m.IngredientDTO{}, []m.IngredientDTO{exact.ConvertToDTO()}, errors.New("duplicate")
	}

	if !force {
		similar, err := s.repo.FindSimilar(ingredientDTO.Name, similarityCandidateLimit)
		if err != nil {
			return m.IngredientDTO{}, nil, errors.New("internal server error")
		}
		if len(similar) > 0 {
			return m.IngredientDTO{}, m.Ingredient{}.ConvertAllToDTO(similar), errors.New("similar ingredients exist")
		}
	}

	ingredient := ingredientDTO.ConvertFromDTO()
	ingredient.CreatedBy = createdBy

	created, err := s.repo.Create(ingredient)
	if err != nil {
		if err.Error() == "duplicate" {
			return m.IngredientDTO{}, nil, err
		}
		return m.IngredientDTO{}, nil, errors.New("internal server error")
	}

	return created.ConvertToDTO(), nil, nil
}

// Suggest returns up to 5 ingredients similar to q, for a type-ahead
// "pick an existing ingredient" UI. Reuses FindSimilar's similarity >
// 0.6 floor (the same one Create's near-duplicate check uses), so a
// weak match returns an empty list rather than forcing 5 results.
func (s IngredientService) Suggest(q string) ([]m.IngredientDTO, error) {
	if q == "" {
		return nil, errors.New("query is required")
	}

	similar, err := s.repo.FindSimilar(q, similarityCandidateLimit)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	return m.Ingredient{}.ConvertAllToDTO(similar), nil
}

func (s IngredientService) Update(ingredientDTO m.IngredientDTO) (m.IngredientDTO, error) {
	var ingredient m.Ingredient

	_, err := s.FindSingle(ingredientDTO)
	if err != nil {
		return m.IngredientDTO{}, errors.New("ingredient does not exist. nothing to update")
	}

	ingredient, err = s.repo.Update(ingredientDTO.ConvertFromDTO())
	if err != nil {
		return m.IngredientDTO{}, err
	}

	return ingredient.ConvertToDTO(), nil
}

func (s IngredientService) Delete(ingredientDTO m.IngredientDTO) error {

	_, err := s.FindSingle(ingredientDTO)
	if err != nil {
		return errors.New("ingredient does not exist. nothing to delete")
	}

	// TODO: check if there are recipes using the ingredient. If so, an error should be returned and the ingredient should not be deleted.
	err = s.repo.Delete(ingredientDTO.ConvertFromDTO())
	if err != nil {
		return err
	}

	return nil
}
