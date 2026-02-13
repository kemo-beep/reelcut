-- Multi-language: link translated transcriptions to source and optional caption language per clip style
ALTER TABLE transcriptions ADD COLUMN IF NOT EXISTS source_transcription_id UUID REFERENCES transcriptions(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_transcriptions_source ON transcriptions(source_transcription_id);

ALTER TABLE clip_styles ADD COLUMN IF NOT EXISTS caption_language VARCHAR(10) DEFAULT NULL;
