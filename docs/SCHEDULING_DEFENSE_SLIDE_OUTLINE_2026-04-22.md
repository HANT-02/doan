# DÀN Ý SLIDE BẢO VỆ CHO PHẦN SCHEDULING

## Slide 1. Scheduling là điểm nhấn kỹ thuật

- Bài toán không chỉ là CRUD mà là tối ưu hóa xếp lịch.
- Dữ liệu đầu vào gồm lớp, giáo viên, phòng, `Shift`, `ClassSchedule`.
- Đầu ra là preview và `lesson` thật sau commit.

## Slide 2. Kiến trúc solver

- Abstraction `SchedulingSolver`.
- Ba solver:
  - `Graph Coloring + Heuristic`
  - `CP-SAT`
  - `Tabu Search`
- Lợi ích: benchmark công bằng, dễ thay solver, dễ mở rộng.

## Slide 3. Hard constraints và soft constraints

- Không trùng giáo viên, phòng, lớp.
- Chỉ dùng `Shift` active.
- Chỉ dùng slot từ `ClassSchedule`.
- Soft score dùng để phân biệt chất lượng nghiệm.

## Slide 4. Thiết kế benchmark

- 3 scenario: `small`, `medium`, `large`
- 7 lần chạy cho mỗi scenario
- Metric:
  - feasibility
  - hard violation
  - soft score
  - runtime
  - stability

## Slide 5. Kết quả benchmark

- Cả 3 solver đều sạch hard constraints.
- `Graph Coloring` nhanh nhất.
- `CP-SAT` cho tín hiệu chất lượng nghiệm tốt nhất ở scenario `small`.
- `Tabu Search` không thắng được về chất lượng so với runtime bỏ ra.

## Slide 6. Quyết định chọn solver chính

- Chọn `CP-SAT`.
- Lý do:
  - đúng hard constraints;
  - chất lượng nghiệm tốt hơn ở case phân biệt;
  - runtime vẫn chấp nhận được trong phạm vi đồ án.

## Slide 7. Luồng production-like

- Chạy preview
- Xem assignments/conflicts
- Commit preview
- Sinh `lesson`
- Lesson tiếp tục phục vụ attendance, summary, academic record

## Slide 8. Kết luận và hướng mở rộng

- Scheduling đã có benchmark thật, không chỉ có mô phỏng.
- Có thể mở rộng benchmark với dữ liệu DB thật hơn.
- Có thể mở rộng soft score và workload balancing trong giai đoạn sau.
