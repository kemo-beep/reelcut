package service

import (
	"fmt"
	"strings"

	"reelcut/internal/domain"
)

// CaptionBlock is one caption line with time range.
type CaptionBlock struct {
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
	Text      string  `json:"text"`
}

const maxCaptionDurationSec = 7.0

// BlocksFromSegments builds caption blocks from transcript segments and style.
// clipStart/clipEnd optionally filter to a time range (use 0,0 for no filter).
// Uses CaptionMaxWords to split long segments; respects max duration per caption.
func BlocksFromSegments(segments []domain.TranscriptSegment, style *domain.ClipStyle, clipStart, clipEnd float64) []CaptionBlock {
	if style == nil {
		style = &domain.ClipStyle{CaptionMaxWords: 3}
	}
	maxWords := style.CaptionMaxWords
	if maxWords <= 0 {
		maxWords = 3
	}
	var blocks []CaptionBlock
	for _, seg := range segments {
		if clipEnd > 0 && (seg.EndTime < clipStart || seg.StartTime > clipEnd) {
			continue
		}
		words := strings.Fields(seg.Text)
		if len(words) == 0 {
			continue
		}
		dur := seg.EndTime - seg.StartTime
		wordDuration := dur / float64(len(words))
		for i := 0; i < len(words); i += maxWords {
			end := i + maxWords
			if end > len(words) {
				end = len(words)
			}
			chunk := words[i:end]
			startTime := seg.StartTime + float64(i)*wordDuration
			endTime := seg.StartTime + float64(end)*wordDuration
			if clipEnd > 0 {
				if endTime < clipStart || startTime > clipEnd {
					continue
				}
				if startTime < clipStart {
					startTime = clipStart
				}
				if endTime > clipEnd {
					endTime = clipEnd
				}
			}
			text := strings.Join(chunk, " ")
			if endTime-startTime > maxCaptionDurationSec && len(chunk) > 1 {
				endTime = startTime + maxCaptionDurationSec
			}
			blocks = append(blocks, CaptionBlock{StartTime: startTime, EndTime: endTime, Text: text})
		}
	}
	return blocks
}

// ToSRT returns SRT-formatted string (sequence, time range, text).
func ToSRT(blocks []CaptionBlock) string {
	var b strings.Builder
	for i, blk := range blocks {
		b.WriteString(fmt.Sprintf("%d\n", i+1))
		b.WriteString(formatSRTTime(blk.StartTime) + " --> " + formatSRTTime(blk.EndTime) + "\n")
		b.WriteString(blk.Text + "\n\n")
	}
	return b.String()
}

// ToVTT returns WebVTT-formatted string.
func ToVTT(blocks []CaptionBlock) string {
	var b strings.Builder
	b.WriteString("WEBVTT\n\n")
	for i, blk := range blocks {
		b.WriteString(fmt.Sprintf("%d\n", i+1))
		b.WriteString(formatVTTTime(blk.StartTime) + " --> " + formatVTTTime(blk.EndTime) + "\n")
		b.WriteString(blk.Text + "\n\n")
	}
	return b.String()
}

func formatSRTTime(sec float64) string {
	h := int(sec / 3600)
	m := int((sec - float64(h*3600)) / 60)
	s := int(sec - float64(h*3600) - float64(m*60))
	ms := int((sec - float64(int(sec))) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d,%03d", h, m, s, ms)
}

func formatVTTTime(sec float64) string {
	h := int(sec / 3600)
	m := int((sec - float64(h*3600)) / 60)
	s := int(sec - float64(h*3600) - float64(m*60))
	ms := int((sec - float64(int(sec))) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}
