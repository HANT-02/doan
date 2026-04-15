# Quản lý Task - Hệ thống quản lý trung tâm dạy thêm

**Đề tài hiện tại:** Xây dựng hệ thống quản lý trung tâm dạy thêm, trọng tâm là xếp lịch thông minh và dự báo sinh viên có nguy cơ học kém

**Mô hình:** Single-tenant

**Stack chính:** Golang (Backend), React + TypeScript (Frontend), PostgreSQL

**Kiến trúc:** Clean Architecture + Wire DI + Gin + GORM

**Tuân thủ:** Thông tư 29/2024/TT-BGDĐT

---

## 1. Mục tiêu dự án đã chốt

Mục tiêu của đồ án ở giai đoạn hiện tại được chốt lại như sau:

1. Hoàn thiện hệ thống quản lý đào tạo cốt lõi:
   - giáo viên,
   - học sinh,
   - lớp học,
   - phòng học,
   - chương trình/khóa học,
   - tài khoản.

2. Nâng scheduling thành điểm nhấn kỹ thuật chính của đồ án:
   - thêm module `Shift`,
   - refactor `class_schedule` sang `shift_id`,
   - benchmark 3 thuật toán,
   - chọn 1 solver tốt nhất để dùng cho API chính.

3. Bổ sung một bài toán AI đúng nghĩa machine learning:
   - `AT_RISK classification`,
   - huấn luyện trong backend hiện tại,
   - có metric đánh giá,
   - có endpoint dự báo và màn hình cảnh báo.

4. Giữ phạm vi đủ sâu nhưng không dàn trải:
   - trong nhóm module phụ phía admin chỉ giữ **Quản lý tài khoản**,
   - **không làm module AI Audit**,
   - **không làm nhánh Compliance/Audit/Chatbot/DevTools** trong backlog chính,
   - ưu tiên scheduling trước, predictive analytics sau.

---

## 2. Rule triển khai đã chốt

Các task bên dưới phải tuân thủ đúng các rule sau:

1. Bám Clean Architecture hiện có của project.
2. Không thêm flow hoặc module đi ngược cấu trúc hiện tại nếu chưa có lý do rõ ràng.
3. Với scheduling:
   - thuật toán phải đi qua **interface ở tầng service**,
   - mỗi thuật toán là một implementation riêng,
   - use case scheduling chỉ dùng abstraction của service solver,
   - benchmark sẽ chạy qua admin API nội bộ.
4. `Shift` là thực thể chuẩn cho ca học.
5. `class_schedule` sẽ dùng **hoàn toàn `shift_id`**.
6. Frontend production-like không cho người dùng cuối chọn thuật toán.
7. Benchmark chỉ dành cho admin/nghiên cứu nội bộ.
8. Predictive analytics không dùng prompt để thay thế mô hình.
9. Chỉ giữ những task còn nằm trong scope thực tế của đồ án.
10. Ở nhóm module phụ trong sidebar admin:
   - chỉ giữ `Quản lý tài khoản`,
   - không tiếp tục đầu tư `Kiểm duyệt tài liệu (AI)`,
   - không tiếp tục đầu tư `DevTools`,
   - các placeholder còn sót phải được dọn để UI khớp với backlog chính.

---

## 3. Ngoài phạm vi

Các nhánh dưới đây được bỏ khỏi file task để giữ backlog rõ ràng:

- Module AI Audit / kiểm duyệt tài liệu
- Compliance dashboard cho AI Audit
- DevTools page / công cụ kiểm thử nội bộ không phục vụ trực tiếp cho scope đồ án
- Legal & Profiles placeholder không nằm trong luồng nghiệp vụ chính
- Chatbot trợ lý ảo
- Các nhánh prompt/OCR/Gemini phục vụ audit tài liệu
- Các task mở rộng không phục vụ trực tiếp cho mục tiêu scheduling + predictive analytics

Ghi chú:

- Code cũ nếu đã tồn tại trong repo thì không bị xóa tự động chỉ vì ra khỏi scope.
- Từ thời điểm này trở đi, backlog chính **không tiếp tục đầu tư** vào nhánh AI Audit.

---

## 4. Nền tảng hiện có

Những phần dưới đây đã có nền tảng triển khai và được xem là baseline hiện tại:

### 4.1. Core admin modules

- [x] Auth + role gating admin cơ bản
- [x] Quản lý giáo viên
- [x] Quản lý học sinh
- [x] Quản lý lớp học
- [x] Quản lý phòng học
- [x] Quản lý chương trình/khóa học
- [x] Quản lý tài khoản mức cơ bản

Ghi chú:

- `Quản lý tài khoản` hiện mới dừng ở auth flow và profile cơ bản.
- Màn hình admin quản lý user / role / trạng thái tài khoản vẫn chưa được triển khai đầy đủ.

### 4.2. Scheduling baseline hiện có

- [x] Có màn hình scheduling preview trên frontend
- [x] Có preview backend và commit xuống `lesson`
- [x] Có conflict messaging cơ bản
- [x] Có tài liệu mô tả logic hiện tại tại:
  - [x] `docs/modeling/scheduling/SCHEDULING_CORE_LOGIC_CURRENT.md`

### 4.3. Tài liệu định hướng mới

