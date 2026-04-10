ALTER TABLE class_schedules
ADD COLUMN IF NOT EXISTS shift_id UUID;

INSERT INTO shifts (
    id,
    code,
    name,
    start_time,
    end_time,
    duration_minutes,
    session_type,
    is_active,
    notes,
    created_at,
    updated_at
)
SELECT
    uuid_generate_v4(),
    CONCAT('LEGACY-', REPLACE(legacy.start_time, ':', ''), '-', REPLACE(legacy.end_time, ':', '')),
    CONCAT('Ca ', legacy.start_time, ' - ', legacy.end_time),
    legacy.start_time,
    legacy.end_time,
    GREATEST(
        1,
        (
            EXTRACT(
                EPOCH FROM (
                    legacy.end_time::time - legacy.start_time::time
                )
            ) / 60
        )::INTEGER
    ),
    'CUSTOM',
    TRUE,
    'Migrated from legacy class_schedules.start_time/end_time',
    NOW(),
    NOW()
FROM (
    SELECT DISTINCT start_time, end_time
    FROM class_schedules
    WHERE start_time IS NOT NULL
      AND end_time IS NOT NULL
      AND TRIM(start_time) <> ''
      AND TRIM(end_time) <> ''
) AS legacy
ON CONFLICT (code) DO NOTHING;

UPDATE class_schedules AS cs
SET shift_id = sh.id
FROM shifts AS sh
WHERE sh.code = CONCAT('LEGACY-', REPLACE(cs.start_time, ':', ''), '-', REPLACE(cs.end_time, ':', ''))
  AND cs.shift_id IS NULL;

ALTER TABLE class_schedules
ALTER COLUMN shift_id SET NOT NULL;

ALTER TABLE class_schedules
ADD CONSTRAINT fk_class_schedules_shift
FOREIGN KEY (shift_id) REFERENCES shifts(id);

CREATE INDEX IF NOT EXISTS idx_class_schedules_shift_id ON class_schedules(shift_id);

ALTER TABLE class_schedules
DROP COLUMN IF EXISTS start_time,
DROP COLUMN IF EXISTS end_time;
