package config

import (
	"errors"
	"os"
)

// ValidateEnv checks required environment variables
func ValidateEnv() error {
	requiredVars := []string{
		"DATABASE_URL",
		"DB_PORT",
		"DB_NAME",
		"DB_PASSWORD",
		"DISCORD_CLIENT_ID",
		"DISCORD_CLIENT_SECRET",
	}

	for _, key := range requiredVars {
		if value := os.Getenv(key); value == "" {
			return errors.New("required environment variable " + key + " is not set")
		}
	}

	return nil
}
