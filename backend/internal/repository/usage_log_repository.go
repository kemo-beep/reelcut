package repository

import (
	"context"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type usageLogRepository struct {
	pool *pgxpool.Pool
}

func NewUsageLogRepository(pool *pgxpool.Pool) UsageLogRepository {
	return &usageLogRepository{pool: pool}
}

func (r *usageLogRepository) Create(ctx context.Context, u *domain.UsageLog) error {
	query := `INSERT INTO usage_logs (id, user_id, action, credits_used, video_duration_seconds, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query, u.ID, u.UserID, u.Action, u.CreditsUsed, u.VideoDurationSeconds, u.Metadata)
	return err
}

func (r *usageLogRepository) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.UsageLog, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM usage_logs WHERE user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 20
	}
	query := `SELECT id, user_id, action, credits_used, video_duration_seconds, metadata, created_at
		FROM usage_logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*domain.UsageLog
	for rows.Next() {
		var u domain.UsageLog
		if err := rows.Scan(&u.ID, &u.UserID, &u.Action, &u.CreditsUsed, &u.VideoDurationSeconds, &u.Metadata, &u.CreatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &u)
	}
	return list, total, rows.Err()
}