- [x] Đã chốt ý tưởng hướng mới tại:
  - [x] `docs/ke_hoach_phan_hoi_gvhd_2026-04-09.md`

### 4.4. Khoảng trống so với sơ đồ phân rã chức năng

Các điểm dưới đây hiện vẫn chưa khớp hoàn toàn với sơ đồ phân rã chức năng đã chốt:

- `Quản lý tài khoản` trong admin sidebar vẫn là placeholder, chưa có CRUD user và phân quyền thật.
- `Kiểm duyệt tài liệu (AI)` và `DevTools` vẫn còn xuất hiện ở UI nhưng không còn nằm trong scope backlog chính.
- `Quản lý lịch tuần lớp` chưa thành một module hoàn chỉnh theo sơ đồ nghiệp vụ.
- `Lesson / Điểm danh / Tổng kết buổi học / Kết quả học tập / Đơn xin phép` mới có entity hoặc nền backend rời rạc, chưa thành luồng đầy đủ → **được tách ra thành Giai đoạn I (Teacher Portal) và Giai đoạn J (Student Portal)** với task chi tiết.
- Portal Giáo viên (điểm danh, sổ đầu bài, vào điểm, duyệt đơn) chưa có usecase/controller/frontend page → **xem Task I1–I10**.
- Portal Sinh viên (thời khóa biểu, theo dõi điểm danh, điểm số, xin nghỉ, cảnh báo AT_RISK) chưa có usecase/controller/frontend page → **xem Task J1–J10**.
- Cần một giai đoạn riêng để đối chiếu lại backlog với sơ đồ phân rã chức năng trước khi chuyển sang kiểm thử và báo cáo.

---

## 5. Backlog triển khai chính

### Giai đoạn A - Khóa lại kiến trúc scheduling

#### Task A1. Chuẩn hóa kiến trúc solver scheduling
- [x] Định nghĩa interface `SchedulingSolver` ở tầng `service`
- [x] Chuẩn hóa `SchedulingInput`
- [x] Chuẩn hóa `SchedulingOutput`
- [x] Chuẩn hóa scorer cho hard constraints / soft constraints
- [x] Chuẩn hóa benchmark result model dùng chung cho 3 solver

**Mục tiêu đóng task:**
- Có abstraction thống nhất để cắm 3 solver mà không lặp use case.

#### Task A2. Thiết kế contract admin benchmark API
- [x] Thiết kế request benchmark
- [x] Thiết kế response benchmark cho 3 solver
- [x] Thiết kế metric trả về:
  - [x] feasibility
  - [x] hard violations
  - [x] soft score
  - [x] runtime
  - [x] solver summary
- [x] Quy ước benchmark API chỉ dành cho admin nội bộ

**Mục tiêu đóng task:**
- Có contract rõ trước khi code benchmark API và solver.

---

### Giai đoạn B - Module `Shift`

#### Task B1. Thiết kế dữ liệu `Shift`
- [x] Chốt schema `shifts`
- [x] Chốt field:
  - [x] `id`
  - [x] `code`
  - [x] `name`
  - [x] `start_time`
  - [x] `end_time`
  - [x] `duration_minutes`
  - [x] `session_type`
  - [x] `is_active`
  - [x] `notes`
- [x] Tài liệu hóa ràng buộc dữ liệu của `Shift`

**Mục tiêu đóng task:**
- Có schema rõ để viết migration và CRUD.

#### Task B2. Backend CRUD `Shift`
- [x] Tạo migration bảng `shifts`
- [x] Tạo entity `Shift`
- [x] Tạo repository cho `Shift`
- [x] Tạo service/use case CRUD `Shift`
- [x] Tạo controller / route CRUD `Shift`
- [x] Wire dependency cho module `Shift`

**Mục tiêu đóng task:**
- Admin API cho `Shift` dùng được end-to-end.

#### Task B3. Frontend quản lý `Shift`
- [x] Tạo page danh sách `Shift`
- [x] Tạo form thêm/sửa `Shift`
- [x] Tạo thao tác bật/tắt `is_active`
- [x] Tích hợp API `Shift`
- [x] Kiểm tra UX list/form/detail theo pattern project

**Mục tiêu đóng task:**
- Admin quản lý được ca học trên UI.

---

### Giai đoạn C - Refactor scheduling sang `Shift`

#### Task C1. Refactor `class_schedule` sang `shift_id`
- [x] Thiết kế migration chuyển dữ liệu `class_schedule`
- [x] Bỏ phụ thuộc chính vào `start_time/end_time`
- [x] Cập nhật entity/repository/use case liên quan
- [x] Cập nhật controller/API liên quan
- [x] Kiểm tra backward impact với dữ liệu hiện có

**Mục tiêu đóng task:**
- `class_schedule` dùng chuẩn `shift_id` trong luồng chính.

#### Task C2. Cập nhật scheduling domain theo `Shift`
- [x] Sửa scheduling input để lấy `Shift`
- [x] Sửa domain generation để tạo assignment theo `Shift`
- [x] Sửa hard constraint check theo `Shift`
- [x] Sửa soft score theo `Shift`
- [x] Sửa preview output để hiển thị theo `Shift`

