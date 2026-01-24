package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DBURL             string
	RedisURL          string
	UsageTrackWorkers int
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	return &Config{
		Port:              getEnv("WORKER_PORT", "8081"),
		DBURL:             os.Getenv("DATABASE_URL"),
		RedisURL:          os.Getenv("REDIS_URL"), // Optional - leave empty to disable rate limiting
		UsageTrackWorkers: getEnvInt("USAGE_TRACK_WORKERS", 3),
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
