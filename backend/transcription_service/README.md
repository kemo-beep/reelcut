# WhisperLiveKit ASR service

FastAPI app that exposes a **WebSocket `/asr`** for real-time transcription using [WhisperLiveKit](https://pypi.org/project/whisperlivekit/). The Reelcut Go backend connects to this service when `TRANSCRIPTION_WS_URL` is set.

**Important:** OpenAPI/Swagger does not document WebSocket endpoints. So at **http://localhost:8000/docs** you will see only HTTP routes (`/health`, `/asr` info). The actual transcription happens over **WebSocket** at `ws://localhost:8000/asr` — the Go worker connects there and sends binary audio.

**Persistence:** This service is stateless and does not use a database. The Go backend transcription worker receives segments over the WebSocket and saves them to Postgres (segments and words); it also updates the transcription row to `completed`.

## Run locally

```bash
cd backend/transcription_service
pip install -r requirements.txt
uvicorn app:app --host 0.0.0.0 --port 8000
```

- **Health:** `curl http://localhost:8000/health`
- **Confirm ASR service:** `curl http://localhost:8000/asr` (returns JSON describing the WebSocket)

Then set in the backend `.env`:

```env
TRANSCRIPTION_WS_URL=ws://localhost:8000
```

**Transcription will not work** until this Python service is running on port 8000. If only another app (e.g. with just `/health`) is on 8000, the Go worker will fail to transcribe.

## Protocol

- **Client** connects to `ws://host:port/asr` and sends **binary** audio chunks (e.g. WAV).
- **Server** streams JSON messages with `status`, `lines` (segments: `text`, `start`, `end`), and when done sends `{"type": "ready_to_stop"}`.
- To signal end-of-stream, the client sends an **empty binary message** after the last audio chunk.

## Optional env (service)

| Variable         | Default   | Description              |
|------------------|-----------|--------------------------|
| `WLK_MODEL`      | `medium`  | Whisper model size       |
| `WLK_DIARIZATION`| `true`    | Enable speaker diarization |
| `WLK_LANGUAGE`   | `en`      | Language code            |
