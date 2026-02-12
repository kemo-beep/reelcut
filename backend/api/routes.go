package api

import (
	"os"

	"reelcut/internal/handler"
	"reelcut/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(r *gin.Engine, h *handler.Handler, m *middleware.AuthMiddleware, rdb *redis.Client) {
	public := r.Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
			auth.POST("/refresh", h.Auth.RefreshToken)
			auth.POST("/forgot-password", h.Auth.ForgotPassword)
			auth.POST("/reset-password", h.Auth.ResetPassword)
			auth.POST("/verify-email", h.Auth.VerifyEmail)
		}
		public.GET("/templates/public", h.Template.GetPublicTemplates)
		public.POST("/webhooks/stripe", h.Webhook.Stripe)
		public.GET("/download-shared-object/:encoded", h.Video.DownloadSharedObject)
	}

	protected := r.Group("/api/v1")
	protected.Use(m.Authenticate())
	// Rate limit only when explicitly enabled (e.g. production). Off by default so local dev never hits 429.
	if rdb != nil && os.Getenv("ENABLE_RATE_LIMIT") == "1" {
		protected.Use(middleware.RateLimit(rdb, middleware.TierFromUser))
	}
	{
		users := protected.Group("/users")
		{
			users.GET("/me", h.User.GetProfile)
			users.PUT("/me", h.User.UpdateProfile)
			users.PUT("/me/password", h.User.ChangePassword)
			users.GET("/me/avatar", h.User.GetAvatar)
			users.POST("/me/avatar", h.User.UploadAvatar)
			users.GET("/me/usage", h.User.GetUsageStats)
			users.GET("/me/subscription", h.Subscription.GetMySubscription)
			users.DELETE("/me", h.User.DeleteAccount)
		}

		subscriptions := protected.Group("/subscriptions")
		{
			subscriptions.POST("/create", h.Subscription.CreateSubscription)
			subscriptions.POST("/cancel", h.Subscription.CancelSubscription)
			subscriptions.POST("/update", h.Subscription.UpdateSubscription)
		}

		projects := protected.Group("/projects")
		{
			projects.GET("", h.Project.List)
			projects.POST("", h.Project.Create)
			projects.GET("/:id", h.Project.GetByID)
			projects.PUT("/:id", h.Project.Update)
			projects.DELETE("/:id", h.Project.Delete)
		}

		videos := protected.Group("/videos")
		{
			videos.POST("/upload", h.Video.Upload)
			videos.POST("/upload/resumable", h.Video.InitiateResumableUpload)
			videos.POST("/upload/url", h.Video.UploadFromURL)
			videos.PUT("/:id/upload/parts/:partNumber", h.Video.UploadPartResumable)
			videos.POST("/:id/upload/complete", h.Video.CompleteResumableUpload)
			videos.POST("/:id/confirm", h.Video.ConfirmUpload)
			videos.GET("", h.Video.List)
			// More specific GET routes first so they are not matched by /:id
			videos.GET("/:id/metadata", h.Video.GetMetadata)
			videos.GET("/:id/playback-url", h.Video.GetPlaybackURL)
			videos.POST("/:id/auto-cut", h.Video.AutoCut)
			videos.GET("/:id", h.Video.GetByID)
			videos.DELETE("/:id", h.Video.Delete)
		}

		transcriptions := protected.Group("/transcriptions")
		{
			// More specific routes first so GET .../videos/:videoId is matched before GET /:id (avoids 404 when no transcription).
			transcriptions.POST("/videos/:videoId", h.Transcription.Create)
			transcriptions.GET("/videos/:videoId", h.Transcription.GetByVideoID)
			transcriptions.GET("/:id", h.Transcription.GetByID)
			transcriptions.PUT("/:id/segments/:segmentId", h.Transcription.UpdateSegment)
		}

		analysis := protected.Group("/analysis")
		{
			analysis.POST("/videos/:videoId", h.Analysis.Analyze)
			analysis.GET("/videos/:videoId", h.Analysis.GetByVideoID)
			analysis.POST("/videos/:videoId/suggest-clips", h.Analysis.SuggestClips)
		}

		clips := protected.Group("/clips")
		{
			clips.POST("", h.Clip.Create)
			clips.GET("", h.Clip.List)
			clips.GET("/:id/playback-url", h.Clip.GetPlaybackURL)
			clips.GET("/:id", h.Clip.GetByID)
			clips.PUT("/:id", h.Clip.Update)
			clips.DELETE("/:id", h.Clip.Delete)
			clips.GET("/:id/captions/srt", h.Clip.GetCaptionsSRT)
			clips.GET("/:id/captions/vtt", h.Clip.GetCaptionsVTT)
			clips.POST("/:id/render", h.Clip.Render)
			clips.POST("/:id/cancel", h.Clip.CancelRender)
			clips.GET("/:id/status", h.Clip.GetRenderStatus)
			clips.GET("/:id/download", h.Clip.Download)
			clips.POST("/:id/duplicate", h.Clip.Duplicate)
		}

		clipStyles := protected.Group("/clips/:id/style")
		{
			clipStyles.GET("", h.Clip.GetStyle)
			clipStyles.PUT("", h.Clip.UpdateStyle)
			clipStyles.POST("/apply-template/:templateId", h.Clip.ApplyTemplate)
		}

		templates := protected.Group("/templates")
		{
			templates.GET("", h.Template.List)
			templates.POST("", h.Template.Create)
			templates.GET("/:id", h.Template.GetByID)
			templates.PUT("/:id", h.Template.Update)
			templates.DELETE("/:id", h.Template.Delete)
		}

		jobs := protected.Group("/jobs")
		{
			jobs.GET("", h.Job.List)
			jobs.GET("/:id", h.Job.GetByID)
			jobs.POST("/:id/cancel", h.Job.Cancel)
		}

		webhooks := protected.Group("/webhooks")
		{
			webhooks.POST("/processing-complete", h.Webhook.ProcessingComplete)
		}
	}

	// Thumbnail: allow Bearer or ?token= so <img src="...?token=..."> works (no CORS to MinIO)
	videoThumbnail := r.Group("/api/v1/videos")
	videoThumbnail.Use(m.AuthenticateThumbnailOrBearer())
	{
		videoThumbnail.GET("/:id/thumbnail", h.Video.GetThumbnail)
	}

	r.GET("/ws", m.AuthenticateWS(), h.WebSocket.Handle)
}
