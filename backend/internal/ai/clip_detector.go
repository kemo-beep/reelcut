package ai

import (
	"reelcut/internal/domain"
)

type ClipSuggestion struct {
	StartTime     float64  `json:"start_time"`
	EndTime       float64  `json:"end_time"`
	Duration      float64  `json:"duration"`
	ViralityScore float64  `json:"virality_score"`
	Reason        string   `json:"reason"`
	Highlights    []string `json:"highlights"`
	Transcript    string   `json:"transcript"`
}

func SuggestClips(transcription *domain.Transcription, minDur, maxDur float64, maxSuggestions int) []ClipSuggestion {
	if transcription == nil || len(transcription.Segments) == 0 {
		return nil
	}
	var out []ClipSuggestion
	for i, seg := range transcription.Segments {
		if maxSuggestions > 0 && len(out) >= maxSuggestions {
			break
		}
		dur := seg.EndTime - seg.StartTime
		if dur < minDur || dur > maxDur {
			continue
		}
		out = append(out, ClipSuggestion{
			StartTime:     seg.StartTime,
			EndTime:       seg.EndTime,
			Duration:      dur,
			ViralityScore: 70 + float64(i%30),
			Reason:        "Segment from transcript",
			Highlights:    []string{seg.Text},
			Transcript:    seg.Text,
		})
	}
	return out
}
