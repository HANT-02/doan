# Quản lý Task - Hệ thống Quản lý Trung tâm Dạy thêm tích hợp AI

**Đề tài:** Xây dựng hệ thống quản lý trung tâm dạy thêm tích hợp AI hỗ trợ kiểm soát chất lượng đào tạo

**Mô hình:** Single-tenant (đơn thuê bao)

**Công nghệ:** Golang (Backend), ReactJS (Frontend)

**Tuân thủ:** Thông tư 29/2024/TT-BGDĐT

---

## 📋 Giai đoạn 1: Phân tích & Đặc tả hệ thống (Tuần 3-4)

### Task 1.1: Đặc tả chi tiết bài toán Xếp lịch (CSP) *(Hoàn thành: Implemented in docs/modeling/scheduling/)*
- [x] Xác định tập hợp biến số
  - [x] Định nghĩa biến Lớp học (Class)
  - [x] Định nghĩa biến Giáo viên (Teacher)
  - [x] Định nghĩa biến Phòng học (Room)
  - [x] Định nghĩa biến Khung giờ (TimeSlot)
- [x] Định nghĩa chi tiết Ràng buộc cứng (Hard Constraints)
  - [x] Chống trùng lịch giáo viên
  - [x] Chống trùng lịch phòng học
  - [x] Giới hạn sĩ số phòng học
  - [x] Chặn khung giờ sau 22h (Thông tư 29)
  - [x] Ràng buộc về thời gian làm việc tối đa của giáo viên
- [x] Định nghĩa Ràng buộc mềm (Soft Constraints)
  - [x] Ưu tiên lịch dạy liên tiếp cho giáo viên
  - [x] Tối ưu khoảng cách giữa các buổi học
  - [x] Ưu tiên phòng học phù hợp với môn học

### Task 1.2: Thiết kế sơ đồ Use Case & Luồng nghiệp vụ *(Hoàn thành: Implemented in docs/modeling/usecase/*)*
- [x] Vẽ sơ đồ Use Case tổng quan
  - [x] Use case cho Admin
  - [x] Use case cho Giáo viên
  - [x] Use case cho Học sinh/Phụ huynh
  - [x] Use case cho Compliance Officer
  - [x] Use case tương tác với AI Agent
- [x] Mô tả luồng nghiệp vụ phê duyệt tài liệu
  - [x] Luồng tải file của Giáo viên
  - [x] Luồng AI quét (OCR & Inference)
  - [x] Luồng gán nhãn tự động
  - [x] Luồng Compliance Officer phê duyệt
  - [x] Luồng phản hồi kết quả cho Giáo viên

### Task 1.3: Thiết kế cơ sở dữ liệu (ERD)
- [ ] Thiết kế bảng dữ liệu Core
  - [ ] Bảng User (Người dùng)
  - [ ] Bảng Teacher (Giáo viên)
  - [ ] Bảng Student (Học sinh)
  - [ ] Bảng Class (Lớp học)
  - [ ] Bảng Room (Phòng học)
  - [ ] Bảng Program (Chương trình đào tạo)
  - [ ] Bảng Course (Khóa học)
  - [ ] Bảng Lesson (Buổi học)
  - [ ] Bảng Schedule (Thời khóa biểu)
- [ ] Thiết kế bảng phục vụ AI Audit
  - [ ] Bảng Material (Tài liệu)
  - [ ] Bảng AuditLog (Lịch sử quét)
  - [ ] Bảng Label (Nhãn an toàn/cảnh báo)
  - [ ] Bảng AIAnalysisResult (Kết quả phân tích AI)
- [ ] Thiết kế bảng bổ trợ
  - [ ] Bảng Attendance (Điểm danh)
  - [ ] Bảng Grade (Điểm số)
  - [ ] Bảng Payment (Học phí)
  - [ ] Bảng Notification (Thông báo)
- [ ] Vẽ sơ đồ ERD hoàn chỉnh
- [ ] Tài liệu hóa các mối quan hệ và ràng buộc

---

## 🔧 Giai đoạn 2: Phát triển Backend (Golang) & Cơ sở dữ liệu (Tuần 5-11)

### Task 2.1: Khởi tạo Project & Cấu trúc thư mục
- [ ] Thiết lập cấu trúc Clean Architecture
  - [ ] Tạo thư mục `entities` (Domain layer)
  - [ ] Tạo thư mục `repositories` (Data layer)
  - [ ] Tạo thư mục `services` (Business logic)
  - [ ] Tạo thư mục `usecases` (Application layer)
  - [ ] Tạo thư mục `infrastructure` (External services)
  - [ ] Tạo thư mục `delivery` (API handlers)
- [ ] Cấu hình PostgreSQL
  - [ ] Thiết lập connection pool
  - [ ] Cấu hình migration tool (golang-migrate)
  - [ ] Tạo file migration cho các bảng
- [ ] Dockerization
  - [ ] Tạo Dockerfile cho Backend
  - [ ] Tạo docker-compose.yml (Backend + PostgreSQL)
  - [ ] Cấu hình environment variables
- [ ] Setup CI/CD cơ bản
  - [ ] Cấu hình linting (golangci-lint)
  - [ ] Cấu hình testing framework
- [x] Cấu hình Dependency Injection (Google Wire) *(Hoàn thành: Implemented in cmd/http/wire.go)*

### Task 2.2: Xây dựng Module Quản lý Đào tạo (Core)
- [x] API quản lý Giáo viên *(Hoàn thành: 2026-02-06, Implemented in: cmd/http/controllers/teacher)*
  - [x] CRUD Teacher (POST/GET/PUT/DELETE `/api/v1/teachers`)
  - [x] API lấy danh sách giáo viên theo bộ lọc (GET `/api/v1/teachers?search=&status=&employment_type=&page=&limit=`)
  - [x] API lấy lịch dạy của giáo viên (GET `/api/v1/teachers/:id/timetable?from=&to=`)
  - [x] API thống kê giờ dạy (GET `/api/v1/teachers/:id/stats/teaching-hours?from=&to=&group_by=`)
- [x] API quản lý Học sinh *(Hoàn thành: 2026-02-25, Implemented in: cmd/http/controllers/student)*
  - [x] CRUD Student
  - [ ] API đăng ký khóa học
  - [ ] API xem lịch học
  - [ ] API xem điểm và chuyên cần
