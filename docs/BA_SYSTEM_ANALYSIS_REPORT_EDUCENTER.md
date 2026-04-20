# BÁO CÁO KHẢO SÁT VÀ PHÂN TÍCH HỆ THỐNG EDUCENTER

**Tên hệ thống:** EduCenter  
**Loại tài liệu:** Báo cáo khảo sát và phân tích hệ thống / tài liệu đầu vào BA-SA  
**Nguồn phân tích:** kế thừa từ các gói reverse-engineer trước đó và được chuẩn hóa lại từ mã nguồn tại `/Users/hant/golang/doan`  
**Ngày tổng hợp:** 13/04/2026  

## Quy ước mức độ chắc chắn

- **Confirmed from code**: xác nhận trực tiếp từ entity, migration, DTO, controller, middleware, use case, repository, route hoặc màn hình.
- **Strongly inferred from code**: suy luận mạnh từ cấu trúc dữ liệu, luồng xử lý, tên trường, placeholder UI, comment hoặc side effect.
- **Assumption / Needs BA validation**: giả định hợp lý nhưng code chưa đủ bằng chứng hoặc đang có mâu thuẫn giữa các lớp triển khai.

---

# 1. Giới thiệu

## 1.1 Bối cảnh hệ thống

EduCenter là hệ thống hỗ trợ vận hành cho một trung tâm dạy thêm hoặc trung tâm học thêm. Hệ thống đang được xây dựng theo hướng số hóa các nghiệp vụ cốt lõi của trung tâm:

- quản lý tài khoản và truy cập;
- quản lý học viên, giáo viên, khóa học, chương trình đào tạo;
- mở lớp, ghi danh, phân công giáo viên;
- quản lý phòng học và ca học;
- xếp lịch học thông minh và sinh buổi học thực tế;
- kiểm soát tài liệu giảng dạy;
- chuẩn bị dữ liệu cho các bài toán phân tích nâng cao như dự báo học viên có nguy cơ học kém.

## 1.2 Mục tiêu khảo sát

Mục tiêu của báo cáo là:

1. Xác định hiện trạng nghiệp vụ đang được hiện thực trong codebase.
2. Chuyển đổi logic kỹ thuật thành ngôn ngữ nghiệp vụ phù hợp cho BA, SA, Product Owner và Tech Lead.
3. Làm cơ sở cho:
   - sơ đồ use case,
   - cây phân rã chức năng,
   - ERD,
   - BPMN/activity flow,
   - tài liệu SRS/BRD sơ bộ.
4. Chỉ ra các khoảng trống giữa dữ liệu, API, UI và quy trình vận hành kỳ vọng.

## 1.3 Phạm vi phân tích

Phạm vi phân tích bao gồm:

- backend Go:
  - entity/model,
  - migration,
  - DTO,
  - route/controller,
  - middleware auth/role,
  - use case/service/repository;
- frontend React:
  - route,
  - màn hình admin/teacher/student/compliance,
  - API client,
  - dấu vết vai trò người dùng;
- tài liệu nội bộ đã sinh trước đó:
  - BA reverse-engineer package,
  - use case specification package,
  - ERD drafting package.

## 1.4 Phương pháp tiếp cận từ mã nguồn

Báo cáo sử dụng phương pháp **reverse-engineering business analysis from code**, gồm các bước:

1. Xác định các miền nghiệp vụ từ route, controller và màn hình.
2. Đọc entity/migration để hiểu thực thể và quan hệ dữ liệu.
3. Đọc use case/service để suy ra luồng nghiệp vụ thật.
4. Đối chiếu middleware/role để đánh giá phân quyền thực tế.
5. Đối chiếu dữ liệu hiện hữu với workflow kỳ vọng để tìm khoảng trống.

## 1.5 Giới hạn của tài liệu

- Codebase là nguồn sự thật chính, nhưng không phải mọi ý đồ nghiệp vụ đều đã được hiện thực.
- Một số module có schema dữ liệu nhưng chưa có API hoặc UI, vì vậy phần nghiệp vụ liên quan chỉ có thể đánh dấu ở mức `Strongly inferred from code` hoặc `Assumption / Needs BA validation`.
- Một số phần có mâu thuẫn giữa entity, migration, DTO và UI. Báo cáo sẽ chỉ ra rõ thay vì cố hợp nhất bằng giả định thiếu căn cứ.

---

# 2. Tổng quan nghiệp vụ hệ thống EduCenter

## 2.1 Mô hình hoạt động của trung tâm dạy thêm

Từ góc nhìn nghiệp vụ, EduCenter phản ánh mô hình vận hành của một trung tâm dạy thêm theo chuỗi giá trị sau:

1. Trung tâm quản lý **nguồn lực đầu vào**:
   - học viên,
   - giáo viên,
   - khóa học,
   - chương trình đào tạo,
   - phòng học,
   - ca học.
2. Trung tâm tạo **đơn vị vận hành** là lớp học:
   - lớp gắn với khóa học hoặc chương trình,
   - có giáo viên phụ trách,
   - có sĩ số tối đa,
   - có lịch tuần chuẩn.
3. Học viên được **ghi danh** vào lớp.
4. Hệ thống dùng dữ liệu lớp, ca, phòng, giáo viên để **xếp lịch** và sinh preview.
5. Khi preview đạt điều kiện, hệ thống **commit** để sinh `Lesson` thực tế.
6. Sau mỗi buổi học, mô hình dữ liệu cho thấy trung tâm kỳ vọng sẽ:
   - điểm danh,
   - tạo tổng kết buổi học,
   - ghi nhận academic record,
   - xử lý leave request.
7. Giáo viên có thể **upload tài liệu giảng dạy**; tài liệu được gắn nhãn và phê duyệt.
8. Hệ thống mở hướng tới **lead intake/tư vấn** và **predictive analytics**.

**Mức độ chắc chắn:**  
- phần từ bước 1 đến bước 5: Confirmed from code  
- phần bước 6 và bước 8: Strongly inferred from code  

## 2.2 Giá trị nghiệp vụ mà hệ thống hỗ trợ

| Mã | Giá trị nghiệp vụ | Diễn giải |
|---|---|---|
| BV-01 | Chuẩn hóa dữ liệu vận hành | gom thông tin học viên, giáo viên, khóa học, lớp, phòng, ca vào một mô hình thống nhất |
| BV-02 | Giảm công xếp lịch thủ công | preview scheduling giúp phát hiện xung đột trước khi sinh lesson thật |
| BV-03 | Tăng khả năng kiểm soát lớp học | lớp có sĩ số tối đa, có roster, có giáo viên và lịch tuần |
| BV-04 | Tăng khả năng truy vết vận hành | lesson là sản phẩm vận hành rõ ràng sau scheduling commit |
| BV-05 | Tạo nền dữ liệu cho phân tích học tập | attendance, lesson summary, academic record, leave request cho thấy định hướng theo dõi tiến bộ học tập |
| BV-06 | Tăng kiểm soát nội dung giảng dạy | material moderation cho phép gắn nhãn và duyệt tài liệu |

## 2.3 Các nhóm nghiệp vụ chính

1. Xác thực và quản lý tài khoản.
2. Quản lý danh mục đào tạo và nguồn lực.
3. Quản lý lớp học và ghi danh.
4. Quản lý lịch học và lesson.
5. Quản lý dữ liệu học tập.
6. Kiểm soát tài liệu giảng dạy.
7. Hỗ trợ tuyển sinh và phân tích nâng cao.

## 2.4 Các nhóm người dùng chính

| Nhóm người dùng | Vai trò nghiệp vụ | Mức độ chắc chắn |
|---|---|---|
| Quản trị viên | vận hành hầu hết chức năng lõi | Confirmed from code |
| Giáo viên | cung cấp tài liệu, xem lịch dạy, theo dõi giờ dạy | Confirmed from code |
| Học viên | actor đích cho cổng học viên, điểm danh, đơn xin phép, kết quả học tập | Strongly inferred from code |
| Compliance/reviewer | duyệt tài liệu bị gắn cờ | Strongly inferred from code |
| Nhân viên vận hành / giáo vụ | persona nghiệp vụ của admin hiện trạng | Strongly inferred from code |
| Phụ huynh | actor liên quan gián tiếp qua guardian_phone và leave/learning communication | Assumption / Needs BA validation |

---

# 3. Phạm vi chức năng

## 3.1 Chức năng đã triển khai

| Nhóm | Chức năng | Ghi chú |
|---|---|---|
| Auth | đăng ký, verify OTP, đăng nhập, refresh, logout, quên mật khẩu, reset mật khẩu, đổi mật khẩu | `/v1/auth/*`, một phần `/v2/auth/*` |
| Teacher | CRUD giáo viên, xem lịch dạy, thống kê giờ dạy | backend và UI admin đã có |
| Student | CRUD học viên, tìm kiếm/lọc | có auth gap |
| Course | CRUD khóa học | dùng cho class và scheduling |
| Program | CRUD chương trình, gán/gỡ khóa học | lifecycle publish/archive chưa hoàn chỉnh |
| Room | CRUD phòng học | contract dữ liệu còn mâu thuẫn |
| Shift | CRUD ca học | đã trở thành đầu vào chuẩn cho scheduling |
| Class | CRUD lớp, gán giáo viên | room mặc định ở entity nhưng chưa expose rõ ở DTO |
| Enrollment | thêm/xóa học viên khỏi lớp, xem roster | lifecycle chưa hoàn chỉnh |
| Scheduling | tạo preview, xem latest/by id, commit preview, benchmark contract | benchmark chưa chạy thật |
| Material | upload, list, detail, download, review | moderation đang dùng OCR/AI stub |

## 3.2 Chức năng triển khai một phần

