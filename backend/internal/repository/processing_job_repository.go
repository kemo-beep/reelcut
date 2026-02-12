package repository

import (
	"context"
	"fmt"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type processingJobRepository struct {
	pool *pgxpool.Pool
}

func NewProcessingJobRepository(pool *pgxpool.Pool) ProcessingJobRepository {
	return &processingJobRepository{pool: pool}
}

func (r *processingJobRepository) Create(ctx context.Context, j *domain.ProcessingJob) error {
	query := `INSERT INTO processing_jobs (id, user_id, job_type, entity_type, entity_id, priority, status, progress, error_message, retry_count, max_retries, metadata, started_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	_, err := r.pool.Exec(ctx, query, j.ID, j.UserID, j.JobType, j.EntityType, j.EntityID, j.Priority, j.Status, j.Progress, j.ErrorMessage, j.RetryCount, j.MaxRetries, j.Metadata, j.StartedAt, j.CompletedAt)
	return err
}

func (r *processingJobRepository) GetByID(ctx context.Context, id string) (*domain.ProcessingJob, error) {
	query := `SELECT id, user_id, job_type, entity_type, entity_id, priority, status, progress, error_message, retry_count, max_retries, metadata, started_at, completed_at, created_at, updated_at
		FROM processing_jobs WHERE id = $1`
	var j domain.ProcessingJob
	err := r.pool.QueryRow(ctx, query, id).Scan(&j.ID, &j.UserID, &j.JobType, &j.EntityType, &j.EntityID, &j.Priority, &j.Status, &j.Progress, &j.ErrorMessage, &j.RetryCount, &j.MaxRetries, &j.Metadata, &j.StartedAt, &j.CompletedAt, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *processingJobRepository) ListByUserID(ctx context.Context, userID string, status *string, limit, offset int) ([]*domain.ProcessingJob, int, error) {
	countQuery := `SELECT COUNT(*) FROM processing_jobs WHERE user_id = $1`
	countArgs := []interface{}{userID}
	if status != nil {
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
	query := `SELECT id, user_id, job_type, entity_type, entity_id, priority, status, progress, error_message, retry_count, max_retries, metadata, started_at, completed_at, created_at, updated_at
		FROM processing_jobs WHERE user_id = $1`
	queryArgs := []interface{}{userID}
	pos := 2
	if status != nil {
		query += fmt.Sprintf(` AND status = $%d`, pos)
		queryArgs = append(queryArgs, *status)
		pos++
	}
	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, pos, pos+1)
	queryArgs = append(queryArgs, limit, offset)
	rows, err := r.pool.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*domain.ProcessingJob
	for rows.Next() {
		var j domain.ProcessingJob
		if err := rows.Scan(&j.ID, &j.UserID, &j.JobType, &j.EntityType, &j.EntityID, &j.Priority, &j.Status, &j.Progress, &j.ErrorMessage, &j.RetryCount, &j.MaxRetries, &j.Metadata, &j.StartedAt, &j.CompletedAt, &j.CreatedAt, &j.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &j)
	}
	return list, total, rows.Err()
}

func (r *processingJobRepository) GetByEntity(ctx context.Context, entityType, entityID string) (*domain.ProcessingJob, error) {
	query := `SELECT id, user_id, job_type, entity_type, entity_id, priority, status, progress, error_message, retry_count, max_retries, metadata, started_at, completed_at, created_at, updated_at
		FROM processing_jobs WHERE entity_type = $1 AND entity_id = $2 ORDER BY created_at DESC LIMIT 1`
	var j domain.ProcessingJob
	err := r.pool.QueryRow(ctx, query, entityType, entityID).Scan(&j.ID, &j.UserID, &j.JobType, &j.EntityType, &j.EntityID, &j.Priority, &j.Status, &j.Progress, &j.ErrorMessage, &j.RetryCount, &j.MaxRetries, &j.Metadata, &j.StartedAt, &j.CompletedAt, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *processingJobRepository) Update(ctx context.Context, j *domain.ProcessingJob) error {
	query := `UPDATE processing_jobs SET status = $2, progress = $3, error_message = $4, retry_count = $5, metadata = $6, started_at = $7, completed_at = $8, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, j.ID, j.Status, j.Progress, j.ErrorMessage, j.RetryCount, j.Metadata, j.StartedAt, j.CompletedAt)
	return err
}
