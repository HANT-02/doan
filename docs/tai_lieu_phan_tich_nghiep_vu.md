# TÀI LIỆU PHÂN TÍCH NGHIỆP VỤ — HỆ THỐNG EDUCENTER

**Hệ thống quản lý trung tâm dạy thêm**
**Dự án:** Đồ án tốt nghiệp — Nguyễn Thế Hà 61165 CS2
**Phiên bản:** 1.0 · **Ngày:** 13/04/2026
**Tuân thủ:** Thông tư 29/2024/TT-BGDĐT

---

## Mục lục

1. [Tổng quan hệ thống](#1-tổng-quan-hệ-thống)
2. [Phạm vi & Bản đồ miền nghiệp vụ](#2-phạm-vi--bản-đồ-miền-nghiệp-vụ)
3. [Danh mục tác nhân](#3-danh-mục-tác-nhân)
4. [Phân rã chức năng](#4-phân-rã-chức-năng)
5. [Danh mục Use Case](#5-danh-mục-use-case)
6. [Luồng nghiệp vụ chính](#6-luồng-nghiệp-vụ-chính)
7. [Quy tắc nghiệp vụ & Ràng buộc](#7-quy-tắc-nghiệp-vụ--ràng-buộc)
8. [Danh mục thực thể (ERD)](#8-danh-mục-thực-thể-erd)
9. [Ma trận quan hệ](#9-ma-trận-quan-hệ)
10. [Vòng đời trạng thái](#10-vòng-đời-trạng-thái)
11. [Câu hỏi cần xác nhận với BA](#11-câu-hỏi-cần-xác-nhận-với-ba)
12. [Phụ lục — Bản đồ bằng chứng kỹ thuật](#12-phụ-lục--bản-đồ-bằng-chứng-kỹ-thuật)

---

## 1. Tổng quan hệ thống

### 1.1 Mô tả

EduCenter là hệ thống quản lý vận hành dành cho **trung tâm dạy thêm** (single-tenant), số hóa các hoạt động cốt lõi bao gồm: quản lý học viên, giáo viên, khóa học, lớp học, phòng học, ca học; **xếp lịch thông minh** bằng thuật toán tối ưu hóa; và **dự báo học viên có nguy cơ học kém** (đang phát triển).

### 1.2 Giá trị nghiệp vụ

| # | Giá trị | Mô tả |
|---|---------|-------|
| 1 | Tự động hóa xếp lịch | Giảm tải cho quản trị viên, tránh xung đột phòng/giáo viên |
| 2 | So sánh thuật toán | Cung cấp bằng chứng thực nghiệm để chọn phương pháp tối ưu |
| 3 | Theo dõi học tập | Điểm danh, đánh giá, phát hiện sớm nguy cơ học kém |
| 4 | Tuân thủ pháp quy | Đáp ứng yêu cầu Thông tư 29/2024/TT-BGDĐT |

### 1.3 Nhóm người dùng chính

- **Quản trị viên (ADMIN):** Toàn quyền quản lý hệ thống
- **Giáo viên (TEACHER):** Tải tài liệu, giảng dạy, ghi nhận kết quả
- **Học viên (STUDENT):** Xem thông tin lớp, khóa học, xin phép nghỉ
- **Khách (Anonymous):** Xem danh mục công khai

---

## 2. Phạm vi & Bản đồ miền nghiệp vụ

| Mã | Miền nghiệp vụ | Mục đích | Trạng thái |
|----|----------------|----------|------------|
| D-01 | Xác thực & Quản lý tài khoản | Đăng ký, đăng nhập, quên mật khẩu, OTP | ✅ Hoàn thành |
| D-02 | Quản lý học viên | Hồ sơ, thông tin liên hệ, cấp lớp | ✅ Hoàn thành |
| D-03 | Quản lý giáo viên | Hồ sơ, loại hình tuyển dụng, lịch dạy, thống kê giờ dạy | ✅ Hoàn thành |
| D-04 | Quản lý khóa học | Danh mục môn học, số buổi, giá | ✅ Hoàn thành |
| D-05 | Quản lý chương trình đào tạo | Nhóm khóa học theo track, mục tiêu, chuẩn đầu ra | ✅ Hoàn thành |
| D-06 | Quản lý lớp học & Đăng ký | Mở lớp, ghi danh, phân công giáo viên | ✅ Hoàn thành |
| D-07 | Quản lý phòng học | Danh mục phòng, sức chứa | ✅ Hoàn thành |
| D-08 | Quản lý ca học (Shift) | Định nghĩa khung giờ chuẩn | ✅ Hoàn thành |
| D-09 | Xếp lịch thông minh | Preview → Commit → Lesson, 3 thuật toán, benchmark | ✅ Hoàn thành |
| D-10 | Buổi học & Điểm danh | Ghi nhận buổi dạy, điểm danh từng học viên | ⚠️ Thực thể có, API chưa đầy đủ |
| D-11 | Kết quả học tập | Điểm bài tập, thái độ, tham gia, tổng điểm | ⚠️ Thực thể có, API chưa đầy đủ |
| D-12 | Xin phép nghỉ/đi muộn/về sớm | Đơn xin, duyệt/từ chối | ⚠️ Thực thể có, chưa có API |
| D-13 | Tài liệu giảng dạy & Kiểm duyệt | Upload, OCR, phân loại AI, duyệt thủ công | ⚠️ Triển khai một phần (stub) |
| D-14 | Dự báo nguy cơ học kém | ML classification, cảnh báo sớm | ⬜ Đang lên kế hoạch |

---

## 3. Danh mục tác nhân

| Mã | Tác nhân | Loại | Mô tả | Hành động chính | Phạm vi quyền |
|----|----------|------|-------|-----------------|---------------|
| A-01 | Quản trị viên (ADMIN) | Nội bộ | Người quản lý trung tâm | CRUD toàn bộ, xếp lịch, ghi danh, benchmark | Toàn quyền |
| A-02 | Giáo viên (TEACHER) | Nội bộ | Giảng viên / gia sư | Upload tài liệu, xem lớp & học viên | Upload tài liệu, đọc hầu hết module |
| A-03 | Học viên (STUDENT) | Nội bộ | Người học đã đăng ký | Xem lớp, khóa học, xin phép nghỉ | Chủ yếu đọc; xin phép (chưa có API) |
| A-04 | Khách (Anonymous) | Bên ngoài | Người truy cập chưa đăng nhập | Xem danh mục khóa học, giáo viên, lớp, chương trình | Chỉ đọc các endpoint công khai |
| A-05 | Hệ thống gửi email | Hệ thống | Dịch vụ SMTP | Gửi OTP, link đặt lại mật khẩu | Gửi email ra ngoài |
| A-06 | Engine xếp lịch | Hệ thống | Bộ giải thuật toán | Tạo preview, benchmark solver | Xử lý nội bộ |

---

## 4. Phân rã chức năng

```
EduCenter
│
├── 1. Xác thực & Quản lý tài khoản
│   ├── 1.1 Đăng ký & Xác minh email (OTP)
│   ├── 1.2 Đăng nhập / Đăng xuất / Làm mới token
│   └── 1.3 Quên / Đặt lại / Đổi mật khẩu
│
├── 2. Dữ liệu danh mục
│   ├── 2.1 Quản lý học viên (CRUD + tìm kiếm)
│   ├── 2.2 Quản lý giáo viên (CRUD + lịch dạy + thống kê giờ)
│   ├── 2.3 Quản lý khóa học (CRUD)
│   ├── 2.4 Quản lý chương trình (CRUD + gán/gỡ khóa học)
│   ├── 2.5 Quản lý phòng học (CRUD)
│   └── 2.6 Cấu hình ca học (CRUD)
│
├── 3. Vận hành lớp học
│   ├── 3.1 Vòng đời lớp học (tạo/sửa/xóa/danh sách)
│   ├── 3.2 Ghi danh học viên (ghi danh/rút tên hàng loạt)
│   └── 3.3 Phân công giáo viên
│
├── 4. Xếp lịch thông minh ★
│   ├── 4.1 Tạo preview lịch (solver mặc định)
│   ├── 4.2 Xem / lấy preview (theo ID hoặc mới nhất)
│   ├── 4.3 Xác nhận preview → tạo buổi học (Lesson)
│   └── 4.4 Benchmark thuật toán (GraphColoring, CP-SAT, TabuSearch)
│
├── 5. Hoạt động học thuật (thực thể có, API hạn chế)
│   ├── 5.1 Quản lý buổi học (tạo tự động từ commit)
│   ├── 5.2 Điểm danh
│   ├── 5.3 Tổng kết buổi học
│   ├── 5.4 Kết quả học tập & chấm điểm
│   └── 5.5 Xin phép nghỉ (chưa có API)
│
├── 6. Tài liệu giảng dạy (phạm vi phụ)
│   ├── 6.1 Upload & lưu trữ
│   ├── 6.2 Kiểm duyệt nội dung (OCR + AI stub)
│   └── 6.3 Duyệt thủ công
│
└── 7. Dự báo phân tích (chưa triển khai)
    ├── 7.1 Phân loại nguy cơ học kém
    ├── 7.2 Pipeline ML
    └── 7.3 Dashboard cảnh báo
```

### Ma trận quyền theo nhánh chức năng

| Nhánh | Đọc công khai | Đọc (auth) | Ghi (ADMIN) | Ghi (TEACHER) |
|-------|:---:|:---:|:---:|:---:|
| 1. Xác thực | — | ✅ | — | — |
| 2.1 Học viên | — | ✅ CRUD† | ✅† | ✅† |
| 2.2 Giáo viên | ✅ | — | ✅ CUD | — |
| 2.3 Khóa học | ✅ | — | ✅ CUD | — |
| 2.4 Chương trình | ✅ | — | ✅ CUD | — |
| 2.5 Phòng học | ✅ | — | ✅ CUD | — |
| 2.6 Ca học | — | — | ✅ Toàn bộ | — |
| 3. Lớp học | ✅ | — | ✅ CUD + ghi danh | — |
| 4. Xếp lịch | — | — | ✅ Toàn bộ | — |
| 6. Tài liệu | — | ✅ | — | ✅ Upload |

> † Không có kiểm tra vai trò — bất kỳ người dùng đã đăng nhập đều thao tác được (xem Q-03)

---

## 5. Danh mục Use Case

| Mã | Tên | Tác nhân | Miền | Độ tin cậy |
|----|-----|----------|------|:---:|
| UC-01 | Đăng nhập | Tất cả | D-01 | Cao |
| UC-02 | Đăng ký tài khoản | Khách | D-01 | Cao |
| UC-03 | Xác minh OTP email | Người dùng | D-01 | Cao |
| UC-04 | Quên mật khẩu | Tất cả | D-01 | Cao |
| UC-05 | Đặt lại mật khẩu | Tất cả | D-01 | Cao |
| UC-06 | Đổi mật khẩu | Đã đăng nhập | D-01 | Cao |
| UC-07 | Xem hồ sơ cá nhân | Đã đăng nhập | D-01 | Cao |
| UC-08 | Tạo học viên | ADMIN† | D-02 | Cao |
| UC-09 | Tìm kiếm / liệt kê học viên | ADMIN† | D-02 | Cao |
| UC-10 | Cập nhật học viên | ADMIN† | D-02 | Cao |
| UC-11 | Xóa học viên | ADMIN† | D-02 | Cao |
| UC-12 | Tạo giáo viên | ADMIN | D-03 | Cao |
| UC-13 | Xem lịch dạy giáo viên | Tất cả | D-03 | Cao |
| UC-14 | Xem thống kê giờ dạy | Tất cả | D-03 | Cao |
| UC-15 | Tạo khóa học | ADMIN | D-04 | Cao |
| UC-16 | Tạo chương trình đào tạo | ADMIN | D-05 | Cao |
| UC-17 | Gán / gỡ khóa học vào chương trình | ADMIN | D-05 | Cao |
| UC-18 | Tạo lớp học | ADMIN | D-06 | Cao |
| UC-19 | Ghi danh học viên (hàng loạt) | ADMIN | D-06 | Cao |
| UC-20 | Rút học viên khỏi lớp (hàng loạt) | ADMIN | D-06 | Cao |
| UC-21 | Phân công giáo viên cho lớp | ADMIN | D-06 | Cao |
| UC-22 | Xem danh sách lớp | Tất cả | D-06 | Cao |
| UC-23 | Tạo phòng học | ADMIN | D-07 | Cao |
| UC-24 | Tạo ca học | ADMIN | D-08 | Cao |
| UC-25 | Tạo preview xếp lịch | ADMIN | D-09 | Cao |
| UC-26 | Xem preview xếp lịch | ADMIN | D-09 | Cao |
| UC-27 | Xác nhận lịch → tạo buổi học | ADMIN | D-09 | Cao |
| UC-28 | Benchmark thuật toán xếp lịch | ADMIN | D-09 | Cao |
| UC-29 | Upload tài liệu giảng dạy | TEACHER | D-13 | Cao |
| UC-30 | Duyệt tài liệu bị gắn cờ | Đã đăng nhập | D-13 | Trung bình |
| UC-31 | Xin phép nghỉ / đi muộn / về sớm | STUDENT | D-12 | Trung bình (chưa có API) |
| UC-32 | Duyệt / từ chối đơn xin phép | ADMIN | D-12 | Trung bình (chưa có API) |
| UC-33 | Điểm danh buổi học | TEACHER | D-10 | Trung bình (API hạn chế) |
| UC-34 | Ghi tổng kết buổi học | TEACHER | D-10 | Trung bình (API hạn chế) |
| UC-35 | Ghi kết quả học tập | TEACHER | D-11 | Trung bình (API hạn chế) |

> † Thực tế không kiểm tra vai trò — xem Q-03

### 5.1 Luồng chi tiết — UC-25: Tạo preview xếp lịch

| Mục | Nội dung |
|-----|---------|
| **Tác nhân** | ADMIN |
| **Tiền điều kiện** | Lớp học có ClassSchedule (ngày + ca); phòng, ca học đang hoạt động |
| **Đầu vào** | Khoảng ngày (từ–đến), tùy chọn lọc theo lớp/giáo viên/phòng |
| **Luồng chính** | 1) Xác thực khoảng ngày hợp lệ → 2) Lấy dữ liệu lớp, lịch, phòng, ca → 3) Xây dựng input cho solver → 4) Chạy solver mặc định → 5) Lưu kết quả vào bộ nhớ (PreviewStore) → 6) Trả về assignments, conflicts, summary |
| **Luồng ngoại lệ** | Sai định dạng ngày → lỗi 400; lỗi solver → lỗi 500 |
| **Đầu ra** | `PreviewResult {runID, status, assignments[], conflicts[], summary}` |
| **Quy tắc liên quan** | BR-16, BR-23, BR-24, BR-25 |

### 5.2 Luồng chi tiết — UC-27: Xác nhận lịch

| Mục | Nội dung |
|-----|---------|
| **Tác nhân** | ADMIN |
| **Tiền điều kiện** | Preview hợp lệ tồn tại trong PreviewStore |
| **Đầu vào** | `run_id` của preview đã tạo |
| **Luồng chính** | 1) Lấy preview theo runID → 2) Tạo bản ghi Lesson cho mỗi assignment (lớp, giáo viên, ngày, phòng) → 3) Trả kết quả commit |
| **Luồng ngoại lệ** | Không tìm thấy runID → lỗi 404 |
| **Tác dụng phụ** | Tạo bản ghi Lesson trong cơ sở dữ liệu |

---

## 6. Luồng nghiệp vụ chính

### WF-01: Vòng đời xếp lịch lớp học

| Bước | Hành động | Tác nhân | Miền |
|:---:|----------|:---:|:---:|
| 1 | Tạo phòng học | ADMIN | D-07 |
| 2 | Tạo ca học (khung giờ) | ADMIN | D-08 |
| 3 | Tạo khóa học (môn, giá) | ADMIN | D-04 |
| 4 | Tạo chương trình → gán khóa học | ADMIN | D-05 |
| 5 | Mở lớp (chương trình + khóa + giáo viên + phòng + lịch tuần) | ADMIN | D-06 |
| 6 | Ghi danh học viên vào lớp | ADMIN | D-06 |
| 7 | Tạo preview xếp lịch | ADMIN | D-09 |
| 8 | Kiểm tra xung đột, điều chỉnh nếu cần | ADMIN | D-09 |
| 9 | Xác nhận preview → tạo buổi học | ADMIN | D-09 |
| 10 | Giảng dạy, điểm danh, tổng kết | TEACHER | D-10 |

**Điểm lỗi:** Bước 7 có thể tạo buổi không xếp được (xung đột). Admin cần điều chỉnh (đổi phòng/ca/giáo viên) và chạy lại.

### WF-02: Đăng ký & Xác minh tài khoản

1. Người dùng gửi đăng ký (email, họ tên, mật khẩu)
2. Hệ thống tạo User (vai trò mặc định: STUDENT), mã hóa mật khẩu
3. Tạo OTP, gửi email xác minh (hết hạn 5 phút)
4. Người dùng nhập OTP → hệ thống xác nhận
5. Tài khoản sẵn sàng đăng nhập

### WF-03: Đặt lại mật khẩu

1. Người dùng yêu cầu quên mật khẩu → hệ thống gửi email chứa token
2. Người dùng gửi mật khẩu mới kèm token
3. Hệ thống xác thực token, cập nhật mật khẩu

### WF-04: Upload & Kiểm duyệt tài liệu

1. Giáo viên upload file → lưu vào hệ thống tập tin
2. OCR trích xuất văn bản (hiện tại: stub)
3. AI phân loại nội dung → gán nhãn AN TOÀN / CẢNH BÁO / NGUY HIỂM (hiện tại: stub)
4. Nếu CẢNH BÁO/NGUY HIỂM → hiển thị trong danh sách cần duyệt
5. Người duyệt ra quyết định phê duyệt/từ chối

---

## 7. Quy tắc nghiệp vụ & Ràng buộc

| Mã | Quy tắc | Thực thể | Mức độ |
|----|---------|----------|--------|
| BR-01 | Email phải là duy nhất trong hệ thống | User | Lỗi |
| BR-02 | Mã (code) của mỗi thực thể phải là duy nhất | User, Student, Teacher, Course, Program, Class, Shift | Lỗi |
| BR-03 | Vai trò mặc định khi đăng ký là STUDENT | User | Thông tin |
| BR-04 | Trạng thái mặc định ghi danh là ĐÃ NỘP ĐƠN (APPLIED) | Enrollment | Nghiệp vụ |
| BR-05 | Trạng thái mặc định lớp học là MỞ (OPEN) | Class | Thông tin |
| BR-06 | Trạng thái mặc định đơn xin phép là CHỜ DUYỆT (PENDING) | LeaveRequest | Nghiệp vụ |
| BR-07 | Trạng thái mặc định tài liệu là ĐÃ TẢI LÊN (UPLOADED) | Material | Nghiệp vụ |
| BR-08 | Mật khẩu tối thiểu 8 ký tự, gồm chữ hoa, chữ thường, số, ký tự đặc biệt | User | Lỗi |
| BR-09 | Số điện thoại định dạng Việt Nam (bắt đầu 0 hoặc +84, 9-10 chữ số) | Student, Teacher | Lỗi |
| BR-10 | Ngày kết thúc phải ≥ ngày bắt đầu trong xếp lịch | Scheduling | Lỗi (400) |
| BR-11 | Chỉ tài khoản đang hoạt động (isActive = true) mới được đăng nhập | User | Từ chối truy cập |
| BR-12 | Mật khẩu xác thực bằng bcrypt | User | Bảo mật |
| BR-13 | Token truy cập mặc định 24h, token làm mới mặc định 7 ngày | User | Cấu hình |
| BR-14 | OTP hết hạn sau 5 phút | UserOTP | Bảo mật |
| BR-15 | Lịch học tuần (ClassSchedule) bắt buộc có ca học (shiftID) | ClassSchedule | Lỗi |
| BR-16 | Ghi danh xóa theo cascade khi xóa lớp hoặc học viên | Enrollment | Dữ liệu |
| BR-17 | Chương trình có 3 track: HỖ TRỢ / CƠ BẢN / NÂNG CAO | Program | Enum |
| BR-18 | Giáo viên có 2 loại: TOÀN THỜI GIAN / BÁN THỜI GIAN | Teacher | Enum |
| BR-19 | Loại xin phép: NGHỈ / ĐI MUỘN / VỀ SỚM | LeaveRequest | Enum |
| BR-20 | Không được xếp cùng giáo viên vào hai phòng cùng lúc | Scheduling | Ràng buộc cứng |
| BR-21 | Không được xếp hai lớp vào cùng phòng cùng lúc | Scheduling | Ràng buộc cứng |
| BR-22 | Sức chứa phòng phải ≥ sĩ số lớp | Scheduling | Ràng buộc cứng |
| BR-23 | Lớp học và khóa học đều có trường giá, nhưng ưu tiên giá nào chưa rõ | Class, Course | ⚠️ Cần xác nhận (Q-02) |

## 8. Danh mục thực thể (ERD)

### 8.1 Phạm vi chính (17 thực thể)

| Mã | Thực thể | Bảng | Thuộc tính chính | Trạng thái / Enum | Xóa mềm |
|----|----------|------|------------------|-------------------|:---:|
| E-01 | Tài khoản | `users` | id, code★, fullName, email★, password, role, isActive | role: ADMIN/TEACHER/STUDENT | ✅ |
| E-02 | Học viên | `students` | id, code★, fullName, email, phone, guardianPhone, gradeLevel, dob, gender, address | status: ACTIVE | ✅ |
| E-03 | Giáo viên | `teachers` | id, code★, fullName, email, phone, isSchoolTeacher, schoolName, employmentType | status: ACTIVE; empType: FULL_TIME/PART_TIME | ✅ |
| E-04 | Khóa học | `courses` | id, code★, name, description, gradeLevel, subject, sessionCount, sessionDurationMin, totalHours, price | status: ACTIVE | ✅ |
| E-05 | Chương trình | `programs` | id, code★, name, track, effectiveFrom/To, createdByID, approvedByID, approvalNote, publishedAt, archivedAt | track: SUPPORT/BASIC/ADVANCED | ✅ |
| E-06 | CT–KH (mapping) | `program_courses` | id, programID, courseID | — | ❌ |
| E-07 | Phòng học | `rooms` | id, code★, name, capacity, address | — | ❌ |
| E-08 | Ca học | `shifts` | id, code★, name, startTime, endTime, durationMinutes, sessionType, isActive | sessionType: (chưa rõ enum) | ❌ |
| E-09 | Lớp học | `classes` | id, code★, name, notes, startDate, endDate, maxStudents, price, programID?, courseID?, teacherID?, roomID? | status: OPEN | ✅ |
| E-10 | Lịch tuần | `class_schedules` | id, classID, dayOfWeek, shiftID, roomID? | — | ❌ |
| E-11 | Ghi danh | `enrollments` | id, classID, studentID, approvedAt, rejectedAt | status: APPLIED | ❌ |
| E-12 | Buổi học | `lessons` | id, classID, teacherID, dateStart, dateEnd, roomID, notes | — | ❌ |
| E-13 | Điểm danh | `attendances` | id, lessonID, studentID, status(int), note, markedAt | status: int (giá trị chưa rõ) | ❌ |
| E-14 | Tổng kết buổi | `lesson_summaries` | id, lessonID (unique), topic, lessonContent, classFeedback, homework, homeworkDeadline, teacherNotes, createdByID | — | ❌ |
| E-15 | Kết quả học tập | `academic_records` | id, lessonSummaryID, studentID, homeworkCompleted, homeworkScore, attitudeRating, participationScore, totalScore, isCompleted | — | ❌ |
| E-16 | Đơn xin phép | `leave_requests` | id, studentID, leaveType, applyDate, lateMin, earlyMin, reason, documents[], classID?, lessonID?, approvedByID?, rejectionReason | status: PENDING; leaveType: LEAVE/LATE/EARLY | ❌ |
| E-17 | Tư vấn tuyển sinh | `consultations` | id, fullName, phone, gradeLevel, notes | status: PENDING | ❌ |

> ★ = ràng buộc unique

### 8.2 Phạm vi mở rộng (7 thực thể)

| Mã | Thực thể | Bảng | Mục đích |
|----|----------|------|----------|
| E-18 | OTP người dùng | `user_otps` | Mã xác minh email (hết hạn 5 phút) |
| E-19 | Đặt lại mật khẩu | `password_resets` | Token reset (có hạn, ghi IP/user agent) |
| E-20 | Mục tiêu CT | `objectives` | Mục tiêu học tập của chương trình |
| E-21 | Chuẩn đầu ra | `outcomes` | Chuẩn đầu ra gắn với mục tiêu |
| E-22 | Tài liệu | `materials` | File giảng dạy do giáo viên upload |
| E-23 | Nhật ký kiểm duyệt | `audit_logs` | Kết quả OCR + AI cho tài liệu |
| E-24 | Quyết định duyệt | `approval_decisions` | Quyết định con người về tài liệu |
| E-25 | Nhãn phân loại | `labels` | Nhãn: AN TOÀN / CẢNH BÁO / NGUY HIỂM |

---

## 9. Ma trận quan hệ

| Thực thể A | Thực thể B | Loại | Ý nghĩa nghiệp vụ | Bắt buộc | Cascade |
|-----------|-----------|------|---------------------|:---:|:---:|
| E-01 Tài khoản | E-18 OTP | 1→N | Tài khoản có nhiều mã OTP | ✅ | — |
| E-01 Tài khoản | E-19 Reset | 1→N | Tài khoản có nhiều token reset | ✅ | — |
| E-05 Chương trình | E-04 Khóa học | N↔N | Chương trình chứa nhiều khóa (qua E-06) | — | — |
| E-05 Chương trình | E-20 Mục tiêu | 1→N | CT định nghĩa mục tiêu | ✅ | CASCADE |
| E-05 Chương trình | E-21 Chuẩn đầu ra | 1→N | CT định nghĩa chuẩn đầu ra | ✅ | CASCADE |
| E-09 Lớp | E-05 Chương trình | N→1 | Lớp thuộc chương trình | Tùy chọn | — |
| E-09 Lớp | E-04 Khóa học | N→1 | Lớp dạy khóa học | Tùy chọn | — |
| E-09 Lớp | E-03 Giáo viên | N→1 | Lớp phân công cho GV | Tùy chọn | — |
| E-09 Lớp | E-07 Phòng | N→1 | Phòng mặc định | Tùy chọn | — |
| E-09 Lớp | E-10 Lịch tuần | 1→N | Lớp có mẫu lịch hàng tuần | ✅ | CASCADE |
| E-10 Lịch tuần | E-08 Ca học | N→1 | Lịch dùng ca chuẩn | ✅ | — |
| E-09 Lớp | E-11 Ghi danh | 1→N | Lớp có học viên ghi danh | ✅ | CASCADE |
| E-02 Học viên | E-11 Ghi danh | 1→N | HV ghi danh vào lớp | ✅ | CASCADE |
| E-09 Lớp | E-12 Buổi học | 1→N | Lớp có nhiều buổi | ✅ | CASCADE |
| E-12 Buổi học | E-03 Giáo viên | N→1 | Buổi do GV dạy | Tùy chọn | — |
| E-12 Buổi học | E-07 Phòng | N→1 | Buổi tại phòng | Tùy chọn | — |
| E-12 Buổi học | E-13 Điểm danh | 1→N | Buổi ghi nhận điểm danh | ✅ | CASCADE |
| E-02 Học viên | E-13 Điểm danh | 1→N | HV có lịch sử điểm danh | ✅ | CASCADE |
| E-12 Buổi học | E-14 Tổng kết | 1→1 | Mỗi buổi có 1 tổng kết | ✅ | CASCADE |
| E-14 Tổng kết | E-15 Kết quả HT | 1→N | Tổng kết chứa điểm từng HV | ✅ | CASCADE |
| E-02 Học viên | E-16 Đơn xin phép | 1→N | HV có đơn xin phép | ✅ | CASCADE |
| E-03 Giáo viên | E-22 Tài liệu | 1→N | GV upload tài liệu | ✅ | — |
| E-22 Tài liệu | E-23 Nhật ký KD | 1→N | Tài liệu có lịch sử kiểm duyệt | ✅ | CASCADE |
| E-22 Tài liệu | E-24 Quyết định | 1→N | Tài liệu có quyết định duyệt | ✅ | CASCADE |

---

## 10. Vòng đời trạng thái

### Ghi danh (Enrollment)
```
APPLIED (đã nộp đơn) ──duyệt──→ APPROVED (approvedAt được set)
APPLIED              ──từ chối──→ REJECTED (rejectedAt được set)
```

### Đơn xin phép (LeaveRequest)
```
PENDING (chờ duyệt) ──duyệt──→ APPROVED (approvedByID, approvedAt)
PENDING              ──từ chối──→ REJECTED (rejectionReason)
```

### Tài liệu (Material) — suy luận
```
UPLOADED ──kiểm duyệt AI──→ gán nhãn (AN TOÀN / CẢNH BÁO / NGUY HIỂM)
         ──duyệt thủ công──→ ApprovalDecision (phê duyệt / từ chối)
```

### Chương trình (Program) — suy luận từ trường timestamp
```
[tạo] ──xuất bản──→ publishedAt ──lưu trữ──→ archivedAt
```
> ⚠️ Không có trường status rõ ràng. Vòng đời suy luận từ publishedAt/archivedAt.

---

## 11. Câu hỏi cần xác nhận với BA

### 🔴 Nghiêm trọng — ảnh hưởng đến đặc tả chính

| Mã | Câu hỏi | Lý do | Bên liên quan |
|----|---------|-------|---------------|
| Q-01 | Ghi danh chuyển từ APPLIED → APPROVED/REJECTED bằng cách nào? Không tìm thấy API duyệt. | Học viên không thể tham gia lớp nếu không được duyệt — luồng nghiệp vụ cốt lõi bị thiếu. | Product Owner |
| Q-02 | Giá nào là giá thanh toán — giá trên Lớp hay trên Khóa? Hai trường cùng tồn tại với giá trị khác nhau. | Ảnh hưởng trực tiếp đến doanh thu, không thể định nghĩa UC thanh toán. | Tài chính / PO |
| Q-03 | Các route CRUD Học viên không kiểm tra vai trò — bất kỳ ai đã đăng nhập đều thao tác được. Có chủ đích? | Lỗ hổng bảo mật — STUDENT có thể xóa hồ sơ học viên khác. | Tech Lead / PO |
| Q-04 | Không có module thanh toán/hóa đơn. Học phí được thu và theo dõi như thế nào? | Luồng doanh thu chính của trung tâm không được thể hiện trong hệ thống. | Tài chính / PO |

### 🟡 Quan trọng — ảnh hưởng đến tính đầy đủ

| Mã | Câu hỏi | Lý do | Bên liên quan |
|----|---------|-------|---------------|
| Q-05 | Thực thể Đơn xin phép (LeaveRequest) không có API. Có nằm trong phạm vi hiện tại? | Thực thể đầy đủ (3 loại, workflow duyệt) nhưng không truy cập được. | Product Owner |
| Q-06 | Thực thể Tư vấn tuyển sinh (Consultation) không có API. Quản lý lead có trong kế hoạch? | Mất khả năng theo dõi học viên tiềm năng. | Marketing / PO |
| Q-07 | Chương trình có trường duyệt (approvedByID, approvalNote) nhưng không có endpoint duyệt. | Quy trình quản trị chương trình đào tạo chưa rõ. | Giám đốc đào tạo |
| Q-08 | Preview xếp lịch lưu trong bộ nhớ tạm — khởi động lại server sẽ mất. Chấp nhận được? | Rủi ro mất dữ liệu giữa việc tạo preview và xác nhận. | Tech Lead |
| Q-09 | Mục tiêu (Objective) và Chuẩn đầu ra (Outcome) không có API. Quản lý thế nào? | Cấu trúc chương trình đào tạo tồn tại trong DB nhưng không quản lý được qua hệ thống. | Giám đốc đào tạo |
| Q-10 | Trạng thái Lớp (Class.Status) ngoài OPEN còn giá trị nào? (CLOSED? CANCELLED?) | Cần định nghĩa vòng đời lớp học hoàn chỉnh. | Product Owner |
| Q-11 | Trạng thái Điểm danh (Attendance.Status) là số nguyên — các giá trị hợp lệ là gì? | Không có bảng mã — không định nghĩa được quy tắc điểm danh. | Product Owner |

### 🔵 Thông tin — làm rõ thiết kế

| Mã | Câu hỏi | Lý do | Bên liên quan |
|----|---------|-------|---------------|
| Q-12 | Hệ thống tuân thủ Thông tư 29/2024 cụ thể ở điểm nào? | Yêu cầu pháp quy được đề cập nhưng không có logic tuân thủ rõ ràng. | Pháp chế |
| Q-13 | API V2 auth là stub (panic). Chiến lược phiên bản API? | V2 chỉ có Login/Logout/Refresh, các method còn lại chưa triển khai. | Tech Lead |
| Q-14 | Phòng học không có soft delete và không có cờ isActive. Ngừng sử dụng phòng thế nào? | Chỉ có thể xóa cứng — không thể đánh dấu phòng tạm ngừng. | Vận hành / PO |
| Q-15 | Không có hệ thống thông báo ngoài email SMTP. Push/SMS/in-app có trong kế hoạch? | Kênh liên lạc với học viên/phụ huynh bị hạn chế. | Product Owner |
| Q-16 | Endpoint duyệt tài liệu không kiểm tra vai trò — ai được phép duyệt? | Bất kỳ người đã đăng nhập đều ra quyết định duyệt được. | PO / Bảo mật |
| Q-17 | Trường isCompleted trên AcademicRecord — điều gì kích hoạt hoàn thành? | Không tìm thấy logic set cờ này — ảnh hưởng báo cáo kết quả. | Product Owner |
| Q-18 | Cờ isSchoolTeacher trên Giáo viên — có ảnh hưởng đến ràng buộc xếp lịch không? | GV đang dạy ở trường chính quy có thể bị giới hạn giờ dạy, nhưng solver không xét cờ này. | Giám đốc đào tạo |

---

## 12. Phụ lục — Bản đồ bằng chứng kỹ thuật

| Phát hiện | Tập tin nguồn | Giải thích |
|-----------|--------------|------------|
| 25 thực thể nghiệp vụ | `internal/entities/*.go` | Struct GORM với tag DB |
| Phân quyền theo vai trò (ADMIN/TEACHER/STUDENT) | `cmd/http/middleware/auth.go` | JWT chứa role; `RoleMiddleware` kiểm tra |
| 4 solver xếp lịch | `internal/services/scheduling/solver_catalog.go` | Interface `SchedulingSolver` + 4 implementation |
| Phân chia route public vs ADMIN | `cmd/http/controllers/*/controller.go` | Pattern áp dụng middleware |
| Email SMTP cho OTP | `internal/services/mailer/mailer.go` | SMTP trực tiếp với hỗ trợ TLS |
| Dịch vụ AI kiểm duyệt là stub | `internal/services/audit/services.go` | `StubOCRService`, `StubGeminiService` — keyword matching |
| Không có API cho LeaveRequest, Consultation | `cmd/http/controllers/` | Không tìm thấy package controller tương ứng |
| Validation phone VN | `pkg/constants/validation.go` | `RegexPhoneNumber: ^(?:\+84\|0)\d{9,10}$` |
| Seed data 25 record/bảng | `seed_data.sql` | Dữ liệu test production-like |
| Trường duyệt Program không có API | `internal/entities/program.go` vs controllers | `ApprovedByID` tồn tại nhưng không có endpoint duyệt |
| Preview lưu in-memory | `internal/services/scheduling/preview_store.go` | Map trong bộ nhớ, mất khi restart |
| Mật khẩu mã hóa bcrypt | `internal/services/user/auth.go` | `bcrypt.CompareHashAndPassword` |
| Token JWT mặc định 24h / 7 ngày | `internal/services/user/auth.go` | Fallback duration nếu config thiếu |

---

*Tài liệu này được tạo tự động từ phân tích mã nguồn. Cần được BA rà soát và xác nhận các mục đánh dấu ⚠️ hoặc tham chiếu đến câu hỏi (Q-xx) trước khi chuyển thành tài liệu chính thức.*
