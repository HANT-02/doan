# Kế Hoạch Đo Đạc Đối Sánh Xếp Lịch Theo Chuẩn Học Thuật Và Thực Tiễn

## 1. Mục tiêu của kế hoạch

Tài liệu này là bản cập nhật kế hoạch đo đạc đối sánh trước khi triển khai chính thức, với mục tiêu:

- chuẩn hóa bộ dữ liệu kiểm thử cho bài toán xếp lịch;
- bổ sung cách đo có giá trị học thuật, không chỉ dừng ở “chạy được hay không”;
- đưa kế hoạch đo đạc tiến gần hơn tới chuẩn vận hành thực tế;
- bảo đảm ba thuật toán hiện có được đo trong cùng một khuôn khổ:
  - thuật toán tô màu đồ thị
  - thuật toán CP-SAT
  - thuật toán tìm kiếm Tabu

Tài liệu này bám theo luồng xử lý thực tế trong mã nguồn:

- nạp dữ liệu và dựng bản xem trước: [preview.go](/Users/hant/golang/doan/internal/usecases/scheduling/preview.go)
- lọc điều kiện mở lớp và ngưỡng chỉnh tay: [policy.go](/Users/hant/golang/doan/internal/usecases/scheduling/policy.go)
- hợp đồng đầu vào và đầu ra của bộ giải: [contracts.go](/Users/hant/golang/doan/internal/services/scheduling/contracts.go)
- tô màu đồ thị: [graph_coloring_solver.go](/Users/hant/golang/doan/internal/services/scheduling/graph_coloring_solver.go)
- CP-SAT hiện tại: [cp_sat_solver.go](/Users/hant/golang/doan/internal/services/scheduling/cp_sat_solver.go)
- thuật toán tìm kiếm Tabu: [tabu_search_solver.go](/Users/hant/golang/doan/internal/services/scheduling/tabu_search_solver.go)
- mã đo đạc tổng hợp hiện tại: [benchmark_study.go](/Users/hant/golang/doan/internal/services/scheduling/benchmark_study.go)

## 1.1 Quy ước ngôn ngữ trình bày

Toàn bộ phần dùng để đưa vào báo cáo, bảng số liệu, biểu đồ, nhật ký minh chứng và tên chỉ số hiển thị ra ngoài phải dùng tiếng Việt hoàn toàn.

Quy ước áp dụng:

- không dùng tên chỉ số tiếng Anh trong bảng báo cáo;
- không dùng tên cột tiếng Anh trong tệp số liệu nộp kèm cho hội đồng;
- không dùng tên thư mục hay tên tệp minh chứng bằng tiếng Anh nếu tệp đó được trích dẫn trong báo cáo;
- chỉ giữ nguyên tiếng Anh ở:
  - đường dẫn file mã nguồn;
  - tên bảng và tên trường kỹ thuật nếu cần đối chiếu trực tiếp với cơ sở dữ liệu;
  - tên thuật toán chuẩn như `CP-SAT`.

Nói ngắn gọn:

- bên trong mã nguồn có thể còn định danh kỹ thuật;
- nhưng mọi thứ xuất ra để đo đạc, tổng hợp, báo cáo và bảo vệ đều phải Việt hóa hoàn toàn.

## 2. Khung tư duy đo đạc đối sánh mới

Kế hoạch đo đạc của đồ án này không chỉ trả lời câu hỏi “thuật toán nào nhanh hơn”, mà phải trả lời đủ bốn câu hỏi:

1. Bài toán lớn đến mức nào.
2. Không gian tìm kiếm của nó lớn ra sao.
3. Hệ thống đã giảm độ khó của bài toán như thế nào trước khi vào bộ giải.
4. Trong cùng một điều kiện dừng, thuật toán nào cho nghiệm tốt hơn, ổn định hơn và gần với vận hành thực tế hơn.

Vì vậy, kế hoạch đo đạc phải có hai lớp:

- lớp phân tích phễu dữ liệu:
  - đo từng bước tiền xử lý;
  - chứng minh hệ thống giảm tải cho bộ giải từ sớm.
- lớp đo đạc đối sánh thuật toán:
  - đo chất lượng nghiệm;
  - đo hiệu năng;
  - đo độ ổn định;
  - đo khả năng mở rộng.

## 3. Điều kiện tiên quyết trước khi đo đạc chính thức

Trước khi triển khai đo đạc hoàn chỉnh, cần chốt bốn điều kiện kỹ thuật sau.

### 3.1 Chuẩn hóa điều kiện dừng

Ba thuật toán hiện tại chưa dùng cùng một cơ chế dừng:

