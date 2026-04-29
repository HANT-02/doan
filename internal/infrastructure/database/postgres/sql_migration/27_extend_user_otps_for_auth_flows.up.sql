ALTER TABLE user_otps
    ADD COLUMN IF NOT EXISTS purpose VARCHAR(50) NOT NULL DEFAULT 'REGISTER',
    ADD COLUMN IF NOT EXISTS pending_password_hash TEXT NULL;

CREATE INDEX IF NOT EXISTS idx_user_otps_user_purpose ON user_otps(user_id, purpose);
