DELETE FROM campuses WHERE code = 'MAIN';

DROP INDEX IF EXISTS idx_campuses_deleted_at;
DROP TABLE IF EXISTS campuses;
