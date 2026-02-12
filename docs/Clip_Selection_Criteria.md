# Clip Selection Criteria

This document defines the criteria used to select and suggest clips from long-form video (e.g. podcasts, webinars) for Reels, TikTok, and Shorts. These criteria drive the **AI clip suggestions** (Gemini via `POST /api/v1/analysis/videos/:videoId/suggest-clips`) and should be reflected in prompts and validation. The full cut/clip flow is specified in [Cut/Clip feature specification](cutclip.md).

---

## 1. Primary criteria

| Criterion | Description |
|-----------|-------------|
| **Hook within first few seconds** | The clip should open with a hook or engaging phrase within the first ~3 seconds so viewers are retained. Prefer segments where the transcript starts with a question, bold claim, or attention-grabbing line. |
| **Complete thought or valuable insight** | Each clip should contain one complete idea, tip, or takeaway—not a fragment. Avoid clips that start or end mid-sentence or mid-thought. |
| **Natural start and end points** | Boundaries must align to **natural break points** in the transcript: sentence boundaries, topic shifts, or clear pauses. No mid-word or mid-sentence cuts when avoidable. |
| **Optimal duration** | Target **7–15 seconds** per clip for maximum engagement; **maximum 60 seconds** per clip. Duration is configurable via suggest-clips API (`min_duration`, `max_duration`). |
| **Variety across clips** | When suggesting multiple clips from the same video, favor **variety in topics or angles** so the set feels diverse and non-repetitive. |

---

## 2. Transcript-aware constraints

- **Source of truth:** Clip boundaries are derived from the **video transcript** (word/segment timestamps). The AI must use the transcript to propose `start_time` and `end_time`; it must not suggest arbitrary cuts that ignore speech boundaries.
- **Alignment:** Cuts should align to sentence or clause boundaries and, where possible, to topic changes so each clip is self-contained and coherent.
- **No mid-word / mid-sentence cuts:** When avoidable, do not start or end a clip in the middle of a word or sentence. Extend or shorten to the nearest natural break within the allowed duration range.

These constraints ensure that when FFmpeg cuts the source video and the app attaches a transcript slice to each clip, the result is watchable and readable (see [cutclip.md §2–3](cutclip.md)).

---

## 3. Optional / future enhancements

- **Virality / engagement scoring:** Score segments by engagement potential (e.g. hooks, emotional peaks, clarity) and surface higher-scoring clips first or filter by `min_virality_score` (see Application design document).
- **Emotional peaks:** Use sentiment or emphasis in the transcript to favor moments with stronger emotional or rhetorical impact.
- **Topic and scene context:** Use topic changes and, if available, scene or speaker boundaries to improve natural break detection.

---

## 4. Summary for implementation

| What | Detail |
|------|--------|
| **Target length** | 7–15 s (configurable); max 60 s per clip |
| **Boundaries** | Natural break points only; transcript-aware; no mid-word/mid-sentence when avoidable |
| **Content** | Hook early; complete thought; variety across suggested clips |
| **API** | `POST /api/v1/analysis/videos/:videoId/suggest-clips` with `min_duration`, `max_duration`, and preferences; backend uses Gemini and these criteria |

---

## 5. Related docs

- [Cut/Clip feature specification](cutclip.md) – Full flow: Gemini suggest-clips → FFmpeg cut → timeline, persistence, export.
- [Manual test: suggest-clips](Manual_test_suggest_clips.md) – How to test suggest-clips with real videos and transcription.
- [Application design document](Application_design_document.md) – AI clip detection, virality scoring, natural break identification.
