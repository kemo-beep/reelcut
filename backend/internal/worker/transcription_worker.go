package worker

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"reelcut/internal/ai"
	"reelcut/internal/domain"
	"reelcut/internal/queue"
	"reelcut/internal/repository"
	"reelcut/internal/video"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type TranscriptionWorker struct {
	transcriptionRepo repository.TranscriptionRepository
	segmentRepo       repository.TranscriptSegmentRepository
	wordRepo          repository.TranscriptWordRepository
	videoRepo         repository.VideoRepository
	storage           StorageDownloader
	transcriber       ai.Transcriber
	queue             *queue.QueueClient
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
	transcriber ai.Transcriber,
	queue *queue.QueueClient,
) *TranscriptionWorker {
	return &TranscriptionWorker{
		transcriptionRepo: transcriptionRepo,
		segmentRepo:      segmentRepo,
		wordRepo:         wordRepo,
		videoRepo:        videoRepo,
		storage:          storage,
		transcriber:      transcriber,
		queue:            queue,
	}
}

func (w *TranscriptionWorker) Register(mux *asynq.ServeMux) {
	mux.Handle(queue.TypeTranscription, asynq.HandlerFunc(w.Handle))
}

const transcriptionChunkSeconds = 60.0

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
	tmpDir := filepath.Join(os.TempDir(), "reelcut", v.ID.String(), tr.ID.String())
	os.MkdirAll(tmpDir, 0755)
	videoPath := filepath.Join(tmpDir, "video.mp4")
	rc, err := w.storage.Download(ctx, v.StoragePath)
	if err != nil {
		w.updateStatusWithError(ctx, payload.TranscriptionID, "failed", err.Error())
		return err
	}
	f, _ := os.Create(videoPath)
	io.Copy(f, rc)
	rc.Close()
	f.Close()
	defer os.RemoveAll(tmpDir)

	meta, err := video.GetMetadata(ctx, videoPath)
	if err != nil || meta == nil {
		msg := "video metadata: "
		if err != nil {
			msg += err.Error()
		} else {
			msg += "no metadata"
		}
		w.updateStatusWithError(ctx, payload.TranscriptionID, "failed", msg)
		return fmt.Errorf("video metadata: %w", err)
	}
	duration := meta.DurationSeconds
	if duration <= 0 {
		duration = transcriptionChunkSeconds
	}

	sequenceOrder := 0
	for chunkStart := 0.0; chunkStart < duration; chunkStart += transcriptionChunkSeconds {
		chunkDur := transcriptionChunkSeconds
		if chunkStart+chunkDur > duration {
			chunkDur = duration - chunkStart
		}
		if chunkDur <= 0 {
			break
		}
		chunkPath := filepath.Join(tmpDir, fmt.Sprintf("chunk_%.0f.wav", chunkStart))
		if err := video.ExtractAudioChunkFromVideo(ctx, videoPath, chunkPath, chunkStart, chunkDur); err != nil {
			w.updateStatusWithError(ctx, payload.TranscriptionID, "failed", fmt.Sprintf("extract chunk at %.0fs: %v", chunkStart, err))
			return fmt.Errorf("extract chunk at %.0fs: %w", chunkStart, err)
		}
		result, err := w.transcriber.TranscribeFile(ctx, chunkPath, tr.Language)
		os.Remove(chunkPath)
		if err != nil {
			w.updateStatusWithError(ctx, payload.TranscriptionID, "failed", fmt.Sprintf("transcribe chunk at %.0fs: %v", chunkStart, err))
			return fmt.Errorf("transcribe chunk at %.0fs: %w", chunkStart, err)
		}
		if len(result.Segments) > 0 {
			segments := make([]*domain.TranscriptSegment, 0, len(result.Segments))
			wordsBySegment := make([][]*domain.TranscriptWord, 0, len(result.Segments))
			for i, seg := range result.Segments {
				globalStart := seg.Start + chunkStart
				globalEnd := seg.End + chunkStart
				if globalEnd <= globalStart {
					globalEnd = globalStart + 0.001
				}
				s := &domain.TranscriptSegment{
					ID:             uuid.New(),
					TranscriptionID: tr.ID,
					StartTime:      globalStart,
					EndTime:        globalEnd,
					Text:           seg.Text,
					SequenceOrder:  sequenceOrder + i,
				}
				segments = append(segments, s)
				var segmentWords []*domain.TranscriptWord
				for _, w := range result.Words {
					mid := (w.Start + w.End) / 2
					if mid >= seg.Start && mid <= seg.End {
						segmentWords = append(segmentWords, &domain.TranscriptWord{
							ID:            uuid.New(),
							SegmentID:     s.ID,
							Word:          w.Word,
							StartTime:     w.Start + chunkStart,
							EndTime:       w.End + chunkStart,
							SequenceOrder: len(segmentWords),
						})
					}
				}
				wordsBySegment = append(wordsBySegment, segmentWords)
			}
			if err := w.segmentRepo.CreateBatch(ctx, segments); err != nil {
				w.updateStatusWithError(ctx, payload.TranscriptionID, "failed", "persist segments: "+err.Error())
				return fmt.Errorf("persist segments: %w", err)
			}
			for _, segmentWords := range wordsBySegment {
				if len(segmentWords) > 0 {
					if err := w.wordRepo.CreateBatch(ctx, segmentWords); err != nil {
						w.updateStatusWithError(ctx, payload.TranscriptionID, "failed", "persist words: "+err.Error())
						return fmt.Errorf("persist words: %w", err)
					}
				}
			}
			sequenceOrder += len(segments)
		}
	}

	tr.Status = "completed"
	if err := w.transcriptionRepo.Update(ctx, tr); err != nil {
		return err
	}
	if w.queue != nil {
		if envVal := os.Getenv("AUTO_CUT_AFTER_TRANSCRIPTION"); envVal == "1" || envVal == "true" || envVal == "yes" {
			_ = w.queue.EnqueueAutoCut(v.ID)
		}
	}
	return nil
}

func (w *TranscriptionWorker) updateStatus(ctx context.Context, transcriptionID, status string) {
	w.updateStatusWithError(ctx, transcriptionID, status, "")
}

func (w *TranscriptionWorker) updateStatusWithError(ctx context.Context, transcriptionID, status, errMsg string) {
	tr, _ := w.transcriptionRepo.GetByID(ctx, transcriptionID)
	if tr != nil {
		tr.Status = status
		if errMsg != "" {
			tr.ErrorMessage = &errMsg
		}
		w.transcriptionRepo.Update(ctx, tr)
	}
}
