package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Video struct {
	ID                uuid.UUID       `json:"id"`
	ProjectID         uuid.UUID       `json:"project_id"`
	UserID            uuid.UUID       `json:"user_id"`
	OriginalFilename  string          `json:"original_filename"`
	StoragePath       string          `json:"storage_path"`
	ThumbnailURL      *string         `json:"thumbnail_url,omitempty"`
	DurationSeconds   *float64        `json:"duration_seconds,omitempty"`
	Width             *int            `json:"width,omitempty"`
	Height            *int            `json:"height,omitempty"`
	FPS               *float64        `json:"fps,omitempty"`
	FileSizeBytes     *int64          `json:"file_size_bytes,omitempty"`
	Codec             *string         `json:"codec,omitempty"`
	Bitrate           *int            `json:"bitrate,omitempty"`
	Status            string          `json:"status"`
	ErrorMessage      *string         `json:"error_message,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         *time.Time      `json:"-"`
}
