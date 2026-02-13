package video

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunFFmpeg(ctx context.Context, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	return cmd.CombinedOutput()
}

func RunFFprobe(ctx context.Context, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "ffprobe", args...)
	return cmd.CombinedOutput()
}

func ExtractFrame(ctx context.Context, inputPath string, timestampSec float64, outputPath string) error {
	args := []string{
		"-y",
		"-ss", fmt.Sprintf("%.3f", timestampSec),
		"-i", inputPath,
		"-vframes", "1",
		"-q:v", "2",
		outputPath,
	}
	out, err := RunFFmpeg(ctx, args...)
	if err != nil {
		return fmt.Errorf("ffmpeg extract frame: %w (output: %s)", err, string(out))
	}
	return nil
}

func ExtractAudio(ctx context.Context, inputPath, outputPath string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	args := []string{
		"-y",
		"-i", inputPath,
		"-vn",
		"-acodec", "pcm_s16le",
		"-ar", "16000",
		"-ac", "1",
		outputPath,
	}
	out, err := RunFFmpeg(ctx, args...)
	if err != nil {
		return fmt.Errorf("ffmpeg extract audio: %w (output: %s)", err, string(out))
	}
	return nil
}

// ExtractAudioChunk extracts a time range [startSec, startSec+durationSec) from an audio file
// into outputPath, using the same format as ExtractAudio (pcm_s16le, 16kHz, mono).
func ExtractAudioChunk(ctx context.Context, inputPath, outputPath string, startSec, durationSec float64) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	args := []string{
		"-y",
		"-ss", fmt.Sprintf("%.3f", startSec),
		"-i", inputPath,
		"-t", fmt.Sprintf("%.3f", durationSec),
		"-vn",
		"-acodec", "pcm_s16le",
		"-ar", "16000",
		"-ac", "1",
		outputPath,
	}
	out, err := RunFFmpeg(ctx, args...)
	if err != nil {
		return fmt.Errorf("ffmpeg extract audio chunk: %w (output: %s)", err, string(out))
	}
	return nil
}

// ExtractAudioChunkFromVideo extracts a time range [startSec, startSec+durationSec) from a video file
// into outputPath, using the same format as ExtractAudio (pcm_s16le, 16kHz, mono).
func ExtractAudioChunkFromVideo(ctx context.Context, videoPath, outputPath string, startSec, durationSec float64) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	args := []string{
		"-y",
		"-ss", fmt.Sprintf("%.3f", startSec),
		"-i", videoPath,
		"-t", fmt.Sprintf("%.3f", durationSec),
		"-vn",
		"-acodec", "pcm_s16le",
		"-ar", "16000",
		"-ac", "1",
		outputPath,
	}
	out, err := RunFFmpeg(ctx, args...)
	if err != nil {
		return fmt.Errorf("ffmpeg extract audio chunk from video: %w (output: %s)", err, string(out))
	}
	return nil
}

// Cut trims video to [start, end] seconds.
func Cut(ctx context.Context, inputPath, outputPath string, start, end float64) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	dur := end - start
	args := []string{
		"-y",
		"-ss", fmt.Sprintf("%.3f", start),
		"-i", inputPath,
		"-t", fmt.Sprintf("%.3f", dur),
		"-c", "copy",
		outputPath,
	}
	out, err := RunFFmpeg(ctx, args...)
	if err != nil {
		return fmt.Errorf("ffmpeg cut: %w (output: %s)", err, string(out))
	}
	return nil
}

// ResizeCrop scales and crops to aspect ratio (e.g. "9:16", "1:1", "16:9").
func ResizeCrop(ctx context.Context, inputPath, outputPath, aspectRatio string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	// scale and crop: scale to cover then crop to aspect
	var vf string
	switch aspectRatio {
	case "9:16":
		vf = "scale=-2:ih*2,crop=ih*9/16:ih"
	case "1:1":
		vf = "scale=min(iw\\,ih):min(iw\\,ih),crop=min(iw\\,ih):min(iw\\,ih)"
	case "16:9":
		vf = "scale=iw*2:-2,crop=iw:iw*9/16"
	default:
		vf = "scale=-2:ih*2,crop=ih*9/16:ih"
	}
	args := []string{
		"-y", "-i", inputPath,
		"-vf", vf,
		"-c:a", "copy",
		outputPath,
	}
	out, err := RunFFmpeg(ctx, args...)
	if err != nil {
		return fmt.Errorf("ffmpeg resize: %w (output: %s)", err, string(out))
	}
	return nil
}

// ResizeCropToSize scales to cover and center-crops to exact width x height.
// Use for platform presets (e.g. 1080x1920 for TikTok/Reels).
func ResizeCropToSize(ctx context.Context, inputPath, outputPath string, width, height int) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	// scale to cover (force_original_aspect_ratio=increase) then center crop
	vf := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=increase,crop=%d:%d:(iw-%d)/2:(ih-%d)/2",
		width, height, width, height, width, height)
	args := []string{
		"-y", "-i", inputPath,
		"-vf", vf,
		"-c:a", "copy",
		outputPath,
	}
	out, err := RunFFmpeg(ctx, args...)
	if err != nil {
		return fmt.Errorf("ffmpeg resize to size: %w (output: %s)", err, string(out))
	}
	return nil
}

