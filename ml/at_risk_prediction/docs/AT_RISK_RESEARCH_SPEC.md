# Đặc tả nghiên cứu cho hệ thống dự báo AT_RISK

**Phiên bản:** 1.0  
**Ngày cập nhật:** 2026-04-25  
**Thuộc giai đoạn:** `Task K4`  
**Thuộc project:** [ml/at_risk_prediction](/Users/hant/golang/doan/ml/at_risk_prediction)

## 1. Mục tiêu nghiên cứu

Mục tiêu của nhánh nghiên cứu này là xây dựng một hệ thống dự báo sớm học viên có nguy cơ học tập kém (`AT_RISK`) trong môi trường trung tâm giáo dục EduCenter.

Hệ thống cần đạt được các yêu cầu sau:

- tận dụng dữ liệu học vụ và vận hành đã có trong hệ thống,
- dự báo ở mức sớm, không chờ đến khi học viên rớt hẳn,
- có thể giải thích được lý do cảnh báo,
- có thể tái lập thực nghiệm để đưa vào báo cáo đồ án,
- đủ nhẹ để huấn luyện trên máy cá nhân,
- tích hợp được với backend Go hiện tại.

## 2. Bài toán đầu vào -> đầu ra

### 2.1. Đầu vào

Đầu vào của hệ thống là dữ liệu bảng (tabular data) trích xuất từ:

- hồ sơ học viên (`students`),
- ghi danh học viên vào lớp (`enrollments`),
- lịch/buổi học (`lessons`),
- điểm danh (`attendance`),
- tổng kết buổi học (`lesson_summaries`),
- kết quả học tập (`academic_records`),
- đơn xin phép (`leave_requests`),
- hoặc dữ liệu đã được export thành CSV.

### 2.2. Đầu ra

Đầu ra của hệ thống gồm hai lớp:

1. **Đầu ra mô hình học máy**
   - `risk_label ∈ {AT_RISK, NOT_AT_RISK}`
   - `risk_score ∈ [0,1]`
   - `model_version`
   - `top_features`
   - `primary_reason`

2. **Đầu ra phục vụ hệ thống**
   - `model_metadata.json`
   - `latest_predictions.json`
   - `metrics.json`
   - `classification_report.md`
   - các hình như confusion matrix, feature importance

## 3. Đối tượng dự báo

Đối tượng dự báo không phải là toàn bộ học viên nói chung, mà là một đơn vị quan sát dạng:

```text
student_class_snapshot
```

Tức là:

- một học viên,
- trong một lớp cụ thể,
- tại một thời điểm snapshot cụ thể.

Việc chọn đơn vị này có ý nghĩa quan trọng:

- một học viên có thể học nhiều lớp,
- nguy cơ học tập có thể khác nhau theo từng lớp,
- dữ liệu attendance/academic record đều gắn với lớp và buổi học.

## 4. Định nghĩa nhãn `AT_RISK`

Một mẫu được gán nhãn `AT_RISK` nếu trong cửa sổ dự báo tương lai học viên có ít nhất một dấu hiệu học tập đáng lo ngại.

### 4.1. Cửa sổ quan sát và cửa sổ dự báo

- **Observation window:** `28 ngày` trước thời điểm snapshot
- **Prediction horizon:** `28 ngày` sau thời điểm snapshot

Diễn giải:

- Feature chỉ được tạo từ dữ liệu trong 28 ngày trước snapshot.
- Label chỉ được xác định từ dữ liệu trong 28 ngày sau snapshot.

### 4.2. Luật gán nhãn

Một snapshot được gán `AT_RISK` nếu thỏa ít nhất một trong các điều kiện:

1. `future_attendance_rate < 0.80`
2. `future_average_total_score < 5.00`
3. `future_homework_completion_rate < 0.60`

Ngược lại, snapshot được gán `NOT_AT_RISK`.

### 4.3. Điều kiện đủ dữ liệu để gán nhãn

Một snapshot chỉ được đưa vào tập huấn luyện nếu có đủ dữ liệu trong cửa sổ tương lai:

- ít nhất `4` bản ghi attendance
  hoặc
- ít nhất `2` bản ghi academic record

