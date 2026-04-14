# Scheduling Benchmark Report

**Ngày chạy benchmark:** 2026-04-14

**Lệnh chạy chuẩn:**

```bash
go run ./cmd/cli/scheduling_benchmark
```

## 1. Mục tiêu

Task E2 cần tạo ra số liệu thực nghiệm để chọn solver chính cho scheduling production-like ở giai đoạn E3.

Benchmark được chạy trên cùng bộ input và cùng hard/soft constraint cho 3 solver:

- `graph_coloring`
- `cp_sat`
- `tabu_search`

Các tiêu chí so sánh bám theo định hướng đã chốt trong `docs/ke_hoach_phan_hoi_gvhd_2026-04-09.md`:

- feasibility rate
- hard violations
- soft score
- runtime
- scalability
- stability

## 2. Bộ dữ liệu benchmark

Benchmark study hiện dùng bộ dữ liệu tổng hợp có thể chạy lặp lại, được định nghĩa trong:

- [benchmark_study.go](/Users/hant/golang/doan/internal/services/scheduling/benchmark_study.go)

Ba scenario chuẩn:

1. `small`
   - 6 lớp
   - 4 giáo viên
   - 3 phòng
   - 3 ca
   - 12 requested sessions
   - 7 iterations

2. `medium`
   - 10 lớp
   - 5 giáo viên
   - 4 phòng
   - 3 ca
   - 30 requested sessions
   - 7 iterations

3. `large`
   - 16 lớp
   - 7 giáo viên
   - 5 phòng
   - 3 ca
   - 64 requested sessions
   - 7 iterations

Các scenario đều dùng:

- `Shift` chuẩn hóa theo sáng / chiều / tối
- `ClassSchedule` để giới hạn slot theo ngày + ca
- teacher reuse và preferred room có kiểm soát để tạo áp lực constraint
- dữ liệu deterministic để stability có thể đối chiếu giữa nhiều lần chạy

## 3. Kết quả benchmark

### 3.1 Scenario `small`

| Solver | Avg feasibility | Avg hard violations | Avg soft score | Avg runtime (ms) | Runtime range (ms) | Stable | Status |
| --- | ---: | ---: | ---: | ---: | --- | --- | --- |
| CP-SAT | 1.000 | 0.000 | 0.000 | 101.000 | 95-109 | true | COMPLETED |
| Tabu Search | 1.000 | 0.000 | -30.000 | 28.571 | 28-31 | true | COMPLETED |
| Graph Coloring + Heuristic | 1.000 | 0.000 | -45.000 | 0.000 | 0-0 | true | COMPLETED |

### 3.2 Scenario `medium`

| Solver | Avg feasibility | Avg hard violations | Avg soft score | Avg runtime (ms) | Runtime range (ms) | Stable | Status |
| --- | ---: | ---: | ---: | ---: | --- | --- | --- |
| Graph Coloring + Heuristic | 1.000 | 0.000 | 0.000 | 0.000 | 0-0 | true | COMPLETED |
| Tabu Search | 1.000 | 0.000 | 0.000 | 236.857 | 218-297 | true | COMPLETED |
| CP-SAT | 1.000 | 0.000 | 0.000 | 294.286 | 270-419 | true | COMPLETED |

### 3.3 Scenario `large`

| Solver | Avg feasibility | Avg hard violations | Avg soft score | Avg runtime (ms) | Runtime range (ms) | Stable | Status |
| --- | ---: | ---: | ---: | ---: | --- | --- | --- |
| Graph Coloring + Heuristic | 1.000 | 0.000 | 0.000 | 1.143 | 1-2 | true | COMPLETED |
| CP-SAT | 1.000 | 0.000 | 0.000 | 871.571 | 810-1038 | true | COMPLETED |
| Tabu Search | 1.000 | 0.000 | 0.000 | 1587.429 | 1561-1668 | true | COMPLETED |

## 4. Nhận xét

### 4.1 Feasibility và hard constraints

- Cả 3 solver đều đạt `feasibility = 1.000` ở cả 3 scenario.
- Không solver nào sinh hard violation trong benchmark hiện tại.
- Điều này cho thấy abstraction solver + hard constraints hiện đã đồng nhất đủ tốt để benchmark công bằng.

### 4.2 Soft score

- Ở scenario `small`, `CP-SAT` cho soft score tốt nhất.
- `Tabu Search` thấp hơn `CP-SAT`.
- `Graph Coloring` thấp nhất ở scenario `small`.
- Ở `medium` và `large`, chênh lệch soft score không còn xuất hiện trong bộ dữ liệu hiện tại.

### 4.3 Runtime và scalability

- `Graph Coloring` là solver nhanh nhất, gần như tức thời trên cả 3 scenario.
- `CP-SAT` chậm hơn rõ rệt nhưng vẫn còn trong ngưỡng chấp nhận được với scenario `large` hiện tại.
- `Tabu Search` là solver chậm nhất ở dữ liệu `large`.

### 4.4 Stability

- Cả 3 solver đều ổn định giữa 7 lần chạy trong từng scenario.
- Signature kết quả không đổi theo các tiêu chí status / scheduled / unscheduled / conflict / soft score.
- Runtime có dao động nhỏ nhưng không làm thay đổi ranking theo chất lượng nghiệm.

## 5. Quyết định chọn solver chính

**Solver được chọn cho giai đoạn E3:** `cp_sat`

### Lý do chọn

1. `CP-SAT` đạt cùng mức feasibility và hard-constraint correctness như các solver còn lại.
2. `CP-SAT` cho quality tốt nhất ở scenario `small`, là nơi bắt đầu thấy khác biệt về soft score.
3. Runtime của `CP-SAT` chậm hơn `Graph Coloring`, nhưng vẫn nằm trong mức chấp nhận được với dataset `large` hiện tại của đồ án.
4. `Tabu Search` không tạo ra lợi thế đủ rõ về quality để bù chi phí runtime.
5. Với scope đồ án hiện tại, ưu tiên hợp lý là:
   - giữ nghiệm sạch hard constraints,
   - chọn solver có tín hiệu tốt hơn về chất lượng nghiệm,
   - chấp nhận runtime cao hơn baseline heuristic nếu vẫn đủ để demo và benchmark báo cáo.

## 6. Hướng dùng cho giai đoạn E3

Từ kết quả benchmark này:

- `CP-SAT` đã được inject làm default solver cho preview scheduling chính.
- `Graph Coloring` vẫn nên giữ lại như baseline benchmark và fallback nghiên cứu.
- `Tabu Search` tiếp tục giữ trong benchmark admin API để so sánh metaheuristic.

## 7. Kết luận ngắn

Task E2 được xem là hoàn thành khi:

- có benchmark dataset lặp lại được,
- có công cụ chạy benchmark nhiều lần,
- có bảng số liệu so sánh 3 solver,
- có quyết định chọn solver chính,
- có tài liệu đủ dùng cho báo cáo và bước E3.
