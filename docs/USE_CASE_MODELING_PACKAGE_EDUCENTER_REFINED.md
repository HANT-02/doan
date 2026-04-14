# GÓI PHÂN TÍCH USE CASE PHỤC VỤ VẼ SƠ ĐỒ HỆ THỐNG EDUCENTER

**Mục đích:** tài liệu chuyên dùng để:

- liệt kê toàn bộ use case hiện trạng;
- phân rã use case tổng thành use case con;
- chuẩn hóa tên use case sang tiếng Việt;
- chỉ rõ chỗ nên dùng `<<include>>`, `<<extend>>`, `generalization`;
- cung cấp đủ chi tiết để vẽ sequence diagram.

**Nguồn baseline:** kế thừa từ các phát hiện trước trong:

- [BA_PACKAGE_EDUCENTER_REVERSE_ENGINEER.md](/Users/hant/golang/doan/docs/BA_PACKAGE_EDUCENTER_REVERSE_ENGINEER.md)
- [USE_CASE_SPEC_PACKAGE_EDUCENTER.md](/Users/hant/golang/doan/docs/USE_CASE_SPEC_PACKAGE_EDUCENTER.md)
- [BA_SYSTEM_ANALYSIS_REPORT_EDUCENTER.md](/Users/hant/golang/doan/docs/BA_SYSTEM_ANALYSIS_REPORT_EDUCENTER.md)

**Nguyên tắc:** không dựng lại từ đầu; chỉ tổ chức lại để phục vụ modeling tốt hơn.

## Quy ước mức độ chắc chắn

- **Confirmed from code**: xác nhận trực tiếp từ route, controller, use case, entity, repository, UI.
- **Strongly inferred from code**: suy luận mạnh từ data model, tên module, UI placeholder, side effect.
- **Assumption / Needs BA validation**: cần xác nhận với BA/PO.

---

# 1. Cách đọc tài liệu này

## 1.1 Mục tiêu của tài liệu

Tài liệu này không đi theo kiểu “mỗi use case một đoạn mô tả dài” như tài liệu đặc tả đầy đủ, mà tập trung trả lời 5 câu hỏi BA cần khi đi vẽ:

1. **Toàn hệ thống hiện có những use case nào?**
2. **Use case nào là use case tổng, use case nào là use case con?**
3. **Khi nào nên vẽ `include`, khi nào nên vẽ `extend`?**
4. **CRUD một module nên tách thành mấy use case?**
5. **Vẽ sequence diagram thì mỗi use case cần các lifeline và message gì?**

## 1.2 Nguyên tắc đặt tên

- Dùng **tiếng Việt hoàn toàn** cho tên use case.
- Một use case chỉ nên đặt theo **mục tiêu nghiệp vụ**, không đặt theo endpoint kỹ thuật.
- Tên use case cấp tổng nên là **danh từ + động từ bao quát**, ví dụ:
  - `Quản lý học viên`
  - `Quản lý giáo viên`
  - `Xếp lịch học`
- Tên use case con nên là **hành động cụ thể**, ví dụ:
  - `Tạo học viên`
  - `Cập nhật học viên`
  - `Xóa học viên`
  - `Xem danh sách học viên`

## 1.3 Quy tắc dùng quan hệ UML

### Dùng `<<include>>` khi:

- một use case **luôn luôn** phải gọi một hành vi con;
- hành vi con được **tái sử dụng** bởi nhiều use case;
- bỏ hành vi con đó thì use case cha **không hoàn thành**.

Ví dụ:
- `Đăng ký tài khoản` `<<include>>` `Tạo mã OTP xác minh`
- `Tạo preview xếp lịch` `<<include>>` `Nạp dữ liệu xếp lịch`
- `Xác nhận preview để tạo buổi học` `<<include>>` `Kiểm tra điều kiện commit`

### Dùng `<<extend>>` khi:

- hành vi chỉ xảy ra ở một nhánh tùy chọn;
- hành vi chỉ kích hoạt khi có điều kiện bổ sung;
- use case chính vẫn hoàn chỉnh mà **không cần** nhánh đó.

Ví dụ:
- `Giải quyết xung đột lịch học` `<<extend>>` `Xem preview xếp lịch`
- `Từ chối tài liệu` `<<extend>>` `Ra quyết định duyệt tài liệu`
- `Tải tài liệu` `<<extend>>` `Xem chi tiết tài liệu`

### Dùng **generalization** khi:

- một use case tổng có nhiều biến thể cùng bản chất;
- các biến thể khác nhau về hành động cuối hoặc luật nghiệp vụ riêng.

Ví dụ:
- `Ra quyết định duyệt tài liệu`
  - kế thừa bởi `Phê duyệt tài liệu`
  - kế thừa bởi `Từ chối tài liệu`

---

# 2. Danh sách toàn bộ use case hiện trạng theo cụm

## 2.1 Nhóm xác thực và tài khoản

| Mã | Tên use case tiếng Việt | Cấp | Cha | Trạng thái | Mức độ chắc chắn |
|---|---|---|---|---|---|
| UC-AT-01 | Quản lý truy cập tài khoản | Tổng | ROOT | Partial | Confirmed from code |
| UC-AT-01.01 | Đăng ký tài khoản | Con | Quản lý truy cập tài khoản | Implemented | Confirmed from code |
| UC-AT-01.02 | Xác minh email bằng OTP | Con | Quản lý truy cập tài khoản | Implemented | Confirmed from code |
| UC-AT-01.03 | Đăng nhập | Con | Quản lý truy cập tài khoản | Implemented | Confirmed from code |
| UC-AT-01.04 | Làm mới phiên đăng nhập | Con | Quản lý truy cập tài khoản | Implemented | Confirmed from code |
| UC-AT-01.05 | Đăng xuất | Con | Quản lý truy cập tài khoản | Implemented | Confirmed from code |
| UC-AT-01.06 | Yêu cầu quên mật khẩu | Con | Quản lý truy cập tài khoản | Implemented | Confirmed from code |
| UC-AT-01.07 | Đặt lại mật khẩu | Con | Quản lý truy cập tài khoản | Implemented | Confirmed from code |
| UC-AT-01.08 | Đổi mật khẩu | Con | Quản lý truy cập tài khoản | Implemented | Confirmed from code |
| UC-AT-01.09 | Xem hồ sơ cá nhân | Con | Quản lý truy cập tài khoản | Partial | Confirmed from code |

## 2.2 Nhóm học viên

| Mã | Tên use case tiếng Việt | Cấp | Cha | Trạng thái | Mức độ chắc chắn |
|---|---|---|---|---|---|
| UC-HV-01 | Quản lý học viên | Tổng | ROOT | Implemented | Confirmed from code |
| UC-HV-01.01 | Tạo học viên | Con | Quản lý học viên | Implemented | Confirmed from code |
| UC-HV-01.02 | Xem danh sách học viên | Con | Quản lý học viên | Implemented | Confirmed from code |
| UC-HV-01.03 | Tìm kiếm học viên | Con | Quản lý học viên | Implemented | Confirmed from code |
| UC-HV-01.04 | Xem chi tiết học viên | Con | Quản lý học viên | Strongly inferred from UI | Strongly inferred from code |
| UC-HV-01.05 | Cập nhật học viên | Con | Quản lý học viên | Implemented | Confirmed from code |
| UC-HV-01.06 | Xóa học viên | Con | Quản lý học viên | Implemented | Confirmed from code |

