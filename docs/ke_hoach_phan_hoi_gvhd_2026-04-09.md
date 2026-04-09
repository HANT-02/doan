# Kế hoạch chi tiết xử lý phản hồi GVHD sau báo cáo ĐATN

**Ngày cập nhật:** 2026-04-09

## 1. Các quyết định đã chốt

Các quyết định sau đã được khóa để làm căn cứ triển khai:

1. Thực thể quản lý ca học sẽ dùng tên **`Shift`**.
2. `class_schedule` sẽ **chuyển hoàn toàn sang `shift_id`**.
3. Phần scheduling sẽ benchmark đúng **3 thuật toán**:
   - `Graph Coloring + heuristic`
   - `CP-SAT`
   - `Tabu Search`
4. Benchmark scheduling sẽ được mở thành **admin API nội bộ**.
5. Sau benchmark sẽ chọn **1 solver tốt nhất** để inject vào use case scheduling chính.
6. Predictive analytics chọn bài toán **classification**:
   - dự báo sinh viên có nguy cơ học kém
   - nhãn đầu ra: `AT_RISK` / `NOT_AT_RISK`
7. Predictive analytics sẽ được **huấn luyện bằng backend hiện tại**.
8. `System prompt` **không** được dùng để thay thế mô hình dự báo; prompt chỉ dùng để giải thích kết quả hoặc sinh khuyến nghị.
9. **Ưu tiên hoàn thiện scheduling trước**, predictive analytics làm sau.

---

## 2. Mục tiêu của đợt điều chỉnh

Theo phản hồi của GVHD, đồ án cần được nâng cấp theo 3 trục chính:

1. Scheduling phải có giá trị nghiên cứu:
   - có nhiều thuật toán,
   - có benchmark,
   - có tiêu chí đánh giá,
   - có lựa chọn thuật toán tốt nhất.

2. Hệ thống phải có **quản lý ca học** để dữ liệu thời gian cho scheduling được chuẩn hóa.

3. Hệ thống phải có một phần **AI dự báo** đủ đúng bản chất của machine learning, ở đây là dự báo sinh viên có nguy cơ học kém.

Mục tiêu của kế hoạch này là:

- giữ scheduling làm trục chính của đồ án,
- thêm `Shift` để scheduling có nền dữ liệu đúng,
- làm predictive analytics ở phạm vi đủ mạnh nhưng không quá rộng,
- đảm bảo vừa triển khai được, vừa viết báo cáo được.

---

## 3. Kiến trúc tổng thể sau điều chỉnh

Sau đợt nâng cấp này, hệ thống được tổ chức thành 3 khối chính:

1. **Khối quản lý ca học (`Shift`)**
   - định nghĩa các ca học chuẩn,
   - làm nguồn dữ liệu thời gian cho scheduling.

2. **Khối scheduling**
   - nhận dữ liệu lớp, giáo viên, phòng, ca học, ràng buộc,
   - hỗ trợ benchmark 3 solver qua admin API nội bộ,
   - chọn 1 solver chính cho API scheduling dùng bởi frontend.

3. **Khối predictive analytics**
   - lấy dữ liệu sinh viên, điểm, điểm danh, vận hành học tập,
   - huấn luyện classifier trong backend hiện tại,
   - trả nhãn `AT_RISK` và insight giải thích.

---

## 4. Workstream A - Scheduling là ưu tiên số 1

## 4.1. Mục tiêu

Biến scheduling từ mức "có preview/commit" thành một module có:

- mô hình đầu vào chuẩn hóa theo `Shift`,
- 3 thuật toán đại diện để benchmark,
- benchmark có số liệu,
- 1 solver chính được chọn bằng thực nghiệm.

## 4.2. Ba thuật toán scheduling sẽ benchmark

### Thuật toán 1. Graph Coloring + heuristic

Vai trò:

- baseline dễ giải thích,
- chạy nhanh,
- phù hợp làm mốc so sánh đầu tiên.

Ý nghĩa trong báo cáo:

- thể hiện hướng heuristic cơ bản,
- giúp chứng minh rằng nghiệm nhanh chưa chắc là nghiệm tốt nhất.

### Thuật toán 2. CP-SAT

Vai trò:

- đại diện cho hướng tối ưu hóa ràng buộc hiện đại,
- mạnh khi hard constraints được mô hình hóa rõ,
- dễ trình bày thành mô hình toán.

Ý nghĩa trong báo cáo:

