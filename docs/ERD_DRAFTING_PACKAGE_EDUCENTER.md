# GÓI ERD DRAFTING PACKAGE VÀ BUSINESS DATA DICTIONARY HỆ THỐNG EDUCENTER

**Phiên bản:** refined from previous reverse-engineering baseline  
**Nguồn sự thật:** codebase EduCenter tại `/Users/hant/golang/doan`  
**Mục đích sử dụng:** hỗ trợ BA/SA vẽ logical ERD, gợi ý physical ERD ban đầu, lifecycle/state diagram và phụ lục data dictionary cho báo cáo

## Quy ước mức độ chắc chắn

- **Confirmed from code**: xác nhận trực tiếp từ entity/model, migration, DTO, controller, use case, repository, validation hoặc UI hiện hữu.
- **Strongly inferred from code**: suy luận mạnh từ tên trường, quan hệ, side effect trong use case, route hoặc màn hình.
- **Assumption / Needs BA validation**: giả định nghiệp vụ hợp lý nhưng code chưa đủ bằng chứng hoặc có mâu thuẫn giữa các lớp triển khai.

---

# 1. Phạm vi mô hình dữ liệu

## 1.1 Mô hình dữ liệu lõi của EduCenter là gì

EduCenter đang có mô hình dữ liệu xoay quanh một trung tâm dạy thêm vận hành theo 5 trục lớn:

1. **Quản trị truy cập và danh tính**
   - tài khoản hệ thống
   - xác minh OTP
   - reset mật khẩu
2. **Danh mục đào tạo và nguồn lực**
   - học viên
   - giáo viên
   - khóa học
   - chương trình đào tạo
   - phòng học
   - ca học
3. **Tổ chức lớp và vận hành học vụ**
   - lớp học
   - lịch tuần lớp
   - ghi danh
   - lesson sinh từ xếp lịch
4. **Theo dõi học tập**
   - điểm danh
   - tổng kết buổi học
   - academic record
   - đơn xin phép
5. **Nghiệp vụ hỗ trợ và kiểm soát**
   - upload tài liệu và kiểm duyệt
   - lead tư vấn
   - preview scheduling và benchmark solver

**Đánh giá BA:** mô hình dữ liệu hiện tại đủ tốt để mô tả một **core tutoring center management system**, nhưng phần học thuật nâng cao và phê duyệt nghiệp vụ đang ở trạng thái **mô hình dữ liệu đi trước use case vận hành**.

**Mức độ chắc chắn:** Confirmed from code

## 1.2 Nhóm dữ liệu master

| Nhóm | Thực thể | Vai trò nghiệp vụ |
|---|---|---|
| Danh tính và tài khoản | `User` | định danh người dùng và quyền truy cập |
| Người học và giảng dạy | `Student`, `Teacher` | master dữ liệu vận hành học vụ |
| Danh mục đào tạo | `Course`, `Program`, `Objective`, `Outcome` | khung đào tạo và cấu trúc học thuật |
| Tài nguyên vận hành | `Room`, `Shift` | đầu vào để mở lớp và xếp lịch |
| Bảng mã kiểm duyệt | `Label` | chuẩn hóa kết quả moderation |

**Mức độ chắc chắn:** Confirmed from code

## 1.3 Nhóm dữ liệu giao dịch

| Nhóm | Thực thể | Vai trò nghiệp vụ |
|---|---|---|
| An toàn tài khoản | `UserOTP`, `PasswordReset` | xác minh và khôi phục truy cập |
| Tổ chức lớp | `Class`, `ClassSchedule`, `Enrollment` | mở lớp, lịch tuần, roster |
| Vận hành thời khóa biểu | `Lesson`, `Preview Scheduling` | từ preview đến lesson thực |
| Học tập | `Attendance`, `LessonSummary`, `AcademicRecord`, `LeaveRequest` | theo dõi chuyên cần, nội dung dạy và kết quả |
| Nội dung số | `Material`, `AuditLog`, `ApprovalDecision` | upload, audit và duyệt tài liệu |
| Tuyển sinh / đầu mối | `Consultation` | tiếp nhận nhu cầu học |

**Mức độ chắc chắn:** Confirmed from code

## 1.4 Nhóm dữ liệu mapping

| Nhóm | Thực thể | Ý nghĩa nghiệp vụ |
|---|---|---|
| Chương trình - khóa học | `ProgramCourse` | một chương trình gồm nhiều khóa học và một khóa học có thể thuộc nhiều chương trình |
| Học viên - lớp | `Enrollment` | quan hệ tham gia lớp có lifecycle riêng |

**Mức độ chắc chắn:** Confirmed from code

## 1.5 Nhóm dữ liệu lịch sử / audit

| Nhóm | Thực thể | Ý nghĩa nghiệp vụ |
|---|---|---|
| Xác minh | `UserOTP`, `PasswordReset` | lịch sử token có hạn, có used/expired |
| Kiểm duyệt nội dung | `AuditLog`, `ApprovalDecision` | giải thích và truy vết kết quả duyệt tài liệu |
| Xóa mềm | `DeletedAt` trên một số entity | giữ lịch sử logic thay vì xóa cứng |

**Mức độ chắc chắn:** Confirmed from code

## 1.6 Nhóm dữ liệu hỗ trợ kỹ thuật

| Thành phần | Vai trò kỹ thuật |
|---|---|
| `Preview Scheduling` | object logic lưu in-memory, chưa có bảng DB |
| `CreatedAt`, `UpdatedAt`, `DeletedAt` | timestamp và audit kỹ thuật |
| storage key của material | lưu relative path thay vì absolute path |
| hash/token trong OTP, reset password | bảo mật ứng dụng |

**Mức độ chắc chắn:** Confirmed from code

## 1.7 Cách đọc mô hình dữ liệu khi vẽ ERD

- **ERD lõi nghiệp vụ** nên tập trung vào `User`, `Student`, `Teacher`, `Course`, `Program`, `Class`, `Shift`, `Room`, `Enrollment`, `Lesson`.
- **ERD vận hành học vụ** nên tách riêng `Attendance`, `LessonSummary`, `AcademicRecord`, `LeaveRequest`.
- **ERD hỗ trợ / kỹ thuật** nên tách riêng `UserOTP`, `PasswordReset`, `Material`, `AuditLog`, `ApprovalDecision`, `Label`, `Consultation`.
- Các thực thể như `Preview Scheduling`, `Benchmark result` hiện **không phải entity DB**; chỉ nên vẽ như object logic ngoài physical ERD.

**Mức độ chắc chắn:** Strongly inferred from code

---

# 2. Danh mục thực thể chi tiết

## ENT-001. User

- **Tên bảng / model:** `users` / `User`
- **Định nghĩa nghiệp vụ:** Tài khoản đăng nhập của người dùng hệ thống.
- **Vai trò trong hệ thống:** Cổng truy cập, phân vai actor, nền cho authentication/authorization.
- **Loại thực thể:** Master
- **Khóa chính:** `id` (UUID)
- **Khóa ngoại:** Không có FK trực tiếp trong model; bị tham chiếu bởi nhiều thực thể khác theo hướng nghiệp vụ.
- **Người tạo dữ liệu:** người dùng tự đăng ký; seed dữ liệu; có thể admin nội bộ trong tương lai.
- **Use case tạo:** đăng ký tài khoản.
- **Use case cập nhật:** verify OTP, reset password, change password, có thể cập nhật hồ sơ trong tương lai.
- **Use case hoàn tất / vô hiệu:** chưa có luồng deactivate rõ ràng.
- **Soft delete:** Có (`DeletedAt`)
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã user | UUID | Có | `uuid_generate_v4()` | PK |  | technical PK |
| `code` | mã nội bộ user | varchar(50) | Không rõ |  | Có |  | ít được dùng ở API hiện tại |
| `full_name` | tên hiển thị | varchar(255) | Có trong register |  | Không | register yêu cầu min length | business display name |
| `email` | email đăng nhập | varchar(255) | Có |  | Có | email format | khóa đăng nhập thực tế |
| `password` | mật khẩu băm | text | Có |  | Không | không trả ra API | technical secure field |
| `role` | vai trò hệ thống | varchar(50) | Không | `STUDENT` | Không | ADMIN/TEACHER/STUDENT... | chưa có bảng role riêng |
| `is_active` | tài khoản đã kích hoạt chưa | bool | Không | `true` ở entity | Không | register thực tế set `false` | có mâu thuẫn default vs use case |
| `created_at` | thời điểm tạo | timestamp | Không | `now()` | Không |  | audit |
| `updated_at` | thời điểm sửa | timestamp | Không |  | Không |  | audit |
| `deleted_at` | xóa mềm | timestamp | Không |  | Index |  | soft delete |

- **Trường trạng thái / vòng đời:** `is_active`, `role`
- **Trường audit:** `created_at`, `updated_at`, `deleted_at`
- **Rủi ro / điểm chưa rõ:**
  - `is_active` default ở entity là `true`, nhưng register luôn set `false`.
  - Chưa có mô hình role master hoặc permission table.
  - API `/auth/me` hiện có dấu hiệu sai kiểu ID giữa context và use case.

## ENT-002. UserOTP

- **Tên bảng / model:** `user_otps` / `UserOTP`
- **Định nghĩa nghiệp vụ:** OTP dùng cho xác minh người dùng, chủ yếu ở luồng kích hoạt tài khoản.
- **Vai trò:** kiểm soát an toàn xác minh tài khoản và có thể tái sử dụng cho flow khác.
- **Loại thực thể:** Transaction / security
- **Khóa chính:** `id`
- **Khóa ngoại:** `user_id -> users.id`
- **Người tạo dữ liệu:** hệ thống
- **Use case tạo:** đăng ký tài khoản
- **Use case cập nhật:** verify OTP
- **Soft delete:** Có
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã record OTP | UUID | Có | uuid | PK |  |  |
| `user_id` | user sở hữu OTP | UUID | Có |  | Index | FK logic tới `users` |  |
| `otp_hash` | hash OTP | text | Có |  | Không | không lưu plain OTP | security |
| `expired_at` | thời điểm hết hạn | timestamp | Có |  | Không | phải lớn hơn thời điểm tạo |  |
| `used_at` | thời điểm đã dùng | timestamp nullable | Không |  | Không | null = chưa dùng | lifecycle marker |
| `created_at` | thời điểm tạo | timestamp | Không | now() | Không |  |  |
| `deleted_at` | xóa mềm | timestamp nullable | Không |  | Index |  |  |

- **Lifecycle logic:** active nếu `used_at is null` và `expired_at > now`
- **Rủi ro / điểm chưa rõ:** chưa thấy rate limit resend OTP và chưa thấy rule lock sau nhiều lần nhập sai.

## ENT-003. PasswordReset

