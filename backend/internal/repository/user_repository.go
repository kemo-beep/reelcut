package repository

import (
	"context"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users (id, email, password_hash, full_name, avatar_url, subscription_tier, credits_remaining, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query,
		u.ID, u.Email, u.PasswordHash, u.FullName, u.AvatarURL, u.SubscriptionTier, u.CreditsRemaining, u.EmailVerified)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, full_name, avatar_url, subscription_tier, credits_remaining, email_verified, created_at, updated_at
		FROM users WHERE id = $1 AND deleted_at IS NULL`
	var u domain.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.AvatarURL, &u.SubscriptionTier, &u.CreditsRemaining, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, full_name, avatar_url, subscription_tier, credits_remaining, email_verified, created_at, updated_at
		FROM users WHERE email = $1 AND deleted_at IS NULL`
	var u domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.AvatarURL, &u.SubscriptionTier, &u.CreditsRemaining, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Update(ctx context.Context, u *domain.User) error {
	query := `UPDATE users SET full_name = $2, avatar_url = $3, subscription_tier = $4, credits_remaining = $5, email_verified = $6, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, u.ID, u.FullName, u.AvatarURL, u.SubscriptionTier, u.CreditsRemaining, u.EmailVerified)
	return err
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID string, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, userID, passwordHash)
	return err
}

func (r *userRepository) SetEmailVerified(ctx context.Context, userID string, verified bool) error {
	query := `UPDATE users SET email_verified = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, userID, verified)
	return err
}

func (r *userRepository) DeductCredits(ctx context.Context, userID string, amount int) error {
	query := `UPDATE users SET credits_remaining = credits_remaining - $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL AND credits_remaining >= $2`
	tag, err := r.pool.Exec(ctx, query, userID, amount)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrInsufficientCredits
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE users SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