## 2.3 Nhóm giáo viên

| Mã | Tên use case tiếng Việt | Cấp | Cha | Trạng thái | Mức độ chắc chắn |
|---|---|---|---|---|---|
| UC-GV-01 | Quản lý giáo viên | Tổng | ROOT | Implemented | Confirmed from code |
| UC-GV-01.01 | Tạo giáo viên | Con | Quản lý giáo viên | Implemented | Confirmed from code |
| UC-GV-01.02 | Xem danh sách giáo viên | Con | Quản lý giáo viên | Implemented | Confirmed from code |
| UC-GV-01.03 | Tìm kiếm giáo viên | Con | Quản lý giáo viên | Implemented | Confirmed from code |
| UC-GV-01.04 | Xem chi tiết giáo viên | Con | Quản lý giáo viên | Implemented | Confirmed from code |
| UC-GV-01.05 | Cập nhật giáo viên | Con | Quản lý giáo viên | Implemented | Confirmed from code |
| UC-GV-01.06 | Xóa giáo viên | Con | Quản lý giáo viên | Implemented | Confirmed from code |
| UC-GV-01.07 | Xem lịch dạy giáo viên | Con | Quản lý giáo viên | Implemented | Confirmed from code |
| UC-GV-01.08 | Xem thống kê giờ dạy | Con | Quản lý giáo viên | Implemented | Confirmed from code |

## 2.4 Nhóm khóa học và chương trình

| Mã | Tên use case tiếng Việt | Cấp | Cha | Trạng thái | Mức độ chắc chắn |
|---|---|---|---|---|---|
| UC-KH-01 | Quản lý khóa học | Tổng | ROOT | Implemented | Confirmed from code |
| UC-KH-01.01 | Tạo khóa học | Con | Quản lý khóa học | Implemented | Confirmed from code |
| UC-KH-01.02 | Xem danh sách khóa học | Con | Quản lý khóa học | Implemented | Confirmed from code |
| UC-KH-01.03 | Tìm kiếm khóa học | Con | Quản lý khóa học | Implemented | Confirmed from code |
| UC-KH-01.04 | Xem chi tiết khóa học | Con | Quản lý khóa học | Strongly inferred from UI | Strongly inferred from code |
| UC-KH-01.05 | Cập nhật khóa học | Con | Quản lý khóa học | Implemented | Confirmed from code |
| UC-KH-01.06 | Xóa khóa học | Con | Quản lý khóa học | Implemented | Confirmed from code |
| UC-CT-01 | Quản lý chương trình đào tạo | Tổng | ROOT | Partial | Confirmed from code |
| UC-CT-01.01 | Tạo chương trình đào tạo | Con | Quản lý chương trình đào tạo | Implemented | Confirmed from code |
| UC-CT-01.02 | Xem danh sách chương trình đào tạo | Con | Quản lý chương trình đào tạo | Implemented | Confirmed from code |
| UC-CT-01.03 | Tìm kiếm chương trình đào tạo | Con | Quản lý chương trình đào tạo | Implemented | Confirmed from code |
| UC-CT-01.04 | Xem chi tiết chương trình đào tạo | Con | Quản lý chương trình đào tạo | Implemented | Confirmed from code |
| UC-CT-01.05 | Cập nhật chương trình đào tạo | Con | Quản lý chương trình đào tạo | Implemented | Confirmed from code |
| UC-CT-01.06 | Xóa chương trình đào tạo | Con | Quản lý chương trình đào tạo | Implemented | Confirmed from code |
| UC-CT-01.07 | Liên kết khóa học vào chương trình | Con | Quản lý chương trình đào tạo | Implemented | Confirmed from code |
| UC-CT-01.08 | Gỡ khóa học khỏi chương trình | Con | Quản lý chương trình đào tạo | Implemented | Confirmed from code |
| UC-CT-01.09 | Xuất bản chương trình đào tạo | Con | Quản lý chương trình đào tạo | Missing API | Strongly inferred from code |
| UC-CT-01.10 | Lưu trữ chương trình đào tạo | Con | Quản lý chương trình đào tạo | Missing API | Strongly inferred from code |

## 2.5 Nhóm phòng học và ca học

| Mã | Tên use case tiếng Việt | Cấp | Cha | Trạng thái | Mức độ chắc chắn |
|---|---|---|---|---|---|
| UC-PH-01 | Quản lý phòng học | Tổng | ROOT | Implemented | Confirmed from code |
| UC-PH-01.01 | Tạo phòng học | Con | Quản lý phòng học | Implemented | Confirmed from code |
| UC-PH-01.02 | Xem danh sách phòng học | Con | Quản lý phòng học | Implemented | Confirmed from code |
| UC-PH-01.03 | Tìm kiếm phòng học | Con | Quản lý phòng học | Implemented | Confirmed from code |
| UC-PH-01.04 | Xem chi tiết phòng học | Con | Quản lý phòng học | Strongly inferred from UI | Strongly inferred from code |
| UC-PH-01.05 | Cập nhật phòng học | Con | Quản lý phòng học | Implemented | Confirmed from code |
| UC-PH-01.06 | Xóa phòng học | Con | Quản lý phòng học | Implemented | Confirmed from code |
| UC-CA-01 | Quản lý ca học | Tổng | ROOT | Implemented | Confirmed from code |
| UC-CA-01.01 | Tạo ca học | Con | Quản lý ca học | Implemented | Confirmed from code |
| UC-CA-01.02 | Xem danh sách ca học | Con | Quản lý ca học | Implemented | Confirmed from code |
| UC-CA-01.03 | Tìm kiếm ca học | Con | Quản lý ca học | Strongly inferred from UI | Strongly inferred from code |
| UC-CA-01.04 | Xem chi tiết ca học | Con | Quản lý ca học | Strongly inferred from UI | Strongly inferred from code |
| UC-CA-01.05 | Cập nhật ca học | Con | Quản lý ca học | Implemented | Confirmed from code |
| UC-CA-01.06 | Xóa ca học | Con | Quản lý ca học | Implemented | Confirmed from code |

## 2.6 Nhóm lớp học và ghi danh

