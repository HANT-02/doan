# GÓI ĐẶC TẢ USE CASE HỆ THỐNG EDUCENTER

**Nguồn gốc tài liệu:** tổng hợp từ codebase EduCenter tại thời điểm 13/04/2026  
**Mục đích:** phục vụ báo cáo BA, annotation sơ đồ use case, activity/BPMN, rà soát yêu cầu với stakeholder  
**Nguyên tắc:** business-first, code-second; các hành vi chưa đủ chắc chắn sẽ được đánh dấu rõ

---

# 1. Hướng dẫn đọc tài liệu

## 1.1 Mức độ ưu tiên use case

- **Rất cao**: use case cốt lõi để hệ thống vận hành trung tâm ở mức tối thiểu; thiếu là mất giá trị sản phẩm.
- **Cao**: use case quan trọng cho vận hành thực tế nhưng chưa nhất thiết chặn demo/go-live lõi.
- **Trung bình**: use case hỗ trợ, mở rộng hoặc cần để hoàn thiện vòng đời nghiệp vụ.
- **Thấp / tương lai**: use case mới có dấu vết định hướng, data model hoặc roadmap.

## 1.2 Mức độ chắc chắn

- **Confirmed from code**: xác nhận trực tiếp từ route, controller, use case, entity, repository hoặc màn hình.
- **Strongly inferred from code**: suy luận mạnh từ mô hình dữ liệu, tên module, placeholder UI, migration hoặc comments.
- **Giả định / cần xác nhận BA**: chưa đủ bằng chứng hoặc đang có mâu thuẫn giữa các phần của code.

## 1.3 Ký hiệu và quy ước

- **Actor chính**: tác nhân khởi tạo use case và chịu trách nhiệm chính.
- **Actor phụ**: hệ thống hoặc vai trò hỗ trợ use case.
- **Luồng chính**: đường đi thành công chuẩn.
- **Luồng thay thế**: đường đi hợp lệ nhưng khác nhánh chính.
- **Luồng ngoại lệ**: lỗi, vi phạm điều kiện, hoặc hành vi hệ thống không thể tiếp tục.
- **Dữ liệu bị tác động**: dữ liệu tạo mới, cập nhật, xóa mềm hoặc chỉ đọc.
- **Ràng buộc phân quyền**: quyền hiện trạng theo code, không phải quyền mong muốn tương lai.

---

# 2. Ma trận tác nhân - use case

## 2.1 Ma trận Actor -> Use Case

| Mã use case | Tên use case | Guest | Admin | Teacher | Student | Compliance | Hệ thống SMTP | Scheduling Engine | OCR/AI Stub |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| UC-AUTH-01 | Đăng ký tài khoản | X |  |  |  |  | X |  |  |
| UC-AUTH-02 | Xác minh email OTP | X |  |  |  |  |  |  |  |
| UC-AUTH-03 | Đăng nhập | X | X | X | X | X |  |  |  |
| UC-AUTH-04 | Refresh token |  | X | X | X | X |  |  |  |
| UC-AUTH-05 | Quên mật khẩu | X | X | X | X | X | X |  |  |
| UC-AUTH-06 | Đặt lại mật khẩu | X | X | X | X | X |  |  |  |
| UC-AUTH-07 | Đổi mật khẩu |  | X | X | X | X |  |  |  |
| UC-AUTH-08 | Xem hồ sơ cá nhân |  | X | X | X | X |  |  |  |
| UC-STU-01 | Tạo học viên |  | X | ? | ? | ? |  |  |  |
| UC-STU-02 | Tìm kiếm / xem danh sách học viên |  | X | ? | ? | ? |  |  |  |
| UC-STU-03 | Cập nhật học viên |  | X | ? | ? | ? |  |  |  |
| UC-STU-04 | Xóa học viên |  | X | ? | ? | ? |  |  |  |
| UC-TCH-01 | Tạo giáo viên |  | X |  |  |  |  |  |  |
| UC-TCH-02 | Xem lịch dạy giáo viên | ? | X | X |  |  |  |  |  |
| UC-TCH-03 | Xem thống kê giờ dạy | ? | X | X |  |  |  |  |  |
| UC-CRS-01 | Tạo khóa học |  | X |  |  |  |  |  |  |
| UC-CRS-02 | Cập nhật khóa học |  | X |  |  |  |  |  |  |
| UC-PRG-01 | Tạo chương trình |  | X |  |  |  |  |  |  |
| UC-PRG-02 | Gán khóa học vào chương trình |  | X |  |  |  |  |  |  |
| UC-PRG-03 | Cập nhật / xuất bản / lưu trữ chương trình |  | X |  |  |  |  |  |  |
| UC-CLS-01 | Tạo lớp học |  | X |  |  |  |  |  |  |
| UC-CLS-02 | Cập nhật lớp học |  | X |  |  |  |  |  |  |
| UC-ENR-01 | Ghi danh học viên vào lớp |  | X |  |  |  |  |  |  |
| UC-ENR-02 | Rút học viên khỏi lớp |  | X |  |  |  |  |  |  |
| UC-CLS-03 | Phân công giáo viên cho lớp |  | X |  |  |  |  |  |  |
| UC-SCHCFG-01 | Cấu hình lịch tuần cho lớp |  | X |  |  |  |  |  |  |
| UC-SOL-01 | Tạo preview xếp lịch |  | X |  |  |  |  | X |  |
| UC-SOL-02 | Xem preview xếp lịch |  | X |  |  |  |  | X |  |
| UC-SOL-03 | Benchmark solver |  | X |  |  |  |  | X |  |
| UC-SOL-04 | Xác nhận preview để tạo lesson |  | X |  |  |  |  | X |  |
| UC-LSN-01 | Xem danh sách lesson |  | X | X | X? |  |  |  |  |
| UC-ATD-01 | Điểm danh |  | X? | X |  |  |  |  |  |
| UC-SUM-01 | Tạo tổng kết buổi học |  | X? | X |  |  |  |  |  |
| UC-ACR-01 | Ghi nhận academic record |  | X? | X |  |  |  |  |  |
| UC-LVE-01 | Tạo đơn xin phép |  | X? |  | X |  |  |  |  |
| UC-LVE-02 | Duyệt / từ chối đơn xin phép |  | X |  |  |  |  |  |  |
| UC-MAT-01 | Upload tài liệu |  |  | X |  |  |  |  | X |
| UC-MAT-02 | Duyệt tài liệu bị gắn cờ |  | X? | ? | ? | X |  |  | X |
| UC-CNS-01 | Tạo / theo dõi consultation | X? | X? |  | X? |  |  |  |  |

Ghi chú:
- `?` nghĩa là hiện trạng code cho thấy có thể thao tác được hoặc có nhu cầu nghiệp vụ, nhưng chưa có thiết kế role/use case chốt.

## 2.2 Ma trận Actor -> Permission

| Actor | Xem | Tạo | Cập nhật | Xóa | Duyệt | Chạy xử lý | Ghi chú |
|---|---|---|---|---|---|---|---|
| Guest | Có trên một số danh mục công khai, auth flows | Đăng ký tài khoản | Không | Không | Không | Không | Public GET đang mở cho teacher/course/program/class/room |
| Admin | Hầu hết module | Hầu hết module lõi | Hầu hết module lõi | Hầu hết module lõi | Có thể duyệt tài liệu theo business, nhưng code chưa khóa riêng | Có scheduling preview/commit/benchmark | Actor nội bộ quyền cao nhất hiện trạng |
| Teacher | Xem portal giáo viên, xem/tải tài liệu, lịch dạy | Upload tài liệu | Chưa rõ nhiều module khác | Không | Theo code có thể review material vì chỉ cần auth, nhưng không phù hợp nghiệp vụ | Không | Cần siết quyền material |
| Student | Xem portal học viên placeholder | Chưa rõ, có thể tạo đơn nghỉ trong tương lai | Chưa rõ | Không | Không | Không | Thực tế nhiều màn còn placeholder |
| Compliance | Xem queue tài liệu | Không | Không | Không | Duyệt / từ chối tài liệu | Không | Chưa có role middleware riêng trong backend |
| SMTP | Không | Không | Không | Không | Không | Gửi email OTP/reset | Actor hệ thống |
| Scheduling Engine | Không | Tạo preview nội bộ, lesson từ commit | Không | Không | Không | Chạy solver / benchmark contract | Actor hệ thống |
| OCR/AI Stub | Xem file bytes đầu vào ở tầng xử lý | Tạo audit log và nhãn | Cập nhật material status gián tiếp | Không | Không | Phân tích OCR/AI | Actor hệ thống |

---

# 3. Danh sách use case ưu tiên

## 3.1 Rất cao

1. UC-AUTH-01 Đăng ký tài khoản  
2. UC-AUTH-02 Xác minh email OTP  
3. UC-AUTH-03 Đăng nhập  
4. UC-AUTH-05 Quên mật khẩu  
5. UC-AUTH-06 Đặt lại mật khẩu  
6. UC-STU-01 Tạo học viên  
7. UC-TCH-01 Tạo giáo viên  
8. UC-CRS-01 Tạo khóa học  
9. UC-PRG-01 Tạo chương trình  
10. UC-PRG-02 Gán khóa học vào chương trình  
11. UC-CLS-01 Tạo lớp học  
12. UC-ENR-01 Ghi danh học viên vào lớp  
13. UC-CLS-03 Phân công giáo viên cho lớp  
14. UC-SCHCFG-01 Cấu hình lịch tuần cho lớp  
15. UC-SOL-01 Tạo preview xếp lịch  
16. UC-SOL-02 Xem preview xếp lịch  
17. UC-SOL-04 Xác nhận preview để tạo lesson

## 3.2 Cao

1. UC-AUTH-04 Refresh token  
2. UC-AUTH-07 Đổi mật khẩu  
3. UC-AUTH-08 Xem hồ sơ cá nhân  
4. UC-STU-02 Tìm kiếm / xem danh sách học viên  
5. UC-STU-03 Cập nhật học viên  
6. UC-STU-04 Xóa học viên  
7. UC-TCH-02 Xem lịch dạy giáo viên  
8. UC-TCH-03 Xem thống kê giờ dạy  
9. UC-CRS-02 Cập nhật khóa học  
10. UC-CLS-02 Cập nhật lớp học  
11. UC-ENR-02 Rút học viên khỏi lớp  
12. UC-SOL-03 Benchmark solver  
13. UC-MAT-01 Upload tài liệu  
14. UC-MAT-02 Duyệt tài liệu bị gắn cờ

## 3.3 Trung bình

1. UC-PRG-03 Cập nhật / xuất bản / lưu trữ chương trình  
2. UC-LSN-01 Xem danh sách lesson  
3. UC-ATD-01 Điểm danh  
4. UC-SUM-01 Tạo tổng kết buổi học  
5. UC-ACR-01 Ghi nhận academic record  
6. UC-LVE-01 Tạo đơn xin phép  
7. UC-LVE-02 Duyệt / từ chối đơn xin phép

## 3.4 Thấp / tương lai

1. UC-CNS-01 Tạo / theo dõi consultation  
2. Các use case dự báo học viên có nguy cơ học kém  
3. Các use case dashboard/báo cáo dữ liệu thật ngoài teaching-hours và scheduling summary

---

# 4. Đặc tả chi tiết use case

## UC-AUTH-01. Đăng ký tài khoản

- **Mã use case:** UC-AUTH-01
- **Tên use case:** Đăng ký tài khoản
- **Mục tiêu nghiệp vụ:** Tạo tài khoản mới để người dùng có thể tham gia hệ thống và bước sang giai đoạn xác minh email.
- **Phạm vi:** Xác thực và tài khoản
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Người dùng chưa đăng nhập
- **Tác nhân phụ:** Hệ thống gửi email SMTP
- **Mô tả ngắn:** Người dùng cung cấp email, họ tên và mật khẩu đã được mã hóa từ frontend; hệ thống tạo tài khoản ở trạng thái chưa kích hoạt và gửi OTP xác minh.
- **Tiền điều kiện:** Chưa có tài khoản sử dụng email đó.
- **Điều kiện kích hoạt:** Người dùng gửi form đăng ký.
- **Hậu điều kiện thành công:** User mới được tạo; tài khoản ở trạng thái chưa kích hoạt; OTP active được tạo; email OTP được gửi hoặc được log ra trong môi trường dev.
- **Hậu điều kiện thất bại:** Không tạo user; không tạo OTP.
- **Dữ liệu đầu vào:** `email`, `full_name`, `password_enc`
- **Dữ liệu đầu ra:** `user_id`
- **Thực thể bị tác động:** `User`, `UserOTP`
- **Quy tắc nghiệp vụ áp dụng:**
  - Email user phải duy nhất.
  - Mật khẩu từ frontend phải giải mã được.
  - Tài khoản đăng ký mới không được login ngay nếu chưa xác minh OTP.
