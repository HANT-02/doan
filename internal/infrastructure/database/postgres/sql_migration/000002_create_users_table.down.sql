ALTER TABLE users DROP CONSTRAINT IF EXISTS uni_users_username;
-- Thêm lại constraint mới
ALTER TABLE users ADD CONSTRAINT uni_users_username UNIQUE (username);