Các snapshot không đủ dữ liệu sẽ bị loại để tránh tạo nhãn nhiễu.

## 5. Mô hình toán học của bài toán

Ta ký hiệu:

- `x_i ∈ R^d`: vector đặc trưng của mẫu thứ `i`
- `y_i ∈ {0,1}`: nhãn nhị phân của mẫu thứ `i`
  - `y_i = 1` tương ứng `AT_RISK`
  - `y_i = 0` tương ứng `NOT_AT_RISK`

Với mỗi mẫu `i`, bài toán cần tìm một hàm dự báo:

```math
f(x_i) -> p_i
```

Trong đó:

- `p_i` là xác suất học viên thuộc lớp `AT_RISK`
- `0 <= p_i <= 1`

Từ xác suất `p_i`, hệ thống dùng một ngưỡng quyết định `τ` để phân lớp:

```math
\hat{y_i} =
\begin{cases}
1, & \text{nếu } p_i \ge \tau \\
0, & \text{nếu } p_i < \tau
\end{cases}
```

Trong baseline ban đầu của đồ án:

- ngưỡng mặc định là `τ = 0.5`

Trong các vòng tối ưu sau này, ngưỡng có thể được tinh chỉnh để ưu tiên `Recall`.

## 6. Danh sách feature và ý nghĩa nghiệp vụ

### 6.1. Nhóm feature chuyên cần

1. `attendance_rate_28d`
   - tỷ lệ có mặt trong 28 ngày quan sát
   - phản ánh mức độ tham dự học tập

2. `absence_count_28d`
   - số buổi vắng trong 28 ngày quan sát
   - phản ánh dấu hiệu bỏ học hoặc mất kết nối

3. `approved_leave_count_28d`
   - số đơn xin phép đã được duyệt
   - bổ sung ngữ cảnh để phân biệt nghỉ có phép và nghỉ bất thường

### 6.2. Nhóm feature kết quả học tập

4. `average_total_score_28d`
   - điểm tổng kết trung bình trong cửa sổ quan sát
   - phản ánh chất lượng học tập gần đây

5. `homework_completion_rate_28d`
   - tỷ lệ hoàn thành bài tập
   - phản ánh mức độ theo bài và tự học

### 6.3. Nhóm feature tải học

6. `active_enrollment_count_28d`
   - số lớp học viên đang theo học tại thời điểm snapshot
   - phản ánh độ phân tán nguồn lực học tập

7. `weekly_lesson_load_28d`
   - số buổi học trung bình mỗi tuần trong 28 ngày quan sát
   - phản ánh cường độ học

### 6.4. Nhóm feature vận hành theo thời gian

8. `days_since_last_lesson`
   - số ngày từ buổi học gần nhất tới snapshot
   - phản ánh mức độ gián đoạn học tập

## 7. Mô hình baseline rule-based

Rule-based baseline không phải là mô hình học máy học trọng số từ dữ liệu, mà là tập luật nghiệp vụ do con người đặt ra.

### 7.1. Ý tưởng

Nếu một học viên có nhiều dấu hiệu xấu cùng lúc, hệ thống coi học viên đó là `AT_RISK`.

Ví dụ các luật cơ bản:

- chuyên cần thấp,
- điểm trung bình thấp,
- bài tập không hoàn thành đủ.

### 7.2. Dạng mô hình

Ta có thể mô tả một hàm điểm rủi ro:

```math
R(x) = w_1 I(attendance\_rate < 0.8)
     + w_2 I(avg\_score < 5.0)
     + w_3 I(homework\_completion < 0.6)
```

Trong đó:

- `I(.)` là hàm chỉ báo nhận giá trị `1` nếu điều kiện đúng, `0` nếu sai
- `w_1, w_2, w_3` là trọng số nghiệp vụ

Sau đó:

```math
\hat{y} =
\begin{cases}
AT\_RISK, & \text{nếu } R(x) \ge \theta \\
NOT\_AT\_RISK, & \text{ngược lại}
\end{cases}
```

### 7.3. Vai trò trong nghiên cứu

Rule-based baseline có vai trò:

- làm chuẩn so sánh tối thiểu,
- thể hiện tri thức nghiệp vụ,
- giúp kiểm tra xem mô hình ML có thực sự tốt hơn luật đơn giản hay không.

