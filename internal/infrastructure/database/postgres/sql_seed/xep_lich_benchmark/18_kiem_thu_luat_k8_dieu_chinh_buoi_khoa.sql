BEGIN;

DELETE FROM lessons
WHERE notes = 'GOI_K8_XLB_NHO';

INSERT INTO lessons (
    class_id,
    teacher_id,
    date_start,
    date_end,
    room_id,
    status,
    published_at,
    source_preview_run_id,
    change_reason,
    notes,
    created_at,
    updated_at
)
SELECT
    class_ref.id,
    class_ref.teacher_id,
    seed.date_start,
    seed.date_end,
    room_ref.id,
    'PUBLISHED',
    NOW(),
    'XLB-K8',
    'INITIAL_SCHEDULING_COMMIT',
    'GOI_K8_XLB_NHO',
    NOW(),
    NOW()
FROM (
    VALUES
        ('XLB-NHO-L10', TIMESTAMP '2026-06-08 08:00:00', TIMESTAMP '2026-06-08 10:00:00', 'XLB-NHO-P02'),
        ('XLB-NHO-L11', TIMESTAMP '2026-06-09 13:30:00', TIMESTAMP '2026-06-09 15:30:00', 'XLB-NHO-P04')
) AS seed(class_code, date_start, date_end, room_code)
JOIN classes AS class_ref ON class_ref.code = seed.class_code
JOIN rooms AS room_ref ON room_ref.code = seed.room_code;

COMMIT;