- đại diện cho hướng exact/exact-like mạnh,
- dùng để so với heuristic và local search.

### Thuật toán 3. Tabu Search

Vai trò:

- đại diện cho local search/metaheuristic,
- phù hợp với bài toán timetabling,
- có thể cải thiện nghiệm từ lời giải khởi tạo ban đầu.

Ý nghĩa trong báo cáo:

- thể hiện hướng tìm kiếm lân cận có nhớ trạng thái,
- thực dụng hơn so với mở rộng quá nhanh sang RL hoặc mô hình quá nặng.

## 4.3. Vì sao không chọn Genetic Algorithm ở vòng triển khai chính

GA vẫn nên xuất hiện trong phần tổng quan tài liệu, nhưng chưa đưa vào nhóm triển khai chính vì:

- cần thiết kế chromosome, crossover, mutation, repair strategy,
- dễ tốn nhiều thời gian tuning,
- khó benchmark công bằng nếu cơ chế repair chưa đủ tốt,
- trong giai đoạn này scheduling còn phải đồng thời xử lý `Shift`, refactor `class_schedule`, API benchmark và API production.

Kết luận:

- **có nhắc GA trong phần nghiên cứu liên quan**,
- **không đưa GA vào nhóm 3 solver triển khai chính** ở vòng này.

## 4.4. Tiêu chí benchmark scheduling

Ba thuật toán sẽ được so sánh trên cùng input và cùng bộ ràng buộc theo các tiêu chí:

1. **Feasibility rate**
   - có tạo được lịch thỏa hard constraints hay không.

2. **Hard constraint violations**
   - trùng giáo viên,
   - trùng phòng,
   - trùng lớp,
   - vượt sức chứa,
   - xếp ngoài ca học hợp lệ.

3. **Soft score**
   - phân bố lịch đều hơn,
   - giảm khoảng trống,
   - hạn chế dồn lịch,
   - tôn trọng ưu tiên phòng hoặc ca nếu có.

4. **Runtime**
   - thời gian chạy trên từng tập dữ liệu.

5. **Scalability**
   - khi tăng số lớp, số ca, số phòng thì chất lượng và thời gian thay đổi thế nào.

6. **Stability**
   - chạy nhiều lần có ra nghiệm ổn định hay không.

## 4.5. Kết quả đầu ra kỳ vọng của scheduling benchmark

Workstream scheduling cần tạo ra:

- bảng benchmark 3 thuật toán,
- biểu đồ runtime/chất lượng nghiệm,
- mô hình hard constraints và soft constraints rõ ràng,
- lựa chọn 1 thuật toán tối ưu nhất cho hệ thống chính.

---

## 5. Workstream B - Quản lý ca học bằng thực thể `Shift`

## 5.1. Mục tiêu

Hiện tại scheduling còn bám nhiều vào giờ học trực tiếp. Cần thêm `Shift` để:

- chuẩn hóa không gian thời gian,
- giúp frontend/backend cùng hiểu chung khái niệm ca học,
- để solver làm việc trên tập ca học thay vì trên giờ rời rạc.

## 5.2. Thực thể `Shift`

Thực thể được chốt là `Shift` với các field đề xuất:

- `id`
- `code`
- `name`
- `start_time`
- `end_time`
- `duration_minutes`
- `session_type`
- `is_active`
- `notes`

Có thể mở rộng nếu cần:

- `campus`
- `building`
- `applies_for_grade`
- `applies_for_mode`

## 5.3. Vai trò của `Shift` trong scheduling

Sau khi có `Shift`, luồng scheduling sẽ là:

1. Chọn tập lớp cần xếp.
2. Chọn tập `Shift` đang hoạt động.
3. Sinh domain theo `Shift`.
4. Gán lớp vào `Shift`.
5. Gán giáo viên và phòng phù hợp.
6. Kiểm tra hard constraints.
7. Tính soft score.

Tức là:

- `Shift` trở thành đơn vị thời gian chuẩn,
- solver không còn làm việc trực tiếp trên mốc giờ rời rạc như trước.

## 5.4. Quyết định về `class_schedule`

`class_schedule` sẽ **chuyển hoàn toàn sang `shift_id`**.

Điều này có nghĩa:

- tầng dữ liệu chính của scheduling sẽ không còn xem `start_time/end_time` là nguồn chuẩn nữa,
- nếu cần hiển thị giờ bắt đầu/kết thúc, hệ thống sẽ suy ra từ `Shift`,
- dữ liệu lịch mẫu của lớp sẽ trở nên thống nhất và dễ kiểm tra hơn.