| Nhóm | Hiện trạng | Vấn đề |
|---|---|---|
| Program lifecycle | có field `approved_by_id`, `approval_note`, `published_at`, `archived_at` | chưa có use case publish/archive rõ |
| Scheduling benchmark | có API và contract | chưa có benchmark thực thi đầy đủ 3 solver |
| Lesson management | lesson sinh từ commit | chưa có API/màn hình lesson management độc lập |
| Dashboard/report | có admin overview và teacher stats | overview dùng nhiều dữ liệu mock |
| Material moderation | workflow có chạy | phân quyền người duyệt chưa chặt, dùng stub thay vì AI thật |

## 3.3 Chức năng có dấu hiệu trong dữ liệu nhưng chưa có API

| Nhóm | Entity/Dấu hiệu | Nhận xét |
|---|---|---|
| Objective/Outcome | đã có entity quan hệ với Program | chưa có API/UI |
| Attendance | entity đã có | chưa có API, chưa có enum status |
| LessonSummary | entity đã có | chưa có API/UI |
| AcademicRecord | entity đã có | chưa có API/UI |
| LeaveRequest | entity đã có | chưa có API/UI |
| Consultation | entity đã có | chưa có API/UI |
| ClassSchedule admin CRUD | scheduling dùng dữ liệu này | chưa thấy API quản trị riêng |

## 3.4 Chức năng dự kiến / tương lai

| Nhóm | Dấu vết | Đánh giá |
|---|---|---|
| Predictive analytics học viên nguy cơ học kém | kế hoạch và dữ liệu nền đã có | future-state |
| Học vụ hậu buổi học đầy đủ | attendance + summary + academic record + leave request | target-state khá rõ |
| Benchmark và chọn solver tối ưu | đã có kiến trúc solver service | đang chuyển từ architecture sang operational use case |
| Consultation / lead intake | có entity | chưa thành module hoàn chỉnh |

## 3.5 Chức năng ngoài phạm vi

Không thấy bằng chứng triển khai cho:

- billing / payment / công nợ;
- quản lý chi nhánh / multi-tenant;
- chấm công nhân sự;
- học liệu nâng cao có versioning;
- notification center nội bộ hoàn chỉnh;
- tích hợp LMS bên ngoài.

---

# 4. Danh mục tác nhân và bên liên quan

## 4.1 Actor catalog

| Mã actor | Tên actor | Loại | Mục tiêu chính | Trách nhiệm | Mức độ chắc chắn |
|---|---|---|---|---|---|
| ACT-01 | Quản trị viên | Nội bộ | vận hành toàn bộ trung tâm | quản lý master data, lớp, scheduling, material review | Confirmed from code |
| ACT-02 | Giáo viên | Nội bộ | giảng dạy và cung cấp học liệu | upload material, xem lịch dạy, theo dõi giờ dạy | Confirmed from code |
| ACT-03 | Học viên | Nội bộ | tham gia học tập | dùng portal học viên, gửi leave request, xem kết quả học tập trong tương lai | Strongly inferred from code |
| ACT-04 | Người dùng chưa đăng nhập | Bên ngoài | đăng ký và khôi phục truy cập | register, verify OTP, login, forgot/reset | Confirmed from code |
| ACT-05 | Reviewer / compliance officer | Nội bộ | kiểm duyệt tài liệu giảng dạy | approve/reject material | Strongly inferred from code |
| ACT-06 | Nhân viên giáo vụ / vận hành | Nội bộ | vận hành học vụ hàng ngày | quản lý lớp, roster, shift, room, scheduling | Strongly inferred from code |
| ACT-07 | Quản lý đào tạo | Nội bộ | quản trị chương trình đào tạo | quản lý program, objective, outcome, lịch học | Strongly inferred from code |
| ACT-08 | Phụ huynh | Bên ngoài liên quan nghiệp vụ | hỗ trợ quản lý học viên | liên hệ, theo dõi, xin phép, nhận kết quả | Assumption / Needs BA validation |
| ACT-09 | SMTP | Hệ thống | gửi mail hệ thống | gửi OTP, gửi reset password | Confirmed from code |
| ACT-10 | Scheduling engine | Hệ thống | sinh preview và lesson | chạy solver, đánh giá conflict, commit lesson | Confirmed from code |
| ACT-11 | OCR stub | Hệ thống | trích xuất văn bản từ tài liệu | tạo raw OCR text | Confirmed from code |
| ACT-12 | AI moderation stub | Hệ thống | gán nhãn tài liệu | SAFE/WARNING/DANGER | Confirmed from code |

## 4.2 Stakeholder catalog

| Mã stakeholder | Nhóm | Quan tâm chính |
|---|---|---|
| STK-01 | Product Owner | phạm vi sản phẩm, thứ tự ưu tiên module |
| STK-02 | BA | quy trình nghiệp vụ, actor, use case, rule, lifecycle |
| STK-03 | Tech Lead / SA | kiến trúc service, domain boundary, dữ liệu và API |
| STK-04 | Vận hành trung tâm | roster, lớp, lịch, lesson, tình trạng học viên |
| STK-05 | Giáo viên | lịch dạy, tài liệu, đánh giá học tập |
| STK-06 | Học viên / phụ huynh | thời khóa biểu, nghỉ phép, kết quả học tập |
| STK-07 | Reviewer nội dung | danh sách tài liệu cần duyệt, quyết định duyệt/từ chối |

## 4.3 Vai trò và trách nhiệm

| Actor | Trách nhiệm nghiệp vụ |
|---|---|
| Admin | tạo và quản trị danh mục; mở lớp; ghi danh; phân công giáo viên; chạy scheduling; commit lesson |
| Giáo viên | tham gia giảng dạy; cung cấp tài liệu; trong tương lai tạo lesson summary và academic record |
| Học viên | là đối tượng ghi danh, điểm danh, đánh giá học tập, leave request |
| Reviewer | duyệt material sau khi AI/audit gắn nhãn |
| SMTP/OCR/AI/Scheduling | actor hệ thống, không phải người dùng, nhưng là swimlane quan trọng khi vẽ BPMN |

## 4.4 Quan hệ giữa tác nhân và phân hệ

| Phân hệ | Actor chính | Actor phụ |
|---|---|---|
| Xác thực | Guest, User | SMTP |
| Học viên | Admin | Student, Parent |
| Giáo viên | Admin | Teacher |
| Khóa học / chương trình | Admin, Quản lý đào tạo |  |
| Lớp / ghi danh | Admin, Giáo vụ | Teacher, Student |
| Scheduling | Admin | Scheduling engine |
| Học vụ sau lesson | Teacher, Admin | Student |
| Material moderation | Teacher, Reviewer | OCR stub, AI stub |
| Consultation | Tư vấn viên, Guest | Parent, Student |

---

# 5. Phân rã chức năng hệ thống

## 5.1 Cây phân rã chức năng 4 tầng

### FH-01. Quản lý truy cập và danh tính

- `FH-01.01` Đăng ký tài khoản
- `FH-01.02` Xác minh email OTP
- `FH-01.03` Đăng nhập
- `FH-01.04` Refresh token
- `FH-01.05` Đăng xuất
- `FH-01.06` Quên mật khẩu
- `FH-01.07` Đặt lại mật khẩu
- `FH-01.08` Đổi mật khẩu
- `FH-01.09` Xem hồ sơ người dùng hiện tại

### FH-02. Quản lý danh mục đào tạo và nguồn lực

- `FH-02.01` Quản lý học viên
  - `FH-02.01.01` Tạo học viên
  - `FH-02.01.02` Tìm kiếm / xem danh sách học viên
  - `FH-02.01.03` Cập nhật học viên
  - `FH-02.01.04` Xóa học viên
- `FH-02.02` Quản lý giáo viên
  - `FH-02.02.01` Tạo giáo viên
  - `FH-02.02.02` Cập nhật giáo viên
  - `FH-02.02.03` Xóa giáo viên
  - `FH-02.02.04` Xem lịch dạy giáo viên
  - `FH-02.02.05` Xem thống kê giờ dạy
- `FH-02.03` Quản lý khóa học
  - `FH-02.03.01` Tạo khóa học
  - `FH-02.03.02` Cập nhật khóa học
  - `FH-02.03.03` Xóa khóa học
- `FH-02.04` Quản lý chương trình đào tạo
  - `FH-02.04.01` Tạo chương trình
  - `FH-02.04.02` Cập nhật chương trình
  - `FH-02.04.03` Gán khóa học vào chương trình
  - `FH-02.04.04` Gỡ khóa học khỏi chương trình
  - `FH-02.04.05` Xuất bản / lưu trữ chương trình
- `FH-02.05` Quản lý phòng học
  - `FH-02.05.01` Tạo phòng
  - `FH-02.05.02` Cập nhật phòng
  - `FH-02.05.03` Xóa phòng
- `FH-02.06` Quản lý ca học
  - `FH-02.06.01` Tạo ca học
  - `FH-02.06.02` Cập nhật ca học
  - `FH-02.06.03` Xóa ca học
  - `FH-02.06.04` Kích hoạt / vô hiệu ca học

### FH-03. Quản lý lớp học và ghi danh

- `FH-03.01` Quản lý lớp học
  - `FH-03.01.01` Tạo lớp
  - `FH-03.01.02` Cập nhật lớp
  - `FH-03.01.03` Đóng / hủy lớp
- `FH-03.02` Quản lý giáo viên của lớp
  - `FH-03.02.01` Gán giáo viên phụ trách
- `FH-03.03` Quản lý roster
  - `FH-03.03.01` Ghi danh học viên vào lớp
  - `FH-03.03.02` Rút học viên khỏi lớp
  - `FH-03.03.03` Xem roster lớp
- `FH-03.04` Quản lý lịch tuần lớp
  - `FH-03.04.01` Tạo slot lịch tuần
  - `FH-03.04.02` Cập nhật slot lịch tuần
  - `FH-03.04.03` Gỡ slot lịch tuần

### FH-04. Xếp lịch và lesson