- thuật toán tô màu đồ thị chạy tham lam một lượt;
- thuật toán CP-SAT hiện bị chặn bởi số nút tối đa;
- thuật toán tìm kiếm Tabu hiện bị chặn bởi số vòng lặp tối đa.

Để đo đạc công bằng, cần bổ sung một cấu hình thống nhất, ví dụ:

```text
Ngân sách thời gian
Số vòng lặp tối đa
Số nút tối đa
Hạt giống ngẫu nhiên
Bật làm nóng
Bật ghi đo đạc
```

Khuyến nghị cho đồ án:

- chuẩn chính thức dùng ngân sách thời gian làm điều kiện dừng công bằng;
- số nút tối đa và số vòng lặp tối đa chỉ là giới hạn bảo vệ nội bộ;
- hạt giống ngẫu nhiên phải được ghi lại ở mỗi lần chạy.

### 3.2 Bổ sung lớp ghi đo đạc

Hiện tại hệ thống đã có phần tóm tắt kết quả, nhưng chưa đủ cho đo đạc học thuật. Cần bổ sung một cấu trúc ghi đo đạc, ví dụ:

```text
Mã lượt chạy
Thời điểm bắt đầu
Thời điểm kết thúc
Thời gian chạy mili giây
Số nút đã duyệt
Số nhánh bị cắt
Số lần cập nhật nghiệm tốt nhất
Số phương án đã đánh giá
Số bước chuyển được chấp nhận
Số bước chuyển bị danh sách cấm loại bỏ
Số gán được sửa ở bước phục hồi
Số lần tính điểm
Tốc độ tính điểm mỗi giây
Bộ nhớ đỉnh theo byte
```

Ngoài ra, để làm minh chứng rằng số liệu đã được đo thật chứ không ghi tay lại sau, mỗi lần chạy cần lưu thêm:

```text
Thời điểm ghi nhận
Múi giờ
Tên máy
Mã tiến trình
Mã phiên bản mã nguồn
Tên kịch bản
Phiên bản kịch bản
Khóa thuật toán
Chỉ số lần lặp
Đường dẫn nhật ký gốc
```

### 3.3 Chuẩn hóa cách đo bộ nhớ

Đề xuất dùng hàm `runtime.ReadMemStats` của Go ở ba thời điểm:

- trước bước xử lý;
- sau bước xử lý;
- sau khi giải phóng tạm và cưỡng bức cơ chế gom rác cho phiên đo đạc.

Chỉ số nên lưu:

- bộ nhớ cấp phát hiện tại theo byte
- tổng bộ nhớ đã cấp phát theo byte
- bộ nhớ heap đang dùng theo byte
- số đối tượng trên heap
- số lần gom rác

### 3.4 Chuẩn hóa cách đo tốc độ tính toán điểm số

Đây là chỉ số rất quan trọng để nâng chất lượng báo cáo. Cần đo:

- tổng số lần đánh giá nghiệm hoặc đánh giá phương án;
- tốc độ tính toán:
  - tốc độ tính điểm mỗi giây = số lần tính điểm chia cho số giây chạy thực tế

Với mã nguồn hiện tại, chỉ số này chưa có sẵn, nên bắt buộc phải gắn bộ đếm trong từng bộ giải.

## 4. Quy mô bộ dữ liệu thực nghiệm

### 4.1 Kịch bản đo đạc chính thức cho báo cáo

Đây là bộ đo đạc chính thức nên xuất hiện trong báo cáo.

| Kịch bản | Số lớp | Số giáo viên | Số phòng | Tổng số buổi cần xếp | Không gian tìm kiếm mục tiêu | Đối tượng mô phỏng |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| Nhỏ | 10 - 20 | 5 - 10 | 5 | 50 - 100 | khoảng `10^60` | trung tâm nhỏ, kiểm thử nhanh tính đúng đắn |
| Trung bình | 50 - 100 | 20 - 40 | 15 - 20 | 300 - 500 | khoảng `10^1000` | chi nhánh trung tâm ngoại ngữ tiêu chuẩn |
| Lớn | 200 - 500 | 100 - 180 | 50 - 80 | 1.000 - 2.500 | khoảng `10^3000+` | chuỗi trung tâm hoặc khoa, trường quy mô lớn |

### 4.2 Không gian tìm kiếm phải được ghi như thế nào

Trong đồ án, không nên chỉ ghi một con số “ước đoán bằng mắt”. Cần dùng hai dạng số liệu:

#### a. Không gian tìm kiếm danh nghĩa của kịch bản

Đây là số liệu cấp báo cáo, dùng để mô tả độ lớn của bài toán theo quy mô mô phỏng.

#### b. Không gian tìm kiếm quan sát được sau tiền xử lý

Đây là số liệu thật phải đo từ dữ liệu sau khi đã sinh miền giá trị.