| Mã | Tên use case tiếng Việt | Cấp | Cha | Trạng thái | Mức độ chắc chắn |
|---|---|---|---|---|---|
| UC-LH-01 | Quản lý lớp học | Tổng | ROOT | Implemented | Confirmed from code |
| UC-LH-01.01 | Tạo lớp học | Con | Quản lý lớp học | Implemented | Confirmed from code |
| UC-LH-01.02 | Xem danh sách lớp học | Con | Quản lý lớp học | Implemented | Confirmed from code |
| UC-LH-01.03 | Tìm kiếm lớp học | Con | Quản lý lớp học | Implemented | Confirmed from code |
| UC-LH-01.04 | Lọc lớp học theo trạng thái | Con | Quản lý lớp học | Implemented | Confirmed from code |
| UC-LH-01.05 | Xem chi tiết lớp học | Con | Quản lý lớp học | Implemented | Confirmed from code |
| UC-LH-01.06 | Cập nhật lớp học | Con | Quản lý lớp học | Implemented | Confirmed from code |
| UC-LH-01.07 | Xóa lớp học | Con | Quản lý lớp học | Implemented | Confirmed from code |
| UC-LH-01.08 | Phân công giáo viên cho lớp | Con | Quản lý lớp học | Implemented | Confirmed from code |
| UC-LH-02 | Quản lý ghi danh lớp học | Tổng | ROOT | Partial | Confirmed from code |
| UC-LH-02.01 | Xem danh sách học viên trong lớp | Con | Quản lý ghi danh lớp học | Implemented | Confirmed from code |
| UC-LH-02.02 | Tìm kiếm học viên khả dụng để ghi danh | Con | Quản lý ghi danh lớp học | Implemented | Confirmed from code |
| UC-LH-02.03 | Ghi danh học viên vào lớp | Con | Quản lý ghi danh lớp học | Implemented | Confirmed from code |
| UC-LH-02.04 | Rút học viên khỏi lớp | Con | Quản lý ghi danh lớp học | Implemented | Confirmed from code |
| UC-LH-03 | Quản lý lịch tuần lớp học | Tổng | ROOT | Missing API | Strongly inferred from code |
| UC-LH-03.01 | Tạo lịch tuần cho lớp | Con | Quản lý lịch tuần lớp học | Missing API | Strongly inferred from code |
| UC-LH-03.02 | Cập nhật lịch tuần cho lớp | Con | Quản lý lịch tuần lớp học | Missing API | Strongly inferred from code |
| UC-LH-03.03 | Xóa lịch tuần của lớp | Con | Quản lý lịch tuần lớp học | Missing API | Strongly inferred from code |

## 2.7 Nhóm xếp lịch và buổi học

| Mã | Tên use case tiếng Việt | Cấp | Cha | Trạng thái | Mức độ chắc chắn |
|---|---|---|---|---|---|
| UC-XL-01 | Xếp lịch học | Tổng | ROOT | Partial | Confirmed from code |
| UC-XL-01.01 | Tạo preview xếp lịch | Con | Xếp lịch học | Implemented | Confirmed from code |
| UC-XL-01.02 | Xem preview xếp lịch | Con | Xếp lịch học | Implemented | Confirmed from code |
| UC-XL-01.03 | Xem xung đột xếp lịch | Con | Xếp lịch học | Implemented | Confirmed from code |
| UC-XL-01.04 | Giải quyết xung đột xếp lịch | Con | Xếp lịch học | Partial | Strongly inferred from code |
| UC-XL-01.05 | Xác nhận preview để tạo buổi học | Con | Xếp lịch học | Implemented | Confirmed from code |
| UC-XL-01.06 | Benchmark solver xếp lịch | Con | Xếp lịch học | Partial | Confirmed from code |
| UC-BH-01 | Quản lý buổi học | Tổng | ROOT | Partial | Strongly inferred from code |
| UC-BH-01.01 | Xem danh sách buổi học | Con | Quản lý buổi học | Missing API/UI đầy đủ | Strongly inferred from code |
| UC-BH-01.02 | Xem chi tiết buổi học | Con | Quản lý buổi học | Missing API/UI đầy đủ | Strongly inferred from code |

## 2.8 Nhóm học vụ sau buổi học

| Mã | Tên use case tiếng Việt | Cấp | Cha | Trạng thái | Mức độ chắc chắn |
|---|---|---|---|---|---|
| UC-HVU-01 | Quản lý điểm danh | Tổng | ROOT | Missing API | Strongly inferred from code |
| UC-HVU-01.01 | Chấm điểm danh buổi học | Con | Quản lý điểm danh | Missing API | Strongly inferred from code |
| UC-HVU-01.02 | Cập nhật điểm danh buổi học | Con | Quản lý điểm danh | Missing API | Strongly inferred from code |
| UC-HVU-01.03 | Xem lịch sử điểm danh | Con | Quản lý điểm danh | Missing API/UI | Strongly inferred from code |
| UC-HVU-02 | Quản lý tổng kết buổi học | Tổng | ROOT | Missing API | Strongly inferred from code |
| UC-HVU-02.01 | Tạo tổng kết buổi học | Con | Quản lý tổng kết buổi học | Missing API | Strongly inferred from code |
| UC-HVU-02.02 | Cập nhật tổng kết buổi học | Con | Quản lý tổng kết buổi học | Missing API | Strongly inferred from code |
| UC-HVU-02.03 | Xem tổng kết buổi học | Con | Quản lý tổng kết buổi học | Missing API/UI | Strongly inferred from code |
| UC-HVU-03 | Quản lý kết quả học tập | Tổng | ROOT | Missing API | Strongly inferred from code |
| UC-HVU-03.01 | Ghi nhận kết quả học tập | Con | Quản lý kết quả học tập | Missing API | Strongly inferred from code |
| UC-HVU-03.02 | Cập nhật kết quả học tập | Con | Quản lý kết quả học tập | Missing API | Strongly inferred from code |
| UC-HVU-03.03 | Xem kết quả học tập | Con | Quản lý kết quả học tập | Missing API/UI | Strongly inferred from code |
| UC-HVU-04 | Quản lý đơn xin phép | Tổng | ROOT | Missing API | Strongly inferred from code |
| UC-HVU-04.01 | Tạo đơn xin phép | Con | Quản lý đơn xin phép | Missing API | Strongly inferred from code |
| UC-HVU-04.02 | Xem danh sách đơn xin phép | Con | Quản lý đơn xin phép | Missing API | Strongly inferred from code |
| UC-HVU-04.03 | Duyệt đơn xin phép | Con | Quản lý đơn xin phép | Missing API | Strongly inferred from code |
| UC-HVU-04.04 | Từ chối đơn xin phép | Con | Quản lý đơn xin phép | Missing API | Strongly inferred from code |

## 2.9 Nhóm tài liệu giảng dạy và kiểm duyệt

| Mã | Tên use case tiếng Việt | Cấp | Cha | Trạng thái | Mức độ chắc chắn |
|---|---|---|---|---|---|
| UC-TL-01 | Quản lý tài liệu giảng dạy | Tổng | ROOT | Partial | Confirmed from code |
| UC-TL-01.01 | Tải lên tài liệu giảng dạy | Con | Quản lý tài liệu giảng dạy | Implemented | Confirmed from code |
| UC-TL-01.02 | Xem danh sách tài liệu | Con | Quản lý tài liệu giảng dạy | Implemented | Confirmed from code |
| UC-TL-01.03 | Xem chi tiết tài liệu | Con | Quản lý tài liệu giảng dạy | Implemented | Confirmed from code |
| UC-TL-01.04 | Tải xuống tài liệu | Con | Quản lý tài liệu giảng dạy | Implemented | Confirmed from code |
| UC-TL-01.05 | Xem hàng chờ tài liệu cần duyệt | Con | Quản lý tài liệu giảng dạy | Implemented | Strongly inferred from code |
| UC-TL-01.06 | Xem lịch sử audit tài liệu | Con | Quản lý tài liệu giảng dạy | Partial | Confirmed from code |
| UC-TL-01.07 | Ra quyết định duyệt tài liệu | Con tổng | Quản lý tài liệu giảng dạy | Implemented | Confirmed from code |
| UC-TL-01.07.01 | Phê duyệt tài liệu | Con chuyên biệt | Ra quyết định duyệt tài liệu | Implemented | Confirmed from code |
| UC-TL-01.07.02 | Từ chối tài liệu | Con chuyên biệt | Ra quyết định duyệt tài liệu | Implemented | Confirmed from code |

