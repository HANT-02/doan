# GÓI TÀI LIỆU BA REVERSE-ENGINEER HỆ THỐNG EDUCENTER

**Nguồn sự thật chính:** codebase EduCenter tại thời điểm phân tích ngày 13/04/2026  
**Vai trò tài liệu:** BA / Product Owner / Tech Lead / Solution Architect  
**Nguyên tắc phân tích:** business-first, code-second; khi chưa đủ bằng chứng từ code sẽ đánh dấu rõ mức độ chắc chắn.

## Quy ước mức độ chắc chắn

- **Confirmed from code**: xác nhận trực tiếp từ entity, route, controller, use case, repository, UI hoặc migration đang có.
- **Strongly inferred from code**: suy luận mạnh từ cấu trúc dữ liệu, tên module, màn hình, route, seed hoặc comment.
- **Assumption / Needs BA validation**: giả định hợp lý nhưng chưa có bằng chứng đủ mạnh hoặc có dấu hiệu mâu thuẫn trong code.

---

# 1. Tổng quan hệ thống

## 1.1 Mục tiêu nghiệp vụ của hệ thống

EduCenter là hệ thống quản lý vận hành cho một trung tâm dạy thêm / trung tâm học thêm, tập trung vào số hóa các nghiệp vụ vận hành cốt lõi:

1. Quản lý người dùng và xác thực tài khoản.
2. Quản lý danh mục vận hành: học viên, giáo viên, khóa học, chương trình đào tạo, lớp học, phòng học, ca học.
3. Tổ chức lớp học và ghi danh học viên.
4. Xếp lịch học thông minh bằng nhiều solver và sinh buổi học (lesson) từ kết quả xếp lịch.
5. Theo dõi tài liệu giảng dạy và quy trình kiểm duyệt nội dung ở mức demo/stub.
6. Mở đường cho các bài toán phân tích dự báo học viên có nguy cơ học kém.

**Mức độ chắc chắn:** Confirmed from code

## 1.2 Phạm vi nghiệp vụ hiện tại

| ID | Phạm vi | Mô tả nghiệp vụ | Mức độ chắc chắn |
|---|---|---|---|
| SCOPE-01 | Xác thực và tài khoản | Đăng ký, OTP xác minh email, đăng nhập, refresh token, quên mật khẩu, đặt lại mật khẩu, đổi mật khẩu, xem hồ sơ | Confirmed from code |
| SCOPE-02 | Quản lý giáo viên | CRUD giáo viên, lịch dạy theo lesson, thống kê giờ dạy | Confirmed from code |
| SCOPE-03 | Quản lý học viên | CRUD học viên, tìm kiếm, lọc trạng thái | Confirmed from code |
| SCOPE-04 | Quản lý khóa học | CRUD khóa học, quản lý số buổi, thời lượng buổi, học phí | Confirmed from code |
| SCOPE-05 | Quản lý chương trình đào tạo | CRUD chương trình, gán/gỡ khóa học vào chương trình | Confirmed from code |
| SCOPE-06 | Quản lý lớp học | CRUD lớp, gán giáo viên, xem roster, thêm/xóa học viên khỏi lớp | Confirmed from code |
| SCOPE-07 | Quản lý phòng học | CRUD phòng học, quản lý sức chứa | Confirmed from code |
| SCOPE-08 | Quản lý ca học (Shift) | CRUD ca học chuẩn để phục vụ xếp lịch | Confirmed from code |
| SCOPE-09 | Xếp lịch thông minh | Tạo preview, xem preview, commit preview thành lesson, kiến trúc 3 solver + benchmark contract | Confirmed from code |
| SCOPE-10 | Lesson generation | Sinh buổi học từ preview scheduling | Confirmed from code |
| SCOPE-11 | Tài liệu giảng dạy và kiểm duyệt | Upload tài liệu, OCR/AI stub, hàng chờ duyệt, quyết định duyệt/từ chối | Confirmed from code |
| SCOPE-12 | Cổng người dùng theo vai trò | Admin portal, Teacher portal, Student portal, Compliance portal | Confirmed from code |

## 1.3 Phạm vi ngoài hệ thống hoặc chưa triển khai

| ID | Hạng mục | Hiện trạng | Ghi chú BA |
|---|---|---|---|
| OUT-01 | Attendance vận hành thật | Chỉ có entity/migration, chưa có API/controller/use case | Cần xác nhận có còn trong phase hiện tại không |
| OUT-02 | Lesson summary vận hành thật | Chỉ có entity/migration, chưa có API/controller/use case | Nếu là phạm vi báo cáo học tập thì đang thiếu toàn bộ flow |
| OUT-03 | Academic record vận hành thật | Chỉ có entity/migration, chưa có API/controller/use case | Chưa có workflow chấm điểm/đánh giá |
| OUT-04 | Leave request vận hành thật | Có entity/migration, không có API/màn hình thật | Frontend chỉ có placeholder ở student area |
| OUT-05 | Consultation / lead intake | Có entity/migration, không có API/màn hình | Miền tuyển sinh mới ở mức data model |
| OUT-06 | Reporting / dashboard thực | Màn overview chủ yếu dùng mock data | Không nên mô tả là báo cáo sản xuất |
| OUT-07 | ML dự báo học kém | Mới có dấu vết ở tài liệu/kế hoạch, chưa có code backend/frontend | Chưa được coi là implemented |
| OUT-08 | Benchmark scheduling thực thi thật | Có endpoint contract, chưa thấy metric benchmark thật ở use case benchmark | Cần tách rõ “đã có kiến trúc” và “chưa có kết quả benchmark” |
| OUT-09 | AI moderation production | OCR/AI hiện là stub keyword-based, không phải model production | Chỉ nên mô tả demo workflow |
| OUT-10 | Payment / billing / công nợ | Không thấy code nghiệp vụ tài chính | Ngoài phạm vi hiện tại |

**Mức độ chắc chắn:** Confirmed from code

## 1.4 Giá trị nghiệp vụ mang lại

| ID | Giá trị | Giải thích business |
|---|---|---|
| VAL-01 | Chuẩn hóa vận hành trung tâm | Gom quản lý học viên, giáo viên, khóa học, chương trình, lớp, phòng và ca vào một hệ thống |
| VAL-02 | Giảm công xếp lịch thủ công | Cho phép preview trước khi commit, phát hiện xung đột giáo viên/phòng/lớp |
| VAL-03 | Tăng khả năng kiểm soát dữ liệu lớp học | Roster, sĩ số tối đa, gán giáo viên, khóa học và chương trình được gắn với lớp |
| VAL-04 | Tăng khả năng truy vết lịch dạy | Lesson là sản phẩm vận hành đầu ra từ scheduling và là nền cho attendance/tổng kết |
| VAL-05 | Chuẩn hóa quy trình tài liệu giảng dạy | Có upload, nhãn AI, audit log, quyết định duyệt/từ chối |
| VAL-06 | Chuẩn bị nền dữ liệu cho phân tích nâng cao | Entities lesson, attendance, academic_record, leave_request cho thấy định hướng phân tích học tập và dự báo |

**Mức độ chắc chắn:** Strongly inferred from code

## 1.5 Bên liên quan chính

- Quản trị viên trung tâm.
- Nhân sự vận hành đào tạo.
- Giáo viên.
- Học viên.
- Phụ huynh: chỉ có dấu vết gián tiếp qua `guardian_phone`, chưa thấy actor riêng.
- Người duyệt tài liệu / compliance officer.
- Chủ dự án / Product Owner.
- Tech Lead / đội kỹ thuật.
- Hệ thống SMTP.
- Engine scheduling.

**Mức độ chắc chắn:**  
- Admin, Teacher, Student, Compliance, SMTP, scheduling engine: Confirmed from code  
- Phụ huynh, vận hành đào tạo: Strongly inferred from code

## 1.6 Những giả định lớn về mô hình vận hành trung tâm dạy thêm

| ID | Giả định | Giải thích | Mức độ chắc chắn |
|---|---|---|---|
| ASM-01 | Trung tâm vận hành theo mô hình single-tenant | Không thấy tenant, branch, campus, organization ownership trong data model | Strongly inferred from code |
| ASM-02 | Một lớp gắn tối đa một khóa học chính và một giáo viên chính tại một thời điểm | `classes` có `course_id`, `teacher_id` kiểu 1-1 | Confirmed from code |
| ASM-03 | Một lớp có thể có nhiều lịch tuần chuẩn | `class_schedules` là quan hệ 1-N với class | Confirmed from code |
| ASM-04 | Một buổi học (lesson) là phiên bản lịch thực thi sau khi scheduling commit | `lessons` được tạo từ preview commit | Confirmed from code |
| ASM-05 | Trung tâm muốn tối ưu xếp lịch dựa trên phòng, ca, giáo viên, lớp và thời lượng khóa học | Chính là input của scheduling preview | Confirmed from code |
| ASM-06 | Có nhu cầu quy trình duyệt tài liệu giảng dạy nội bộ | Có materials, labels, audit_logs, approval_decisions, compliance queue page | Confirmed from code |
| ASM-07 | Hệ thống hướng tới theo dõi học tập sâu hơn trong tương lai | Có attendance, lesson_summary, academic_record, leave_request nhưng chưa triển khai API | Strongly inferred from code |

---

# 2. Bản đồ miền nghiệp vụ

| Mã miền | Tên miền | Mục đích nghiệp vụ | Bài toán thực tế giải quyết | Tác nhân liên quan | Thực thể cốt lõi | API/Entrypoint chính | Dịch vụ / job / integration liên quan | Trạng thái triển khai | Ghi chú cho BA | Mức độ chắc chắn |
|---|---|---|---|---|---|---|---|---|---|---|
| DOM-01 | Xác thực và tài khoản | Quản lý vòng đời tài khoản người dùng | Tạo tài khoản, xác minh email, đăng nhập, khôi phục truy cập | Guest, User, System SMTP | User, UserOTP, PasswordReset | `/v1/auth/*`, `/v2/auth/*` | Mailer SMTP, password cipher/hasher, JWT | Implemented nhưng có lỗi `/me` và V2 incomplete | `/v1/auth/me` đang sai kiểu ID; v2 mới có login/logout/refresh | Confirmed from code |
| DOM-02 | Học viên | Quản lý hồ sơ học viên | Lưu thông tin cá nhân, phụ huynh, khối lớp, trạng thái | Admin, Student, Parent (gián tiếp) | Student | `/v1/students` | Student use cases, student repo | Implemented nhưng auth gap | CRUD chỉ cần auth, không check ADMIN | Confirmed from code |
| DOM-03 | Giáo viên | Quản lý hồ sơ giảng dạy | Theo dõi đội ngũ giáo viên và lịch/giờ dạy | Admin, Teacher | Teacher, Lesson | `/v1/teachers`, `/v1/teachers/:id/timetable`, `/v1/teachers/:id/stats/teaching-hours` | Teacher use cases, lesson queries | Implemented | Có lịch dạy và thống kê giờ dạy nhưng chưa có attendance/journal teacher thật | Confirmed from code |
| DOM-04 | Khóa học | Quản lý đơn vị học tập cơ sở | Lưu môn học, số buổi, thời lượng, học phí | Admin | Course | `/v1/courses` | Course use cases | Implemented partial | Search trong use case list chưa dùng rõ; validation còn mỏng | Confirmed from code |
| DOM-05 | Chương trình đào tạo | Gom khóa học theo track/chương trình | Xây khung đào tạo và liên kết nhiều khóa học | Admin, quản lý đào tạo | Program, ProgramCourse, Objective, Outcome | `/v1/programs`, `/v1/programs/:id/courses` | Program use cases | Program implemented; Objective/Outcome missing API | Có mâu thuẫn field `status` giữa migration và entity | Confirmed from code |
| DOM-06 | Lớp học | Quản lý đơn vị vận hành lớp | Mở lớp, quy định sĩ số, gán khóa học/chương trình/giáo viên/phòng | Admin | Class, ClassSchedule | `/v1/classes` | Class use cases | Implemented | Room gắn với class là optional; lifecycle chỉ OPEN/CLOSED/CANCELLED | Confirmed from code |
| DOM-07 | Ghi danh | Đưa học viên vào lớp | Ghi danh hàng loạt và lấy roster lớp | Admin | Enrollment | `/v1/classes/:id/students` | Enroll/remove/roster use cases | Partial | Hiện chưa đếm enrollment hiện có, chưa có approval flow APPLIED->ENROLLED | Confirmed from code |
| DOM-08 | Phòng học | Quản lý tài nguyên phòng | Quản lý sức chứa phòng cho lớp và scheduling | Admin | Room | `/v1/rooms` | Room use cases | Partial | DTO có `status/location` nhưng entity là `code/address`, use case không dùng status/location | Confirmed from code |
| DOM-09 | Ca học | Chuẩn hóa khung giờ vận hành | Tạo các ca sáng/chiều/tối/custom để lịch lớp tham chiếu | Admin | Shift | `/v1/shifts` | Shift use cases | Implemented | Đây là nền dữ liệu bắt buộc cho scheduling mới | Confirmed from code |
| DOM-10 | Xếp lịch thông minh | Sinh lịch buổi học tối ưu | Tránh trùng lớp, giáo viên, phòng; dùng shift và class_schedule | Admin, scheduling engine | Class, Shift, Room, Lesson, PreviewResult | `/v1/scheduling/preview`, `/preview/latest`, `/preview/:id`, `/commit`, `/benchmark` | Scheduling solver services, PreviewStore | Implemented partial | Có 3 solver service + commit thật; benchmark use case mới trả contract | Confirmed from code |
| DOM-11 | Buổi học | Đại diện cho phiên học thực tế | Là đầu ra được commit từ preview scheduling | Admin, Teacher | Lesson | Không có controller riêng | Lesson repository, commit preview | Partial | Chưa có CRUD/manual ops hay màn lesson management riêng | Confirmed from code |
| DOM-12 | Điểm danh | Theo dõi chuyên cần từng buổi | Ghi nhận tham dự / vắng / đi muộn / có phép | Teacher, Student, Admin | Attendance | Không có API/controller | Chỉ có entity/migration | Missing API | Chưa có enum status, chưa có quy trình sử dụng | Confirmed from code |
| DOM-13 | Tổng kết buổi học | Ghi nhận nội dung đã dạy | Ghi topic, homework, feedback, note giáo viên | Teacher | LessonSummary | Không có API/controller | Chỉ có entity/migration | Missing API | Có quan hệ 1-1 lesson nhưng chưa có flow | Confirmed from code |
| DOM-14 | Kết quả học tập | Đánh giá học viên sau buổi học | Lưu điểm bài tập, thái độ, tham gia, tổng điểm | Teacher, Student, phụ huynh (gián tiếp) | AcademicRecord | Không có API/controller | Chỉ có entity/migration | Missing API | Là nền tốt cho predictive analytics nhưng hiện chưa dùng | Confirmed from code |
| DOM-15 | Đơn xin phép | Xin nghỉ/đi muộn/về sớm | Gắn lý do và quyết định duyệt/từ chối | Student, Admin | LeaveRequest | Không có API/controller | Chỉ có entity/migration | Missing API | Có thiết kế workflow trong data model nhưng chưa triển khai | Confirmed from code |
| DOM-16 | Tài liệu giảng dạy và kiểm duyệt | Kiểm soát tài liệu dùng trong giảng dạy | Upload, phân tích OCR/AI, duyệt thủ công, tải file | Teacher, Compliance, Admin | Material, Label, AuditLog, ApprovalDecision | `/v1/materials/*` | Local file storage, OCR stub, Gemini stub, download use case | Partial | Quyền duyệt và ownership còn lỏng; runtime migration không auto tạo bảng module này | Confirmed from code |
| DOM-17 | Tư vấn tuyển sinh | Quản lý lead / học viên tiềm năng | Lưu đầu mối tư vấn đầu vào | Tư vấn viên, marketing | Consultation | Không có API/controller | Chỉ có entity/migration | Missing API | Miền này mới ở mức data model | Confirmed from code |
| DOM-18 | Báo cáo/thống kê | Hỗ trợ điều hành | Dashboard tổng quan, xung đột, số liệu học tập | Admin, quản lý | Mock stats, teaching hours | Admin overview, teacher stats | Recharts mock, teaching hours use case | Partial / mock | Dashboard chính chưa dùng dữ liệu thật | Confirmed from code |
| DOM-19 | Tích hợp hệ thống ngoài | Email, queue, AI/OCR stub | Gửi OTP/reset mail, hỗ trợ gửi message nền, AI/OCR giả lập | SMTP, queue, audit provider | Mail, queue adapter | N/A | SMTP, RabbitMQ/Kafka abstractions, audit stub | Partial | Queue có abstraction nhưng nhiều phần chưa dùng trong use case thực | Confirmed from code |
| DOM-20 | ML / predictive analytics | Dự báo học viên nguy cơ học kém | Chưa có implementation, mới có dấu vết trong docs/kế hoạch | Admin, quản lý đào tạo | (dự kiến) AcademicRecord, Attendance, LessonSummary, Student | Chưa có API | Chưa có service/model code | Planned only | Không nên ghi implemented trong hồ sơ BA hiện trạng | Confirmed from code |

---

# 3. Danh mục tác nhân và bên liên quan

## 3.1 Danh mục tác nhân

