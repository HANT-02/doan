# ĐỐI CHIẾU SƠ ĐỒ PHÂN RÃ CHỨC NĂNG VỚI CODEBASE EDUCENTER

**Ngày chốt:** 2026-04-15

**Mục tiêu**

- đối chiếu từng nhánh trong sơ đồ phân rã chức năng với backend/frontend hiện có;
- đánh dấu rõ `Implemented / Partial / Missing / Out of main scope`;
- chốt nhánh nào tiếp tục nằm trong scope chính của đồ án.

**Nguồn đối chiếu**

- [PROJECT_TASKS.md](/Users/hant/golang/doan/PROJECT_TASKS.md)
- [docs/USECASE_MASTER_CURRENT_AND_PLAN.md](/Users/hant/golang/doan/docs/USECASE_MASTER_CURRENT_AND_PLAN.md)
- route/controller trong `cmd/http/controllers`
- route/nav/page trong `frontend/src`

## 1. Bảng map tổng

| Nhánh trong sơ đồ | Backend/API hiện có | Frontend/UI hiện có | Trạng thái | Quyết định scope | Ghi chú |
|---|---|---|---|---|---|
| Quản lý tài khoản | `GET/POST/PUT /v1/users`, reset password | `/app/admin/accounts` | Implemented | Giữ trong scope chính | Chỉ `ADMIN`, `SUPER_ADMIN` thao tác |
| Quản lý học sinh | CRUD `students` | `/app/admin/students` | Implemented | Giữ trong scope chính | |
| Quản lý giáo viên | CRUD `teachers`, timetable, teaching hours | `/app/admin/teachers`, detail/edit | Implemented | Giữ trong scope chính | |
| Quản lý khóa học | CRUD `courses` | `/app/admin/courses` | Implemented | Giữ trong scope chính | |
| Quản lý chương trình đào tạo | CRUD `programs`, link/unlink courses | `/app/admin/programs` | Partial | Giữ trong scope chính | Chưa có lifecycle publish/archive |
| Quản lý phòng học | CRUD `rooms` | `/app/admin/rooms` | Implemented | Giữ trong scope chính | |
| Quản lý ca học | CRUD `shifts`, toggle active | `/app/admin/shifts` | Implemented | Giữ trong scope chính | |
| Quản lý lớp học | CRUD `classes`, detail, assign teacher | `/app/admin/classes` | Implemented | Giữ trong scope chính | |
| Quản lý ghi danh lớp học | add/remove enrollment, available students | từ detail lớp | Implemented | Giữ trong scope chính | |
| Quản lý lịch tuần lớp | `GET/POST/DELETE /v1/classes/:id/schedules` | tab lịch tuần trong detail lớp | Implemented | Giữ trong scope chính | Chưa có endpoint update riêng; đang dùng create/delete |
| Xếp lịch học | preview, conflict, commit | `/app/admin/scheduling` | Implemented | Giữ trong scope chính | Solver mặc định đã chốt là `CP-SAT` |
| Benchmark thuật toán xếp lịch | benchmark API + CLI study | chưa có trang riêng, dùng admin scheduling/report docs | Partial | Giữ trong scope chính | Backend hoàn chỉnh, FE chuyên biệt chưa làm |
| Quản lý buổi học | `GET /v1/lessons`, `GET /v1/lessons/:id` | `/app/admin/lessons`, detail lesson | Implemented | Giữ trong scope chính | Dữ liệu sinh từ commit scheduling |
| Quản lý điểm danh | `GET/PUT /v1/lessons/:id/attendance` | admin lesson detail, `/app/teacher/attendance` | Implemented | Giữ trong scope chính | Teacher bị giới hạn theo lesson phụ trách |
| Quản lý tổng kết buổi học | `GET/PUT /v1/lessons/:id/summary` | admin lesson detail, `/app/teacher/journal` | Implemented | Giữ trong scope chính | |
| Quản lý kết quả học tập | `GET/PUT/POST finalize /v1/lessons/:id/records`, `GET /v1/academic-records/my` | admin lesson detail, teacher journal, `/app/student/results` | Implemented | Giữ trong scope chính | |
| Quản lý đơn xin phép | `GET/POST/DELETE /v1/leave-requests`, approve/reject | `/app/admin/leaves`, `/app/teacher/leaves`, `/app/student/leaves` | Implemented | Giữ trong scope chính | Flow actor-based |
| Dự báo học sinh AT_RISK | training runtime tối giản + `GET /v1/predictive/at-risk/*` | `/app/admin/predictive` | Implemented | Giữ trong scope chính | Dùng DB loader + model nhẹ |
| Lịch giảng dạy giáo viên | `GET /v1/teachers/:id/timetable` | `/app/teacher/schedule` còn placeholder | Partial | Ghi nhận làm tiếp | Backend có, FE teacher chưa nối |
| Thời khóa biểu học sinh/phụ huynh | chưa có endpoint riêng theo actor | `/app/student/timetable` còn placeholder | Missing | Ghi nhận làm tiếp | Nằm ở giai đoạn J |
| Tài liệu giảng dạy / kiểm duyệt compliance | backend + một phần UI vẫn còn | teacher documents, compliance pages | Partial | Ngoài scope đồ án chính | Không ưu tiên trong backlog chính sau G |
| AI Audit / DevTools / AI Assistant / Consulting | chỉ còn vết placeholder rời rạc | placeholder hoặc đã bỏ khỏi admin flow | Out of main scope | Loại khỏi scope chính | Không dùng để vẽ sơ đồ chính thức |

## 2. Kết luận chốt cho G7

### 2.1 Nhánh đã khớp tốt giữa sơ đồ, backlog và code

- quản lý tài khoản;
- học sinh, giáo viên, khóa học, chương trình, phòng học, ca học, lớp học, ghi danh;
- lịch tuần lớp;
- scheduling preview/commit;
- lesson;
- attendance, lesson summary, academic record, leave request;
- predictive alerts.

### 2.2 Nhánh còn lệch nhưng đã xác định rõ nguyên nhân

- `Benchmark thuật toán xếp lịch`: backend/CLI đã xong, chưa có UI riêng.
- `Lịch giảng dạy giáo viên`: backend có, FE teacher vẫn đang dùng placeholder.
- `Thời khóa biểu học sinh/phụ huynh`: chưa có API actor-based và chưa có page thật.
- `Chương trình đào tạo`: lifecycle publish/archive chưa làm.

### 2.3 Nhánh ghi nhận nhưng không đưa vào scope chính thức để vẽ sơ đồ đồ án

- AI Audit / compliance mở rộng;
- DevTools;
- AI Assistant / chatbot;
- tư vấn tuyển sinh / lead intake.

## 3. Đề xuất dùng cho giai đoạn tiếp theo

- Khi vẽ sơ đồ overview chính thức, giữ các nhánh trong `2.1` và có thể thêm `Benchmark thuật toán xếp lịch` như use case hỗ trợ quản trị.
- Khi vẽ sơ đồ chi tiết cho teacher/student portal, ghi rõ `teacher schedule` và `student timetable` đang là nhánh `Partial/Missing`, tránh mô tả như đã hoàn tất.
- Khi chuẩn bị báo cáo/bảo vệ, dùng file này như bảng đối chiếu cuối cùng giữa sơ đồ và code thay vì suy từ backlog cũ.
