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
	docs.SwaggerInfo.Host = "http://localhost:8080"
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
	v1.Use(c.Middleware.OidcMW.Middleware())
	{
		recipe := v1.Group("/recipe")
		{
			recipe.GET("", c.RecipeEndpoints.GetAll)
			recipe.GET(":recipeID", c.RecipeEndpoints.Get)
			recipe.POST("", c.RecipeEndpoints.Create)
			recipe.POST(":recipeID/cover", c.RecipeEndpoints.ImageUpload)
			recipe.PUT(":recipeID", c.RecipeEndpoints.Update)
			recipe.DELETE(":recipeID", c.RecipeEndpoints.Delete)
		}

		ingredient := v1.Group("/ingredient")
		{
			ingredient.GET("", c.IngredientEndpoints.GetAll)
			ingredient.GET(":ingredientID", c.IngredientEndpoints.GetSingle)
			ingredient.POST("", c.IngredientEndpoints.Create)
			ingredient.PUT(":ingredientID", c.IngredientEndpoints.Update)
			ingredient.DELETE(":ingredientID", c.IngredientEndpoints.Delete)
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