Với mỗi biến quyết định `i`, gọi `|D_i|` là kích thước miền giá trị. Khi đó:

```text
SearchSpaceApprox = Π |D_i|
Lôgarit cơ số 10 của không gian tìm kiếm ước lượng = tổng lôgarit cơ số 10 của từng miền giá trị
```

Khuyến nghị:

- trong bảng báo cáo, ghi lôgarit cơ số 10 của không gian tìm kiếm thay vì ghi trực tiếp số quá lớn;
- ngoài ra vẫn lưu miền nhỏ nhất, miền trung bình, miền trung vị và miền lớn nhất.

Ví dụ cách trình bày:

| Kịch bản | Số biến `V` | Miền trung bình | Lôgarit cơ số 10 của không gian tìm kiếm |
| :---- | :---- | :---- | :---- |
| Nhỏ | 78 | 6.4 | 62.7 |
| Trung bình | 412 | 9.8 | 1008.1 |
| Lớn | 1685 | 11.7 | 3064.2 |

### 4.3 Bộ thử nhanh để phát triển trước khi chạy đo đạc lớn

Để tránh bị kẹt tiến độ triển khai, nên có thêm bộ thử nhanh nhỏ hơn chỉ dùng cho phát triển:

| Bộ thử nhanh | Số lớp | Số giáo viên | Số phòng | Tổng số buổi |
| :---- | :---- | :---- | :---- | :---- |
| P1 | 6 | 4 | 3 | 12 |
| P2 | 10 | 5 | 4 | 30 |
| P3 | 16 | 7 | 5 | 64 |

Lưu ý:

- bộ thử nhanh chỉ dùng để rà lỗi và kiểm tra phần ghi đo đạc;
- không dùng bộ này làm đo đạc chính thức trong báo cáo nếu bạn đã có khả năng chạy bộ `nhỏ`, `trung bình`, `lớn` ở trên.

## 5. Dữ liệu nền bắt buộc phải khởi tạo

### 5.1 Danh mục bảng tối thiểu

- `teachers`
- `courses`
- `rooms`
- `shifts`
- `classes`
- `class_schedules`
- `students`
- `enrollments`
- `lessons`
- `campus_travel_times`

### 5.2 Trường tối thiểu cần có theo bảng

| Bảng | Trường tối thiểu |
| :---- | :---- |
| `teachers` | `id`, `code`, `full_name`, `status`, `skills` |
| `courses` | `id`, `code`, `name`, `session_count`, `session_duration_minutes`, `required_skills`, `status` |
| `rooms` | `id`, `code`, `name`, `capacity`, `campus_id` |
| `shifts` | `id`, `code`, `name`, `start_time`, `end_time`, `duration_minutes`, `session_type`, `is_active` |
| `classes` | `id`, `code`, `name`, `start_date`, `end_date`, `max_students`, `status`, `course_id`, `teacher_id`, `room_id` |
| `class_schedules` | `id`, `class_id`, `day_of_week`, `shift_id`, `room_id` |
| `students` | `id`, `code`, `full_name`, `status` |
| `enrollments` | `id`, `class_id`, `student_id`, `status` |
| `lessons` | `id`, `class_id`, `teacher_id`, `room_id`, `date_start`, `date_end`, `status` |
| `campus_travel_times` | `from_campus_id`, `to_campus_id`, `travel_minutes`, `is_active` |

### 5.3 Dữ liệu nền cố định

#### a. Cơ sở

| Mã cơ sở | Tên gợi ý |
| :---- | :---- |
| `CS1` | Cơ sở trung tâm |
| `CS2` | Cơ sở vệ tinh |
| `CS3` | Cơ sở mở rộng, chỉ dùng ở kịch bản lớn |

#### b. Ca học

| Mã ca | Giờ bắt đầu | Giờ kết thúc | Thời lượng | Trạng thái |
| :---- | :---- | :---- | :---- | :---- |
| `S1` | 08:00 | 10:00 | 120 | hoạt động |
| `S2` | 10:15 | 12:15 | 120 | hoạt động |
| `S3` | 13:30 | 15:30 | 120 | hoạt động |
| `S4` | 15:45 | 17:45 | 120 | hoạt động |
| `S5` | 18:00 | 20:00 | 120 | hoạt động |
| `S6` | 20:15 | 22:15 | 120 | tùy chọn, chỉ mở ở kịch bản lớn hoặc kiểm thử |

#### c. Phòng học

Nên chia rõ theo cụm sức chứa:

- cụm nhỏ:
  - 18 chỗ
  - 20 chỗ
- cụm trung bình:
  - 24 chỗ
  - 30 chỗ
- cụm lớn:
  - 36 chỗ
  - 45 chỗ

