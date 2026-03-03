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
  - [ ] API thêm/xóa học sinh khỏi lớp
  - [ ] API gán giáo viên cho lớp
  - [ ] API kiểm tra sĩ số
- [ ] API quản lý Chương trình đào tạo *(Missing: usecases, repositories, controllers)*
  - [ ] CRUD Program
  - [ ] CRUD Course
  - [ ] API liên kết Course với Program
- [x] API quản lý Phòng học *(Hoàn thành: 2026-02-25, Implemented in: cmd/http/controllers/room)*
  - [x] CRUD Room
  - [x] API kiểm tra sức chứa
  - [x] API kiểm tra tình trạng phòng học


### Task 2.3: Triển khai thuật toán Xếp lịch tự động (CSP)
- [ ] Thiết kế cấu trúc dữ liệu CSP
  - [ ] Định nghĩa Variable (Biến)
  - [ ] Định nghĩa Domain (Miền giá trị)
  - [ ] Định nghĩa Constraint (Ràng buộc)
- [ ] Triển khai giải thuật Backtracking
  - [ ] Implement thuật toán backtracking cơ bản
  - [ ] Implement heuristic MRV (Minimum Remaining Values)
  - [ ] Implement heuristic Degree Heuristic
  - [ ] Implement heuristic LCV (Least Constraining Value)
- [ ] Tích hợp Forward Checking
  - [ ] Implement logic loại bỏ sớm các giá trị xung đột
  - [ ] Optimize performance với pruning
- [ ] Implement Hard Constraints
  - [ ] Kiểm tra trùng lịch giáo viên
  - [ ] Kiểm tra trùng lịch phòng học
  - [ ] Kiểm tra khung giờ cấm (sau 22h)
  - [ ] Kiểm tra sức chứa phòng học
- [ ] Implement Soft Constraints
  - [ ] Tính điểm ưu tiên lịch liên tiếp
  - [ ] Tối ưu hóa khoảng cách giữa các buổi học
- [ ] API Scheduling
  - [ ] API trigger xếp lịch tự động
  - [ ] API lấy kết quả xếp lịch
  - [ ] API chỉnh sửa lịch thủ công
  - [ ] API kiểm tra xung đột khi chỉnh sửa

### Task 2.4: Tích hợp AI Pipeline (Kiểm soát chất lượng)
- [ ] Xây dựng module OCR
  - [ ] Tích hợp thư viện OCR (Tesseract/Cloud Vision API)
  - [ ] API upload tài liệu
  - [ ] API trích xuất văn bản từ PDF/Image
  - [ ] Xử lý và làm sạch văn bản
- [ ] Kết nối Google Gemini API
  - [ ] Setup API credentials
  - [ ] Thiết kế prompt phát hiện nội dung không phù hợp
  - [ ] Thiết kế prompt phát hiện sai lệch kiến thức
  - [ ] API gửi văn bản đến Gemini và nhận kết quả
- [ ] Xây dựng logic gán nhãn tự động
  - [ ] Định nghĩa các loại nhãn (An toàn/Cảnh báo/Nguy hiểm)
  - [ ] Logic phân loại dựa trên kết quả AI
  - [ ] Lưu kết quả vào database
- [ ] API AI Audit
  - [ ] API quét tài liệu
  - [ ] API lấy lịch sử audit
  - [ ] API phê duyệt/từ chối tài liệu (Compliance Officer)
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
  - [ ] Danh sách tài liệu chờ phê duyệt
  - [ ] Chi tiết kết quả AI audit
  - [ ] Thao tác phê duyệt/từ chối
  - [ ] Báo cáo tuân thủ

### Task 3.2: Module Thời khóa biểu & Xếp lịch
- [ ] Giao diện hiển thị lịch học
  - [ ] Calendar view (ngày/tuần/tháng)
  - [ ] Hiển thị thông tin lớp học trên calendar
  - [ ] Filter theo giáo viên/phòng/lớp
  - [ ] Color coding theo trạng thái
- [ ] Chức năng xếp lịch tự động
  - [ ] Form cấu hình tham số xếp lịch
  - [ ] Button trigger xếp lịch
  - [ ] Hiển thị progress/loading
  - [ ] Hiển thị kết quả xếp lịch
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
- [x] Quản lý Học sinh *(Đã fix getRowId & Normalize response)*
  - [x] Danh sách học sinh
  - [x] Form thêm/sửa học sinh
  - [x] Xem chi tiết hồ sơ học sinh
- [x] Quản lý hồ sơ giáo viên *(Đã fix getRowId & Normalize response)*
  - [x] Danh sách giáo viên
  - [x] Form thêm/sửa giáo viên
  - [x] Xem chi tiết hồ sơ
  - [ ] Quản lý chứng chỉ/bằng cấp
- [ ] Tải và quản lý tài liệu
  - [ ] Upload giáo án/tài liệu giảng dạy
  - [ ] Danh sách tài liệu đã tải
  - [ ] Xem trạng thái kiểm duyệt AI
  - [ ] Xem chi tiết phản hồi từ AI
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

## 🚀 Next Tasks (Top 10 Ưu tiên Demo)

| Ưu tiên | Trạng thái | Module | Tên Task | Phân loại |
|---|---|---|---|---|
| 1 | [/] | Program/Course | CRUD Program/Course + Link Course->Program | Backend |
| 2 | [ ] | Program/Course | Giao diện quản lý Program/Course (Table + Dialog) | Frontend |
| 3 | [ ] | Class | API thêm/xóa học sinh, gán giáo viên, check sức chứa | Backend |
| 4 | [ ] | Class | Giao diện quản lý roster (danh sách học sinh) + gán GV | Frontend |
| 5 | [x] | Student | Cố định màn hình Student (lỗi hiện thị danh sách, trắng trang) | Frontend |
| 6 | [x] | Room | Xử lý trạng thái deleted_at đồng bộ BE/FE | Frontend |
| 7 | [ ] | Scheduling (CSP) | Scaffold cấu trúc dữ liệu + Hard constraints (base) | Backend |
| 8 | [ ] | Scheduling (CSP) | Giao diện Scheduling Trigger & Preview | Frontend |
| 9 | [ ] | AI Audit | Scaffold luồng Upload tài liệu + Audit log + Phê duyệt | Backend |
| 10 | [ ] | AI Audit | Giao diện màn hình upload cho Giáo viên & Compliance queue | Frontend |