| Mã tác nhân | Tên | Loại | Mục tiêu nghiệp vụ | Trách nhiệm | Các use case tham gia | Quyền thao tác | Dấu hiệu nhận biết trong code | Mức độ chắc chắn |
|---|---|---|---|---|---|---|---|---|
| ACT-01 | Quản trị viên | Nội bộ | Vận hành toàn bộ trung tâm | Quản lý danh mục, mở lớp, ghi danh, xếp lịch, commit lesson | Hầu hết use case admin | Full trên teachers/courses/programs/classes/rooms/shifts/scheduling | `RoleMiddleware("ADMIN")`, admin routes, admin pages | Confirmed from code |
| ACT-02 | Giáo viên | Nội bộ | Giảng dạy và cung cấp tài liệu | Upload tài liệu, xem lớp/lịch dạy | Upload material, xem timetable/stats | Upload material; đọc một số màn teacher | Teacher routes/pages, teacher role in nav | Confirmed from code |
| ACT-03 | Học viên | Nội bộ | Theo dõi việc học | Xem dashboard học viên; tương lai xem thời khóa biểu, kết quả, xin phép | Student overview; planned leave/results/timetable | Hiện gần như đọc UI placeholder | Student overview route, default role STUDENT | Confirmed from code |
| ACT-04 | Người dùng chưa đăng nhập | Bên ngoài | Truy cập chức năng công khai và đăng ký tài khoản | Đăng ký, login, forgot/reset; đọc một số danh mục công khai | Register/login/forgot/reset; public GET list/detail | Không cần auth cho nhiều GET route | Public GET routes, auth register/login pages | Confirmed from code |
| ACT-05 | Compliance officer / người duyệt nội dung | Nội bộ | Duyệt tài liệu bị gắn cờ | Approve/reject tài liệu sau AI audit | Review material | Theo business thì phải có quyền duyệt | `approval_decisions`, compliance queue page, `compliance_officer_id` | Strongly inferred from code |
| ACT-06 | Quản lý đào tạo | Nội bộ | Quản trị chương trình, khóa học, chuẩn đầu ra, lịch học | Có thể là persona nghiệp vụ của Admin | Programs, courses, scheduling | Chưa có role riêng | Program/Objective/Outcome domain | Strongly inferred from code |
| ACT-07 | Nhân viên vận hành | Nội bộ | Vận hành lớp, roster, phòng, ca | Có thể dùng cùng role Admin ở hiện trạng | Student/class/enrollment/shift/room | Chưa có role riêng | Admin pages/use cases | Strongly inferred from code |
| ACT-08 | Phụ huynh | Bên ngoài liên quan nghiệp vụ | Theo dõi liên hệ và hỗ trợ học viên | Chưa có use case trực tiếp trong code | Chưa có | Không có role riêng | `guardian_phone` trên Student | Strongly inferred from code |
| ACT-09 | Hệ thống SMTP | Hệ thống | Gửi OTP và link đặt lại mật khẩu | Gửi mail ra ngoài | Register, forgot password | Không phải user role | `internal/services/mailer/mailer.go` | Confirmed from code |
| ACT-10 | Scheduling engine | Hệ thống | Sinh phương án xếp lịch và lesson | Chạy solver, lưu preview, commit lesson | Preview, benchmark, commit | Nội bộ use case/service | `internal/services/scheduling/*` | Confirmed from code |
| ACT-11 | OCR stub | Hệ thống | Trích xuất text từ tài liệu | Tạo raw OCR text cho audit log | Upload material | Nội bộ service | `StubOCRService` | Confirmed from code |
| ACT-12 | AI moderation stub | Hệ thống | Gán nhãn SAFE/WARNING/DANGER | Tạo audit reasoning | Upload material | Nội bộ service | `StubGeminiService` | Confirmed from code |
| ACT-13 | Queue infrastructure | Hệ thống | Hạ tầng hỗ trợ async | Chưa thấy use case core dùng thực tế | Tiềm năng | N/A | RabbitMQ/Kafka abstractions | Strongly inferred from code |
| ACT-14 | Product Owner / BA | Bên liên quan | Chốt phạm vi, rule, lifecycle, quyền | Xác nhận điểm mơ hồ và gap | Không phải actor hệ thống | N/A | Không có trong code | Assumption / Needs BA validation |

## 3.2 Ma trận Actor -> Domain

| Actor\Domain | Auth | Student | Teacher | Course | Program | Class | Enrollment | Room | Shift | Scheduling | Lesson | Attendance | Summary | Academic | Leave | Material | Consultation | Reports | ML |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| ACT-01 Admin | X | X | X | X | X | X | X | X | X | X | X | X? | X? | X? | X? | X | ? | X | X? |
| ACT-02 Teacher | X |  | X |  |  |  |  |  |  |  | X | X? | X? | X? |  | X |  |  |  |
| ACT-03 Student | X | X |  |  |  | X? |  |  |  |  | X? | X? | X? | X? | X? |  |  |  |  |
| ACT-04 Guest | X |  | X? | X? | X? | X? |  | X? |  |  |  |  |  |  |  |  |  |  |  |
| ACT-05 Compliance | X |  |  |  |  |  |  |  |  |  |  |  |  |  |  | X |  |  |  |

Ghi chú:
- `X` = tham gia rõ từ code.
- `X?` = suy luận mạnh từ code/UI nhưng chưa có workflow đủ rõ.
- `?` = data model có nhưng chưa có role/use case thật.

## 3.3 Ma trận Actor -> Use Case sơ bộ

| Mã UC | Tên use case | Guest | Admin | Teacher | Student | Compliance | Hệ thống |
|---|---|---|---|---|---|---|---|
| UC-001 | Đăng ký tài khoản | X |  |  |  |  | SMTP |
| UC-002 | Xác minh OTP email | X |  |  |  |  | SMTP |
| UC-003 | Đăng nhập | X | X | X | X | X |  |
| UC-004 | Refresh token | X | X | X | X | X |  |
| UC-005 | Đăng xuất |  | X | X | X | X |  |
| UC-006 | Quên mật khẩu | X | X | X | X | X | SMTP |
| UC-007 | Đặt lại mật khẩu | X | X | X | X | X |  |
| UC-008 | Đổi mật khẩu |  | X | X | X | X |  |
| UC-009 | Xem hồ sơ hiện tại |  | X | X | X | X |  |
| UC-010 | Quản lý giáo viên |  | X |  |  |  |  |
| UC-011 | Quản lý học viên |  | X | ? | ? | ? |  |
| UC-012 | Quản lý khóa học | ? | X |  |  |  |  |
| UC-013 | Quản lý chương trình | ? | X |  |  |  |  |
| UC-014 | Gán khóa học vào chương trình |  | X |  |  |  |  |
| UC-015 | Quản lý lớp học | ? | X |  |  |  |  |
| UC-016 | Ghi danh học viên |  | X |  |  |  |  |
| UC-017 | Gỡ học viên khỏi lớp |  | X |  |  |  |  |
| UC-018 | Gán giáo viên cho lớp |  | X |  |  |  |  |
| UC-019 | Quản lý phòng học | ? | X |  |  |  |  |
| UC-020 | Quản lý ca học |  | X |  |  |  |  |
| UC-021 | Tạo preview scheduling |  | X |  |  |  | Scheduling engine |
| UC-022 | Xem preview scheduling |  | X |  |  |  | Scheduling engine |
| UC-023 | Commit preview thành lesson |  | X |  |  |  | Scheduling engine |
| UC-024 | Benchmark solver scheduling |  | X |  |  |  | Scheduling engine |
| UC-025 | Xem timetable giáo viên | ? | X | X |  |  |  |
| UC-026 | Xem thống kê giờ dạy | ? | X | X |  |  |  |
| UC-027 | Upload tài liệu |  |  | X |  |  | OCR/AI stub |
| UC-028 | Xem tài liệu |  | X | X | ? | X |  |
| UC-029 | Tải tài liệu |  | X | X | ? | X | File storage |
| UC-030 | Duyệt / từ chối tài liệu |  | ? | ? | ? | X | OCR/AI stub |

---

# 4. Cây phân rã chức năng

## 4.1 Cấu trúc phân rã chức năng

| Mã chức năng | Tên chức năng | Chức năng cha | Cấp | Mô tả nghiệp vụ | Tác nhân thực hiện | Trigger | Input chính | Output chính | Dữ liệu liên quan | Điều kiện / quy tắc liên quan | Trạng thái triển khai | API / service / module liên quan |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| FUN-01 | Quản lý truy cập hệ thống | ROOT | 1 | Quản lý vòng đời đăng ký, xác minh, đăng nhập, khôi phục truy cập | Guest, User | Người dùng gửi form auth | Email, mật khẩu, OTP, token | Token, profile, trạng thái tài khoản | User, UserOTP, PasswordReset | OTP hết hạn 5 phút; user phải active mới login | Implemented partial | Auth controllers, auth use cases, mailer |
| FUN-01.01 | Đăng ký tài khoản | FUN-01 | 2 | Tạo user mới ở trạng thái chưa kích hoạt | Guest | Submit form register | Email, full_name, password_enc | user_id mới | User | Email unique, tạo OTP, inactive | Implemented | `/v1/auth/register` |
| FUN-01.02 | Xác minh email OTP | FUN-01 | 2 | Kích hoạt tài khoản qua OTP | Guest/User | Submit OTP | user_id, otp | Account active | UserOTP, User | OTP phải còn hạn và chưa dùng | Implemented | `/v1/auth/verify-otp` |
| FUN-01.03 | Đăng nhập / refresh / logout | FUN-01 | 2 | Thiết lập và gia hạn session JWT | User | Login/refresh/logout request | email, password / refresh_token / token | access_token, refresh_token | User | User phải active | Implemented partial | `/v1/auth/login`, `/refresh`, `/logout` |
| FUN-01.04 | Khôi phục và đổi mật khẩu | FUN-01 | 2 | Quên mật khẩu, reset, đổi mật khẩu khi đã login | Guest/User | Forgot/reset/change request | email / token / old+new password | trạng thái thành công | PasswordReset, User | Tránh user enumeration; token có hạn | Implemented | `/forgot-password`, `/reset-password`, `/change-password` |
| FUN-01.05 | Xem hồ sơ hiện tại | FUN-01 | 2 | Trả về profile người dùng hiện tại | User | Open profile page | JWT token | User profile | User | Hiện tại bị lỗi do kiểu ID int/uuid | Partial / defective | `/v1/auth/me`, `ProfilePage` |
| FUN-02 | Quản lý danh mục vận hành | ROOT | 1 | Quản lý master data cho hoạt động đào tạo | Admin | Truy cập màn hình admin | form CRUD | bản ghi mới/cập nhật/xóa mềm | Teacher, Student, Course, Program, Room, Shift | Quyền chủ yếu ADMIN | Implemented partial | Admin pages + CRUD APIs |
| FUN-02.01 | Quản lý giáo viên | FUN-02 | 2 | Quản lý hồ sơ giảng dạy | Admin | CRUD giáo viên | code, full_name, email, phone, loại hình, trạng thái | Teacher | Teacher | Email/code unique nếu có; mặc định PART_TIME/ACTIVE | Implemented | `/v1/teachers` |
| FUN-02.01.01 | Xem lịch dạy giáo viên | FUN-02.01 | 3 | Xem lesson của giáo viên theo khoảng ngày | Admin, Teacher | Mở chi tiết/timetable | teacher_id, from, to | lessons | Lesson, Teacher | teacher phải tồn tại | Implemented | `/v1/teachers/:id/timetable` |
| FUN-02.01.02 | Xem thống kê giờ dạy | FUN-02.01 | 3 | Thống kê tổng giờ theo ngày/tuần/tháng | Admin, Teacher | Gọi stats | teacher_id, from, to, group_by | total_hours, breakdown | Lesson | group_by chỉ day/week/month | Implemented | `/v1/teachers/:id/stats/teaching-hours` |
| FUN-02.02 | Quản lý học viên | FUN-02 | 2 | Quản lý hồ sơ học viên và liên hệ phụ huynh | Admin | CRUD học viên | code, name, grade, phone, guardian_phone, status | Student | Student | Validation mỏng, thiếu role check | Implemented with auth gap | `/v1/students` |
| FUN-02.03 | Quản lý khóa học | FUN-02 | 2 | Quản lý khóa học/môn học | Admin | CRUD khóa học | code, name, subject, session_count, duration, price | Course | Course | code+name required, status mặc định ACTIVE | Implemented partial | `/v1/courses` |
| FUN-02.04 | Quản lý chương trình đào tạo | FUN-02 | 2 | Quản lý tập khóa học theo track | Admin | CRUD chương trình | code, name, track, hiệu lực | Program | Program, ProgramCourse | code+name required | Implemented partial | `/v1/programs` |
| FUN-02.04.01 | Liên kết khóa học với chương trình | FUN-02.04 | 3 | Thêm/gỡ nhiều khóa học cho một chương trình | Admin | Chọn khóa học | program_id, course_ids[] | mapping được cập nhật | ProgramCourse | Cần chương trình tồn tại | Implemented | `/v1/programs/:id/courses` |
| FUN-02.05 | Quản lý phòng học | FUN-02 | 2 | Quản lý tài nguyên phòng cho lớp và scheduling | Admin | CRUD phòng | name, capacity | Room | Room | Dữ liệu DTO/entity đang lệch | Implemented partial | `/v1/rooms` |
| FUN-02.06 | Quản lý ca học | FUN-02 | 2 | Định nghĩa ca học chuẩn | Admin | CRUD shift | code, name, start/end, duration, type, active | Shift | Shift | session_type enum MORNING/AFTERNOON/EVENING/CUSTOM | Implemented | `/v1/shifts` |
| FUN-03 | Vận hành lớp học | ROOT | 1 | Tạo và quản lý đơn vị lớp thực tế | Admin | Mở lớp / cập nhật / roster | form lớp, danh sách học viên | Class, Enrollment | Class, Enrollment, Teacher, Room | Class status OPEN/CLOSED/CANCELLED | Implemented partial | Class APIs + dialogs |
| FUN-03.01 | Tạo/cập nhật/xóa lớp | FUN-03 | 2 | Quản lý thông tin lớp | Admin | CRUD lớp | code, name, dates, max_students, course, program, teacher | Class | Class | status mặc định OPEN | Implemented | `/v1/classes` |
| FUN-03.02 | Ghi danh học viên vào lớp | FUN-03 | 2 | Thêm nhiều học viên vào roster lớp | Admin | Add students | class_id, student_ids[] | Enrollment records | Enrollment, Class, Student | Chưa đếm enrollment hiện có, chưa kiểm duplicate | Partial | `/v1/classes/:id/students` |
| FUN-03.03 | Gỡ học viên khỏi lớp | FUN-03 | 2 | Xóa học viên khỏi roster lớp | Admin | Remove students | class_id, student_ids[] | enrollment deleted | Enrollment | Không có validate lifecycle | Partial | `DELETE /v1/classes/:id/students` |
| FUN-03.04 | Gán giáo viên cho lớp | FUN-03 | 2 | Chọn giáo viên phụ trách chính của lớp | Admin | Assign teacher | class_id, teacher_id | class.teacher_id updated | Class, Teacher | Không thấy rule về teacher workload tại bước gán | Implemented | `PUT /v1/classes/:id/teacher` |
| FUN-03.05 | Xem roster lớp | FUN-03 | 2 | Xem danh sách học viên lớp và snapshot capacity | Admin | Open class detail | class_id | students[], capacity_limit | Class, Enrollment, Student, Room | capacity_limit = min(max_students, room.capacity) | Implemented | `GET /v1/classes/:id/students` |
| FUN-04 | Xếp lịch thông minh | ROOT | 1 | Tạo phương án lịch học và sinh lesson | Admin, Scheduling engine | Run preview / commit / benchmark | filters lớp, giáo viên, phòng, ngày | preview, conflicts, lessons | Class, Shift, Room, Lesson | Hard constraints scheduling | Implemented partial | Scheduling APIs, scheduling services |
| FUN-04.01 | Tạo preview xếp lịch | FUN-04 | 2 | Chạy solver trên tập lớp OPEN | Admin | Click run preview | date range, class_ids, teacher_ids, room_ids | PreviewResult | Class, Shift, Room | date_to >= date_from; cần shift active | Implemented | `POST /v1/scheduling/preview` |
| FUN-04.02 | Xem preview hiện tại / mới nhất | FUN-04 | 2 | Lấy preview theo ID hoặc bản mới nhất | Admin | Refresh or open preview | run_id | PreviewResult | PreviewStore | Preview hiện lưu in-memory | Implemented partial | `GET /preview/:id`, `/preview/latest` |
| FUN-04.03 | Commit preview thành lesson | FUN-04 | 2 | Ghi kết quả preview thành lesson thật trong DB | Admin | Click commit | run_id | số lesson tạo được | Lesson | Chỉ commit khi COMPLETED và không còn unscheduled/conflict | Implemented | `POST /v1/scheduling/commit` |
| FUN-04.04 | Benchmark solver scheduling | FUN-04 | 2 | Chuẩn bị API benchmark cho 3 solver | Admin | Call benchmark API | same filters as preview | contract benchmark | Solver catalog | Chưa có benchmark metric thật | Partial | `POST /v1/scheduling/benchmark` |
| FUN-05 | Quản lý tài liệu giảng dạy | ROOT | 1 | Upload, phân loại và duyệt tài liệu | Teacher, Compliance | Upload/review/download | file + metadata, review decision | Material + audit trail | Material, Label, AuditLog, ApprovalDecision | File type/size validation; AI stub | Implemented partial | Material APIs, teacher/compliance pages |
| FUN-05.01 | Upload tài liệu | FUN-05 | 2 | Giáo viên tải tài liệu lên hệ thống | Teacher | Submit form upload | teacher_id, title, description, file | Material mới + audit log | Material | Chưa derive teacher_id từ auth context | Partial | `POST /v1/materials/upload` |
| FUN-05.02 | Phân tích OCR + AI | FUN-05 | 2 | OCR text và gán nhãn nội dung | System | Trigger sau upload | file bytes | audit log, label, material status AI_REVIEWED | AuditLog, Label | OCR/AI hiện là stub keyword-based | Partial/demo | `internal/services/audit/services.go` |
| FUN-05.03 | Xem danh sách / chi tiết tài liệu | FUN-05 | 2 | Xem lịch sử audit và quyết định duyệt | Teacher, Compliance, Admin | Open list/detail | filters teacher/status/queue | materials/detail | Material, AuditLog, ApprovalDecision | Ownership chưa enforce | Implemented partial | `GET /v1/materials`, `GET /v1/materials/:id` |
| FUN-05.04 | Tải tài liệu | FUN-05 | 2 | Download file qua endpoint nội bộ | Teacher, Compliance, Admin | Click download | material_id | file attachment | Material, file storage | Quyền tải chưa chặt | Implemented | `GET /v1/materials/:id/download` |
| FUN-05.05 | Duyệt / từ chối tài liệu | FUN-05 | 2 | Ghi quyết định cuối cùng | Compliance, Admin | Submit review form | material_id, compliance_officer_id, approved, note | status APPROVED/REJECTED | ApprovalDecision, Material | Hiện bất kỳ user auth cũng có thể gọi review endpoint | Partial / control gap | `POST /v1/materials/:id/review` |
| FUN-06 | Theo dõi học tập và vận hành buổi học | ROOT | 1 | Điểm danh, tổng kết, kết quả, xin phép, lead | Teacher, Student, Admin | Các tác vụ học thuật sau lesson | attendance, summary, score, leave | records vận hành học tập | Attendance, LessonSummary, AcademicRecord, LeaveRequest, Consultation | Chưa có API/màn hình thật | Missing API | Chỉ có entities/migrations |