- **Luồng chính:**
  1. Người dùng nhập email, họ tên và mật khẩu trên màn đăng ký.
  2. Frontend gửi `password_enc` thay vì plain password.
  3. Backend kiểm tra email đã tồn tại hay chưa.
  4. Backend giải mã mật khẩu.
  5. Backend băm mật khẩu.
  6. Backend sinh OTP 6 số.
  7. Backend băm OTP.
  8. Backend tạo `User` với `is_active = false`.
  9. Backend tạo bản ghi `UserOTP` với hạn dùng 5 phút.
  10. Backend commit transaction.
  11. Backend gửi email OTP bất đồng bộ.
  12. Hệ thống trả về `user_id` cho frontend để phục vụ bước xác minh OTP.
- **Luồng thay thế:**
  - A1. Môi trường dev/local chưa cấu hình SMTP: hệ thống vẫn tạo user và OTP, nhưng thay vì gửi mail thật sẽ log nội dung mail.
  - A2. Email `test@gmail.com` trong code hiện tại có OTP cố định `123456` để phục vụ test/demo.
- **Luồng ngoại lệ:**
  - E1. Email đã tồn tại: trả lỗi đăng ký thất bại.
  - E2. `password_enc` không giải mã được: trả lỗi payload mật khẩu không hợp lệ.
  - E3. Lỗi DB trong transaction: không tạo user, không tạo OTP.
  - E4. Gửi mail lỗi sau khi transaction đã commit: user vẫn được tạo nhưng người dùng có thể không nhận được OTP.
- **Ràng buộc phân quyền:** Public, không cần đăng nhập.
- **API / module / màn hình liên quan:** `POST /api/v1/auth/register`, `RegisterPage`, `internal/usecases/user/register.go`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Chưa có use case resend OTP.
  - Chưa có giới hạn số lần đăng ký hoặc số lần gửi OTP.
  - **Độ tin cậy:** Confirmed from code

## UC-AUTH-02. Xác minh email OTP

- **Mã use case:** UC-AUTH-02
- **Tên use case:** Xác minh email OTP
- **Mục tiêu nghiệp vụ:** Kích hoạt tài khoản sau khi người dùng xác nhận email.
- **Phạm vi:** Xác thực và tài khoản
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Người dùng chưa đăng nhập / người dùng vừa đăng ký
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Người dùng nhập OTP; hệ thống kiểm tra OTP còn hạn, chưa dùng, đúng với user tương ứng và kích hoạt tài khoản.
- **Tiền điều kiện:** User đã đăng ký; tồn tại OTP active cho user.
- **Điều kiện kích hoạt:** Người dùng submit OTP.
- **Hậu điều kiện thành công:** OTP bị đánh dấu đã dùng; tài khoản được kích hoạt.
- **Hậu điều kiện thất bại:** Tài khoản vẫn chưa kích hoạt.
- **Dữ liệu đầu vào:** `user_id`, `otp`
- **Dữ liệu đầu ra:** `message`
- **Thực thể bị tác động:** `User`, `UserOTP`
- **Quy tắc nghiệp vụ áp dụng:**
  - OTP phải tồn tại và chưa được dùng.
  - OTP phải còn thời hạn.
  - OTP phải khớp với hash đã lưu.
- **Luồng chính:**
  1. Người dùng nhập OTP đã nhận được.
  2. Backend tra OTP active của `user_id`.
  3. Backend kiểm tra thời điểm hiện tại có vượt `expired_at` hay không.
  4. Backend so sánh OTP đầu vào với hash OTP đã lưu.
  5. Backend đánh dấu OTP là đã dùng.
  6. Backend cập nhật user sang trạng thái active.
  7. Backend trả thông báo xác minh thành công.
- **Luồng thay thế:**
  - A1. Người dùng nhập lại OTP sau lần đầu sai nhưng OTP vẫn còn hạn: hệ thống cho phép xác minh lại.
- **Luồng ngoại lệ:**
  - E1. Không tìm thấy OTP active: báo OTP không tồn tại hoặc đã hết hạn.
  - E2. OTP hết hạn: báo OTP expired.
  - E3. OTP sai: báo invalid OTP.
  - E4. Lỗi transaction khi mark OTP used hoặc activate user: không kích hoạt tài khoản.
- **Ràng buộc phân quyền:** Public.
- **API / module / màn hình liên quan:** `POST /api/v1/auth/verify-otp`, `internal/usecases/user/verify_otp.go`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Chưa thấy cơ chế khóa tài khoản hoặc giới hạn số lần nhập OTP sai.
  - **Độ tin cậy:** Confirmed from code

## UC-AUTH-03. Đăng nhập

- **Mã use case:** UC-AUTH-03
- **Tên use case:** Đăng nhập
- **Mục tiêu nghiệp vụ:** Cho phép người dùng truy cập hệ thống bằng tài khoản hợp lệ.
- **Phạm vi:** Xác thực và tài khoản
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Người dùng
- **Tác nhân phụ:** Dịch vụ JWT
- **Mô tả ngắn:** Người dùng đăng nhập bằng email và mật khẩu; hệ thống kiểm tra tài khoản hoạt động, xác thực mật khẩu rồi phát hành access token và refresh token.
- **Tiền điều kiện:** User tồn tại và đã active.
- **Điều kiện kích hoạt:** Gửi form login.
- **Hậu điều kiện thành công:** Người dùng nhận token và thông tin profile cơ bản.
- **Hậu điều kiện thất bại:** Không cấp token.
- **Dữ liệu đầu vào:** `username` (thực chất là email), `password`
- **Dữ liệu đầu ra:** `access_token`, `refresh_token`, `user`
- **Thực thể bị tác động:** Không cập nhật dữ liệu domain
- **Quy tắc nghiệp vụ áp dụng:**
  - User được tra theo `email`.
  - User phải active.
  - Password phải khớp bcrypt hash.
- **Luồng chính:**
  1. Người dùng nhập email và mật khẩu.
  2. Backend tra cứu user theo email.
  3. Backend kiểm tra `is_active`.
  4. Backend so sánh mật khẩu.
  5. Backend đọc cấu hình JWT.
  6. Backend sinh access token theo thời hạn cấu hình.
  7. Backend sinh refresh token theo thời hạn cấu hình.
  8. Backend trả token và thông tin user.
- **Luồng thay thế:**
  - A1. Không parse được thời hạn JWT từ config: hệ thống dùng mặc định 24h cho access token và 7 ngày cho refresh token.
- **Luồng ngoại lệ:**
  - E1. Không tìm thấy user.
  - E2. User chưa active.
  - E3. Sai mật khẩu.
  - E4. Lỗi sinh token.
- **Ràng buộc phân quyền:** Public.
- **API / module / màn hình liên quan:** `POST /api/v1/auth/login`, `LoginPage`, `internal/services/user/auth.go`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Trường `username` thực tế đang được hiểu là email, không phải login name linh hoạt.
  - **Độ tin cậy:** Confirmed from code

## UC-AUTH-04. Refresh token

- **Mã use case:** UC-AUTH-04
- **Tên use case:** Refresh token
- **Mục tiêu nghiệp vụ:** Gia hạn access token khi access token sắp hết hạn hoặc frontend reload.
- **Phạm vi:** Xác thực và tài khoản
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Người dùng đã có refresh token
- **Tác nhân phụ:** Dịch vụ JWT
- **Mô tả ngắn:** Backend nhận refresh token hợp lệ và phát hành access token mới.
- **Tiền điều kiện:** Người dùng đã từng login và giữ refresh token hợp lệ.
- **Điều kiện kích hoạt:** Frontend gọi refresh flow.
- **Hậu điều kiện thành công:** Có access token mới.
- **Hậu điều kiện thất bại:** Người dùng cần đăng nhập lại.
- **Dữ liệu đầu vào:** `refresh_token`
- **Dữ liệu đầu ra:** `access_token`
- **Thực thể bị tác động:** Không cập nhật dữ liệu domain
- **Quy tắc nghiệp vụ áp dụng:** Refresh token được validate dựa trên JWT signature và expiry.
- **Luồng chính:**
  1. Frontend phát hiện cần refresh.
  2. Gửi refresh token tới API.
  3. Backend đọc cấu hình JWT.
  4. Backend validate refresh token.
  5. Backend sinh access token mới.
  6. Trả về access token mới.
- **Luồng thay thế:** Không có notable business branch.
- **Luồng ngoại lệ:**
  - E1. Refresh token không hợp lệ.
  - E2. Refresh token hết hạn.
  - E3. Lỗi cấu hình JWT.
- **Ràng buộc phân quyền:** Public.
- **API / module / màn hình liên quan:** `POST /api/v1/auth/refresh`, `internal/usecases/user/refresh_token.go`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Chưa có persistent refresh session; backend không revoke refresh token theo user/session.
  - **Độ tin cậy:** Confirmed from code

## UC-AUTH-05. Quên mật khẩu

- **Mã use case:** UC-AUTH-05
- **Tên use case:** Quên mật khẩu
- **Mục tiêu nghiệp vụ:** Cho phép người dùng yêu cầu cấp quyền đặt lại mật khẩu qua email.
- **Phạm vi:** Xác thực và tài khoản
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Người dùng chưa nhớ mật khẩu
- **Tác nhân phụ:** SMTP
- **Mô tả ngắn:** Người dùng nhập email; hệ thống tạo reset token nếu tài khoản tồn tại và gửi link đặt lại mật khẩu.
- **Tiền điều kiện:** Người dùng có email hợp lệ; có thể tồn tại hoặc không.
- **Điều kiện kích hoạt:** Gửi form forgot password.
- **Hậu điều kiện thành công:** Nếu email tồn tại, một reset record được tạo và email được gửi; nếu không tồn tại, hệ thống vẫn trả success message chung.
- **Hậu điều kiện thất bại:** Không tạo reset record khi lỗi DB.
- **Dữ liệu đầu vào:** `email`
- **Dữ liệu đầu ra:** `message`
- **Thực thể bị tác động:** `PasswordReset`
- **Quy tắc nghiệp vụ áp dụng:**
  - Không được lộ thông tin email có tồn tại hay không.
  - Reset token phải có hạn dùng.
  - Có lưu IP và user agent nếu được truyền.
- **Luồng chính:**
  1. Người dùng nhập email.
  2. Backend normalize email.
  3. Backend tra user theo email.
  4. Nếu user tồn tại, backend sinh token đặt lại mật khẩu.
  5. Backend băm token.
  6. Backend tạo `PasswordReset` trong transaction với `expires_at`.
  7. Backend tạo reset link và gửi email bất đồng bộ.
  8. Backend trả message chung.
- **Luồng thay thế:**
  - A1. Email không tồn tại: backend vẫn trả message chung, không tạo reset record.
  - A2. SMTP không cấu hình trong dev: mail được log thay vì gửi thật.
- **Luồng ngoại lệ:**
  - E1. Email trống hoặc không hợp lệ: báo payload invalid.
  - E2. Lỗi DB khi tra user hoặc tạo reset record.
- **Ràng buộc phân quyền:** Public.
- **API / module / màn hình liên quan:** `POST /api/v1/auth/forgot-password`, `ForgotPasswordPage`, `forgot_reset_change.go`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Chưa thấy chính sách vô hiệu các reset token cũ.
  - **Độ tin cậy:** Confirmed from code

## UC-AUTH-06. Đặt lại mật khẩu

- **Mã use case:** UC-AUTH-06
- **Tên use case:** Đặt lại mật khẩu
- **Mục tiêu nghiệp vụ:** Cho phép người dùng đặt mật khẩu mới sau khi có reset token hợp lệ.
- **Phạm vi:** Xác thực và tài khoản
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Người dùng
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Người dùng gửi token và mật khẩu mới; hệ thống xác minh token, đánh dấu token đã dùng và cập nhật mật khẩu.
- **Tiền điều kiện:** Có reset token hợp lệ, chưa dùng, chưa hết hạn.
- **Điều kiện kích hoạt:** Người dùng submit form reset password.
- **Hậu điều kiện thành công:** Mật khẩu mới được lưu; token reset bị mark used.
- **Hậu điều kiện thất bại:** Không đổi mật khẩu.
- **Dữ liệu đầu vào:** `token`, `new_password_enc`
- **Dữ liệu đầu ra:** `message`
- **Thực thể bị tác động:** `PasswordReset`, `User`
- **Quy tắc nghiệp vụ áp dụng:**
  - Token phải hợp lệ và chưa hết hạn.
  - Token reset chỉ dùng một lần.
  - Mật khẩu mới tối thiểu 8 ký tự ở use case hiện tại.
- **Luồng chính:**
  1. Người dùng mở reset link và nhập mật khẩu mới.
  2. Backend nhận token và mật khẩu mới đã mã hóa.
  3. Backend băm token đầu vào để tra `PasswordReset`.
  4. Backend transaction:
     1. Lấy reset record theo token hash.
     2. Kiểm tra record chưa used.
     3. Mark token used.
  5. Backend giải mã mật khẩu mới.
  6. Backend kiểm tra độ dài tối thiểu.
  7. Backend băm mật khẩu mới.
  8. Backend cập nhật password của user.
  9. Trả thông báo thành công.
- **Luồng thay thế:** Không có nhánh business hợp lệ khác.
- **Luồng ngoại lệ:**
  - E1. Token không hợp lệ hoặc hết hạn.
  - E2. Token đã dùng rồi.
  - E3. Mật khẩu mới không giải mã được.
  - E4. Mật khẩu mới quá ngắn.
  - E5. Lỗi DB cập nhật password.
