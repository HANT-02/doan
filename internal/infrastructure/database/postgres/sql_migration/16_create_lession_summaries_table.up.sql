CREATE TABLE IF NOT EXISTS lesson_summaries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lesson_id UUID UNIQUE REFERENCES lessons(id) ON DELETE CASCADE,
    topic TEXT,
    lesson_content TEXT,
    class_feedback TEXT,
    homework TEXT,
    homework_deadline TIMESTAMP,
    teacher_notes TEXT,
    created_by_id UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);