- **Tên bảng / model:** `password_resets` / `PasswordReset`
- **Định nghĩa nghiệp vụ:** yêu cầu đặt lại mật khẩu thông qua token có hạn.
- **Vai trò:** khôi phục truy cập cho người dùng quên mật khẩu.
- **Loại thực thể:** Transaction / security
- **Khóa chính:** `id`
- **Khóa ngoại:** `user_id -> users.id`
- **Người tạo dữ liệu:** hệ thống
- **Use case tạo:** quên mật khẩu
- **Use case cập nhật:** đặt lại mật khẩu
- **Soft delete:** Không thấy trong entity
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã request reset | UUID | Có | uuid | PK |  |  |
| `user_id` | user cần reset | UUID | Có |  | Index | FK logic |  |
| `token_hash` | hash token reset | text | Có |  | Không rõ | chỉ so sánh hash |  |
| `expires_at` | hết hạn token | timestamp | Có |  | Không | phải kiểm tra trước reset |  |
| `used_at` | thời điểm token đã dùng | timestamp nullable | Không |  | Không | null = chưa dùng | lifecycle marker |
| `requested_ip` | IP yêu cầu | string | Không |  | Không |  | security/audit |
| `user_agent` | user agent | string | Không |  | Không |  | security/audit |
| `created_at` | thời điểm tạo | timestamp | Không | now() | Không |  |  |

- **Rủi ro / điểm chưa rõ:** chưa có chính sách revoke token cũ khi tạo token mới hoặc chống spam.

## ENT-004. Student

- **Tên bảng / model:** `students` / `Student`
- **Định nghĩa nghiệp vụ:** hồ sơ học viên của trung tâm.
- **Vai trò:** nguồn dữ liệu cho ghi danh, điểm danh, kết quả học tập, đơn xin phép.
- **Loại thực thể:** Master
- **Khóa chính:** `id`
- **Khóa ngoại:** bị tham chiếu bởi `Enrollment`, `Attendance`, `AcademicRecord`, `LeaveRequest`
- **Người tạo dữ liệu:** admin/vận hành
- **Use case tạo:** tạo học viên
- **Use case cập nhật:** cập nhật học viên
- **Use case kết thúc:** xóa mềm học viên
- **Soft delete:** Có
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã học viên | UUID | Có | uuid | PK |  |  |
| `code` | mã học viên | varchar(50) | Không rõ |  | Có |  | mã vận hành |
| `full_name` | họ tên học viên | varchar(255) | Có ở create/update |  | Không | required ở request |  |
| `email` | email liên hệ | varchar(255) | Không |  | Không rõ | email nếu có | không dùng login hiện tại |
| `phone` | số điện thoại học viên | varchar(20) | Không |  | Không |  |  |
| `guardian_phone` | số điện thoại phụ huynh | varchar(20) | Không |  | Không |  | dấu hiệu có actor phụ huynh |
| `grade_level` | khối lớp | varchar(50) | Không |  | Không |  | quan trọng cho course/class |
| `status` | trạng thái học viên | varchar(50) | Không | `ACTIVE` | Không | enum chưa chuẩn hóa | thiếu bảng mã trạng thái |
| `date_of_birth` | ngày sinh | timestamp nullable | Không |  | Không |  |  |
| `gender` | giới tính | varchar(20) | Không |  | Không | enum chưa chuẩn hóa |  |
| `address` | địa chỉ | text | Không |  | Không |  |  |
| `created_at` | tạo lúc nào | timestamp | Không | now() | Không |  |  |
| `updated_at` | sửa lúc nào | timestamp | Không |  | Không |  |  |
| `deleted_at` | xóa mềm | timestamp | Không |  | Index |  |  |

- **Trường trạng thái / vòng đời:** `status`
- **Rủi ro / điểm chưa rõ:**
  - student CRUD hiện chỉ yêu cầu auth, chưa khóa `ADMIN`.
  - chưa có rule trùng lặp dữ liệu học viên rõ ràng.
  - chưa có liên kết trực tiếp `Student -> User` cho portal học viên.

## ENT-005. Teacher

- **Tên bảng / model:** `teachers` / `Teacher`
- **Định nghĩa nghiệp vụ:** hồ sơ giáo viên/nhân sự giảng dạy.
- **Vai trò:** nguồn lực dạy học, gắn với lớp, lesson, scheduling và material.
- **Loại thực thể:** Master
- **Khóa chính:** `id`
- **Khóa ngoại:** bị tham chiếu bởi `Class`, `Lesson`, `Material`
- **Người tạo dữ liệu:** admin
- **Use case tạo:** tạo giáo viên
- **Use case cập nhật:** cập nhật giáo viên
- **Use case kết thúc:** xóa mềm giáo viên
- **Soft delete:** Có
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã giáo viên | UUID | Có | uuid | PK |  |  |
| `code` | mã giáo viên | varchar(50) | Không |  | Có | checked unique ở create |  |
| `full_name` | họ tên giáo viên | varchar(255) | Có |  | Không | required |  |
| `email` | email giáo viên | varchar(255) | Không |  | Có nếu nhập | checked unique ở create |  |
| `phone` | số điện thoại | varchar(20) | Không |  | Không |  |  |
| `is_school_teacher` | có phải giáo viên trường không | bool | Không | `false` | Không | true/false | phân loại nguồn tuyển GV |
| `school_name` | trường đang công tác | varchar(255) | Không |  | Không |  |  |
| `employment_type` | hình thức làm việc | varchar(50) | Không | `PART_TIME` | Không | PART_TIME/FULL_TIME | UI dùng cả 2 giá trị |
| `status` | trạng thái giáo viên | varchar(50) | Không | `ACTIVE` | Không | ACTIVE/INACTIVE | dùng trong UI |
| `notes` | ghi chú | text | Không |  | Không |  |  |
| `created_at`/`updated_at`/`deleted_at` | audit | timestamp | Không |  |  |  |  |

- **Trường trạng thái / vòng đời:** `status`, `employment_type`
- **Rủi ro / điểm chưa rõ:**
  - teacher không tự động sinh `User` account.
  - timetable/statistics teacher đang có API đọc nhưng không thấy ownership control cho từng teacher portal.

## ENT-006. Course

- **Tên bảng / model:** `courses` / `Course`
- **Định nghĩa nghiệp vụ:** đơn vị đào tạo cơ sở dùng để mở lớp.
- **Vai trò:** nguồn dữ liệu chính cho số buổi, thời lượng buổi, học phí và phân loại môn học.
- **Loại thực thể:** Master
- **Khóa chính:** `id`
- **Khóa ngoại:** bị tham chiếu bởi `Class`; liên kết N-N với `Program`
- **Người tạo dữ liệu:** admin
- **Use case tạo:** tạo khóa học
- **Use case cập nhật:** cập nhật khóa học
- **Use case kết thúc:** xóa mềm khóa học
- **Soft delete:** Có
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã khóa học | UUID | Có | uuid | PK |  |  |
| `code` | mã khóa học | varchar(50) | Có |  | Có | not null |  |
| `name` | tên khóa học | varchar(255) | Có ở use case |  | Không | required |  |
| `description` | mô tả khóa học | text | Không |  | Không |  |  |
| `grade_level` | khối áp dụng | varchar(50) | Không |  | Không |  |  |
| `subject` | môn học | varchar(255) | Không |  | Không |  |  |
| `session_count` | số buổi cần dạy | int | Không rõ ở DTO nhưng quan trọng nghiệp vụ |  | Không | nên > 0 | solver dùng trực tiếp |
| `session_duration_minutes` | số phút mỗi buổi | int | Không rõ ở DTO nhưng quan trọng nghiệp vụ |  | Không | nên > 0 | solver dùng trực tiếp |
| `total_hours` | tổng giờ khóa học | numeric(8,2) | Không |  | Không |  | chưa thấy rule sync với session_count * duration |
| `price` | học phí khóa học | numeric(10,2) | Không |  | Không |  | có thể mâu thuẫn với `Class.price` |
| `status` | trạng thái khóa học | varchar(50) | Không | `ACTIVE` | Không | ACTIVE/... | enum chưa khóa hoàn toàn |
| `created_at`/`updated_at`/`deleted_at` | audit | timestamp | Không |  |  |  |  |

- **Trường trạng thái / vòng đời:** `status`
- **Rủi ro / điểm chưa rõ:**
  - chưa có validation chặt cho `session_count`, `session_duration_minutes`, `total_hours`.
  - cần BA xác nhận giá chuẩn ở `Course` hay ở `Class`.

## ENT-007. Program

- **Tên bảng / model:** `programs` / `Program`
- **Định nghĩa nghiệp vụ:** khung chương trình đào tạo gồm nhiều khóa học theo track.
- **Vai trò:** chuẩn hóa lộ trình đào tạo và giúp mở lớp theo chương trình.
- **Loại thực thể:** Master
- **Khóa chính:** `id`
- **Khóa ngoại:** `created_by_id`, `approved_by_id` (logical FK); N-N với `Course`
- **Người tạo dữ liệu:** admin/quản lý đào tạo
- **Use case tạo:** tạo chương trình
- **Use case cập nhật:** cập nhật chương trình, thêm/gỡ khóa học
- **Use case publish/archive:** chưa có API rõ
- **Soft delete:** Có
- **Mức độ chắc chắn:** Confirmed from code, nhưng lifecycle có phần Strongly inferred

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã chương trình | UUID | Có | uuid | PK |  |  |
| `code` | mã chương trình | varchar(50) | Có |  | Có | not null |  |
| `name` | tên chương trình | varchar(255) | Có ở create/update |  | Không | required |  |
| `track` | loại chương trình | varchar(50) | Không |  | Không | SUPPORT/BASIC/ADVANCED | UI dùng rõ nhưng backend chưa khóa enum |
| `effective_from` | hiệu lực từ | timestamp nullable | Không |  | Không |  |  |
| `effective_to` | hiệu lực đến | timestamp nullable | Không |  | Không |  |  |
| `created_by_id` | người tạo | UUID nullable | Không |  | Không | logical FK tới `users` |  |
| `approved_by_id` | người phê duyệt | UUID nullable | Không |  | Không | logical FK |  |
| `approval_note` | ghi chú phê duyệt | text | Không |  | Không |  |  |
| `published_at` | thời điểm xuất bản | timestamp nullable | Không |  | Không |  | lifecycle marker |
| `archived_at` | thời điểm lưu trữ | timestamp nullable | Không |  | Không |  | lifecycle marker |
| `created_at`/`updated_at`/`deleted_at` | audit | timestamp | Không |  |  |  |  |

- **Trường trạng thái / vòng đời:** `published_at`, `archived_at`
- **Rủi ro / điểm chưa rõ:**
  - migration cũ có `status`, entity hiện tại không có.
  - publish/archive là dấu vết dữ liệu, chưa có use case đầy đủ.
  - chưa có rule “program active” để chặn mở lớp hay không.

## ENT-008. ProgramCourse

- **Tên bảng / model:** `program_courses` / `ProgramCourse`
- **Định nghĩa nghiệp vụ:** bảng liên kết N-N giữa chương trình và khóa học.
- **Vai trò:** xây dựng cấu trúc chương trình từ các khóa học thành phần.
- **Loại thực thể:** Mapping
- **Khóa chính:** `id`
- **Khóa ngoại:** `program_id`, `course_id`
- **Người tạo dữ liệu:** admin
- **Use case tạo:** gán khóa học vào chương trình
- **Use case xóa:** gỡ khóa học khỏi chương trình
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã mapping | UUID | Có | uuid | PK |  | khác với cách dùng composite key thuần |
| `program_id` | chương trình cha | UUID | Có |  | Không | FK logic |  |
| `course_id` | khóa học thành phần | UUID | Có |  | Không | FK logic |  |

