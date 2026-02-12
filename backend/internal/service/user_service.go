package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"reelcut/internal/domain"
	"reelcut/internal/repository"
)

const (
	maxAvatarSizeBytes = 5 * 1024 * 1024 // 5MB
)

var allowedAvatarTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type UserService struct {
	userRepo repository.UserRepository
	storage  *StorageService
}

func NewUserService(userRepo repository.UserRepository, storage *StorageService) *UserService {
	return &UserService{userRepo: userRepo, storage: storage}
}

func (s *UserService) UploadAvatar(ctx context.Context, userID string, file io.Reader, contentType string, contentLength int64) (*domain.User, error) {
	ext, ok := allowedAvatarTypes[contentType]
	if !ok {
		return nil, &domain.ValidationError{Field: "file", Message: "avatar must be image/jpeg, image/png, or image/webp"}
	}
	if contentLength > maxAvatarSizeBytes {
		return nil, &domain.ValidationError{Field: "file", Message: "avatar must be at most 5MB"}
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, domain.ErrNotFound
	}
	key := fmt.Sprintf("avatars/%s/avatar%s", userID, ext)
	if err := s.storage.Upload(ctx, key, file, contentType); err != nil {
		return nil, fmt.Errorf("upload avatar: %w", err)
	}
	user.AvatarURL = &key
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetAvatarURL(ctx context.Context, userID string) (string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil || user.AvatarURL == nil || strings.TrimSpace(*user.AvatarURL) == "" {
		return "", domain.ErrNotFound
	}
	return s.storage.GeneratePresignedGet(ctx, *user.AvatarURL, 15*time.Minute)
}
