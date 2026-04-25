# Báo cáo huấn luyện và đánh giá mô hình AT_RISK

- Thời điểm sinh báo cáo: `2026-04-25T02:24:04.202421Z`
- Dataset: `at_risk_dataset_demo`
- Nguồn dữ liệu: `csv:at_risk_dataset_demo.csv`
- Seed: `42`
- Tỉ lệ test: `0.2`

## 1. Tóm tắt dataset

- Số dòng: `20`
- Số học viên: `20`
- Số lớp: `2`
- Feature columns: `attendance_rate_28d, absence_count_28d, average_total_score_28d, homework_completion_rate_28d, active_enrollment_count_28d, weekly_lesson_load_28d, approved_leave_count_28d, days_since_last_lesson`
- Phân phối nhãn: `AT_RISK=10`, `NOT_AT_RISK=10`

## 2. Bảng so sánh mô hình

| Mô hình | Accuracy | Precision | Recall | F1 | TP | FP | FN | TN |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `rule_based` | 1.0000 | 1.0000 | 1.0000 | 1.0000 | 2 | 0 | 0 | 2 |
| `logistic_regression` | 1.0000 | 1.0000 | 1.0000 | 1.0000 | 2 | 0 | 0 | 2 |
| `random_forest` | 1.0000 | 1.0000 | 1.0000 | 1.0000 | 2 | 0 | 0 | 2 |

## 3. Nhận xét chính

- Mô hình có `F1-score` cao nhất: `rule_based`
- Mô hình có `Recall` cao nhất: `rule_based`
- Mô hình được dùng để vẽ confusion matrix: `rule_based`

## 4. Cấu hình baseline rule-based

```json
{
  "name": "rule_based",
  "thresholds": {
    "attendance_rate_28d": 0.8,
    "average_total_score_28d": 5.0,
    "homework_completion_rate_28d": 0.6
  },
  "decision_threshold": 0.5,
  "weights": {
    "attendance": 0.4,
    "score": 0.35,
    "homework": 0.25
  },
  "generated_at": "2026-04-25T02:24:03.918629Z"
}
```

## 5. Phân tích chi tiết từng mô hình

### 5.1. rule_based

- Accuracy: `1.0000`
- Precision: `1.0000`
- Recall: `1.0000`
- F1-score: `1.0000`
- Confusion matrix: `[[2, 0], [0, 2]]`

### 5.2. logistic_regression

- Accuracy: `1.0000`
- Precision: `1.0000`
- Recall: `1.0000`
- F1-score: `1.0000`
- Confusion matrix: `[[2, 0], [0, 2]]`

### 5.3. random_forest

- Accuracy: `1.0000`
- Precision: `1.0000`
- Recall: `1.0000`
- F1-score: `1.0000`
- Confusion matrix: `[[2, 0], [0, 2]]`

## 6. Kết luận chọn mô hình chính

- Mô hình được chọn: `logistic_regression`
- Các mô hình có metric tương đương: `rule_based, logistic_regression, random_forest`
- Lý do chọn: Các mô hình có metric tương đương trên tập kiểm thử hiện tại; ưu tiên Logistic Regression vì vẫn là mô hình học máy chính thức, nhẹ khi huấn luyện/suy luận và có khả năng giải thích tốt hơn Random Forest.
- Tiêu chí ưu tiên: `Recall -> F1-score -> Precision -> Accuracy -> explainability -> lightweight`
