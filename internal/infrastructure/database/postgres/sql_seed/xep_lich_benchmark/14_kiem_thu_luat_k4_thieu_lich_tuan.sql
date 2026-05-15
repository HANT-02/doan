BEGIN;

DELETE FROM class_schedules
WHERE class_id = (
    SELECT id
    FROM classes
    WHERE code = 'XLB-NHO-L05'
    LIMIT 1
);

COMMIT;
