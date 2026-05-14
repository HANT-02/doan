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
- ngưỡng mở lớp tối thiểu để xác định lớp có đủ điều kiện được xếp lịch hay không;
- danh sách phòng học và sức chứa;
- danh sách ca học `Shift`;
- lịch tuần mẫu `ClassSchedule`;
- khoảng thời gian chạy xếp lịch;
- các lesson đã tồn tại để tránh đè lịch;
- trạng thái vòng đời của từng buổi học: `History`, `Published`, `Draft`, `Unplanned`;
- dữ liệu kỹ năng giáo viên, campus/phòng học, thời gian di chuyển và các ngoại lệ vận hành nếu có;
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

### 5.4 Điều kiện mở lớp trước khi đưa vào scheduling

Không phải lớp nào được tạo ra cũng nên được đưa ngay vào bài toán xếp lịch.

Nên bổ sung rule vận hành:

- chỉ các lớp đạt ngưỡng sĩ số tối thiểu mới được đưa vào scheduling;
- ngưỡng này có thể cấu hình theo trung tâm hoặc theo khóa học;
- mốc khuyến nghị ban đầu là **`80% - 90%` sức chứa mục tiêu** hoặc tối thiểu theo cấu hình của lớp;
- lớp chưa đủ sĩ số thì giữ ở trạng thái chờ mở lớp, không tham gia preview và không làm nhiễu benchmark.

Lợi ích:

- tránh sinh lịch cho lớp chưa chắc mở;
- giảm số assignment vô ích;
- làm dữ liệu benchmark sát bài toán vận hành thực hơn.

### 5.5 Vòng đời lịch học trong production

Trong production, lịch học không nên chỉ có hai trạng thái "preview" và "committed". Cần chia vùng lịch để kiểm soát mức độ được phép thay đổi:

- `History`: buổi đã diễn ra, chỉ dùng để xem lại, điểm danh, báo cáo; không cho solver thay đổi.
- `Published`: lịch đã công bố cho giáo viên và học viên; chỉ được đổi khi có nghiệp vụ thay thế hoặc tái lập lịch có kiểm soát.
- `Draft`: vùng nháp để solver chạy tự động, kéo thả, tính điểm và chuẩn bị công bố.
- `Unplanned`: nhu cầu học chưa được gán vào ngày/ca/phòng/giáo viên cụ thể.

Nguyên tắc production nên chốt:

- solver tự động chỉ chạy trực tiếp trên vùng `Draft`;
- `Published` là vùng bán cố định, nếu phải đổi thì cần ghi nhận lý do và cảnh báo ảnh hưởng;
- `History` không được đưa vào tái lập lịch;
- mọi thay đổi sau khi lịch đã công bố phải sinh log để truy vết.

### 5.6 Nhóm nghiệp vụ thời gian thực

Sau khi lịch đã commit, hệ thống cần xử lý các biến động vận hành hằng ngày.

#### Giáo viên dạy thay

Khi giáo viên báo nghỉ đột xuất, hệ thống nên có luồng đề xuất giáo viên thay thế thay vì bắt giáo vụ tìm thủ công.

Luồng gợi ý:

1. chọn lesson cần thay giáo viên;
2. hệ thống lọc giáo viên có kỹ năng phù hợp với khóa học/môn học;
3. loại giáo viên đang bận trong cùng khung giờ;
4. tính tải công việc hiện tại để tránh dồn lịch;
5. nếu có nhiều campus thì tính thêm thời gian di chuyển từ buổi trước hoặc tới buổi sau;
6. trả về danh sách đề xuất trong thời gian ngắn, mục tiêu vận hành khoảng `30 giây`.

Tiêu chí xếp hạng giáo viên thay thế:

- khớp chuyên môn hoặc chứng chỉ bắt buộc;
- đang trống lịch;
- tải dạy trong tuần chưa quá cao;
- thời gian di chuyển giữa các cơ sở còn đủ;
- ưu tiên giáo viên từng dạy lớp hoặc cùng chương trình nếu có.

#### Bảo lưu, chuyển lớp và học bù

Production cần quản lý sức chứa theo thời gian thực, không chỉ kiểm tra sức chứa lúc tạo lịch.

Các nghiệp vụ cần có:

- bảo lưu: học viên tạm dừng, slot của lớp có thể giảm tải hoặc mở chỗ học bù;
- chuyển lớp: kiểm tra lớp đích còn sức chứa, cùng trình độ, cùng chương trình hoặc khóa tương đương;
- học bù: tìm `spot` trống trong các lớp tương thích mà không vượt sức chứa phòng;
- tỷ lệ lấp đầy: theo dõi `capacity_utilization` theo lớp, phòng, ca học và tuần.

