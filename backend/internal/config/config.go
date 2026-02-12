package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Email    EmailConfig
	S3       S3Config
	Asynq    AsynqConfig
	Whisper  WhisperConfig
	Stripe   StripeConfig
}

// EmailConfig for transactional email (password reset, verification). Use SMTP (e.g. SendGrid SMTP relay).
type EmailConfig struct {
	From             string
	FrontendBaseURL  string
	SMTPHost         string
	SMTPPort         int
	SMTPUser         string
	SMTPPassword     string
	SMTPUseTLS       bool
	TokenSecret      string // secret for one-time JWT (reset/verify)
	TokenExpiryReset time.Duration
	TokenExpiryVerify time.Duration
}

type StripeConfig struct {
	SecretKey     string
	WebhookSecret string
	PriceIDPro    string
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type RedisConfig struct {
	URL string
}

type JWTConfig struct {
	Secret           string
	AccessExpiry     time.Duration
	RefreshExpiry    time.Duration
	RefreshSecret    string
}

type S3Config struct {
	Endpoint        string // MinIO/S3 API endpoint (e.g. http://localhost:9000 or http://minio:9000 in Docker)
	PublicEndpoint  string // If set, used for presigned URLs so the browser can reach MinIO (e.g. http://localhost:9000 when backend runs in Docker)
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	UsePathStyle    bool
}

type AsynqConfig struct {
	RedisURL      string
	QueueName     string
	Concurrency   int
}

type WhisperConfig struct {
	APIKey      string // OpenAI API key (used when WebSocketURL is empty)
	WebSocketURL string // WhisperLiveKit ASR WebSocket base URL (e.g. ws://localhost:8000). If set, used instead of OpenAI.
}

func Load() (*Config, error) {
	_ = os.Setenv("TZ", "UTC")
	if err := godotenvLoad(); err != nil {
		// .env is optional
	}

	accessExpiry := 1 * time.Hour
	if v := os.Getenv("JWT_ACCESS_EXPIRY"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			accessExpiry = d
		}
	}
	refreshExpiry := 168 * time.Hour // 7 days
	if v := os.Getenv("JWT_REFRESH_EXPIRY"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			refreshExpiry = d
		}
	}

	port := 8080
	if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	}

	maxOpenConns := 25
	if v := os.Getenv("DB_MAX_OPEN_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			maxOpenConns = n
		}
	}
	maxIdleConns := 5
	if v := os.Getenv("DB_MAX_IDLE_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			maxIdleConns = n
		}
	}
	connMaxLifetime := 5 * time.Minute
	if v := os.Getenv("DB_CONN_MAX_LIFETIME"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			connMaxLifetime = d
		}
	}
	connMaxIdleTime := 10 * time.Minute
	if v := os.Getenv("DB_CONN_MAX_IDLE_TIME"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			connMaxIdleTime = d
		}
	}

	asynqConcurrency := 5
	if v := os.Getenv("ASYNQ_CONCURRENCY"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			asynqConcurrency = n
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:         port,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/reelcut?sslmode=disable"),
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    maxIdleConns,
			ConnMaxLifetime: connMaxLifetime,
			ConnMaxIdleTime: connMaxIdleTime,
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379/0"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "change-me-in-production"),
			AccessExpiry:  accessExpiry,
			RefreshExpiry:  refreshExpiry,
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", "change-me-refresh"),
		},
		S3: S3Config{
			Endpoint:       getEnv("S3_ENDPOINT", "http://localhost:9002"),
			PublicEndpoint: getEnv("S3_PUBLIC_ENDPOINT", ""), // optional; if set, presigned URLs use this (e.g. http://localhost:9000 when backend is in Docker)
			Region:         getEnv("S3_REGION", "us-east-1"),
			Bucket:         getEnv("S3_BUCKET", "reelcut"),
			AccessKeyID:    getEnv("S3_ACCESS_KEY_ID", "minioadmin"),
			SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", "minioadmin"),
			UsePathStyle:   getEnv("S3_USE_PATH_STYLE", "true") == "true",
		},
		Asynq: AsynqConfig{
			RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379/0"),
			QueueName:   getEnv("ASYNQ_QUEUE", "default"),
			Concurrency: asynqConcurrency,
		},
		Whisper: WhisperConfig{
			APIKey:        getEnv("OPENAI_API_KEY", ""),
			WebSocketURL:  getEnv("TRANSCRIPTION_WS_URL", ""),
		},
		Stripe: StripeConfig{
			SecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
			WebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
			PriceIDPro:    getEnv("STRIPE_PRICE_ID_PRO", ""),
		},
		Email: EmailConfig{
			From:              getEnv("EMAIL_FROM", "noreply@reelcut.local"),
			FrontendBaseURL:   getEnv("FRONTEND_BASE_URL", "http://localhost:5173"),
			SMTPHost:          getEnv("SMTP_HOST", ""),
			SMTPPort:          getEnvInt("SMTP_PORT", 587),
			SMTPUser:          getEnv("SMTP_USER", ""),
			SMTPPassword:      getEnv("SMTP_PASSWORD", ""),
			SMTPUseTLS:        getEnv("SMTP_USE_TLS", "true") == "true",
			TokenSecret:       getEnv("EMAIL_TOKEN_SECRET", getEnv("JWT_SECRET", "change-me-in-production")),
			TokenExpiryReset:  getEnvDuration("EMAIL_RESET_EXPIRY", 1*time.Hour),
			TokenExpiryVerify: getEnvDuration("EMAIL_VERIFY_EXPIRY", 24*time.Hour),
		},
	}

	if cfg.Database.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWT.Secret == "" || cfg.JWT.Secret == "change-me-in-production" {
		// Allow for dev; in prod caller should validate
	}
	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultVal
}

var godotenvLoad = func() error { return nil }
