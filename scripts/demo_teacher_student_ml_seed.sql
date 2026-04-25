-- Demo seed for teacher portal + student portal + Python ML pipeline
-- Project: EduCenter
-- Purpose:
--   1. Create 1 TEACHER login and 1 primary STUDENT login
--   2. Create a teacher-managed class with enrolled students
--   3. Create enough lesson / attendance / summary / academic data for
--      ml/at_risk_prediction/scripts/train_from_db.py and predict_from_db.py
--
-- Primary login accounts
--   Teacher: teacher_demo@educenter.local / password123
--   Student: student_demo01@educenter.local / password123
--
-- Password hash below is bcrypt hash for "password123".

BEGIN;

SET TIME ZONE 'Asia/Ho_Chi_Minh';

-- ============================================================================
-- 1. USERS
-- ============================================================================
INSERT INTO users (
    id, code, full_name, email, password, role, is_active, created_at, updated_at, deleted_at
) VALUES
    (
        '11111111-1111-1111-1111-000000000101',
        'USR-TEACHER-DEMO',
        'Nguyen Van Giao Vien Demo',
        'teacher_demo@educenter.local',
        '$2y$10$xyikCXU91i9HSB52vxhp2.42gxxLUSakYUJCERqMRfebIS7wP6dou',
        'TEACHER',
        TRUE,
        NOW(),
        NOW(),
        NULL
    ),
    (
        '11111111-1111-1111-1111-000000000102',
        'USR-STUDENT-DEMO-01',
        'Tran Thi Hoc Vien Demo 01',
        'student_demo01@educenter.local',
        '$2y$10$xyikCXU91i9HSB52vxhp2.42gxxLUSakYUJCERqMRfebIS7wP6dou',
        'STUDENT',
        TRUE,
        NOW(),
        NOW(),
        NULL
    )
ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code,
    full_name = EXCLUDED.full_name,
    email = EXCLUDED.email,
    password = EXCLUDED.password,
    role = EXCLUDED.role,
    is_active = EXCLUDED.is_active,
    updated_at = EXCLUDED.updated_at,
    deleted_at = EXCLUDED.deleted_at;

-- ============================================================================
-- 2. TEACHER PROFILE
-- ============================================================================
INSERT INTO teachers (
    id, code, full_name, email, phone, is_school_teacher, school_name,
    employment_type, status, notes, created_at, updated_at, deleted_at
) VALUES (
    '22222222-2222-2222-2222-000000000201',
    'TCH-DEMO-01',
    'Nguyen Van Giao Vien Demo',
    'teacher_demo@educenter.local',
    '0901000201',
    TRUE,
    'THPT Demo EduCenter',
    'PART_TIME',
    'ACTIVE',
    'Giáo viên demo dùng để kiểm thử teacher portal, attendance, summary và academic record.',
    NOW(),
    NOW(),
    NULL
)
ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code,
    full_name = EXCLUDED.full_name,
    email = EXCLUDED.email,
    phone = EXCLUDED.phone,
    is_school_teacher = EXCLUDED.is_school_teacher,
    school_name = EXCLUDED.school_name,
    employment_type = EXCLUDED.employment_type,
    status = EXCLUDED.status,
    notes = EXCLUDED.notes,
    updated_at = EXCLUDED.updated_at,
    deleted_at = EXCLUDED.deleted_at;

-- ============================================================================
-- 3. STUDENT PROFILES
-- ============================================================================
INSERT INTO students (
    id, code, full_name, email, phone, guardian_phone, grade_level,
    status, date_of_birth, gender, address, created_at, updated_at, deleted_at
) VALUES
    (
        '77777777-7777-7777-7777-000000000701',
        'STD-DEMO-01',
        'Tran Thi Hoc Vien Demo 01',
        'student_demo01@educenter.local',
        '0907000701',
        '0917000701',
        '10',
        'ACTIVE',
        '2010-05-12 00:00:00+07',
        'FEMALE',
        'Quan 1, TP.HCM',
        NOW(),
        NOW(),
        NULL
    ),
    (
        '77777777-7777-7777-7777-000000000702',
        'STD-DEMO-02',
        'Le Minh Hoc Vien 02',
        'student_demo02@educenter.local',
        '0907000702',
        '0917000702',
        '10',
        'ACTIVE',
        '2010-07-21 00:00:00+07',
        'MALE',
        'Quan 3, TP.HCM',
        NOW(),
        NOW(),
        NULL
    ),
    (
        '77777777-7777-7777-7777-000000000703',
        'STD-DEMO-03',
        'Pham Gia Hoc Vien 03',
        'student_demo03@educenter.local',
        '0907000703',
        '0917000703',
        '10',
        'ACTIVE',
        '2010-03-08 00:00:00+07',
        'FEMALE',
        'Quan 7, TP.HCM',
        NOW(),
        NOW(),
        NULL
    ),
    (
        '77777777-7777-7777-7777-000000000704',
        'STD-DEMO-04',
        'Hoang Nam Hoc Vien 04',
        'student_demo04@educenter.local',
        '0907000704',
        '0917000704',
        '10',
        'ACTIVE',
        '2010-01-26 00:00:00+07',
        'MALE',
        'Thu Duc, TP.HCM',
        NOW(),
        NOW(),
        NULL
    ),
    (
        '77777777-7777-7777-7777-000000000705',
        'STD-DEMO-05',
        'Vo Anh Hoc Vien 05',
        'student_demo05@educenter.local',
        '0907000705',
        '0917000705',
        '10',
        'ACTIVE',
        '2010-10-15 00:00:00+07',
        'FEMALE',
        'Binh Thanh, TP.HCM',
        NOW(),
        NOW(),
        NULL
    ),
    (
        '77777777-7777-7777-7777-000000000706',
        'STD-DEMO-06',
        'Dang Khoa Hoc Vien 06',
        'student_demo06@educenter.local',
        '0907000706',
        '0917000706',
        '10',
        'ACTIVE',
        '2010-09-02 00:00:00+07',
        'MALE',
        'Go Vap, TP.HCM',
        NOW(),
        NOW(),
        NULL
    )
ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code,
    full_name = EXCLUDED.full_name,
    email = EXCLUDED.email,
    phone = EXCLUDED.phone,
    guardian_phone = EXCLUDED.guardian_phone,
    grade_level = EXCLUDED.grade_level,
    status = EXCLUDED.status,
    date_of_birth = EXCLUDED.date_of_birth,
    gender = EXCLUDED.gender,
    address = EXCLUDED.address,
    updated_at = EXCLUDED.updated_at,
    deleted_at = EXCLUDED.deleted_at;

