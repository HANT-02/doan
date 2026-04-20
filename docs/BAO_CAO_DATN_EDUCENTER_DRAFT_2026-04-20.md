# ĐẠI HỌC XÂY DỰNG HÀ NỘI

# KHOA CÔNG NGHỆ THÔNG TIN

**ĐỒ ÁN TỐT NGHIỆP**

## Xây dựng hệ thống quản lý trung tâm dạy thêm EduCenter, trọng tâm là xếp lịch thông minh và dự báo học viên có nguy cơ học kém

**Sinh viên thực hiện:** Nguyễn Thế Hà  
**Mã số sinh viên:** 61165  
**Lớp quản lý:** CS2  
**Giảng viên hướng dẫn:** [Bổ sung GVHD]  

Hà Nội, 20/04/2026

---

# LỜI NÓI ĐẦU

Trong bối cảnh chuyển đổi số giáo dục ngày càng được quan tâm, nhu cầu xây dựng các hệ thống hỗ trợ vận hành cho trung tâm dạy thêm trở nên thiết thực. Một hệ thống không chỉ cần đáp ứng các chức năng quản trị cơ bản như quản lý học viên, giáo viên, lớp học, chương trình đào tạo mà còn cần giải quyết được các bài toán vận hành phức tạp, trong đó nổi bật là bài toán xếp lịch.

Đồ án này tập trung xây dựng hệ thống EduCenter theo hướng vừa phục vụ vận hành thực tế của trung tâm dạy thêm, vừa tạo ra điểm nhấn kỹ thuật cho báo cáo tốt nghiệp thông qua hai hướng chính. Hướng thứ nhất là xây dựng mô-đun xếp lịch thông minh với nhiều thuật toán khác nhau để so sánh, benchmark và lựa chọn solver phù hợp cho hệ thống chính. Hướng thứ hai là chuẩn bị nền tảng cho bài toán dự báo học viên có nguy cơ học kém nhằm phục vụ cảnh báo sớm trong tương lai.

Tài liệu này được biên soạn lại từ mã nguồn, tài liệu phân tích nghiệp vụ, tài liệu benchmark và các artifact mô hình hóa hiện có trong project. Các hình ảnh giao diện, biểu đồ và minh họa chi tiết sẽ được bổ sung sau trong quá trình hoàn thiện báo cáo.

---

# MỤC LỤC

1. Giới thiệu đề tài  
2. Phân tích hệ thống  
3. Thiết kế hệ thống  
4. Xếp lịch thông minh và đánh giá thuật toán  
5. Kết luận và hướng phát triển  
Tài liệu tham khảo  

---

# CHƯƠNG 1. GIỚI THIỆU ĐỀ TÀI

## 1.1 Lý do chọn đề tài

Hoạt động của một trung tâm dạy thêm thường bao gồm nhiều đầu việc liên quan chặt chẽ với nhau như quản lý học viên, quản lý giáo viên, tổ chức lớp học, quản lý chương trình đào tạo, quản lý ca học và phòng học. Khi số lượng lớp và nguồn lực tăng lên, việc vận hành thủ công dễ dẫn tới sai sót, trùng lịch và khó theo dõi trạng thái học tập của học viên.

Trong số các nghiệp vụ trên, bài toán xếp lịch là một bài toán có độ phức tạp cao do phải đồng thời thỏa mãn nhiều ràng buộc như tránh trùng giáo viên, tránh trùng phòng, tôn trọng lịch tuần của lớp, chỉ sử dụng ca học hợp lệ và bảo đảm tính khả thi trong một khoảng thời gian xác định. Đây là nội dung có giá trị cả về mặt ứng dụng lẫn học thuật, phù hợp để trở thành điểm nhấn kỹ thuật của đồ án tốt nghiệp.

Bên cạnh đó, hệ thống giáo dục hiện nay cũng có nhu cầu theo dõi sát hơn tiến độ học tập của học viên. Việc chuẩn bị nền tảng dữ liệu cho các bài toán phân tích và cảnh báo sớm học viên có nguy cơ học kém là hướng phát triển có ý nghĩa thực tiễn, giúp hệ thống không chỉ dừng ở mức quản trị dữ liệu mà còn tiến tới hỗ trợ ra quyết định.

Từ các lý do trên, đề tài "Xây dựng hệ thống quản lý trung tâm dạy thêm EduCenter, trọng tâm là xếp lịch thông minh và dự báo học viên có nguy cơ học kém" được lựa chọn để triển khai.

## 1.2 Mục tiêu đề tài

Mục tiêu của đồ án là xây dựng một hệ thống quản lý trung tâm dạy thêm theo mô hình single-tenant, có khả năng hỗ trợ các nghiệp vụ quản trị cốt lõi và tạo ra chiều sâu kỹ thuật ở hai trục chính là xếp lịch thông minh và dự báo học tập.

Cụ thể, hệ thống hướng tới các mục tiêu sau:

- Hoàn thiện các chức năng quản trị cốt lõi gồm quản lý tài khoản, học viên, giáo viên, khóa học, chương trình đào tạo, phòng học, ca học và lớp học.
- Xây dựng luồng vận hành lớp học đầy đủ từ tạo lớp, ghi danh học viên, cấu hình lịch tuần cho lớp đến tạo preview xếp lịch và commit thành lesson thực tế.
- Chuẩn hóa dữ liệu thời gian bằng thực thể `Shift`, giúp việc xếp lịch dựa trên ca học tiêu chuẩn thay vì xử lý các khung giờ rời rạc.
- Nghiên cứu, cài đặt và benchmark ba thuật toán xếp lịch gồm `Graph Coloring + heuristic`, `CP-SAT` và `Tabu Search`.
- Lựa chọn solver phù hợp nhất để dùng làm mặc định cho API scheduling của hệ thống.
- Xây dựng nền dữ liệu và kiến trúc cho bài toán dự báo học viên có nguy cơ học kém trong các giai đoạn tiếp theo.

## 1.3 Đối tượng và phạm vi đề tài

### 1.3.1 Đối tượng hướng tới

Đối tượng mà hệ thống hướng tới là các trung tâm dạy thêm hoặc trung tâm học thêm có quy mô nhỏ đến trung bình, cần một công cụ thống nhất để quản lý học viên, giáo viên, lớp học và vận hành lịch học.

### 1.3.2 Phạm vi triển khai

Trong phạm vi hiện tại của đồ án, hệ thống tập trung vào:

