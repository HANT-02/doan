# KẾ HOẠCH CHỐT SCOPE ĐỒ ÁN THEO BÀI TOÁN XẾP LỊCH

**Ngày chốt:** 2026-05-13

## 1. Quyết định phạm vi

Đồ án ở giai đoạn chấm điểm sẽ tập trung vào **bài toán xếp lịch học**.

Cụ thể:

- Phần **học máy / at-risk prediction** không còn nằm trong scope chính để bảo vệ.
- Nếu cần nhắc tới học máy trong báo cáo, chỉ đặt ở mục **hướng phát triển tương lai**.
- Trọng tâm kỹ thuật cần chứng minh bằng số liệu là:
  - mô hình hóa bài toán xếp lịch;
  - bộ luật cứng và luật mềm;
  - benchmark nhiều thuật toán trên cùng một bộ dữ liệu chuẩn;
  - lý do chọn `CP-SAT` làm thuật toán chính.

## 2. Mục tiêu cần bảo vệ được

Khi trình bày, hệ thống cần trả lời rõ 5 câu hỏi:

1. Bài toán xếp lịch trong trung tâm dạy thêm được phát biểu như thế nào?
2. Dữ liệu đầu vào chuẩn gồm những gì và vì sao đủ sát thực tế?
3. Bộ luật cứng nào bắt buộc phải đúng trước khi lịch được dùng?
4. Khi mọi thuật toán đều đúng về luật cứng thì so tiếp bằng luật mềm và hiệu năng ra sao?
5. Vì sao `CP-SAT` được chọn làm thuật toán mặc định thay vì chọn theo cảm tính?

## 3. Phạm vi nghiệp vụ nên giữ lại

Chỉ giữ các phần phục vụ trực tiếp cho scheduling:

- `Teacher`
- `Student`
- `Course`
- `Class`
- `Room`
- `Shift`
- `ClassSchedule`
- `Scheduling Preview`
- `Scheduling Commit`
- `Lesson`

Các phần khác chỉ giữ ở mức đủ làm dữ liệu nền nếu đang có trong hệ thống, không mở rộng thành nhánh tài liệu riêng.

## 4. Phát biểu lại bài toán scheduling

### 4.1 Đầu vào chuẩn

Một bài toán scheduling chuẩn cần có:

- danh sách lớp cần xếp;
- số buổi cần học của từng lớp;
- giáo viên phụ trách hoặc tập giáo viên có thể dạy;
- sĩ số lớp;
- danh sách phòng học và sức chứa;
- danh sách ca học `Shift`;
- lịch tuần mẫu `ClassSchedule`;
- khoảng thời gian chạy xếp lịch;
- các lesson đã tồn tại để tránh đè lịch;
- các tham số ưu tiên dùng cho luật mềm.

### 4.2 Đầu ra chuẩn

Đầu ra của solver phải là:

- danh sách assignment hợp lệ theo từng buổi;
- danh sách conflict hoặc session chưa xếp được;
- thống kê vi phạm luật cứng;
- điểm luật mềm;
- thời gian xử lý;
- dữ liệu đủ rõ để người dùng preview và commit.

## 5. Mô hình dữ liệu nên chuyên nghiệp hơn

### 5.1 Chuẩn hóa thời gian bằng `Shift`

Không nên mô hình hóa thời gian bằng giờ rời rạc nhập tay ở nhiều nơi. Nên giữ `Shift` là đơn vị chuẩn của hệ thống:

- `shift_id`
- `code`
- `name`
- `start_time`
- `end_time`
- `duration_minutes`
- `session_type`
- `is_active`

Lợi ích:

- backend và frontend dùng chung một khái niệm ca học;
- solver làm việc trên miền thời gian ổn định;
- giảm lỗi nhập giờ thủ công;
- dễ mở rộng ca sáng, chiều, tối, cuối tuần, ca online, ca tăng cường.

### 5.2 Lịch tuần nên lưu gì?

`ClassSchedule` nên đóng vai trò là **mẫu lịch tuần**, không phải lesson thực tế.

Tối thiểu nên lưu:

- `class_id`
- `day_of_week`
- `shift_id`
- tùy chọn `preferred_room_id` nếu có
- cờ hiệu lực nếu sau này cần bật/tắt một mẫu lịch

Như vậy câu trả lời cho câu hỏi "lịch tuần chỉ cài đặt ngày/thứ học thì sao?" là:

- **Không nên chỉ lưu ngày/thứ đơn thuần.**
- Nên lưu ít nhất **`day_of_week + shift_id`**.
- Nếu cần linh động hơn, có thể mở rộng thêm ưu tiên phòng hoặc campus, nhưng không nên biến `ClassSchedule` thành lesson thật.

### 5.3 Lesson thực tế sinh ở đâu?

`Lesson` chỉ nên được sinh sau khi:

1. chạy solver tạo preview;
2. preview không vi phạm luật cứng;
3. người dùng xác nhận commit.

