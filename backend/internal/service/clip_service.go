package service

import (
	"context"
	"encoding/json"

	"reelcut/internal/domain"
	"reelcut/internal/queue"
	"reelcut/internal/repository"

	"github.com/google/uuid"
)

type ClipService struct {
	clipRepo         repository.ClipRepository
	clipStyleRepo    repository.ClipStyleRepository
	videoRepo        repository.VideoRepository
	transcriptionSvc *TranscriptionService
	jobRepo          repository.ProcessingJobRepository
	queue            *queue.QueueClient
	templateRepo     repository.TemplateRepository
	userRepo         repository.UserRepository
	usageLogRepo     repository.UsageLogRepository
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
) *ClipService {
	return &ClipService{
		clipRepo:         clipRepo,
		clipStyleRepo:    clipStyleRepo,
		videoRepo:        videoRepo,
		transcriptionSvc: transcriptionSvc,
		jobRepo:          jobRepo,
		queue:            queue,
		templateRepo:     templateRepo,
		userRepo:         userRepo,
		usageLogRepo:     usageLogRepo,
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
	t, err := s.transcriptionSvc.GetByVideoID(ctx, c.VideoID.String())
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
	t, err := s.transcriptionSvc.GetByVideoID(ctx, c.VideoID.String())
	if err != nil || t == nil {
		return "", domain.ErrNotFound
	}
	blocks := BlocksFromSegments(t.Segments, c.Style, c.StartTime, c.EndTime)
	return ToVTT(blocks), nil
}

func (s *ClipService) StartRender(ctx context.Context, clipID, userID string) (jobID string, err error) {
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
	if err := s.queue.EnqueueRender(c.ID, job.ID); err != nil {
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
