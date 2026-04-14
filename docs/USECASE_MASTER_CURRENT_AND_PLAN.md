# USE CASE MASTER DUY NHẤT CỦA HỆ THỐNG EDUCENTER

**Mục đích của file này**

Đây là **file chuẩn duy nhất** để mô tả use case của EduCenter trong phạm vi hiện tại và kế hoạch đã chốt.

File này chỉ giữ:

1. **Use case đang có trong hệ thống**
2. **Use case đã nằm trong kế hoạch chính thức**

File này **không đưa vào**:

- nhánh `AI Audit / Compliance / Chatbot`
- các use case chỉ có entity nhưng chưa nằm trong scope chính
- các chức năng ngoài mục tiêu đồ án hiện tại

**Căn cứ scope**

- [PROJECT_TASKS.md](/Users/hant/golang/doan/PROJECT_TASKS.md)
- [ke_hoach_phan_hoi_gvhd_2026-04-09.md](/Users/hant/golang/doan/docs/ke_hoach_phan_hoi_gvhd_2026-04-09.md)

**Scope hiện tại đã chốt**

1. Core admin modules:
   - tài khoản
   - học sinh
   - giáo viên
   - lớp học
   - phòng học
   - chương trình / khóa học
2. Scheduling là trọng tâm:
   - `Shift`
   - `class_schedule` dùng `shift_id`
   - 3 solver
   - benchmark
   - chọn solver chính
3. Predictive analytics:
   - `AT_RISK classification`
   - train trong backend
   - có prediction API
   - có màn hình cảnh báo

---

# 1. Nguyên tắc dùng file này để vẽ sơ đồ

## 1.1 Chỉ vẽ những gì có trong hệ thống hoặc đã có kế hoạch chính thức

### Được phép đưa vào sơ đồ

- module đang chạy được;
- module đang có backlog chính thức trong `PROJECT_TASKS.md`;
- use case phục vụ trực tiếp cho scheduling hoặc predictive analytics.

### Không đưa vào sơ đồ

- `AI Audit`
- `Compliance dashboard`
- `OCR/Gemini audit`
- `Chatbot`
- `Consultation / lead intake`
- `Attendance`, `Lesson Summary`, `Academic Record`, `Leave Request` dưới dạng use case vận hành độc lập, vì hiện chưa nằm trong backlog triển khai chính.

## 1.2 Cách hiểu “use case tổng” và “use case con”

- **Use case tổng**: nhóm chức năng lớn để đặt ở sơ đồ overview.
- **Use case con**: hành vi cụ thể để vẽ ở sơ đồ chi tiết.

Ví dụ:

- `Quản lý học sinh` là use case tổng
- `Tạo học sinh`, `Xem danh sách học sinh`, `Cập nhật học sinh`, `Xóa học sinh` là use case con

## 1.3 Quy tắc dùng `<<include>>` và `<<extend>>`

### Dùng `<<include>>` khi:

- use case cha **luôn phải gọi** use case con;
- use case con là bước bắt buộc;
- thiếu bước đó thì use case cha không hoàn thành.

### Dùng `<<extend>>` khi:

- use case chỉ xảy ra trong tình huống bổ sung;
- là nhánh tùy chọn;
- hoặc chỉ kích hoạt khi có điều kiện.

---

# 2. Danh sách use case tổng của hệ thống

Đây là danh sách use case tổng duy nhất nên dùng cho **sơ đồ overview**.

