package worker

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"reelcut/internal/ai"
	"reelcut/internal/domain"
	"reelcut/internal/queue"
	"reelcut/internal/repository"
	"reelcut/internal/service"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type AnalysisWorker struct {
	videoAnalysisRepo repository.VideoAnalysisRepository
	videoRepo         repository.VideoRepository
	transcriptionRepo repository.TranscriptionRepository
	segmentRepo       repository.TranscriptSegmentRepository
	storageSvc        *service.StorageService
}

func NewAnalysisWorker(
	videoAnalysisRepo repository.VideoAnalysisRepository,
	videoRepo repository.VideoRepository,
	transcriptionRepo repository.TranscriptionRepository,
	segmentRepo repository.TranscriptSegmentRepository,
	storageSvc *service.StorageService,
) *AnalysisWorker {
	return &AnalysisWorker{
		videoAnalysisRepo: videoAnalysisRepo,
		videoRepo:         videoRepo,
		transcriptionRepo: transcriptionRepo,
		segmentRepo:       segmentRepo,
		storageSvc:        storageSvc,
	}
}

func (w *AnalysisWorker) Register(mux *asynq.ServeMux) {
	mux.Handle(queue.TypeAnalysis, asynq.HandlerFunc(w.Handle))
}

func (w *AnalysisWorker) Handle(ctx context.Context, t *asynq.Task) error {
	payload, err := queue.ParseAnalysisPayload(t.Payload())
	if err != nil {
		return err
	}
	video, err := w.videoRepo.GetByID(ctx, payload.VideoID)
	if err != nil {
		return err
	}

	var scenesJSON json.RawMessage = []byte("[]")
	if video.StoragePath != "" {
		tmpFile, err := os.CreateTemp("", "analysis-video-*"+filepath.Ext(video.StoragePath))
		if err == nil {
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()
			rc, err := w.storageSvc.Download(ctx, video.StoragePath)
			if err == nil {
				_, _ = io.Copy(tmpFile, rc)
				rc.Close()
				_ = tmpFile.Sync()
				scenes, _ := ai.DetectScenes(ctx, tmpFile.Name())
				scenesJSON, _ = ai.ScenesToJSON(scenes)
			}
		}
	}

	var sentimentJSON json.RawMessage = []byte("[]")
	transcription, _ := w.transcriptionRepo.GetByVideoID(ctx, payload.VideoID)
	if transcription != nil {
		segments, _ := w.segmentRepo.GetByTranscriptionID(ctx, transcription.ID.String())
		if len(segments) > 0 {
			sentiment := ai.SentimentFromSegments(segments)
			sentimentJSON, _ = ai.SentimentToJSON(sentiment)
		}
	}

	facesJSON, _ := ai.FacesToJSON(nil)

	a := &domain.VideoAnalysis{
		ID:                uuid.New(),
		VideoID:           uuid.MustParse(payload.VideoID),
		ScenesDetected:    scenesJSON,
		FacesDetected:     facesJSON,
		Topics:            json.RawMessage(`[]`),
		SentimentAnalysis: sentimentJSON,
		EngagementScores:  json.RawMessage(`[]`),
	}
	return w.videoAnalysisRepo.Upsert(ctx, a)
}