- **Rủi ro / điểm chưa rõ:**
  - chưa có unique constraint rõ ở entity cho cặp `program_id + course_id`.
  - chưa có thứ tự môn học trong chương trình.

## ENT-009. Objective

- **Tên bảng / model:** `objectives` / `Objective`
- **Định nghĩa nghiệp vụ:** mục tiêu đào tạo cấp chương trình.
- **Vai trò:** mô hình hóa cấu trúc học thuật của program.
- **Loại thực thể:** Master / academic structure
- **Khóa chính:** `id`
- **Khóa ngoại:** `program_id -> programs.id`
- **Người tạo dữ liệu:** chưa có use case
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code về schema; use case là Assumption / Needs BA validation

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã objective | UUID | Có | uuid | PK |  |  |
| `code` | mã mục tiêu | varchar(50) | Có |  | Có | not null |  |
| `name` | tên mục tiêu | text | Có |  | Không | not null |  |
| `program_id` | chương trình cha | UUID | Có |  | Không | OnDelete CASCADE |  |

- **Rủi ro / điểm chưa rõ:** có entity nhưng chưa có API, UI, use case.

## ENT-010. Outcome

- **Tên bảng / model:** `outcomes` / `Outcome`
- **Định nghĩa nghiệp vụ:** chuẩn đầu ra của chương trình hoặc của objective.
- **Vai trò:** mô hình hóa thành phần học thuật đầu ra.
- **Loại thực thể:** Master / academic structure
- **Khóa chính:** `id`
- **Khóa ngoại:** `program_id`, `objective_id`
- **Người tạo dữ liệu:** chưa có use case
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code về schema; use case là Assumption / Needs BA validation

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã outcome | UUID | Có | uuid | PK |  |  |
| `code` | mã outcome | varchar(50) | Có |  | Có | not null |  |
| `name` | tên outcome | text | Có |  | Không | not null |  |
| `program_id` | chương trình cha | UUID | Có |  | Không | OnDelete CASCADE |  |
| `objective_id` | objective cha | UUID nullable | Không |  | Không | OnDelete SET NULL | optional nesting |

- **Rủi ro / điểm chưa rõ:** chưa có workflow quản lý thực thể này.

## ENT-011. Room

- **Tên bảng / model:** `rooms` / `Room`
- **Định nghĩa nghiệp vụ:** phòng học vật lý dùng để tổ chức lớp và lesson.
- **Vai trò:** nguồn lực sức chứa và vị trí.
- **Loại thực thể:** Master / resource
- **Khóa chính:** `id`
- **Khóa ngoại:** bị tham chiếu bởi `Class`, `ClassSchedule`, `Lesson`
- **Người tạo dữ liệu:** admin
- **Use case tạo:** tạo phòng
- **Use case cập nhật:** cập nhật phòng
- **Use case xóa:** xóa phòng
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code, nhưng contract API có mâu thuẫn

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã phòng | UUID | Có | uuid | PK |  |  |
| `code` | mã phòng | varchar(50) | Có |  | Có | not null |  |
| `name` | tên phòng | varchar(255) | Có |  | Không | not null |  |
| `capacity` | sức chứa tối đa | int | Không rõ ở entity, nhưng DTO validate |  | Không | min = 1 ở request | dùng cho scheduling/enrollment |
| `address` | địa chỉ / vị trí | text | Không |  | Không |  | UI gọi là location/address |
| `created_at`/`updated_at` | audit | timestamp | Không |  |  |  |  |

- **Rủi ro / điểm chưa rõ:**
  - DTO và UI từng tham chiếu `location`, `status`, nhưng entity chỉ có `address`.
  - chưa có status hoạt động của phòng.

## ENT-012. Shift

- **Tên bảng / model:** `shifts` / `Shift`
- **Định nghĩa nghiệp vụ:** ca học chuẩn để lớp và scheduling sử dụng.
- **Vai trò:** chuẩn hóa khung giờ, thay cho cách cấu hình tự do từng lớp.
- **Loại thực thể:** Configuration / master
- **Khóa chính:** `id`
- **Khóa ngoại:** bị tham chiếu bởi `ClassSchedule`
- **Người tạo dữ liệu:** admin
- **Use case tạo:** tạo ca học
- **Use case cập nhật:** cập nhật ca học
- **Use case xóa:** xóa ca học
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã shift | UUID | Có | uuid | PK |  |  |
| `code` | mã ca học | varchar(50) | Có |  | Có | uniqueIndex, not null |  |
| `name` | tên ca học | varchar(255) | Có |  | Không | not null |  |
| `start_time` | giờ bắt đầu | varchar(10) | Có |  | Không | required | đang lưu dạng string |
| `end_time` | giờ kết thúc | varchar(10) | Có |  | Không | required | đang lưu dạng string |
| `duration_minutes` | thời lượng ca | int | Có |  | Không | > 0 |  |
| `session_type` | loại ca | varchar(50) | Có |  | Không | MORNING/AFTERNOON/EVENING/CUSTOM | locked in DTO |
| `is_active` | ca có được phép dùng không | bool | Không | `true` | Không | true/false | scheduling chỉ dùng active shifts |
| `notes` | ghi chú | text | Không |  | Không |  |  |
| `created_at`/`updated_at` | audit | timestamp | Không |  |  |  |  |

- **Trường trạng thái / vòng đời:** `is_active`, `session_type`
- **Rủi ro / điểm chưa rõ:** chưa thấy chặn xóa shift đang được lớp sử dụng.

## ENT-013. Class

- **Tên bảng / model:** `classes` / `Class`
- **Định nghĩa nghiệp vụ:** lớp học thực tế đang vận hành tại trung tâm.
- **Vai trò:** trung tâm của flow học vụ, ghi danh và xếp lịch.
- **Loại thực thể:** Transaction / operational master
- **Khóa chính:** `id`
- **Khóa ngoại:** `program_id`, `course_id`, `teacher_id`, `room_id`
- **Người tạo dữ liệu:** admin
- **Use case tạo:** tạo lớp học
- **Use case cập nhật:** cập nhật lớp, phân công giáo viên
- **Use case kết thúc:** đóng lớp/hủy lớp/xóa mềm
- **Soft delete:** Có
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã lớp | UUID | Có | uuid | PK |  |  |
| `code` | mã lớp | varchar(50) | Có |  | Có | not null |  |
| `name` | tên lớp | varchar(255) | Có |  | Không | not null |  |
| `notes` | ghi chú lớp | text | Không |  | Không |  |  |
| `start_date` | ngày bắt đầu | timestamp | Có |  | Không | required |  |
| `end_date` | ngày kết thúc | timestamp nullable | Không |  | Không |  |  |
| `max_students` | sĩ số tối đa | int | Có ở DTO |  | Không | min = 1 |  |
| `status` | trạng thái lớp | varchar(50) | Không | `OPEN` | Không | OPEN/CLOSED/CANCELLED | scheduling chỉ lấy OPEN |
| `price` | giá lớp / giá bán thực tế | numeric(10,2) | Không |  | Không |  | mâu thuẫn tiềm ẩn với `Course.price` |
| `program_id` | chương trình áp dụng | UUID nullable | Không |  | Không |  | optional |
| `course_id` | khóa học chính | UUID nullable | Không |  | Không |  | scheduling phụ thuộc mạnh |
| `teacher_id` | giáo viên phụ trách | UUID nullable | Không |  | Không |  | scheduling phụ thuộc mạnh |
| `room_id` | phòng mặc định | UUID nullable | Không |  | Không |  | DTO create/update chưa expose rõ |
| `created_at`/`updated_at`/`deleted_at` | audit | timestamp | Không |  |  |  |  |

- **Trường trạng thái / vòng đời:** `status`
- **Rủi ro / điểm chưa rõ:**
  - `room_id` có ở entity nhưng luồng tạo/cập nhật lớp chưa expose rõ.
  - `price` ở lớp chưa rõ là override hay copy từ khóa học.
  - chưa có policy đóng lớp khi còn lesson chưa hoàn tất.

## ENT-014. ClassSchedule

- **Tên bảng / model:** `class_schedules` / `ClassSchedule`
- **Định nghĩa nghiệp vụ:** lịch tuần chuẩn của lớp theo thứ và ca học.
- **Vai trò:** nguồn ràng buộc cứng cho scheduling.
- **Loại thực thể:** Configuration / operational
- **Khóa chính:** `id`
- **Khóa ngoại:** `class_id`, `shift_id`, `room_id` optional
- **Người tạo dữ liệu:** dự kiến admin/điều phối đào tạo
- **Use case tạo/cập nhật:** hiện chưa có API quản trị rõ
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code về schema; current-state use case là Strongly inferred

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã lịch tuần | UUID | Có | uuid | PK |  |  |
| `class_id` | lớp áp dụng | UUID | Có |  | Không | OnDelete CASCADE |  |
| `day_of_week` | thứ trong tuần | varchar(20) | Có |  | Không | cần chốt domain MONDAY/TUESDAY hay 2/3/4 | entity hiện dùng string |
| `shift_id` | ca học áp dụng | UUID | Có |  | Không | FK tới `Shift` |  |
| `room_id` | phòng cố định của slot | UUID nullable | Không |  | Không | FK tới `Room` | optional fixed room |

- **Rủi ro / điểm chưa rõ:**
  - rất quan trọng cho scheduling nhưng thiếu API CRUD.
  - cần BA chốt domain chính thức của `day_of_week`.

## ENT-015. Enrollment

- **Tên bảng / model:** `enrollments` / `Enrollment`
- **Định nghĩa nghiệp vụ:** ghi nhận học viên thuộc lớp nào.
- **Vai trò:** roster lớp, nền cho attendance và academic record.
- **Loại thực thể:** Mapping có lifecycle nghiệp vụ
- **Khóa chính:** `id`
- **Khóa ngoại:** `class_id`, `student_id`
- **Người tạo dữ liệu:** admin/vận hành
- **Use case tạo:** ghi danh học viên
- **Use case xóa:** rút học viên khỏi lớp
- **Use case approve/reject:** có dấu vết ở dữ liệu nhưng chưa có flow thật
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã enrollment | UUID | Có | uuid | PK |  |  |
| `class_id` | lớp | UUID | Có |  | Không | OnDelete CASCADE |  |
| `student_id` | học viên | UUID | Có |  | Không | OnDelete CASCADE |  |
| `status` | trạng thái ghi danh | varchar(50) | Không | `APPLIED` | Không | APPLIED/ENROLLED/... | use case hiện tạo `ENROLLED` trực tiếp |
| `approved_at` | thời điểm duyệt | timestamp nullable | Không |  | Không |  | dấu vết approval |
| `rejected_at` | thời điểm từ chối | timestamp nullable | Không |  | Không |  | dấu vết rejection |
| `created_at`/`updated_at` | audit | timestamp | Không |  |  |  |  |

- **Trường trạng thái / vòng đời:** `status`, `approved_at`, `rejected_at`
- **Rủi ro / điểm chưa rõ:**
  - lifecycle dữ liệu không nhất quán với use case.
  - chưa có chống duplicate enrollment và kiểm sĩ số thực theo roster hiện tại.

## ENT-016. Lesson

