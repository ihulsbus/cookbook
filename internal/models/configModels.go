package models

import "gorm.io/gorm"

type Config struct {
	Global         GlobalConfig
	Cors           CorsConfig
	Oidc           OidcConfig
	Database       DatabaseConfig
	DatabaseClient *gorm.DB
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
