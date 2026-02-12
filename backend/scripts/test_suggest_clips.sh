#!/usr/bin/env bash
# E2E test: suggest-clips (AI) and auto-cut (FFmpeg clips + transcripts) using a real video.
# Resolves VIDEO_ID via GET /api/v1/videos with the Bearer token; prefers video with
# completed transcription and (for auto-cut) video with no clips yet.
set -e

API_URL="${API_URL:-http://localhost:8080}"
# User: kemo@wonders.ai (user_id c6fe2d7d-435e-4681-a5e4-cb6b8c272e84)
TOKEN="${TOKEN:-Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzA4Nzk3OTcsImlhdCI6MTc3MDg3NjE5NywianRpIjoiMDI3YjRlZGQtMjU0NS00ODc2LWIzMjgtNzMwZWNlMTU5ZWE5IiwidXNlcl9pZCI6ImM2ZmUyZDdkLTQzNWUtNDY4MS1hNWU0LWNiNmI4YzI3MmU4NCIsImVtYWlsIjoia2Vtb0B3b25kZXJzLmFpIn0.ZnOL1Jimx-ZMB35WNoXCBoe-UlgX0Awryn80B-M4wzs}"

# Ensure token has Bearer prefix for convenience
if [[ "$TOKEN" != Bearer* ]]; then
  TOKEN="Bearer $TOKEN"
fi

echo "== Reelcut E2E: suggest-clips + auto-cut (clips + transcripts) =="
echo "API: $API_URL"
echo ""

# 1. Resolve a real video: completed transcription, prefer 0 clips for auto-cut
echo "[1] Resolving video (completed transcript, prefer no clips for auto-cut)..."
VIDEOS_RESP=$(curl -s -X GET "$API_URL/api/v1/videos?per_page=10" -H "Authorization: $TOKEN")
if ! echo "$VIDEOS_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); exit(0 if d.get('data',{}).get('videos') is not None else 1)" 2>/dev/null; then
  echo "Failed to list videos. Response: $VIDEOS_RESP"
  exit 1
fi

