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
   - **không làm module AI Audit**,
   - **không làm nhánh Compliance/Audit/Chatbot** trong backlog chính,
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

---

## 3. Ngoài phạm vi

Các nhánh dưới đây được bỏ khỏi file task để giữ backlog rõ ràng:

- Module AI Audit / kiểm duyệt tài liệu
- Compliance dashboard cho AI Audit
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

### 4.2. Scheduling baseline hiện có

- [x] Có màn hình scheduling preview trên frontend
- [x] Có preview backend và commit xuống `lesson`
- [x] Có conflict messaging cơ bản
- [x] Có tài liệu mô tả logic hiện tại tại:
  - [x] `docs/modeling/scheduling/SCHEDULING_CORE_LOGIC_CURRENT.md`

### 4.3. Tài liệu định hướng mới

- [x] Đã chốt ý tưởng hướng mới tại:
  - [x] `docs/ke_hoach_phan_hoi_gvhd_2026-04-09.md`

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
- [ ] Tạo admin API benchmark
- [ ] Cho phép chạy cùng input qua 3 solver
- [ ] Trả bảng metric so sánh
- [ ] Kiểm tra quyền admin nội bộ

**Mục tiêu đóng task:**
- Có admin API benchmark dùng được để lấy số liệu báo cáo.

#### Task E2. Chạy benchmark và chọn solver chính
- [ ] Chuẩn bị bộ dữ liệu benchmark
- [ ] Chạy benchmark nhiều lần
- [ ] Ghi feasibility / hard violations / soft score / runtime
- [ ] Chọn solver tốt nhất
- [ ] Tài liệu hóa quyết định chọn solver

**Mục tiêu đóng task:**
- Có bằng chứng thực nghiệm để chọn solver chính.

#### Task E3. Gắn solver tốt nhất vào scheduling API chính
- [ ] Inject solver chính vào use case scheduling
- [ ] Giữ benchmark API tách riêng
- [ ] Kiểm tra lại preview
- [ ] Kiểm tra lại commit
- [ ] Kiểm tra luồng FE-BE sau khi đổi solver

**Mục tiêu đóng task:**
- Scheduling production-like dùng đúng solver đã chọn.

---

### Giai đoạn F - Predictive analytics

#### Task F1. Chốt dữ liệu đầu vào cho `AT_RISK classification`
- [ ] Chốt nguồn dữ liệu:
  - [ ] student
  - [ ] attendance
  - [ ] grade
  - [ ] class/course enrollment
  - [ ] tín hiệu vận hành cần thiết
- [ ] Chốt label `AT_RISK`
- [ ] Chốt feature set ban đầu

**Mục tiêu đóng task:**
- Có dataset definition rõ để bắt đầu pipeline ML.

#### Task F2. Pipeline ML trong backend hiện tại
- [ ] Tạo pipeline feature engineering
- [ ] Tạo pipeline train/test split
- [ ] Tạo huấn luyện `Logistic Regression`
- [ ] Tạo huấn luyện `Random Forest`
- [ ] Tạo huấn luyện `LightGBM/XGBoost`
- [ ] Tạo metric evaluation

**Mục tiêu đóng task:**
- Backend hiện tại train/evaluate được các mô hình classification cơ bản.

#### Task F3. Prediction API + model metadata
- [ ] Thiết kế lưu model metadata
- [ ] Tạo endpoint dự báo `AT_RISK`
- [ ] Tạo output explanation cơ bản
- [ ] Tạo cơ chế versioning tối thiểu cho model

**Mục tiêu đóng task:**
- Có API dự báo dùng được trong hệ thống.

#### Task F4. Frontend cảnh báo sinh viên nguy cơ học kém
- [ ] Tạo màn hình danh sách sinh viên `AT_RISK`
- [ ] Hiển thị score/label
- [ ] Hiển thị insight hoặc lý do chính
- [ ] Kiểm tra flow end-to-end

**Mục tiêu đóng task:**
- Có UI demo được giá trị của predictive analytics.

---

### Giai đoạn G - Kiểm thử, báo cáo, bảo vệ

#### Task G1. Kiểm thử scheduling
- [ ] Test dữ liệu nhỏ
- [ ] Test dữ liệu trung bình
- [ ] Test dữ liệu lớn hơn
- [ ] Test conflict cases
- [ ] Test benchmark API
- [ ] Test preview/commit sau khi chọn solver chính

#### Task G2. Kiểm thử predictive analytics
- [ ] Kiểm tra pipeline train/evaluate
- [ ] Kiểm tra metric trên tập validation/test
- [ ] Kiểm tra prediction API
- [ ] Kiểm tra UI cảnh báo

#### Task G3. Cập nhật báo cáo đồ án
- [ ] Viết lại mục tiêu dự án theo hướng mới
- [ ] Bổ sung chương benchmark scheduling
- [ ] Bổ sung phần `Shift` và mô hình thời gian
- [ ] Bổ sung phần predictive analytics classification
- [ ] Bổ sung đánh giá và lựa chọn thuật toán

#### Task G4. Chuẩn bị bảo vệ
- [ ] Tạo demo flow scheduling benchmark
- [ ] Tạo demo flow scheduling production-like
- [ ] Tạo demo flow predictive analytics
- [ ] Chuẩn bị slide và câu hỏi phản biện

---

## 6. Thứ tự ưu tiên triển khai

Thứ tự thực hiện từ bây giờ:

1. `Task A1 -> A2`
2. `Task B1 -> B3`
3. `Task C1 -> C3`
4. `Task D1 -> D3`
5. `Task E1 -> E3`
6. `Task F1 -> F4`
7. `Task G1 -> G4`

Nguyên tắc:

- Không làm predictive analytics trước khi scheduling ổn định.
- Không mở rộng thêm module ngoài scope khi backlog chính chưa xong.
- Mỗi task lớn chỉ triển khai sau khi task trước đã được review/chấp nhận.

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
