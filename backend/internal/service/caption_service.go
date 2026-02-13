package service

import (
	"fmt"
	"regexp"
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

// AllowedCaptionFonts is the fallback chain and allowed set for caption fonts.
// Workers should have these available or use Arial as final fallback.
// No custom font file uploads; use only these to avoid licensing/security issues.
var AllowedCaptionFonts = []string{"Arial", "Inter", "Montserrat", "Open Sans", "Roboto", "Helvetica"}

// DefaultCaptionFont is used when style font is missing or not in allowed set.
const DefaultCaptionFont = "Arial"

func allowedFont(style *domain.ClipStyle) string {
	if style == nil || style.CaptionFont == "" {
		return DefaultCaptionFont
	}
	for _, f := range AllowedCaptionFonts {
		if strings.EqualFold(style.CaptionFont, f) {
			return style.CaptionFont
		}
	}
	return DefaultCaptionFont
}

// assAlignment returns ASS alignment number: 2=bottom center, 5=middle, 8=top center.
func assAlignment(position string) int {
	switch strings.ToLower(position) {
	case "top":
		return 8
	case "center", "centre":
		return 5
	case "bottom":
	default:
	}
	return 2
}

// hexColorToASS converts #RRGGBB or #AARRGGBB to ASS &HAABBGGRR (BGR order, alpha optional).
func hexColorToASS(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 6 {
		// No alpha: use BBGGRR, alpha 0 = opaque in ASS
		return "&H00" + hex[4:6] + hex[2:4] + hex[0:2]
	}
	if len(hex) == 8 {
		// AARRGGBB -> AABBGGRR
		a, r, g, b := hex[0:2], hex[2:4], hex[4:6], hex[6:8]
		return "&H" + a + b + g + r
	}
	return "&H00FFFFFF" // white
}

// ToASS produces ASS subtitles with styling from clip style (font, size, color, position, background).
// Used for burning styled captions into video; FFmpeg's ass filter respects these styles.
func ToASS(blocks []CaptionBlock, style *domain.ClipStyle) string {
	font := allowedFont(style)
	size := 48
	if style != nil && style.CaptionSize > 0 {
		size = style.CaptionSize
	}
	if size > 120 {
		size = 120
	}
	primary := "&H00FFFFFF"
	border := "&H80000000" // semi-transparent black outline for readability
	backColour := "&H80000000"
	if style != nil && style.CaptionColor != "" {
		primary = hexColorToASS(style.CaptionColor)
	}
	if style != nil && style.CaptionBgColor != nil && *style.CaptionBgColor != "" {
		backColour = hexColorToASS(*style.CaptionBgColor)
	}
	align := assAlignment("bottom")
	if style != nil && style.CaptionPosition != "" {
		align = assAlignment(style.CaptionPosition)
	}

	var b strings.Builder
	b.WriteString("[Script Info]\n")
	b.WriteString("ScriptType: v4.00+\n\n")
	b.WriteString("[V4+ Styles]\n")
	b.WriteString("Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding\n")
	b.WriteString(fmt.Sprintf("Style: Default,%s,%d,%s,&H000000FF,%s,%s,0,0,0,0,100,100,0,0,1,2,1,%d,10,10,30,1\n\n", font, size, primary, border, backColour, align))
	b.WriteString("[Events]\n")
	b.WriteString("Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n")

	for _, blk := range blocks {
		start := formatASSTime(blk.StartTime)
		end := formatASSTime(blk.EndTime)
		text := escapeASSText(blk.Text)
		b.WriteString(fmt.Sprintf("Dialogue: 0,%s,%s,Default,,0,0,0,,%s\n", start, end, text))
	}
	return b.String()
}

func formatASSTime(sec float64) string {
	h := int(sec / 3600)
	m := int((sec - float64(h*3600)) / 60)
	s := int(sec - float64(h*3600) - float64(m*60))
	cs := int((sec - float64(int(sec))) * 100)
	return fmt.Sprintf("%d:%02d:%02d.%02d", h, m, s, cs)
}

var assSpecialRe = regexp.MustCompile(`[{}]`)

func escapeASSText(s string) string {
	// ASS uses \N for newline, \{ and \} for literal braces
	return assSpecialRe.ReplaceAllStringFunc(s, func(c string) string {
		if c == "{" {
			return "\\{"
		}
		return "\\}"
	})
}
