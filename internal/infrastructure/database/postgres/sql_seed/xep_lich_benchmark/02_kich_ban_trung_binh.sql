BEGIN;

DELETE FROM enrollments
WHERE class_id IN (SELECT id FROM classes WHERE code LIKE 'XLB-TB-%')
   OR student_id IN (SELECT id FROM students WHERE code LIKE 'XLB-TB-%');

DELETE FROM lessons
WHERE class_id IN (SELECT id FROM classes WHERE code LIKE 'XLB-TB-%');

DELETE FROM class_schedules
WHERE class_id IN (SELECT id FROM classes WHERE code LIKE 'XLB-TB-%');

DELETE FROM classes WHERE code LIKE 'XLB-TB-%';
DELETE FROM courses WHERE code LIKE 'XLB-TB-%';
DELETE FROM teachers WHERE code LIKE 'XLB-TB-%';
DELETE FROM students WHERE code LIKE 'XLB-TB-%';
DELETE FROM rooms WHERE code LIKE 'XLB-TB-%';

INSERT INTO rooms (code, name, capacity, address, campus_id, created_at, updated_at)
SELECT
    FORMAT('XLB-TB-P%s', LPAD(seq::text, 2, '0')),
    FORMAT('Phòng benchmark trung bình %s', seq),
    CASE (seq - 1) % 6
        WHEN 0 THEN 18
        WHEN 1 THEN 20
        WHEN 2 THEN 24
        WHEN 3 THEN 30
        WHEN 4 THEN 36
        ELSE 45
    END,
    FORMAT('Cụm phòng benchmark trung bình %s', seq),
    CASE
        WHEN seq <= 8 THEN (SELECT id FROM campuses WHERE code = 'XLB-CHUNG-CS1')
        WHEN seq <= 14 THEN (SELECT id FROM campuses WHERE code = 'XLB-CHUNG-CS2')
        ELSE (SELECT id FROM campuses WHERE code = 'XLB-CHUNG-CS3')
    END,
    NOW(),
    NOW()
FROM generate_series(1, 18) AS seq;

INSERT INTO teachers (
    code,
    full_name,
    email,
    phone,
    is_school_teacher,
    school_name,
    employment_type,
    status,
    skills,
    notes,
    created_at,
    updated_at
)
SELECT
    FORMAT('XLB-TB-GV%s', LPAD(seq::text, 3, '0')),
    FORMAT('Giáo viên benchmark trung bình %s', seq),
    FORMAT('xlb.tb.gv%s@example.com', LPAD(seq::text, 3, '0')),
    FORMAT('092500%s', LPAD(seq::text, 4, '0')),
    FALSE,
    NULL,
    CASE WHEN seq <= 18 THEN 'FULL_TIME' ELSE 'PART_TIME' END,
    'ACTIVE',
    ARRAY['TOAN', 'ANH', 'LY', 'HOA', 'VAN', 'TIN']::TEXT[],
    'Bộ giáo viên chuẩn cho benchmark trung bình',
    NOW(),
    NOW()
FROM generate_series(1, 30) AS seq;

INSERT INTO courses (
    code,
    name,
    description,
    grade_level,
    subject,
    session_count,
    session_duration_minutes,
    total_hours,
    price,
    status,
    required_skills,
    created_at,
    updated_at
)
SELECT
    FORMAT('XLB-TB-KH%s', LPAD(seq::text, 3, '0')),
    FORMAT('Khóa benchmark trung bình %s', seq),
    'Khóa học dùng cho đo đạc đối sánh xếp lịch quy mô trung bình',
    CASE (seq - 1) % 5
        WHEN 0 THEN 'Khối 8'
        WHEN 1 THEN 'Khối 9'
        WHEN 2 THEN 'Khối 10'
        WHEN 3 THEN 'Khối 11'
        ELSE 'Khối 12'
    END,
    CASE (seq - 1) % 6
        WHEN 0 THEN 'Toán'
        WHEN 1 THEN 'Tiếng Anh'
        WHEN 2 THEN 'Vật Lý'
        WHEN 3 THEN 'Hóa Học'
        WHEN 4 THEN 'Ngữ Văn'
        ELSE 'Tin Học'
    END,
    6 + (seq % 3),
    120,
    14 + (seq % 3) * 2,
    1900000 + seq * 20000,
    'ACTIVE',
    ARRAY[
        CASE (seq - 1) % 6
            WHEN 0 THEN 'TOAN'
            WHEN 1 THEN 'ANH'
            WHEN 2 THEN 'LY'
            WHEN 3 THEN 'HOA'
            WHEN 4 THEN 'VAN'
            ELSE 'TIN'
        END
    ]::TEXT[],
    NOW(),
    NOW()
FROM generate_series(1, 72) AS seq;

INSERT INTO students (
    code,
    full_name,
    email,
    phone,
    guardian_phone,
    grade_level,
    status,
    gender,
    address,
    created_at,
    updated_at
)
SELECT
    FORMAT('XLB-TB-HS%s', LPAD(seq::text, 4, '0')),
    FORMAT('Học sinh benchmark trung bình %s', seq),
    FORMAT('xlb.tb.hs%s@example.com', LPAD(seq::text, 4, '0')),
    FORMAT('083600%s', LPAD(seq::text, 4, '0')),
    FORMAT('093600%s', LPAD(seq::text, 4, '0')),
    CASE (seq - 1) % 5
        WHEN 0 THEN 'Khối 8'
        WHEN 1 THEN 'Khối 9'
        WHEN 2 THEN 'Khối 10'
        WHEN 3 THEN 'Khối 11'
        ELSE 'Khối 12'
    END,
    'ACTIVE',
    CASE WHEN seq % 2 = 0 THEN 'FEMALE' ELSE 'MALE' END,
    'Cụm dân cư benchmark trung bình',
    NOW(),
    NOW()