- [x] API quản lý Lớp học *(Hoàn thành: 2026-02-25, Implemented in: cmd/http/controllers/class)*
  - [x] CRUD Class
  - [x] API thêm/xóa học sinh khỏi lớp
  - [x] API gán giáo viên cho lớp
  - [x] API kiểm tra sĩ số
- [ ] API quản lý Chương trình đào tạo *(Missing: usecases, repositories, controllers)*
  - [ ] CRUD Program
  - [ ] CRUD Course
  - [ ] API liên kết Course với Program
- [x] API quản lý Phòng học *(Hoàn thành: 2026-02-25, Implemented in: cmd/http/controllers/room)*
  - [x] CRUD Room
  - [x] API kiểm tra sức chứa
  - [x] API kiểm tra tình trạng phòng học


### Task 2.3: Triển khai thuật toán Xếp lịch tự động (CSP)
- [x] Thiết kế cấu trúc dữ liệu CSP *(Hoàn thành: 2026-03-06, Implemented in: internal/usecases/scheduling/types.go, internal/usecases/scheduling/preview.go)*
  - [x] Định nghĩa Variable (Biến)
  - [x] Định nghĩa Domain (Miền giá trị)
  - [x] Định nghĩa Constraint (Ràng buộc)
- [x] Triển khai giải thuật Backtracking *(Hoàn thành: 2026-03-06, Implemented in: internal/usecases/scheduling/preview.go)*
  - [x] Implement thuật toán backtracking cơ bản
  - [x] Implement heuristic MRV (Minimum Remaining Values)
  - [ ] Implement heuristic Degree Heuristic
  - [ ] Implement heuristic LCV (Least Constraining Value)
- [x] Tích hợp Forward Checking *(Hoàn thành: 2026-03-06, Implemented in: internal/usecases/scheduling/preview.go)*
  - [x] Implement logic loại bỏ sớm các giá trị xung đột
  - [ ] Optimize performance với pruning
- [x] Implement Hard Constraints *(Hoàn thành: 2026-03-06, Implemented in: internal/usecases/scheduling/preview.go)*
  - [x] Kiểm tra trùng lịch giáo viên
  - [x] Kiểm tra trùng lịch phòng học
  - [x] Kiểm tra khung giờ cấm (sau 22h)
  - [x] Kiểm tra sức chứa phòng học
- [ ] Implement Soft Constraints
  - [ ] Tính điểm ưu tiên lịch liên tiếp
  - [ ] Tối ưu hóa khoảng cách giữa các buổi học
- [x] API Scheduling *(Hoàn thành: 2026-03-06, Implemented in: cmd/http/controllers/scheduling, internal/usecases/scheduling, internal/services/scheduling)*
  - [x] API trigger xếp lịch tự động
  - [x] API lấy kết quả xếp lịch
  - [ ] API chỉnh sửa lịch thủ công
  - [ ] API kiểm tra xung đột khi chỉnh sửa

### Task 2.4: Tích hợp AI Pipeline (Kiểm soát chất lượng)
- [ ] Xây dựng module OCR
  - [ ] Tích hợp thư viện OCR (Tesseract/Cloud Vision API)
  - [x] API upload tài liệu *(Hoàn thành: 2026-03-06, Implemented in: cmd/http/controllers/material, internal/usecases/material/upload_material.go, internal/services/audit/services.go)*
  - [x] API trích xuất văn bản từ PDF/Image *(Stub OCR: Hoàn thành 2026-03-06, Implemented in: internal/services/audit/services.go)*
  - [x] Xử lý và làm sạch văn bản *(Stub OCR: Hoàn thành 2026-03-06, Implemented in: internal/services/audit/services.go)*
- [ ] Kết nối Google Gemini API
  - [ ] Setup API credentials
  - [ ] Thiết kế prompt phát hiện nội dung không phù hợp
  - [ ] Thiết kế prompt phát hiện sai lệch kiến thức
  - [x] API gửi văn bản đến Gemini và nhận kết quả *(Stub Gemini: Hoàn thành 2026-03-06, Implemented in: internal/services/audit/services.go, internal/usecases/material/upload_material.go)*
- [ ] Xây dựng logic gán nhãn tự động
  - [x] Định nghĩa các loại nhãn (An toàn/Cảnh báo/Nguy hiểm) *(Hoàn thành: 2026-03-06, Implemented in: internal/infrastructure/database/postgres/sql_migration/23_create_material_audit_tables.up.sql, internal/entities/label.go)*
  - [x] Logic phân loại dựa trên kết quả AI *(Stub: Hoàn thành 2026-03-06, Implemented in: internal/services/audit/services.go)*
  - [x] Lưu kết quả vào database *(Hoàn thành: 2026-03-06, Implemented in: internal/usecases/material, internal/infrastructure/database/postgres/implement/material_repository.go, internal/infrastructure/database/postgres/implement/audit_log_repository.go, internal/infrastructure/database/postgres/implement/approval_decision_repository.go)*
- [x] API AI Audit *(Hoàn thành: 2026-03-06, Implemented in: cmd/http/controllers/material, internal/usecases/material, internal/services/audit, internal/infrastructure/database/postgres/sql_migration/23_create_material_audit_tables.up.sql)*
  - [x] API quét tài liệu
  - [x] API lấy lịch sử audit
  - [x] API phê duyệt/từ chối tài liệu (Compliance Officer)
  - [ ] API thống kê chất lượng tài liệu

---

## 🎨 Giai đoạn 3: Phát triển Frontend (ReactJS + TypeScript) (Tuần 7-11)

### Task 3.1: Giao diện Quản trị & Dashboard
- [x] Setup project ReactJS + TypeScript
  - [x] Cấu hình Vite/Create React App
  - [x] Cài đặt Material UI (MUI)
  - [x] Cấu hình Redux Toolkit
  - [x] Setup React Router
  - [x] Cấu hình Axios/Fetch API (RTK Query)
- [x] Phase 0: ENDPOINT MAP *(Hoàn thành: Implemented in frontend/src/api/ENDPOINTS.md)*
- [x] Phase 1: FOUNDATION UI *(Hoàn thành)*
  - [x] Cấu hình Theme MUI, Typography, Component overrides
  - [x] AppLayout, PageHeader, ErrorBoundary, Loader
- [x] Phase 2: AUTH & ROUTING *(Hoàn thành)*
  - [x] Store access_token/refresh_token in localStorage
  - [x] Tích hợp API call với Bearer token (RTK Query)
  - [x] Handle unauth redirection