## 4.2 Gợi ý sơ đồ phân rã chức năng BA nên vẽ

1. **Phân hệ Tài khoản và Truy cập**
2. **Phân hệ Danh mục Đào tạo**
3. **Phân hệ Vận hành Lớp học**
4. **Phân hệ Xếp lịch và Sinh Buổi học**
5. **Phân hệ Tài liệu Giảng dạy**
6. **Phân hệ Theo dõi Học tập**
7. **Phân hệ Báo cáo và Phân tích**

---

# 5. Danh mục use case tổng thể

## 5.1 Catalog use case

| Mã use case | Tên use case | Mục tiêu | Tác nhân chính | Tác nhân phụ / hệ thống hỗ trợ | Miền nghiệp vụ | Mô tả ngắn | Tiền điều kiện | Trigger | Hậu điều kiện | Dữ liệu liên quan | Quy tắc nghiệp vụ chính | API / màn hình / module liên quan | Độ tin cậy |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| UC-001 | Đăng ký tài khoản | Tạo tài khoản mới | Guest | SMTP | Auth | Tạo user mới và gửi OTP xác minh | Email chưa tồn tại | Gửi form register | User tạo ở trạng thái inactive, OTP được tạo | User, UserOTP | Email unique; OTP 5 phút | `/v1/auth/register`, RegisterPage | Confirmed from code |
| UC-002 | Xác minh OTP email | Kích hoạt tài khoản | Guest/User |  | Auth | Người dùng nhập OTP để activate account | Có user_id và OTP active | Gửi form OTP | User được active, OTP used | UserOTP, User | OTP phải đúng và chưa hết hạn | `/v1/auth/verify-otp` | Confirmed from code |
| UC-003 | Đăng nhập | Nhận token truy cập | User | JWT service | Auth | Login bằng email + password | User active | Submit login form | access_token + refresh_token được cấp | User | Chỉ user active mới login | `/v1/auth/login`, LoginPage | Confirmed from code |
| UC-004 | Refresh token | Gia hạn access token | User | JWT service | Auth | Dùng refresh token để lấy access token mới | Refresh token hợp lệ | Token refresh request | Access token mới | User | Không có session store; phụ thuộc chữ ký JWT | `/v1/auth/refresh` | Confirmed from code |
| UC-005 | Đăng xuất | Kết thúc phiên làm việc logic | User | JWT service | Auth | Xác nhận token hợp lệ rồi trả success | Có token | Logout request | Client xóa session local | JWT token | Không có blacklist/revoke thật | `/v1/auth/logout`, ProfilePage | Confirmed from code |
| UC-006 | Quên mật khẩu | Khởi tạo flow reset password | Guest/User | SMTP | Auth | Gửi link reset nếu email tồn tại | Email được nhập | Forgot password form | PasswordReset record được tạo nếu email tồn tại | PasswordReset, User | Không được lộ user enumeration | `/v1/auth/forgot-password`, ForgotPasswordPage | Confirmed from code |
| UC-007 | Đặt lại mật khẩu | Cập nhật mật khẩu mới | Guest/User |  | Auth | Dùng token reset để đổi password | Token reset còn hạn | Reset password form | Password mới được lưu, token marked used | PasswordReset, User | Password tối thiểu 8 ký tự ở use case reset | `/v1/auth/reset-password`, ResetPasswordPage | Confirmed from code |
| UC-008 | Đổi mật khẩu | Đổi mật khẩu khi đang đăng nhập | User |  | Auth | So sánh mật khẩu cũ rồi cập nhật mật khẩu mới | User đã auth | Change password form | Password mới được lưu | User | old password phải đúng | `/v1/auth/change-password`, ChangePasswordPage | Confirmed from code |
| UC-009 | Xem hồ sơ hiện tại | Xem profile account | User |  | Auth | Trả về user hiện tại theo JWT | User đã auth | Open profile | Profile được trả về | User | Hiện lỗi do ép UUID sang int | `/v1/auth/me`, ProfilePage | Confirmed from code |
| UC-010 | Tạo giáo viên | Thêm giáo viên mới | Admin |  | Teacher | Tạo hồ sơ giáo viên | Admin auth | Submit form teacher | Teacher mới được tạo | Teacher | FullName required; email/code unique nếu nhập; default PART_TIME/ACTIVE | `/v1/teachers`, TeachersPage/TeacherFormPage | Confirmed from code |
| UC-011 | Xem danh sách giáo viên | Tìm/lọc giáo viên | Admin/Guest |  | Teacher | Danh sách giáo viên với filter và phân trang | None | Open teacher list | List giáo viên | Teacher | Search name/email/phone/code, filter status/employment type | `GET /v1/teachers`, TeachersPage | Confirmed from code |
| UC-012 | Cập nhật giáo viên | Sửa hồ sơ giáo viên | Admin |  | Teacher | Cập nhật một phần teacher | Teacher tồn tại | Edit teacher | Teacher được cập nhật | Teacher | Role ADMIN | `/v1/teachers/:id`, TeacherDetailPage | Confirmed from code |
| UC-013 | Xóa mềm giáo viên | Ngừng sử dụng giáo viên | Admin |  | Teacher | Soft delete teacher | Teacher tồn tại | Delete action | Teacher bị soft delete | Teacher | Role ADMIN | `DELETE /v1/teachers/:id` | Confirmed from code |
| UC-014 | Xem lịch dạy giáo viên | Xem lesson theo ngày | Admin/Teacher |  | Teacher | Trả về lessons của giáo viên | Teacher tồn tại | Query timetable | List lesson | Teacher, Lesson | from/to optional, format YYYY-MM-DD | `/v1/teachers/:id/timetable` | Confirmed from code |
| UC-015 | Xem thống kê giờ dạy | Tổng hợp giờ dạy | Admin/Teacher |  | Teacher | Tính tổng giờ và breakdown | Teacher tồn tại | Query stats | total_hours, breakdown | Lesson | group_by day/week/month | `/v1/teachers/:id/stats/teaching-hours` | Confirmed from code |
| UC-016 | Tạo học viên | Thêm hồ sơ học viên | Admin |  | Student | Tạo học viên mới | Auth required | Submit form student | Student mới | Student | Hiện không bắt ADMIN, chỉ cần auth | `/v1/students`, StudentsPage | Confirmed from code |
| UC-017 | Xem danh sách học viên | Tìm/lọc học viên | Admin / auth user |  | Student | Trả về danh sách học viên | Auth required | Open page | List students | Student | Search/status filter có nhưng query xây hơi bất thường | `GET /v1/students` | Confirmed from code |
| UC-018 | Cập nhật học viên | Sửa hồ sơ học viên | Admin / auth user |  | Student | Cập nhật student | Auth required | Edit student | Student updated | Student | Không có role check | `/v1/students/:id` | Confirmed from code |
| UC-019 | Xóa mềm học viên | Ngừng sử dụng hồ sơ học viên | Admin / auth user |  | Student | Soft delete student | Auth required | Delete student | Student deleted | Student | Không có role check | `DELETE /v1/students/:id` | Confirmed from code |
| UC-020 | Tạo khóa học | Tạo course mới | Admin |  | Course | Lưu môn học, số buổi, thời lượng, giá | Admin auth | Submit form | Course mới | Course | code+name required; status default ACTIVE | `/v1/courses`, CoursePage | Confirmed from code |
| UC-021 | Xem danh sách khóa học | Xem catalog khóa học | Guest/Admin |  | Course | Danh sách khóa học | None | Open page/API | Courses list | Course | Public GET; list use case bỏ qua search trong thực thi hiện tại | `/v1/courses` | Confirmed from code |
| UC-022 | Cập nhật khóa học | Sửa course | Admin |  | Course | Sửa thông tin course | Course tồn tại | Edit course | Course updated | Course | Role ADMIN | `/v1/courses/:id` | Confirmed from code |
| UC-023 | Xóa mềm khóa học | Ngừng sử dụng course | Admin |  | Course | Soft delete course | Course tồn tại | Delete course | Course deleted | Course | Role ADMIN | `DELETE /v1/courses/:id` | Confirmed from code |
| UC-024 | Tạo chương trình đào tạo | Tạo khung chương trình | Admin |  | Program | Tạo program theo track và hiệu lực | Admin auth | Submit program form | Program mới | Program | code+name required | `/v1/programs`, ProgramPage | Confirmed from code |
| UC-025 | Xem danh sách chương trình | Tra cứu chương trình | Guest/Admin |  | Program | Danh sách program | None | Open list | Programs list | Program | Public GET | `/v1/programs` | Confirmed from code |
| UC-026 | Cập nhật chương trình | Sửa program | Admin |  | Program | Cập nhật program | Program tồn tại | Edit program | Program updated | Program | Role ADMIN | `/v1/programs/:id` | Confirmed from code |
| UC-027 | Xóa mềm chương trình | Ngừng dùng program | Admin |  | Program | Soft delete program | Program tồn tại | Delete | Program deleted | Program | Role ADMIN | `DELETE /v1/programs/:id` | Confirmed from code |
| UC-028 | Gán khóa học vào chương trình | Liên kết course-program | Admin |  | Program | Thêm nhiều course vào program | Program tồn tại | Select courses | Mapping được tạo | ProgramCourse | Role ADMIN | `POST /v1/programs/:id/courses` | Confirmed from code |
| UC-029 | Gỡ khóa học khỏi chương trình | Bỏ liên kết course-program | Admin |  | Program | Xóa mapping | Program tồn tại | Remove courses | Mapping được xóa | ProgramCourse | Role ADMIN | `DELETE /v1/programs/:id/courses` | Confirmed from code |
| UC-030 | Tạo lớp học | Mở lớp vận hành | Admin |  | Class | Tạo lớp có thể gắn program/course/teacher | Admin auth | Submit class form | Class mới | Class | status default OPEN | `/v1/classes`, ClassesPage | Confirmed from code |
| UC-031 | Xem danh sách lớp | Tra cứu lớp | Guest/Admin |  | Class | Danh sách lớp với filter status | None | Open list | Classes list | Class | Public GET | `/v1/classes` | Confirmed from code |
| UC-032 | Xem chi tiết lớp / roster | Xem lớp và roster học viên | Admin |  | Class/Enrollment | Trả roster và capacity snapshot | Class tồn tại | Open class detail | Student list, capacity snapshot | Class, Enrollment, Room | capacity_limit = min(max_students, room.capacity) | `/v1/classes/:id/students`, ClassDetailDialog | Confirmed from code |
| UC-033 | Cập nhật lớp học | Sửa thông tin lớp | Admin |  | Class | Cập nhật class | Class tồn tại | Edit class | Class updated | Class | Role ADMIN | `/v1/classes/:id` | Confirmed from code |
| UC-034 | Xóa mềm lớp | Dừng sử dụng lớp | Admin |  | Class | Soft delete class | Class tồn tại | Delete class | Class deleted | Class | Role ADMIN | `DELETE /v1/classes/:id` | Confirmed from code |
| UC-035 | Ghi danh học viên vào lớp | Thêm học viên vào roster | Admin |  | Enrollment | Tạo enrollments hàng loạt | Class tồn tại | Add students | Enrollment created | Enrollment, Student, Class | Chưa đếm enrollment hiện có, chưa duplicate check | `POST /v1/classes/:id/students` | Confirmed from code |
| UC-036 | Gỡ học viên khỏi lớp | Loại học viên khỏi roster | Admin |  | Enrollment | Xóa enrollment theo danh sách | Enrollments tồn tại | Remove students | Enrollment deleted | Enrollment | Không có check trạng thái | `DELETE /v1/classes/:id/students` | Confirmed from code |
| UC-037 | Phân công giáo viên cho lớp | Chọn GV phụ trách chính | Admin |  | Class | Set `teacher_id` cho class | Class tồn tại | Assign teacher | Class updated | Class, Teacher | Role ADMIN | `PUT /v1/classes/:id/teacher` | Confirmed from code |
| UC-038 | Tạo phòng học | Thêm tài nguyên phòng | Admin |  | Room | Tạo room mới | Admin auth | Submit room form | Room mới | Room | DTO và use case đang lệch field | `/v1/rooms`, RoomsPage | Confirmed from code |
| UC-039 | Xem danh sách phòng | Tra cứu phòng | Guest/Admin |  | Room | Danh sách phòng | None | Open list | Rooms list | Room | Public GET | `/v1/rooms` | Confirmed from code |
| UC-040 | Cập nhật phòng | Sửa room | Admin |  | Room | Cập nhật room | Room tồn tại | Edit room | Room updated | Room | Role ADMIN | `/v1/rooms/:id` | Confirmed from code |
| UC-041 | Xóa mềm phòng | Ngừng dùng phòng | Admin |  | Room | Soft delete room | Room tồn tại | Delete | Room deleted | Room | Role ADMIN | `DELETE /v1/rooms/:id` | Confirmed from code |
| UC-042 | Tạo ca học | Định nghĩa shift | Admin |  | Shift | Tạo khung giờ chuẩn | Admin auth | Submit shift form | Shift mới | Shift | session_type enum, min duration 1 | `/v1/shifts`, ShiftsPage | Confirmed from code |
| UC-043 | Xem danh sách ca học | Tra cứu shift | Admin |  | Shift | Danh sách shift | Admin auth | Open list | shifts list | Shift | Chỉ ADMIN mới xem qua API | `/v1/shifts` | Confirmed from code |
| UC-044 | Cập nhật / xóa ca học | Quản lý lifecycle shift | Admin |  | Shift | Sửa hoặc xóa shift | Shift tồn tại | Edit/delete shift | Shift updated/deleted | Shift | Role ADMIN | `/v1/shifts/:id` | Confirmed from code |
| UC-045 | Tạo preview xếp lịch | Sinh phương án lịch | Admin | Scheduling engine | Scheduling | Chạy solver trên lớp OPEN trong khoảng ngày | date range hợp lệ | Run preview | Preview stored in memory, assignments/conflicts generated | Class, Shift, Room | date range valid; shifts active; classes OPEN | `/v1/scheduling/preview`, SchedulingPage | Confirmed from code |
| UC-046 | Xem preview xếp lịch | Kiểm tra kết quả preview | Admin | Scheduling engine | Scheduling | Lấy preview theo run_id hoặc latest | Preview tồn tại trong memory | Load preview | Preview shown | PreviewStore | Latest mất khi restart backend | `GET /preview/:id`, `/preview/latest` | Confirmed from code |
| UC-047 | Commit preview thành lesson | Sinh lesson thật | Admin | Scheduling engine | Scheduling/Lesson | Tạo lessons từ assignments | Preview COMPLETED và không conflict | Commit action | Lessons được tạo | Lesson | Chặn commit nếu status != COMPLETED hoặc còn conflict/unscheduled | `/v1/scheduling/commit` | Confirmed from code |
| UC-048 | Benchmark solver scheduling | So sánh solver | Admin | Scheduling engine | Scheduling | Lấy contract benchmark cho 3 solver | Admin auth | Call benchmark | Solver descriptors returned | Solver catalog | Chưa có execution metrics thật | `/v1/scheduling/benchmark` | Confirmed from code |
| UC-049 | Upload tài liệu giảng dạy | Đưa file vào hệ thống | Teacher | OCR stub, AI stub | Material | Upload file + metadata | Teacher auth | Submit upload | Material + audit log + label stub | Material, AuditLog | Chỉ hỗ trợ pdf/doc/docx/png/jpg/jpeg, max 10MB | `/v1/materials/upload`, TeacherDocumentsPage | Confirmed from code |
| UC-050 | Xem danh sách tài liệu | Tra cứu materials | Teacher/Compliance/Admin/auth user |  | Material | Xem list material theo teacher/status/queue | Auth required | Open page | Material list | Material | Ownership chưa enforce | `/v1/materials`, TeacherDocumentsPage, ComplianceQueuePage | Confirmed from code |
| UC-051 | Xem chi tiết tài liệu | Xem audit trail | Teacher/Compliance/Admin/auth user |  | Material | Mở material detail dialog | Auth required | Open detail | detail + audit logs + last decision | Material, AuditLog, ApprovalDecision | Ownership chưa enforce | `/v1/materials/:id` | Confirmed from code |
| UC-052 | Tải file tài liệu | Download file | Teacher/Compliance/Admin/auth user | File storage | Material | Tải file gốc qua endpoint nội bộ | Auth required | Click download | File attachment | Material | Quyền tải chưa chặt | `/v1/materials/:id/download` | Confirmed from code |
| UC-053 | Duyệt / từ chối tài liệu | Ra quyết định cuối cùng | Compliance / auth user |  | Material | Tạo ApprovalDecision và update status | Auth required | Submit review | Material APPROVED/REJECTED | Material, ApprovalDecision | Hiện không enforce compliance role; `compliance_officer_id` nhập từ request body | `/v1/materials/:id/review`, ComplianceQueuePage | Confirmed from code |
| UC-054 | Điểm danh buổi học | Ghi nhận chuyên cần | Teacher |  | Attendance | Có dấu vết data model nhưng chưa có flow | Lesson tồn tại | Chưa có | Attendance record | Attendance | Chưa có API và enum status | Missing API | Confirmed from code |
| UC-055 | Tổng kết buổi học | Lưu nội dung đã dạy | Teacher |  | Lesson Summary | Có data model nhưng chưa có flow | Lesson tồn tại | Chưa có | Lesson summary | LessonSummary | Chưa có API | Missing API | Confirmed from code |
| UC-056 | Ghi nhận kết quả học tập | Lưu kết quả chi tiết học viên | Teacher |  | Academic | Có data model nhưng chưa có flow | Lesson summary tồn tại | Chưa có | AcademicRecord | AcademicRecord | Chưa có API | Missing API | Confirmed from code |
| UC-057 | Xin phép nghỉ / đi muộn / về sớm | Tạo đơn nghỉ phép | Student | Admin | Leave | Có data model nhưng chưa có flow | Student/lesson/class liên quan | Chưa có | LeaveRequest | LeaveRequest | Chưa có API | Missing API | Confirmed from code |
| UC-058 | Tiếp nhận lead tư vấn | Lưu đầu mối tư vấn | Staff |  | Consultation | Có data model nhưng chưa có flow | Khách hàng tiềm năng liên hệ | Chưa có | Consultation | Consultation | Chưa có API | Missing API | Confirmed from code |
| UC-059 | Dự báo học viên nguy cơ học kém | Phân loại AT_RISK | Admin/Manager | ML engine | Predictive | Chưa có code, mới có kế hoạch | Có dữ liệu học tập đầy đủ | Chưa có | Prediction output | Attendance, AcademicRecord, LessonSummary, Student | Chưa implemented | Planned only | Confirmed from code |