-- ============================================================================
-- 4. COURSE / PROGRAM / ROOM / SHIFT / CLASS
-- ============================================================================
INSERT INTO courses (
    id, code, name, description, grade_level, subject, session_count,
    session_duration_minutes, total_hours, price, status, created_at, updated_at, deleted_at
) VALUES (
    '55555555-5555-5555-5555-000000000502',
    'COURSE-MATH10-DEMO',
    'Toán 10 Nâng cao - Demo ML',
    'Khóa học demo phục vụ kiểm thử scheduling, lesson và dự báo AT_RISK.',
    '10',
    'Toán',
    8,
    120,
    16.00,
    2400000,
    'ACTIVE',
    NOW(),
    NOW(),
    NULL
)
ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    grade_level = EXCLUDED.grade_level,
    subject = EXCLUDED.subject,
    session_count = EXCLUDED.session_count,
    session_duration_minutes = EXCLUDED.session_duration_minutes,
    total_hours = EXCLUDED.total_hours,
    price = EXCLUDED.price,
    status = EXCLUDED.status,
    updated_at = EXCLUDED.updated_at,
    deleted_at = EXCLUDED.deleted_at;

INSERT INTO programs (
    id, code, name, track, effective_from, effective_to, created_by_id,
    approved_by_id, approval_note, published_at, archived_at, status, created_at, updated_at, deleted_at
) VALUES (
    '55555555-5555-5555-5555-000000000501',
    'PROGRAM-MATH10-DEMO',
    'Chương trình Toán 10 Demo',
    'BASIC',
    '2026-03-01 00:00:00+07',
    '2026-06-30 00:00:00+07',
    NULL,
    NULL,
    'Program demo để phục vụ lớp học và dữ liệu predictive.',
    '2026-03-01 08:00:00+07',
    NULL,
    'ACTIVE',
    NOW(),
    NOW(),
    NULL
)
ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code,
    name = EXCLUDED.name,
    track = EXCLUDED.track,
    effective_from = EXCLUDED.effective_from,
    effective_to = EXCLUDED.effective_to,
    approval_note = EXCLUDED.approval_note,
    published_at = EXCLUDED.published_at,
    status = EXCLUDED.status,
    updated_at = EXCLUDED.updated_at,
    deleted_at = EXCLUDED.deleted_at;

INSERT INTO program_courses (id, program_id, course_id) VALUES
    (
        '55555555-5555-5555-5555-000000000503',
        '55555555-5555-5555-5555-000000000501',
        '55555555-5555-5555-5555-000000000502'
    )
ON CONFLICT (id) DO UPDATE SET
    program_id = EXCLUDED.program_id,
    course_id = EXCLUDED.course_id;

INSERT INTO rooms (id, code, name, capacity, address, created_at, updated_at, deleted_at) VALUES
    (
        '33333333-3333-3333-3333-000000000301',
        'ROOM-DEMO-01',
        'Phòng học Demo 01',
        30,
        'Lầu 2, chi nhánh EduCenter TP.HCM',
        NOW(),
        NOW(),
        NULL
    )
ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code,
    name = EXCLUDED.name,
    capacity = EXCLUDED.capacity,
    address = EXCLUDED.address,
    updated_at = EXCLUDED.updated_at,
    deleted_at = EXCLUDED.deleted_at;

