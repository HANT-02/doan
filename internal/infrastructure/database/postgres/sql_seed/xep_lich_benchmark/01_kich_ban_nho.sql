BEGIN;

DELETE FROM enrollments
WHERE class_id IN (SELECT id FROM classes WHERE code LIKE 'XLB-NHO-%')
   OR student_id IN (SELECT id FROM students WHERE code LIKE 'XLB-NHO-%');

DELETE FROM lessons
WHERE class_id IN (SELECT id FROM classes WHERE code LIKE 'XLB-NHO-%');

DELETE FROM class_schedules
WHERE class_id IN (SELECT id FROM classes WHERE code LIKE 'XLB-NHO-%');

DELETE FROM classes WHERE code LIKE 'XLB-NHO-%';
DELETE FROM courses WHERE code LIKE 'XLB-NHO-%';
DELETE FROM teachers WHERE code LIKE 'XLB-NHO-%';
DELETE FROM students WHERE code LIKE 'XLB-NHO-%';
DELETE FROM rooms WHERE code LIKE 'XLB-NHO-%';

INSERT INTO rooms (code, name, capacity, address, campus_id, created_at, updated_at)
SELECT
    room_data.code,
    room_data.name,
    room_data.capacity,
    room_data.address,
    campus.id,
    NOW(),
    NOW()
FROM (
    VALUES
        ('XLB-NHO-P01', 'Phòng nhỏ 01', 18, 'Tầng 1 cơ sở trung tâm', 'XLB-CHUNG-CS1'),
        ('XLB-NHO-P02', 'Phòng nhỏ 02', 20, 'Tầng 1 cơ sở trung tâm', 'XLB-CHUNG-CS1'),
        ('XLB-NHO-P03', 'Phòng trung bình 01', 24, 'Tầng 2 cơ sở trung tâm', 'XLB-CHUNG-CS1'),
        ('XLB-NHO-P04', 'Phòng trung bình 02', 30, 'Tầng 1 cơ sở vệ tinh', 'XLB-CHUNG-CS2'),
        ('XLB-NHO-P05', 'Phòng lớn 01', 36, 'Tầng 2 cơ sở vệ tinh', 'XLB-CHUNG-CS2')
) AS room_data(code, name, capacity, address, campus_code)
JOIN campuses AS campus ON campus.code = room_data.campus_code;

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
    FORMAT('XLB-NHO-GV%s', LPAD(seq::text, 2, '0')),
    FORMAT('Giáo viên benchmark nhỏ %s', seq),
    FORMAT('xlb.nho.gv%s@example.com', LPAD(seq::text, 2, '0')),
    FORMAT('0909000%s', LPAD(seq::text, 3, '0')),
    FALSE,
    NULL,
    CASE WHEN seq <= 4 THEN 'FULL_TIME' ELSE 'PART_TIME' END,
    'ACTIVE',
    ARRAY['TOAN', 'ANH', 'LY', 'HOA', 'VAN', 'TIN']::TEXT[],
    'Bộ giáo viên chuẩn cho benchmark nhỏ',
    NOW(),
    NOW()
FROM generate_series(1, 8) AS seq;

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
    FORMAT('XLB-NHO-KH%s', LPAD(seq::text, 2, '0')),
    FORMAT('Khóa benchmark nhỏ %s', seq),
    'Khóa học dùng cho đo đạc đối sánh xếp lịch quy mô nhỏ',
    CASE (seq - 1) % 5
        WHEN 0 THEN 'Khối 8'
        WHEN 1 THEN 'Khối 9'
        WHEN 2 THEN 'Khối 10'
        WHEN 3 THEN 'Khối 11'
        ELSE 'Khối 12'
    END,
    CASE (seq - 1) % 4
        WHEN 0 THEN 'Toán'
        WHEN 1 THEN 'Tiếng Anh'
        WHEN 2 THEN 'Vật Lý'
        ELSE 'Hóa Học'
    END,
    6,
    120,
    12,
    1800000 + seq * 50000,
    'ACTIVE',
    ARRAY[
        CASE (seq - 1) % 4
            WHEN 0 THEN 'TOAN'
            WHEN 1 THEN 'ANH'
            WHEN 2 THEN 'LY'
            ELSE 'HOA'
        END
    ]::TEXT[],
    NOW(),
    NOW()