- Quản lý tài khoản và xác thực người dùng.
- Quản lý học viên, giáo viên, khóa học, chương trình đào tạo.
- Quản lý lớp học, ghi danh, phân công giáo viên, quản lý lịch tuần.
- Quản lý phòng học, ca học.
- Xếp lịch học thông minh, benchmark thuật toán, commit lesson.
- Theo dõi học vụ sau lesson như điểm danh, tổng kết buổi học, kết quả học tập, đơn xin phép.
- Chuẩn bị dữ liệu và kiến trúc cho bài toán dự báo `AT_RISK`.

### 1.3.3 Ngoài phạm vi chính

Một số nhánh tuy còn tồn tại trong codebase nhưng không được xem là phạm vi chính thức của báo cáo ở thời điểm hiện tại, bao gồm:

- AI Audit / kiểm duyệt tài liệu mở rộng.
- DevTools và các nhánh công cụ nội bộ.
- Chatbot hoặc trợ lý ảo.
- Tư vấn tuyển sinh / lead intake ở mức hoàn chỉnh.

## 1.4 Công nghệ sử dụng

Hệ thống được xây dựng trên nền tảng full-stack với các công nghệ chính sau:

| Nhóm | Công nghệ | Mục đích sử dụng |
| --- | --- | --- |
| Backend | Golang | Xây dựng API và các use case nghiệp vụ |
| Web framework | Gin | Tổ chức REST API |
| ORM | GORM | Ánh xạ thực thể và truy cập cơ sở dữ liệu |
| Dependency Injection | Google Wire | Quản lý phụ thuộc theo Clean Architecture |
| Cơ sở dữ liệu | PostgreSQL | Lưu trữ dữ liệu hệ thống |
| Frontend | React + TypeScript | Xây dựng giao diện quản trị và giao diện người dùng |
| UI framework | MUI | Chuẩn hóa thành phần giao diện |
| Quản lý trạng thái | Redux Toolkit / RTK Query | Tương tác dữ liệu giữa frontend và backend |
| Tài liệu API | Swagger | Hỗ trợ kiểm thử và mô tả API |
| Mô hình hóa | PlantUML, draw.io | Thiết kế use case, sequence, class diagram, BPMN, ERD |

## 1.5 Ý nghĩa khoa học và thực tiễn

Về mặt thực tiễn, đề tài hỗ trợ số hóa các nghiệp vụ quan trọng trong vận hành trung tâm dạy thêm. Về mặt học thuật, đề tài tạo ra giá trị nghiên cứu ở việc chuẩn hóa bài toán xếp lịch, tách abstraction solver, benchmark nhiều thuật toán trên cùng bộ ràng buộc và lựa chọn chiến lược phù hợp cho hệ thống thực tế.

---

# CHƯƠNG 2. PHÂN TÍCH HỆ THỐNG

## 2.1 Tổng quan bài toán nghiệp vụ

EduCenter được xây dựng như một hệ thống hỗ trợ vận hành cho trung tâm dạy thêm. Dưới góc nhìn nghiệp vụ, chuỗi hoạt động chính của hệ thống bao gồm:

1. Quản lý nguồn lực đầu vào: học viên, giáo viên, khóa học, chương trình đào tạo, phòng học, ca học.
2. Tạo lớp học thực tế, gắn khóa học hoặc chương trình, gán giáo viên phụ trách và cấu hình lịch tuần.
3. Ghi danh học viên vào lớp.
4. Sử dụng dữ liệu lớp học, phòng học, ca học và giáo viên để tạo preview xếp lịch.
5. Khi preview đạt điều kiện, commit để sinh các buổi học thực tế (`Lesson`).
6. Sau mỗi buổi học, tiếp tục thực hiện các nghiệp vụ như điểm danh, tổng kết buổi học, ghi nhận kết quả học tập và xử lý đơn xin phép.

## 2.2 Các tác nhân chính

Các tác nhân tham gia vào hệ thống gồm:

| Tác nhân | Vai trò |
| --- | --- |
| Quản trị viên | Quản lý dữ liệu lõi, lớp học, ghi danh, scheduling, benchmark |
| Giáo viên | Xem lịch dạy, điểm danh, tổng kết buổi học, ghi nhận kết quả học tập |
| Học viên | Xem thông tin học tập, theo dõi kết quả, gửi đơn xin phép |
| Khách | Truy cập các luồng công khai như đăng ký tài khoản |
| Hệ thống gửi email | Gửi OTP và hỗ trợ xác thực |
| Scheduling Engine | Tạo preview và hỗ trợ commit lesson |

**[Chèn Hình 2.1: Sơ đồ use case tổng quan của hệ thống EduCenter]**

## 2.3 Phân rã chức năng

Hệ thống được chia thành các nhóm chức năng chính sau:

### 2.3.1 Xác thực và quản lý tài khoản

- Đăng ký tài khoản.
- Xác minh OTP.
- Đăng nhập, đăng xuất, làm mới token.
- Quên mật khẩu, đặt lại mật khẩu, đổi mật khẩu.
- Xem hồ sơ người dùng.

### 2.3.2 Quản lý danh mục đào tạo và nguồn lực

- Quản lý học viên.
- Quản lý giáo viên.
- Quản lý khóa học.
- Quản lý chương trình đào tạo.
- Quản lý phòng học.
- Quản lý ca học `Shift`.

### 2.3.3 Quản lý lớp học và ghi danh

- Tạo, cập nhật, xóa lớp học.
- Ghi danh học viên theo lớp.
- Rút học viên khỏi lớp.
- Phân công giáo viên cho lớp.
- Quản lý lịch tuần lớp.

### 2.3.4 Xếp lịch thông minh

- Tạo preview scheduling.
- Xem preview và các conflict.
- Benchmark nhiều solver.
- Commit preview thành lesson.

### 2.3.5 Học vụ sau lesson

- Điểm danh.
- Tạo lesson summary.
- Ghi nhận academic record.
- Quản lý leave request.

### 2.3.6 Dự báo học viên có nguy cơ học kém

- Thu thập và chuẩn hóa dữ liệu học tập.
- Xây dựng đặc trưng phục vụ mô hình.
- Huấn luyện và đánh giá mô hình `AT_RISK`.
- Cảnh báo sớm cho quản trị viên.

## 2.4 Các use case quan trọng

Trong số các use case của hệ thống, nhóm use case có giá trị nhất đối với báo cáo tốt nghiệp là:

- Tạo lớp học.
- Ghi danh học viên.
- Phân công giáo viên.
- Cấu hình lịch tuần.
- Tạo preview xếp lịch.
- Commit preview thành lesson.
- Benchmark thuật toán xếp lịch.
- Điểm danh và tổng kết buổi học.
- Ghi nhận kết quả học tập.
- Dự báo học viên có nguy cơ học kém.

