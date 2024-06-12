package config

import (
	e "instruction-service/internal/endpoints"
	h "instruction-service/internal/handlers"
	m "instruction-service/internal/models"
	r "instruction-service/internal/repositories"
	s "instruction-service/internal/services"

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
	InstructionRepository *r.InstructionRepository

	// Services
	InstructionService *s.InstructionService

	// Handlers
	InstructionHandlers *h.InstructionHandlers

	// Endpoints
	InstructionEndpoints *e.InstructionEndpoints
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
	InstructionRepository = r.NewInstructionRepository(DatabaseClient)

	// Init services
	InstructionService = s.NewInstructionService(InstructionRepository)

	// Init handlers
	InstructionHandlers = h.NewInstructionHandlers(InstructionService, Logger)

	// Init endpoints
	InstructionEndpoints = e.NewInstructionEndpoints(InstructionHandlers)
}
