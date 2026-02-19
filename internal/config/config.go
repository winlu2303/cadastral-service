package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port              string
	DatabaseURL       string
	LogLevel          string
	Environment       string
	DocsEnabled       bool
	Auth              AuthConfig
	ExternalServerURL string
}

type AuthConfig struct {
	Enabled   bool
	JWTSecret string
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://cadastral:cadastral123@localhost:5432/cadastral_db?sslmode=disable"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		Environment:       getEnv("ENVIRONMENT", "development"),
		DocsEnabled:       getEnvBool("DOCS_ENABLED", true),
		ExternalServerURL: getEnv("EXTERNAL_SERVER_URL", ""),
		Auth: AuthConfig{
			Enabled:   getEnvBool("AUTH_ENABLED", false),
			JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}
