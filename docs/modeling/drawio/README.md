# EduCenter Draw.io Diagrams

Các sơ đồ trong thư mục này được dựng từ:

- [BA_SYSTEM_ANALYSIS_REPORT_EDUCENTER.md](/Users/hant/golang/doan/docs/BA_SYSTEM_ANALYSIS_REPORT_EDUCENTER.md)
- [ERD_DRAFTING_PACKAGE_EDUCENTER.md](/Users/hant/golang/doan/docs/ERD_DRAFTING_PACKAGE_EDUCENTER.md)

Phạm vi:

- `educenter_usecase.drawio.xml`: use case tổng quan cho các actor chính và các luồng cốt lõi.
- `educenter_bpmn_activity_open_class_to_lesson.drawio.xml`: BPMN-level activity cho quy trình mở lớp tới commit lesson.
- `educenter_erd_logic.drawio.xml`: logical ERD lõi và vận hành học vụ.

Ghi chú:

- Sơ đồ ưu tiên phản ánh current-state đã xác nhận từ code và các target-state đã được BA report đánh dấu `Strongly inferred from code`.
- Logical ERD tập trung vào miền đào tạo, lớp học, scheduling và học vụ sau lesson; không mở rộng toàn bộ cụm auth/audit/material để tránh quá tải.
