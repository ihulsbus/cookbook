package services

import (
	"fmt"
	"github.com/google/uuid"
	m "metadata-service/internal/models"
)

type MetadataRepository interface {
	FindAll() ([]m.RecipeMetadata, error)
	FindSingle(recipeID uuid.UUID) (m.RecipeMetadata, error)
	Create(meta m.RecipeMetadata) (m.RecipeMetadata, error)
	Update(meta m.RecipeMetadata) (m.RecipeMetadata, error)
	Delete(meta m.RecipeMetadata) error
}
type MetadataService struct {
	repo MetadataRepository
}

// NewMetadataService creates a new MetadataService instance
func NewMetadataService(metadataRepo MetadataRepository) *MetadataService {
	return &MetadataService{
		repo: metadataRepo,
	}
}

func (m *MetadataService) FindAll() (*[]m.RecipeMetadataDTO, error) {
	var metadataResult []m.RecipeMetadataDTO

	result, err := m.repo.FindAll()
	if err != nil {
		return nil, err
	}

	metadataResult = m.RecipeMetadata{}.ConvertAllToDTO(result)
	return &metadataResult, nil
}

func (m *MetadataService) Find(recipeID uuid.UUID) (*m.RecipeMetadataDTO, error) {
	var metadataResult m.RecipeMetadataDTO
	if recipeID == uuid.Nil {
		return nil, fmt.Errorf("recipe id should not be nil")
	}

	result, err := m.repo.FindSingle(recipeID)
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

//
//func (m *MetadataService) Create(meta m.RecipeMetadataDTO) (*m.RecipeMetadataDTO, error) {}
//
//func (m *MetadataService) Update(meta m.RecipeMetadataDTO) (*m.RecipeMetadataDTO, error) {}
//
//func (m *MetadataService) Delete(meta m.RecipeMetadataDTO) (*m.RecipeMetadataDTO, error) {}
