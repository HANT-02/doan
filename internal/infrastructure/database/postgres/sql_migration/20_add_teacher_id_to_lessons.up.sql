-- Add teacher_id to lessons table for direct teacher-lesson relationship
ALTER TABLE lessons ADD COLUMN teacher_id UUID REFERENCES teachers(id);

-- Create index for performance
CREATE INDEX idx_lessons_teacher_id ON lessons(teacher_id);

-- Update existing lessons to use class.teacher_id
UPDATE lessons l 
SET teacher_id = c.teacher_id 
FROM classes c 
WHERE l.class_id = c.id AND c.teacher_id IS NOT NULL;