- **Ràng buộc phân quyền:** Public.
- **API / module / màn hình liên quan:** `POST /api/v1/auth/reset-password`, `ResetPasswordPage`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Chưa có cơ chế revoke toàn bộ session/token sau khi reset.
  - **Độ tin cậy:** Confirmed from code

## UC-AUTH-07. Đổi mật khẩu

- **Mã use case:** UC-AUTH-07
- **Tên use case:** Đổi mật khẩu
- **Mục tiêu nghiệp vụ:** Cho phép người dùng đã đăng nhập đổi mật khẩu chủ động.
- **Phạm vi:** Xác thực và tài khoản
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Người dùng đã đăng nhập
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Người dùng nhập mật khẩu cũ và mật khẩu mới; hệ thống xác minh mật khẩu cũ rồi cập nhật mật khẩu mới.
- **Tiền điều kiện:** Có JWT hợp lệ; user tồn tại.
- **Điều kiện kích hoạt:** Người dùng chọn “Đổi mật khẩu”.
- **Hậu điều kiện thành công:** Password được cập nhật.
- **Hậu điều kiện thất bại:** Password không đổi.
- **Dữ liệu đầu vào:** `old_password_enc`, `new_password_enc`
- **Dữ liệu đầu ra:** `message`
- **Thực thể bị tác động:** `User`
- **Quy tắc nghiệp vụ áp dụng:** Mật khẩu cũ phải đúng.
- **Luồng chính:**
  1. Người dùng mở trang đổi mật khẩu.
  2. Frontend gửi request kèm JWT và 2 trường mật khẩu đã mã hóa.
  3. Controller lấy `user_id` từ auth context.
  4. Use case lấy user theo ID.
  5. Use case giải mã mật khẩu cũ.
  6. So sánh với password hash hiện tại.
  7. Giải mã mật khẩu mới.
  8. Băm mật khẩu mới.
  9. Cập nhật password.
- **Luồng thay thế:** Không có.
- **Luồng ngoại lệ:**
  - E1. Không có `user_id` trong context.
  - E2. User không tồn tại.
  - E3. Mật khẩu cũ sai.
  - E4. Payload mật khẩu không giải mã được.
- **Ràng buộc phân quyền:** Yêu cầu đăng nhập.
- **API / module / màn hình liên quan:** `POST /api/v1/auth/change-password`, `ChangePasswordPage`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Route không được gắn auth middleware rõ ràng tại layer route; controller kỳ vọng context đã có auth.
  - **Độ tin cậy:** Confirmed from code

## UC-AUTH-08. Xem hồ sơ cá nhân

- **Mã use case:** UC-AUTH-08
- **Tên use case:** Xem hồ sơ cá nhân
- **Mục tiêu nghiệp vụ:** Cho người dùng xem thông tin tài khoản hiện tại.
- **Phạm vi:** Xác thực và tài khoản
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Người dùng đã đăng nhập
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Hệ thống lấy thông tin user hiện tại để hiển thị ở profile page.
- **Tiền điều kiện:** Có JWT hợp lệ.
- **Điều kiện kích hoạt:** Người dùng mở màn hình hồ sơ.
- **Hậu điều kiện thành công:** Trả về profile user.
- **Hậu điều kiện thất bại:** Không lấy được hồ sơ; có thể buộc logout hoặc báo lỗi.
- **Dữ liệu đầu vào:** JWT context
- **Dữ liệu đầu ra:** `id`, `code`, `full_name`, `email`, `role`, `is_active`
- **Thực thể bị tác động:** Chỉ đọc `User`
- **Quy tắc nghiệp vụ áp dụng:** Phải lấy đúng user hiện tại từ token.
- **Luồng chính (mong muốn nghiệp vụ):**
  1. Người dùng mở trang profile.
  2. Frontend gọi API `/auth/me`.
  3. Backend xác định user hiện tại từ JWT.
  4. Backend lấy thông tin user.
  5. Trả về profile.
- **Luồng thay thế:** Frontend dùng dữ liệu user đã lưu trong state local thay vì gọi API.
- **Luồng ngoại lệ:**
  - E1. Route không có auth middleware nên `user_id` không có trong context.
  - E2. Controller đang ép `user_id` UUID sang `int`, gây lỗi runtime logic.
- **Ràng buộc phân quyền:** Yêu cầu đăng nhập theo thiết kế nghiệp vụ; hiện trạng route chưa khóa đúng.
- **API / module / màn hình liên quan:** `GET /api/v1/auth/me`, `ProfilePage`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Đây là use case đang **defective** trong code hiện tại.
  - **Độ tin cậy:** Confirmed from code

## UC-STU-01. Tạo học viên

- **Mã use case:** UC-STU-01
- **Tên use case:** Tạo học viên
- **Mục tiêu nghiệp vụ:** Tạo hồ sơ học viên để phục vụ ghi danh, quản lý lớp và theo dõi học tập.
- **Phạm vi:** Quản lý học viên
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Quản trị viên / vận hành
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Tạo record học viên với thông tin cá nhân, liên hệ phụ huynh và khối lớp.
- **Tiền điều kiện:** Người thao tác có quyền truy cập module học viên.
- **Điều kiện kích hoạt:** Người dùng submit form tạo học viên.
- **Hậu điều kiện thành công:** Student record mới được tạo.
- **Hậu điều kiện thất bại:** Không tạo student.
- **Dữ liệu đầu vào:** `code`, `full_name`, `email`, `phone`, `guardian_phone`, `grade_level`, `status`, `date_of_birth`, `gender`, `address`
- **Dữ liệu đầu ra:** Student detail
- **Thực thể bị tác động:** `Student`
- **Quy tắc nghiệp vụ áp dụng:** Hiện code chỉ bind field, chưa có validation business mạnh.
- **Luồng chính:**
  1. Người vận hành mở màn hình học viên.
  2. Nhập thông tin cơ bản học viên.
  3. Gửi request tạo học viên.
  4. Backend bind request.
  5. Use case tạo entity `Student`.
  6. Repository lưu vào DB.
  7. Trả student detail cho frontend.
- **Luồng thay thế:** Không có notable business branch.
- **Luồng ngoại lệ:**
  - E1. Lỗi bind request body.
  - E2. Lỗi DB khi create.
- **Ràng buộc phân quyền:** Theo business nên là Admin/Vận hành; theo code hiện tại là bất kỳ user đã auth.
- **API / module / màn hình liên quan:** `POST /api/v1/students`, `StudentsPage`, `StudentDialog`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Chưa thấy validation unique cho mã/email học viên.
  - Cần BA xác nhận ai được quản lý student.
  - **Độ tin cậy:** Confirmed from code

## UC-STU-02. Tìm kiếm / xem danh sách học viên

- **Mã use case:** UC-STU-02
- **Tên use case:** Tìm kiếm / xem danh sách học viên
- **Mục tiêu nghiệp vụ:** Tra cứu học viên phục vụ ghi danh, quản lý lớp và vận hành.
- **Phạm vi:** Quản lý học viên
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Quản trị viên / vận hành
- **Tác nhân phụ:** Có thể là giáo viên hoặc user khác theo hiện trạng code
- **Mô tả ngắn:** Hệ thống trả danh sách học viên theo bộ lọc.
- **Tiền điều kiện:** Người dùng đã đăng nhập.
- **Điều kiện kích hoạt:** Mở màn hình học viên hoặc nhập điều kiện tìm kiếm.
- **Hậu điều kiện thành công:** Trả danh sách học viên.
- **Hậu điều kiện thất bại:** Không có dữ liệu hoặc lỗi truy vấn.
- **Dữ liệu đầu vào:** `search`, `status`, `page`, `limit`, `sortBy`, `sortOrder`
- **Dữ liệu đầu ra:** Danh sách student + pagination
- **Thực thể bị tác động:** Chỉ đọc `Student`
- **Quy tắc nghiệp vụ áp dụng:** Search và filter status.
- **Luồng chính:**
  1. Người dùng mở danh sách học viên.
  2. Chọn tìm kiếm theo tên/mã và trạng thái.
  3. Frontend gửi query params.
  4. Backend build điều kiện tìm kiếm.
  5. Repository lấy dữ liệu và phân trang.
  6. Trả danh sách học viên.
- **Luồng thay thế:**
  - A1. Không có điều kiện tìm kiếm: trả danh sách mặc định.
- **Luồng ngoại lệ:**
  - E1. Lỗi DB khi truy vấn.
- **Ràng buộc phân quyền:** Hiện chỉ cần auth.
- **API / module / màn hình liên quan:** `GET /api/v1/students`, `StudentsPage`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Câu lệnh search ở use case đang ghép điều kiện theo cách không thực sự chuẩn, cần xác nhận behavior thực tế khi search.
  - **Độ tin cậy:** Confirmed from code

## UC-STU-03. Cập nhật học viên

- **Mã use case:** UC-STU-03
- **Tên use case:** Cập nhật học viên
- **Mục tiêu nghiệp vụ:** Chỉnh sửa thông tin học viên khi có thay đổi.
- **Phạm vi:** Quản lý học viên
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Quản trị viên / vận hành
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Cập nhật thông tin student theo ID.
- **Tiền điều kiện:** Student tồn tại.
- **Điều kiện kích hoạt:** Người dùng submit form sửa học viên.
- **Hậu điều kiện thành công:** Student được cập nhật.
- **Hậu điều kiện thất bại:** Student giữ nguyên.
- **Dữ liệu đầu vào:** ID học viên + các field cập nhật
- **Dữ liệu đầu ra:** Student updated
- **Thực thể bị tác động:** `Student`
- **Quy tắc nghiệp vụ áp dụng:** Chưa có validation business đáng kể.
- **Luồng chính:**
  1. Người dùng chọn một học viên cần sửa.
  2. Frontend hiển thị dữ liệu hiện tại.
  3. Người dùng sửa dữ liệu.
  4. Frontend gửi request cập nhật.
  5. Backend bind body.
  6. Use case update student theo ID.
  7. Trả student đã cập nhật.
- **Luồng thay thế:** Cập nhật một phần field.
- **Luồng ngoại lệ:**
  - E1. Student không tồn tại.
  - E2. Lỗi DB update.
- **Ràng buộc phân quyền:** Theo business nên là Admin/Vận hành; hiện trạng là mọi auth user.
- **API / module / màn hình liên quan:** `PUT /api/v1/students/:id`, `StudentsPage`, `StudentDialog`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Cần khóa role và audit trail thay đổi.
- **Độ tin cậy:** Confirmed from code

## UC-STU-04. Xóa học viên

- **Mã use case:** UC-STU-04
- **Tên use case:** Xóa học viên
- **Mục tiêu nghiệp vụ:** Loại bỏ học viên khỏi danh mục vận hành khi không còn sử dụng hồ sơ.
- **Phạm vi:** Quản lý học viên
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Quản trị viên / vận hành
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Xóa mềm hồ sơ học viên.
- **Tiền điều kiện:** Student tồn tại.
- **Điều kiện kích hoạt:** Người dùng xác nhận xóa.
- **Hậu điều kiện thành công:** Student bị soft delete.
- **Hậu điều kiện thất bại:** Student còn nguyên.
- **Dữ liệu đầu vào:** `student_id`
- **Dữ liệu đầu ra:** message
- **Thực thể bị tác động:** `Student`
- **Quy tắc nghiệp vụ áp dụng:** Hiện chưa thấy chặn xóa student đang có enrollment/lesson history.
- **Luồng chính:**
  1. Người dùng chọn xóa học viên.
  2. Frontend gửi request delete.
  3. Backend gọi use case delete.
  4. Repository soft delete student.
  5. Hệ thống trả thông báo thành công.
- **Luồng ngoại lệ:**
  - E1. Student không tồn tại.
  - E2. Lỗi DB xóa.
- **Ràng buộc phân quyền:** Theo business nên là Admin/Vận hành; hiện trạng là mọi auth user.
- **API / module / màn hình liên quan:** `DELETE /api/v1/students/:id`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Có cho xóa học viên đã học xong hay chỉ chuyển trạng thái?
- **Độ tin cậy:** Confirmed from code

## UC-TCH-01. Tạo giáo viên

- **Mã use case:** UC-TCH-01
- **Tên use case:** Tạo giáo viên
- **Mục tiêu nghiệp vụ:** Tạo hồ sơ giảng dạy để dùng cho lớp học và scheduling.
- **Phạm vi:** Quản lý giáo viên
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Quản trị viên
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Tạo teacher với mã, thông tin liên hệ, loại hình làm việc, trạng thái và ghi chú.
- **Tiền điều kiện:** Admin đã đăng nhập.
- **Điều kiện kích hoạt:** Admin submit form giáo viên.
- **Hậu điều kiện thành công:** Teacher mới được tạo.
- **Hậu điều kiện thất bại:** Không tạo teacher.
- **Dữ liệu đầu vào:** `code`, `full_name`, `email`, `phone`, `is_school_teacher`, `school_name`, `employment_type`, `status`, `notes`
- **Dữ liệu đầu ra:** Teacher detail
- **Thực thể bị tác động:** `Teacher`
- **Quy tắc nghiệp vụ áp dụng:**
  - `full_name` là bắt buộc.
  - Email unique nếu có nhập.
  - Code unique nếu có nhập.
  - `employment_type` mặc định `PART_TIME` nếu trống.
  - `status` mặc định `ACTIVE` nếu trống.