- `FH-04.01` Chuẩn bị dữ liệu xếp lịch
  - `FH-04.01.01` Chọn phạm vi lớp
  - `FH-04.01.02` Chọn tập giáo viên
  - `FH-04.01.03` Chọn tập phòng
  - `FH-04.01.04` Chọn khoảng ngày
- `FH-04.02` Tạo preview xếp lịch
  - `FH-04.02.01` Xây dựng bài toán xếp lịch
  - `FH-04.02.02` Chạy solver
  - `FH-04.02.03` Tổng hợp conflict
- `FH-04.03` Xem và đánh giá preview
  - `FH-04.03.01` Xem assignments
  - `FH-04.03.02` Xem conflicts
  - `FH-04.03.03` Xem summary
- `FH-04.04` Commit preview
  - `FH-04.04.01` Kiểm tra điều kiện commit
  - `FH-04.04.02` Sinh lesson
- `FH-04.05` Benchmark solver
  - `FH-04.05.01` Chạy benchmark nội bộ
  - `FH-04.05.02` So sánh solver

### FH-05. Vận hành học tập sau buổi học

- `FH-05.01` Quản lý lesson
  - `FH-05.01.01` Xem danh sách lesson
  - `FH-05.01.02` Xem chi tiết lesson
- `FH-05.02` Điểm danh
  - `FH-05.02.01` Chấm điểm danh
  - `FH-05.02.02` Cập nhật trạng thái chuyên cần
- `FH-05.03` Tổng kết buổi học
  - `FH-05.03.01` Ghi chủ đề
  - `FH-05.03.02` Ghi nội dung dạy
  - `FH-05.03.03` Giao bài tập
- `FH-05.04` Ghi nhận kết quả học tập
  - `FH-05.04.01` Đánh giá bài tập
  - `FH-05.04.02` Đánh giá thái độ / tham gia
  - `FH-05.04.03` Chốt academic record
- `FH-05.05` Xử lý đơn xin phép
  - `FH-05.05.01` Tạo đơn
  - `FH-05.05.02` Duyệt / từ chối đơn

### FH-06. Quản lý tài liệu giảng dạy và kiểm duyệt

- `FH-06.01` Upload tài liệu
- `FH-06.02` Lưu trữ file và metadata
- `FH-06.03` OCR / AI audit
- `FH-06.04` Xem hàng chờ duyệt
- `FH-06.05` Duyệt / từ chối tài liệu
- `FH-06.06` Tải tài liệu

### FH-07. Tuyển sinh và phân tích nâng cao

- `FH-07.01` Tiếp nhận consultation / lead
- `FH-07.02` Theo dõi trạng thái lead
- `FH-07.03` Dự báo học viên nguy cơ học kém

## 5.2 Diễn giải từng nhánh chính

| Nhánh | Diễn giải |
|---|---|
| FH-01 | bảo vệ cổng vào hệ thống và định danh actor |
| FH-02 | chuẩn hóa danh mục nền cho mọi hoạt động đào tạo |
| FH-03 | hình thành đơn vị vận hành lớp học và roster |
| FH-04 | chuyển dữ liệu lớp thành lịch học thực thi |
| FH-05 | theo dõi quá trình học tập sau khi lesson đã tồn tại |
| FH-06 | kiểm soát học liệu được sử dụng trong giảng dạy |
| FH-07 | mở rộng sang tuyển sinh và phân tích dự báo |

## 5.3 Chức năng trọng yếu

Các chức năng trọng yếu nhất theo hiện trạng code và giá trị vận hành là:

1. Quản lý lớp học.
2. Ghi danh học viên.
3. Quản lý ca học và phòng học.
4. Tạo preview scheduling.
5. Commit preview thành lesson.
6. Upload và duyệt material.

---

# 6. Danh mục use case

## 6.1 Danh sách use case theo nhóm tác nhân

### Nhóm Guest / người dùng chưa đăng nhập

- UC-01 Đăng ký tài khoản
- UC-02 Xác minh email OTP
- UC-03 Đăng nhập
- UC-04 Quên mật khẩu
- UC-05 Đặt lại mật khẩu

### Nhóm Admin / vận hành

- UC-06 Xem hồ sơ hiện tại
- UC-07 Tạo học viên
- UC-08 Tìm kiếm / xem danh sách học viên
- UC-09 Cập nhật học viên
- UC-10 Xóa học viên
- UC-11 Tạo giáo viên
- UC-12 Cập nhật giáo viên
- UC-13 Xóa giáo viên
- UC-14 Xem lịch dạy giáo viên
- UC-15 Xem thống kê giờ dạy
- UC-16 Tạo khóa học
- UC-17 Cập nhật khóa học
- UC-18 Tạo chương trình
- UC-19 Gán khóa học vào chương trình
- UC-20 Cập nhật / xuất bản / lưu trữ chương trình
- UC-21 Tạo lớp học
- UC-22 Cập nhật lớp học
- UC-23 Ghi danh học viên vào lớp
- UC-24 Rút học viên khỏi lớp
- UC-25 Phân công giáo viên cho lớp
- UC-26 Cấu hình lịch tuần cho lớp
- UC-27 Tạo preview xếp lịch
- UC-28 Xem preview xếp lịch
- UC-29 Benchmark solver
- UC-30 Xác nhận preview để tạo lesson
- UC-31 Xem lesson
- UC-32 Duyệt / từ chối tài liệu

### Nhóm Teacher

- UC-33 Đăng nhập / refresh / đổi mật khẩu
- UC-34 Xem hồ sơ hiện tại
- UC-35 Xem lịch dạy
- UC-36 Xem thống kê giờ dạy
- UC-37 Upload tài liệu
- UC-38 Tải tài liệu
- UC-39 Tạo điểm danh
- UC-40 Tạo lesson summary
- UC-41 Ghi nhận academic record

### Nhóm Student

- UC-42 Đăng nhập / refresh / đổi mật khẩu
- UC-43 Xem hồ sơ hiện tại
- UC-44 Xem thông tin học tập
- UC-45 Tạo đơn xin phép
- UC-46 Xem kết quả học tập

### Nhóm Reviewer / Compliance

- UC-47 Xem hàng chờ tài liệu
- UC-48 Xem chi tiết tài liệu
- UC-49 Duyệt / từ chối tài liệu

### Nhóm Hệ thống

- UC-50 Gửi OTP / reset password qua SMTP
- UC-51 OCR tài liệu
- UC-52 Gắn nhãn AI
- UC-53 Chạy solver scheduling
- UC-54 Commit lesson từ preview

## 6.2 Danh sách use case theo miền nghiệp vụ

| Miền | Use case chính |
|---|---|
| Xác thực | UC-01 đến UC-06 |
| Học viên | UC-07 đến UC-10 |
| Giáo viên | UC-11 đến UC-15 |
| Khóa học / chương trình | UC-16 đến UC-20 |
| Lớp học / ghi danh | UC-21 đến UC-26 |
| Scheduling | UC-27 đến UC-30 |
| Lesson và học vụ | UC-31, UC-39, UC-40, UC-41, UC-45, UC-46 |
| Material moderation | UC-37, UC-38, UC-47, UC-48, UC-49 |
| Hệ thống nền | UC-50 đến UC-54 |

## 6.3 Top use case ưu tiên cao

1. UC-01 Đăng ký tài khoản
2. UC-02 Xác minh email OTP
3. UC-03 Đăng nhập
4. UC-07 Tạo học viên
5. UC-11 Tạo giáo viên
6. UC-16 Tạo khóa học
7. UC-18 Tạo chương trình
8. UC-19 Gán khóa học vào chương trình
9. UC-21 Tạo lớp học
10. UC-23 Ghi danh học viên vào lớp
11. UC-25 Phân công giáo viên cho lớp
12. UC-26 Cấu hình lịch tuần cho lớp
13. UC-27 Tạo preview xếp lịch
14. UC-28 Xem preview xếp lịch
15. UC-30 Xác nhận preview để tạo lesson
16. UC-37 Upload tài liệu
17. UC-49 Duyệt / từ chối tài liệu

---

# 7. Đặc tả các use case trọng yếu

## 7.1 Xác thực tài khoản

### UC-AUTH-01. Đăng ký tài khoản

- **Mục tiêu:** tạo tài khoản mới và khởi tạo luồng xác minh email.
- **Tác nhân:** Guest.
- **Tiền điều kiện:** email chưa tồn tại.
- **Luồng chính:**
  1. Người dùng nhập email, họ tên và mật khẩu.
  2. Frontend gửi `password_enc`.
  3. Backend kiểm tra email trùng.
  4. Backend giải mã và băm mật khẩu.
  5. Backend tạo `User` với `is_active = false`.
  6. Backend tạo `UserOTP`.
  7. Hệ thống gửi OTP qua email.
  8. Trả về `user_id`.
- **Luồng thay thế:**
  - môi trường dev không gửi mail thật mà log mail.
  - email test đặc biệt dùng OTP cố định để phục vụ demo.
- **Luồng ngoại lệ:**
  - email đã tồn tại;
  - lỗi giải mã mật khẩu;
  - lỗi transaction DB;
  - lỗi gửi mail sau khi đã tạo user.
- **Dữ liệu liên quan:** `User`, `UserOTP`.
- **Quy tắc nghiệp vụ:** email phải unique; user chưa active thì chưa được login.
- **Ghi chú:** Confirmed from code.

### UC-AUTH-02. Xác minh email OTP

- **Mục tiêu:** kích hoạt tài khoản.
- **Tác nhân:** Guest/User vừa đăng ký.
- **Tiền điều kiện:** user đã có OTP active.
- **Luồng chính:**
  1. Người dùng nhập OTP.
  2. Backend lấy OTP gần nhất còn hiệu lực cho user.
  3. So khớp hash OTP.
  4. Đánh dấu `used_at`.
  5. Cập nhật `User.is_active = true`.
