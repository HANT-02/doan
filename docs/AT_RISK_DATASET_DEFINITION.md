# AT_RISK Dataset Definition

**Ngày chốt:** 2026-04-14

**Trạng thái:** khóa cho `Task F1`

## 1. Mục tiêu

Khóa dataset definition ban đầu cho bài toán:

- phân loại học viên `AT_RISK`
- theo hướng supervised learning
- dùng backend hiện tại để chuẩn bị cho `F2` pipeline train/evaluate

## 2. Bài toán ML được chốt

- Problem type: `classification`
- Target classes:
  - `AT_RISK`
  - `NOT_AT_RISK`
- Prediction unit: `student_class_snapshot`
- Observation window: `28 ngày` trước thời điểm snapshot
- Prediction horizon: `28 ngày` sau thời điểm snapshot

Ý nghĩa:

- Mỗi sample đại diện cho một học viên trong một lớp tại một thời điểm quan sát.
- Feature chỉ được lấy từ dữ liệu đã xảy ra trước snapshot.
- Label được xác định từ hành vi/kết quả trong 28 ngày tiếp theo.

## 3. Nguồn dữ liệu được chốt

### 3.1. Nguồn bắt buộc

1. `student`
   - entity: `Student`
   - field dùng ban đầu:
     - `id`
     - `grade_level`
     - `date_of_birth`
     - `gender`
     - `status`

2. `class/course enrollment`
   - entity: `Enrollment`
   - field dùng ban đầu:
     - `student_id`
     - `class_id`
     - `status`
     - `approved_at`
     - `created_at`

3. `attendance`
   - entity: `Attendance`
   - field dùng ban đầu:
     - `student_id`
     - `lesson_id`
     - `status`
     - `marked_at`

4. `grade`
   - nguồn grade được chốt là `AcademicRecord`
   - field dùng ban đầu:
     - `student_id`
     - `lesson_summary_id`
     - `homework_completed`
     - `homework_score`
     - `attitude_rating`
     - `participation_score`
     - `total_score`
     - `is_completed`
     - `created_at`

5. `lesson`
   - entity: `Lesson`
   - dùng để suy ra tải học và nối attendance / summary / class
   - field dùng ban đầu:
     - `id`
     - `class_id`
     - `date_start`
     - `date_end`

### 3.2. Tín hiệu vận hành được chốt cho baseline

1. `leave_request`
   - field dùng:
     - `leave_type`
     - `apply_date`
     - `late_minutes`
     - `early_minutes`
     - `status`

2. `lesson_summary`
   - field dùng:
     - `lesson_id`
     - `homework_deadline`
     - `created_at`

## 4. Ghi chú về chất lượng dữ liệu

### 4.1. Attendance status

Code hiện tại dùng `Attendance.status` kiểu `int` nhưng chưa có enum chính thức ở nghiệp vụ.

Để F2 có thể triển khai, baseline mapping cho predictive pipeline được chốt tạm như sau:

- `1 = PRESENT`
- `2 = ABSENT`
- `3 = EXCUSED`
- `4 = LATE`
- `5 = EARLY_LEAVE`

Nếu dữ liệu thực tế không khớp mapping này, pipeline phải:

- log số record unmapped,
- loại record lỗi khỏi training set,
- và ghi cảnh báo vào metadata của model.

### 4.2. Scope kiểm soát

Trong `F1`, chỉ khóa baseline tabular dataset.

Không đưa vào baseline:

- `personal_comment`
- `teacher_notes`
- `class_feedback`
- `lesson_content`
- `leave_request.reason`

Lý do:

- dữ liệu text làm tăng độ phức tạp,
- khó chuẩn hóa trong scope đồ án,
- dễ phát sinh leakage hoặc overfit khi dữ liệu còn mỏng.

## 5. Label `AT_RISK` được chốt

Một sample được gán nhãn `AT_RISK` nếu **thỏa ít nhất một** điều kiện trong `prediction horizon = 28 ngày` sau snapshot:

1. `future_attendance_rate < 0.80`
2. `future_average_total_score < 5.00`
3. `future_homework_completion_rate < 0.60`

Điều kiện đủ dữ liệu để gán nhãn:

- có ít nhất `4` attendance rows trong horizon
  hoặc
- có ít nhất `2` academic records trong horizon

Sample bị loại khỏi training set nếu:

- học viên không ở trạng thái `ACTIVE` tại snapshot,
- không có active enrollment tại snapshot,
- horizon không đủ dữ liệu theo ngưỡng tối thiểu ở trên.

## 6. Feature set ban đầu

### 6.1. Nhóm profile học viên

- `student_grade_level`
- `student_age_years`
- `student_gender`

### 6.2. Nhóm tải học / enrollment

- `active_enrollment_count_28d`
- `weekly_lesson_load_28d`

### 6.3. Nhóm attendance

- `attendance_rate_28d`
- `absence_count_28d`
- `late_or_early_incident_count_28d`
- `approved_leave_count_28d`

### 6.4. Nhóm grade / academic performance

- `average_total_score_28d`
- `minimum_total_score_28d`
- `homework_completion_rate_28d`
- `average_homework_score_28d`
- `average_participation_score_28d`
- `average_attitude_rating_28d`
- `completed_record_ratio_28d`

### 6.5. Nhóm vận hành theo thời gian

- `days_since_last_lesson`

## 7. Leakage guard được chốt

1. Không dùng bất kỳ dữ liệu nào sau snapshot để tạo feature.
2. Chỉ dùng dữ liệu tương lai để gán label.
3. Không dùng text free-form trong baseline F2.
4. Không dùng trực tiếp kết quả label horizon làm feature ở observation window.

## 8. Tài sản code đã khóa cho F1

Contract code đã được thêm tại:

- [definition.go](/Users/hant/golang/doan/internal/services/predictive/definition.go)
- [definition_test.go](/Users/hant/golang/doan/internal/services/predictive/definition_test.go)

Mục đích:

- làm chuẩn nguồn dữ liệu cho data loader ở `F2`
- làm chuẩn label/feature cho feature engineering service
- giảm tranh cãi lại scope khi sang train/test split và model training

## 9. Kết luận ngắn

`Task F1` được xem là hoàn tất khi:

- đã chốt nguồn dữ liệu bắt buộc,
- đã chốt label `AT_RISK`,
- đã chốt feature set ban đầu,
- đã có contract code để triển khai `F2`.
