package ai

import (
	"encoding/json"
	"regexp"

	"reelcut/internal/domain"
)

// SentimentSegment is a segment with sentiment score (stub: rule-based).
type SentimentSegment struct {
	Start     float64 `json:"start"`
	End       float64 `json:"end"`
	Sentiment string  `json:"sentiment"` // positive, negative, neutral
	Score     float64 `json:"score"`
}

var (
	positiveWords = regexp.MustCompile(`(?i)\b(great|awesome|love|amazing|excellent|happy|best|good|wonderful|fantastic)\b`)
	negativeWords = regexp.MustCompile(`(?i)\b(bad|terrible|hate|worst|awful|sad|poor|wrong|fail|horrible)\b`)
)

// AnalyzeSentiment returns sentiment per segment (simple keyword-based stub).
func AnalyzeSentiment(segments []*domain.TranscriptSegment) []SentimentSegment {
	out := make([]SentimentSegment, 0, len(segments))
	for _, s := range segments {
		if s == nil {
			continue
		}
		text := s.Text
		pos := len(positiveWords.FindAllString(text, -1))
		neg := len(negativeWords.FindAllString(text, -1))
		sentiment := "neutral"
		score := 0.5
		if pos > neg {
			sentiment = "positive"
			score = 0.5 + float64(pos-neg)*0.1
			if score > 1 {
				score = 1
			}
		} else if neg > pos {
			sentiment = "negative"
			score = 0.5 - float64(neg-pos)*0.1
			if score < 0 {
				score = 0
			}
		}
		out = append(out, SentimentSegment{
			Start:     s.StartTime,
			End:       s.EndTime,
			Sentiment: sentiment,
			Score:     score,
		})
	}
	return out
}

// SentimentFromSegments accepts transcript segments.
func SentimentFromSegments(segments []*domain.TranscriptSegment) []SentimentSegment {
	return AnalyzeSentiment(segments)
}

// SentimentToJSON returns JSON for video_analysis.sentiment_analysis.
func SentimentToJSON(segments []SentimentSegment) ([]byte, error) {
	return json.Marshal(segments)
}