**Mục tiêu đóng task:**
- Scheduling domain không còn phụ thuộc vào time slot hardcode kiểu cũ.

#### Task C3. Cập nhật UI scheduling theo `Shift`
- [x] Hiển thị thông tin `Shift` trên preview/result
- [x] Hiển thị lịch mẫu theo `Shift`
- [x] Kiểm tra flow preview/commit sau refactor
- [x] Giữ UI production-like gọn, không để lộ logic benchmark

**Mục tiêu đóng task:**
- UI scheduling chạy đúng trên dữ liệu `Shift`.

---

### Giai đoạn D - 3 solver scheduling

#### Task D1. Implement `GraphColoringSolver`
- [x] Tạo implementation `GraphColoringSolver`
- [x] Tạo logic heuristic baseline
- [x] Nối với scorer chung
- [x] Test được trên input scheduling chuẩn

**Mục tiêu đóng task:**
- Có solver baseline chạy được qua abstraction chung.

#### Task D2. Implement `CPSATSolver`
- [x] Tạo implementation `CPSATSolver`
- [x] Mô hình hóa biến/ràng buộc cho CP-SAT
- [x] Nối với scorer chung
- [x] Test được trên input scheduling chuẩn

**Mục tiêu đóng task:**
- Có solver CP-SAT chạy được qua abstraction chung.

#### Task D3. Implement `TabuSearchSolver`
- [x] Tạo implementation `TabuSearchSolver`
- [x] Thiết kế neighborhood move
- [x] Thiết kế tabu memory / stopping criteria
- [x] Nối với scorer chung
- [x] Test được trên input scheduling chuẩn

**Mục tiêu đóng task:**
- Có solver Tabu Search chạy được qua abstraction chung.

---

### Giai đoạn E - Benchmark scheduling

#### Task E1. Benchmark admin API
- [x] Tạo admin API benchmark
- [x] Cho phép chạy cùng input qua 3 solver
- [x] Trả bảng metric so sánh
- [x] Kiểm tra quyền admin nội bộ

**Mục tiêu đóng task:**
- Có admin API benchmark dùng được để lấy số liệu báo cáo.

#### Task E2. Chạy benchmark và chọn solver chính
- [x] Chuẩn bị bộ dữ liệu benchmark
- [x] Chạy benchmark nhiều lần
- [x] Ghi feasibility / hard violations / soft score / runtime
- [x] Chọn solver tốt nhất
- [x] Tài liệu hóa quyết định chọn solver

**Mục tiêu đóng task:**
- Có bằng chứng thực nghiệm để chọn solver chính.

#### Task E3. Gắn solver tốt nhất vào scheduling API chính
- [x] Inject solver chính vào use case scheduling
- [x] Giữ benchmark API tách riêng
- [x] Kiểm tra lại preview
- [x] Kiểm tra lại commit
- [ ] Kiểm tra luồng FE-BE sau khi đổi solver

**Mục tiêu đóng task:**
- Scheduling production-like dùng đúng solver đã chọn.

---

### Giai đoạn F - Predictive analytics

#### Task F1. Chốt dữ liệu đầu vào cho `AT_RISK classification`
- [x] Chốt nguồn dữ liệu:
  - [x] student
  - [x] attendance
  - [x] grade
  - [x] class/course enrollment
  - [x] tín hiệu vận hành cần thiết
- [x] Chốt label `AT_RISK`
- [x] Chốt feature set ban đầu

**Mục tiêu đóng task:**
- Có dataset definition rõ để bắt đầu pipeline ML.

#### Task F2. Pipeline ML trong backend hiện tại
- [x] Tạo pipeline feature engineering
- [x] Tạo pipeline train/test split
- [x] Tạo nạp dataset từ PostgreSQL hiện tại
- [x] Tạo baseline `Rule-based`
- [x] Tạo huấn luyện `Logistic Regression`
- [x] Tạo metric evaluation

**Mục tiêu đóng task:**
- Backend hiện tại train/evaluate được các mô hình classification cơ bản.

#### Task F3. Prediction API + model metadata
- [x] Thiết kế lưu model metadata
- [x] Tạo endpoint dự báo `AT_RISK`
- [x] Tạo output explanation cơ bản
- [x] Tạo cơ chế versioning tối thiểu cho model

**Mục tiêu đóng task:**
- Có API dự báo dùng được trong hệ thống.

#### Task F4. Frontend cảnh báo sinh viên nguy cơ học kém
- [x] Tạo màn hình danh sách sinh viên `AT_RISK`
- [x] Hiển thị score/label
- [x] Hiển thị insight hoặc lý do chính
- [x] Kiểm tra flow end-to-end

**Mục tiêu đóng task:**
- Có UI demo được giá trị của predictive analytics.

---

### Giai đoạn G - Khớp sơ đồ phân rã chức năng và thu gọn scope phụ

#### Task G1. Thu gọn module phụ trong admin sidebar
- [x] Bỏ menu/route/page `Kiểm duyệt tài liệu (AI)` khỏi luồng admin chính
- [x] Bỏ menu/route/page `DevTools` khỏi luồng admin chính
- [x] Bỏ placeholder `Legal & Profiles` khỏi điều hướng chính
- [x] Giữ lại và đổi tên rõ module `Quản lý tài khoản`
- [x] Rà lại tài liệu backlog để scope UI khớp với scope đồ án
- [x] Gộp các màn dữ liệu cơ bản thành cụm `Thông tin cơ bản` trong sidebar và khôi phục lối vào riêng cho `Khóa học`