## 5.5. Phần cần triển khai trong hệ thống

1. CRUD `Shift`.
2. Màn quản lý `Shift` cho admin.
3. Refactor `class_schedule` để dùng `shift_id`.
4. Cập nhật scheduling solver để domain dùng `Shift`.
5. Cập nhật preview UI để hiển thị theo `Shift`.

## 5.6. Thứ tự triển khai module `Shift`

1. Thiết kế schema `shifts`.
2. Viết migration.
3. Bổ sung entity, repository, use case, controller.
4. Tạo UI quản lý `Shift`.
5. Refactor `class_schedule` sang `shift_id`.
6. Sửa scheduling input để dùng `Shift`.
7. Kiểm thử lại preview và commit.

---

## 6. Workstream C - Kiến trúc code cho benchmark scheduling

## 6.1. Cách tổ chức backend

Để tránh code lặp và benchmark công bằng, scheduling nên dùng một interface solver chung.

Khuyến nghị:

1. Tạo interface chung:
   - `SchedulingSolver`
   - `Solve(ctx, input) -> output`

2. Cài đặt 3 solver riêng:
   - `GraphColoringSolver`
   - `CPSATSolver`
   - `TabuSearchSolver`

3. Chuẩn hóa `SchedulingInput` và `SchedulingOutput` dùng chung cho:
   - benchmark,
   - scheduling API chính.

## 6.2. Admin API nội bộ cho benchmark

Benchmark sẽ được mở thành **admin API nội bộ**.

Vai trò của API này:

- nhận cùng một input benchmark,
- chạy qua cả 3 solver,
- ghi metric,
- trả bảng so sánh để phục vụ nghiên cứu và demo.

Khuyến nghị:

- frontend chính của scheduling **không** cho người dùng cuối chọn solver,
- benchmark API chỉ phục vụ admin/nghiên cứu,
- API scheduling chính chỉ dùng **1 solver đã được chọn** sau benchmark.

## 6.3. Cách dùng solver trong API chính

Sau benchmark:

- chọn 1 solver tốt nhất,
- inject solver đó vào use case scheduling chính,
- giữ `benchmark admin API` riêng cho mục đích nghiên cứu,
- frontend chính chỉ gọi **một API scheduling production-like**.

## 6.4. Lợi ích của kiến trúc này

1. Dễ benchmark công bằng.
2. Dễ viết báo cáo vì cùng một input và cùng một scorer.
3. Không làm frontend phức tạp không cần thiết.
4. Có thể thay solver chính mà ít ảnh hưởng contract UI.

---

## 7. Workstream D - Predictive Analytics

## 7.1. Bài toán đã chốt

Bài toán được chọn là:

**Dự báo sinh viên có nguy cơ học kém (`AT_RISK`)**

Đây là bài toán **classification**.

Đầu ra mong muốn:

- `AT_RISK`
- `NOT_AT_RISK`

## 7.2. Tập dữ liệu đầu vào nên dùng

Feature nên lấy từ các nhóm dữ liệu:

1. **Thông tin sinh viên**
   - khối lớp,
   - độ tuổi,
   - giới tính nếu có.

2. **Thông tin học tập**
   - điểm đầu vào,
   - điểm quá khứ,
   - điểm quiz/bài tập,
   - số môn hoặc lớp đang học.

3. **Thông tin điểm danh**
   - số buổi vắng,
   - tỷ lệ chuyên cần,
   - số lần đến muộn.

4. **Thông tin vận hành**
   - đổi lớp,
   - học phí chậm,
   - xin nghỉ học,
   - cường độ học theo tuần.

## 7.3. Các mô hình nên so sánh

Để vừa sức, nên so sánh:

1. `Logistic Regression`
   - baseline dễ giải thích.

2. `Random Forest`
   - baseline mạnh, ổn định.

3. `LightGBM` hoặc `XGBoost`
   - ứng viên chính cho dữ liệu tabular.

## 7.4. Mô hình chính khuyến nghị

Nếu dữ liệu ở dạng bảng là chủ yếu, mô hình chính nên là:

- `LightGBM` hoặc `XGBoost`

và dùng thêm:

- `SHAP` hoặc feature importance để giải thích vì sao sinh viên bị gắn nhãn `AT_RISK`.

## 7.5. Thước đo đánh giá

Vì đây là classification, cần đánh giá bằng:

