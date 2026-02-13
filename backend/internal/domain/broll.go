package domain

import (
	"time"

	"github.com/google/uuid"
)

// BrollAsset is an uploaded B-roll file (reusable across clips).
type BrollAsset struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"user_id"`
	ProjectID        *uuid.UUID `json:"project_id,omitempty"`
	OriginalFilename string     `json:"original_filename"`
	StoragePath      string     `json:"storage_path"`
	DurationSeconds  *float64   `json:"duration_seconds,omitempty"`
	Width            *int       `json:"width,omitempty"`
	Height           *int       `json:"height,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ClipBrollSegment places a B-roll asset on the clip timeline.
// start_time/end_time are in clip-relative seconds (0 = start of clip).
// position: "cut_in" = replace A-roll for that span; "overlay" = picture-in-picture.
type ClipBrollSegment struct {
	ID            uuid.UUID  `json:"id"`
	ClipID        uuid.UUID  `json:"clip_id"`
	BrollAssetID  uuid.UUID  `json:"broll_asset_id"`
	StartTime     float64    `json:"start_time"`
	EndTime       float64    `json:"end_time"`
	Position      string     `json:"position"` // "cut_in" | "overlay"
	Scale         float64    `json:"scale"`
	Opacity       float64    `json:"opacity"`
	SequenceOrder int        `json:"sequence_order"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	// Joined for API responses
	Asset *BrollAsset `json:"asset,omitempty"`
}