**Mục tiêu đóng task:**
- Sidebar admin chỉ còn các module đúng scope đồ án; không còn nhánh AI Audit / DevTools gây lệch phạm vi.

#### Task G2. Hoàn thiện module `Quản lý tài khoản`
- [x] Tạo trang danh sách user cho `ADMIN` và `SUPER_ADMIN`
- [x] Hỗ trợ tìm kiếm / lọc theo role / status
- [x] Xem chi tiết user và hồ sơ cơ bản
- [x] Tạo tài khoản nội bộ mới
- [x] Cập nhật role của user
- [x] Kích hoạt / vô hiệu hóa tài khoản
- [x] Thêm thao tác reset mật khẩu ở phía admin
- [x] Kiểm tra role guard để chỉ `ADMIN` / `SUPER_ADMIN` được thao tác

**Mục tiêu đóng task:**
- `Quản lý tài khoản` trở thành module thật thay vì placeholder; admin và super admin quản lý được user + phân quyền.

#### Task G3. Hoàn thiện `Quản lý lịch tuần lớp`
- [x] Chốt lại scope nghiệp vụ của `class_schedule` theo sơ đồ phân rã
- [x] Tạo API CRUD lịch tuần lớp
- [x] Gắn `shift_id`, `day_of_week`, `room_id` đúng chuẩn scheduling mới
- [x] Tạo UI quản lý lịch tuần lớp từ màn hình lớp học
- [x] Kiểm tra dữ liệu lịch tuần lớp dùng được cho preview scheduling

**Mục tiêu đóng task:**
- Nhánh `Quản lý lịch tuần lớp` trong sơ đồ có implementation rõ ràng và dùng chung được với scheduling.

#### Task G4. Hoàn thiện `Lesson` sau khi commit scheduling
- [x] Tạo màn hình danh sách `lesson`
- [x] Tạo màn hình chi tiết `lesson`
- [x] Hỗ trợ lọc theo lớp / giáo viên / khoảng ngày / trạng thái
- [x] Nối điều hướng từ scheduling hoặc class sang lesson
- [x] Kiểm tra dữ liệu lesson sinh ra từ commit preview hiển thị đúng trên UI
- [x] Sửa lỗi list `lesson` sau review 2026-04-14: preload dữ liệu liên quan, sửa sort backend và làm lại UI filter/table

**Mục tiêu đóng task:**
- `Lesson` không chỉ là dữ liệu nền mà trở thành module vận hành có thể xem và theo dõi được.

#### Task G5. Hoàn thiện `Điểm danh` và `Tổng kết buổi học`
- [x] Tạo API chấm điểm danh theo lesson
- [x] Tạo API cập nhật trạng thái chuyên cần
- [x] Tạo API ghi chú buổi học
- [x] Tạo API ghi nội dung dạy
- [x] Tạo API giao bài tập
- [x] Tạo UI teacher/admin cho điểm danh và tổng kết buổi học
- [x] Bổ sung phân quyền lesson operations cho `ADMIN` / `SUPER_ADMIN` / giáo viên phụ trách lesson
- [x] Tái sử dụng cùng contract `lesson_id` cho admin detail và teacher portal để tránh lệch FE-BE

**Mục tiêu đóng task:**
- Sau mỗi lesson, hệ thống ghi nhận được chuyên cần và summary theo đúng nhánh nghiệp vụ trong sơ đồ.

#### Task G6. Hoàn thiện `Kết quả học tập` và `Đơn xin phép`
- [x] Tạo API đánh giá bài tập
- [x] Tạo API đánh giá thái độ / tham gia
- [x] Tạo API chốt `academic_record`
- [x] Tạo API tạo đơn xin phép
- [x] Tạo API duyệt / từ chối đơn xin phép
- [x] Tạo UI phù hợp cho teacher/admin/student theo từng tác vụ
- [x] Tái sử dụng `lesson_id` làm trục teacher/admin cho `academic_record`
- [x] Tạo actor-based flow cho `leave_request` để student tạo, teacher/admin duyệt, student hủy

**Mục tiêu đóng task:**
- Nhánh vận hành sau buổi học khép kín: từ lesson -> attendance -> summary -> academic record -> leave request.

#### Task G7. Đối chiếu backlog với sơ đồ phân rã chức năng
- [x] Lập bảng map giữa từng nhánh trong sơ đồ và API/UI hiện có
- [x] Đánh dấu rõ `Implemented / Partial / Missing`
- [x] Chốt phần nào thuộc scope làm tiếp và phần nào chỉ ghi nhận ngoài scope
- [x] Cập nhật lại file task sau khi review scope