## 5.2 Nhóm use case theo tác nhân

- **Guest/User chưa đăng nhập**: UC-001, UC-002, UC-003, UC-006, UC-007, UC-021, UC-025, UC-031, UC-039.
- **Admin**: UC-010 đến UC-048, UC-050 đến UC-053.
- **Teacher**: UC-014, UC-015, UC-027 đến UC-029, UC-049 đến UC-052, và dự kiến UC-054 đến UC-056.
- **Student**: UC-003 đến UC-009, dự kiến UC-057 và các use case timetable/results sau này.
- **Compliance**: UC-050, UC-051, UC-053.

## 5.3 Nhóm use case theo miền

- **Auth**: UC-001 -> UC-009
- **Teacher**: UC-010 -> UC-015
- **Student**: UC-016 -> UC-019
- **Course/Program**: UC-020 -> UC-029
- **Class/Enrollment**: UC-030 -> UC-037
- **Room/Shift**: UC-038 -> UC-044
- **Scheduling/Lesson**: UC-045 -> UC-048
- **Materials/Moderation**: UC-049 -> UC-053
- **Academic Operations (planned/missing)**: UC-054 -> UC-059

## 5.4 Gợi ý chia cụm sơ đồ use case

1. **Cụm Tài khoản và Truy cập**
2. **Cụm Danh mục đào tạo**
3. **Cụm Vận hành lớp học**
4. **Cụm Xếp lịch và Lesson**
5. **Cụm Tài liệu giảng dạy và Kiểm duyệt**
6. **Cụm Theo dõi học tập**
7. **Cụm Phân tích và Dự báo**

## 5.5 Danh sách 25 use case ưu tiên cao nhất

UC-001, UC-002, UC-003, UC-006, UC-007, UC-008, UC-010, UC-011, UC-016, UC-017, UC-020, UC-024, UC-028, UC-030, UC-032, UC-035, UC-037, UC-038, UC-042, UC-045, UC-046, UC-047, UC-049, UC-051, UC-053.

---

# 6. Luồng nghiệp vụ đầu-cuối

## WF-01. Đăng ký và xác minh email

- **Mã workflow:** WF-01
- **Mục tiêu nghiệp vụ:** Tạo tài khoản mới và kích hoạt hợp lệ qua email OTP
- **Điểm bắt đầu:** Người dùng điền form đăng ký
- **Điểm kết thúc:** Tài khoản được active và có thể đăng nhập
- **Swimlane:** Guest -> Auth API -> User/OTP Repository -> SMTP
- **Các bước chi tiết:**
  1. Guest gửi `email`, `full_name`, `password_enc`.
  2. Hệ thống kiểm tra email đã tồn tại hay chưa.
  3. Hệ thống giải mã password FE gửi lên, băm mật khẩu.
  4. Hệ thống tạo `User` với `is_active = false`.
  5. Hệ thống sinh OTP 6 số, băm OTP, tạo bản ghi `UserOTP` với hạn 5 phút.
  6. Hệ thống gửi email OTP bất đồng bộ qua mailer.
  7. Người dùng nhập OTP.
  8. Hệ thống kiểm tra OTP active, chưa hết hạn, so khớp hash.
  9. Hệ thống mark OTP used và activate user.
- **Điểm quyết định:** Email đã tồn tại? OTP còn hạn? OTP đúng?
- **Nhánh thay thế:** SMTP không cấu hình trong môi trường dev thì log mail ra console.
- **Nhánh lỗi / ngoại lệ:** Email trùng, OTP sai, OTP hết hạn.
- **Dữ liệu tạo/cập nhật:** User, UserOTP.
- **Quy tắc nghiệp vụ liên quan:** BR-DATA-001, BR-LC-001, BR-INT-001.
- **Hệ thống / phân hệ tham gia:** Auth, Mailer.
- **Điểm chưa rõ cần BA xác nhận:** Có cần resend OTP và giới hạn số lần gửi không?
- **Mức độ chắc chắn:** Confirmed from code

## WF-02. Đăng nhập / refresh / logout

- **Mã workflow:** WF-02
- **Mục tiêu nghiệp vụ:** Thiết lập và duy trì truy cập hệ thống
- **Điểm bắt đầu:** User đăng nhập
- **Điểm kết thúc:** Có access token hợp lệ hoặc người dùng đăng xuất
- **Swimlane:** User -> Auth API -> AuthService -> JWT
- **Các bước chi tiết:**
  1. User gửi `username` và `password`.
  2. Hệ thống tìm user theo email.
  3. Hệ thống chặn login nếu `is_active = false`.
  4. Hệ thống so khớp mật khẩu băm.
  5. Hệ thống sinh access token và refresh token.
  6. Khi refresh, hệ thống dùng refresh token để sinh access token mới.
  7. Khi logout, hệ thống chỉ validate token và trả success.
- **Điểm quyết định:** User tồn tại? active? password đúng? refresh token hợp lệ?
- **Nhánh thay thế:** Không có session store; logout chỉ là logic client-side + validate token.
- **Nhánh lỗi / ngoại lệ:** Invalid credentials, invalid token.
- **Dữ liệu tạo/cập nhật:** Không có session persistence.
- **Quy tắc nghiệp vụ liên quan:** BR-AUTH-002, BR-AUTH-003.
- **Điểm chưa rõ cần BA xác nhận:** Có cần blacklist refresh token và multi-device session không?
- **Mức độ chắc chắn:** Confirmed from code

## WF-03. Quên mật khẩu / đặt lại mật khẩu

- **Mã workflow:** WF-03
- **Mục tiêu nghiệp vụ:** Khôi phục truy cập an toàn
- **Điểm bắt đầu:** User quên mật khẩu
- **Điểm kết thúc:** Mật khẩu mới được lưu
- **Swimlane:** User -> Auth API -> PasswordReset repo -> SMTP
- **Các bước chi tiết:**
  1. User nhập email.
  2. Hệ thống normalize email và tra user.
  3. Nếu không thấy user, vẫn trả success để tránh enumeration.
  4. Nếu có user, hệ thống sinh reset token, băm token, lưu `PasswordReset` kèm IP/UA và hạn dùng.
  5. Hệ thống gửi email chứa reset link.
  6. User mở link, nhập mật khẩu mới.
  7. Hệ thống băm token đầu vào để tra `PasswordReset`.
  8. Hệ thống mark token used trong transaction.
  9. Hệ thống cập nhật password mới.
- **Nhánh lỗi / ngoại lệ:** Token sai/hết hạn; mật khẩu mới quá ngắn.
- **Điểm chưa rõ cần BA xác nhận:** Có cần revoke toàn bộ session sau reset password không?
- **Mức độ chắc chắn:** Confirmed from code

## WF-04. Tạo học viên

- **Mã workflow:** WF-04
- **Mục tiêu nghiệp vụ:** Tạo hồ sơ học viên phục vụ ghi danh
- **Điểm bắt đầu:** Người vận hành nhập form học viên
- **Điểm kết thúc:** Student record được tạo
- **Swimlane:** Admin/Auth user -> Student API -> Student repo
- **Các bước chi tiết:**
  1. Người dùng auth mở màn hình học viên.
  2. Nhập mã, họ tên, email, điện thoại, khối lớp, số phụ huynh, trạng thái...
  3. Hệ thống bind request.
  4. Use case tạo `Student`.
  5. Repository lưu vào DB.
- **Nhánh lỗi / ngoại lệ:** Lỗi DB/unique nếu schema có ràng buộc.
- **Điểm chưa rõ cần BA xác nhận:** Có cần role ADMIN hoặc phân quyền tinh hơn không?
- **Mức độ chắc chắn:** Confirmed from code

## WF-05. Tạo giáo viên

- **Mã workflow:** WF-05
- **Mục tiêu nghiệp vụ:** Tạo hồ sơ giảng viên
- **Điểm bắt đầu:** Admin mở form giáo viên
- **Điểm kết thúc:** Teacher được tạo
- **Các bước chi tiết:**
  1. Admin nhập thông tin giáo viên.
  2. Hệ thống kiểm tra `full_name`.
  3. Nếu có email thì check uniqueness.
  4. Nếu có code thì check uniqueness.
  5. Hệ thống gán mặc định `employment_type = PART_TIME`, `status = ACTIVE` nếu trống.
  6. Lưu teacher vào DB.
- **Mức độ chắc chắn:** Confirmed from code

## WF-06. Tạo khóa học

- **Mã workflow:** WF-06
- **Mục tiêu nghiệp vụ:** Tạo đơn vị đào tạo cơ sở
- **Điểm bắt đầu:** Admin nhập khóa học
- **Điểm kết thúc:** Course được lưu
- **Các bước chi tiết:** nhập code, tên, mô tả, môn, khối, số buổi, thời lượng, tổng giờ, giá; hệ thống gán `status=ACTIVE` nếu trống rồi lưu.
- **Điểm chưa rõ:** Không thấy validate chặt `session_count`, `duration`, `total_hours`, quan hệ logic giữa chúng.
- **Mức độ chắc chắn:** Confirmed from code

## WF-07. Tạo chương trình và gán khóa học

- **Mã workflow:** WF-07
- **Mục tiêu nghiệp vụ:** Xây chương trình đào tạo từ nhiều khóa học
- **Điểm bắt đầu:** Admin tạo program
- **Điểm kết thúc:** Program có danh sách course liên kết
- **Các bước chi tiết:**
  1. Tạo program với code, name, track, khoảng hiệu lực, approval_note.
  2. Mở chi tiết program.
  3. Gửi danh sách `course_ids` để add vào program.
  4. Có thể gửi `course_ids` để remove.
- **Điểm chưa rõ:** Objective/Outcome có trong DB nhưng chưa có luồng quản trị.
- **Mức độ chắc chắn:** Confirmed from code

## WF-08. Mở lớp

- **Mã workflow:** WF-08
- **Mục tiêu nghiệp vụ:** Tạo lớp thực thi từ khóa học/chương trình
- **Điểm bắt đầu:** Admin tạo class
- **Điểm kết thúc:** Class ở trạng thái OPEN
- **Các bước chi tiết:**
  1. Nhập mã lớp, tên lớp, ngày bắt đầu/kết thúc, max_students, price, notes.
  2. Chọn `program_id`, `course_id`, `teacher_id` nếu có.
  3. Hệ thống gán `status = OPEN` nếu để trống.
  4. Lưu class vào DB.
- **Điểm chưa rõ:** Tại bước mở lớp chưa thấy API quản lý `class_schedule`; dữ liệu này có thể được seed/chỉnh trực tiếp ở nơi khác.
- **Mức độ chắc chắn:** Confirmed from code

## WF-09. Ghi danh học viên

- **Mã workflow:** WF-09
- **Mục tiêu nghiệp vụ:** Đưa học viên vào roster lớp
- **Điểm bắt đầu:** Admin chọn lớp và danh sách học viên
- **Điểm kết thúc:** Enrollment records được tạo
- **Các bước chi tiết:**
  1. Admin mở class detail.
  2. Chọn nhiều `student_id`.
  3. Use case lấy class.
  4. Tính `maxAllowed = min(class.max_students, room.capacity nếu có)`.
  5. Hiện tại use case giả định chưa có enrollment hiện hữu và kiểm tra trên danh sách đầu vào.
  6. Tạo enrollment cho từng student với status `ENROLLED`.
- **Nhánh lỗi / ngoại lệ:** Nếu số student đầu vào vượt `maxAllowed` thì lỗi.
- **Điểm chưa rõ / gap:** Không đếm enrollment đang có; không duplicate check; không lifecycle APPLIED -> APPROVED -> ENROLLED.
- **Mức độ chắc chắn:** Confirmed from code

## WF-10. Phân công giáo viên

- **Mã workflow:** WF-10
- **Mục tiêu nghiệp vụ:** Gán giáo viên phụ trách cho lớp
- **Điểm bắt đầu:** Admin chọn teacher cho class
- **Điểm kết thúc:** `class.teacher_id` được cập nhật
- **Mức độ chắc chắn:** Confirmed from code

## WF-11. Cấu hình lịch tuần

- **Mã workflow:** WF-11
- **Mục tiêu nghiệp vụ:** Gắn lịch chuẩn để solver biết ngày trong tuần và shift nào lớp có thể học
- **Điểm bắt đầu:** Class đã tồn tại
- **Điểm kết thúc:** `class_schedules` có `day_of_week`, `shift_id`, optional `room_id`
- **Các bước chi tiết:** Chưa thấy controller/use case quản lý `class_schedules`; chỉ thấy scheduling preload bảng này.
- **Điểm chưa rõ cần BA xác nhận:** Lớp cấu hình lịch tuần ở màn nào? quản lý bởi admin hay seed dữ liệu?
- **Mức độ chắc chắn:** Confirmed from code về data model, Assumption / Needs BA validation về workflow UI

## WF-12. Tạo preview xếp lịch

- **Mã workflow:** WF-12
- **Mục tiêu nghiệp vụ:** Sinh phương án lịch trước khi commit
- **Điểm bắt đầu:** Admin chạy scheduling preview
- **Điểm kết thúc:** PreviewResult được lưu vào PreviewStore và hiển thị trên UI
- **Swimlane:** Admin -> Scheduling API -> load classes/rooms/shifts -> solver -> PreviewStore
- **Các bước chi tiết:**
  1. Validate `date_to >= date_from`.
  2. Load classes `status = OPEN`, preload teacher/course/room/class_schedules/shift.
  3. Load rooms theo filter.
  4. Load shifts `is_active = true`.
  5. Nếu thiếu class/room/shift thì đẩy conflicts tổng quan.
  6. Gọi solver service qua abstraction `SchedulingSolver`.
  7. Nhận `assignments`, `conflicts`, `summary`.
  8. Lưu preview vào `PreviewStore` in-memory.
- **Nhánh lỗi / ngoại lệ:** date range invalid, DB lỗi, preview not found ở latest nếu chưa có preview thành công.
- **Mức độ chắc chắn:** Confirmed from code

## WF-13. Xác nhận preview để sinh lesson

- **Mã workflow:** WF-13
- **Mục tiêu nghiệp vụ:** Ghi lịch preview thành lesson thật
- **Điểm bắt đầu:** Admin bấm commit
- **Điểm kết thúc:** lessons được tạo trong DB
- **Các bước chi tiết:**
  1. Nhận `run_id`.
  2. Lấy preview từ store.
  3. Chặn nếu không có assignment.
  4. Chặn nếu preview status khác `COMPLETED`.
  5. Tính cửa sổ thời gian preview và load lesson overlap hiện có.
  6. So conflict với class/teacher/room hiện hữu.
  7. Nếu không conflict thì tạo `Lesson` cho từng assignment trong transaction.
- **Điểm chưa rõ:** Chưa có idempotency hoặc cờ preview đã commit.
- **Mức độ chắc chắn:** Confirmed from code

## WF-14. Vận hành buổi học

- **Mã workflow:** WF-14
- **Mục tiêu nghiệp vụ:** Diễn ra buổi học sau khi lesson được tạo
- **Hiện trạng:** lesson tồn tại nhưng chưa có màn/API vận hành buổi học thật.
- **Mức độ chắc chắn:** Confirmed from code về lesson; Missing API cho flow vận hành

## WF-15. Điểm danh

- **Mã workflow:** WF-15
- **Mục tiêu nghiệp vụ:** Ghi chuyên cần theo lesson/student
- **Hiện trạng:** chỉ có entity `Attendance`; không có API/controller/use case.
- **Mức độ chắc chắn:** Confirmed from code

## WF-16. Tổng kết buổi học

- **Mã workflow:** WF-16
- **Mục tiêu nghiệp vụ:** Lưu topic, homework, feedback sau buổi học
- **Hiện trạng:** chỉ có entity `LessonSummary`; không có API/controller/use case.
- **Mức độ chắc chắn:** Confirmed from code

## WF-17. Ghi nhận kết quả học tập

- **Mã workflow:** WF-17
- **Mục tiêu nghiệp vụ:** Ghi nhận điểm/tổng kết từng học viên sau lesson summary
- **Hiện trạng:** chỉ có entity `AcademicRecord`; không có API/controller/use case.
- **Mức độ chắc chắn:** Confirmed from code

## WF-18. Xin phép nghỉ / đi muộn / về sớm

