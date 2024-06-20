package services

import (
	"errors"

	m "ingredient-service/internal/models"

	"github.com/google/uuid"
)

type IngredientRepository interface {
	FindAll() ([]m.Ingredient, error)
	FindUnits() ([]m.Unit, error)
	FindSingle(ingredient m.Ingredient) (m.Ingredient, error)
	Create(ingredient m.Ingredient) (m.Ingredient, error)
	Update(ingredient m.Ingredient) (m.Ingredient, error)
	Delete(ingredient m.Ingredient) error
}
type IngredientService struct {
	repo IngredientRepository
}

// NewIngredientService creates a new IngredientService instance
func NewIngredientService(ingredientRepo IngredientRepository) *IngredientService {
	return &IngredientService{
		repo: ingredientRepo,
	}
}

func (s IngredientService) FindAll() ([]m.IngredientDTO, error) {
	var ingredients []m.Ingredient

	ingredients, err := s.repo.FindAll()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	return m.Ingredient{}.ConvertAllToDTO(ingredients), nil
}

func (s IngredientService) FindUnits() ([]m.UnitDTO, error) {
	var units []m.Unit

	units, err := s.repo.FindUnits()
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, errors.New("internal server error")
		}
	}

	return m.Unit{}.ConvertAllToDTO(units), nil
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

func (s IngredientService) Create(ingredientDTO m.IngredientDTO) (m.IngredientDTO, error) {
	var ingredient m.Ingredient

	if ingredientDTO.ID != uuid.Nil {
		return m.IngredientDTO{}, errors.New("existing id on new element is not allowed")
	}

	if ingredientDTO.Name == "" {
		return m.IngredientDTO{}, errors.New("name is empty")
	}

	found, err := s.FindSingle(ingredientDTO)
	if err == nil || found.ID != uuid.Nil {
		return m.IngredientDTO{}, errors.New("ingredient already exists")
	}

	ingredient, err = s.repo.Create(ingredientDTO.ConvertFromDTO())
	if err != nil {
		return m.IngredientDTO{}, err
	}

	return ingredient.ConvertToDTO(), nil
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

	// TODO: check if there are recipies using the ingredient. If so, an error should be returned and the ingredient should not be deleted.
	err = s.repo.Delete(ingredientDTO.ConvertFromDTO())
	if err != nil {
		return err
	}

	return nil
}
