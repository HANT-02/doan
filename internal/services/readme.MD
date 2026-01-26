### Service

1. **Định nghĩa**: Service trong kiến trúc phần mềm thường đề cập đến một lớp hoặc mô-đun chứa logic không đặc thù cho một usecase cụ thể nào mà phục vụ các nghiệp vụ chung trong ứng dụng. Services thường được thiết kế để làm việc ở mức độ thấp hơn, cung cấp các chức năng hỗ trợ cho nhiều usecase khác nhau.

2. **Trách nhiệm**:
    - Thực hiện các chức năng chung như xử lý mã hóa, giao tiếp mạng, thao tác dữ liệu, v.v.
    - Cung cấp một giao diện đơn giản cho các usecase để tương tác với nó mà không phải quan tâm đến chi tiết thực hiện bên trong.
    - Có thể được tái sử dụng bởi nhiều usecase khác nhau trong toàn bộ ứng dụng.

3. **Ví dụ**: Một service có thể là "Authentication Service" mà xử lý việc xác thực người dùng, hoặc "Email Service" mà gửi email thông báo tới khách hàng.

### Kết Luận

- **Usecase** tập trung vào việc thực hiện các yêu cầu cụ thể của kinh doanh và là cầu nối giữa người dùng cuối và hệ thống.
- **Service** cung cấp các chức năng hỗ trợ dùng chung, giúp các usecase thực hiện nhiệm vụ của chúng mà không phải quan tâm đến chi tiết kỹ thuật.

Trong khi usecase là về "cái gì cần làm" để đáp ứng yêu cầu kinh doanh, thì service là về "làm thế nào" để hỗ trợ các yêu cầu đó ở một cấp độ kỹ thuật. Sự phân biệt này giúp trong việc tổ chức và phân chia trách nhiệm trong ứng dụng một cách rõ ràng.