**Kết quả chốt 2026-04-15:**
- [x] Tạo bảng đối chiếu chính thức tại `docs/FUNCTION_DECOMPOSITION_ALIGNMENT_2026-04-15.md`
- [x] Cập nhật lại `USECASE_MASTER_CURRENT_AND_PLAN.md` để phản ánh `class_schedule`, `lesson`, `attendance`, `lesson_summary`, `academic_record`, `leave_request`, `predictive`
- [x] Cập nhật lại `USE_CASE_MODELING_PACKAGE_EDUCENTER_REFINED.md` cho các trạng thái vừa hoàn tất ở `G3 -> G6`
- [x] Chốt nhánh còn `Partial/Missing` cần làm tiếp: `teacher schedule FE`, `student timetable API + FE`, lifecycle `program publish/archive`, benchmark UI riêng
- [x] Chốt nhánh ngoài scope chính: `AI Audit`, `Compliance mở rộng`, `DevTools`, `Consultation`, `AI Assistant`

**Mục tiêu đóng task:**
- Backlog, codebase và sơ đồ phân rã chức năng khớp nhau trước khi bước sang giai đoạn test/báo cáo.

---

### Giai đoạn H - Kiểm thử, báo cáo, bảo vệ

#### Task H1. Kiểm thử scheduling
- [ ] Test dữ liệu nhỏ
- [ ] Test dữ liệu trung bình
- [ ] Test dữ liệu lớn hơn
- [ ] Test conflict cases
- [ ] Test benchmark API
- [ ] Test preview/commit sau khi chọn solver chính

#### Task H2. Kiểm thử predictive analytics
- [ ] Kiểm tra pipeline train/evaluate
- [ ] Kiểm tra metric trên tập validation/test
- [ ] Kiểm tra prediction API
- [ ] Kiểm tra UI cảnh báo

#### Task H3. Cập nhật báo cáo đồ án
- [ ] Viết lại mục tiêu dự án theo hướng mới
- [ ] Bổ sung chương benchmark scheduling
- [ ] Bổ sung phần `Shift` và mô hình thời gian
- [ ] Bổ sung phần predictive analytics classification
- [ ] Bổ sung đánh giá và lựa chọn thuật toán

#### Task H4. Chuẩn bị bảo vệ
- [ ] Tạo demo flow scheduling benchmark
- [ ] Tạo demo flow scheduling production-like
- [ ] Tạo demo flow predictive analytics
- [ ] Chuẩn bị slide và câu hỏi phản biện

---

### Giai đoạn I - Portal Giáo viên (Teacher Portal)

> **Entity nền đã có:** `Lesson`, `Attendance`, `LessonSummary`, `AcademicRecord`, `LeaveRequest`
>
> **Nav stub đã có:** `teacher-schedule`, `teacher-attendance`, `teacher-journal`
>
> **Chưa có:** usecase, controller, repository logic, và frontend page thực tế cho tất cả các tác vụ bên dưới.

---

#### Task I1. Backend - Lịch giảng dạy của giáo viên

**Mục tiêu:** Giáo viên xem được lịch dạy cá nhân theo tuần/tháng.

- [ ] Tạo usecase `GetTeacherLessons` (lọc theo `teacher_id`, khoảng ngày, lớp)
- [ ] Tạo repository query join `lessons` ← `classes` ← `class_schedules` theo teacher
- [ ] Tạo controller + route `GET /api/teacher/lessons` (role guard: TEACHER)
- [ ] Trả về thông tin ca học `shift`, phòng, lớp, ngày giờ

**Mục tiêu đóng task:** API trả được danh sách buổi học cho giáo viên đang đăng nhập.

---

#### Task I2. Backend - Điểm danh sinh viên theo buổi học

**Mục tiêu:** Giáo viên chấm điểm danh cho từng sinh viên trong buổi học.

- [ ] Tạo usecase `MarkAttendance` (tạo hoặc cập nhật `Attendance` theo `lesson_id` + `student_id`)
- [ ] Định nghĩa `AttendanceStatus`: `0=Vắng`, `1=Có mặt`, `2=Muộn`, `3=Xin phép`
- [ ] Tạo usecase `GetAttendanceByLesson` (lấy toàn bộ danh sách điểm danh của một buổi)
- [ ] Tạo usecase `GetAttendanceSummaryByStudent` (thống kê chuyên cần theo học sinh trong lớp)
- [ ] Tạo controller + route:
  - [ ] `GET /api/teacher/lessons/:lesson_id/attendance`
  - [ ] `POST /api/teacher/lessons/:lesson_id/attendance` (submit cả buổi dạng mảng)
  - [ ] `PUT /api/teacher/lessons/:lesson_id/attendance/:student_id` (cập nhật từng dòng)
- [ ] Kiểm tra: chỉ giáo viên phụ trách buổi học mới được điểm danh

**Mục tiêu đóng task:** Sau mỗi buổi, `attendance` có dữ liệu đầy đủ cho toàn bộ học sinh.

---

#### Task I3. Backend - Tổng kết buổi học (Sổ đầu bài)

**Mục tiêu:** Giáo viên ghi nhận nội dung dạy, bài tập, ghi chú sau mỗi buổi.

- [ ] Tạo usecase `UpsertLessonSummary` (tạo mới hoặc cập nhật `LessonSummary` theo `lesson_id`)
- [ ] Tạo usecase `GetLessonSummary` (lấy theo `lesson_id`)
- [ ] Tạo controller + route:
  - [ ] `GET /api/teacher/lessons/:lesson_id/summary`
  - [ ] `PUT /api/teacher/lessons/:lesson_id/summary`
