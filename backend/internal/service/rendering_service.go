package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"reelcut/internal/domain"
	"reelcut/internal/repository"
	"reelcut/internal/video"
)

type RenderingService struct {
	clipRepo         repository.ClipRepository
	clipStyleRepo    repository.ClipStyleRepository
	videoRepo        repository.VideoRepository
	transcriptionSvc *TranscriptionService
	storage          *StorageService
}

func NewRenderingService(
	clipRepo repository.ClipRepository,
	clipStyleRepo repository.ClipStyleRepository,
	videoRepo repository.VideoRepository,
	transcriptionSvc *TranscriptionService,
	storage *StorageService,
) *RenderingService {
	return &RenderingService{
		clipRepo:         clipRepo,
		clipStyleRepo:     clipStyleRepo,
		videoRepo:         videoRepo,
		transcriptionSvc:  transcriptionSvc,
		storage:           storage,
	}
}

// Render produces the output video for the clip and uploads to storage.
func (s *RenderingService) Render(ctx context.Context, clipID string) error {
	c, err := s.clipRepo.GetByID(ctx, clipID)
	if err != nil || c == nil {
		return domain.ErrNotFound
	}
	style, _ := s.clipStyleRepo.GetByClipID(ctx, clipID)
	v, err := s.videoRepo.GetByID(ctx, c.VideoID.String())
	if err != nil || v == nil {
		return domain.ErrNotFound
	}
	tmpDir := filepath.Join(os.TempDir(), "reelcut", "render", clipID)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	sourcePath := filepath.Join(tmpDir, "source.mp4")
	rc, err := s.storage.Download(ctx, v.StoragePath)
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

	stepPath := filepath.Join(tmpDir, "step1_cut.mp4")
	if err := video.Cut(ctx, sourcePath, stepPath, c.StartTime, c.EndTime); err != nil {
		return err
	}
	current := stepPath

	stepPath = filepath.Join(tmpDir, "step2_crop.mp4")
	if err := video.ResizeCrop(ctx, current, stepPath, c.AspectRatio); err != nil {
		return err
	}
	current = stepPath

	if style != nil && style.CaptionEnabled {
		t, _ := s.transcriptionSvc.GetByVideoID(ctx, c.VideoID.String())
		if t != nil {
			blocks := BlocksFromSegments(t.Segments, style, c.StartTime, c.EndTime)
			srtPath := filepath.Join(tmpDir, "captions.srt")
			if err := os.WriteFile(srtPath, []byte(ToSRT(blocks)), 0644); err != nil {
				return err
			}
			stepPath = filepath.Join(tmpDir, "step3_subs.mp4")
			if err := video.BurnSubtitles(ctx, current, srtPath, stepPath); err != nil {
				return err
			}
			current = stepPath
		}
	}

	outputKey := filepath.Join("renders", clipID, "output.mp4")
	outPath := filepath.Join(tmpDir, "output.mp4")
	if err := copyFile(current, outPath); err != nil {
		return err
	}
	outFile, err := os.Open(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	if err := s.storage.Upload(ctx, outputKey, outFile, "video/mp4"); err != nil {
		return fmt.Errorf("upload render: %w", err)
	}

	c.StoragePath = &outputKey
	c.Status = "ready"
	return s.clipRepo.Update(ctx, c)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