## 2.10 Nhóm tư vấn / lead intake

| Mã | Tên use case tiếng Việt | Cấp | Cha | Trạng thái | Mức độ chắc chắn |
|---|---|---|---|---|---|
| UC-TV-01 | Quản lý tư vấn tuyển sinh | Tổng | ROOT | Missing API | Strongly inferred from code |
| UC-TV-01.01 | Tiếp nhận yêu cầu tư vấn | Con | Quản lý tư vấn tuyển sinh | Missing API | Strongly inferred from code |
| UC-TV-01.02 | Xem danh sách lead tư vấn | Con | Quản lý tư vấn tuyển sinh | Missing API | Strongly inferred from code |
| UC-TV-01.03 | Cập nhật trạng thái lead tư vấn | Con | Quản lý tư vấn tuyển sinh | Missing API | Strongly inferred from code |

---

# 3. Phân rã CRUD chi tiết để vẽ use case

## 3.1 Nguyên tắc phân rã CRUD

Khi vẽ use case cho các module CRUD, **không nên vẽ một oval duy nhất tên “CRUD X”**. Nên tách thành:

1. **Quản lý X**: use case tổng
2. **Tạo X**
3. **Xem danh sách X**
4. **Tìm kiếm / lọc X**
5. **Xem chi tiết X**
6. **Cập nhật X**
7. **Xóa X**

Nếu module có hành vi riêng ngoài CRUD, tách thêm:

- `Gán ...`
- `Liên kết ...`
- `Xuất bản ...`
- `Duyệt ...`
- `Commit ...`
- `Benchmark ...`

## 3.2 Phân rã CRUD học viên

### Use case tổng

- `Quản lý học viên`

### Use case con

- `Tạo học viên`
- `Xem danh sách học viên`
- `Tìm kiếm học viên`
- `Xem chi tiết học viên`
- `Cập nhật học viên`
- `Xóa học viên`

### Gợi ý vẽ quan hệ

- `Quản lý học viên` là use case tổng, **không bắt buộc** nối `include` tới toàn bộ use case con nếu sơ đồ của bạn là sơ đồ phân rã chức năng.
- Nếu vẽ **sơ đồ tổng hợp**, có thể nối actor `Quản trị viên` tới `Quản lý học viên`, rồi trong **sơ đồ con** mới tách 6 use case trên.
- `Tìm kiếm học viên` có thể `<<extend>>` `Xem danh sách học viên`.
- `Xem chi tiết học viên` có thể `<<extend>>` `Xem danh sách học viên`.

## 3.3 Phân rã CRUD giáo viên

- `Quản lý giáo viên`
  - `Tạo giáo viên`
  - `Xem danh sách giáo viên`
  - `Tìm kiếm giáo viên`
  - `Xem chi tiết giáo viên`
  - `Cập nhật giáo viên`
  - `Xóa giáo viên`
  - `Xem lịch dạy giáo viên`
  - `Xem thống kê giờ dạy`

### Gợi ý include/extend

- `Xem lịch dạy giáo viên` `<<extend>>` `Xem chi tiết giáo viên`.
- `Xem thống kê giờ dạy` `<<extend>>` `Xem chi tiết giáo viên`.
- `Tìm kiếm giáo viên` `<<extend>>` `Xem danh sách giáo viên`.

## 3.4 Phân rã CRUD khóa học

- `Quản lý khóa học`
  - `Tạo khóa học`
  - `Xem danh sách khóa học`
  - `Tìm kiếm khóa học`
  - `Xem chi tiết khóa học`
  - `Cập nhật khóa học`
  - `Xóa khóa học`

## 3.5 Phân rã CRUD chương trình đào tạo

- `Quản lý chương trình đào tạo`
  - `Tạo chương trình đào tạo`
  - `Xem danh sách chương trình đào tạo`
  - `Tìm kiếm chương trình đào tạo`
  - `Xem chi tiết chương trình đào tạo`
  - `Cập nhật chương trình đào tạo`
  - `Xóa chương trình đào tạo`
  - `Liên kết khóa học vào chương trình`
  - `Gỡ khóa học khỏi chương trình`
  - `Xuất bản chương trình đào tạo`
  - `Lưu trữ chương trình đào tạo`

### Gợi ý include/extend

- `Liên kết khóa học vào chương trình` `<<extend>>` `Xem chi tiết chương trình đào tạo`.
- `Gỡ khóa học khỏi chương trình` `<<extend>>` `Xem chi tiết chương trình đào tạo`.
- `Xuất bản chương trình đào tạo` và `Lưu trữ chương trình đào tạo` là **hành vi trạng thái**, nên vẽ riêng thay vì gộp vào `Cập nhật chương trình đào tạo`.

## 3.6 Phân rã CRUD lớp học

- `Quản lý lớp học`
  - `Tạo lớp học`
  - `Xem danh sách lớp học`
  - `Tìm kiếm lớp học`
  - `Lọc lớp học theo trạng thái`
  - `Xem chi tiết lớp học`
  - `Cập nhật lớp học`
  - `Xóa lớp học`
  - `Phân công giáo viên cho lớp`

### Use case liên quan nhưng nên tách thành nhóm riêng

- `Quản lý ghi danh lớp học`
  - `Xem danh sách học viên trong lớp`
  - `Tìm kiếm học viên khả dụng để ghi danh`
  - `Ghi danh học viên vào lớp`
  - `Rút học viên khỏi lớp`

- `Quản lý lịch tuần lớp học`
  - `Tạo lịch tuần cho lớp`
  - `Cập nhật lịch tuần cho lớp`
  - `Xóa lịch tuần của lớp`

### Gợi ý include/extend

- `Phân công giáo viên cho lớp` `<<extend>>` `Xem chi tiết lớp học`.
- `Ghi danh học viên vào lớp` `<<extend>>` `Xem chi tiết lớp học`.
- `Rút học viên khỏi lớp` `<<extend>>` `Xem chi tiết lớp học`.
- `Tạo lịch tuần cho lớp` `<<extend>>` `Xem chi tiết lớp học`.

## 3.7 Phân rã CRUD phòng học và ca học

### Quản lý phòng học

- `Tạo phòng học`
- `Xem danh sách phòng học`
- `Tìm kiếm phòng học`
- `Xem chi tiết phòng học`
- `Cập nhật phòng học`
- `Xóa phòng học`

### Quản lý ca học

- `Tạo ca học`
- `Xem danh sách ca học`
- `Tìm kiếm ca học`
- `Xem chi tiết ca học`
- `Cập nhật ca học`
- `Xóa ca học`

## 3.8 Phân rã nhóm tài liệu giảng dạy

- `Quản lý tài liệu giảng dạy`
  - `Tải lên tài liệu giảng dạy`
  - `Xem danh sách tài liệu`
  - `Xem chi tiết tài liệu`
  - `Tải xuống tài liệu`
  - `Xem hàng chờ tài liệu cần duyệt`
  - `Xem lịch sử audit tài liệu`
  - `Ra quyết định duyệt tài liệu`
    - `Phê duyệt tài liệu`
    - `Từ chối tài liệu`

