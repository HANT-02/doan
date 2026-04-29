# Dự thảo mục báo cáo: Dự báo học viên có nguy cơ học tập kém

**Ngày soạn:** 2026-04-28  
**Nguồn số liệu:** dữ liệu PostgreSQL demo của hệ thống EduCenter  
**Mục đích:** dùng để cập nhật phần báo cáo học thuật về bài toán dự báo học viên có nguy cơ học tập kém trước khi sửa trực tiếp Google Doc

> Ghi chú về thuật ngữ:
> - Trong bản thảo này, cụm “nhãn cảnh báo nguy cơ học tập kém” được dùng thay cho việc lặp lại quá nhiều ký hiệu `AT_RISK`.
> - Nếu cần giữ các ký hiệu tiếng Anh trong bản báo cáo chính, nên bổ sung vào bảng chữ viết tắt của tài liệu các mục sau: `ML`, `Python`, `scikit-learn`, `Random Forest`, `Logistic Regression`.

---

## 4.8 Dự báo học viên có nguy cơ học tập kém bằng học máy

### 4.8.1 Bối cảnh nghiên cứu và mục tiêu bài toán

Trong hệ thống EduCenter, ngoài bài toán xếp lịch thông minh, một hướng mở rộng có giá trị thực tiễn cao là dự báo sớm học viên có nguy cơ học tập kém. Đây là một bài toán có ý nghĩa lớn trong quản trị học vụ vì nếu phát hiện được sớm, trung tâm có thể:

- chủ động hỗ trợ học viên bằng phụ đạo hoặc tư vấn học tập;
- cảnh báo giáo viên về những trường hợp có dấu hiệu suy giảm chuyên cần;
- hỗ trợ phụ huynh theo dõi tiến độ học tập của học viên;
- hình thành cơ sở dữ liệu cho các quyết định điều hành theo hướng định lượng.

Ngay từ giai đoạn thiết kế phần nghiên cứu này, Python đã được lựa chọn là ngôn ngữ chính cho nhánh học máy. Lý do lựa chọn là vì Python có hệ sinh thái phù hợp với các bài toán phân tích dữ liệu giáo dục và phân lớp nhị phân, đặc biệt với các thư viện:

- `pandas`: xử lý dữ liệu bảng;
- `numpy`: tính toán số học;
- `scikit-learn`: huấn luyện và đánh giá mô hình học máy;
- `sqlalchemy` và `psycopg`: kết nối và trích xuất dữ liệu từ PostgreSQL;
- `joblib`: lưu mô hình đã huấn luyện;
- `matplotlib` và `seaborn`: trực quan hóa kết quả thực nghiệm.

Việc lựa chọn Python ngay từ đầu giúp bài toán được phát triển theo đúng tinh thần nghiên cứu: có dữ liệu, có quy trình tiền xử lý, có mô hình, có thực nghiệm và có đánh giá định lượng.

**[Chèn Hình 4.4.1: Vai trò của nhánh dự báo học tập kém trong tổng thể hệ thống EduCenter]**

### 4.8.2 Đầu vào, đầu ra và đơn vị dự báo

Đầu vào của hệ thống là dữ liệu bảng được trích xuất từ các bảng nghiệp vụ đã có của EduCenter:

- `students`
- `enrollments`
- `lessons`
- `attendance`
- `lesson_summaries`
- `academic_records`
- `leave_requests`

Do một học viên có thể đồng thời theo học nhiều lớp khác nhau, đơn vị dự báo không phải là “một học viên” theo nghĩa chung, mà là:

```text
student_class_snapshot
```

Nói cách khác, mỗi mẫu trong tập dữ liệu biểu diễn:

- một học viên;
- trong một lớp cụ thể;
- tại một thời điểm snapshot cụ thể;
- với các đặc trưng được tổng hợp từ lịch sử học tập gần nhất.

Đầu ra của mô hình gồm:

- nhãn phân lớp: học viên có nguy cơ hay không có nguy cơ học tập kém;
- xác suất rủi ro `p ∈ [0,1]`;
- mức cảnh báo;
- lý do cảnh báo chính;
- các đặc trưng có ảnh hưởng lớn đến quyết định dự báo.

