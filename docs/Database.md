-- Users and Authentication
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    avatar_url TEXT,
    subscription_tier VARCHAR(50) DEFAULT 'free', -- free, pro, enterprise
    credits_remaining INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Projects and Videos
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    original_filename VARCHAR(255) NOT NULL,
    storage_path TEXT NOT NULL, -- S3 path
    thumbnail_url TEXT,
    duration_seconds DECIMAL(10, 2),
    width INT,
    height INT,
    file_size_bytes BIGINT,
    status VARCHAR(50) DEFAULT 'uploading', -- uploading, processing, ready, failed
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Transcription
CREATE TABLE transcriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    language VARCHAR(10) DEFAULT 'en',
    status VARCHAR(50) DEFAULT 'pending', -- pending, processing, completed, failed
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE transcript_segments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transcription_id UUID REFERENCES transcriptions(id) ON DELETE CASCADE,
    start_time DECIMAL(10, 3) NOT NULL,
    end_time DECIMAL(10, 3) NOT NULL,
    text TEXT NOT NULL,
    confidence DECIMAL(5, 4),
    speaker_id INT, -- for diarization
    sequence_order INT NOT NULL
);

CREATE TABLE transcript_words (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    segment_id UUID REFERENCES transcript_segments(id) ON DELETE CASCADE,
    word VARCHAR(255) NOT NULL,
    start_time DECIMAL(10, 3) NOT NULL,
    end_time DECIMAL(10, 3) NOT NULL,
    confidence DECIMAL(5, 4),
    sequence_order INT NOT NULL
);

-- AI Analysis
CREATE TABLE video_analysis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    scenes_detected JSONB, -- [{start: 0, end: 10, type: 'intro'}, ...]
    faces_detected JSONB, -- [{timestamp: 5, bbox: {...}, confidence: 0.9}, ...]
    topics JSONB, -- [{topic: 'marketing', relevance: 0.8, timestamps: [...]}]
    sentiment_analysis JSONB, -- [{start: 0, end: 10, sentiment: 'positive', score: 0.8}]
    engagement_scores JSONB, -- [{start: 0, end: 10, score: 0.75}]
    created_at TIMESTAMP DEFAULT NOW()
);

-- Clips
CREATE TABLE clips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    start_time DECIMAL(10, 3) NOT NULL,
    end_time DECIMAL(10, 3) NOT NULL,
    duration_seconds DECIMAL(10, 2),
    aspect_ratio VARCHAR(10) DEFAULT '9:16', -- 9:16, 1:1, 16:9
    virality_score DECIMAL(5, 2), -- 0-100
    status VARCHAR(50) DEFAULT 'draft', -- draft, rendering, ready, failed
    storage_path TEXT, -- rendered clip path
    thumbnail_url TEXT,
    is_ai_suggested BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Clip Customization
CREATE TABLE clip_styles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clip_id UUID REFERENCES clips(id) ON DELETE CASCADE,
    caption_enabled BOOLEAN DEFAULT true,
    caption_font VARCHAR(100) DEFAULT 'Inter',
    caption_size INT DEFAULT 48,
    caption_color VARCHAR(7) DEFAULT '#FFFFFF',
    caption_bg_color VARCHAR(9), -- with alpha
    caption_position VARCHAR(50) DEFAULT 'bottom', -- top, center, bottom
    caption_animation VARCHAR(50), -- fade, slide, bounce
    brand_logo_url TEXT,
    brand_logo_position VARCHAR(50), -- top-left, top-right, bottom-left, bottom-right
    overlay_template VARCHAR(100),
    background_music_url TEXT,
    background_music_volume DECIMAL(3, 2) DEFAULT 0.3,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Processing Jobs
CREATE TABLE processing_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    job_type VARCHAR(50) NOT NULL, -- transcription, analysis, rendering
    entity_type VARCHAR(50) NOT NULL, -- video, clip
    entity_id UUID NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', -- pending, processing, completed, failed
    progress INT DEFAULT 0, -- 0-100
    error_message TEXT,
    metadata JSONB, -- job-specific data
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Templates
CREATE TABLE templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- NULL for global templates
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100), -- educational, podcast, social, business
    is_public BOOLEAN DEFAULT false,
    preview_url TEXT,
    style_config JSONB NOT NULL, -- complete style configuration
    usage_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Usage Tracking
