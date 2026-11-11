package models

type Config struct {
	Global   GlobalConfig
	Cors     CorsConfig
	Oauth    OauthConfig
	Database DatabaseConfig
	S3       S3Config
	RabbitMQ RabbitMQConfig
}

type GlobalConfig struct {
	LogLevel   string
	ListenPort int
	// InternalListenPort, when set, serves a service's internal-only API
	// not meant to be reachable through the ingress.
	InternalListenPort int
}

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
	// Issuer and Audience validate inbound tokens.
	Issuer   string
	Audience string
	// ClientID/ClientSecret authenticate this service's own outbound
	// requests via a client-credentials grant; not used for inbound
	// validation.
	ClientID     string
	ClientSecret string
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

type ApiClient struct {
	BaseURL string
}

type LoggerInterface interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