- **Mã workflow:** WF-18
- **Mục tiêu nghiệp vụ:** Ghi nhận đơn nghỉ học hoặc đi muộn/về sớm
- **Hiện trạng:** chỉ có entity `LeaveRequest`; không có API/controller/use case.
- **Mức độ chắc chắn:** Confirmed from code

## WF-19. Upload và duyệt tài liệu

- **Mã workflow:** WF-19
- **Mục tiêu nghiệp vụ:** Kiểm soát tài liệu giảng dạy
- **Swimlane:** Teacher -> Material API -> Storage + OCR stub + AI stub -> Compliance/Admin review
- **Các bước chi tiết:**
  1. Teacher upload file và metadata.
  2. Backend validate file type/size.
  3. Backend lưu file local theo path có cấu trúc.
  4. Tạo material status `SCANNING`.
  5. Chạy OCR stub và AI stub.
  6. Tạo audit log và update material `AI_REVIEWED`.
  7. Compliance mở queue flagged.
  8. Gửi quyết định approve/reject.
  9. Hệ thống tạo ApprovalDecision và update material status.
- **Điểm chưa rõ:** Quyền duyệt hiện không enforce role; `teacher_id` và `compliance_officer_id` chưa derive từ auth.
- **Mức độ chắc chắn:** Confirmed from code

## WF-20. Tiếp nhận tư vấn tuyển sinh

- **Mã workflow:** WF-20
- **Hiện trạng:** Có data model `consultations`, chưa có API/controller/UI.
- **Mức độ chắc chắn:** Confirmed from code

---

# 7. Danh mục quy tắc nghiệp vụ

## 7.1 Quy tắc dữ liệu

| Mã rule | Tên rule | Phát biểu rule bằng ngôn ngữ nghiệp vụ | Loại rule | Ngữ cảnh áp dụng | Use case bị ảnh hưởng | Thực thể bị ảnh hưởng | Điều gì xảy ra nếu vi phạm | Bằng chứng kỹ thuật | Mức độ chắc chắn | Ghi chú cho BA |
|---|---|---|---|---|---|---|---|---|---|---|
| BR-DATA-001 | Email người dùng là duy nhất | Không được tạo 2 tài khoản cùng email | validation | Register | UC-001 | User | Đăng ký lỗi | `register.go`, unique on user email | Confirmed from code |  |
| BR-DATA-002 | OTP phải còn hiệu lực | OTP chỉ dùng được trong khoảng thời gian cho phép | validation/lifecycle | Verify OTP | UC-002 | UserOTP | Kích hoạt thất bại | `verify_otp.go` | Confirmed from code |  |
| BR-DATA-003 | Email giáo viên là duy nhất nếu có nhập | Không được trùng email giáo viên | validation | Create teacher | UC-010 | Teacher | Tạo lỗi | `create_teacher.go` | Confirmed from code | Không thấy check khi update |
| BR-DATA-004 | Code giáo viên là duy nhất nếu có nhập | Không được trùng mã giáo viên | validation | Create teacher | UC-010 | Teacher | Tạo lỗi | `create_teacher.go` | Confirmed from code | Không thấy check khi update |
| BR-DATA-005 | Khóa học phải có mã và tên | Course không được tạo nếu thiếu code hoặc name | validation | Create course | UC-020 | Course | API lỗi | `create_course.go` | Confirmed from code |  |
| BR-DATA-006 | Chương trình phải có mã và tên | Program không được tạo nếu thiếu code hoặc name | validation | Create program | UC-024 | Program | API lỗi | `create_program.go` | Confirmed from code |  |
| BR-DATA-007 | Lớp phải có sĩ số tối đa > 0 | `max_students` tối thiểu là 1 | validation | Create/update class | UC-030, UC-033 | Class | Bind request lỗi | `class/dto.go` | Confirmed from code |  |
| BR-DATA-008 | Shift phải có thời lượng > 0 | Ca học phải có `duration_minutes >= 1` | validation | Create/update shift | UC-042, UC-044 | Shift | Bind request lỗi | `shift/dto.go` | Confirmed from code |  |
| BR-DATA-009 | Review material phải có quyết định | Một quyết định duyệt cần approved hoặc reject cùng thông tin người duyệt | validation | Review material | UC-053 | ApprovalDecision | API lỗi | `material/dto.go`, `review_material.go` | Confirmed from code | Nhưng người duyệt nhập tay ID |
| BR-DATA-010 | Upload material bị giới hạn loại file và dung lượng | Chỉ chấp nhận pdf/doc/docx/png/jpg/jpeg và <= 10 MB | validation | Upload material | UC-049 | Material | Upload lỗi | `upload_material.go`, `TeacherDocumentsPage.tsx` | Confirmed from code |  |

## 7.2 Quy tắc mặc định

| Mã rule | Tên rule | Phát biểu | Loại rule | Bằng chứng kỹ thuật | Mức độ chắc chắn |
|---|---|---|---|---|---|
| BR-DEF-001 | Role user mặc định là STUDENT | Nếu không chỉ định role, user mặc định là STUDENT | default value | `entities/user.go` | Confirmed from code |
| BR-DEF-002 | User entity mặc định active nhưng register override inactive | Entity default là true nhưng luồng register tạo user inactive | default value / contradiction | `entities/user.go` và `register.go` | Confirmed from code |
| BR-DEF-003 | Teacher employment type mặc định PART_TIME | Nếu không nhập loại tuyển dụng, mặc định bán thời gian | default value | `create_teacher.go` | Confirmed from code |
| BR-DEF-004 | Teacher status mặc định ACTIVE | Nếu không nhập trạng thái, mặc định ACTIVE | default value | `create_teacher.go` | Confirmed from code |
| BR-DEF-005 | Course status mặc định ACTIVE | Nếu không nhập trạng thái, course mặc định ACTIVE | default value | `create_course.go` | Confirmed from code |
| BR-DEF-006 | Class status mặc định OPEN | Lớp mới mặc định OPEN | default value | `create_class.go`, entity `Class` | Confirmed from code |
| BR-DEF-007 | Teaching hours group_by mặc định day | Nếu không truyền group_by, hệ thống nhóm theo ngày | default value | `get_teaching_hours.go` | Confirmed from code |

## 7.3 Quy tắc vòng đời

| Mã rule | Tên rule | Phát biểu | Loại rule | Ngữ cảnh | Vi phạm thì sao | Bằng chứng kỹ thuật | Mức độ chắc chắn |
|---|---|---|---|---|---|---|---|
| BR-LC-001 | User phải được kích hoạt mới đăng nhập | Tài khoản chưa verify OTP không được login | lifecycle | Login | Lỗi `user is not active` | `auth.go` | Confirmed from code |
| BR-LC-002 | Password reset token chỉ dùng một lần | Token reset sau khi dùng phải mark used | lifecycle | Reset password | Reset thất bại nếu token đã dùng | `forgot_reset_change.go` | Confirmed from code |
| BR-LC-003 | Preview scheduling chỉ commit khi hoàn tất | Chỉ preview `COMPLETED` mới commit được | lifecycle/scheduling | Commit preview | Commit bị chặn | `commit_preview.go` | Confirmed from code |
| BR-LC-004 | Material đi qua trạng thái scanning -> reviewed -> approved/rejected | Tài liệu sau upload được scan và sau đó có thể được duyệt thủ công | lifecycle | Material flow | Không có state machine formal nhưng có status update thực tế | `upload_material.go`, `review_material.go` | Confirmed from code |
| BR-LC-005 | Enrollment hiện bị đẩy thẳng sang ENROLLED | Khi thêm vào lớp, enrollment không qua bước approval | lifecycle | Enroll students | Không có bước duyệt | `enroll_students.go` | Confirmed from code |

## 7.4 Quy tắc phân quyền

| Mã rule | Tên rule | Phát biểu | Loại | Ngữ cảnh | Bằng chứng kỹ thuật | Mức độ chắc chắn | Ghi chú BA |
|---|---|---|---|---|---|---|---|
| BR-AUTH-001 | Teachers/Courses/Programs/Rooms/Classes read APIs đang public | Nhiều API đọc không yêu cầu auth | authorization | List/detail | route controllers | Confirmed from code | Cần xác nhận có chủ đích public không |
| BR-AUTH-002 | Teacher CRUD chỉ dành cho ADMIN | Tạo/sửa/xóa giáo viên có role ADMIN | authorization | Teacher CRUD | `teacher/controller.go` | Confirmed from code |  |
| BR-AUTH-003 | Course/Program/Room/Class/Shift CRUD chủ yếu dành cho ADMIN | Tạo/sửa/xóa các master data và scheduling là ADMIN | authorization | CRUD admin | route controllers | Confirmed from code |  |
| BR-AUTH-004 | Student CRUD hiện chỉ yêu cầu auth, không yêu cầu ADMIN | Bất kỳ user đăng nhập đều thao tác student | authorization gap | Student CRUD | `student/controller.go` | Confirmed from code | Rủi ro nghiệp vụ cao |
| BR-AUTH-005 | Review material hiện chỉ yêu cầu auth, không yêu cầu compliance role | Bất kỳ user đăng nhập đều có thể gọi review | authorization gap | Material review | `material/controller.go` | Confirmed from code | Rủi ro cao |
| BR-AUTH-006 | Upload material yêu cầu TEACHER | Teacher upload tài liệu qua role TEACHER | authorization | Upload material | `material/controller.go` | Confirmed from code |  |
| BR-AUTH-007 | `/v1/auth/me` không được bảo vệ bằng auth middleware ở route | Route không gắn middleware nhưng controller lại kỳ vọng `user_id` trong context | authorization / technical defect | GetMe | `user/controller.go`, `user/v1.go` | Confirmed from code | Vừa là auth gap vừa là defect |

## 7.5 Quy tắc xếp lịch

| Mã rule | Tên rule | Phát biểu rule | Loại | Ngữ cảnh | Vi phạm thì sao | Bằng chứng kỹ thuật | Mức độ chắc chắn |
|---|---|---|---|---|---|---|---|
| BR-SCH-001 | Khoảng ngày preview phải hợp lệ | `date_to >= date_from` | scheduling/validation | Preview, benchmark | Preview FAILED hoặc 400 | `scheduling/v1.go`, `preview.go` | Confirmed from code |
| BR-SCH-002 | Chỉ lớp OPEN mới được xét xếp lịch | Solver chỉ load class `status = OPEN` | scheduling | Preview | Lớp khác trạng thái bị bỏ qua | `preview.go` | Confirmed from code |
| BR-SCH-003 | Phải có ít nhất một Shift active để sinh slot | Nếu không có ca học active thì không có slot xếp lịch | scheduling | Preview | Conflict `NO_ACTIVE_SHIFT` | `preview.go` | Confirmed from code |
| BR-SCH-004 | Solver dùng class_schedule theo shift_id | Lịch tuần lớp xác định ngày trong tuần và ca học được phép | scheduling | Preview | Conflict no-slot/no-domain | `preview.go`, solver services | Confirmed from code |
| BR-SCH-005 | Commit phải tránh trùng lesson hiện hữu | Không được commit nếu trùng lớp, giáo viên hoặc phòng với lesson đang tồn tại | scheduling/lifecycle | Commit | Commit fail | `commit_preview.go`, lesson repo overlap | Confirmed from code |
| BR-SCH-006 | Preview latest chỉ tồn tại trong bộ nhớ tiến trình | Restart backend sẽ mất preview latest | technical/scheduling | Get latest preview | 404 preview not found | `preview_store.go` | Confirmed from code |
| BR-SCH-007 | Benchmark 3 solver đã có catalog nhưng chưa có metric benchmark thật | Hệ thống mới chuẩn bị contract benchmark | scheduling | Benchmark | Chưa có số liệu so sánh thực | `benchmark.go` | Confirmed from code |

## 7.6 Quy tắc tích hợp

| Mã rule | Tên rule | Phát biểu | Loại | Bằng chứng kỹ thuật | Mức độ chắc chắn |
|---|---|---|---|---|---|
| BR-INT-001 | Môi trường dev có thể không cấu hình SMTP | Nếu thiếu SMTP host/port và app.env là dev/local thì mail được log ra thay vì fail | integration | `mailer.go` | Confirmed from code |
| BR-INT-002 | OCR và Gemini hiện là stub | Audit nội dung tài liệu hiện không gọi dịch vụ AI production | integration | `internal/services/audit/services.go` | Confirmed from code |
| BR-INT-003 | Queue infrastructure có nhưng chưa là đường đi chính của mail/material flow | Mail đang gửi trực tiếp qua SMTP trong `Send`, queue chỉ là dependency nằm sẵn | integration | `mailer.go`, queue packages | Confirmed from code |

## 7.7 Quy tắc còn thiếu hoặc mâu thuẫn

| Mã rule | Mô tả mâu thuẫn / thiếu | Ảnh hưởng | Bằng chứng kỹ thuật | Mức độ chắc chắn |
|---|---|---|---|---|
| BR-GAP-001 | Program migration có `status`, entity/controller hiện không dùng `status` | Lifecycle chương trình mơ hồ | SQL migration 7 vs entity `Program` | Confirmed from code |
| BR-GAP-002 | Room DTO nói `status/location`, entity/use case dùng `code/address` | Hợp đồng API room không nhất quán với model | `room/dto.go`, `entities/room.go`, room use cases | Confirmed from code |
| BR-GAP-003 | `/auth/me` dùng UUID nhưng ép sang int | Không thể lấy profile đúng trong runtime | `user/v1.go`, `get_user_by_id.go` | Confirmed from code |
| BR-GAP-004 | Enrollment lifecycle không rõ | Có entity default `APPLIED` nhưng use case thêm vào lớp lại tạo `ENROLLED` trực tiếp | `entities/enrollment.go`, `enroll_students.go` | Confirmed from code |
| BR-GAP-005 | Attendance status là số nguyên nhưng không có bảng mã | Không thể viết rule chuyên cần đầy đủ | `entities/attendance.go` | Confirmed from code |

---

# 8. Danh mục thực thể nghiệp vụ

| Mã thực thể | Tên thực thể | Định nghĩa nghiệp vụ | Tên bảng / model | Loại thực thể | Vai trò trong vận hành trung tâm | Nguồn phát sinh dữ liệu | Use case liên quan | Trường trạng thái / lifecycle | Soft delete / audit | Ghi chú / điểm mơ hồ |
|---|---|---|---|---|---|---|---|---|---|---|
| ENT-001 | Người dùng | Tài khoản đăng nhập hệ thống | `users` / `User` | Master | Cấp truy cập cho admin/teacher/student | Register, seed, có thể tạo nội bộ | UC-001..UC-009 | `role`, `is_active` | Soft delete | `role` default STUDENT |
| ENT-002 | OTP người dùng | OTP xác minh email | `user_otps` / `UserOTP` | Transaction/security | Kích hoạt tài khoản và xác minh hành động | Register | UC-002 | `expired_at`, `used_at` | Soft delete |  |
| ENT-003 | Yêu cầu đặt lại mật khẩu | Token reset mật khẩu | `password_resets` / `PasswordReset` | Transaction/security | Hỗ trợ quên mật khẩu | Forgot password | UC-006, UC-007 | `expires_at`, `used_at` | Không soft delete rõ | Lưu IP/UA |
| ENT-004 | Giáo viên | Hồ sơ giảng dạy | `teachers` / `Teacher` | Master | Nguồn lực giảng dạy | CRUD teacher | UC-010..UC-015 | `status`, `employment_type` | Soft delete | Có cờ giáo viên trường |
| ENT-005 | Học viên | Hồ sơ người học | `students` / `Student` | Master | Đầu vào cho ghi danh và theo dõi học tập | CRUD student | UC-016..UC-019 | `status` | Soft delete | Có `guardian_phone` |
| ENT-006 | Khóa học | Đơn vị học tập cơ sở | `courses` / `Course` | Master | Định nghĩa số buổi, môn, thời lượng, học phí | CRUD course | UC-020..UC-023 | `status` | Soft delete | Là input quan trọng cho scheduling |
| ENT-007 | Chương trình đào tạo | Tập hợp khóa học theo track | `programs` / `Program` | Master | Khung chương trình đào tạo | CRUD program | UC-024..UC-029 | Không có status ở entity, có `published_at`, `archived_at` | Soft delete | Mâu thuẫn migration/entity |
| ENT-008 | Liên kết chương trình-khóa học | Mapping N-N giữa program và course | `program_courses` / `ProgramCourse` | Mapping | Gắn course vào program | Add/remove courses | UC-028, UC-029 | Không có status | Không soft delete rõ |  |
| ENT-009 | Mục tiêu đào tạo | Objective của chương trình | `objectives` / `Objective` | Master | Cấu trúc mục tiêu chương trình | Chưa có use case | Planned | Không rõ | Không rõ | Chưa có API |
| ENT-010 | Chuẩn đầu ra | Outcome của chương trình | `outcomes` / `Outcome` | Master | Cấu trúc outcome chương trình | Chưa có use case | Planned | Không rõ | Không rõ | Chưa có API |
| ENT-011 | Lớp học | Đơn vị vận hành lớp | `classes` / `Class` | Master/operational | Nơi gắn học viên, giáo viên, khóa học, lịch tuần | CRUD class | UC-030..UC-037 | `status` OPEN/CLOSED/CANCELLED | Soft delete | Có optional room/program/course/teacher |
| ENT-012 | Lịch tuần lớp | Mẫu lịch cho lớp | `class_schedules` / `ClassSchedule` | Configuration/operational | Ràng buộc ngày trong tuần và shift cho scheduling | Back-office config | UC-045 | `day_of_week`, `shift_id` | Không thấy soft delete | Chưa có API quản trị |
| ENT-013 | Ghi danh | Quan hệ học viên vào lớp | `enrollments` / `Enrollment` | Transaction/mapping | Quản lý roster lớp | Enroll/remove students | UC-035, UC-036 | `status`, `approved_at`, `rejected_at` | Không soft delete rõ ở entity nhưng xóa qua repo | Use case không dùng hết lifecycle |
| ENT-014 | Phòng học | Tài nguyên phòng | `rooms` / `Room` | Master | Ràng buộc sức chứa và lịch học | CRUD room | UC-038..UC-041 | Không có status ở entity | Không soft delete | DTO lệch entity |
| ENT-015 | Ca học | Khung giờ chuẩn | `shifts` / `Shift` | Configuration/master | Sinh time slot cho scheduling | CRUD shift | UC-042..UC-044 | `session_type`, `is_active` | Không soft delete |  |
| ENT-016 | Buổi học | Phiên học thực tế | `lessons` / `Lesson` | Transaction/operational | Lịch thực thi sau commit | Scheduling commit | UC-014, UC-047 | Không có status riêng | Không soft delete | Cần cho attendance/summary |
| ENT-017 | Điểm danh | Chuyên cần theo lesson/student | `attendances` / `Attendance` | Transaction | Theo dõi hiện diện | Chưa có API | UC-054 | `status` int | Không rõ | Thiếu bảng mã status |
| ENT-018 | Tổng kết buổi học | Nội dung dạy sau lesson | `lesson_summaries` / `LessonSummary` | Transaction | Ghi topic, homework, feedback | Chưa có API | UC-055 | Không rõ | Không rõ | Quan hệ 1-1 lesson |
| ENT-019 | Kết quả học tập | Đánh giá chi tiết học viên sau buổi học | `academic_records` / `AcademicRecord` | Transaction | Nền dữ liệu đánh giá và dự báo | Chưa có API | UC-056 | `is_completed` | Không rõ | Dùng `lesson_summary_id` + `student_id` |
| ENT-020 | Đơn xin phép | Ghi nhận nghỉ, đi muộn, về sớm | `leave_requests` / `LeaveRequest` | Transaction | Theo dõi ngoại lệ tham gia học | Chưa có API | UC-057 | `status`, `leave_type` | Soft delete | Workflow duyệt chưa hiện thực |
| ENT-021 | Lead tư vấn | Đầu mối tuyển sinh | `consultations` / `Consultation` | Transaction | Theo dõi khách hàng tiềm năng | Chưa có API | UC-058 | `status` | Không rõ | Chưa có màn hình |
| ENT-022 | Tài liệu giảng dạy | File nội dung dùng trong giảng dạy | `materials` / `Material` | Transaction/content | Lưu file, metadata, trạng thái duyệt | Upload material | UC-049..UC-053 | `status`, `latest_label_id` | Soft delete | Runtime migration hiện không auto tạo |
| ENT-023 | Nhãn kiểm duyệt | Bảng mã SAFE/WARNING/DANGER | `labels` / `Label` | Master/reference | Chuẩn hóa mức độ cảnh báo tài liệu | SQL seed + material flow | UC-049..UC-053 | `code`, `severity` | Không rõ |  |
| ENT-024 | Nhật ký audit tài liệu | Lưu OCR text, reasoning, confidence | `audit_logs` / `AuditLog` | History/audit | Truy vết AI/OCR | Upload material | UC-049..UC-053 | `status`, timestamps | History |  |
| ENT-025 | Quyết định duyệt tài liệu | Quyết định thủ công cuối cùng | `approval_decisions` / `ApprovalDecision` | Transaction/history | Ghi approve/reject của compliance | Review material | UC-053 | `approved`, `decided_at` | History | Người duyệt hiện nhập tay ID |

