package handler

import (
	"errors"
	"net/http"
	"strconv"

	"reelcut/internal/domain"
	"reelcut/internal/middleware"
	"reelcut/internal/repository"
	"reelcut/internal/service"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo     repository.UserRepository
	usageLogRepo repository.UsageLogRepository
	authSvc      *service.AuthService
	userSvc      *service.UserService
}

func NewUserHandler(userRepo repository.UserRepository, usageLogRepo repository.UsageLogRepository, authSvc *service.AuthService, userSvc *service.UserService) *UserHandler {
	return &UserHandler{userRepo: userRepo, usageLogRepo: usageLogRepo, authSvc: authSvc, userSvc: userSvc}
}

// GetProfile godoc
// @Summary		Get current user profile
// @Tags			users
// @Produce		json
// @Security	BearerAuth
// @Success	200	{object}	object
// @Failure	401	{object}	utils.ErrorResponse
// @Router		/api/v1/users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		utils.Unauthorized(c, "")
		return
	}
	sanitizeUser(user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateProfile godoc
// @Summary		Update current user profile
// @Tags			users
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		body	body		object	true	"full_name, avatar_url"
// @Success	200	{object}	object
// @Failure	401	{object}	utils.ErrorResponse
// @Router		/api/v1/users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		utils.Unauthorized(c, "")
		return
	}
	var body struct {
		FullName *string `json:"full_name"`
		AvatarURL *string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, nil)
		return
	}
	if body.FullName != nil {
		user.FullName = body.FullName
	}
	if body.AvatarURL != nil {
		user.AvatarURL = body.AvatarURL
	}
	if err := h.userRepo.Update(c.Request.Context(), user); err != nil {
		utils.Internal(c, "")
		return
	}
	sanitizeUser(user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetUsageStats godoc
// @Summary		Get current user usage stats
// @Tags			users
// @Produce		json
// @Security	BearerAuth
// @Param		page		query		int	false	"Page"
// @Param		per_page	query		int	false	"Per page"
// @Success	200	{object}	object
// @Failure	401	{object}	utils.ErrorResponse
// @Router		/api/v1/users/me/usage [get]
func (h *UserHandler) GetUsageStats(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	logs, total, err := h.usageLogRepo.ListByUserID(c.Request.Context(), userID, perPage, offset)
	if err != nil {
		utils.Internal(c, "")
		return
	}
	utils.JSONPaginated(c, http.StatusOK, gin.H{"usage_logs": logs}, page, perPage, total)
}

// ChangePassword godoc
// @Summary		Change password (logged-in user)
// @Tags			users
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		body	body		object	true	"current_password, new_password"
// @Success	200	{object}	object
// @Failure	400	{object}	utils.ErrorResponse
// @Failure	401	{object}	utils.ErrorResponse
// @Router		/api/v1/users/me/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	var body struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, []utils.ErrorDetail{{Field: "body", Message: "current_password and new_password required"}})
		return
	}
	err := h.authSvc.ChangePassword(c.Request.Context(), userID, body.CurrentPassword, body.NewPassword)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.Unauthorized(c, "Invalid current password")
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
	c.JSON(http.StatusOK, gin.H{"message": "Password updated"})
}

// UploadAvatar godoc
// @Summary		Upload avatar image
// @Tags			users
// @Accept		multipart/form-data
// @Produce		json
// @Security	BearerAuth
// @Param		file	formData	file	true	"Image (JPEG, PNG, WebP; max 5MB)"
// @Success	200	{object}	object
// @Failure	400	{object}	utils.ErrorResponse
// @Failure	401	{object}	utils.ErrorResponse
// @Router		/api/v1/users/me/avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		utils.ValidationError(c, []utils.ErrorDetail{{Field: "file", Message: "file is required"}})
		return
	}
	f, err := file.Open()
	if err != nil {
		utils.Internal(c, "")
		return
	}
	defer f.Close()
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	user, err := h.userSvc.UploadAvatar(c.Request.Context(), userID, f, contentType, file.Size)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			utils.NotFound(c, "User not found")
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
	sanitizeUser(user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetAvatar godoc
// @Summary		Redirect to avatar image URL
// @Tags			users
// @Security	BearerAuth
// @Success	302	"Redirect to presigned URL"
// @Failure	404	"No avatar"
// @Router		/api/v1/users/me/avatar [get]
func (h *UserHandler) GetAvatar(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.Unauthorized(c, "")
		return
	}
	url, err := h.userSvc.GetAvatarURL(c.Request.Context(), userID)
	if err != nil {
		utils.NotFound(c, "Avatar not found")
		return
	}
	c.Redirect(http.StatusFound, url)
}

// DeleteAccount godoc
// @Summary		Delete current user account
// @Tags			users
// @Produce		json
// @Security	BearerAuth
// @Success	200	{object}	object
// @Failure	401	{object}	utils.ErrorResponse
// @Router		/api/v1/users/me [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		utils.Unauthorized(c, "")
		return
	}
	if err := h.userRepo.Delete(c.Request.Context(), user.ID.String()); err != nil {
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted"})
}
