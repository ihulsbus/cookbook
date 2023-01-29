package services

import (
	"errors"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/sunshineplan/imgconv"

	m "github.com/ihulsbus/cookbook/internal/models"
	r "github.com/ihulsbus/cookbook/internal/repositories"
	u "github.com/ihulsbus/cookbook/internal/utils"
)

type RecipeService struct {
	repo        *r.RecipeRepository
	logger      *log.Logger
	imageFolder string
}

// NewRecipeService creates a new RecipeService instance
func NewRecipeService(recipeRepo *r.RecipeRepository, ImageStorePath string, logger *log.Logger) *RecipeService {
	return &RecipeService{
		repo:        recipeRepo,
		logger:      logger,
		imageFolder: ImageStorePath,
	}
}

// Find contains the business logic to get all recipes
func (s RecipeService) FindAllRecipes() ([]m.RecipeDTO, error) {
	var recipes []m.Recipe

	recipes, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return m.Recipe{}.ConvertAllToDTO(recipes), nil
}

// Find contains the business logic to get a specific recipe
func (s RecipeService) FindSingleRecipe(recipeID uint) (m.RecipeDTO, error) {
	var recipe m.Recipe

	recipe, err := s.repo.Find(recipeID)
	if err != nil {
		return recipe.ConvertToDTO(), err
	}

	return recipe.ConvertToDTO(), nil
}

// Create handles the business logic for the creation of a recipe and passes the recipe object to the recipe repo for processing
func (s RecipeService) CreateRecipe(recipe m.Recipe) (m.RecipeDTO, error) {

	recipe, err := s.repo.Create(recipe)
	if err != nil {
		return recipe.ConvertToDTO(), err
	}

	return recipe.ConvertToDTO(), nil
}

func (s RecipeService) UpdateRecipe(recipe m.Recipe) (m.RecipeDTO, error) {
	var updatedRecipe m.Recipe
	var originalRecipe m.Recipe

	originalRecipe, err := s.repo.Find(recipe.ID)
	if err != nil {
		return updatedRecipe.ConvertToDTO(), err
	}

	if recipe.Title == "" {
		recipe.Title = originalRecipe.Title
	}

	if recipe.Description == "" {
		recipe.Description = originalRecipe.Description
	}

	if recipe.PrepTime == 0 {
		recipe.PrepTime = originalRecipe.PrepTime
	}

	if recipe.CookTime == 0 {
		recipe.CookTime = originalRecipe.CookTime
	}

	if recipe.Persons == 0 {
		recipe.Persons = originalRecipe.Persons
	}

	updatedRecipe, err = s.repo.Update(recipe)
	if err != nil {
		return updatedRecipe.ConvertToDTO(), err
	}

	return updatedRecipe.ConvertToDTO(), nil
}

func (s RecipeService) UploadRecipeImages(files []m.RecipeFile) error {

	for i := range files {
		filePath := fmt.Sprintf("%s/%d/", s.imageFolder, files[i].ID)

		err := u.InitFolder(filePath)
		if err != nil {
			s.logger.Fatalf("Unable to create or detect image folder: %v", err)
		}

		file, err := files[i].File.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		out, err := os.Create(filePath + "cover.orig")
		if err != nil {
			return err
		}

		defer out.Close()
		if err != nil {
			return errors.New("unable to create the file for writing. Check your write access privilege")
		}

		_, err = io.Copy(out, file) // file not files[i] !

		if err != nil {
			return err
		}

		src, err := imgconv.Open(filePath + "cover.orig")
		if err != nil {
			os.Remove(filePath + "cover.orig")
			return errors.New("unable to open image for conversion")
		}

		coverFile, err := os.Create(filePath + "cover.jpg")
		if err != nil {
			return err
		}

		defer coverFile.Close()
		err = imgconv.Write(coverFile, src, imgconv.FormatOption{Format: imgconv.JPEG})
		if err != nil {
			return errors.New("could not convert image")
		}

		os.Remove(filePath + "cover.orig")
	}

	return nil
}

func (s RecipeService) DeleteRecipe(recipe m.Recipe) error {

	if err := s.repo.Delete(recipe); err != nil {
		return err
	}

	return nil
}
