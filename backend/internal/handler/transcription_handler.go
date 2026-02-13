package handler

import (
	"errors"
	"net/http"

	"reelcut/internal/domain"
	"reelcut/internal/middleware"
	"reelcut/internal/service"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TranscriptionHandler struct {
	transcriptionSvc *service.TranscriptionService
}

func NewTranscriptionHandler(transcriptionSvc *service.TranscriptionService) *TranscriptionHandler {
	return &TranscriptionHandler{transcriptionSvc: transcriptionSvc}
}

// Create godoc
// @Summary		Start transcription for a video
// @Tags			transcriptions
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		videoId	path		string	true	"Video ID"
// @Param		body	body		object	false	"language, enable_diarization"
// @Success	202	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/transcriptions/videos/{videoId} [post]
func (h *TranscriptionHandler) Create(c *gin.Context) {
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
	var body struct {
		Language          string `json:"language"`
		EnableDiarization bool   `json:"enable_diarization"`
	}
	body.Language = "en"
	c.ShouldBindJSON(&body)
	t, err := h.transcriptionSvc.Create(c.Request.Context(), videoID, body.Language, body.EnableDiarization)
	if err != nil {
		if err == domain.ErrNotFound {
			utils.NotFound(c, "Video not found")
			return
		}
		if err == domain.ErrValidation {
			utils.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Video is not ready for transcription", nil)
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"transcription": t,
		"job":           gin.H{"status_url": "/api/v1/transcriptions/" + t.ID.String()},
	})
}

// GetByID godoc
// @Summary		Get transcription by ID
// @Tags			transcriptions
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Transcription ID"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/transcriptions/{id} [get]
func (h *TranscriptionHandler) GetByID(c *gin.Context) {
	if middleware.GetUserID(c) == "" {
		utils.Unauthorized(c, "")
		return
	}
	id := c.Param("id")
	t, err := h.transcriptionSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.NotFound(c, "Transcription not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"transcription": t})
}

// GetByVideoID godoc
// @Summary		Get transcription by video ID
// @Tags			transcriptions
// @Produce		json
// @Security	BearerAuth
// @Param		videoId	path		string	true	"Video ID"
// @Success	200	{object}	object	"Returns { transcription: <object> } or { transcription: null } when none"
// @Router		/api/v1/transcriptions/videos/{videoId} [get]
func (h *TranscriptionHandler) GetByVideoID(c *gin.Context) {
	if middleware.GetUserID(c) == "" {
		utils.Unauthorized(c, "")
		return
	}
	videoID := c.Param("videoId")
	language := c.Query("language")
	var t *domain.Transcription
	var err error
	if language != "" {
		t, err = h.transcriptionSvc.GetByVideoIDAndLanguage(c.Request.Context(), videoID, language)
	} else {
		t, err = h.transcriptionSvc.GetByVideoID(c.Request.Context(), videoID)
	}
	if err != nil {
		if err == domain.ErrNotFound || errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusOK, gin.H{"transcription": nil})
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"transcription": t})
}

// UpdateSegment godoc
// @Summary		Update a transcript segment
// @Tags			transcriptions
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id			path		string	true	"Transcription ID"
// @Param		segmentId	path		string	true	"Segment ID"
// @Param		body		body		object	true	"text, start_time, end_time"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/transcriptions/{id}/segments/{segmentId} [put]
func (h *TranscriptionHandler) UpdateSegment(c *gin.Context) {
	if middleware.GetUserID(c) == "" {
		utils.Unauthorized(c, "")
		return
	}
	transcriptionID := c.Param("id")
	segmentID := c.Param("segmentId")
	var body struct {
		Text      string   `json:"text"`
		StartTime float64  `json:"start_time"`
		EndTime   float64  `json:"end_time"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, nil)
		return
	}
	if err := h.transcriptionSvc.UpdateSegment(c.Request.Context(), transcriptionID, segmentID, body.Text, body.StartTime, body.EndTime); err != nil {
		if err == domain.ErrNotFound {
			utils.NotFound(c, "Not found")
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// Translate godoc
// @Summary		Translate transcription to another language
// @Tags			transcriptions
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Source transcription ID"
// @Param		body	body		object	true	"target_language (ISO 639-1)"
// @Success	201	{object}	object	"Returns new transcription with same timestamps, translated text"
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/transcriptions/{id}/translate [post]
func (h *TranscriptionHandler) Translate(c *gin.Context) {
	if middleware.GetUserID(c) == "" {
		utils.Unauthorized(c, "")
		return
	}
	id := c.Param("id")
	var body struct {
		TargetLanguage string `json:"target_language"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.TargetLanguage == "" {
		utils.ValidationError(c, []utils.ErrorDetail{{Field: "target_language", Message: "required"}})
		return
	}
	t, err := h.transcriptionSvc.Translate(c.Request.Context(), id, body.TargetLanguage)
	if err != nil {
		if err == domain.ErrNotFound {
			utils.NotFound(c, "Transcription not found")
			return
		}
		if err == domain.ErrValidation {
			utils.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Cannot translate this transcription", nil)
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"transcription": t})
}

// ListByVideoID godoc
// @Summary		List all completed transcriptions for a video (for caption language selection)
// @Tags			transcriptions
// @Produce		json
// @Security	BearerAuth
// @Param		videoId	path		string	true	"Video ID"
// @Success	200	{object}	object	"Returns { transcriptions: [ { id, language, ... } ] }"
// @Router		/api/v1/transcriptions/videos/{videoId}/list [get]
func (h *TranscriptionHandler) ListByVideoID(c *gin.Context) {
	if middleware.GetUserID(c) == "" {
		utils.Unauthorized(c, "")
		return
	}
	videoID := c.Param("videoId")
	list, err := h.transcriptionSvc.ListCompletedByVideoID(c.Request.Context(), videoID)
	if err != nil {
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"transcriptions": list})
}