FROM generate_series(1, 720) AS seq;

INSERT INTO classes (
    code,
    name,
    notes,
    start_date,
    end_date,
    max_students,
    status,
    price,
    program_id,
    course_id,
    teacher_id,
    room_id,
    created_at,
    updated_at
)
SELECT
    FORMAT('XLB-TB-L%s', LPAD(seq::text, 3, '0')),
    FORMAT('Lớp benchmark trung bình %s', seq),
    'Lớp chuẩn cho benchmark xếp lịch quy mô trung bình',
    DATE '2026-06-01' + ((seq - 1) % 7) * INTERVAL '1 day',
    DATE '2026-09-15',
    CASE (seq - 1) % 6
        WHEN 0 THEN 18
        WHEN 1 THEN 20
        WHEN 2 THEN 24
        WHEN 3 THEN 28
        WHEN 4 THEN 32
        ELSE 36
    END,
    'OPEN',
    2500000,
    NULL,
    course_ref.id,
    teacher_ref.id,
    room_ref.id,
    NOW(),
    NOW()
FROM generate_series(1, 72) AS seq
JOIN courses AS course_ref ON course_ref.code = FORMAT('XLB-TB-KH%s', LPAD(seq::text, 3, '0'))
JOIN teachers AS teacher_ref ON teacher_ref.code = FORMAT('XLB-TB-GV%s', LPAD((((seq - 1) % 30) + 1)::text, 3, '0'))
JOIN rooms AS room_ref ON room_ref.code = FORMAT('XLB-TB-P%s', LPAD((((seq - 1) % 18) + 1)::text, 2, '0'));

INSERT INTO class_schedules (class_id, day_of_week, shift_id, room_id)
SELECT
    class_ref.id,
    schedule_data.day_of_week,
    shift_ref.id,
    room_ref.id
FROM classes AS class_ref
JOIN LATERAL (
    VALUES
        (
            CASE (RIGHT(class_ref.code, 3)::INT - 1) % 5
                WHEN 0 THEN 'MONDAY'
                WHEN 1 THEN 'TUESDAY'
                WHEN 2 THEN 'WEDNESDAY'
                WHEN 3 THEN 'THURSDAY'
                ELSE 'FRIDAY'
            END,
            CASE (RIGHT(class_ref.code, 3)::INT - 1) % 6
                WHEN 0 THEN 'XLB-CHUNG-S1'
                WHEN 1 THEN 'XLB-CHUNG-S2'
                WHEN 2 THEN 'XLB-CHUNG-S3'
                WHEN 3 THEN 'XLB-CHUNG-S4'
                WHEN 4 THEN 'XLB-CHUNG-S5'
                ELSE 'XLB-CHUNG-S6'
            END,
            FORMAT('XLB-TB-P%s', LPAD((((RIGHT(class_ref.code, 3)::INT - 1) % 18) + 1)::text, 2, '0'))
        ),
        (
            CASE (RIGHT(class_ref.code, 3)::INT - 1) % 5
                WHEN 0 THEN 'WEDNESDAY'
                WHEN 1 THEN 'THURSDAY'
                WHEN 2 THEN 'FRIDAY'
                WHEN 3 THEN 'SATURDAY'
                ELSE 'SUNDAY'
            END,
            CASE (RIGHT(class_ref.code, 3)::INT) % 6
                WHEN 0 THEN 'XLB-CHUNG-S1'
                WHEN 1 THEN 'XLB-CHUNG-S2'
                WHEN 2 THEN 'XLB-CHUNG-S3'
                WHEN 3 THEN 'XLB-CHUNG-S4'
                WHEN 4 THEN 'XLB-CHUNG-S5'
                ELSE 'XLB-CHUNG-S6'
            END,
            FORMAT('XLB-TB-P%s', LPAD((((RIGHT(class_ref.code, 3)::INT) % 18) + 1)::text, 2, '0'))
        )
) AS schedule_data(day_of_week, shift_code, room_code) ON TRUE
JOIN shifts AS shift_ref ON shift_ref.code = schedule_data.shift_code
JOIN rooms AS room_ref ON room_ref.code = schedule_data.room_code
WHERE class_ref.code LIKE 'XLB-TB-%';

INSERT INTO enrollments (
    class_id,
    student_id,
    status,
    approved_at,
    created_at,
    updated_at
)
SELECT
    class_ref.id,
    student_ref.id,
    'APPROVED',
    NOW(),
    NOW(),
    NOW()
FROM generate_series(1, 72) AS class_seq
JOIN classes AS class_ref ON class_ref.code = FORMAT('XLB-TB-L%s', LPAD(class_seq::text, 3, '0'))
JOIN LATERAL generate_series(1, 10 + (class_seq % 8)) AS enrollment_seq(seq) ON TRUE
JOIN students AS student_ref
    ON student_ref.code = FORMAT(
        'XLB-TB-HS%s',
        LPAD((((class_seq - 1) * 11 + enrollment_seq.seq - 1) % 720 + 1)::text, 4, '0')
    );

COMMIT;
