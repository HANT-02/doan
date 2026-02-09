-- Rollback: Remove teacher_id from lessons table
DROP INDEX IF EXISTS idx_lessons_teacher_id;
ALTER TABLE lessons DROP COLUMN IF EXISTS teacher_id;
