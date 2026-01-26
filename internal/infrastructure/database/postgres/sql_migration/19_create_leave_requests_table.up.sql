CREATE TABLE IF NOT EXISTS leave_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    leave_type VARCHAR(50) NOT NULL, -- LEAVE, LATE, EARLY
    apply_date TIMESTAMP NOT NULL,
    late_minutes INTEGER,
    early_minutes INTEGER,
    reason TEXT NOT NULL,
    documents TEXT[],
    class_id UUID REFERENCES classes(id) ON DELETE SET NULL,
    lesson_id UUID REFERENCES lessons(id) ON DELETE SET NULL,
    subject VARCHAR(255),
    status VARCHAR(50) DEFAULT 'PENDING',
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    rejection_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);

-- Add Comments
COMMENT ON TABLE users IS 'Table 3.1: User accounts';
COMMENT ON TABLE teachers IS 'Table 3.2: Teacher profiles';
COMMENT ON TABLE students IS 'Table 3.3: Student profiles';
COMMENT ON TABLE classes IS 'Table 3.10: Class management';
