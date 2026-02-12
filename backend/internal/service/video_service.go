package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"reelcut/internal/domain"
	"reelcut/internal/queue"
	"reelcut/internal/repository"

	"github.com/google/uuid"
)

var allowedVideoExtensions = map[string]bool{
	".mp4":  true,
	".mov":  true,
	".webm": true,
}

const (
	JobTypeVideoProcessing = "video_processing"
)

type VideoService struct {
	videoRepo    repository.VideoRepository
	projectRepo  repository.ProjectRepository
	jobRepo      repository.ProcessingJobRepository
	storage      *StorageService
	queue        *queue.QueueClient
	userRepo     repository.UserRepository
	usageLogRepo repository.UsageLogRepository
}

func NewVideoService(videoRepo repository.VideoRepository, projectRepo repository.ProjectRepository, jobRepo repository.ProcessingJobRepository, storage *StorageService, queue *queue.QueueClient, userRepo repository.UserRepository, usageLogRepo repository.UsageLogRepository) *VideoService {
	return &VideoService{
		videoRepo:    videoRepo,
		projectRepo:  projectRepo,
		jobRepo:      jobRepo,
		storage:      storage,
		queue:        queue,
		userRepo:     userRepo,
		usageLogRepo: usageLogRepo,
	}
}

func (s *VideoService) CreateVideo(ctx context.Context, userID, projectID uuid.UUID, originalFilename, storagePath string) (*domain.Video, error) {
	p, err := s.projectRepo.GetByID(ctx, projectID.String())
	if err != nil || p == nil || p.UserID != userID {
		return nil, domain.ErrNotFound
	}
	v := &domain.Video{
		ID:               uuid.New(),
		ProjectID:        projectID,
		UserID:           userID,
		OriginalFilename: originalFilename,
		StoragePath:      storagePath,
		Status:           "uploading",
	}
	if err := s.videoRepo.Create(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *VideoService) ConfirmUpload(ctx context.Context, videoID uuid.UUID) error {
	v, err := s.videoRepo.GetByID(ctx, videoID.String())
	if err != nil || v == nil {
		return domain.ErrNotFound
	}
	if err := s.userRepo.DeductCredits(ctx, v.UserID.String(), 1); err != nil {
		return domain.ErrInsufficientCredits
	}
	usageLog := &domain.UsageLog{ID: uuid.New(), UserID: v.UserID, Action: "video_upload", CreditsUsed: 1}
	_ = s.usageLogRepo.Create(ctx, usageLog)
	// Mark video as ready immediately after confirming upload so the UI and
	// transcription flow don't get stuck in a perpetual "processing" state
	// if background metadata/thumbnail jobs are slow or fail.
	//
	// Background workers will still enrich the video record with metadata
	// and thumbnails, but the core actions (playback, transcription, clips)
	// can proceed as soon as the upload is confirmed.
	v.Status = "ready"
	if err := s.videoRepo.Update(ctx, v); err != nil {
		return err
	}
	job := &domain.ProcessingJob{
		ID:         uuid.New(),
		UserID:     v.UserID,
		JobType:    JobTypeVideoProcessing,
		EntityType: "video",
		EntityID:   videoID,
		Status:     "processing",
		Progress:   0,
	}
	if err := s.jobRepo.Create(ctx, job); err != nil {
		return err
	}
	if err := s.queue.EnqueueVideoMetadata(videoID); err != nil {
		return err
	}
	if err := s.queue.EnqueueVideoThumbnail(videoID); err != nil {
		return err
	}
	return nil
}

func (s *VideoService) GetByID(ctx context.Context, videoID, userID string) (*domain.Video, error) {
	v, err := s.videoRepo.GetByID(ctx, videoID)
	if err != nil || v == nil {
		return nil, domain.ErrNotFound
	}
	if v.UserID.String() != userID {
		return nil, domain.ErrNotFound
	}
	return v, nil
}

func (s *VideoService) List(ctx context.Context, userID string, projectID *string, status *string, page, perPage int, sortBy, sortOrder string) ([]*domain.Video, int, error) {
	if perPage <= 0 {
		perPage = 20
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	return s.videoRepo.List(ctx, userID, projectID, status, perPage, offset, sortBy, sortOrder)
}

func (s *VideoService) Update(ctx context.Context, v *domain.Video) error {
	return s.videoRepo.Update(ctx, v)
}

func (s *VideoService) Delete(ctx context.Context, videoID, userID string) error {
	v, err := s.videoRepo.GetByID(ctx, videoID)
	if err != nil || v == nil || v.UserID.String() != userID {
		return domain.ErrNotFound
	}
	if err := s.storage.Delete(ctx, v.StoragePath); err != nil {
		// Log but continue to soft-delete in DB
	}
	return s.videoRepo.Delete(ctx, videoID)
}

func (s *VideoService) GetPresignedUploadURL(ctx context.Context, userID, projectID, filename string) (uploadURL string, videoID string, err error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedVideoExtensions[ext] {
		return "", "", &domain.ValidationError{Field: "filename", Message: "allowed formats: .mp4, .mov, .webm"}
	}
	pid, err := uuid.Parse(projectID)
	if err != nil {
		return "", "", domain.ErrValidation
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return "", "", domain.ErrValidation
	}
	key := filepath.Join("videos", uid.String(), uuid.New().String(), filepath.Base(filename))
	v, err := s.CreateVideo(ctx, uid, pid, filepath.Base(filename), key)
	if err != nil {
		return "", "", err
	}
	uploadURL, err = s.storage.GeneratePresignedPut(ctx, key, "video/mp4", 60*time.Minute)
	if err != nil {
		return "", "", fmt.Errorf("presigned put: %w", err)
	}
	return uploadURL, v.ID.String(), nil
}

func (s *VideoService) GetPresignedDownloadURL(ctx context.Context, path string, expirySec int) (string, error) {
	dur := time.Duration(expirySec) * time.Second
	if dur <= 0 {
		dur = 15 * time.Minute
	}
	return s.storage.GeneratePresignedGet(ctx, path, dur)
}

// Resumable upload

type ResumablePart struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
}

func (s *VideoService) InitiateResumableUpload(ctx context.Context, userID, projectID, filename string) (uploadID, videoID string, err error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedVideoExtensions[ext] {
		return "", "", &domain.ValidationError{Field: "filename", Message: "allowed formats: .mp4, .mov, .webm"}
	}
	pid, err := uuid.Parse(projectID)
	if err != nil {
		return "", "", domain.ErrValidation
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return "", "", domain.ErrValidation
	}
	p, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil || p == nil || p.UserID != uid {
		return "", "", domain.ErrNotFound
	}
	vid := uuid.New()
	key := filepath.Join("videos", uid.String(), vid.String(), "video"+ext)
	v := &domain.Video{
		ID:               vid,
		ProjectID:        pid,
		UserID:           uid,
		OriginalFilename: filepath.Base(filename),
		StoragePath:      key,
		Status:           "uploading",
	}
	if err := s.videoRepo.Create(ctx, v); err != nil {
		return "", "", err
	}
	uploadID, err = s.storage.CreateMultipartUpload(ctx, key)
	if err != nil {
		return "", "", err
	}
	v, _ = s.videoRepo.GetByID(ctx, vid.String())
	if v != nil {
		v.Metadata, _ = json.Marshal(map[string]string{"upload_id": uploadID})
		_ = s.videoRepo.Update(ctx, v)
	}
	return uploadID, vid.String(), nil
}

func (s *VideoService) UploadPart(ctx context.Context, videoID, userID, uploadID string, partNumber int, body io.Reader) (etag string, err error) {
	v, err := s.videoRepo.GetByID(ctx, videoID)
	if err != nil || v == nil || v.UserID.String() != userID {
		return "", domain.ErrNotFound
	}
	var meta struct {
		UploadID string `json:"upload_id"`
	}
	_ = json.Unmarshal(v.Metadata, &meta)
	if meta.UploadID != uploadID {
		return "", domain.ErrValidation
	}
	return s.storage.UploadPart(ctx, v.StoragePath, uploadID, partNumber, body)
}

func (s *VideoService) CompleteResumableUpload(ctx context.Context, videoID, userID string, parts []ResumablePart) error {
	v, err := s.videoRepo.GetByID(ctx, videoID)
	if err != nil || v == nil || v.UserID.String() != userID {
		return domain.ErrNotFound
	}
	var meta struct {
		UploadID string `json:"upload_id"`
	}
	_ = json.Unmarshal(v.Metadata, &meta)
	if meta.UploadID == "" {
		return domain.ErrValidation
	}
	var storageParts []CompletedPart
	for _, p := range parts {
		storageParts = append(storageParts, CompletedPart{PartNumber: p.PartNumber, ETag: p.ETag})
	}
	if err := s.storage.CompleteMultipartUpload(ctx, v.StoragePath, meta.UploadID, storageParts); err != nil {
		return err
	}
	vid, _ := uuid.Parse(videoID)
	return s.ConfirmUpload(ctx, vid)
}
