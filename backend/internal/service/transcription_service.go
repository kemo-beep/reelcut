package service

import (
	"context"

	"reelcut/internal/ai"
	"reelcut/internal/domain"
	"reelcut/internal/queue"
	"reelcut/internal/repository"

	"github.com/google/uuid"
)

type TranscriptionService struct {
	transcriptionRepo repository.TranscriptionRepository
	segmentRepo       repository.TranscriptSegmentRepository
	wordRepo          repository.TranscriptWordRepository
	videoRepo         repository.VideoRepository
	queue             *queue.QueueClient
}

func NewTranscriptionService(
	transcriptionRepo repository.TranscriptionRepository,
	segmentRepo repository.TranscriptSegmentRepository,
	wordRepo repository.TranscriptWordRepository,
	videoRepo repository.VideoRepository,
	queue *queue.QueueClient,
) *TranscriptionService {
	return &TranscriptionService{
		transcriptionRepo: transcriptionRepo,
		segmentRepo:      segmentRepo,
		wordRepo:         wordRepo,
		videoRepo:        videoRepo,
		queue:            queue,
	}
}

func (s *TranscriptionService) Create(ctx context.Context, videoID uuid.UUID, language string, _ bool) (*domain.Transcription, error) {
	v, err := s.videoRepo.GetByID(ctx, videoID.String())
	if err != nil || v == nil {
		return nil, domain.ErrNotFound
	}
	if v.Status != "ready" {
		return nil, domain.ErrValidation
	}
	t := &domain.Transcription{
		ID:       uuid.New(),
		VideoID:  videoID,
		Language: language,
		Status:   "pending",
	}
	if err := s.transcriptionRepo.Create(ctx, t); err != nil {
		return nil, err
	}
	if err := s.queue.EnqueueTranscription(videoID, t.ID); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TranscriptionService) GetByID(ctx context.Context, id string) (*domain.Transcription, error) {
	t, err := s.transcriptionRepo.GetByID(ctx, id)
	if err != nil || t == nil {
		return nil, domain.ErrNotFound
	}
	segments, _ := s.segmentRepo.GetByTranscriptionID(ctx, id)
	for _, seg := range segments {
		words, _ := s.wordRepo.GetBySegmentID(ctx, seg.ID.String())
		for _, w := range words {
			seg.Words = append(seg.Words, *w)
		}
		t.Segments = append(t.Segments, *seg)
	}
	return t, nil
}

func (s *TranscriptionService) GetByVideoID(ctx context.Context, videoID string) (*domain.Transcription, error) {
	t, err := s.transcriptionRepo.GetByVideoID(ctx, videoID)
	if err != nil || t == nil {
		return nil, domain.ErrNotFound
	}
	return s.GetByID(ctx, t.ID.String())
}

func (s *TranscriptionService) GetByVideoIDAndLanguage(ctx context.Context, videoID, language string) (*domain.Transcription, error) {
	t, err := s.transcriptionRepo.GetByVideoIDAndLanguage(ctx, videoID, language)
	if err != nil || t == nil {
		return nil, domain.ErrNotFound
	}
	return s.GetByID(ctx, t.ID.String())
}

func (s *TranscriptionService) ListCompletedByVideoID(ctx context.Context, videoID string) ([]*domain.Transcription, error) {
	return s.transcriptionRepo.ListCompletedByVideoID(ctx, videoID)
}

// Translate creates a new transcription with the same timestamps but segment text translated to targetLang.
// Uses Gemini to translate. The new transcription has source_transcription_id set to the source.
func (s *TranscriptionService) Translate(ctx context.Context, sourceTranscriptionID string, targetLang string) (*domain.Transcription, error) {
	source, err := s.GetByID(ctx, sourceTranscriptionID)
	if err != nil || source == nil || source.Status != "completed" {
		return nil, domain.ErrNotFound
	}
	if len(source.Segments) == 0 {
		return nil, domain.ErrValidation
	}
	texts := make([]string, 0, len(source.Segments))
	for _, seg := range source.Segments {
		texts = append(texts, seg.Text)
	}
	translated, err := ai.TranslateSegments(ctx, texts, targetLang)
	if err != nil {
		return nil, err
	}
	sourceID, _ := uuid.Parse(sourceTranscriptionID)
	newT := &domain.Transcription{
		ID:                    uuid.New(),
		VideoID:               source.VideoID,
		Language:              targetLang,
		Status:                "completed",
		SourceTranscriptionID: &sourceID,
		WordCount:             source.WordCount,
		DurationSeconds:       source.DurationSeconds,
		ConfidenceAvg:         source.ConfidenceAvg,
	}
	newSegments := make([]*domain.TranscriptSegment, 0, len(source.Segments))
	for i, seg := range source.Segments {
		translatedText := seg.Text
		if i < len(translated) {
			translatedText = translated[i]
		}
		newSeg := &domain.TranscriptSegment{
			ID:             uuid.New(),
			TranscriptionID: newT.ID,
			StartTime:      seg.StartTime,
			EndTime:        seg.EndTime,
			Text:           translatedText,
			Confidence:     seg.Confidence,
			SpeakerID:      seg.SpeakerID,
			SequenceOrder:  seg.SequenceOrder,
			Words: []domain.TranscriptWord{
				{ID: uuid.New(), SegmentID: uuid.Nil, Word: translatedText, StartTime: seg.StartTime, EndTime: seg.EndTime, SequenceOrder: 0},
			},
		}
		newSeg.Words[0].SegmentID = newSeg.ID
		newSeg.Words[0].ID = uuid.New()
		newSegments = append(newSegments, newSeg)
	}
	if err := s.transcriptionRepo.CreateWithSegments(ctx, newT, newSegments); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, newT.ID.String())
}

func (s *TranscriptionService) UpdateSegment(ctx context.Context, transcriptionID, segmentID string, text string, startTime, endTime float64) error {
	t, err := s.transcriptionRepo.GetByID(ctx, transcriptionID)
	if err != nil || t == nil {
		return domain.ErrNotFound
	}
	segments, err := s.segmentRepo.GetByTranscriptionID(ctx, transcriptionID)
	if err != nil {
		return err
	}
	for _, seg := range segments {
		if seg.ID.String() == segmentID {
			seg.Text = text
			if startTime > 0 {
				seg.StartTime = startTime
			}
			if endTime > 0 {
				seg.EndTime = endTime
			}
			return s.segmentRepo.Update(ctx, seg)
		}
	}
	return domain.ErrNotFound
}
