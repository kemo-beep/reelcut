# Cut/Clip feature specification

This document specifies how **auto-suggested clips** and the **clip timeline** should work. The current auto-cuts behavior is incorrect; the intended flow is described below.

## Goal

**Repurposing:** Take one **long-form video** (e.g. podcast, webinar, vlog) and cut it into **several short clips** optimized for:

- **Instagram Reels**
- **TikTok**
- **YouTube Shorts**

Users get multiple platform-ready clips from a single upload without re-recording. Clip lengths (e.g. 7–15 s, max 1 min) are chosen to fit these formats.

## Overview

**Important:** The original video is never modified or deleted. All cutting is **read-only** from the source; FFmpeg produces new clip files without altering the original.

The cut/clip flow must:

1. Use the **video transcript** (not arbitrary cuts).
2. Call **Gemini** to suggest clip boundaries (7–15 seconds per clip, max 1 minute per clip).
3. Use **FFmpeg** to cut the original video into the suggested clips.
4. Present clips on a **timeline** with a synced main viewer; the user can adjust clip durations (which updates transcript and timestamps) and export.

All clip and timeline state must **persist** across page refresh (stored server-side and/or synced on load).

---

## 1. Prerequisites

- **Transcription** must be completed for the video.  
  If there is no completed transcription, the suggest-clips step must not run (see [Manual test: suggest-clips](Manual_test_suggest_clips.md)).
- Use **real video IDs** for the owner’s video (e.g. from `GET /api/v1/videos?per_page=10`); do not use fake UUIDs for the primary test target.

---

## 2. AI clip suggestions (Gemini)

- **Input:** Full transcript (with word/segment timestamps) for the video.
- **Output:** A list of suggested clips, each with:
  - `start_time`, `end_time` (seconds)
  - Optional: short reason or label per clip
- **Constraints to send to Gemini:**
  - Target clip length: **7–15 seconds** (or configurable in this range).
  - **Max length per clip: 1 minute (60 seconds).**
  - Cuts should align to **natural break points** (sentence/topic boundaries) using the transcript; no mid-word or mid-sentence cuts when avoidable.

Implementation can go through the existing **suggest-clips** pipeline (`POST /api/v1/analysis/videos/:videoId/suggest-clips`) with backend updated to use **Gemini** and the above constraints (min/max duration, transcript-aware boundaries). The API already supports `min_duration`, `max_duration`, and preferences.

---

## 3. Creating clips from suggestions (FFmpeg)

After receiving suggestions from Gemini (via suggest-clips API):

1. For each suggested segment `(start_time, end_time)`:
   - Use **FFmpeg** to **extract** that range from the original video (the source file is only read, never modified or deleted).
   - Produce one **video clip file** per suggestion.
   - Attach the **subset of the transcript** that falls in `[start_time, end_time]`, with **timestamps** relative to the clip (e.g. 0.0 to duration) or absolute; store both if needed for the editor.
2. Persist each clip (e.g. via `POST /api/v1/clips`) with:
   - `start_time`, `end_time`, `duration_seconds`
   - Reference to transcript segment(s) or clip-level transcript + word timestamps
   - Storage path for the cut clip file (or trigger a render job that uses FFmpeg to produce it).

So: **Gemini suggests ranges → backend/worker uses FFmpeg to cut the original into clips → each clip has its own video file, transcript slice, and timestamps.** Clip **names** are derived from the **actual transcript** (segments overlapping the clip’s `[start_time, end_time]`), not from Gemini’s freeform transcript field. Clips with a cut file are playable via `GET /api/v1/clips/:id/playback-url`.

---

## 4. Timeline and main viewer

- **Timeline:** All clips for the video are arranged in order (by `start_time` or a dedicated order field) on a single timeline.
- **Main viewer:** One primary video player showing either:
  - The **full source video** with playhead synced to the timeline, or
  - The **currently selected clip** (or concatenated clip sequence).  
    Design doc assumes a timeline with a synced main view; exact UX can follow the existing [Application design document](Application_design_document.md) (e.g. TimelineEditor, VideoPlayer with markers).
- **Sync:** When the user selects a clip or moves the playhead on the timeline, the main viewer seeks to the corresponding time (or switches to the clip). When the user plays, the main viewer and timeline progress together.

---

## 5. User adjustments (durations, transcript, timestamps)

- The user can **change clip boundaries** (e.g. drag in/out on the timeline).
- When a clip’s **duration** is adjusted:
  - **Transcript** shown for that clip must be updated to the new `[start_time, end_time]` (subset of the full transcript).
  - **Timestamps** (word/segment) must be updated to reflect the new range (relative to clip or absolute, consistent with the rest of the app).
- These updates must **sync to the main view**: the main viewer reflects the new in/out points and, when playing, stays in sync with the timeline and transcript.

Persistence: clip boundary changes are saved via `PUT /api/v1/clips/:id` (and optionally trigger re-cut/re-render if we store pre-rendered clip files).

---

## 6. Export

- **Export as shorts:** User can export the current set of clips (e.g. batch) in formats/settings suitable for **Instagram Reels, TikTok, and YouTube Shorts** (e.g. 9:16 vertical, 7–15 s or up to 1 min depending on platform).
- **Export one clip:** User can export a single clip individually (same render pipeline, e.g. `POST /api/v1/clips/:id/render` then download).

Export uses the existing render pipeline (FFmpeg, aspect ratio, quality) as in the design doc and API.

---

## 7. Persistence and refresh

- Clip list, timeline order, and clip boundaries (start/end) must **persist** when the user refreshes the page.
- Implementation:
  - Store clips and metadata in the backend (existing clip and transcription APIs).
  - On load, fetch video → transcription → clips for that video; restore timeline and selection so the experience continues after refresh.

---

## 8. Summary table

| Step | What happens                                                                                                          |
| ---- | --------------------------------------------------------------------------------------------------------------------- |
| 1    | Ensure video has **completed transcription** (required for suggest-clips).                                            |
| 2    | Call suggest-clips (backed by **Gemini**) with **7–15 s** target, **max 60 s**, transcript-aware.                     |
| 3    | For each suggestion, use **FFmpeg** to cut the original into a clip; store **video + transcript slice + timestamps**. |
| 4    | Show clips on a **timeline** with **main viewer** synced to timeline/selection.                                       |
| 5    | User can **adjust clip durations**; transcript and timestamps update and **sync to main view**.                       |
| 6    | User can **export** all as shorts or **one clip** individually.                                                       |
| 7    | **State persists** on refresh (backend-stored clips + transcript).                                                    |

---

## 9. Related docs

- [Manual test: suggest-clips](Manual_test_suggest_clips.md) – Use real videos and completed transcription; script `./scripts/test_suggest_clips.sh`.
- [Application design document](Application_design_document.md) – AI clip suggestion flow, timeline editor, rendering with FFmpeg.
- [API routes](API_Routes.md) – `POST /api/v1/analysis/videos/:videoId/suggest-clips`, clip CRUD, render.