- [ ] Dashboard Admin
  - [ ] Trang tổng quan thống kê
  - [ ] Biểu đồ số lượng học sinh/giáo viên
  - [ ] Biểu đồ doanh thu
  - [ ] Thống kê lớp học đang hoạt động
- [ ] Quản lý tài khoản
  - [x] Danh sách người dùng (Placeholder/Scaffold)
  - [x] Thêm/Sửa/Xóa tài khoản (Placeholder/Scaffold)
  - [x] Phân quyền người dùng (Role-based)
  - [ ] Reset mật khẩu (DEFERRED - NOT IN SCOPE)
- [ ] Compliance Dashboard
  - [x] Danh sách tài liệu chờ phê duyệt *(Hoàn thành: 2026-03-06, Implemented in: frontend/src/pages/compliance/ComplianceQueuePage.tsx, frontend/src/api/materialApi.ts)*
  - [x] Chi tiết kết quả AI audit *(Hoàn thành: 2026-03-06, Implemented in: frontend/src/components/material/MaterialDetailDialog.tsx)*
  - [x] Thao tác phê duyệt/từ chối *(Hoàn thành: 2026-03-06, Implemented in: frontend/src/pages/compliance/ComplianceQueuePage.tsx, frontend/src/api/materialApi.ts)*
  - [ ] Báo cáo tuân thủ

### Task 3.2: Module Thời khóa biểu & Xếp lịch
- [ ] Giao diện hiển thị lịch học
  - [ ] Calendar view (ngày/tuần/tháng)
  - [ ] Hiển thị thông tin lớp học trên calendar
  - [ ] Filter theo giáo viên/phòng/lớp
  - [ ] Color coding theo trạng thái
- [x] Chức năng xếp lịch tự động *(Hoàn thành: 2026-03-06, Implemented in: frontend/src/pages/admin/SchedulingPage.tsx, frontend/src/api/schedulingApi.ts, frontend/src/App.tsx)*
  - [x] Form cấu hình tham số xếp lịch
  - [x] Button trigger xếp lịch
  - [x] Hiển thị progress/loading
  - [x] Hiển thị kết quả xếp lịch
- [ ] Chức năng chỉnh sửa thủ công
  - [ ] Drag & drop để di chuyển lịch
  - [ ] Modal chỉnh sửa chi tiết buổi học
  - [ ] Kiểm tra xung đột real-time
  - [ ] Confirm và lưu thay đổi

### Task 3.3: Module Giáo viên & Học vụ (Phase 3 MVP)
- [x] Phase 4: DEV TEST PANEL *(Hoàn thành: frontend/src/pages/admin/DevToolsPage.tsx)*
- [x] Quản lý Phòng học *(Đã fix getRowId & Normalize response)*
  - [x] Danh sách phòng học
  - [x] Form thêm/sửa phòng học
  - [x] Xem chi tiết phòng học
- [x] Quản lý Lớp học *(Đã fix getRowId & Normalize response)*
  - [x] Danh sách lớp học
  - [x] Form thêm/sửa lớp học
  - [x] Xem chi tiết lớp học
  - [x] Chi tiết lớp với roster học sinh + gán giáo viên *(Hoàn thành: 2026-03-06, Implemented in: frontend/src/pages/admin/ClassesPage.tsx, frontend/src/components/admin/ClassDetailDialog.tsx, frontend/src/api/classApi.ts, cmd/http/controllers/class, internal/usecases/class/get_class_roster.go)*
- [x] Quản lý Học sinh *(Đã fix getRowId & Normalize response)*
  - [x] Danh sách học sinh
  - [x] Form thêm/sửa học sinh
  - [x] Xem chi tiết hồ sơ học sinh
- [x] Quản lý hồ sơ giáo viên *(Đã fix getRowId & Normalize response)*
  - [x] Danh sách giáo viên
  - [x] Form thêm/sửa giáo viên
  - [x] Xem chi tiết hồ sơ
- [ ] Tải và quản lý tài liệu
  - [x] Upload giáo án/tài liệu giảng dạy *(Hoàn thành: 2026-03-06, Implemented in: frontend/src/pages/teacher/TeacherDocumentsPage.tsx, frontend/src/api/materialApi.ts)*
  - [x] Danh sách tài liệu đã tải *(Hoàn thành: 2026-03-06, Implemented in: frontend/src/pages/teacher/TeacherDocumentsPage.tsx, frontend/src/api/materialApi.ts)*
  - [x] Xem trạng thái kiểm duyệt AI *(Hoàn thành: 2026-03-06, Implemented in: frontend/src/pages/teacher/TeacherDocumentsPage.tsx, frontend/src/components/material/MaterialDetailDialog.tsx)*
  - [x] Xem chi tiết phản hồi từ AI *(Hoàn thành: 2026-03-06, Implemented in: frontend/src/components/material/MaterialDetailDialog.tsx, frontend/src/pages/compliance/ComplianceQueuePage.tsx)*
  - [ ] Download/Preview tài liệu
- [ ] Sổ đầu bài điện tử
  - [ ] Danh sách buổi học
  - [ ] Form nhập nội dung bài giảng
  - [ ] Nhập nhận xét học sinh
  - [ ] Đánh giá sau mỗi buổi học
- [ ] Điểm danh
  - [ ] Giao diện điểm danh nhanh
  - [ ] Đánh dấu có mặt/vắng/muộn
  - [ ] Ghi chú lý do vắng
  - [ ] Lịch sử điểm danh

### Task 3.4: Module Học sinh & Phụ huynh
- [ ] Cổng thông tin học sinh
  - [ ] Dashboard học sinh
  - [ ] Xem lịch học cá nhân
  - [ ] Xem kết quả học tập
  - [ ] Xem chuyên cần
  - [ ] Xem thông báo
- [ ] Quản lý học phí
  - [ ] Xem công nợ
  - [ ] Lịch sử thanh toán
  - [ ] Thông báo nhắc nộp học phí
- [ ] Tích hợp Chatbot trợ lý ảo
  - [ ] Setup chatbot UI component
  - [ ] Kết nối với AI backend (Gemini/GPT)
  - [ ] Thiết kế prompt trả lời câu hỏi về khóa học
  - [ ] Thiết kế prompt trả lời câu hỏi về học phí
  - [ ] Thiết kế prompt trả lời câu hỏi về lịch học
  - [ ] Lưu lịch sử chat