- **Luồng chính:**
  1. Admin mở form tạo giáo viên.
  2. Nhập thông tin giáo viên.
  3. Gửi request tạo.
  4. Backend validate `full_name`.
  5. Nếu có email, kiểm tra trùng email.
  6. Nếu có code, kiểm tra trùng code.
  7. Gán giá trị mặc định cho loại hình và trạng thái nếu cần.
  8. Tạo teacher và lưu DB.
  9. Trả thông tin teacher mới tạo.
- **Luồng thay thế:**
  - A1. Không nhập code: hệ thống vẫn tạo nếu schema cho phép.
  - A2. Không nhập email: hệ thống vẫn tạo.
- **Luồng ngoại lệ:**
  - E1. Thiếu `full_name`.
  - E2. Email đã tồn tại.
  - E3. Code đã tồn tại.
  - E4. Lỗi DB.
- **Ràng buộc phân quyền:** Chỉ ADMIN.
- **API / module / màn hình liên quan:** `POST /api/v1/teachers`, `TeachersPage`, `TeacherFormPage`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Có cần đồng thời tạo `User` cho teacher không? Hiện teacher là entity độc lập, chưa thấy auto tạo account.
- **Độ tin cậy:** Confirmed from code

## UC-TCH-02. Xem lịch dạy giáo viên

- **Mã use case:** UC-TCH-02
- **Tên use case:** Xem lịch dạy giáo viên
- **Mục tiêu nghiệp vụ:** Theo dõi các buổi dạy đã được sinh cho giáo viên trong một khoảng thời gian.
- **Phạm vi:** Quản lý giáo viên / lesson
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Quản trị viên hoặc giáo viên
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Trả danh sách lesson của giáo viên theo from/to.
- **Tiền điều kiện:** Giáo viên tồn tại; lesson đã được sinh.
- **Điều kiện kích hoạt:** Người dùng mở tab timetable hoặc gọi API lịch dạy.
- **Hậu điều kiện thành công:** Trả lesson list.
- **Hậu điều kiện thất bại:** Không có dữ liệu hoặc lỗi tra cứu.
- **Dữ liệu đầu vào:** `teacher_id`, `from`, `to`
- **Dữ liệu đầu ra:** danh sách lesson gồm class, room, thời gian, notes
- **Thực thể bị tác động:** Chỉ đọc `Teacher`, `Lesson`, `Class`, `Room`
- **Quy tắc nghiệp vụ áp dụng:** `teacher_id` bắt buộc; `from` và `to` phải đúng format nếu truyền.
- **Luồng chính:**
  1. Người dùng chọn giáo viên.
  2. Tùy chọn chọn khoảng ngày.
  3. Frontend gọi API lịch dạy.
  4. Backend kiểm tra teacher có tồn tại.
  5. Repository lấy lesson của teacher trong khoảng thời gian.
  6. Hệ thống map sang dữ liệu timetable và trả về.
- **Luồng thay thế:**
  - A1. Không truyền `from`/`to`: lấy theo mặc định repository hiện tại.
- **Luồng ngoại lệ:**
  - E1. Teacher ID trống.
  - E2. Sai format ngày.
  - E3. Teacher không tồn tại.
  - E4. Lỗi truy vấn lesson.
- **Ràng buộc phân quyền:** Read API đang public; theo business nên ít nhất teacher của chính mình hoặc admin.
- **API / module / màn hình liên quan:** `GET /api/v1/teachers/:id/timetable`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Có cần ràng buộc teacher chỉ xem timetable của chính mình không?
- **Độ tin cậy:** Confirmed from code

## UC-TCH-03. Xem thống kê giờ dạy

- **Mã use case:** UC-TCH-03
- **Tên use case:** Xem thống kê giờ dạy
- **Mục tiêu nghiệp vụ:** Tổng hợp thời lượng giảng dạy của giáo viên theo ngày/tuần/tháng.
- **Phạm vi:** Quản lý giáo viên / báo cáo
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Quản trị viên hoặc giáo viên
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Tính tổng số giờ dạy và breakdown theo period.
- **Tiền điều kiện:** Giáo viên tồn tại; có lesson dữ liệu.
- **Điều kiện kích hoạt:** Mở thống kê giờ dạy hoặc gọi API.
- **Hậu điều kiện thành công:** Trả total_hours và breakdown.
- **Hậu điều kiện thất bại:** Không có kết quả hoặc lỗi.
- **Dữ liệu đầu vào:** `teacher_id`, `from`, `to`, `group_by`
- **Dữ liệu đầu ra:** `total_hours`, `breakdown[]`
- **Thực thể bị tác động:** Chỉ đọc `Teacher`, `Lesson`
- **Quy tắc nghiệp vụ áp dụng:** `group_by` chỉ được là `day`, `week`, `month`.
- **Luồng chính:**
  1. Người dùng chọn giáo viên và khoảng thời gian.
  2. Chọn cách nhóm thống kê.
  3. Backend xác minh teacher tồn tại.
  4. Repository tính giờ dạy theo period.
  5. Use case cộng dồn `total_hours`.
  6. Trả kết quả.
- **Luồng thay thế:** Không truyền `group_by` thì mặc định `day`.
- **Luồng ngoại lệ:**
  - E1. `group_by` không hợp lệ.
  - E2. Teacher không tồn tại.
  - E3. Lỗi truy vấn thống kê.
- **Ràng buộc phân quyền:** Read API đang public; nên cần xem lại.
- **API / module / màn hình liên quan:** `GET /api/v1/teachers/:id/stats/teaching-hours`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Báo cáo này dùng cho chấm công hay chỉ tham khảo?
- **Độ tin cậy:** Confirmed from code

## UC-CRS-01. Tạo khóa học

- **Mã use case:** UC-CRS-01
- **Tên use case:** Tạo khóa học
- **Mục tiêu nghiệp vụ:** Khởi tạo một đơn vị đào tạo với số buổi, thời lượng và học phí làm đầu vào cho lớp học và scheduling.
- **Phạm vi:** Quản lý khóa học
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Quản trị viên
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Tạo course mới.
- **Tiền điều kiện:** Admin đăng nhập.
- **Điều kiện kích hoạt:** Submit form khóa học.
- **Hậu điều kiện thành công:** Course được tạo.
- **Hậu điều kiện thất bại:** Không tạo course.
- **Dữ liệu đầu vào:** `code`, `name`, `description`, `grade_level`, `subject`, `session_count`, `session_duration_minutes`, `total_hours`, `price`, `status`
- **Dữ liệu đầu ra:** Course detail
- **Thực thể bị tác động:** `Course`
- **Quy tắc nghiệp vụ áp dụng:**
  - `code` và `name` bắt buộc.
  - `status` mặc định `ACTIVE`.
- **Luồng chính:**
  1. Admin nhập thông tin khóa học.
  2. Gửi request tạo.
  3. Backend kiểm tra code và name.
  4. Gán `status = ACTIVE` nếu trống.
  5. Tạo và lưu course.
  6. Trả về course vừa tạo.
- **Luồng ngoại lệ:**
  - E1. Thiếu code hoặc name.
  - E2. Lỗi DB.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `POST /api/v1/courses`, `CoursePage`, `CourseDialog`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Chưa có validation chặt giữa `session_count`, `session_duration_minutes`, `total_hours`.
- **Độ tin cậy:** Confirmed from code

## UC-CRS-02. Cập nhật khóa học

- **Mã use case:** UC-CRS-02
- **Tên use case:** Cập nhật khóa học
- **Mục tiêu nghiệp vụ:** Chỉnh sửa thông tin học vụ của course.
- **Phạm vi:** Quản lý khóa học
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Quản trị viên
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Cập nhật một phần thông tin course theo ID.
- **Tiền điều kiện:** Course tồn tại.
- **Điều kiện kích hoạt:** Admin sửa khóa học.
- **Hậu điều kiện thành công:** Course được cập nhật.
- **Hậu điều kiện thất bại:** Course không đổi.
- **Dữ liệu đầu vào:** `course_id` + body cập nhật
- **Dữ liệu đầu ra:** Course updated
- **Thực thể bị tác động:** `Course`
- **Quy tắc nghiệp vụ áp dụng:** Chưa thấy validation business sâu ngoài bind.
- **Luồng chính:**
  1. Admin mở chỉnh sửa course.
  2. Sửa các trường cần thiết.
  3. Gửi request.
  4. Backend bind body và gọi update use case.
  5. Repository cập nhật course.
  6. Trả course đã cập nhật.
- **Luồng ngoại lệ:**
  - E1. Course ID trống.
  - E2. Body không hợp lệ.
  - E3. Course không tồn tại / lỗi DB.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `PUT /api/v1/courses/:id`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Khi course đang được class sử dụng, có cho sửa `session_count` hoặc `duration` không?
- **Độ tin cậy:** Confirmed from code

## UC-PRG-01. Tạo chương trình

- **Mã use case:** UC-PRG-01
- **Tên use case:** Tạo chương trình đào tạo
- **Mục tiêu nghiệp vụ:** Tạo khung chương trình để gom nhóm khóa học và định nghĩa track đào tạo.
- **Phạm vi:** Quản lý chương trình đào tạo
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Quản trị viên
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Tạo program với code, name, track, hiệu lực và ghi chú phê duyệt.
- **Tiền điều kiện:** Admin đăng nhập.
- **Điều kiện kích hoạt:** Submit form program.
- **Hậu điều kiện thành công:** Program được tạo.
- **Hậu điều kiện thất bại:** Không tạo program.
- **Dữ liệu đầu vào:** `code`, `name`, `track`, `effective_from`, `effective_to`, `approval_note`
- **Dữ liệu đầu ra:** Program detail
- **Thực thể bị tác động:** `Program`
- **Quy tắc nghiệp vụ áp dụng:** `code` và `name` bắt buộc.
- **Luồng chính:**
  1. Admin nhập thông tin chương trình.
  2. Gửi request tạo.
  3. Backend validate `code`, `name`.
  4. Tạo program và lưu DB.
  5. Trả chi tiết program.
- **Luồng ngoại lệ:**
  - E1. Thiếu code hoặc name.
  - E2. Lỗi DB.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `POST /api/v1/programs`, `ProgramPage`, `ProgramDialog`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Chưa rõ `approval_note` là ghi chú nội bộ hay bằng chứng phê duyệt thật.
- **Độ tin cậy:** Confirmed from code

## UC-PRG-02. Gán khóa học vào chương trình

- **Mã use case:** UC-PRG-02
- **Tên use case:** Gán khóa học vào chương trình
- **Mục tiêu nghiệp vụ:** Thiết lập cấu trúc chương trình bằng danh sách course thành phần.
- **Phạm vi:** Quản lý chương trình đào tạo
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Quản trị viên
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Add hoặc remove nhiều course cho một program.
- **Tiền điều kiện:** Program tồn tại; course tồn tại.
- **Điều kiện kích hoạt:** Admin chọn khóa học và xác nhận liên kết.
- **Hậu điều kiện thành công:** Mapping `program_courses` được cập nhật.
- **Hậu điều kiện thất bại:** Mapping không thay đổi.
- **Dữ liệu đầu vào:** `program_id`, `course_ids[]`
- **Dữ liệu đầu ra:** message
- **Thực thể bị tác động:** `ProgramCourse`
- **Quy tắc nghiệp vụ áp dụng:** Danh sách `course_ids` là bắt buộc trong request.
- **Luồng chính:**
  1. Admin mở chi tiết chương trình.
  2. Chọn các khóa học muốn liên kết.
  3. Frontend gửi danh sách `course_ids`.
  4. Backend gọi addCourses use case.
  5. Hệ thống tạo mapping.
  6. Trả message thành công.
- **Luồng thay thế:**
  - A1. Gỡ khóa học khỏi chương trình bằng API remove courses.
- **Luồng ngoại lệ:**
  - E1. Program ID không hợp lệ.
  - E2. Request body không có `course_ids`.
  - E3. Course không tồn tại hoặc lỗi DB.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `POST /api/v1/programs/:id/courses`, `DELETE /api/v1/programs/:id/courses`, `ProgramDetailDialog`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Có cần thứ tự khóa học trong chương trình hay không? Hiện mapping chưa thấy field ordering.
- **Độ tin cậy:** Confirmed from code

## UC-PRG-03. Cập nhật / xuất bản / lưu trữ chương trình