| Mã UC tổng | Tên use case tổng | Trạng thái | Actor chính | Ghi chú |
|---|---|---|---|---|
| UC-TONG-01 | Quản lý tài khoản người dùng | Đang có | Người dùng, Quản trị viên | Auth hiện chạy thật |
| UC-TONG-02 | Quản lý học sinh | Đang có | Quản trị viên | CRUD đã có |
| UC-TONG-03 | Quản lý giáo viên | Đang có | Quản trị viên | CRUD + lịch dạy + giờ dạy |
| UC-TONG-04 | Quản lý khóa học | Đang có | Quản trị viên | CRUD đã có |
| UC-TONG-05 | Quản lý chương trình đào tạo | Đang có | Quản trị viên | CRUD + liên kết khóa học |
| UC-TONG-06 | Quản lý phòng học | Đang có | Quản trị viên | CRUD đã có |
| UC-TONG-07 | Quản lý ca học | Đang có | Quản trị viên | CRUD `Shift` đã có |
| UC-TONG-08 | Quản lý lớp học | Đang có | Quản trị viên | CRUD + chi tiết lớp |
| UC-TONG-09 | Quản lý ghi danh lớp học | Đang có | Quản trị viên | add/remove học sinh trong lớp |
| UC-TONG-10 | Xếp lịch học | Đang có / đang mở rộng | Quản trị viên | preview + commit đã có |
| UC-TONG-11 | Benchmark thuật toán xếp lịch | Trong kế hoạch gần | Quản trị viên | admin API nội bộ |
| UC-TONG-12 | Dự báo học sinh có nguy cơ học kém | Trong kế hoạch | Quản trị viên | predictive analytics |
| UC-TONG-13 | Xem cảnh báo học sinh có nguy cơ học kém | Trong kế hoạch | Quản trị viên | UI cảnh báo |

---

# 3. Phân rã chi tiết use case theo từng phân hệ

## 3.1 Quản lý tài khoản người dùng

### Use case tổng

- `Quản lý tài khoản người dùng`

### Use case con đang có

| Mã | Tên use case | Trạng thái |
|---|---|---|
| UC-TK-01 | Đăng ký tài khoản | Đang có |
| UC-TK-02 | Xác minh email bằng OTP | Đang có |
| UC-TK-03 | Đăng nhập | Đang có |
| UC-TK-04 | Làm mới phiên đăng nhập | Đang có |
| UC-TK-05 | Đăng xuất | Đang có |
| UC-TK-06 | Yêu cầu quên mật khẩu | Đang có |
| UC-TK-07 | Đặt lại mật khẩu | Đang có |
| UC-TK-08 | Đổi mật khẩu | Đang có |
| UC-TK-09 | Xem hồ sơ cá nhân | Đang có nhưng backend còn lỗi `/me` |

### Quan hệ đề xuất để vẽ

#### `include`

- `Đăng ký tài khoản` `<<include>>` `Kiểm tra email đã tồn tại`
- `Đăng ký tài khoản` `<<include>>` `Tạo mã OTP xác minh`
- `Đăng ký tài khoản` `<<include>>` `Gửi email OTP`
- `Xác minh email bằng OTP` `<<include>>` `Kiểm tra hiệu lực OTP`
- `Yêu cầu quên mật khẩu` `<<include>>` `Tạo yêu cầu đặt lại mật khẩu`
- `Yêu cầu quên mật khẩu` `<<include>>` `Gửi email đặt lại mật khẩu`
- `Đăng nhập` `<<include>>` `Kiểm tra trạng thái tài khoản`

#### `extend`

- `Xem hồ sơ cá nhân` `<<extend>>` `Đăng nhập`

### Gợi ý sơ đồ sequence nên vẽ

1. `Đăng ký tài khoản`
2. `Xác minh email bằng OTP`
3. `Đăng nhập`
4. `Quên mật khẩu / Đặt lại mật khẩu`

## 3.2 Quản lý học sinh

### Use case tổng

- `Quản lý học sinh`

### Use case con đang có

| Mã | Tên use case | Trạng thái |
|---|---|---|
| UC-HS-01 | Tạo học sinh | Đang có |
| UC-HS-02 | Xem danh sách học sinh | Đang có |
| UC-HS-03 | Tìm kiếm học sinh | Đang có |
| UC-HS-04 | Xem chi tiết học sinh | Có thể suy ra từ UI/detail flow |
| UC-HS-05 | Cập nhật học sinh | Đang có |
| UC-HS-06 | Xóa học sinh | Đang có |

### Quan hệ đề xuất để vẽ

#### `extend`

- `Tìm kiếm học sinh` `<<extend>>` `Xem danh sách học sinh`
- `Xem chi tiết học sinh` `<<extend>>` `Xem danh sách học sinh`
- `Cập nhật học sinh` `<<extend>>` `Xem chi tiết học sinh`
- `Xóa học sinh` `<<extend>>` `Xem chi tiết học sinh`

### Gợi ý sequence nên vẽ

1. `Tạo học sinh`
2. `Xem danh sách học sinh`
3. `Cập nhật học sinh`
4. `Xóa học sinh`

