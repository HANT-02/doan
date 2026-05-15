BEGIN;

UPDATE courses
SET
    required_skills = ARRAY['KY_NANG_KHONG_TON_TAI']::TEXT[],
    updated_at = NOW()
WHERE id = (
    SELECT course_id
    FROM classes
    WHERE code = 'XLB-NHO-L04'
    LIMIT 1
);

COMMIT;
