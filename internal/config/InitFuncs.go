package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	m "github.com/ihulsbus/cookbook/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connectS3(endpoint string, secret, key, region string) *s3.S3 {
	sess, err := session.NewSession(
		&aws.Config{
			Credentials:      credentials.NewStaticCredentials(key, secret, ""),
			Endpoint:         aws.String(endpoint),
			Region:           aws.String(region),
			S3ForcePathStyle: aws.Bool(false),
		})
	if err != nil {
		panic(err)
	}

	client := s3.New(sess)
	return client
}

func initDatabase(host string, user string, password string, dbname string, port int, sslmode string, timezone string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslmode, timezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		Logger.Fatalf("Unable to connect to the database. Exiting..\n%v\n", err)
	}

	err = db.SetupJoinTable(&m.Recipe{}, "Ingredient", &m.RecipeIngredient{})
	if err != nil {
		Logger.Errorf("Error while creating RecipeIngredient join tables: %s", err.Error())
	}

	err = db.AutoMigrate(
		&m.Recipe{},
		&m.Ingredient{},
		&m.Instruction{},
		&m.Category{},
		&m.Tag{},
	)

	if err != nil {
		Logger.Errorf("Error while automigrating database: %s", err.Error())
	}

	return db
}
