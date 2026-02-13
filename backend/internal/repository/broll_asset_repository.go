package repository

import (
	"context"
	"fmt"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type brollAssetRepository struct {
	pool *pgxpool.Pool
}

func NewBrollAssetRepository(pool *pgxpool.Pool) BrollAssetRepository {
	return &brollAssetRepository{pool: pool}
}

func (r *brollAssetRepository) Create(ctx context.Context, a *domain.BrollAsset) error {
	query := `INSERT INTO broll_assets (id, user_id, project_id, original_filename, storage_path, duration_seconds, width, height)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query, a.ID, a.UserID, a.ProjectID, a.OriginalFilename, a.StoragePath, a.DurationSeconds, a.Width, a.Height)
	return err
}

func (r *brollAssetRepository) GetByID(ctx context.Context, id string) (*domain.BrollAsset, error) {
	query := `SELECT id, user_id, project_id, original_filename, storage_path, duration_seconds, width, height, created_at, updated_at
		FROM broll_assets WHERE id = $1`
	var a domain.BrollAsset
	err := r.pool.QueryRow(ctx, query, id).Scan(&a.ID, &a.UserID, &a.ProjectID, &a.OriginalFilename, &a.StoragePath, &a.DurationSeconds, &a.Width, &a.Height, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *brollAssetRepository) ListByUserID(ctx context.Context, userID string, projectID *string, limit, offset int) ([]*domain.BrollAsset, int, error) {
	countQuery := `SELECT COUNT(*) FROM broll_assets WHERE user_id = $1`
	countArgs := []interface{}{userID}
	if projectID != nil {
		countQuery += ` AND project_id = $2`
		countArgs = append(countArgs, *projectID)
	}
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 50
	}
	listQuery := `SELECT id, user_id, project_id, original_filename, storage_path, duration_seconds, width, height, created_at, updated_at
		FROM broll_assets WHERE user_id = $1`
	listArgs := []interface{}{userID}
	if projectID != nil {
		listQuery += ` AND project_id = $2`
		listArgs = append(listArgs, *projectID)
	}
	listArgs = append(listArgs, limit, offset)
	listQuery += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(listArgs)-1) + ` OFFSET $` + fmt.Sprintf("%d", len(listArgs))
	rows, err := r.pool.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*domain.BrollAsset
	for rows.Next() {
		var a domain.BrollAsset
		if err := rows.Scan(&a.ID, &a.UserID, &a.ProjectID, &a.OriginalFilename, &a.StoragePath, &a.DurationSeconds, &a.Width, &a.Height, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &a)
	}
	return list, total, rows.Err()
}

func (r *brollAssetRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM broll_assets WHERE id = $1`, id)
	return err
}
