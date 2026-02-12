package handler

import (
	"net/http"
	"strconv"

	"reelcut/internal/domain"
	"reelcut/internal/middleware"
	"reelcut/internal/service"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClipHandler struct {
	clipSvc  *service.ClipService
	videoSvc *service.VideoService
}

func NewClipHandler(clipSvc *service.ClipService, videoSvc *service.VideoService) *ClipHandler {
	return &ClipHandler{clipSvc: clipSvc, videoSvc: videoSvc}
}

// Create godoc
// @Summary		Create a clip
// @Tags			clips
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		body	body		object	true	"video_id, name, start_time, end_time, aspect_ratio, ..."
// @Success	201	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/clips [post]
func (h *ClipHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	uid, _ := uuid.Parse(userID)
	var body struct {
		VideoID       string   `json:"video_id" binding:"required"`
		Name          string   `json:"name" binding:"required"`
		StartTime     float64  `json:"start_time" binding:"required"`
		EndTime       float64  `json:"end_time" binding:"required"`
		AspectRatio   string   `json:"aspect_ratio"`
		ViralityScore *float64 `json:"virality_score"`
		FromSuggestion string  `json:"from_suggestion"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, []utils.ErrorDetail{{Message: err.Error()}})
		return
	}
	clip, err := h.clipSvc.Create(c.Request.Context(), uid, body.VideoID, body.Name, body.StartTime, body.EndTime, body.AspectRatio, body.ViralityScore, body.FromSuggestion != "")
	if err != nil {
		if err == domain.ErrNotFound {
			utils.NotFound(c, "Video not found")
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"clip": clip})
}

// List godoc
// @Summary		List clips
// @Tags			clips
// @Produce		json
// @Security	BearerAuth
// @Param		page		query		int		false	"Page"
// @Param		per_page	query		int		false	"Per page"
// @Param		video_id	query		string	false	"Filter by video"
// @Param		status		query		string	false	"Filter by status"
// @Success	200	{object}	object
// @Router		/api/v1/clips [get]
func (h *ClipHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	videoID := c.Query("video_id")
	status := c.Query("status")
	var vID, st *string
	if videoID != "" {
		vID = &videoID
	}
	if status != "" {
		st = &status
	}
	list, total, err := h.clipSvc.List(c.Request.Context(), userID, vID, st, page, perPage, c.DefaultQuery("sort_by", "created_at"), c.DefaultQuery("sort_order", "desc"))
	if err != nil {
		utils.Internal(c, "")
		return
	}
	utils.JSONPaginated(c, http.StatusOK, gin.H{"clips": list}, page, perPage, total)
}

// GetByID godoc
// @Summary		Get clip by ID
// @Tags			clips
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Clip ID"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/clips/{id} [get]
func (h *ClipHandler) GetByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	clipID := c.Param("id")
	clip, err := h.clipSvc.GetByID(c.Request.Context(), clipID, userID)
	if err != nil {
		utils.NotFound(c, "Clip not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"clip": clip})
}

// GetPlaybackURL godoc
// @Summary		Get presigned URL for clip video playback (cut file)
// @Tags			clips
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Clip ID"
// @Success	200	{object}	object	"url"
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/clips/{id}/playback-url [get]
func (h *ClipHandler) GetPlaybackURL(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	clipID := c.Param("id")
	clip, err := h.clipSvc.GetByID(c.Request.Context(), clipID, userID)
	if err != nil {
		utils.NotFound(c, "Clip not found")
		return
	}
	if clip.StoragePath == nil || *clip.StoragePath == "" {
		c.JSON(http.StatusOK, gin.H{"url": nil})
		return
	}
	url, err := h.videoSvc.GetPresignedDownloadURL(c.Request.Context(), *clip.StoragePath, 3600)
	if err != nil {
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// Update godoc
// @Summary		Update a clip
// @Tags			clips
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id		path		string	true	"Clip ID"
// @Param		body	body		object	true	"name, start_time, end_time, aspect_ratio"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/clips/{id} [put]
func (h *ClipHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	clipID := c.Param("id")
	clip, err := h.clipSvc.GetByID(c.Request.Context(), clipID, userID)
	if err != nil {
		utils.NotFound(c, "Clip not found")
		return
	}
	var body struct {
		Name        *string  `json:"name"`
		StartTime   *float64 `json:"start_time"`
		EndTime     *float64 `json:"end_time"`
		AspectRatio *string  `json:"aspect_ratio"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, nil)
		return
	}
	if body.Name != nil {
		clip.Name = *body.Name
	}
	if body.StartTime != nil {
		clip.StartTime = *body.StartTime
	}
	if body.EndTime != nil {
		clip.EndTime = *body.EndTime
	}
	if body.AspectRatio != nil {
		clip.AspectRatio = *body.AspectRatio
	}
	if clip.EndTime > clip.StartTime {
		dur := clip.EndTime - clip.StartTime
		clip.DurationSeconds = &dur
	}
	if err := h.clipSvc.Update(c.Request.Context(), clip); err != nil {
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"clip": clip})
}

// Delete godoc
// @Summary		Delete a clip
// @Tags			clips
// @Security	BearerAuth
// @Param		id	path		string	true	"Clip ID"
// @Success	204	"No Content"
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/clips/{id} [delete]
func (h *ClipHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	if err := h.clipSvc.Delete(c.Request.Context(), c.Param("id"), userID); err != nil {
		utils.NotFound(c, "Clip not found")
		return
	}
	c.Status(http.StatusNoContent)
}