## 3.3 Quản lý giáo viên

### Use case tổng

- `Quản lý giáo viên`

### Use case con đang có

| Mã | Tên use case | Trạng thái |
|---|---|---|
| UC-GV-01 | Tạo giáo viên | Đang có |
| UC-GV-02 | Xem danh sách giáo viên | Đang có |
| UC-GV-03 | Tìm kiếm giáo viên | Đang có |
| UC-GV-04 | Xem chi tiết giáo viên | Đang có |
| UC-GV-05 | Cập nhật giáo viên | Đang có |
| UC-GV-06 | Xóa giáo viên | Đang có |
| UC-GV-07 | Xem lịch dạy giáo viên | Đang có |
| UC-GV-08 | Xem thống kê giờ dạy giáo viên | Đang có |

### Quan hệ đề xuất để vẽ

#### `extend`

- `Tìm kiếm giáo viên` `<<extend>>` `Xem danh sách giáo viên`
- `Xem chi tiết giáo viên` `<<extend>>` `Xem danh sách giáo viên`
- `Cập nhật giáo viên` `<<extend>>` `Xem chi tiết giáo viên`
- `Xóa giáo viên` `<<extend>>` `Xem chi tiết giáo viên`
- `Xem lịch dạy giáo viên` `<<extend>>` `Xem chi tiết giáo viên`
- `Xem thống kê giờ dạy giáo viên` `<<extend>>` `Xem chi tiết giáo viên`

### Gợi ý sequence nên vẽ

1. `Tạo giáo viên`
2. `Xem lịch dạy giáo viên`
3. `Xem thống kê giờ dạy giáo viên`

## 3.4 Quản lý khóa học

### Use case tổng

- `Quản lý khóa học`

### Use case con đang có

| Mã | Tên use case | Trạng thái |
|---|---|---|
| UC-KH-01 | Tạo khóa học | Đang có |
| UC-KH-02 | Xem danh sách khóa học | Đang có |
| UC-KH-03 | Tìm kiếm khóa học | Đang có |
| UC-KH-04 | Xem chi tiết khóa học | Có thể suy ra từ UI |
| UC-KH-05 | Cập nhật khóa học | Đang có |
| UC-KH-06 | Xóa khóa học | Đang có |

### Quan hệ đề xuất để vẽ

#### `extend`

- `Tìm kiếm khóa học` `<<extend>>` `Xem danh sách khóa học`
- `Xem chi tiết khóa học` `<<extend>>` `Xem danh sách khóa học`
- `Cập nhật khóa học` `<<extend>>` `Xem chi tiết khóa học`
- `Xóa khóa học` `<<extend>>` `Xem chi tiết khóa học`

### Gợi ý sequence nên vẽ

1. `Tạo khóa học`
2. `Cập nhật khóa học`

## 3.5 Quản lý chương trình đào tạo

### Use case tổng

- `Quản lý chương trình đào tạo`

### Use case con đang có

| Mã | Tên use case | Trạng thái |
|---|---|---|
| UC-CT-01 | Tạo chương trình đào tạo | Đang có |
| UC-CT-02 | Xem danh sách chương trình đào tạo | Đang có |
| UC-CT-03 | Tìm kiếm chương trình đào tạo | Đang có |
| UC-CT-04 | Xem chi tiết chương trình đào tạo | Đang có |
| UC-CT-05 | Cập nhật chương trình đào tạo | Đang có |
| UC-CT-06 | Xóa chương trình đào tạo | Đang có |
| UC-CT-07 | Liên kết khóa học vào chương trình | Đang có |
| UC-CT-08 | Gỡ khóa học khỏi chương trình | Đang có |

### Use case trong kế hoạch nhưng chưa làm

| Mã | Tên use case | Trạng thái |
|---|---|---|
| UC-CT-09 | Xuất bản chương trình đào tạo | Trong kế hoạch dữ liệu/lifecycle, chưa có API |
| UC-CT-10 | Lưu trữ chương trình đào tạo | Trong kế hoạch dữ liệu/lifecycle, chưa có API |

### Quan hệ đề xuất để vẽ

#### `extend`

