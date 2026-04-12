package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	m "github.com/ihulsbus/cookbook/shared/models"
)

type MetadataRepository interface {
	FindAll() (*[]m.RecipeMetadata, error)
	FindSingle(recipeID uuid.UUID) (*m.RecipeMetadata, error)
	Create(meta *m.RecipeMetadata) (*m.RecipeMetadata, error)
	Update(meta *m.RecipeMetadata) (*m.RecipeMetadata, error)
	Delete(meta *m.RecipeMetadata) error
}

type RecipeRepository interface {
	RecipeExists(recipeID string) (bool, error)
}

type MetadataService struct {
	repo   MetadataRepository
	recipe RecipeRepository
}

// NewMetadataService creates a new MetadataService instance
func NewMetadataService(metadataRepo MetadataRepository, recipe RecipeRepository) *MetadataService {
	return &MetadataService{
		repo:   metadataRepo,
		recipe: recipe,
	}
}

func (s *MetadataService) FindAll() (*[]m.RecipeMetadataDTO, error) {
	var metadataResult []m.RecipeMetadataDTO

	result, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	metadataResult = m.RecipeMetadata{}.ConvertAllToDTO(*result)

	return &metadataResult, nil
}

func (s *MetadataService) Find(recipeID uuid.UUID) (*m.RecipeMetadataDTO, error) {
	var metadataResult m.RecipeMetadataDTO
	if recipeID == uuid.Nil {
		return nil, fmt.Errorf("recipe id should not be nil")
	}

	result, err := s.repo.FindSingle(recipeID)
	if err != nil {
		switch err.Error() {
		case "not found":
			return nil, fmt.Errorf("recipe not found")
		default:
			return nil, err
		}
	}

	metadataResult = result.ConvertToDTO()

	return &metadataResult, nil
}

func (s *MetadataService) Create(recipeID uuid.UUID, meta *m.RecipeMetadataDTO) (*m.RecipeMetadataDTO, error) {
	var metadata = meta.ConvertFromDTO()
	if recipeID == uuid.Nil {
		return nil, fmt.Errorf("recipe id should not be nil")
	}

	ok, err := s.recipe.RecipeExists(recipeID.String())
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("provided recipe does not exist. Cannot create metadata for a recipe that does not exist")
	}

	metadata.RecipeID = recipeID

	response, err := s.repo.Create(&metadata)
	if err != nil {
		return nil, err
	}
	metadataResponse := response.ConvertToDTO()
	return &metadataResponse, nil
}

func (s *MetadataService) Update(recipeID uuid.UUID, meta *m.RecipeMetadataDTO) (*m.RecipeMetadataDTO, error) {
	var metadata = meta.ConvertFromDTO()

	if recipeID == uuid.Nil {
		return nil, fmt.Errorf("recipe id should not be nil")
	}

	_, err := s.repo.FindSingle(recipeID)
	if err != nil {
		return nil, errors.New("provided recipe does not exist")
	}

	metadata.RecipeID = recipeID

	response, err := s.repo.Update(&metadata)
	if err != nil {
		return nil, err
	}

	metadataResponse := response.ConvertToDTO()
	return &metadataResponse, nil
}

func (s *MetadataService) Delete(recipeID uuid.UUID) error {
	if recipeID == uuid.Nil {
		return fmt.Errorf("recipe id should not be nil")
	}

	existingMetadata, err := s.repo.FindSingle(recipeID)
	if err != nil {
		return err
	}

	err = s.repo.Delete(existingMetadata)
	if err != nil {
		return err
	}

	return nil
}