Luồng `Find Spot` cho học bù:

1. xác định học viên, lớp gốc, buổi vắng và nội dung cần học bù;
2. tìm các lớp cùng trình độ/chương trình đang chạy;
3. lọc các lesson còn chỗ theo sức chứa phòng và sĩ số thực tế;
4. kiểm tra học viên không bị trùng lịch;
5. xếp hạng slot theo mức phù hợp nội dung, thời gian, cơ sở và độ đầy lớp.

### 5.7 Lập lịch liên tục và tái lập lịch

Hệ thống production không nên chỉ chạy lịch một lần cho cả kỳ. Nên có cơ chế lập lịch cuốn chiếu:

- lập lịch tự động theo từng cửa sổ thời gian, ví dụ 2-4 tuần phía trước;
- vùng gần hiện tại được công bố sớm cho giáo viên/học viên;
- vùng xa hơn vẫn giữ ở `Draft` để dễ điều chỉnh;
- các nhu cầu mới phát sinh được đưa vào `Unplanned` rồi xử lý ở vòng kế tiếp.

Khi có biến động như phòng học hỏng, giáo viên nghỉ dài ngày hoặc đổi giáo viên giữa chừng, hệ thống cần hỗ trợ tái lập lịch ít xáo trộn nhất.

Nguyên tắc `Nonvolatile Replanning`:

- giữ nguyên tối đa các buổi đã `Published`;
- nếu buộc phải đổi lịch đã công bố thì phạt điểm rất nặng;
- ưu tiên thay đổi cục bộ quanh các lesson bị ảnh hưởng;
- chỉ mở rộng vùng tái lập lịch khi không tìm được phương án hợp lệ trong vùng nhỏ;
- mọi thay đổi so với lịch cũ cần được thống kê thành `change set`.

Các chỉ số cần đo khi tái lập lịch:

- số lesson bị đổi ngày/ca/phòng/giáo viên;
- số học viên bị ảnh hưởng;
- số giáo viên bị ảnh hưởng;
- điểm giảm xáo trộn so với phương án chạy lại toàn bộ;
- thời gian xử lý.

## 6. Cách xử lý xung đột để ít phải sửa tay nhất

Mục tiêu không phải triệt tiêu 100% thao tác tay, mà là giảm tối đa số lần phải sửa bằng tay.

Nguyên tắc nên chốt:

- solver chịu trách nhiệm tạo phương án tự động tốt nhất;
- preview cho phép người dùng kiểm tra trước khi commit;
- nếu còn conflict hoặc unscheduled session thì không commit;
- chỉnh tay chỉ là bước tinh chỉnh cuối, không phải bước vận hành chính.
- nếu mức độ xung đột quá dày đặc thì nên xem là dữ liệu đầu vào chưa tốt, không đẩy gánh nặng sang người xếp thủ công.

### 6.1 Có nên kéo thả trên UI không?

Có, nhưng kéo thả chỉ nên là **lớp điều chỉnh trên preview**, không thay thế solver.

Luồng đúng nên là:

1. người dùng chạy auto scheduling;
2. hệ thống hiển thị preview theo tuần;
3. nếu một buổi cần đổi slot/phòng thì người dùng kéo thả **các chip assignment** hoặc chọn lại candidate;
4. hệ thống re-validate ngay luật cứng;
5. chỉ cho commit khi toàn bộ preview hợp lệ.

Nói ngắn gọn:

- **UI kéo thả là hợp lý**, nhưng phải đi kèm kiểm tra ràng buộc tức thời.
- Không nên cho sửa tay trực tiếp trên `Lesson` sau commit như một cách vá lỗi chính.

### 6.1.1 Tính điểm delta khi kéo thả

Ở mức production, kéo thả không chỉ dừng ở kiểm tra luật cứng.

Mỗi lần người dùng kéo một chip assignment sang slot khác, UI nên hiển thị ngay:

- hard validation: có trùng giáo viên, phòng, lớp, học viên hoặc thời gian di chuyển hay không;
- `score_delta`: điểm luật mềm tăng/giảm bao nhiêu so với trạng thái trước khi kéo;
- nguyên nhân điểm giảm, ví dụ giáo viên bị dạy liên tiếp quá nhiều, phải di chuyển xa, lớp bị đổi ca không ổn định;
- cảnh báo chất lượng nếu thao tác hợp lệ nhưng làm lịch xấu đi rõ rệt.