- **Luồng thay thế:** có thể nhập lại OTP khác nếu OTP cũ sai nhưng còn hạn.
- **Luồng ngoại lệ:** OTP hết hạn, OTP sai, OTP đã dùng, user không tồn tại.
- **Dữ liệu:** `UserOTP`, `User`.
- **Quy tắc:** OTP chỉ dùng một lần; OTP phải còn hạn.
- **Ghi chú:** Confirmed from code.

### UC-AUTH-03. Đăng nhập

- **Mục tiêu:** cấp access token và refresh token cho user hợp lệ.
- **Tác nhân:** User.
- **Tiền điều kiện:** user tồn tại, mật khẩu đúng, account active.
- **Luồng chính:**
  1. Người dùng nhập email và mật khẩu.
  2. Backend tìm user theo email.
  3. So khớp bcrypt.
  4. Kiểm tra `is_active`.
  5. Sinh JWT access token và refresh token.
  6. Trả profile + token.
- **Luồng thay thế:** dùng `v2/auth/login` với shape response hơi khác.
- **Luồng ngoại lệ:** email không tồn tại, sai mật khẩu, account inactive.
- **Dữ liệu:** `User`.
- **Quy tắc:** chỉ user active mới login được.
- **Ghi chú:** Confirmed from code.

## 7.2 Quản lý học viên

### UC-STU-01. Tạo học viên

- **Mục tiêu:** tạo hồ sơ học viên để có thể ghi danh.
- **Tác nhân:** Admin/vận hành.
- **Tiền điều kiện:** đã đăng nhập.
- **Luồng chính:**
  1. Admin mở màn quản lý học viên.
  2. Nhập mã, họ tên, email, điện thoại, SĐT phụ huynh, khối lớp, trạng thái.
  3. Backend tạo `Student`.
  4. UI refresh danh sách.
- **Luồng thay thế:** có thể tạo với nhiều trường tùy chọn rỗng.
- **Luồng ngoại lệ:** payload không hợp lệ, lỗi DB, mã học viên trùng nếu DB constraint bị vi phạm.
- **Dữ liệu:** `Student`.
- **Quy tắc:** học viên là master data cho enrollment.
- **Ghi chú:** chức năng có thật, nhưng phân quyền hiện chỉ dùng auth, chưa khóa `ADMIN`. Confirmed from code.

### UC-STU-02. Tìm kiếm / xem danh sách học viên

- **Mục tiêu:** tra cứu học viên phục vụ vận hành lớp.
- **Tác nhân:** Admin/vận hành.
- **Tiền điều kiện:** đã đăng nhập.
- **Luồng chính:**
  1. Người dùng vào trang học viên.
  2. Nhập từ khóa hoặc dùng bộ lọc.
  3. Hệ thống trả danh sách có phân trang.
- **Luồng thay thế:** xem chi tiết bằng dialog.
- **Luồng ngoại lệ:** lỗi gọi API.
- **Dữ liệu:** `Student`.
- **Quy tắc:** search by code/name/email/phone tùy implementation.
- **Ghi chú:** Confirmed from code.

## 7.3 Quản lý giáo viên

### UC-TCH-01. Tạo giáo viên

- **Mục tiêu:** tạo hồ sơ giáo viên để gán lớp và xếp lịch.
- **Tác nhân:** Admin.
- **Tiền điều kiện:** tài khoản admin hợp lệ.
- **Luồng chính:**
  1. Admin nhập mã, họ tên, email, phone, loại làm việc, trạng thái, ghi chú.
  2. Backend kiểm tra unique code/email nếu có.
  3. Tạo `Teacher` với default hợp lý nếu thiếu.
  4. UI refresh danh sách.
- **Luồng thay thế:** teacher có thể là giáo viên trường với `is_school_teacher = true`.
- **Luồng ngoại lệ:** email/code trùng, lỗi DB.
- **Dữ liệu:** `Teacher`.
- **Quy tắc:** teacher là nguồn lực bắt buộc cho scheduling.
- **Ghi chú:** Confirmed from code.

### UC-TCH-02. Xem lịch dạy giáo viên

- **Mục tiêu:** tra cứu lesson theo giáo viên.
- **Tác nhân:** Admin, Teacher.
- **Tiền điều kiện:** có dữ liệu lesson.
- **Luồng chính:**
  1. Người dùng chọn giáo viên.
  2. Hệ thống truy vấn lesson theo `teacher_id`.
  3. Trả danh sách lesson theo thời gian.
- **Luồng thay thế:** giáo viên xem chính lịch dạy của mình.
- **Luồng ngoại lệ:** chưa có lesson nào.
- **Dữ liệu:** `Lesson`, `Teacher`.
- **Quy tắc:** lesson là snapshot đã commit, không phải preview.
- **Ghi chú:** Confirmed from code.

## 7.4 Quản lý khóa học/chương trình

### UC-CRS-01. Tạo khóa học

- **Mục tiêu:** tạo đơn vị đào tạo cơ sở.
- **Tác nhân:** Admin.
- **Tiền điều kiện:** admin đã đăng nhập.
- **Luồng chính:**
  1. Admin nhập mã khóa học, tên, mô tả, môn học, khối, số buổi, thời lượng, giá.
  2. Backend tạo `Course`.
  3. UI refresh danh sách.
- **Luồng ngoại lệ:** mã khóa học trùng, lỗi DB.
- **Dữ liệu:** `Course`.
- **Quy tắc:** `session_count` và `session_duration_minutes` là đầu vào quan trọng cho scheduling.
- **Ghi chú:** Confirmed from code.

### UC-PRG-01. Tạo chương trình

- **Mục tiêu:** tạo khung chương trình đào tạo.
- **Tác nhân:** Admin/quản lý đào tạo.
- **Tiền điều kiện:** đã đăng nhập.
- **Luồng chính:**
  1. Nhập mã chương trình, tên, track, khoảng hiệu lực.
  2. Backend tạo `Program`.
  3. UI refresh danh sách.
- **Luồng thay thế:** có thể để trống khoảng hiệu lực.
- **Luồng ngoại lệ:** mã chương trình trùng.
- **Dữ liệu:** `Program`.
- **Quy tắc:** program có thể chưa có course ngay lúc tạo.
- **Ghi chú:** Confirmed from code.

### UC-PRG-02. Gán khóa học vào chương trình

- **Mục tiêu:** xây chương trình từ nhiều khóa học.
- **Tác nhân:** Admin/quản lý đào tạo.
- **Tiền điều kiện:** program và course tồn tại.
- **Luồng chính:**
  1. Mở chi tiết program.
  2. Chọn danh sách course khả dụng.
  3. Gọi mutation add/remove mapping.
  4. Hệ thống cập nhật `ProgramCourse`.
- **Luồng thay thế:** gỡ khóa học khỏi chương trình.
- **Luồng ngoại lệ:** mapping trùng, lỗi DB.
- **Dữ liệu:** `Program`, `Course`, `ProgramCourse`.
- **Quy tắc:** một program có thể có nhiều course.
- **Ghi chú:** Confirmed from code.

## 7.5 Quản lý lớp và ghi danh

### UC-CLS-01. Tạo lớp học

- **Mục tiêu:** tạo đơn vị vận hành lớp thực tế.
- **Tác nhân:** Admin/vận hành.
- **Tiền điều kiện:** có course/program/teacher nếu muốn gán ngay.
- **Luồng chính:**
  1. Nhập mã lớp, tên lớp, ngày bắt đầu, sĩ số tối đa, giá, ghi chú.
  2. Chọn program, course, teacher nếu có.
  3. Backend tạo `Class`.
  4. UI refresh danh sách.
- **Luồng thay thế:** có thể tạo lớp trước, gán course/teacher sau.
- **Luồng ngoại lệ:** mã lớp trùng, max_students không hợp lệ.
- **Dữ liệu:** `Class`.
- **Quy tắc:** scheduling chỉ xử lý class `OPEN`.
- **Ghi chú:** Confirmed from code.

### UC-ENR-01. Ghi danh học viên vào lớp

- **Mục tiêu:** thêm học viên vào roster lớp.
- **Tác nhân:** Admin/vận hành.
- **Tiền điều kiện:** lớp tồn tại; học viên tồn tại.
- **Luồng chính:**
  1. Mở chi tiết lớp.
  2. Chọn học viên từ danh sách khả dụng.
  3. Backend tạo các `Enrollment`.
  4. UI refresh roster.
- **Luồng thay thế:** thêm nhiều học viên một lần.
- **Luồng ngoại lệ:**
  - vượt quá max_students;
  - phòng mặc định có sức chứa thấp hơn giới hạn lớp;
  - lỗi DB.
- **Dữ liệu:** `Enrollment`, `Class`, `Student`.
- **Quy tắc:** use case hiện set `ENROLLED` trực tiếp, chưa đi qua `APPLIED`.
- **Ghi chú:** Confirmed from code; lifecycle enrollment còn chưa hoàn chỉnh.

### UC-CLS-03. Phân công giáo viên cho lớp

- **Mục tiêu:** gắn giáo viên phụ trách chính cho lớp.
- **Tác nhân:** Admin.
- **Tiền điều kiện:** class và teacher tồn tại.
- **Luồng chính:**
  1. Admin chọn lớp.
  2. Chọn giáo viên.
  3. Backend cập nhật `class.teacher_id`.
  4. UI refresh chi tiết lớp.
- **Luồng ngoại lệ:** teacher không tồn tại hoặc lỗi DB.
- **Dữ liệu:** `Class`, `Teacher`.
- **Quy tắc:** thiếu giáo viên sẽ dẫn đến conflict scheduling.
- **Ghi chú:** Confirmed from code.

### UC-SCHCFG-01. Cấu hình lịch tuần cho lớp

