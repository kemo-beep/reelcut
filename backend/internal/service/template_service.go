package service

import (
	"context"
	"encoding/json"

	"reelcut/internal/domain"
	"reelcut/internal/repository"

	"github.com/google/uuid"
)

type TemplateService struct {
	templateRepo repository.TemplateRepository
}

func NewTemplateService(templateRepo repository.TemplateRepository) *TemplateService {
	return &TemplateService{templateRepo: templateRepo}
}

func (s *TemplateService) Create(ctx context.Context, userID *uuid.UUID, name, category string, isPublic bool, styleConfig json.RawMessage) (*domain.Template, error) {
	t := &domain.Template{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		Category:    &category,
		IsPublic:    isPublic,
		StyleConfig: styleConfig,
	}
	if err := s.templateRepo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TemplateService) GetByID(ctx context.Context, id string) (*domain.Template, error) {
	return s.templateRepo.GetByID(ctx, id)
}

func (s *TemplateService) List(ctx context.Context, userID *string, publicOnly bool, page, perPage int) ([]*domain.Template, int, error) {
	if perPage <= 0 {
		perPage = 20
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	return s.templateRepo.List(ctx, userID, publicOnly, perPage, offset)
}

func (s *TemplateService) Update(ctx context.Context, t *domain.Template) error {
	return s.templateRepo.Update(ctx, t)
}

func (s *TemplateService) Delete(ctx context.Context, id string) error {
	return s.templateRepo.Delete(ctx, id)
}

func (s *TemplateService) ApplyToClip(ctx context.Context, clipID, templateID, userID string) error {
	// Load template and clip style, copy style_config to clip_style
	t, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil || t == nil {
		return domain.ErrNotFound
	}
	// Clip style update would be done by clip service
	_ = clipID
	_ = userID
	_ = t
	return nil
}