Mỗi cơ sở nên có đủ cả ba cụm để đo được tỷ lệ lấp đầy phòng học.

#### d. Thời gian di chuyển liên cơ sở

| Từ | Đến | Phút |
| :---- | :---- | :---- |
| `CS1` | `CS2` | 30 |
| `CS2` | `CS1` | 30 |
| `CS1` | `CS3` | 45 |
| `CS3` | `CS1` | 45 |
| `CS2` | `CS3` | 20 |
| `CS3` | `CS2` | 20 |

#### e. Giáo viên

Mỗi giáo viên cần được khởi tạo theo nhóm kỹ năng để tạo được:

- ca hợp lệ;
- ca sai kỹ năng;
- ca hiếm kỹ năng;
- ca gây xung đột tài nguyên.

## 6. Cấu trúc dữ liệu đo đạc chuẩn

### 6.1 Bộ đo đạc sạch

Đây là bộ để so sánh ba thuật toán.

Yêu cầu:

- mọi lớp đều ở trạng thái mở;
- mọi lớp đủ ngưỡng sĩ số mở lớp;
- mọi lớp có giáo viên;
- mọi lớp có khóa học;
- mọi lớp có lịch tuần;
- giáo viên đáp ứng kỹ năng bắt buộc;
- dữ liệu có một lượng xung đột tự nhiên vừa đủ:
  - trùng giáo viên tiềm năng;
  - trùng học sinh tiềm năng;
  - phòng ưu tiên gây áp lực;
  - ràng buộc di chuyển liên cơ sở.

### 6.2 Bộ kiểm thử luật và tiền xử lý

Đây là bộ để đo phễu dữ liệu và chứng minh hệ thống loại dữ liệu không đạt từ sớm.

Giữ lại các gói sau:

- `K1`: sĩ số không đạt ngưỡng mở lớp
- `K2`: thiếu giáo viên phụ trách
- `K3`: giáo viên không đúng kỹ năng
- `K4`: thiếu lịch học tuần
- `K5`: không đủ ca để xếp
- `K6`: áp lực sức chứa phòng
- `K7`: xung đột thời gian di chuyển liên cơ sở
- `K8`: điều chỉnh lại với buổi học đã khóa
- `K9`: điều chỉnh lại với buổi học nháp

### 6.3 Bộ dữ liệu xung đột có kiểm soát để đo thuật toán

Ngoài hai nhóm trên, nên thêm một nhóm xung đột có kiểm soát với đặc điểm:

- dữ liệu hợp lệ để vào bộ giải;
- nhưng đủ khó để thuật toán phải thật sự tối ưu.

Ví dụ:

- 15% lớp có phòng ưu tiên;
- 20% cặp lớp chia sẻ học sinh;
- 25% giáo viên dạy từ 2 lớp trở lên;
- 10% buổi có nguy cơ vi phạm di chuyển liên cơ sở;
- 15% lớp nằm sát giới hạn sĩ số phòng.

## 7. Có cần ghi dữ liệu output ở từng bước hay không

Câu trả lời là có, và đây là phần rất nên làm. Đây là bằng chứng cho thấy hệ thống đã thực hiện “thu hẹp bài toán” trước khi vào engine xếp lịch.

Bạn nên xem toàn bộ quy trình theo dạng phễu:

```text
đầu vào thô
-> lọc điều kiện mở lớp
-> kiểm tra tiền xử lý
-> dựng cửa sổ xếp lịch
-> sinh biến
-> sinh miền giá trị
-> chạy thuật toán
-> đối chiếu buổi học đã có
-> chỉnh tay
-> lưu thật
```

Ở mỗi bước, phải ghi lại cả:

- chỉ số nghiệp vụ;
- chỉ số hệ thống.

## 8. Bộ chỉ số cần đo ở từng bước

### 8.1 Bảng chỉ số chuẩn hóa

| Bước xử lý | Chỉ số nghiệp vụ | Chỉ số hệ thống |
| :---- | :---- | :---- |
| Nạp dữ liệu | số lớp, số GV, số phòng, số ca, số buổi học đã có | thời gian truy vấn cơ sở dữ liệu, số bản ghi nạp, bộ nhớ sau nạp |
| Lọc điều kiện mở lớp | số lớp bị loại do sĩ số thấp | thời gian lọc, số truy vấn ghi danh |
| Kiểm tra tiền xử lý | số lỗi thiếu GV, thiếu khóa học, sai kỹ năng, thiếu lịch tuần | bộ nhớ đệm, thời gian kiểm tra dữ liệu |
| Dựng cửa sổ xếp lịch | số lớp còn lại, số buổi yêu cầu | thời gian sinh cửa sổ, số ngày cửa sổ trung bình |
| Sinh miền giá trị | số biến, số miền rỗng, kích thước miền nhỏ nhất, trung bình, trung vị, lớn nhất | bộ nhớ tăng thêm, số phương án sinh ra |
| Chạy thuật toán | điểm cứng, điểm mềm, tỷ lệ xếp thành công | thời gian chạy, tốc độ tính điểm mỗi giây, bộ nhớ đỉnh |
| Đối chiếu buổi học đã có | số xung đột theo nhóm | thời gian đối chiếu, số buổi học khóa được quét |
| Chỉnh tay và lưu | số ca phải chỉnh tay, tỷ lệ lấp đầy phòng học, số ca lưu thành công | thời gian ghi xuống cơ sở dữ liệu, số câu lệnh ghi |