### Gợi ý modeling

- `Phê duyệt tài liệu` và `Từ chối tài liệu` nên là **generalization** của `Ra quyết định duyệt tài liệu`.
- `Xem lịch sử audit tài liệu` `<<extend>>` `Xem chi tiết tài liệu`.
- `Tải xuống tài liệu` `<<extend>>` `Xem chi tiết tài liệu`.

---

# 4. Bảng quan hệ `include`, `extend`, `generalization`

## 4.1 Bảng `include` đề xuất

| UC nguồn | UC đích | Loại | Lý do dùng | Mức độ chắc chắn |
|---|---|---|---|---|
| Đăng ký tài khoản | Kiểm tra trùng email | `<<include>>` | luôn phải kiểm tra trước khi tạo user | Confirmed from code |
| Đăng ký tài khoản | Giải mã và băm mật khẩu | `<<include>>` | bắt buộc để tạo tài khoản | Confirmed from code |
| Đăng ký tài khoản | Tạo mã OTP xác minh | `<<include>>` | luôn tạo OTP khi đăng ký | Confirmed from code |
| Đăng ký tài khoản | Gửi email OTP | `<<include>>` | là bước nghiệp vụ mặc định sau khi tạo user | Confirmed from code |
| Xác minh email bằng OTP | Kiểm tra hiệu lực OTP | `<<include>>` | luôn cần trước khi activate user | Confirmed from code |
| Đăng nhập | Kiểm tra trạng thái tài khoản | `<<include>>` | account inactive không được login | Confirmed from code |
| Yêu cầu quên mật khẩu | Tạo yêu cầu đặt lại mật khẩu | `<<include>>` | hành vi cốt lõi | Confirmed from code |
| Yêu cầu quên mật khẩu | Gửi email đặt lại mật khẩu | `<<include>>` | hành vi bắt buộc của flow | Confirmed from code |
| Tạo preview xếp lịch | Nạp dữ liệu xếp lịch | `<<include>>` | phải load class/shift/room/teacher | Confirmed from code |
| Tạo preview xếp lịch | Xây dựng bài toán xếp lịch | `<<include>>` | là phần lõi của use case | Confirmed from code |
| Tạo preview xếp lịch | Chạy solver xếp lịch | `<<include>>` | preview không có nếu không chạy solver | Confirmed from code |
| Tạo preview xếp lịch | Tổng hợp xung đột xếp lịch | `<<include>>` | preview luôn có summary/conflicts | Confirmed from code |
| Xác nhận preview để tạo buổi học | Kiểm tra điều kiện commit | `<<include>>` | luôn cần trước khi tạo lesson | Confirmed from code |
| Xác nhận preview để tạo buổi học | Kiểm tra trùng buổi học | `<<include>>` | luôn cần trước khi insert lessons | Confirmed from code |
| Xác nhận preview để tạo buổi học | Tạo buổi học từ preview | `<<include>>` | hành vi cốt lõi | Confirmed from code |
| Tải lên tài liệu giảng dạy | Kiểm tra định dạng và dung lượng file | `<<include>>` | luôn bắt buộc | Confirmed from code |
| Tải lên tài liệu giảng dạy | Lưu file và metadata tài liệu | `<<include>>` | luôn bắt buộc | Confirmed from code |
| Tải lên tài liệu giảng dạy | Quét OCR tài liệu | `<<include>>` | flow hiện tại luôn chạy | Confirmed from code |
| Tải lên tài liệu giảng dạy | Gắn nhãn AI cho tài liệu | `<<include>>` | flow hiện tại luôn chạy | Confirmed from code |
| Tải lên tài liệu giảng dạy | Ghi nhật ký audit tài liệu | `<<include>>` | luôn tạo audit log | Confirmed from code |

## 4.2 Bảng `extend` đề xuất

| UC nguồn | UC đích | Loại | Khi nào dùng | Mức độ chắc chắn |
|---|---|---|---|---|
| Tìm kiếm học viên | Xem danh sách học viên | `<<extend>>` | chỉ xảy ra khi người dùng nhập tiêu chí tìm kiếm | Confirmed from code |
| Xem chi tiết học viên | Xem danh sách học viên | `<<extend>>` | từ danh sách click vào một dòng cụ thể | Strongly inferred from UI |
| Tìm kiếm giáo viên | Xem danh sách giáo viên | `<<extend>>` | lọc trên danh sách | Confirmed from code |
| Xem chi tiết giáo viên | Xem danh sách giáo viên | `<<extend>>` | mở detail từ danh sách | Confirmed from code |
| Xem lịch dạy giáo viên | Xem chi tiết giáo viên | `<<extend>>` | chỉ khi muốn đi sâu vào một giáo viên | Strongly inferred from UI |
| Xem thống kê giờ dạy | Xem chi tiết giáo viên | `<<extend>>` | chỉ khi xem thống kê | Strongly inferred from UI |
| Liên kết khóa học vào chương trình | Xem chi tiết chương trình đào tạo | `<<extend>>` | thao tác từ màn chi tiết program | Confirmed from code |
| Gỡ khóa học khỏi chương trình | Xem chi tiết chương trình đào tạo | `<<extend>>` | thao tác từ màn chi tiết program | Confirmed from code |
| Ghi danh học viên vào lớp | Xem chi tiết lớp học | `<<extend>>` | thao tác từ màn chi tiết lớp | Confirmed from code |
| Rút học viên khỏi lớp | Xem chi tiết lớp học | `<<extend>>` | thao tác từ roster của lớp | Confirmed from code |
| Phân công giáo viên cho lớp | Xem chi tiết lớp học | `<<extend>>` | thao tác từ chi tiết lớp | Confirmed from code |
| Tạo lịch tuần cho lớp | Xem chi tiết lớp học | `<<extend>>` | nếu hệ thống bổ sung màn chi tiết lịch tuần | Strongly inferred from code |
| Giải quyết xung đột xếp lịch | Xem preview xếp lịch | `<<extend>>` | chỉ xảy ra khi preview có conflict | Strongly inferred from code |
| Tải xuống tài liệu | Xem chi tiết tài liệu | `<<extend>>` | tùy chọn khi cần lấy file | Confirmed from code |
| Xem lịch sử audit tài liệu | Xem chi tiết tài liệu | `<<extend>>` | tùy chọn khi cần xem log audit | Confirmed from code |
| Từ chối tài liệu | Ra quyết định duyệt tài liệu | `<<extend>>` hoặc generalization | nếu bạn muốn nhấn mạnh nhánh tùy điều kiện reject | Confirmed from code |

## 4.3 Bảng `generalization` đề xuất

| UC cha | UC con | Loại | Lý do |
|---|---|---|---|
| Ra quyết định duyệt tài liệu | Phê duyệt tài liệu | Generalization | cùng mục tiêu nghiệp vụ, khác kết quả cuối |
| Ra quyết định duyệt tài liệu | Từ chối tài liệu | Generalization | cùng mục tiêu nghiệp vụ, khác kết quả cuối |
| Quản lý truy cập tài khoản | Đăng ký tài khoản / Đăng nhập / Quên mật khẩu / Đặt lại mật khẩu / Đổi mật khẩu | Không khuyến nghị generalization | nên tách theo cụm hơn là kế thừa |
| Quản lý học viên | Tạo / Xem / Cập nhật / Xóa học viên | Không khuyến nghị generalization | đây là phân rã chức năng, không phải kế thừa hành vi |