Ở tầng hệ thống, đầu ra này được lưu thành các tệp kết quả mô hình để backend sử dụng lại cho giao diện quản trị và giao diện học viên.

**[Chèn Hình 4.4.2: Sơ đồ đầu vào - đầu ra của bài toán dự báo học viên có nguy cơ học tập kém]**

### 4.8.3 Quy trình tiền xử lý dữ liệu

Vì dữ liệu giáo dục có tính thời gian rõ rệt, phần tiền xử lý đóng vai trò đặc biệt quan trọng. Quy trình tiền xử lý trong nghiên cứu này gồm các bước:

1. **Lọc dữ liệu hợp lệ**
   - chỉ giữ học viên có trạng thái hoạt động;
   - chỉ giữ các ghi danh hợp lệ trong lớp;
   - chỉ giữ các buổi học đã tồn tại trong hệ thống.

2. **Ghép dữ liệu liên bảng**
   - ghép học viên với lớp học;
   - ghép lớp học với các buổi học;
   - ghép buổi học với dữ liệu điểm danh;
   - ghép tổng kết buổi học với kết quả học tập;
   - ghép đơn xin phép để bổ sung ngữ cảnh chuyên cần.

3. **Chuẩn hóa thời gian**
   - chuyển các trường ngày giờ sang cùng định dạng thời gian;
   - sắp xếp các buổi học theo thứ tự thời gian;
   - xác định các mốc thời gian phục vụ cửa sổ quan sát và cửa sổ dự báo.

4. **Tạo mẫu snapshot**
   - tại mỗi buổi học đủ điều kiện, tạo một mẫu `student_class_snapshot`;
   - gom các bản ghi lịch sử gần nhất của học viên trong lớp tương ứng.

5. **Loại bỏ mẫu không đủ dữ liệu**
   - nếu không có đủ dữ liệu tương lai để gán nhãn thì mẫu bị loại bỏ;
   - mục tiêu là giảm nhiễu nhãn và giữ tính tin cậy của thực nghiệm.

Có thể biểu diễn quy trình tiền xử lý tổng quát như sau:

```text
Dữ liệu thô từ CSDL
-> làm sạch và chuẩn hóa kiểu dữ liệu
-> ghép liên bảng
-> tạo snapshot theo từng học viên trong từng lớp
-> tính đặc trưng
-> gán nhãn
-> tạo tập dữ liệu huấn luyện
```

**[Chèn Hình 4.4.3: Quy trình tiền xử lý dữ liệu cho bài toán dự báo học viên có nguy cơ học tập kém]**

### 4.8.4 Cửa sổ quan sát, cửa sổ dự báo và nguyên tắc gán nhãn

Để bảo đảm tính đúng đắn về mặt phương pháp, nghiên cứu này áp dụng hai cửa sổ thời gian:

- **Cửa sổ quan sát:** `28 ngày` trước thời điểm snapshot;
- **Cửa sổ dự báo:** `28 ngày` sau thời điểm snapshot.

Ý nghĩa của hai cửa sổ:

- mọi đặc trưng đầu vào phải được tính từ dữ liệu nằm trong cửa sổ quan sát;
- nhãn đầu ra chỉ được xác định từ dữ liệu tương lai trong cửa sổ dự báo.

Nguyên tắc này giúp tránh hiện tượng “rò rỉ dữ liệu tương lai”, tức mô hình vô tình nhìn thấy thông tin của tương lai ngay trong lúc huấn luyện.

Một mẫu được gán nhãn “có nguy cơ học tập kém” nếu trong cửa sổ dự báo thỏa ít nhất một điều kiện:

```text
Tỷ lệ chuyên cần tương lai < 0.80
hoặc
Điểm tổng hợp trung bình tương lai < 5.00
hoặc
Tỷ lệ hoàn thành bài tập tương lai < 0.60
```

Ngược lại, mẫu được gán nhãn “không có nguy cơ học tập kém”.

Ngoài ra, để một mẫu được đưa vào tập dữ liệu huấn luyện, hệ thống yêu cầu:

- có ít nhất `4` bản ghi điểm danh trong cửa sổ dự báo;
  hoặc
- có ít nhất `2` bản ghi kết quả học tập trong cửa sổ dự báo.

