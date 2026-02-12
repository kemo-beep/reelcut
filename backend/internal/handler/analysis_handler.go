package handler

import (
	"net/http"

	"reelcut/internal/domain"
	"reelcut/internal/middleware"
	"reelcut/internal/service"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalysisHandler struct {
	analysisSvc *service.AnalysisService
	videoSvc    *service.VideoService
}

func NewAnalysisHandler(analysisSvc *service.AnalysisService, videoSvc *service.VideoService) *AnalysisHandler {
	return &AnalysisHandler{analysisSvc: analysisSvc, videoSvc: videoSvc}
}

// Analyze godoc
// @Summary		Start analysis for a video
// @Tags			analysis
// @Produce		json
// @Security	BearerAuth
// @Param		videoId	path		string	true	"Video ID"
// @Success	202	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/analysis/videos/{videoId} [post]
func (h *AnalysisHandler) Analyze(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	videoIDStr := c.Param("videoId")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		utils.NotFound(c, "Video not found")
		return
	}
	if _, err := h.videoSvc.GetByID(c.Request.Context(), videoIDStr, userID); err != nil {
		utils.NotFound(c, "Video not found")
		return
	}
	if err := h.analysisSvc.Analyze(c.Request.Context(), videoID); err != nil {
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "Analysis started"})
}

// GetByVideoID godoc
// @Summary		Get analysis by video ID
// @Tags			analysis
// @Produce		json
// @Security	BearerAuth
// @Param		videoId	path		string	true	"Video ID"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/analysis/videos/{videoId} [get]
func (h *AnalysisHandler) GetByVideoID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	videoID := c.Param("videoId")
	if _, err := h.videoSvc.GetByID(c.Request.Context(), videoID, userID); err != nil {
		utils.NotFound(c, "Video not found")
		return
	}
	a, err := h.analysisSvc.GetByVideoID(c.Request.Context(), videoID)
	if err != nil || a == nil {
		utils.NotFound(c, "Analysis not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"analysis": a})
}

// SuggestClips godoc
// @Summary		Get AI-suggested clip segments
// @Tags			analysis
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		videoId	path		string	true	"Video ID"
// @Param		body	body		object	false	"min_duration, max_duration, max_suggestions"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/analysis/videos/{videoId}/suggest-clips [post]
func (h *AnalysisHandler) SuggestClips(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	videoID := c.Param("videoId")
	if _, err := h.videoSvc.GetByID(c.Request.Context(), videoID, userID); err != nil {
		utils.NotFound(c, "Video not found")
		return
	}
	var body struct {
		MinDuration    float64 `json:"min_duration"`
		MaxDuration    float64 `json:"max_duration"`
		MaxSuggestions int    `json:"max_suggestions"`
	}
	body.MinDuration = 7
	body.MaxDuration = 60
	body.MaxSuggestions = 20
	c.ShouldBindJSON(&body)
	if body.MaxSuggestions <= 0 {
		body.MaxSuggestions = 20
	}
	suggestions, err := h.analysisSvc.SuggestClips(c.Request.Context(), videoID, body.MinDuration, body.MaxDuration, body.MaxSuggestions)
	if err != nil {
		if err == domain.ErrNotFound {
			utils.NotFound(c, "Transcription not found")
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}
