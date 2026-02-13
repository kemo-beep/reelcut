package repository

import (
	"context"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type clipStyleRepository struct {
	pool *pgxpool.Pool
}

func NewClipStyleRepository(pool *pgxpool.Pool) ClipStyleRepository {
	return &clipStyleRepository{pool: pool}
}

func (r *clipStyleRepository) Create(ctx context.Context, s *domain.ClipStyle) error {
	query := `INSERT INTO clip_styles (id, clip_id, caption_enabled, caption_font, caption_size, caption_color, caption_bg_color, caption_position, caption_animation, caption_max_words, caption_language,
		brand_logo_url, brand_logo_position, brand_logo_scale, brand_watermark_opacity, overlay_template, transition_effect, background_music_url, background_music_volume, original_audio_volume)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`
	_, err := r.pool.Exec(ctx, query, s.ID, s.ClipID, s.CaptionEnabled, s.CaptionFont, s.CaptionSize, s.CaptionColor, s.CaptionBgColor, s.CaptionPosition, s.CaptionAnimation, s.CaptionMaxWords, s.CaptionLanguage,
		s.BrandLogoURL, s.BrandLogoPosition, s.BrandLogoScale, s.BrandWatermarkOpacity, s.OverlayTemplate, s.TransitionEffect, s.BackgroundMusicURL, s.BackgroundMusicVolume, s.OriginalAudioVolume)
	return err
}

func (r *clipStyleRepository) GetByClipID(ctx context.Context, clipID string) (*domain.ClipStyle, error) {
	query := `SELECT id, clip_id, caption_enabled, caption_font, caption_size, caption_color, caption_bg_color, caption_position, caption_animation, caption_max_words, caption_language,
		brand_logo_url, brand_logo_position, brand_logo_scale, brand_watermark_opacity, overlay_template, transition_effect, background_music_url, background_music_volume, original_audio_volume, created_at, updated_at
		FROM clip_styles WHERE clip_id = $1`
	var s domain.ClipStyle
	err := r.pool.QueryRow(ctx, query, clipID).Scan(&s.ID, &s.ClipID, &s.CaptionEnabled, &s.CaptionFont, &s.CaptionSize, &s.CaptionColor, &s.CaptionBgColor, &s.CaptionPosition, &s.CaptionAnimation, &s.CaptionMaxWords, &s.CaptionLanguage,
		&s.BrandLogoURL, &s.BrandLogoPosition, &s.BrandLogoScale, &s.BrandWatermarkOpacity, &s.OverlayTemplate, &s.TransitionEffect, &s.BackgroundMusicURL, &s.BackgroundMusicVolume, &s.OriginalAudioVolume, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *clipStyleRepository) Update(ctx context.Context, s *domain.ClipStyle) error {
	query := `UPDATE clip_styles SET caption_enabled = $2, caption_font = $3, caption_size = $4, caption_color = $5, caption_bg_color = $6, caption_position = $7, caption_animation = $8, caption_max_words = $9, caption_language = $10,
		brand_logo_url = $11, brand_logo_position = $12, brand_logo_scale = $13, brand_watermark_opacity = $14, overlay_template = $15, transition_effect = $16, background_music_url = $17, background_music_volume = $18, original_audio_volume = $19, updated_at = NOW()
		WHERE clip_id = $1`
	_, err := r.pool.Exec(ctx, query, s.ClipID, s.CaptionEnabled, s.CaptionFont, s.CaptionSize, s.CaptionColor, s.CaptionBgColor, s.CaptionPosition, s.CaptionAnimation, s.CaptionMaxWords, s.CaptionLanguage,
		s.BrandLogoURL, s.BrandLogoPosition, s.BrandLogoScale, s.BrandWatermarkOpacity, s.OverlayTemplate, s.TransitionEffect, s.BackgroundMusicURL, s.BackgroundMusicVolume, s.OriginalAudioVolume)
	return err
}
