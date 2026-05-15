BEGIN;

DELETE FROM lessons
WHERE notes = 'GOI_K9_XLB_NHO';

INSERT INTO lessons (
    class_id,
    teacher_id,
    date_start,
    date_end,
    room_id,
    status,
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
    'DRAFT',
    'XLB-K9',
    'REPLAN_REPLACEMENT',
    'GOI_K9_XLB_NHO',
    NOW(),
    NOW()
FROM (
    VALUES
        ('XLB-NHO-L12', TIMESTAMP '2026-06-10 18:00:00', TIMESTAMP '2026-06-10 20:00:00', 'XLB-NHO-P05')
) AS seed(class_code, date_start, date_end, room_code)
JOIN classes AS class_ref ON class_ref.code = seed.class_code
JOIN rooms AS room_ref ON room_ref.code = seed.room_code;

COMMIT;
