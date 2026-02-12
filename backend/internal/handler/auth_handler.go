package handler

import (
	"errors"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"reelcut/internal/domain"
	"reelcut/internal/service"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Register godoc
// @Summary		Register a new user
// @Tags			auth
// @Accept		json
// @Produce		json
// @Param		body	body		service.RegisterInput	true	"Registration payload"
// @Success	201		{object}	service.TokenResponse
// @Failure	400		{object}	utils.ErrorResponse
// @Router		/api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[auth] register panic: %v\n%s", r, debug.Stack())
			if !c.Writer.Written() {
				utils.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Registration failed", nil)
			}
		}
	}()

	var in service.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.ValidationError(c, []utils.ErrorDetail{{Field: "body", Message: err.Error()}})
		return
	}
	out, err := h.auth.Register(c.Request.Context(), in)
	if err != nil {
		switch {
		case err == service.ErrEmailExists:
			utils.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Email already exists", []utils.ErrorDetail{{Field: "email", Message: "Email already exists"}})
		case errors.Is(err, domain.ErrValidation):
			utils.ValidationError(c, []utils.ErrorDetail{{Field: "email", Message: "Invalid email format"}})
		case errors.As(err, new(*domain.ValidationError)):
			var ve *domain.ValidationError
			errors.As(err, &ve)
			utils.ValidationError(c, []utils.ErrorDetail{{Field: ve.Field, Message: ve.Message}})
		default:
			log.Printf("[auth] register error: %v", err)
			msg := "Internal server error"
			if os.Getenv("APP_ENV") == "development" || os.Getenv("DEBUG") != "" {
				msg = err.Error()
			}
			if !c.Writer.Written() {
				utils.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", msg, nil)
			}
		}
		return
	}
	sanitizeUser(out.User)
	c.JSON(http.StatusCreated, out)
}

// Login godoc
// @Summary		Login with email and password
// @Tags			auth
// @Accept		json
// @Produce		json
// @Param		body	body		service.LoginInput	true	"Login credentials"
// @Success	200	{object}	service.TokenResponse
// @Failure	401	{object}	utils.ErrorResponse
// @Router		/api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var in service.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.ValidationError(c, []utils.ErrorDetail{{Field: "body", Message: err.Error()}})
		return
	}
	out, err := h.auth.Login(c.Request.Context(), in)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			utils.Unauthorized(c, "Invalid email or password")
			return
		}
		log.Printf("[auth] login error: %v", err)
		utils.Internal(c, "")
		return
	}
	sanitizeUser(out.User)
	c.JSON(http.StatusOK, out)
}

// RefreshToken godoc
// @Summary		Refresh access token
// @Tags			auth
// @Accept		json
// @Produce		json
// @Param		body	body		object	true	"{\"refresh_token\":\"string\"}"
// @Success	200	{object}	service.TokenResponse
// @Failure	401	{object}	utils.ErrorResponse
// @Router		/api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, []utils.ErrorDetail{{Field: "refresh_token", Message: "required"}})
		return
	}
	out, err := h.auth.RefreshToken(c.Request.Context(), body.RefreshToken)
	if err != nil {
		utils.Unauthorized(c, "Invalid or expired refresh token")
		return
	}
	sanitizeUser(out.User)
	c.JSON(http.StatusOK, out)
}

// ForgotPassword godoc
// @Summary		Request password reset email
// @Tags			auth
// @Accept		json
// @Produce		json
// @Param		body	body		object	true	"{\"email\":\"string\"}"
// @Success	200	{object}	object
// @Router		/api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, nil)
		return
	}
	_ = h.auth.ForgotPassword(c.Request.Context(), body.Email)
	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link has been sent"})
}

// ResetPassword godoc
// @Summary		Reset password with token
// @Tags			auth
// @Accept		json
// @Produce		json
// @Param		body	body		service.ResetPasswordInput	true	"Token and new password"
// @Success	200	{object}	object
// @Failure	400	{object}	utils.ErrorResponse
// @Router		/api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var in service.ResetPasswordInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.ValidationError(c, nil)
		return
	}
	if err := h.auth.ResetPassword(c.Request.Context(), in); err != nil {
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			utils.ValidationError(c, []utils.ErrorDetail{{Field: ve.Field, Message: ve.Message}})
			return
		}
		if err == service.ErrInvalidToken {
			utils.Unauthorized(c, "Invalid or expired reset token")
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset"})
}

// VerifyEmail godoc
// @Summary		Verify email with token
// @Tags			auth
// @Accept		json
// @Produce		json
// @Param		body	body		object	true	"{\"token\":\"string\"}"
// @Success	200	{object}	object
// @Failure	401	{object}	utils.ErrorResponse
// @Router		/api/v1/auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var body struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, nil)
		return
	}
	if err := h.auth.VerifyEmail(c.Request.Context(), body.Token); err != nil {
		if err == service.ErrInvalidToken {
			utils.Unauthorized(c, "Invalid or expired verification token")
			return
		}
		utils.Internal(c, "")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email verified"})
}

func sanitizeUser(u *domain.User) {
	if u != nil {
		u.PasswordHash = ""
	}
}
