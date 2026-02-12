package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"reelcut/internal/domain"

	"github.com/google/uuid"
)

func TestParseGeminiClipsJSON(t *testing.T) {
	tests := []struct {
		name   string
		raw    string
		wantN  int
		first  geminiClipItem
	}{
		{
			name:  "valid json array",
			raw:   `[{"start_time":10,"end_time":45,"reason":"hook","virality_score":85,"transcript":"Hello world"}]`,
			wantN: 1,
			first: geminiClipItem{StartTime: 10, EndTime: 45, Reason: "hook", ViralityScore: 85, Transcript: "Hello world"},
		},
		{
			name:  "markdown code block",
			raw:   "Here are the clips:\n```json\n[{\"start_time\":0.5,\"end_time\":32,\"reason\":\"viral\",\"virality_score\":90,\"transcript\":\"Clip one\"}]\n```",
			wantN: 1,
			first: geminiClipItem{StartTime: 0.5, EndTime: 32, Reason: "viral", ViralityScore: 90, Transcript: "Clip one"},
		},
		{
			name:  "multiple clips",
			raw:   `[{"start_time":0,"end_time":30,"reason":"a","virality_score":70,"transcript":""},{"start_time":60,"end_time":95,"reason":"b","virality_score":80,"transcript":"Second"}]`,
			wantN: 2,
			first: geminiClipItem{StartTime: 0, EndTime: 30, Reason: "a", ViralityScore: 70, Transcript: ""},
		},
		{
			name:  "empty array",
			raw:   `[]`,
			wantN: 0,
		},
		{
			name:  "invalid json",
			raw:   `not json`,
			wantN: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseGeminiClipsJSON(tt.raw)
			if len(got) != tt.wantN {
				t.Errorf("parseGeminiClipsJSON: got %d items, want %d", len(got), tt.wantN)
			}
			if tt.wantN > 0 && len(got) > 0 {
				g := got[0]
				if g.StartTime != tt.first.StartTime || g.EndTime != tt.first.EndTime || g.Reason != tt.first.Reason ||
					g.ViralityScore != tt.first.ViralityScore || g.Transcript != tt.first.Transcript {
					t.Errorf("first item: got %+v, want %+v", g, tt.first)
				}
			}
		})
	}
}

func TestBuildTranscriptWithTimestamps(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		if got := buildTranscriptWithTimestamps(nil); got != "" {
			t.Errorf("want empty, got %q", got)
		}
		if got := buildTranscriptWithTimestamps(&domain.Transcription{Segments: []domain.TranscriptSegment{}}); got != "" {
			t.Errorf("want empty, got %q", got)
		}
	})
	t.Run("with segments", func(t *testing.T) {
		tr := &domain.Transcription{
			Segments: []domain.TranscriptSegment{
				{StartTime: 0, EndTime: 5.2, Text: "Hello"},
				{StartTime: 5.2, EndTime: 12.1, Text: "World"},
			},
		}
		got := buildTranscriptWithTimestamps(tr)
		if got == "" {
			t.Fatal("expected non-empty transcript")
		}
		if !strings.Contains(got, "0.0") || !strings.Contains(got, "5.2") || !strings.Contains(got, "12.1") || !strings.Contains(got, "Hello") || !strings.Contains(got, "World") {
			t.Errorf("transcript should contain timestamps and text: %s", got)
		}
	})
}

// TestSuggestClipsWithGemini_MockServer tests the full path: HTTP request to "Gemini" (mock server),
// response parsed into segments with timestamps, and conversion to ClipSuggestion (min/max duration filtering).
func TestSuggestClipsWithGemini_MockServer(t *testing.T) {
	payload := `[{"start_time":10,"end_time":45,"reason":"hook","virality_score":85,"transcript":"First clip"},
		{"start_time":60,"end_time":95,"reason":"second","virality_score":78,"transcript":"Second clip"}]`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		// Gemini API response shape
		enc := json.NewEncoder(w)
		enc.Encode(struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
		}{
			Candidates: []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			}{
				{Content: struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				}{Parts: []struct {
					Text string `json:"text"`
				}{{Text: payload}}},
			},
			},
		})
	}))
	defer server.Close()

	oldBase := os.Getenv("GEMINI_API_BASE")
	oldKey := os.Getenv("GEMINI_API_KEY")
	defer func() {
		if oldBase != "" {
			os.Setenv("GEMINI_API_BASE", oldBase)
		} else {
			os.Unsetenv("GEMINI_API_BASE")
		}
		if oldKey != "" {
			os.Setenv("GEMINI_API_KEY", oldKey)
		} else {
			os.Unsetenv("GEMINI_API_KEY")
		}
	}()
	os.Setenv("GEMINI_API_BASE", server.URL)
	os.Setenv("GEMINI_API_KEY", "test-key")

	tr := &domain.Transcription{
		ID: uuid.New(),
		Segments: []domain.TranscriptSegment{
			{StartTime: 0, EndTime: 100, Text: "Sample"},
		},
	}
	suggestions, err := SuggestClipsWithGemini(context.Background(), tr, 30, 90, 10)
	if err != nil {
		t.Fatalf("SuggestClipsWithGemini: %v", err)
	}
	if len(suggestions) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(suggestions))
	}
	// First: 10–45 (35s), second: 60–95 (35s); both within 30–90
	if suggestions[0].StartTime != 10 || suggestions[0].EndTime != 45 || suggestions[0].Transcript != "First clip" {
		t.Errorf("first suggestion: got %+v", suggestions[0])
	}
	if suggestions[1].StartTime != 60 || suggestions[1].EndTime != 95 || suggestions[1].Transcript != "Second clip" {
		t.Errorf("second suggestion: got %+v", suggestions[1])
	}
}