- **Tên bảng / model:** `lessons` / `Lesson`
- **Định nghĩa nghiệp vụ:** buổi học thực tế đã được lên lịch.
- **Vai trò:** trục giao dịch chính sau scheduling.
- **Loại thực thể:** Transaction / operational
- **Khóa chính:** `id`
- **Khóa ngoại:** `class_id`, `teacher_id`, `room_id`
- **Người tạo dữ liệu:** scheduling commit
- **Use case tạo:** xác nhận preview để tạo lesson
- **Use case cập nhật:** chưa có lesson management riêng
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã lesson | UUID | Có | uuid | PK |  |  |
| `class_id` | lớp của lesson | UUID | Có |  | Không | OnDelete CASCADE |  |
| `teacher_id` | giáo viên dạy | UUID nullable | Không |  | Không | FK optional | snapshot nguồn lực |
| `date_start` | thời điểm bắt đầu | timestamp | Có |  | Không | not null |  |
| `date_end` | thời điểm kết thúc | timestamp | Có |  | Không | not null |  |
| `room_id` | phòng dạy | UUID nullable | Không |  | Không | FK optional | snapshot tài nguyên |
| `notes` | ghi chú buổi học | text | Không |  | Không |  | preview commit ghi note solver |
| `created_at`/`updated_at` | audit | timestamp | Không |  |  |  |  |

- **Rủi ro / điểm chưa rõ:**
  - chưa có API quản lý danh sách/chi tiết lesson thật.
  - lesson hiện thiên về dữ liệu được sinh tự động hơn là vận hành thủ công.

## ENT-017. Attendance

- **Tên bảng / model:** `attendances` / `Attendance`
- **Định nghĩa nghiệp vụ:** điểm danh học viên theo từng lesson.
- **Vai trò:** dữ liệu chuyên cần, nền cho leave request và predictive analytics.
- **Loại thực thể:** Transaction
- **Khóa chính:** `id`
- **Khóa ngoại:** `lesson_id`, `student_id`
- **Người tạo dữ liệu:** dự kiến giáo viên/admin
- **Use case:** chưa có flow hiện thực
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code về schema; business flow là Strongly inferred

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã attendance | UUID | Có | uuid | PK |  |  |
| `lesson_id` | lesson được chấm | UUID | Có |  | Không | OnDelete CASCADE |  |
| `student_id` | học viên được chấm | UUID | Có |  | Không | OnDelete CASCADE |  |
| `status` | trạng thái điểm danh | int | Có |  | Không | chưa có bảng mã | gap lớn nhất của entity này |
| `note` | ghi chú điểm danh | text | Không |  | Không |  |  |
| `marked_at` | thời điểm chấm | timestamp | Không |  | Không |  |  |
| `created_at`/`updated_at` | audit | timestamp | Không |  |  |  |  |

- **Rủi ro / điểm chưa rõ:** thiếu enum trạng thái nên chưa thể chốt data dictionary nghiệp vụ.

## ENT-018. LessonSummary

- **Tên bảng / model:** `lesson_summaries` / `LessonSummary`
- **Định nghĩa nghiệp vụ:** tổng kết nội dung dạy của một lesson.
- **Vai trò:** cầu nối từ lesson sang đánh giá học viên và giao bài tập.
- **Loại thực thể:** Transaction
- **Khóa chính:** `id`
- **Khóa ngoại:** `lesson_id`, `created_by_id`
- **Người tạo dữ liệu:** dự kiến giáo viên
- **Use case:** chưa có API hiện thực
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code về schema; flow là Strongly inferred

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã summary | UUID | Có | uuid | PK |  |  |
| `lesson_id` | lesson được tổng kết | UUID | Có |  | Unique | OnDelete CASCADE | 1 lesson - 1 summary |
| `topic` | chủ đề buổi học | text | Không |  | Không |  |  |
| `lesson_content` | nội dung đã dạy | text | Không |  | Không |  |  |
| `class_feedback` | phản hồi chung về lớp | text | Không |  | Không |  |  |
| `homework` | bài tập giao | text | Không |  | Không |  |  |
| `homework_deadline` | hạn nộp bài | timestamp | Không |  | Không |  | business due date |
| `teacher_notes` | ghi chú chuyên môn | text | Không |  | Không |  |  |
| `created_by_id` | người lập summary | UUID nullable | Không |  | Không | logical FK tới user |  |
| `created_at`/`updated_at` | audit | timestamp | Không |  |  |  |  |

- **Rủi ro / điểm chưa rõ:** chưa có workflow xác nhận summary là “đã chốt”.

## ENT-019. AcademicRecord

- **Tên bảng / model:** `academic_records` / `AcademicRecord`
- **Định nghĩa nghiệp vụ:** đánh giá kết quả học tập của từng học viên sau lesson summary.
- **Vai trò:** dữ liệu theo dõi tiến bộ học tập và nền cho classification `AT_RISK`.
- **Loại thực thể:** Transaction
- **Khóa chính:** `id`
- **Khóa ngoại:** `lesson_summary_id`, `student_id`
- **Người tạo dữ liệu:** dự kiến giáo viên
- **Use case:** chưa có API hiện thực
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code về schema; flow là Strongly inferred

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã record | UUID | Có | uuid | PK |  |  |
| `lesson_summary_id` | summary cha | UUID | Có |  | Không | OnDelete CASCADE |  |
| `student_id` | học viên được đánh giá | UUID | Có |  | Không | OnDelete CASCADE |  |
| `homework_completed` | đã làm bài tập chưa | bool | Không | `false` | Không | true/false |  |
| `homework_score` | điểm bài tập | numeric(5,2) | Không |  | Không |  |  |
| `attitude_rating` | điểm/level thái độ học tập | int | Không |  | Không | scale chưa rõ |  |
| `participation_score` | điểm tham gia | numeric(5,2) | Không |  | Không |  |  |
| `personal_comment` | nhận xét cá nhân | text | Không |  | Không |  |  |
| `total_score` | tổng điểm | numeric(5,2) | Không |  | Không | công thức chưa rõ |  |
| `is_completed` | đã chốt record chưa | bool | Không | `false` | Không | true/false | lifecycle nhẹ |
| `created_at`/`updated_at` | audit | timestamp | Không |  |  |  | `created-at` json tag hiện có lỗi gõ |

- **Rủi ro / điểm chưa rõ:**
  - chưa có công thức tính `total_score`.
  - chưa có rule khi nào `is_completed = true`.

## ENT-020. LeaveRequest

- **Tên bảng / model:** `leave_requests` / `LeaveRequest`
- **Định nghĩa nghiệp vụ:** đơn xin nghỉ học, đi muộn hoặc về sớm.
- **Vai trò:** xử lý ngoại lệ tham gia học tập.
- **Loại thực thể:** Transaction
- **Khóa chính:** `id`
- **Khóa ngoại:** `student_id`, `class_id`, `lesson_id`, `approved_by_id`
- **Người tạo dữ liệu:** dự kiến học viên hoặc phụ huynh
- **Use case duyệt:** dự kiến admin/giáo vụ
- **Soft delete:** Không thấy `DeletedAt`
- **Mức độ chắc chắn:** Confirmed from code về schema; flow là Strongly inferred

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã đơn | UUID | Có | uuid | PK |  |  |
| `student_id` | người gửi đơn | UUID | Có |  | Không | OnDelete CASCADE |  |
| `leave_type` | loại đơn | varchar(50) | Có |  | Không | LEAVE/LATE/EARLY |  |
| `apply_date` | ngày áp dụng | timestamp | Có |  | Không | not null |  |
| `late_minutes` | số phút đi muộn | int | Không |  | Không | dùng cho `LATE` |  |
| `early_minutes` | số phút về sớm | int | Không |  | Không | dùng cho `EARLY` |  |
| `reason` | lý do | text | Có |  | Không | not null |  |
| `documents` | tài liệu đính kèm | text[] | Không |  | Không |  | dùng `pq.StringArray` |
| `class_id` | lớp liên quan | UUID nullable | Không |  | Không | OnDelete SET NULL |  |
| `lesson_id` | lesson liên quan | UUID nullable | Không |  | Không | OnDelete SET NULL |  |
| `subject` | tiêu đề đơn | varchar(255) | Không |  | Không |  |  |
| `status` | trạng thái đơn | varchar(50) | Không | `PENDING` | Không | PENDING/APPROVED/REJECTED? | enum đầy đủ chưa thấy |
| `approved_by_id` | người duyệt | UUID nullable | Không |  | Không | logical FK tới user |  |
| `approved_at` | duyệt lúc nào | timestamp nullable | Không |  | Không |  |  |
| `rejection_reason` | lý do từ chối | text | Không |  | Không |  |  |
| `created_at`/`updated_at` | audit | timestamp | Không |  |  |  |  |

- **Trường trạng thái / vòng đời:** `status`, `approved_by_id`, `approved_at`
- **Rủi ro / điểm chưa rõ:** chưa có API, chưa nối tự động vào `Attendance`.

## ENT-021. Material

- **Tên bảng / model:** `materials` / `Material`
- **Định nghĩa nghiệp vụ:** tài liệu giảng dạy do giáo viên tải lên.
- **Vai trò:** đối tượng chính của quy trình audit và phê duyệt nội dung.
- **Loại thực thể:** Transaction / content
- **Khóa chính:** `id`
- **Khóa ngoại:** `teacher_id`, `latest_label_id`
- **Người tạo dữ liệu:** giáo viên
- **Use case tạo:** upload tài liệu
- **Use case cập nhật:** audit tự động, review thủ công
- **Soft delete:** Có
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã tài liệu | UUID | Có | uuid | PK |  |  |
| `teacher_id` | giáo viên sở hữu | UUID | Có |  | Không | FK logic tới `Teacher` | ownership đang lấy từ body |
| `title` | tiêu đề | varchar(255) | Có |  | Không | not null |  |
| `description` | mô tả | text | Không |  | Không |  |  |
| `file_name` | tên file gốc | varchar(255) | Có |  | Không | not null |  |
| `file_path` | storage key/relative path | text | Có |  | Không | not null | không trả raw path ra UI nữa |
| `file_type` | loại file | varchar(100) | Có |  | Không | pdf/doc/docx/png/jpg/jpeg | validated in use case |
| `file_size` | kích thước file | int64 | Không | `0` | Không | <= 10MB | validated in use case |
| `status` | trạng thái xử lý | varchar(50) | Có | `UPLOADED` | Không | UPLOADED/SCANNING/AI_REVIEWED/APPROVED/REJECTED | lifecycle hiện rõ |
| `latest_label_id` | nhãn AI mới nhất | UUID nullable | Không |  | Không | FK optional tới `Label` |  |
| `uploaded_at` | thời điểm upload | timestamp | Không | now() | Không |  |  |
| `created_at`/`updated_at`/`deleted_at` | audit | timestamp | Không |  |  |  |  |

- **Trường trạng thái / vòng đời:** `status`, `latest_label_id`
- **Rủi ro / điểm chưa rõ:**
  - ownership và reviewer hiện chưa derive từ JWT.
  - cần xác nhận module này còn trong scope sản phẩm chính hay là scope phụ.

## ENT-022. AuditLog

