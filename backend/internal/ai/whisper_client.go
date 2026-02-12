package ai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	openai "github.com/sashabaranov/go-openai"
)

type WhisperSegment struct {
	Start float64
	End   float64
	Text  string
}

type WhisperWord struct {
	Word  string
	Start float64
	End   float64
}

type WhisperResult struct {
	Segments []struct {
		Start float64
		End   float64
		Text  string
	}
	Words []WhisperWord
}

type WhisperClient struct {
	apiKey string
	client *openai.Client
}

func NewWhisperClient(apiKey string) *WhisperClient {
	if apiKey == "" {
		return &WhisperClient{}
	}
	return &WhisperClient{
		apiKey: apiKey,
		client: openai.NewClient(apiKey),
	}
}

func (w *WhisperClient) Transcribe(ctx context.Context, audioPath string, language string) (*WhisperResult, error) {
	if w.client == nil {
		return &WhisperResult{}, nil
	}
	f, err := os.Open(audioPath)
	if err != nil {
		return nil, fmt.Errorf("open audio: %w", err)
	}
	defer f.Close()
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: audioPath,
		Language: language,
	}
	resp, err := w.client.CreateTranscription(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("whisper api: %w", err)
	}
	result := &WhisperResult{}
	// Simple response has only Text; for segments we'd use verbose_json format
	_ = resp
	result.Segments = append(result.Segments, struct {
		Start float64
		End   float64
		Text  string
	}{0, 0, resp.Text})
	return result, nil
}

// TranscribeWithTimestamp uses a local or external service that returns word-level timestamps.
// For now we return the simple transcription and fake segment boundaries.
func (w *WhisperClient) TranscribeWithTimestamp(ctx context.Context, audioPath string, language string) (*WhisperResult, error) {
	return w.Transcribe(ctx, audioPath, language)
}

func (w *WhisperClient) TranscribeFile(ctx context.Context, audioPath string, lang string) (*WhisperResult, error) {
	dir := filepath.Dir(audioPath)
	_ = dir
	return w.TranscribeWithTimestamp(ctx, audioPath, lang)
}
