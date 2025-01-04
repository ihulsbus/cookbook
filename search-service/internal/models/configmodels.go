package models

type Config struct {
	Global GlobalConfig
	Cors   CorsConfig
	Oauth  OauthConfig
}

// GlobalConfig holds global configuration items
type GlobalConfig struct {
	LogLevel string
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

type LoggerInterface interface {
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}
