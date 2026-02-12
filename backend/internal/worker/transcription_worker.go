package worker

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"reelcut/internal/ai"
	"reelcut/internal/queue"
	"reelcut/internal/repository"
	"reelcut/internal/video"

	"github.com/hibiken/asynq"
)

type TranscriptionWorker struct {
	transcriptionRepo repository.TranscriptionRepository
	segmentRepo       repository.TranscriptSegmentRepository
	wordRepo          repository.TranscriptWordRepository
	videoRepo         repository.VideoRepository
	storage           StorageDownloader
	whisper           *ai.WhisperClient
}

type StorageDownloader interface {
	Download(ctx context.Context, key string) (io.ReadCloser, error)
}

func NewTranscriptionWorker(
	transcriptionRepo repository.TranscriptionRepository,
	segmentRepo repository.TranscriptSegmentRepository,
	wordRepo repository.TranscriptWordRepository,
	videoRepo repository.VideoRepository,
	storage StorageDownloader,
	whisper *ai.WhisperClient,
) *TranscriptionWorker {
	return &TranscriptionWorker{
		transcriptionRepo: transcriptionRepo,
		segmentRepo:      segmentRepo,
		wordRepo:         wordRepo,
		videoRepo:        videoRepo,
		storage:          storage,
		whisper:          whisper,
	}
}

func (w *TranscriptionWorker) Register(mux *asynq.ServeMux) {
	mux.Handle(queue.TypeTranscription, asynq.HandlerFunc(w.Handle))
}

func (w *TranscriptionWorker) Handle(ctx context.Context, t *asynq.Task) error {
	payload, err := queue.ParseTranscriptionPayload(t.Payload())
	if err != nil {
		return err
	}
	v, err := w.videoRepo.GetByID(ctx, payload.VideoID)
	if err != nil || v == nil {
		return fmt.Errorf("video not found: %s", payload.VideoID)
	}
	tr, err := w.transcriptionRepo.GetByID(ctx, payload.TranscriptionID)
	if err != nil || tr == nil {
		return fmt.Errorf("transcription not found: %s", payload.TranscriptionID)
	}
	tmpDir := os.TempDir()
	videoPath := filepath.Join(tmpDir, "reelcut", v.ID.String()+".mp4")
	os.MkdirAll(filepath.Dir(videoPath), 0755)
	rc, err := w.storage.Download(ctx, v.StoragePath)
	if err != nil {
		w.updateStatus(ctx, payload.TranscriptionID, "failed")
		return err
	}
	f, _ := os.Create(videoPath)
	io.Copy(f, rc)
	rc.Close()
	f.Close()
	defer os.Remove(videoPath)

	audioPath := filepath.Join(tmpDir, "reelcut", v.ID.String()+".wav")
	if err := video.ExtractAudio(ctx, videoPath, audioPath); err != nil {
		w.updateStatus(ctx, payload.TranscriptionID, "failed")
		return err
	}
	defer os.Remove(audioPath)

	result, err := w.whisper.TranscribeFile(ctx, audioPath, tr.Language)
	if err != nil {
		w.updateStatus(ctx, payload.TranscriptionID, "failed")
		return err
	}
	if len(result.Segments) > 0 {
		// TODO: persist segments and words; for now just mark completed
	}
	tr.Status = "completed"
	if err := w.transcriptionRepo.Update(ctx, tr); err != nil {
		return err
	}
	return nil
}

func (w *TranscriptionWorker) updateStatus(ctx context.Context, transcriptionID, status string) {
	tr, _ := w.transcriptionRepo.GetByID(ctx, transcriptionID)
	if tr != nil {
		tr.Status = status
		w.transcriptionRepo.Update(ctx, tr)
	}
}