---

# 9. Gói dữ liệu phục vụ vẽ ERD

## 9.1 Logical relationships

| Mã quan hệ | Thực thể A | Thực thể B | Ý nghĩa nghiệp vụ của quan hệ | Bội số | Tính bắt buộc | FK / liên kết trung gian | Cascade / delete behavior | Ghi chú |
|---|---|---|---|---|---|---|---|---|
| REL-001 | User | UserOTP | Một user có nhiều OTP theo thời gian | 1-N | OTP bắt buộc thuộc một user | `user_otps.user_id` | Không rõ từ entity, logic use case có mark used |  |
| REL-002 | User | PasswordReset | Một user có nhiều request reset | 1-N | PasswordReset bắt buộc thuộc user | `password_resets.user_id` | Không rõ |  |
| REL-003 | Program | Course | Một chương trình gồm nhiều khóa học và một khóa học có thể thuộc nhiều chương trình | N-N | Optional cả hai phía | `program_courses` | Liên kết mapping |  |
| REL-004 | Program | Objective | Một chương trình có nhiều objective | 1-N | Objective thuộc program | `objectives.program_id` | OnDelete cascade theo migration/entity? cần verify | Chưa có API |
| REL-005 | Program | Outcome | Một chương trình có nhiều outcome | 1-N | Outcome thuộc program | `outcomes.program_id` | OnDelete cascade/SET NULL với objective | Chưa có API |
| REL-006 | Objective | Outcome | Outcome có thể gắn về objective | 1-N | Objective là optional | `outcomes.objective_id` | OnDelete SET NULL |  |
| REL-007 | Class | Program | Một lớp có thể gắn một chương trình | N-1 | Optional | `classes.program_id` | Không rõ |  |
| REL-008 | Class | Course | Một lớp có thể gắn một khóa học chính | N-1 | Optional | `classes.course_id` | Không rõ |  |
| REL-009 | Class | Teacher | Một lớp có thể gắn một giáo viên chính | N-1 | Optional | `classes.teacher_id` | Không rõ |  |
| REL-010 | Class | Room | Một lớp có thể gắn một phòng mặc định | N-1 | Optional | `classes.room_id` | Không rõ |  |
| REL-011 | Class | ClassSchedule | Một lớp có nhiều lịch tuần chuẩn | 1-N | Optional nhưng scheduling phụ thuộc mạnh vào dữ liệu này | `class_schedules.class_id` | Không rõ |  |
| REL-012 | ClassSchedule | Shift | Một lịch tuần tham chiếu một ca học | N-1 | Bắt buộc | `class_schedules.shift_id` | Không rõ |  |
| REL-013 | ClassSchedule | Room | Một lịch tuần có thể khóa vào một phòng cụ thể | N-1 | Optional | `class_schedules.room_id` | Không rõ |  |
| REL-014 | Class | Enrollment | Một lớp có nhiều enrollment | 1-N | Optional | `enrollments.class_id` | OnDelete cascade theo entity? |  |
| REL-015 | Student | Enrollment | Một học viên có nhiều enrollment theo lớp | 1-N | Optional | `enrollments.student_id` | OnDelete cascade theo entity? |  |
| REL-016 | Class | Lesson | Một lớp có nhiều lesson | 1-N | Optional | `lessons.class_id` | OnDelete CASCADE | Entity có constraint CASCADE |
| REL-017 | Teacher | Lesson | Một giáo viên có nhiều lesson | 1-N | Optional trên lesson | `lessons.teacher_id` | Không rõ |  |
| REL-018 | Room | Lesson | Một phòng có nhiều lesson | 1-N | Optional trên lesson | `lessons.room_id` | Không rõ |  |
| REL-019 | Lesson | Attendance | Một lesson có nhiều attendance | 1-N | Optional | `attendances.lesson_id` | Cascade qua FK migration/entity? | Chưa có API |
| REL-020 | Student | Attendance | Một student có nhiều attendance | 1-N | Optional | `attendances.student_id` | Cascade? |  |
| REL-021 | Lesson | LessonSummary | Một lesson có tối đa một summary | 1-1 | Summary bắt buộc thuộc lesson | `lesson_summaries.lesson_id` unique | OnDelete CASCADE |  |
| REL-022 | LessonSummary | AcademicRecord | Một summary có nhiều academic record theo học viên | 1-N | AcademicRecord bắt buộc thuộc summary | `academic_records.lesson_summary_id` | OnDelete CASCADE |  |
| REL-023 | Student | AcademicRecord | Một student có nhiều academic record | 1-N | Bắt buộc | `academic_records.student_id` | OnDelete CASCADE |  |
| REL-024 | Student | LeaveRequest | Một học viên có nhiều đơn xin phép | 1-N | Bắt buộc | `leave_requests.student_id` | OnDelete CASCADE theo migration |  |
| REL-025 | Class | LeaveRequest | Đơn xin phép có thể gắn lớp | 1-N | Optional | `leave_requests.class_id` | SET NULL? |  |
| REL-026 | Lesson | LeaveRequest | Đơn xin phép có thể gắn lesson | 1-N | Optional | `leave_requests.lesson_id` | SET NULL? |  |
| REL-027 | Teacher | Material | Một giáo viên có nhiều tài liệu | 1-N | Bắt buộc | `materials.teacher_id` | Không rõ |  |
| REL-028 | Material | Label | Material có latest label | N-1 | Optional | `materials.latest_label_id` | SET NULL |  |
| REL-029 | Material | AuditLog | Một material có nhiều audit log | 1-N | Optional | `audit_logs.material_id` | OnDelete CASCADE |  |
| REL-030 | AuditLog | Label | Một audit log có thể gắn nhãn | N-1 | Optional | `audit_logs.label_id` | SET NULL |  |
| REL-031 | Material | ApprovalDecision | Một material có nhiều quyết định duyệt theo lịch sử | 1-N | Optional | `approval_decisions.material_id` | OnDelete CASCADE |  |

## 9.2 Phân loại thực thể

- **Thực thể lõi nghiệp vụ**: User, Teacher, Student, Course, Program, Class, Room, Shift.
- **Thực thể vận hành**: ClassSchedule, Enrollment, Lesson, Material.
- **Thực thể hỗ trợ / kỹ thuật**: UserOTP, PasswordReset, Label.
- **Thực thể history / audit**: AuditLog, ApprovalDecision.
- **Thực thể mapping**: ProgramCourse, Enrollment.
- **Thực thể học thuật mở rộng**: Attendance, LessonSummary, AcademicRecord, LeaveRequest, Consultation, Objective, Outcome.

## 9.3 Gợi ý phạm vi ERD

- **ERD lõi BA nên vẽ trước**:
  - User
  - Teacher
  - Student
  - Course
  - Program
  - ProgramCourse
  - Class
  - ClassSchedule
  - Room
  - Shift
  - Enrollment
  - Lesson
- **ERD mở rộng**:
  - Attendance
  - LessonSummary
  - AcademicRecord
  - LeaveRequest
  - Consultation
- **ERD kỹ thuật / kiểm duyệt**:
  - Material
  - Label
  - AuditLog
  - ApprovalDecision
  - UserOTP
  - PasswordReset

---

# 10. Vòng đời và trạng thái

| Thực thể | Danh sách trạng thái | Ý nghĩa từng trạng thái | Điều kiện vào trạng thái | Điều kiện ra trạng thái | Chuyển trạng thái hợp lệ | Chuyển trạng thái không hợp lệ / chưa rõ | Ai được phép chuyển | Tác dụng phụ sau chuyển trạng thái | Điểm thiếu / chưa rõ trong code |
|---|---|---|---|---|---|---|---|---|---|
| User | `is_active=true/false`, `role` | false = chưa kích hoạt hoặc bị vô hiệu; true = có thể login | Register tạo false; OTP verify chuyển true | Chưa thấy deactivate flow | false -> true qua verify OTP | true -> false chưa có use case | System / user flow register | Cho phép login khi true | Thiếu deactivate/suspend |
| Teacher | `status` ACTIVE/INACTIVE (theo UI/use case) | GV đang hoạt động hoặc không | Create mặc định ACTIVE | Update teacher | ACTIVE <-> INACTIVE | Chưa có rule tác động tới scheduling/class | Admin | Không thấy side effect | Mức enum chưa khóa chặt ở backend update |
| Student | `status` tự do từ UI | Học viên đang học / nghỉ / khác | Create/update student | Update student | Không rõ | Không rõ | Auth user hiện tại | Không thấy side effect | Enum/status chưa chuẩn hóa |
| Course | `status` ACTIVE/... | Khóa học có dùng được không | Create default ACTIVE | Update | Không rõ | Không rõ | Admin | Solver phụ thuộc session_count/duration chứ chưa check status tại class/course mapping | Chưa có state machine |
| Program | `published_at`, `archived_at`, có thể từng có `status` | Gợi ý draft/published/archived nhưng chưa được dùng rõ | Tạo program | Chưa có publish/archive API | Không rõ | Không rõ | Admin | Không rõ | Mâu thuẫn giữa migration và entity |
| Class | OPEN/CLOSED/CANCELLED | OPEN = cho scheduling và vận hành; CLOSED/CANCELLED = ngừng | Create default OPEN | Update class | OPEN -> CLOSED/CANCELLED | Chưa thấy rule reopen hoặc effects | Admin | Scheduling chỉ load OPEN | Chưa có automation khi class end_date qua hạn |
| Enrollment | entity default APPLIED; use case tạo ENROLLED | APPLIED = chờ duyệt; ENROLLED = vào lớp | Tạo enrollment theo entity hoặc use case | Chưa có approve/reject flow | APPLIED -> ENROLLED? -> REJECTED? | Chưa có thực thi đầy đủ | Admin theo hiện trạng | Roster lớp thay đổi | Lifecycle không đồng nhất |
| Shift | `is_active` true/false, session_type | active = được scheduling dùng | Create/update shift | Update shift | true <-> false | Delete shift có thể làm class_schedule mồ côi? chưa thấy guard | Admin | Scheduling preview chỉ lấy active shifts | Thiếu rule chặn xóa shift đang dùng |
| Preview scheduling | FAILED/PARTIAL/COMPLETED | tình trạng preview | Tạo preview | Commit hoặc re-run | FAILED/PARTIAL/COMPLETED | Không thấy trạng thái EXPIRED | Admin/System | Lưu in-memory | Mất khi restart |
| Commit scheduling | COMMITTED | Đã tạo lesson từ preview | Commit thành công | Không thấy rollback business | COMPLETED preview -> COMMITTED | Re-commit cùng preview có thể gây duplicate nếu không bị overlap | Admin | Lessons được tạo | Thiếu cờ “preview đã commit” |
| Material | UPLOADED, SCANNING, AI_REVIEWED, APPROVED, REJECTED | Trạng thái tài liệu | Upload tạo SCANNING rồi AI_REVIEWED; review chuyển APPROVED/REJECTED | Review tiếp | SCANNING -> AI_REVIEWED -> APPROVED/REJECTED | Không rõ APPROVED -> REJECTED lại hay re-review | Teacher/System/Compliance | Tạo audit log và approval history | Chưa có state machine formal |
| LeaveRequest | PENDING, APPROVED, REJECTED? | Chờ duyệt / đã duyệt / bị từ chối | Chưa có use case | Chưa có use case | PENDING -> APPROVED/REJECTED (suy luận) | Chưa có code thực thi | Student/Admin (suy luận) | Có thể ảnh hưởng attendance | Missing API |

---

# 11. Ma trận phân quyền và kiểm soát truy cập

## 11.1 Ma trận quyền chi tiết

| Chức năng / Module | Guest | Auth User bất kỳ | Teacher | Admin | Compliance | Ownership ràng buộc | Bằng chứng kỹ thuật | Rủi ro nếu thiếu kiểm soát |
|---|---|---|---|---|---|---|---|---|
| Register/Login/Forgot/Reset | Có | Có | Có | Có | Có | Không | auth routes | Thấp |
| `/v1/auth/me` | Không nên public | Route hiện public nhưng controller cần auth context | Có thể lỗi | Có thể lỗi | Có thể lỗi | Không | `user/controller.go`, `user/v1.go` | Lỗi profile + rò logic auth |
| Teacher list/detail | Có | Có | Có | Có | Có | Không | teacher routes GET public | Có thể lộ dữ liệu nội bộ giáo viên |
| Teacher create/update/delete | Không | Không | Không | Có | Không | Không | teacher controller routes | Thấp |
| Student list/detail/create/update/delete | Không | Có | Có | Có | Có | Không | student controller uses only AuthMiddleware | Rất cao: user không phù hợp có thể CRUD học viên |
| Course list/detail | Có | Có | Có | Có | Có | Không | course controller GET public | Trung bình |
| Course create/update/delete | Không | Không | Không | Có | Không | Không | course controller | Thấp |
| Program list/detail | Có | Có | Có | Có | Có | Không | program controller GET public | Trung bình |
| Program create/update/delete/add/remove course | Không | Không | Không | Có | Không | Không | program controller | Thấp |
| Class list/detail/roster | Có | Có | Có | Có | Có | Không | class controller GET public | Trung bình-cao nếu roster lộ |
| Class create/update/delete/enroll/remove/assign teacher | Không | Không | Không | Có | Không | Không | class controller | Thấp |
| Room list/detail | Có | Có | Có | Có | Có | Không | room controller GET public | Trung bình |
| Room create/update/delete | Không | Không | Không | Có | Không | Không | room controller | Thấp |
| Shift list/detail/CRUD | Không | Không | Không | Có | Không | Không | shift controller uses auth+admin for all routes | Thấp |
| Scheduling preview/get/commit/benchmark | Không | Không | Không | Có | Không | Không | scheduling controller auth+admin | Thấp |
| Material upload | Không | Không | Có | Không | Không | Không, teacher_id nhập tay | material controller | Cao: teacher có thể upload thay teacher khác nếu role check không phối hợp ownership |
| Material list/get/download | Không | Có | Có | Có | Có | Không | material controller auth only | Cao: người không liên quan vẫn xem/tải tài liệu |
| Material review | Không | Có | Có | Có | Có | Không, compliance_officer_id nhập tay | material controller auth only | Rất cao: bất kỳ user đăng nhập đều approve/reject tài liệu |

## 11.2 Endpoint không có role check / có nguy cơ lộ

| ID | Endpoint / nhóm endpoint | Vấn đề |
|---|---|---|
| PERM-GAP-01 | `/v1/students/*` | Chỉ có auth, không role-specific |
| PERM-GAP-02 | `/v1/materials/:id/review` | Chỉ auth, không compliance/admin role |
| PERM-GAP-03 | `/v1/materials`, `/v1/materials/:id`, `/v1/materials/:id/download` | Không ràng buộc ownership teacher/compliance |
| PERM-GAP-04 | `/v1/teachers`, `/v1/teachers/:id` | Public read, có thể không phù hợp |
| PERM-GAP-05 | `/v1/classes/:id/students` | Roster lớp public read qua class GET routes | 
| PERM-GAP-06 | `/v1/auth/me` | Route không gắn auth middleware |

