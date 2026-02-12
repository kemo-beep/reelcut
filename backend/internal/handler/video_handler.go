package handler

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"reelcut/internal/domain"
	"reelcut/internal/middleware"
	"reelcut/internal/service"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const thumbnailTokenExpiry = 5 * time.Minute

type VideoHandler struct {
	videoSvc       *service.VideoService
	jwtSecret      string
	s3RedirectHost string // allowed host for download-shared-object redirects (e.g. localhost:9002)
}

func NewVideoHandler(videoSvc *service.VideoService, jwtSecret string, s3Endpoint string) *VideoHandler {
	redirectHost := ""
	if u, err := url.Parse(s3Endpoint); err == nil {
		redirectHost = u.Host
	}
	return &VideoHandler{videoSvc: videoSvc, jwtSecret: jwtSecret, s3RedirectHost: redirectHost}
}

// Upload godoc
// @Summary		Get presigned URL for video upload
// @Tags			videos
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		body	body		object	true	"project_id, filename"
// @Success	202	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/videos/upload [post]
func (h *VideoHandler) Upload(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	var body struct {
		ProjectID string `json:"project_id" binding:"required"`
		Filename  string `json:"filename" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, []utils.ErrorDetail{{Message: err.Error()}})
		return
	}
	uploadURL, videoID, err := h.videoSvc.GetPresignedUploadURL(c.Request.Context(), userID, body.ProjectID, body.Filename)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			utils.NotFound(c, "Project not found")
			return
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			utils.ValidationError(c, []utils.ErrorDetail{{Field: ve.Field, Message: ve.Message}})
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"video": gin.H{
			"id":                 videoID,
			"project_id":         body.ProjectID,
			"original_filename":  body.Filename,
			"status":             "uploading",
		},
		"upload": gin.H{
			"upload_url": uploadURL,
			"method":     "PUT",
		},
	})
}

// ConfirmUpload godoc
// @Summary		Confirm video upload completed
// @Tags			videos
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Video ID"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/videos/{id}/confirm [post]
func (h *VideoHandler) ConfirmUpload(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	videoID := c.Param("id")
	vid, err := uuid.Parse(videoID)
	if err != nil {
		utils.NotFound(c, "Video not found")
		return
	}
	if err := h.videoSvc.ConfirmUpload(c.Request.Context(), vid); err != nil {
		if err == domain.ErrNotFound {
			utils.NotFound(c, "Video not found")
			return
		}
		if err == domain.ErrInsufficientCredits {
			utils.Error(c, http.StatusPaymentRequired, "INSUFFICIENT_CREDITS", "Insufficient credits", nil)
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Upload confirmed", "video_id": videoID})
}

// UploadFromURL godoc
// @Summary		Upload video from URL (not implemented)
// @Tags			videos
// @Security	BearerAuth
// @Failure	501	{object}	object
// @Router		/api/v1/videos/upload/url [post]
func (h *VideoHandler) UploadFromURL(c *gin.Context) {
	if middleware.GetUserID(c) == "" {
		utils.Unauthorized(c, "")
		return
	}
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// InitiateResumableUpload godoc
// @Summary		Start resumable multipart video upload
// @Tags			videos
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		body	body		object	true	"project_id, filename"
// @Success	200	{object}	object	"video_id, upload_id"
// @Router		/api/v1/videos/upload/resumable [post]
func (h *VideoHandler) InitiateResumableUpload(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	var body struct {
		ProjectID string `json:"project_id" binding:"required"`
		Filename  string `json:"filename" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, []utils.ErrorDetail{{Message: err.Error()}})
		return
	}
	uploadID, videoID, err := h.videoSvc.InitiateResumableUpload(c.Request.Context(), userID, body.ProjectID, body.Filename)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			utils.NotFound(c, "Project not found")
			return
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			utils.ValidationError(c, []utils.ErrorDetail{{Field: ve.Field, Message: ve.Message}})
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"video_id": videoID, "upload_id": uploadID})
}

