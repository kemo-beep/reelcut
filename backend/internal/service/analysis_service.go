package service

import (
	"context"
	"log/slog"
	"os"

	"reelcut/internal/ai"
	"reelcut/internal/domain"
	"reelcut/internal/queue"
	"reelcut/internal/repository"

	"github.com/google/uuid"
)

type AnalysisService struct {
	videoAnalysisRepo repository.VideoAnalysisRepository
	transcriptionRepo repository.TranscriptionRepository
	segmentRepo       repository.TranscriptSegmentRepository
	videoRepo         repository.VideoRepository
	queue             *queue.QueueClient
}

func NewAnalysisService(
	videoAnalysisRepo repository.VideoAnalysisRepository,
	transcriptionRepo repository.TranscriptionRepository,
	segmentRepo repository.TranscriptSegmentRepository,
	videoRepo repository.VideoRepository,
	queue *queue.QueueClient,
) *AnalysisService {
	return &AnalysisService{
		videoAnalysisRepo: videoAnalysisRepo,
		transcriptionRepo: transcriptionRepo,
		segmentRepo:      segmentRepo,
		videoRepo:        videoRepo,
		queue:            queue,
	}
}

func (s *AnalysisService) Analyze(ctx context.Context, videoID uuid.UUID) error {
	if err := s.queue.EnqueueAnalysis(videoID); err != nil {
		return err
	}
	return nil
}

func (s *AnalysisService) GetByVideoID(ctx context.Context, videoID string) (*domain.VideoAnalysis, error) {
	return s.videoAnalysisRepo.GetByVideoID(ctx, videoID)
}

func (s *AnalysisService) SuggestClips(ctx context.Context, videoID string, minDur, maxDur float64, maxSuggestions int) ([]ai.ClipSuggestion, error) {
	t, err := s.transcriptionRepo.GetByVideoID(ctx, videoID)
	if err != nil || t == nil {
		return nil, domain.ErrNotFound
	}
	segments, _ := s.segmentRepo.GetByTranscriptionID(ctx, t.ID.String())
	for _, seg := range segments {
		t.Segments = append(t.Segments, *seg)
	}
	if os.Getenv("GEMINI_API_KEY") != "" {
		suggestions, geminiErr := ai.SuggestClipsWithGemini(ctx, t, minDur, maxDur, maxSuggestions)
		if geminiErr != nil {
			slog.Warn("gemini suggest clips failed, falling back to rule-based", "err", geminiErr)
			return ai.SuggestClips(t, minDur, maxDur, maxSuggestions), nil
		}
		if len(suggestions) > 0 {
			return suggestions, nil
		}
	}
	return ai.SuggestClips(t, minDur, maxDur, maxSuggestions), nil
}