- **Tên bảng / model:** `audit_logs` / `AuditLog`
- **Định nghĩa nghiệp vụ:** lịch sử OCR/AI audit cho material.
- **Vai trò:** giải thích vì sao tài liệu bị gắn nhãn nào, mức tự tin ra sao.
- **Loại thực thể:** History / audit
- **Khóa chính:** `id`
- **Khóa ngoại:** `material_id`, `label_id`
- **Người tạo dữ liệu:** hệ thống OCR/AI
- **Use case tạo:** upload material
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã audit log | UUID | Có | uuid | PK |  |  |
| `material_id` | material bị audit | UUID | Có |  | Không | OnDelete CASCADE |  |
| `label_id` | nhãn kết quả | UUID nullable | Không |  | Không | optional FK |  |
| `status` | trạng thái audit | varchar(50) | Có | `COMPLETED` | Không | COMPLETED/... | hiện stub flow gần như luôn complete |
| `provider` | công cụ/nguồn audit | varchar(100) | Có |  | Không | ví dụ `STUB_OCR_GEMINI` |  |
| `raw_ocr_text` | text OCR | text | Không |  | Không |  |  |
| `confidence_score` | độ tin cậy | numeric(5,4) | Không |  | Không |  |  |
| `reasoning` | giải thích AI | text | Không |  | Không |  |  |
| `detected_issues` | danh sách vấn đề | jsonb/string | Không | `[]` | Không |  | entity đang map string |
| `triggered_at` | bắt đầu audit | timestamp | Không | now() | Không |  |  |
| `completed_at` | hoàn thành audit | timestamp nullable | Không |  | Không |  |  |
| `created_at`/`updated_at` | audit | timestamp | Không |  |  |  |  |

- **Rủi ro / điểm chưa rõ:** audit hiện là stub, chưa phải provider thật.

## ENT-023. ApprovalDecision

- **Tên bảng / model:** `approval_decisions` / `ApprovalDecision`
- **Định nghĩa nghiệp vụ:** quyết định phê duyệt thủ công với material.
- **Vai trò:** chốt kết quả cuối cùng sau AI review.
- **Loại thực thể:** History / approval
- **Khóa chính:** `id`
- **Khóa ngoại:** `material_id`, `audit_log_id`, `compliance_officer_id`
- **Người tạo dữ liệu:** người duyệt
- **Use case tạo:** review material
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã quyết định | UUID | Có | uuid | PK |  |  |
| `material_id` | material bị duyệt | UUID | Có |  | Không | OnDelete CASCADE |  |
| `audit_log_id` | log audit tham chiếu | UUID nullable | Không |  | Không | optional FK |  |
| `compliance_officer_id` | người duyệt | UUID | Có |  | Không | logical FK | hiện gửi từ client |
| `approved` | kết quả duyệt | bool | Có |  | Không | true/false |  |
| `reject_reason` | lý do từ chối | text | Không |  | Không |  |  |
| `notes` | ghi chú duyệt | text | Không |  | Không |  |  |
| `decided_at` | thời điểm quyết định | timestamp | Không | now() | Không |  |  |
| `created_at`/`updated_at` | audit | timestamp | Không |  |  |  |  |

- **Rủi ro / điểm chưa rõ:** chưa có role check đủ chặt cho reviewer.

## ENT-024. Label

- **Tên bảng / model:** `labels` / `Label`
- **Định nghĩa nghiệp vụ:** bảng mã phân loại mức độ an toàn tài liệu.
- **Vai trò:** chuẩn hóa output audit.
- **Loại thực thể:** Master / reference
- **Khóa chính:** `id`
- **Khóa ngoại:** bị tham chiếu bởi `Material`, `AuditLog`
- **Người tạo dữ liệu:** seed dữ liệu
- **Use case cập nhật:** chưa có
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã label | UUID | Có | uuid | PK |  |  |
| `code` | mã nhãn | varchar(50) | Có |  | Có | SAFE/WARNING/DANGER |  |
| `name` | tên hiển thị | varchar(100) | Có |  | Không | not null |  |
| `severity` | mức độ rủi ro | varchar(50) | Có |  | Không | SAFE/WARNING/DANGER | có thể trùng code |
| `description` | mô tả | text | Không |  | Không |  |  |
| `created_at`/`updated_at` | audit | timestamp | Không |  |  |  |  |

- **Rủi ro / điểm chưa rõ:** bảng reference này nên được quản trị read-only.

## ENT-025. Consultation

- **Tên bảng / model:** `consultations` / `Consultation`
- **Định nghĩa nghiệp vụ:** lead/yêu cầu tư vấn từ phụ huynh hoặc khách hàng tiềm năng.
- **Vai trò:** mở rộng sang tuyển sinh và chăm sóc lead.
- **Loại thực thể:** Transaction
- **Khóa chính:** `id`
- **Khóa ngoại:** Không thấy
- **Người tạo dữ liệu:** dự kiến khách hàng hoặc nhân viên tư vấn
- **Use case:** chưa có API hiện thực
- **Soft delete:** Không thấy
- **Mức độ chắc chắn:** Confirmed from code về schema; current-state flow là Assumption / Needs BA validation

| Trường | Ý nghĩa nghiệp vụ | Kiểu suy ra | Bắt buộc | Default | Unique | Domain/Validation | Ghi chú |
|---|---|---|---|---|---|---|---|
| `id` | mã lead | UUID | Có | uuid | PK |  |  |
| `full_name` | tên người quan tâm | varchar(255) | Có |  | Không | not null |  |
| `phone` | số điện thoại lead | varchar(20) | Có |  | Không | not null |  |
| `grade_level` | khối lớp quan tâm | varchar(50) | Có |  | Không | not null |  |
| `notes` | nhu cầu/ghi chú tư vấn | text | Không |  | Không |  |  |
| `status` | trạng thái lead | varchar(50) | Không | `PENDING` | Không | PENDING/... | enum đầy đủ chưa rõ |
| `created_at`/`updated_at` | audit | timestamp | Không |  |  |  |  |

- **Rủi ro / điểm chưa rõ:** entity đã có nhưng chưa có use case, chưa biết có còn trong roadmap chính không.

## ENT-026. Preview Scheduling (thực thể logic)

- **Tên bảng / model:** không có bảng DB; object in-memory
- **Định nghĩa nghiệp vụ:** kết quả preview xếp lịch trước khi tạo lesson thật.
- **Vai trò:** lớp đệm kiểm tra khả thi, xung đột và quyết định commit.
- **Loại thực thể:** Technical / transient operational object
- **Khóa logic:** `run_id`
- **Người tạo dữ liệu:** scheduling engine
- **Use case tạo:** tạo preview xếp lịch
- **Use case kết thúc:** commit preview
- **Mức độ chắc chắn:** Confirmed from code

| Trường | Ý nghĩa nghiệp vụ |
|---|---|
| `run_id` | mã preview run |
| `status` | FAILED / PARTIAL / COMPLETED |
| `generated_at` | thời điểm sinh preview |
| `filters` | đầu vào preview |
| `assignments` | các lesson candidate |
| `conflicts` | xung đột chưa xếp được |
| `summary` | số lớp, số buổi, số conflict, score |

- **Rủi ro / điểm chưa rõ:** mất khi restart backend vì chưa persist DB.

---

# 3. Phân tích quan hệ dữ liệu

