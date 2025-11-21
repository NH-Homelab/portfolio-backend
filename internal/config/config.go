package config

import (
	"github.com/joho/godotenv"

	"fmt"
	"log"
	"os"
	"strconv"
)

type BackendConfig struct {
	db_host string
	db_port int
}

func Load() (*BackendConfig, error) {
	err := godotenv.Load()
	if err != nil {
		// Not a fatal error - just means we'll use environment variables
		log.Println("WARNING: No .env file found, using environment variables")
	}

	db_port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse integer database port value: %w", err)
	}

	return &BackendConfig{
		db_host: getEnv("DB_HOST", "localhost"),
		db_port: db_port,
	}, nil
}

// getEnv retrieves an environment variable or returns a fallback value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
