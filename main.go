package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	c "github.com/ihulsbus/cookbook/internal/config"
)

func main() {
	router := gin.New()

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
	v1 := router.Group("/v1")
	v1.Use(c.Middleware.OidcMW.Middleware())

	/*~~~~~~~~~~~~~~~~~~~ Image folder ~~~~~~~~~~~~~~~~~~~~~*/
	router.Static("/images", c.Configuration.Global.ImageFolder)

	/*~~~~~~~~~~~~~~~~~~~ All GET routes ~~~~~~~~~~~~~~~~~~~*/
	// Recipes
	v1.GET("/recipe", c.Endpoints.RecipeGetAll)
	v1.GET("/recipe/:recipeID", c.Endpoints.RecipeGet)

	// Ingredients
	v1.GET("/ingredients", c.Endpoints.IngredientGetAll)
	v1.GET("/ingredients/:ingredientID", c.Endpoints.IngredientGetSingle)

	/*~~~~~~~~~~~~~~~~~~~ All PUT routes ~~~~~~~~~~~~~~~~~~~*/
	// Recipes
	v1.PUT("/recipe", c.Endpoints.RecipeUpdate)

	/*~~~~~~~~~~~~~~~~~~~ All POST routes ~~~~~~~~~~~~~~~~~~*/
	// Recipes
	v1.POST("/recipe", c.Endpoints.RecipeCreate)
	v1.POST("/recipe/:recipeID/upload", c.Endpoints.RecipeImageUpload)

	// Ingredients
	v1.POST("/ingredients", c.Endpoints.IngredientCreate)

	/*~~~~~~~~~~~~~~~~~~~ All DELETE routes ~~~~~~~~~~~~~~~~*/
	// Recipes
	v1.DELETE("/recipe", c.Endpoints.RecipeDelete)

	// Ingredients
	v1.DELETE("/ingredients", c.Endpoints.IngredientDelete)

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
