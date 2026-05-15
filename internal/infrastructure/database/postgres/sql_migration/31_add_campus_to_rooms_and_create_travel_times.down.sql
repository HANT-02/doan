DROP INDEX IF EXISTS idx_campus_travel_times_from_to;
DROP INDEX IF EXISTS idx_campus_travel_times_deleted_at;
DROP TABLE IF EXISTS campus_travel_times;

DROP INDEX IF EXISTS idx_rooms_campus_id;
ALTER TABLE rooms DROP COLUMN IF EXISTS campus_id;
