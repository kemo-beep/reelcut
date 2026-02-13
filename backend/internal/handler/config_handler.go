package handler

import (
	"net/http"

	"reelcut/internal/config"
	"reelcut/internal/service"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct{}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

// GetCaptionFonts returns the list of allowed caption fonts for styling.
// Used by the editor to populate font dropdown; only allowed fonts are accepted by the backend when burning ASS.
func (h *ConfigHandler) GetCaptionFonts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"fonts": service.AllowedCaptionFonts})
}

// GetExportPresets returns platform export presets (TikTok/Reels, Instagram, YouTube Shorts, etc.).
func (h *ConfigHandler) GetExportPresets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"presets": config.ExportPresets})
}
