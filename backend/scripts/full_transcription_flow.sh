#!/usr/bin/env bash
# Full flow: Docker (Postgres, Redis, MinIO) -> ASR service -> Go API -> upload sample video -> start transcription -> poll until done.
set -e
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
BACKEND="$REPO_ROOT/backend"
SAMPLE_VIDEO="${1:-$REPO_ROOT/samples/_samplevideo.mp4}"
API_URL="${API_URL:-http://localhost:8080}"
ASR_PORT="${ASR_PORT:-8000}"

echo "== Reelcut full transcription flow =="
echo "Sample video: $SAMPLE_VIDEO"
echo "API: $API_URL"
echo ""

if [[ ! -f "$SAMPLE_VIDEO" ]]; then
  echo "Sample video not found: $SAMPLE_VIDEO"
  exit 1
fi

# Free ports (optional)
for port in 8080 "$ASR_PORT"; do
  pid=$(lsof -ti ":$port" 2>/dev/null) || true
  if [[ -n "$pid" ]]; then
    echo "Stopping process on port $port (PID $pid)..."
    kill $pid 2>/dev/null || true
    sleep 2
  fi
done

# 1. Docker
echo "[1/7] Starting Postgres, Redis, MinIO..."
cd "$BACKEND"
docker compose up -d
echo "Waiting for services..."
sleep 5
for i in 1 2 3 4 5 6 7 8 9 10; do
  if docker compose exec -T postgres pg_isready -U postgres -q 2>/dev/null; then break; fi
  sleep 2
done
docker compose exec -T postgres pg_isready -U postgres

# Create MinIO bucket if missing
echo "Ensuring MinIO bucket 'reelcut' exists..."
docker run --rm --add-host=host.docker.internal:host-gateway minio/mc alias set local http://host.docker.internal:9002 minioadmin minioadmin 2>/dev/null || true
docker run --rm --add-host=host.docker.internal:host-gateway minio/mc mb local/reelcut 2>/dev/null || true

# Ensure dev user has known password (password123) — hash from bcrypt cost 10
DEV_HASH='$2a$10$hsXIJ2gxRL/Q8NKjQBqciu2TEBCMgr6DDoCJeDhNIKoW4/zWgD9ie'
echo "Setting dev user password..."
echo "UPDATE users SET password_hash = '$DEV_HASH' WHERE email = 'dev@reelcut.local';" | docker compose exec -T postgres psql -U postgres -d reelcut -f - 2>/dev/null || true

# 2. ASR service (background)
echo "[2/7] Starting WhisperLiveKit ASR (port $ASR_PORT)..."
cd "$BACKEND/transcription_service"
WLK_MODEL=tiny WLK_DIARIZATION=false python3 -m uvicorn app:app --host 127.0.0.1 --port "$ASR_PORT" &
ASR_PID=$!
# Wait for ASR to be up (model load can take 30-60s)
for i in $(seq 1 90); do
  if curl -s "http://127.0.0.1:$ASR_PORT/health" 2>/dev/null | grep -q '"status":"ok"'; then
    echo "ASR ready."
    break
  fi
  if [[ $i -eq 90 ]]; then
    echo "ASR did not become ready in time."
    kill $ASR_PID 2>/dev/null || true
    exit 1
  fi
  sleep 2
done

# 3. Go API (background; migrations run on startup)
echo "[3/7] Starting Go API..."
cd "$BACKEND"
export TRANSCRIPTION_WS_URL="ws://127.0.0.1:$ASR_PORT"
go run ./cmd/api &
API_PID=$!
for i in $(seq 1 30); do
  if curl -s "$API_URL/health" >/dev/null 2>&1; then break; fi
  sleep 1
done
if ! curl -s "$API_URL/health" >/dev/null 2>&1; then
  echo "API failed to become healthy"
  kill $ASR_PID 2>/dev/null || true
  exit 1
fi
echo "API ready."

# 4. Login
echo "[4/7] Login and get token..."
LOGIN_RESP=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@reelcut.local","password":"password123"}')
TOKEN=$(echo "$LOGIN_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('token',{}).get('access_token',''))")
if [[ -z "$TOKEN" ]]; then
  echo "Login failed. Response: $LOGIN_RESP"
  kill $API_PID $ASR_PID 2>/dev/null || true
  exit 1
fi
echo "Token obtained."

# 5. Upload video
echo "[5/7] Uploading sample video..."
UPLOAD_RESP=$(curl -s -X POST "$API_URL/api/v1/videos/upload" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"project_id\":\"b0000001-0000-4000-8000-000000000001\",\"filename\":\"sample.mp4\"}")
UPLOAD_URL=$(echo "$UPLOAD_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); v=d.get('video',{}); u=d.get('upload',{}); print(u.get('upload_url',''))")
VIDEO_ID=$(echo "$UPLOAD_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); v=d.get('video',{}); print(v.get('id',''))")
if [[ -z "$UPLOAD_URL" || -z "$VIDEO_ID" ]]; then
  echo "Upload init failed. Response: $UPLOAD_RESP"
  kill $API_PID $ASR_PID 2>/dev/null || true
  exit 1
fi
# PUT file to presigned URL (escape URL for curl)
curl -s -X PUT -T "$SAMPLE_VIDEO" -H "Content-Type: video/mp4" "$UPLOAD_URL" >/dev/null
curl -s -X POST "$API_URL/api/v1/videos/$VIDEO_ID/confirm" -H "Authorization: Bearer $TOKEN" >/dev/null
echo "Video ID: $VIDEO_ID"

# 6. Start transcription and poll
echo "[6/7] Starting transcription and polling..."
TR_RESP=$(curl -s -X POST "$API_URL/api/v1/transcriptions/videos/$VIDEO_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"language":"en"}')
TR_ID=$(echo "$TR_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); t=d.get('transcription',{}); print(t.get('id',''))")
if [[ -z "$TR_ID" ]]; then
  echo "Start transcription failed. Response: $TR_RESP"
  kill $API_PID $ASR_PID 2>/dev/null || true
  exit 1
fi
echo "Transcription ID: $TR_ID"

echo "Polling for completion (may take a few minutes)..."
for i in $(seq 1 60); do
  STATUS_RESP=$(curl -s "$API_URL/api/v1/transcriptions/$TR_ID" -H "Authorization: Bearer $TOKEN")
  STATUS=$(echo "$STATUS_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); t=d.get('transcription',{}); print(t.get('status',''))")
  echo "  [$i] status: $STATUS"
  if [[ "$STATUS" == "completed" ]]; then
    echo ""
    echo "=== Transcription completed ==="
    echo "$STATUS_RESP" | python3 -c "
import sys, json
d = json.load(sys.stdin)
t = d.get('transcription') or {}
for s in t.get('segments', []):
    print(\"  {:.1f}s - {:.1f}s: {}\".format(s.get('start_time',0), s.get('end_time',0), (s.get('text') or '')[:60]))
" 2>/dev/null || echo "$STATUS_RESP"
    break
  fi
  if [[ "$STATUS" == "failed" ]]; then
    echo "Transcription failed."
    echo "$STATUS_RESP"
    kill $API_PID $ASR_PID 2>/dev/null || true
    exit 1
  fi
  sleep 10
done

if [[ "$STATUS" != "completed" ]]; then
  echo "Timeout waiting for completion (last status: $STATUS)"
fi

echo ""
echo "Done. API PID=$API_PID, ASR PID=$ASR_PID (leave running or kill them)."