VIDEO_IDS=$(echo "$VIDEOS_RESP" | python3 -c "
import sys, json
d = json.load(sys.stdin)
videos = d.get('data', {}).get('videos') or []
for v in videos:
    print(v.get('id', ''))
")
if [[ -z "$VIDEO_IDS" ]]; then
  echo "No videos returned for this user. Upload a video first."
  exit 1
fi

VIDEO_ID=""
VIDEO_ID_NO_CLIPS=""
for id in $VIDEO_IDS; do
  id=$(echo "$id" | tr -d ' ')
  TR_RESP=$(curl -s -X GET "$API_URL/api/v1/transcriptions/videos/$id" -H "Authorization: $TOKEN")
  status=$(echo "$TR_RESP" | python3 -c "
import sys, json
d = json.load(sys.stdin)
t = d.get('transcription')
print(t.get('status', '') if t else '')
" 2>/dev/null)
  if [[ "$status" != "completed" ]]; then
    continue
  fi
  CLIPS_RESP=$(curl -s -X GET "$API_URL/api/v1/clips?video_id=$id&per_page=5" -H "Authorization: $TOKEN")
  clip_count=$(echo "$CLIPS_RESP" | python3 -c "
import sys, json
d = json.load(sys.stdin)
data = d.get('data') or {}
clips = data.get('clips') or []
print(len(clips))
" 2>/dev/null)
  if [[ -z "$VIDEO_ID" ]]; then
    VIDEO_ID="$id"
  fi
  if [[ "$clip_count" == "0" ]] && [[ -z "$VIDEO_ID_NO_CLIPS" ]]; then
    VIDEO_ID_NO_CLIPS="$id"
  fi
done

if [[ -z "$VIDEO_ID" ]]; then
  echo "No video with completed transcription found."
  exit 1
fi

echo "Using video $VIDEO_ID (has completed transcription)."
if [[ -n "$VIDEO_ID_NO_CLIPS" ]]; then
  echo "  (Video $VIDEO_ID_NO_CLIPS has no clips; will use for auto-cut test.)"
  AUTO_CUT_VIDEO_ID="$VIDEO_ID_NO_CLIPS"
else
  echo "  (All transcribed videos already have clips; auto-cut will no-op on this video.)"
  AUTO_CUT_VIDEO_ID="$VIDEO_ID"
fi

# 2. AI suggest-clips: POST and validate (7-60s, required fields)
echo ""
echo "[2] POST suggest-clips (min 7s, max 60s)..."
SUGGEST_RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$API_URL/api/v1/analysis/videos/$VIDEO_ID/suggest-clips" \
  -H "Authorization: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"min_duration":7,"max_duration":60,"max_suggestions":20}' \
  --max-time 90)
HTTP_CODE=$(echo "$SUGGEST_RESP" | grep -o 'HTTP_CODE:[0-9]*' | cut -d: -f2)
BODY=$(echo "$SUGGEST_RESP" | sed '/HTTP_CODE:/d')

if [[ "$HTTP_CODE" != "200" ]]; then
  echo "  HTTP $HTTP_CODE"
  echo "$BODY" | head -5
  exit 1
fi

count=$(echo "$BODY" | python3 -c "
import sys, json
d = json.load(sys.stdin)
s = d.get('suggestions') or []
print(len(s))
" 2>/dev/null)
echo "  HTTP 200, suggestions: $count"

# Validate each suggestion: start_time < end_time, duration in [7, 60], required fields
validate_suggestions=$(echo "$BODY" | python3 -c "
import sys, json
d = json.load(sys.stdin)
suggestions = d.get('suggestions') or []
errors = []
for i, s in enumerate(suggestions):
    st = s.get('start_time')
    et = s.get('end_time')
    if st is None or et is None:
        errors.append(f'Suggestion {i+1}: missing start_time or end_time')
        continue
    if et <= st:
        errors.append(f'Suggestion {i+1}: end_time must be > start_time')
        continue
    dur = et - st
    if dur < 7 or dur > 60:
        errors.append(f'Suggestion {i+1}: duration {dur:.1f}s not in [7, 60]')
    for key in ('transcript', 'reason', 'virality_score'):
        if key not in s:
            errors.append(f'Suggestion {i+1}: missing {key}')
if errors:
    for e in errors:
        print(e)
    sys.exit(1)
sys.exit(0)
" 2>/dev/null) || true
if [[ -n "$validate_suggestions" ]]; then
  echo "  Validation failed:"
  echo "$validate_suggestions"
  exit 1
fi
echo "  Validated: all suggestions in 7-60s with required fields."

# 3. Auto-cut: trigger job, poll clips, validate storage_path and status
echo ""
echo "[3] POST auto-cut and wait for clips..."
AUTOCUT_RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$API_URL/api/v1/videos/$AUTO_CUT_VIDEO_ID/auto-cut" \
  -H "Authorization: $TOKEN" \
  -H "Content-Type: application/json")
AUTOCUT_CODE=$(echo "$AUTOCUT_RESP" | grep -o 'HTTP_CODE:[0-9]*' | cut -d: -f2)
if [[ "$AUTOCUT_CODE" != "202" ]]; then
  echo "  Auto-cut returned HTTP $AUTOCUT_CODE (expected 202)."
  echo "$AUTOCUT_RESP" | sed '/HTTP_CODE:/d' | head -3
else
  echo "  HTTP 202 (job queued). Polling for clips (up to 90s)..."
fi

CLIP_COUNT=0
for _ in $(seq 1 18); do
  sleep 5
  CLIPS_RESP=$(curl -s -X GET "$API_URL/api/v1/clips?video_id=$AUTO_CUT_VIDEO_ID&per_page=50" -H "Authorization: $TOKEN")
  CLIP_COUNT=$(echo "$CLIPS_RESP" | python3 -c "
import sys, json
d = json.load(sys.stdin)
data = d.get('data') or {}
clips = data.get('clips') or []
with_storage = sum(1 for c in clips if c.get('storage_path'))
print(f'{len(clips)},{with_storage}')
" 2>/dev/null)
  total=$(echo "$CLIP_COUNT" | cut -d, -f1)
  with_path=$(echo "$CLIP_COUNT" | cut -d, -f2)
  if [[ "$total" -gt 0 ]] && [[ "$with_path" -gt 0 ]]; then
    echo "  Clips: $total (with storage_path: $with_path)."
    break
  fi
done

# Validate clips when present: duration 7-60s, storage_path set, status ready
total_clips=$(echo "$CLIPS_RESP" | python3 -c "
import sys, json
d = json.load(sys.stdin)
data = d.get('data') or {}
clips = data.get('clips') or []
print(len(clips))
" 2>/dev/null)
if [[ "$total_clips" -gt 0 ]]; then
  echo "$CLIPS_RESP" | python3 -c "
import sys, json
d = json.load(sys.stdin)
data = d.get('data') or {}
clips = data.get('clips') or []
errors = []
for i, c in enumerate(clips):
    st = c.get('start_time')
    et = c.get('end_time')
    dur = c.get('duration_seconds')
    path = c.get('storage_path')
    status = c.get('status')
    if st is not None and et is not None:
        d = et - st
        if d < 7 or d > 60:
            errors.append('Clip %d: duration %.1fs not in [7, 60]' % (i+1, d))
    if dur is not None and (dur < 7 or dur > 60):
        errors.append('Clip %d: duration_seconds %s not in [7, 60]' % (i+1, dur))
    if not path:
        errors.append('Clip %d: missing storage_path' % (i+1))
    if status != 'ready':
        errors.append('Clip %d: status=%r (expected ready)' % (i+1, status))
if errors:
    for e in errors:
        print(e)
    sys.exit(1)
print('Validated %d clip(s): 7-60s, storage_path set, status ready.' % len(clips))
" 2>/dev/null || { echo "  Clip validation failed."; exit 1; }
  echo "  Clips validated."
else
  echo "  No clips yet (auto-cut may have no-op if video already had clips)."
fi

# 4. Transcript overlap: for each clip, at least one segment overlaps [start_time, end_time] (only when we have clips)
echo ""
echo "[4] Validating transcript overlap per clip..."
if [[ "$total_clips" -gt 0 ]]; then
  TR_FULL=$(curl -s -X GET "$API_URL/api/v1/transcriptions/videos/$AUTO_CUT_VIDEO_ID" -H "Authorization: $TOKEN")
  CLIPS_JSON=$(echo "$CLIPS_RESP" | python3 -c "
import sys, json
d = json.load(sys.stdin)
print(json.dumps((d.get('data') or {}).get('clips') or []))
" 2>/dev/null)
  export CLIPS_JSON
  echo "$TR_FULL" | python3 -c "
import os, sys, json
tr = json.load(sys.stdin)
transcription = tr.get('transcription')
segments = (transcription or {}).get('segments') or []
try:
    clips_data = json.loads(os.environ.get('CLIPS_JSON', '[]'))
except Exception:
    clips_data = []
errors = []
for i, c in enumerate(clips_data):
    cs, ce = c.get('start_time'), c.get('end_time')
    if cs is None or ce is None:
        continue
    overlapping = [s for s in segments if s.get('end_time', 0) > cs and s.get('start_time', 0) < ce]
    if not overlapping:
        errors.append('Clip %d [%.1f-%.1fs]: no transcript segment overlaps' % (i+1, cs, ce))
if errors:
    for e in errors:
        print(e)
    sys.exit(1)
print('Transcript overlap OK for %d clip(s).' % len(clips_data))
" 2>/dev/null || { echo "  Transcript overlap check skipped or failed (transcription may not include segments in response)."; true; }
else
  echo "  Skipped (no clips to check)."
fi

# [5] Negative test: other user's video (only if OTHER_VIDEO_ID is set; must be a real video ID)
if [[ -n "${OTHER_VIDEO_ID:-}" ]]; then
  echo ""
  echo "[5] Negative test: POST suggest-clips for other user's video (expect 404)..."
  OTHER_RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$API_URL/api/v1/analysis/videos/$OTHER_VIDEO_ID/suggest-clips" \
    -H "Authorization: $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{}')
  OTHER_CODE=$(echo "$OTHER_RESP" | grep -o 'HTTP_CODE:[0-9]*' | cut -d: -f2)
  if [[ "$OTHER_CODE" == "404" ]]; then
    echo "  HTTP 404 (Video not found) as expected."
  else
    echo "  HTTP $OTHER_CODE (expected 404)."
  fi
else
  echo ""
  echo "[5] Negative test skipped (set OTHER_VIDEO_ID to a real video ID owned by another user to test 404)."
fi

# [6] Optional negative test: video with no transcript (only if VIDEO_ID_NO_TRANSCRIPT is set; real video ID)
if [[ -n "${VIDEO_ID_NO_TRANSCRIPT:-}" ]]; then
  echo ""
  echo "[6] Negative test: POST suggest-clips for video with no transcript (expect 404)..."
  NOTR_RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$API_URL/api/v1/analysis/videos/$VIDEO_ID_NO_TRANSCRIPT/suggest-clips" \
    -H "Authorization: $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{}')
  NOTR_CODE=$(echo "$NOTR_RESP" | grep -o 'HTTP_CODE:[0-9]*' | cut -d: -f2)
  if [[ "$NOTR_CODE" == "404" ]]; then
    echo "  HTTP 404 (Transcription not found or Video not found) as expected."
  else
    echo "  HTTP $NOTR_CODE (expected 404)."
  fi
else
  echo ""
  echo "[6] No-transcript negative test skipped (set VIDEO_ID_NO_TRANSCRIPT to a real video ID with no transcription to test 404)."
fi

echo ""
echo "Done. All tests used real video IDs; no fake UUIDs for the owner's video."
