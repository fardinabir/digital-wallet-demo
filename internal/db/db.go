// Package db provides the database connection and migration functionality.
package db

import (
	"fmt"
	"github.com/fardinabir/digital-wallet-demo/internal/model"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New creates a new database connection
func New(cfg model.PostgreSQL) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = Migrate(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// NewTestDB creates a new test database connection
func NewTestDB() (*gorm.DB, error) {
	cfg, err := initTestConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize test config: %w", err)
	}

	// Connect to test database
	db, err := New(cfg.PostgreSQL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// initTestConfig loads test configuration from config.test.yaml
func initTestConfig() (*model.Config, error) {
	v := viper.New()
	v.SetConfigName("config.test")
	v.SetConfigType("yaml")

	// Add config path relative to the current directory
	v.AddConfigPath("../../")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read test config: %w", err)
	}

	var cfg model.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal test config: %w", err)
	}

	return &cfg, nil
}
