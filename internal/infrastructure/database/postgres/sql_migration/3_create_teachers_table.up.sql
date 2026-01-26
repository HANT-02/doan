CREATE TABLE IF NOT EXISTS teachers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE,
    full_name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(20),
    is_school_teacher BOOLEAN DEFAULT FALSE,
    school_name VARCHAR(255),
    employment_type VARCHAR(50) DEFAULT 'PART_TIME',
    status VARCHAR(50) DEFAULT 'ACTIVE',
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);