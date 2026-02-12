package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Template struct {
	ID          uuid.UUID       `json:"id"`
	UserID      *uuid.UUID      `json:"user_id,omitempty"`
	Name        string          `json:"name"`
	Category    *string         `json:"category,omitempty"`
	IsPublic    bool            `json:"is_public"`
	PreviewURL  *string         `json:"preview_url,omitempty"`
	StyleConfig json.RawMessage `json:"style_config"`
	UsageCount  int             `json:"usage_count"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
