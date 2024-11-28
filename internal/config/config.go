package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	DBUsername string
	DBPassword string
	DBHost     string
	DBName     string
	Port       string
	CA         string
	UseSSL     string
}

// LoadConfig reads configuration from .env file and environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's often okay if .env doesn't exist in production
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	cfg := &Config{
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBName:     os.Getenv("DB_NAME"),
		Port:       os.Getenv("PORT"),
		CA:         os.Getenv("CA"),
		UseSSL:     os.Getenv("USE_SSL"),
	}

	// Validate required configurations
	if cfg.DBUsername == "" || cfg.DBPassword == "" ||
		cfg.DBHost == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("missing required database configuration")
	}

	return cfg, nil
}
