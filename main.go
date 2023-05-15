package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	docs "github.com/ihulsbus/cookbook/docs"
	c "github.com/ihulsbus/cookbook/internal/config"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @contact.name	Ian Hulsbus
// @contact.url	https://github.com/ihulsbus/cookbook
// @license.name	GNU Affero General Public License v3.0
// @license.url	https://www.gnu.org/licenses/agpl-3.0.en.html
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Bearer Token
func main() {
	docs.SwaggerInfo.Title = "Cookbook API"
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.Description = "Backend API of the cookbook application. Source code and support can be found at [https://github.com/ihulsbus/cookbook](https://github.com/ihulsbus/cookbook)"
	docs.SwaggerInfo.Host = "cookbook-backend.hulsbus.be"
	docs.SwaggerInfo.BasePath = "/api/v1"

	router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	// Panic recovery
	router.Use(gin.Recovery())

	// Logging
	router.Use(c.Middleware.LoggingMiddleware())

	// CORS handler
	router.Use(cors.New(cors.Config{
		AllowOrigins:     c.Configuration.Cors.AllowedOrigins,
		AllowMethods:     c.Configuration.Cors.AllowedMethods,
		AllowHeaders:     c.Configuration.Cors.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: c.Configuration.Cors.AllowCredentials,
		MaxAge:           12 * time.Hour,
	}))

	// API versioning setup
	v1 := router.Group("/api/v1")
	v1.Use(c.Middleware.AuthMW.EnsureValidToken())
	{
		recipe := v1.Group("/recipe")
		{
			readRecipe := recipe.Group("")
			readRecipe.Use(c.Middleware.AuthMW.EnsureValidScope("read:recipes"))
			{
				readRecipe.GET("", c.RecipeEndpoints.GetAll)
				readRecipe.GET(":id", c.RecipeEndpoints.Get)
				readRecipe.GET(":id/ingredients", c.RecipeEndpoints.GetIngredientLink)
				readRecipe.GET(":id/instruction", c.RecipeEndpoints.GetInstruction)
			}

			createRecipe := recipe.Group("")
			createRecipe.Use(c.Middleware.AuthMW.EnsureValidScope("create:recipes"))
			{
				createRecipe.POST("", c.RecipeEndpoints.Create)
				createRecipe.POST(":id/instruction", c.RecipeEndpoints.CreateInstruction)
				createRecipe.POST(":id/ingredients", c.RecipeEndpoints.CreateIngredientLink)
				createRecipe.POST(":id/cover", c.RecipeEndpoints.ImageUpload)
			}

			updateRecipe := recipe.Group("")
			updateRecipe.Use(c.Middleware.AuthMW.EnsureValidScope("update:recipes"))
			{

				updateRecipe.PUT(":id", c.RecipeEndpoints.Update)
				updateRecipe.PUT(":id/instruction", c.RecipeEndpoints.UpdateInstruction)
				updateRecipe.PUT(":id/ingredients", c.RecipeEndpoints.UpdateIngredientLink)
			}

			adminRecipe := recipe.Group("")
			adminRecipe.Use(c.Middleware.AuthMW.EnsureValidScope("delete:recipes"))
			{
				adminRecipe.DELETE(":id", c.RecipeEndpoints.Delete)
				adminRecipe.DELETE(":id/instruction", c.RecipeEndpoints.DeleteInstruction)
				adminRecipe.DELETE(":id/ingredients", c.RecipeEndpoints.DeleteIngredientLink)
			}
		}

		ingredient := v1.Group("/ingredient")
		{
			readIngredient := ingredient.Group("")
			readIngredient.Use(c.Middleware.AuthMW.EnsureValidScope("read:ingredients"))
			{
				readIngredient.GET("", c.IngredientEndpoints.GetAll)
				readIngredient.GET(":id", c.IngredientEndpoints.GetSingle)
				readIngredient.GET("unit", c.IngredientEndpoints.GetUnits)
			}

			createIngredient := ingredient.Group("")
			readIngredient.Use(c.Middleware.AuthMW.EnsureValidScope("create:ingredients"))
			{
				createIngredient.POST("", c.IngredientEndpoints.Create)
			}

			updateIngredient := ingredient.Group("")
			readIngredient.Use(c.Middleware.AuthMW.EnsureValidScope("update:ingredients"))
			{
				updateIngredient.PUT(":id", c.IngredientEndpoints.Update)
			}

			adminIngredient := ingredient.Group("")
			readIngredient.Use(c.Middleware.AuthMW.EnsureValidScope("delete:ingredients"))
			{
				adminIngredient.DELETE(":id", c.IngredientEndpoints.Delete)
			}

		}

		tag := v1.Group("/tag")
		{
			readTag := tag.Group("")
			readTag.Use(c.Middleware.AuthMW.EnsureValidScope("read:tags"))
			{
				readTag.GET("", c.TagEndpoints.GetAll)
				readTag.GET(":id", c.TagEndpoints.GetSingle)
			}

			createTag := tag.Group("")
			createTag.Use(c.Middleware.AuthMW.EnsureValidScope("create:tags"))
			{
				createTag.POST("", c.TagEndpoints.Create)
			}

			updateTag := tag.Group("")
			updateTag.Use(c.Middleware.AuthMW.EnsureValidScope("update:tags"))
			{
				updateTag.PUT(":id", c.TagEndpoints.Update)
			}

			deleteTag := tag.Group("")
			deleteTag.Use(c.Middleware.AuthMW.EnsureValidScope("delete:tags"))
			{
				deleteTag.DELETE(":id", c.TagEndpoints.Delete)
			}
		}

		category := v1.Group("/category")
		{
			readCategory := category.Group("")
			readCategory.Use(c.Middleware.AuthMW.EnsureValidScope("read:categories"))
			{
				readCategory.GET("", c.CategoryEndpoints.GetAll)
				readCategory.GET(":id", c.CategoryEndpoints.GetSingle)
			}

			createCategory := category.Group("")
			createCategory.Use(c.Middleware.AuthMW.EnsureValidScope("create:categories"))
			{
				createCategory.POST("", c.CategoryEndpoints.Create)
			}

			updateCategory := category.Group("")
			updateCategory.Use(c.Middleware.AuthMW.EnsureValidScope("update:categories"))
			{
				updateCategory.PUT(":id", c.CategoryEndpoints.Update)
			}

			deleteCategory := category.Group("")
			deleteCategory.Use(c.Middleware.AuthMW.EnsureValidScope("delete:categories"))
			{
				deleteCategory.DELETE(":id", c.CategoryEndpoints.Delete)
			}
		}

	}

	// Swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	/*~~~~~~~~~~~~~~~~~~~*/

	// Server startup
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	c.Logger.Info("server available on port 8080")
	c.Logger.Fatal(srv.ListenAndServe())

}
