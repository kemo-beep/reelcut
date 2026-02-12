package repository

import (
	"context"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type subscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) SubscriptionRepository {
	return &subscriptionRepository{pool: pool}
}

func (r *subscriptionRepository) Create(ctx context.Context, s *domain.Subscription) error {
	query := `INSERT INTO subscriptions (id, user_id, tier, status, stripe_subscription_id, current_period_start, current_period_end)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query,
		s.ID, s.UserID, s.Tier, s.Status, s.StripeSubscriptionID, s.CurrentPeriodStart, s.CurrentPeriodEnd)
	return err
}

func (r *subscriptionRepository) GetByUserID(ctx context.Context, userID string) (*domain.Subscription, error) {
	query := `SELECT id, user_id, tier, status, stripe_subscription_id, current_period_start, current_period_end, created_at, updated_at
		FROM subscriptions WHERE user_id = $1 ORDER BY updated_at DESC LIMIT 1`
	var s domain.Subscription
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&s.ID, &s.UserID, &s.Tier, &s.Status, &s.StripeSubscriptionID, &s.CurrentPeriodStart, &s.CurrentPeriodEnd, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *subscriptionRepository) GetByStripeID(ctx context.Context, stripeSubscriptionID string) (*domain.Subscription, error) {
	query := `SELECT id, user_id, tier, status, stripe_subscription_id, current_period_start, current_period_end, created_at, updated_at
		FROM subscriptions WHERE stripe_subscription_id = $1`
	var s domain.Subscription
	err := r.pool.QueryRow(ctx, query, stripeSubscriptionID).Scan(
		&s.ID, &s.UserID, &s.Tier, &s.Status, &s.StripeSubscriptionID, &s.CurrentPeriodStart, &s.CurrentPeriodEnd, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, s *domain.Subscription) error {
	query := `UPDATE subscriptions SET tier = $2, status = $3, stripe_subscription_id = $4, current_period_start = $5, current_period_end = $6, updated_at = NOW()
		WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, s.ID, s.Tier, s.Status, s.StripeSubscriptionID, s.CurrentPeriodStart, s.CurrentPeriodEnd)
	return err
}