## 8. Mô hình Logistic Regression

### 8.1. Lý do chọn

Logistic Regression là mô hình phù hợp cho bài toán classification nhị phân khi:

- dữ liệu là dạng bảng,
- số lượng feature vừa phải,
- cần tính giải thích cao,
- cần huấn luyện nhẹ trên CPU.

### 8.2. Công thức mô hình

Với vector đặc trưng `x`, mô hình tính:

```math
z = w^T x + b
```

Sau đó đưa qua hàm sigmoid:

```math
\sigma(z) = \frac{1}{1 + e^{-z}}
```

Xác suất học viên thuộc lớp `AT_RISK` là:

```math
P(y=1|x) = \sigma(w^T x + b)
```

### 8.3. Hàm mất mát

Logistic Regression được huấn luyện bằng cách tối thiểu hóa hàm mất mát cross-entropy:

```math
L = - \frac{1}{N} \sum_{i=1}^{N}
\left[
y_i \log(\hat{y_i}) + (1-y_i)\log(1-\hat{y_i})
\right]
```

Trong đó:

- `y_i` là nhãn thật,
- `\hat{y_i}` là xác suất dự báo.

Nếu dùng regularization `L2`, hàm mất mát có thể mở rộng:

```math
L_{reg} = L + \lambda ||w||_2^2
```

### 8.4. Ý nghĩa nghiên cứu

Ưu điểm của Logistic Regression:

- dễ giải thích trọng số feature,
- dễ huấn luyện,
- phù hợp làm baseline ML chính thức,
- phù hợp để trình bày trong báo cáo đồ án.

Hạn chế:

- giả định quan hệ gần tuyến tính giữa feature và log-odds,
- khó nắm bắt tương tác phi tuyến phức tạp.

## 9. Mô hình Random Forest

### 9.1. Lý do chọn

Random Forest là mô hình tổ hợp từ nhiều cây quyết định, phù hợp với dữ liệu bảng và có khả năng học quan hệ phi tuyến tốt hơn Logistic Regression.

### 9.2. Nguyên lý hoạt động

Random Forest gồm nhiều cây quyết định:

1. Mỗi cây được huấn luyện trên một tập bootstrap sampling của dữ liệu.
2. Tại mỗi node, mô hình chỉ xét một tập con ngẫu nhiên của feature.
3. Mỗi cây đưa ra một dự báo riêng.
4. Kết quả cuối cùng là tổng hợp của toàn bộ cây.

### 9.3. Mô tả toán học

Giả sử rừng có `T` cây:

```math
\hat{p}(x) = \frac{1}{T} \sum_{t=1}^{T} p_t(x)
```

Trong đó:

- `p_t(x)` là xác suất dự báo của cây thứ `t`
- `\hat{p}(x)` là xác suất trung bình của toàn bộ rừng

Quy tắc phân lớp:

```math
\hat{y} =
\begin{cases}
1, & \text{nếu } \hat{p}(x) \ge \tau \\
0, & \text{ngược lại}
\end{cases}
```

### 9.4. Vai trò trong nghiên cứu

Ưu điểm:

- học được quan hệ phi tuyến,
- ít nhạy hơn với scale của feature,
- thường cho kết quả tốt trên dữ liệu tabular.

Hạn chế:

- khó giải thích hơn Logistic Regression,
- nặng hơn về tài nguyên,
- kích thước model artifact thường lớn hơn.

## 10. Metric đánh giá

Với bài toán cảnh báo sớm, chỉ dùng Accuracy là chưa đủ.  
Hệ thống cần đánh giá theo nhiều chiều.

### 10.1. Accuracy

```math
Accuracy = \frac{TP + TN}{TP + TN + FP + FN}
```

Ý nghĩa:

- tỷ lệ dự báo đúng trên toàn bộ mẫu.

### 10.2. Precision

```math
Precision = \frac{TP}{TP + FP}
```

Ý nghĩa:

- trong số các học viên bị cảnh báo là `AT_RISK`, có bao nhiêu học viên thực sự có nguy cơ.

### 10.3. Recall

```math
Recall = \frac{TP}{TP + FN}
```

Ý nghĩa:

