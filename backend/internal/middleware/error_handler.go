package middleware

import (
	"errors"
	"log"
	"net/http"

	"reelcut/internal/domain"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last()
		if err == nil {
			return
		}
		if c.Writer.Written() {
			return
		}
		unwrap := err.Err
		for u, ok := unwrap.(interface{ Unwrap() error }); ok; u, ok = unwrap.(interface{ Unwrap() error }) {
			unwrap = u.Unwrap()
		}
		switch {
		case errors.Is(unwrap, domain.ErrNotFound):
			utils.NotFound(c, unwrap.Error())
		case errors.Is(unwrap, domain.ErrUnauthorized):
			utils.Unauthorized(c, unwrap.Error())
		case errors.Is(unwrap, domain.ErrForbidden):
			utils.Forbidden(c, unwrap.Error())
		case errors.Is(unwrap, domain.ErrValidation):
			utils.ValidationError(c, nil)
		case errors.Is(unwrap, domain.ErrConflict):
			utils.Error(c, http.StatusConflict, "CONFLICT", unwrap.Error(), nil)
		case errors.Is(unwrap, domain.ErrRateLimitExceeded):
			utils.Error(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", unwrap.Error(), nil)
		default:
			log.Printf("unhandled error: %v", unwrap)
			utils.Internal(c, "")
		}
	}
}
