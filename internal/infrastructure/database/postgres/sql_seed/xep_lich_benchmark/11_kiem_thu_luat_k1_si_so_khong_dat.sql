BEGIN;

DELETE FROM enrollments
WHERE class_id IN (
    SELECT id
    FROM classes
    WHERE code IN ('XLB-NHO-L01', 'XLB-NHO-L02')
);

INSERT INTO enrollments (class_id, student_id, status, approved_at, created_at, updated_at)
SELECT
    class_ref.id,
    student_ref.id,
    'APPROVED',
    NOW(),
    NOW(),
    NOW()
FROM (
    VALUES
        ('XLB-NHO-L01', 'XLB-NHO-HS001'),
        ('XLB-NHO-L01', 'XLB-NHO-HS002'),
        ('XLB-NHO-L02', 'XLB-NHO-HS003'),
        ('XLB-NHO-L02', 'XLB-NHO-HS004')
) AS seed(class_code, student_code)
JOIN classes AS class_ref ON class_ref.code = seed.class_code
JOIN students AS student_ref ON student_ref.code = seed.student_code;

COMMIT;