| Mã quan hệ | Thực thể nguồn | Thực thể đích | Ý nghĩa nghiệp vụ | Cardinality | Optionality | FK / bằng chứng | Cascade / delete | Hệ quả nghiệp vụ khi xóa / cập nhật | Mức độ chắc chắn | Điểm chưa rõ |
|---|---|---|---|---|---|---|---|---|---|---|
| REL-001 | `User` | `UserOTP` | một user có thể sinh nhiều OTP theo thời gian | 1-N | OTP bắt buộc thuộc user | `user_id` | chưa thấy explicit FK ở migration runtime | xóa user có thể làm mất nghĩa audit OTP | Confirmed from code | physical FK cần xác nhận |
| REL-002 | `User` | `PasswordReset` | một user có thể có nhiều token reset | 1-N | reset bắt buộc thuộc user | `user_id` | chưa thấy rõ | mất history reset khi xóa user | Confirmed from code |  |
| REL-003 | `Program` | `Course` qua `ProgramCourse` | chương trình gồm nhiều khóa học, khóa học có thể tái sử dụng ở nhiều chương trình | N-N | optional cả hai phía | many2many + mapping | nên xóa mapping khi xóa 1 đầu | thay đổi program ảnh hưởng cấu trúc đào tạo | Confirmed from code | thiếu thứ tự course |
| REL-004 | `Program` | `Objective` | chương trình có nhiều mục tiêu | 1-N | objective bắt buộc thuộc program | `program_id` | OnDelete CASCADE | xóa program sẽ mất mục tiêu | Confirmed from code | chưa có use case |
| REL-005 | `Program` | `Outcome` | chương trình có nhiều chuẩn đầu ra | 1-N | outcome bắt buộc thuộc program | `program_id` | OnDelete CASCADE | xóa program sẽ mất outcome | Confirmed from code | chưa có use case |
| REL-006 | `Objective` | `Outcome` | một objective có thể nhóm nhiều outcome | 1-N | outcome có thể không gắn objective | `objective_id` | OnDelete SET NULL | outcome vẫn tồn tại nếu objective bị xóa | Confirmed from code |  |
| REL-007 | `Class` | `Program` | lớp có thể mở theo một chương trình | N-1 | optional | `program_id` | chưa rõ | nếu program bị thay đổi có thể làm lệch ngữ nghĩa lớp | Confirmed from code | rule đổi program giữa chừng chưa rõ |
| REL-008 | `Class` | `Course` | lớp gắn một khóa học chính | N-1 | optional | `course_id` | chưa rõ | nếu thiếu course, scheduling fail | Confirmed from code |  |
| REL-009 | `Class` | `Teacher` | lớp có một giáo viên phụ trách chính | N-1 | optional | `teacher_id` | chưa rõ | nếu thiếu teacher, scheduling conflict | Confirmed from code | chưa hỗ trợ nhiều giáo viên |
| REL-010 | `Class` | `Room` | lớp có thể có phòng mặc định | N-1 | optional | `room_id` | chưa rõ | room ảnh hưởng capacity và lesson default | Confirmed from code | create/update class chưa expose rõ |
| REL-011 | `Class` | `ClassSchedule` | lớp có nhiều slot lịch tuần | 1-N | optional trong schema nhưng gần như bắt buộc cho scheduling tốt | `class_id` | OnDelete CASCADE | xóa class sẽ mất lịch tuần | Confirmed from code | thiếu API quản trị |
| REL-012 | `ClassSchedule` | `Shift` | slot lịch tuần dùng một ca học | N-1 | bắt buộc | `shift_id` | chưa thấy delete guard | xóa shift có thể làm lịch lớp invalid | Confirmed from code | thiếu guard ở use case |
| REL-013 | `ClassSchedule` | `Room` | slot lịch tuần có thể khóa về một phòng | N-1 | optional | `room_id` | chưa rõ | phòng mất khả dụng sẽ tạo conflict scheduling | Confirmed from code |  |
| REL-014 | `Enrollment` | `Student` | một enrollment thuộc một học viên | N-1 | bắt buộc | `student_id` | OnDelete CASCADE | xóa student làm mất roster/history | Confirmed from code | giữ history hay không cần BA chốt |
| REL-015 | `Enrollment` | `Class` | một enrollment thuộc một lớp | N-1 | bắt buộc | `class_id` | OnDelete CASCADE | xóa class làm mất roster/history | Confirmed from code |  |
| REL-016 | `Lesson` | `Class` | lesson là phiên học của một lớp | N-1 | bắt buộc | `class_id` | OnDelete CASCADE | xóa class sẽ xóa lessons | Confirmed from code |  |
| REL-017 | `Lesson` | `Teacher` | lesson có teacher thực thi buổi dạy | N-1 | optional | `teacher_id` | chưa rõ | phục vụ timetable/statistics | Confirmed from code |  |
| REL-018 | `Lesson` | `Room` | lesson diễn ra tại một phòng | N-1 | optional | `room_id` | chưa rõ | phục vụ kiểm tra trùng phòng | Confirmed from code |  |
| REL-019 | `Attendance` | `Lesson` | điểm danh phát sinh theo lesson | N-1 | bắt buộc | `lesson_id` | OnDelete CASCADE | xóa lesson sẽ mất attendance | Confirmed from code |  |
| REL-020 | `Attendance` | `Student` | mỗi bản ghi attendance của một học viên | N-1 | bắt buộc | `student_id` | OnDelete CASCADE | xóa student sẽ mất attendance history | Confirmed from code |  |
| REL-021 | `LessonSummary` | `Lesson` | một lesson có tối đa một summary | 1-1 | summary bắt buộc thuộc lesson | `lesson_id` unique | OnDelete CASCADE | xóa lesson sẽ mất summary | Confirmed from code |  |
| REL-022 | `LessonSummary` | `User` | summary có người lập | N-1 | optional | `created_by_id` | chưa rõ | cần audit trách nhiệm | Confirmed from code | user nào được lập cần chốt |
| REL-023 | `AcademicRecord` | `LessonSummary` | một summary có nhiều đánh giá học viên | 1-N | bắt buộc | `lesson_summary_id` | OnDelete CASCADE | xóa summary sẽ mất record | Confirmed from code |  |
| REL-024 | `AcademicRecord` | `Student` | một học viên có nhiều academic record | 1-N | bắt buộc | `student_id` | OnDelete CASCADE | xóa student sẽ mất learning history | Confirmed from code |  |
| REL-025 | `LeaveRequest` | `Student` | đơn xin phép thuộc một học viên | N-1 | bắt buộc | `student_id` | OnDelete CASCADE | xóa student sẽ mất đơn lịch sử | Confirmed from code |  |
| REL-026 | `LeaveRequest` | `Class` | đơn có thể gắn lớp | N-1 | optional | `class_id` | OnDelete SET NULL | còn đơn nhưng mất lớp tham chiếu | Confirmed from code |  |
| REL-027 | `LeaveRequest` | `Lesson` | đơn có thể gắn lesson cụ thể | N-1 | optional | `lesson_id` | OnDelete SET NULL | đơn còn nhưng mất lesson tham chiếu | Confirmed from code |  |
| REL-028 | `LeaveRequest` | `User` | đơn có thể ghi approver | N-1 | optional | `approved_by_id` | chưa rõ | cần audit duyệt đơn | Confirmed from code | chưa có API duyệt |
| REL-029 | `Material` | `Teacher` | giáo viên có nhiều tài liệu | 1-N | bắt buộc | `teacher_id` | chưa rõ | nếu teacher bị xóa cần giữ tài liệu hay không | Confirmed from code | ownership chưa derive từ JWT |
| REL-030 | `Material` | `Label` | material giữ nhãn AI mới nhất | N-1 | optional | `latest_label_id` | optional/set null logic | UI hiển thị nhãn cuối cùng | Confirmed from code |  |
| REL-031 | `Material` | `AuditLog` | material có nhiều log audit | 1-N | optional | `material_id` | OnDelete CASCADE | xóa material sẽ mất trace audit | Confirmed from code |  |
| REL-032 | `Material` | `ApprovalDecision` | material có thể có nhiều quyết định review theo lịch sử | 1-N | optional | `material_id` | OnDelete CASCADE | xóa material sẽ mất history duyệt | Confirmed from code | current UI dùng lần review gần nhất |
| REL-033 | `AuditLog` | `Label` | mỗi audit log có thể gắn một nhãn | N-1 | optional | `label_id` | optional | giải thích vì sao bị gắn cờ | Confirmed from code |  |
| REL-034 | `ApprovalDecision` | `AuditLog` | quyết định duyệt có thể tham chiếu log audit cụ thể | N-1 | optional | `audit_log_id` | optional | tăng khả năng truy xuất quyết định | Confirmed from code | chưa chắc UI có dùng |
| REL-035 | `ApprovalDecision` | `User` | một user đóng vai reviewer | N-1 | bắt buộc về schema logic | `compliance_officer_id` | chưa rõ | audit người quyết định | Confirmed from code | role reviewer chưa khóa chặt |
| REL-036 | `Preview Scheduling` | `Class` / `Teacher` / `Room` / `Shift` | object preview dùng snapshot đầu vào từ master và operational data | logical multi-ref | optional theo filter | use case scheduling | không có DB cascade | preview mất khi restart | Confirmed from code | không vẽ vào physical ERD |

---

# 4. Phân tích lifecycle / state transition

## 4.1 User

| Nội dung | Phân tích |
|---|---|
| Field lifecycle | `is_active`, gián tiếp là `role` |
| Trạng thái chắc chắn | `is_active = false`, `is_active = true` |
| Ý nghĩa | chưa kích hoạt / đã kích hoạt |
| Chuyển trạng thái hợp lệ | `false -> true` qua verify OTP |
| Tác nhân | hệ thống trong use case verify OTP |
| Điều kiện | OTP đúng, chưa hết hạn, chưa dùng |
| Side effects | cho phép login thành công |
| Điểm thiếu | chưa có suspend/lock account, chưa có deactivate by admin |
| Mức độ chắc chắn | Confirmed from code |

## 4.2 Program

| Nội dung | Phân tích |
|---|---|
| Field lifecycle | `published_at`, `archived_at`, `approved_by_id`, `approval_note` |
| Trạng thái suy luận | Draft, Published, Archived |
| Ý nghĩa | chương trình nháp / đã công bố / ngừng sử dụng |
| Chuyển trạng thái hợp lệ | Draft -> Published -> Archived |
| Tác nhân | admin hoặc quản lý đào tạo |
| Điều kiện | chưa có use case rõ |
| Side effects | ảnh hưởng khả năng mở lớp theo program |
| Điểm thiếu | migration từng có `status` nhưng entity không còn; thiếu API publish/archive |
| Mức độ chắc chắn | Strongly inferred from code |

## 4.3 Class

| Nội dung | Phân tích |
|---|---|
| Field lifecycle | `status` |
| Trạng thái chắc chắn | `OPEN`, `CLOSED`, `CANCELLED` |
| Ý nghĩa | đang vận hành / đã đóng / bị hủy |
| Chuyển trạng thái hợp lệ | `OPEN -> CLOSED`, `OPEN -> CANCELLED`; chiều ngược lại chưa rõ |
| Tác nhân | admin |
| Điều kiện | update class |
| Side effects | scheduling preview chỉ load `OPEN` |
| Điểm thiếu | chưa rõ rule khi class đã có lesson thì có được cancel hay không |
| Mức độ chắc chắn | Confirmed from code |

## 4.4 Enrollment

| Nội dung | Phân tích |
|---|---|
| Field lifecycle | `status`, `approved_at`, `rejected_at` |
| Trạng thái chắc chắn | `APPLIED` default ở entity, `ENROLLED` ở use case hiện tại |
| Trạng thái suy luận | `REJECTED` |
| Ý nghĩa | chờ xử lý / đã vào lớp / bị từ chối |
| Chuyển trạng thái mong muốn | `APPLIED -> ENROLLED`, `APPLIED -> REJECTED` |
| Thực tế code | tạo trực tiếp `ENROLLED` |
| Tác nhân | admin/vận hành |
| Side effects | tăng roster lớp |
| Điểm thiếu | thiếu approval flow, thiếu duplicate check, thiếu tính lại capacity theo roster hiện tại |
| Mức độ chắc chắn | Confirmed from code + Strongly inferred for intended lifecycle |

## 4.5 LeaveRequest

| Nội dung | Phân tích |
|---|---|
| Field lifecycle | `status`, `approved_by_id`, `approved_at`, `rejection_reason` |
| Trạng thái chắc chắn | `PENDING` |
| Trạng thái suy luận | `APPROVED`, `REJECTED` |
| Ý nghĩa | chờ duyệt / được duyệt / bị từ chối |
| Chuyển trạng thái hợp lệ | `PENDING -> APPROVED`, `PENDING -> REJECTED` |
| Tác nhân | admin, giáo vụ hoặc người quản lý học vụ |
| Side effects | có thể ảnh hưởng attendance hoặc giải thích absence |
| Điểm thiếu | chưa có API, chưa có enum status đầy đủ |
| Mức độ chắc chắn | Strongly inferred from code |

## 4.6 Material

| Nội dung | Phân tích |
|---|---|
| Field lifecycle | `status`, `latest_label_id` |
| Trạng thái chắc chắn | `UPLOADED`, `SCANNING`, `AI_REVIEWED`, `APPROVED`, `REJECTED` |
| Ý nghĩa | mới tải lên / đang audit / đã có AI kết quả / duyệt / từ chối |
| Chuyển trạng thái hợp lệ | `UPLOADED -> SCANNING -> AI_REVIEWED -> APPROVED/REJECTED` |
| Tác nhân | teacher, OCR/AI service, reviewer |
| Side effects | sinh `AuditLog`, `ApprovalDecision`, cập nhật `latest_label_id` |
| Điểm thiếu | chưa có re-review flow, chưa rõ tài liệu SAFE có cần duyệt tay hay không |
| Mức độ chắc chắn | Confirmed from code |

## 4.7 AcademicRecord

| Nội dung | Phân tích |
|---|---|
| Field lifecycle | `is_completed` |
| Trạng thái chắc chắn | `false`, `true` |
| Ý nghĩa | bản ghi đánh giá nháp / đã chốt |
| Chuyển trạng thái hợp lệ | `false -> true` |
| Tác nhân | dự kiến giáo viên |
| Side effects | có thể trở thành nguồn training cho predictive analytics |
| Điểm thiếu | chưa có use case, chưa có điều kiện chốt record |
| Mức độ chắc chắn | Strongly inferred from code |

## 4.8 Attendance

| Nội dung | Phân tích |
|---|---|
| Field lifecycle | `status` kiểu int |
| Trạng thái chắc chắn | chưa xác định được domain |
| Ý nghĩa | dự kiến hiện diện/vắng/có phép/đi muộn... |
| Chuyển trạng thái hợp lệ | chưa xác định |
| Tác nhân | dự kiến giáo viên/admin |
| Side effects | ảnh hưởng học tập, leave request, predictive analytics |
| Điểm thiếu | thiếu bảng mã status, thiếu API, thiếu ownership rule |
| Mức độ chắc chắn | Confirmed gap |

## 4.9 Preview Scheduling

