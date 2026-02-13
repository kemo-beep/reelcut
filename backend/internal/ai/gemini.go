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
	suggestLine := "Analyze the transcript and suggest 5-10 clips that would work well as short viral clips."
	if maxSuggestions > 0 {
		suggestLine = fmt.Sprintf("Analyze the transcript and suggest clips that would work well as short viral clips. Suggest at most %d clips.", maxSuggestions)
	}
	return fmt.Sprintf(`You are an expert at identifying viral short-form video clips for TikTok, Instagram Reels, and YouTube Shorts.

Below is a transcript of a long-form video with timestamps in seconds in the format [start - end] text.

%s

Requirements:
- Hook within first ~3 seconds: Each clip should open with a hook or engaging phrase within the first ~3 seconds. Prefer segments that start with a question, bold claim, or attention-grabbing line.
- Complete thought or valuable insight: Each clip must contain one complete idea, tip, or takeaway—not a fragment. Do not start or end clips mid-sentence or mid-thought.
- Natural start and end points: Boundaries must align to natural break points in the transcript (sentence boundaries, topic shifts, or clear pauses). Use only exact start_time and end_time from the transcript; no mid-word or mid-sentence cuts.
- Optimal duration: Target 7-15 seconds per clip for maximum engagement. Each clip must be between %.0f and %.0f seconds (maximum 60 seconds).
- Variety across clips: When suggesting multiple clips from the same video, favor variety in topics or angles so the set feels diverse and non-repetitive.

Return ONLY a valid JSON array of objects. No other text or markdown. Each object must have:
  - "start_time" (number, seconds)
  - "end_time" (number, seconds)
  - "reason" (string, short explanation why this clip is viral-worthy)
  - "virality_score" (number, 0-100)
  - "transcript" (string, the exact text snippet for this clip)

Transcript:
%s`, suggestLine, minDur, maxDur, transcriptText)
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

// TranslateSegments translates segment texts to targetLang using Gemini. Returns translated texts in same order, or nil on error.
func TranslateSegments(ctx context.Context, segmentTexts []string, targetLang string) ([]string, error) {
	if len(segmentTexts) == 0 {
		return nil, nil
	}
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}
	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = geminiDefaultModel
	}
	var b strings.Builder
	for i, t := range segmentTexts {
		fmt.Fprintf(&b, "%d. %s\n", i+1, t)
	}
	prompt := fmt.Sprintf(`Translate the following numbered lines to %s. Preserve the exact same number of lines and order. Return ONLY the translated lines in the same numbering format: "1. translation" then "2. translation" etc. No other text.\n\n%s`, targetLang, b.String())
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
	// Parse "1. text" lines
	lines := strings.Split(rawText, "\n")
	out := make([]string, 0, len(segmentTexts))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, ".")
		if idx <= 0 {
			out = append(out, line)
			continue
		}
		rest := strings.TrimSpace(line[idx+1:])
		out = append(out, rest)
	}
	if len(out) < len(segmentTexts) {
		return nil, fmt.Errorf("gemini: got %d lines, need %d", len(out), len(segmentTexts))
	}
	return out[:len(segmentTexts)], nil
}
