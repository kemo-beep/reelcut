package worker

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"reelcut/internal/ai"
	"reelcut/internal/domain"
	"reelcut/internal/queue"
	"reelcut/internal/repository"
	"reelcut/internal/service"
	videopkg "reelcut/internal/video"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// analysisSuggestClips and clipCreator allow tests to inject mocks.
type analysisSuggestClips interface {
	SuggestClips(ctx context.Context, videoID string, minDur, maxDur float64, maxSuggestions int) ([]ai.ClipSuggestion, error)
}
type clipCreator interface {
	Create(ctx context.Context, userID uuid.UUID, videoID, name string, startTime, endTime float64, aspectRatio string, viralityScore *float64, isAISuggested bool) (*domain.Clip, error)
}

const (
	autocutMinDur = 7
	autocutMaxDur = 60
	autocutMax    = 20
	clipNameLimit = 40
)

type AutoCutWorker struct {
	videoRepo        repository.VideoRepository
	clipRepo         repository.ClipRepository
	analysisSvc      analysisSuggestClips
	clipSvc          clipCreator
	storage          *service.StorageService
	transcriptionRepo repository.TranscriptionRepository
	segmentRepo      repository.TranscriptSegmentRepository
}

// NewAutoCutWorker builds an AutoCut worker. analysisSvc and clipSvc can be *service.AnalysisService and *service.ClipService or test mocks.
func NewAutoCutWorker(
	videoRepo repository.VideoRepository,
	clipRepo repository.ClipRepository,
	analysisSvc analysisSuggestClips,
	clipSvc clipCreator,
	storage *service.StorageService,
	transcriptionRepo repository.TranscriptionRepository,
	segmentRepo repository.TranscriptSegmentRepository,
) *AutoCutWorker {
	return &AutoCutWorker{
		videoRepo:         videoRepo,
		clipRepo:          clipRepo,
		analysisSvc:       analysisSvc,
		clipSvc:           clipSvc,
		storage:           storage,
		transcriptionRepo: transcriptionRepo,
		segmentRepo:       segmentRepo,
	}
}

func (w *AutoCutWorker) Register(mux *asynq.ServeMux) {
	mux.Handle(queue.TypeAutoCut, asynq.HandlerFunc(w.Handle))
}

func (w *AutoCutWorker) Handle(ctx context.Context, t *asynq.Task) error {
	payload, err := queue.ParseAutoCutPayload(t.Payload())
	if err != nil {
		return err
	}
	videoID := payload.VideoID
	video, err := w.videoRepo.GetByID(ctx, videoID)
	if err != nil || video == nil {
		return fmt.Errorf("video not found: %s", videoID)
	}

	// Skip if video already has clips (idempotent: only auto-cut once)
	videoIDPtr := &videoID
	_, total, err := w.clipRepo.List(ctx, video.UserID.String(), videoIDPtr, nil, 1, 0, "created_at", "asc")
	if err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	suggestions, err := w.analysisSvc.SuggestClips(ctx, videoID, autocutMinDur, autocutMaxDur, autocutMax)
	if err != nil {
		return err
	}
	if len(suggestions) == 0 {
		return nil
	}

	// Load actual transcript segments so clip names use the real text from the video, not Gemini's output.
	var segments []*domain.TranscriptSegment
	if t, err := w.transcriptionRepo.GetByVideoID(ctx, videoID); err == nil && t != nil {
		segments, _ = w.segmentRepo.GetByTranscriptionID(ctx, t.ID.String())
	}

	if w.storage == nil {
		// Test path: only create clip records (no FFmpeg cut or upload)
		for i, s := range suggestions {
			name := clipNameFromTranscript(transcriptSlice(segments, s.StartTime, s.EndTime), i+1)
			_, err := w.clipSvc.Create(ctx, video.UserID, videoID, name, s.StartTime, s.EndTime, "9:16", &s.ViralityScore, true)
			if err != nil {
				return fmt.Errorf("create clip %d: %w", i+1, err)
			}
		}
		return nil
	}

	// Download source video to temp (non-destructive: original in storage is never modified)
	tmpDir := filepath.Join(os.TempDir(), "reelcut", "autocut", videoID)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)
	sourcePath := filepath.Join(tmpDir, "source.mp4")
	rc, err := w.storage.Download(ctx, video.StoragePath)
	if err != nil {
		return fmt.Errorf("download source: %w", err)
	}
	f, err := os.Create(sourcePath)
	if err != nil {
		rc.Close()
		return err
	}
	_, err = io.Copy(f, rc)
	rc.Close()
	f.Close()
	if err != nil {
		return err
	}

	for i, s := range suggestions {
		name := clipNameFromTranscript(transcriptSlice(segments, s.StartTime, s.EndTime), i+1)
		c, err := w.clipSvc.Create(ctx, video.UserID, videoID, name, s.StartTime, s.EndTime, "9:16", &s.ViralityScore, true)
		if err != nil {
			return fmt.Errorf("create clip %d: %w", i+1, err)
		}
		clipPath := filepath.Join(tmpDir, c.ID.String()+".mp4")
		if err := videopkg.Cut(ctx, sourcePath, clipPath, s.StartTime, s.EndTime); err != nil {
			return fmt.Errorf("cut clip %d: %w", i+1, err)
		}
		clipFile, err := os.Open(clipPath)
		if err != nil {
			return fmt.Errorf("open cut file %d: %w", i+1, err)
		}
		storageKey := filepath.Join("clips", c.ID.String(), "cut.mp4")
		if err := w.storage.Upload(ctx, storageKey, clipFile, "video/mp4"); err != nil {
			clipFile.Close()
			return fmt.Errorf("upload clip %d: %w", i+1, err)
		}
		clipFile.Close()
		c.StoragePath = &storageKey
		c.Status = "ready"
		if err := w.clipRepo.Update(ctx, c); err != nil {
			return fmt.Errorf("update clip %d storage_path: %w", i+1, err)
		}
	}
	return nil
}

// transcriptSlice returns the concatenated text of segments that overlap [start, end] (subset of the real video transcript).
func transcriptSlice(segments []*domain.TranscriptSegment, start, end float64) string {
	var parts []string
	for _, seg := range segments {
		if seg.EndTime <= start || seg.StartTime >= end {
			continue
		}
		parts = append(parts, strings.TrimSpace(seg.Text))
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}

func clipNameFromTranscript(transcript string, index int) string {
	trimmed := strings.TrimSpace(transcript)
	if trimmed == "" {
		return fmt.Sprintf("Clip %d", index)
	}
	if len(trimmed) > clipNameLimit {
		return trimmed[:clipNameLimit] + "…"
	}
	return trimmed
}