INSERT INTO shifts (
    id, code, name, start_time, end_time, duration_minutes,
    session_type, is_active, notes, created_at, updated_at
) VALUES
    (
        '44444444-4444-4444-4444-000000000401',
        'SHIFT-DEMO-18H',
        'Ca tối 18h00 - 20h00',
        '18:00',
        '20:00',
        120,
        'EVENING',
        TRUE,
        'Ca demo phục vụ lớp học và benchmark lesson/predictive.',
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code,
    name = EXCLUDED.name,
    start_time = EXCLUDED.start_time,
    end_time = EXCLUDED.end_time,
    duration_minutes = EXCLUDED.duration_minutes,
    session_type = EXCLUDED.session_type,
    is_active = EXCLUDED.is_active,
    notes = EXCLUDED.notes,
    updated_at = EXCLUDED.updated_at;

INSERT INTO classes (
    id, code, name, notes, start_date, end_date, max_students, status, price,
    program_id, course_id, teacher_id, room_id, created_at, updated_at, deleted_at
) VALUES (
    '66666666-6666-6666-6666-000000000601',
    'CLASS-MATH10-DEMO-A',
    'Lớp Toán 10 Demo A',
    'Lớp demo dùng để kiểm thử teacher portal, student portal và pipeline ML từ DB.',
    '2026-03-30 00:00:00+07',
    '2026-06-30 23:59:59+07',
    20,
    'OPEN',
    2400000,
    '55555555-5555-5555-5555-000000000501',
    '55555555-5555-5555-5555-000000000502',
    '22222222-2222-2222-2222-000000000201',
    '33333333-3333-3333-3333-000000000301',
    NOW(),
    NOW(),
    NULL
)
ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code,
    name = EXCLUDED.name,
    notes = EXCLUDED.notes,
    start_date = EXCLUDED.start_date,
    end_date = EXCLUDED.end_date,
    max_students = EXCLUDED.max_students,
    status = EXCLUDED.status,
    price = EXCLUDED.price,
    program_id = EXCLUDED.program_id,
    course_id = EXCLUDED.course_id,
    teacher_id = EXCLUDED.teacher_id,
    room_id = EXCLUDED.room_id,
    updated_at = EXCLUDED.updated_at,
    deleted_at = EXCLUDED.deleted_at;

INSERT INTO class_schedules (id, class_id, day_of_week, start_time, end_time, shift_id, room_id) VALUES
    (
        '66666666-6666-6666-6666-000000000602',
        '66666666-6666-6666-6666-000000000601',
        'MONDAY',
        '18:00',
        '20:00',
        '44444444-4444-4444-4444-000000000401',
        '33333333-3333-3333-3333-000000000301'
    ),
    (
        '66666666-6666-6666-6666-000000000603',
        '66666666-6666-6666-6666-000000000601',
        'WEDNESDAY',
        '18:00',
        '20:00',
        '44444444-4444-4444-4444-000000000401',
        '33333333-3333-3333-3333-000000000301'
    )
ON CONFLICT (id) DO UPDATE SET
    class_id = EXCLUDED.class_id,
    day_of_week = EXCLUDED.day_of_week,
    start_time = EXCLUDED.start_time,
    end_time = EXCLUDED.end_time,
    shift_id = EXCLUDED.shift_id,
    room_id = EXCLUDED.room_id;

-- ============================================================================
-- 5. ENROLLMENTS
-- ============================================================================
INSERT INTO enrollments (
    id, class_id, student_id, status, approved_at, rejected_at, created_at, updated_at
) VALUES
    ('88888888-8888-8888-8888-000000000801', '66666666-6666-6666-6666-000000000601', '77777777-7777-7777-7777-000000000701', 'ENROLLED', '2026-03-28 09:00:00+07', NULL, '2026-03-28 09:00:00+07', NOW()),
    ('88888888-8888-8888-8888-000000000802', '66666666-6666-6666-6666-000000000601', '77777777-7777-7777-7777-000000000702', 'ENROLLED', '2026-03-28 09:05:00+07', NULL, '2026-03-28 09:05:00+07', NOW()),
    ('88888888-8888-8888-8888-000000000803', '66666666-6666-6666-6666-000000000601', '77777777-7777-7777-7777-000000000703', 'ENROLLED', '2026-03-28 09:10:00+07', NULL, '2026-03-28 09:10:00+07', NOW()),
    ('88888888-8888-8888-8888-000000000804', '66666666-6666-6666-6666-000000000601', '77777777-7777-7777-7777-000000000704', 'ENROLLED', '2026-03-28 09:15:00+07', NULL, '2026-03-28 09:15:00+07', NOW()),
    ('88888888-8888-8888-8888-000000000805', '66666666-6666-6666-6666-000000000601', '77777777-7777-7777-7777-000000000705', 'ENROLLED', '2026-03-28 09:20:00+07', NULL, '2026-03-28 09:20:00+07', NOW()),
    ('88888888-8888-8888-8888-000000000806', '66666666-6666-6666-6666-000000000601', '77777777-7777-7777-7777-000000000706', 'ENROLLED', '2026-03-28 09:25:00+07', NULL, '2026-03-28 09:25:00+07', NOW())
ON CONFLICT (id) DO UPDATE SET
    class_id = EXCLUDED.class_id,
    student_id = EXCLUDED.student_id,
    status = EXCLUDED.status,
    approved_at = EXCLUDED.approved_at,
    rejected_at = EXCLUDED.rejected_at,
    created_at = EXCLUDED.created_at,
    updated_at = EXCLUDED.updated_at;

-- ============================================================================
-- 6. LESSONS
-- ============================================================================
INSERT INTO lessons (
    id, class_id, teacher_id, date_start, date_end, room_id, notes, created_at, updated_at
) VALUES
    ('99999999-9999-9999-9999-000000000901', '66666666-6666-6666-6666-000000000601', '22222222-2222-2222-2222-000000000201', '2026-03-30 18:00:00+07', '2026-03-30 20:00:00+07', '33333333-3333-3333-3333-000000000301', 'Buổi 1 - Khởi động đại số', NOW(), NOW()),
    ('99999999-9999-9999-9999-000000000902', '66666666-6666-6666-6666-000000000601', '22222222-2222-2222-2222-000000000201', '2026-04-01 18:00:00+07', '2026-04-01 20:00:00+07', '33333333-3333-3333-3333-000000000301', 'Buổi 2 - Phương trình bậc nhất', NOW(), NOW()),
    ('99999999-9999-9999-9999-000000000903', '66666666-6666-6666-6666-000000000601', '22222222-2222-2222-2222-000000000201', '2026-04-06 18:00:00+07', '2026-04-06 20:00:00+07', '33333333-3333-3333-3333-000000000301', 'Buổi 3 - Hệ phương trình', NOW(), NOW()),
    ('99999999-9999-9999-9999-000000000904', '66666666-6666-6666-6666-000000000601', '22222222-2222-2222-2222-000000000201', '2026-04-08 18:00:00+07', '2026-04-08 20:00:00+07', '33333333-3333-3333-3333-000000000301', 'Buổi 4 - Đồ thị hàm số', NOW(), NOW()),
    ('99999999-9999-9999-9999-000000000905', '66666666-6666-6666-6666-000000000601', '22222222-2222-2222-2222-000000000201', '2026-04-13 18:00:00+07', '2026-04-13 20:00:00+07', '33333333-3333-3333-3333-000000000301', 'Buổi 5 - Bài tập ứng dụng', NOW(), NOW()),
    ('99999999-9999-9999-9999-000000000906', '66666666-6666-6666-6666-000000000601', '22222222-2222-2222-2222-000000000201', '2026-04-15 18:00:00+07', '2026-04-15 20:00:00+07', '33333333-3333-3333-3333-000000000301', 'Buổi 6 - Luyện đề cơ bản', NOW(), NOW()),
    ('99999999-9999-9999-9999-000000000907', '66666666-6666-6666-6666-000000000601', '22222222-2222-2222-2222-000000000201', '2026-04-20 18:00:00+07', '2026-04-20 20:00:00+07', '33333333-3333-3333-3333-000000000301', 'Buổi 7 - Luyện đề nâng cao', NOW(), NOW()),
    ('99999999-9999-9999-9999-000000000908', '66666666-6666-6666-6666-000000000601', '22222222-2222-2222-2222-000000000201', '2026-04-22 18:00:00+07', '2026-04-22 20:00:00+07', '33333333-3333-3333-3333-000000000301', 'Buổi 8 - Tổng ôn và kiểm tra', NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
    class_id = EXCLUDED.class_id,
    teacher_id = EXCLUDED.teacher_id,
    date_start = EXCLUDED.date_start,
    date_end = EXCLUDED.date_end,
    room_id = EXCLUDED.room_id,
    notes = EXCLUDED.notes,
    updated_at = EXCLUDED.updated_at;

-- ============================================================================
-- 7. ATTENDANCE
-- Internal status mapping:
--   1 = PRESENT, 2 = ABSENT, 3 = EXCUSED, 4 = LATE, 5 = EARLY
-- ============================================================================
INSERT INTO attendances (
    id, lesson_id, student_id, status, note, marked_at, created_at, updated_at
) VALUES
    -- Lesson 1
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000001001', '99999999-9999-9999-9999-000000000901', '77777777-7777-7777-7777-000000000701', 1, 'Có mặt đầy đủ', '2026-03-30 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000001002', '99999999-9999-9999-9999-000000000901', '77777777-7777-7777-7777-000000000702', 1, 'Có mặt', '2026-03-30 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000001003', '99999999-9999-9999-9999-000000000901', '77777777-7777-7777-7777-000000000703', 1, 'Có mặt', '2026-03-30 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000001004', '99999999-9999-9999-9999-000000000901', '77777777-7777-7777-7777-000000000704', 1, 'Có mặt', '2026-03-30 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000001005', '99999999-9999-9999-9999-000000000901', '77777777-7777-7777-7777-000000000705', 1, 'Có mặt', '2026-03-30 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000001006', '99999999-9999-9999-9999-000000000901', '77777777-7777-7777-7777-000000000706', 1, 'Có mặt', '2026-03-30 20:05:00+07', NOW(), NOW()),

    -- Lesson 2
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000002001', '99999999-9999-9999-9999-000000000902', '77777777-7777-7777-7777-000000000701', 1, 'Có mặt đầy đủ', '2026-04-01 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000002002', '99999999-9999-9999-9999-000000000902', '77777777-7777-7777-7777-000000000702', 2, 'Vắng không phép', '2026-04-01 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000002003', '99999999-9999-9999-9999-000000000902', '77777777-7777-7777-7777-000000000703', 1, 'Có mặt', '2026-04-01 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000002004', '99999999-9999-9999-9999-000000000902', '77777777-7777-7777-7777-000000000704', 1, 'Có mặt', '2026-04-01 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000002005', '99999999-9999-9999-9999-000000000902', '77777777-7777-7777-7777-000000000705', 1, 'Có mặt', '2026-04-01 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000002006', '99999999-9999-9999-9999-000000000902', '77777777-7777-7777-7777-000000000706', 1, 'Có mặt', '2026-04-01 20:05:00+07', NOW(), NOW()),

    -- Lesson 3
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000003001', '99999999-9999-9999-9999-000000000903', '77777777-7777-7777-7777-000000000701', 1, 'Có mặt', '2026-04-06 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000003002', '99999999-9999-9999-9999-000000000903', '77777777-7777-7777-7777-000000000702', 2, 'Vắng không phép', '2026-04-06 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000003003', '99999999-9999-9999-9999-000000000903', '77777777-7777-7777-7777-000000000703', 4, 'Đi muộn 10 phút', '2026-04-06 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000003004', '99999999-9999-9999-9999-000000000903', '77777777-7777-7777-7777-000000000704', 1, 'Có mặt', '2026-04-06 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000003005', '99999999-9999-9999-9999-000000000903', '77777777-7777-7777-7777-000000000705', 1, 'Có mặt', '2026-04-06 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000003006', '99999999-9999-9999-9999-000000000903', '77777777-7777-7777-7777-000000000706', 1, 'Có mặt', '2026-04-06 20:05:00+07', NOW(), NOW()),

    -- Lesson 4
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000004001', '99999999-9999-9999-9999-000000000904', '77777777-7777-7777-7777-000000000701', 4, 'Đi muộn 5 phút', '2026-04-08 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000004002', '99999999-9999-9999-9999-000000000904', '77777777-7777-7777-7777-000000000702', 1, 'Có mặt', '2026-04-08 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000004003', '99999999-9999-9999-9999-000000000904', '77777777-7777-7777-7777-000000000703', 1, 'Có mặt', '2026-04-08 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000004004', '99999999-9999-9999-9999-000000000904', '77777777-7777-7777-7777-000000000704', 1, 'Có mặt', '2026-04-08 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000004005', '99999999-9999-9999-9999-000000000904', '77777777-7777-7777-7777-000000000705', 1, 'Có mặt', '2026-04-08 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000004006', '99999999-9999-9999-9999-000000000904', '77777777-7777-7777-7777-000000000706', 2, 'Vắng không phép', '2026-04-08 20:05:00+07', NOW(), NOW()),

    -- Lesson 5
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000005001', '99999999-9999-9999-9999-000000000905', '77777777-7777-7777-7777-000000000701', 1, 'Có mặt', '2026-04-13 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000005002', '99999999-9999-9999-9999-000000000905', '77777777-7777-7777-7777-000000000702', 2, 'Vắng không phép', '2026-04-13 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000005003', '99999999-9999-9999-9999-000000000905', '77777777-7777-7777-7777-000000000703', 1, 'Có mặt', '2026-04-13 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000005004', '99999999-9999-9999-9999-000000000905', '77777777-7777-7777-7777-000000000704', 1, 'Có mặt', '2026-04-13 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000005005', '99999999-9999-9999-9999-000000000905', '77777777-7777-7777-7777-000000000705', 4, 'Đi muộn 7 phút', '2026-04-13 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000005006', '99999999-9999-9999-9999-000000000905', '77777777-7777-7777-7777-000000000706', 3, 'Có phép', '2026-04-13 20:05:00+07', NOW(), NOW()),

    -- Lesson 6
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000006001', '99999999-9999-9999-9999-000000000906', '77777777-7777-7777-7777-000000000701', 1, 'Có mặt', '2026-04-15 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000006002', '99999999-9999-9999-9999-000000000906', '77777777-7777-7777-7777-000000000702', 2, 'Vắng không phép', '2026-04-15 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000006003', '99999999-9999-9999-9999-000000000906', '77777777-7777-7777-7777-000000000703', 3, 'Nghỉ có phép', '2026-04-15 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000006004', '99999999-9999-9999-9999-000000000906', '77777777-7777-7777-7777-000000000704', 1, 'Có mặt', '2026-04-15 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000006005', '99999999-9999-9999-9999-000000000906', '77777777-7777-7777-7777-000000000705', 1, 'Có mặt', '2026-04-15 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000006006', '99999999-9999-9999-9999-000000000906', '77777777-7777-7777-7777-000000000706', 2, 'Vắng không phép', '2026-04-15 20:05:00+07', NOW(), NOW()),

    -- Lesson 7
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000007001', '99999999-9999-9999-9999-000000000907', '77777777-7777-7777-7777-000000000701', 1, 'Có mặt', '2026-04-20 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000007002', '99999999-9999-9999-9999-000000000907', '77777777-7777-7777-7777-000000000702', 2, 'Vắng không phép', '2026-04-20 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000007003', '99999999-9999-9999-9999-000000000907', '77777777-7777-7777-7777-000000000703', 1, 'Có mặt', '2026-04-20 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000007004', '99999999-9999-9999-9999-000000000907', '77777777-7777-7777-7777-000000000704', 1, 'Có mặt', '2026-04-20 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000007005', '99999999-9999-9999-9999-000000000907', '77777777-7777-7777-7777-000000000705', 1, 'Có mặt', '2026-04-20 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000007006', '99999999-9999-9999-9999-000000000907', '77777777-7777-7777-7777-000000000706', 2, 'Vắng không phép', '2026-04-20 20:05:00+07', NOW(), NOW()),

    -- Lesson 8
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000008001', '99999999-9999-9999-9999-000000000908', '77777777-7777-7777-7777-000000000701', 1, 'Có mặt', '2026-04-22 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000008002', '99999999-9999-9999-9999-000000000908', '77777777-7777-7777-7777-000000000702', 2, 'Vắng không phép', '2026-04-22 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000008003', '99999999-9999-9999-9999-000000000908', '77777777-7777-7777-7777-000000000703', 1, 'Có mặt', '2026-04-22 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000008004', '99999999-9999-9999-9999-000000000908', '77777777-7777-7777-7777-000000000704', 1, 'Có mặt', '2026-04-22 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000008005', '99999999-9999-9999-9999-000000000908', '77777777-7777-7777-7777-000000000705', 1, 'Có mặt', '2026-04-22 20:05:00+07', NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-000000008006', '99999999-9999-9999-9999-000000000908', '77777777-7777-7777-7777-000000000706', 1, 'Có mặt', '2026-04-22 20:05:00+07', NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
    lesson_id = EXCLUDED.lesson_id,
    student_id = EXCLUDED.student_id,
    status = EXCLUDED.status,
    note = EXCLUDED.note,
    marked_at = EXCLUDED.marked_at,
    updated_at = EXCLUDED.updated_at;

-- ============================================================================
-- 8. LESSON SUMMARIES
-- ============================================================================
INSERT INTO lesson_summaries (
    id, lesson_id, topic, lesson_content, class_feedback, homework, homework_deadline,
    teacher_notes, created_by_id, created_at, updated_at
) VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-000000000a01', '99999999-9999-9999-9999-000000000901', 'Khởi động đại số', 'Ôn tập biểu thức, hằng đẳng thức và đặt ẩn phụ đơn giản.', 'Lớp tập trung, nhiều học viên nắm chắc phần biến đổi cơ bản.', 'Làm 10 bài biểu thức trang 12.', '2026-03-31 22:00:00+07', 'Nhóm học tốt, cần quan sát thêm học viên có nền tảng yếu.', '11111111-1111-1111-1111-000000000101', NOW(), NOW()),
    ('aaaaaaaa-aaaa-aaaa-aaaa-000000000a02', '99999999-9999-9999-9999-000000000902', 'Phương trình bậc nhất', 'Giải phương trình và bài toán ứng dụng cơ bản.', 'Một số học viên bắt đầu có dấu hiệu chậm nhịp.', 'Bài tập 1 đến 8 phần phương trình bậc nhất.', '2026-04-02 22:00:00+07', 'Cần nhắc các em làm bài đầy đủ trước buổi sau.', '11111111-1111-1111-1111-000000000101', NOW(), NOW()),
    ('aaaaaaaa-aaaa-aaaa-aaaa-000000000a03', '99999999-9999-9999-9999-000000000903', 'Hệ phương trình', 'Phương pháp cộng đại số và thế.', 'Nhóm khá làm tốt, nhóm yếu mất nhiều thời gian.', 'Hoàn thành phiếu bài tập hệ phương trình.', '2026-04-07 22:00:00+07', 'Nên chia nhóm hỗ trợ ở đầu giờ kế tiếp.', '11111111-1111-1111-1111-000000000101', NOW(), NOW()),
    ('aaaaaaaa-aaaa-aaaa-aaaa-000000000a04', '99999999-9999-9999-9999-000000000904', 'Đồ thị hàm số', 'Vẽ đồ thị hàm số bậc nhất và đọc dữ liệu từ đồ thị.', 'Lớp tương tác tốt nhưng vẫn còn chênh lệch năng lực.', 'Làm worksheet đồ thị hàm số.', '2026-04-09 22:00:00+07', 'Cần theo dõi học viên thường xuyên bỏ bài tập.', '11111111-1111-1111-1111-000000000101', NOW(), NOW()),
    ('aaaaaaaa-aaaa-aaaa-aaaa-000000000a05', '99999999-9999-9999-9999-000000000905', 'Bài tập ứng dụng', 'Giải các bài toán tổng hợp từ 4 buổi đầu.', 'Bắt đầu phân hóa rõ nhóm theo kịp và nhóm có nguy cơ tụt lại.', 'Hoàn thành đề mini test số 1.', '2026-04-14 22:00:00+07', 'Các học viên có chuyên cần thấp đang ảnh hưởng trực tiếp tới kết quả.', '11111111-1111-1111-1111-000000000101', NOW(), NOW()),
    ('aaaaaaaa-aaaa-aaaa-aaaa-000000000a06', '99999999-9999-9999-9999-000000000906', 'Luyện đề cơ bản', 'Chữa đề và củng cố phương pháp làm bài.', 'Lớp giữ nhịp ổn, nhưng một số em còn vắng không phép.', 'Chữa lại lỗi sai trong đề mini test số 1.', '2026-04-16 22:00:00+07', 'Cần trao đổi phụ huynh với nhóm học viên yếu.', '11111111-1111-1111-1111-000000000101', NOW(), NOW()),
    ('aaaaaaaa-aaaa-aaaa-aaaa-000000000a07', '99999999-9999-9999-9999-000000000907', 'Luyện đề nâng cao', 'Áp dụng chiến lược giải nhanh và xử lý câu vận dụng.', 'Nhóm khá tiến bộ rõ, nhóm yếu vẫn chưa ổn định.', 'Làm bộ đề nâng cao số 2.', '2026-04-21 22:00:00+07', 'Cần ưu tiên can thiệp sớm cho học viên thuộc nhóm rủi ro.', '11111111-1111-1111-1111-000000000101', NOW(), NOW()),
    ('aaaaaaaa-aaaa-aaaa-aaaa-000000000a08', '99999999-9999-9999-9999-000000000908', 'Tổng ôn và kiểm tra', 'Tổng kết nội dung chính và làm bài kiểm tra ngắn cuối chuỗi.', 'Đã quan sát rõ sự khác biệt giữa nhóm ổn định và nhóm AT_RISK.', 'Xem lại toàn bộ lỗi sai và ôn tập chương tiếp theo.', '2026-04-23 22:00:00+07', 'Dữ liệu buổi này dùng tốt để thử nghiệm dashboard cảnh báo.', '11111111-1111-1111-1111-000000000101', NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
    lesson_id = EXCLUDED.lesson_id,
    topic = EXCLUDED.topic,
    lesson_content = EXCLUDED.lesson_content,
    class_feedback = EXCLUDED.class_feedback,
    homework = EXCLUDED.homework,
    homework_deadline = EXCLUDED.homework_deadline,
    teacher_notes = EXCLUDED.teacher_notes,
    created_by_id = EXCLUDED.created_by_id,
    updated_at = EXCLUDED.updated_at;

-- ============================================================================
-- 9. ACADEMIC RECORDS
-- ============================================================================
INSERT INTO academic_records (
    id, lesson_summary_id, student_id, homework_completed, homework_score,
    attitude_rating, participation_score, personal_comment, total_score,
    is_completed, created_at, updated_at
) VALUES
    -- Lesson 1
    ('cccccccc-cccc-cccc-cccc-000000001001', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a01', '77777777-7777-7777-7777-000000000701', TRUE, 8.5, 5, 8.5, 'Học tốt và chủ động phát biểu.', 8.5, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000001002', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a01', '77777777-7777-7777-7777-000000000702', TRUE, 5.0, 3, 5.4, 'Nền tảng còn yếu, cần kèm thêm.', 5.2, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000001003', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a01', '77777777-7777-7777-7777-000000000703', TRUE, 7.4, 4, 7.8, 'Theo bài khá tốt.', 7.6, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000001004', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a01', '77777777-7777-7777-7777-000000000704', TRUE, 6.0, 4, 6.0, 'Cần thêm thời gian làm bài.', 6.0, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000001005', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a01', '77777777-7777-7777-7777-000000000705', TRUE, 6.8, 4, 6.8, 'Ổn định.', 6.8, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000001006', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a01', '77777777-7777-7777-7777-000000000706', TRUE, 6.4, 4, 6.6, 'Mức độ theo bài tạm ổn.', 6.5, TRUE, NOW(), NOW()),

    -- Lesson 2
    ('cccccccc-cccc-cccc-cccc-000000002001', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a02', '77777777-7777-7777-7777-000000000701', TRUE, 8.0, 5, 8.4, 'Giữ phong độ tốt.', 8.2, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000002002', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a02', '77777777-7777-7777-7777-000000000702', FALSE, 4.5, 2, 5.0, 'Bỏ bài tập và vắng học.', 4.8, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000002003', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a02', '77777777-7777-7777-7777-000000000703', TRUE, 7.8, 4, 7.8, 'Làm bài tốt.', 7.8, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000002004', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a02', '77777777-7777-7777-7777-000000000704', TRUE, 5.8, 4, 5.8, 'Bắt đầu có dấu hiệu chậm.', 5.8, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000002005', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a02', '77777777-7777-7777-7777-000000000705', TRUE, 6.8, 4, 7.0, 'Ổn định.', 6.9, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000002006', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a02', '77777777-7777-7777-7777-000000000706', TRUE, 6.2, 4, 6.4, 'Theo bài nhưng chưa bền.', 6.3, TRUE, NOW(), NOW()),

    -- Lesson 3
    ('cccccccc-cccc-cccc-cccc-000000003001', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a03', '77777777-7777-7777-7777-000000000701', TRUE, 8.2, 5, 8.6, 'Làm tốt dạng hệ phương trình.', 8.4, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000003002', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a03', '77777777-7777-7777-7777-000000000702', FALSE, 4.0, 2, 4.8, 'Mất gốc và thiếu bài tập.', 4.4, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000003003', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a03', '77777777-7777-7777-7777-000000000703', TRUE, 7.2, 4, 7.6, 'Đi muộn nhưng vẫn hoàn thành bài.', 7.4, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000003004', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a03', '77777777-7777-7777-7777-000000000704', TRUE, 5.2, 3, 5.2, 'Cần tăng tốc độ làm bài.', 5.2, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000003005', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a03', '77777777-7777-7777-7777-000000000705', TRUE, 7.0, 4, 7.2, 'Duy trì tốt.', 7.1, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000003006', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a03', '77777777-7777-7777-7777-000000000706', TRUE, 6.0, 4, 6.2, 'Nắm bài nhưng chưa chủ động.', 6.1, TRUE, NOW(), NOW()),

    -- Lesson 4
    ('cccccccc-cccc-cccc-cccc-000000004001', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a04', '77777777-7777-7777-7777-000000000701', TRUE, 8.4, 5, 8.8, 'Đi muộn nhẹ nhưng vẫn theo kịp.', 8.6, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000004002', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a04', '77777777-7777-7777-7777-000000000702', FALSE, 4.4, 2, 4.6, 'Cần can thiệp sớm.', 4.5, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000004003', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a04', '77777777-7777-7777-7777-000000000703', TRUE, 7.6, 4, 7.8, 'Đồ thị làm đúng phần lớn.', 7.7, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000004004', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a04', '77777777-7777-7777-7777-000000000704', FALSE, 4.8, 3, 5.0, 'Hiểu ý nhưng làm chưa chắc.', 4.9, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000004005', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a04', '77777777-7777-7777-7777-000000000705', TRUE, 6.6, 4, 6.8, 'Tiến độ ổn.', 6.7, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000004006', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a04', '77777777-7777-7777-7777-000000000706', TRUE, 5.8, 3, 6.0, 'Vắng 1 buổi nên hụt nhịp.', 5.9, TRUE, NOW(), NOW()),

    -- Lesson 5
    ('cccccccc-cccc-cccc-cccc-000000005001', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a05', '77777777-7777-7777-7777-000000000701', TRUE, 8.8, 5, 8.8, 'Giữ nhóm dẫn đầu.', 8.8, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000005002', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a05', '77777777-7777-7777-7777-000000000702', FALSE, 4.2, 2, 4.2, 'Bắt đầu rơi khỏi nhịp lớp.', 4.2, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000005003', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a05', '77777777-7777-7777-7777-000000000703', TRUE, 7.8, 4, 8.0, 'Ổn định.', 7.9, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000005004', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a05', '77777777-7777-7777-7777-000000000704', FALSE, 4.6, 3, 4.8, 'Kết quả giảm liên tiếp.', 4.7, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000005005', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a05', '77777777-7777-7777-7777-000000000705', TRUE, 7.0, 4, 7.0, 'Đi muộn nhưng hoàn thành tốt.', 7.0, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000005006', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a05', '77777777-7777-7777-7777-000000000706', TRUE, 5.6, 3, 5.8, 'Có phép nhưng cần củng cố.', 5.7, TRUE, NOW(), NOW()),

    -- Lesson 6
    ('cccccccc-cccc-cccc-cccc-000000006001', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a06', '77777777-7777-7777-7777-000000000701', TRUE, 8.6, 5, 8.8, 'Duy trì kết quả rất tốt.', 8.7, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000006002', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a06', '77777777-7777-7777-7777-000000000702', FALSE, 4.0, 2, 4.0, 'Chưa cải thiện được chuyên cần.', 4.0, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000006003', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a06', '77777777-7777-7777-7777-000000000703', TRUE, 7.4, 4, 7.6, 'Nghỉ có phép nhưng vẫn theo được bài.', 7.5, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000006004', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a06', '77777777-7777-7777-7777-000000000704', FALSE, 4.2, 3, 4.6, 'Cần phụ đạo riêng.', 4.4, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000006005', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a06', '77777777-7777-7777-7777-000000000705', TRUE, 7.2, 4, 7.2, 'Giữ nhịp ổn.', 7.2, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000006006', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a06', '77777777-7777-7777-7777-000000000706', TRUE, 5.4, 3, 5.6, 'Vắng không phép làm giảm kết quả.', 5.5, TRUE, NOW(), NOW()),

    -- Lesson 7
    ('cccccccc-cccc-cccc-cccc-000000007001', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a07', '77777777-7777-7777-7777-000000000701', TRUE, 8.8, 5, 9.0, 'Hoàn thành tốt đề nâng cao.', 8.9, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000007002', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a07', '77777777-7777-7777-7777-000000000702', FALSE, 3.8, 2, 3.8, 'Nguy cơ học kém rất rõ.', 3.8, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000007003', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a07', '77777777-7777-7777-7777-000000000703', FALSE, 7.6, 4, 8.0, 'Bỏ 1 bài tập nhưng vẫn đạt.', 7.8, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000007004', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a07', '77777777-7777-7777-7777-000000000704', FALSE, 4.2, 2, 4.4, 'Điểm tiếp tục giảm.', 4.3, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000007005', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a07', '77777777-7777-7777-7777-000000000705', TRUE, 6.8, 4, 7.0, 'Ổn định.', 6.9, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000007006', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a07', '77777777-7777-7777-7777-000000000706', TRUE, 5.2, 3, 5.6, 'Vắng kéo dài, cần cảnh báo.', 5.4, TRUE, NOW(), NOW()),

    -- Lesson 8
    ('cccccccc-cccc-cccc-cccc-000000008001', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a08', '77777777-7777-7777-7777-000000000701', TRUE, 8.4, 5, 8.6, 'Kết thúc chuỗi học rất tốt.', 8.5, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000008002', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a08', '77777777-7777-7777-7777-000000000702', FALSE, 4.0, 2, 4.2, 'Cần phụ đạo và theo dõi sát.', 4.1, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000008003', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a08', '77777777-7777-7777-7777-000000000703', TRUE, 7.4, 4, 7.8, 'Hoàn thành ổn.', 7.6, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000008004', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a08', '77777777-7777-7777-7777-000000000704', FALSE, 4.0, 2, 4.4, 'Thuộc nhóm cần can thiệp sớm.', 4.2, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000008005', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a08', '77777777-7777-7777-7777-000000000705', TRUE, 7.0, 4, 7.2, 'Kết quả ổn định.', 7.1, TRUE, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-000000008006', 'aaaaaaaa-aaaa-aaaa-aaaa-000000000a08', '77777777-7777-7777-7777-000000000706', FALSE, 5.0, 3, 5.6, 'Đã có tiến triển nhẹ ở buổi cuối.', 5.3, TRUE, NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
    lesson_summary_id = EXCLUDED.lesson_summary_id,
    student_id = EXCLUDED.student_id,
    homework_completed = EXCLUDED.homework_completed,
    homework_score = EXCLUDED.homework_score,
    attitude_rating = EXCLUDED.attitude_rating,
    participation_score = EXCLUDED.participation_score,
    personal_comment = EXCLUDED.personal_comment,
    total_score = EXCLUDED.total_score,
    is_completed = EXCLUDED.is_completed,
    updated_at = EXCLUDED.updated_at;

-- ============================================================================
-- 10. LEAVE REQUESTS
-- ============================================================================
INSERT INTO leave_requests (
    id, student_id, leave_type, apply_date, late_minutes, early_minutes, reason,
    documents, class_id, lesson_id, subject, status, approved_by, approved_at,
    rejection_reason, created_at, updated_at
) VALUES
    (
        'dddddddd-dddd-dddd-dddd-000000000d01',
        '77777777-7777-7777-7777-000000000703',
        'LEAVE',
        '2026-04-15 08:00:00+07',
        0,
        0,
        'Xin nghỉ có phép do tham gia hoạt động trường.',
        ARRAY['https://example.com/docs/leave-student-03.pdf']::text[],
        '66666666-6666-6666-6666-000000000601',
        '99999999-9999-9999-9999-000000000906',
        'Xin nghỉ buổi 6',
        'APPROVED',
        '11111111-1111-1111-1111-000000000101',
        '2026-04-15 12:00:00+07',
        '',
        NOW(),
        NOW()
    ),
    (
        'dddddddd-dddd-dddd-dddd-000000000d02',
        '77777777-7777-7777-7777-000000000706',
        'LATE',
        '2026-04-13 09:00:00+07',
        20,
        0,
        'Xin đi muộn do kẹt xe.',
        ARRAY['https://example.com/docs/late-student-06.jpg']::text[],
        '66666666-6666-6666-6666-000000000601',
        '99999999-9999-9999-9999-000000000905',
        'Xin đi muộn buổi 5',
        'APPROVED',
        '11111111-1111-1111-1111-000000000101',
        '2026-04-13 13:00:00+07',
        '',
        NOW(),
        NOW()
    ),
    (
        'dddddddd-dddd-dddd-dddd-000000000d03',
        '77777777-7777-7777-7777-000000000701',
        'LEAVE',
        '2026-04-22 09:30:00+07',
        0,
        0,
        'Xin nghỉ ôn thi học kỳ ở trường.',
        ARRAY['https://example.com/docs/pending-student-01.pdf']::text[],
        '66666666-6666-6666-6666-000000000601',
        '99999999-9999-9999-9999-000000000908',
        'Xin nghỉ buổi 8',
        'PENDING',
        NULL,
        NULL,
        '',
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO UPDATE SET
    student_id = EXCLUDED.student_id,
    leave_type = EXCLUDED.leave_type,
    apply_date = EXCLUDED.apply_date,
    late_minutes = EXCLUDED.late_minutes,
    early_minutes = EXCLUDED.early_minutes,
    reason = EXCLUDED.reason,
    documents = EXCLUDED.documents,
    class_id = EXCLUDED.class_id,
    lesson_id = EXCLUDED.lesson_id,
    subject = EXCLUDED.subject,
    status = EXCLUDED.status,
    approved_by = EXCLUDED.approved_by,
    approved_at = EXCLUDED.approved_at,
    rejection_reason = EXCLUDED.rejection_reason,
    updated_at = EXCLUDED.updated_at;

COMMIT;
