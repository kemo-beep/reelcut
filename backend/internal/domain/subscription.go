package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID                     uuid.UUID  `json:"id"`
	UserID                 uuid.UUID  `json:"user_id"`
	Tier                   string     `json:"tier"`
	Status                 string     `json:"status"`
	StripeSubscriptionID   *string    `json:"stripe_subscription_id,omitempty"`
	CurrentPeriodStart      *time.Time `json:"current_period_start,omitempty"`
	CurrentPeriodEnd        *time.Time `json:"current_period_end,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}