Đây là điều kiện tối thiểu để nhãn mang ý nghĩa thống kê chấp nhận được.

**[Chèn Hình 4.4.4: Minh họa cửa sổ quan sát 28 ngày và cửa sổ dự báo 28 ngày]**

### 4.8.5 Tập đặc trưng và các phép tính được sử dụng

Đặc trưng được xây dựng theo hướng gọn, dễ diễn giải, nhưng vẫn đủ phản ánh hành vi học tập. Bảng 4.7 trình bày các đặc trưng chính.

**Bảng 4.7. Tập đặc trưng dùng cho mô hình dự báo**

| Nhóm | Tên đặc trưng | Giải thích |
| --- | --- | --- |
| Chuyên cần | `attendance_rate_28d` | Tỷ lệ có mặt trong 28 ngày gần nhất |
| Chuyên cần | `absence_count_28d` | Số buổi vắng trong 28 ngày gần nhất |
| Học tập | `average_total_score_28d` | Điểm tổng hợp trung bình của 28 ngày quan sát |
| Học tập | `homework_completion_rate_28d` | Tỷ lệ hoàn thành bài tập |
| Tải học | `active_enrollment_count_28d` | Số lớp học viên đang theo học |
| Tải học | `weekly_lesson_load_28d` | Số buổi học trung bình mỗi tuần |
| Vận hành | `approved_leave_count_28d` | Số đơn xin phép đã được duyệt |
| Thời gian | `days_since_last_lesson` | Số ngày kể từ buổi học gần nhất |

Các phép tính đặc trưng chính được dùng trong nghiên cứu gồm:

1. **Tỷ lệ chuyên cần**

```text
attendance_rate_28d = số buổi có mặt / tổng số buổi trong cửa sổ quan sát
```

2. **Số buổi vắng**

```text
absence_count_28d = tổng số buổi có trạng thái vắng
```

3. **Điểm trung bình**

```text
average_total_score_28d = (tổng total_score của các bản ghi hợp lệ) / số bản ghi hợp lệ
```

4. **Tỷ lệ hoàn thành bài tập**

```text
homework_completion_rate_28d = số bản ghi homework_completed = TRUE / tổng số bản ghi học tập
```

5. **Mật độ buổi học theo tuần**

```text
weekly_lesson_load_28d = tổng số buổi học trong 28 ngày / (28 / 7)
```

6. **Số ngày từ buổi học gần nhất**

```text
days_since_last_lesson = snapshot_at - thời điểm buổi học gần nhất
```

Những phép tính này đều mang tính diễn giải cao, do đó dễ được giải thích trong bối cảnh quản trị giáo dục.

**[Chèn Hình 4.4.5: Bảng hoặc sơ đồ minh họa cách tính các đặc trưng chính]**

### 4.8.6 Mô hình toán học của bài toán

Xét một mẫu dữ liệu thứ `i` với:

- `x_i ∈ R^d`: vector đặc trưng;
- `y_i ∈ {0,1}`: nhãn phân lớp nhị phân.

Trong đó:

- `y_i = 1` nếu học viên thuộc nhóm có nguy cơ học tập kém;
- `y_i = 0` nếu học viên không thuộc nhóm này.

Mục tiêu của bài toán là học một hàm:

```text
f(x_i) -> p_i
```

trong đó `p_i` là xác suất mẫu thứ `i` thuộc lớp dương. Sau đó dùng ngưỡng quyết định `τ = 0.5`:

```text
y_hat_i = 1 nếu p_i >= 0.5
y_hat_i = 0 nếu p_i < 0.5
```

Để đánh giá mô hình, nghiên cứu sử dụng các đại lượng cơ bản của ma trận nhầm lẫn:

- `TP`: dự báo đúng mẫu có nguy cơ;
- `TN`: dự báo đúng mẫu không có nguy cơ;
- `FP`: dự báo nhầm mẫu không có nguy cơ thành có nguy cơ;
- `FN`: bỏ sót mẫu có nguy cơ.

Từ đó tính các chỉ số:

1. **Accuracy**

```text
Accuracy = (TP + TN) / (TP + TN + FP + FN)
```

2. **Precision**

```text
Precision = TP / (TP + FP)
```

3. **Recall**

