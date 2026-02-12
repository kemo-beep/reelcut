package repository

import (
	"context"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userSessionRepository struct {
	pool *pgxpool.Pool
}

func NewUserSessionRepository(pool *pgxpool.Pool) UserSessionRepository {
	return &userSessionRepository{pool: pool}
}

func (r *userSessionRepository) Create(ctx context.Context, s *domain.UserSession) error {
	query := `INSERT INTO user_sessions (id, user_id, token, expires_at) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, s.ID, s.UserID, s.Token, s.ExpiresAt)
	return err
}

func (r *userSessionRepository) GetByToken(ctx context.Context, token string) (*domain.UserSession, error) {
	query := `SELECT id, user_id, token, expires_at, created_at FROM user_sessions WHERE token = $1`
	var s domain.UserSession
	err := r.pool.QueryRow(ctx, query, token).Scan(&s.ID, &s.UserID, &s.Token, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *userSessionRepository) DeleteByToken(ctx context.Context, token string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM user_sessions WHERE token = $1`, token)
	return err
}

func (r *userSessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM user_sessions WHERE user_id = $1`, userID)
	return err
}