## 2.5 Luồng nghiệp vụ chính

Luồng nghiệp vụ cốt lõi của EduCenter có thể tóm tắt như sau:

1. Quản trị viên tạo các dữ liệu nền như phòng học, ca học, khóa học, chương trình, giáo viên.
2. Quản trị viên mở lớp học và gán các thông tin liên quan.
3. Quản trị viên ghi danh học viên và cấu hình lịch tuần cho lớp.
4. Hệ thống tạo preview xếp lịch dựa trên dữ liệu lớp, `Shift`, giáo viên và phòng học.
5. Quản trị viên kiểm tra preview, xử lý xung đột nếu có.
6. Preview hợp lệ được commit để sinh `Lesson`.
7. Giáo viên sử dụng dữ liệu lesson để điểm danh, tổng kết buổi học và nhập kết quả học tập.

**[Chèn Hình 2.2: BPMN hoặc activity diagram cho luồng mở lớp - ghi danh - phân công giáo viên - xếp lịch]**

## 2.6 Yêu cầu và ràng buộc nghiệp vụ

Một số ràng buộc quan trọng của hệ thống bao gồm:

- Không được trùng giáo viên ở cùng một thời điểm.
- Không được trùng phòng ở cùng một thời điểm.
- Một lớp chỉ được xếp vào các ngày và ca học đã được khai báo trong `ClassSchedule`.
- Chỉ các `Shift` đang hoạt động mới được dùng cho scheduling.
- Preview chỉ được commit khi đạt điều kiện hợp lệ.
- Lesson là đầu ra chính thức sau khi scheduling được xác nhận.

## 2.7 Nhận xét hiện trạng triển khai

Tại thời điểm tổng hợp báo cáo, phần lớn các module quan trọng phục vụ mục tiêu của đồ án đã có mặt trong codebase. Đặc biệt, tài liệu đối chiếu chức năng với mã nguồn cho thấy các nhánh như lesson, attendance, lesson summary, academic record, leave request và predictive alerts đã được triển khai ở mức đủ để đưa vào phạm vi báo cáo chính, dù một số nhánh frontend vẫn còn cần hoàn thiện thêm.

---

# CHƯƠNG 3. THIẾT KẾ HỆ THỐNG

## 3.1 Kiến trúc tổng thể

Hệ thống EduCenter được xây dựng theo hướng tách biệt rõ các thành phần chính:

- Frontend React + TypeScript đảm nhận giao diện quản trị và một phần giao diện người dùng.
- Backend Golang với Gin chịu trách nhiệm xử lý nghiệp vụ, điều phối use case và cung cấp REST API.
- PostgreSQL lưu trữ dữ liệu hoạt động của hệ thống.
- Scheduling layer được tổ chức qua abstraction `SchedulingSolver`, cho phép thay thế hoặc benchmark nhiều thuật toán.

Điểm nổi bật của hệ thống là việc áp dụng `Clean Architecture`, trong đó:

- `entities` biểu diễn thực thể nghiệp vụ,
- `repositories` xử lý truy cập dữ liệu,
- `usecases` đóng vai trò điều phối logic nghiệp vụ,
- `infrastructure` đảm nhận tích hợp với database và thành phần ngoài.

**[Chèn Hình 3.1: Kiến trúc tổng thể hệ thống EduCenter]**

## 3.2 Thiết kế dữ liệu

Các thực thể trọng tâm của hệ thống gồm:

- `User`
- `Student`
- `Teacher`
- `Course`
- `Program`
- `ProgramCourse`
- `Room`
- `Shift`
- `Class`
- `ClassSchedule`
- `Enrollment`
- `Lesson`
- `Attendance`
- `LessonSummary`
- `AcademicRecord`
- `LeaveRequest`

Trong đó:

- `Class` là thực thể trung tâm của vận hành học vụ.
- `ClassSchedule` lưu ràng buộc lịch tuần chuẩn theo `day_of_week` và `shift_id`.
- `Enrollment` biểu diễn roster giữa học viên và lớp.
- `Lesson` là kết quả chính thức sau khi commit preview scheduling.
- `Attendance`, `LessonSummary`, `AcademicRecord`, `LeaveRequest` phục vụ chuỗi nghiệp vụ sau lesson.

**[Chèn Hình 3.2: ERD logic của EduCenter]**

## 3.2.1 Nguyên tắc đặc tả cơ sở dữ liệu

Phần đặc tả dữ liệu trong báo cáo này được tổng hợp từ ba nguồn chính của hệ thống: `entities` trong backend Golang, các file migration PostgreSQL và các use case đang vận hành thực tế. Cách làm này giúp phần mô tả cơ sở dữ liệu bám sát trạng thái code hiện tại thay vì chỉ mô tả ở mức ý tưởng.

Một số quy ước thiết kế dữ liệu đang được dùng thống nhất trong EduCenter:

- Khóa chính của hầu hết các bảng sử dụng UUID.
- Cơ sở dữ liệu mục tiêu là PostgreSQL.
- Các bảng giao dịch quan trọng dùng `created_at`, `updated_at` để truy vết thay đổi.
- Một số bảng dùng `deleted_at` để xóa mềm nhằm bảo toàn lịch sử.
- Các trường trạng thái đang được cài theo kiểu chuỗi hoặc số nguyên; một số domain đã chốt rõ trong code, một số vẫn cần chuẩn hóa thêm ở bước hoàn thiện sản phẩm.

## 3.2.2 Danh mục các bảng trọng tâm của hệ thống

| Nhóm dữ liệu | Bảng | Vai trò |
| --- | --- | --- |
| Tài khoản và bảo mật | `users`, `user_otps`, `password_resets` | quản lý đăng nhập, kích hoạt tài khoản, quên mật khẩu |
| Danh mục đào tạo | `students`, `teachers`, `courses`, `programs`, `program_courses` | quản lý học viên, giáo viên, khóa học và chương trình |
| Tài nguyên vận hành | `rooms`, `shifts` | chuẩn hóa phòng học và ca học |
| Tổ chức lớp | `classes`, `class_schedules`, `enrollments` | mở lớp, lịch tuần, roster học viên |
| Vận hành buổi học | `lessons`, `attendances`, `lesson_summaries`, `academic_records`, `leave_requests` | quản lý lesson sau scheduling và học vụ sau lesson |
| Hỗ trợ mở rộng | `consultations`, `materials`, `audit_logs`, `approval_decisions`, `labels`, `objectives`, `outcomes` | tư vấn tuyển sinh, kiểm duyệt tài liệu, mục tiêu học tập |