```text
Recall = TP / (TP + FN)
```

4. **F1-score**

```text
F1 = 2 * Precision * Recall / (Precision + Recall)
```

Trong bối cảnh quản trị học tập, `Recall` và `F1-score` được ưu tiên hơn `Accuracy` vì việc bỏ sót học viên có nguy cơ là một lỗi nghiêm trọng hơn việc cảnh báo dư.

**[Chèn Hình 4.4.6: Sơ đồ minh họa bài toán phân lớp nhị phân và các chỉ số đánh giá]**

### 4.8.7 Các mô hình được thử nghiệm

Để cân bằng giữa độ chính xác, tính giải thích và tài nguyên phần cứng, nghiên cứu thử nghiệm ba mô hình:

1. mô hình luật nghiệp vụ cơ sở;
2. Logistic Regression;
3. Random Forest.

#### a. Mô hình luật nghiệp vụ cơ sở

Mô hình này không học tham số từ dữ liệu mà dựa trên các ngưỡng nghiệp vụ đã biết. Một hàm điểm rủi ro đơn giản có thể viết:

```text
R(x) = w1 * I(attendance_rate_28d < 0.8)
     + w2 * I(average_total_score_28d < 5.0)
     + w3 * I(homework_completion_rate_28d < 0.6)
```

Trong đó:

- `I(.)` là hàm chỉ báo;
- `w1, w2, w3` là trọng số nghiệp vụ.

Mô hình này được dùng như một baseline để so sánh với mô hình học máy.

#### b. Logistic Regression

Với Logistic Regression, xác suất của lớp dương được tính theo:

```text
z = w^T x + b
```

và:

```text
sigma(z) = 1 / (1 + e^(-z))
```

do đó:

```text
p(y = 1 | x) = sigma(w^T x + b)
```

Hàm mất mát được sử dụng là binary cross-entropy:

```text
L = - [ y log(p) + (1 - y) log(1 - p) ]
```

Logistic Regression có ưu điểm:

- nhẹ;
- phù hợp với dữ liệu bảng;
- dễ giải thích;
- thuận lợi khi trình bày trong báo cáo học thuật.

#### c. Random Forest

Random Forest là mô hình tổ hợp nhiều cây quyết định. Giả sử có `T` cây trong rừng, xác suất cuối cùng có thể biểu diễn:

```text
p(y = 1 | x) = (1 / T) * sum_{t=1..T} p_t(y = 1 | x)
```

Ưu điểm chính:

- biểu diễn tốt hơn quan hệ phi tuyến;
- ổn định hơn cây quyết định đơn;
- có thể trích xuất độ quan trọng của đặc trưng.

Tuy nhiên, mô hình này thường khó giải thích hơn Logistic Regression.

**[Chèn Hình 4.4.7: Bảng tóm tắt nguyên lý hoạt động của ba mô hình]**

### 4.8.8 Công nghệ và thư viện sử dụng

Phần nghiên cứu học máy sử dụng các công nghệ sau:

**Bảng 4.8. Công nghệ và thư viện dùng cho bài toán dự báo**

| Công nghệ / thư viện | Vai trò |
| --- | --- |
| Python | Ngôn ngữ chính cho nhánh học máy |
| pandas | Xử lý dữ liệu bảng |
| numpy | Tính toán số học |
| scikit-learn | Huấn luyện và đánh giá mô hình |
| sqlalchemy | Tầng kết nối dữ liệu từ PostgreSQL |
| psycopg | Driver PostgreSQL cho Python |
| joblib | Lưu và nạp mô hình đã huấn luyện |
| matplotlib | Vẽ biểu đồ, confusion matrix |
| seaborn | Trực quan hóa thống kê |

Việc sử dụng bộ thư viện này làm cho quy trình nghiên cứu có tính chuẩn hóa cao, dễ tái lập và phù hợp với cách tiếp cận khoa học dữ liệu hiện đại.

**[Chèn Hình 4.4.8: Sơ đồ công nghệ sử dụng cho bài toán học máy]**

### 4.8.9 Thiết lập thực nghiệm