- Accuracy
- Precision
- Recall
- F1-score
- ROC-AUC

Khuyến nghị:

- ưu tiên `Recall` và `F1` cho lớp `AT_RISK`,
- vì bỏ sót sinh viên nguy cơ cao gây hại hơn so với cảnh báo dư.

## 7.6. Predictive analytics trong backend hiện tại

Do đã chốt huấn luyện bằng backend hiện tại, kế hoạch sẽ theo hướng:

1. Xây dựng pipeline feature engineering trong backend.
2. Tạo job huấn luyện nội bộ trong backend cho các mô hình baseline và mô hình chính.
3. Lưu metadata mô hình, metric, version mô hình trong DB hoặc file nội bộ.
4. Expose prediction endpoint từ backend hiện tại.
5. Nếu cần giải thích, backend trả thêm feature importance hoặc kết quả giải thích đã được tính trước.

Lưu ý:

- đây là hướng nhất quán về kiến trúc,
- nhưng phải giữ phạm vi mô hình vừa phải để phù hợp với backend hiện tại.

## 7.7. Vì sao không dùng system prompt thay cho mô hình dự báo

`System prompt` không nên thay thế mô hình predictive analytics vì:

1. Không có bước huấn luyện trên dữ liệu lịch sử.
2. Không có metric đánh giá chuẩn như Accuracy, Recall, F1, ROC-AUC.
3. Khó chứng minh tính ổn định và khả năng tổng quát hóa.
4. Kết quả phụ thuộc mạnh vào prompt.
5. Khó bảo vệ trước GVHD nếu gọi đó là "mô hình dự báo".

Kết luận:

- lõi dự báo phải là classifier thực sự,
- prompt/LLM chỉ nên là lớp giải thích và sinh khuyến nghị.

## 7.8. Kết quả đầu ra cần có

1. Bảng so sánh 3 mô hình.
2. Confusion matrix.
3. Feature importance hoặc SHAP.
4. Danh sách sinh viên nguy cơ cao.
5. Màn hình cảnh báo/insight ở mức demo.

---

## 8. Kế hoạch triển khai theo giai đoạn

## Giai đoạn 1. Chốt kiến trúc scheduling và benchmark framework

Mục tiêu:

- khóa kiến trúc solver,
- khóa 3 thuật toán benchmark,
- khóa benchmark admin API.

Việc làm:

1. Chốt `Graph Coloring + CP-SAT + Tabu Search`.
2. Thiết kế interface `SchedulingSolver`.
3. Thiết kế `SchedulingInput` và `SchedulingOutput`.
4. Thiết kế metric benchmark.
5. Thiết kế contract cho benchmark admin API.

Deliverable:

- tài liệu scope scheduling,
- sơ đồ kiến trúc solver,
- contract benchmark API.

## Giai đoạn 2. Bổ sung module `Shift`

Mục tiêu:

- có module `Shift` dùng được trong admin và sẵn sàng cho scheduling.

Việc làm:

1. Thiết kế bảng `shifts`.
2. Viết migration.
3. CRUD backend.
4. UI quản lý `Shift`.
5. Refactor `class_schedule` sang `shift_id`.

Deliverable:

- module quản lý `Shift` hoàn chỉnh.

## Giai đoạn 3. Chuẩn hóa scheduling theo `Shift`

Mục tiêu:

- đưa scheduling sang chạy trên dữ liệu `Shift`.

Việc làm:

1. Sửa solver input để nhận `Shift`.
2. Chuẩn hóa hard constraints theo `Shift`.
3. Chuẩn hóa soft constraints.
4. Chuẩn hóa scorer và violation report.

Deliverable:

- scheduling domain model mới dựa trên `Shift`.

## Giai đoạn 4. Cài đặt 3 solver scheduling

Mục tiêu:

- có 3 solver chạy được trên cùng input.

Việc làm:

1. Cài `GraphColoringSolver`.
2. Cài `CPSATSolver`.
3. Cài `TabuSearchSolver`.
4. Nối cả 3 solver vào benchmark admin API.

Deliverable:

- 3 solver có thể chạy độc lập,
- benchmark admin API dùng được.

## Giai đoạn 5. Benchmark và chọn solver chính

Mục tiêu:

- có số liệu để chọn thuật toán tối ưu nhất.

Việc làm:

1. Chạy benchmark trên cùng bộ dữ liệu.
2. Ghi runtime, feasibility, hard violations, soft score.
3. So sánh kết quả.
4. Chọn 1 solver tốt nhất.