- `Tìm kiếm chương trình đào tạo` `<<extend>>` `Xem danh sách chương trình đào tạo`
- `Xem chi tiết chương trình đào tạo` `<<extend>>` `Xem danh sách chương trình đào tạo`
- `Cập nhật chương trình đào tạo` `<<extend>>` `Xem chi tiết chương trình đào tạo`
- `Xóa chương trình đào tạo` `<<extend>>` `Xem chi tiết chương trình đào tạo`
- `Liên kết khóa học vào chương trình` `<<extend>>` `Xem chi tiết chương trình đào tạo`
- `Gỡ khóa học khỏi chương trình` `<<extend>>` `Xem chi tiết chương trình đào tạo`

### Gợi ý sequence nên vẽ

1. `Tạo chương trình đào tạo`
2. `Liên kết khóa học vào chương trình`

## 3.6 Quản lý phòng học

### Use case tổng

- `Quản lý phòng học`

### Use case con đang có

| Mã | Tên use case | Trạng thái |
|---|---|---|
| UC-PH-01 | Tạo phòng học | Đang có |
| UC-PH-02 | Xem danh sách phòng học | Đang có |
| UC-PH-03 | Tìm kiếm phòng học | Đang có |
| UC-PH-04 | Xem chi tiết phòng học | Có thể suy ra từ UI |
| UC-PH-05 | Cập nhật phòng học | Đang có |
| UC-PH-06 | Xóa phòng học | Đang có |

### Quan hệ đề xuất để vẽ

#### `extend`

- `Tìm kiếm phòng học` `<<extend>>` `Xem danh sách phòng học`
- `Xem chi tiết phòng học` `<<extend>>` `Xem danh sách phòng học`
- `Cập nhật phòng học` `<<extend>>` `Xem chi tiết phòng học`
- `Xóa phòng học` `<<extend>>` `Xem chi tiết phòng học`

## 3.7 Quản lý ca học

### Use case tổng

- `Quản lý ca học`

### Use case con đang có

| Mã | Tên use case | Trạng thái |
|---|---|---|
| UC-CA-01 | Tạo ca học | Đang có |
| UC-CA-02 | Xem danh sách ca học | Đang có |
| UC-CA-03 | Tìm kiếm ca học | Có thể suy ra từ UI |
| UC-CA-04 | Xem chi tiết ca học | Có thể suy ra từ UI |
| UC-CA-05 | Cập nhật ca học | Đang có |
| UC-CA-06 | Xóa ca học | Đang có |

### Quan hệ đề xuất để vẽ

#### `extend`

- `Tìm kiếm ca học` `<<extend>>` `Xem danh sách ca học`
- `Xem chi tiết ca học` `<<extend>>` `Xem danh sách ca học`
- `Cập nhật ca học` `<<extend>>` `Xem chi tiết ca học`
- `Xóa ca học` `<<extend>>` `Xem chi tiết ca học`

### Gợi ý sequence nên vẽ

1. `Tạo ca học`
2. `Cập nhật ca học`

## 3.8 Quản lý lớp học

### Use case tổng

- `Quản lý lớp học`

### Use case con đang có

| Mã | Tên use case | Trạng thái |
|---|---|---|
| UC-LH-01 | Tạo lớp học | Đang có |
| UC-LH-02 | Xem danh sách lớp học | Đang có |
| UC-LH-03 | Tìm kiếm lớp học | Đang có |
| UC-LH-04 | Lọc lớp học theo trạng thái | Đang có |
| UC-LH-05 | Xem chi tiết lớp học | Đang có |
| UC-LH-06 | Cập nhật lớp học | Đang có |
| UC-LH-07 | Xóa lớp học | Đang có |
| UC-LH-08 | Phân công giáo viên cho lớp | Đang có |

### Quan hệ đề xuất để vẽ

#### `extend`

- `Tìm kiếm lớp học` `<<extend>>` `Xem danh sách lớp học`
- `Lọc lớp học theo trạng thái` `<<extend>>` `Xem danh sách lớp học`
- `Xem chi tiết lớp học` `<<extend>>` `Xem danh sách lớp học`
- `Cập nhật lớp học` `<<extend>>` `Xem chi tiết lớp học`
- `Xóa lớp học` `<<extend>>` `Xem chi tiết lớp học`
- `Phân công giáo viên cho lớp` `<<extend>>` `Xem chi tiết lớp học`