- **Mã use case:** UC-PRG-03
- **Tên use case:** Cập nhật / xuất bản / lưu trữ chương trình
- **Mục tiêu nghiệp vụ:** Duy trì chương trình đào tạo trong suốt vòng đời sử dụng.
- **Phạm vi:** Quản lý chương trình đào tạo
- **Mức độ ưu tiên:** Trung bình
- **Tác nhân chính:** Quản trị viên / quản lý đào tạo
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Code hiện có cập nhật program; dấu vết entity cho thấy khả năng có publish/archive, nhưng chưa có API riêng và đang mâu thuẫn với migration `status`.
- **Tiền điều kiện:** Program tồn tại.
- **Điều kiện kích hoạt:** Admin sửa thông tin program hoặc thay đổi trạng thái sử dụng.
- **Hậu điều kiện thành công:** Program được cập nhật.
- **Hậu điều kiện thất bại:** Program không đổi.
- **Dữ liệu đầu vào:** `program_id`, các field cần sửa
- **Dữ liệu đầu ra:** Program updated
- **Thực thể bị tác động:** `Program`
- **Quy tắc nghiệp vụ áp dụng:** Chưa có state machine publish/archive rõ ràng trong API.
- **Luồng chính (đã xác nhận từ code):**
  1. Admin mở form sửa program.
  2. Sửa code, name, track, hiệu lực, approval_note.
  3. Gửi request update.
  4. Backend cập nhật program.
- **Luồng thay thế (suy luận mạnh từ code):**
  - A1. Chương trình có thể được “xuất bản” bằng cách set `published_at`.
  - A2. Chương trình có thể được “lưu trữ” bằng cách set `archived_at`.
- **Luồng ngoại lệ:**
  - E1. Program không tồn tại.
  - E2. Mâu thuẫn contract giữa migration và entity làm BA khó xác định lifecycle thật.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `PUT /api/v1/programs/:id`, `ProgramPage`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Đây là use case có **Giả định / cần xác nhận BA** ở phần publish/archive.
  - Cần chốt program có `status` hay chỉ `published_at` / `archived_at`.
- **Độ tin cậy:** Mixed - update là Confirmed from code; publish/archive là Strongly inferred from code

## UC-CLS-01. Tạo lớp học

- **Mã use case:** UC-CLS-01
- **Tên use case:** Tạo lớp học
- **Mục tiêu nghiệp vụ:** Mở một lớp học thực tế để nhận học viên, gắn giáo viên, khóa học và đưa vào xếp lịch.
- **Phạm vi:** Quản lý lớp học
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Quản trị viên
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Tạo class với mã lớp, tên lớp, thời gian, sĩ số, học phí, course/program/teacher tùy chọn.
- **Tiền điều kiện:** Admin đăng nhập; khóa học/chương trình/giáo viên có thể đã tồn tại nếu muốn gắn.
- **Điều kiện kích hoạt:** Submit form class.
- **Hậu điều kiện thành công:** Class mới được tạo với trạng thái OPEN mặc định nếu không chỉ định khác.
- **Hậu điều kiện thất bại:** Không tạo class.
- **Dữ liệu đầu vào:** `code`, `name`, `notes`, `start_date`, `end_date`, `max_students`, `status`, `price`, `program_id`, `course_id`, `teacher_id`
- **Dữ liệu đầu ra:** Class detail
- **Thực thể bị tác động:** `Class`
- **Quy tắc nghiệp vụ áp dụng:**
  - `max_students >= 1`
  - `status` thuộc `OPEN`, `CLOSED`, `CANCELLED` ở level bind request
  - nếu không nhập `status` thì mặc định `OPEN`
- **Luồng chính:**
  1. Admin mở form tạo lớp.
  2. Nhập thông tin lớp.
  3. Chọn các liên kết tùy chọn: chương trình, khóa học, giáo viên.
  4. Gửi request.
  5. Backend bind request và set mặc định `OPEN` nếu cần.
  6. Tạo class và lưu DB.
  7. Trả class vừa tạo.
- **Luồng ngoại lệ:**
  - E1. Thiếu field bắt buộc như code, name, start_date, max_students.
  - E2. `max_students` không hợp lệ.
  - E3. Lỗi DB.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `POST /api/v1/classes`, `ClassesPage`, `ClassDialog`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Chưa thấy bước nhập `room_id` ở DTO hiện tại dù entity class có `room_id`.
- **Độ tin cậy:** Confirmed from code

## UC-CLS-02. Cập nhật lớp học

- **Mã use case:** UC-CLS-02
- **Tên use case:** Cập nhật lớp học
- **Mục tiêu nghiệp vụ:** Điều chỉnh thông tin lớp trong quá trình vận hành.
- **Phạm vi:** Quản lý lớp học
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Quản trị viên
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Sửa class theo ID.
- **Tiền điều kiện:** Class tồn tại.
- **Điều kiện kích hoạt:** Admin sửa lớp.
- **Hậu điều kiện thành công:** Class updated.
- **Hậu điều kiện thất bại:** Class không đổi.
- **Dữ liệu đầu vào:** `class_id` + dữ liệu cập nhật
- **Dữ liệu đầu ra:** Class updated
- **Thực thể bị tác động:** `Class`
- **Quy tắc nghiệp vụ áp dụng:** Tương tự UC-CLS-01.
- **Luồng chính:**
  1. Admin mở chi tiết hoặc form sửa lớp.
  2. Chỉnh sửa thông tin.
  3. Gửi request update.
  4. Backend bind body.
  5. Use case update class.
  6. Trả class đã cập nhật.
- **Luồng ngoại lệ:**
  - E1. Class không tồn tại.
  - E2. Body invalid.
  - E3. Lỗi DB update.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `PUT /api/v1/classes/:id`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Có nên chặn sửa lớp đã `CLOSED` hoặc `CANCELLED` không?
- **Độ tin cậy:** Confirmed from code

## UC-ENR-01. Ghi danh học viên vào lớp

- **Mã use case:** UC-ENR-01
- **Tên use case:** Ghi danh học viên vào lớp
- **Mục tiêu nghiệp vụ:** Thêm một hoặc nhiều học viên vào roster lớp.
- **Phạm vi:** Ghi danh
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Quản trị viên / vận hành
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Nhận `student_ids[]`, kiểm tra giới hạn sức chứa và tạo enrollments.
- **Tiền điều kiện:** Class tồn tại; student tồn tại.
- **Điều kiện kích hoạt:** Người dùng xác nhận thêm học viên vào lớp.
- **Hậu điều kiện thành công:** Enrollment records được tạo; roster lớp tăng lên.
- **Hậu điều kiện thất bại:** Không tạo enrollment mới.
- **Dữ liệu đầu vào:** `class_id`, `student_ids[]`
- **Dữ liệu đầu ra:** Danh sách enrollment hoặc output use case
- **Thực thể bị tác động:** `Enrollment`
- **Quy tắc nghiệp vụ áp dụng:**
  - Danh sách student phải có ít nhất 1 phần tử.
  - `maxAllowed = min(class.max_students, room.capacity nếu class có room)`
  - Hiện use case không đếm enrollment đã tồn tại trong DB.
  - Hiện use case gán `status = ENROLLED` trực tiếp.
- **Luồng chính:**
  1. Người dùng mở chi tiết lớp.
  2. Chọn nhiều học viên từ danh sách khả dụng.
  3. Frontend gửi `student_ids`.
  4. Backend lấy class.
  5. Backend xác định giới hạn sức chứa thực tế.
  6. Backend kiểm tra số lượng đầu vào không vượt giới hạn đang tính.
  7. Backend tạo enrollment cho từng student với trạng thái `ENROLLED`.
  8. Trả kết quả thành công.
- **Luồng thay thế:**
  - A1. Class không có phòng: giới hạn lấy theo `max_students`.
- **Luồng ngoại lệ:**
  - E1. Class không tồn tại.
  - E2. Danh sách student rỗng.
  - E3. Số lượng student đầu vào vượt `maxAllowed`.
  - E4. Lỗi DB khi tạo enrollment.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `POST /api/v1/classes/:id/students`, `ClassDetailDialog`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Đây là use case **có gap nghiệp vụ quan trọng**: chưa tính roster hiện hữu, chưa chống duplicate, chưa có approval flow.
  - **Độ tin cậy:** Confirmed from code

## UC-ENR-02. Rút học viên khỏi lớp

- **Mã use case:** UC-ENR-02
- **Tên use case:** Rút học viên khỏi lớp
- **Mục tiêu nghiệp vụ:** Loại học viên ra khỏi roster lớp.
- **Phạm vi:** Ghi danh
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Quản trị viên / vận hành
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Xóa enrollment tương ứng giữa class và student.
- **Tiền điều kiện:** Class tồn tại; enrollments tương ứng tồn tại.
- **Điều kiện kích hoạt:** Người dùng bấm xóa học viên khỏi lớp.
- **Hậu điều kiện thành công:** Enrollment bị xóa.
- **Hậu điều kiện thất bại:** Roster giữ nguyên.
- **Dữ liệu đầu vào:** `class_id`, `student_ids[]`
- **Dữ liệu đầu ra:** message
- **Thực thể bị tác động:** `Enrollment`
- **Quy tắc nghiệp vụ áp dụng:** Hiện không kiểm tra trạng thái enrollment.
- **Luồng chính:**
  1. Người dùng chọn học viên cần loại khỏi lớp.
  2. Gửi request remove.
  3. Backend lấy danh sách enrollments theo class.
  4. Backend tìm enrollment khớp `student_id`.
  5. Xóa enrollment tương ứng.
  6. Trả thông báo thành công.
- **Luồng ngoại lệ:**
  - E1. Request body không hợp lệ.
  - E2. Lỗi truy vấn hoặc xóa DB.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `DELETE /api/v1/classes/:id/students`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Có cần lịch sử rút lớp thay vì xóa mềm mapping hay không?
- **Độ tin cậy:** Confirmed from code

## UC-CLS-03. Phân công giáo viên cho lớp

- **Mã use case:** UC-CLS-03
- **Tên use case:** Phân công giáo viên cho lớp
- **Mục tiêu nghiệp vụ:** Gắn giáo viên phụ trách chính cho lớp để phục vụ vận hành và xếp lịch.
- **Phạm vi:** Quản lý lớp học
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Quản trị viên
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Cập nhật `teacher_id` trên lớp.
- **Tiền điều kiện:** Class tồn tại; giáo viên tồn tại theo business.
- **Điều kiện kích hoạt:** Admin chọn giáo viên và xác nhận.
- **Hậu điều kiện thành công:** Class gắn với teacher mới.
- **Hậu điều kiện thất bại:** Class giữ teacher cũ.
- **Dữ liệu đầu vào:** `class_id`, `teacher_id`
- **Dữ liệu đầu ra:** class updated / message
- **Thực thể bị tác động:** `Class`
- **Quy tắc nghiệp vụ áp dụng:** Theo route là ADMIN; chưa thấy validate workload hay teacher availability tại bước gán.
- **Luồng chính:**
  1. Admin mở lớp.
  2. Chọn giáo viên phụ trách.
  3. Gửi request assign teacher.
  4. Backend update `teacher_id`.
  5. Trả thành công.
- **Luồng ngoại lệ:**
  - E1. `teacher_id` rỗng.
  - E2. Class không tồn tại.
  - E3. Lỗi DB.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `PUT /api/v1/classes/:id/teacher`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Có cần phân công nhiều giáo viên/co-teacher không?
- **Độ tin cậy:** Confirmed from code

## UC-SCHCFG-01. Cấu hình lịch tuần cho lớp

- **Mã use case:** UC-SCHCFG-01
- **Tên use case:** Cấu hình lịch tuần cho lớp
- **Mục tiêu nghiệp vụ:** Xác định các ngày trong tuần và ca học mà lớp được phép học, làm đầu vào cho solver scheduling.
- **Phạm vi:** Lớp học / scheduling
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Quản trị viên / vận hành đào tạo
- **Tác nhân phụ:** Không có
- **Mô tả ngắn:** Về mặt dữ liệu, hệ thống dùng `class_schedules` gồm `class_id`, `day_of_week`, `shift_id`, optional `room_id`. Tuy nhiên hiện chưa có API/controller/màn hình rõ ràng cho use case này.
- **Tiền điều kiện:** Class đã tồn tại; Shift đã được tạo.
- **Điều kiện kích hoạt:** Người vận hành cần khai báo lịch chuẩn cho lớp trước khi xếp lịch.
- **Hậu điều kiện thành công:** Class có tối thiểu một bản ghi `class_schedule`.
- **Hậu điều kiện thất bại:** Scheduling có thể không tạo được domain/slot phù hợp.
- **Dữ liệu đầu vào:** `class_id`, `day_of_week`, `shift_id`, `room_id?`
- **Dữ liệu đầu ra:** ClassSchedule records
- **Thực thể bị tác động:** `ClassSchedule`
- **Quy tắc nghiệp vụ áp dụng:**
  - Mỗi lịch tuần phải tham chiếu một `shift_id`.
  - `room_id` trong lịch tuần là tùy chọn.
  - Scheduling preview phụ thuộc rất mạnh vào dữ liệu này.
- **Luồng chính (mong muốn nghiệp vụ):**
  1. Admin mở cấu hình lịch tuần của lớp.
  2. Chọn các ngày trong tuần.
  3. Với mỗi ngày, chọn ca học (`shift_id`).
  4. Có thể gắn phòng cố định nếu cần.
  5. Lưu class schedules.
  6. Hệ thống dùng dữ liệu này cho scheduling preview.