- trong số các học viên thực sự có nguy cơ, mô hình phát hiện được bao nhiêu.

### 10.4. F1-score

```math
F1 = 2 \cdot \frac{Precision \cdot Recall}{Precision + Recall}
```

Ý nghĩa:

- cân bằng giữa Precision và Recall.

### 10.5. Confusion Matrix

Confusion matrix gồm 4 thành phần:

- `TP`: dự báo nguy cơ đúng
- `TN`: dự báo an toàn đúng
- `FP`: cảnh báo nhầm
- `FN`: bỏ sót học viên nguy cơ

### 10.6. Vì sao ưu tiên Recall và F1-score

Trong bài toán cảnh báo sớm, bỏ sót học viên thực sự có nguy cơ (`FN`) nguy hiểm hơn việc cảnh báo dư một số học viên (`FP`).

Vì vậy:

- `Recall` là metric rất quan trọng,
- `F1-score` giúp cân bằng giữa khả năng phát hiện và độ chính xác cảnh báo.

Trong báo cáo, việc chọn mô hình chính nên ưu tiên:

1. `Recall`
2. `F1-score`
3. khả năng giải thích
4. độ nhẹ khi triển khai trên máy cá nhân

## 11. Quy trình thực nghiệm

Quy trình thực nghiệm chuẩn của project Python:

1. Kết nối PostgreSQL hoặc nạp CSV.
2. Xây dựng dataset `student_class_snapshot`.
3. Làm sạch và chuẩn hóa feature.
4. Chia train/test.
5. Huấn luyện:
   - rule-based baseline,
   - Logistic Regression,
   - Random Forest.
6. Tính metric:
   - Accuracy,
   - Precision,
   - Recall,
   - F1-score,
   - confusion matrix.
7. So sánh mô hình.
8. Chọn mô hình chính.
9. Sinh:
   - model artifact,
   - metrics,
   - report,
   - prediction artifact cho Go backend.

## 12. Tái lập thực nghiệm

Để đảm bảo khả năng tái lập:

- dùng cùng dataset export hoặc cùng snapshot DB,
- cố định random seed,
- ghi rõ phiên bản thư viện,
- lưu model metadata,
- lưu train/test split hoặc CSV đã chuẩn hóa.

Các artifact cần giữ:

- `*_full.csv`
- `*_train.csv`
- `*_test.csv`
- `metrics.json`
- `classification_report.md`
- `model_metadata.json`

## 13. Hạn chế hiện tại

1. Dữ liệu hiện tại vẫn phụ thuộc mức độ đầy đủ của:
   - `lessons`
   - `attendance`
   - `academic_records`

2. Label `AT_RISK` hiện được xây dựng theo heuristic nghiệp vụ, chưa phải nhãn do cố vấn học tập xác nhận thủ công.

3. Feature text như:
   - `teacher_notes`
   - `class_feedback`
   - `personal_comment`
   hiện chưa được khai thác để tránh làm tăng độ phức tạp.

4. Mô hình hiện tập trung vào dữ liệu tabular, chưa xét chuỗi thời gian nâng cao hay deep learning.

## 14. Hướng phát triển

1. Bổ sung calibration cho xác suất dự báo.
2. Tối ưu ngưỡng quyết định `τ` theo mục tiêu Recall.
3. Thêm `feature importance` và `SHAP` nếu thời gian cho phép.
4. Đánh giá chéo theo nhiều khóa dữ liệu khác nhau.
5. Mở rộng sang dự báo đa mức nguy cơ thay vì chỉ nhị phân.

## 15. Kết luận

Tài liệu này khóa hướng nghiên cứu chính thức của nhánh predictive analytics:

- dùng Python làm nền tảng ML,
- dùng dữ liệu thực từ hệ thống EduCenter,
- xây dựng bài toán classification nhị phân `AT_RISK`,
- đánh giá bằng các metric phù hợp cho bài toán cảnh báo sớm,
- sinh artifact để Go backend và frontend tích hợp ổn định.

Tài liệu này là cơ sở cho:

- `K5` huấn luyện và đánh giá mô hình,
- `K6` chọn mô hình chính,
- `K8` cập nhật báo cáo đồ án chính thức.
