# Job queue (Asynq)

## Task types

- `video:metadata` — Extract duration/codec after upload
- `video:thumbnail` — Generate video thumbnail (1s offset)
- `video:fetch_url` — Fetch video from URL
- `transcription` — Whisper/WhisperLive transcription
- `analysis` — Video analysis
- `auto_cut` — Create clips from suggestions
- `render` — Render clip (cut, crop, captions, B-roll overlay, upload)
- `clip:thumbnail` — Generate clip thumbnail

## Recommended priorities

When multiple job types pile up, preferred order: **render** > **clip:thumbnail** > **video:thumbnail** > **transcription** > **analysis** > **auto_cut**. Configure Asynq queues (e.g. `critical`, `default`, `low`) and assign task types to queues in worker registration if needed.

## Validation and limits

- **B-roll uploads:** Max file size (e.g. 500 MB), segment times within clip duration.
- **Translation:** Rate limiting per user (beyond generic API limit); consider usage caps per tier.
- **Rendering:** B-roll compositing is CPU/memory intensive; size worker pool or use dedicated workers for heavy jobs.

## Failure handling

Failed jobs (after retries) can be moved to a dead-letter queue; monitor queue depth and failure rate. Document runbooks for common failures (storage unavailable, FFmpeg errors, etc.).