## 3.2.3 Đặc tả các bảng lõi

### a. Bảng `users`

| Trường | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK, default `uuid_generate_v4()` | mã định danh tài khoản |
| `code` | varchar(50) | unique | mã nội bộ người dùng |
| `full_name` | varchar(255) |  | họ tên hiển thị |
| `email` | varchar(255) | not null, unique | email đăng nhập chính |
| `password` | text | not null | mật khẩu đã băm |
| `role` | varchar(50) | default `STUDENT` | vai trò hệ thống, đang dùng các giá trị như `ADMIN`, `TEACHER`, `STUDENT` |
| `is_active` | boolean | default `true` ở entity | cờ kích hoạt tài khoản; luồng đăng ký hiện đang override thành `false` trước khi verify OTP |
| `created_at`, `updated_at` | timestamp | audit | thời gian tạo/cập nhật |
| `deleted_at` | timestamp | soft delete | phục vụ xóa mềm |

### b. Bảng `user_otps`

| Trường | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã OTP record |
| `user_id` | UUID | not null, index | tham chiếu `users.id` |
| `otp_hash` | text | not null | mã OTP đã băm |
| `expired_at` | timestamp | not null | thời điểm hết hạn |
| `used_at` | timestamp | nullable | đánh dấu OTP đã dùng |
| `created_at`, `deleted_at` | timestamp | audit / soft delete | lịch sử xác minh |

### c. Bảng `password_resets`

| Trường | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã yêu cầu reset |
| `user_id` | UUID | not null | tham chiếu `users.id` |
| `token_hash` | text | not null | token reset đã băm |
| `expires_at` | timestamp | not null | hết hạn token |
| `used_at` | timestamp | nullable | token đã dùng hay chưa |
| `requested_ip` | varchar(45) |  | IP gửi yêu cầu |
| `user_agent` | text |  | thông tin thiết bị |
| `created_at`, `updated_at`, `deleted_at` | timestamp | audit | truy vết bảo mật |

### d. Bảng `students`

| Trường | Kiểu dữ liệu | Ràng buộc / miền giá trị | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã học viên |
| `code` | varchar(50) | unique | mã học viên |
| `full_name` | varchar(255) |  | họ tên học viên |
| `email` | varchar(255) |  | email liên hệ |
| `phone` | varchar(20) |  | số điện thoại học viên |
| `guardian_phone` | varchar(20) |  | số điện thoại phụ huynh |
| `grade_level` | varchar(50) |  | khối lớp của học viên |
| `status` | varchar(50) | default `ACTIVE` | trạng thái học viên |
| `date_of_birth` | timestamp | nullable | ngày sinh |
| `gender` | varchar(20) |  | giới tính |
| `address` | text |  | địa chỉ |
| `created_at`, `updated_at`, `deleted_at` | timestamp | audit / soft delete | quản lý vòng đời hồ sơ |

### e. Bảng `teachers`

| Trường | Kiểu dữ liệu | Ràng buộc / miền giá trị | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã giáo viên |
| `code` | varchar(50) | unique | mã giáo viên |
| `full_name` | varchar(255) |  | họ tên giáo viên |
| `email` | varchar(255) |  | email liên hệ |
| `phone` | varchar(20) |  | số điện thoại |
| `is_school_teacher` | boolean | default `false` | có phải giáo viên trường phổ thông hay không |
| `school_name` | varchar(255) |  | tên trường công tác |
| `employment_type` | varchar(50) | default `PART_TIME`; domain đang dùng `PART_TIME`, `FULL_TIME` | hình thức cộng tác |
| `status` | varchar(50) | default `ACTIVE`; domain đang dùng `ACTIVE`, `INACTIVE` | trạng thái làm việc |
| `notes` | text |  | ghi chú bổ sung |
| `created_at`, `updated_at`, `deleted_at` | timestamp | audit / soft delete | quản lý vòng đời hồ sơ |

### f. Bảng `courses`

| Trường | Kiểu dữ liệu | Ràng buộc / miền giá trị | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã khóa học |
| `code` | varchar(50) | unique, not null | mã khóa học |
| `name` | varchar(255) |  | tên khóa học |
| `description` | text |  | mô tả nội dung |
| `grade_level` | varchar(50) |  | khối lớp áp dụng |
| `subject` | varchar(255) |  | môn học |
| `session_count` | int |  | số buổi học |
| `session_duration_minutes` | int |  | thời lượng mỗi buổi |
| `total_hours` | numeric(8,2) |  | tổng số giờ học |
| `price` | numeric(10,2) |  | học phí |
| `status` | varchar(50) | default `ACTIVE` | trạng thái khóa học |
| `created_at`, `updated_at`, `deleted_at` | timestamp | audit / soft delete | vòng đời dữ liệu |

### g. Bảng `programs`

| Trường | Kiểu dữ liệu | Ràng buộc / miền giá trị | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã chương trình |
| `code` | varchar(50) | unique, not null | mã chương trình |
| `name` | varchar(255) |  | tên chương trình |
| `track` | varchar(50) | dấu vết domain `SUPPORT`, `BASIC`, `ADVANCED` | nhánh hoặc cấp độ chương trình |
| `effective_from`, `effective_to` | timestamp | nullable | thời gian hiệu lực |
| `created_by_id` | UUID | nullable | người khởi tạo |
| `approved_by_id` | UUID | nullable | người duyệt |
| `approval_note` | text |  | ghi chú duyệt |
| `published_at` | timestamp | nullable | thời điểm xuất bản |
| `archived_at` | timestamp | nullable | thời điểm lưu trữ |
| `created_at`, `updated_at`, `deleted_at` | timestamp | audit / soft delete | vòng đời dữ liệu |

Lưu ý rằng migration cũ của bảng `programs` từng có trường `status`, nhưng entity hiện tại đã chuyển sang mô hình lifecycle bằng `published_at` và `archived_at`. Đây là một điểm cần nêu rõ khi bảo vệ để tránh nhầm lẫn giữa trạng thái logical và trạng thái nghiệp vụ cũ.

### h. Bảng `program_courses`

| Trường | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã mapping |
| `program_id` | UUID | not null | tham chiếu `programs.id` |
| `course_id` | UUID | not null | tham chiếu `courses.id` |

Bảng này dùng để cài đặt quan hệ nhiều-nhiều giữa chương trình đào tạo và khóa học.

### i. Bảng `rooms`

