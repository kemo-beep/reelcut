package ai

import (
	"context"
	"encoding/json"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// SceneRange is a detected scene [Start, End] in seconds.
type SceneRange struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// DetectScenes uses FFmpeg scene filter to get scene change timestamps.
// Returns nil if video path is empty or FFmpeg fails.
func DetectScenes(ctx context.Context, videoPath string) ([]SceneRange, error) {
	if videoPath == "" {
		return nil, nil
	}
	// FFmpeg: select scene changes + showinfo to get pts_time
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", videoPath,
		"-vf", "select='gt(scene,0.4)',showinfo",
		"-vsync", "vfr",
		"-f", "null", "-",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, nil
	}
	// Parse stderr for "pts_time:1.234" or similar
	ptsRe := regexp.MustCompile(`pts_time:([\d.]+)`)
	matches := ptsRe.FindAllStringSubmatch(string(out), -1)
	if len(matches) == 0 {
		return []SceneRange{{Start: 0, End: 0}}, nil
	}
	var times []float64
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		t, _ := strconv.ParseFloat(strings.TrimSpace(m[1]), 64)
		times = append(times, t)
	}
	if len(times) == 0 {
		return []SceneRange{{Start: 0, End: 0}}, nil
	}
	var ranges []SceneRange
	for i := 0; i < len(times); i++ {
		end := 0.0
		if i+1 < len(times) {
			end = times[i+1]
		} else {
			end = times[i] + 10
		}
		ranges = append(ranges, SceneRange{Start: times[i], End: end})
	}
	return ranges, nil
}

// ScenesToJSON returns JSON array for video_analysis.scenes_detected.
func ScenesToJSON(ranges []SceneRange) ([]byte, error) {
	return json.Marshal(ranges)
}