Điều này giúp tránh sửa tay dữ liệu thật quá sớm.

## 6. Cách xử lý xung đột để ít phải sửa tay nhất

Mục tiêu không phải triệt tiêu 100% thao tác tay, mà là giảm tối đa số lần phải sửa bằng tay.

Nguyên tắc nên chốt:

- solver chịu trách nhiệm tạo phương án tự động tốt nhất;
- preview cho phép người dùng kiểm tra trước khi commit;
- nếu còn conflict hoặc unscheduled session thì không commit;
- chỉnh tay chỉ là bước tinh chỉnh cuối, không phải bước vận hành chính.

### 6.1 Có nên kéo thả trên UI không?

Có, nhưng kéo thả chỉ nên là **lớp điều chỉnh trên preview**, không thay thế solver.

Luồng đúng nên là:

1. người dùng chạy auto scheduling;
2. hệ thống hiển thị preview theo tuần;
3. nếu một buổi cần đổi slot/phòng thì người dùng kéo thả hoặc chọn lại candidate;
4. hệ thống re-validate ngay luật cứng;
5. chỉ cho commit khi toàn bộ preview hợp lệ.

Nói ngắn gọn:

- **UI kéo thả là hợp lý**, nhưng phải đi kèm kiểm tra ràng buộc tức thời.
- Không nên cho sửa tay trực tiếp trên `Lesson` sau commit như một cách vá lỗi chính.

## 7. Bộ luật đánh giá thuật toán

## 7.1 Luật cứng

Đây là lớp điều kiện bắt buộc. Thuật toán nào vi phạm luật cứng thì bị loại trước khi so tiếp.

Danh sách đề xuất:

- không trùng giáo viên ở cùng thời điểm;
- không trùng phòng ở cùng thời điểm;
- không trùng lớp ở cùng thời điểm;
- chỉ dùng `Shift` đang hoạt động;
- chỉ xếp vào slot sinh ra từ `ClassSchedule`;
- phòng phải đủ sức chứa;
- số buổi xếp được không vượt nhu cầu của khóa học;
- buổi học phải nằm trong khoảng ngày cho phép;
- không đè lên lesson đã tồn tại hoặc slot đã khóa;
- nếu có ràng buộc giáo viên nghỉ cố định thì phải tôn trọng.

## 7.2 Luật mềm

Chỉ so luật mềm sau khi thuật toán đã sạch luật cứng.

Danh sách đề xuất:

- hạn chế khoảng trống lớn giữa các ca dạy của cùng giáo viên;
- hạn chế dồn quá nhiều buổi liên tiếp cho một giáo viên;
- ưu tiên phòng quen thuộc hoặc phòng tối ưu sức chứa;
- ưu tiên phân bố lịch đều trong tuần cho lớp;
- hạn chế đổi ca liên tục giữa các buổi của cùng lớp;
- hạn chế dùng các ca ít mong muốn nếu có ca ưu tiên hơn;
- giảm số buổi chưa xếp phải xử lý tay.

## 8. Thiết kế bộ dữ liệu benchmark chuẩn

## 8.1 Mục tiêu

Cần một **bộ dữ liệu lớn, cố định, tái sử dụng được** cho tất cả thuật toán.

Bộ dữ liệu này là nền để:

- so sánh công bằng;
- tái lập kết quả;
- dùng chung cho benchmark CLI, API và báo cáo;
- tránh việc mỗi thuật toán dùng một input khác nhau.

## 8.2 Thành phần dữ liệu

Bộ dữ liệu nên bao gồm:

- nhiều lớp với số buổi học khác nhau;
- nhiều giáo viên, trong đó có giáo viên dạy nhiều lớp;
- nhiều học sinh để phản ánh quy mô vận hành;
- nhiều phòng với sức chứa khác nhau;
- nhiều ca học `Shift`;
- lịch tuần đa dạng giữa các lớp;
- một số slot bị khóa hoặc lesson có sẵn để mô phỏng dữ liệu thực;
- các case dễ, vừa, khó.

## 8.3 Cấp độ dataset

Nên chốt ít nhất 3 mức:

- `small`: để debug, demo, đọc kết quả nhanh;
- `medium`: để so hành vi solver ở quy mô trung bình;
- `large`: để làm số liệu chính khi bảo vệ vì đủ nhiều lớp, giáo viên, phòng và ràng buộc chéo.

Nếu đủ thời gian, có thể thêm:

- `stress`: dữ liệu dày đặc, ít phòng, nhiều lớp cạnh tranh cùng ca.

## 8.4 Yêu cầu với dataset

Dataset phải:

- deterministic;
- có seed cố định nếu sinh tự động;
- export được ra file để tái sử dụng;
- mô tả rõ số lượng lớp, giáo viên, học sinh, phòng, ca, số session;
- được tất cả solver dùng chung không chỉnh tay riêng.