---

## ✅ Giai đoạn 4: Kiểm thử, Tối ưu & Hoàn thiện (Tuần 12-15)

### Task 4.1: Kiểm thử chức năng (Functional Testing)
- [ ] Kiểm thử module Quản lý đào tạo
  - [ ] Test CRUD operations
  - [ ] Test validation rules
  - [ ] Test business logic
- [ ] Kiểm thử thuật toán Xếp lịch
  - [ ] Test với tập dữ liệu nhỏ (10 lớp)
  - [ ] Test với tập dữ liệu trung bình (50 lớp)
  - [ ] Test với tập dữ liệu lớn (100+ lớp)
  - [ ] Test các trường hợp xung đột
  - [ ] Test ràng buộc khung giờ cấm
  - [ ] Verify không có trùng lịch
- [ ] Kiểm thử AI Audit
  - [ ] Test OCR với các định dạng file khác nhau
  - [ ] Test độ chính xác nhãn dán với tài liệu mẫu
  - [ ] Test với nội dung an toàn
  - [ ] Test với nội dung có cảnh báo
  - [ ] Test với nội dung không phù hợp
- [ ] Kiểm thử tích hợp Frontend-Backend
  - [ ] Test authentication flow
  - [ ] Test authorization (phân quyền)
  - [ ] Test API calls
  - [ ] Test error handling

### Task 4.2: Kiểm thử hiệu năng (Load Testing)
- [ ] Setup JMeter
  - [ ] Cài đặt và cấu hình JMeter
  - [ ] Tạo test plan
- [ ] Load test API endpoints
  - [ ] Test GET danh sách người dùng (100-200 req/s)
  - [ ] Test GET danh sách lớp học (100-200 req/s)
  - [ ] Test API xếp lịch với concurrent requests
  - [ ] Test API upload tài liệu
- [ ] Phân tích kết quả
  - [ ] Đo thời gian phản hồi trung bình
  - [ ] Đo throughput
  - [ ] Xác định bottleneck
  - [ ] Tối ưu hóa performance
- [ ] Database optimization
  - [ ] Tạo indexes phù hợp
  - [ ] Optimize queries
  - [ ] Connection pooling tuning

### Task 4.3: Triển khai & Bảo vệ
- [ ] Triển khai hệ thống
  - [ ] Hoàn thiện Docker Compose
  - [ ] Cấu hình production environment
  - [ ] Setup reverse proxy (Nginx)
  - [ ] Cấu hình SSL/TLS
  - [ ] Deploy lên server/cloud
- [ ] Tài liệu hóa
  - [ ] Viết User Manual
  - [ ] Viết API Documentation (Swagger)
  - [ ] Viết Developer Guide
  - [ ] Tạo video demo
- [ ] Hoàn thiện báo cáo đồ án
  - [ ] Chương 1: Giới thiệu đề tài
  - [ ] Chương 2: Cơ sở lý thuyết
    - [ ] Thuật toán CSP
    - [ ] AI trong kiểm soát chất lượng
    - [ ] Thông tư 29/2024/TT-BGDĐT
  - [ ] Chương 3: Phân tích & Thiết kế
  - [ ] Chương 4: Triển khai hệ thống
  - [ ] Chương 5: Kiểm thử & Đánh giá
  - [ ] Chương 6: Kết luận & Hướng phát triển
- [ ] Chuẩn bị bảo vệ
  - [ ] Tạo slide thuyết trình
  - [ ] Chuẩn bị demo
  - [ ] Dự đoán câu hỏi và chuẩn bị trả lời

---

## 📊 Theo dõi tiến độ

### Tuần 3-4: Phân tích & Đặc tả
- [x] Hoàn thành Task 1.1
- [x] Hoàn thành Task 1.2
- [ ] Hoàn thành Task 1.3

### Tuần 5-6: Backend Core
- [ ] Hoàn thành Task 2.1
- [ ] Hoàn thành Task 2.2

### Tuần 7-9: Backend Advanced & Frontend Core
- [ ] Hoàn thành Task 2.3
- [ ] Hoàn thành Task 2.4
- [ ] Hoàn thành Task 3.1
- [ ] Hoàn thành Task 3.2

### Tuần 10-11: Frontend Advanced
- [ ] Hoàn thành Task 3.3
- [ ] Hoàn thành Task 3.4

### Tuần 12-13: Testing
- [ ] Hoàn thành Task 4.1
- [ ] Hoàn thành Task 4.2

### Tuần 14-15: Deployment & Documentation
- [ ] Hoàn thành Task 4.3

---

## 🎯 Điểm nhấn của đề tài

> **Lời khuyên từ BA:** Do bạn làm mô hình Single-tenant, hãy tập trung tối ưu vào:
> 1. **Deep UX** - Trải nghiệm người dùng sâu cho từng vai trò
> 2. **Độ chính xác tuyệt đối** - Thuật toán xếp lịch không được có lỗi
> 3. **AI Integration** - Tận dụng AI để tạo giá trị khác biệt
> 4. **Compliance** - Tuân thủ nghiêm ngặt Thông tư 29

---

## 📝 Ghi chú

- Đánh dấu `[x]` cho task đã hoàn thành
- Đánh dấu `[/]` cho task đang thực hiện
- Để trống `[ ]` cho task chưa bắt đầu
- Cập nhật file này thường xuyên để theo dõi tiến độ

*Auth Standardization Implemented: `frontend/src/contexts/AuthContext.tsx`, `frontend/src/App.tsx`, `frontend/src/pages/LoginPage.tsx`. (OTP/Email flows DEFERRED/NOT IN SCOPE for demo).*

**Ngày tạo:** 2026-02-06
**Người thực hiện:** [Tên của bạn]
**Giảng viên hướng dẫn:** [Tên GVHD]

---

## 🚨 Current Fix Queue (User Review 2026-03-31)

Nguyen tac thuc hien cho dot fix nay:
- Lam tuan tu theo thu tu duoi day, khong mo song song nhieu module.
- Moi task chi bat dau sau khi task truoc da duoc review/chap nhan.
- Uu tien sua luong dang loi that, bo placeholder va thong diep "dang phat trien" truoc khi lam them tinh nang moi.
- Cac task UI text/can le chi dong task khi da kiem tra lai tren man hinh that, khong dong task chi vi da sua CSS.

