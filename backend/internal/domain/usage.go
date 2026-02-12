package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type UsageLog struct {
	ID                   uuid.UUID       `json:"id"`
	UserID               uuid.UUID       `json:"user_id"`
	Action               string          `json:"action"`
	CreditsUsed          int             `json:"credits_used"`
	VideoDurationSeconds *float64         `json:"video_duration_seconds,omitempty"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
}