### Gợi ý sequence nên vẽ

1. `Tạo lớp học`
2. `Phân công giáo viên cho lớp`

## 3.9 Quản lý ghi danh lớp học

### Use case tổng

- `Quản lý ghi danh lớp học`

### Use case con đang có

| Mã | Tên use case | Trạng thái |
|---|---|---|
| UC-GD-01 | Xem danh sách học sinh trong lớp | Đang có |
| UC-GD-02 | Tìm kiếm học sinh khả dụng để ghi danh | Đang có |
| UC-GD-03 | Ghi danh học sinh vào lớp | Đang có |
| UC-GD-04 | Rút học sinh khỏi lớp | Đang có |

### Quan hệ đề xuất để vẽ

#### `extend`

- `Xem danh sách học sinh trong lớp` `<<extend>>` `Xem chi tiết lớp học`
- `Tìm kiếm học sinh khả dụng để ghi danh` `<<extend>>` `Ghi danh học sinh vào lớp`
- `Ghi danh học sinh vào lớp` `<<extend>>` `Xem chi tiết lớp học`
- `Rút học sinh khỏi lớp` `<<extend>>` `Xem chi tiết lớp học`

### Gợi ý sequence nên vẽ

1. `Ghi danh học sinh vào lớp`
2. `Rút học sinh khỏi lớp`

## 3.10 Xếp lịch học

### Use case tổng

- `Xếp lịch học`

### Use case con đang có

| Mã | Tên use case | Trạng thái |
|---|---|---|
| UC-XL-01 | Tạo preview xếp lịch | Đang có |
| UC-XL-02 | Xem preview xếp lịch | Đang có |
| UC-XL-03 | Xem xung đột xếp lịch | Đang có |
| UC-XL-04 | Xác nhận preview để tạo buổi học | Đang có |

### Use case trong kế hoạch đã chốt

| Mã | Tên use case | Trạng thái |
|---|---|---|
| UC-XL-05 | Benchmark thuật toán xếp lịch | Trong kế hoạch gần |
| UC-XL-06 | Chọn solver xếp lịch chính | Trong kế hoạch |
| UC-XL-07 | Dùng solver chính cho API xếp lịch | Trong kế hoạch |

### Quan hệ đề xuất để vẽ

#### `include`

- `Tạo preview xếp lịch` `<<include>>` `Nạp dữ liệu xếp lịch`
- `Tạo preview xếp lịch` `<<include>>` `Xây dựng bài toán xếp lịch`
- `Tạo preview xếp lịch` `<<include>>` `Chạy solver xếp lịch`
- `Tạo preview xếp lịch` `<<include>>` `Tổng hợp xung đột xếp lịch`
- `Xác nhận preview để tạo buổi học` `<<include>>` `Kiểm tra điều kiện commit`
- `Xác nhận preview để tạo buổi học` `<<include>>` `Kiểm tra trùng buổi học`
- `Xác nhận preview để tạo buổi học` `<<include>>` `Tạo buổi học từ preview`
- `Benchmark thuật toán xếp lịch` `<<include>>` `Chạy solver Graph Coloring`
- `Benchmark thuật toán xếp lịch` `<<include>>` `Chạy solver CP-SAT`
- `Benchmark thuật toán xếp lịch` `<<include>>` `Chạy solver Tabu Search`
- `Benchmark thuật toán xếp lịch` `<<include>>` `Tổng hợp kết quả benchmark`

#### `extend`

- `Xem xung đột xếp lịch` `<<extend>>` `Xem preview xếp lịch`
- `Xác nhận preview để tạo buổi học` `<<extend>>` `Xem preview xếp lịch`
- `Benchmark thuật toán xếp lịch` `<<extend>>` `Xếp lịch học`
- `Chọn solver xếp lịch chính` `<<extend>>` `Benchmark thuật toán xếp lịch`
- `Dùng solver chính cho API xếp lịch` `<<extend>>` `Chọn solver xếp lịch chính`

### Gợi ý sequence nên vẽ

1. `Tạo preview xếp lịch`
2. `Xác nhận preview để tạo buổi học`
3. `Benchmark thuật toán xếp lịch`

## 3.11 Predictive analytics

