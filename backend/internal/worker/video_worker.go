package worker

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"reelcut/internal/notifier"
	"reelcut/internal/queue"
	"reelcut/internal/repository"
	"reelcut/internal/service"
	"reelcut/internal/video"

	"github.com/hibiken/asynq"
)

type VideoWorker struct {
	videoRepo repository.VideoRepository
	jobRepo   repository.ProcessingJobRepository
	storage   *service.StorageService
	notifier  notifier.JobNotifier
}

func NewVideoWorker(videoRepo repository.VideoRepository, jobRepo repository.ProcessingJobRepository, storage *service.StorageService, jobNotifier notifier.JobNotifier) *VideoWorker {
	return &VideoWorker{videoRepo: videoRepo, jobRepo: jobRepo, storage: storage, notifier: jobNotifier}
}

func (w *VideoWorker) Register(mux *asynq.ServeMux) {
	mux.Handle(queue.TypeVideoMetadata, asynq.HandlerFunc(w.HandleMetadata))
	mux.Handle(queue.TypeVideoThumbnail, asynq.HandlerFunc(w.HandleThumbnail))
}

func (w *VideoWorker) HandleMetadata(ctx context.Context, t *asynq.Task) error {
	payload, err := queue.ParseVideoMetadataPayload(t.Payload())
	if err != nil {
		return err
	}
	v, err := w.videoRepo.GetByID(ctx, payload.VideoID)
	if err != nil || v == nil {
		return fmt.Errorf("video not found: %s", payload.VideoID)
	}
	// Download to temp file
	tmpDir := os.TempDir()
	localPath := filepath.Join(tmpDir, "reelcut", v.ID.String()+".mp4")
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return err
	}
	rc, err := w.storage.Download(ctx, v.StoragePath)
	if err != nil {
		return fmt.Errorf("download video: %w", err)
	}
	f, err := os.Create(localPath)
	if err != nil {
		rc.Close()
		return err
	}
	_, err = io.Copy(f, rc)
	rc.Close()
	f.Close()
	if err != nil {
		os.Remove(localPath)
		return err
	}
	defer os.Remove(localPath)

	meta, err := video.GetMetadata(ctx, localPath)
	if err != nil {
		return fmt.Errorf("get metadata: %w", err)
	}
	v.DurationSeconds = &meta.DurationSeconds
	v.Width = &meta.Width
	v.Height = &meta.Height
	v.FPS = &meta.FPS
	v.Codec = &meta.Codec
	v.Bitrate = &meta.Bitrate
	v.FileSizeBytes = &meta.FileSizeBytes
	v.Status = "ready"
	if err := w.videoRepo.Update(ctx, v); err != nil {
		return err
	}
	if job, _ := w.jobRepo.GetByEntity(ctx, "video", payload.VideoID); job != nil {
		job.Progress = 50
		_ = w.jobRepo.Update(ctx, job)
		if w.notifier != nil {
			w.notifier.NotifyJob(ctx, job)
		}
	}
	return nil
}

func (w *VideoWorker) HandleThumbnail(ctx context.Context, t *asynq.Task) error {
	payload, err := queue.ParseVideoThumbnailPayload(t.Payload())
	if err != nil {
		return err
	}
	v, err := w.videoRepo.GetByID(ctx, payload.VideoID)
	if err != nil || v == nil {
		return fmt.Errorf("video not found: %s", payload.VideoID)
	}
	tmpDir := os.TempDir()
	localPath := filepath.Join(tmpDir, "reelcut", v.ID.String()+".mp4")
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return err
	}
	rc, err := w.storage.Download(ctx, v.StoragePath)
	if err != nil {
		return fmt.Errorf("download video: %w", err)
	}
	f, err := os.Create(localPath)
	if err != nil {
		rc.Close()
		return err
	}
	_, err = io.Copy(f, rc)
	rc.Close()
	f.Close()
	if err != nil {
		os.Remove(localPath)
		return err
	}
	defer os.Remove(localPath)

	thumbPath := filepath.Join(tmpDir, "reelcut", v.ID.String()+"_thumb.jpg")
	// Use 1s offset to skip black slate at start when possible
	timestampSec := 0.0
	if v.DurationSeconds != nil && *v.DurationSeconds > 1 {
		timestampSec = 1.0
	}
	if err := video.ExtractFrame(ctx, localPath, timestampSec, thumbPath); err != nil {
		return fmt.Errorf("extract frame: %w", err)
	}
	defer os.Remove(thumbPath)

	thumbKey := filepath.Join("thumbnails", v.ID.String()+".jpg")
	thumbFile, err := os.Open(thumbPath)
	if err != nil {
		return err
	}
	defer thumbFile.Close()
	if err := w.storage.Upload(ctx, thumbKey, thumbFile, "image/jpeg"); err != nil {
		return fmt.Errorf("upload thumbnail: %w", err)
	}
	v.ThumbnailURL = &thumbKey
	if err := w.videoRepo.Update(ctx, v); err != nil {
		return err
	}
	if job, _ := w.jobRepo.GetByEntity(ctx, "video", payload.VideoID); job != nil {
		job.Progress = 100
		job.Status = "completed"
		now := time.Now()
		job.CompletedAt = &now
		_ = w.jobRepo.Update(ctx, job)
		if w.notifier != nil {
			w.notifier.NotifyJob(ctx, job)
		}
	}
	return nil
}