| Nội dung | Phân tích |
|---|---|
| Field lifecycle | `status` |
| Trạng thái chắc chắn | `FAILED`, `PARTIAL`, `COMPLETED` |
| Ý nghĩa | không xếp được / xếp được một phần / xếp xong và có thể commit |
| Chuyển trạng thái | preview sinh ra là 1 trạng thái cuối, không update tiếp |
| Tác nhân | scheduling engine |
| Side effects | quyết định có cho phép commit hay không |
| Điểm thiếu | không persist DB, mất khi restart |
| Mức độ chắc chắn | Confirmed from code |

---

# 5. Quy tắc dữ liệu và ràng buộc

## 5.1 Unique constraints

| Mã rule | Ràng buộc | Bằng chứng | Mức độ chắc chắn |
|---|---|---|---|
| DR-UNQ-01 | `users.email` unique | entity `User` | Confirmed from code |
| DR-UNQ-02 | `users.code` unique | entity `User` | Confirmed from code |
| DR-UNQ-03 | `students.code` unique | entity `Student` | Confirmed from code |
| DR-UNQ-04 | `teachers.code` unique | entity `Teacher` | Confirmed from code |
| DR-UNQ-05 | `teachers.email` được check unique khi tạo | use case teacher create | Confirmed from code |
| DR-UNQ-06 | `courses.code` unique | entity `Course` | Confirmed from code |
| DR-UNQ-07 | `programs.code` unique | entity `Program` | Confirmed from code |
| DR-UNQ-08 | `objectives.code` unique | entity `Objective` | Confirmed from code |
| DR-UNQ-09 | `outcomes.code` unique | entity `Outcome` | Confirmed from code |
| DR-UNQ-10 | `rooms.code` unique | entity `Room` | Confirmed from code |
| DR-UNQ-11 | `shifts.code` unique | entity `Shift` | Confirmed from code |
| DR-UNQ-12 | `classes.code` unique | entity `Class` | Confirmed from code |
| DR-UNQ-13 | `lesson_summaries.lesson_id` unique | entity `LessonSummary` | Confirmed from code |
| DR-UNQ-14 | `labels.code` unique | entity `Label` | Confirmed from code |
| DR-UNQ-15 | `program_courses(program_id, course_id)` nên unique logic | mapping semantics | Strongly inferred from code |

## 5.2 Required fields

| Mã rule | Trường bắt buộc | Bối cảnh nghiệp vụ | Mức độ chắc chắn |
|---|---|---|---|
| DR-REQ-01 | `User.email`, `User.password`, `User.full_name` | đăng ký tài khoản | Confirmed from code |
| DR-REQ-02 | `Teacher.full_name` | tạo giáo viên | Confirmed from code |
| DR-REQ-03 | `Course.code`, `Course.name` | tạo khóa học | Confirmed from code |
| DR-REQ-04 | `Program.code`, `Program.name` | tạo chương trình | Confirmed from code |
| DR-REQ-05 | `Class.code`, `Class.name`, `Class.start_date`, `Class.max_students` | tạo lớp | Confirmed from code |
| DR-REQ-06 | `Shift.code`, `Shift.name`, `Shift.start_time`, `Shift.end_time`, `Shift.duration_minutes`, `Shift.session_type` | tạo ca học | Confirmed from code |
| DR-REQ-07 | `Material.title`, `file`, `teacher_id` | upload material | Confirmed from code |
| DR-REQ-08 | `Consultation.full_name`, `phone`, `grade_level` | lead intake | Confirmed from code về schema |

## 5.3 Default values

| Mã rule | Field | Default | Ý nghĩa nghiệp vụ | Mức độ chắc chắn |
|---|---|---|---|---|
| DR-DEF-01 | `User.role` | `STUDENT` | mặc định người dùng thường | Confirmed from code |
| DR-DEF-02 | `User.is_active` | `true` ở entity | nhưng register override thành `false` | Confirmed from code |
| DR-DEF-03 | `Student.status` | `ACTIVE` | học viên đang hoạt động | Confirmed from code |
| DR-DEF-04 | `Teacher.is_school_teacher` | `false` | mặc định không phải GV trường | Confirmed from code |
| DR-DEF-05 | `Teacher.employment_type` | `PART_TIME` | mặc định cộng tác | Confirmed from code |
| DR-DEF-06 | `Teacher.status` | `ACTIVE` | có thể nhận lớp | Confirmed from code |
| DR-DEF-07 | `Course.status` | `ACTIVE` | khóa học đang dùng | Confirmed from code |
| DR-DEF-08 | `Shift.is_active` | `true` | ca học được phép dùng | Confirmed from code |
| DR-DEF-09 | `Class.status` | `OPEN` | lớp sẵn sàng vận hành | Confirmed from code |
| DR-DEF-10 | `Enrollment.status` | `APPLIED` | có ý đồ approval flow | Confirmed from code |
| DR-DEF-11 | `AcademicRecord.homework_completed` | `false` | chưa làm bài là mặc định | Confirmed from code |
| DR-DEF-12 | `AcademicRecord.is_completed` | `false` | chưa chốt record | Confirmed from code |
| DR-DEF-13 | `LeaveRequest.status` | `PENDING` | chờ duyệt | Confirmed from code |
| DR-DEF-14 | `Material.status` | `UPLOADED` | vừa nhận tài liệu | Confirmed from code |
| DR-DEF-15 | `AuditLog.status` | `COMPLETED` | stub audit hiện trả ngay | Confirmed from code |
| DR-DEF-16 | `Consultation.status` | `PENDING` | lead mới tiếp nhận | Confirmed from code |

## 5.4 Domain / enum restrictions

| Mã rule | Field | Domain giá trị | Mức độ chắc chắn | Ghi chú |
|---|---|---|---|---|
| DR-DOM-01 | `Teacher.employment_type` | `PART_TIME`, `FULL_TIME` | Confirmed from code |  |
| DR-DOM-02 | `Teacher.status` | `ACTIVE`, `INACTIVE` | Confirmed from code |  |
| DR-DOM-03 | `Shift.session_type` | `MORNING`, `AFTERNOON`, `EVENING`, `CUSTOM` | Confirmed from code | DTO đã validate |
| DR-DOM-04 | `Class.status` | `OPEN`, `CLOSED`, `CANCELLED` | Confirmed from code | DTO đã validate |
| DR-DOM-05 | `LeaveRequest.leave_type` | `LEAVE`, `LATE`, `EARLY` | Confirmed from code | comment trong entity |
| DR-DOM-06 | `Material.status` | `UPLOADED`, `SCANNING`, `AI_REVIEWED`, `APPROVED`, `REJECTED` | Confirmed from code | flow hiện có |
| DR-DOM-07 | `Label.code`, `Label.severity` | `SAFE`, `WARNING`, `DANGER` | Confirmed from code |  |
| DR-DOM-08 | `Program.track` | `SUPPORT`, `BASIC`, `ADVANCED` | Strongly inferred from code | UI và sample data dùng rõ |
| DR-DOM-09 | `Attendance.status` | chưa xác định | Confirmed gap | cần BA chốt |
| DR-DOM-10 | `Consultation.status` | chưa thấy enum đầy đủ | Assumption / Needs BA validation |  |

## 5.5 Validation rules

| Mã rule | Rule | Bằng chứng | Tác động nghiệp vụ |
|---|---|---|---|
| DR-VAL-01 | room capacity >= 1 | room DTO | không cho tạo phòng vô nghĩa |
| DR-VAL-02 | shift duration_minutes >= 1 | shift DTO | ca học phải có độ dài dương |
| DR-VAL-03 | material file type chỉ cho pdf/doc/docx/png/jpg/jpeg | material use case | giới hạn định dạng tài liệu |
| DR-VAL-04 | material file size <= 10MB | material use case | kiểm soát lưu trữ và an toàn |
| DR-VAL-05 | register email phải đúng format | auth DTO/use case | dữ liệu đăng nhập hợp lệ |
| DR-VAL-06 | only active user can login | auth use case | chặn login trước verify OTP |
| DR-VAL-07 | preview scheduling phải có date range hợp lệ | scheduling controller/use case | tránh chạy solver với input lỗi |
| DR-VAL-08 | commit preview chỉ khi status = COMPLETED và không conflict hard | scheduling use case | tránh sinh lesson sai |
| DR-VAL-09 | scheduling chỉ dùng class `OPEN` | scheduling use case | tránh xếp lớp không còn hoạt động |
| DR-VAL-10 | scheduling chỉ dùng shift `is_active = true` | scheduling use case | tránh dùng ca đã vô hiệu |

## 5.6 Cascade delete / optionality rules

| Mã rule | Quan hệ | Behavior | Mức độ chắc chắn |
|---|---|---|---|
| DR-CAS-01 | `Class -> Lesson` | OnDelete CASCADE | Confirmed from code |
| DR-CAS-02 | `Lesson -> LessonSummary` | OnDelete CASCADE | Confirmed from code |
| DR-CAS-03 | `LessonSummary -> AcademicRecord` | OnDelete CASCADE | Confirmed from code |
| DR-CAS-04 | `Student -> AcademicRecord` | OnDelete CASCADE | Confirmed from code |
| DR-CAS-05 | `Material -> AuditLog` | OnDelete CASCADE | Confirmed from code |
| DR-CAS-06 | `Material -> ApprovalDecision` | OnDelete CASCADE | Confirmed from code |
| DR-CAS-07 | `Objective -> Outcome` | OnDelete SET NULL | Confirmed from code |
| DR-CAS-08 | `LeaveRequest -> Class/Lesson` | OnDelete SET NULL | Confirmed from code |
| DR-CAS-09 | `Enrollment -> Student/Class` | OnDelete CASCADE | Confirmed from code |

## 5.7 Soft delete behavior

| Mã rule | Entity | Soft delete | Mức độ chắc chắn |
|---|---|---|---|
| DR-SD-01 | `User` | Có | Confirmed from code |
| DR-SD-02 | `UserOTP` | Có | Confirmed from code |
| DR-SD-03 | `Student` | Có | Confirmed from code |
| DR-SD-04 | `Teacher` | Có | Confirmed from code |
| DR-SD-05 | `Course` | Có | Confirmed from code |
| DR-SD-06 | `Program` | Có | Confirmed from code |
| DR-SD-07 | `Class` | Có | Confirmed from code |
| DR-SD-08 | `Material` | Có | Confirmed from code |
| DR-SD-09 | `Room`, `Shift`, `Enrollment`, `Lesson`, `Attendance`, `LessonSummary`, `AcademicRecord`, `LeaveRequest`, `ApprovalDecision`, `Label`, `Consultation` | Không thấy `DeletedAt` | Confirmed from code |

## 5.8 Approval / publish / archive fields

| Mã rule | Cụm field | Ý nghĩa nghiệp vụ | Mức độ chắc chắn |
|---|---|---|---|
| DR-LFC-01 | `Program.approved_by_id`, `approval_note`, `published_at`, `archived_at` | dấu vết quy trình duyệt/xuất bản chương trình | Strongly inferred from code |
| DR-LFC-02 | `Enrollment.approved_at`, `rejected_at` | dấu vết quy trình duyệt ghi danh | Strongly inferred from code |
| DR-LFC-03 | `LeaveRequest.approved_by_id`, `approved_at`, `rejection_reason` | quy trình duyệt đơn | Strongly inferred from code |
| DR-LFC-04 | `ApprovalDecision.*` | quyết định duyệt material | Confirmed from code |