| Trường | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã phòng |
| `code` | varchar(50) | unique, not null | mã phòng |
| `name` | varchar(255) | not null | tên phòng |
| `capacity` | int | nên >= 1 | sức chứa |
| `address` | text |  | địa điểm phòng học |
| `created_at`, `updated_at` | timestamp | audit | truy vết tạo/sửa |

### j. Bảng `shifts`

| Trường | Kiểu dữ liệu | Ràng buộc / miền giá trị | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã ca học |
| `code` | varchar(50) | unique, not null | mã ca |
| `name` | varchar(255) | not null | tên ca |
| `start_time` | varchar(10) | not null | giờ bắt đầu, ví dụ `08:00` |
| `end_time` | varchar(10) | not null | giờ kết thúc |
| `duration_minutes` | int | not null, nên >= 1 | độ dài ca |
| `session_type` | varchar(50) | domain đang dùng `MORNING`, `AFTERNOON`, `EVENING`, `CUSTOM` | loại ca học |
| `is_active` | boolean | default `true` | ca có được phép scheduling hay không |
| `notes` | text |  | ghi chú |
| `created_at`, `updated_at`, `deleted_at` | timestamp | audit / soft delete | vòng đời dữ liệu |

### k. Bảng `classes`

| Trường | Kiểu dữ liệu | Ràng buộc / miền giá trị | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã lớp học |
| `code` | varchar(50) | unique, not null | mã lớp |
| `name` | varchar(255) | not null | tên lớp |
| `notes` | text |  | ghi chú lớp |
| `start_date` | timestamp | not null | ngày bắt đầu |
| `end_date` | timestamp | nullable | ngày kết thúc |
| `max_students` | int |  | sĩ số tối đa |
| `status` | varchar(50) | default `OPEN`; domain đang dùng `OPEN`, `CLOSED`, `CANCELLED` | trạng thái vận hành lớp |
| `price` | numeric(10,2) |  | học phí lớp |
| `program_id` | UUID | nullable | lớp gắn chương trình nào |
| `course_id` | UUID | nullable | lớp gắn khóa học nào |
| `teacher_id` | UUID | nullable | giáo viên phụ trách |
| `room_id` | UUID | nullable | phòng mặc định |
| `created_at`, `updated_at`, `deleted_at` | timestamp | audit / soft delete | vòng đời lớp |

Trong hệ thống hiện tại, scheduling chỉ tải các lớp có `status = OPEN`, vì vậy đây là trường rất quan trọng trong logic nghiệp vụ.

### l. Bảng `class_schedules`

| Trường | Kiểu dữ liệu | Ràng buộc / miền giá trị | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã lịch tuần |
| `class_id` | UUID | not null | tham chiếu `classes.id` |
| `day_of_week` | varchar(20) | not null; thực tế đang dùng `MONDAY` đến `SUNDAY` | thứ học trong tuần |
| `shift_id` | UUID | not null | tham chiếu `shifts.id` |
| `room_id` | UUID | nullable | phòng cố định cho slot này nếu có |

`class_schedules` là bảng đặc biệt quan trọng vì nó quy định miền slot hợp lệ trước khi chạy solver. Dữ liệu preview scheduling hiện đọc lịch tuần theo cặp `day_of_week + shift_id`.

### m. Bảng `enrollments`

| Trường | Kiểu dữ liệu | Ràng buộc / miền giá trị | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã ghi danh |
| `class_id` | UUID | not null | tham chiếu lớp |
| `student_id` | UUID | not null | tham chiếu học viên |
| `status` | varchar(50) | default `APPLIED`; use case hiện tạo `ENROLLED` trực tiếp | trạng thái ghi danh |
| `approved_at` | timestamp | nullable | thời điểm duyệt |
| `rejected_at` | timestamp | nullable | thời điểm từ chối |
| `created_at`, `updated_at` | timestamp | audit | truy vết thay đổi |

Điểm cần lưu ý là entity định nghĩa mặc định `APPLIED`, nhưng use case ghi danh hiện đang chèn bản ghi với `status = ENROLLED`. Đây là một khoảng lệch cần ghi nhận trong báo cáo như một debt nghiệp vụ.

### n. Bảng `lessons`

| Trường | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã buổi học |
| `class_id` | UUID | not null | lớp tương ứng |
| `teacher_id` | UUID | nullable | giáo viên của lesson |
| `date_start` | timestamp | not null | thời điểm bắt đầu |
| `date_end` | timestamp | not null | thời điểm kết thúc |
| `room_id` | UUID | nullable | phòng dạy |
| `notes` | text |  | ghi chú sinh từ preview hoặc vận hành |
| `created_at`, `updated_at`, `deleted_at` | timestamp | audit / soft delete | vòng đời lesson |

`lessons` là đầu ra chính thức của quá trình commit preview scheduling.

### o. Bảng `attendances`

| Trường | Kiểu dữ liệu | Ràng buộc / miền giá trị | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã điểm danh |
| `lesson_id` | UUID | not null | tham chiếu lesson |
| `student_id` | UUID | not null | tham chiếu học viên |
| `status` | int | not null; code nội bộ hiện dùng `1=PRESENT`, `2=ABSENT`, `3=EXCUSED`, `4=LATE`, `5=EARLY` | trạng thái chuyên cần |
| `note` | text |  | ghi chú điểm danh |
| `marked_at` | timestamp | nullable | thời điểm chấm |
| `created_at`, `updated_at`, `deleted_at` | timestamp | audit / soft delete | truy vết lịch sử |

Teacher portal hiện dùng tập giá trị hiển thị riêng là `0=ABSENT`, `1=PRESENT`, `2=LATE`, `3=EXCUSED`, còn giá trị `EARLY` được giữ như legacy ở tầng use case. Đây là bảng có rủi ro đồng bộ domain lớn nhất và nên được chuẩn hóa sớm.

### p. Bảng `lesson_summaries`

| Trường | Kiểu dữ liệu | Ràng buộc | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã tổng kết |
| `lesson_id` | UUID | unique, not null | mỗi lesson có tối đa một summary |
| `topic` | text |  | chủ đề buổi học |
| `lesson_content` | text |  | nội dung đã dạy |
| `class_feedback` | text |  | phản hồi chung |
| `homework` | text |  | bài tập giao |
| `homework_deadline` | timestamp | nullable | hạn nộp bài |
| `teacher_notes` | text |  | ghi chú của giáo viên |
| `created_by_id` | UUID | nullable | người tạo summary |
| `created_at`, `updated_at`, `deleted_at` | timestamp | audit / soft delete | vòng đời dữ liệu |

### q. Bảng `academic_records`

