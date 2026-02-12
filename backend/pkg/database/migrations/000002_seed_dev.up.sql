-- Seed dev data (idempotent: only insert if not present).
-- Optional: run only when SEED_DEV=true by convention (e.g. separate seed command); this migration runs with all migrations.
-- Dev user: dev@reelcut.local / password123 (bcrypt hash below)
INSERT INTO users (id, email, password_hash, full_name, subscription_tier, credits_remaining, email_verified)
VALUES (
  'a0000001-0000-4000-8000-000000000001',
  'dev@reelcut.local',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
  'Dev User',
  'free',
  10,
  true
)
ON CONFLICT (id) DO NOTHING;

-- Dev project
INSERT INTO projects (id, user_id, name, description)
VALUES (
  'b0000001-0000-4000-8000-000000000001',
  'a0000001-0000-4000-8000-000000000001',
  'Seed Project',
  'Dev seed project'
)
ON CONFLICT (id) DO NOTHING;

-- Dev template (public)
INSERT INTO templates (id, user_id, name, category, is_public, style_config, usage_count)
VALUES (
  'c0000001-0000-4000-8000-000000000001',
  'a0000001-0000-4000-8000-000000000001',
  'Seed Template',
  'default',
  true,
  '{"caption_font":"Inter","caption_size":24,"caption_color":"#ffffff"}',
  0
)
ON CONFLICT (id) DO NOTHING;
