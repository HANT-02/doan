### Usecase

1. **Định nghĩa**: Usecase là một mô tả về loạt các bước mà hệ thống và các actor tương tác để đạt được một mục tiêu cụ thể. Trong kiến trúc phần mềm, usecase thường đại diện cho một tính năng hoặc chức năng kinh doanh cụ thể mà người dùng cuối mong muốn.

2. **Trách nhiệm**:
    - Diễn giải yêu cầu kinh doanh thành các bước có thể thực thi.
    - Đảm bảo rằng logic kinh doanh cần thiết được thực hiện để đạt được mục tiêu đó.
    - Có thể tương tác trực tiếp với các thành phần khác như entities, repositories, và services để thực hiện nghiệp vụ.

3. **Ví dụ**: Một usecase có thể là "Xác nhận Đơn Hàng", trong đó bao gồm việc kiểm tra tính khả dụng của sản phẩm, xác nhận thông tin thanh toán của khách hàng, và gửi thông tin đơn hàng cho bộ phận giao hàng.

4. **Lưu ý**:
   
   - Usecase không phụ thuộc vào framework (gin, echo, mux,...)
   - Usecase có thể sử dụng services, usecases khác
   - Sử dụng service cho các logic nghiệp vụ phức tạp
   - Sử dụng repository cho các tác vụ CRUD trực tiếp