- [ ] Field cần hỗ trợ: `topic`, `lesson_content`, `class_feedback`, `homework`, `homework_deadline`, `teacher_notes`
- [ ] Kiểm tra: chỉ giáo viên phụ trách mới được ghi

**Mục tiêu đóng task:** Mỗi buổi học có thể có đúng một `LessonSummary` đầy đủ thông tin.

---

#### Task I4. Backend - Vào điểm và đánh giá kết quả học tập

**Mục tiêu:** Giáo viên nhập điểm bài tập, thái độ, kết quả cho từng học sinh sau buổi học.

- [ ] Tạo usecase `UpsertAcademicRecord` (tạo/cập nhật `AcademicRecord` theo `lesson_summary_id` + `student_id`)
- [ ] Tạo usecase `GetAcademicRecordsByLesson` (lấy tất cả kết quả của một buổi)
- [ ] Tạo usecase `GetAcademicRecordsByStudent` (lịch sử điểm của một học sinh trong lớp)
- [ ] Tạo usecase `FinalizeAcademicRecord` (set `is_completed = true`, tính `total_score`)
- [ ] Tạo controller + route:
  - [ ] `GET /api/teacher/lessons/:lesson_id/records`
  - [ ] `PUT /api/teacher/lessons/:lesson_id/records/:student_id`
  - [ ] `POST /api/teacher/lessons/:lesson_id/records/finalize`
- [ ] Field cần hỗ trợ: `homework_completed`, `homework_score`, `attitude_rating`, `participation_score`, `personal_comment`

**Mục tiêu đóng task:** Giáo viên nhập được điểm số và đánh giá cho từng học sinh, khóa kết quả sau khi hoàn tất.

---

#### Task I5. Backend - Duyệt đơn xin nghỉ

**Mục tiêu:** Giáo viên xem danh sách đơn xin nghỉ liên quan đến lớp mình và duyệt/từ chối.

- [ ] Tạo usecase `ListLeaveRequestsForTeacher` (lọc theo lớp do giáo viên phụ trách, lọc theo `status`)
- [ ] Tạo usecase `ApproveLeaveRequest` (set `status=APPROVED`, ghi `approved_by_id`, `approved_at`)
- [ ] Tạo usecase `RejectLeaveRequest` (set `status=REJECTED`, ghi `rejection_reason`)
- [ ] Tạo controller + route:
  - [ ] `GET /api/teacher/leave-requests` (có filter `class_id`, `status`, `student_id`)
  - [ ] `POST /api/teacher/leave-requests/:id/approve`
  - [ ] `POST /api/teacher/leave-requests/:id/reject`
- [ ] Kiểm tra: chỉ duyệt đơn thuộc lớp giáo viên phụ trách

**Mục tiêu đóng task:** Vòng đời đơn xin phép khép kín: sinh viên tạo → giáo viên duyệt/từ chối.

---

#### Task I6. Frontend - Trang lịch giảng dạy (`/app/teacher/schedule`)

- [ ] Hiển thị danh sách buổi học theo tuần (calendar view hoặc list view theo ngày)
- [ ] Hiển thị: tên lớp, tên ca (`shift.name`), phòng, thời gian
- [ ] Click vào buổi học → điều hướng sang chi tiết buổi học (điểm danh + tổng kết)
- [ ] Tích hợp API `GET /api/teacher/lessons`

**Mục tiêu đóng task:** Giáo viên thấy được lịch dạy cá nhân trên UI.

---

#### Task I7. Frontend - Trang điểm danh (`/app/teacher/attendance`)

- [ ] Chọn lớp → chọn buổi học → hiển thị danh sách học sinh trong buổi
- [ ] Cho mỗi học sinh chọn trạng thái: Có mặt / Vắng / Muộn / Xin phép
- [ ] Hiển thị đơn xin phép liên quan (nếu có) bên cạnh tên học sinh
- [ ] Nút "Lưu điểm danh" → submit toàn bộ danh sách
- [ ] Tab "Theo dõi chuyên cần" → bảng tổng hợp số buổi có mặt/vắng/muộn của từng học sinh trong lớp
- [ ] Tích hợp API `GET/POST /api/teacher/lessons/:lesson_id/attendance`

**Mục tiêu đóng task:** Giáo viên điểm danh và xem được tổng hợp chuyên cần trên UI.

---

#### Task I8. Frontend - Trang sổ đầu bài / tổng kết buổi học (`/app/teacher/journal`)

- [ ] Chọn lớp → chọn buổi → hiển thị form tổng kết buổi học
- [ ] Form: `topic`, `lesson_content`, `homework`, `homework_deadline`, `teacher_notes`, `class_feedback`
- [ ] Auto-save hoặc nút Lưu tổng kết
- [ ] Danh sách điều hướng các buổi đã/chưa có tổng kết
- [ ] Tích hợp API `GET/PUT /api/teacher/lessons/:lesson_id/summary`

**Mục tiêu đóng task:** Giáo viên ghi chép sổ đầu bài cho từng buổi học trên UI.

---

#### Task I9. Frontend - Vào điểm học sinh (tích hợp trong chi tiết buổi học)

