package models

import "github.com/google/uuid"

type RecipeMetadata struct {
	RecipeID uuid.UUID  `gorm:"primaryKey"`
	Category []Category `gorm:"many2many:recipe_categories;constraint:OnDelete:CASCADE;"`

	CuisineTypeID uuid.UUID   // belongs-to CuisineType
	CuisineType   CuisineType `gorm:"foreignKey:CuisineTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	DifficultyLevelID uuid.UUID       // belongs-to DifficultyLevel
	DifficultyLevel   DifficultyLevel `gorm:"foreignKey:DifficultyLevelID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Tags              []Tag           `gorm:"many2many:recipe_tags;constraint:OnDelete:CASCADE;"`
	PreparationTime   int
	ServingCount      int
}

func (r RecipeMetadata) ConvertToDTO() RecipeMetadataDTO {
	return RecipeMetadataDTO{
		RecipeID:        r.RecipeID,
		Category:        Category{}.ConvertAllToDTO(r.Category),
		CuisineType:     r.CuisineType.ConvertToDTO(),
		DifficultyLevel: r.DifficultyLevel.ConvertToDTO(),
		Tags:            Tag{}.ConvertAllToDTO(r.Tags),
		PreparationTime: r.PreparationTime,
		ServingCount:    r.ServingCount,
	}
}

func (r RecipeMetadata) ConvertAllToDTO(metadata []RecipeMetadata) []RecipeMetadataDTO {
	var data []RecipeMetadataDTO

	for _, ri := range metadata {
		data = append(data, ri.ConvertToDTO())
	}

	return data
}

type RecipeMetadataDTO struct {
	RecipeID        uuid.UUID          `json:"recipe_id" example:"23582396-12a3-425b-a597-8a22052823da"`
	Category        []CategoryDTO      `gorm:"categories"`
	CuisineType     CuisineTypeDTO     `json:"cuisine_type"`
	Tags            []TagDTO           `json:"tags"`
	DifficultyLevel DifficultyLevelDTO `json:"difficulty_level"`
	PreparationTime int                `json:"preparation_time" example:"55"`
	ServingCount    int                `json:"serving_count" example:"4"`
}

func (r RecipeMetadataDTO) ConvertFromDTO() RecipeMetadata {
	return RecipeMetadata{
		RecipeID:        r.RecipeID,
		Category:        CategoryDTO{}.ConvertAllFromDTO(r.Category),
		CuisineType:     r.CuisineType.ConvertFromDTO(),
		Tags:            TagDTO{}.ConvertAllFromDTO(r.Tags),
		PreparationTime: r.PreparationTime,
		ServingCount:    r.ServingCount,
	}
}

func (r RecipeMetadataDTO) ConvertAllFromDTO(metadataDTO []RecipeMetadataDTO) []RecipeMetadata {
	var data []RecipeMetadata

	for _, ri := range metadataDTO {
		data = append(data, ri.ConvertFromDTO())
	}

	return data
}
