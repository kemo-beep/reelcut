package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"reelcut/internal/domain"
	"reelcut/internal/email"
	"reelcut/internal/repository"
	"reelcut/internal/utils"

	"github.com/google/uuid"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type AuthService struct {
	userRepo       repository.UserRepository
	sessionRepo    repository.UserSessionRepository
	emailSender    email.Sender
	emailFrom      string
	frontendURL    string
	tokenSecret    string
	tokenExpReset  time.Duration
	tokenExpVerify time.Duration
	jwtSecret      string
	jwtRefresh     string
	accessExp      time.Duration
	refreshExp     time.Duration
}

type AuthServiceOpts struct {
	UserRepo       repository.UserRepository
	SessionRepo    repository.UserSessionRepository
	EmailSender    email.Sender
	EmailFrom      string
	FrontendBaseURL string
	TokenSecret    string
	TokenExpiryReset  time.Duration
	TokenExpiryVerify time.Duration
	JWTSecret      string
	JWTRefresh     string
	AccessExpiry   time.Duration
	RefreshExpiry  time.Duration
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.UserSessionRepository, jwtSecret, jwtRefresh string, accessExp, refreshExp time.Duration) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		emailSender: &email.NoOpSender{},
		jwtSecret:   jwtSecret,
		jwtRefresh:  jwtRefresh,
		accessExp:   accessExp,
		refreshExp:  refreshExp,
	}
}

func NewAuthServiceWithEmail(opts AuthServiceOpts) *AuthService {
	sender := opts.EmailSender
	if sender == nil {
		sender = &email.NoOpSender{}
	}
	return &AuthService{
		userRepo:        opts.UserRepo,
		sessionRepo:     opts.SessionRepo,
		emailSender:     sender,
		emailFrom:       opts.EmailFrom,
		frontendURL:     opts.FrontendBaseURL,
		tokenSecret:     opts.TokenSecret,
		tokenExpReset:   opts.TokenExpiryReset,
		tokenExpVerify:  opts.TokenExpiryVerify,
		jwtSecret:       opts.JWTSecret,
		jwtRefresh:      opts.JWTRefresh,
		accessExp:       opts.AccessExpiry,
		refreshExp:      opts.RefreshExpiry,
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
	// Optionally send verification email (non-blocking)
	if s.tokenSecret != "" {
		go func() {
			tok, _ := utils.IssueEmailVerifyToken(user.ID.String(), user.Email, s.tokenSecret, s.tokenExpVerify)
			resetURL := fmt.Sprintf("%s/verify-email?token=%s", s.frontendURL, tok)
			_ = s.emailSender.Send(user.Email, "Verify your Reelcut email",
				fmt.Sprintf("Welcome! Please verify your email by visiting:\n\n%s\n\nThis link expires in 24 hours.", resetURL))
		}()
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

func (s *AuthService) ForgotPassword(ctx context.Context, addr string) error {
	user, err := s.userRepo.GetByEmail(ctx, addr)
	if err != nil || user == nil {
		return nil // don't reveal whether email exists
	}
	tok, err := utils.IssuePasswordResetToken(user.ID.String(), user.Email, s.tokenSecret, s.tokenExpReset)
	if err != nil {
		return err
	}
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, tok)
	body := fmt.Sprintf("You requested a password reset. Click the link below to set a new password:\n\n%s\n\nThis link expires in 1 hour. If you didn't request this, ignore this email.", resetURL)
	return s.emailSender.Send(addr, "Reset your Reelcut password", body)
}

type ResetPasswordInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (s *AuthService) ResetPassword(ctx context.Context, in ResetPasswordInput) error {
	if ok, msg := utils.ValidatePassword(in.NewPassword); !ok {
		return &domain.ValidationError{Field: "new_password", Message: msg}
	}
	claims, err := utils.ParseOneTimeToken(in.Token, s.tokenSecret, "password_reset")
	if err != nil {
		return ErrInvalidToken
	}
	hash, err := utils.HashPassword(in.NewPassword)
	if err != nil {
		return err
	}
	return s.userRepo.UpdatePassword(ctx, claims.UserID, hash)
}

func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	claims, err := utils.ParseOneTimeToken(token, s.tokenSecret, "email_verify")
	if err != nil {
		return ErrInvalidToken
	}
	return s.userRepo.SetEmailVerified(ctx, claims.UserID, true)
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
