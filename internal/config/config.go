package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	ServerPort   int
	DatabaseURL  string
	KafkaBrokers []string
	KafkaTopic   string
}

// Load reads configuration from environment variables or .env file
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Default configuration
	config := &Config{
		ServerPort:   8080,
		DatabaseURL:  "postgres://todos:todos@db:5432/todos?sslmode=disable",
		KafkaBrokers: []string{"kafka:9092"},
		KafkaTopic:   "todos",
	}

	// Override with environment variables if present
	if port := os.Getenv("SERVER_PORT"); port != "" {
		portInt, err := strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
		}
		config.ServerPort = portInt
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.DatabaseURL = dbURL
	}

	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		config.KafkaBrokers = []string{brokers}
	}

	if topic := os.Getenv("KAFKA_TOPIC"); topic != "" {
		config.KafkaTopic = topic
	}

	return config, nil
}
