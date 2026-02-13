package config

// ExportPreset defines platform-specific export settings (aspect ratio, dimensions, codec, bitrate).
// Used for render so output matches platform best practices (e.g. TikTok/Reels 9:16 1080x1920).
type ExportPreset struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	AspectRatio  string `json:"aspect_ratio"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	VideoBitrate int    `json:"video_bitrate_kbps"` // kbps, 0 = default
	AudioBitrate int    `json:"audio_bitrate_kbps"` // kbps, 0 = default
	FPS          int    `json:"fps"`                // 0 = source
}

// ExportPresets are the built-in platform presets. Can be overridden or extended via config file later.
var ExportPresets = []ExportPreset{
	{
		ID:           "tiktok",
		Name:         "TikTok / Reels",
		Description:  "9:16 vertical, 1080×1920, H.264, 30fps, 5–8 Mbps",
		AspectRatio:  "9:16",
		Width:        1080,
		Height:       1920,
		VideoBitrate: 6000,
		AudioBitrate: 128,
		FPS:          30,
	},
	{
		ID:           "reels",
		Name:         "Instagram Reels",
		Description:  "9:16 vertical, 1080×1920, H.264, 30fps",
		AspectRatio:  "9:16",
		Width:        1080,
		Height:       1920,
		VideoBitrate: 6000,
		AudioBitrate: 128,
		FPS:          30,
	},
	{
		ID:           "instagram_feed",
		Name:         "Instagram Feed",
		Description:  "1:1 square, 1080×1080, H.264, 30fps, ~3.5 Mbps",
		AspectRatio:  "1:1",
		Width:        1080,
		Height:       1080,
		VideoBitrate: 3500,
		AudioBitrate: 128,
		FPS:          30,
	},
	{
		ID:           "youtube_shorts",
		Name:         "YouTube Shorts",
		Description:  "9:16 vertical, 1080×1920, H.264, 24–60fps, 5–10 Mbps",
		AspectRatio:  "9:16",
		Width:        1080,
		Height:       1920,
		VideoBitrate: 8000,
		AudioBitrate: 128,
		FPS:          30,
	},
}

// GetExportPresetByID returns the preset with the given id or nil.
func GetExportPresetByID(id string) *ExportPreset {
	for i := range ExportPresets {
		if ExportPresets[i].ID == id {
			return &ExportPresets[i]
		}
	}
	return nil
}
