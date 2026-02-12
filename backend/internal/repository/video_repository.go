package repository

import (
	"context"
	"fmt"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type videoRepository struct {
	pool *pgxpool.Pool
}

func NewVideoRepository(pool *pgxpool.Pool) VideoRepository {
	return &videoRepository{pool: pool}
}

func (r *videoRepository) Create(ctx context.Context, v *domain.Video) error {
	query := `INSERT INTO videos (id, project_id, user_id, original_filename, storage_path, thumbnail_url, duration_seconds, width, height, fps, file_size_bytes, codec, bitrate, status, error_message, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	_, err := r.pool.Exec(ctx, query,
		v.ID, v.ProjectID, v.UserID, v.OriginalFilename, v.StoragePath, v.ThumbnailURL, v.DurationSeconds, v.Width, v.Height, v.FPS, v.FileSizeBytes, v.Codec, v.Bitrate, v.Status, v.ErrorMessage, v.Metadata)
	return err
}

func (r *videoRepository) GetByID(ctx context.Context, id string) (*domain.Video, error) {
	query := `SELECT id, project_id, user_id, original_filename, storage_path, thumbnail_url, duration_seconds, width, height, fps, file_size_bytes, codec, bitrate, status, error_message, metadata, created_at, updated_at
		FROM videos WHERE id = $1 AND deleted_at IS NULL`
	var v domain.Video
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&v.ID, &v.ProjectID, &v.UserID, &v.OriginalFilename, &v.StoragePath, &v.ThumbnailURL, &v.DurationSeconds, &v.Width, &v.Height, &v.FPS, &v.FileSizeBytes, &v.Codec, &v.Bitrate, &v.Status, &v.ErrorMessage, &v.Metadata, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *videoRepository) List(ctx context.Context, userID string, projectID *string, status *string, limit, offset int, sortBy, sortOrder string) ([]*domain.Video, int, error) {
	countQuery := `SELECT COUNT(*) FROM videos WHERE user_id = $1 AND deleted_at IS NULL`
	countArgs := []interface{}{userID}
	if projectID != nil {
		countQuery += ` AND project_id = $2`
		countArgs = append(countArgs, *projectID)
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
	allowedSort := map[string]bool{"created_at": true, "updated_at": true, "original_filename": true, "duration_seconds": true}
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	query := `SELECT id, project_id, user_id, original_filename, storage_path, thumbnail_url, duration_seconds, width, height, fps, file_size_bytes, codec, bitrate, status, error_message, metadata, created_at, updated_at
		FROM videos WHERE user_id = $1 AND deleted_at IS NULL`
	queryArgs := []interface{}{userID}
	pos := 2
	if projectID != nil {
		query += fmt.Sprintf(` AND project_id = $%d`, pos)
		queryArgs = append(queryArgs, *projectID)
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
	var list []*domain.Video
	for rows.Next() {
		var v domain.Video
		if err := rows.Scan(&v.ID, &v.ProjectID, &v.UserID, &v.OriginalFilename, &v.StoragePath, &v.ThumbnailURL, &v.DurationSeconds, &v.Width, &v.Height, &v.FPS, &v.FileSizeBytes, &v.Codec, &v.Bitrate, &v.Status, &v.ErrorMessage, &v.Metadata, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &v)
	}
	return list, total, rows.Err()
}

func (r *videoRepository) Update(ctx context.Context, v *domain.Video) error {
	query := `UPDATE videos SET thumbnail_url = $2, duration_seconds = $3, width = $4, height = $5, fps = $6, file_size_bytes = $7, codec = $8, bitrate = $9, status = $10, error_message = $11, metadata = $12, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, v.ID, v.ThumbnailURL, v.DurationSeconds, v.Width, v.Height, v.FPS, v.FileSizeBytes, v.Codec, v.Bitrate, v.Status, v.ErrorMessage, v.Metadata)
	return err
}

func (r *videoRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE videos SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
