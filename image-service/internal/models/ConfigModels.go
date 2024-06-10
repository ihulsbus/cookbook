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
	Domain   string
	ClientID string
	Audience string
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
