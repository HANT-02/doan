# CÂU HỎI PHẢN BIỆN DỰ KIẾN CHO PHẦN SCHEDULING

## 1. Tại sao phải benchmark ba solver?

Vì mục tiêu của đề tài không chỉ là “có thuật toán chạy được”, mà là lựa chọn solver chính dựa trên dữ liệu thực nghiệm. Benchmark giúp chứng minh quyết định chọn `CP-SAT` có cơ sở, thay vì chọn theo cảm tính.

## 2. Tại sao không chọn Graph Coloring khi nó nhanh nhất?

`Graph Coloring + Heuristic` đúng là nhanh nhất, nhưng benchmark cho thấy ở scenario `small`, `CP-SAT` tạo nghiệm có soft score tốt hơn. Với đồ án này, tốc độ không phải tiêu chí duy nhất; cần cân bằng giữa tính đúng, chất lượng nghiệm và khả năng mở rộng.

## 3. Tại sao không chọn Tabu Search?

`Tabu Search` có giá trị học thuật vì đại diện cho metaheuristic, nhưng trong benchmark hiện tại nó không tạo ra lợi thế chất lượng đủ mạnh để bù cho runtime cao hơn, đặc biệt ở scenario `large`.

## 4. Tại sao benchmark dùng dữ liệu tổng hợp mà vẫn có giá trị?

Vì dữ liệu benchmark được thiết kế deterministic, có kiểm soát ràng buộc và bám theo logic domain thật của hệ thống như `Shift`, `ClassSchedule`, `Teacher`, `Room`, `Course`. Điều này giúp so sánh công bằng giữa các solver và tái lập được kết quả.

## 5. Hard constraints và soft constraints khác nhau như thế nào?

Hard constraints là các điều kiện bắt buộc, ví dụ không trùng giáo viên, phòng hoặc lớp ở cùng thời điểm. Nếu vi phạm hard constraint thì nghiệm không thể dùng. Soft constraints là các ưu tiên vận hành, ví dụ nghiệm “đẹp” hơn hoặc thuận tiện hơn. Benchmark dùng soft score để phân biệt chất lượng nghiệm khi các solver đều đã đúng về hard constraints.

## 6. Tại sao phải có bước preview rồi mới commit?

Preview giúp người vận hành kiểm tra kết quả trước khi ghi thật xuống `lesson`. Đây là chốt an toàn quan trọng vì lesson là dữ liệu đầu vào cho attendance, lesson summary và academic record.

## 7. Nếu dữ liệu thật khác dữ liệu benchmark thì kết luận còn đúng không?

Kết luận hiện tại đúng trong phạm vi benchmark của đồ án. Đây là cơ sở hợp lý để chọn solver mặc định ở giai đoạn hiện tại. Trong tương lai, khi có dữ liệu lesson và class thực hơn, benchmark có thể được mở rộng để kiểm tra lại kết luận trên dữ liệu gần production hơn.