- **Mục tiêu:** xác định thứ học và ca học định kỳ của lớp.
- **Tác nhân:** Admin/vận hành.
- **Tiền điều kiện:** class và shift tồn tại.
- **Luồng chính mong muốn:**
  1. Chọn lớp.
  2. Chọn thứ trong tuần.
  3. Chọn shift.
  4. Chọn phòng cố định nếu cần.
  5. Tạo `ClassSchedule`.
- **Luồng thay thế:** một lớp có nhiều slot lịch tuần.
- **Luồng ngoại lệ:** shift không hợp lệ, phòng không hợp lệ, trùng slot.
- **Dữ liệu:** `ClassSchedule`, `Shift`, `Room`, `Class`.
- **Quy tắc:** scheduling mới đọc `shift_id` từ class_schedule.
- **Ghi chú:** Strongly inferred from code. Đây là use case quan trọng nhưng thiếu API quản trị rõ.

## 7.6 Xếp lịch và commit lesson

### UC-SOL-01. Tạo preview xếp lịch

- **Mục tiêu:** tạo phương án lịch học thử trước khi commit lesson.
- **Tác nhân:** Admin.
- **Tiền điều kiện:**
  - có class `OPEN`;
  - có shift active;
  - date range hợp lệ.
- **Luồng chính:**
  1. Admin chọn `date_from`, `date_to`.
  2. Chọn danh sách class, teacher, room nếu cần giới hạn phạm vi.
  3. Backend tải classes, class schedules, teacher, room, shift.
  4. Hệ thống xây dựng bài toán scheduling.
  5. Solver chạy và trả assignments + conflicts + summary.
  6. Kết quả lưu vào preview store in-memory.
  7. UI hiển thị preview.
- **Luồng thay thế:** chọn tập con class/teacher/room.
- **Luồng ngoại lệ:**
  - class thiếu course;
  - class thiếu teacher;
  - class không có class_schedule;
  - room không đủ sức chứa;
  - ngày không hợp lệ;
  - lỗi solver.
- **Dữ liệu:** `Class`, `ClassSchedule`, `Shift`, `Room`, `Course`, `Teacher`, preview object.
- **Quy tắc:**
  - chỉ dùng class `OPEN`;
  - chỉ dùng shift active;
  - phải tôn trọng lịch tuần lớp;
  - solver phải tránh trùng lớp, giáo viên, phòng;
  - kết quả có thể là `FAILED`, `PARTIAL`, `COMPLETED`.
- **Ghi chú:** Confirmed from code.

### UC-SOL-04. Xác nhận preview để tạo lesson

- **Mục tiêu:** biến preview hợp lệ thành lesson thực tế.
- **Tác nhân:** Admin.
- **Tiền điều kiện:** preview tồn tại và status `COMPLETED`.
- **Luồng chính:**
  1. Admin mở preview.
  2. Bấm commit.
  3. Backend tải preview từ memory store.
  4. Kiểm tra preview có assignments và không có hard conflicts.
  5. Kiểm tra không trùng lesson đã tồn tại trong DB.
  6. Sinh `Lesson` bằng transaction.
  7. Trả kết quả thành công.
- **Luồng thay thế:** nếu preview `PARTIAL`, admin xem conflict rồi chạy lại preview.
- **Luồng ngoại lệ:**
  - preview không tồn tại;
  - preview không `COMPLETED`;
  - preview không có assignments;
  - DB phát hiện overlap lesson;
  - lỗi transaction.
- **Dữ liệu:** preview object, `Lesson`.
- **Quy tắc:** chỉ preview đạt điều kiện mới được commit.
- **Ghi chú:** Confirmed from code.

## 7.7 Điểm danh và tổng kết buổi học

### UC-ATD-01. Điểm danh

- **Mục tiêu:** ghi nhận tình trạng tham dự của học viên theo lesson.
- **Tác nhân:** Teacher/Admin.
- **Tiền điều kiện:** lesson đã tồn tại; roster lớp đã xác định.
- **Luồng chính mong muốn:**
  1. Mở lesson.
  2. Hệ thống nạp danh sách học viên của lớp.
  3. Giáo viên chọn trạng thái attendance cho từng học viên.
  4. Lưu `Attendance`.
- **Luồng thay thế:** cập nhật attendance sau khi đã chấm lần đầu.
- **Luồng ngoại lệ:** lesson không tồn tại, student không thuộc lớp, trạng thái attendance không hợp lệ.
- **Dữ liệu:** `Attendance`, `Lesson`, `Enrollment`, `Student`.
- **Quy tắc:** cần enum status chính thức.
- **Ghi chú:** Strongly inferred from code; chưa có API thực.

### UC-SUM-01. Tạo tổng kết buổi học

- **Mục tiêu:** ghi lại nội dung dạy, bài tập, feedback.
- **Tác nhân:** Teacher.
- **Tiền điều kiện:** lesson đã tồn tại.
- **Luồng chính mong muốn:**
  1. Chọn lesson.
  2. Nhập topic, nội dung đã dạy, feedback, homework, deadline, notes.
  3. Lưu `LessonSummary`.
- **Luồng ngoại lệ:** lesson đã có summary; dữ liệu không hợp lệ.
- **Dữ liệu:** `LessonSummary`, `Lesson`.
- **Quy tắc:** một lesson chỉ có một summary.
- **Ghi chú:** schema đã có, API chưa có.

## 7.8 Kết quả học tập

### UC-ACR-01. Ghi nhận academic record

- **Mục tiêu:** đánh giá kết quả học tập của từng học viên theo lesson summary.
- **Tác nhân:** Teacher.
- **Tiền điều kiện:** lesson summary đã tồn tại.
- **Luồng chính mong muốn:**
  1. Chọn lesson summary.
  2. Tải danh sách học viên của lớp.
  3. Nhập homework score, participation, attitude, comment.
  4. Hệ thống tính hoặc nhận `total_score`.
  5. Lưu `AcademicRecord`.
- **Luồng thay thế:** lưu nháp trước khi chốt `is_completed = true`.
- **Luồng ngoại lệ:** học viên không thuộc lớp, summary không tồn tại, total score không hợp lệ.
- **Dữ liệu:** `AcademicRecord`, `LessonSummary`, `Student`.
- **Quy tắc:** công thức `total_score` chưa được định nghĩa trong code.
- **Ghi chú:** Strongly inferred from code.

## 7.9 Đơn xin phép

### UC-LVE-01. Tạo đơn xin phép

- **Mục tiêu:** cho phép học viên xin nghỉ, xin đi muộn hoặc về sớm.
- **Tác nhân:** Student.
- **Tiền điều kiện:** học viên tồn tại.
- **Luồng chính mong muốn:**
  1. Học viên chọn loại đơn.
  2. Nhập ngày áp dụng, lý do, phút đi muộn/về sớm nếu cần.
  3. Chọn lớp hoặc lesson liên quan nếu có.
  4. Đính kèm tài liệu.
  5. Hệ thống lưu `LeaveRequest` với `PENDING`.
- **Luồng ngoại lệ:** dữ liệu không hợp lệ, thiếu lý do, loại đơn không phù hợp.
- **Dữ liệu:** `LeaveRequest`, `Student`, `Lesson`, `Class`.
- **Quy tắc:** `late_minutes` chỉ hợp lệ cho loại `LATE`, `early_minutes` chỉ hợp lệ cho `EARLY`.
- **Ghi chú:** Strongly inferred from code.

### UC-LVE-02. Duyệt / từ chối đơn xin phép

- **Mục tiêu:** quyết định chấp thuận hay từ chối đơn.
- **Tác nhân:** Admin/giáo vụ.
- **Tiền điều kiện:** tồn tại leave request `PENDING`.
- **Luồng chính mong muốn:**
  1. Người duyệt xem đơn.
  2. Quyết định approve hoặc reject.
  3. Hệ thống cập nhật `status`, `approved_by_id`, `approved_at` hoặc `rejection_reason`.
- **Luồng ngoại lệ:** đơn không tồn tại, đơn đã xử lý.
- **Dữ liệu:** `LeaveRequest`, `User`.
- **Quy tắc:** chỉ đơn `PENDING` mới được xử lý.
- **Ghi chú:** Strongly inferred from code.

## 7.10 Tài liệu giảng dạy và kiểm duyệt

### UC-MAT-01. Upload tài liệu

- **Mục tiêu:** cho phép giáo viên tải tài liệu lên hệ thống.
- **Tác nhân:** Teacher.
- **Tiền điều kiện:** teacher đã đăng nhập.
- **Luồng chính:**
  1. Giáo viên nhập title/description.
  2. Chọn file hợp lệ.
  3. Backend validate loại file và dung lượng.
  4. Lưu file vào local structured storage.
  5. Tạo `Material` status `UPLOADED`.
  6. Chạy OCR stub và AI moderation stub.
  7. Tạo `AuditLog`.
  8. Gắn `latest_label_id`.
  9. Cập nhật material sang `AI_REVIEWED`.
- **Luồng thay thế:** nếu OCR/AI stub phát hiện an toàn, material vẫn có thể đi tiếp tới review.
- **Luồng ngoại lệ:** file quá lớn, loại file không hợp lệ, lỗi storage, lỗi DB.
- **Dữ liệu:** `Material`, `AuditLog`, `Label`.
- **Quy tắc:** chỉ chấp nhận `pdf/doc/docx/png/jpg/jpeg`, tối đa 10MB.
- **Ghi chú:** Confirmed from code.

### UC-MAT-02. Duyệt tài liệu bị gắn cờ

- **Mục tiêu:** quyết định chấp thuận hay từ chối material sau AI review.
- **Tác nhân:** Reviewer/Compliance/Admin.
- **Tiền điều kiện:** material tồn tại.
- **Luồng chính:**
  1. Người duyệt mở chi tiết material.
  2. Xem file metadata, label, audit reasoning.
  3. Chọn approve hoặc reject.
  4. Hệ thống tạo `ApprovalDecision`.
  5. Cập nhật `Material.status` tương ứng.
