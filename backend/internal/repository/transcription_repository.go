package repository

import (
	"context"

	"reelcut/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type transcriptionRepository struct {
	pool *pgxpool.Pool
}

func NewTranscriptionRepository(pool *pgxpool.Pool) TranscriptionRepository {
	return &transcriptionRepository{pool: pool}
}

func (r *transcriptionRepository) Create(ctx context.Context, t *domain.Transcription) error {
	query := `INSERT INTO transcriptions (id, video_id, language, status, source_transcription_id, error_message, word_count, duration_seconds, confidence_avg)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, query, t.ID, t.VideoID, t.Language, t.Status, t.SourceTranscriptionID, t.ErrorMessage, t.WordCount, t.DurationSeconds, t.ConfidenceAvg)
	return err
}

func (r *transcriptionRepository) GetByID(ctx context.Context, id string) (*domain.Transcription, error) {
	query := `SELECT id, video_id, language, status, source_transcription_id, error_message, word_count, duration_seconds, confidence_avg, created_at, updated_at
		FROM transcriptions WHERE id = $1`
	var t domain.Transcription
	err := r.pool.QueryRow(ctx, query, id).Scan(&t.ID, &t.VideoID, &t.Language, &t.Status, &t.SourceTranscriptionID, &t.ErrorMessage, &t.WordCount, &t.DurationSeconds, &t.ConfidenceAvg, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *transcriptionRepository) GetByVideoID(ctx context.Context, videoID string) (*domain.Transcription, error) {
	queryCompleted := `SELECT id, video_id, language, status, source_transcription_id, error_message, word_count, duration_seconds, confidence_avg, created_at, updated_at
		FROM transcriptions WHERE video_id = $1 AND status = 'completed' ORDER BY created_at DESC LIMIT 1`
	var t domain.Transcription
	err := r.pool.QueryRow(ctx, queryCompleted, videoID).Scan(&t.ID, &t.VideoID, &t.Language, &t.Status, &t.SourceTranscriptionID, &t.ErrorMessage, &t.WordCount, &t.DurationSeconds, &t.ConfidenceAvg, &t.CreatedAt, &t.UpdatedAt)
	if err == nil {
		return &t, nil
	}
	queryLatest := `SELECT id, video_id, language, status, source_transcription_id, error_message, word_count, duration_seconds, confidence_avg, created_at, updated_at
		FROM transcriptions WHERE video_id = $1 ORDER BY created_at DESC LIMIT 1`
	err = r.pool.QueryRow(ctx, queryLatest, videoID).Scan(&t.ID, &t.VideoID, &t.Language, &t.Status, &t.SourceTranscriptionID, &t.ErrorMessage, &t.WordCount, &t.DurationSeconds, &t.ConfidenceAvg, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *transcriptionRepository) GetByVideoIDAndLanguage(ctx context.Context, videoID, language string) (*domain.Transcription, error) {
	query := `SELECT id, video_id, language, status, source_transcription_id, error_message, word_count, duration_seconds, confidence_avg, created_at, updated_at
		FROM transcriptions WHERE video_id = $1 AND language = $2 AND status = 'completed' ORDER BY created_at DESC LIMIT 1`
	var t domain.Transcription
	err := r.pool.QueryRow(ctx, query, videoID, language).Scan(&t.ID, &t.VideoID, &t.Language, &t.Status, &t.SourceTranscriptionID, &t.ErrorMessage, &t.WordCount, &t.DurationSeconds, &t.ConfidenceAvg, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *transcriptionRepository) ListCompletedByVideoID(ctx context.Context, videoID string) ([]*domain.Transcription, error) {
	query := `SELECT id, video_id, language, status, source_transcription_id, error_message, word_count, duration_seconds, confidence_avg, created_at, updated_at
		FROM transcriptions WHERE video_id = $1 AND status = 'completed' ORDER BY language`
	rows, err := r.pool.Query(ctx, query, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Transcription
	for rows.Next() {
		var t domain.Transcription
		if err := rows.Scan(&t.ID, &t.VideoID, &t.Language, &t.Status, &t.SourceTranscriptionID, &t.ErrorMessage, &t.WordCount, &t.DurationSeconds, &t.ConfidenceAvg, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	return out, rows.Err()
}

func (r *transcriptionRepository) Update(ctx context.Context, t *domain.Transcription) error {
	query := `UPDATE transcriptions SET status = $2, error_message = $3, word_count = $4, duration_seconds = $5, confidence_avg = $6, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, t.ID, t.Status, t.ErrorMessage, t.WordCount, t.DurationSeconds, t.ConfidenceAvg)
	return err
}

func (r *transcriptionRepository) CreateWithSegments(ctx context.Context, t *domain.Transcription, segments []*domain.TranscriptSegment) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	query := `INSERT INTO transcriptions (id, video_id, language, status, source_transcription_id, error_message, word_count, duration_seconds, confidence_avg)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	if _, err := tx.Exec(ctx, query, t.ID, t.VideoID, t.Language, t.Status, t.SourceTranscriptionID, t.ErrorMessage, t.WordCount, t.DurationSeconds, t.ConfidenceAvg); err != nil {
		return err
	}
	for _, s := range segments {
		if _, err := tx.Exec(ctx, `INSERT INTO transcript_segments (id, transcription_id, start_time, end_time, text, confidence, speaker_id, sequence_order)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			s.ID, s.TranscriptionID, s.StartTime, s.EndTime, s.Text, s.Confidence, s.SpeakerID, s.SequenceOrder); err != nil {
			return err
		}
		for _, w := range s.Words {
			if _, err := tx.Exec(ctx, `INSERT INTO transcript_words (id, segment_id, word, start_time, end_time, confidence, sequence_order)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				w.ID, w.SegmentID, w.Word, w.StartTime, w.EndTime, w.Confidence, w.SequenceOrder); err != nil {
				return err
			}
		}
	}
	return tx.Commit(ctx)
}
