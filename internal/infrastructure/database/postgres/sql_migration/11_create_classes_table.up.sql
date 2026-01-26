CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    notes TEXT,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    max_students INTEGER,
    status VARCHAR(50) DEFAULT 'OPEN',
    price NUMERIC(10,2),
    program_id UUID REFERENCES programs(id),
    course_id UUID REFERENCES courses(id),
    teacher_id UUID REFERENCES teachers(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);