Thực nghiệm hiện tại được chạy trực tiếp trên dữ liệu PostgreSQL demo đã seed vào hệ thống. Đây là điểm quan trọng vì dữ liệu không còn là ví dụ tách rời bên ngoài mà đã gắn với chính các bảng nghiệp vụ của EduCenter.

Thông số thực nghiệm:

- nguồn dữ liệu: `database`
- số mẫu: `36`
- số học viên: `6`
- số lớp: `1`
- phân phối nhãn:
  - nhóm có nguy cơ: `21`
  - nhóm không có nguy cơ: `15`
- kích thước tập huấn luyện: `28`
- kích thước tập kiểm thử: `8`
- tỉ lệ chia train/test: `0.8 / 0.2`
- random seed: `42`

Các bước thực nghiệm:

1. trích xuất dữ liệu từ DB;
2. tiền xử lý và tạo tập đặc trưng;
3. chia train/test;
4. huấn luyện ba mô hình;
5. đánh giá trên cùng tập kiểm thử;
6. chọn mô hình chính theo tiêu chí đã xác định.

**[Chèn Hình 4.4.9: Tổng quan bộ dữ liệu thực nghiệm và cách chia train/test]**

### 4.8.10 Kết quả thực nghiệm

Kết quả đánh giá ba mô hình được trình bày trong Bảng 4.9.

**Bảng 4.9. Kết quả đánh giá ba mô hình trên dữ liệu DB demo**

| Mô hình | Accuracy | Precision | Recall | F1-score | TN | FP | FN | TP | Nhận xét |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| Luật nghiệp vụ cơ sở | 0.375 | 0.000 | 0.000 | 0.000 | 3 | 0 | 5 | 0 | Không phát hiện được mẫu có nguy cơ trong tập kiểm thử |
| Logistic Regression | 0.750 | 1.000 | 0.600 | 0.750 | 3 | 0 | 2 | 3 | Cân bằng tốt giữa độ chính xác và khả năng giải thích |
| Random Forest | 0.750 | 1.000 | 0.600 | 0.750 | 3 | 0 | 2 | 3 | Hiệu năng tương đương Logistic Regression |

Từ bảng trên có thể rút ra các kết luận sau:

1. Mô hình luật nghiệp vụ cơ sở không đủ mạnh để phát hiện đúng các trường hợp có nguy cơ trong tập kiểm thử.
2. Hai mô hình học máy vượt trội rõ rệt so với baseline.
3. Cả Logistic Regression và Random Forest đều đạt `Accuracy = 0.75`, `Precision = 1.0`, `Recall = 0.6`, `F1-score = 0.75`.
4. Dù chưa đạt `Recall` tuyệt đối, hai mô hình học máy đã phát hiện được phần lớn các mẫu có nguy cơ và không tạo ra cảnh báo sai (`FP = 0`).

**[Chèn Hình 4.4.10: Confusion matrix của mô hình Logistic Regression]**

**[Chèn Hình 4.4.11: Biểu đồ mức độ ảnh hưởng của các đặc trưng]**

Ngoài metric phân lớp, điểm rủi ro trung bình của từng mô hình cho thấy cách lượng hóa rủi ro của các mô hình là khác nhau:

- mô hình luật nghiệp vụ cơ sở: `0.05`
- Logistic Regression: `0.416093`
- Random Forest: `0.454597`

Điều này có ý nghĩa quan trọng khi xây dựng dashboard cảnh báo, vì không chỉ cần biết một học viên có thuộc nhóm nguy cơ hay không, mà còn cần biết mức độ nguy cơ tương đối của học viên đó.

### 4.8.11 Lựa chọn mô hình chính

Trong hai mô hình học máy, Logistic Regression và Random Forest đang cho kết quả tương đương trên tập kiểm thử. Vì vậy nghiên cứu áp dụng thứ tự ưu tiên sau khi chọn mô hình:

1. `Recall`
2. `F1-score`
3. `Precision`
4. `Accuracy`
5. khả năng giải thích
6. chi phí tính toán

Theo tiêu chí này, mô hình được chọn là:

```text
Logistic Regression
```

Lý do lựa chọn:

- đạt metric tương đương Random Forest;
- nhẹ hơn về chi phí huấn luyện và suy luận;
- dễ diễn giải hơn;
- phù hợp hơn với mục tiêu đồ án khi cần trình bày mối liên hệ giữa đặc trưng và quyết định dự báo;
- thuận lợi khi tích hợp vào giao diện cảnh báo cho quản trị viên và học viên.