| Thu tu | Trang thai | Module | Van de hien tai | Muc tieu dong task |
|---|---|---|---|---|
| 1 | [x] | Auth + Teacher | Man danh sach giao vien chua co chuc nang them giao vien tren UI; kha nang cao do check role dang dung chu thuong trong khi BE/store dang luu role hoa (`ADMIN`). | Nut `Them giao vien` hien dung voi tai khoan admin; route tao/sua giao vien di duoc tu danh sach; kiem tra lai role gating de khong lap lai loi an nut o cac man admin khac. |
| 2 | [x] | Class - list filter + session stability | Man quan ly lop hoc bi vo toolbar sau khi chon filter trang thai; sau khi reload thi bi logout; text ghi chu "Chua co ghi chu" dang lech baseline. | Filter search + trang thai hoat dong on dinh sau moi lan chon/xoa filter; reload trang khong lam mat session admin; subtitle/secondary text trong list can hang dung voi ten lop. |
| 3 | [x] | Class - detail roster + add student UX | `Chi tiet lop hoc` khong hien du thong tin hoc sinh trong roster; dialog them hoc sinh chua search duoc hoc sinh de chon; can tich hop search ngay trong danh sach kha dung. | Tab `Danh sach hoc sinh` hien day du ten, ma/khối lop, SDT, SDT phu huynh, trang thai; dialog them hoc sinh cho phep tim theo ten/ma/khối lop ngay tren danh sach chon va them thanh cong vao lop. |
| 4 | [x] | Room | Text phu trong man phong hoc dang bi thap/le baseline, vi du `Tang 1 Toa A`. | Text ten phong + dong thong tin phu can hang thang tren danh sach va chi tiet lien quan; khong con cam giac "dut dong" khi scroll danh sach. |
| 5 | [x] | Program/Course - data contract + layout | Man chuong trinh/khoa hoc dang lech text detail; dialog chi tiet dang hien thong diep ky thuat `Contract backend... Chua co field...`; hien trang nay khong phu hop voi user flow. | Layout ten + detail text can hang dung; bo thong diep contract ky thuat khoi UI nguoi dung; mapping data detail dua tren field that backend dang tra ve, neu thieu field thi degrade gracefully. |
| 6 | [x] | Program/Course - action flow | Lien ket khoa hoc chua dung duoc; nut chinh sua record chua dung; flow quan ly chuong trinh/khoa hoc chua hoan tat du du API da co. | Co the mo form chinh sua chuong trinh/khoa hoc, luu thanh cong; co the lien ket/bo lien ket khoa hoc vao chuong trinh tu UI; refresh lai danh sach/chi tiet thay doi dung. |
| 7 | [x] | Scheduling (CSP) | Audit ngay 2026-04-01 cho thay module xep lich dang o muc scaffold/demo: FE co man preview nhung contract date bi lech, BE commit chua persist lesson, preview store dang luu in-memory va solver chua bam sat du lieu hoc vu that. | Da hoan tat phase audit + implementation uu tien: sua contract input/output, nang cap preview bam course/session/class_schedule, hien conflict/actionable message, va mo commit that xuong lesson co check trung lich. |
| 8 | [/] | AI Audit | Module kiem duyet tai lieu hien tai moi o muc scaffold/stub, chua dung duoc cho luong giao vien/compliance. | Da audit luong hien tai: tai lieu duoc upload tu FE giao vien, file luu local, metadata luu DB, nhan do stub OCR + stub Gemini sinh ra; can tach task con de bien scaffold nay thanh luong co the demo duoc on dinh. |

### Chi tiet task theo thu tu implement

1. Auth + Teacher CTA
- Ra soat role gating o man `TeachersPage` va cac man admin lien quan.
- Doi chieu role dang luu trong store/localStorage voi dieu kien hien nut thao tac.
- Xac nhan route tao moi/chinh sua giao vien da mo duong di tu danh sach.

2. Class list stability
- Sua logic state/filter de toolbar khong bi bien mat sau khi chon trang thai.
- Kiem tra nguyen nhan reload bi logout va dua ve mot hanh vi session on dinh.
- Chinh lai typography/spacing cua dong note duoi ten lop.

3. Class detail roster + add student
- Ra soat mapping du lieu roster tu API sang bang hien thi.
- Hien day du cot thong tin hoc sinh thay vi chi co avatar.
- Dua search vao danh sach hoc sinh kha dung, uu tien tim ngay trong control chon hoc sinh.

4. Room alignment polish
- Chinh typography stack o ten phong/dia chi.
- Kiem tra lai row height, icon alignment va baseline tren cac kich thuoc man hinh.

5. Program/Course contract cleanup
- Loai bo thong diep ky thuat dang lo ra o dialog chi tiet.
- Dung cac field backend dang tra ve de hien thi thong tin co nghia voi user.
- Chinh can le text detail trong list va detail dialog.

6. Program/Course action flow completion
- Mo lai flow them/sua chuong trinh.
- Mo flow them/sua khoa hoc neu dang bi placeholder.
- Noi lai thao tac lien ket/bo lien ket khoa hoc voi chuong trinh bang cac mutation da co.

7. Scheduling audit -> implementation
- Audit completed 2026-04-01. Ket qua chinh:
- FE dang gui `date_from`/`date_to` dang `YYYY-MM-DD` tu form date ([frontend/src/pages/admin/SchedulingPage.tsx]), trong khi BE bind truc tiep vao `time.Time` ([cmd/http/controllers/scheduling/dto.go]). Day la blocker contract rat de lam preview fail ngay tu request body.
- FE cho phep bam `Commit scaffold` sau moi preview va hien toast thanh cong ([frontend/src/pages/admin/SchedulingPage.tsx]), nhung BE commit chi tra message `Persisting lessons is TODO` va khong tao `Lesson`/`ClassSchedule` nao ([internal/usecases/scheduling/commit_preview.go]).
- Preview result dang luu bang in-memory store ([internal/services/scheduling/preview_store.go]), nen restart server la mat het run; `preview/latest` khong co gia tri van hanh thuc te.
- Solver hien tai moi tao 1 bien cho moi lop va dem `scheduled_lessons = len(assignments)` ([internal/usecases/scheduling/preview.go]), trong khi du lieu hoc vu co `course.session_count`, `course.session_duration_minutes`, `class_schedule`, `lesson`. Nghia la preview chua xep "cac buoi hoc that", moi xep 1 slot demo cho tung lop.
- `TeacherLabel` dang lay truc tiep tu `teacher_id` thay vi ten giao vien ([internal/usecases/scheduling/preview.go]), nen ket qua preview kho doc voi admin.
- Domain slot dang fix cung 6 khung gio/ngay va duration 120 phut ([internal/usecases/scheduling/preview.go]), chua bam `course`, chua bam `class_schedule`, chua co logic fallback/no-data message cho cac truong hop thieu `teacher`, `room`, `course`, `session_count`.
- FE da co danh sach conflict, nhung message hien tai van o muc tong quat; chua huong user phai sua du lieu nao truoc de chay lai preview cho dung.
- Thu tu implement de xong task 7:
- 7.1. [x] Sua contract request/response giua FE-BE cho preview:
  DTO BE da doi sang nhan string va parse duoc `YYYY-MM-DD`/RFC3339; FE tiep tuc gui plain date tu date picker va co validate `date_to >= date_from`.