### 8.2 Chỉ số phễu dữ liệu bắt buộc

Các chỉ số này phải có trong báo cáo:

- `L0`: số lớp ban đầu đọc được
- `L1`: số lớp còn lại sau lọc sĩ số
- `L2`: số lớp còn lại sau tiền xử lý cấu trúc
- `V`: tổng số biến quyết định
- `R0`: tổng số phương án thô
- `R1`: tổng số phương án còn lại sau lọc cứng
- lôgarit cơ số 10 của không gian tìm kiếm ước lượng

### 8.3 Chỉ số nghiệp vụ đầu ra phải có

- số lớp yêu cầu
- số buổi yêu cầu
- số buổi xếp được
- số buổi chưa xếp được
- số xung đột
- điểm mềm
- số lần đổi lịch
- số lần đổi giáo viên
- số lần đổi phòng
- tỷ lệ sử dụng sức chứa phòng trung bình
- số ca cần chỉnh tay
- ngưỡng chỉnh tay
- cờ vượt ngưỡng chỉnh tay
- tỷ lệ lấp đầy phòng học

### 8.4 Chỉ số hệ thống đầu ra phải có

- thời gian chạy mili giây
- thời điểm bắt đầu lượt chạy
- thời điểm kết thúc lượt chạy
- thời điểm ghi nhận
- thời gian truy vấn cơ sở dữ liệu
- thời gian ghi xuống cơ sở dữ liệu
- bộ nhớ đỉnh theo byte
- bộ nhớ heap đang dùng theo byte
- số lần tính điểm
- tốc độ tính điểm mỗi giây
- cờ ổn định qua nhiều lần chạy
- độ lệch chuẩn thời gian chạy
- độ lệch chuẩn điểm mềm
- độ lệch chuẩn số vi phạm cứng

### 8.5 Trường minh chứng bắt buộc cho mỗi lần đo

Mỗi lần đo phải sinh ra một bản ghi thô độc lập, không chỉ ghi số tổng hợp. Bản ghi này là bằng chứng đối chiếu cho báo cáo.

Các trường bắt buộc nên lưu:

- mã lượt chạy
- tên kịch bản
- phiên bản kịch bản
- khóa thuật toán
- chỉ số lần lặp
- chế độ chạy
- hạt giống ngẫu nhiên
- thời điểm bắt đầu lượt chạy
- thời điểm kết thúc lượt chạy
- thời điểm ghi nhận
- múi giờ
- tên máy
- mã phiên bản mã nguồn
- đường dẫn ảnh chụp đầu vào
- đường dẫn số liệu thô
- đường dẫn kết quả xem trước thô
- đường dẫn nhật ký hệ thống thô

Khuyến nghị định dạng thời gian:

- dùng `RFC3339` hoặc `RFC3339Nano`;
- ví dụ: `2026-05-15T14:32:18.245+07:00`.

Khuyến nghị đặt tên thư mục kết quả:

```text
minh_chung_do_dac/
  2026-05-15/
    nho/
      thuat_toan_cp_sat/
        lan-chay-001/
          dau_vao.json
          so_lieu.json
          ket_qua_xem_truoc.json
          nhat_ky_he_thong.txt
```

Như vậy, mỗi dòng số liệu đưa vào bảng báo cáo đều truy ngược được về một thư mục kết quả cụ thể.

## 9. Chỉ số riêng cho từng thuật toán

### 9.1 Tô màu đồ thị

- số biến sau sắp thứ tự
- số nhóm slot được xét
- số phương án bị loại vì xung đột
- tỷ lệ gán thành công ngay lượt đầu
- tổng số phương án đã đánh giá

### 9.2 CP-SAT hiện tại

- số nút đã duyệt
- số nhánh bị cắt
- số lần cập nhật nghiệm tốt nhất
- số nghiệm lá
- cờ chạm ngưỡng số nút tối đa
- tổng số phương án đã thử

### 9.3 Tìm kiếm Tabu

