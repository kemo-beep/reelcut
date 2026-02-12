"""
WhisperLiveKit ASR service for Reelcut.
Exposes WebSocket /asr: client sends binary audio; server streams JSON results and sends {"type": "ready_to_stop"} when done.
"""
import asyncio
import logging
import os
from contextlib import asynccontextmanager

from fastapi import FastAPI, WebSocket, WebSocketDisconnect
from fastapi.middleware.cors import CORSMiddleware

from whisperlivekit import AudioProcessor, TranscriptionEngine

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)

transcription_engine = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    global transcription_engine
    model = os.getenv("WLK_MODEL", "medium")
    diarization = os.getenv("WLK_DIARIZATION", "true").lower() in ("1", "true", "yes")
    language = os.getenv("WLK_LANGUAGE", "en")
    transcription_engine = TranscriptionEngine(
        model=model,
        diarization=diarization,
        lan=language,
    )
    yield


app = FastAPI(lifespan=lifespan)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


async def handle_websocket_results(websocket: WebSocket, results_generator):
    try:
        async for response in results_generator:
            # response may have .to_dict() in some versions
            payload = response.to_dict() if hasattr(response, "to_dict") else response
            await websocket.send_json(payload)
    except WebSocketDisconnect:
        logger.info("Client disconnected while sending results")
    await websocket.send_json({"type": "ready_to_stop"})


@app.websocket("/asr")
async def websocket_endpoint(websocket: WebSocket):
    global transcription_engine

    audio_processor = AudioProcessor(transcription_engine=transcription_engine)
    results_generator = await audio_processor.create_tasks()
    results_task = asyncio.create_task(
        handle_websocket_results(websocket, results_generator)
    )
    await websocket.accept()

    try:
        while True:
            message = await websocket.receive_bytes()
            await audio_processor.process_audio(message)
    except WebSocketDisconnect:
        logger.info("WebSocket disconnected by client")
    except Exception as e:
        logger.warning("WebSocket receive error: %s", e)
    finally:
        if not results_task.done():
            results_task.cancel()
        try:
            await results_task
        except asyncio.CancelledError:
            pass
        await audio_processor.cleanup()


@app.get("/health")
async def health():
    return {"status": "ok"}


@app.get(
    "/asr",
    summary="ASR WebSocket endpoint (documentation)",
    description="Transcription is done via WebSocket, not HTTP. Connect to ws://host:port/asr and send binary audio chunks (WAV). Server streams JSON segments and sends {\"type\": \"ready_to_stop\"} when done. This GET route only documents the endpoint for /docs.",
)
async def asr_info():
    """Return service info so /docs shows that the WebSocket exists."""
    return {
        "service": "WhisperLiveKit ASR",
        "websocket_url": "/asr",
        "usage": "Connect to ws://<this_host>/asr and send binary audio; server responds with JSON (lines, type: ready_to_stop).",
    }
