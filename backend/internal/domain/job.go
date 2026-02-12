package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProcessingJob struct {
	ID           uuid.UUID       `json:"id"`
	UserID       uuid.UUID       `json:"user_id"`
	JobType      string          `json:"job_type"`
	EntityType   string          `json:"entity_type"`
	EntityID     uuid.UUID       `json:"entity_id"`
	Priority     int             `json:"priority"`
	Status       string          `json:"status"`
	Progress     int             `json:"progress"`
	ErrorMessage *string         `json:"error_message,omitempty"`
	RetryCount   int             `json:"retry_count"`
	MaxRetries   int             `json:"max_retries"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	StartedAt    *time.Time      `json:"started_at,omitempty"`
	CompletedAt  *time.Time      `json:"completed_at,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
