ALTER TABLE lessons
    ADD COLUMN IF NOT EXISTS status VARCHAR(32) NOT NULL DEFAULT 'PUBLISHED',
    ADD COLUMN IF NOT EXISTS published_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS source_preview_run_id VARCHAR(255),
    ADD COLUMN IF NOT EXISTS change_reason TEXT;

UPDATE lessons
SET
    status = CASE
        WHEN date_end < NOW() THEN 'HISTORY'
        ELSE 'PUBLISHED'
    END,
    published_at = COALESCE(published_at, created_at, NOW()),
    change_reason = COALESCE(NULLIF(change_reason, ''), 'LEGACY_BACKFILL')
WHERE status IS NULL
   OR status = ''
   OR published_at IS NULL
   OR change_reason IS NULL
   OR change_reason = '';

CREATE INDEX IF NOT EXISTS idx_lessons_status ON lessons(status);
CREATE INDEX IF NOT EXISTS idx_lessons_published_at ON lessons(published_at);