## 4.4 Cách chọn `extend` hay `generalization` cho material review

Nếu bạn muốn sơ đồ **dễ đọc cho BA**:

- dùng `Ra quyết định duyệt tài liệu` làm UC cha
- `Phê duyệt tài liệu`
- `Từ chối tài liệu`
- nối **generalization**

Nếu bạn muốn nhấn mạnh điều kiện “từ chối chỉ xảy ra khi reviewer quyết định reject”:

- `Từ chối tài liệu` có thể vẽ `<<extend>>` `Ra quyết định duyệt tài liệu`

**Khuyến nghị thực tế:** với EduCenter, nên dùng **generalization** cho cặp `Phê duyệt` / `Từ chối`, vì 2 nhánh này ngang hàng và rõ nghĩa hơn.

---

# 5. Ma trận actor - use case để vẽ sơ đồ

## 5.1 Ma trận rút gọn theo sơ đồ tổng

| Actor | Use case tổng nên nối trực tiếp |
|---|---|
| Người dùng chưa đăng nhập | Quản lý truy cập tài khoản |
| Quản trị viên | Quản lý học viên, Quản lý giáo viên, Quản lý khóa học, Quản lý chương trình đào tạo, Quản lý phòng học, Quản lý ca học, Quản lý lớp học, Quản lý ghi danh lớp học, Quản lý lịch tuần lớp học, Xếp lịch học, Quản lý tài liệu giảng dạy |
| Giáo viên | Quản lý truy cập tài khoản, Quản lý giáo viên, Quản lý tài liệu giảng dạy, Quản lý điểm danh, Quản lý tổng kết buổi học, Quản lý kết quả học tập |
| Học viên | Quản lý truy cập tài khoản, Quản lý đơn xin phép, Quản lý kết quả học tập |
| Reviewer / Compliance | Quản lý tài liệu giảng dạy |
| Scheduling Engine | Tạo preview xếp lịch, Xác nhận preview để tạo buổi học, Benchmark solver xếp lịch |
| SMTP | Đăng ký tài khoản, Yêu cầu quên mật khẩu |
| OCR / AI stub | Tải lên tài liệu giảng dạy |

## 5.2 Ma trận chi tiết theo cụm

| Cụm sơ đồ | Actor chính | Use case nên xuất hiện |
|---|---|---|
| Sơ đồ xác thực | Guest, User, SMTP | Đăng ký tài khoản, Xác minh email bằng OTP, Đăng nhập, Làm mới phiên đăng nhập, Yêu cầu quên mật khẩu, Đặt lại mật khẩu, Đổi mật khẩu, Xem hồ sơ cá nhân |
| Sơ đồ học viên | Admin | Tạo học viên, Xem danh sách học viên, Tìm kiếm học viên, Xem chi tiết học viên, Cập nhật học viên, Xóa học viên |
| Sơ đồ giáo viên | Admin, Teacher | Tạo giáo viên, Xem danh sách giáo viên, Xem chi tiết giáo viên, Cập nhật giáo viên, Xóa giáo viên, Xem lịch dạy giáo viên, Xem thống kê giờ dạy |
| Sơ đồ khóa học/chương trình | Admin | Tạo khóa học, Cập nhật khóa học, Tạo chương trình đào tạo, Cập nhật chương trình đào tạo, Liên kết khóa học vào chương trình, Gỡ khóa học khỏi chương trình |
| Sơ đồ lớp và ghi danh | Admin | Tạo lớp học, Cập nhật lớp học, Xem chi tiết lớp học, Phân công giáo viên cho lớp, Ghi danh học viên vào lớp, Rút học viên khỏi lớp, Quản lý lịch tuần lớp học |
| Sơ đồ xếp lịch | Admin, Scheduling Engine | Tạo preview xếp lịch, Xem preview xếp lịch, Xem xung đột xếp lịch, Giải quyết xung đột xếp lịch, Xác nhận preview để tạo buổi học, Benchmark solver xếp lịch |
| Sơ đồ học vụ sau buổi học | Teacher, Student, Admin | Chấm điểm danh buổi học, Tạo tổng kết buổi học, Ghi nhận kết quả học tập, Tạo đơn xin phép, Duyệt đơn xin phép, Từ chối đơn xin phép |
| Sơ đồ tài liệu | Teacher, Reviewer, OCR/AI | Tải lên tài liệu giảng dạy, Xem chi tiết tài liệu, Xem lịch sử audit tài liệu, Phê duyệt tài liệu, Từ chối tài liệu, Tải xuống tài liệu |

---

# 6. Danh sách use case nên gộp và không nên gộp

## 6.1 Những use case đang bị gộp trong tài liệu cũ và nên tách

| Cụm đang gộp | Nên tách thành |
|---|---|
| Quản lý học viên | Tạo / Xem danh sách / Tìm kiếm / Xem chi tiết / Cập nhật / Xóa |
| Quản lý giáo viên | Tạo / Xem danh sách / Tìm kiếm / Xem chi tiết / Cập nhật / Xóa / Xem lịch dạy / Xem thống kê |
| Quản lý khóa học | Tạo / Xem danh sách / Tìm kiếm / Xem chi tiết / Cập nhật / Xóa |
| Quản lý chương trình | Tạo / Xem danh sách / Tìm kiếm / Xem chi tiết / Cập nhật / Xóa / Liên kết khóa học / Gỡ khóa học / Xuất bản / Lưu trữ |
| Quản lý lớp học | Tạo / Xem danh sách / Tìm kiếm / Lọc trạng thái / Xem chi tiết / Cập nhật / Xóa / Phân công giáo viên |
| Quản lý ghi danh | Xem roster / Tìm học viên khả dụng / Ghi danh / Rút học viên |
| Xếp lịch học | Tạo preview / Xem preview / Xem conflict / Giải quyết conflict / Commit / Benchmark |
| Quản lý tài liệu | Upload / Xem danh sách / Xem chi tiết / Tải xuống / Xem audit log / Phê duyệt / Từ chối |

## 6.2 Những use case có thể giữ ở cấp tổng

Chỉ nên giữ ở cấp tổng trong sơ đồ overview:

- Quản lý truy cập tài khoản
- Quản lý học viên
- Quản lý giáo viên
- Quản lý khóa học
- Quản lý chương trình đào tạo
- Quản lý phòng học
- Quản lý ca học
- Quản lý lớp học
- Quản lý ghi danh lớp học
- Xếp lịch học
- Quản lý tài liệu giảng dạy
- Quản lý học vụ sau buổi học

---

# 7. Đặc tả sequence-ready cho các use case trọng yếu

Phần này viết theo cách để bạn có thể cầm lên và vẽ sequence ngay.

## 7.1 Đăng ký tài khoản

### Tên use case

- `Đăng ký tài khoản`

### Lifeline nên có

