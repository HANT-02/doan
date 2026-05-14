CREATE TABLE IF NOT EXISTS campuses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_campuses_deleted_at ON campuses(deleted_at);

INSERT INTO campuses (code, name, address)
SELECT 'MAIN', 'Main Campus', 'Auto-generated default campus for legacy room data'
WHERE NOT EXISTS (
    SELECT 1 FROM campuses WHERE code = 'MAIN'
);
