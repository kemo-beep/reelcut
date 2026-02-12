// One-off test: run with TRANSCRIPTION_WS_URL=ws://localhost:8000 and audio path.
// Example: go run ./cmd/transcribe_test /tmp/sample_30s.wav
package main

import (
	"context"
	"fmt"
	"os"

	"reelcut/internal/ai"
)

func main() {
	audioPath := "/tmp/sample_30s.wav"
	if len(os.Args) > 1 {
		audioPath = os.Args[1]
	}
	wsURL := os.Getenv("TRANSCRIPTION_WS_URL")
	if wsURL == "" {
		wsURL = "ws://localhost:8000"
	}

	client := ai.NewWhisperLiveClient(wsURL)
	result, err := client.TranscribeFile(context.Background(), audioPath, "en")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Transcribe error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Segments: %d\n", len(result.Segments))
	for i, s := range result.Segments {
		fmt.Printf("  [%d] %.2f - %.2f: %q\n", i+1, s.Start, s.End, s.Text)
	}
}
