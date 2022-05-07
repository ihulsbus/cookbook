package main

import (
	"net/http"
	"time"

	c "github.com/ihulsbus/cookbook/internal/config"
	h "github.com/ihulsbus/cookbook/internal/handlers"
	"github.com/rs/cors"

	httplogger "github.com/gleicon/go-httplogger"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := mux.NewRouter()
	omw := h.OidcMW{}

	// API versioning setup
	v1 := router.PathPrefix("/v1").Subrouter()

	/*~~~~~~~~~~~~~~~~~~~ Image folder ~~~~~~~~~~~~~~~~~~~~~*/
	imageRouter := router.PathPrefix("/images/").Subrouter()
	// imageRouter.Use(omw.Middleware)

	fs := http.FileServer(http.Dir(c.Configuration.Global.ImageFolder))
	imageRouter.NewRoute().Handler(http.StripPrefix("/images/", fs))

	/*~~~~~~~~~~~~~~~~~~~ All GET routes ~~~~~~~~~~~~~~~~~~~*/
	v1get := v1.Methods("GET").Subrouter()
	v1get.Use(omw.Middleware)

	// Recipes
	v1get.Path("/recipe").HandlerFunc(h.RecipeGetAll)
	v1get.Path("/recipe/{recipeID}").HandlerFunc(h.RecipeGet)

	// Ingredients
	v1get.Path("/ingredients").HandlerFunc(h.IngredientGetAll)
	v1get.Path("/ingredients/{ingredientID}").HandlerFunc(h.IngredientGetSingle)

	/*~~~~~~~~~~~~~~~~~~~ All PUT routes ~~~~~~~~~~~~~~~~~~~*/
	v1put := v1.Methods("PUT").Subrouter()
	v1put.Use(omw.Middleware)

	// Recipes
	v1put.Path("/recipe").HandlerFunc(h.RecipeUpdate)

	/*~~~~~~~~~~~~~~~~~~~ All POST routes ~~~~~~~~~~~~~~~~~~~*/
	v1post := v1.Methods("POST").Subrouter()
	v1post.Use(omw.Middleware)

	// Recipes
	v1post.Path("/recipe").HandlerFunc(h.RecipeCreate)
	v1post.Path("/recipe/{recipeID}/upload").HandlerFunc(h.RecipeImageUpload)

	// Ingredients
	v1post.Path("/ingredients").HandlerFunc(h.IngredientCreate)

	/*~~~~~~~~~~~~~~~~~~~ All DELETE routes ~~~~~~~~~~~~~~~~~~~*/
	v1del := v1.Methods("DELETE").Subrouter()
	v1del.Use(omw.Middleware)

	// Recipes
	v1del.Path("/recipe").HandlerFunc(h.RecipeDelete)

	// Ingredients
	v1del.Path("/ingredients").HandlerFunc(h.NotImplemented)

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

	log.Info("server available on port 8080")
	log.Fatal(srv.ListenAndServe())

}