- số vòng lặp đã thực hiện
- điểm phạt ban đầu
- điểm phạt tốt nhất
- số bước chuyển được chấp nhận
- số bước chuyển bị danh sách cấm loại bỏ
- số gán được sửa ở bước phục hồi
- tổng số phương án đã duyệt trong lân cận

## 10. Nguyên tắc đo đạc chuẩn xác

### 10.1 Nguyên tắc 1: Làm nóng trước khi đo

Với Go, không có cơ chế biên dịch tức thời như trong máy ảo Java, nhưng vẫn cần làm nóng vì:

- bộ nhớ và bộ gom rác cần ổn định;
- bộ nhớ đệm dữ liệu, bộ nhớ đệm trang và cơ chế dự đoán nhánh của CPU cần ổn định;
- kết nối cơ sở dữ liệu cần qua giai đoạn làm nóng;
- dữ liệu có thể được nạp vào bộ nhớ đệm ứng dụng.

Hành động:

- chạy nháp 30 giây đến 60 giây trước khi bắt đầu ghi số liệu chính thức;
- hoặc chạy trước 3 đến 5 lần và bỏ các lần đó khỏi báo cáo.

### 10.2 Nguyên tắc 2: Chạy nhiều lần và tính độ lệch chuẩn

Mỗi kịch bản và mỗi thuật toán phải chạy ít nhất 10 lần độc lập.

Với mỗi cặp `kịch bản x thuật toán`, cần tính:

- trung bình thời gian chạy;
- độ lệch chuẩn thời gian chạy;
- trung bình điểm mềm;
- độ lệch chuẩn điểm mềm;
- trung bình số vi phạm cứng;
- độ lệch chuẩn số vi phạm cứng.

Lưu ý:

- dù thuật toán tô màu đồ thị có thể gần như xác định, vẫn phải chạy 10 lần để đo biến động hiệu năng hệ thống;
- với thuật toán có ngẫu nhiên về sau, phải lưu cả hạt giống ngẫu nhiên.

Ngoài số liệu trung bình và độ lệch chuẩn, cần giữ nguyên bảng nhật ký của từng lần chạy đơn lẻ để hội đồng hoặc người đọc có thể kiểm tra:

- thời điểm bắt đầu;
- thời điểm kết thúc;
- thời lượng thực;
- bộ dữ liệu đã dùng;
- thuật toán đã dùng;
- tệp kết quả gốc.

### 10.3 Nguyên tắc 3: So sánh trên cùng điều kiện dừng

Để công bằng, cần khóa một biến.

Khuyến nghị chính thức cho đồ án:

- dùng khóa thời gian.

Ví dụ:

- kịch bản nhỏ: 10 giây
- kịch bản trung bình: 60 giây
- kịch bản lớn: 180 giây

Sau khi hết thời gian, lấy nghiệm tốt nhất mà thuật toán đã tìm được.

Lý do chọn cách này:

- sát thực tế vận hành;
- dễ so sánh giữa thuật toán gần đúng và bộ giải dựa trên tìm kiếm;
- tránh tình huống một thuật toán chạy quá lâu chỉ để tối thêm rất ít.

### 10.4 Nguyên tắc 4: Tốc độ tính toán điểm số

Đây là chỉ số rất nên đưa vào báo cáo.

Phải đo:

- số lần tính điểm
- tốc độ tính điểm mỗi giây

Nếu sau này bạn áp dụng chấm điểm tăng dần hoặc đánh giá delta tốt hơn, đây sẽ là bằng chứng rất mạnh về tối ưu hiệu năng.

## 11. Khuôn đo đạc đề xuất cho báo cáo

### 11.1 Bảng mô tả bộ dữ liệu thực nghiệm

| Kịch bản | Số lớp | Số GV | Số phòng | Tổng số buổi | Lôgarit cơ số 10 của không gian tìm kiếm | Mô tả |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| Nhỏ |  |  |  |  |  |  |
| Trung bình |  |  |  |  |  |  |
| Lớn |  |  |  |  |  |  |

### 11.2 Bảng phễu dữ liệu

| Kịch bản | `L0` | `L1` | `L2` | `V` | `R0` | `R1` | Miền rỗng | Miền nhỏ nhất | Miền trung bình | Miền trung vị | Miền lớn nhất |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| Nhỏ |  |  |  |  |  |  |  |  |  |  |  |
| Trung bình |  |  |  |  |  |  |  |  |  |  |  |
| Lớn |  |  |  |  |  |  |  |  |  |  |  |

### 11.3 Bảng đo đạc thuật toán

