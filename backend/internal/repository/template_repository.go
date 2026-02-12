package repository

import (
	"context"
	"fmt"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type templateRepository struct {
	pool *pgxpool.Pool
}

func NewTemplateRepository(pool *pgxpool.Pool) TemplateRepository {
	return &templateRepository{pool: pool}
}

func (r *templateRepository) Create(ctx context.Context, t *domain.Template) error {
	query := `INSERT INTO templates (id, user_id, name, category, is_public, preview_url, style_config, usage_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query, t.ID, t.UserID, t.Name, t.Category, t.IsPublic, t.PreviewURL, t.StyleConfig, t.UsageCount)
	return err
}

func (r *templateRepository) GetByID(ctx context.Context, id string) (*domain.Template, error) {
	query := `SELECT id, user_id, name, category, is_public, preview_url, style_config, usage_count, created_at, updated_at FROM templates WHERE id = $1`
	var t domain.Template
	err := r.pool.QueryRow(ctx, query, id).Scan(&t.ID, &t.UserID, &t.Name, &t.Category, &t.IsPublic, &t.PreviewURL, &t.StyleConfig, &t.UsageCount, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *templateRepository) List(ctx context.Context, userID *string, publicOnly bool, limit, offset int) ([]*domain.Template, int, error) {
	var countQuery string
	var countArgs []interface{}
	if publicOnly {
		countQuery = `SELECT COUNT(*) FROM templates WHERE is_public = true`
	} else if userID != nil {
		countQuery = `SELECT COUNT(*) FROM templates WHERE user_id = $1 OR is_public = true`
		countArgs = append(countArgs, *userID)
	} else {
		countQuery = `SELECT COUNT(*) FROM templates`
	}
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 20
	}
	query := `SELECT id, user_id, name, category, is_public, preview_url, style_config, usage_count, created_at, updated_at FROM templates`
	queryArgs := []interface{}{}
	pos := 1
	if publicOnly {
		query += ` WHERE is_public = true`
	} else if userID != nil {
		query += ` WHERE user_id = $1 OR is_public = true`
		queryArgs = append(queryArgs, *userID)
		pos++
	}
	query += fmt.Sprintf(` ORDER BY updated_at DESC LIMIT $%d OFFSET $%d`, pos, pos+1)
	queryArgs = append(queryArgs, limit, offset)
	rows, err := r.pool.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*domain.Template
	for rows.Next() {
		var t domain.Template
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Category, &t.IsPublic, &t.PreviewURL, &t.StyleConfig, &t.UsageCount, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &t)
	}
	return list, total, rows.Err()
}

func (r *templateRepository) Update(ctx context.Context, t *domain.Template) error {
	query := `UPDATE templates SET name = $2, category = $3, is_public = $4, preview_url = $5, style_config = $6, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, t.ID, t.Name, t.Category, t.IsPublic, t.PreviewURL, t.StyleConfig)
	return err
}

func (r *templateRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM templates WHERE id = $1`, id)
	return err
}

func (r *templateRepository) IncrementUsageCount(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE templates SET usage_count = usage_count + 1, updated_at = NOW() WHERE id = $1`, id)
	return err
}
