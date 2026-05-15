BEGIN;

UPDATE classes
SET
    teacher_id = (
        SELECT teacher_id
        FROM classes
        WHERE code = 'XLB-NHO-L08'
        LIMIT 1
    ),
    updated_at = NOW()
WHERE code = 'XLB-NHO-L09';

DELETE FROM class_schedules
WHERE class_id IN (
    SELECT id
    FROM classes
    WHERE code IN ('XLB-NHO-L08', 'XLB-NHO-L09')
);

INSERT INTO class_schedules (class_id, day_of_week, shift_id, room_id)
SELECT
    class_ref.id,
    seed.day_of_week,
    shift_ref.id,
    room_ref.id
FROM (
    VALUES
        ('XLB-NHO-L08', 'MONDAY', 'XLB-CHUNG-S1', 'XLB-NHO-P01'),
        ('XLB-NHO-L08', 'WEDNESDAY', 'XLB-CHUNG-S3', 'XLB-NHO-P01'),
        ('XLB-NHO-L09', 'MONDAY', 'XLB-CHUNG-S2', 'XLB-NHO-P05'),
        ('XLB-NHO-L09', 'WEDNESDAY', 'XLB-CHUNG-S4', 'XLB-NHO-P05')
) AS seed(class_code, day_of_week, shift_code, room_code)
JOIN classes AS class_ref ON class_ref.code = seed.class_code
JOIN shifts AS shift_ref ON shift_ref.code = seed.shift_code
JOIN rooms AS room_ref ON room_ref.code = seed.room_code;

UPDATE campus_travel_times
SET
    travel_minutes = 150,
    updated_at = NOW()
WHERE from_campus_id = (SELECT id FROM campuses WHERE code = 'XLB-CHUNG-CS1')
  AND to_campus_id = (SELECT id FROM campuses WHERE code = 'XLB-CHUNG-CS2');

UPDATE campus_travel_times
SET
    travel_minutes = 150,
    updated_at = NOW()
WHERE from_campus_id = (SELECT id FROM campuses WHERE code = 'XLB-CHUNG-CS2')
  AND to_campus_id = (SELECT id FROM campuses WHERE code = 'XLB-CHUNG-CS1');

COMMIT;
