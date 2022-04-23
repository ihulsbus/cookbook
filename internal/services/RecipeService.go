package services

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	m "github.com/ihulsbus/cookbook/internal/models"
	r "github.com/ihulsbus/cookbook/internal/repositories"
	u "github.com/ihulsbus/cookbook/internal/utils"
)

type FindAllRecipes func(s *RecipeService) ([]m.RecipeDTO, error)
type FindSingleRecipe func(s *RecipeService, recipeID int) (m.Recipe, error)
type CreateRecipe func(s *RecipeService, recipe m.RecipeDTO) (m.RecipeDTO, error)
type UpdateRecipe func(s *RecipeService, recipe m.RecipeDTO) (m.RecipeDTO, error)
type UploadRecipeImages func(s *RecipeService, files []m.RecipeFile) error
type DeleteRecipe func(s *RecipeService, recipe m.RecipeDTO) error

type RecipeService struct {
	repo        *r.RecipeRepository
	imageFolder string

	FindAllRecipes     FindAllRecipes
	FindSingleRecipe   FindSingleRecipe
	CreateRecipe       CreateRecipe
	UpdateRecipe       UpdateRecipe
	UploadRecipeImages UploadRecipeImages
	DeleteRecipe       DeleteRecipe
}

// NewRecipeService creates a new RecipeService instance
func NewRecipeService(recipeRepo *r.RecipeRepository, ImageStorePath string) *RecipeService {
	return &RecipeService{
		repo:        recipeRepo,
		imageFolder: ImageStorePath,

		FindAllRecipes:     findAllRecipes,
		FindSingleRecipe:   findSingleRecipe,
		CreateRecipe:       createRecipe,
		UpdateRecipe:       updateRecipe,
		UploadRecipeImages: uploadRecipeImages,
		DeleteRecipe:       deleteRecipe,
	}
}

// Find contains the business logic to get all recipes
func findAllRecipes(s *RecipeService) ([]m.RecipeDTO, error) {
	var recipes []m.Recipe

	recipes, err := s.repo.FindAll(s.repo)
	if err != nil {
		return nil, err
	}

	return m.Recipe{}.ConvertAllToDTO(recipes), nil
}

// Find contains the business logic to get a specific recipe
func findSingleRecipe(s *RecipeService, recipeID int) (m.Recipe, error) {
	var recipe m.Recipe

	recipe, err := s.repo.Find(s.repo, recipeID)
	if err != nil {
		return recipe, err
	}

	return recipe, nil
}

// Create handles the business logic for the creation of a recipe and passes the recipe object to the recipe repo for processing
func createRecipe(s *RecipeService, recipe m.RecipeDTO) (m.RecipeDTO, error) {

	response, err := s.repo.Create(s.repo, recipe.ConvertFromDTO())
	if err != nil {
		return response.ConvertToDTO(), err
	}

	return response.ConvertToDTO(), nil
}

func updateRecipe(s *RecipeService, recipe m.RecipeDTO) (m.RecipeDTO, error) {
	var updatedRecipe m.Recipe
	var originalRecipe m.Recipe

	originalRecipe, err := s.repo.Find(s.repo, recipe.ID)
	if err != nil {
		return recipe, err
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

	if recipe.Amount_Persons == 0 {
		recipe.Amount_Persons = originalRecipe.Amount_Persons
	}

	updatedRecipe, err = s.repo.Update(s.repo, recipe.ConvertFromDTO())
	if err != nil {
		return updatedRecipe.ConvertToDTO(), err
	}

	return updatedRecipe.ConvertToDTO(), nil
}

func uploadRecipeImages(s *RecipeService, files []m.RecipeFile) error {

	for i := range files {
		filePath := fmt.Sprintf("%s/%d/", s.imageFolder, files[i].ID)

		err := u.InitFolder(filePath)
		if err != nil {
			log.Fatal("Unable to create or detect image folder: %v", err)
		}

		file, err := files[i].File.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		out, err := os.Create(filePath + files[i].File.Filename)
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
	}
	// defer file.
	// 	initRecipeImageFolder(filePath)

	// f, err := os.OpenFile(filePath+fileHeader.Filename, os.O_WRONLY|os.O_CREATE, 0666)

	return nil
}

func deleteRecipe(s *RecipeService, recipe m.RecipeDTO) error {

	if err := s.repo.Delete(s.repo, recipe.ConvertFromDTO()); err != nil {
		return err
	}

	return nil
}
