BEGIN;

UPDATE classes
SET
    max_students = 80,
    updated_at = NOW()
WHERE code = 'XLB-NHO-L07';

COMMIT;
