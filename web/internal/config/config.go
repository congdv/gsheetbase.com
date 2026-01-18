package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	DBURL               string
	JWTAccessSecret     string
	JWTAccessTTLMinutes int
	FrontendOrigin      string
	CookieDomain        string
	CookieSecure        bool

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectUrl  string
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	cfg := &Config{
		Port:                env("PORT", "8080"),
		DBURL:               env("DATABASE_URL", env("DB_URL", "")),
		JWTAccessSecret:     env("JWT_ACCESS_SECRET", ""),
		FrontendOrigin:      env("FRONTEND_ORIGIN", ""),
		CookieDomain:        env("COOKIE_DOMAIN", "localhost"),
		CookieSecure:        envBool("COOKIE_SECURE", false),
		JWTAccessTTLMinutes: envInt("JWT_ACCESS_TTL_MINUTES", 15),

		GoogleClientID:     env("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: env("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectUrl:  env("GOOGLE_REDIRECT_URL", ""),
	}
	return cfg, nil
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func envBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		switch v {
		case "1", "true", "TRUE", "True", "yes":
			return true
		case "0", "false", "FALSE", "False", "no":
			return false
		}
	}
	return def
}
