package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config holds all configuration values
type Config struct {
	Port string `envconfig:"PORT" default:"8000"`

	// OAuth configuration
	TwitterClientID     string `envconfig:"TWITTER_CLIENT_ID" required:"true"`
	TwitterClientSecret string `envconfig:"TWITTER_CLIENT_SECRET" required:"true"`

	// API configuration
	APIKey string `envconfig:"API_KEY" required:"true"`
}

// Load reads the configuration from environment variables
func Load() (Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	var c Config
	err = envconfig.Process("", &c)
	if err != nil {
		return Config{}, fmt.Errorf("failed to process config: %w", err)
	}

	return c, nil
}
