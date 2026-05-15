BEGIN;

UPDATE classes
SET
    end_date = start_date + INTERVAL '7 day',
    updated_at = NOW()
WHERE code = 'XLB-NHO-L06';

DELETE FROM class_schedules
WHERE class_id = (
    SELECT id
    FROM classes
    WHERE code = 'XLB-NHO-L06'
    LIMIT 1
)
AND day_of_week <> 'MONDAY';

UPDATE courses
SET
    session_count = 12,
    updated_at = NOW()
WHERE id = (
    SELECT course_id
    FROM classes
    WHERE code = 'XLB-NHO-L06'
    LIMIT 1
);

COMMIT;
