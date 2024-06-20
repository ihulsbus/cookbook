package config

import (
	h "ingredient-service/internal/handlers"
	m "ingredient-service/internal/models"
	r "ingredient-service/internal/repositories"
	s "ingredient-service/internal/services"

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
	IngredientRepository *r.IngredientRepository

	// Services
	IngredientService *s.IngredientService

	// Handlers
	IngredientHandlers *h.IngredientHandlers
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
	IngredientRepository = r.NewIngredientRepository(DatabaseClient)

	// Init services
	IngredientService = s.NewIngredientService(IngredientRepository)

	// Init handlers
	IngredientHandlers = h.NewIngredientHandlers(IngredientService, Logger)
}