- **Luồng thay thế:** reviewer tải file xuống trước khi quyết định.
- **Luồng ngoại lệ:** material không tồn tại, lỗi DB.
- **Dữ liệu:** `Material`, `ApprovalDecision`, `AuditLog`.
- **Quy tắc:** reviewer cần có quyền phù hợp.
- **Ghi chú:** workflow có thật nhưng authorization còn lỏng; compliance officer id hiện lấy từ body thay vì JWT. Confirmed from code.

---

# 8. Luồng nghiệp vụ đầu-cuối

## WF-01. Đăng ký và kích hoạt tài khoản

**Mục tiêu:** đưa người dùng mới từ trạng thái chưa có tài khoản sang tài khoản active.

**Luồng chuẩn:**
1. Guest đăng ký tài khoản.
2. Hệ thống tạo `User` inactive và `UserOTP`.
3. SMTP gửi OTP.
4. Guest nhập OTP.
5. Hệ thống verify OTP và activate account.
6. User có thể login.

**Nhánh thay thế:**
- mail không gửi được nhưng tài khoản và OTP đã tạo;
- người dùng dùng OTP test trong môi trường dev.

**Nhánh lỗi:**
- email trùng;
- OTP sai/hết hạn/đã dùng.

## WF-02. Mở lớp và chuẩn bị cho xếp lịch

**Mục tiêu:** chuẩn bị dữ liệu lớp đủ điều kiện để scheduling.

**Luồng chuẩn:**
1. Admin tạo course/program nếu chưa có.
2. Admin tạo teacher, room, shift.
3. Admin tạo class.
4. Admin gán teacher cho class.
5. Admin ghi danh học viên.
6. Admin cấu hình class schedule theo `day_of_week + shift_id + room_id(optional)`.
7. Class sẵn sàng cho scheduling.

**Nhánh thay thế:**
- tạo class trước, gán course/teacher sau;
- class có nhiều slot lịch tuần.

**Nhánh lỗi:**
- teacher/course thiếu;
- room capacity không hợp lý;
- roster vượt sĩ số.

## WF-03. Tạo preview scheduling và commit lesson

**Mục tiêu:** sinh lesson thực tế từ dữ liệu lớp.

**Luồng chuẩn:**
1. Admin chọn khoảng ngày và tập class.
2. Hệ thống tải class `OPEN`, shifts active, room, teacher, class_schedule.
3. Solver tạo assignments và conflicts.
4. Hệ thống trả preview.
5. Admin xem summary và conflict.
6. Nếu preview `COMPLETED`, admin commit.
7. Hệ thống kiểm tra overlap với lessons hiện có.
8. Sinh `Lesson`.

**Nhánh thay thế:**
- preview `PARTIAL`: admin sửa dữ liệu rồi chạy lại;
- benchmark solver để so sánh kết quả.

**Nhánh lỗi:**
- class thiếu teacher hoặc course;
- class không có class_schedule;
- phòng không đủ sức chứa;
- preview không tồn tại hoặc không đạt điều kiện commit.

## WF-04. Upload và duyệt tài liệu giảng dạy

**Mục tiêu:** kiểm soát tài liệu trước khi dùng trong giảng dạy.

**Luồng chuẩn:**
1. Giáo viên upload material.
2. Hệ thống validate file.
3. Lưu file local và metadata.
4. OCR stub trích xuất văn bản.
5. AI stub gắn nhãn.
6. Tạo audit log.
7. Reviewer xem material và audit reasoning.
8. Reviewer approve hoặc reject.
9. Material cập nhật trạng thái cuối.

**Nhánh thay thế:**
- reviewer tải file xuống để kiểm tra thủ công.

**Nhánh lỗi:**
- file không hợp lệ;
- lỗi storage;
- lỗi DB;
- reviewer không có quyền nhưng backend chưa chặn đúng.

## WF-05. Chuỗi học vụ sau lesson

**Mục tiêu:** theo dõi kết quả vận hành và chất lượng học tập sau khi lesson tồn tại.

**Luồng mục tiêu nghiệp vụ:**
1. Lesson đã được sinh từ scheduling.
2. Giáo viên điểm danh cho học viên.
3. Giáo viên tạo lesson summary.
4. Giáo viên ghi academic record cho từng học viên.
5. Học viên có thể tạo leave request liên quan.
6. Admin/giáo vụ xử lý leave request.

**Trạng thái hiện tại:** phần dữ liệu đã có nhưng API/UI chưa hoàn thiện.

---

# 9. Phân tích dữ liệu nghiệp vụ

## 9.1 Danh mục thực thể

Nhóm thực thể cốt lõi:

- `User`, `UserOTP`, `PasswordReset`
- `Student`, `Teacher`
- `Course`, `Program`, `ProgramCourse`, `Objective`, `Outcome`
- `Room`, `Shift`
- `Class`, `ClassSchedule`, `Enrollment`
- `Lesson`
- `Attendance`, `LessonSummary`, `AcademicRecord`, `LeaveRequest`
- `Material`, `AuditLog`, `ApprovalDecision`, `Label`
- `Consultation`

## 9.2 Mô tả dữ liệu lõi

| Thực thể | Vai trò nghiệp vụ |
|---|---|
| `Class` | trung tâm của vận hành học vụ |
| `ClassSchedule` | ràng buộc lịch tuần cố định cho lớp |
| `Shift` | chuẩn hóa khung giờ xếp lịch |
| `Enrollment` | roster giữa student và class |
| `Lesson` | đầu ra thực thi sau scheduling commit |
| `Attendance` | dữ liệu chuyên cần |
| `LessonSummary` | dữ liệu nội dung dạy |
| `AcademicRecord` | dữ liệu đánh giá học tập |
| `Material` | đầu mối content moderation |

## 9.3 Quan hệ dữ liệu chính

| Quan hệ | Ý nghĩa |
|---|---|
| Program - Course | N-N qua ProgramCourse |
| Class - Course | lớp học chính thuộc một khóa học |
| Class - Teacher | lớp có giáo viên phụ trách |
| Class - ClassSchedule | lớp có nhiều slot lịch tuần |
| ClassSchedule - Shift | mỗi slot tuần gắn một ca học |
| Enrollment - Student/Class | ghi danh học viên vào lớp |
| Lesson - Class/Teacher/Room | lesson là snapshot thực thi |
| Attendance - Lesson/Student | điểm danh theo buổi |
| LessonSummary - Lesson | 1 lesson có tối đa 1 summary |
| AcademicRecord - LessonSummary/Student | đánh giá từng học viên |
| LeaveRequest - Student/Class/Lesson | ngoại lệ học tập |
| Material - Teacher | tài liệu thuộc giáo viên |
| Material - AuditLog/ApprovalDecision/Label | chuỗi moderation |

## 9.4 Vòng đời dữ liệu

| Thực thể | Lifecycle chính |
|---|---|
| `User` | inactive -> active |
| `Class` | OPEN -> CLOSED/CANCELLED |
| `Enrollment` | APPLIED -> ENROLLED/REJECTED (dự kiến), nhưng hiện code tạo ENROLLED trực tiếp |
| `Material` | UPLOADED -> SCANNING -> AI_REVIEWED -> APPROVED/REJECTED |
| `LeaveRequest` | PENDING -> APPROVED/REJECTED (dự kiến) |
| `Program` | Draft -> Published -> Archived (dự kiến, chưa chốt) |

## 9.5 ERD scope recommendation

- **ERD lõi:** User, Student, Teacher, Course, Program, ProgramCourse, Room, Shift, Class, ClassSchedule, Enrollment, Lesson
- **ERD vận hành học vụ:** Attendance, LessonSummary, AcademicRecord, LeaveRequest
- **ERD hỗ trợ/kỹ thuật:** UserOTP, PasswordReset, Material, Label, AuditLog, ApprovalDecision, Consultation

---

# 10. Quy tắc nghiệp vụ và ràng buộc

## 10.1 Quy tắc xác thực

| Mã | Quy tắc |
|---|---|
| BR-AUTH-01 | email phải unique khi đăng ký |
| BR-AUTH-02 | user mới tạo phải xác minh OTP mới được login |
| BR-AUTH-03 | OTP phải còn hạn và chưa dùng |
| BR-AUTH-04 | token reset password phải còn hạn và chưa dùng |
| BR-AUTH-05 | chỉ user active mới đăng nhập được |

## 10.2 Quy tắc dữ liệu

| Mã | Quy tắc |
|---|---|
| BR-DATA-01 | `Teacher.code`, `Course.code`, `Program.code`, `Class.code`, `Room.code`, `Shift.code` phải unique |
| BR-DATA-02 | `Class.max_students` phải dương |
| BR-DATA-03 | `Room.capacity` phải dương |
| BR-DATA-04 | `Shift.duration_minutes` phải dương |
| BR-DATA-05 | `Material.file_type` chỉ thuộc tập định dạng được phép |
| BR-DATA-06 | `Material.file_size` không vượt 10MB |

## 10.3 Quy tắc vòng đời

| Mã | Quy tắc |
|---|---|
| BR-LFC-01 | class chỉ được scheduling khi `status = OPEN` |
| BR-LFC-02 | preview chỉ commit khi `status = COMPLETED` |
| BR-LFC-03 | material sau AI review mới được approve/reject |
| BR-LFC-04 | enrollment có dấu vết approval flow nhưng code hiện tại bỏ qua bước `APPLIED` |

## 10.4 Quy tắc phân quyền

| Mã | Quy tắc mong muốn / hiện trạng |
|---|---|
| BR-SEC-01 | teacher/course/program/class/room/shift/scheduling chủ yếu là admin functions |
| BR-SEC-02 | material upload là teacher-only ở route |
| BR-SEC-03 | material review theo nghiệp vụ phải là reviewer/admin, nhưng backend hiện chỉ yêu cầu auth |
| BR-SEC-04 | student CRUD theo nghiệp vụ nên là admin/vận hành, nhưng backend hiện chỉ yêu cầu auth |

