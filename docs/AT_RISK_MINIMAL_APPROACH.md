# Minimal AT_RISK Approach

## Mục tiêu

Giảm tối đa độ phức tạp triển khai `F2` để phù hợp:

- thời gian đồ án ngắn
- máy cá nhân không đủ để train model nặng
- không muốn phụ thuộc AI runtime hoặc cloud model ngoài

## Hướng triển khai được chốt

Thay vì làm pipeline nặng với nhiều dependency ML, hướng tối giản là:

1. `Rule-based baseline`
2. `Logistic Regression` rất nhẹ, train CPU
3. Toàn bộ chạy bằng Go, không cần `sklearn`, `pandas`, `numpy`

## Vì sao chọn hướng này

- Không cần GPU
- Không cần cài thêm thư viện ML ngoài
- Không cần model AI có sẵn từ cloud
- Dễ train, dễ giải thích, dễ bảo vệ
- Vẫn đủ tính chất machine learning để trình bày trong đồ án

## Thành phần đã có trong code

### Code chính

- [minimal_pipeline.go](/Users/hant/golang/doan/internal/services/predictive/minimal_pipeline.go)
- [minimal_pipeline_test.go](/Users/hant/golang/doan/internal/services/predictive/minimal_pipeline_test.go)
- [db_loader.go](/Users/hant/golang/doan/internal/services/predictive/db_loader.go)
- [db_loader_test.go](/Users/hant/golang/doan/internal/services/predictive/db_loader_test.go)
- [runtime.go](/Users/hant/golang/doan/internal/services/predictive/runtime.go)
- [main.go](/Users/hant/golang/doan/cmd/cli/predictive_train/main.go)
- [v1.go](/Users/hant/golang/doan/cmd/http/controllers/predictive/v1.go)
- [PredictiveAlertsPage.tsx](/Users/hant/golang/doan/frontend/src/pages/admin/PredictiveAlertsPage.tsx)

### Lệnh chạy

```bash
make predictive-train
```

Hoặc dùng dataset CSV riêng:

```bash
go run ./cmd/cli/predictive_train --source csv --dataset /path/to/at_risk_dataset.csv --name real_dataset
```

Hoặc lấy trực tiếp từ PostgreSQL của hệ thống:

```bash
go run ./cmd/cli/predictive_train --source db --config-file-path ./configs --config-file config
```

## Output hiện tại

Pipeline tối giản sẽ:

1. nạp dataset demo, CSV thật, hoặc dữ liệu từ DB thật
2. tách train/test
3. train `Logistic Regression`
4. chạy `Rule-based baseline`
5. tính:
   - Accuracy
   - Precision
   - Recall
   - F1
6. đề xuất model nên dùng

## Cách lấy dữ liệu từ DB

Loader DB hiện tại dùng đúng pattern backend sẵn có:

- load config bằng `Viper`
- mở kết nối Postgres bằng `postgres.GetDBContext(...)`
- nạp dữ liệu từ các bảng:
  - `students`
  - `enrollments`
  - `lessons`
  - `attendances`
  - `lesson_summaries`
  - `academic_records`
  - `leave_requests`
- dựng `student_class_snapshot` bằng cửa sổ:
  - quan sát `28 ngày`
  - dự báo `28 ngày`

Các feature tối giản vẫn giữ nguyên:

- `attendance_rate_28d`
- `absence_count_28d`
- `average_total_score_28d`
- `homework_completion_rate_28d`
- `active_enrollment_count_28d`
- `weekly_lesson_load_28d`
- `approved_leave_count_28d`
- `days_since_last_lesson`

## Prediction API và UI

Phần F3/F4 hiện đã có bản tối giản dùng được:

- API danh sách cảnh báo:
  - `GET /api/v1/predictive/at-risk/students`
- API metadata model:
  - `GET /api/v1/predictive/at-risk/model-metadata`
- UI admin:
  - `/app/admin/predictive`

API hiện chạy theo cơ chế:

1. load dữ liệu từ PostgreSQL
2. train nhẹ trong memory
3. cache model metadata trong runtime HTTP app
4. suy luận danh sách học sinh/lớp có nguy cơ
5. trả thêm `risk_score`, `label`, `reasons`, `model_version`

Giới hạn thực tế cần lưu ý:

- Nếu DB chưa có `lesson` cho các lớp đang học thì predictive từ DB sẽ chưa chạy được.
- Khi đó cần commit scheduling preview xuống `lesson` hoặc seed dữ liệu lesson trước.

## Giới hạn đã chấp nhận

- Chưa làm `Random Forest`
- Chưa làm `XGBoost/LightGBM`
- Chưa làm feature importance nâng cao
- Model metadata hiện lưu tối giản trong memory runtime, chưa persist xuống bảng riêng

Những phần này có thể làm sau nếu còn thời gian.

## Kết luận ngắn

Đây là hướng phù hợp nhất nếu mục tiêu là:

- có predictive analytics thật
- đủ nhẹ để chạy trên máy cá nhân
- đủ nhanh để hoàn thành trong thời gian còn lại
- không làm dự án bị phình quá mức
