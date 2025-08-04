// Package model provides the data models for the application.
package model

// Config is the configuration for the application.
type Config struct {
	APIServer     Server
	SwaggerServer Server
	PostgreSQL    PostgreSQL
	Services      Services
}

// Services is the configuration for external services.
type Services struct {
	Transaction Service
}

// Service is the configuration for the transaction service.
type Service struct {
	BaseURL string `yaml:"baseURL"`
}

// Server is the configuration for the server.
type Server struct {
	Enable bool
	Port   int
}

// PostgreSQL is the configuration for the PostgreSQL database.
type PostgreSQL struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DBName   string `validate:"required"`
	SSLMode  string `validate:"required"`
}
