package repository

import (
	"context"

	"reelcut/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	UpdatePassword(ctx context.Context, userID string, passwordHash string) error
	SetEmailVerified(ctx context.Context, userID string, verified bool) error
	DeductCredits(ctx context.Context, userID string, amount int) error
	Delete(ctx context.Context, id string) error
}

type UserSessionRepository interface {
	Create(ctx context.Context, s *domain.UserSession) error
	GetByToken(ctx context.Context, token string) (*domain.UserSession, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID string) error
}

type ProjectRepository interface {
	Create(ctx context.Context, p *domain.Project) error
	GetByID(ctx context.Context, id string) (*domain.Project, error)
	ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Project, int, error)
	Update(ctx context.Context, p *domain.Project) error
	Delete(ctx context.Context, id string) error
}

type VideoRepository interface {
	Create(ctx context.Context, v *domain.Video) error
	GetByID(ctx context.Context, id string) (*domain.Video, error)
	List(ctx context.Context, userID string, projectID *string, status *string, limit, offset int, sortBy, sortOrder string) ([]*domain.Video, int, error)
	Update(ctx context.Context, v *domain.Video) error
	Delete(ctx context.Context, id string) error
}

type TranscriptionRepository interface {
	Create(ctx context.Context, t *domain.Transcription) error
	GetByID(ctx context.Context, id string) (*domain.Transcription, error)
	GetByVideoID(ctx context.Context, videoID string) (*domain.Transcription, error)
	Update(ctx context.Context, t *domain.Transcription) error
	CreateWithSegments(ctx context.Context, t *domain.Transcription, segments []*domain.TranscriptSegment) error
}

type TranscriptSegmentRepository interface {
	GetByTranscriptionID(ctx context.Context, transcriptionID string) ([]*domain.TranscriptSegment, error)
	Update(ctx context.Context, s *domain.TranscriptSegment) error
	CreateBatch(ctx context.Context, segments []*domain.TranscriptSegment) error
}

type TranscriptWordRepository interface {
	GetBySegmentID(ctx context.Context, segmentID string) ([]*domain.TranscriptWord, error)
	CreateBatch(ctx context.Context, words []*domain.TranscriptWord) error
}

type VideoAnalysisRepository interface {
	GetByVideoID(ctx context.Context, videoID string) (*domain.VideoAnalysis, error)
	Upsert(ctx context.Context, a *domain.VideoAnalysis) error
}

type ClipRepository interface {
	Create(ctx context.Context, c *domain.Clip) error
	GetByID(ctx context.Context, id string) (*domain.Clip, error)
	List(ctx context.Context, userID string, videoID *string, status *string, limit, offset int, sortBy, sortOrder string) ([]*domain.Clip, int, error)
	Update(ctx context.Context, c *domain.Clip) error
	Delete(ctx context.Context, id string) error
}

type ClipStyleRepository interface {
	Create(ctx context.Context, s *domain.ClipStyle) error
	GetByClipID(ctx context.Context, clipID string) (*domain.ClipStyle, error)
	Update(ctx context.Context, s *domain.ClipStyle) error
}

type TemplateRepository interface {
	Create(ctx context.Context, t *domain.Template) error
	GetByID(ctx context.Context, id string) (*domain.Template, error)
	List(ctx context.Context, userID *string, publicOnly bool, limit, offset int) ([]*domain.Template, int, error)
	Update(ctx context.Context, t *domain.Template) error
	Delete(ctx context.Context, id string) error
	IncrementUsageCount(ctx context.Context, id string) error
}

type ProcessingJobRepository interface {
	Create(ctx context.Context, j *domain.ProcessingJob) error
	GetByID(ctx context.Context, id string) (*domain.ProcessingJob, error)
	ListByUserID(ctx context.Context, userID string, status *string, limit, offset int) ([]*domain.ProcessingJob, int, error)
	GetByEntity(ctx context.Context, entityType, entityID string) (*domain.ProcessingJob, error)
	Update(ctx context.Context, j *domain.ProcessingJob) error
}

type UsageLogRepository interface {
	Create(ctx context.Context, u *domain.UsageLog) error
	ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.UsageLog, int, error)
}

type SubscriptionRepository interface {
	Create(ctx context.Context, s *domain.Subscription) error
	GetByUserID(ctx context.Context, userID string) (*domain.Subscription, error)
	GetByStripeID(ctx context.Context, stripeSubscriptionID string) (*domain.Subscription, error)
	Update(ctx context.Context, s *domain.Subscription) error
}
