DROP INDEX IF EXISTS idx_lessons_published_at;
DROP INDEX IF EXISTS idx_lessons_status;

ALTER TABLE lessons
    DROP COLUMN IF EXISTS change_reason,
    DROP COLUMN IF EXISTS source_preview_run_id,
    DROP COLUMN IF EXISTS published_at,
    DROP COLUMN IF EXISTS status;
