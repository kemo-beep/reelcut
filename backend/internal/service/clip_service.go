package service

import (
	"context"
	"encoding/json"
	"io"

	"reelcut/internal/domain"
	"reelcut/internal/queue"
	"reelcut/internal/repository"

	"github.com/google/uuid"
)

const (
	MaxBrollFileSizeMB = 500
	MaxBrollDurationSec = 300
)

type ClipService struct {
	clipRepo              repository.ClipRepository
	clipStyleRepo         repository.ClipStyleRepository
	videoRepo             repository.VideoRepository
	transcriptionSvc      *TranscriptionService
	jobRepo               repository.ProcessingJobRepository
	queue                 *queue.QueueClient
	templateRepo          repository.TemplateRepository
	userRepo              repository.UserRepository
	usageLogRepo          repository.UsageLogRepository
	brollAssetRepo        repository.BrollAssetRepository
	clipBrollSegmentRepo  repository.ClipBrollSegmentRepository
	storage               *StorageService
}

func NewClipService(
	clipRepo repository.ClipRepository,
	clipStyleRepo repository.ClipStyleRepository,
	videoRepo repository.VideoRepository,
	transcriptionSvc *TranscriptionService,
	jobRepo repository.ProcessingJobRepository,
	queue *queue.QueueClient,
	templateRepo repository.TemplateRepository,
	userRepo repository.UserRepository,
	usageLogRepo repository.UsageLogRepository,
	brollAssetRepo repository.BrollAssetRepository,
	clipBrollSegmentRepo repository.ClipBrollSegmentRepository,
	storage *StorageService,
) *ClipService {
	return &ClipService{
		clipRepo:             clipRepo,
		clipStyleRepo:        clipStyleRepo,
		videoRepo:            videoRepo,
		transcriptionSvc:     transcriptionSvc,
		jobRepo:              jobRepo,
		queue:                queue,
		templateRepo:         templateRepo,
		userRepo:             userRepo,
		usageLogRepo:         usageLogRepo,
		brollAssetRepo:       brollAssetRepo,
		clipBrollSegmentRepo: clipBrollSegmentRepo,
		storage:              storage,
	}
}

func (s *ClipService) Create(ctx context.Context, userID uuid.UUID, videoID, name string, startTime, endTime float64, aspectRatio string, viralityScore *float64, isAISuggested bool) (*domain.Clip, error) {
	vid, err := uuid.Parse(videoID)
	if err != nil {
		return nil, domain.ErrValidation
	}
	v, err := s.videoRepo.GetByID(ctx, videoID)
	if err != nil || v == nil || v.UserID != userID {
		return nil, domain.ErrNotFound
	}
	dur := endTime - startTime
	if dur <= 0 {
		return nil, domain.ErrValidation
	}
	if aspectRatio == "" {
		aspectRatio = "9:16"
	}
	c := &domain.Clip{
		ID:              uuid.New(),
		VideoID:         vid,
		UserID:          userID,
		Name:            name,
		StartTime:       startTime,
		EndTime:         endTime,
		DurationSeconds: &dur,
		AspectRatio:     aspectRatio,
		ViralityScore:   viralityScore,
		Status:          "draft",
		IsAISuggested:   isAISuggested,
	}
	if err := s.clipRepo.Create(ctx, c); err != nil {
		return nil, err
	}
	style := &domain.ClipStyle{
		ID:                    uuid.New(),
		ClipID:                c.ID,
		CaptionEnabled:        true,
		CaptionFont:           "Inter",
		CaptionSize:           48,
		CaptionColor:          "#FFFFFF",
		CaptionPosition:       "bottom",
		CaptionMaxWords:        3,
		BrandLogoScale:        1,
		BrandWatermarkOpacity: 0.8,
		BackgroundMusicVolume: 0.3,
		OriginalAudioVolume:   1,
	}
	if err := s.clipStyleRepo.Create(ctx, style); err != nil {
		return nil, err
	}
	_ = s.queue.EnqueueClipThumbnail(c.ID)
	return c, nil
}