1. `Người dùng chưa đăng nhập`
2. `Màn hình đăng ký`
3. `Bộ điều khiển xác thực`
4. `Use case đăng ký tài khoản`
5. `Kho dữ liệu người dùng`
6. `Kho dữ liệu OTP`
7. `Dịch vụ mã hóa mật khẩu`
8. `Dịch vụ gửi email`

### Thông điệp chính để vẽ sequence

1. Người dùng -> Màn hình đăng ký: nhập thông tin đăng ký
2. Màn hình đăng ký -> Bộ điều khiển xác thực: gửi yêu cầu đăng ký
3. Bộ điều khiển xác thực -> Use case đăng ký tài khoản: gọi xử lý đăng ký
4. Use case đăng ký tài khoản -> Kho dữ liệu người dùng: kiểm tra email đã tồn tại
5. Use case đăng ký tài khoản -> Dịch vụ mã hóa mật khẩu: giải mã và băm mật khẩu
6. Use case đăng ký tài khoản -> Kho dữ liệu người dùng: tạo user mới
7. Use case đăng ký tài khoản -> Kho dữ liệu OTP: tạo OTP xác minh
8. Use case đăng ký tài khoản -> Dịch vụ gửi email: gửi email OTP
9. Use case đăng ký tài khoản -> Bộ điều khiển xác thực: trả `user_id`
10. Bộ điều khiển xác thực -> Màn hình đăng ký: phản hồi thành công

### Nhánh thay thế / lỗi nên vẽ

- Email đã tồn tại
- Lỗi giải mã mật khẩu
- Lỗi tạo user / OTP
- Lỗi gửi email sau khi đã commit transaction

## 7.2 Xác minh email bằng OTP

### Lifeline

1. Người dùng
2. Màn hình xác minh OTP
3. Bộ điều khiển xác thực
4. Use case xác minh OTP
5. Kho dữ liệu OTP
6. Kho dữ liệu người dùng

### Thông điệp chính

1. Người dùng gửi OTP
2. Controller gọi use case xác minh
3. Use case đọc OTP active theo `user_id`
4. Use case kiểm tra hạn OTP và trạng thái used
5. Use case so khớp OTP
6. Use case cập nhật OTP used
7. Use case cập nhật user active
8. Trả kết quả thành công

### Nhánh lỗi

- không tìm thấy OTP
- OTP hết hạn
- OTP sai
- lỗi transaction khi activate user

## 7.3 Đăng nhập

### Lifeline

1. Người dùng
2. Màn hình đăng nhập
3. Bộ điều khiển xác thực
4. Use case đăng nhập
5. Kho dữ liệu người dùng
6. Dịch vụ JWT

### Thông điệp chính

1. Người dùng gửi email + mật khẩu
2. Controller gọi use case đăng nhập
3. Use case tìm user theo email
4. Use case kiểm tra mật khẩu
5. Use case kiểm tra `is_active`
6. Use case gọi dịch vụ JWT tạo access/refresh token
7. Trả token + thông tin user

### Nhánh lỗi

- user không tồn tại
- mật khẩu sai
- account inactive

## 7.4 Tạo học viên

### Lifeline

1. Quản trị viên
2. Màn hình quản lý học viên
3. Bộ điều khiển học viên
4. Use case tạo học viên
5. Kho dữ liệu học viên

### Thông điệp chính

1. Admin mở form tạo học viên
2. Nhập thông tin học viên
3. UI gửi request tạo
4. Controller bind request và gọi use case
5. Use case tạo student
6. Repo ghi DB
7. Trả học viên mới tạo

### Nhánh lỗi

- dữ liệu không hợp lệ
- mã học viên trùng
- lỗi DB

## 7.5 Tạo giáo viên

### Lifeline

1. Quản trị viên
2. Màn hình quản lý giáo viên
3. Bộ điều khiển giáo viên
4. Use case tạo giáo viên
5. Kho dữ liệu giáo viên

### Thông điệp chính

1. Admin nhập thông tin giáo viên
2. Controller bind request
3. Use case kiểm tra code/email trùng
4. Use case tạo teacher
5. Repo ghi DB
6. Trả kết quả

### Nhánh lỗi

- code trùng
- email trùng
- lỗi DB

## 7.6 Tạo chương trình đào tạo và liên kết khóa học

### Lifeline

1. Quản trị viên
2. Màn hình chương trình đào tạo
3. Bộ điều khiển chương trình
4. Use case tạo chương trình
5. Kho dữ liệu chương trình
6. Use case liên kết khóa học vào chương trình
7. Kho dữ liệu mapping chương trình-khóa học

### Sequence phần 1: tạo chương trình

1. Admin nhập thông tin chương trình
2. Controller gọi use case tạo chương trình
3. Repo tạo `Program`
4. Trả `program_id`

### Sequence phần 2: liên kết khóa học

1. Admin mở chi tiết chương trình
2. Chọn một hoặc nhiều khóa học
3. Controller gọi use case add courses
4. Use case kiểm tra program tồn tại
5. Use case kiểm tra course tồn tại
6. Repo tạo `ProgramCourse`
7. Trả danh sách cập nhật

### Nhánh lỗi

- chương trình không tồn tại
- khóa học không tồn tại
- mapping trùng

## 7.7 Tạo lớp học

### Lifeline

1. Quản trị viên
2. Màn hình quản lý lớp học
3. Bộ điều khiển lớp học
4. Use case tạo lớp học
5. Kho dữ liệu lớp học

### Thông điệp chính

1. Admin nhập mã lớp, tên, sĩ số, ngày bắt đầu, course/program/teacher nếu có
2. UI gửi request tạo lớp
3. Controller bind request
4. Use case validate input
5. Repo tạo `Class`
6. Trả lớp mới tạo

### Nhánh lỗi

- mã lớp trùng
- `max_students` không hợp lệ
- course/program/teacher không tồn tại nếu có validate

## 7.8 Ghi danh học viên vào lớp

### Lifeline

1. Quản trị viên
2. Màn hình chi tiết lớp học
3. Bộ điều khiển lớp học
4. Use case ghi danh học viên
5. Kho dữ liệu lớp học
6. Kho dữ liệu học viên
7. Kho dữ liệu ghi danh

### Thông điệp chính

1. Admin mở roster lớp
2. Tìm học viên khả dụng
3. Chọn nhiều học viên
4. UI gửi `student_ids`
5. Controller gọi use case ghi danh
6. Use case tải class
7. Use case tính giới hạn sĩ số
8. Use case tạo các bản ghi `Enrollment`
9. Trả roster cập nhật

### Nhánh lỗi

- lớp không tồn tại
- học viên không tồn tại
- vượt sĩ số
- lỗi DB

### Ghi chú modeling

- nếu muốn vẽ đủ chi tiết, nên tách `Tìm kiếm học viên khả dụng để ghi danh` thành một use case riêng ở use case diagram, và một sequence riêng nếu UX search là trọng tâm.

## 7.9 Phân công giáo viên cho lớp

### Lifeline

1. Quản trị viên
2. Màn hình chi tiết lớp học
3. Bộ điều khiển lớp học
4. Use case phân công giáo viên
5. Kho dữ liệu lớp học
6. Kho dữ liệu giáo viên

### Thông điệp chính