### Use case tổng

- `Dự báo học sinh có nguy cơ học kém`

### Use case con trong kế hoạch đã chốt

| Mã | Tên use case | Trạng thái |
|---|---|---|
| UC-DB-01 | Xác định dữ liệu đầu vào dự báo | Trong kế hoạch |
| UC-DB-02 | Huấn luyện mô hình phân loại nguy cơ học kém | Trong kế hoạch |
| UC-DB-03 | Đánh giá mô hình dự báo | Trong kế hoạch |
| UC-DB-04 | Dự báo nhãn nguy cơ học kém | Trong kế hoạch |
| UC-DB-05 | Xem danh sách học sinh có nguy cơ học kém | Trong kế hoạch |
| UC-DB-06 | Xem điểm số / mức độ rủi ro | Trong kế hoạch |
| UC-DB-07 | Xem giải thích cơ bản cho kết quả dự báo | Trong kế hoạch |

### Quan hệ đề xuất để vẽ

#### `include`

- `Huấn luyện mô hình phân loại nguy cơ học kém` `<<include>>` `Xác định dữ liệu đầu vào dự báo`
- `Huấn luyện mô hình phân loại nguy cơ học kém` `<<include>>` `Tách tập train/test`
- `Huấn luyện mô hình phân loại nguy cơ học kém` `<<include>>` `Tạo đặc trưng đầu vào`
- `Đánh giá mô hình dự báo` `<<include>>` `Tính metric đánh giá`
- `Dự báo nhãn nguy cơ học kém` `<<include>>` `Nạp model dự báo`
- `Dự báo nhãn nguy cơ học kém` `<<include>>` `Sinh điểm rủi ro`

#### `extend`

- `Xem điểm số / mức độ rủi ro` `<<extend>>` `Xem danh sách học sinh có nguy cơ học kém`
- `Xem giải thích cơ bản cho kết quả dự báo` `<<extend>>` `Xem danh sách học sinh có nguy cơ học kém`

### Gợi ý sequence nên vẽ

1. `Huấn luyện mô hình phân loại nguy cơ học kém`
2. `Dự báo nhãn nguy cơ học kém`
3. `Xem danh sách học sinh có nguy cơ học kém`

---

# 4. Danh sách use case NÊN vẽ và KHÔNG NÊN vẽ

## 4.1 Nên vẽ

### Sơ đồ tổng quan

- Quản lý tài khoản người dùng
- Quản lý học sinh
- Quản lý giáo viên
- Quản lý khóa học
- Quản lý chương trình đào tạo
- Quản lý phòng học
- Quản lý ca học
- Quản lý lớp học
- Quản lý ghi danh lớp học
- Xếp lịch học
- Benchmark thuật toán xếp lịch
- Dự báo học sinh có nguy cơ học kém
- Xem cảnh báo học sinh có nguy cơ học kém

### Sơ đồ chi tiết

- cụm tài khoản
- cụm học sinh
- cụm giáo viên
- cụm khóa học / chương trình
- cụm lớp học / ghi danh
- cụm xếp lịch
- cụm predictive analytics

## 4.2 Không nên vẽ vào scope hiện tại

Không nên đưa vào sơ đồ chính thức của đồ án ở thời điểm này:

- `Quản lý tài liệu giảng dạy`
- `Phê duyệt tài liệu`
- `Từ chối tài liệu`
- `Quản lý đơn xin phép`
- `Quản lý điểm danh`
- `Quản lý tổng kết buổi học`
- `Quản lý kết quả học tập`
- `Quản lý tư vấn tuyển sinh`

**Lý do:**

1. `AI Audit/Compliance` đã bị loại khỏi backlog chính.
2. `Attendance`, `Lesson Summary`, `Academic Record`, `Leave Request` chưa nằm trong luồng triển khai chính thức hiện tại của backlog.
3. `Consultation` mới có entity, chưa có module vận hành.

---

# 5. Gợi ý chi tiết để vẽ sequence diagram

## 5.1 Sequence nên vẽ ở scope hiện tại

### Nhóm tài khoản

1. `Đăng ký tài khoản`
2. `Xác minh email bằng OTP`
3. `Đăng nhập`
4. `Yêu cầu quên mật khẩu`
5. `Đặt lại mật khẩu`

