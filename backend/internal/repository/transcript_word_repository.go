package repository

import (
	"context"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type transcriptWordRepository struct {
	pool *pgxpool.Pool
}

func NewTranscriptWordRepository(pool *pgxpool.Pool) TranscriptWordRepository {
	return &transcriptWordRepository{pool: pool}
}

func (r *transcriptWordRepository) GetBySegmentID(ctx context.Context, segmentID string) ([]*domain.TranscriptWord, error) {
	query := `SELECT id, segment_id, word, start_time, end_time, confidence, sequence_order
		FROM transcript_words WHERE segment_id = $1 ORDER BY sequence_order`
	rows, err := r.pool.Query(ctx, query, segmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.TranscriptWord
	for rows.Next() {
		var w domain.TranscriptWord
		if err := rows.Scan(&w.ID, &w.SegmentID, &w.Word, &w.StartTime, &w.EndTime, &w.Confidence, &w.SequenceOrder); err != nil {
			return nil, err
		}
		list = append(list, &w)
	}
	return list, rows.Err()
}

func (r *transcriptWordRepository) CreateBatch(ctx context.Context, words []*domain.TranscriptWord) error {
	for _, w := range words {
		_, err := r.pool.Exec(ctx, `INSERT INTO transcript_words (id, segment_id, word, start_time, end_time, confidence, sequence_order)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			w.ID, w.SegmentID, w.Word, w.StartTime, w.EndTime, w.Confidence, w.SequenceOrder)
		if err != nil {
			return err
		}
	}
	return nil
}
