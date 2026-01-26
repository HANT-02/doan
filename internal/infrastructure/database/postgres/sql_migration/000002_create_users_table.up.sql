-- Create users table
CREATE TABLE IF NOT EXISTS users
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username   VARCHAR(255) NOT NULL UNIQUE,
    password   TEXT NOT NULL,
    status     VARCHAR(255) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_user_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_user_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_user_created_at ON users(created_at);

-- Add comments
COMMENT ON TABLE users IS 'User accounts table';
COMMENT ON COLUMN users.id IS 'User unique identifier (UUID)';
COMMENT ON COLUMN users.username IS 'Unique username for login';
COMMENT ON COLUMN users.password IS 'Hashed password';
COMMENT ON COLUMN users.status IS 'User status: active, inactive, suspended';

