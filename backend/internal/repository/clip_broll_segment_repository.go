package repository

import (
	"context"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type clipBrollSegmentRepository struct {
	pool *pgxpool.Pool
}

func NewClipBrollSegmentRepository(pool *pgxpool.Pool) ClipBrollSegmentRepository {
	return &clipBrollSegmentRepository{pool: pool}
}

func (r *clipBrollSegmentRepository) Create(ctx context.Context, s *domain.ClipBrollSegment) error {
	query := `INSERT INTO clip_broll_segments (id, clip_id, broll_asset_id, start_time, end_time, position, scale, opacity, sequence_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, query, s.ID, s.ClipID, s.BrollAssetID, s.StartTime, s.EndTime, s.Position, s.Scale, s.Opacity, s.SequenceOrder)
	return err
}

func (r *clipBrollSegmentRepository) GetByID(ctx context.Context, id string) (*domain.ClipBrollSegment, error) {
	query := `SELECT id, clip_id, broll_asset_id, start_time, end_time, position, scale, opacity, sequence_order, created_at, updated_at
		FROM clip_broll_segments WHERE id = $1`
	var s domain.ClipBrollSegment
	err := r.pool.QueryRow(ctx, query, id).Scan(&s.ID, &s.ClipID, &s.BrollAssetID, &s.StartTime, &s.EndTime, &s.Position, &s.Scale, &s.Opacity, &s.SequenceOrder, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *clipBrollSegmentRepository) GetByClipID(ctx context.Context, clipID string) ([]*domain.ClipBrollSegment, error) {
	query := `SELECT s.id, s.clip_id, s.broll_asset_id, s.start_time, s.end_time, s.position, s.scale, s.opacity, s.sequence_order, s.created_at, s.updated_at,
		a.id, a.user_id, a.project_id, a.original_filename, a.storage_path, a.duration_seconds, a.width, a.height, a.created_at, a.updated_at
		FROM clip_broll_segments s
		JOIN broll_assets a ON a.id = s.broll_asset_id
		WHERE s.clip_id = $1 ORDER BY s.sequence_order`
	rows, err := r.pool.Query(ctx, query, clipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.ClipBrollSegment
	for rows.Next() {
		var s domain.ClipBrollSegment
		var a domain.BrollAsset
		err := rows.Scan(&s.ID, &s.ClipID, &s.BrollAssetID, &s.StartTime, &s.EndTime, &s.Position, &s.Scale, &s.Opacity, &s.SequenceOrder, &s.CreatedAt, &s.UpdatedAt,
			&a.ID, &a.UserID, &a.ProjectID, &a.OriginalFilename, &a.StoragePath, &a.DurationSeconds, &a.Width, &a.Height, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		s.Asset = &a
		list = append(list, &s)
	}
	return list, rows.Err()
}

func (r *clipBrollSegmentRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM clip_broll_segments WHERE id = $1`, id)
	return err
}

func (r *clipBrollSegmentRepository) DeleteByClipID(ctx context.Context, clipID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM clip_broll_segments WHERE clip_id = $1`, clipID)
	return err
}

func (r *clipBrollSegmentRepository) NextSequenceOrder(ctx context.Context, clipID string) (int, error) {
	var max *int
	err := r.pool.QueryRow(ctx, `SELECT MAX(sequence_order) FROM clip_broll_segments WHERE clip_id = $1`, clipID).Scan(&max)
	if err != nil || max == nil {
		return 0, err
	}
	return *max + 1, nil
}