CREATE TABLE usage_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL, -- video_upload, clip_render, transcription
    credits_used INT DEFAULT 0,
    video_duration_seconds DECIMAL(10, 2),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Subscriptions
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    tier VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'active', -- active, canceled, expired
    stripe_subscription_id VARCHAR(255),
    current_period_start TIMESTAMP,
    current_period_end TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_videos_user_id ON videos(user_id);
CREATE INDEX idx_videos_status ON videos(status);
CREATE INDEX idx_clips_video_id ON clips(video_id);
CREATE INDEX idx_clips_user_id ON clips(user_id);
CREATE INDEX idx_transcript_segments_transcription_id ON transcript_segments(transcription_id);
CREATE INDEX idx_transcript_words_segment_id ON transcript_words(segment_id);
CREATE INDEX idx_processing_jobs_status ON processing_jobs(status);
CREATE INDEX idx_processing_jobs_entity ON processing_jobs(entity_type, entity_id);
```

## Backend Structure (Go Gin)
```
backend/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   ├── config.go              # Configuration management
│   │   └── env.go                 # Environment variables
│   ├── domain/
│   │   ├── user.go                # User domain models
│   │   ├── video.go               # Video domain models
│   │   ├── clip.go                # Clip domain models
│   │   ├── transcription.go       # Transcription models
│   │   ├── template.go            # Template models
│   │   └── errors.go              # Custom error types
│   ├── repository/
│   │   ├── interfaces.go          # Repository interfaces
│   │   ├── user_repository.go
│   │   ├── video_repository.go
│   │   ├── clip_repository.go
│   │   ├── transcription_repository.go
│   │   └── template_repository.go
│   ├── service/
│   │   ├── auth_service.go        # Authentication logic
│   │   ├── video_service.go       # Video upload, management
│   │   ├── transcription_service.go # Transcription orchestration
│   │   ├── analysis_service.go    # AI analysis
│   │   ├── clip_service.go        # Clip creation, management
│   │   ├── rendering_service.go   # Video rendering
│   │   ├── storage_service.go     # S3/file storage
│   │   └── queue_service.go       # Job queue management
│   ├── handler/
│   │   ├── auth_handler.go        # Auth endpoints
│   │   ├── user_handler.go        # User endpoints
│   │   ├── project_handler.go     # Project CRUD
│   │   ├── video_handler.go       # Video upload/management
│   │   ├── clip_handler.go        # Clip creation/management
│   │   ├── transcription_handler.go
│   │   ├── template_handler.go
│   │   └── webhook_handler.go     # External webhooks
│   ├── middleware/
│   │   ├── auth.go                # JWT authentication
│   │   ├── cors.go                # CORS configuration
│   │   ├── logger.go              # Request logging
│   │   ├── rate_limit.go          # Rate limiting
│   │   └── error_handler.go       # Global error handling
│   ├── worker/
│   │   ├── worker.go              # Worker pool manager
│   │   ├── transcription_worker.go # Transcription jobs
│   │   ├── analysis_worker.go     # AI analysis jobs
│   │   └── rendering_worker.go    # Video rendering jobs
│   ├── ai/
│   │   ├── whisper_client.go      # Whisper API integration
│   │   ├── clip_detector.go       # AI clip detection
│   │   ├── scene_detector.go      # Scene detection
│   │   ├── face_tracker.go        # Face tracking
│   │   └── sentiment_analyzer.go  # Sentiment analysis
│   ├── video/
│   │   ├── ffmpeg.go              # FFmpeg wrapper
│   │   ├── processor.go           # Video processing utilities
│   │   ├── frame_extractor.go     # Extract frames
│   │   └── metadata.go            # Extract video metadata
│   ├── queue/
│   │   ├── redis_queue.go         # Redis queue implementation
│   │   └── job.go                 # Job definitions
│   └── utils/
│       ├── jwt.go                 # JWT utilities
│       ├── password.go            # Password hashing
│       ├── validator.go           # Input validation
│       └── response.go            # Standard API responses
├── pkg/
│   ├── database/
│   │   ├── postgres.go            # PostgreSQL connection
│   │   └── migrations/
│   │       └── *.sql              # Migration files
│   ├── redis/
│   │   └── client.go              # Redis connection
│   └── logger/
│       └── logger.go              # Structured logging
├── api/
│   └── routes.go                  # Route definitions
├── .air.toml                      # Air hot reload config
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum