CREATE TABLE IF NOT EXISTS outcomes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    program_id UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    objective_id UUID REFERENCES objectives(id) ON DELETE SET NULL
);