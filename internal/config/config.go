package config

import (
	"github.com/joho/godotenv"

	"log"
	"os"
)

type BackendConfig struct {
	Db_host     string
	Db_port     string
	Db_user     string
	Db_password string
	Db_name     string
}

func Load() (*BackendConfig, error) {
	err := godotenv.Load()
	if err != nil {
		// Not a fatal error - just means we'll use environment variables
		log.Println("WARNING: No .env file found, using environment variables")
	}

	return &BackendConfig{
		Db_host:     getEnv("DB_HOST", "localhost"),
		Db_port:     getEnv("DB_PORT", "5432"),
		Db_user:     getEnv("DB_USER", "postgres"),
		Db_password: getEnv("DB_PASSWORD", "password"),
		Db_name:     getEnv("DB_NAME", "postgres"),
	}, nil
}

// getEnv retrieves an environment variable or returns a fallback value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
