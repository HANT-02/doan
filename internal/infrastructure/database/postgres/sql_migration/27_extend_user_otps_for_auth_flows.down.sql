DROP INDEX IF EXISTS idx_user_otps_user_purpose;

ALTER TABLE user_otps
    DROP COLUMN IF EXISTS pending_password_hash,
    DROP COLUMN IF EXISTS purpose;