- 7.2. [x] Khoa hanh vi commit gia:
  FE da bo CTA `commit scaffold` khoi luong thao tac chinh, doi thanh trang thai disabled + thong bao ro day moi la preview, chua persist `lesson`.
- 7.3. [x] Nang cap conflict messaging:
  BE da tra conflict cu the hon theo nhom nguyen nhan (`NO_CLASS_INPUT`, `NO_ACTIVE_ROOM`, `PREFERRED_ROOM_UNAVAILABLE`, `ROOM_CAPACITY_BLOCK`, `NO_SLOT_IN_RANGE`, ...); FE da hien severity + goi y xu ly de admin biet can sua du lieu nao truoc.
- 7.4. [x] Dua preview bam du lieu hoc vu that:
  Da doi tu `1 lop = 1 bien demo` sang sinh nhieu buoi theo `course.session_count`, duration theo `course.session_duration_minutes`, preload `class_schedule`, va chi sinh slot theo `day_of_week/start_time/end_time` thuc te cua lop. Neu `class_schedule` co `room_id` thi preview cung ton trong rang buoc phong nay ngay tu domain.
- 7.5. [x] Nang cap chat luong ket qua preview:
  Da hien ten giao vien that neu preload co du lieu, hien `Buoi X/Y` tren assignment/conflict, severity cua conflict, action hint, tong `requested_sessions`, progress theo tung lop, va bo sung conflict ro hon cho cac tinh huong `CLASS_SCHEDULE_NO_SLOT`, `CLASS_SCHEDULE_ROOM_UNAVAILABLE`, `ROOM_CAPACITY_BLOCK`, `NO_DOMAIN`...
- 7.6. [x] Commit that sau khi preview on dinh:
  BE da persist assignment xuong bang `lessons` bang transaction, co check trung lich theo lop/giao vien/phong voi lesson da ton tai, va FE da mo lai CTA commit khi preview dat `COMPLETED`. Preview store van luu in-memory cho muc dich xem lai preview, nhung khong con chan luong commit xuong `lesson`.

8. AI Audit audit -> implementation
- Audit completed 2026-04-01. Ket qua chinh:
- Nguon tai lieu hien tai den tu giao dien giao vien qua `POST /api/v1/materials/upload` ([frontend/src/pages/teacher/TeacherDocumentsPage.tsx], [cmd/http/controllers/material/v1.go]). He thong khong tu lay tai lieu tu kho ngoai, Google Drive hay DB khac.
- File vat ly dang duoc luu local tren server o thu muc `storage/materials` qua `localStorageService.Save(...)` ([internal/services/audit/services.go]). Duong dan luu vao cot `materials.file_path`.
- Metadata tai lieu luu o bang `materials`; ket qua OCR/AI luu o bang `audit_logs`; nhan hien tai cua tai lieu duoc tro den bang `materials.latest_label_id`; quyet dinh phe duyet luu o bang `approval_decisions` ([internal/infrastructure/database/postgres/sql_migration/23_create_material_audit_tables.up.sql]).
- Co che gan nhan hien tai la stub:
  - OCR stub lay preview noi dung file bang cach doc byte va cat toi da 500 ky tu.
  - Gemini stub khong goi model that; no chi so tu khoa trong text (`violence`, `gambling`, `danger`, `exam`, `cheat`, `warning`) de gan `SAFE/WARNING/DANGER` ([internal/services/audit/services.go]).
- Luong upload hien tai la dong bo trong 1 request:
  - luu file
  - tao `material`
  - OCR stub
  - Gemini stub
  - tao `audit_log`
  - cap nhat `material.status = AI_REVIEWED`
  Nghia la chua co queue/background job that.
- Frontend dang cho giao vien chon `teacher_id` thu cong tu combobox ([frontend/src/pages/teacher/TeacherDocumentsPage.tsx]). Dieu nay phu hop demo nhung chua dung voi luong thuc te, vi tai khoan giao vien dang dang nhap phai tu suy ra `teacher_id`.
- Man chi tiet hien duong dan file (`file_path`) va raw OCR text, nhung chua co endpoint tai/xem file an toan. Nghia la nguoi dung thay duong dan local, nhung chua mo/xem tep ngay tren UI duoc.
- Hien chua co validate manh cho loai file, kich thuoc file, scan loi, retry, hay timeout AI.
- Hien chua co "manual relabel/re-audit", chua co thong ke chat luong tai lieu, chua co tich hop model AI that, va chua co rule/prompt theo Thong tu 29 hay rubric compliance that.
- Thu tu implement de xong task 8:
- 8.1. [ ] Chuan hoa nguon tai lieu va ownership:
  upload se lay `teacher_id` tu user/role dang nhap hoac map ro tu tai khoan, bo chon tay tren man giao vien; validate file type/size ngay tu FE + BE.
- 8.2. [/] Chuan hoa luu tru va truy cap file:
  Da bo hien raw local path tren UI, them endpoint download an toan theo `material_id`, validate upload, va doi quy uoc luu file local theo cau truc thu muc. Phan cleanup/retention van de lai cho buoc sau.
