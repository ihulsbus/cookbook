package config

import (
	"fmt"
	"ingredient-service/internal/helpers"
	"strings"
	"time"

	az "github.com/ihulsbus/cookbook/shared/authz"
	"github.com/ihulsbus/cookbook/shared/models"
	"gorm.io/gorm/clause"

	"github.com/gin-contrib/cors"
	"github.com/ihulsbus/cookbook/shared/oidc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initLogging() {
	log.Info("setting up the logging framework")

	Logger.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	logLevels := helpers.SetupLogLevels()

	if i, found := logLevels[strings.ToUpper(Configuration.Global.LogLevel)]; found {
		Logger.SetLevel(i)
		Logger.Infof("loglevel set to %s", strings.ToUpper(Logger.Level.String()))

	} else {
		Logger.Warn("no or invalid loglevel specified. Assuming default value of INFO. \n valid loglevels are: PANIC FATAL ERROR WARN INFO DEBUG TRACE")
		Logger.SetLevel(logLevels["INFO"])
	}
}

func initViper() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
}

func initConfig() {
	Logger.Info("loading config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Logger.Fatalf("config file not found: %v", err)
		} else {
			Logger.Fatalf("unknown error occured while reading config. error: %v", err)
		}
	}

	if err := viper.Unmarshal(&Configuration); err != nil {
		Logger.Fatalf("error unmarshaling config: %v", err)
	}

	Logger.Info("config file loaded")
}

func initDatabase() {
	Logger.Info("connecting to the database")
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		Configuration.Database.Host,
		Configuration.Database.Username,
		Configuration.Database.Password,
		Configuration.Database.Database,
		Configuration.Database.Port,
		Configuration.Database.SSLMode,
		Configuration.Database.Timezone)

	DatabaseClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		Logger.Fatalf("Unable to connect to the database. Exiting..\n%v\n", err)
	}

	Logger.Info("performing database migrations")
	if err := DatabaseClient.AutoMigrate(
		&models.Ingredient{},
		&models.Unit{},
		&models.Amount{},
	); err != nil {
		Logger.Fatalf("Error while automigrating database: %s", err.Error())
	}
	// recipe_access is kept current by Projector's RabbitMQ consumer, not
	// by calling recipe-service at startup.
	if err := az.AutoMigrate(DatabaseClient); err != nil {
		Logger.Fatalf("Error while automigrating recipe_access table: %s", err.Error())
	}

	initIngredientDuplicatePrevention()

	Logger.Info("connected!")
}

// initIngredientDuplicatePrevention uses raw SQL, not GORM AutoMigrate,
// since a STORED GENERATED column isn't something GORM's migrator can
// declare. normalized_name is deliberately not a Go struct field on
// models.Ingredient; the repository queries it by raw column name.
//
// Idempotent (every statement is CREATE ... IF NOT EXISTS / CREATE OR
// REPLACE), so it's safe on every startup. Must run after AutoMigrate
// (needs the "ingredients" table and "name" column to exist).
func initIngredientDuplicatePrevention() {
	Logger.Info("setting up ingredient duplicate prevention")

	// unaccent and pg_trgm are "trusted" extensions (Postgres 13+), so
	// this service's own DB role can create them without superuser.
	if err := DatabaseClient.Exec(`CREATE EXTENSION IF NOT EXISTS unaccent`).Error; err != nil {
		Logger.Fatalf("Error creating unaccent extension: %s", err.Error())
	}
	if err := DatabaseClient.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm`).Error; err != nil {
		Logger.Fatalf("Error creating pg_trgm extension: %s", err.Error())
	}

	// unaccent() is only STABLE, not IMMUTABLE, so Postgres refuses it in
	// a generated column expression. This wrapper pins the dictionary and
	// asserts IMMUTABLE - the standard documented workaround.
	if err := DatabaseClient.Exec(`
		CREATE OR REPLACE FUNCTION immutable_unaccent(text) RETURNS text AS $$
			SELECT unaccent('unaccent', $1)
		$$ LANGUAGE sql IMMUTABLE PARALLEL SAFE
	`).Error; err != nil {
		Logger.Fatalf("Error creating immutable_unaccent function: %s", err.Error())
	}

	// Lowercased, unaccented, whitespace-collapsed - so "Tomaat ",
	// "tomaat", and "TOMAAT" all normalize to the same value.
	if err := DatabaseClient.Exec(`
		ALTER TABLE ingredients ADD COLUMN IF NOT EXISTS normalized_name text
			GENERATED ALWAYS AS (lower(immutable_unaccent(trim(regexp_replace(name, '\s+', ' ', 'g'))))) STORED
	`).Error; err != nil {
		Logger.Fatalf("Error adding normalized_name column: %s", err.Error())
	}

	// DB-enforced exact-duplicate prevention; IngredientRepository.Create
	// also handles the race this closes.
	if err := DatabaseClient.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS ingredients_normalized_name_key ON ingredients (normalized_name)
	`).Error; err != nil {
		Logger.Fatalf("Error creating normalized_name unique index: %s", err.Error())
	}

	// GIN trigram index so the "did you mean ...?" similarity() queries
	// don't scan the whole table.
	if err := DatabaseClient.Exec(`
		CREATE INDEX IF NOT EXISTS ingredients_normalized_name_trgm_idx ON ingredients USING gin (normalized_name gin_trgm_ops)
	`).Error; err != nil {
		Logger.Fatalf("Error creating normalized_name trigram index: %s", err.Error())
	}
}

func initUnits() {
	if err := DatabaseClient.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(models.DefaultUnits).Error; err != nil {
			return err
		}

		return nil

	}); err != nil {
		Logger.Fatalf("Error while creating units: %v", err)
	}
}

func initCors() {
	Cors = cors.Config{
		AllowOrigins:     Configuration.Cors.AllowedOrigins,
		AllowMethods:     Configuration.Cors.AllowedMethods,
		AllowHeaders:     Configuration.Cors.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: Configuration.Cors.AllowCredentials,
		MaxAge:           12 * time.Hour,
	}
}

func initOauth() (*oidc.Module, error) {
	module, err := oidc.NewModule(Ctx, oidc.Config{
		Issuer: Configuration.Oauth.Issuer,
		// Accepts both the frontend's public client ID (end-user tokens)
		// and this service's own confidential client ID (service calls).
		Audiences: []string{Configuration.Oauth.Audience, Configuration.Oauth.ClientID},
	})
	if err != nil {
		return nil, err
	}

	return module, nil
}