## 11.3 Rủi ro nghiệp vụ từ authorization gaps

1. Người dùng bất kỳ đã đăng nhập có thể sửa/xóa hồ sơ học viên.
2. Người dùng bất kỳ đã đăng nhập có thể duyệt hoặc từ chối tài liệu.
3. Tài liệu giảng dạy có thể bị xem/tải không đúng đối tượng.
4. Roster lớp và dữ liệu giáo viên/lớp/phòng có thể bị công khai ngoài dự kiến.
5. Flow profile hiện tại không đáng tin vì route và kiểu ID không đúng.

---

# 12. Danh mục API / màn hình / báo cáo

## 12.1 API Catalog

| Endpoint | Method | Module | Actor | Mục đích nghiệp vụ | Input | Output | Thực thể bị tác động | Side effects | Permission | Trạng thái triển khai |
|---|---|---|---|---|---|---|---|---|---|---|
| `/api/v1/auth/register` | POST | Auth | Guest | Tạo tài khoản mới | email, full_name, password_enc | user_id | User, UserOTP | Gửi OTP async | Public | Implemented |
| `/api/v1/auth/verify-otp` | POST | Auth | Guest/User | Kích hoạt tài khoản | user_id, otp | message | User, UserOTP | Activate user | Public | Implemented |
| `/api/v1/auth/login` | POST | Auth | User | Đăng nhập | username(email), password | access/refresh token + user | User | Issue JWT | Public | Implemented |
| `/api/v1/auth/refresh` | POST | Auth | User | Refresh access token | refresh_token | access_token | N/A | New JWT | Public | Implemented |
| `/api/v1/auth/logout` | POST | Auth | User | Logout logic | token | message | N/A | None persistent | Public | Implemented |
| `/api/v1/auth/forgot-password` | POST | Auth | Guest/User | Gửi reset link | email | message | PasswordReset | Gửi email | Public | Implemented |
| `/api/v1/auth/reset-password` | POST | Auth | Guest/User | Đặt lại mật khẩu | token, new_password_enc | message | PasswordReset, User | Update password | Public | Implemented |
| `/api/v1/auth/change-password` | POST | Auth | User | Đổi mật khẩu | old_password_enc, new_password_enc | message | User | Update password | Route public nhưng controller cần auth context | Partial / risky |
| `/api/v1/auth/me` | GET | Auth | User | Lấy profile hiện tại | JWT context | UserResponse | User | None | Route public, controller expects auth | Defective |
| `/api/v1/teachers` | GET | Teacher | Guest/Admin | Danh sách giáo viên | search, status, employment_type, page... | teacher list + pagination | Teacher | None | Public | Implemented |
| `/api/v1/teachers/:id` | GET | Teacher | Guest/Admin | Chi tiết giáo viên | id | teacher detail | Teacher | None | Public | Implemented |
| `/api/v1/teachers` | POST | Teacher | Admin | Tạo giáo viên | CreateTeacherRequest | teacher detail | Teacher | Create | ADMIN | Implemented |
| `/api/v1/teachers/:id` | PUT | Teacher | Admin | Cập nhật giáo viên | UpdateTeacherRequest | teacher detail | Teacher | Update | ADMIN | Implemented |
| `/api/v1/teachers/:id` | DELETE | Teacher | Admin | Xóa mềm giáo viên | id | message | Teacher | Soft delete | ADMIN | Implemented |
| `/api/v1/teachers/:id/timetable` | GET | Teacher | Admin/Teacher | Lịch dạy GV | from, to | lessons | Lesson | None | Public | Implemented |
| `/api/v1/teachers/:id/stats/teaching-hours` | GET | Teacher | Admin/Teacher | Thống kê giờ dạy | from, to, group_by | total_hours + breakdown | Lesson | None | Public | Implemented |
| `/api/v1/students` | GET/POST | Student | Auth user | Danh sách / tạo học viên | search/status hoặc CreateStudentRequest | list / detail | Student | Create | Auth only | Implemented with auth gap |
| `/api/v1/students/:id` | GET/PUT/DELETE | Student | Auth user | Xem/sửa/xóa học viên | id + update body | detail / message | Student | Update/delete | Auth only | Implemented with auth gap |
| `/api/v1/courses` | GET | Course | Guest/Admin | Danh sách khóa học | search/status/subject/page | list | Course | None | Public | Implemented |
| `/api/v1/courses/:id` | GET | Course | Guest/Admin | Chi tiết khóa học | id | detail | Course | None | Public | Implemented |
| `/api/v1/courses` | POST | Course | Admin | Tạo khóa học | CreateCourseRequest | course detail | Course | Create | ADMIN | Implemented |
| `/api/v1/courses/:id` | PUT/DELETE | Course | Admin | Sửa/xóa khóa học | id + body | detail / message | Course | Update/delete | ADMIN | Implemented |
| `/api/v1/programs` | GET/POST | Program | Guest/Admin | Danh sách / tạo chương trình | track/page hoặc create body | list / detail | Program | Create | GET public, POST ADMIN | Implemented |
| `/api/v1/programs/:id` | GET/PUT/DELETE | Program | Guest/Admin | Chi tiết/sửa/xóa chương trình | id + body | detail / message | Program | Update/delete | GET public, write ADMIN | Implemented |
| `/api/v1/programs/:id/courses` | POST/DELETE | Program | Admin | Thêm/gỡ khóa học | course_ids[] | message | ProgramCourse | Create/delete mapping | ADMIN | Implemented |
| `/api/v1/classes` | GET/POST | Class | Guest/Admin | Danh sách / tạo lớp | filters hoặc create body | list / detail | Class | Create | GET public, POST ADMIN | Implemented |
| `/api/v1/classes/:id` | GET/PUT/DELETE | Class | Guest/Admin | Chi tiết/sửa/xóa lớp | id + update body | detail / message | Class | Update/delete | GET public, write ADMIN | Implemented |
| `/api/v1/classes/:id/students` | GET | Enrollment | Guest/Admin | Lấy roster lớp | class_id | roster snapshot | Enrollment, Student | None | Public | Implemented |
| `/api/v1/classes/:id/students` | POST/DELETE | Enrollment | Admin | Thêm/xóa học viên khỏi lớp | student_ids[] | output/message | Enrollment | create/delete | ADMIN | Implemented partial |
| `/api/v1/classes/:id/teacher` | PUT | Class | Admin | Gán giáo viên cho lớp | teacher_id | class updated | Class | update | ADMIN | Implemented |
| `/api/v1/rooms` | GET/POST | Room | Guest/Admin | Danh sách / tạo phòng | search/page hoặc create body | list / detail | Room | Create | GET public, POST ADMIN | Implemented partial |
| `/api/v1/rooms/:id` | GET/PUT/DELETE | Room | Guest/Admin | Chi tiết/sửa/xóa phòng | id + body | detail / message | Room | Update/delete | GET public, write ADMIN | Implemented partial |
| `/api/v1/shifts` | GET/POST | Shift | Admin | Danh sách / tạo ca học | filters hoặc create body | list / detail | Shift | Create | ADMIN | Implemented |
| `/api/v1/shifts/:id` | GET/PUT/DELETE | Shift | Admin | Chi tiết/sửa/xóa ca học | id + body | detail / message | Shift | Update/delete | ADMIN | Implemented |
| `/api/v1/scheduling/preview` | POST | Scheduling | Admin | Tạo preview xếp lịch | date_from, date_to, class_ids, teacher_ids, room_ids | PreviewResult | N/A trực tiếp | Save preview in memory | ADMIN | Implemented |
| `/api/v1/scheduling/preview/latest` | GET | Scheduling | Admin | Lấy preview mới nhất | none | PreviewResult | N/A | None | ADMIN | Implemented partial |
| `/api/v1/scheduling/preview/:id` | GET | Scheduling | Admin | Lấy preview theo run_id | id | PreviewResult | N/A | None | ADMIN | Implemented |
| `/api/v1/scheduling/commit` | POST | Scheduling | Admin | Commit preview thành lessons | run_id | commit output | Lesson | Create lessons | ADMIN | Implemented |
| `/api/v1/scheduling/benchmark` | POST | Scheduling | Admin | Lấy contract benchmark solver | same filters | benchmark contract | N/A | None | ADMIN | Partial |
| `/api/v1/materials/upload` | POST | Material | Teacher | Upload tài liệu | multipart form: teacher_id, title, description, file | material output | Material, AuditLog | Save file + scan stub | TEACHER | Implemented partial |
| `/api/v1/materials` | GET | Material | Auth user | Danh sách tài liệu | teacher_id, status, queue | materials | Material | None | Auth only | Implemented with gap |
| `/api/v1/materials/flagged` | GET | Material | Auth user | Danh sách cần duyệt | none | materials | Material | None | Auth only | Implemented with gap |
| `/api/v1/materials/:id` | GET | Material | Auth user | Chi tiết tài liệu | id | material detail | Material | None | Auth only | Implemented with gap |
| `/api/v1/materials/:id/download` | GET | Material | Auth user | Tải file | id | file attachment | Material/file storage | Download file | Auth only | Implemented with gap |
| `/api/v1/materials/:id/review` | POST | Material | Auth user / Compliance | Duyệt tài liệu | compliance_officer_id, approved... | material output | Material, ApprovalDecision | Update status | Auth only | Implemented with gap |

## 12.2 Screen/UI Candidate Catalog

| Mã màn hình | Tên màn hình | Mục đích | Actor sử dụng | Trường dữ liệu chính | Hành động chính | Use case liên quan | Ghi chú |
|---|---|---|---|---|---|---|---|
| SCR-001 | LoginPage | Đăng nhập | Guest/User | username, password | login | UC-003 | Implemented |
| SCR-002 | RegisterPage | Đăng ký tài khoản | Guest | email, full_name, password_enc | register | UC-001 | Implemented |
| SCR-003 | ForgotPasswordPage | Quên mật khẩu | Guest/User | email | send reset | UC-006 | Implemented |
| SCR-004 | ResetPasswordPage | Đặt lại mật khẩu | Guest/User | token, new_password_enc | reset | UC-007 | Implemented |
| SCR-005 | ChangePasswordPage | Đổi mật khẩu | User | old_password_enc, new_password_enc | change password | UC-008 | Implemented |
| SCR-006 | ProfilePage | Hồ sơ tài khoản | User | name, email, code, role | logout, navigate change password | UC-009 | Implemented |
| SCR-007 | AdminOverview | Tổng quan admin | Admin | mock stats/charts | view | reporting candidate | Dữ liệu mock |
| SCR-008 | TeachersPage | Danh sách giáo viên | Admin | filter search/status/employment | view/create/edit/detail | UC-010..013 | Implemented |
| SCR-009 | TeacherFormPage | Tạo/sửa giáo viên | Admin | teacher fields | save | UC-010, UC-012 | Implemented |
| SCR-010 | TeacherDetailPage | Chi tiết giáo viên | Admin | teacher detail, timetable summary | view/edit | UC-012, UC-014, UC-015 | Implemented |
| SCR-011 | StudentsPage | Quản lý học viên | Admin/Auth user | student fields, filter search/status | create/edit/delete | UC-016..019 | Implemented |
| SCR-012 | CoursePage | Quản lý khóa học | Admin | course fields | CRUD | UC-020..023 | Implemented |
| SCR-013 | ProgramPage | Quản lý chương trình | Admin | program fields, track, linked courses | CRUD, open detail, link/unlink courses | UC-024..029 | Implemented |
| SCR-014 | ProgramDetailDialog | Chi tiết chương trình | Admin | course links, approval note, effectivity | link/unlink courses | UC-028, UC-029 | Implemented |
| SCR-015 | ClassesPage | Quản lý lớp học | Admin | class list, search/status | create/edit/delete/view detail | UC-030..037 | Implemented |
| SCR-016 | ClassDetailDialog | Chi tiết lớp học | Admin | roster, teacher, class info | add/remove students | UC-032, UC-035, UC-036 | Implemented |
| SCR-017 | RoomsPage | Quản lý phòng học | Admin | room list, capacity, address | CRUD | UC-038..041 | Implemented partial |
| SCR-018 | ShiftsPage | Quản lý ca học | Admin | code, name, time, duration, type, active | CRUD | UC-042..044 | Implemented |
| SCR-019 | SchedulingPage | Xếp lịch thông minh | Admin | date range, selected classes/teachers/rooms, preview | run preview, view conflicts, commit | UC-045..048 | Implemented |
| SCR-020 | TeacherOverview | Dashboard giáo viên | Teacher | placeholder stats | view | teacher portal | Partial/mock |
| SCR-021 | TeacherDocumentsPage | Tài liệu giảng dạy | Teacher | teacher_id, title, description, file, material list | upload, view detail, download | UC-049..052 | Implemented partial |
| SCR-022 | ComplianceOverview | Dashboard compliance | Compliance | placeholder status | view | material moderation | Partial/mock |
| SCR-023 | ComplianceQueuePage | Hàng chờ duyệt tài liệu | Compliance | queue flagged, reasoning, decision form | view, review | UC-050..053 | Implemented partial |
| SCR-024 | MaterialDetailDialog | Chi tiết tài liệu | Teacher/Compliance/Admin | file info, latest label, audit logs, decision | download | UC-051, UC-052 | Implemented |
| SCR-025 | StudentOverview | Dashboard học viên | Student | mock GPA, credits, attendance | view | student portal | Partial/mock |
| SCR-026 | PlaceholderPage Teacher Attendance | Placeholder điểm danh giáo viên | Teacher | none | none | UC-054 | Placeholder |
| SCR-027 | PlaceholderPage Student Timetable | Placeholder thời khóa biểu học viên | Student | none | none | future timetable | Placeholder |
| SCR-028 | PlaceholderPage Student Results | Placeholder kết quả học tập | Student | none | none | UC-056 | Placeholder |
| SCR-029 | PlaceholderPage Student Leaves | Placeholder đơn xin phép | Student | none | none | UC-057 | Placeholder |
| SCR-030 | PlaceholderPage Consulting | Placeholder tư vấn tuyển sinh | Student/Guest | none | none | UC-058 | Placeholder |

## 12.3 Report / Statistics / Dashboard Catalog

| Mã báo cáo | Tên | Mục đích | Actor | Nguồn dữ liệu | Bộ lọc | Output | Giá trị nghiệp vụ |
|---|---|---|---|---|---|---|---|
| REP-001 | Admin Overview Mock Dashboard | Minh họa tổng quan hệ thống | Admin | Mock JSON trong UI | Không | cards + charts | Chỉ có giá trị demo giao diện |
| REP-002 | Teacher Teaching Hours Stats | Thống kê giờ dạy GV | Admin/Teacher | Lesson | from, to, group_by | total_hours + breakdown | Giá trị vận hành thật |
| REP-003 | Scheduling Preview Summary | Tóm tắt số buổi xếp được/chưa xếp/xung đột | Admin | PreviewResult | date range, class/teacher/room filter | summary + conflicts + assignments | Giá trị vận hành cao |
| REP-004 | Compliance Queue | Danh sách tài liệu cần duyệt | Compliance/Admin | Material, Label, AuditLog | queue/status/search | table queue | Giá trị kiểm soát nội dung |

---

# 13. CRUD Matrix

| Entity \ Functional Group | Auth | Teacher Mgmt | Student Mgmt | Course Mgmt | Program Mgmt | Class Mgmt | Enrollment | Room Mgmt | Shift Mgmt | Scheduling | Material | Academic Ops |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| User | C/R/U |  |  |  |  |  |  |  |  |  |  |  |
| UserOTP | C/R/U |  |  |  |  |  |  |  |  |  |  |  |
| PasswordReset | C/R/U |  |  |  |  |  |  |  |  |  |  |  |
| Teacher |  | C/R/U/D |  |  |  | R |  |  |  | R | R (owner implied) |  |
| Student |  |  | C/R/U/D |  |  | R | R |  |  |  |  | R? |
| Course |  |  |  | C/R/U/D | R | R |  |  |  | R |  |  |
| Program |  |  |  |  | C/R/U/D | R |  |  |  | R |  |  |
| ProgramCourse |  |  |  |  | C/R/U/D |  |  |  |  |  |  |  |
| Class |  |  |  |  |  | C/R/U/D | R | R |  | R |  |  |
| ClassSchedule |  |  |  |  |  | R? |  |  |  | R |  |  |
| Enrollment |  |  |  |  |  | R | C/R/D |  |  |  |  |  |
| Room |  |  |  |  |  | R |  | C/R/U/D |  | R |  |  |
| Shift |  |  |  |  |  |  |  |  | C/R/U/D | R |  |  |
| Lesson |  | R |  |  |  |  |  |  |  | G/C/R |  | R? |
| Attendance |  |  |  |  |  |  |  |  |  |  |  | C/R/U? |
| LessonSummary |  |  |  |  |  |  |  |  |  |  |  | C/R/U? |
| AcademicRecord |  |  |  |  |  |  |  |  |  |  |  | C/R/U? |
| LeaveRequest |  |  |  |  |  |  |  |  |  |  |  | C/R/U/Approve? |
| Material |  |  |  |  |  |  |  |  |  |  | C/R/Approve |  |
| AuditLog |  |  |  |  |  |  |  |  |  |  | C/R |  |
| ApprovalDecision |  |  |  |  |  |  |  |  |  |  | C/R |  |
| Consultation |  |  |  |  |  |  |  |  |  |  |  | C/R/U? |

Ký hiệu bổ sung:
- `Approve` tương ứng quyết định duyệt/từ chối.
- `Generate` = hệ thống sinh dữ liệu tự động, ví dụ Lesson từ Scheduling.

---