- Ke hoach de xuat theo huong "production don gian cho do an":
  - 8.2.a. [x] Doi mo hinh luu tru tu `file_path` raw sang `storage_key/relative_path`:
    file van luu local disk, nhung theo cau truc co to chuc nhu `storage/materials/YYYY/MM/<material_id>/<original_file>` thay vi 1 thu muc phang.
  - 8.2.b. [x] Bo sung metadata file toi thieu trong bang `materials`:
    luu them `mime_type`, `file_size`, va neu can `original_file_name`; muc tieu la du thong tin de hien thi, validate va debug ma khong can lo duong dan that.
  - 8.2.c. [x] Tao abstraction download/view file:
    them endpoint backend kieu `GET /api/v1/materials/:id/download` de doc file theo `material_id`, kiem tra quyen, roi moi stream file; frontend chi goi endpoint nay, khong doc local path truc tiep.
  - 8.2.d. [x] Don dep UI hien thi file:
    bo hien `file_path` raw khoi dialog chi tiet; doi thanh thong tin an toan hon nhu ten file, loai file, kich thuoc, va nut `Tai tep`/`Xem tep`.
  - 8.2.e. [x] Validate upload theo huong an toan toi thieu:
    chan file qua lon, gioi han mime type/phan mo rong (pdf/docx/png/jpg tuy pham vi demo), va tra thong bao loi ro rang neu file khong hop le.
  - 8.2.f. [ ] Chuan bi cleanup/backup don gian:
    quy uoc ro file local nao thuoc `material_id` nao, de sau nay co the xoa mem metadata ma van co cach doi chieu hoac viet job cleanup neu can.
- Ly do chon huong nay:
  - khong can them S3/MinIO nen setup nhanh cho do an;
  - van giong production hon cach luu phang + lo raw path hien tai;
  - de demo duoc luong upload, xem file, audit, review ma khong tang nhieu ha tang.
- Thu tu implement de tranh sua lan man:
  - buoc 1: schema + entity metadata file
  - buoc 2: storage service doi sang path co cau truc
  - buoc 3: endpoint download/view file
  - buoc 4: frontend bo raw path, them nut tai/xem
  - buoc 5: validate file upload FE + BE
- 8.3. [ ] Lam ro pipeline audit:
  tach status `UPLOADED -> SCANNING -> AI_REVIEWED -> APPROVED/REJECTED`, xu ly loi OCR/AI ro rang hon, va hien retry/manual re-audit neu can.
- 8.4. [ ] Thay stub label bang logic AI that:
  tich hop OCR/model that, dung prompt/schema/rule engine phu hop nghiep vu compliance.
- 8.5. [ ] Hoan thien compliance review:
  review form dung voi user compliance dang dang nhap, luu lich su day du, va cap nhat queue/status nhat quan.
- 8.6. [ ] Bo sung bao cao/thong ke:
  tong hop so tai lieu theo nhan, severity, trang thai phe duyet, giao vien, khoang thoi gian.

---

## 📚 Historical Top 10 Demo (giu lai de doi chieu)

| Ưu tiên | Trạng thái | Module | Tên Task | Phân loại |
|---|---|---|---|---|
| 1 | [x] | Program/Course | CRUD Program/Course + Link Course->Program (Implemented in: internal/entities, repositories, usecases/program, cmd/http/controllers/program) | Backend |
| 2 | [x] | Program/Course | Giao diện quản lý Program/Course (Table + Dialog) (Implemented in: frontend/src/api/programApi.ts, courseApi.ts, pages/admin/ProgramPage.tsx, pages/admin/CoursePage.tsx) | Frontend |
| 3 | [x] | Class | API thêm/xóa học sinh, gán giáo viên, check sức chứa (Implemented in internal/usecases/class, cmd/http/controllers/class) | Backend |
| 4 | [x] | Class | Giao diện quản lý roster (danh sách học sinh) + gán GV (Implemented in: frontend/src/pages/admin/ClassesPage.tsx, frontend/src/components/admin/ClassDetailDialog.tsx, frontend/src/api/classApi.ts, frontend/src/components/admin/ClassDialog.tsx, cmd/http/controllers/class, internal/usecases/class/get_class_roster.go, internal/infrastructure/database/postgres/implement/enrollment_repository.go) | Frontend |
| 5 | [x] | Student | Cố định màn hình Student (lỗi hiện thị danh sách, trắng trang) | Frontend |
| 6 | [x] | Room | Xử lý trạng thái deleted_at đồng bộ BE/FE | Frontend |
| 7 | [x] | Scheduling (CSP) | Scaffold cấu trúc dữ liệu + Hard constraints (base) (Implemented in: internal/usecases/scheduling, internal/services/scheduling, cmd/http/controllers/scheduling, cmd/http/main.go) | Backend |
| 8 | [x] | Scheduling (CSP) | Giao diện Scheduling Trigger & Preview (Implemented in: frontend/src/pages/admin/SchedulingPage.tsx, frontend/src/api/schedulingApi.ts, frontend/src/App.tsx, frontend/src/api/baseApi.ts) | Frontend |
| 9 | [x] | AI Audit | Scaffold luồng Upload tài liệu + Audit log + Phê duyệt (Implemented in: internal/entities/material.go, internal/entities/label.go, internal/entities/audit_log.go, internal/entities/approval_decision.go, internal/usecases/material, internal/services/audit, cmd/http/controllers/material, internal/infrastructure/database/postgres/sql_migration/23_create_material_audit_tables.up.sql) | Backend |
| 10 | [x] | AI Audit | Giao diện màn hình upload cho Giáo viên & Compliance queue (Implemented in: frontend/src/pages/teacher/TeacherDocumentsPage.tsx, frontend/src/pages/compliance/ComplianceQueuePage.tsx, frontend/src/components/material/MaterialDetailDialog.tsx, frontend/src/api/materialApi.ts, frontend/src/App.tsx) | Frontend |

### Ghi chú triển khai gần nhất

