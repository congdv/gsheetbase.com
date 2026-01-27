package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DBURL              string
	RedisURL           string
	UsageTrackWorkers  int
	GoogleClientID     string
	GoogleClientSecret string
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	return &Config{
		Port:               getEnv("WORKER_PORT", "8081"),
		DBURL:              os.Getenv("DATABASE_URL"),
		RedisURL:           os.Getenv("REDIS_URL"),
		UsageTrackWorkers:  getEnvInt("USAGE_TRACK_WORKERS", 3),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
