package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Logger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set(KeyRequestID, reqID)
		c.Header("X-Request-ID", reqID)
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		err := c.Errors.ByType(gin.ErrorTypePrivate).Last()
		if log != nil {
			attrs := []any{
				slog.Int("status", status),
				slog.Duration("latency", latency),
				slog.String("method", method),
				slog.String("path", path),
				slog.String("ip", clientIP),
				slog.String("request_id", reqID),
			}
			if err != nil {
				attrs = append(attrs, slog.String("error", err.Error()))
			}
			log.Info("request", attrs...)
		}
	}
}
