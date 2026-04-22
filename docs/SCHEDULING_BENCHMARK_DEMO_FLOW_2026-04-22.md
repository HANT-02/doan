# KỊCH BẢN DEMO BENCHMARK VÀ XẾP LỊCH PRODUCTION-LIKE

## 1. Mục tiêu

Tài liệu này dùng để hỗ trợ `Task H4` trong phần scheduling, giúp trình bày ngắn gọn nhưng đủ thuyết phục hai luồng demo:

1. demo benchmark ba solver;
2. demo luồng preview -> commit -> lesson theo solver mặc định `CP-SAT`.

---

## 2. Demo flow benchmark scheduling

### Bước 1. Giới thiệu bài toán

- Nêu ngắn gọn scheduling là điểm nhấn kỹ thuật của đồ án.
- Nhấn mạnh hệ thống không chỉ có CRUD mà còn benchmark nhiều solver trên cùng một chuẩn input.
- Giới thiệu ba solver:
  - `Graph Coloring + Heuristic`
  - `CP-SAT`
  - `Tabu Search`

### Bước 2. Chạy artifact benchmark

Lệnh demo:

```bash
env GOCACHE=/tmp/go-build go run ./cmd/cli/scheduling_benchmark
```

Điểm cần nói khi chạy:

- benchmark dùng ba scenario `small`, `medium`, `large`;
- mỗi scenario chạy lặp `7 lần`;
- kết quả ghi nhận `feasibility`, `hard violation`, `soft score`, `runtime`, `stability`.

### Bước 3. Giải thích bảng kết quả

Các ý chính cần nhấn:

- cả ba solver đều sạch hard constraints;
- `Graph Coloring` là baseline rất nhanh;
- `CP-SAT` cho chất lượng nghiệm tốt nhất ở scenario `small`;
- `Tabu Search` không tạo ra ưu thế đủ mạnh để thắng `CP-SAT`;
- vì vậy `CP-SAT` được chọn làm solver chính cho API scheduling.

### Bước 4. Chốt thông điệp

Thông điệp kết luận:

> Đề tài không chọn thuật toán theo cảm tính. Hệ thống có benchmark độc lập, có dữ liệu small/medium/large, có số liệu lặp lại được và từ đó chọn `CP-SAT` làm solver chính.

---

## 3. Demo flow scheduling production-like

### Bước 1. Mở màn hình scheduling preview

Các điểm cần chỉ:

- bộ lọc lớp / giáo viên / phòng / khoảng ngày;
- dữ liệu thời gian đã chuẩn hóa qua `Shift`;
- `ClassSchedule` là nền để sinh slot hợp lệ.

### Bước 2. Chạy preview

Điểm cần nói:

- preview sinh assignments và conflicts;
- solver mặc định hiện tại là `CP-SAT`;
- hệ thống không commit trực tiếp mà luôn đi qua bước preview để kiểm soát rủi ro vận hành.

### Bước 3. Giải thích kết quả preview

Các thành phần cần chỉ:

- số buổi đã xếp;
- số conflict / unscheduled lesson;
- danh sách assignment theo lớp, ca, phòng, giáo viên.

### Bước 4. Commit preview

Điểm cần nói:

- chỉ commit khi preview ở trạng thái `COMPLETED`;
- hệ thống kiểm tra lesson trùng lịch trước khi ghi vào DB;
- sau commit, `lesson` trở thành dữ liệu đầu vào cho attendance, lesson summary và academic record.

### Bước 5. Mở màn hình lesson

Mục đích:

- chứng minh scheduling không dừng ở benchmark;
- kết quả preview đã trở thành dữ liệu vận hành thật của hệ thống.

---

## 4. Danh sách hình minh chứng cần chụp

### Cho benchmark

1. Ảnh terminal chạy `go run ./cmd/cli/scheduling_benchmark`
2. Ảnh bảng benchmark có đủ `small / medium / large`
3. Ảnh biểu đồ runtime
4. Ảnh biểu đồ soft score scenario `small`

### Cho production-like scheduling

1. Ảnh màn hình scheduling preview trước khi chạy
2. Ảnh màn hình scheduling preview sau khi có assignments
3. Ảnh commit preview thành công
4. Ảnh màn hình lesson list sau commit

---

## 5. Khuyến nghị về thời lượng trình bày

- Benchmark scheduling: `2-3 phút`
- Preview -> commit -> lesson: `2-3 phút`
- Tổng phần scheduling: `5-6 phút`

Đây là thời lượng đủ để thể hiện chiều sâu kỹ thuật mà vẫn giữ nhịp bảo vệ gọn.
