package config

import (
	m "instruction-service/internal/models"

	ih "instruction-service/internal/handlers/instructions"
	sh "instruction-service/internal/handlers/search"

	ir "instruction-service/internal/repositories/instructions"
	sr "instruction-service/internal/repositories/search"

	is "instruction-service/internal/services/instructions"
	ss "instruction-service/internal/services/search"

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
	InstructionRepository *ir.InstructionRepository
	SearchRepository      *sr.SearchRepository

	// Services
	InstructionService *is.InstructionService
	SearchService      *ss.SearchService

	// Handlers
	InstructionHandlers *ih.InstructionHandlers
	SearchHandlers      *sh.SearchHandlers
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
	InstructionRepository = ir.NewInstructionRepository(DatabaseClient)
	SearchRepository = sr.NewSearchRepository(DatabaseClient)

	// Init services
	InstructionService = is.NewInstructionService(InstructionRepository)
	SearchService = ss.NewSearchService(SearchRepository)

	// Init handlers
	InstructionHandlers = ih.NewInstructionHandlers(InstructionService, Logger)
	SearchHandlers = sh.NewSearchHandlers(SearchService, Logger)
}
