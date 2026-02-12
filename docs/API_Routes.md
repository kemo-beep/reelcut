// api/routes.go
package api

import (
    "github.com/gin-gonic/gin"
    "yourapp/internal/handler"
    "yourapp/internal/middleware"
)

func SetupRoutes(r *gin.Engine, h *handler.Handler, m *middleware.Middleware) {
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // Public routes
    public := r.Group("/api/v1")
    {
        // Auth
        auth := public.Group("/auth")
        {
            auth.POST("/register", h.Auth.Register)
            auth.POST("/login", h.Auth.Login)
            auth.POST("/refresh", h.Auth.RefreshToken)
            auth.POST("/forgot-password", h.Auth.ForgotPassword)
            auth.POST("/reset-password", h.Auth.ResetPassword)
        }

        // Public templates
        public.GET("/templates/public", h.Template.GetPublicTemplates)
    }

    // Protected routes
    protected := r.Group("/api/v1")
    protected.Use(m.Auth.Authenticate())
    {
        // User
        users := protected.Group("/users")
        {
            users.GET("/me", h.User.GetProfile)
            users.PUT("/me", h.User.UpdateProfile)
            users.GET("/me/usage", h.User.GetUsageStats)
            users.DELETE("/me", h.User.DeleteAccount)
        }

        // Projects
        projects := protected.Group("/projects")
        {
            projects.GET("", h.Project.List)
            projects.POST("", h.Project.Create)
            projects.GET("/:id", h.Project.GetByID)
            projects.PUT("/:id", h.Project.Update)
            projects.DELETE("/:id", h.Project.Delete)
        }

        // Videos
        videos := protected.Group("/videos")
        {
            videos.POST("/upload", h.Video.Upload)
            videos.POST("/upload/url", h.Video.UploadFromURL)
            videos.GET("", h.Video.List)
            videos.GET("/:id", h.Video.GetByID)
            videos.DELETE("/:id", h.Video.Delete)
            videos.GET("/:id/metadata", h.Video.GetMetadata)
            videos.GET("/:id/thumbnail", h.Video.GetThumbnail)
        }

        // Transcription
        transcriptions := protected.Group("/transcriptions")
        {
            transcriptions.POST("/videos/:videoId", h.Transcription.Create)
            transcriptions.GET("/:id", h.Transcription.GetByID)
            transcriptions.GET("/videos/:videoId", h.Transcription.GetByVideoID)
            transcriptions.PUT("/:id/segments/:segmentId", h.Transcription.UpdateSegment)
        }

        // AI Analysis
        analysis := protected.Group("/analysis")
        {
            analysis.POST("/videos/:videoId", h.Analysis.Analyze)
            analysis.GET("/videos/:videoId", h.Analysis.GetByVideoID)
            analysis.POST("/videos/:videoId/suggest-clips", h.Analysis.SuggestClips)
        }

        // Clips
        clips := protected.Group("/clips")
        {
            clips.POST("", h.Clip.Create)
            clips.GET("", h.Clip.List)
            clips.GET("/:id", h.Clip.GetByID)
            clips.PUT("/:id", h.Clip.Update)
            clips.DELETE("/:id", h.Clip.Delete)
            clips.POST("/:id/render", h.Clip.Render)
            clips.GET("/:id/status", h.Clip.GetRenderStatus)
            clips.GET("/:id/download", h.Clip.Download)
            clips.POST("/:id/duplicate", h.Clip.Duplicate)
        }

        // Clip Styles
        clipStyles := protected.Group("/clips/:clipId/style")
        {
            clipStyles.GET("", h.Clip.GetStyle)
            clipStyles.PUT("", h.Clip.UpdateStyle)
            clipStyles.POST("/apply-template/:templateId", h.Clip.ApplyTemplate)
        }

        // Templates
        templates := protected.Group("/templates")
        {
            templates.GET("", h.Template.List)
            templates.POST("", h.Template.Create)
            templates.GET("/:id", h.Template.GetByID)
            templates.PUT("/:id", h.Template.Update)
            templates.DELETE("/:id", h.Template.Delete)
        }

        // Jobs
        jobs := protected.Group("/jobs")
        {
            jobs.GET("", h.Job.List)
            jobs.GET("/:id", h.Job.GetByID)
            jobs.POST("/:id/cancel", h.Job.Cancel)
        }

        // Webhooks (for processing updates)
        webhooks := protected.Group("/webhooks")
        {
            webhooks.POST("/processing-complete", h.Webhook.ProcessingComplete)
        }
    }

    // WebSocket for real-time updates
    r.GET("/ws", m.Auth.AuthenticateWS(), h.WebSocket.Handle)
}
```

## Frontend Structure (TanStack Start + TypeScript)
```
frontend/
├── app/
│   ├── routes/
│   │   ├── __root.tsx             # Root layout
│   │   ├── index.tsx              # Landing page
│   │   ├── auth/
│   │   │   ├── login.tsx
│   │   │   ├── register.tsx
│   │   │   └── reset-password.tsx
│   │   ├── dashboard/
│   │   │   ├── index.tsx          # Dashboard home
│   │   │   ├── projects/
│   │   │   │   ├── index.tsx      # Projects list
│   │   │   │   └── $projectId.tsx # Project detail
│   │   │   ├── videos/
│   │   │   │   ├── index.tsx      # Videos list
│   │   │   │   ├── upload.tsx     # Video upload
│   │   │   │   └── $videoId/
│   │   │   │       ├── index.tsx  # Video detail
│   │   │   │       ├── clips.tsx  # Clips from video
│   │   │   │       └── transcription.tsx
│   │   │   ├── clips/
│   │   │   │   ├── index.tsx      # All clips
│   │   │   │   └── $clipId/
│   │   │   │       ├── edit.tsx   # Clip editor
│   │   │   │       └── preview.tsx
│   │   │   ├── templates/
│   │   │   │   ├── index.tsx      # Templates library
│   │   │   │   └── $templateId.tsx
│   │   │   └── settings/
│   │   │       ├── profile.tsx
│   │   │       ├── billing.tsx
│   │   │       └── usage.tsx
│   │   └── editor/
│   │       └── $clipId.tsx        # Full-screen clip editor
│   ├── components/
│   │   ├── ui/                    # Base UI components
│   │   │   ├── button.tsx
│   │   │   ├── input.tsx
│   │   │   ├── select.tsx
│   │   │   ├── modal.tsx
│   │   │   ├── dropdown.tsx
│   │   │   ├── tabs.tsx
│   │   │   ├── card.tsx
│   │   │   ├── badge.tsx
│   │   │   ├── progress.tsx
│   │   │   ├── slider.tsx
│   │   │   ├── switch.tsx
│   │   │   └── toast.tsx
│   │   ├── layout/
│   │   │   ├── Header.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   ├── Footer.tsx
│   │   │   └── DashboardLayout.tsx
│   │   ├── auth/
│   │   │   ├── LoginForm.tsx
│   │   │   ├── RegisterForm.tsx
│   │   │   └── ProtectedRoute.tsx
│   │   ├── video/
│   │   │   ├── VideoPlayer.tsx
│   │   │   ├── VideoUploader.tsx
│   │   │   ├── VideoCard.tsx
│   │   │   ├── VideoGrid.tsx
│   │   │   ├── VideoTimeline.tsx
│   │   │   └── ThumbnailGenerator.tsx
│   │   ├── clip/
│   │   │   ├── ClipCard.tsx
│   │   │   ├── ClipGrid.tsx
│   │   │   ├── ClipPreview.tsx
│   │   │   ├── ClipEditor.tsx
│   │   │   ├── ClipTimeline.tsx
│   │   │   └── ViralityScore.tsx
│   │   ├── editor/
│   │   │   ├── EditorCanvas.tsx
│   │   │   ├── EditorToolbar.tsx
│   │   │   ├── EditorSidebar.tsx
│   │   │   ├── TimelineEditor.tsx
│   │   │   ├── CaptionEditor.tsx
│   │   │   ├── StylePanel.tsx
│   │   │   ├── LayersPanel.tsx
│   │   │   └── ExportPanel.tsx
│   │   ├── transcription/
│   │   │   ├── TranscriptViewer.tsx
│   │   │   ├── TranscriptEditor.tsx
│   │   │   ├── WordTimeline.tsx
│   │   │   └── SpeakerLabels.tsx
│   │   ├── template/
│   │   │   ├── TemplateCard.tsx
│   │   │   ├── TemplateGrid.tsx
│   │   │   ├── TemplatePreview.tsx
│   │   │   └── TemplateEditor.tsx
│   │   ├── processing/
│   │   │   ├── JobStatus.tsx
│   │   │   ├── ProgressBar.tsx
│   │   │   └── ProcessingQueue.tsx
│   │   └── common/
│   │       ├── LoadingSpinner.tsx
│   │       ├── ErrorBoundary.tsx
│   │       ├── EmptyState.tsx
│   │       ├── Pagination.tsx
│   │       └── SearchBar.tsx
│   ├── lib/
│   │   ├── api/
│   │   │   ├── client.ts          # Axios/fetch client
│   │   │   ├── auth.ts            # Auth API calls
│   │   │   ├── videos.ts          # Video API calls
│   │   │   ├── clips.ts           # Clip API calls
│   │   │   ├── transcriptions.ts
│   │   │   ├── templates.ts
│   │   │   └── jobs.ts
│   │   ├── hooks/
│   │   │   ├── useAuth.ts
│   │   │   ├── useVideos.ts
│   │   │   ├── useClips.ts
│   │   │   ├── useTranscription.ts
│   │   │   ├── useWebSocket.ts
│   │   │   ├── useUpload.ts
│   │   │   ├── useDebounce.ts
│   │   │   └── useLocalStorage.ts
│   │   ├── utils/
│   │   │   ├── format.ts          # Date, time, file size formatting
│   │   │   ├── validation.ts      # Form validation
│   │   │   ├── video.ts           # Video utilities
│   │   │   ├── time.ts            # Time conversion utilities
│   │   │   └── download.ts        # File download helpers
│   │   └── constants/
│   │       ├── config.ts
│   │       ├── routes.ts
│   │       └── styles.ts
│   ├── stores/
│   │   ├── authStore.ts           # Zustand auth store
│   │   ├── videoStore.ts          # Video state
│   │   ├── clipStore.ts           # Clip state
│   │   ├── editorStore.ts         # Editor state
│   │   ├── uiStore.ts             # UI state (modals, toasts)
│   │   └── uploadStore.ts         # Upload progress
│   ├── types/
│   │   ├── auth.ts
│   │   ├── user.ts
│   │   ├── video.ts
│   │   ├── clip.ts
│   │   ├── transcription.ts
│   │   ├── template.ts
│   │   ├── job.ts
│   │   └── api.ts
│   └── styles/
│       ├── globals.css
│       └── tailwind.css
├── public/
│   ├── fonts/
│   ├── images/
│   └── icons/
├── .env.example
├── .env.local
├── tsconfig.json
├── tailwind.config.ts
├── postcss.config.js
└── package.json