package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	c "github.com/ihulsbus/cookbook/internal/config"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	httplogger "github.com/gleicon/go-httplogger"
)

func main() {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(loggingMiddleware())

	// API versioning setup
	v1 := router.Group("/v1")

	/*~~~~~~~~~~~~~~~~~~~ Image folder ~~~~~~~~~~~~~~~~~~~~~*/
	// imageRouter := router.PathPrefix("/images/").Subrouter()

	// fs := http.FileServer(http.Dir(c.Configuration.Global.ImageFolder))
	// imageRouter.NewRoute().Handler(http.StripPrefix("/images/", fs))

	/*~~~~~~~~~~~~~~~~~~~ All GET routes ~~~~~~~~~~~~~~~~~~~*/
	// Recipes
	v1.GET("/recipe", c.Endpoints.RecipeGetAll)
	v1.GET("/recipe/{recipeID}", c.Endpoints.RecipeGet)

	// Ingredients
	v1.GET("/ingredients", c.Endpoints.IngredientGetAll)
	v1.GET("/ingredients/{ingredientID}", c.Endpoints.IngredientGetSingle)

	/*~~~~~~~~~~~~~~~~~~~ All PUT routes ~~~~~~~~~~~~~~~~~~~*/
	// Recipes
	v1.PUT("/recipe", c.Endpoints.RecipeUpdate)

	/*~~~~~~~~~~~~~~~~~~~ All POST routes ~~~~~~~~~~~~~~~~~~~*/
	// Recipes
	v1.POST("/recipe", c.Endpoints.RecipeCreate)
	v1.POST("/recipe/{recipeID}/upload", c.Endpoints.RecipeImageUpload)

	// Ingredients
	v1.POST("/ingredients", c.Endpoints.IngredientCreate)

	/*~~~~~~~~~~~~~~~~~~~ All DELETE routes ~~~~~~~~~~~~~~~~~~~*/
	// Recipes
	v1.POST("/recipe", c.Endpoints.RecipeDelete)

	// Ingredients
	v1.POST("/ingredients", c.Endpoints.IngredientDelete)

	/*~~~~~~~~~~~~~~~~~~~*/

	// CORS handler
	crs := cors.New(cors.Options{
		AllowedOrigins:   c.Configuration.Cors.AllowedOrigins,
		AllowedHeaders:   c.Configuration.Cors.AllowedHeaders,
		AllowedMethods:   c.Configuration.Cors.AllowedMethods,
		AllowCredentials: c.Configuration.Cors.AllowCredentials,
		Debug:            c.Configuration.Cors.Debug,
	})
	handler := httplogger.HTTPLogger(crs.Handler(router))

	// Server startup
	srv := &http.Server{
		Handler:      handler,
		Addr:         ":8080",
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	c.Logger.Info("server available on port 8080")
	c.Logger.Fatal(srv.ListenAndServe())

}

func loggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c.Logger.WithFields(log.Fields{"remote": ctx.Request.RemoteAddr, "method": ctx.Request.Method, "uri": ctx.Request.RequestURI})

		ctx.Next()
	}
}
