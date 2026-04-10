ALTER TABLE class_schedules
ADD COLUMN IF NOT EXISTS start_time VARCHAR(10),
ADD COLUMN IF NOT EXISTS end_time VARCHAR(10);

UPDATE class_schedules AS cs
SET start_time = sh.start_time,
    end_time = sh.end_time
FROM shifts AS sh
WHERE cs.shift_id = sh.id;

DROP INDEX IF EXISTS idx_class_schedules_shift_id;

ALTER TABLE class_schedules
DROP CONSTRAINT IF EXISTS fk_class_schedules_shift;

ALTER TABLE class_schedules
ALTER COLUMN start_time SET NOT NULL,
ALTER COLUMN end_time SET NOT NULL;

ALTER TABLE class_schedules
DROP COLUMN IF EXISTS shift_id;
