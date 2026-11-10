package config

import (
	"context"

	healthh "github.com/ihulsbus/cookbook/shared/healthchecks"
	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/ihulsbus/cookbook/shared/oidc"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"user-service/internal/authentik"
	handlers "user-service/internal/handlers/http"
	repositories "user-service/internal/repositories/database"
	"user-service/internal/services"
)

// userConfig embeds the common shared/models.Config (global/cors/oauth/
// database/rabbitmq) and adds the one field specific to this service:
// the Authentik service-account API token (KIT-74) user-service will use
// (KIT-33) to call Authentik's admin API and resolve user/group info.
// user-service is deliberately the only service configured with this
// token - every other service only ever validates inbound tokens or
// authenticates its own outbound calls, never calls the Authentik
// management API itself.
type userConfig struct {
	m.Config  `mapstructure:",squash"`
	Authentik struct {
		APIToken string
	}
}

var (
	Configuration userConfig
	Ctx           context.Context
	err           error

	Logger            *log.Logger = log.New()
	DatabaseClient    *gorm.DB
	OIDCModule        *oidc.Module
	Cors              cors.Config
	RabbitMQClient    *rmq.RabbitMQ
	RabbitMQPublisher *rmq.Publisher
	AuthentikClient   *authentik.Client
	UserService       *services.UserService

	// Handlers
	HealthHandler *healthh.Handlers
	HttpHandler   *handlers.HttpHandlers
)

func init() {
	Ctx = context.Background()

	initViper()
	initConfig()
	initLogging()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("config file changed: %s", e.Name)

		initConfig()
		initLogging()
	})

	if Configuration.Global.ListenPort == 0 {
		Logger.Warn("Listen port is empty. Defaulting to 8080")
		Configuration.Global.ListenPort = 8080
	}

	initDatabase()
	initCors()

	OIDCModule, err = initOauth()
	if err != nil {
		Logger.Panicf("error initialising oauth: %v", err)
	}

	RabbitMQClient, err = rmq.NewRabbitMQConnection(
		Configuration.RabbitMQ.Username,
		Configuration.RabbitMQ.Password,
		Configuration.RabbitMQ.Host,
		Logger,
	)

	// RabbitMQPublisher is for household.member.added/removed (KIT-74) -
	// published directly, not through an outbox, since nothing consumes
	// them yet and there's no in-flight request whose outcome depends on
	// delivery (same reasoning recipe-service's republish-access tool
	// uses - see recipe-service/internal/replay/replay.go).
	RabbitMQPublisher, err = rmq.NewPublisher(RabbitMQClient.Connection, "cookbook")
	if err != nil {
		Logger.Errorf("Error setting up RabbitMQ publisher: %v", err)
		Logger.Fatal("Encountered fatal error. Exiting.")
	}

	HealthHandler = healthh.NewHealthHandlers(DatabaseClient, Logger)

	// AuthentikClient calls Authentik's core API (KIT-33/34) with this
	// service's service-account token - the only service configured
	// with it (KIT-74). Its base URL is derived from the OIDC issuer
	// rather than a separate config field: Authentik always serves its
	// REST API at the same host as the issuer it hands out tokens for.
	AuthentikClient, err = authentik.NewClient(Configuration.Oauth.Issuer, Configuration.Authentik.APIToken)
	if err != nil {
		Logger.Fatalf("error initialising Authentik client: %v", err)
	}

	userRepo := repositories.NewDatabaseRepository(DatabaseClient)
	UserService = services.NewUserService(AuthentikClient, userRepo, Logger)
	HttpHandler = handlers.NewHttpHandlers(UserService, Logger)
}
