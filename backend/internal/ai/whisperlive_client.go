package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	whisperLiveChunkSize = 32 * 1024
	whisperLiveDialWait  = 30 * time.Second
	whisperLiveReadWait  = 5 * time.Minute
)

// WhisperLiveClient transcribes audio via a WhisperLiveKit WebSocket (/asr).
// Send binary audio chunks; server streams JSON (status, lines) and sends {"type":"ready_to_stop"} when done.
type WhisperLiveClient struct {
	baseURL string
}

// NewWhisperLiveClient creates a client for the WhisperLiveKit ASR WebSocket.
// baseURL is the base of the service (e.g. "ws://localhost:8000"); the client connects to baseURL + "/asr".
func NewWhisperLiveClient(baseURL string) *WhisperLiveClient {
	baseURL = strings.TrimSuffix(baseURL, "/")
	if baseURL == "" {
		baseURL = "ws://localhost:8000"
	}
	return &WhisperLiveClient{baseURL: baseURL}
}

// TranscribeFile reads the audio file, sends it over WebSocket to the ASR service, and returns segments.
func (c *WhisperLiveClient) TranscribeFile(ctx context.Context, audioPath string, _ string) (*WhisperResult, error) {
	data, err := os.ReadFile(audioPath)
	if err != nil {
		return nil, fmt.Errorf("read audio file: %w", err)
	}

	wsURL := c.baseURL + "/asr"
	if strings.HasPrefix(c.baseURL, "http://") {
		wsURL = "ws" + strings.TrimPrefix(c.baseURL, "http") + "/asr"
	} else if strings.HasPrefix(c.baseURL, "https://") {
		wsURL = "wss" + strings.TrimPrefix(c.baseURL, "https") + "/asr"
	}

	dialer := websocket.Dialer{HandshakeTimeout: whisperLiveDialWait}
	conn, _, err := dialer.DialContext(ctx, wsURL, http.Header{})
	if err != nil {
		return nil, fmt.Errorf("websocket dial: %w", err)
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(whisperLiveReadWait))

	// Send audio in chunks, then empty message to signal end-of-stream
	go func() {
		for i := 0; i < len(data); i += whisperLiveChunkSize {
			end := i + whisperLiveChunkSize
			if end > len(data) {
				end = len(data)
			}
			if err := conn.WriteMessage(websocket.BinaryMessage, data[i:end]); err != nil {
				return
			}
		}
		conn.WriteMessage(websocket.BinaryMessage, []byte{})
	}()

	var segments []struct {
		Start float64
		End   float64
		Text  string
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return nil, fmt.Errorf("websocket read: %w", err)
		}
		var obj map[string]interface{}
		if err := json.Unmarshal(msg, &obj); err != nil {
			continue
		}
		if t, _ := obj["type"].(string); t == "ready_to_stop" {
			break
		}
		if t, _ := obj["type"].(string); t == "config" {
			continue
		}
		lines, _ := obj["lines"].([]interface{})
		if len(lines) == 0 {
			continue
		}
		// Build new segment list from this message; server may send string or number for start/end
		var next []struct {
			Start float64
			End   float64
			Text  string
		}
		for _, l := range lines {
			line, _ := l.(map[string]interface{})
			text, _ := line["text"].(string)
			start := parseStartEnd(line["start"])
			end := parseStartEnd(line["end"])
			if end <= start {
				end = start + 0.001
			}
			next = append(next, struct {
				Start float64
				End   float64
				Text  string
			}{Start: start, End: end, Text: text})
		}
		// Only replace segments when we got at least one; otherwise keep previous (avoid losing data on empty updates)
		if len(next) > 0 {
			segments = next
		}
	}

	result := &WhisperResult{
		Segments: segments,
		Words:    nil,
	}
	return result, nil
}

// parseStartEnd parses start/end from server: either a number (float) or "H:MM:SS" string.
func parseStartEnd(v interface{}) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case string:
		return parseHHMMSS(x)
	default:
		return 0
	}
}

// parseHHMMSS parses "H:MM:SS" or "M:SS" or "S" to seconds (float).
func parseHHMMSS(s string) float64 {
	parts := strings.Split(strings.TrimSpace(s), ":")
	if len(parts) == 0 {
		return 0
	}
	var h, m, sec int
	if len(parts) == 3 {
		h, _ = strconv.Atoi(parts[0])
		m, _ = strconv.Atoi(parts[1])
		sec, _ = strconv.Atoi(parts[2])
	} else if len(parts) == 2 {
		m, _ = strconv.Atoi(parts[0])
		sec, _ = strconv.Atoi(parts[1])
	} else {
		sec, _ = strconv.Atoi(parts[0])
	}
	return float64(h*3600 + m*60 + sec)
}
