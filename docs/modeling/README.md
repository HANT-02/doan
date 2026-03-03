# System Modeling - Hệ thống Quản lý Trung tâm Dạy thêm tích hợp AI

Thư mục này chứa bộ tài liệu Mô hình hóa Toàn vẹn (System Modeling) của hệ thống Quản lý Trung tâm. Các biểu đồ được thiết kế với [PlantUML](https://plantuml.com/), giúp đội ngũ Development dễ dàng tham chiếu logic Core Domain, luồng UI và lõi Thuật toán Xếp lịch của Dự án (*Bám theo Thông tư 29/BGDĐT*).

---

## 📂 Kiến trúc Thư mục Hệ thống

Khu vực này được chia làm 4 thành phần thiết kế chính (A/B/C/D):

### A. Use Case Diagrams (`/usecase`)
*(Nhóm sơ đồ phân rã tính năng theo Roles (vai trò) và Flow (chuỗi hoạt động))*
- `01_overview.puml`: Bức tranh toàn cảnh về System Actions.
- `02_admin.puml`: Vai trò Admin - Quản trị hệ thống, dữ liệu nguồn, xếp lịch.
- `03_teacher.puml`: Vai trò Giáo viên - Nộp giáo án (AI Audit), điểm danh, sổ đầu bài.
- `04_student_parent.puml`: Vai trò Học sinh/Phụ huynh (Portal & Chatbot AI scaffolding).
- `05_compliance.puml`: Vai trò Kiểm soát viên (Compliance).
- `06_ai_agent.puml`: Các Task chạy ngầm của AI (System Actor).
- `07_scheduling_flow.puml`: Workflow chạy Xếp lịch tự động.
- `08_material_audit_flow.puml`: Workflow chạy OCR & LLM xác minh tài liệu.

### B. Class Diagrams (`/class-diagrams`)
*(Nhóm sơ đồ thiết kế Thực thể (Domain Entities))*
- `01_domain_core.puml`: Sơ đồ các Table/Class chuẩn (User, Role, Class, Program, Course, Room, Student, Teacher). Tích hợp constraint từ Thông tư 29.
- `02_ai_audit.puml`: Sơ đồ các Entity mô phỏng quá trình kiểm định Tài liệu (Material, AuditLog, AIAnalysisResult, ApprovalDecision).
- `03_scheduling.puml`: Cấu trúc code thiết kế (Strategy Pattern) cho Background Solver giải quyết bài toán Xếp lịch thông qua CSP (Variables, Domains, Constraints, Heuristics).

### C. Sequence Diagrams (`/sequence-diagrams`)
*(Nhóm biểu đồ luồng Hệ thống từ Frontend → DB (Clean Architecture))*
- `01_login.puml`: Luồng Đăng nhập (Auth & JWT).
- `02_teacher_crud_admin.puml`: Luồng Admin tạo mới Giáo viên (Validate -> Insert DB).
- `03_class_enroll_assign_teacher.puml`: Chức năng Xếp lớp Học sinh & Gán Giáo viên (Check Capacity -> Relationships Update).
- `04_auto_scheduling_csp.puml`: Luồng hệ thống kích hoạt Scheduling Solver -> Trả về kết quả dự thảo (Preview) & Lưu.
- `05_material_upload_audit_approval.puml`: Luồng Teacher Upload -> Background Queue gọi OCR/Gemini API -> Compliance xem/duyệt -> Alert Alert.
- `06_student_leave_request.puml`: Luồng Học sinh xin phép nghỉ phép trực tuyến & Teacher xem sự thay đổi (Scaffold).

### D. Scheduling Modeling (CSP) (`/scheduling`)
*(Đặc tả riêng về Lõi trí tuệ hệ thống: Bài toán thuật toán Xếp Lịch)*
- `01_csp_model.md`: Định nghĩa Toán học của Variables $V$, Domains $D$ và Các bộ phận Constraints $C$ cứng/mềm.
- `02_backtracking_flow.puml`: Sơ đồ Nhánh Hoạt động (Activity Diagram) của vòng lặp Đệ quy Tìm kiếm Quay lui (DFS/Forward Checking).
- `03_constraints.puml`: Khối Component đánh giá ràng buộc (Checkers).
- `04_pseudocode.md`: Mã giả chi tiết cho logic Backtrack, Forward Check & Hàm chấm điểm Heuristic (MRV/LCV).

---

## 🛠 Hướng dẫn Render PlantUML 

Để trực quan hóa các bản vẽ, sử dụng một trong các cách sau:

### 1. Dành cho Developer / BA đang dùng IDE: VS Code
- Tải Plugin có tên: **[PlantUML](https://marketplace.visualstudio.com/items?itemName=jebbs.plantuml)** (của jebbs).
- Mở bất kỳ file có đuôi mở rộng `.puml` trong thư mục này.
- Bấm **Alt + D** (Windows/Linux) hoặc **Option + D** (MacOS) để hệ thống tự render ảnh và chạy Preview split-view ngay trong biên tập.

### 2. Dành cho người Review Nhanh (Online Viewer)
Bạn không cần tải phần mềm.
- Copy mã Nguồn nội dung (Raw Source) của `.puml`.
- Truy cập Text Compiler như [PlantText](https://www.planttext.com/) hoặc [PlantUML Server](http://www.plantuml.com/plantuml/uml/).
- Dán Code -> Bấm **Refresh**.

### 3. Xuất Hàng Loạt (Export Image Files Bulk) bằng Terminal
Có thể xuất ảnh đồng loạt dùng chuẩn bị báo cáo Luận văn Word / Latex.
- Yêu cầu môi trường có cài sẵn **Java** và gói hỗ trợ **Graphviz** (Tùy chọn cho Windows/Mac).
- Tải xuống tệp thực thi Standalone `plantuml.jar` tử Internet.
- Đặt Jar file vào vị trí thư mục hiện hành và thực thi lệnh Terminal:
```bash
  $ java -jar plantuml.jar "*/**/*.puml"
```
Hệ thống sẽ duyệt các Folder và cấp phát ra các tệp ảnh (`.png`) tương ứng. 
