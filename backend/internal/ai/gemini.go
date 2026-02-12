package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"

	"reelcut/internal/domain"
)

const (
	geminiAPIBase = "https://generativelanguage.googleapis.com/v1beta"
	geminiDefaultModel = "gemini-2.0-flash"
)

// geminiGenerateRequest and geminiGenerateResponse match the Gemini REST API.
type geminiGenerateRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type geminiGenerateResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// geminiClipItem is the shape we ask Gemini to return (and parse).
type geminiClipItem struct {
	StartTime     float64 `json:"start_time"`
	EndTime      float64 `json:"end_time"`
	Reason       string  `json:"reason"`
	ViralityScore float64 `json:"virality_score"`
	Transcript   string  `json:"transcript"`
}

var jsonBlockRe = regexp.MustCompile("(?s)```(?:json)?\\s*([\\s\\S]*?)```")

// SuggestClipsWithGemini calls the Gemini API to suggest viral clips from the transcript.
// Returns nil, nil if GEMINI_API_KEY is not set. Returns error on API or parse failure.
func SuggestClipsWithGemini(ctx context.Context, transcription *domain.Transcription, minDur, maxDur float64, maxSuggestions int) ([]ClipSuggestion, error) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		return nil, nil
	}
	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = geminiDefaultModel
	}

	transcriptText := buildTranscriptWithTimestamps(transcription)
	if transcriptText == "" {
		return nil, nil
	}

	prompt := buildClipSuggestionPrompt(transcriptText, minDur, maxDur, maxSuggestions)

	reqBody := geminiGenerateRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{Parts: []struct {
				Text string `json:"text"`
			}{{Text: prompt}}},
		},
	}
	body, _ := json.Marshal(reqBody)
	base := os.Getenv("GEMINI_API_BASE")
	if base == "" {
		base = geminiAPIBase
	}
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", base, model, key)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API: %s", resp.Status)
	}

	var apiResp geminiGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	if len(apiResp.Candidates) == 0 || len(apiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini: empty response")
	}

	rawText := apiResp.Candidates[0].Content.Parts[0].Text
	parsed := parseGeminiClipsJSON(rawText)
	if len(parsed) == 0 {
		return nil, fmt.Errorf("gemini: no clips parsed from response")
	}

	out := make([]ClipSuggestion, 0, len(parsed))
	for _, p := range parsed {
		if p.EndTime <= p.StartTime {
			continue
		}
		dur := p.EndTime - p.StartTime
		if dur < minDur || dur > maxDur {
			continue
		}
		if maxSuggestions > 0 && len(out) >= maxSuggestions {
			break
		}
		out = append(out, ClipSuggestion{
			StartTime:     p.StartTime,
			EndTime:       p.EndTime,
			Duration:      dur,
			ViralityScore: p.ViralityScore,
			Reason:        p.Reason,
			Transcript:    p.Transcript,
			Highlights:    []string{p.Transcript},
		})
	}
	return out, nil
}

func buildTranscriptWithTimestamps(t *domain.Transcription) string {
	if t == nil || len(t.Segments) == 0 {
		return ""
	}
	var b strings.Builder
	for _, seg := range t.Segments {
		fmt.Fprintf(&b, "[%.1f - %.1f] %s\n", seg.StartTime, seg.EndTime, seg.Text)
	}
	return b.String()
}

func buildClipSuggestionPrompt(transcriptText string, minDur, maxDur float64, maxSuggestions int) string {
	return fmt.Sprintf(`You are an expert at identifying viral short-form video clips for TikTok, Instagram Reels, and YouTube Shorts.

Below is a transcript of a long-form video with timestamps in seconds in the format [start - end] text.

Analyze the transcript and suggest 5-10 clips that would work well as short viral clips. Requirements:
- Each clip must be between %.0f and %.0f seconds long. Prefer 7-15 seconds for Reels/Shorts/TikTok (max 60 seconds).
- Start and end at natural phrase or sentence boundaries (use the exact start_time and end_time from the transcript).
- Prefer moments with hooks, strong statements, or high engagement potential.
- Return ONLY a valid JSON array of objects. No other text or markdown. Each object must have:
  - "start_time" (number, seconds)
  - "end_time" (number, seconds)
  - "reason" (string, short explanation why this clip is viral-worthy)
  - "virality_score" (number, 0-100)
  - "transcript" (string, the exact text snippet for this clip)

Transcript:
%s`, minDur, maxDur, transcriptText)
}

// parseGeminiClipsJSON extracts a JSON array from the model output (handles markdown code blocks).
func parseGeminiClipsJSON(raw string) []geminiClipItem {
	raw = strings.TrimSpace(raw)
	if sub := jsonBlockRe.FindStringSubmatch(raw); len(sub) >= 2 {
		raw = strings.TrimSpace(sub[1])
	}
	var list []geminiClipItem
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		slog.Warn("gemini: failed to parse clips JSON", "err", err, "preview", truncate(raw, 200))
		return nil
	}
	return list
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
