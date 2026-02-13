CREATE TABLE broll_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
    original_filename VARCHAR(255) NOT NULL,
    storage_path TEXT NOT NULL,
    duration_seconds DECIMAL(10, 2),
    width INT,
    height INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_broll_assets_user ON broll_assets(user_id);
CREATE INDEX idx_broll_assets_project ON broll_assets(project_id);

CREATE TABLE clip_broll_segments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clip_id UUID NOT NULL REFERENCES clips(id) ON DELETE CASCADE,
    broll_asset_id UUID NOT NULL REFERENCES broll_assets(id) ON DELETE CASCADE,
    start_time DECIMAL(10, 3) NOT NULL,
    end_time DECIMAL(10, 3) NOT NULL,
    position VARCHAR(20) NOT NULL DEFAULT 'cut_in',
    scale DECIMAL(3, 2) DEFAULT 1.0,
    opacity DECIMAL(3, 2) DEFAULT 1.0,
    sequence_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT check_broll_time CHECK (end_time > start_time),
    UNIQUE(clip_id, sequence_order)
);
CREATE INDEX idx_clip_broll_segments_clip ON clip_broll_segments(clip_id);
CREATE INDEX idx_clip_broll_segments_asset ON clip_broll_segments(broll_asset_id);
