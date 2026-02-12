# Manual test: suggest-clips (real user video only)

Tests for `POST /api/v1/analysis/videos/:videoId/suggest-clips` use **only real videos** owned by the authenticated user. **No fake video IDs** (e.g. `00000000-0000-0000-0000-000000000001`) are used for the owner's video or as the primary test target.

## Resolving the test video

1. **API (recommended)**  
   Call `GET /api/v1/videos?per_page=10` with the same Bearer token. The response is `{ "data": { "videos": [ { "id": "<uuid>", ... } ] }, "pagination": { ... } }`. Use the first video `id`, or the first that has a completed transcription for suggest-clips.

2. **Script**  
   Run the backend script; it resolves a real video and runs suggest-clips:

   ```bash
   cd backend
   ./scripts/test_suggest_clips.sh
   ```

   Set `API_URL` if the API is not on `http://localhost:8080`:

   ```bash
   API_URL=http://localhost:8080 ./scripts/test_suggest_clips.sh
   ```

   The script uses token `TOKEN` (default: Bearer token for kemo@wonders.ai). It:
   - Lists videos via `GET /api/v1/videos`
   - Prefers a video with **completed transcription** so suggest-clips returns 200
   - Calls `POST .../videos/:videoId/suggest-clips` with that real `videoId`
   - Fails or skips if the user has no videos

## Suggest-clips and transcription

- **Suggest-clips** requires a completed transcription for the video. If the chosen video has none, the API returns 404 "Transcription not found".
- **Suggest-clips tests** therefore require at least one video with completed transcription for a successful positive test. The script prefers such a video when available; otherwise it uses the first video and documents that suggest-clips may 404.

## Negative tests (real IDs only)

- **Other user's video:** Use a **real** video ID that belongs to a different user. Set `OTHER_VIDEO_ID` when running the script to test 404 "Video not found" (ownership check). Do not use a fake UUID.
- **No transcript:** Use a **real** video owned by this user that has not been transcribed yet (if any); expect 404 "Transcription not found". Set `VIDEO_ID_NO_TRANSCRIPT` to that video's ID to run this negative test. Do not use a fake or non-existent video UUID for this.

## E2E: suggest-clips + auto-cut

The same script runs a full E2E flow when the API and token are valid:

1. **[1] Resolve video** – Prefers a video with completed transcription and, when available, a video with **no clips** (so auto-cut creates new clips).
2. **[2] Suggest-clips** – POST suggest-clips (7–60 s, max 20); validates response shape and duration.
3. **[3] Auto-cut** – POST auto-cut, poll `GET /api/v1/clips?video_id=...` until clips have `storage_path` (up to 90 s); validates each clip: 7–60 s, `storage_path` set, `status` ready.
4. **[4] Transcript overlap** – For each clip, asserts at least one transcription segment overlaps `[start_time, end_time]`.
5. **[5]/[6]** – Optional negative tests (other user’s video, video with no transcript) when `OTHER_VIDEO_ID` / `VIDEO_ID_NO_TRANSCRIPT` are set.

For a full auto-cut run, use a video that currently has **no clips**; the script prefers such a video when resolving. Token default: Bearer for **kemo@wonders.ai**; override with `TOKEN=...`.

## Summary

| What            | Action                                                                                        |
|-----------------|----------------------------------------------------------------------------------------------|
| Resolve video   | `GET /api/v1/videos` with Bearer token → use `data.videos[0].id` or first with transcript.  |
| Positive tests  | Use that real `VIDEO_ID` for suggest-clips.                                                  |
| Other user test | Set `OTHER_VIDEO_ID` to a real video ID from another user.                                    |
| No fake video   | Do not use `00000000-...` or other fake UUIDs for the owner's video or primary test target. |
