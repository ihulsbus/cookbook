package config

import (
	e "recipe-service/internal/endpoints"
	h "recipe-service/internal/handlers"
	m "recipe-service/internal/models"
	r "recipe-service/internal/repositories"
	s "recipe-service/internal/services"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	Configuration m.Config

	Logger         *log.Logger = log.New()
	DatabaseClient *gorm.DB
	Cors           cors.Config

	// Repositories
	RecipeRepository *r.RecipeRepository

	// Services
	RecipeService *s.RecipeService

	// Handlers
	RecipeHandlers *h.RecipeHandlers

	// Endpoints
	RecipeEndpoints *e.RecipeEndpoints
)

func init() {
	initViper()
	initConfig()
	initLogging()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("config file changed: %s", e.Name)

		initConfig()
		initLogging()
	})

	initDatabase()
	initCors()

	// Init repositories
	RecipeRepository = r.NewRecipeRepository(DatabaseClient)

	// Init services
	RecipeService = s.NewRecipeService(RecipeRepository)

	// Init handlers
	RecipeHandlers = h.NewRecipeHandlers(RecipeService, Logger)

	// Init endpoints
	RecipeEndpoints = e.NewRecipeEndpoints(RecipeHandlers)
}
