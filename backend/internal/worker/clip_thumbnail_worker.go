package worker

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"reelcut/internal/queue"
	"reelcut/internal/repository"
	"reelcut/internal/service"
	"reelcut/internal/video"

	"github.com/hibiken/asynq"
)

type ClipThumbnailWorker struct {
	clipRepo  repository.ClipRepository
	videoRepo repository.VideoRepository
	storage   *service.StorageService
}

func NewClipThumbnailWorker(
	clipRepo repository.ClipRepository,
	videoRepo repository.VideoRepository,
	storage *service.StorageService,
) *ClipThumbnailWorker {
	return &ClipThumbnailWorker{clipRepo: clipRepo, videoRepo: videoRepo, storage: storage}
}

func (w *ClipThumbnailWorker) Register(mux *asynq.ServeMux) {
	mux.Handle(queue.TypeClipThumbnail, asynq.HandlerFunc(w.Handle))
}

func (w *ClipThumbnailWorker) Handle(ctx context.Context, t *asynq.Task) error {
	payload, err := queue.ParseClipThumbnailPayload(t.Payload())
	if err != nil {
		return err
	}
	c, err := w.clipRepo.GetByID(ctx, payload.ClipID)
	if err != nil || c == nil {
		return fmt.Errorf("clip not found: %s", payload.ClipID)
	}
	v, err := w.videoRepo.GetByID(ctx, c.VideoID.String())
	if err != nil || v == nil {
		return fmt.Errorf("video not found for clip: %s", payload.ClipID)
	}
	tmpDir := filepath.Join(os.TempDir(), "reelcut", "clip_thumb", payload.ClipID)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	localPath := filepath.Join(tmpDir, "source.mp4")
	rc, err := w.storage.Download(ctx, v.StoragePath)
	if err != nil {
		return fmt.Errorf("download source video: %w", err)
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
		return err
	}

	dur := c.EndTime - c.StartTime
	if dur <= 0 {
		dur = 1
	}
	// Frame at 10% into the clip to avoid black slate at start
	timestampSec := c.StartTime + 0.1*dur
	thumbPath := filepath.Join(tmpDir, "thumb.jpg")
	if err := video.ExtractFrame(ctx, localPath, timestampSec, thumbPath); err != nil {
		return fmt.Errorf("extract frame: %w", err)
	}

	thumbKey := filepath.Join("thumbnails", "clips", c.ID.String()+".jpg")
	thumbFile, err := os.Open(thumbPath)
	if err != nil {
		return err
	}
	defer thumbFile.Close()
	if err := w.storage.Upload(ctx, thumbKey, thumbFile, "image/jpeg"); err != nil {
		return fmt.Errorf("upload clip thumbnail: %w", err)
	}
	c.ThumbnailURL = &thumbKey
	return w.clipRepo.Update(ctx, c)
}
