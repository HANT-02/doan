# AT_RISK Prediction Python Project

Project này là nhánh chính thức cho bài toán dự báo học viên có nguy cơ học tập kém (`AT_RISK`) trong đồ án EduCenter.

Mục tiêu của project Python:

- lấy dữ liệu từ PostgreSQL hoặc CSV,
- xây dựng dataset cho bài toán classification,
- huấn luyện và đánh giá mô hình,
- sinh artifact để Go backend tích hợp vào API và frontend,
- tạo ra phần thực nghiệm có thể đưa trực tiếp vào báo cáo đồ án.

## 1. Bài toán của project

### 1.1. Bài toán dự báo

Mỗi mẫu dữ liệu đại diện cho một `student_class_snapshot`, tức là ảnh chụp trạng thái của một học viên trong một lớp tại một thời điểm quan sát.

Đầu vào:

- dữ liệu hồ sơ học viên,
- dữ liệu ghi danh,
- dữ liệu buổi học,
- dữ liệu điểm danh,
- dữ liệu tổng kết buổi học,
- dữ liệu kết quả học tập,
- dữ liệu đơn xin phép,
- dữ liệu trích xuất từ PostgreSQL hoặc CSV.

Đầu ra:

- nhãn `AT_RISK` hoặc `NOT_AT_RISK`,
- `risk_score`,
- `model_metadata`,
- các yếu tố ảnh hưởng chính,
- các artifact phục vụ Go backend và báo cáo đồ án.

### 1.2. Vai trò của Python và Go

Python chịu trách nhiệm:

- data loading,
- feature engineering,
- train/test split,
- huấn luyện mô hình,
- đánh giá metric,
- sinh prediction artifact.

Go backend chịu trách nhiệm:

- đọc artifact do Python sinh ra,
- cung cấp API cho admin/student,
- giữ ổn định contract với frontend hiện có.

## 2. Công nghệ sử dụng

| Công nghệ | Vai trò |
|---|---|
| `Python 3.11+` | Ngôn ngữ chính cho pipeline ML |
| `pandas` | Xử lý bảng dữ liệu, feature engineering |
| `numpy` | Tính toán số học, vector/matrix |
| `scikit-learn` | Huấn luyện và đánh giá model |
| `sqlalchemy` / `psycopg` | Kết nối PostgreSQL |
| `python-dotenv` | Nạp biến môi trường từ `.env` |
| `joblib` | Lưu model artifact |
| `matplotlib` / `seaborn` | Vẽ confusion matrix, feature importance, chart |
| `Jupyter` | Khám phá dữ liệu và kiểm tra thực nghiệm khi cần |

## 3. Cấu trúc thư mục

```text
ml/at_risk_prediction/
├── README.md
├── requirements.txt
├── .env.example
├── src/
│   ├── __init__.py
│   ├── config.py
│   └── ...
├── scripts/
│   └── ...
├── data/
│   ├── raw/
│   └── processed/
├── artifacts/
│   ├── models/
│   ├── reports/
│   └── figures/
└── docs/
```

Ý nghĩa:

- `src/`: mã nguồn chính của pipeline ML.
- `scripts/`: các entrypoint chạy export/train/predict.
- `data/raw/`: dữ liệu gốc lấy từ DB hoặc CSV.
- `data/processed/`: dữ liệu sau khi chuẩn hóa/feature engineering.
- `artifacts/models/`: model và metadata.
- `artifacts/reports/`: metric, report, prediction output.
- `artifacts/figures/`: hình phục vụ báo cáo đồ án.
- `docs/`: tài liệu nghiên cứu riêng của nhánh predictive.

## 4. Khởi tạo môi trường Python

### 4.1. Tạo virtual environment

Từ thư mục gốc repo:

```bash
cd /Users/hant/golang/doan/ml/at_risk_prediction
python3 -m venv .venv
```

Kích hoạt môi trường:

```bash
source .venv/bin/activate
```

### 4.2. Cài dependencies

```bash
pip install --upgrade pip
pip install -r requirements.txt
```

### 4.3. Kiểm tra nhanh

```bash
python --version
pip list
```

Kỳ vọng:

- Python từ `3.11` trở lên,
- cài được đầy đủ `pandas`, `numpy`, `scikit-learn`, `sqlalchemy`, `psycopg`, `joblib`, `matplotlib`, `seaborn`.

