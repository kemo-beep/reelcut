package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"reelcut/api"
	"reelcut/internal/config"
	"reelcut/internal/email"
	"reelcut/internal/handler"
	"reelcut/internal/middleware"
	"reelcut/internal/ai"
	"reelcut/internal/queue"
	"reelcut/internal/repository"
	"reelcut/internal/service"
	"reelcut/internal/worker"
	"reelcut/pkg/database"
	"reelcut/pkg/logger"
	"reelcut/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "reelcut/docs"
)

// @title			Reelcut API
// @version		1.0
// @description	API for Reelcut — clip creation from long-form video (auth, projects, videos, transcriptions, analysis, clips, templates, jobs).
// @termsOfService	https://example.com/terms

// @contact.name	API Support
// @contact.url	https://example.com/support

// @license.name	MIT
// @license.url	https://opensource.org/licenses/MIT

// @host		localhost:8080
// @BasePath	/

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger := logger.New(os.Getenv("LOG_LEVEL"))
	_ = logger

	if err := database.RunMigrations(cfg.Database.URL); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Print("database migrations applied")

	ctx := context.Background()
	pool, err := database.NewPostgresPool(ctx, cfg.Database.URL, database.PoolConfig{
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
	})
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	rdb, err := redis.NewClient(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer rdb.Close()

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	sessionRepo := repository.NewUserSessionRepository(pool)
	projectRepo := repository.NewProjectRepository(pool)
	videoRepo := repository.NewVideoRepository(pool)
	transcriptionRepo := repository.NewTranscriptionRepository(pool)
	segmentRepo := repository.NewTranscriptSegmentRepository(pool)
	wordRepo := repository.NewTranscriptWordRepository(pool)
	videoAnalysisRepo := repository.NewVideoAnalysisRepository(pool)
	clipRepo := repository.NewClipRepository(pool)
	clipStyleRepo := repository.NewClipStyleRepository(pool)
	templateRepo := repository.NewTemplateRepository(pool)
	jobRepo := repository.NewProcessingJobRepository(pool)
	usageLogRepo := repository.NewUsageLogRepository(pool)
	subscriptionRepo := repository.NewSubscriptionRepository(pool)

	// Storage
	storageSvc, err := service.NewStorageService(service.S3Config{
		Endpoint:       cfg.S3.Endpoint,
		PublicEndpoint: cfg.S3.PublicEndpoint,
		Region:         cfg.S3.Region,
		Bucket:         cfg.S3.Bucket,
		AccessKeyID:    cfg.S3.AccessKeyID,
		SecretAccessKey: cfg.S3.SecretAccessKey,
		UsePathStyle:   cfg.S3.UsePathStyle,
	})
	if err != nil {
		log.Fatalf("storage: %v", err)
	}

	// Queue
	queueClient, err := queue.NewQueueClient(cfg.Asynq.RedisURL)
	if err != nil {
		log.Fatalf("queue: %v", err)
	}
	defer queueClient.Close()

	// Email sender (SMTP or no-op)
	var emailSender email.Sender = &email.NoOpSender{}
	if cfg.Email.SMTPHost != "" && cfg.Email.SMTPPort > 0 {
		emailSender = email.NewSMTPSender(email.SMTPConfig{
			From:     cfg.Email.From,
			Host:     cfg.Email.SMTPHost,
			Port:     cfg.Email.SMTPPort,
			Username: cfg.Email.SMTPUser,
			Password: cfg.Email.SMTPPassword,
			UseTLS:   cfg.Email.SMTPUseTLS,
		})
	}
	authSvc := service.NewAuthServiceWithEmail(service.AuthServiceOpts{
		UserRepo:           userRepo,
		SessionRepo:        sessionRepo,
		EmailSender:        emailSender,
		EmailFrom:          cfg.Email.From,
		FrontendBaseURL:    cfg.Email.FrontendBaseURL,
		TokenSecret:        cfg.Email.TokenSecret,
		TokenExpiryReset:   cfg.Email.TokenExpiryReset,
		TokenExpiryVerify:  cfg.Email.TokenExpiryVerify,
		JWTSecret:          cfg.JWT.Secret,
		JWTRefresh:         cfg.JWT.RefreshSecret,
		AccessExpiry:       cfg.JWT.AccessExpiry,
		RefreshExpiry:      cfg.JWT.RefreshExpiry,
	})
	userSvc := service.NewUserService(userRepo, storageSvc)
	videoSvc := service.NewVideoService(videoRepo, projectRepo, jobRepo, storageSvc, queueClient, userRepo, usageLogRepo)
	transcriptionSvc := service.NewTranscriptionService(transcriptionRepo, segmentRepo, wordRepo, videoRepo, queueClient)
	analysisSvc := service.NewAnalysisService(videoAnalysisRepo, transcriptionRepo, segmentRepo, videoRepo, queueClient)
	renderingSvc := service.NewRenderingService(clipRepo, clipStyleRepo, videoRepo, transcriptionSvc, storageSvc)
	clipSvc := service.NewClipService(clipRepo, clipStyleRepo, videoRepo, transcriptionSvc, jobRepo, queueClient, templateRepo, userRepo, usageLogRepo)
	templateSvc := service.NewTemplateService(templateRepo)
	subscriptionSvc := service.NewSubscriptionService(subscriptionRepo, userRepo, cfg.Stripe.SecretKey, cfg.Stripe.PriceIDPro)
	var transcriber ai.Transcriber
	if cfg.Whisper.WebSocketURL != "" {
		transcriber = ai.NewWhisperLiveClient(cfg.Whisper.WebSocketURL)
	} else {
		transcriber = ai.NewWhisperClient(cfg.Whisper.APIKey)
	}

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret, userRepo, sessionRepo)

	// WebSocket hub and job notifier
	wsHub := handler.NewHub()
	jobNotifier := handler.NewJobNotifier(wsHub)

	// Asynq worker (video metadata + thumbnail)
	asynqOpt, _ := asynq.ParseRedisURI(cfg.Asynq.RedisURL)
	asynqSrv := asynq.NewServer(asynqOpt, asynq.Config{Concurrency: cfg.Asynq.Concurrency})
	mux := asynq.NewServeMux()
	videoWorker := worker.NewVideoWorker(videoRepo, jobRepo, storageSvc, jobNotifier)
	videoWorker.Register(mux)
	transcriptionWorker := worker.NewTranscriptionWorker(transcriptionRepo, segmentRepo, wordRepo, videoRepo, storageSvc, transcriber)
	transcriptionWorker.Register(mux)
	analysisWorker := worker.NewAnalysisWorker(videoAnalysisRepo, videoRepo, transcriptionRepo, segmentRepo, storageSvc)
	analysisWorker.Register(mux)
	renderingWorker := worker.NewRenderingWorker(renderingSvc, clipRepo, jobRepo, jobNotifier)
	renderingWorker.Register(mux)
	go func() {
		if err := asynqSrv.Run(mux); err != nil {
			log.Printf("asynq worker: %v", err)
		}
	}()

	// Handlers
	handlers := &handler.Handler{
		Auth:         handler.NewAuthHandler(authSvc),
		User:         handler.NewUserHandler(userRepo, usageLogRepo, authSvc, userSvc),
		Project:      handler.NewProjectHandler(projectRepo),
		Video:        handler.NewVideoHandler(videoSvc, cfg.JWT.Secret, cfg.S3.Endpoint),
		Transcription: handler.NewTranscriptionHandler(transcriptionSvc),
		Analysis:     handler.NewAnalysisHandler(analysisSvc),
		Clip:         handler.NewClipHandler(clipSvc),
		Template:     handler.NewTemplateHandler(templateSvc),
		Job:          handler.NewJobHandler(jobRepo),
		Subscription: handler.NewSubscriptionHandler(subscriptionSvc),
		Webhook:      handler.NewWebhookHandler(cfg.Stripe.WebhookSecret, subscriptionRepo, userRepo),
		WebSocket:    handler.NewWebSocketHandler(wsHub),
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS(getCORSOrigins()))
	r.Use(middleware.ErrorHandler())

	api.SetupRoutes(r, handlers, authMiddleware, rdb)

	r.GET("/swagger", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/swagger/index.html") })
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check (override the one in routes to add DB/Redis check)
	r.GET("/health", func(c *gin.Context) {
		if err := pool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		if rdb.Ping(ctx).Err() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": "redis"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	portStr := fmt.Sprintf("%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         ":" + portStr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func getCORSOrigins() []string {
	s := os.Getenv("CORS_ORIGINS")
	if s == "" {
		return []string{
			"http://localhost:3000", "http://localhost:3001",
			"http://localhost:3002", "http://localhost:3003",
			"http://localhost:5173",
		}
	}
	var out []string
	for _, part := range strings.Split(s, ",") {
		if t := strings.TrimSpace(part); t != "" {
			out = append(out, t)
		}
	}
	return out
}