1. Admin chọn giáo viên
2. UI gửi `teacher_id`
3. Controller gọi use case
4. Use case kiểm tra class và teacher tồn tại
5. Repo cập nhật `class.teacher_id`
6. Trả class mới

## 7.10 Tạo preview xếp lịch

### Lifeline

1. Quản trị viên
2. Màn hình xếp lịch
3. Bộ điều khiển xếp lịch
4. Use case tạo preview xếp lịch
5. Kho dữ liệu lớp học
6. Kho dữ liệu phòng học
7. Kho dữ liệu ca học
8. Solver service
9. Preview store

### Thông điệp chính

1. Admin nhập date range và bộ lọc class/teacher/room
2. UI gọi API preview
3. Controller bind request và validate date range
4. Use case tải classes `OPEN` và preload teacher/course/class_schedule
5. Use case tải rooms và shifts active
6. Use case xây `SchedulingInput`
7. Use case gọi solver
8. Solver trả assignments/conflicts/summary
9. Use case lưu preview vào preview store
10. Trả preview response cho UI

### Nhánh lỗi

- lớp không đủ dữ liệu
- không có shift active
- class không có class_schedule
- room thiếu hoặc không đủ sức chứa
- solver fail

## 7.11 Xác nhận preview để tạo buổi học

### Lifeline

1. Quản trị viên
2. Màn hình xếp lịch
3. Bộ điều khiển xếp lịch
4. Use case commit preview
5. Preview store
6. Kho dữ liệu lesson
7. DB transaction

### Thông điệp chính

1. Admin chọn commit preview
2. Controller gọi use case commit
3. Use case lấy preview theo `run_id`
4. Use case kiểm tra `status = COMPLETED`
5. Use case kiểm tra preview có assignments
6. Use case kiểm tra trùng lesson hiện có
7. Use case mở transaction
8. Repo batch insert `Lesson`
9. Commit transaction
10. Trả kết quả thành công

### Nhánh lỗi

- preview không tồn tại
- preview `FAILED` hoặc `PARTIAL`
- preview không có assignment
- lesson overlap
- transaction fail

## 7.12 Tải lên tài liệu giảng dạy

### Lifeline

1. Giáo viên
2. Màn hình tài liệu giảng dạy
3. Bộ điều khiển tài liệu
4. Use case tải lên tài liệu
5. Dịch vụ lưu file
6. Kho dữ liệu tài liệu
7. OCR service
8. AI moderation service
9. Kho dữ liệu audit log

### Thông điệp chính

1. Giáo viên chọn file và nhập metadata
2. UI gửi multipart/form-data
3. Controller nhận file
4. Use case validate file type/size
5. Dịch vụ lưu file vào storage local
6. Repo tạo `Material`
7. OCR service xử lý file
8. AI service gắn nhãn
9. Repo tạo `AuditLog`
10. Repo cập nhật `Material.latest_label_id` và `status`
11. Trả material detail

### Nhánh lỗi

- file quá 10MB
- file type không hợp lệ
- lỗi lưu file
- lỗi DB

## 7.13 Phê duyệt / từ chối tài liệu

### Lifeline

1. Reviewer / Admin
2. Màn hình hàng chờ duyệt tài liệu
3. Bộ điều khiển tài liệu
4. Use case duyệt tài liệu
5. Kho dữ liệu tài liệu
6. Kho dữ liệu quyết định duyệt

### Thông điệp chính

1. Reviewer mở chi tiết tài liệu
2. Xem metadata, label, audit reasoning
3. Chọn approve hoặc reject
4. Controller bind request review
5. Use case tải material
6. Use case tạo `ApprovalDecision`
7. Use case cập nhật `Material.status`
8. Trả kết quả review

### Nhánh lỗi

- material không tồn tại
- lỗi DB
- actor không đúng quyền, dù hiện code chưa chặn đủ

---

# 8. Gợi ý bộ sơ đồ use case nên vẽ

## 8.1 Sơ đồ tổng quan hệ thống

### Actors

- Người dùng chưa đăng nhập
- Quản trị viên
- Giáo viên
- Học viên
- Reviewer / Compliance
- SMTP
- Scheduling Engine
- OCR / AI Stub

### Use case tổng

- Quản lý truy cập tài khoản
- Quản lý học viên
- Quản lý giáo viên
- Quản lý khóa học
- Quản lý chương trình đào tạo
- Quản lý phòng học
- Quản lý ca học
- Quản lý lớp học
- Quản lý ghi danh lớp học
- Quản lý lịch tuần lớp học
- Xếp lịch học
- Quản lý buổi học
- Quản lý điểm danh
- Quản lý tổng kết buổi học
- Quản lý kết quả học tập
- Quản lý đơn xin phép
- Quản lý tài liệu giảng dạy
- Quản lý tư vấn tuyển sinh

## 8.2 Sơ đồ chi tiết nên tách riêng

1. Sơ đồ xác thực và tài khoản
2. Sơ đồ quản lý học viên
3. Sơ đồ quản lý giáo viên
4. Sơ đồ khóa học và chương trình đào tạo
5. Sơ đồ lớp học, ghi danh và lịch tuần
6. Sơ đồ xếp lịch và tạo buổi học
7. Sơ đồ tài liệu giảng dạy và kiểm duyệt
8. Sơ đồ học vụ sau buổi học

---

# 9. Các điểm mơ hồ cần ghi chú ngay trên sơ đồ

| Mã | Vấn đề | Gợi ý ghi chú trên sơ đồ |
|---|---|---|
| NOTE-01 | `Quản lý lịch tuần lớp học` chưa có API | đánh dấu `Future-state / Missing API` |
| NOTE-02 | `Điểm danh`, `Tổng kết buổi học`, `Kết quả học tập`, `Đơn xin phép` mới có data model | đánh dấu `Data model exists, API/UI pending` |
| NOTE-03 | `Xuất bản chương trình đào tạo`, `Lưu trữ chương trình đào tạo` chưa có use case chạy thật | đánh dấu `Lifecycle inferred` |
| NOTE-04 | `Reviewer / Compliance` chưa có role check backend riêng | đánh dấu `Role needs validation` |
| NOTE-05 | `Tạo / theo dõi consultation` chỉ mới có entity | đánh dấu `Inferred from data model` |

---

# 10. Kết luận modeling

Để tài liệu use case của EduCenter “đúng ý BA” và dễ vẽ, nên chốt theo nguyên tắc sau:

1. **Không gộp CRUD thành một oval duy nhất**.
2. **Tách use case tổng và use case con**.
3. **Dùng tiếng Việt hoàn toàn cho tên use case**.
4. **Chỉ dùng `<<include>>` cho bước bắt buộc, tái sử dụng**.
5. **Chỉ dùng `<<extend>>` cho nhánh tùy chọn hoặc hành vi phát sinh theo điều kiện**.
6. **Sequence diagram phải bám theo use case con quan trọng**, không bám theo use case tổng.

Nếu đi tiếp từ tài liệu này, bước tốt nhất là:

- vẽ sơ đồ use case tổng trước;
- sau đó tách 4 sơ đồ chi tiết quan trọng nhất:
  1. xác thực,
  2. lớp học và ghi danh,
  3. xếp lịch,
  4. tài liệu giảng dạy và kiểm duyệt.

