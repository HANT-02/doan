CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255),
    description TEXT,
    grade_level VARCHAR(50),
    subject VARCHAR(255),
    session_count INTEGER,
    session_duration_minutes INTEGER,
    total_hours NUMERIC(8,2),
    price NUMERIC(10,2),
    status VARCHAR(50) DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);