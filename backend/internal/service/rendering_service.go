package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"reelcut/internal/config"
	"reelcut/internal/domain"
	"reelcut/internal/repository"
	"reelcut/internal/video"
)

type RenderingService struct {
	clipRepo              repository.ClipRepository
	clipStyleRepo         repository.ClipStyleRepository
	videoRepo             repository.VideoRepository
	transcriptionSvc       *TranscriptionService
	storage               *StorageService
	clipBrollSegmentRepo  repository.ClipBrollSegmentRepository
	brollAssetRepo        repository.BrollAssetRepository
}

func NewRenderingService(
	clipRepo repository.ClipRepository,
	clipStyleRepo repository.ClipStyleRepository,
	videoRepo repository.VideoRepository,
	transcriptionSvc *TranscriptionService,
	storage *StorageService,
	clipBrollSegmentRepo repository.ClipBrollSegmentRepository,
	brollAssetRepo repository.BrollAssetRepository,
) *RenderingService {
	return &RenderingService{
		clipRepo:             clipRepo,
		clipStyleRepo:        clipStyleRepo,
		videoRepo:            videoRepo,
		transcriptionSvc:     transcriptionSvc,
		storage:              storage,
		clipBrollSegmentRepo: clipBrollSegmentRepo,
		brollAssetRepo:       brollAssetRepo,
	}
}

// Render produces the output video for the clip and uploads to storage.
// If preset is non-empty (e.g. "tiktok", "instagram_feed"), uses preset dimensions and aspect ratio.
func (s *RenderingService) Render(ctx context.Context, clipID, preset string) error {
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

	// B-roll overlay: apply each segment in sequence (each step reads previous output).
	segments, _ := s.clipBrollSegmentRepo.GetByClipID(ctx, clipID)
	clipDur := c.EndTime - c.StartTime
	for i, seg := range segments {
		if seg.Asset == nil {
			asset, _ := s.brollAssetRepo.GetByID(ctx, seg.BrollAssetID.String())
			if asset == nil {
				continue
			}
			seg.Asset = asset
		}
		brollPath := filepath.Join(tmpDir, fmt.Sprintf("broll_%s.mp4", seg.BrollAssetID.String()))
		rc, err := s.storage.Download(ctx, seg.Asset.StoragePath)
		if err != nil {
			continue
		}
		f, _ := os.Create(brollPath)
		io.Copy(f, rc)
		rc.Close()
		f.Close()
		nextPath := filepath.Join(tmpDir, fmt.Sprintf("step1b_broll_%d.mp4", i))
		scale := seg.Scale
		if scale <= 0 {
			scale = 0.5
		}
		opacity := seg.Opacity
		if opacity <= 0 {
			opacity = 1
		}
		if err := video.OverlayVideo(ctx, current, brollPath, nextPath, seg.StartTime, seg.EndTime, scale, opacity); err != nil {
			continue
		}
		current = nextPath
	}
	// If we have segments and clip duration, ensure overlay times are within clip duration (already in clip-relative 0..clipDur)
	_ = clipDur

	stepPath = filepath.Join(tmpDir, "step2_crop.mp4")
	if preset != "" {
		if p := config.GetExportPresetByID(preset); p != nil && p.Width > 0 && p.Height > 0 {
			if err := video.ResizeCropToSize(ctx, current, stepPath, p.Width, p.Height); err != nil {
				return err
			}
		} else {
			if err := video.ResizeCrop(ctx, current, stepPath, c.AspectRatio); err != nil {
				return err
			}
		}
	} else {
		if err := video.ResizeCrop(ctx, current, stepPath, c.AspectRatio); err != nil {
			return err
		}
	}
	current = stepPath

	if style != nil && style.CaptionEnabled {
		var t *domain.Transcription
		if style.CaptionLanguage != nil && *style.CaptionLanguage != "" {
			t, _ = s.transcriptionSvc.GetByVideoIDAndLanguage(ctx, c.VideoID.String(), *style.CaptionLanguage)
		}
		if t == nil {
			t, _ = s.transcriptionSvc.GetByVideoID(ctx, c.VideoID.String())
		}
		if t != nil {
			blocks := BlocksFromSegments(t.Segments, style, c.StartTime, c.EndTime)
			// Use ASS with full styling (font, color, position) when burning captions.
			assPath := filepath.Join(tmpDir, "captions.ass")
			if err := os.WriteFile(assPath, []byte(ToASS(blocks, style)), 0644); err != nil {
				return err
			}
			stepPath = filepath.Join(tmpDir, "step3_subs.mp4")
			if err := video.BurnASS(ctx, current, assPath, stepPath); err != nil {
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