// BurnSubtitles burns SRT file into video (plain styling).
func BurnSubtitles(ctx context.Context, inputPath, srtPath, outputPath string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	// Escape path for filter (Windows/Unix)
	escaped := srtPath
	if filepath.Separator == '\\' {
		escaped = strings.ReplaceAll(srtPath, "\\", "\\\\")
	}
	escaped = strings.ReplaceAll(escaped, ":", "\\:")
	args := []string{
		"-y", "-i", inputPath,
		"-vf", "subtitles=" + escaped,
		"-c:a", "copy",
		outputPath,
	}
	out, err := RunFFmpeg(ctx, args...)
	if err != nil {
		return fmt.Errorf("ffmpeg subtitles: %w (output: %s)", err, string(out))
	}
	return nil
}

// BurnASS burns an ASS subtitle file into video with full style (font, color, position).
// ASS supports styled captions; use this when clip_styles caption options are set.
func BurnASS(ctx context.Context, inputPath, assPath, outputPath string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	escaped := assPath
	if filepath.Separator == '\\' {
		escaped = strings.ReplaceAll(assPath, "\\", "\\\\")
	}
	escaped = strings.ReplaceAll(escaped, ":", "\\:")
	args := []string{
		"-y", "-i", inputPath,
		"-vf", "ass=" + escaped,
		"-c:a", "copy",
		outputPath,
	}
	out, err := RunFFmpeg(ctx, args...)
	if err != nil {
		return fmt.Errorf("ffmpeg ass: %w (output: %s)", err, string(out))
	}
	return nil
}

// OverlayImage overlays image on video at position (e.g. "bottom-right").
func OverlayImage(ctx context.Context, inputPath, imagePath, outputPath, position string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	var overlay string
	switch position {
	case "top-left":
		overlay = "overlay=10:10"
	case "top-right":
		overlay = "overlay=main_w-overlay_w-10:10"
	case "bottom-left":
		overlay = "overlay=10:main_h-overlay_h-10"
	case "bottom-right":
		overlay = "overlay=main_w-overlay_w-10:main_h-overlay_h-10"
	default:
		overlay = "overlay=main_w-overlay_w-10:main_h-overlay_h-10"
	}
	args := []string{
		"-y", "-i", inputPath, "-i", imagePath,
		"-filter_complex", "[1]scale=iw*0.2:-1[logo];[0][logo]" + overlay,
		"-c:a", "copy",
		outputPath,
	}
	out, err := RunFFmpeg(ctx, args...)
	if err != nil {
		return fmt.Errorf("ffmpeg overlay: %w (output: %s)", err, string(out))
	}
	return nil
}

// OverlayVideo overlays a B-roll video on the main video between startTime and endTime (clip-relative seconds).
// scale is the factor for the overlay size (e.g. 0.5 = half width); opacity is 0-1.
func OverlayVideo(ctx context.Context, mainPath, brollPath, outputPath string, startTime, endTime, scale, opacity float64) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	if scale <= 0 {
		scale = 0.5
	}
	if opacity <= 0 {
		opacity = 1
	}
	// [1:v] scale and opacity -> [ov]; [0:v][ov] overlay with enable=between(t,start,end), centered
	// Escape single quotes in filter for enable=
	enable := fmt.Sprintf("between(t\\,%.3f\\,%.3f)", startTime, endTime)
	scaleStr := fmt.Sprintf("scale=iw*%.2f:-1", scale)
	alphaStr := fmt.Sprintf("format=rgba,colorchannelmixer=aa=%.2f", opacity)
	overlayPos := "(main_w-overlay_w)/2:(main_h-overlay_h)/2"
	filter := "[1:v]" + scaleStr + "," + alphaStr + "[ov];[0:v][ov]overlay=" + overlayPos + ":enable='" + enable + "'"
	args := []string{
		"-y", "-i", mainPath, "-i", brollPath,
		"-filter_complex", filter,
		"-c:a", "copy",
		outputPath,
	}
	out, err := RunFFmpeg(ctx, args...)
	if err != nil {
		return fmt.Errorf("ffmpeg overlay video: %w (output: %s)", err, string(out))
	}
	return nil
}

// MixAudio mixes background audio with video at volume (0-1).
func MixAudio(ctx context.Context, videoInput, audioInput, outputPath string, audioVolume float64) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	vol := fmt.Sprintf("%.2f", audioVolume)
	args := []string{
		"-y", "-i", videoInput, "-i", audioInput,
		"-filter_complex", "[1]volume=" + vol + "[a1];[0:a][a1]amix=inputs=2:duration=first[aout]",
		"-map", "0:v", "-map", "[aout]",
		"-c:v", "copy",
		outputPath,
	}
	out, err := RunFFmpeg(ctx, args...)
	if err != nil {
		return fmt.Errorf("ffmpeg mix audio: %w (output: %s)", err, string(out))
	}
	return nil
}
