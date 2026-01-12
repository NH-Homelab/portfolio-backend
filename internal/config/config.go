package config

import (
	"github.com/joho/godotenv"

	"log"
	"os"
	"strings"
)

type BackendConfig struct {
	Db_host     string
	Db_port     string
	Db_user     string
	Db_password string
	Db_name     string
	Allowed_origins []string
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
		Allowed_origins: parseAllowedOrigins("ALLOWED_ORIGINS", "https://www.nickhenley.dev"),
	}, nil
}

// getEnv retrieves an environment variable or returns a fallback value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// parseAllowedOrigins reads an env var and returns a slice of origins.
// The env var may be a comma-separated list. If not set, returns the
// provided default as a single-element slice.
func parseAllowedOrigins(key, defaultOrigin string) []string {
	if v, ok := os.LookupEnv(key); ok {
		// split by comma and trim spaces
		parts := strings.Split(v, ",")
		var out []string
		for _, p := range parts {
			t := strings.TrimSpace(p)
			if t != "" {
				out = append(out, t)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	return []string{defaultOrigin}
}
