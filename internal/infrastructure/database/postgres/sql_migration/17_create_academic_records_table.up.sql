CREATE TABLE IF NOT EXISTS academic_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lesson_summary_id UUID REFERENCES lesson_summaries(id) ON DELETE CASCADE,
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    homework_completed BOOLEAN DEFAULT FALSE,
    homework_score NUMERIC(5,2),
    attitude_rating INTEGER,
    participation_score NUMERIC(5,2),
    personal_comment TEXT,
    total_score NUMERIC(5,2),
    is_completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);