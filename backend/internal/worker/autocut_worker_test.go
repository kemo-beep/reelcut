package worker

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"reelcut/internal/ai"
	"reelcut/internal/domain"
	"reelcut/internal/queue"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// mockVideoRepo returns a fixed video for the target ID.
type mockVideoRepo struct {
	video *domain.Video
	err   error
}

func (m *mockVideoRepo) Create(ctx context.Context, v *domain.Video) error { return nil }
func (m *mockVideoRepo) GetByID(ctx context.Context, id string) (*domain.Video, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.video, nil
}
func (m *mockVideoRepo) List(ctx context.Context, userID string, projectID *string, status *string, limit, offset int, sortBy, sortOrder string) ([]*domain.Video, int, error) {
	return nil, 0, nil
}
func (m *mockVideoRepo) Update(ctx context.Context, v *domain.Video) error { return nil }
func (m *mockVideoRepo) Delete(ctx context.Context, id string) error       { return nil }

// mockClipRepo returns no clips (total 0) so auto-cut proceeds.
type mockClipRepo struct {
	total int
	err   error
}

func (m *mockClipRepo) Create(ctx context.Context, c *domain.Clip) error { return nil }
func (m *mockClipRepo) GetByID(ctx context.Context, id string) (*domain.Clip, error) {
	return nil, nil
}
func (m *mockClipRepo) List(ctx context.Context, userID string, videoID *string, status *string, limit, offset int, sortBy, sortOrder string) ([]*domain.Clip, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return nil, m.total, nil
}
func (m *mockClipRepo) Update(ctx context.Context, c *domain.Clip) error { return nil }
func (m *mockClipRepo) Delete(ctx context.Context, id string) error       { return nil }

// mockAnalysisSvc returns fixed suggestions.
type mockAnalysisSvc struct {
	suggestions []ai.ClipSuggestion
	err         error
}

func (m *mockAnalysisSvc) SuggestClips(ctx context.Context, videoID string, minDur, maxDur float64, maxSuggestions int) ([]ai.ClipSuggestion, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.suggestions, nil
}

// createCall records one Create call for assertions.
type createCall struct {
	UserID         uuid.UUID
	VideoID        string
	Name           string
	StartTime      float64
	EndTime        float64
	AspectRatio    string
	ViralityScore  *float64
	IsAISuggested  bool
}

type mockClipCreator struct {
	calls []createCall
	err   error
}

func (m *mockClipCreator) Create(ctx context.Context, userID uuid.UUID, videoID, name string, startTime, endTime float64, aspectRatio string, viralityScore *float64, isAISuggested bool) (*domain.Clip, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.calls = append(m.calls, createCall{
		UserID:        userID,
		VideoID:       videoID,
		Name:          name,
		StartTime:     startTime,
		EndTime:       endTime,
		AspectRatio:   aspectRatio,
		ViralityScore: viralityScore,
		IsAISuggested: isAISuggested,
	})
	return &domain.Clip{ID: uuid.New()}, nil
}

type mockTranscriptionRepo struct {
	t *domain.Transcription
}

func (m *mockTranscriptionRepo) Create(ctx context.Context, t *domain.Transcription) error { return nil }
func (m *mockTranscriptionRepo) GetByID(ctx context.Context, id string) (*domain.Transcription, error) {
	return m.t, nil
}
func (m *mockTranscriptionRepo) GetByVideoID(ctx context.Context, videoID string) (*domain.Transcription, error) {
	return m.t, nil
}
func (m *mockTranscriptionRepo) Update(ctx context.Context, t *domain.Transcription) error { return nil }
func (m *mockTranscriptionRepo) CreateWithSegments(ctx context.Context, t *domain.Transcription, segments []*domain.TranscriptSegment) error {
	return nil
}

type mockSegmentRepo struct {
	segments []*domain.TranscriptSegment
}

func (m *mockSegmentRepo) GetByTranscriptionID(ctx context.Context, transcriptionID string) ([]*domain.TranscriptSegment, error) {
	return m.segments, nil
}
func (m *mockSegmentRepo) Update(ctx context.Context, s *domain.TranscriptSegment) error { return nil }
func (m *mockSegmentRepo) CreateBatch(ctx context.Context, segments []*domain.TranscriptSegment) error {
	return nil
}

