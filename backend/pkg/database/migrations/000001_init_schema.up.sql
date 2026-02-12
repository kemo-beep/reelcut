-- Users and Authentication
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    avatar_url TEXT,
    subscription_tier VARCHAR(50) DEFAULT 'free',
    credits_remaining INT DEFAULT 0,
    email_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_subscription_tier ON users(subscription_tier);

CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_token ON user_sessions(token);

-- Projects
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_projects_user_id ON projects(user_id);

-- Videos
CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    original_filename VARCHAR(255) NOT NULL,
    storage_path TEXT NOT NULL,
    thumbnail_url TEXT,
    duration_seconds DECIMAL(10, 2),
    width INT,
    height INT,
    fps DECIMAL(5, 2),
    file_size_bytes BIGINT,
    codec VARCHAR(50),
    bitrate INT,
    status VARCHAR(50) DEFAULT 'uploading',
    error_message TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_videos_user_id ON videos(user_id);
CREATE INDEX idx_videos_project_id ON videos(project_id);
CREATE INDEX idx_videos_status ON videos(status);
CREATE INDEX idx_videos_created_at ON videos(created_at DESC);

-- Transcriptions
CREATE TABLE transcriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    language VARCHAR(10) DEFAULT 'en',
    status VARCHAR(50) DEFAULT 'pending',
    error_message TEXT,
    word_count INT,
    duration_seconds DECIMAL(10, 2),
    confidence_avg DECIMAL(5, 4),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_transcriptions_video_id ON transcriptions(video_id);

CREATE TABLE transcript_segments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transcription_id UUID REFERENCES transcriptions(id) ON DELETE CASCADE,
    start_time DECIMAL(10, 3) NOT NULL,
    end_time DECIMAL(10, 3) NOT NULL,
    text TEXT NOT NULL,
    confidence DECIMAL(5, 4),
    speaker_id INT,
    sequence_order INT NOT NULL,
    CONSTRAINT check_time_order CHECK (end_time > start_time)
);

CREATE INDEX idx_transcript_segments_transcription_id ON transcript_segments(transcription_id);
CREATE INDEX idx_segments_time ON transcript_segments(start_time, end_time);

CREATE TABLE transcript_words (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    segment_id UUID REFERENCES transcript_segments(id) ON DELETE CASCADE,
    word VARCHAR(255) NOT NULL,
    start_time DECIMAL(10, 3) NOT NULL,
    end_time DECIMAL(10, 3) NOT NULL,
    confidence DECIMAL(5, 4),
    sequence_order INT NOT NULL
);

CREATE INDEX idx_transcript_words_segment_id ON transcript_words(segment_id);
CREATE INDEX idx_words_time ON transcript_words(start_time);

-- AI Analysis
CREATE TABLE video_analysis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    scenes_detected JSONB,
    faces_detected JSONB,
    topics JSONB,
    sentiment_analysis JSONB,
    engagement_scores JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_video_analysis_video_id ON video_analysis(video_id);

-- Clips
CREATE TABLE clips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    start_time DECIMAL(10, 3) NOT NULL,
    end_time DECIMAL(10, 3) NOT NULL,
    duration_seconds DECIMAL(10, 2),
    aspect_ratio VARCHAR(10) DEFAULT '9:16',
    virality_score DECIMAL(5, 2),
    status VARCHAR(50) DEFAULT 'draft',
    storage_path TEXT,
    thumbnail_url TEXT,
    is_ai_suggested BOOLEAN DEFAULT false,
    suggestion_reason TEXT,
    view_count INT DEFAULT 0,
    download_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT check_clip_time CHECK (end_time > start_time),
    CONSTRAINT check_virality_score CHECK (virality_score IS NULL OR (virality_score >= 0 AND virality_score <= 100))
);

CREATE INDEX idx_clips_video_id ON clips(video_id);
CREATE INDEX idx_clips_user_id ON clips(user_id);
CREATE INDEX idx_clips_status ON clips(status);
CREATE INDEX idx_clips_is_ai_suggested ON clips(is_ai_suggested);
CREATE INDEX idx_clips_virality_score ON clips(virality_score DESC);

-- Clip Styles
CREATE TABLE clip_styles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clip_id UUID UNIQUE REFERENCES clips(id) ON DELETE CASCADE,
    caption_enabled BOOLEAN DEFAULT true,
    caption_font VARCHAR(100) DEFAULT 'Inter',
    caption_size INT DEFAULT 48,
    caption_color VARCHAR(7) DEFAULT '#FFFFFF',
    caption_bg_color VARCHAR(9),
    caption_position VARCHAR(50) DEFAULT 'bottom',
    caption_animation VARCHAR(50),
    caption_max_words INT DEFAULT 3,
    brand_logo_url TEXT,
    brand_logo_position VARCHAR(50),
    brand_logo_scale DECIMAL(3, 2) DEFAULT 1.0,
    brand_watermark_opacity DECIMAL(3, 2) DEFAULT 0.8,
    overlay_template VARCHAR(100),
    transition_effect VARCHAR(50),
    background_music_url TEXT,
    background_music_volume DECIMAL(3, 2) DEFAULT 0.3,
    original_audio_volume DECIMAL(3, 2) DEFAULT 1.0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_clip_styles_clip_id ON clip_styles(clip_id);

-- Processing Jobs
CREATE TABLE processing_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    job_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    priority INT DEFAULT 5,
    status VARCHAR(50) DEFAULT 'pending',
    progress INT DEFAULT 0,
    error_message TEXT,
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    metadata JSONB,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_processing_jobs_status ON processing_jobs(status);
CREATE INDEX idx_processing_jobs_entity ON processing_jobs(entity_type, entity_id);
CREATE INDEX idx_jobs_status_priority ON processing_jobs(status, priority DESC, created_at);
CREATE INDEX idx_jobs_user_id ON processing_jobs(user_id);
CREATE INDEX idx_jobs_created_at ON processing_jobs(created_at DESC);

-- Templates
CREATE TABLE templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    is_public BOOLEAN DEFAULT false,
    preview_url TEXT,
    style_config JSONB NOT NULL,
    usage_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_templates_user_id ON templates(user_id);
CREATE INDEX idx_templates_is_public ON templates(is_public);

-- Usage Tracking
CREATE TABLE usage_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL,
    credits_used INT DEFAULT 0,
    video_duration_seconds DECIMAL(10, 2),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_usage_logs_user_id ON usage_logs(user_id);

-- Subscriptions
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    tier VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    stripe_subscription_id VARCHAR(255),
    current_period_start TIMESTAMP,
    current_period_end TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