// UploadPartResumable godoc
// @Summary		Upload one part of a resumable upload
// @Tags			videos
// @Accept		application/octet-stream
// @Produce		json
// @Security	BearerAuth
// @Param		id			path		string	true	"Video ID"
// @Param		upload_id	query		string	true	"Upload ID from initiate"
// @Param		partNumber	path		int		true	"Part number (1-based)"
// @Param		body		body		binary	true	"Part binary data"
// @Success	200	{object}	object	"etag"
// @Router		/api/v1/videos/{id}/upload/parts/{partNumber} [put]
func (h *VideoHandler) UploadPartResumable(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	videoID := c.Param("id")
	uploadID := c.Query("upload_id")
	if uploadID == "" {
		utils.ValidationError(c, []utils.ErrorDetail{{Field: "upload_id", Message: "required"}})
		return
	}
	partNumStr := c.Param("partNumber")
	partNum, err := strconv.Atoi(partNumStr)
	if err != nil || partNum < 1 {
		utils.ValidationError(c, []utils.ErrorDetail{{Field: "partNumber", Message: "must be positive integer"}})
		return
	}
	etag, err := h.videoSvc.UploadPart(c.Request.Context(), videoID, userID, uploadID, partNum, c.Request.Body)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			utils.NotFound(c, "Video not found")
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			utils.ValidationError(c, []utils.ErrorDetail{{Message: "invalid upload_id or request"}})
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"etag": etag})
}

// CompleteResumableUpload godoc
// @Summary		Complete resumable upload and start processing
// @Tags			videos
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Video ID"
// @Param		body	body		object	true	"parts: [{part_number, etag}]"
// @Success	202	{object}	object
// @Router		/api/v1/videos/{id}/upload/complete [post]
func (h *VideoHandler) CompleteResumableUpload(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	videoID := c.Param("id")
	var body struct {
		Parts []service.ResumablePart `json:"parts" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, []utils.ErrorDetail{{Message: err.Error()}})
		return
	}
	if err := h.videoSvc.CompleteResumableUpload(c.Request.Context(), videoID, userID, body.Parts); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			utils.NotFound(c, "Video not found")
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			utils.ValidationError(c, []utils.ErrorDetail{{Message: "invalid or missing upload_id/parts"}})
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "Upload complete", "video_id": videoID})
}

// List godoc
// @Summary		List videos
// @Tags			videos
// @Produce		json
// @Security	BearerAuth
// @Param		page		query		int		false	"Page"
// @Param		per_page	query		int		false	"Per page"
// @Param		project_id	query		string	false	"Filter by project"
// @Param		status		query		string	false	"Filter by status"
// @Param		sort_by		query		string	false	"created_at"
// @Param		sort_order	query		string	false	"asc|desc"
// @Success	200	{object}	object
// @Router		/api/v1/videos [get]
func (h *VideoHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	projectID := c.Query("project_id")
	status := c.Query("status")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	var pID *string
	if projectID != "" {
		pID = &projectID
	}
	var st *string
	if status != "" {
		st = &status
	}
	list, total, err := h.videoSvc.List(c.Request.Context(), userID, pID, st, page, perPage, sortBy, sortOrder)
	if err != nil {
		utils.Internal(c, "")
		return
	}
	// Use download-shared-object proxy URL so <img> loads without CORS (base64-encoded presigned URL)
	baseURL := requestBaseURL(c)
	entries := make([]videoListEntry, 0, len(list))
	for _, v := range list {
		e := videoListEntry{Video: v}
		if v.ThumbnailURL != nil && *v.ThumbnailURL != "" {
			if presigned, err := h.videoSvc.GetPresignedDownloadURL(c.Request.Context(), *v.ThumbnailURL, 300); err == nil {
				encoded := base64.RawURLEncoding.EncodeToString([]byte(presigned))
				e.ThumbnailDisplayURL = baseURL + "/api/v1/download-shared-object/" + encoded
			}
		}
		entries = append(entries, e)
	}
	utils.JSONPaginated(c, http.StatusOK, gin.H{"videos": entries}, page, perPage, total)
}

// requestBaseURL returns the request's base URL (scheme + host) for building absolute thumbnail URLs.
func requestBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}
	return scheme + "://" + host
}

// videoListEntry adds thumbnail_display_url (API URL with token) for list responses
type videoListEntry struct {
	*domain.Video
	ThumbnailDisplayURL string `json:"thumbnail_display_url,omitempty"`
}

// GetByID godoc
// @Summary		Get video by ID
// @Tags			videos
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Video ID"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/videos/{id} [get]
func (h *VideoHandler) GetByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	videoID := c.Param("id")
	v, err := h.videoSvc.GetByID(c.Request.Context(), videoID, userID)
	if err != nil {
		utils.NotFound(c, "Video not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"video": v})
}

// Delete godoc
// @Summary		Delete a video
// @Tags			videos
// @Security	BearerAuth
// @Param		id	path		string	true	"Video ID"
// @Success	204	"No Content"
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/videos/{id} [delete]
func (h *VideoHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	videoID := c.Param("id")
	if err := h.videoSvc.Delete(c.Request.Context(), videoID, userID); err != nil {
		if err == domain.ErrNotFound {
			utils.NotFound(c, "Video not found")
			return
		}
		utils.Internal(c, "")
		return
	}
	c.Status(http.StatusNoContent)
}

// GetMetadata godoc
// @Summary		Get video metadata
// @Tags			videos
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Video ID"
// @Success	200	{object}	object
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/videos/{id}/metadata [get]
func (h *VideoHandler) GetMetadata(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	videoID := c.Param("id")
	v, err := h.videoSvc.GetByID(c.Request.Context(), videoID, userID)
	if err != nil {
		utils.NotFound(c, "Video not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"duration_seconds": v.DurationSeconds,
		"width":             v.Width,
		"height":            v.Height,
		"fps":               v.FPS,
		"codec":             v.Codec,
		"bitrate":           v.Bitrate,
		"file_size_bytes":   v.FileSizeBytes,
	})
}

// GetPlaybackURL godoc
// @Summary		Get presigned URL for video playback
// @Tags			videos
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Video ID"
// @Success	200	{object}	object	"url"
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/videos/{id}/playback-url [get]
func (h *VideoHandler) GetPlaybackURL(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	videoID := c.Param("id")
	v, err := h.videoSvc.GetByID(c.Request.Context(), videoID, userID)
	if err != nil {
		utils.NotFound(c, "Video not found")
		return
	}
	if v.StoragePath == "" {
		// Video exists but not ready for playback (e.g. still processing)
		status := v.Status
		if status == "" {
			status = "processing"
		}
		c.JSON(http.StatusOK, gin.H{"url": nil, "status": status})
		return
	}
	url, err := h.videoSvc.GetPresignedDownloadURL(c.Request.Context(), v.StoragePath, 3600) // 1h
	if err != nil {
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// GetThumbnail godoc
// @Summary		Redirect to video thumbnail URL
// @Tags			videos
// @Security	BearerAuth
// @Param		id	path		string	true	"Video ID"
// @Success	302	"Redirect to thumbnail"
// @Failure	404	{object}	utils.ErrorResponse
// @Router		/api/v1/videos/{id}/thumbnail [get]
func (h *VideoHandler) GetThumbnail(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	videoID := c.Param("id")
	v, err := h.videoSvc.GetByID(c.Request.Context(), videoID, userID)
	if err != nil || v.ThumbnailURL == nil || *v.ThumbnailURL == "" {
		utils.NotFound(c, "Video or thumbnail not found")
		return
	}
	// Redirect to presigned URL (thumbnail is stored as S3 key)
	url, err := h.videoSvc.GetPresignedDownloadURL(c.Request.Context(), *v.ThumbnailURL, 900) // 15 min
	if err != nil {
		utils.Internal(c, "")
		return
	}
	c.Redirect(http.StatusFound, url)
}

// DownloadSharedObject decodes a base64-encoded URL and redirects to it. Used so <img> can load
// thumbnails via the API (no CORS to MinIO). Only redirects to the configured S3 endpoint host.
// @Summary		Redirect to a shared object (presigned URL passed as base64 path)
// @Tags			videos
// @Param		encoded	path		string	true	"Base64-encoded presigned URL"
// @Success	302	"Redirect to object"
// @Failure	400	{object}	utils.ErrorResponse
// @Router		/api/v1/download-shared-object/{encoded} [get]
func (h *VideoHandler) DownloadSharedObject(c *gin.Context) {
	encoded := strings.TrimPrefix(c.Param("encoded"), "/")
	if encoded == "" {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Missing encoded URL", nil)
		return
	}
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid encoded URL", nil)
		return
	}
	redirectURL, err := url.Parse(string(decoded))
	if err != nil || !redirectURL.IsAbs() {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid URL", nil)
		return
	}
	// Only allow redirect to our S3/MinIO host (e.g. localhost:9002 or 127.0.0.1:9000)
	if h.s3RedirectHost == "" {
		utils.Internal(c, "")
		return
	}
	redirectHost := redirectURL.Hostname()
	redirectPort := redirectURL.Port()
	allowedPort := ""
	if idx := strings.Index(h.s3RedirectHost, ":"); idx >= 0 {
		allowedPort = h.s3RedirectHost[idx+1:]
	}
	// Allow localhost or 127.0.0.1 with our S3 port
	allowedHosts := []string{"localhost", "127.0.0.1", strings.Split(h.s3RedirectHost, ":")[0]}
	hostOK := false
	for _, ah := range allowedHosts {
		if redirectHost == ah {
			hostOK = true
			break
		}
	}
	if !hostOK || (allowedPort != "" && redirectPort != allowedPort) {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Redirect not allowed", nil)
		return
	}
	c.Redirect(http.StatusFound, redirectURL.String())
}
