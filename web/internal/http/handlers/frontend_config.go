package handlers

import (
	"net/http"

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

func NewFrontendConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{cfg: cfg}
}

func (h *ConfigHandler) GetConfig(c *gin.Context) {

	config := FrontendConfig{
		APIBaseUrl:     h.cfg.FrontendApiBaseUrl,
		WorkerBaseUrl:  h.cfg.FrontendWorkerBaseUrl,
		LandingPageUrl: h.cfg.FrontendLandingPageUrl,
		ForceProd:      h.cfg.FrontEendForceProd,
	}

	c.JSON(http.StatusOK, config)
}