FROM generate_series(1, 12) AS seq;

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
    FORMAT('XLB-NHO-HS%s', LPAD(seq::text, 3, '0')),
    FORMAT('Học sinh benchmark nhỏ %s', seq),
    FORMAT('xlb.nho.hs%s@example.com', LPAD(seq::text, 3, '0')),
    FORMAT('081800%s', LPAD(seq::text, 4, '0')),
    FORMAT('091800%s', LPAD(seq::text, 4, '0')),
    CASE (seq - 1) % 5
        WHEN 0 THEN 'Khối 8'
        WHEN 1 THEN 'Khối 9'
        WHEN 2 THEN 'Khối 10'
        WHEN 3 THEN 'Khối 11'
        ELSE 'Khối 12'
    END,
    'ACTIVE',
    CASE WHEN seq % 2 = 0 THEN 'FEMALE' ELSE 'MALE' END,
    'Cụm dân cư benchmark nhỏ',
    NOW(),
    NOW()
FROM generate_series(1, 108) AS seq;

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
    FORMAT('XLB-NHO-L%s', LPAD(seq::text, 2, '0')),
    FORMAT('Lớp benchmark nhỏ %s', seq),
    'Lớp chuẩn cho benchmark xếp lịch quy mô nhỏ',
    DATE '2026-06-01' + ((seq - 1) % 5) * INTERVAL '1 day',
    DATE '2026-07-31',
    CASE (seq - 1) % 4
        WHEN 0 THEN 18
        WHEN 1 THEN 20
        WHEN 2 THEN 24
        ELSE 28
    END,
    'OPEN',
    2200000,
    NULL,
    course_ref.id,
    teacher_ref.id,
    room_ref.id,
    NOW(),
    NOW()
FROM generate_series(1, 12) AS seq
JOIN courses AS course_ref ON course_ref.code = FORMAT('XLB-NHO-KH%s', LPAD(seq::text, 2, '0'))
JOIN teachers AS teacher_ref ON teacher_ref.code = FORMAT('XLB-NHO-GV%s', LPAD((((seq - 1) % 8) + 1)::text, 2, '0'))
JOIN rooms AS room_ref ON room_ref.code = FORMAT('XLB-NHO-P%s', LPAD((((seq - 1) % 5) + 1)::text, 2, '0'));

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
            CASE (RIGHT(class_ref.code, 2)::INT - 1) % 4
                WHEN 0 THEN 'MONDAY'
                WHEN 1 THEN 'TUESDAY'
                WHEN 2 THEN 'WEDNESDAY'
                ELSE 'THURSDAY'
            END,
            CASE (RIGHT(class_ref.code, 2)::INT - 1) % 4
                WHEN 0 THEN 'XLB-CHUNG-S1'
                WHEN 1 THEN 'XLB-CHUNG-S3'
                WHEN 2 THEN 'XLB-CHUNG-S5'
                ELSE 'XLB-CHUNG-S2'
            END,
            FORMAT('XLB-NHO-P%s', LPAD((((RIGHT(class_ref.code, 2)::INT - 1) % 5) + 1)::text, 2, '0'))
        ),
        (
            CASE (RIGHT(class_ref.code, 2)::INT - 1) % 4
                WHEN 0 THEN 'WEDNESDAY'
                WHEN 1 THEN 'THURSDAY'
                WHEN 2 THEN 'FRIDAY'
                ELSE 'SATURDAY'
            END,
            CASE (RIGHT(class_ref.code, 2)::INT - 1) % 4
                WHEN 0 THEN 'XLB-CHUNG-S3'
                WHEN 1 THEN 'XLB-CHUNG-S5'
                WHEN 2 THEN 'XLB-CHUNG-S1'
                ELSE 'XLB-CHUNG-S4'
            END,
            FORMAT('XLB-NHO-P%s', LPAD((((RIGHT(class_ref.code, 2)::INT) % 5) + 1)::text, 2, '0'))
        )
) AS schedule_data(day_of_week, shift_code, room_code) ON TRUE
JOIN shifts AS shift_ref ON shift_ref.code = schedule_data.shift_code
JOIN rooms AS room_ref ON room_ref.code = schedule_data.room_code
WHERE class_ref.code LIKE 'XLB-NHO-%';

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
FROM generate_series(1, 12) AS class_seq
JOIN classes AS class_ref ON class_ref.code = FORMAT('XLB-NHO-L%s', LPAD(class_seq::text, 2, '0'))
JOIN LATERAL generate_series(1, 8 + (class_seq % 6)) AS enrollment_seq(seq) ON TRUE
JOIN students AS student_ref
    ON student_ref.code = FORMAT(
        'XLB-NHO-HS%s',
        LPAD((((class_seq - 1) * 7 + enrollment_seq.seq - 1) % 108 + 1)::text, 3, '0')
    );

COMMIT;
