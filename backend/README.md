# Reelcut Backend

Go/Gin backend for the Reelcut platform: AI-powered clip generation from long-form video.

## Stack

- **Go** 1.21+
- **Gin** – HTTP API
- **PostgreSQL** – primary database
- **Redis** – cache, rate limiting, job queue (Asynq)
- **MinIO** – object storage (S3-compatible; videos, thumbnails, avatars, rendered clips). Default in docker-compose; create bucket `reelcut` via MinIO console at http://localhost:9001.
- **FFmpeg** – video metadata, thumbnails, audio extraction (and rendering when implemented)
- **OpenAI Whisper** – transcription (optional, set `OPENAI_API_KEY`)

## Quick start

### 1. Start dependencies

```bash
docker-compose up -d
```

This runs PostgreSQL (5432), Redis (6379), and MinIO (9000/9001). **Object storage is MinIO only** (no AWS S3 required). Create the bucket `reelcut` in MinIO before first upload (MinIO console: http://localhost:9001).

### 2. Environment

```bash
cp .env.example .env
# Edit .env: DATABASE_URL, REDIS_URL, JWT_SECRET, S3_* (MinIO defaults in .env.example)
```

### 3. Run migrations and server

Migrations run on startup (including optional seed migration `000002_seed_dev` which inserts a dev user, project, and template when present). Start the API (and embedded Asynq worker):

```bash
go run ./cmd/api
```

Server listens on `PORT` (default 8080). Health:

```bash
curl http://localhost:8080/health
```

### 4. Development with hot reload

```bash
air
```

(Uses `.air.toml` in the repo root.)

## API overview

- **Auth**: `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `POST /api/v1/auth/refresh`
- **Users**: `GET/PUT /api/v1/users/me`, `GET /api/v1/users/me/usage`
- **Projects**: `GET/POST /api/v1/projects`, `GET/PUT/DELETE /api/v1/projects/:id`
- **Videos**: `POST /api/v1/videos/upload` (returns presigned URL), `POST /api/v1/videos/:id/confirm`, `GET /api/v1/videos`, `GET/DELETE /api/v1/videos/:id`
- **Transcriptions**: `POST /api/v1/transcriptions/videos/:videoId`, `GET /api/v1/transcriptions/:id`, `GET /api/v1/transcriptions/videos/:videoId`
- **Analysis**: `POST /api/v1/analysis/videos/:videoId`, `GET /api/v1/analysis/videos/:videoId`, `POST /api/v1/analysis/videos/:videoId/suggest-clips`
- **Clips**: CRUD at `/api/v1/clips`, style at `/api/v1/clips/:clipId/style`
- **Templates**: CRUD at `/api/v1/templates`, public list at `GET /api/v1/templates/public`
- **Jobs**: `GET /api/v1/jobs`, `GET /api/v1/jobs/:id`, `POST /api/v1/jobs/:id/cancel`

Protected routes require `Authorization: Bearer <access_token>`.

## Swagger

Interactive API docs are served at:

- **Swagger UI**: http://localhost:8080/swagger/index.html

When using **Air** for development, Swagger docs are regenerated automatically before each build (via `pre_cmd` in `.air.toml`). Ensure the swag CLI is installed:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

To regenerate docs manually (e.g. when not using Air):

```bash
swag init -g cmd/api/main.go -o docs --parseInternal
```

## Project layout

- `cmd/api` – entrypoint, wiring, HTTP server and Asynq worker
- `api/` – route registration
- `internal/config` – env-based config
- `internal/domain` – domain models and errors
- `internal/repository` – repository interfaces and Postgres implementations
- `internal/service` – business logic (auth, video, transcription, analysis, clip, template, storage)
- `internal/handler` – HTTP handlers
- `internal/middleware` – auth, CORS, logger, rate limit, error handler
- `internal/worker` – Asynq handlers (video metadata/thumbnail, transcription, analysis)
- `internal/queue` – Asynq task types and queue client
- `internal/ai` – Whisper client, clip suggestion
- `internal/video` – FFmpeg helpers (metadata, frame extraction, audio)
- `pkg/database` – Postgres pool and migrations
- `pkg/redis` – Redis client
- `pkg/logger` – structured logger

## Video upload flow

1. `POST /api/v1/videos/upload` with `{ "project_id", "filename" }` → returns `upload_url` and `video.id`.
2. Client `PUT` the file to `upload_url` (presigned S3).
3. `POST /api/v1/videos/:id/confirm` → enqueues metadata and thumbnail jobs.
4. Workers update the video record (duration, dimensions, thumbnail path, status `ready`).

## Requirements

- FFmpeg and ffprobe on `PATH` for video processing and thumbnail jobs.
- OpenAI API key in `OPENAI_API_KEY` for transcription (otherwise transcription jobs complete with no segments).
