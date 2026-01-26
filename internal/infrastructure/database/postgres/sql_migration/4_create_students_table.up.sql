CREATE TABLE IF NOT EXISTS students (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE,
    full_name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(20),
    guardian_phone VARCHAR(20),
    grade_level VARCHAR(50),
    status VARCHAR(50) DEFAULT 'ACTIVE',
    date_of_birth TIMESTAMP,
    gender VARCHAR(20),
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);