## 10.5 Quy tắc xếp lịch

| Mã | Quy tắc |
|---|---|
| BR-SCH-01 | scheduling dùng classes `OPEN` |
| BR-SCH-02 | scheduling dùng shift `is_active = true` |
| BR-SCH-03 | scheduling dùng class_schedule theo `day_of_week + shift_id` |
| BR-SCH-04 | scheduling phải tránh trùng lớp, giáo viên, phòng |
| BR-SCH-05 | phải kiểm tra phòng đủ sức chứa |
| BR-SCH-06 | preview là kết quả tạm, commit mới sinh lesson thật |

## 10.6 Quy tắc còn mơ hồ / mâu thuẫn

| Mã | Mô tả |
|---|---|
| BR-GAP-01 | `Program.status` mâu thuẫn giữa migration và entity |
| BR-GAP-02 | `Room` có contract field không thống nhất giữa DTO và entity |
| BR-GAP-03 | `Course.price` và `Class.price` cùng tồn tại, chưa rõ giá nghiệp vụ chuẩn |
| BR-GAP-04 | `Attendance.status` là int nhưng chưa có bảng mã |
| BR-GAP-05 | `Enrollment` vừa là mapping vừa có dấu hiệu approval flow |
| BR-GAP-06 | `ClassSchedule` là dữ liệu lõi của scheduling nhưng thiếu API vận hành |

---

# 11. Ma trận phân quyền

## 11.1 Role vs function

| Function | Guest | Admin | Teacher | Student | Reviewer/Compliance | Hiện trạng |
|---|---|---|---|---|---|---|
| Register/login/reset | X | X | X | X | X | đúng |
| Student CRUD |  | X |  |  |  | backend hiện chỉ yêu cầu auth |
| Teacher CRUD |  | X |  |  |  | đúng |
| Course CRUD |  | X |  |  |  | đúng |
| Program CRUD |  | X |  |  |  | đúng |
| Room CRUD |  | X |  |  |  | đúng |
| Shift CRUD |  | X |  |  |  | đúng |
| Class CRUD |  | X |  |  |  | đúng |
| Enrollment ops |  | X |  |  |  | đúng về route |
| Scheduling preview/commit/benchmark |  | X |  |  |  | đúng |
| Upload material |  |  | X |  |  | đúng |
| Review material |  | X? | ? | ? | X | backend chưa chặn bằng role cụ thể |

## 11.2 Risk points

| Mã | Rủi ro | Tác động |
|---|---|---|
| SEC-R01 | student CRUD chỉ cần auth | user không phù hợp có thể tạo/sửa/xóa học viên |
| SEC-R02 | material review chỉ cần auth | teacher/student có thể review nếu biết API |
| SEC-R03 | `teacher_id` upload material lấy từ request body | giả mạo ownership |
| SEC-R04 | `compliance_officer_id` review material lấy từ request body | giả mạo người duyệt |
| SEC-R05 | `/auth/me` triển khai lỗi | thông tin profile không ổn định và có thể phát sinh bug logic |

## 11.3 Endpoints thiếu role check

| Endpoint | Vấn đề |
|---|---|
| `/api/v1/students/*` | chỉ dùng AuthMiddleware, không có RoleMiddleware("ADMIN") |
| `/api/v1/materials/review` | chỉ dùng AuthMiddleware, thiếu role reviewer/admin |
| `/api/v1/materials/*` download/detail/list | cần rà ownership nếu mở cho nhiều vai trò |

## 11.4 Nhận xét về kiểm soát truy cập hiện tại

Hệ thống đã có nền middleware auth và role check, nhưng phân quyền hiện tại vẫn mang tính **module-based coarse control**. Một số module quan trọng đã khóa admin đúng hướng, nhưng:

- chưa có ownership enforcement đầy đủ;
- chưa có role reviewer riêng ở backend;
- chưa có policy truy cập chéo giữa admin, teacher, student, compliance.

---

# 12. Danh mục API / màn hình / báo cáo

## 12.1 Các API nghiệp vụ chính

| Module | API tiêu biểu | Mục đích |
|---|---|---|
| Auth | `/api/v1/auth/register`, `/verify-otp`, `/login`, `/refresh`, `/forgot-password`, `/reset-password`, `/change-password` | xác thực và vòng đời tài khoản |
| Student | `/api/v1/students` | CRUD học viên |
| Teacher | `/api/v1/teachers`, `/:id/timetable`, `/:id/stats/teaching-hours` | CRUD + lịch dạy + thống kê |
| Course | `/api/v1/courses` | CRUD khóa học |
| Program | `/api/v1/programs`, `/:id/courses` | CRUD chương trình + gán khóa học |
| Room | `/api/v1/rooms` | CRUD phòng |
| Shift | `/api/v1/shifts` | CRUD ca học |
| Class | `/api/v1/classes`, `/:id/students`, `/:id/assign-teacher` | CRUD lớp, roster, assign teacher |
| Scheduling | `/api/v1/scheduling/preview`, `/preview/latest`, `/preview/:id`, `/commit`, `/benchmark` | preview, commit, benchmark |
| Material | `/api/v1/materials/upload`, `/materials`, `/materials/:id`, `/materials/:id/download`, `/materials/review` | upload, list/detail/download, review |

## 12.2 Màn hình nghiệp vụ suy luận được

| Mã màn hình | Tên màn hình | Actor | Hiện trạng |
|---|---|---|---|
| SCR-01 | Đăng ký / đăng nhập / quên mật khẩu / reset | Guest | có |
| SCR-02 | Hồ sơ cá nhân | User | có nhưng `/me` backend lỗi |
| SCR-03 | Quản lý học viên | Admin | có |
| SCR-04 | Quản lý giáo viên | Admin | có |
| SCR-05 | Quản lý khóa học | Admin | có |
| SCR-06 | Quản lý chương trình đào tạo | Admin | có |
| SCR-07 | Quản lý phòng học | Admin | có |
| SCR-08 | Quản lý ca học | Admin | có |
| SCR-09 | Quản lý lớp học | Admin | có |
| SCR-10 | Chi tiết lớp / roster | Admin | có |
| SCR-11 | Xếp lịch | Admin | có |
| SCR-12 | Tài liệu giảng dạy của giáo viên | Teacher | có |
| SCR-13 | Hàng chờ duyệt tài liệu | Compliance/Admin | có |
| SCR-14 | Tổng quan admin | Admin | có nhưng nhiều mock data |
| SCR-15 | Tổng quan học viên | Student | placeholder |

## 12.3 Báo cáo / dashboard / thống kê

| Mã | Tên | Hiện trạng | Giá trị nghiệp vụ |
|---|---|---|---|
| REP-01 | Teaching hours statistics | có API thật | theo dõi tải giảng dạy giáo viên |
| REP-02 | Scheduling preview summary | có trong preview response | đánh giá khả thi của lịch học |
| REP-03 | Admin overview dashboard | chủ yếu mock | định hướng dashboard tổng quan |

---

# 13. Phân tích khoảng trống và rủi ro

## 13.1 Gap giữa code và vận hành nghiệp vụ kỳ vọng

| Mã gap | Miền | Mô tả | Mức độ |
|---|---|---|---|
| GAP-01 | Student | CRUD học viên chưa khóa role admin | Cao |
| GAP-02 | Material | reviewer chưa có role check thực | Cao |
| GAP-03 | Material | ownership/reviewer identity lấy từ payload thay vì JWT | Cao |
| GAP-04 | Program | lifecycle publish/archive chưa đầy đủ dù đã có field | Trung bình |
| GAP-05 | ClassSchedule | thực thể lõi của scheduling nhưng thiếu API quản trị | Cao |
| GAP-06 | Enrollment | lifecycle APPLIED/ENROLLED/REJECTED không thống nhất | Cao |
| GAP-07 | Lesson | lesson đã sinh thật nhưng chưa có module quản trị riêng | Trung bình |
| GAP-08 | Attendance/Summary/Academic/Leave | có schema nhưng thiếu API/UI | Cao |
| GAP-09 | Consultation | entity có nhưng chưa thành module | Thấp |
| GAP-10 | Benchmark solver | có contract nhưng chưa là benchmark vận hành hoàn chỉnh | Trung bình |

## 13.2 Rủi ro nghiệp vụ

| Mã | Rủi ro |
|---|---|
| RSK-BIZ-01 | dữ liệu lớp có thể được tạo nhưng không đủ điều kiện scheduling do thiếu teacher/course/class_schedule |
| RSK-BIZ-02 | roster có thể không phản ánh đúng approval lifecycle của enrollment |
| RSK-BIZ-03 | thiếu module attendance/summary/academic khiến lesson kết thúc nhưng không có hậu kiểm học vụ |
| RSK-BIZ-04 | material moderation có thể bị thao tác bởi actor không phù hợp |

## 13.3 Rủi ro dữ liệu

| Mã | Rủi ro |
|---|---|
| RSK-DATA-01 | `Course.price` và `Class.price` gây hiểu sai giá chuẩn |
| RSK-DATA-02 | `Attendance.status` không có domain chuẩn |
| RSK-DATA-03 | `Room` contract lệch giữa DTO/entity/UI |
| RSK-DATA-04 | `Program.status` không nhất quán giữa migration và entity |

## 13.4 Rủi ro bảo mật

| Mã | Rủi ro |
|---|---|
| RSK-SEC-01 | authorization gap ở student CRUD |
| RSK-SEC-02 | authorization gap ở material review |
| RSK-SEC-03 | identity spoofing qua `teacher_id` và `compliance_officer_id` trong request body |
| RSK-SEC-04 | `/auth/me` lỗi logic ID có thể gây sai profile hoặc phát sinh lỗi runtime |