func TestAutoCutWorker_Handle(t *testing.T) {
	vid := uuid.New()
	uid := uuid.New()
	video := &domain.Video{ID: vid, UserID: uid}

	t.Run("creates clips from suggestions", func(t *testing.T) {
		analysisSvc := &mockAnalysisSvc{
			suggestions: []ai.ClipSuggestion{
				{StartTime: 10, EndTime: 45, ViralityScore: 0.8, Transcript: "First highlight"},
				{StartTime: 60, EndTime: 95, ViralityScore: 0.7, Transcript: ""},
			},
		}
		clipCreator := &mockClipCreator{}
		// Segments from actual transcript: name is derived from real text, not Gemini's Transcript field.
		trRepo := &mockTranscriptionRepo{t: &domain.Transcription{ID: uuid.New()}}
		segRepo := &mockSegmentRepo{
			segments: []*domain.TranscriptSegment{
				{StartTime: 10, EndTime: 45, Text: "First highlight"},
			},
		}
		w := NewAutoCutWorker(
			&mockVideoRepo{video: video},
			&mockClipRepo{total: 0},
			analysisSvc,
			clipCreator,
			nil, // no storage in test: worker only creates clip records
			trRepo,
			segRepo,
		)
		payload, _ := json.Marshal(queue.AutoCutPayload{VideoID: vid.String()})
		task := asynq.NewTask(queue.TypeAutoCut, payload)

		err := w.Handle(context.Background(), task)
		if err != nil {
			t.Fatalf("Handle: %v", err)
		}
		if n := len(clipCreator.calls); n != 2 {
			t.Fatalf("expected 2 Create calls, got %d", n)
		}
		if clipCreator.calls[0].Name != "First highlight" {
			t.Errorf("first clip name (from transcript slice): got %q", clipCreator.calls[0].Name)
		}
		if clipCreator.calls[1].Name != "Clip 2" {
			t.Errorf("second clip name (no overlapping segment): got %q", clipCreator.calls[1].Name)
		}
		if clipCreator.calls[0].AspectRatio != "9:16" || !clipCreator.calls[0].IsAISuggested {
			t.Errorf("first clip: aspect 9:16 and isAISuggested true expected")
		}
		// Segment timestamps from suggestions are used for each clip (RenderingService will use these to cut the main video).
		if clipCreator.calls[0].StartTime != 10 || clipCreator.calls[0].EndTime != 45 {
			t.Errorf("first clip segment times: got start=%.0f end=%.0f, want 10–45", clipCreator.calls[0].StartTime, clipCreator.calls[0].EndTime)
		}
		if clipCreator.calls[1].StartTime != 60 || clipCreator.calls[1].EndTime != 95 {
			t.Errorf("second clip segment times: got start=%.0f end=%.0f, want 60–95", clipCreator.calls[1].StartTime, clipCreator.calls[1].EndTime)
		}
	})

	t.Run("skips when video already has clips", func(t *testing.T) {
		clipCreator := &mockClipCreator{}
		w := NewAutoCutWorker(
			&mockVideoRepo{video: video},
			&mockClipRepo{total: 1},
			&mockAnalysisSvc{suggestions: []ai.ClipSuggestion{{Transcript: "x"}}},
			clipCreator,
			nil,
			&mockTranscriptionRepo{},
			&mockSegmentRepo{},
		)
		payload, _ := json.Marshal(queue.AutoCutPayload{VideoID: vid.String()})
		task := asynq.NewTask(queue.TypeAutoCut, payload)

		err := w.Handle(context.Background(), task)
		if err != nil {
			t.Fatalf("Handle: %v", err)
		}
		if len(clipCreator.calls) != 0 {
			t.Errorf("expected no Create calls when video has clips, got %d", len(clipCreator.calls))
		}
	})

	t.Run("returns error when video not found", func(t *testing.T) {
		w := NewAutoCutWorker(
			&mockVideoRepo{video: nil},
			&mockClipRepo{},
			&mockAnalysisSvc{},
			&mockClipCreator{},
			nil,
			&mockTranscriptionRepo{},
			&mockSegmentRepo{},
		)
		payload, _ := json.Marshal(queue.AutoCutPayload{VideoID: vid.String()})
		task := asynq.NewTask(queue.TypeAutoCut, payload)

		err := w.Handle(context.Background(), task)
		if err == nil {
			t.Fatal("expected error when video not found")
		}
	})

	t.Run("returns error on invalid payload", func(t *testing.T) {
		w := NewAutoCutWorker(
			&mockVideoRepo{video: video},
			&mockClipRepo{},
			&mockAnalysisSvc{},
			&mockClipCreator{},
			nil,
			&mockTranscriptionRepo{},
			&mockSegmentRepo{},
		)
		task := asynq.NewTask(queue.TypeAutoCut, []byte("invalid json"))

		err := w.Handle(context.Background(), task)
		if err == nil {
			t.Fatal("expected error on invalid payload")
		}
	})
}

func TestClipNameFromTranscript(t *testing.T) {
	tests := []struct {
		transcript string
		index      int
		want       string
	}{
		{"Short", 1, "Short"},
		{"", 2, "Clip 2"},
		{"  trimmed  ", 1, "trimmed"},
		{strings.Repeat("x", 50), 1, strings.Repeat("x", 40) + "…"},
	}
	for _, tt := range tests {
		got := clipNameFromTranscript(tt.transcript, tt.index)
		if got != tt.want {
			t.Errorf("clipNameFromTranscript(%q, %d) = %q, want %q", tt.transcript, tt.index, got, tt.want)
		}
	}
}
