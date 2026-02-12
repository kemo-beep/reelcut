package repository

import (
	"context"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type projectRepository struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) ProjectRepository {
	return &projectRepository{pool: pool}
}

func (r *projectRepository) Create(ctx context.Context, p *domain.Project) error {
	query := `INSERT INTO projects (id, user_id, name, description) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, p.ID, p.UserID, p.Name, p.Description)
	return err
}

func (r *projectRepository) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	query := `SELECT id, user_id, name, description, created_at, updated_at FROM projects WHERE id = $1`
	var p domain.Project
	err := r.pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.UserID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *projectRepository) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Project, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM projects WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 20
	}
	query := `SELECT id, user_id, name, description, created_at, updated_at FROM projects WHERE user_id = $1 ORDER BY updated_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []*domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &p)
	}
	return list, total, rows.Err()
}

func (r *projectRepository) Update(ctx context.Context, p *domain.Project) error {
	query := `UPDATE projects SET name = $2, description = $3, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, p.ID, p.Name, p.Description)
	return err
}

func (r *projectRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	return err
}