## 9. Phương pháp benchmark

## 9.1 Thuật toán cần so

Giữ 3 thuật toán:

- `Graph Coloring + Heuristic`
- `CP-SAT`
- `Tabu Search`

## 9.2 Chỉ số cần đo

Mỗi thuật toán cần đo ít nhất:

- `feasibility rate`
- `hard violation count`
- `unscheduled session count`
- `soft score`
- `runtime`
- `runtime stability` qua nhiều lần chạy

## 9.3 Quy trình so sánh

Quy trình chuẩn:

1. Chạy tất cả solver trên cùng dataset.
2. Loại solver vi phạm luật cứng hoặc không đạt feasibility.
3. Với nhóm còn lại, so tiếp:
   - số session xếp được;
   - điểm luật mềm;
   - thời gian xử lý;
   - độ ổn định.
4. Chọn solver mặc định dựa trên số liệu tổng hợp, không chọn cảm tính.

## 10. Lý do chọn `CP-SAT` cần được chứng minh thế nào

Không nên viết "chọn `CP-SAT` vì nổi tiếng" hoặc "vì mạnh hơn".

Lập luận nên đi theo thứ tự:

1. `Graph Coloring + Heuristic` nhanh nhưng soft score không tốt bằng trong các case phân biệt.
2. `Tabu Search` có ý nghĩa học thuật nhưng runtime cao hơn hoặc không tạo lợi thế đủ lớn về chất lượng nghiệm.
3. `CP-SAT` giữ được:
   - nghiệm sạch luật cứng;
   - chất lượng nghiệm ổn định;
   - soft score tốt hơn hoặc cân bằng tốt nhất;
   - runtime vẫn chấp nhận được trong quy mô đồ án.

Kết luận chốt:

- `CP-SAT` là solver chính cho production-like scheduling.
- Hai solver còn lại là baseline và đối chứng thực nghiệm.

## 11. Kế hoạch triển khai đề xuất

## Giai đoạn 1 - Chốt lại tài liệu và scope

- bỏ tài liệu BA/predictive cũ khỏi nhánh tài liệu chính;
- giữ lại nhóm tài liệu scheduling;
- viết lại mô tả đề tài theo hướng "hệ thống quản lý trung tâm dạy thêm, trọng tâm là xếp lịch học".

## Giai đoạn 2 - Chuẩn hóa domain scheduling

- rà lại `Shift`;
- rà lại `ClassSchedule`;
- xác nhận đầu vào chuẩn cho solver;
- xác định rõ hard constraints và soft constraints bằng tài liệu.

## Giai đoạn 3 - Chuẩn hóa benchmark dataset

- tạo dataset `small`, `medium`, `large`;
- thêm số lượng lớp, giáo viên, học sinh, phòng lớn hơn để phản ánh vận hành thật hơn;
- cố định seed hoặc file input dùng chung;
- bổ sung case có conflict chéo và slot khan hiếm.

## Giai đoạn 4 - Đo đạc thực nghiệm

- chạy benchmark nhiều lần;
- ghi lại runtime, feasibility, vi phạm luật cứng, soft score;
- tổng hợp bảng và biểu đồ;
- chốt solver thắng theo tiêu chí đã công bố trước.

## Giai đoạn 5 - Hoàn thiện UX preview và chỉnh tay tối thiểu

- cho phép chỉnh candidate ở màn preview;
- nếu cần thì hỗ trợ kéo thả theo tuần;
- mỗi lần chỉnh phải re-check hard constraints ngay;
- khóa commit khi còn conflict.

## Giai đoạn 6 - Chốt nội dung báo cáo và bảo vệ

- mô tả bài toán;
- mô tả mô hình dữ liệu;
- mô tả hard/soft constraints;
- mô tả dataset benchmark;
- đưa bảng số liệu thực tế;
- kết luận vì sao chọn `CP-SAT`.

## 12. Deliverable cần có để chấm điểm

Các đầu ra nên chốt thành checklist:

- tài liệu scope chỉ còn trục scheduling;
- bộ luật hard/soft rõ ràng;
- bộ dataset benchmark chuẩn dùng chung;
- benchmark report có số liệu lặp lại được;
- solver mặc định `CP-SAT`;
- màn preview đủ để xem, chỉnh và commit;
- phần học máy chỉ còn ở mục hướng phát triển.

## 13. Kết luận định hướng

Từ thời điểm này, toàn bộ câu chuyện của đồ án nên xoay quanh một luận điểm chính:

> Hệ thống không chỉ quản lý dữ liệu trung tâm dạy thêm, mà giải quyết bài toán xếp lịch học bằng nhiều thuật toán, đo đạc được bằng dữ liệu thực nghiệm, và lựa chọn `CP-SAT` làm solver chính dựa trên kết quả benchmark.

Đó là trục thuyết phục nhất để viết báo cáo, demo hệ thống và trả lời phản biện.
