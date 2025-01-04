package config

import (
	m "instruction-service/internal/models"
	"time"

	ih "instruction-service/internal/handlers/instructions"
	sh "instruction-service/internal/handlers/search"

	ir "instruction-service/internal/repositories/instructions"
	sr "instruction-service/internal/repositories/search"

	hr "github.com/ihulsbus/cookbook/shared/httpclient"
	imr "github.com/ihulsbus/cookbook/shared/imageclient"
	rr "github.com/ihulsbus/cookbook/shared/recipeclient"

	is "instruction-service/internal/services/instructions"
	ss "instruction-service/internal/services/search"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/ihulsbus/cookbook/shared/keycloak"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	err           error
	Configuration m.Config

	Logger         *log.Logger = log.New()
	DatabaseClient *gorm.DB
	KeycloakModule *keycloak.KeycloakModule
	Cors           cors.Config

	// Repositories
	HttpClient            *hr.HTTPClient
	InstructionRepository *ir.InstructionRepository
	SearchRepository      *sr.SearchRepository
	RecipeRepository      *rr.RecipeAPIClient
	ImageRepository       *imr.ImageAPIClient

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
	KeycloakModule, err = initOauth()
	if err != nil {
		Logger.Panicf("error initialising oauth: %v", err)
	}

	// Init repositories
	HttpClient = hr.NewHTTPClient(5*time.Second, Configuration.Oauth.Url, Configuration.Oauth.Realm, Configuration.Oauth.ClientID, Configuration.Oauth.ClientSecret, Logger)
	InstructionRepository = ir.NewInstructionRepository(DatabaseClient)
	SearchRepository = sr.NewSearchRepository(DatabaseClient)
	RecipeRepository = rr.NewRecipeAPIClient("http://localhost:8081/api/v2", HttpClient)
	ImageRepository = imr.NewImageAPIClient("http://localhost:8082/api/v2", HttpClient)

	// Init services
	InstructionService = is.NewInstructionService(InstructionRepository, RecipeRepository, ImageRepository)
	SearchService = ss.NewSearchService(SearchRepository)

	// Init handlers
	InstructionHandlers = ih.NewInstructionHandlers(InstructionService, Logger)
	SearchHandlers = sh.NewSearchHandlers(SearchService, Logger)
}
