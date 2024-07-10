package models

type Config struct {
	Global   GlobalConfig
	Cors     CorsConfig
	Oauth    OauthConfig
	Database DatabaseConfig
	S3       S3Config
}

// GlobalConfig holds global configuration items
type GlobalConfig struct {
	LogLevel string
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
	Service              string
	Url                  string
	Realm                string
	FullCertsPath        *string
	DisableSecurityCheck bool
}

type CorsConfig struct {
	AllowedOrigins   []string
	AllowCredentials bool
	AllowedHeaders   []string
	AllowedMethods   []string
}

type S3Config struct {
	AWSRegion       string
	AWSAccessKey    string
	AWSAccessSecret string
	BucketName      string
	Endpoint        string
}

type LoggerInterface interface {
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}
