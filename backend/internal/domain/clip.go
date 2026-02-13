package domain

import (
	"time"

	"github.com/google/uuid"
)

type Clip struct {
	ID               uuid.UUID  `json:"id"`
	VideoID          uuid.UUID  `json:"video_id"`
	UserID           uuid.UUID  `json:"user_id"`
	Name             string     `json:"name"`
	StartTime        float64    `json:"start_time"`
	EndTime          float64    `json:"end_time"`
	DurationSeconds  *float64   `json:"duration_seconds,omitempty"`
	AspectRatio      string     `json:"aspect_ratio"`
	ViralityScore    *float64   `json:"virality_score,omitempty"`
	Status           string     `json:"status"`
	StoragePath      *string    `json:"storage_path,omitempty"`
	ThumbnailURL     *string    `json:"thumbnail_url,omitempty"`
	IsAISuggested    bool       `json:"is_ai_suggested"`
	SuggestionReason *string    `json:"suggestion_reason,omitempty"`
	ViewCount        int        `json:"view_count"`
	DownloadCount    int        `json:"download_count"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"-"`
	Style            *ClipStyle `json:"style,omitempty"`
}

type ClipStyle struct {
	ID                    uuid.UUID  `json:"id"`
	ClipID                uuid.UUID  `json:"clip_id"`
	CaptionEnabled        bool       `json:"caption_enabled"`
	CaptionFont           string     `json:"caption_font"`
	CaptionSize           int        `json:"caption_size"`
	CaptionColor          string     `json:"caption_color"`
	CaptionBgColor        *string    `json:"caption_bg_color,omitempty"`
	CaptionPosition      string     `json:"caption_position"`
	CaptionAnimation      *string    `json:"caption_animation,omitempty"`
	CaptionMaxWords       int        `json:"caption_max_words"`
	CaptionLanguage       *string    `json:"caption_language,omitempty"`
	BrandLogoURL          *string    `json:"brand_logo_url,omitempty"`
	BrandLogoPosition     *string    `json:"brand_logo_position,omitempty"`
	BrandLogoScale        float64    `json:"brand_logo_scale"`
	BrandWatermarkOpacity float64    `json:"brand_watermark_opacity"`
	OverlayTemplate       *string    `json:"overlay_template,omitempty"`
	TransitionEffect      *string    `json:"transition_effect,omitempty"`
	BackgroundMusicURL    *string    `json:"background_music_url,omitempty"`
	BackgroundMusicVolume float64    `json:"background_music_volume"`
	OriginalAudioVolume   float64    `json:"original_audio_volume"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
