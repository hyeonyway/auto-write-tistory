package config

import (
	"os"
	"strings"
)

type Config struct {
	AppEnv             string
	ServerPort         string
	DatabaseURL        string
	SeleniumRemoteURL  string
	CORSAllowedOrigins []string
}

func Load() Config {
	return Config{
		AppEnv:             getEnv("APP_ENV", "local"),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://devlog:devlog@localhost:5432/devlog?sslmode=disable"),
		SeleniumRemoteURL:  getEnv("SELENIUM_REMOTE_URL", "http://localhost:4444/wd/hub"),
		CORSAllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")),
	}
}

func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