| Trường | Kiểu dữ liệu | Ràng buộc / miền giá trị | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã record |
| `lesson_summary_id` | UUID | not null | tham chiếu summary cha |
| `student_id` | UUID | not null | học viên được đánh giá |
| `homework_completed` | boolean | default `false` | đã làm bài tập hay chưa |
| `homework_score` | numeric(5,2) |  | điểm bài tập |
| `attitude_rating` | int | scale chưa chốt | đánh giá thái độ học tập |
| `participation_score` | numeric(5,2) |  | điểm tham gia |
| `personal_comment` | text |  | nhận xét cá nhân |
| `total_score` | numeric(5,2) | công thức chưa khóa cứng | tổng hợp kết quả |
| `is_completed` | boolean | default `false` | đã chốt record hay chưa |
| `created_at`, `updated_at`, `deleted_at` | timestamp | audit / soft delete | vòng đời dữ liệu |

### r. Bảng `leave_requests`

| Trường | Kiểu dữ liệu | Ràng buộc / miền giá trị | Ý nghĩa |
| --- | --- | --- | --- |
| `id` | UUID | PK | mã đơn |
| `student_id` | UUID | not null | học viên gửi đơn |
| `leave_type` | varchar(50) | domain đang dùng `LEAVE`, `LATE`, `EARLY` | loại đơn |
| `apply_date` | timestamp | not null | ngày áp dụng |
| `late_minutes` | int | dùng cho `LATE` | số phút đi muộn |
| `early_minutes` | int | dùng cho `EARLY` | số phút về sớm |
| `reason` | text | not null | lý do |
| `documents` | text[] | nullable | minh chứng đính kèm |
| `class_id` | UUID | nullable | lớp liên quan |
| `lesson_id` | UUID | nullable | lesson liên quan |
| `subject` | varchar(255) |  | tiêu đề đơn |
| `status` | varchar(50) | default `PENDING`; flow hiện dùng thêm `APPROVED`, `REJECTED` | trạng thái xử lý |
| `approved_by_id` | UUID | nullable | người duyệt |
| `approved_at` | timestamp | nullable | thời điểm duyệt |
| `rejection_reason` | text |  | lý do từ chối |
| `created_at`, `updated_at`, `deleted_at` | timestamp | audit / soft delete | vòng đời dữ liệu |

## 3.2.4 Các bảng hỗ trợ và mở rộng

Ngoài nhóm bảng lõi phục vụ trực tiếp cho bài toán quản lý trung tâm dạy thêm, hệ thống còn có một số bảng hỗ trợ:

- `consultations`: lưu lead tư vấn với các trường chính `full_name`, `phone`, `grade_level`, `notes`, `status`.
- `materials`: lưu tài liệu do giáo viên tải lên cùng metadata file và trạng thái xử lý.
- `audit_logs`: lưu kết quả OCR/AI audit cho tài liệu, gồm `provider`, `confidence_score`, `reasoning`, `detected_issues`.
- `approval_decisions`: lưu quyết định duyệt/từ chối tài liệu của người phụ trách.
- `labels`: bảng mã nhãn kiểm duyệt như `SAFE`, `WARNING`, `DANGER`.
- `objectives`, `outcomes`: hỗ trợ biểu diễn mục tiêu và chuẩn đầu ra của chương trình đào tạo.

Nhóm bảng này chưa phải là trọng tâm của phạm vi bảo vệ, nhưng thể hiện kiến trúc dữ liệu của hệ thống đã tính tới khả năng mở rộng về kiểm duyệt nội dung, đánh giá học thuật và chuẩn đầu ra.

## 3.2.5 Giá trị mặc định và domain dữ liệu quan trọng

Một số giá trị mặc định đang có vai trò trực tiếp trong nghiệp vụ:

| Trường | Giá trị mặc định hiện tại | Ý nghĩa |
| --- | --- | --- |
| `users.role` | `STUDENT` | người đăng ký mới mặc định là học viên |
| `students.status` | `ACTIVE` | học viên đang hoạt động |
| `teachers.employment_type` | `PART_TIME` | giáo viên cộng tác mặc định |
| `teachers.status` | `ACTIVE` | giáo viên có thể nhận lớp |
| `courses.status` | `ACTIVE` | khóa học sẵn sàng sử dụng |
| `shifts.is_active` | `true` | ca học được phép scheduling |
| `classes.status` | `OPEN` | lớp được phép tham gia scheduling |
| `enrollments.status` | `APPLIED` | có ý đồ approval flow, dù code hiện đang tạo `ENROLLED` trực tiếp |
| `academic_records.homework_completed` | `false` | mặc định chưa hoàn thành bài tập |
| `academic_records.is_completed` | `false` | mặc định chưa chốt record |
| `leave_requests.status` | `PENDING` | đơn chờ duyệt |

Các domain dữ liệu hiện đã nhìn thấy tương đối rõ trong hệ thống gồm:

| Trường | Domain giá trị |
| --- | --- |
| `users.role` | `ADMIN`, `TEACHER`, `STUDENT`, ngoài ra một số use case có dấu vết `SUPER_ADMIN`, `PARENT` |
| `teachers.employment_type` | `PART_TIME`, `FULL_TIME` |
| `teachers.status` | `ACTIVE`, `INACTIVE` |
| `shifts.session_type` | `MORNING`, `AFTERNOON`, `EVENING`, `CUSTOM` |
| `classes.status` | `OPEN`, `CLOSED`, `CANCELLED` |
| `class_schedules.day_of_week` | `MONDAY`, `TUESDAY`, `WEDNESDAY`, `THURSDAY`, `FRIDAY`, `SATURDAY`, `SUNDAY` |
| `leave_requests.leave_type` | `LEAVE`, `LATE`, `EARLY` |
| `leave_requests.status` | `PENDING`, `APPROVED`, `REJECTED` |
| `materials.status` | `UPLOADED`, `SCANNING`, `AI_REVIEWED`, `APPROVED`, `REJECTED` |
| `attendances.status` | `1=PRESENT`, `2=ABSENT`, `3=EXCUSED`, `4=LATE`, `5=EARLY` ở tầng nội bộ |

## 3.2.6 Nhận xét về thiết kế dữ liệu hiện tại

Thiết kế dữ liệu của EduCenter đã đủ tốt để phục vụ một hệ thống quản lý trung tâm dạy thêm có chiều sâu vận hành. Điểm mạnh là mô hình dữ liệu đã bao phủ được toàn bộ chuỗi từ dữ liệu master, mở lớp, xếp lịch, lesson cho tới học vụ sau lesson. Ngoài ra, việc tách `Shift` thành bảng riêng và dùng `ClassSchedule` theo `day_of_week + shift_id` là một bước thiết kế rất phù hợp cho bài toán scheduling.

