package config

import (
	ih "ingredient-service/internal/handlers/ingredients"
	uh "ingredient-service/internal/handlers/units"
	m "ingredient-service/internal/models"
	ir "ingredient-service/internal/repositories/ingredients"
	ur "ingredient-service/internal/repositories/units"
	is "ingredient-service/internal/services/ingredients"
	us "ingredient-service/internal/services/units"

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
	IngredientRepository *ir.IngredientRepository
	UnitRepository       *ur.UnitRepository
	// Services
	IngredientService *is.IngredientService
	UnitService       *us.UnitService

	// Handlers
	IngredientHandlers *ih.IngredientHandlers
	UnitHandlers       *uh.UnitHandlers
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
	IngredientRepository = ir.NewIngredientRepository(DatabaseClient)
	UnitRepository = ur.NewUnitRepository(DatabaseClient)

	// Init services
	IngredientService = is.NewIngredientService(IngredientRepository)
	UnitService = us.NewUnitService(UnitRepository)

	// Init handlers
	IngredientHandlers = ih.NewIngredientHandlers(IngredientService, Logger)
	UnitHandlers = uh.NewUnitHandlers(UnitService, Logger)
}
