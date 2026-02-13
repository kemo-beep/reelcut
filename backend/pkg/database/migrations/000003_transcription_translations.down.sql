ALTER TABLE clip_styles DROP COLUMN IF EXISTS caption_language;
DROP INDEX IF EXISTS idx_transcriptions_source;
ALTER TABLE transcriptions DROP COLUMN IF EXISTS source_transcription_id;
