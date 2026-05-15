BEGIN;

UPDATE classes
SET
    teacher_id = NULL,
    updated_at = NOW()
WHERE code = 'XLB-NHO-L03';

COMMIT;