| Kịch bản | Thuật toán | Số lần chạy | Tỷ lệ xếp thành công trung bình | Vi phạm cứng trung bình | Điểm mềm trung bình | Thời gian chạy trung bình ms | Độ lệch chuẩn thời gian | Tốc độ tính toán trung bình | Ổn định |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| Nhỏ | Tô màu đồ thị | 10 |  |  |  |  |  |  |  |
| Nhỏ | CP-SAT | 10 |  |  |  |  |  |  |  |
| Nhỏ | Tìm kiếm Tabu | 10 |  |  |  |  |  |  |  |
| Trung bình | Tô màu đồ thị | 10 |  |  |  |  |  |  |  |
| Trung bình | CP-SAT | 10 |  |  |  |  |  |  |  |
| Trung bình | Tìm kiếm Tabu | 10 |  |  |  |  |  |  |  |
| Lớn | Tô màu đồ thị | 10 |  |  |  |  |  |  |  |
| Lớn | CP-SAT | 10 |  |  |  |  |  |  |  |
| Lớn | Tìm kiếm Tabu | 10 |  |  |  |  |  |  |  |

### 11.4 Bảng chỉ số hệ thống

| Kịch bản | Thuật toán | Truy vấn cơ sở dữ liệu ms | Kiểm tra dữ liệu ms | Giải bài toán ms | Đối chiếu buổi học ms | Ghi cơ sở dữ liệu ms | Bộ nhớ đỉnh MB |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| Nhỏ |  |  |  |  |  |  |  |
| Trung bình |  |  |  |  |  |  |  |
| Lớn |  |  |  |  |  |  |  |

### 11.5 Bảng nhật ký từng lần chạy làm minh chứng

| Kịch bản | Thuật toán | Lần chạy | Bắt đầu đo | Kết thúc đo | Thời gian chạy ms | Đường dẫn nhật ký gốc |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| Nhỏ | Tô màu đồ thị | 1 |  |  |  |  |
| Nhỏ | Tô màu đồ thị | 2 |  |  |  |  |
| Nhỏ | CP-SAT | 1 |  |  |  |  |
| Nhỏ | Tìm kiếm Tabu | 1 |  |  |  |  |
| Trung bình | CP-SAT | 1 |  |  |  |  |
| Lớn | CP-SAT | 1 |  |  |  |  |

Lưu ý:

- bảng này không thay thế bảng tổng hợp;
- bảng này là bằng chứng nguồn cho các số trung bình và độ lệch chuẩn ở bảng đo đạc chính.

### 11.6 Bảng chỉnh tay và tỷ lệ lấp đầy

| Kịch bản | Thuật toán | Số ca cần chỉnh tay | Ngưỡng chỉnh tay | Vượt ngưỡng | Tỷ lệ lấp đầy phòng học | Số buổi lưu thành công |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| Nhỏ |  |  |  |  |  |  |
| Trung bình |  |  |  |  |  |  |
| Lớn |  |  |  |  |  |  |

## 12. Kế hoạch chạy đo đạc tuần tự

### Giai đoạn 1. Chuẩn hóa hạ tầng đo đạc

1. Bổ sung cấu hình đo đạc thống nhất.
2. Bổ sung cấu trúc ghi đo đạc cho từng lần chạy.
3. Gắn bộ đếm cho từng thuật toán.
4. Gắn đo bộ nhớ bằng hàm `runtime.ReadMemStats`.
5. Gắn nhật ký thời gian cho từng bước của quy trình.

### Giai đoạn 2. Dựng bộ khởi tạo dữ liệu

1. Khởi tạo dữ liệu nền cố định.
2. Khởi tạo bộ thử nhanh.
3. Khởi tạo bộ `nhỏ`, `trung bình`, `lớn`.
4. Khởi tạo các gói `K1` đến `K9`.

### Giai đoạn 3. Kiểm thử luật trước

1. Chạy `K1` đến `K9`.
2. Xác nhận từng lỗi sinh đúng loại xung đột.
3. Ghi bảng tiền xử lý và xung đột theo loại.

### Giai đoạn 4. Làm nóng

1. Với từng kịch bản và từng thuật toán, chạy nháp 3 đến 5 lần.
2. Không ghi các lần này vào báo cáo.

### Giai đoạn 5. Đo đạc chính thức

Với từng kịch bản `nhỏ`, `trung bình`, `lớn`:

1. Ghi ảnh chụp dữ liệu đầu vào.
2. Ghi phễu dữ liệu:
   - `L0`
   - `L1`
   - `L2`
   - `V`
   - `R0`
   - `R1`
3. Chạy thuật toán tô màu đồ thị 10 lần.
4. Chạy thuật toán CP-SAT 10 lần.
5. Chạy thuật toán tìm kiếm Tabu 10 lần.
6. Ở mỗi lần chạy phải lưu:
   - thời điểm bắt đầu;
   - thời điểm kết thúc;
   - thời lượng;
   - tệp nhật ký gốc;
   - tệp kết quả thô dạng dữ liệu có cấu trúc.