## 13.5 Rủi ro do module chưa hoàn thiện

| Mã | Rủi ro |
|---|---|
| RSK-PART-01 | BA có thể vẽ thừa quy trình attendance/academic/leave nếu không tách rõ current-state |
| RSK-PART-02 | benchmark scheduling có thể bị hiểu nhầm là hoàn tất dù mới là contract |
| RSK-PART-03 | dashboard overview có thể bị hiểu là số liệu thật dù đang dùng mock |

---

# 14. Câu hỏi cần xác nhận với BA/PO/stakeholder

## 14.1 Product Owner / Quản lý đào tạo

| Ưu tiên | Câu hỏi |
|---|---|
| Cao | `Program` có lifecycle chính thức là Draft/Published/Archived hay không? |
| Cao | Giá bán nghiệp vụ chuẩn nằm ở `Course.price` hay `Class.price`? |
| Cao | `Enrollment` có cần approval flow thật hay chỉ là mapping trực tiếp vào lớp? |
| Cao | `ClassSchedule` được quản trị bởi admin hay giáo vụ? |
| Trung bình | `Objective` và `Outcome` có nằm trong scope hiện tại không? |

## 14.2 Vận hành / giáo vụ

| Ưu tiên | Câu hỏi |
|---|---|
| Cao | Khi ghi danh học viên, có cần bước duyệt hay chỉ cần thêm thẳng vào roster? |
| Cao | Sĩ số lớp tính theo `max_students`, theo `room.capacity`, hay theo min của cả hai? |
| Cao | Khi class đã có lesson, có được đổi teacher/course/room hay hủy lớp không? |
| Trung bình | `day_of_week` trong class_schedule nên chuẩn hóa theo chữ hay số? |

## 14.3 Giáo viên

| Ưu tiên | Câu hỏi |
|---|---|
| Cao | Giáo viên có cần tự chấm attendance và tạo lesson summary hay không? |
| Cao | Có cần teacher portal xem học viên lớp mình, leave request, academic record không? |
| Trung bình | Tài liệu SAFE có cần duyệt thủ công hay được auto-approve? |

## 14.4 Học viên / phụ huynh

| Ưu tiên | Câu hỏi |
|---|---|
| Cao | Học viên hay phụ huynh là người tạo leave request? |
| Trung bình | Học viên có cần xem timetable, attendance, academic record trên portal hay không? |
| Trung bình | Có cần actor phụ huynh độc lập hay chỉ lưu thông tin liên hệ? |

## 14.5 Tech Lead / Security

| Ưu tiên | Câu hỏi |
|---|---|
| Cao | Có chốt role reviewer/compliance riêng ở backend không? |
| Cao | `teacher_id` và `compliance_officer_id` có được chuyển sang derive từ JWT không? |
| Cao | `Attendance.status` chuẩn enum là gì? |
| Trung bình | Có persist preview scheduling xuống DB hay vẫn giữ in-memory? |

---

# 15. Đề xuất artifact BA cần vẽ tiếp theo

## 15.1 Nhóm sơ đồ use case cần vẽ

1. Sơ đồ xác thực tài khoản.
2. Sơ đồ quản lý danh mục đào tạo và nguồn lực.
3. Sơ đồ lớp học và ghi danh.
4. Sơ đồ xếp lịch và commit lesson.
5. Sơ đồ tài liệu giảng dạy và kiểm duyệt.
6. Sơ đồ học vụ hậu lesson.

## 15.2 Cây phân rã chức năng chính thức

Nên chuẩn hóa lại cây gồm 7 phân hệ:

1. Auth
2. Master data
3. Class & enrollment
4. Scheduling & lesson
5. Academic operations
6. Material moderation
7. Consultation & analytics

## 15.3 ERD lõi và ERD mở rộng

- **ERD lõi:** User, Student, Teacher, Course, Program, ProgramCourse, Room, Shift, Class, ClassSchedule, Enrollment, Lesson
- **ERD mở rộng:** Attendance, LessonSummary, AcademicRecord, LeaveRequest, Material, AuditLog, ApprovalDecision, Label, Consultation

## 15.4 BPMN ưu tiên

1. Đăng ký -> OTP -> login
2. Mở lớp -> ghi danh -> phân công giáo viên -> cấu hình lịch tuần
3. Tạo preview scheduling -> xử lý conflict -> commit lesson
4. Upload material -> AI audit -> review decision
5. Chuỗi lesson -> attendance -> summary -> academic record

## 15.5 CRUD matrix

BA nên dựng CRUD matrix ít nhất cho:

- User
- Student
- Teacher
- Course
- Program
- Room
- Shift
- Class
- ClassSchedule
- Enrollment
- Lesson
- Material

## 15.6 Actor-use case matrix

Nên chốt lại ma trận actor-use case với 5 actor chính:

- Guest
- Admin
- Teacher
- Student
- Reviewer/Compliance

## 15.7 State transition diagram cần thiết

1. User
2. Program
3. Class
4. Enrollment
5. Material
6. LeaveRequest
7. Preview Scheduling

---

# 16. Kết luận

## 16.1 Tóm tắt hiện trạng hệ thống

EduCenter hiện đã có một lõi vận hành đủ rõ cho trung tâm dạy thêm:

- quản lý master data;
- tạo lớp và roster;
- quản lý room/shift;
- xếp lịch preview và commit thành lesson;
- upload và duyệt tài liệu ở mức demo workflow.

## 16.2 Điểm mạnh

- Mô hình lớp học, shift, room, class_schedule khá phù hợp với bài toán trung tâm dạy thêm.
- Scheduling đã được tách thành kiến trúc solver service, tạo nền tốt cho benchmark và tối ưu hóa.
- Material moderation đã có đủ chuỗi metadata, audit log, decision history.
- Data model cho học vụ hậu lesson và predictive analytics đã được đặt nền.

## 16.3 Điểm thiếu

- Nhiều module học vụ sau lesson mới dừng ở entity.
- Có authorization gap ở student và material review.
- Có mâu thuẫn dữ liệu/lifecycle ở Program, Room, Enrollment, Attendance.
- Một số API quan trọng như class schedule management, attendance, lesson summary, academic record, leave request chưa có.

## 16.4 Các bước BA nên thực hiện tiếp

1. Chốt current-state và target-state theo từng module.
2. Xác nhận các điểm mâu thuẫn dữ liệu và lifecycle.
3. Vẽ ERD lõi trước, rồi ERD mở rộng.
4. Vẽ BPMN cho 4 workflow trọng yếu.
5. Chốt permission matrix chính thức.
6. Chuyển báo cáo này thành BRD/SRS sơ bộ theo module.

---

# 17. Phụ lục

## 17.1 Glossary

| Thuật ngữ | Nghĩa |
|---|---|
| Program | khung chương trình đào tạo gồm nhiều khóa học |
| Course | khóa học cơ sở dùng để mở lớp |
| Class | lớp học thực tế |
| Shift | ca học chuẩn |
| ClassSchedule | lịch tuần chuẩn của lớp |
| Enrollment | ghi danh học viên vào lớp |
| Lesson | buổi học thực tế sau scheduling commit |
| LessonSummary | tổng kết buổi học |
| AcademicRecord | đánh giá học tập theo học viên |
| Material | tài liệu giảng dạy |
| AuditLog | lịch sử OCR/AI audit |
| ApprovalDecision | quyết định duyệt tài liệu |

## 17.2 Rule catalog summary

| Nhóm rule | Nội dung chính |
|---|---|
| Auth | email unique, OTP hết hạn, user inactive không login được |
| Data | unique code cho nhiều master entity, validate room/shift/material |
| Lifecycle | class OPEN mới scheduling, material phải qua AI_REVIEWED rồi mới review |
| Scheduling | tránh trùng lớp/giáo viên/phòng, dùng shift active và class_schedule |
| Security | nhiều module cần siết role và ownership |

## 17.3 Entity summary

| Entity | Vai trò |
|---|---|
| User | tài khoản hệ thống |
| Student | master học viên |
| Teacher | master giáo viên |
| Course | master khóa học |
| Program | master chương trình |
| Room | master phòng |
| Shift | master ca học |
| Class | trung tâm vận hành học vụ |
| ClassSchedule | ràng buộc lịch tuần |
| Enrollment | roster lớp |
| Lesson | đầu ra của scheduling |
| Attendance | chuyên cần |
| LessonSummary | nội dung dạy |
| AcademicRecord | kết quả học tập |
| LeaveRequest | ngoại lệ học tập |
| Material | nội dung giảng dạy cần kiểm soát |

## 17.4 Open assumptions register

| Mã | Giả định mở |
|---|---|
| ASM-01 | trung tâm đang là single-tenant |
| ASM-02 | admin hiện kiêm vai trò giáo vụ và quản lý đào tạo |
| ASM-03 | student portal và leave/academic modules là target-state gần |
| ASM-04 | reviewer/compliance là actor riêng, nhưng chưa có backend role riêng |
| ASM-05 | predictive analytics sẽ dùng dữ liệu attendance + academic record + lesson summary |

## 17.5 Evidence map

| Nhóm bằng chứng | Vị trí tiêu biểu |
|---|---|
| Route/controller | `cmd/http/controllers/*`, `cmd/http/main.go` |
| Middleware auth/role | `cmd/http/middleware/auth.go` |
| Entity/model | `internal/entities/*` |
| Use case | `internal/usecases/*` |
| Scheduling service | `internal/services/scheduling/*` |
| Material audit service | `internal/services/audit/services.go` |
| Frontend route | `frontend/src/App.tsx` |
| Frontend admin pages | `frontend/src/pages/admin/*` |
| Frontend material pages | `frontend/src/pages/teacher/TeacherDocumentsPage.tsx`, `frontend/src/pages/compliance/ComplianceQueuePage.tsx` |

