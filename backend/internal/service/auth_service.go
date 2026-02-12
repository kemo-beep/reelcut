package service

import (
	"context"
	"errors"
	"time"

	"reelcut/internal/domain"
	"reelcut/internal/repository"
	"reelcut/internal/utils"

	"github.com/google/uuid"
)

var (
	ErrEmailExists    = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.UserSessionRepository
	jwtSecret   string
	jwtRefresh  string
	accessExp   time.Duration
	refreshExp  time.Duration
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.UserSessionRepository, jwtSecret, jwtRefresh string, accessExp, refreshExp time.Duration) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
		jwtRefresh:  jwtRefresh,
		accessExp:   accessExp,
		refreshExp:  refreshExp,
	}
}

type RegisterInput struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	FullName *string `json:"full_name,omitempty"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	User  *domain.User  `json:"user"`
	Token *utils.TokenPair `json:"token"`
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*TokenResponse, error) {
	if !utils.ValidateEmail(in.Email) {
		return nil, &domain.ValidationError{Field: "email", Message: "Invalid email format"}
	}
	ok, msg := utils.ValidatePassword(in.Password)
	if !ok {
		return nil, &domain.ValidationError{Field: "password", Message: msg}
	}
	existing, _ := s.userRepo.GetByEmail(ctx, in.Email)
	if existing != nil {
		return nil, ErrEmailExists
	}
	hash, err := utils.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		ID:               uuid.New(),
		Email:            in.Email,
		PasswordHash:     hash,
		FullName:         in.FullName,
		SubscriptionTier: "free",
		CreditsRemaining: 0,
		EmailVerified:    false,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return s.issueTokenPair(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (*TokenResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, in.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}
	if !utils.ComparePassword(user.PasswordHash, in.Password) {
		return nil, ErrInvalidCredentials
	}
	return s.issueTokenPair(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	userID, err := utils.ParseRefreshToken(refreshToken, s.jwtRefresh)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, domain.ErrUnauthorized
	}
	return s.issueTokenPair(ctx, user)
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	_, _ = s.userRepo.GetByEmail(ctx, email)
	// Don't reveal whether email exists; in production send email with reset link
	return nil
}

type ResetPasswordInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (s *AuthService) ResetPassword(ctx context.Context, in ResetPasswordInput) error {
	// Token would be validated (e.g. from email); for MVP we could accept a short-lived JWT
	if ok, msg := utils.ValidatePassword(in.NewPassword); !ok {
		return &domain.ValidationError{Field: "new_password", Message: msg}
	}
	// Stub: in production parse token and get user ID, then update password
	_ = in.Token
	return nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	// Stub: in production validate token and set user email_verified = true
	_ = token
	return nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return ErrInvalidCredentials
	}
	if !utils.ComparePassword(user.PasswordHash, oldPassword) {
		return ErrInvalidCredentials
	}
	ok, msg := utils.ValidatePassword(newPassword)
	if !ok {
		return &domain.ValidationError{Field: "new_password", Message: msg}
	}
	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.userRepo.UpdatePassword(ctx, userID, hash)
}

func (s *AuthService) issueTokenPair(ctx context.Context, user *domain.User) (*TokenResponse, error) {
	access, expAt, err := utils.IssueAccessToken(user.ID.String(), user.Email, s.jwtSecret, s.accessExp)
	if err != nil {
		return nil, err
	}
	refresh, _, err := utils.IssueRefreshToken(user.ID.String(), s.jwtRefresh, s.refreshExp)
	if err != nil {
		return nil, err
	}
	expiresIn := int(s.accessExp.Seconds())
	return &TokenResponse{
		User: user,
		Token: &utils.TokenPair{
			AccessToken:  access,
			RefreshToken: refresh,
			ExpiresIn:    expiresIn,
			ExpiresAt:    expAt,
		},
	}, nil
}
