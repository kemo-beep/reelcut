package ai

import "context"

// Transcriber can transcribe an audio file to segments (and optional words).
// Implemented by WhisperClient (OpenAI) and WhisperLiveClient (WhisperLiveKit WebSocket).
type Transcriber interface {
	TranscribeFile(ctx context.Context, audioPath string, lang string) (*WhisperResult, error)
}
