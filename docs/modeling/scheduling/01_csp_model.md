# Scheduling Algorithm: Constraint Satisfaction Problem (CSP) Model

Mục tiêu của tài liệu này là định nghĩa chính thức bài toán Xếp lịch của hệ thống dưới góc độ một Bài toán Thỏa mãn Ràng buộc (CSP).

---

## 1. Biến số (Variables - $V$)
Mỗi Variable $v \in V$ đại diện cho một buổi học (ClassSession/Lesson) cần được xếp lịch.
Một Class có thể bao gồm nhiều Sessions (ví dụ: Lớp Toán 12 yêu cầu 2 buổi/tuần).
- **$V = \{v_1, v_2, ..., v_n\}$**
- Trong đó, mỗi $v_i$ mang các thuộc tính:
  - `class_id`: Định danh lớp học.
  - `teacher_id`: (Có thể đã fix sẵn hoặc cần tìm).
  - `duration`: Thời lượng của buổi học (VD: 90 phút, 120 phút).
  - `capacity`: Sĩ số dự kiến của lớp học.

## 2. Miền giá trị (Domains - $D$)
Mỗi biến $v_i$ sẽ lấy một giá trị $d \in D_i$.
Giá trị $d$ là một bộ tuple kết hợp của **Thời gian (TimeSlot)** và **Không gian (Room)**.
- **$d = (TimeSlot, Room)$**
- **$D_i = \{ (t_1, r_1), (t_1, r_2), ..., (t_k, r_m) \}$**

*(Chú ý: Nếu `teacher_id` chưa được fix cứng cho Lớp, giá trị $d$ sẽ mở rộng thành tuple 3: `(TimeSlot, Room, Teacher)`)*

## 3. Ràng buộc cứng (Hard Constraints - $C_{hard}$)
Một giải pháp (Assignment) chỉ hợp lệ nếu thỏa mãn TẤT CẢ các ràng buộc cứng:
1. **Teacher Conflict ($C_{teacher}$):** Một giáo viên không thể dạy 2 lớp cùng một lúc.
   - $\forall v_i, v_j \in V: (teacher(v_i) == teacher(v_j)) \implies (timeslot(v_i) \neq timeslot(v_j))$
2. **Room Conflict ($C_{room}$):** Một phòng học không thể diễn ra 2 lớp cùng một lúc.
   - $\forall v_i, v_j \in V: (room(v_i) == room(v_j)) \implies (timeslot(v_i) \neq timeslot(v_j))$
3. **Room Capacity ($C_{capacity}$):** Sức chứa của phòng học phải $\ge$ Sĩ số lớp.
   - $\forall v_i \in V: capacity(room(v_i)) \ge capacity(v_i)$
4. **Time Rule - Thông tư 29 ($C_{tt29}$):** Không có ca học nào kết thúc sau 22h00.
   - $\forall v_i \in V: end\_time(timeslot(v_i)) \le 22:00$

## 4. Ràng buộc mềm (Soft Constraints - $C_{soft}$) & Hàm Mục Tiêu (Objective Function)
Thay vì chỉ tìm ra một lịch hợp lệ, Backtracking Solver sẽ cố gắng tối ưu hóa điểm số dựa trên các ràng buộc mềm:
1. **Liên tiếp ($SC_{consecutive}$):** Ưu tiên xếp lịch cho giáo viên dạy các ca liên tiếp nhau (tránh khoảng trống vô ích).
2. **Tránh Gap dài ($SC_{gap}$):** Phạt điểm (Penalty) nếu giáo viên có khoảng trống $> 2$ tiếng giữa 2 ca dạy trong ngày.
3. **Phân bổ đều ($SC_{spread}$):** Các buổi học của cùng 1 Class nên cách nhau hợp lý (VD: Thứ 2 - Thứ 5, thay vì Thứ 2 - Thứ 3 liền nhau).

**Hàm mục tiêu:**  
Maximize $Score(Assignment) = \sum (w_k \times SC_k)$
- $w_k$: Trọng số (weight) của từng loại Soft Constraint.

---
## 5. Chiến lược tối ưu (Heuristics & Pruning)
- **MRV (Minimum Remaining Values):** Luôn chọn Variable có ít Domain values nhất để gán trước (Fail-First). Thường là các lớp học có yêu cầu phòng lab đặc biệt hoặc giáo viên cực kỳ bận.
- **LCV (Least Constraining Value):** Với Variable đã chọn, thử assign Domain value ít gây cản trở nhất cho các Variable khác chưa được gán.
- **Forward Checking:** Sau mỗi lần assign thành công $v_i = d$, duyệt qua tất cả các $v_{unassigned}$ liên quan, loại bỏ các giá trị ra khỏi miền $D_{unassigned}$ nếu chúng vi phạm Hard Constraints với $d$. Nếu bất kỳ $D_{unassigned}$ nào rỗng, lập tức Backtrack.