- 2026-03-06 — Top 10 #4 Class roster UI + gán giáo viên
  - Trạng thái: [x]
  - Implemented in: `frontend/src/pages/admin/ClassesPage.tsx`, `frontend/src/components/admin/ClassDetailDialog.tsx`, `frontend/src/components/admin/ClassDialog.tsx`, `frontend/src/api/classApi.ts`, `frontend/src/api/teacherApi.ts`, `cmd/http/controllers/class/controller.go`, `cmd/http/controllers/class/v1.go`, `internal/usecases/class/get_class_roster.go`, `internal/repositories/interface/enrollment.go`, `internal/infrastructure/database/postgres/implement/enrollment_repository.go`
  - Endpoints: `GET /api/v1/classes/:id/students`, `POST /api/v1/classes/:id/students`, `DELETE /api/v1/classes/:id/students`, `PUT /api/v1/classes/:id/teacher`
  - Cách test nhanh:
    1. Vào `Quản lý lớp học` -> dùng toolbar tìm/lọc -> bấm một dòng bất kỳ để mở `Chi tiết lớp học`.
    2. Tab `Danh sách học sinh` -> `Thêm học sinh` -> chọn nhiều học sinh -> lưu -> kiểm tra DataGrid roster cập nhật ngay.
    3. Tab `Giáo viên phụ trách` -> chọn giáo viên -> `Lưu giáo viên phụ trách` -> đóng dialog -> kiểm tra cột giáo viên ở danh sách lớp.
    4. API nhanh:
       - `curl -X GET http://localhost:9000/api/v1/classes/<class_id>/students`
       - `curl -X POST http://localhost:9000/api/v1/classes/<class_id>/students -H 'Content-Type: application/json' -d '{"student_ids":["<student_id>"]}'`
       - `curl -X PUT http://localhost:9000/api/v1/classes/<class_id>/teacher -H 'Content-Type: application/json' -d '{"teacher_id":"<teacher_id>"}'`

- 2026-03-06 — Top 10 #7 Scheduling backend scaffold
  - Trạng thái: [x]
  - Implemented in: `internal/usecases/scheduling/types.go`, `internal/usecases/scheduling/preview.go`, `internal/usecases/scheduling/get_preview.go`, `internal/usecases/scheduling/commit_preview.go`, `internal/services/scheduling/preview_store.go`, `internal/services/provider.go`, `cmd/http/controllers/scheduling/controller.go`, `cmd/http/controllers/scheduling/v1.go`, `cmd/http/controllers/scheduling/dto.go`, `cmd/http/controllers/provider.go`, `cmd/http/main.go`
  - Endpoints: `POST /api/v1/scheduling/preview`, `GET /api/v1/scheduling/preview/latest`, `GET /api/v1/scheduling/preview/:id`, `POST /api/v1/scheduling/commit`
  - Migration: Không cần migration ở bước scaffold này. Preview/result đang lưu in-memory để phục vụ demo.
  - Cách test nhanh:
    1. `curl -X POST http://localhost:9000/api/v1/scheduling/preview -H 'Content-Type: application/json' -d '{"date_from":"2026-03-06T00:00:00Z","date_to":"2026-03-12T00:00:00Z"}'`
    2. Copy `run_id` từ response rồi gọi `curl http://localhost:9000/api/v1/scheduling/preview/<run_id>`.
    3. Gọi `curl http://localhost:9000/api/v1/scheduling/preview/latest` để kiểm tra latest preview.

- 2026-03-06 — Top 10 #8 Scheduling FE trigger + preview
  - Trạng thái: [x]
  - Implemented in: `frontend/src/pages/admin/SchedulingPage.tsx`, `frontend/src/api/schedulingApi.ts`, `frontend/src/App.tsx`, `frontend/src/api/baseApi.ts`
  - Cách test nhanh:
    1. Vào `Xếp lịch (CSP)` trên sidebar admin.
    2. Chọn khoảng ngày, lọc lớp/GV/phòng nếu cần rồi bấm `Chạy xếp lịch`.
    3. Kiểm tra `Preview run`, bảng assignment, banner conflict và nút `Commit scaffold`.

- 2026-03-06 — Top 10 #9 AI Audit backend scaffold
  - Trạng thái: [x]
  - Implemented in: `internal/entities/material.go`, `internal/entities/label.go`, `internal/entities/audit_log.go`, `internal/entities/approval_decision.go`, `internal/repositories/interface/material.go`, `internal/repositories/interface/label.go`, `internal/repositories/interface/audit_log.go`, `internal/repositories/interface/approval_decision.go`, `internal/infrastructure/database/postgres/implement/material_repository.go`, `internal/infrastructure/database/postgres/implement/label_repository.go`, `internal/infrastructure/database/postgres/implement/audit_log_repository.go`, `internal/infrastructure/database/postgres/implement/approval_decision_repository.go`, `internal/services/audit/services.go`, `internal/usecases/material`, `cmd/http/controllers/material`, `internal/infrastructure/database/postgres/sql_migration/23_create_material_audit_tables.up.sql`
  - Endpoints: `POST /api/v1/materials/upload`, `GET /api/v1/materials`, `GET /api/v1/materials/flagged`, `GET /api/v1/materials/:id`, `POST /api/v1/materials/:id/review`
  - Migration: `internal/infrastructure/database/postgres/sql_migration/23_create_material_audit_tables.up.sql`
  - Nhắc chạy migrate: `make migrate`
  - Cách test nhanh:
    1. Chạy migrate theo quy trình repo: `make migrate` (tôi chưa chạy).
    2. Upload bằng multipart: `curl -X POST http://localhost:9000/api/v1/materials/upload -F teacher_id=<teacher_id> -F title='Giáo án Demo' -F description='demo' -F file=@/path/to/file.txt`
    3. Kiểm tra queue: `curl http://localhost:9000/api/v1/materials?queue=flagged&status=AI_REVIEWED`
    4. Duyệt tài liệu: `curl -X POST http://localhost:9000/api/v1/materials/<material_id>/review -H 'Content-Type: application/json' -d '{"compliance_officer_id":"<user_id>","approved":true,"notes":"OK cho demo"}'`

- 2026-03-06 — Top 10 #10 AI Audit frontend
  - Trạng thái: [x]
  - Implemented in: `frontend/src/api/materialApi.ts`, `frontend/src/components/material/MaterialDetailDialog.tsx`, `frontend/src/pages/teacher/TeacherDocumentsPage.tsx`, `frontend/src/pages/compliance/ComplianceQueuePage.tsx`, `frontend/src/App.tsx`
  - Cách test nhanh:
    1. Vào `Tài liệu giảng dạy` ở menu giáo viên, chọn giáo viên demo, nhập tiêu đề và upload một file.
    2. Mở dialog chi tiết để xem label AI, reasoning, raw OCR text và lịch sử audit.
    3. Vào `Tài liệu cần duyệt` của Compliance/Admin, mở chi tiết, nhập `compliance_officer_id`, chọn Approve/Reject rồi lưu quyết định.
