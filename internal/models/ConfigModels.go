package models

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"gorm.io/gorm"
)

type Config struct {
	Global          GlobalConfig
	Cors            CorsConfig
	Oidc            OidcConfig
	Database        DatabaseConfig
	S3              S3Config
	DatabaseClient  *gorm.DB
	S3ClientSession *s3.S3
}

// GlobalConfig holds global configuration items
type GlobalConfig struct {
	Debug       bool
	ImageFolder string
}

// DatabaseConfig holds database configuration items
type DatabaseConfig struct {
	Host     string
	Username string
	Password string
	Database string
	Port     int
	SSLMode  string
	Timezone string
}

type OidcConfig struct {
	URL               string
	ClientID          string
	SigningAlgs       []string
	SkipClientIDCheck bool
	SkipExpiryCheck   bool
	SkipIssuerCheck   bool
}

type CorsConfig struct {
	AllowedOrigins   []string
	AllowCredentials bool
	AllowedHeaders   []string
	AllowedMethods   []string
	Debug            bool
}

type S3Config struct {
	AWSRegion       string
	AWSAccessKey    string
	AWSAccessSecret string
	BucketName      string
	Endpoint        string
}