Như vậy:

- mô hình luật nghiệp vụ cơ sở đóng vai trò baseline;
- Random Forest đóng vai trò mô hình đối chứng phi tuyến;
- Logistic Regression được chọn làm mô hình chính.

**[Chèn Hình 4.4.12: Sơ đồ lập luận lựa chọn Logistic Regression làm mô hình chính]**

### 4.8.12 Đầu ra hệ thống và tích hợp vận hành

Sau khi mô hình được huấn luyện, hệ thống sinh ra các tệp kết quả mô hình dùng cho tầng backend và giao diện. Các tệp chính bao gồm:

- tệp metadata mô hình;
- tệp metric đánh giá;
- tệp dự báo mới nhất cho từng học viên;
- báo cáo phân loại;
- các hình trực quan hóa.

Từ các tệp này, hệ thống cung cấp:

- danh sách học viên có nguy cơ cho màn quản trị;
- thông tin cảnh báo cho học viên ở cổng thông tin học viên;
- các lý do giải thích ngắn gọn cho cảnh báo.

Nhờ đó, phần học máy không dừng ở mức phân tích ngoại tuyến mà đã được gắn trực tiếp vào luồng vận hành của EduCenter.

**[Chèn Hình 4.4.13: Luồng đưa kết quả dự báo vào dashboard quản trị và cổng học viên]**

### 4.8.13 Hạn chế và hướng phát triển

Dù quy trình dự báo hiện tại đã vận hành được, vẫn còn một số hạn chế:

1. số lượng mẫu còn ít (`36` mẫu), do đó kết quả đánh giá chưa phản ánh đầy đủ độ phức tạp của môi trường thực tế;
2. `Recall = 0.6` cho thấy mô hình vẫn còn bỏ sót một phần học viên có nguy cơ;
3. dữ liệu hiện tại là dữ liệu demo từ hệ thống, chưa đủ đa dạng như dữ liệu dài hạn trong vận hành thực tế;
4. nhóm mô hình thử nghiệm mới dừng ở các thuật toán nhẹ.

Các hướng mở rộng trong tương lai gồm:

- tăng quy mô dữ liệu thật từ hệ thống;
- điều chỉnh ngưỡng quyết định để ưu tiên `Recall`;
- thử nghiệm các mô hình boosting như `XGBoost` hoặc `LightGBM`;
- bổ sung kỹ thuật giải thích mô hình sâu hơn;
- theo dõi sự biến đổi điểm rủi ro theo thời gian cho từng học viên.

**[Chèn Hình 4.4.14: Định hướng mở rộng nhánh dự báo học tập trong tương lai]**

### 4.8.14 Kết luận của nhánh nghiên cứu

Nhánh dự báo học viên có nguy cơ học tập kém là một phần mở rộng có chiều sâu học thuật của đồ án EduCenter. Phần này thể hiện đầy đủ các thành phần cốt lõi của một nghiên cứu ứng dụng học máy:

- xác định bài toán;
- xây dựng dữ liệu;
- tiền xử lý;
- mô hình hóa toán học;
- thử nghiệm nhiều mô hình;
- đánh giá định lượng;
- lựa chọn mô hình chính;
- và tích hợp kết quả vào hệ thống thực tế.

Kết quả hiện tại cho thấy:

- mô hình luật nghiệp vụ cơ sở chưa đủ mạnh;
- mô hình học máy mang lại cải thiện rõ rệt;
- Logistic Regression là lựa chọn phù hợp nhất ở giai đoạn hiện tại;
- hướng tiếp cận bằng Python và dữ liệu từ hệ thống EduCenter là đủ chặt chẽ để trở thành một điểm nhấn nghiên cứu của đồ án.

---

## Gợi ý chèn vào báo cáo chính

Nếu bạn duyệt bản này, mình sẽ dùng nó để:

1. thay phần tương ứng trong báo cáo local;
2. chuẩn hóa lại theo số thứ tự hình/bảng của toàn chương;
3. sau đó cập nhật trực tiếp vào Google Doc bạn đã chỉ định.