### Nhóm quản trị dữ liệu

6. `Tạo học sinh`
7. `Tạo giáo viên`
8. `Tạo khóa học`
9. `Tạo chương trình đào tạo`
10. `Liên kết khóa học vào chương trình`
11. `Tạo lớp học`
12. `Phân công giáo viên cho lớp`
13. `Ghi danh học sinh vào lớp`
14. `Tạo ca học`

### Nhóm scheduling

15. `Tạo preview xếp lịch`
16. `Xác nhận preview để tạo buổi học`
17. `Benchmark thuật toán xếp lịch`

### Nhóm predictive analytics

18. `Huấn luyện mô hình phân loại nguy cơ học kém`
19. `Dự báo nhãn nguy cơ học kém`
20. `Xem danh sách học sinh có nguy cơ học kém`

## 5.2 Lifeline chuẩn nên dùng cho nhóm CRUD

Với các use case CRUD như `Tạo học sinh`, `Tạo giáo viên`, `Tạo khóa học`, `Tạo lớp học`, sequence nên có:

1. Actor
2. Màn hình quản lý
3. Controller
4. Use case
5. Repository
6. Database

### Message chuẩn

1. Actor nhập dữ liệu
2. UI gửi request
3. Controller bind request
4. Controller gọi use case
5. Use case validate
6. Use case gọi repository
7. Repository ghi DB
8. Trả response
9. UI refresh danh sách hoặc detail

## 5.3 Lifeline chuẩn cho scheduling

### Với `Tạo preview xếp lịch`

1. Quản trị viên
2. Màn hình xếp lịch
3. Controller scheduling
4. Use case tạo preview
5. Repository lớp học
6. Repository phòng học
7. Repository ca học
8. Scheduling solver service
9. Preview store

### Message chính

1. Gửi request preview
2. Validate date range
3. Nạp classes / rooms / shifts
4. Build scheduling input
5. Run solver
6. Nhận assignments / conflicts / summary
7. Lưu preview
8. Trả response

### Với `Xác nhận preview để tạo buổi học`

1. Quản trị viên
2. Màn hình xếp lịch
3. Controller scheduling
4. Use case commit preview
5. Preview store
6. Lesson repository
7. Database transaction

### Message chính

1. Gửi commit request
2. Lấy preview theo `run_id`
3. Kiểm tra điều kiện commit
4. Kiểm tra lesson overlap
5. Insert lessons
6. Commit transaction
7. Trả kết quả

## 5.4 Lifeline chuẩn cho benchmark solver

1. Quản trị viên
2. Màn hình benchmark nội bộ hoặc công cụ test admin
3. Controller benchmark
4. Use case benchmark
5. Solver catalog
6. Graph Coloring Solver
7. CP-SAT Solver
8. Tabu Search Solver
9. Metric aggregator

### Message chính

1. Gửi input benchmark
2. Nạp dữ liệu benchmark
3. Chạy solver 1
4. Chạy solver 2
5. Chạy solver 3
6. Tổng hợp feasibility / hard violations / soft score / runtime
7. Trả bảng so sánh

## 5.5 Lifeline chuẩn cho predictive analytics

### Với `Huấn luyện mô hình phân loại nguy cơ học kém`

1. Quản trị viên
2. Màn hình quản trị predictive hoặc công cụ admin nội bộ
3. Controller predictive
4. Use case huấn luyện mô hình
5. Data loader
6. Feature engineering service
7. Training service
8. Model repository / metadata store

### Message chính

1. Gửi yêu cầu train
2. Tải dữ liệu đầu vào
3. Tạo feature set
4. Chia train/test
5. Train các mô hình
6. Tính metric
7. Lưu metadata model
8. Trả kết quả huấn luyện

### Với `Dự báo nhãn nguy cơ học kém`

1. Quản trị viên
2. Màn hình cảnh báo
3. Controller predictive
4. Use case dự báo
5. Model loader
6. Feature builder
7. Prediction service

### Message chính

1. Gửi yêu cầu dự báo
2. Nạp model hiện hành
3. Lấy dữ liệu học sinh
4. Sinh feature
5. Dự báo label + score
6. Trả kết quả

---

# 6. Chốt scope use case cuối cùng để vẽ

