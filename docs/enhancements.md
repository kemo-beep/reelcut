Future Enhancements (implementation status)

1. Auto-caption generation with styling — Implemented. ASS generation in caption_service; BurnASS in ffmpeg; rendering uses clip style (font, size, color, bg, position, animation). GET /api/v1/config/caption-fonts. StylePanel + CaptionEditor in frontend.

2. Thumbnail extraction and optimization — Implemented. Video thumbnail at 1s offset; clip thumbnail job (TypeClipThumbnail) enqueued on clip create and when clip times change. thumbnails/clips/{clipId}.jpg. ClipCard and editor use thumbnail_url.

3. Multi-language support — Implemented. Migration 000003: source_transcription_id, caption_language on clip_styles. POST /transcriptions/:id/translate (Gemini); GET /transcriptions/videos/:videoId?language=; GET /transcriptions/videos/:videoId/list. Rendering and SRT/VTT use caption_language. Frontend: language on create, Translate to…, caption language in StylePanel.

4. Platform-specific optimization (9:16, 1:1, etc.) — Implemented. Export presets (tiktok, reels, instagram_feed, youtube_shorts) in config; GET /api/v1/config/export-presets. Render accepts preset; ResizeCropToSize for exact dimensions. ExportPanel preset dropdown.

5. B-roll insertion at strategic points — Implemented. Migration 000004: broll_assets, clip_broll_segments. POST/GET /api/v1/broll/assets; GET/POST/DELETE /api/v1/clips/:id/broll. Overlay in render (OverlayVideo). Editor BrollPanel: list segments, add (asset, start, end, position), delete.

