# Use Case Diagrams - Hệ thống Quản lý Trung tâm Dạy thêm tích hợp AI

Thư mục này chứa toàn bộ các sơ đồ Use Case và Flow Diagram dưới định dạng **PlantUML (.puml)** đã đáp ứng yêu cầu Task 1.2 trong file `PROJECT_TASKS.md`.

## Danh sách Sơ đồ

1. **[01_overview.puml](./01_overview.puml)**: Sơ đồ Use Case tổng quan hiển thị nhóm phân hệ (Modules) và các tác nhân (Actors) chính.
2. **[02_admin.puml](./02_admin.puml)**: Đặc tả chi tiết cho Admin liên quan đến Core Management (Tài khoản, Hồ sơ, Class, Room, Program...) và tương tác với hệ Scheduling.
3. **[03_teacher.puml](./03_teacher.puml)**: Các Use case học vụ hàng ngày của Giáo viên (thời khóa biểu, điểm danh, sổ đầu bài) và tải tài liệu tương tác với AI Audit.
4. **[04_student_parent.puml](./04_student_parent.puml)**: Use case cho Cổng thông tin (Portal) của Học sinh/Phụ huynh và tích hợp AI Chatbot.
5. **[05_compliance.puml](./05_compliance.puml)**: Các chức năng phê duyệt kiểm duyệt rủi ro của nhân sự Compliance Officer.
6. **[06_ai_agent.puml](./06_ai_agent.puml)**: Use case ngầm định (System/Backstage) của AI Agent để tự động OCR, quét ngữ nghĩa qua Gemini trích nhãn, và thuật toán Xếp lịch (Backtracking CSP).
7. **[07_scheduling_flow.puml](./07_scheduling_flow.puml)**: Luồng Use Case nghiệp vụ từ lúc Admin kích hoạt Xếp lịch cho tới khi hiển thị Preview, giải quyết xung đột (conflict resolution) và chốt lịch.
8. **[08_material_audit_flow.puml](./08_material_audit_flow.puml)**: Luồng Use Case xoay quanh quá trình lọc tài liệu (Upload -> OCR -> Call Gemini -> Labeling -> Compliance Review -> Notification).

## Hướng dẫn Render PlantUML

### Cách 1: Sử dụng VS Code Extension
- Cài đặt extension: **PlantUML** (của *jebbs*).
- Mở một file `.puml` bất kỳ trong thư mục này.
- Bấm tổ hợp phím `Alt + D` (hoặc `Option + D` trên Mac) để xem trước trực tiếp (Preview). 

### Cách 2: Render trực tuyến
- Copy nội dung text bất kỳ trong file `.puml` nào.
- Dán vào các bộ sinh PlantUML chuẩn online như:
  - [PlantText](https://www.planttext.com/)
  - [PlantUML Server](http://www.plantuml.com/plantuml/uml/)

### Cách 3: Sử dụng CLI với `.jar` (Chạy độc lập local)
- Đảm bảo máy có sẵn Java (JRE/JDK) và Graphviz.
- Tải `plantuml.jar` từ trang chủ PlantUML.
- Mở Terminal/CMD ngay ở thư mục hiện tại:
  ```bash
  java -jar plantuml.jar *.puml
  ```
  Lệnh này sẽ render đồng loạt toàn bộ file thiết kế sang định dạng ảnh (thường cấu hình xuất `.png` hoặc `.svg`) để nhúng báo cáo Word/Latex.