## 5.9 Business contradictions and data gaps

| Mã gap | Mô tả | Tác động nghiệp vụ | Mức độ chắc chắn |
|---|---|---|---|
| DG-01 | `Program` migration cũ có `status`, entity hiện không có | khó chốt lifecycle chính thức của chương trình | Confirmed from code |
| DG-02 | `Room` DTO/use case từng dùng `status/location`, entity dùng `code/address` | dictionary room có nguy cơ lệch giữa backend và UI | Confirmed from code |
| DG-03 | `Class.price` và `Course.price` cùng tồn tại | chưa rõ giá bán thật nằm ở đâu | Confirmed from code |
| DG-04 | `Enrollment.status` default `APPLIED` nhưng use case tạo `ENROLLED` luôn | lifecycle ghi danh không nhất quán | Confirmed from code |
| DG-05 | `Attendance.status` là int nhưng không có domain values | không thể chuẩn hóa lifecycle attendance | Confirmed gap |
| DG-06 | `ClassSchedule` là thực thể cốt lõi của scheduling nhưng chưa có API quản trị | quy trình vận hành dữ liệu bị thiếu | Confirmed gap |
| DG-07 | nhiều entity có schema nhưng chưa có use case thật (`Objective`, `Outcome`, `Attendance`, `LessonSummary`, `AcademicRecord`, `LeaveRequest`, `Consultation`) | mô hình dữ liệu đi trước sản phẩm | Confirmed from code |
| DG-08 | material module có workflow thật nhưng quyền reviewer và ownership chưa chặt | rủi ro authorization | Confirmed from code |
| DG-09 | preview scheduling lưu in-memory | không phù hợp nếu cần audit lịch sử preview | Confirmed from code |
| DG-10 | material tables từng có migration riêng nhưng không nằm trong runtime AutoMigrate | deploy không đồng nhất có thể thiếu bảng | Confirmed from code |

---

# 6. Gợi ý vẽ ERD

## 6.1 ERD lõi nghiệp vụ

**Nên vẽ trước:**

1. `User`
2. `Student`
3. `Teacher`
4. `Course`
5. `Program`
6. `ProgramCourse`
7. `Room`
8. `Shift`
9. `Class`
10. `ClassSchedule`
11. `Enrollment`
12. `Lesson`

**Vì sao nên vẽ trước:**

- đây là phần phản ánh rõ nhất nghiệp vụ cốt lõi của trung tâm;
- bao phủ trọn chuỗi danh mục đào tạo -> mở lớp -> ghi danh -> xếp lịch -> sinh lesson;
- là nhóm thực thể có giá trị cao nhất cho báo cáo và demo sản phẩm.

## 6.2 ERD vận hành học vụ

**Nên nhóm riêng:**

1. `Lesson`
2. `Attendance`
3. `LessonSummary`
4. `AcademicRecord`
5. `LeaveRequest`
6. `Student`
7. `Class`

**Vì sao nên tách riêng:**

- cùng một chuỗi nghiệp vụ hậu lesson;
- nhiều thực thể đã có schema nhưng chưa có API, nên BA dễ đánh dấu `current-state` và `future-state`;
- phù hợp để vẽ thêm BPMN/Activity sau này.

## 6.3 ERD hỗ trợ và kỹ thuật

**Nên tách thành ERD phụ:**

1. `UserOTP`
2. `PasswordReset`
3. `Material`
4. `Label`
5. `AuditLog`
6. `ApprovalDecision`
7. `Consultation`
8. `Preview Scheduling` (logical object, không phải physical table)

**Vì sao nên tách:**

- đây không phải trục lõi của classroom operations;
- có nhiều yếu tố security, audit, approval và technical storage;
- giúp sơ đồ chính gọn và tập trung hơn.

## 6.4 Điểm chặn khi vẽ ERD

| Mã chặn | Vấn đề | Ảnh hưởng | Hành động BA nên làm |
|---|---|---|---|
| ERD-BLK-01 | lifecycle `Program` mâu thuẫn giữa entity và migration | khó chốt state diagram | hỏi PO/Tech Lead về draft-publish-archive chính thức |
| ERD-BLK-02 | contract `Room` không thống nhất | dễ vẽ sai field | chốt tên business field là `address` hay `location`, có `status` hay không |
| ERD-BLK-03 | `Attendance.status` chưa có bảng mã | không thể hoàn thiện data dictionary | chốt enum chuẩn |
| ERD-BLK-04 | `ClassSchedule` chưa có API quản trị | khó mô tả nguồn tạo dữ liệu | hỏi quy trình cấu hình lịch tuần |
| ERD-BLK-05 | `Enrollment` vừa giống mapping vừa có dấu vết approval | khó xác định bản chất entity | chốt có approval flow hay không |
| ERD-BLK-06 | nhiều entity học thuật là future-state | dễ trộn current và target architecture | tách rõ khi vẽ |
| ERD-BLK-07 | `Material` có thể không còn nằm trong scope chính | ảnh hưởng phạm vi ERD tổng | xác nhận đưa vào ERD lõi hay phụ lục |

---

# 7. Data dictionary tóm tắt cho phụ lục báo cáo

| Entity | Ý nghĩa nghiệp vụ | Key fields | Quan hệ chính | Ghi chú |
|---|---|---|---|---|
| `User` | tài khoản hệ thống | `email`, `role`, `is_active` | `UserOTP`, `PasswordReset` | nền auth |
| `UserOTP` | OTP xác minh | `user_id`, `expired_at`, `used_at` | `User` | active account flow |
| `PasswordReset` | yêu cầu reset mật khẩu | `user_id`, `token_hash`, `expires_at`, `used_at` | `User` | forgot/reset flow |
| `Student` | hồ sơ học viên | `code`, `full_name`, `grade_level`, `guardian_phone`, `status` | `Enrollment`, `Attendance`, `AcademicRecord`, `LeaveRequest` | student CRUD có gap phân quyền |
| `Teacher` | hồ sơ giáo viên | `code`, `full_name`, `email`, `employment_type`, `status` | `Class`, `Lesson`, `Material` | chưa tự liên kết `User` |
| `Course` | khóa học cơ sở | `code`, `name`, `session_count`, `session_duration_minutes`, `price` | `ProgramCourse`, `Class` | input chính cho scheduling |
| `Program` | chương trình đào tạo | `code`, `name`, `track`, `effective_from`, `effective_to` | `ProgramCourse`, `Objective`, `Outcome`, `Class` | lifecycle mơ hồ |
| `ProgramCourse` | mapping chương trình-khóa học | `program_id`, `course_id` | `Program`, `Course` | nên có unique cặp |
| `Objective` | mục tiêu đào tạo | `code`, `name`, `program_id` | `Program`, `Outcome` | schema có, use case chưa có |
| `Outcome` | chuẩn đầu ra | `code`, `name`, `program_id`, `objective_id` | `Program`, `Objective` | schema có, use case chưa có |
| `Room` | phòng học | `code`, `name`, `capacity`, `address` | `Class`, `ClassSchedule`, `Lesson` | contract field còn lệch |
| `Shift` | ca học chuẩn | `code`, `name`, `start_time`, `end_time`, `session_type`, `is_active` | `ClassSchedule` | đầu vào chuẩn của scheduling |
| `Class` | lớp học vận hành | `code`, `name`, `status`, `max_students`, `course_id`, `teacher_id` | `Program`, `Course`, `Teacher`, `Room`, `Enrollment`, `ClassSchedule`, `Lesson` | trung tâm của học vụ |
| `ClassSchedule` | lịch tuần lớp | `class_id`, `day_of_week`, `shift_id`, `room_id` | `Class`, `Shift`, `Room` | thiếu API quản trị |
| `Enrollment` | ghi danh học viên | `class_id`, `student_id`, `status` | `Class`, `Student` | lifecycle chưa thống nhất |
| `Lesson` | buổi học thực tế | `class_id`, `teacher_id`, `room_id`, `date_start`, `date_end` | `Class`, `Teacher`, `Room`, `Attendance`, `LessonSummary` | sinh từ scheduling commit |
| `Attendance` | điểm danh | `lesson_id`, `student_id`, `status` | `Lesson`, `Student` | thiếu enum status |
| `LessonSummary` | tổng kết buổi học | `lesson_id`, `topic`, `homework`, `teacher_notes` | `Lesson`, `AcademicRecord` | 1 lesson - 1 summary |
| `AcademicRecord` | kết quả học tập từng học viên | `lesson_summary_id`, `student_id`, `total_score`, `is_completed` | `LessonSummary`, `Student` | nền cho predictive analytics |
| `LeaveRequest` | đơn xin phép | `student_id`, `leave_type`, `apply_date`, `status`, `approved_by_id` | `Student`, `Class`, `Lesson`, `User` | chưa có API |
| `Material` | tài liệu giảng dạy | `teacher_id`, `title`, `file_name`, `file_path`, `status`, `latest_label_id` | `Teacher`, `Label`, `AuditLog`, `ApprovalDecision` | module moderation |
| `AuditLog` | log OCR/AI audit | `material_id`, `provider`, `confidence_score`, `label_id` | `Material`, `Label` | lịch sử audit |
| `ApprovalDecision` | quyết định duyệt | `material_id`, `compliance_officer_id`, `approved`, `decided_at` | `Material`, `AuditLog`, `User` | reviewer chưa khóa role đủ chặt |
| `Label` | nhãn audit | `code`, `severity`, `name` | `Material`, `AuditLog` | SAFE/WARNING/DANGER |
| `Consultation` | lead tư vấn | `full_name`, `phone`, `grade_level`, `status` |  | schema có, use case chưa có |
| `Preview Scheduling` | object preview lịch học | `run_id`, `status`, `assignments`, `conflicts`, `summary` | logical refs tới class/teacher/room/shift | không có bảng DB |

---

## Kết luận sử dụng cho BA/SA

Bản refined này nên được dùng như sau:

1. **Vẽ ERD lõi trước**
   - `User`, `Student`, `Teacher`, `Course`, `Program`, `ProgramCourse`, `Room`, `Shift`, `Class`, `ClassSchedule`, `Enrollment`, `Lesson`
2. **Tách ERD học vụ nâng cao**
   - `Attendance`, `LessonSummary`, `AcademicRecord`, `LeaveRequest`
3. **Tách ERD hỗ trợ / kỹ thuật**
   - `UserOTP`, `PasswordReset`, `Material`, `AuditLog`, `ApprovalDecision`, `Label`, `Consultation`
4. **Đánh dấu rõ current-state vs target-state**
   - current-state: auth, master data, class, enrollment, shift, scheduling preview/commit, material upload/review
   - target-state hoặc partial-state: attendance, lesson summary, academic record, leave request, consultation, objective/outcome

Nếu BA cần một sơ đồ sạch để đưa thẳng vào báo cáo, nên dùng 3 tầng:

- **ERD lõi nghiệp vụ**
- **ERD vận hành học vụ**
- **ERD hỗ trợ và kỹ thuật**

Đó sẽ là cách ít rủi ro nhất để phản ánh đúng codebase hiện tại mà không làm báo cáo bị “vẽ quá” so với sản phẩm thực.
