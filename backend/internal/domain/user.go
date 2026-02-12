package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID  `json:"id"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"`
	FullName          *string    `json:"full_name,omitempty"`
	AvatarURL         *string    `json:"avatar_url,omitempty"`
	SubscriptionTier  string     `json:"subscription_tier"`
	CreditsRemaining  int        `json:"credits_remaining"`
	EmailVerified     bool       `json:"email_verified"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"-"`
}

type UserSession struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
