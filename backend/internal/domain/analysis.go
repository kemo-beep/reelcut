package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type VideoAnalysis struct {
	ID               uuid.UUID       `json:"id"`
	VideoID          uuid.UUID       `json:"video_id"`
	ScenesDetected   json.RawMessage `json:"scenes_detected,omitempty"`
	FacesDetected    json.RawMessage `json:"faces_detected,omitempty"`
	Topics           json.RawMessage `json:"topics,omitempty"`
	SentimentAnalysis json.RawMessage `json:"sentiment_analysis,omitempty"`
	EngagementScores json.RawMessage `json:"engagement_scores,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
}