- [ ] Trong màn chi tiết buổi học: tab "Kết quả học tập"
- [ ] Bảng danh sách học sinh với input: `homework_completed`, `homework_score`, `attitude_rating`, `participation_score`, `personal_comment`
- [ ] Nút "Chốt kết quả" → gọi finalize, khóa input sau khi hoàn tất
- [ ] Tích hợp API `GET/PUT /api/teacher/lessons/:lesson_id/records`

**Mục tiêu đóng task:** Giáo viên nhập và confirm điểm số học sinh cho từng buổi học.

---

#### Task I10. Frontend - Trang duyệt đơn xin nghỉ (tích hợp trong portal giáo viên)

- [ ] Thêm nav item `teacher-leaves`: "Đơn xin phép" vào `nav.ts`
- [ ] Trang `/app/teacher/leaves`: danh sách đơn xin phép (filter: lớp, trạng thái, học sinh)
- [ ] Hiển thị: tên học sinh, lớp, ngày xin nghỉ, lý do, tài liệu đính kèm
- [ ] Nút Duyệt / Từ chối + ô nhập lý do từ chối
- [ ] Tích hợp API `GET /api/teacher/leave-requests`, `POST .../approve`, `POST .../reject`

**Mục tiêu đóng task:** Giáo viên duyệt/từ chối đơn xin phép trực tiếp trên UI.

---

### Giai đoạn J - Portal Sinh viên (Student Portal)

> **Entity nền đã có:** `Lesson`, `Attendance`, `AcademicRecord`, `LeaveRequest`, `Enrollment`
>
> **Nav stub đã có:** `student-timetable`, `student-results`, `student-leaves`
>
> **Chưa có:** usecase, controller, frontend page thực tế.

---

#### Task J1. Backend - Thời khóa biểu học sinh

**Mục tiêu:** Học sinh xem được lịch học của mình.

- [ ] Tạo usecase `GetStudentTimetable` (lọc theo `student_id`, khoảng ngày)
- [ ] Query join: `enrollments` → `classes` → `lessons` → `class_schedules` → `shifts` → `rooms`
- [ ] Tạo controller + route `GET /api/student/timetable` (role guard: STUDENT)
- [ ] Trả về danh sách buổi học: tên lớp, giáo viên, ca học, phòng, ngày

**Mục tiêu đóng task:** Học sinh xem được thời khóa biểu cá nhân qua API.

---

#### Task J2. Backend - Điểm danh của sinh viên

**Mục tiêu:** Học sinh xem được lịch sử chuyên cần của bản thân.

- [ ] Tạo usecase `GetMyAttendance` (lọc theo `student_id`, `class_id`, khoảng ngày)
- [ ] Tính tổng hợp: tổng buổi, số buổi có mặt, vắng, muộn, xin phép
- [ ] Thêm cảnh báo nếu số buổi vắng vượt ngưỡng (ví dụ > 20%)
- [ ] Tạo controller + route `GET /api/student/attendance`

**Mục tiêu đóng task:** Học sinh tự theo dõi chuyên cần và nhận cảnh báo nếu cần.

---

#### Task J3. Backend - Kết quả học tập của sinh viên

**Mục tiêu:** Học sinh xem điểm số và đánh giá từng buổi học.

- [ ] Tạo usecase `GetMyAcademicRecords` (lọc theo `student_id`, `class_id`, khoảng ngày)
- [ ] Join: `academic_records` → `lesson_summaries` → `lessons`
- [ ] Tính điểm trung bình tổng hợp theo lớp
- [ ] Tạo controller + route `GET /api/student/academic-records`

**Mục tiêu đóng task:** Học sinh xem được điểm số và nhận xét từng buổi học.

---

#### Task J4. Backend - Tạo đơn xin nghỉ

**Mục tiêu:** Học sinh tự tạo đơn xin phép nghỉ/muộn.

- [ ] Tạo usecase `CreateLeaveRequest` (validate: `apply_date`, `class_id` phải thuộc lớp đã enroll)
- [ ] Tạo usecase `GetMyLeaveRequests` (lịch sử đơn của học sinh, filter theo `status`, `class_id`)
- [ ] Tạo usecase `CancelLeaveRequest` (hủy đơn nếu còn PENDING)
- [ ] Tạo controller + route:
  - [ ] `GET /api/student/leave-requests`
  - [ ] `POST /api/student/leave-requests`
  - [ ] `DELETE /api/student/leave-requests/:id` (cancel nếu PENDING)
- [ ] Hỗ trợ upload `documents` (URL danh sách)

**Mục tiêu đóng task:** Học sinh tạo và theo dõi được đơn xin phép của mình.

---

#### Task J5. Backend - Dự báo kết quả kém (AT_RISK cho sinh viên)

**Mục tiêu:** Học sinh xem kết quả dự báo nguy cơ học kém của bản thân.

- [ ] Tạo usecase `GetMyAtRiskPrediction` (gọi prediction engine từ module `predictive`)
- [ ] Trả về: `risk_label`, `risk_score`, top features ảnh hưởng
- [ ] Tạo controller + route `GET /api/student/at-risk`
- [ ] Kết nối với pipeline ML đã có ở Task F3

