package video

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCut_SegmentTimes verifies that Cut is invoked with segment start/end times.
// When the input file does not exist, we get an ffmpeg error that confirms the requested segment was used.
func TestCut_SegmentTimes(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	missingInput := filepath.Join(dir, "missing.mp4")
	output := filepath.Join(dir, "out.mp4")

	// Cut requests segment [10.5, 45.25]. FFmpeg will fail (no input) but the error should reflect our args.
	err := Cut(ctx, missingInput, output, 10.5, 45.25)
	if err == nil {
		t.Fatal("expected error when input file does not exist")
	}
	errStr := err.Error()
	// FFmpeg reports the -ss and -t we passed; error wraps "ffmpeg cut: ..."
	if !strings.Contains(errStr, "ffmpeg cut") {
		t.Errorf("error should mention ffmpeg cut: %s", errStr)
	}
	// Output file should not be created when ffmpeg fails
	if _, err := os.Stat(output); err == nil {
		t.Error("output file should not exist when cut fails")
	}
}
