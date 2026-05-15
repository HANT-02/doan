CREATE TABLE IF NOT EXISTS skills (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(120) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS teacher_skills (
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (teacher_id, skill_id)
);

CREATE TABLE IF NOT EXISTS course_required_skills (
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (course_id, skill_id)
);

INSERT INTO skills (code, name, status)
SELECT DISTINCT
    normalized.code,
    REPLACE(normalized.code, '_', ' '),
    'ACTIVE'
FROM (
    SELECT UPPER(TRIM(skill_code)) AS code
    FROM teachers, UNNEST(skills) AS skill_code
    UNION
    SELECT UPPER(TRIM(skill_code)) AS code
    FROM courses, UNNEST(required_skills) AS skill_code
) AS normalized
WHERE normalized.code <> ''
ON CONFLICT (code) DO NOTHING;

INSERT INTO teacher_skills (teacher_id, skill_id)
SELECT DISTINCT
    teachers.id,
    skills.id
FROM teachers
JOIN UNNEST(teachers.skills) AS skill_code ON TRUE
JOIN skills ON skills.code = UPPER(TRIM(skill_code))
ON CONFLICT DO NOTHING;

INSERT INTO course_required_skills (course_id, skill_id)
SELECT DISTINCT
    courses.id,
    skills.id
FROM courses
JOIN UNNEST(courses.required_skills) AS skill_code ON TRUE
JOIN skills ON skills.code = UPPER(TRIM(skill_code))
ON CONFLICT DO NOTHING;
