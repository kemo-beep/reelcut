package repository

import (
	"context"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type videoAnalysisRepository struct {
	pool *pgxpool.Pool
}

func NewVideoAnalysisRepository(pool *pgxpool.Pool) VideoAnalysisRepository {
	return &videoAnalysisRepository{pool: pool}
}

func (r *videoAnalysisRepository) GetByVideoID(ctx context.Context, videoID string) (*domain.VideoAnalysis, error) {
	query := `SELECT id, video_id, scenes_detected, faces_detected, topics, sentiment_analysis, engagement_scores, created_at
		FROM video_analysis WHERE video_id = $1`
	var a domain.VideoAnalysis
	err := r.pool.QueryRow(ctx, query, videoID).Scan(&a.ID, &a.VideoID, &a.ScenesDetected, &a.FacesDetected, &a.Topics, &a.SentimentAnalysis, &a.EngagementScores, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *videoAnalysisRepository) Upsert(ctx context.Context, a *domain.VideoAnalysis) error {
	existing, _ := r.GetByVideoID(ctx, a.VideoID.String())
	if existing != nil {
		query := `UPDATE video_analysis SET scenes_detected = $2, faces_detected = $3, topics = $4, sentiment_analysis = $5, engagement_scores = $6 WHERE video_id = $1`
		_, err := r.pool.Exec(ctx, query, a.VideoID, a.ScenesDetected, a.FacesDetected, a.Topics, a.SentimentAnalysis, a.EngagementScores)
		return err
	}
	query := `INSERT INTO video_analysis (id, video_id, scenes_detected, faces_detected, topics, sentiment_analysis, engagement_scores)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query, a.ID, a.VideoID, a.ScenesDetected, a.FacesDetected, a.Topics, a.SentimentAnalysis, a.EngagementScores)
	return err
}
