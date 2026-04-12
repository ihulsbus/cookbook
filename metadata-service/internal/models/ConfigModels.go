package models

type Config struct {
	Global       GlobalConfig
	Cors         CorsConfig
	Oauth        OauthConfig
	Database     DatabaseConfig
	RecipeClient ApiClient
}

// GlobalConfig holds global configuration items
type GlobalConfig struct {
	LogLevel   string
	ListenPort string
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

type OauthConfig struct {
	ClientID     string
	ClientSecret string
	Audience     string
	Url          string
	Realm        string
}

type CorsConfig struct {
	AllowedOrigins   []string
	AllowCredentials bool
	AllowedHeaders   []string
	AllowedMethods   []string
}

type ApiClient struct {
	BaseURL string
}

type LoggerInterface interface {
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}