## 6.1 Use case CURRENT-STATE

### Tài khoản

- Đăng ký tài khoản
- Xác minh email bằng OTP
- Đăng nhập
- Làm mới phiên đăng nhập
- Đăng xuất
- Yêu cầu quên mật khẩu
- Đặt lại mật khẩu
- Đổi mật khẩu
- Xem hồ sơ cá nhân

### Học sinh

- Tạo học sinh
- Xem danh sách học sinh
- Tìm kiếm học sinh
- Xem chi tiết học sinh
- Cập nhật học sinh
- Xóa học sinh

### Giáo viên

- Tạo giáo viên
- Xem danh sách giáo viên
- Tìm kiếm giáo viên
- Xem chi tiết giáo viên
- Cập nhật giáo viên
- Xóa giáo viên
- Xem lịch dạy giáo viên
- Xem thống kê giờ dạy giáo viên

### Khóa học

- Tạo khóa học
- Xem danh sách khóa học
- Tìm kiếm khóa học
- Xem chi tiết khóa học
- Cập nhật khóa học
- Xóa khóa học

### Chương trình đào tạo

- Tạo chương trình đào tạo
- Xem danh sách chương trình đào tạo
- Tìm kiếm chương trình đào tạo
- Xem chi tiết chương trình đào tạo
- Cập nhật chương trình đào tạo
- Xóa chương trình đào tạo
- Liên kết khóa học vào chương trình
- Gỡ khóa học khỏi chương trình

### Phòng học

- Tạo phòng học
- Xem danh sách phòng học
- Tìm kiếm phòng học
- Xem chi tiết phòng học
- Cập nhật phòng học
- Xóa phòng học

### Ca học

- Tạo ca học
- Xem danh sách ca học
- Tìm kiếm ca học
- Xem chi tiết ca học
- Cập nhật ca học
- Xóa ca học

### Lớp học

- Tạo lớp học
- Xem danh sách lớp học
- Tìm kiếm lớp học
- Lọc lớp học theo trạng thái
- Xem chi tiết lớp học
- Cập nhật lớp học
- Xóa lớp học
- Phân công giáo viên cho lớp

### Ghi danh

- Xem danh sách học sinh trong lớp
- Tìm kiếm học sinh khả dụng để ghi danh
- Ghi danh học sinh vào lớp
- Rút học sinh khỏi lớp

### Scheduling

- Tạo preview xếp lịch
- Xem preview xếp lịch
- Xem xung đột xếp lịch
- Xác nhận preview để tạo buổi học

## 6.2 Use case PLAN-STATE đã chốt

### Scheduling

- Benchmark thuật toán xếp lịch
- Chọn solver xếp lịch chính
- Dùng solver chính cho API xếp lịch

### Predictive analytics

- Xác định dữ liệu đầu vào dự báo
- Huấn luyện mô hình phân loại nguy cơ học kém
- Đánh giá mô hình dự báo
- Dự báo nhãn nguy cơ học kém
- Xem danh sách học sinh có nguy cơ học kém
- Xem điểm số / mức độ rủi ro
- Xem giải thích cơ bản cho kết quả dự báo

## 6.3 Những gì đã loại khỏi file này

Các nhánh sau **cố ý không đưa vào**:

- tải lên tài liệu giảng dạy
- phê duyệt / từ chối tài liệu
- AI Audit / Compliance
- OCR / Gemini / audit log
- điểm danh
- tổng kết buổi học
- kết quả học tập
- đơn xin phép
- tư vấn tuyển sinh

---

# 7. Kết luận

Nếu bạn chỉ muốn giữ **một file duy nhất** để tiếp tục vẽ sơ đồ, thì từ giờ nên dùng file này:

- [USECASE_MASTER_CURRENT_AND_PLAN.md](/Users/hant/golang/doan/docs/USECASE_MASTER_CURRENT_AND_PLAN.md)

File này đã được rút gọn đúng scope:

1. chỉ còn use case của hệ thống hiện tại;
2. cộng thêm use case đã nằm trong kế hoạch chính thức;
3. bỏ hẳn nhánh `AI Audit/Compliance`;
4. có sẵn:
   - danh sách use case tổng,
   - phân rã CRUD,
   - bảng `include/extend`,
   - gợi ý sequence diagram.