# 14. Phân tích khoảng trống và yêu cầu còn thiếu

| Mã gap | Mô tả | Miền bị ảnh hưởng | Tác động nghiệp vụ | Mức độ nghiêm trọng | Dấu hiệu trong code | Câu hỏi BA cần hỏi | Stakeholder cần xác nhận |
|---|---|---|---|---|---|---|---|
| GAP-001 | Student CRUD không có role check | Student | Người dùng không phù hợp có thể sửa/xóa dữ liệu học viên | Cao | `student/controller.go` | Ai được phép quản lý học viên? | PO, Security, Admin |
| GAP-002 | Material review không có compliance/admin role check | Material | Người dùng bất kỳ có thể duyệt/từ chối tài liệu | Cao | `material/controller.go` | Ai được quyền review tài liệu? | PO, Compliance |
| GAP-003 | Material ownership không enforced | Material | Giáo viên có thể upload/xem/tải sai ownership | Cao | `TeacherDocumentsPage.tsx`, `material/v1.go` | Tài liệu là của giáo viên hay của trung tâm? | PO, Compliance |
| GAP-004 | `/auth/me` route và implementation lỗi | Auth | Không dùng được profile đúng, gây lỗi xác thực | Cao | `user/controller.go`, `user/v1.go`, `get_user_by_id.go` | Profile hiện tại cần phục vụ use case nào? | Tech Lead, PO |
| GAP-005 | Program có mâu thuẫn status giữa migration và entity | Program | Không rõ lifecycle chương trình đào tạo | Trung bình-cao | `7_create_programs_table.up.sql`, `entities/program.go` | Program có trạng thái nghiệp vụ không? | Đào tạo, PO |
| GAP-006 | Room DTO lệch entity/use case | Room | API contract không thống nhất với dữ liệu nghiệp vụ | Trung bình-cao | `room/dto.go`, `entities/room.go`, room use cases | Phòng có cần status/location hay code/address? | Admin vận hành, PO |
| GAP-007 | Enrollment use case không đếm enrollment hiện có | Enrollment | Có thể vượt sĩ số thực, sai roster | Cao | `enroll_students.go` | Quy tắc sĩ số và duplicate enrollment cụ thể là gì? | Admin vận hành, PO |
| GAP-008 | Không có API quản trị `class_schedule` | Scheduling/Class | Không hoàn thiện flow cấu hình lịch tuần cho lớp | Cao | Có entity nhưng không có controller/use case | Ai cấu hình lịch tuần và ở màn nào? | Admin vận hành, PO |
| GAP-009 | Attendance chưa có API, enum status chưa rõ | Attendance | Không triển khai được chuyên cần thật | Cao | Chỉ có entity/migration | Hệ thống cần mấy trạng thái chuyên cần? | Giáo viên, đào tạo |
| GAP-010 | LessonSummary chưa có API | Academic ops | Không tổng kết được buổi học | Trung bình-cao | Chỉ có entity/migration | Có cần bắt buộc tổng kết sau mỗi lesson không? | Teacher, PO |
| GAP-011 | AcademicRecord chưa có API | Academic ops / ML | Không có dữ liệu thực cho đánh giá và dự báo | Cao | Chỉ có entity/migration | Chấm điểm ở level buổi, khóa hay kỳ? | Đào tạo, PO |
| GAP-012 | LeaveRequest chưa có API | Leave | Không xử lý được nghỉ học/đi muộn thực tế | Trung bình-cao | Chỉ có entity/migration | Quy trình duyệt đơn là gì? | Student service, PO |
| GAP-013 | Consultation chưa có API | Consultation | Không theo dõi lead tuyển sinh | Trung bình | Chỉ có entity/migration | Lead intake có trong scope hiện tại không? | Marketing/PO |
| GAP-014 | Admin overview là mock data | Reports | Báo cáo vận hành chưa có giá trị thật | Trung bình | `AdminOverview.tsx` | Dashboard giai đoạn này cần KPI nào? | PO, Admin |
| GAP-015 | Scheduling benchmark mới là contract | Scheduling | Chưa có số liệu so sánh solver để chứng minh chọn thuật toán | Cao | `benchmark.go` | Tiêu chí benchmark chính là gì? | GVHD, PO, Tech Lead |
| GAP-016 | Preview store là in-memory | Scheduling | Mất preview khi restart, khó audit/review lại | Trung bình | `preview_store.go` | Có cần lưu preview lịch sử vào DB không? | Admin, PO |
| GAP-017 | Không có idempotency cho commit preview | Scheduling/Lesson | Có nguy cơ tạo lesson trùng nếu commit nhiều lần trong điều kiện đặc biệt | Trung bình | `commit_preview.go` | Có cần cờ “đã commit” không? | Tech Lead, PO |
| GAP-018 | Material/audit tables không nằm trong AutoMigrate runtime | Material | Deploy mới có thể thiếu bảng nếu chỉ chạy runtime migrate | Cao | `migration.go` không migrate materials/audit tables | Quy trình migration chuẩn của dự án là gì? | Tech Lead |
| GAP-019 | V2 auth controller còn TODO/panic | Auth | Gây nhầm phạm vi API versioning | Thấp-trung bình | `user/v2.go` | Có thực sự cần v2 trong scope hiện tại không? | Tech Lead, PO |
| GAP-020 | Queue abstractions có nhưng luồng chính chưa dùng | Integration | Kiến trúc có phần dư, dễ gây nhầm đang có async thật | Thấp | queue packages, mailer direct SMTP | Có cần queue production ở scope đồ án không? | Tech Lead |
| GAP-021 | ML/predictive chỉ ở tài liệu, chưa có code | Predictive | Không thể mô tả là chức năng hiện hữu | Cao | Chỉ có docs/kế hoạch | Predictive nằm ở phase nào? | GVHD, PO |
| GAP-022 | Compliance/AI moderation còn trong code dù scope sau này có thể bỏ | Product scope | Dễ lệch giữa hiện trạng code và phạm vi báo cáo | Trung bình | Frontend routes/pages + material module | Báo cáo nên tính module này là hiện trạng hay loại khỏi scope chính? | GVHD, PO |

---

# 15. Bộ câu hỏi phỏng vấn BA

## 15.1 Product Owner

| Câu hỏi | Vì sao cần hỏi | Artifact phụ thuộc | Ưu tiên |
|---|---|---|---|
| Trong phạm vi release chính, những miền nào là “bắt buộc go-live” và những miền nào chỉ là nền dữ liệu? | Hiện code có nhiều entity chưa có API | Scope, Use case, Gap analysis | Cao |
| Tài liệu giảng dạy và compliance moderation có còn trong phạm vi chính không? | Code có module nhưng backlog gần đây có lúc muốn bỏ | Scope, Functional decomposition | Cao |
| Chương trình đào tạo có lifecycle nghiệp vụ riêng không: draft/active/published/archived? | Migration và entity program đang mâu thuẫn | ERD, lifecycle, rules | Cao |
| Có cần phân quyền riêng cho vận hành đào tạo, compliance, tư vấn tuyển sinh không? | Hiện chủ yếu chỉ có ADMIN/TEACHER/STUDENT | Permission matrix, actor catalog | Cao |
| Predictive analytics sẽ là phase tiếp theo hay phạm vi nghiên cứu ngoài sản phẩm chạy? | Chưa có code nhưng có định hướng mạnh | Scope, roadmap | Trung bình-cao |

## 15.2 Admin / Vận hành

| Câu hỏi | Vì sao cần hỏi | Artifact phụ thuộc | Ưu tiên |
|---|---|---|---|
| Ai được quyền tạo/sửa/xóa học viên? | Hiện backend cho mọi auth user | Permission matrix | Cao |
| Khi ghi danh học viên, nghiệp vụ có cần trạng thái chờ duyệt hay ghi danh trực tiếp? | Entity có APPLIED nhưng use case tạo ENROLLED luôn | Use case, lifecycle, rules | Cao |
| Hệ thống có cần chặn học viên ghi danh trùng lớp/trùng khóa không? | Chưa có duplicate check | Rules, workflows | Cao |
| Lịch tuần lớp hiện được cấu hình ở bước nào? | ClassSchedule chưa có API | Workflow, functional decomposition | Cao |
| Có cho phép class dùng nhiều giáo viên hoặc teacher phụ không? | Class hiện chỉ có teacher_id 1-1 | ERD, use case | Trung bình |

## 15.3 Giáo viên

| Câu hỏi | Vì sao cần hỏi | Artifact phụ thuộc | Ưu tiên |
|---|---|---|---|
| Giáo viên có cần tự xem lesson, điểm danh, tổng kết buổi học và học viên theo lớp không? | Teacher portal hiện còn thiếu nhiều flow | Use case, screens | Cao |
| Tài liệu upload là tài sản cá nhân giáo viên hay tài sản chung trung tâm? | Liên quan ownership và quyền tải/xem | Permissions, rules | Cao |
| Có cần versioning tài liệu và re-review khi cập nhật file không? | Hiện material flow chưa có version | Entity, workflow | Trung bình |
| Thống kê giờ dạy cần phục vụ chấm công hay chỉ báo cáo? | Quyết định mức độ chính xác và report | Reports, business rules | Trung bình |

## 15.4 Học viên

| Câu hỏi | Vì sao cần hỏi | Artifact phụ thuộc | Ưu tiên |
|---|---|---|---|
| Học viên cần xem gì trong portal: timetable, kết quả, điểm danh, bài tập, đơn nghỉ? | Student portal hiện mostly placeholder | Screen catalog, use case | Cao |
| Học viên có được gửi đơn nghỉ trực tiếp hay phải qua phụ huynh? | LeaveRequest có model nhưng chưa có flow | Workflow, permissions | Trung bình-cao |
| Kết quả học tập hiển thị theo buổi, theo khóa, hay theo kỳ? | AcademicRecord chưa có flow | ERD, report, screens | Trung bình-cao |

## 15.5 Quản lý đào tạo

| Câu hỏi | Vì sao cần hỏi | Artifact phụ thuộc | Ưu tiên |
|---|---|---|---|
| Objective và Outcome có cần quản trị trong hệ thống không? | Đã có entity nhưng chưa có API | Scope, ERD | Cao |
| Chương trình đào tạo có phải qua bước phê duyệt/publish trước khi dùng mở lớp không? | Có `approved_by_id`, `published_at`, `archived_at` nhưng chưa có flow | Lifecycle, workflow | Cao |
| Scheduling tối ưu theo tiêu chí nào là quan trọng nhất: tránh xung đột, tối ưu phòng, tối ưu phân bổ giờ GV, hay trải đều ca? | Phục vụ benchmark solver | Business rules, benchmark criteria | Cao |
| Có cho phép thay đổi lesson sau khi commit preview không? | Hiện lesson chỉ được sinh tự động | Use case, workflow | Trung bình-cao |

## 15.6 Tài chính

| Câu hỏi | Vì sao cần hỏi | Artifact phụ thuộc | Ưu tiên |
|---|---|---|---|
| `price` ở class và course dùng để làm gì: tham khảo, học phí chuẩn, hay công nợ? | Hiện có field nhưng không có module tài chính | ERD, future scope | Trung bình |
| Có cần báo cáo doanh thu theo lớp/khóa học không? | Admin overview hiện là mock | Reports | Thấp-trung bình |

## 15.7 Bảo mật / Tech Lead

| Câu hỏi | Vì sao cần hỏi | Artifact phụ thuộc | Ưu tiên |
|---|---|---|---|
| Có cần chuẩn hóa RBAC đầy đủ thay vì chỉ role string đơn giản không? | Hiện role check rời rạc và còn lỗ hổng | Permission matrix | Cao |
| Có cần token revocation / session store cho logout/reset password không? | Logout hiện chỉ validate token | Security design | Cao |
| Có cần persist scheduling preview vào DB để audit/history không? | Preview latest đang in-memory | Workflow, architecture | Trung bình-cao |
| Migration chuẩn sẽ dùng AutoMigrate hay SQL migrations? | Hai cơ chế đang song song và lệch phạm vi | Deployment architecture | Cao |
| Có cần derive `teacher_id` / `compliance_officer_id` từ JWT thay vì request body không? | Tránh impersonation | Security rules, API contract | Cao |

---

# 16. Kết luận và gói bàn giao

## 16.1 Top 25 use case ưu tiên

1. UC-001 Đăng ký tài khoản  
2. UC-002 Xác minh OTP email  
3. UC-003 Đăng nhập  
4. UC-006 Quên mật khẩu  
5. UC-007 Đặt lại mật khẩu  
6. UC-008 Đổi mật khẩu  
7. UC-010 Tạo giáo viên  
8. UC-011 Xem danh sách giáo viên  
9. UC-016 Tạo học viên  
10. UC-017 Xem danh sách học viên  
11. UC-020 Tạo khóa học  
12. UC-024 Tạo chương trình đào tạo  
13. UC-028 Gán khóa học vào chương trình  
14. UC-030 Tạo lớp học  
15. UC-032 Xem roster lớp  
16. UC-035 Ghi danh học viên  
17. UC-037 Phân công giáo viên cho lớp  
18. UC-038 Tạo phòng học  
19. UC-042 Tạo ca học  
20. UC-045 Tạo preview xếp lịch  
21. UC-046 Xem preview xếp lịch  
22. UC-047 Commit preview thành lesson  
23. UC-049 Upload tài liệu  
24. UC-051 Xem chi tiết tài liệu  
25. UC-053 Duyệt / từ chối tài liệu

## 16.2 Cây phân rã chức năng đề xuất cuối cùng

- Phân hệ 1: Tài khoản và truy cập
- Phân hệ 2: Danh mục đào tạo
- Phân hệ 3: Vận hành lớp học
- Phân hệ 4: Xếp lịch và lesson
- Phân hệ 5: Tài liệu giảng dạy và kiểm duyệt
- Phân hệ 6: Theo dõi học tập
- Phân hệ 7: Báo cáo và phân tích

## 16.3 Phạm vi ERD BA nên vẽ trước

1. User, Teacher, Student  
2. Course, Program, ProgramCourse  
3. Class, ClassSchedule, Enrollment  
4. Room, Shift  
5. Lesson  

Sau đó mới mở rộng:
- Attendance, LessonSummary, AcademicRecord, LeaveRequest
- Material, Label, AuditLog, ApprovalDecision

## 16.4 Danh sách workflow BA nên vẽ BPMN trước

1. WF-01 Đăng ký và xác minh email  
2. WF-02 Đăng nhập / refresh / logout  
3. WF-03 Quên mật khẩu / đặt lại mật khẩu  
4. WF-08 Mở lớp  
5. WF-09 Ghi danh học viên  
6. WF-11 Cấu hình lịch tuần  
7. WF-12 Tạo preview xếp lịch  
8. WF-13 Commit preview thành lesson  
9. WF-19 Upload và duyệt tài liệu  

## 16.5 Danh sách giả định mở phải xác nhận

1. Program có lifecycle publish/archive/status hay không.
2. Ai được quyền quản lý học viên.
3. Ai được quyền duyệt tài liệu.
4. ClassSchedule được cấu hình ở đâu và bởi ai.
5. Enrollment có cần approval workflow không.
6. Attendance status gồm những giá trị nào.
7. LeaveRequest có còn trong phase hiện tại không.
8. Consultation/lead management có nằm trong phạm vi hay không.
9. Material moderation có là phạm vi sản phẩm chính hay chỉ demo.
10. Scheduling benchmark cần metric nào để chọn solver.

## 16.6 Bản tóm tắt 1 trang cho BA

EduCenter hiện là một hệ thống quản lý trung tâm dạy thêm đã có nền vững cho các miền vận hành cốt lõi: tài khoản, giáo viên, học viên, khóa học, chương trình, lớp học, phòng học, ca học, xếp lịch và sinh lesson. Điểm mạnh nhất của codebase là **trục vận hành đào tạo** và **module scheduling**: hệ thống đã có cấu trúc dữ liệu đủ tốt cho class, shift, room, class_schedule, lesson; có preview scheduling, commit lesson và đã tách kiến trúc solver service để phục vụ nhiều thuật toán.

Tuy vậy, hệ thống vẫn đang ở trạng thái **nửa vận hành, nửa nền dữ liệu**. Nhiều miền nghiệp vụ quan trọng cho vận hành học tập thực tế đã có entity nhưng chưa có API và workflow thật, gồm attendance, lesson summary, academic record, leave request và consultation. Vì vậy, khi dựng sơ đồ BA, cần phân biệt rõ:

- **Miền đã implemented**: auth, teacher, student, course, program, class, enrollment, room, shift, scheduling, materials.
- **Miền có data model nhưng chưa vận hành**: attendance, lesson summary, academic record, leave request, consultation, objective/outcome.
- **Miền mới ở mức kế hoạch**: predictive analytics AT_RISK.

Điểm cần ưu tiên xác nhận với PO/Tech Lead là **phân quyền** và **lifecycle**. Hiện student CRUD chỉ cần auth, material review chỉ cần auth, `/auth/me` đang lỗi auth/ID, program/room đang có mâu thuẫn contract, và enrollment chưa phản ánh đầy đủ vòng đời nghiệp vụ. Đây là các điểm có thể làm BA, kiến trúc và báo cáo sản phẩm lệch nhau nếu không chốt sớm.

Nếu BA cần vẽ nhanh artefact, nên bắt đầu từ:

1. **Use case diagram** theo 5 cụm: Auth, Danh mục đào tạo, Vận hành lớp, Scheduling, Tài liệu.
2. **BPMN** cho 9 workflow ưu tiên đã nêu.
3. **ERD lõi** cho User/Teacher/Student/Course/Program/Class/Enrollment/Room/Shift/Lesson.
4. **Permission matrix** vì code hiện đang có nhiều authorization gap cần business quyết định lại.

Tài liệu này nên được dùng như gói reverse-engineer hiện trạng. Khi viết báo cáo nghiệp vụ chính thức, cần đánh dấu rõ phần nào là **current-state thực sự từ code**, phần nào là **future-state mong muốn**.