**Mục tiêu đóng task:** Học sinh xem được mức độ nguy cơ học kém và lý do cảnh báo.

---

#### Task J6. Frontend - Trang thời khóa biểu sinh viên (`/app/student/timetable`)

- [ ] Hiển thị lịch học theo tuần (calendar grid hoặc danh sách theo ngày)
- [ ] Mỗi buổi: tên lớp, ca học (`shift.name`), phòng, giáo viên, trạng thái (sắp tới / đã qua)
- [ ] Hiển thị đơn xin phép tương ứng nếu có
- [ ] Tích hợp API `GET /api/student/timetable`

**Mục tiêu đóng task:** Học sinh xem thời khóa biểu trực tiếp trên UI.

---

#### Task J7. Frontend - Theo dõi điểm danh cá nhân (tích hợp trong timetable hoặc trang riêng)

- [ ] Bảng tổng hợp chuyên cần: tổng buổi / có mặt / vắng / muộn / xin phép
- [ ] Hiển thị chi tiết từng buổi và trạng thái điểm danh
- [ ] Banner cảnh báo nếu tỷ lệ vắng vượt ngưỡng
- [ ] Tích hợp API `GET /api/student/attendance`

**Mục tiêu đóng task:** Học sinh chủ động theo dõi chuyên cần của mình.

---

#### Task J8. Frontend - Trang kết quả học tập (`/app/student/results`)

- [ ] Danh sách lớp đã đăng ký → chọn lớp → xem điểm theo từng buổi
- [ ] Bảng điểm: `homework_score`, `attitude_rating`, `participation_score`, `total_score`, `personal_comment`
- [ ] Hiển thị điểm trung bình tổng hợp, biểu đồ tiến độ nếu có thể
- [ ] Tích hợp API `GET /api/student/academic-records`

**Mục tiêu đóng task:** Học sinh xem điểm và nhận xét của giáo viên qua UI.

---

#### Task J9. Frontend - Trang đơn xin nghỉ (`/app/student/leaves`)

- [ ] Form tạo đơn: chọn lớp, ngày, loại đơn (`LEAVE`/`LATE`/`EARLY`), lý do, đính kèm file
- [ ] Danh sách đơn đã tạo: hiển thị trạng thái `PENDING / APPROVED / REJECTED`, lý do từ chối
- [ ] Nút Hủy đơn nếu còn PENDING
- [ ] Tích hợp API `GET/POST /api/student/leave-requests`

**Mục tiêu đóng task:** Học sinh tạo và theo dõi đơn xin phép hoàn toàn trên UI.

---

#### Task J10. Frontend - Cảnh báo dự báo kết quả kém (tích hợp vào trang kết quả hoặc overview)

- [ ] Hiển thị badge / card cảnh báo AT_RISK nếu `risk_label = AT_RISK`
- [ ] Hiển thị `risk_score` và các yếu tố ảnh hưởng chính (top features)
- [ ] Gợi ý hành động: liên hệ giáo viên, cải thiện chuyên cần, hoàn thành bài tập
- [ ] Tích hợp API `GET /api/student/at-risk`

**Mục tiêu đóng task:** Học sinh nhận cảnh báo sớm và có thông tin hướng dẫn cải thiện.

---

## 6. Thứ tự ưu tiên triển khai

Thứ tự thực hiện từ bây giờ:

1. `Task A1 -> A2`
2. `Task B1 -> B3`
3. `Task C1 -> C3`
4. `Task D1 -> D3`
5. `Task E1 -> E3`
6. `Task F1 -> F4`
7. `Task G1 -> G7`
8. **`Task I1 -> I10`** _(Portal Giáo viên - thực hiện song song với G sau khi G4-G6 khởi động)_
9. **`Task J1 -> J10`** _(Portal Sinh viên - thực hiện sau I hoặc song song từ J1)_
10. `Task H1 -> H4`

Nguyên tắc:

- Không làm predictive analytics trước khi scheduling ổn định.
- Không mở rộng lại AI Audit / Compliance / DevTools vào backlog chính.
- Không mở rộng thêm module ngoài scope khi backlog chính chưa xong.
- Mỗi task lớn chỉ triển khai sau khi task trước đã được review/chấp nhận.
- Task I (Teacher Portal) phụ thuộc vào dữ liệu `lesson` từ commit scheduling (G4).
- Task J5 (AT_RISK cho sinh viên) phụ thuộc vào pipeline ML từ Task F3.

---

## 7. Ghi chú tiến độ

Ký hiệu:

- `[x]` đã hoàn thành
- `[/]` đang thực hiện
- `[ ]` chưa bắt đầu

Tài liệu liên quan:

- [docs/ke_hoach_phan_hoi_gvhd_2026-04-09.md](/Users/hant/golang/doan/docs/ke_hoach_phan_hoi_gvhd_2026-04-09.md)
- [docs/modeling/scheduling/SCHEDULING_CORE_LOGIC_CURRENT.md](/Users/hant/golang/doan/docs/modeling/scheduling/SCHEDULING_CORE_LOGIC_CURRENT.md)

**Ngày tạo file task gốc:** 2026-02-06

**Ngày tái cấu trúc backlog:** 2026-04-09
