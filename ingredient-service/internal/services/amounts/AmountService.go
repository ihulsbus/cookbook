package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	m "ingredient-service/internal/models"
)

type AmountRepository interface {
	Find(recipeID uuid.UUID) (*[]m.Amount, error)
	Create(amounts *[]m.Amount) (*[]m.Amount, error)
	Delete(amounts *[]m.Amount) error
}
type AmountService struct {
	repo AmountRepository
}

// NewAmountService creates a new AmountService instance
func NewAmountService(amountRepo AmountRepository) *AmountService {
	return &AmountService{
		repo: amountRepo,
	}
}

func (s AmountService) Find(recipeID uuid.UUID) (*[]m.AmountDTO, error) {

	if recipeID == uuid.Nil {
		return nil, errors.New("RecipeID is empty or invalid")
	}

	amounts, err := s.repo.Find(recipeID)
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, err
		default:
			return nil, fmt.Errorf("internal server error: %s", err.Error())
		}
	}

	amountDTO := m.Amount{}.ConvertAllToDTO(*amounts)
	return &amountDTO, nil
}

func (s AmountService) Create(recipeID uuid.UUID, amountsDTO *[]m.AmountDTO) (*[]m.AmountDTO, error) {
	_, err := s.Find(recipeID)

	if err == nil {
		return nil, errors.New("recipe has existing ingredients. Use update or delete first")
	}

	switch err.Error() {
	case "not found":
		amounts := m.AmountDTO{}.ConvertAllFromDTO(*amountsDTO)

		// force set all ingredient recipeID fields to the provided recipeID
		for i := range amounts {
			amounts[i].RecipeID = recipeID
		}

		amountsResponse, err := s.repo.Create(&amounts)
		if err != nil {
			return nil, err
		}

		amountDTO := m.Amount{}.ConvertAllToDTO(*amountsResponse)
		return &amountDTO, nil
	default:
		return nil, err
	}
}

func (s AmountService) Update(recipeID uuid.UUID, amountsDTO *[]m.AmountDTO) (*[]m.AmountDTO, error) {
	var err error

	existingIngredients, err := s.repo.Find(recipeID)
	if err != nil {
		return nil, err
	}

	err = s.Delete(recipeID)
	if err != nil {
		return nil, fmt.Errorf("an error occured deleting existing ingredients: %s", err.Error())
	}

	amounts := m.AmountDTO{}.ConvertAllFromDTO(*amountsDTO)

	// force set all ingredient recipeID fields to the provided recipeID
	for i := range amounts {
		amounts[i].RecipeID = recipeID
	}

	amountsResponse, err := s.repo.Create(&amounts)
	if err != nil {
		err = fmt.Errorf("an error occured creating new ingredients: %s", err.Error())

		_, createErr := s.repo.Create(existingIngredients)
		if createErr != nil {
			err = fmt.Errorf("%s\nadditionally, an second error occured restoring previous ingredients: %s",
				err.Error(),
				createErr.Error())
			return nil, err
		}
	}

	amountDTO := m.Amount{}.ConvertAllToDTO(*amountsResponse)
	return &amountDTO, nil
}

func (s AmountService) Delete(recipeID uuid.UUID) error {
	amounts, err := s.repo.Find(recipeID)
	if err != nil {
		return err
	}

	err = s.repo.Delete(amounts)
	if err != nil {
		return err
	}

	return nil
}