7. Tính trung bình và độ lệch chuẩn.
8. So sánh trên cùng một ngân sách thời gian.

### Giai đoạn 6. Đo lại khi điều chỉnh lịch

1. Chạy `K8` với chế độ giữ nguyên các buổi học đã khóa.
2. Chạy `K9` với chế độ điều chỉnh trên dữ liệu nháp.
3. Ghi:
  - số lần đổi lịch
  - số lần đổi giáo viên
  - số lần đổi phòng
  - xung đột với buổi học đã có

### Giai đoạn 7. Chỉnh tay và lưu thật

1. Chọn các kết quả xếp một phần.
2. Ghi số ca còn nhiều phương án lựa chọn.
3. Ghi số ca không còn phương án nào.
4. Chỉnh tay trên UI.
5. Đo lại điểm mềm.
6. Đo thời gian ghi xuống cơ sở dữ liệu.
7. Ghi số buổi lưu thành công thực tế.

## 13. Điểm cần chèn nhật ký trong mã nguồn

### 13.1 Ở tầng use case

- [preview.go](/Users/hant/golang/doan/internal/usecases/scheduling/preview.go)
  - thời gian nạp lớp
  - thời gian lọc ghi danh
  - thời gian nạp ca học
  - thời gian dựng cửa sổ lớp
  - thời gian nạp phòng
  - thời gian nạp bản đồ di chuyển
  - thời gian nạp buổi học đã có
  - thời gian giải bài toán
  - thời gian hậu xử lý bản xem trước

- [policy.go](/Users/hant/golang/doan/internal/usecases/scheduling/policy.go)
  - số lớp bị loại theo từng nguyên nhân
  - số ca vượt ngưỡng chỉnh tay

### 13.2 Ở tầng bộ giải dùng chung

- [default_preview_solver.go](/Users/hant/golang/doan/internal/services/scheduling/default_preview_solver.go)
  - số biến tạo ra
  - số phương án miền giá trị tạo ra
  - thống kê kích thước miền
  - lôgarit cơ số 10 của không gian tìm kiếm ước lượng

### 13.3 Ở từng thuật toán

- [graph_coloring_solver.go](/Users/hant/golang/doan/internal/services/scheduling/graph_coloring_solver.go)
  - tổng số phương án đã kiểm tra
  - tỷ lệ gán thành công ngay lượt đầu

- [cp_sat_solver.go](/Users/hant/golang/doan/internal/services/scheduling/cp_sat_solver.go)
  - số nút đã duyệt
  - số nhánh bị cắt
  - số lần cập nhật nghiệm tốt nhất
  - cờ chạm ngưỡng số nút tối đa

- [tabu_search_solver.go](/Users/hant/golang/doan/internal/services/scheduling/tabu_search_solver.go)
  - số vòng lặp
  - số bước chuyển được chấp nhận
  - số bước chuyển bị danh sách cấm loại bỏ
  - số gán được sửa ở bước phục hồi

## 14. Kết luận thực thi

Kế hoạch đo đạc cập nhật này nên được triển khai theo đúng thứ tự sau:

1. Chuẩn hóa điều kiện dừng.
2. Bổ sung lớp ghi đo đạc và đo bộ nhớ.
3. Tính không gian tìm kiếm theo lôgarit cơ số 10.
4. Dựng bộ khởi tạo dữ liệu chuẩn.
5. Chạy kiểm thử luật.
6. Làm nóng.
7. Chạy đo đạc chính thức 10 lần cho mỗi cặp `kịch bản x thuật toán`.
8. Tính trung bình, độ lệch chuẩn và tốc độ tính toán.
9. Chạy thêm đo đạc điều chỉnh lịch.
10. Chạy chỉnh tay và lưu thật.

Nếu muốn phần đo đạc có giá trị học thuật thật sự, ba thứ bắt buộc phải xuất hiện trong báo cáo là:

- quy mô dữ liệu và không gian tìm kiếm;
- phễu dữ liệu theo từng bước xử lý;
- đo đạc lặp lại nhiều lần với trung bình và độ lệch chuẩn;
- nhật ký thời gian của từng lần đo để làm minh chứng kiểm tra chéo.

Nếu bạn muốn, bước tiếp theo mình có thể làm tiếp một trong hai việc:

- dựng luôn bộ tập lệnh khởi tạo dữ liệu cơ sở dữ liệu theo ba kịch bản `nhỏ`, `trung bình`, `lớn`;
- viết sẵn bản thiết kế lược đồ ghi đo đạc và cấu hình đo đạc để bạn triển khai trực tiếp vào mã nguồn.