Cách hiển thị nên ưu tiên nhanh và dễ đọc:

- màu xanh nếu điểm tốt hơn;
- màu vàng nếu hợp lệ nhưng điểm giảm nhẹ;
- màu đỏ nếu vi phạm luật cứng hoặc giảm chất lượng mạnh;
- ghi rõ `+/- soft score` trên chip hoặc panel bên cạnh.

Điều này giúp người xếp thủ công không chỉ biết "được hay không", mà còn biết "đổi như vậy có đáng không".

### 6.2 Khi nào xem một preview là "trùng quá nhiều"?

Đây nên được xem là tiêu chí vận hành chứ không chỉ là cảm giác người dùng.

Nên bổ sung ngưỡng cảnh báo:

- nếu số conflict vượt quá ngưỡng cho phép thì preview bị xem là không đạt;
- nếu tỷ lệ session chưa xếp hoặc session đang conflict quá cao thì không cho commit;
- mục tiêu là không để người xếp thủ công phải vá quá nhiều vì như vậy hiệu suất vận hành giảm mạnh và khó bảo đảm tính đúng.

Có thể theo dõi bằng các chỉ số:

- `conflict_count`
- `unscheduled_session_count`
- `manual_adjustment_count`
- `manual_adjustment_ratio`

Nếu các chỉ số này vượt ngưỡng, nên yêu cầu:

- nới khoảng ngày xếp;
- thêm `Shift`;
- thêm phòng;
- đổi giáo viên;
- hoặc điều chỉnh lịch tuần lớp trước khi chạy lại solver.

## 7. Bộ luật đánh giá thuật toán

## 7.1 Luật cứng

Đây là lớp điều kiện bắt buộc. Thuật toán nào vi phạm luật cứng thì bị loại trước khi so tiếp.

Danh sách đề xuất:

- không trùng giáo viên ở cùng thời điểm;
- không trùng phòng ở cùng thời điểm;
- không trùng lớp ở cùng thời điểm;
- giáo viên phải có kỹ năng hoặc chứng chỉ phù hợp với môn học/khóa học;
- nếu có nhiều campus, không xếp giáo viên sang cơ sở khác khi khoảng nghỉ nhỏ hơn thời gian di chuyển tối thiểu;
- chỉ scheduling cho các lớp đã đạt điều kiện mở lớp tối thiểu;
- chỉ dùng `Shift` đang hoạt động;
- chỉ xếp vào slot sinh ra từ `ClassSchedule`;
- phòng phải đủ sức chứa;
- học bù hoặc chuyển lớp không được làm vượt sức chứa phòng/lớp tại thời điểm lesson diễn ra;
- số buổi xếp được không vượt nhu cầu của khóa học;
- buổi học phải nằm trong khoảng ngày cho phép;
- không đè lên lesson đã tồn tại hoặc slot đã khóa;
- nếu có ràng buộc giáo viên nghỉ cố định thì phải tôn trọng.

## 7.2 Luật mềm

Chỉ so luật mềm sau khi thuật toán đã sạch luật cứng.

Danh sách đề xuất:

- hạn chế khoảng trống lớn giữa các ca dạy của cùng giáo viên;
- hạn chế dồn quá nhiều buổi liên tiếp cho một giáo viên;
- cân bằng tải dạy giữa các giáo viên có cùng năng lực;
- hạn chế xếp giáo viên phải di chuyển giữa campus trong thời gian quá sát nhau;
- phạt rất nặng việc thay đổi lesson đã `Published` để tái lập lịch ít xáo trộn;
- hạn chế xếp hai môn nặng cho cùng một nhóm học viên trong cùng một ngày;
- ưu tiên phòng quen thuộc hoặc phòng tối ưu sức chứa;
- ưu tiên phân bố lịch đều trong tuần cho lớp;
- hạn chế đổi ca liên tục giữa các buổi của cùng lớp;
- hạn chế dùng các ca ít mong muốn nếu có ca ưu tiên hơn;
- giảm số conflict phải sửa tay trên preview;
- giảm số thao tác kéo thả/chỉnh tay cần thiết trước khi commit;
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
- kỹ năng/chứng chỉ của giáo viên và yêu cầu kỹ năng của khóa học;
- nhiều học sinh để phản ánh quy mô vận hành;
- nhiều phòng với sức chứa khác nhau;
- nhiều campus hoặc nhóm phòng để mô phỏng thời gian di chuyển;
- nhiều ca học `Shift`;
- lịch tuần đa dạng giữa các lớp;
- một số slot bị khóa hoặc lesson có sẵn để mô phỏng dữ liệu thực;
- trạng thái lesson gồm `Published`, `Draft`, `Unplanned` để kiểm thử tái lập lịch;
- tình huống giáo viên nghỉ, phòng hỏng, học viên cần học bù;
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
- `score_delta` khi chỉnh tay hoặc tái lập lịch;
- `schedule_change_count` khi replanning;
- `manual_adjustment_count`
- `capacity_utilization`
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

