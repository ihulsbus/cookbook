package main

import (
	"net/http"
	"time"

	c "github.com/ihulsbus/cookbook/internal/config"
	h "github.com/ihulsbus/cookbook/internal/handlers"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"remote_addr": r.RemoteAddr, "method": r.Method, "uri": r.RequestURI}).Info()

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()
	omw := h.OidcMW{}
	router.Use(loggingMiddleware)

	// API versioning setup
	v1 := router.PathPrefix("/v1").Subrouter()

	/*~~~~~~~~~~~~~~~~~~~ All GET routes ~~~~~~~~~~~~~~~~~~~*/
	v1get := v1.Methods("GET").Subrouter()
	v1get.Use(omw.Middleware)

	// Recipes
	v1get.Path("/recipe").HandlerFunc(h.RecipeGetAll)
	v1get.Path("/recipe/{recipeID}").HandlerFunc(h.RecipeGet)
	v1get.Path("/recipe/{recipeID}/ingredients").HandlerFunc(h.RecipeIngredientGet)

	// Ingredients
	v1get.Path("/ingredient").HandlerFunc(h.IngredientGetAll)

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

	// Ingredients
	v1post.Path("/ingredient").HandlerFunc(h.IngredientCreate)

	/*~~~~~~~~~~~~~~~~~~~ All DELETE routes ~~~~~~~~~~~~~~~~~~~*/
	v1del := v1.Methods("DELETE").Subrouter()
	v1del.Use(omw.Middleware)

	// Recipes
	v1del.Path("/recipe").HandlerFunc(h.RecipeDelete)

	// Ingredients
	v1del.Path("/ingredients").HandlerFunc(h.NotImplemented)

	/*~~~~~~~~~~~~~~~~~~~*/
	// Server startup
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Info("server available on port 8080")
	log.Debug(c.Configuration)
	log.Fatal(srv.ListenAndServe())

}
