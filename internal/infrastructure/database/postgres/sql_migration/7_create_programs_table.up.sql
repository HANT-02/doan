CREATE TABLE IF NOT EXISTS programs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255),
    track VARCHAR(50), -- SUPPORT, BASIC, ADVANCED
    effective_from TIMESTAMP,
    effective_to TIMESTAMP,
    created_by_id UUID REFERENCES users(id),
    approved_by_id UUID REFERENCES users(id),
    approval_note TEXT,
    published_at TIMESTAMP,
    archived_at TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);