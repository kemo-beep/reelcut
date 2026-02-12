package middleware

import (
	"net/http"
	"strings"

	"reelcut/internal/domain"
	"reelcut/internal/repository"
	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtSecret    string
	userRepo     repository.UserRepository
	sessionRepo  repository.UserSessionRepository
}

func NewAuthMiddleware(jwtSecret string, userRepo repository.UserRepository, sessionRepo repository.UserSessionRepository) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret, userRepo: userRepo, sessionRepo: sessionRepo}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token == "" {
			utils.Unauthorized(c, "Missing or invalid authorization header")
			c.Abort()
			return
		}
		claims, err := utils.ParseAccessToken(token, m.jwtSecret)
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}
		user, err := m.userRepo.GetByID(c.Request.Context(), claims.UserID)
		if err != nil || user == nil {
			utils.Unauthorized(c, "User not found")
			c.Abort()
			return
		}
		c.Set(KeyUserID, claims.UserID)
		c.Set(KeyUser, user)
		c.Next()
	}
}

func (m *AuthMiddleware) AuthenticateWS() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			token = extractBearerToken(c)
		}
		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, err := utils.ParseAccessToken(token, m.jwtSecret)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(KeyUserID, claims.UserID)
		c.Next()
	}
}

func extractBearerToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return ""
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return ""
	}
	return strings.TrimSpace(auth[len(prefix):])
}

func GetUserID(c *gin.Context) string {
	v, _ := c.Get(KeyUserID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func GetUser(c *gin.Context) *domain.User {
	v, _ := c.Get(KeyUser)
	if u, ok := v.(*domain.User); ok {
		return u
	}
	return nil
}

// AuthenticateThumbnailOrBearer validates either a thumbnail token (query "token") or Bearer and sets user ID.
// Used for GET /videos/:id/thumbnail so <img src="...?token=..."> works without sending Authorization.
func (m *AuthMiddleware) AuthenticateThumbnailOrBearer() gin.HandlerFunc {
	return func(c *gin.Context) {
		videoID := c.Param("id")
		if videoID == "" {
			utils.Unauthorized(c, "Missing video id")
			c.Abort()
			return
		}
		if tokenStr := c.Query("token"); tokenStr != "" {
			claims, err := utils.ParseThumbnailToken(tokenStr, m.jwtSecret)
			if err != nil || claims.VideoID != videoID {
				utils.Unauthorized(c, "Invalid or expired thumbnail token")
				c.Abort()
				return
			}
			c.Set(KeyUserID, claims.UserID)
			c.Next()
			return
		}
		token := extractBearerToken(c)
		if token == "" {
			utils.Unauthorized(c, "Missing or invalid authorization header or token")
			c.Abort()
			return
		}
		claims, err := utils.ParseAccessToken(token, m.jwtSecret)
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}
		user, err := m.userRepo.GetByID(c.Request.Context(), claims.UserID)
		if err != nil || user == nil {
			utils.Unauthorized(c, "User not found")
			c.Abort()
			return
		}
		c.Set(KeyUserID, claims.UserID)
		c.Set(KeyUser, user)
		c.Next()
	}
}