// Render godoc
// @Summary		Start clip render job
// @Tags			clips
// @Security	BearerAuth
// @Param		id	path		string	true	"Clip ID"
// @Failure	401	{object}	utils.ErrorResponse
// @Failure	501	{object}	object
// @Router		/api/v1/clips/{id}/render [post]
func (h *ClipHandler) Render(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	clipID := c.Param("id")
	jobID, err := h.clipSvc.StartRender(c.Request.Context(), clipID, userID)
	if err != nil {
		if err == domain.ErrInsufficientCredits {
			utils.Error(c, http.StatusPaymentRequired, "INSUFFICIENT_CREDITS", "Insufficient credits", nil)
			return
		}
		utils.NotFound(c, "Clip not found")
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message":   "Render started",
		"job_id":    jobID,
		"status_url": "/api/v1/jobs/" + jobID,
	})
}

// GetRenderStatus godoc
// @Summary		Get clip render job status
// @Tags			clips
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Clip ID"
// @Success	200	{object}	object
// @Router		/api/v1/clips/{id}/status [get]
func (h *ClipHandler) GetRenderStatus(c *gin.Context) {
	if middleware.GetUserID(c) == "" {
		utils.Unauthorized(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"render_job": gin.H{"status": "pending"}})
}

// Download godoc
// @Summary		Download rendered clip
// @Tags			clips
// @Security	BearerAuth
// @Param		id	path		string	true	"Clip ID"
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/clips/{id}/download [get]
func (h *ClipHandler) Download(c *gin.Context) {
	if middleware.GetUserID(c) == "" {
		utils.Unauthorized(c, "")
		return
	}
	utils.NotFound(c, "Clip not ready for download")
}

// Duplicate godoc
// @Summary		Duplicate a clip
// @Tags			clips
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Clip ID"
// @Success	201	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/clips/{id}/duplicate [post]
func (h *ClipHandler) Duplicate(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	clip, err := h.clipSvc.Duplicate(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		utils.NotFound(c, "Clip not found")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"clip": clip})
}

// GetStyle godoc
// @Summary		Get clip style
// @Tags			clips
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Clip ID"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/clips/{id}/style [get]
func (h *ClipHandler) GetStyle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	clipID := c.Param("id")
	style, err := h.clipSvc.GetStyle(c.Request.Context(), clipID, userID)
	if err != nil {
		utils.NotFound(c, "Clip not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"style": style})
}

// UpdateStyle godoc
// @Summary		Update clip style
// @Tags			clips
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id		path		string	true	"Clip ID"
// @Param		body	body		domain.ClipStyle	true	"Style config"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/clips/{id}/style [put]
func (h *ClipHandler) UpdateStyle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	clipID := c.Param("id")
	var body domain.ClipStyle
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, nil)
		return
	}
	if err := h.clipSvc.UpdateStyle(c.Request.Context(), clipID, userID, &body); err != nil {
		utils.NotFound(c, "Clip not found")
		return
	}
	style, _ := h.clipSvc.GetStyle(c.Request.Context(), clipID, userID)
	c.JSON(http.StatusOK, gin.H{"style": style})
}

// GetCaptionsSRT godoc
// @Summary		Get clip captions as SRT
// @Tags			clips
// @Produce		text/plain
// @Security	BearerAuth
// @Param		id	path		string	true	"Clip ID"
// @Success	200	"SRT file"
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/clips/{id}/captions/srt [get]
func (h *ClipHandler) GetCaptionsSRT(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	clipID := c.Param("id")
	srt, err := h.clipSvc.GetCaptionsSRT(c.Request.Context(), clipID, userID)
	if err != nil {
		utils.NotFound(c, "Clip or transcription not found")
		return
	}
	c.Header("Content-Disposition", "attachment; filename=captions.srt")
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(srt))
}

// GetCaptionsVTT godoc
// @Summary		Get clip captions as WebVTT
// @Tags			clips
// @Produce		text/vtt
// @Security	BearerAuth
// @Param		id	path		string	true	"Clip ID"
// @Success	200	"VTT file"
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/clips/{id}/captions/vtt [get]
func (h *ClipHandler) GetCaptionsVTT(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	clipID := c.Param("id")
	vtt, err := h.clipSvc.GetCaptionsVTT(c.Request.Context(), clipID, userID)
	if err != nil {
		utils.NotFound(c, "Clip or transcription not found")
		return
	}
	c.Header("Content-Disposition", "attachment; filename=captions.vtt")
	c.Data(http.StatusOK, "text/vtt; charset=utf-8", []byte(vtt))
}

// ApplyTemplate godoc
// @Summary		Apply template to clip style
// @Tags			clips
// @Security	BearerAuth
// @Param		id			path		string	true	"Clip ID"
// @Param		templateId	path		string	true	"Template ID"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/clips/{id}/style/apply-template/{templateId} [post]
func (h *ClipHandler) ApplyTemplate(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	clipID := c.Param("id")
	templateID := c.Param("templateId")
	if err := h.clipSvc.ApplyTemplate(c.Request.Context(), clipID, userID, templateID); err != nil {
		utils.NotFound(c, "Clip or template not found")
		return
	}
	style, _ := h.clipSvc.GetStyle(c.Request.Context(), clipID, userID)
	c.JSON(http.StatusOK, gin.H{"style": style})
}

// CancelRender godoc
// @Summary		Cancel clip render job
// @Tags			clips
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Clip ID"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/clips/{id}/cancel [post]
func (h *ClipHandler) CancelRender(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	clipID := c.Param("id")
	if err := h.clipSvc.CancelRender(c.Request.Context(), clipID, userID); err != nil {
		utils.NotFound(c, "Clip or job not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Render cancelled"})
}