Tuy nhiên, qua đối chiếu giữa entity, migration và use case, vẫn còn một số điểm cần hoàn thiện:

- Lifecycle của `Enrollment` chưa thống nhất giữa entity và use case.
- Domain của `Attendance.status` chưa được chuẩn hóa hoàn toàn giữa tầng nội bộ và teacher portal.
- `Program` đang có dấu vết chuyển từ mô hình `status` sang mô hình lifecycle bằng `published_at` và `archived_at`.
- Một số bảng hỗ trợ như `consultations`, `materials` và `approval_decisions` đã có dữ liệu và schema tốt nhưng chưa được đưa vào flow nghiệp vụ chính thức đầy đủ.

## 3.3 Thiết kế lớp nghiệp vụ

Từ góc nhìn thiết kế, có thể chia hệ thống thành các nhóm lớp chính:

- Nhóm xác thực và tài khoản.
- Nhóm dữ liệu đào tạo và nguồn lực.
- Nhóm lớp học và ghi danh.
- Nhóm scheduling.
- Nhóm vận hành học vụ.
- Nhóm phân tích dự báo.

Các lớp scheduling được thiết kế theo hướng mở rộng:

- `SchedulingSolver` là interface chung.
- `GraphColoringSolver`, `CPSATSolver`, `TabuSearchSolver` là các implementation cụ thể.
- `SchedulingInput`, `SchedulingOutput`, `BenchmarkResult` được chuẩn hóa dùng chung.

**[Chèn Hình 3.3: Biểu đồ lớp chính của hệ thống]**

## 3.4 Thiết kế luồng trình tự

Một số luồng trình tự quan trọng cần mô tả trong báo cáo gồm:

- Luồng đăng ký và xác minh OTP.
- Luồng tạo lớp, ghi danh và phân công giáo viên.
- Luồng tạo preview scheduling.
- Luồng commit preview thành lesson.
- Luồng điểm danh và tổng kết buổi học.

**[Chèn Hình 3.4: Sequence diagram đăng nhập / đăng ký]**  
**[Chèn Hình 3.5: Sequence diagram tạo lớp - ghi danh - phân công giáo viên]**  
**[Chèn Hình 3.6: Sequence diagram auto scheduling]**

## 3.5 Thiết kế giao diện

Phần giao diện của hệ thống được xây dựng để phục vụ thao tác quản trị rõ ràng, tập trung vào các module sau:

- Trang đăng nhập / xác thực.
- Quản lý học viên.
- Quản lý giáo viên.
- Quản lý chương trình và khóa học.
- Quản lý lớp học và chi tiết lớp.
- Quản lý ca học `Shift`.
- Scheduling dashboard.
- Lesson detail.
- Trang predictive analytics.

Trong bản báo cáo hiện tại, phần hình ảnh giao diện sẽ được bổ sung sau. Nội dung mô tả giao diện có thể trình bày theo từng nhóm chức năng để đảm bảo tính liền mạch.

**[Chèn Hình 3.7: Giao diện đăng nhập]**  
**[Chèn Hình 3.8: Giao diện quản lý lớp học]**  
**[Chèn Hình 3.9: Giao diện scheduling preview]**  
**[Chèn Hình 3.10: Giao diện lesson detail / attendance / summary]**  
**[Chèn Hình 3.11: Giao diện predictive analytics]**

---

# CHƯƠNG 4. XẾP LỊCH THÔNG MINH VÀ ĐÁNH GIÁ THUẬT TOÁN

## 4.1 Vai trò của bài toán scheduling trong hệ thống

Scheduling là điểm nhấn kỹ thuật chính của đồ án. Khác với các module CRUD thông thường, scheduling yêu cầu hệ thống phải đồng thời xét nhiều ràng buộc và tìm ra lịch học hợp lệ trong một khoảng thời gian xác định. Đây là bài toán tối ưu hóa có tính ứng dụng cao trong môi trường trung tâm dạy thêm.

## 4.2 Chuẩn hóa dữ liệu thời gian bằng `Shift`

Một trong những thay đổi quan trọng của đồ án là chuẩn hóa dữ liệu thời gian bằng thực thể `Shift`.

Thay vì lưu trực tiếp các khung giờ rời rạc trong lịch tuần của lớp, hệ thống sử dụng:

- `Shift` để định nghĩa ca học chuẩn,
- `ClassSchedule` để gắn lớp với `day_of_week` và `shift_id`.

Cách tiếp cận này giúp:

- đơn giản hóa việc xây dựng domain scheduling,
- giảm rủi ro lệch dữ liệu giữa frontend và backend,
- tạo nền tảng tốt hơn cho benchmark và mở rộng thuật toán.

## 4.3 Các ràng buộc của bài toán

Các ràng buộc chính của bài toán scheduling gồm:

### 4.3.1 Hard constraints

- Không trùng giáo viên.
- Không trùng phòng học.
- Không trùng lớp học.
- Chỉ dùng các `Shift` còn active.
- Chỉ xếp vào các slot do `ClassSchedule` cho phép.

### 4.3.2 Soft constraints

- Ưu tiên nghiệm có chất lượng tốt hơn theo scoring.
- Giảm mức độ xung đột hoặc phân bố chưa hợp lý.
- Tối ưu mức sử dụng nguồn lực trong giới hạn bài toán.

## 4.4 Các thuật toán được cài đặt

Để giải bài toán scheduling, hệ thống đã cài đặt ba solver chính:

### 4.4.1 Graph Coloring + heuristic

Đây là baseline heuristic có runtime rất nhanh, phù hợp để dùng làm chuẩn so sánh. Solver này cho kết quả gần như tức thời ở mọi scenario benchmark, tuy nhiên ở dữ liệu nhỏ có thể cho soft score chưa tối ưu bằng các cách tiếp cận mạnh hơn.

### 4.4.2 CP-SAT

`CP-SAT` được lựa chọn như solver có định hướng tối ưu hóa tốt hơn. Kết quả benchmark cho thấy solver này giữ được nghiệm sạch theo hard constraints, có chất lượng nghiệm tốt và runtime vẫn ở mức chấp nhận được trong phạm vi đồ án.

### 4.4.3 Tabu Search

`Tabu Search` đại diện cho nhóm metaheuristic. Solver này được giữ lại trong benchmark để so sánh về khả năng tìm nghiệm, nhưng kết quả thực nghiệm cho thấy chưa tạo ra lợi thế đủ lớn so với chi phí runtime ở scenario lớn.

## 4.5 Bộ dữ liệu benchmark

