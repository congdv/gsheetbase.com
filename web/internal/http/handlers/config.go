package handlers

import (
	"net/http"
	"os"

	"gsheetbase/web/internal/config"

	"github.com/gin-gonic/gin"
)

type FrontendConfig struct {
	APIBaseUrl     string `json:"apiBaseUrl"`
	WorkerBaseUrl  string `json:"workerBaseUrl"`
	LandingPageUrl string `json:"landingPageUrl"`
	ForceProd      bool   `json:"forceProd"`
}

type ConfigHandler struct {
	cfg *config.Config
}

func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{cfg: cfg}
}

func (h *ConfigHandler) GetConfig(c *gin.Context) {
	apiBaseUrl := envOrDefault("API_BASE_URL", "http://localhost:8080/api")
	workerBaseUrl := envOrDefault("WORKER_BASE_URL", "http://localhost:8081")
	landingPageUrl := envOrDefault("LANDING_PAGE_URL", "http://localhost:4321")
	forceProd := envBoolOrDefault("FORCE_PROD", false)

	config := FrontendConfig{
		APIBaseUrl:     apiBaseUrl,
		WorkerBaseUrl:  workerBaseUrl,
		LandingPageUrl: landingPageUrl,
		ForceProd:      forceProd,
	}

	c.JSON(http.StatusOK, config)
}

func envOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func envBoolOrDefault(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	switch val {
	case "1", "true", "TRUE", "True", "yes":
		return true
	case "0", "false", "FALSE", "False", "no":
		return false
	default:
		return defaultValue
	}
}
