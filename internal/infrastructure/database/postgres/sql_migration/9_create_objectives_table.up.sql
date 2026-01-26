CREATE TABLE IF NOT EXISTS objectives (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    program_id UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE
);