- **Luồng thay thế:**
  - A1. Nếu lớp không có class_schedule, solver có thể fallback sang slot từ danh sách shift active chung.
- **Luồng ngoại lệ:**
  - E1. Không có shift active hoặc shift không hợp lệ.
  - E2. Dữ liệu class_schedule không được tạo do hiện chưa có API quản trị chính thức.
- **Ràng buộc phân quyền:** Theo nghiệp vụ nên là Admin/Vận hành đào tạo.
- **API / module / màn hình liên quan:** Chưa có controller/use case riêng; được preload trong scheduling preview.
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Đây là use case **Giả định / cần xác nhận BA** về mặt UI/API.
  - Cần chốt nơi cấu hình class_schedule trong sản phẩm.
- **Độ tin cậy:** Data model Confirmed from code; workflow/UI Assumption / Needs BA validation

## UC-SOL-01. Tạo preview xếp lịch

- **Mã use case:** UC-SOL-01
- **Tên use case:** Tạo preview xếp lịch
- **Mục tiêu nghiệp vụ:** Sinh phương án lịch học khả thi trước khi ghi lesson thật vào hệ thống.
- **Phạm vi:** Xếp lịch thông minh
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Quản trị viên
- **Tác nhân phụ:** Scheduling engine
- **Mô tả ngắn:** Người dùng chọn khoảng ngày và bộ lọc; hệ thống load lớp OPEN, phòng, ca học và chạy solver để sinh preview gồm assignments, conflicts và summary.
- **Tiền điều kiện:** Có shift active; có dữ liệu lớp/phòng phù hợp; khoảng ngày hợp lệ.
- **Điều kiện kích hoạt:** Admin bấm “Chạy xếp lịch”.
- **Hậu điều kiện thành công:** Preview được lưu vào PreviewStore; hiển thị được kết quả.
- **Hậu điều kiện thất bại:** Preview status FAILED hoặc API trả lỗi.
- **Dữ liệu đầu vào:** `date_from`, `date_to`, `class_ids[]`, `teacher_ids[]`, `room_ids[]`
- **Dữ liệu đầu ra:** `PreviewResult` gồm `run_id`, `status`, `assignments[]`, `conflicts[]`, `summary`
- **Thực thể bị tác động:** Không cập nhật entity lõi; chỉ tác động `PreviewStore` in-memory
- **Quy tắc nghiệp vụ áp dụng:**
  - `date_to >= date_from`
  - Chỉ load `Class.status = OPEN`
  - Chỉ load `Shift.is_active = true`
  - Có conflict tổng quan nếu không có class, room hoặc shift phù hợp
  - Solver dùng abstraction `SchedulingSolver`
- **Luồng chính:**
  1. Admin mở màn Scheduling.
  2. Chọn khoảng ngày preview.
  3. Tùy chọn chọn lớp, giáo viên, phòng để giới hạn phạm vi.
  4. Gửi request preview.
  5. Backend validate khoảng ngày.
  6. Backend load classes OPEN và preload teacher/course/room/class_schedules/shift.
  7. Backend load rooms theo filter.
  8. Backend load active shifts.
  9. Backend tạo `PreviewResult` base.
  10. Backend thêm conflict tổng quan nếu thiếu class/room/shift.
  11. Backend gọi solver service.
  12. Nhận output từ solver gồm assignments/conflicts/summary/status.
  13. Backend lưu preview vào PreviewStore bằng `run_id`.
  14. Trả preview cho frontend.
- **Luồng thay thế:**
  - A1. Không chọn class_ids: hệ thống xét toàn bộ lớp OPEN.
  - A2. Không chọn room_ids: hệ thống xét toàn bộ phòng tải được.
  - A3. Không có class_schedule: solver có thể dùng shift chung làm slot fallback theo implementation hiện tại.
- **Luồng ngoại lệ:**
  - E1. `date_to < date_from`: trả lỗi 400 hoặc preview FAILED.
  - E2. Lỗi truy vấn classes/rooms/shifts.
  - E3. Solver báo không thể sinh preview.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `POST /api/v1/scheduling/preview`, `SchedulingPage`, `internal/usecases/scheduling/preview.go`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Hiện benchmark chưa chọn solver thắng; preview chính vẫn đi qua solver đang được inject mặc định.
  - **Độ tin cậy:** Confirmed from code

## UC-SOL-02. Xem preview

- **Mã use case:** UC-SOL-02
- **Tên use case:** Xem preview xếp lịch
- **Mục tiêu nghiệp vụ:** Cho phép admin rà soát kết quả xếp lịch trước khi commit.
- **Phạm vi:** Xếp lịch thông minh
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Quản trị viên
- **Tác nhân phụ:** Scheduling engine
- **Mô tả ngắn:** Lấy preview theo `run_id` hoặc preview gần nhất.
- **Tiền điều kiện:** Đã từng tạo preview trong tiến trình backend hiện tại.
- **Điều kiện kích hoạt:** Admin mở lại preview hoặc refresh trang scheduling.
- **Hậu điều kiện thành công:** Hiển thị assignments, conflicts, summary.
- **Hậu điều kiện thất bại:** Không tìm thấy preview.
- **Dữ liệu đầu vào:** `run_id` hoặc none với latest
- **Dữ liệu đầu ra:** `PreviewResult`
- **Thực thể bị tác động:** Chỉ đọc `PreviewStore`
- **Quy tắc nghiệp vụ áp dụng:** Preview latest đang là dữ liệu in-memory, mất khi restart backend.
- **Luồng chính:**
  1. Admin chọn xem lại preview theo run_id hoặc latest.
  2. Frontend gọi API tương ứng.
  3. Backend tra PreviewStore.
  4. Nếu có preview thì trả dữ liệu đầy đủ.
  5. Frontend hiển thị cột assignments, conflicts, progress, severity.
- **Luồng thay thế:**
  - A1. Lấy latest preview thay vì run_id cụ thể.
- **Luồng ngoại lệ:**
  - E1. Không tìm thấy preview run.
  - E2. Backend vừa restart nên latest preview bị mất.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `GET /api/v1/scheduling/preview/:id`, `GET /api/v1/scheduling/preview/latest`, `SchedulingPage`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Có cần lưu lịch sử preview vào DB cho audit không?
- **Độ tin cậy:** Confirmed from code

## UC-SOL-03. Benchmark solver

- **Mã use case:** UC-SOL-03
- **Tên use case:** Benchmark solver
- **Mục tiêu nghiệp vụ:** So sánh các solver scheduling để chọn thuật toán phù hợp nhất cho sản phẩm hoặc cho báo cáo nghiên cứu.
- **Phạm vi:** Xếp lịch thông minh
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Quản trị viên / người làm nghiên cứu giải thuật
- **Tác nhân phụ:** Scheduling engine
- **Mô tả ngắn:** Endpoint benchmark chạy thật 3 solver `graph_coloring`, `cp_sat`, `tabu_search` trên cùng input để trả metric so sánh và hỗ trợ chọn solver chính.
- **Tiền điều kiện:** Admin đăng nhập.
- **Điều kiện kích hoạt:** Admin gọi benchmark API.
- **Hậu điều kiện thành công:** Nhận bảng so sánh solver gồm feasibility, hard violations, soft score, runtime và summary.
- **Hậu điều kiện thất bại:** Không lấy được dữ liệu benchmark hoặc có solver lỗi khi thực thi.
- **Dữ liệu đầu vào:** `date_from`, `date_to`, `class_ids[]`, `teacher_ids[]`, `room_ids[]`
- **Dữ liệu đầu ra:** `BenchmarkOutput` gồm generatedAt, filters, mode, solvers[]
- **Thực thể bị tác động:** Không cập nhật entity domain
- **Quy tắc nghiệp vụ áp dụng:** Benchmark dùng cùng filter với preview, không cập nhật entity domain, chỉ phục vụ ADMIN nội bộ/nghiên cứu.
- **Luồng chính:**
  1. Admin gửi yêu cầu benchmark với cùng bộ filter như preview.
  2. Backend validate khoảng ngày.
  3. Use case nạp dữ liệu benchmark và lấy danh mục solver từ `SolverCatalog`.
  4. Hệ thống chạy lần lượt 3 solver trên cùng input.
  5. Hệ thống tổng hợp feasibility, hard violations, soft score, runtime và summary.
  6. Trả về bảng benchmark để so sánh solver.
- **Luồng thay thế:** Không có.
- **Luồng ngoại lệ:**
  - E1. Request body invalid.
  - E2. date range invalid.
  - E3. Lỗi thực thi benchmark hoặc solver không khả dụng.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `POST /api/v1/scheduling/benchmark`, `internal/usecases/scheduling/benchmark.go`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Quyết định solver chính đang được tài liệu hóa tại `docs/SCHEDULING_BENCHMARK_REPORT_2026-04-14.md`.
  - **Độ tin cậy:** Confirmed from code

## UC-SOL-04. Xác nhận preview để tạo lesson

- **Mã use case:** UC-SOL-04
- **Tên use case:** Xác nhận preview để tạo lesson
- **Mục tiêu nghiệp vụ:** Biến kết quả preview thành lịch buổi học thực tế trong DB.
- **Phạm vi:** Scheduling / lesson generation
- **Mức độ ưu tiên:** Rất cao
- **Tác nhân chính:** Quản trị viên
- **Tác nhân phụ:** Scheduling engine
- **Mô tả ngắn:** Chỉ preview có trạng thái `COMPLETED` và không còn unscheduled/conflict mới được commit thành lesson.
- **Tiền điều kiện:** Preview tồn tại trong store; preview đã hoàn tất hợp lệ.
- **Điều kiện kích hoạt:** Admin bấm commit.
- **Hậu điều kiện thành công:** Tạo `Lesson` records tương ứng với assignments.
- **Hậu điều kiện thất bại:** Không tạo lesson nào.
- **Dữ liệu đầu vào:** `run_id`
- **Dữ liệu đầu ra:** `message`, `scheduled_lessons`, `status=COMMITTED`
- **Thực thể bị tác động:** `Lesson`
- **Quy tắc nghiệp vụ áp dụng:**
  - `run_id` bắt buộc.
  - Preview phải tồn tại.
  - Preview phải có assignment.
  - Preview phải `COMPLETED`.
  - Không được còn unscheduled/conflict theo điều kiện preview.
  - Không được trùng lesson đã có theo class, teacher hoặc room.
- **Luồng chính:**
  1. Admin chọn preview đã rà soát.
  2. Frontend gửi `run_id`.
  3. Backend lấy preview từ store.
  4. Backend kiểm tra preview có assignments.
  5. Backend kiểm tra `status == COMPLETED`.
  6. Backend tính khoảng thời gian của preview.
  7. Backend tải các lesson hiện hữu có khả năng overlap.
  8. Backend kiểm tra conflict commit theo class/teacher/room.
  9. Nếu không conflict, backend mở transaction.
  10. Tạo `Lesson` cho từng assignment.
  11. Commit transaction.
  12. Trả kết quả.
- **Luồng thay thế:** Không có nhánh hợp lệ khác.
- **Luồng ngoại lệ:**
  - E1. `run_id` rỗng.
  - E2. Không tìm thấy preview.
  - E3. Preview không có assignment.
  - E4. Preview chưa `COMPLETED`.
  - E5. Conflict với lesson hiện hữu.
  - E6. Lỗi transaction tạo lesson.
- **Ràng buộc phân quyền:** ADMIN.
- **API / module / màn hình liên quan:** `POST /api/v1/scheduling/commit`, `SchedulingPage`, `commit_preview.go`
- **Ghi chú / điểm mơ hồ cần xác nhận:** Chưa có cờ “preview đã commit”; cần xác nhận có cho re-commit cùng preview không.
- **Độ tin cậy:** Confirmed from code

## UC-LSN-01. Xem danh sách lesson

- **Mã use case:** UC-LSN-01
- **Tên use case:** Xem danh sách lesson
- **Mục tiêu nghiệp vụ:** Tra cứu các buổi học đã sinh để phục vụ vận hành giảng dạy và theo dõi học tập.
- **Phạm vi:** Lesson
- **Mức độ ưu tiên:** Trung bình
- **Tác nhân chính:** Quản trị viên / giáo viên
- **Tác nhân phụ:** Có thể là học viên trong tương lai
- **Mô tả ngắn:** Code hiện có entity `Lesson` và các use case thống kê/timetable theo lesson, nhưng chưa có controller/màn lesson list tổng quát.
- **Tiền điều kiện:** Lessons đã được tạo từ scheduling commit.
- **Điều kiện kích hoạt:** Người dùng cần tra danh sách lesson.
- **Hậu điều kiện thành công:** Hiển thị được lesson list theo nhu cầu.
- **Hậu điều kiện thất bại:** Không tra cứu được.
- **Dữ liệu đầu vào:** các bộ lọc dự kiến như class, teacher, from, to
- **Dữ liệu đầu ra:** lesson list
- **Thực thể bị tác động:** `Lesson`
- **Quy tắc nghiệp vụ áp dụng:** Chưa có API chính thức.
- **Luồng chính (mong muốn):**
  1. Người dùng mở màn danh sách lesson.
  2. Chọn bộ lọc.
  3. Hệ thống truy vấn lesson.
  4. Hiển thị lesson cùng class, room, teacher, start/end.
