import uuid

def generate_sql():
    sql = "-- Data Seed Script cho EduCenter\n"
    sql += "-- Tự động sinh dữ liệu test chuẩn Edu Prod với ít nhất 25 records mỗi bảng.\n\n"
    
    # 1. Users (Admin để làm created_by_id)
    admin_id = "00000000-0000-0000-0000-000000000001"
    sql += f"INSERT INTO users (id, code, full_name, email, password, role, is_active) VALUES \n"
    sql += f"('{admin_id}', 'ADM-01', 'Quản trị viên', 'admin@educenter.com', 'hashed_pw', 'ADMIN', true)\n"
    sql += "ON CONFLICT DO NOTHING;\n\n"

    # 2. Rooms
    sql += "INSERT INTO rooms (id, code, name, capacity, address) VALUES\n"
    rooms = []
    room_ids = []
    for i in range(1, 26):
        r_id = f"11111111-1111-1111-1111-{str(i).zfill(12)}"
        room_ids.append(r_id)
        rooms.append(f"('{r_id}', 'P{100+i}', 'Phòng Học {100+i}', 30, 'Tầng {1 if i < 15 else 2} Tòa A')")
    sql += ",\n".join(rooms) + "\nON CONFLICT DO NOTHING;\n\n"

    # 3. Courses
    sql += "INSERT INTO courses (id, code, name, description, grade_level, subject, session_count, session_duration_minutes, total_hours, price, status) VALUES\n"
    courses = []
    course_ids = []
    for i in range(1, 26):
        c_id = f"22222222-2222-2222-2222-{str(i).zfill(12)}"
        course_ids.append(c_id)
        subjects = ['Toán', 'Tiếng Anh', 'Vật Lý', 'Hóa Học', 'Ngữ Văn']
        sub = subjects[i % 5]
        courses.append(f"('{c_id}', 'CRS-{str(i).zfill(3)}', 'Khóa học {sub} Tăng Cường {i}', 'Nội dung chuẩn kiến thức {sub}', 'Khối {8 + (i % 5)}', '{sub}', 24, 90, 36.00, {1500000 + (i*10000)}, 'ACTIVE')")
    sql += ",\n".join(courses) + "\nON CONFLICT DO NOTHING;\n\n"

    # 4. Programs
    sql += "INSERT INTO programs (id, code, name, track, created_by_id, status) VALUES\n"
    programs = []
    program_ids = []
    for i in range(1, 26):
        p_id = f"33333333-3333-3333-3333-{str(i).zfill(12)}"
        program_ids.append(p_id)
        tracks = ['BASIC', 'ADVANCED', 'SUPPORT']
        trk = tracks[i % 3]
        programs.append(f"('{p_id}', 'PRG-{str(i).zfill(3)}', 'Chương trình Đào tạo {i} ({trk})', '{trk}', '{admin_id}', 'ACTIVE')")
    sql += ",\n".join(programs) + "\nON CONFLICT DO NOTHING;\n\n"

    # 5. Teachers
    sql += "INSERT INTO teachers (id, code, full_name, email, phone, is_school_teacher, school_name, employment_type, status) VALUES\n"
    teachers = []
    teacher_ids = []
    first_names = ["Nguyễn", "Trần", "Lê", "Phạm", "Hoàng", "Huỳnh", "Phan", "Vũ", "Võ", "Đặng"]
    middle_names = ["Văn", "Thị", "Thanh", "Minh", "Thu", "Ngọc", "Hải", "Đức", "Bảo", "Lan"]
    last_names = ["A", "B", "C", "D", "E", "F", "G", "H", "I", "K", "Long", "Linh", "Hùng", "Hường", "Tùng", "Hoa"]
    for i in range(1, 26):
        t_id = f"44444444-4444-4444-4444-{str(i).zfill(12)}"
        teacher_ids.append(t_id)
        name = f"{first_names[i%10]} {middle_names[i%10]} {last_names[i%15]}"
        teachers.append(f"('{t_id}', 'GV-{str(i).zfill(3)}', 'GV {name}', 'gv{i}@educenter.com', '09{str(i).zfill(8)}', false, NULL, 'FULL_TIME', 'ACTIVE')")
    sql += ",\n".join(teachers) + "\nON CONFLICT DO NOTHING;\n\n"

    # 6. Students
    sql += "INSERT INTO students (id, code, full_name, email, phone, guardian_phone, grade_level, status) VALUES\n"
    students = []
    student_ids = []
    for i in range(1, 26):
        s_id = f"55555555-5555-5555-5555-{str(i).zfill(12)}"
        student_ids.append(s_id)
        name = f"{first_names[(i+2)%10]} {middle_names[(i+3)%10]} {last_names[(i+1)%15]}"
        students.append(f"('{s_id}', 'HS-{str(i).zfill(3)}', 'HS {name}', 'hs{i}@student.com', '08{str(i).zfill(8)}', '09{str(i+10).zfill(8)}', 'Khối {8 + (i % 5)}', 'ACTIVE')")
    sql += ",\n".join(students) + "\nON CONFLICT DO NOTHING;\n\n"

    # 7. Classes
    sql += "INSERT INTO classes (id, code, name, start_date, max_students, status, price, program_id, course_id, teacher_id) VALUES\n"
    classes = []
    for i in range(1, 26):
        cls_id = f"66666666-6666-6666-6666-{str(i).zfill(12)}"
        c_id = course_ids[i-1]
        p_id = program_ids[i-1]
        t_id = teacher_ids[i-1]
        classes.append(f"('{cls_id}', 'L-{str(i).zfill(3)}', 'Lớp Học {i}', NOW(), 20, 'OPEN', 2000000, '{p_id}', '{c_id}', '{t_id}')")
    sql += ",\n".join(classes) + "\nON CONFLICT DO NOTHING;\n\n"

    with open('seed_data.sql', 'w', encoding='utf-8') as f:
        f.write(sql)
    print("Created seed_data.sql successfully!")

if __name__ == '__main__':
    generate_sql()
