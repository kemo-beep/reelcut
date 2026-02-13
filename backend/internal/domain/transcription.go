package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transcription struct {
	ID                    uuid.UUID  `json:"id"`
	VideoID               uuid.UUID  `json:"video_id"`
	Language              string     `json:"language"`
	Status                string     `json:"status"`
	SourceTranscriptionID *uuid.UUID `json:"source_transcription_id,omitempty"`
	ErrorMessage          *string    `json:"error_message,omitempty"`
	WordCount             *int       `json:"word_count,omitempty"`
	DurationSeconds       *float64   `json:"duration_seconds,omitempty"`
	ConfidenceAvg         *float64   `json:"confidence_avg,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	Segments              []TranscriptSegment `json:"segments,omitempty"`
}

type TranscriptSegment struct {
	ID              uuid.UUID        `json:"id"`
	TranscriptionID  uuid.UUID       `json:"transcription_id"`
	StartTime       float64         `json:"start_time"`
	EndTime         float64         `json:"end_time"`
	Text            string          `json:"text"`
	Confidence      *float64        `json:"confidence,omitempty"`
	SpeakerID       *int            `json:"speaker_id,omitempty"`
	SequenceOrder   int             `json:"sequence_order"`
	Words           []TranscriptWord `json:"words,omitempty"`
}

type TranscriptWord struct {
	ID            uuid.UUID `json:"id"`
	SegmentID     uuid.UUID `json:"segment_id"`
	Word          string    `json:"word"`
	StartTime     float64   `json:"start_time"`
	EndTime       float64   `json:"end_time"`
	Confidence    *float64  `json:"confidence,omitempty"`
	SequenceOrder int       `json:"sequence_order"`
}