Benchmark được chạy trên bộ dữ liệu tổng hợp có thể lặp lại, chia thành ba scenario:

| Scenario | Số lớp | Số giáo viên | Số phòng | Số ca | Số session yêu cầu |
| --- | ---: | ---: | ---: | ---: | ---: |
| Small | 6 | 4 | 3 | 3 | 12 |
| Medium | 10 | 5 | 4 | 3 | 30 |
| Large | 16 | 7 | 5 | 3 | 64 |

Việc dùng dữ liệu tổng hợp có kiểm soát giúp benchmark công bằng giữa các solver và cho phép chạy lặp lại nhiều lần để kiểm tra tính ổn định.

## 4.6 Kết quả benchmark

Theo báo cáo benchmark ngày `14/04/2026`, cả ba solver đều đạt `feasibility = 1.000` và không sinh hard violation trong các scenario được kiểm thử.

Một số nhận xét quan trọng:

- `Graph Coloring + heuristic` là solver nhanh nhất.
- `CP-SAT` có chất lượng nghiệm tốt hơn ở scenario nhỏ và giữ runtime chấp nhận được ở scenario lớn.
- `Tabu Search` chậm hơn đáng kể ở scenario lớn nhưng không tạo ưu thế tương xứng về chất lượng nghiệm.

## 4.7 Quyết định lựa chọn solver chính

Từ các số liệu benchmark, `CP-SAT` được chọn làm solver mặc định cho giai đoạn production-like của API scheduling.

Lý do lựa chọn:

1. Giữ được nghiệm sạch về hard constraints.
2. Thể hiện tín hiệu tốt hơn về chất lượng nghiệm ở các scenario có phân biệt soft score.
3. Runtime vẫn phù hợp để demo và trình bày trong phạm vi đồ án.

## 4.8 Hướng phát triển cho bài toán dự báo `AT_RISK`

Bên cạnh scheduling, đồ án định hướng bổ sung một bài toán machine learning có ý nghĩa thực tiễn là dự báo học viên có nguy cơ học kém.

Hướng triển khai của bài toán gồm:

- xác định nhãn `AT_RISK / NOT_AT_RISK`,
- chuẩn hóa dữ liệu attendance, lesson summary, academic record,
- lựa chọn feature set phù hợp,
- huấn luyện mô hình nhẹ ở backend,
- cung cấp API dự báo và cảnh báo sớm cho quản trị viên.

Tại thời điểm viết bản thảo này, scheduling đã là trục hoàn chỉnh nhất để đưa vào nội dung bảo vệ; predictive analytics là nhánh đang tiếp tục được hoàn thiện để bổ sung chiều sâu nghiên cứu cho hệ thống.

**[Chèn Hình 4.1: Bảng benchmark 3 solver]**  
**[Chèn Hình 4.2: Biểu đồ so sánh runtime / soft score giữa các solver]**  
**[Chèn Hình 4.3: Luồng dữ liệu dự báo AT_RISK]**

---

# CHƯƠNG 5. KẾT LUẬN VÀ HƯỚNG PHÁT TRIỂN

## 5.1 Kết quả đạt được

Đồ án đã xây dựng được một hệ thống quản lý trung tâm dạy thêm có cấu trúc rõ ràng, bám sát kiến trúc phần mềm hiện đại và có giá trị thực tiễn trong quản lý đào tạo.

Các kết quả chính đạt được gồm:

- Hoàn thiện phần lớn các module quản trị cốt lõi.
- Chuẩn hóa dữ liệu thời gian bằng thực thể `Shift`.
- Xây dựng thành công luồng scheduling preview -> commit -> lesson.
- Cài đặt và benchmark ba thuật toán scheduling.
- Lựa chọn `CP-SAT` làm solver chính dựa trên số liệu thực nghiệm.
- Tạo nền dữ liệu và kiến trúc để phát triển bài toán dự báo học viên có nguy cơ học kém.

## 5.2 Hạn chế hiện tại

Dù đã đạt được nhiều kết quả, hệ thống vẫn còn một số hạn chế:

- Một số hình ảnh minh họa, biểu đồ và giao diện chưa được bổ sung vào báo cáo.
- Một số nhánh frontend ở teacher/student portal vẫn cần hoàn thiện thêm để đồng bộ với backend mới.
- Module predictive analytics cần được hoàn thiện thêm về dữ liệu, metric và giao diện cảnh báo.
- Một số chức năng ngoài scope chính vẫn còn tồn tại trong codebase và cần được phân tách rõ khi trình bày báo cáo.

## 5.3 Hướng phát triển

Trong giai đoạn tiếp theo, hệ thống có thể được mở rộng theo các hướng sau:

- Hoàn thiện teacher portal và student portal theo actor-based flow.
- Bổ sung dashboard học tập và cảnh báo sớm học viên `AT_RISK`.
- Mở rộng benchmark với nhiều bộ dữ liệu thực tế hơn.
- Tối ưu thêm soft constraints của scheduling để nâng cao chất lượng nghiệm.
- Hoàn thiện thêm tài liệu và bộ hình minh họa phục vụ bảo vệ tốt nghiệp.

---

# TÀI LIỆU THAM KHẢO

1. README của project Doan tại `/Users/hant/golang/doan/README.md`.  
2. `PROJECT_TASKS.md` — tài liệu chốt mục tiêu và phạm vi đồ án.  
3. `docs/tai_lieu_phan_tich_nghiep_vu.md` — tài liệu phân tích nghiệp vụ hệ thống EduCenter.  
4. `docs/BA_SYSTEM_ANALYSIS_REPORT_EDUCENTER.md` — báo cáo khảo sát và phân tích hệ thống.  
5. `docs/FUNCTION_DECOMPOSITION_ALIGNMENT_2026-04-15.md` — đối chiếu chức năng với codebase.  
6. `docs/SCHEDULING_BENCHMARK_REPORT_2026-04-14.md` — báo cáo benchmark scheduling.  
7. Các sơ đồ mô hình hóa trong `docs/modeling/` và `docs/modeling/drawio/`.  

---

# GHI CHÚ BỔ SUNG KHI HOÀN THIỆN BẢN CUỐI

- Bổ sung logo trường và định dạng căn giữa cho bìa.
- Điền tên giảng viên hướng dẫn.
- Chèn tự động mục lục trong Google Docs sau khi ổn định heading.
- Chèn hình use case, BPMN, ERD, class diagram, sequence diagram từ các file đã dựng.
- Chèn ảnh giao diện và bảng benchmark minh họa.
- Rà soát lại số hiệu hình và bảng sau khi thêm ảnh.
