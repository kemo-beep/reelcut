package video

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

type Metadata struct {
	DurationSeconds float64 `json:"duration_seconds"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	FPS            float64 `json:"fps"`
	Codec          string  `json:"codec"`
	Bitrate        int     `json:"bitrate"`
	FileSizeBytes  int64   `json:"file_size_bytes"`
}

var (
	reDuration = regexp.MustCompile(`Duration: (\d{2}):(\d{2}):(\d{2})\.(\d{2})`)
	reStream   = regexp.MustCompile(`Stream #\d+:\d+.*: Video: (\w+).* (\d+)x(\d+).* (\d+(?:\.\d+)?) fps`)
)

func GetMetadata(ctx context.Context, path string) (*Metadata, error) {
	args := []string{"-v", "error", "-show_entries", "format=duration,size:stream=width,height,r_frame_rate,codec_name,bit_rate", "-of", "json", path}
	out, err := RunFFprobe(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("ffprobe: %w", err)
	}
	var probe struct {
		Format struct {
			Duration string `json:"duration"`
			Size     string `json:"size"`
		} `json:"format"`
		Streams []struct {
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			RFrameRate string `json:"r_frame_rate"`
			CodecName string `json:"codec_name"`
			BitRate  string `json:"bit_rate"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(out, &probe); err != nil {
		return nil, fmt.Errorf("parse ffprobe json: %w", err)
	}
	meta := &Metadata{}
	if probe.Format.Duration != "" {
		if d, err := strconv.ParseFloat(probe.Format.Duration, 64); err == nil {
			meta.DurationSeconds = d
		}
	}
	if probe.Format.Size != "" {
		if s, err := strconv.ParseInt(probe.Format.Size, 10, 64); err == nil {
			meta.FileSizeBytes = s
		}
	}
	for _, s := range probe.Streams {
		if s.CodecName != "" && s.Width > 0 {
			meta.Width = s.Width
			meta.Height = s.Height
			meta.Codec = s.CodecName
			if s.BitRate != "" {
				if b, err := strconv.Atoi(s.BitRate); err == nil {
					meta.Bitrate = b / 1000
				}
			}
			if s.RFrameRate != "" {
				var num, den int
				fmt.Sscanf(s.RFrameRate, "%d/%d", &num, &den)
				if den > 0 {
					meta.FPS = float64(num) / float64(den)
				}
			}
			break
		}
	}
	return meta, nil
}
