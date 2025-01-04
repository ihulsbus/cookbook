package models

type Config struct {
	Global   GlobalConfig
	Cors     CorsConfig
	Oauth    OauthConfig
	Database DatabaseConfig
	S3       S3Config
	RabbitMQ RabbitMQConfig
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

type S3Config struct {
	AWSRegion       string
	AWSAccessKey    string
	AWSAccessSecret string
	BucketName      string
	Endpoint        string
}

type RabbitMQConfig struct {
	Username string
	Password string
	Host     string
}

type LoggerInterface interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