## 5. Cấu hình kết nối dữ liệu

### 5.1. File môi trường

Tạo file `.env` từ `.env.example`:

```bash
cp .env.example .env
```

Nội dung hiện tại:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=doan
DB_SSLMODE=disable

DATASET_CSV_PATH=./data/raw/at_risk_dataset.csv
ARTIFACTS_DIR=./artifacts
MODEL_NAME=logistic_regression
```

### 5.2. Nguồn dữ liệu hỗ trợ

Project này được thiết kế để hỗ trợ 2 nguồn dữ liệu:

1. `PostgreSQL trực tiếp`
2. `CSV offline`

### 5.3. Các bảng DB dự kiến dùng

- `students`
- `enrollments`
- `classes`
- `lessons`
- `attendance`
- `lesson_summaries`
- `academic_records`
- `leave_requests`

## 6. Workflow dữ liệu của project

Luồng chuẩn:

```text
PostgreSQL / CSV
    ->
data loader
    ->
feature engineering
    ->
train / test split
    ->
Rule-based / Logistic Regression / Random Forest
    ->
metric evaluation
    ->
model artifact + prediction artifact
    ->
Go backend đọc artifact
    ->
API / frontend admin / frontend student
```

## 7. Các script chuẩn của project

Lưu ý:

- Ở `Task K3`, các script `export_dataset.py`, `train_from_db.py`, `train_from_csv.py` đã được tạo.
- `export_dataset.py` dùng để dựng dataset chuẩn hóa từ DB.
- Từ `K5`, `train_from_db.py` và `train_from_csv.py` đã huấn luyện đầy đủ `rule-based`, `Logistic Regression`, `Random Forest`, đồng thời sinh metric, report, figure và model artifact.

### 7.1. Export dataset từ DB

Script hiện có:

```text
scripts/export_dataset.py
```

Vai trò:

- kết nối PostgreSQL,
- join dữ liệu từ các bảng cần thiết,
- dựng dataset `student_class_snapshot`,
- ghi dữ liệu ra `data/raw/` hoặc `data/processed/`.

Lệnh chuẩn:

```bash
python scripts/export_dataset.py --dataset-name at_risk_dataset
```

### 7.2. Train từ DB

Script hiện có:

```text
scripts/train_from_db.py
```

Vai trò:

- đọc cấu hình từ `.env`,
- lấy dữ liệu trực tiếp từ PostgreSQL,
- feature engineering,
- sinh dataset `full/train/test`,
- sinh dataset summary report vào thư mục `artifacts/`.

Lệnh chuẩn dự kiến:

```bash
python scripts/train_from_db.py --dataset-name at_risk_dataset_db
```

### 7.3. Train từ CSV

Script hiện có:

```text
scripts/train_from_csv.py
```

Vai trò:

- đọc dữ liệu CSV đã export,
- chuẩn hóa dataset và chia train/test mà không cần kết nối DB,
- phù hợp cho việc tái lập thực nghiệm và chia sẻ dataset nội bộ.

CSV demo hiện có sẵn trong repo:

```text
data/raw/at_risk_dataset_demo.csv
```

Lệnh chuẩn:

```bash
python scripts/train_from_csv.py --input ./data/raw/at_risk_dataset_demo.csv --dataset-name at_risk_dataset_demo
```

### 7.4. Sinh prediction artifact

Script hiện có từ `K6`:

```text
scripts/predict_from_db.py
scripts/predict_from_csv.py
```

Vai trò:

- tự chọn mô hình chính từ `metrics.json` và `model_metadata.json`,
- sinh file prediction mà Go backend có thể đọc trực tiếp,
- hỗ trợ cả dữ liệu `DB` và `CSV`,
- nếu dữ liệu DB chưa sinh được snapshot thì vẫn trả artifact rỗng có thông báo rõ ràng.

Lệnh chuẩn:

```bash
python scripts/predict_from_db.py
```

Hoặc:

```bash
python scripts/predict_from_csv.py --input ./data/raw/at_risk_dataset_demo.csv
```

## 8. Artifact đầu ra cần có

### 8.1. Model artifact

Thư mục:

```text
artifacts/models/
```

Các file hiện có sau `K5`:

- `rule_based.json`
- `logistic_regression.joblib`
- `random_forest.joblib`
- `model_metadata.json`

### 8.2. Report artifact

Thư mục:

```text
artifacts/reports/
```

Các file hiện có sau `K5`:

- `metrics.json`
- `classification_report.md`
- `latest_predictions.json`

### 8.3. Figure artifact

Thư mục:

```text
artifacts/figures/
```

Các file hiện có sau `K5`:

- `confusion_matrix.png`
- `feature_importance.png`
- `model_comparison.png`

## 9. Go backend sẽ đọc gì từ project Python

Go backend không còn train model.

Go backend chỉ đọc các artifact sau:

- `artifacts/models/model_metadata.json`
- `artifacts/reports/latest_predictions.json`

Mục tiêu:

- giữ route admin predictive hiện có,
- giữ route student `/api/v1/student/at-risk`,
- không làm frontend phải sửa lớn.

## 10. Hai cách làm việc với dữ liệu

### 10.1. Hướng 1: làm việc trực tiếp với DB

Phù hợp khi:

- muốn dùng dữ liệu thật của hệ thống,
- đã có lesson, attendance, academic record tương đối đủ,
- muốn có kết quả sát với dữ liệu production-like.

Ưu điểm:

- dữ liệu mới nhất,
- không cần export tay nhiều lần.

Nhược điểm:

- phụ thuộc trạng thái DB,
- khó chia sẻ thực nghiệm nếu DB thay đổi liên tục.

### 10.2. Hướng 2: làm việc qua CSV

Phù hợp khi:

- muốn khóa dataset thực nghiệm,
- muốn tái lập nhiều lần cùng một bộ dữ liệu,
- muốn đưa file dataset hoặc snapshot vào báo cáo/phụ lục nội bộ.

Ưu điểm:

- tái lập tốt,
- thuận tiện benchmark giữa các model.

Nhược điểm:

- cần thêm bước export dữ liệu.

## 11. Khuyến nghị sử dụng trong đồ án

Thứ tự thực hành nên là:

1. Export dữ liệu từ DB thật ra CSV.
2. Dùng CSV để train/evaluate và so sánh mô hình.
3. Chọn mô hình chính.
4. Sinh `latest_predictions.json`.
5. Để Go backend đọc artifact này và phục vụ UI.

Lý do:

- vừa bám dữ liệu thật của hệ thống,
- vừa giữ được tính tái lập của thực nghiệm,
- vừa phù hợp với cách viết báo cáo đồ án.

## 12. Ghi chú phát triển theo task

| Task | Mục tiêu |
|---|---|
| `K1` | Tạo skeleton project |
| `K2` | Viết README vận hành |
| `K3` | Tạo data pipeline DB/CSV |
| `K4` | Viết đặc tả nghiên cứu và mô hình toán học |
| `K5` | Huấn luyện, đánh giá, sinh metric |
| `K6` | Chọn mô hình chính, sinh prediction artifact |
| `K7` | Refactor Go backend đọc artifact |
| `K8` | Cập nhật báo cáo đồ án |

## 13. Trạng thái hiện tại

Đã hoàn thành:

- tạo project Python trong repo,
- khóa công nghệ sử dụng,
- khóa cấu trúc thư mục,
- viết README vận hành tổng quát,
- hoàn thành data pipeline DB/CSV,
- viết đặc tả nghiên cứu,
- huấn luyện và đánh giá 3 mô hình,
- sinh model artifact, metric artifact và figure artifact.

Chưa hoàn thành trong bước này:

- tích hợp Go backend.

Đã hoàn thành thêm ở `K6`:

- chọn mô hình chính và ghi vào `artifacts/models/model_metadata.json`
- sinh `artifacts/reports/latest_predictions.json`
- thêm `predict_from_db.py` và `predict_from_csv.py`
- thêm giải thích `risk_band`, `primary_reason`, `top_features`, `feature_summary` cho từng prediction

Ghi chú thực nghiệm hiện tại:

- train từ `CSV demo` đã chạy thành công và sinh đủ artifact tại `artifacts/`
- train từ `DB thật` đã kết nối được, nhưng dữ liệu hiện tại chưa sinh được dòng `student_class_snapshot`, nên cần bổ sung dữ liệu học vụ trước khi benchmark trên dữ liệu production-like