func (s *ClipService) GetByID(ctx context.Context, clipID, userID string) (*domain.Clip, error) {
	c, err := s.clipRepo.GetByID(ctx, clipID)
	if err != nil || c == nil || c.UserID.String() != userID {
		return nil, domain.ErrNotFound
	}
	style, _ := s.clipStyleRepo.GetByClipID(ctx, clipID)
	c.Style = style
	return c, nil
}

func (s *ClipService) List(ctx context.Context, userID string, videoID *string, status *string, page, perPage int, sortBy, sortOrder string) ([]*domain.Clip, int, error) {
	if perPage <= 0 {
		perPage = 20
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	return s.clipRepo.List(ctx, userID, videoID, status, perPage, offset, sortBy, sortOrder)
}

func (s *ClipService) Update(ctx context.Context, c *domain.Clip) error {
	return s.clipRepo.Update(ctx, c)
}

// EnqueueClipThumbnail enqueues a job to generate the clip's thumbnail. Call when clip is created or when start/end time change.
func (s *ClipService) EnqueueClipThumbnail(ctx context.Context, clipID string) error {
	uid, err := uuid.Parse(clipID)
	if err != nil {
		return err
	}
	return s.queue.EnqueueClipThumbnail(uid)
}

func (s *ClipService) Delete(ctx context.Context, clipID, userID string) error {
	c, err := s.clipRepo.GetByID(ctx, clipID)
	if err != nil || c == nil || c.UserID.String() != userID {
		return domain.ErrNotFound
	}
	return s.clipRepo.Delete(ctx, clipID)
}

func (s *ClipService) Duplicate(ctx context.Context, clipID, userID string) (*domain.Clip, error) {
	c, err := s.GetByID(ctx, clipID, userID)
	if err != nil {
		return nil, err
	}
	newName := c.Name + " (copy)"
	return s.Create(ctx, c.UserID, c.VideoID.String(), newName, c.StartTime, c.EndTime, c.AspectRatio, c.ViralityScore, false)
}

func (s *ClipService) GetStyle(ctx context.Context, clipID, userID string) (*domain.ClipStyle, error) {
	c, err := s.clipRepo.GetByID(ctx, clipID)
	if err != nil || c == nil || c.UserID.String() != userID {
		return nil, domain.ErrNotFound
	}
	return s.clipStyleRepo.GetByClipID(ctx, clipID)
}

func (s *ClipService) UpdateStyle(ctx context.Context, clipID, userID string, updates *domain.ClipStyle) error {
	style, err := s.GetStyle(ctx, clipID, userID)
	if err != nil {
		return err
	}
	if updates.CaptionEnabled != style.CaptionEnabled {
		style.CaptionEnabled = updates.CaptionEnabled
	}
	if updates.CaptionFont != "" {
		style.CaptionFont = updates.CaptionFont
	}
	if updates.CaptionSize > 0 {
		style.CaptionSize = updates.CaptionSize
	}
	if updates.CaptionColor != "" {
		style.CaptionColor = updates.CaptionColor
	}
	if updates.CaptionPosition != "" {
		style.CaptionPosition = updates.CaptionPosition
	}
	if updates.CaptionLanguage != nil {
		style.CaptionLanguage = updates.CaptionLanguage
	}
	if updates.BackgroundMusicVolume >= 0 {
		style.BackgroundMusicVolume = updates.BackgroundMusicVolume
	}
	if updates.OriginalAudioVolume >= 0 {
		style.OriginalAudioVolume = updates.OriginalAudioVolume
	}
	return s.clipStyleRepo.Update(ctx, style)
}

func (s *ClipService) GetCaptionsSRT(ctx context.Context, clipID, userID string) (string, error) {
	c, err := s.GetByID(ctx, clipID, userID)
	if err != nil {
		return "", err
	}
	var t *domain.Transcription
	if c.Style != nil && c.Style.CaptionLanguage != nil && *c.Style.CaptionLanguage != "" {
		t, _ = s.transcriptionSvc.GetByVideoIDAndLanguage(ctx, c.VideoID.String(), *c.Style.CaptionLanguage)
	}
	if t == nil {
		t, err = s.transcriptionSvc.GetByVideoID(ctx, c.VideoID.String())
	}
	if err != nil || t == nil {
		return "", domain.ErrNotFound
	}
	blocks := BlocksFromSegments(t.Segments, c.Style, c.StartTime, c.EndTime)
	return ToSRT(blocks), nil
}

func (s *ClipService) GetCaptionsVTT(ctx context.Context, clipID, userID string) (string, error) {
	c, err := s.GetByID(ctx, clipID, userID)
	if err != nil {
		return "", err
	}
	var t *domain.Transcription
	if c.Style != nil && c.Style.CaptionLanguage != nil && *c.Style.CaptionLanguage != "" {
		t, _ = s.transcriptionSvc.GetByVideoIDAndLanguage(ctx, c.VideoID.String(), *c.Style.CaptionLanguage)
	}
	if t == nil {
		t, err = s.transcriptionSvc.GetByVideoID(ctx, c.VideoID.String())
	}
	if err != nil || t == nil {
		return "", domain.ErrNotFound
	}
	blocks := BlocksFromSegments(t.Segments, c.Style, c.StartTime, c.EndTime)
	return ToVTT(blocks), nil
}

func (s *ClipService) StartRender(ctx context.Context, clipID, userID, preset string) (jobID string, err error) {
	c, err := s.clipRepo.GetByID(ctx, clipID)
	if err != nil || c == nil || c.UserID.String() != userID {
		return "", domain.ErrNotFound
	}
	if err := s.userRepo.DeductCredits(ctx, userID, 1); err != nil {
		return "", domain.ErrInsufficientCredits
	}
	usageLog := &domain.UsageLog{ID: uuid.New(), UserID: c.UserID, Action: "render", CreditsUsed: 1}
	_ = s.usageLogRepo.Create(ctx, usageLog)
	job := &domain.ProcessingJob{
		ID:         uuid.New(),
		UserID:     c.UserID,
		JobType:    "rendering",
		EntityType: "clip",
		EntityID:   c.ID,
		Status:     "pending",
		Progress:   0,
	}
	if err := s.jobRepo.Create(ctx, job); err != nil {
		return "", err
	}
	if err := s.queue.EnqueueRender(c.ID, job.ID, preset); err != nil {
		return "", err
	}
	c.Status = "rendering"
	if err := s.clipRepo.Update(ctx, c); err != nil {
		return "", err
	}
	return job.ID.String(), nil
}

func (s *ClipService) CancelRender(ctx context.Context, clipID, userID string) error {
	c, err := s.clipRepo.GetByID(ctx, clipID)
	if err != nil || c == nil || c.UserID.String() != userID {
		return domain.ErrNotFound
	}
	job, err := s.jobRepo.GetByEntity(ctx, "clip", clipID)
	if err != nil || job == nil || job.UserID != c.UserID {
		return domain.ErrNotFound
	}
	if job.Status != "pending" && job.Status != "processing" {
		return nil
	}
	job.Status = "cancelled"
	return s.jobRepo.Update(ctx, job)
}

func (s *ClipService) ApplyTemplate(ctx context.Context, clipID, userID, templateID string) error {
	if _, err := s.GetByID(ctx, clipID, userID); err != nil {
		return err
	}
	tpl, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil || tpl == nil {
		return domain.ErrNotFound
	}
	style, err := s.clipStyleRepo.GetByClipID(ctx, clipID)
	if err != nil || style == nil {
		return domain.ErrNotFound
	}
	var cfg map[string]interface{}
	if len(tpl.StyleConfig) > 0 {
		_ = json.Unmarshal(tpl.StyleConfig, &cfg)
	}
	if cfg != nil {
		if v, ok := cfg["caption_font"].(string); ok && v != "" {
			style.CaptionFont = v
		}
		if v, ok := cfg["caption_size"].(float64); ok {
			style.CaptionSize = int(v)
		}
		if v, ok := cfg["caption_color"].(string); ok && v != "" {
			style.CaptionColor = v
		}
		if v, ok := cfg["caption_position"].(string); ok && v != "" {
			style.CaptionPosition = v
		}
	}
	if err := s.clipStyleRepo.Update(ctx, style); err != nil {
		return err
	}
	return s.templateRepo.IncrementUsageCount(ctx, templateID)
}

func (s *ClipService) CreateBrollAsset(ctx context.Context, userID uuid.UUID, projectID *string, filename string, body io.Reader, contentType string, sizeBytes int64) (*domain.BrollAsset, error) {
	maxBytes := int64(MaxBrollFileSizeMB * 1024 * 1024)
	if sizeBytes > maxBytes || sizeBytes <= 0 {
		return nil, domain.ErrValidation
	}
	id := uuid.New()
	storagePath := "broll/" + id.String() + ".mp4"
	if err := s.storage.Upload(ctx, storagePath, body, contentType); err != nil {
		return nil, err
	}
	var projID *uuid.UUID
	if projectID != nil {
		p, _ := uuid.Parse(*projectID)
		projID = &p
	}
	a := &domain.BrollAsset{
		ID:               id,
		UserID:           userID,
		ProjectID:        projID,
		OriginalFilename: filename,
		StoragePath:      storagePath,
	}
	if err := s.brollAssetRepo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *ClipService) ListBrollAssets(ctx context.Context, userID string, projectID *string, limit, offset int) ([]*domain.BrollAsset, int, error) {
	return s.brollAssetRepo.ListByUserID(ctx, userID, projectID, limit, offset)
}

func (s *ClipService) AddBrollSegment(ctx context.Context, clipID, userID, brollAssetID string, startTime, endTime float64, position string, scale, opacity float64) (*domain.ClipBrollSegment, error) {
	c, err := s.clipRepo.GetByID(ctx, clipID)
	if err != nil || c == nil || c.UserID.String() != userID {
		return nil, domain.ErrNotFound
	}
	asset, err := s.brollAssetRepo.GetByID(ctx, brollAssetID)
	if err != nil || asset == nil || asset.UserID.String() != userID {
		return nil, domain.ErrNotFound
	}
	if endTime <= startTime {
		return nil, domain.ErrValidation
	}
	clipDur := c.EndTime - c.StartTime
	if startTime < 0 || endTime > clipDur {
		return nil, domain.ErrValidation
	}
	if position == "" {
		position = "overlay"
	}
	if position != "overlay" && position != "cut_in" {
		position = "overlay"
	}
	if scale <= 0 {
		scale = 0.5
	}
	if opacity <= 0 {
		opacity = 1
	}
	order, _ := s.clipBrollSegmentRepo.NextSequenceOrder(ctx, clipID)
	assetUUID, _ := uuid.Parse(brollAssetID)
	seg := &domain.ClipBrollSegment{
		ID:            uuid.New(),
		ClipID:        c.ID,
		BrollAssetID:  assetUUID,
		StartTime:     startTime,
		EndTime:       endTime,
		Position:      position,
		Scale:         scale,
		Opacity:       opacity,
		SequenceOrder: order,
	}
	if err := s.clipBrollSegmentRepo.Create(ctx, seg); err != nil {
		return nil, err
	}
	seg.Asset = asset
	return seg, nil
}

func (s *ClipService) ListBrollSegments(ctx context.Context, clipID, userID string) ([]*domain.ClipBrollSegment, error) {
	if _, err := s.GetByID(ctx, clipID, userID); err != nil {
		return nil, err
	}
	return s.clipBrollSegmentRepo.GetByClipID(ctx, clipID)
}

func (s *ClipService) DeleteBrollSegment(ctx context.Context, segmentID, userID string) error {
	seg, err := s.clipBrollSegmentRepo.GetByID(ctx, segmentID)
	if err != nil || seg == nil {
		return domain.ErrNotFound
	}
	c, err := s.clipRepo.GetByID(ctx, seg.ClipID.String())
	if err != nil || c == nil || c.UserID.String() != userID {
		return domain.ErrNotFound
	}
	return s.clipBrollSegmentRepo.Delete(ctx, segmentID)
}