- **Luồng ngoại lệ:** Use case chưa có implementation riêng.
- **Ràng buộc phân quyền:** Theo nghiệp vụ nên là Admin/Teacher; có thể Student chỉ xem lesson của mình.
- **API / module / màn hình liên quan:** Chưa có API/controller; hiện có gián tiếp qua teacher timetable.
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Đây là use case **Giả định / cần xác nhận BA** ở mức vận hành.
  - **Độ tin cậy:** Strongly inferred from code

## UC-ATD-01. Điểm danh

- **Mã use case:** UC-ATD-01
- **Tên use case:** Điểm danh
- **Mục tiêu nghiệp vụ:** Ghi nhận tình trạng tham dự của học viên theo buổi học.
- **Phạm vi:** Học thuật
- **Mức độ ưu tiên:** Trung bình
- **Tác nhân chính:** Giáo viên
- **Tác nhân phụ:** Quản trị viên
- **Mô tả ngắn:** Hệ thống có entity `Attendance` nhưng chưa có API và chưa định nghĩa bảng mã trạng thái rõ ràng.
- **Tiền điều kiện:** Lesson đã tồn tại; roster lớp đã xác định.
- **Điều kiện kích hoạt:** Giáo viên muốn chấm điểm danh sau hoặc trong buổi học.
- **Hậu điều kiện thành công:** Attendance records được tạo/cập nhật.
- **Hậu điều kiện thất bại:** Không lưu điểm danh.
- **Dữ liệu đầu vào:** dự kiến `lesson_id`, `student_id`, `status`, `note`, `marked_at`
- **Dữ liệu đầu ra:** attendance result
- **Thực thể bị tác động:** `Attendance`
- **Quy tắc nghiệp vụ áp dụng:** Chưa xác định được enum `status` từ code.
- **Luồng chính (mong muốn):**
  1. Giáo viên mở danh sách lesson.
  2. Chọn lesson cần điểm danh.
  3. Hệ thống tải roster của lesson/lớp.
  4. Giáo viên chọn trạng thái cho từng học viên.
  5. Lưu attendance records.
- **Luồng thay thế:** Admin sửa điểm danh thay cho giáo viên.
- **Luồng ngoại lệ:** Chưa có API/use case.
- **Ràng buộc phân quyền:** Theo nghiệp vụ nên là Teacher/Admin.
- **API / module / màn hình liên quan:** Placeholder `teacher/attendance`; chưa có backend API.
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Cần chốt trạng thái attendance: có mặt, vắng, có phép, đi muộn, về sớm...
  - **Độ tin cậy:** Entity Confirmed from code; workflow Assumption / Needs BA validation

## UC-SUM-01. Tạo tổng kết buổi học

- **Mã use case:** UC-SUM-01
- **Tên use case:** Tạo tổng kết buổi học
- **Mục tiêu nghiệp vụ:** Ghi nhận nội dung giảng dạy, bài tập và phản hồi sau buổi học.
- **Phạm vi:** Học thuật
- **Mức độ ưu tiên:** Trung bình
- **Tác nhân chính:** Giáo viên
- **Tác nhân phụ:** Quản trị viên
- **Mô tả ngắn:** Hệ thống có entity `LessonSummary` nhưng chưa có API/controller/use case.
- **Tiền điều kiện:** Lesson tồn tại.
- **Điều kiện kích hoạt:** Giáo viên kết thúc buổi học và cần tổng kết.
- **Hậu điều kiện thành công:** Lesson summary được tạo.
- **Hậu điều kiện thất bại:** Không có bản tổng kết.
- **Dữ liệu đầu vào:** `lesson_id`, `topic`, `lesson_content`, `class_feedback`, `homework`, `homework_deadline`, `teacher_notes`
- **Dữ liệu đầu ra:** LessonSummary detail
- **Thực thể bị tác động:** `LessonSummary`
- **Quy tắc nghiệp vụ áp dụng:** Mỗi lesson có tối đa một summary theo thiết kế entity.
- **Luồng chính (mong muốn):**
  1. Giáo viên mở lesson.
  2. Nhập nội dung đã dạy, phản hồi lớp, bài tập và hạn nộp.
  3. Lưu summary.
  4. Hệ thống gắn summary với lesson.
- **Luồng ngoại lệ:** Chưa có API.
- **Ràng buộc phân quyền:** Theo nghiệp vụ nên là Teacher/Admin.
- **API / module / màn hình liên quan:** Chưa có API; teacher journal page hiện là placeholder.
- **Ghi chú / điểm mơ hồ cần xác nhận:** Có bắt buộc summary trước khi ghi academic record không?
- **Độ tin cậy:** Entity Confirmed from code; workflow Assumption / Needs BA validation

## UC-ACR-01. Ghi nhận academic record

- **Mã use case:** UC-ACR-01
- **Tên use case:** Ghi nhận academic record
- **Mục tiêu nghiệp vụ:** Đánh giá kết quả học tập của từng học viên sau một buổi học hoặc một lần tổng kết.
- **Phạm vi:** Học thuật / đánh giá
- **Mức độ ưu tiên:** Trung bình
- **Tác nhân chính:** Giáo viên
- **Tác nhân phụ:** Quản trị viên
- **Mô tả ngắn:** Có entity `AcademicRecord` nhưng chưa có API/use case.
- **Tiền điều kiện:** Có lesson summary; student thuộc lớp tương ứng.
- **Điều kiện kích hoạt:** Giáo viên cần ghi nhận đánh giá học viên.
- **Hậu điều kiện thành công:** Academic record được tạo/cập nhật.
- **Hậu điều kiện thất bại:** Không có record mới.
- **Dữ liệu đầu vào:** `lesson_summary_id`, `student_id`, `homework_completed`, `homework_score`, `attitude_rating`, `participation_score`, `personal_comment`, `total_score`, `is_completed`
- **Dữ liệu đầu ra:** AcademicRecord detail
- **Thực thể bị tác động:** `AcademicRecord`
- **Quy tắc nghiệp vụ áp dụng:** Chưa có quy tắc tính `total_score` hoặc `is_completed`.
- **Luồng chính (mong muốn):**
  1. Giáo viên mở lesson summary.
  2. Chọn học viên cần đánh giá.
  3. Nhập điểm và nhận xét.
  4. Hệ thống lưu academic record.
- **Luồng ngoại lệ:** Chưa có API.
- **Ràng buộc phân quyền:** Theo nghiệp vụ nên là Teacher/Admin.
- **API / module / màn hình liên quan:** Chưa có.
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Có cần tính điểm tự động hay giáo viên nhập tay toàn bộ?
  - `is_completed` được set theo rule nào?
  - **Độ tin cậy:** Entity Confirmed from code; workflow Assumption / Needs BA validation

## UC-LVE-01. Tạo đơn xin phép

- **Mã use case:** UC-LVE-01
- **Tên use case:** Tạo đơn xin phép
- **Mục tiêu nghiệp vụ:** Cho học viên gửi yêu cầu nghỉ học, đi muộn hoặc về sớm.
- **Phạm vi:** Đơn xin phép
- **Mức độ ưu tiên:** Trung bình
- **Tác nhân chính:** Học viên
- **Tác nhân phụ:** Quản trị viên / giáo vụ
- **Mô tả ngắn:** Có entity `LeaveRequest` với `leave_type`, `status`, `approved_by`, nhưng chưa có API triển khai.
- **Tiền điều kiện:** Student tồn tại; có thể gắn class hoặc lesson tương ứng.
- **Điều kiện kích hoạt:** Học viên cần báo nghỉ hoặc điều chỉnh tham gia buổi học.
- **Hậu điều kiện thành công:** LeaveRequest được tạo ở trạng thái chờ duyệt.
- **Hậu điều kiện thất bại:** Không tạo đơn.
- **Dữ liệu đầu vào:** `student_id`, `leave_type`, `apply_date`, `late_minutes`, `early_minutes`, `reason`, `documents[]`, `class_id?`, `lesson_id?`, `subject`
- **Dữ liệu đầu ra:** LeaveRequest detail
- **Thực thể bị tác động:** `LeaveRequest`
- **Quy tắc nghiệp vụ áp dụng:** `leave_type` có dấu vết comment gồm `LEAVE`, `LATE`, `EARLY`; trạng thái mặc định `PENDING`.
- **Luồng chính (mong muốn):**
  1. Học viên mở form xin phép.
  2. Chọn loại yêu cầu: nghỉ, đi muộn hoặc về sớm.
  3. Nhập lý do, ngày áp dụng, các phút muộn/về sớm nếu có.
  4. Có thể gắn lớp/buổi học liên quan.
  5. Gửi đơn.
  6. Hệ thống lưu `LeaveRequest` trạng thái `PENDING`.
- **Luồng thay thế:** Có thể upload chứng từ đính kèm.
- **Luồng ngoại lệ:** Chưa có API hiện thực.
- **Ràng buộc phân quyền:** Theo nghiệp vụ nên là Student hoặc phụ huynh thay mặt.
- **API / module / màn hình liên quan:** Chưa có backend; Student portal có placeholder leaves.
- **Ghi chú / điểm mơ hồ cần xác nhận:** Có cần phụ huynh là actor gửi đơn thay không?
- **Độ tin cậy:** Entity Confirmed from code; workflow Assumption / Needs BA validation

## UC-LVE-02. Duyệt / từ chối đơn xin phép

- **Mã use case:** UC-LVE-02
- **Tên use case:** Duyệt / từ chối đơn xin phép
- **Mục tiêu nghiệp vụ:** Quyết định cho phép nghỉ hoặc điều chỉnh chuyên cần của học viên.
- **Phạm vi:** Đơn xin phép
- **Mức độ ưu tiên:** Trung bình
- **Tác nhân chính:** Quản trị viên / giáo vụ
- **Tác nhân phụ:** Học viên
- **Mô tả ngắn:** Entity `LeaveRequest` có `approved_by_id`, `approved_at`, `rejection_reason`, nhưng chưa có API.
- **Tiền điều kiện:** Có LeaveRequest ở trạng thái chờ duyệt.
- **Điều kiện kích hoạt:** Người duyệt mở danh sách đơn.
- **Hậu điều kiện thành công:** Đơn được approve hoặc reject.
- **Hậu điều kiện thất bại:** Trạng thái đơn giữ nguyên.
- **Dữ liệu đầu vào:** `leave_request_id`, quyết định, lý do từ chối
- **Dữ liệu đầu ra:** LeaveRequest updated
- **Thực thể bị tác động:** `LeaveRequest`
- **Quy tắc nghiệp vụ áp dụng:** Chưa có state machine chính thức ngoài dấu vết field.
- **Luồng chính (mong muốn):**
  1. Người duyệt xem danh sách đơn chờ.
  2. Mở chi tiết đơn.
  3. Chọn approve hoặc reject.
  4. Nếu reject thì nhập lý do.
  5. Hệ thống cập nhật trạng thái và người duyệt.
- **Luồng thay thế:** Có thể đồng bộ sang attendance nếu đơn gắn với lesson.
- **Luồng ngoại lệ:** Chưa có API hiện thực.
- **Ràng buộc phân quyền:** Theo nghiệp vụ nên là Admin/giáo vụ.
- **API / module / màn hình liên quan:** Chưa có.
- **Ghi chú / điểm mơ hồ cần xác nhận:** Có cần nhiều cấp duyệt hay không?
- **Độ tin cậy:** Entity Confirmed from code; workflow Assumption / Needs BA validation

## UC-MAT-01. Upload tài liệu

- **Mã use case:** UC-MAT-01
- **Tên use case:** Upload tài liệu
- **Mục tiêu nghiệp vụ:** Cho giáo viên đưa tài liệu giảng dạy lên hệ thống và chạy kiểm tra nội dung tự động.
- **Phạm vi:** Tài liệu giảng dạy và kiểm duyệt
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Giáo viên
- **Tác nhân phụ:** OCR stub, AI moderation stub, file storage
- **Mô tả ngắn:** Giáo viên upload file, hệ thống lưu file local, tạo material, chạy OCR/AI stub và gán nhãn.
- **Tiền điều kiện:** Giáo viên đăng nhập; file đúng định dạng.
- **Điều kiện kích hoạt:** Giáo viên gửi form upload.
- **Hậu điều kiện thành công:** Material được lưu; audit log được tạo; tài liệu chuyển sang `AI_REVIEWED`.
- **Hậu điều kiện thất bại:** Không có material hợp lệ hoặc không có audit log.
- **Dữ liệu đầu vào:** multipart gồm `teacher_id`, `title`, `description`, `file`
- **Dữ liệu đầu ra:** Material view
- **Thực thể bị tác động:** `Material`, `AuditLog`, `Label`
- **Quy tắc nghiệp vụ áp dụng:**
  - Chỉ hỗ trợ pdf/doc/docx/png/jpg/jpeg.
  - File tối đa 10 MB.
  - Material được lưu trong local storage theo path cấu trúc.
  - OCR/AI hiện là stub.
