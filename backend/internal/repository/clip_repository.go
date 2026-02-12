package repository

import (
	"context"
	"fmt"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type clipRepository struct {
	pool *pgxpool.Pool
}

func NewClipRepository(pool *pgxpool.Pool) ClipRepository {
	return &clipRepository{pool: pool}
}

func (r *clipRepository) Create(ctx context.Context, c *domain.Clip) error {
	query := `INSERT INTO clips (id, video_id, user_id, name, start_time, end_time, duration_seconds, aspect_ratio, virality_score, status, storage_path, thumbnail_url, is_ai_suggested, suggestion_reason, view_count, download_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	_, err := r.pool.Exec(ctx, query, c.ID, c.VideoID, c.UserID, c.Name, c.StartTime, c.EndTime, c.DurationSeconds, c.AspectRatio, c.ViralityScore, c.Status, c.StoragePath, c.ThumbnailURL, c.IsAISuggested, c.SuggestionReason, c.ViewCount, c.DownloadCount)
	return err
}

func (r *clipRepository) GetByID(ctx context.Context, id string) (*domain.Clip, error) {
	query := `SELECT id, video_id, user_id, name, start_time, end_time, duration_seconds, aspect_ratio, virality_score, status, storage_path, thumbnail_url, is_ai_suggested, suggestion_reason, view_count, download_count, created_at, updated_at
		FROM clips WHERE id = $1 AND deleted_at IS NULL`
	var c domain.Clip
	err := r.pool.QueryRow(ctx, query, id).Scan(&c.ID, &c.VideoID, &c.UserID, &c.Name, &c.StartTime, &c.EndTime, &c.DurationSeconds, &c.AspectRatio, &c.ViralityScore, &c.Status, &c.StoragePath, &c.ThumbnailURL, &c.IsAISuggested, &c.SuggestionReason, &c.ViewCount, &c.DownloadCount, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *clipRepository) List(ctx context.Context, userID string, videoID *string, status *string, limit, offset int, sortBy, sortOrder string) ([]*domain.Clip, int, error) {
	countQuery := `SELECT COUNT(*) FROM clips WHERE user_id = $1 AND deleted_at IS NULL`
	countArgs := []interface{}{userID}
	if videoID != nil {
		countQuery += ` AND video_id = $2`
		countArgs = append(countArgs, *videoID)
		if status != nil {
			countQuery += ` AND status = $3`
			countArgs = append(countArgs, *status)
		}
	} else if status != nil {
		countQuery += ` AND status = $2`
		countArgs = append(countArgs, *status)
	}
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 20
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	allowedSort := map[string]bool{"created_at": true, "updated_at": true, "name": true, "virality_score": true}
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	query := `SELECT id, video_id, user_id, name, start_time, end_time, duration_seconds, aspect_ratio, virality_score, status, storage_path, thumbnail_url, is_ai_suggested, suggestion_reason, view_count, download_count, created_at, updated_at
		FROM clips WHERE user_id = $1 AND deleted_at IS NULL`
	queryArgs := []interface{}{userID}
	pos := 2
	if videoID != nil {
		query += fmt.Sprintf(` AND video_id = $%d`, pos)
		queryArgs = append(queryArgs, *videoID)
		pos++
	}
	if status != nil {
		query += fmt.Sprintf(` AND status = $%d`, pos)
		queryArgs = append(queryArgs, *status)
		pos++
	}
	query += fmt.Sprintf(` ORDER BY %s %s LIMIT $%d OFFSET $%d`, sortBy, sortOrder, pos, pos+1)
	queryArgs = append(queryArgs, limit, offset)
	rows, err := r.pool.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*domain.Clip
	for rows.Next() {
		var c domain.Clip
		if err := rows.Scan(&c.ID, &c.VideoID, &c.UserID, &c.Name, &c.StartTime, &c.EndTime, &c.DurationSeconds, &c.AspectRatio, &c.ViralityScore, &c.Status, &c.StoragePath, &c.ThumbnailURL, &c.IsAISuggested, &c.SuggestionReason, &c.ViewCount, &c.DownloadCount, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &c)
	}
	return list, total, rows.Err()
}

func (r *clipRepository) Update(ctx context.Context, c *domain.Clip) error {
	query := `UPDATE clips SET name = $2, start_time = $3, end_time = $4, duration_seconds = $5, aspect_ratio = $6, virality_score = $7, status = $8, storage_path = $9, thumbnail_url = $10, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, c.ID, c.Name, c.StartTime, c.EndTime, c.DurationSeconds, c.AspectRatio, c.ViralityScore, c.Status, c.StoragePath, c.ThumbnailURL)
	return err
}

func (r *clipRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE clips SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
