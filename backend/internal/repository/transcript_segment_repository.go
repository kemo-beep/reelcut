package repository

import (
	"context"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type transcriptSegmentRepository struct {
	pool *pgxpool.Pool
}

func NewTranscriptSegmentRepository(pool *pgxpool.Pool) TranscriptSegmentRepository {
	return &transcriptSegmentRepository{pool: pool}
}

func (r *transcriptSegmentRepository) GetByTranscriptionID(ctx context.Context, transcriptionID string) ([]*domain.TranscriptSegment, error) {
	query := `SELECT id, transcription_id, start_time, end_time, text, confidence, speaker_id, sequence_order
		FROM transcript_segments WHERE transcription_id = $1 ORDER BY sequence_order`
	rows, err := r.pool.Query(ctx, query, transcriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.TranscriptSegment
	for rows.Next() {
		var s domain.TranscriptSegment
		if err := rows.Scan(&s.ID, &s.TranscriptionID, &s.StartTime, &s.EndTime, &s.Text, &s.Confidence, &s.SpeakerID, &s.SequenceOrder); err != nil {
			return nil, err
		}
		list = append(list, &s)
	}
	return list, rows.Err()
}

func (r *transcriptSegmentRepository) Update(ctx context.Context, s *domain.TranscriptSegment) error {
	query := `UPDATE transcript_segments SET text = $2, start_time = $3, end_time = $4, confidence = $5, speaker_id = $6, sequence_order = $7 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, s.ID, s.Text, s.StartTime, s.EndTime, s.Confidence, s.SpeakerID, s.SequenceOrder)
	return err
}