- **Luồng chính:**
  1. Giáo viên chọn giáo viên, nhập tiêu đề và mô tả.
  2. Chọn file tải lên.
  3. Frontend kiểm tra client-side về kích thước và mime type.
  4. Gửi multipart request.
  5. Backend đọc file bytes.
  6. Backend validate file type và file size.
  7. Backend tạo material trạng thái `SCANNING`.
  8. Backend lưu file vào storage local.
  9. Backend chạy OCR stub lấy raw text.
  10. Backend chạy AI stub gán nhãn SAFE/WARNING/DANGER.
  11. Backend tạo audit log.
  12. Backend cập nhật material sang `AI_REVIEWED` và gắn `latest_label_id`.
  13. Trả material detail.
- **Luồng thay thế:**
  - A1. Nếu AI gán SAFE, material vẫn được đưa vào list; có thể không cần review thủ công tùy nghiệp vụ.
- **Luồng ngoại lệ:**
  - E1. Không có file trong request.
  - E2. File không đọc được.
  - E3. Sai định dạng file.
  - E4. File quá lớn.
  - E5. Lỗi lưu file hoặc lưu DB.
- **Ràng buộc phân quyền:** Route yêu cầu role TEACHER.
- **API / module / màn hình liên quan:** `POST /api/v1/materials/upload`, `TeacherDocumentsPage`, `upload_material.go`, `services/audit/services.go`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - `teacher_id` hiện lấy từ form, chưa derive từ token.
  - Ownership tài liệu cần BA xác nhận lại.
  - **Độ tin cậy:** Confirmed from code

## UC-MAT-02. Duyệt tài liệu bị gắn cờ

- **Mã use case:** UC-MAT-02
- **Tên use case:** Duyệt tài liệu bị gắn cờ
- **Mục tiêu nghiệp vụ:** Cho phép con người ra quyết định cuối cùng với tài liệu đã được AI gắn nhãn.
- **Phạm vi:** Tài liệu giảng dạy và kiểm duyệt
- **Mức độ ưu tiên:** Cao
- **Tác nhân chính:** Compliance officer / quản trị viên
- **Tác nhân phụ:** AI stub, giáo viên
- **Mô tả ngắn:** Người duyệt xem reasoning và nhãn AI, rồi quyết định approve hoặc reject tài liệu.
- **Tiền điều kiện:** Tồn tại material `AI_REVIEWED`.
- **Điều kiện kích hoạt:** Người duyệt mở hàng chờ hoặc material detail và chọn duyệt.
- **Hậu điều kiện thành công:** Material chuyển `APPROVED` hoặc `REJECTED`; có ApprovalDecision mới.
- **Hậu điều kiện thất bại:** Material giữ trạng thái cũ.
- **Dữ liệu đầu vào:** `material_id`, `compliance_officer_id`, `approved`, `reject_reason`, `notes`
- **Dữ liệu đầu ra:** Material updated detail
- **Thực thể bị tác động:** `Material`, `ApprovalDecision`
- **Quy tắc nghiệp vụ áp dụng:** Quyết định duyệt ghi lại người duyệt, kết quả và thời điểm.
- **Luồng chính:**
  1. Người duyệt mở màn hàng chờ flagged.
  2. Xem chi tiết tài liệu và reasoning AI.
  3. Chọn approve hoặc reject.
  4. Nếu reject thì nhập lý do.
  5. Gửi request review.
  6. Backend tạo `ApprovalDecision`.
  7. Backend cập nhật status material sang `APPROVED` hoặc `REJECTED`.
  8. Trả kết quả.
- **Luồng thay thế:**
  - A1. Người duyệt chỉ xem detail mà chưa quyết định.
- **Luồng ngoại lệ:**
  - E1. Material không tồn tại.
  - E2. Request body thiếu dữ liệu.
  - E3. Lỗi DB khi tạo decision hoặc update material.
- **Ràng buộc phân quyền:** Theo business nên là Compliance/Admin; theo code hiện tại chỉ cần auth.
- **API / module / màn hình liên quan:** `POST /api/v1/materials/:id/review`, `ComplianceQueuePage`, `MaterialDetailDialog`
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Đây là use case có **lỗ hổng phân quyền lớn**.
  - `compliance_officer_id` hiện do client gửi lên, không lấy từ auth context.
  - **Độ tin cậy:** Confirmed from code

## UC-CNS-01. Tạo / theo dõi consultation

- **Mã use case:** UC-CNS-01
- **Tên use case:** Tạo / theo dõi consultation
- **Mục tiêu nghiệp vụ:** Ghi nhận đầu mối tư vấn tuyển sinh hoặc nhu cầu học của khách hàng tiềm năng.
- **Phạm vi:** Tư vấn / lead intake
- **Mức độ ưu tiên:** Thấp / tương lai
- **Tác nhân chính:** Khách hàng tiềm năng / nhân viên tư vấn
- **Tác nhân phụ:** Quản trị viên
- **Mô tả ngắn:** Có entity `Consultation` với `full_name`, `phone`, `grade_level`, `notes`, `status`, nhưng chưa có API hoặc màn hình.
- **Tiền điều kiện:** Không yêu cầu nếu là lead mới.
- **Điều kiện kích hoạt:** Có nhu cầu tiếp nhận thông tin tư vấn.
- **Hậu điều kiện thành công:** Consultation record được tạo.
- **Hậu điều kiện thất bại:** Không tạo consultation.
- **Dữ liệu đầu vào:** `full_name`, `phone`, `grade_level`, `notes`, `status`
- **Dữ liệu đầu ra:** Consultation detail / danh sách lead
- **Thực thể bị tác động:** `Consultation`
- **Quy tắc nghiệp vụ áp dụng:** Chưa có use case rõ; status mặc định theo entity/migration là `PENDING`.
- **Luồng chính (mong muốn):**
  1. Người dùng hoặc nhân viên tư vấn nhập thông tin lead.
  2. Hệ thống lưu bản ghi consultation.
  3. Nhân viên theo dõi, cập nhật trạng thái xử lý lead.
- **Luồng ngoại lệ:** Chưa có API.
- **Ràng buộc phân quyền:** Cần xác nhận BA.
- **API / module / màn hình liên quan:** Chưa có backend/frontend use case.
- **Ghi chú / điểm mơ hồ cần xác nhận:**
  - Đây là use case **Giả định / cần xác nhận BA** gần như toàn phần.
  - **Độ tin cậy:** Confirmed from entity; workflow Assumption / Needs BA validation

---

# 5. Nhóm use case theo cụm sơ đồ

## 5.1 Sơ đồ xác thực

- **Actors included:** Guest, User, SMTP
- **Use cases included:**
  - UC-AUTH-01 Đăng ký tài khoản
  - UC-AUTH-02 Xác minh email OTP
  - UC-AUTH-03 Đăng nhập
  - UC-AUTH-04 Refresh token
  - UC-AUTH-05 Quên mật khẩu
  - UC-AUTH-06 Đặt lại mật khẩu
  - UC-AUTH-07 Đổi mật khẩu
  - UC-AUTH-08 Xem hồ sơ cá nhân
- **Why grouped that way:** Các use case này cùng xoay quanh vòng đời truy cập và danh tính người dùng, dễ vẽ chung với include/extend cho OTP và reset password.

## 5.2 Sơ đồ quản lý học vụ

- **Actors included:** Admin, Teacher, Student
- **Use cases included:**
  - UC-STU-01, UC-STU-02, UC-STU-03, UC-STU-04
  - UC-TCH-01, UC-TCH-02, UC-TCH-03
  - UC-CRS-01, UC-CRS-02
  - UC-PRG-01, UC-PRG-02, UC-PRG-03
- **Why grouped that way:** Đây là nhóm master data và học vụ nền cho toàn bộ vận hành đào tạo.

## 5.3 Sơ đồ lớp học và ghi danh

- **Actors included:** Admin, Student, Teacher
- **Use cases included:**
  - UC-CLS-01, UC-CLS-02, UC-CLS-03
  - UC-ENR-01, UC-ENR-02
  - UC-SCHCFG-01
- **Why grouped that way:** Cụm này mô tả toàn bộ vòng đời mở lớp, gán giáo viên, xây roster và cấu hình lịch tuần, là phần vận hành trực tiếp trước scheduling.

## 5.4 Sơ đồ xếp lịch

- **Actors included:** Admin, Scheduling Engine
- **Use cases included:**
  - UC-SOL-01 Tạo preview xếp lịch
  - UC-SOL-02 Xem preview
  - UC-SOL-03 Benchmark solver
  - UC-SOL-04 Xác nhận preview để tạo lesson
  - UC-LSN-01 Xem danh sách lesson
- **Why grouped that way:** Cùng một chuỗi logic từ chuẩn bị input -> chạy solver -> rà soát -> commit -> xem đầu ra lesson.

## 5.5 Sơ đồ học thuật

- **Actors included:** Teacher, Admin, Student
- **Use cases included:**
  - UC-LSN-01
  - UC-ATD-01
  - UC-SUM-01
  - UC-ACR-01
  - UC-LVE-01
  - UC-LVE-02
- **Why grouped that way:** Đây là chuỗi nghiệp vụ sau lesson: chuyên cần, tổng kết, đánh giá, xử lý đơn nghỉ.

## 5.6 Sơ đồ tài liệu và kiểm duyệt

- **Actors included:** Teacher, Compliance, Admin, OCR/AI Stub
- **Use cases included:**
  - UC-MAT-01 Upload tài liệu
  - UC-MAT-02 Duyệt tài liệu bị gắn cờ
- **Why grouped that way:** Flow có cả người và hệ thống tự động; dễ thể hiện include giữa upload và AI audit stub.

## 5.7 Sơ đồ hỗ trợ / tuyển sinh

- **Actors included:** Guest/Lead, Staff tư vấn, Admin
- **Use cases included:**
  - UC-CNS-01 Tạo / theo dõi consultation
- **Why grouped that way:** Đây là cụm tiền tuyển sinh, hiện còn sơ khai nên nên tách riêng để không làm rối sơ đồ lõi.

---

# 6. Use case chưa đủ thông tin

| Mã use case | Lý do chưa đủ chắc chắn | Phần nào cần hỏi BA/PO | Mức độ ảnh hưởng |
|---|---|---|---|
| UC-AUTH-08 | `/auth/me` route và use case lỗi kiểu dữ liệu / auth middleware | Profile phải hoạt động ra sao? có dùng API hay chỉ local state? | Cao |
| UC-PRG-03 | Program có dấu vết publish/archive/status nhưng chưa có contract thống nhất | Chương trình có vòng đời chính thức nào? | Cao |
| UC-SCHCFG-01 | Có data model `class_schedules` nhưng không có API/màn hình quản trị | Ai cấu hình lịch tuần? cấu hình ở đâu? | Cao |
| UC-SOL-03 | Benchmark solver mới là contract, chưa có metric benchmark thật | Tiêu chí benchmark và cách hiển thị báo cáo? | Cao |
| UC-LSN-01 | Lesson chỉ có entity và commit từ scheduling, chưa có màn/API tổng quát | Có cần lesson management riêng không? | Trung bình-cao |
| UC-ATD-01 | Không có API, không có enum status attendance | Các trạng thái điểm danh là gì? ai chấm? có chấm bù không? | Cao |
| UC-SUM-01 | Chỉ có entity `LessonSummary`, thiếu workflow và giao diện | Có bắt buộc tổng kết mỗi lesson không? | Trung bình-cao |
| UC-ACR-01 | Chỉ có entity `AcademicRecord`, chưa có rule tính điểm | Điểm thành phần nào là bắt buộc? tính tổng ra sao? | Cao |
| UC-LVE-01 | Chỉ có entity `LeaveRequest`, không có flow | Học viên hay phụ huynh gửi đơn? có cần tài liệu đính kèm? | Trung bình-cao |
| UC-LVE-02 | Chưa có lifecycle duyệt đơn | Ai duyệt? một cấp hay nhiều cấp? ảnh hưởng attendance thế nào? | Trung bình-cao |
| UC-MAT-02 | Quyền duyệt hiện không an toàn | Compliance role có chính thức tồn tại không? | Cao |
| UC-CNS-01 | Chỉ có entity `Consultation`, chưa có use case | Lead intake có nằm trong phạm vi không? | Trung bình |
| UC-STU-01 đến UC-STU-04 | Quyền hiện tại quá rộng | Có cần tách role vận hành / giáo vụ? | Cao |
| UC-ENR-01 | Logic ghi danh chưa đếm enrollment hiện hữu, chưa duplicate check | Quy tắc sĩ số và lifecycle enrollment thật là gì? | Cao |

---

**Kết luận ngắn:** Bộ use case trên đủ để BA bắt đầu viết đặc tả chính thức, vẽ sơ đồ use case và BPMN cho lõi hệ thống. Tuy nhiên, các cụm `class_schedule`, `attendance`, `lesson summary`, `academic record`, `leave request`, `consultation`, cùng các vấn đề `RBAC` và `program lifecycle` cần được xác nhận với PO/BA trước khi chốt tài liệu nghiệp vụ cuối cùng.