## 9.4 Benchmark cho production operations

Ngoài benchmark cold start, nên có thêm nhóm benchmark production-like:

- `substitution`: giáo viên nghỉ đột xuất, đo thời gian tìm giáo viên thay thế và chất lượng đề xuất;
- `makeup_find_spot`: học viên cần học bù, đo khả năng tìm slot còn chỗ và không trùng lịch học viên;
- `replanning_room_failure`: phòng học bị hỏng, đo số lesson bị đổi và mức xáo trộn;
- `published_stability`: kiểm tra solver có giữ ổn định lịch đã công bố hay không.

Các benchmark này giúp chứng minh hệ thống không chỉ biết tạo lịch ban đầu, mà còn có hướng vận hành liên tục sau khi lịch đã đi vào thực tế.

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
- bổ sung rule điều kiện mở lớp theo ngưỡng sĩ số tối thiểu;
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
- hỗ trợ kéo thả các chip assignment theo tuần cho các ca bị trùng;
- mỗi lần chỉnh phải re-check hard constraints ngay;
- khóa commit khi còn conflict.

## Giai đoạn 5.1 - Đồng bộ UI giữa preview và lịch đã xếp

Sau khi commit, giao diện lịch học đã xếp nên giữ cùng ngôn ngữ hiển thị với preview để người dùng không phải học lại cách đọc lịch.

Cần chốt:

- màn hình lịch đã xếp dùng cùng cấu trúc tuần, màu sắc, chip và cách đọc slot như preview;
- `teacher view` và `student view` cũng hiển thị theo layout tương tự;
- các thành phần như ca học, phòng học, giáo viên, trạng thái conflict trước commit và lesson sau commit cần có quy ước UI thống nhất;
- nếu preview hỗ trợ kéo thả chip thì màn lịch đã xếp ít nhất phải đồng bộ về cách hiển thị chip, dù không nhất thiết cho sửa trực tiếp theo cùng quyền.

## Giai đoạn 5.2 - Real-time operations

- bổ sung use case giáo viên dạy thay;
- bổ sung use case tìm chỗ học bù;
- tính `capacity_utilization` theo lesson/lớp/phòng/ca;
- thêm trạng thái vòng đời lịch: `History`, `Published`, `Draft`, `Unplanned`;
- thêm cơ chế ghi nhận `change set` khi thay đổi lịch đã công bố.

## Giai đoạn 5.3 - Continuous replanning

- giới hạn solver chạy tự động trên vùng `Draft`;
- mô hình hóa `Published` như semi-movable entities;
- thêm luật mềm phạt nặng khi thay đổi lịch đã công bố;
- benchmark số lượng thay đổi khi tái lập lịch;
- hiển thị score delta và impact summary cho người xếp lịch trước khi xác nhận.

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
- mô hình vòng đời lịch học `History/Published/Draft/Unplanned`;
- mô tả real-time operations: dạy thay, học bù, chuyển lớp, bảo lưu;
- benchmark report có số liệu lặp lại được;
- benchmark production-like cho substitution, makeup và replanning;
- solver mặc định `CP-SAT`;
- rule mở lớp theo ngưỡng sĩ số được mô tả rõ;
- màn preview đủ để xem, chỉnh và commit;
- UI kéo thả có hard validation và score delta;
- màn lịch đã xếp đồng bộ UI với preview cho admin, teacher và student;
- phần học máy chỉ còn ở mục hướng phát triển.

## 13. Kết luận định hướng

Từ thời điểm này, toàn bộ câu chuyện của đồ án nên xoay quanh một luận điểm chính:

> Hệ thống không chỉ quản lý dữ liệu trung tâm dạy thêm, mà giải quyết bài toán xếp lịch học bằng nhiều thuật toán, đo đạc được bằng dữ liệu thực nghiệm, và lựa chọn `CP-SAT` làm solver chính dựa trên kết quả benchmark.

Đó là trục thuyết phục nhất để viết báo cáo, demo hệ thống và trả lời phản biện.