Deliverable:

- bảng benchmark,
- biểu đồ so sánh,
- quyết định solver chính.

## Giai đoạn 6. Tích hợp solver chính vào API và UI scheduling

Mục tiêu:

- scheduling production-like dùng thuật toán đã chọn.

Việc làm:

1. Inject solver chính vào use case scheduling.
2. Giữ benchmark admin API riêng cho nghiên cứu.
3. Cập nhật preview/commit/timetable.
4. Kiểm thử lại toàn bộ luồng scheduling.

Deliverable:

- module scheduling hoàn chỉnh hơn và dùng được trong demo.

## Giai đoạn 7. Làm predictive analytics trong backend hiện tại

Mục tiêu:

- có một pipeline classification hoàn chỉnh.

Việc làm:

1. Gom dữ liệu lịch sử.
2. Làm sạch dữ liệu.
3. Tạo feature.
4. Huấn luyện `Logistic Regression`.
5. Huấn luyện `Random Forest`.
6. Huấn luyện `LightGBM/XGBoost`.
7. Đánh giá.
8. Lưu model metadata.
9. Expose prediction endpoint.
10. Thêm explanation/insight.

Deliverable:

- mô hình `AT_RISK` classification,
- bảng so sánh mô hình,
- API dự báo,
- UI cảnh báo mức demo.

---

## 9. Thứ tự ưu tiên triển khai thực tế

Do đã chốt **ưu tiên scheduling**, thứ tự triển khai thực tế nên là:

1. Chốt benchmark framework và admin API cho scheduling.
2. Làm module `Shift`.
3. Refactor `class_schedule` sang `shift_id`.
4. Chuẩn hóa scheduling domain theo `Shift`.
5. Cài 3 solver.
6. Benchmark và chọn solver chính.
7. Gắn solver tốt nhất vào API/frontend.
8. Sau khi scheduling ổn định mới triển khai predictive analytics.

---

## 10. Các điểm cần review tiếp trước khi chuyển thành backlog kỹ thuật

Các quyết định lớn hiện đã khóa. Phần còn cần review là mức chi tiết triển khai:

1. Benchmark admin API có cần lưu lịch sử kết quả benchmark vào DB hay chỉ trả response tức thời.
2. `Shift` có cần hỗ trợ phân loại theo cơ sở/khu học hay chưa.
3. Predictive analytics trong backend hiện tại sẽ lưu model theo file nội bộ hay lưu hoàn toàn metadata trong DB.
4. Phần explanation cho `AT_RISK` sẽ dùng feature importance tĩnh hay sinh insight theo từng sinh viên.

---

## 11. Tài liệu tham khảo định hướng

### Scheduling

- Meta-heuristics và university course timetabling:
  [https://doi.org/10.1016/j.iswa.2023.200253](https://doi.org/10.1016/j.iswa.2023.200253)

- Hyper-heuristic cho university course timetabling:
  [https://doi.org/10.1016/j.asoc.2025.114150](https://doi.org/10.1016/j.asoc.2025.114150)

- Tabu Search cho timetabling:
  [https://doi.org/10.1016/j.cor.2020.105007](https://doi.org/10.1016/j.cor.2020.105007)

### Predictive analytics cho giáo dục

- Survey academic performance prediction:
  [https://doi.org/10.1109/TLT.2025.3554174](https://doi.org/10.1109/TLT.2025.3554174)

- LightGBM + SHAP cho educational prediction:
  [https://doi.org/10.3390/app152010875](https://doi.org/10.3390/app152010875)

- So sánh mô hình ML cho student performance prediction:
  [https://doi.org/10.3390/app15158409](https://doi.org/10.3390/app15158409)

---

## 12. Kết luận ngắn

Hướng triển khai tối ưu cho đồ án ở thời điểm này là:

1. thêm module `Shift`,
2. chuyển `class_schedule` sang `shift_id`,
3. benchmark `Graph Coloring + CP-SAT + Tabu Search`,
4. chọn 1 solver tốt nhất để tích hợp,
5. giữ predictive analytics ở phạm vi `AT_RISK classification`,
6. dùng mô hình ML thật trong backend hiện tại, còn prompt chỉ làm lớp giải thích.

Đây là hướng vừa bám sát góp ý của GVHD, vừa đủ chiều sâu kỹ thuật, vừa giữ khối lượng triển khai ở mức hợp lý cho đồ án tốt nghiệp.
