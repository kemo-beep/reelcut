package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaginationMeta struct {
	Page        int  `json:"page"`
	PerPage     int  `json:"per_page"`
	TotalPages  int  `json:"total_pages"`
	TotalCount  int  `json:"total_count"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error struct {
		Code      string        `json:"code"`
		Message   string        `json:"message"`
		Details   []ErrorDetail `json:"details,omitempty"`
		RequestID string        `json:"request_id,omitempty"`
		Timestamp string        `json:"timestamp"`
	} `json:"error"`
}

func Pagination(page, perPage, total int) PaginationMeta {
	if perPage <= 0 {
		perPage = 20
	}
	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}
	return PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		TotalCount: total,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func JSONPaginated(c *gin.Context, status int, data interface{}, page, perPage, total int) {
	c.JSON(status, gin.H{
		"data":       data,
		"pagination": Pagination(page, perPage, total),
	})
}

func Error(c *gin.Context, status int, code, message string, details []ErrorDetail) {
	reqID := c.GetString("request_id")
	if reqID == "" {
		reqID = uuid.New().String()
		c.Set("request_id", reqID)
	}
	c.JSON(status, ErrorResponse{
		Error: struct {
			Code      string        `json:"code"`
			Message   string        `json:"message"`
			Details   []ErrorDetail `json:"details,omitempty"`
			RequestID string        `json:"request_id,omitempty"`
			Timestamp string        `json:"timestamp"`
		}{
			Code:      code,
			Message:   message,
			Details:   details,
			RequestID: reqID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func ValidationError(c *gin.Context, details []ErrorDetail) {
	Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input data", details)
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Authentication required"
	}
	Error(c, http.StatusUnauthorized, "AUTHENTICATION_ERROR", message, nil)
}

func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Insufficient permissions"
	}
	Error(c, http.StatusForbidden, "AUTHORIZATION_ERROR", message, nil)
}

func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	Error(c, http.StatusNotFound, "NOT_FOUND", message, nil)
}

func Internal(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}
